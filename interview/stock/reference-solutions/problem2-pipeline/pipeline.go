// Package pipeline implements a bounded, parallel enrichment pipeline that
// preserves strict input order on output despite out-of-order completion.
package pipeline

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"time"
)

// TimeoutErr is reported in Out.Err when a message's enrichment exceeds
// perMsgTimeout.
var TimeoutErr = errors.New("pipeline: enrichment timeout")

type Msg struct {
	ID  int // strictly increasing on input
	Raw string
}

type Out struct {
	Msg      Msg
	Enriched string
	Err      error // nil on success; TimeoutErr on timeout; enricher error otherwise
}

// Enricher may be slow, may return an error. It is NOT guaranteed to respect
// ctx quickly (it may ignore cancellation and keep running).
type Enricher func(ctx context.Context, m Msg) (string, error)

type task struct {
	seq int // internal submission sequence number, used for ordering
	msg Msg
}

// Pipeline enriches messages in parallel (bounded by `workers`) while
// guaranteeing output emerges from Out() in strict submission order.
type Pipeline struct {
	perMsgTimeout time.Duration
	enricher      Enricher

	inCh  chan task // bounded input queue; Submit blocks when full
	outCh chan Out  // bounded output queue

	sem chan struct{} // bounds concurrently-awaited enrichment attempts to `workers`

	baseCtx context.Context

	mu      sync.Mutex
	buffer  map[int]Out // seq -> result, waiting for its turn to be emitted
	nextSeq int         // next seq expected to flush toward Out()

	seqCounter int64 // atomic: next seq to assign on Submit

	closed    int32 // atomic bool: 1 once Drain has begun; further Submit panics
	submitted int64 // atomic: total accepted via Submit
	emitted   int64 // atomic: total actually sent on outCh

	dispatchWG sync.WaitGroup // the single dispatch loop goroutine
	taskWG     sync.WaitGroup // one count per accepted task, released once resolved

	stopEmit  chan struct{} // closed by Drain once it stops waiting; gates sends
	sendWG    sync.WaitGroup // brackets in-flight attempts to send on outCh
	drainOnce sync.Once
}

// NewPipeline creates a pipeline with the given worker concurrency,
// per-message timeout, and enrichment function.
//
// Bounded memory: the input queue holds up to 2*workers pending messages and
// the output queue holds up to 2*workers pending results before Submit /
// internal delivery block. Combined with the `workers` concurrency cap, this
// puts a hard ceiling on how many messages can be "in flight" at once
// (accepted but not yet emitted), which is what keeps the reorder buffer
// bounded too -- it can never hold more entries than there are in-flight
// tasks.
func NewPipeline(workers int, perMsgTimeout time.Duration, e Enricher) *Pipeline {
	if workers < 1 {
		workers = 1
	}
	p := &Pipeline{
		perMsgTimeout: perMsgTimeout,
		enricher:      e,
		inCh:          make(chan task, workers*2),
		outCh:         make(chan Out, workers*2),
		sem:           make(chan struct{}, workers),
		baseCtx:       context.Background(),
		buffer:        make(map[int]Out),
		stopEmit:      make(chan struct{}),
	}
	p.dispatchWG.Add(1)
	go p.dispatchLoop()
	return p
}

// Submit enqueues a message for enrichment, blocking if internal capacity is
// full. Submit must not be called after Drain has started; doing so panics
// (documented choice -- fail loud rather than silently drop, since a caller
// racing Submit against Drain almost always indicates a shutdown-sequencing
// bug in the caller).
func (p *Pipeline) Submit(m Msg) {
	if atomic.LoadInt32(&p.closed) == 1 {
		panic("pipeline: Submit called after Drain")
	}
	seq := int(atomic.AddInt64(&p.seqCounter, 1) - 1)
	p.taskWG.Add(1)
	p.inCh <- task{seq: seq, msg: m}
	atomic.AddInt64(&p.submitted, 1)
}

// Out returns the channel of results, delivered in strict ascending
// submission order. It is closed once Drain has finished flushing.
func (p *Pipeline) Out() <-chan Out {
	return p.outCh
}

func (p *Pipeline) dispatchLoop() {
	defer p.dispatchWG.Done()
	for t := range p.inCh {
		p.sem <- struct{}{} // acquire a concurrency slot
		go p.runTask(t)
	}
}

// runTask races the enricher against perMsgTimeout. If the enricher doesn't
// return in time, we emit a timeout result immediately and move on --
// we do NOT keep the worker slot pinned waiting for a straggler that may be
// ignoring ctx, because that would let one slow message stall the whole
// pipeline (rule 3). The straggler's own goroutine is left to finish (or
// leak) in the background; its eventual result is discarded.
func (p *Pipeline) runTask(t task) {
	defer p.taskWG.Done()

	ctx, cancel := context.WithTimeout(p.baseCtx, p.perMsgTimeout)
	resCh := make(chan Out, 1) // buffered so the inner goroutine never blocks on send
	go func() {
		enriched, err := p.enricher(ctx, t.msg)
		resCh <- Out{Msg: t.msg, Enriched: enriched, Err: err}
	}()

	var out Out
	select {
	case out = <-resCh:
	case <-ctx.Done():
		out = Out{Msg: t.msg, Enriched: "", Err: TimeoutErr}
	}
	cancel()
	<-p.sem // release the concurrency slot regardless of straggler state

	p.deliver(t.seq, out)
}

// deliver stores a result and flushes any now-contiguous run of results
// (starting at nextSeq) onto Out(), preserving strict order.
//
// IMPORTANT: the actual send to outCh happens WHILE p.mu is still held, not
// after releasing it. If two goroutines' deliver() calls each computed their
// own "ready" slice and then sent after unlocking, a second goroutine could
// race ahead and send its (later) result before the first goroutine gets
// scheduled to send its (earlier) one -- corrupting order even though the
// buffer bookkeeping itself was correct. Holding the lock across the send
// serializes emission in the same order nextSeq is advanced, at the cost of
// blocking other deliverers while Out() has no reader (an accepted
// trade-off: correctness over max throughput under a stalled consumer).
func (p *Pipeline) deliver(seq int, out Out) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.buffer[seq] = out
	for {
		o, ok := p.buffer[p.nextSeq]
		if !ok {
			break
		}
		delete(p.buffer, p.nextSeq)
		p.nextSeq++
		p.send(o)
	}
}

// send delivers a single result to Out(), unless Drain has already given up
// waiting (stopEmit closed), in which case it is dropped and counted via the
// submitted/emitted diff. sendWG lets Drain know when it's safe to close
// Out() without racing a concurrent send.
func (p *Pipeline) send(o Out) {
	p.sendWG.Add(1)
	defer p.sendWG.Done()

	select {
	case <-p.stopEmit:
		return // Drain already gave up; do not touch outCh (may be closed soon).
	default:
	}

	p.outCh <- o
	atomic.AddInt64(&p.emitted, 1)
}

// Drain stops accepting new messages, waits (up to ctx's deadline) for
// in-flight enrichments to finish and be emitted in order, then reports how
// many accepted messages were never emitted, and closes Out().
func (p *Pipeline) Drain(ctx context.Context) (dropped int) {
	p.drainOnce.Do(func() {
		atomic.StoreInt32(&p.closed, 1)
		close(p.inCh)
	})

	allDone := make(chan struct{})
	go func() {
		p.dispatchWG.Wait() // dispatch loop has drained inCh
		p.taskWG.Wait()     // every accepted task has resolved (success or timeout) and been delivered
		close(allDone)
	}()

	select {
	case <-allDone:
	case <-ctx.Done():
		// Deadline hit before everything finished. Some runTask goroutines
		// may still be alive (e.g. blocked in an enricher that ignores ctx
		// past its own perMsgTimeout in a way that's still running its final
		// steps). We stop waiting for them; their eventual deliver() calls
		// will see stopEmit closed below and drop cleanly instead of
		// blocking or panicking on a closed channel.
	}

	close(p.stopEmit)
	p.sendWG.Wait() // ensure no send is mid-flight before we close outCh

	dropped = int(atomic.LoadInt64(&p.submitted) - atomic.LoadInt64(&p.emitted))
	close(p.outCh)
	return dropped
}
