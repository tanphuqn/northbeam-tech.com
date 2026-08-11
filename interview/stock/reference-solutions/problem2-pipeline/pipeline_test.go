package pipeline

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"
)

// run: go test -race ./...

// TestStrictOrderDespiteOutOfOrderCompletion (rule 1): an enricher whose
// delay is deliberately inversely correlated with ID -- early messages
// finish LAST -- must still be emitted in strictly ascending ID order.
func TestStrictOrderDespiteOutOfOrderCompletion(t *testing.T) {
	const n = 12
	enricher := func(ctx context.Context, m Msg) (string, error) {
		delay := time.Duration(n-m.ID) * 4 * time.Millisecond
		time.Sleep(delay)
		return fmt.Sprintf("enriched-%d", m.ID), nil
	}

	p := NewPipeline(6, 500*time.Millisecond, enricher)

	go func() {
		for i := 0; i < n; i++ {
			p.Submit(Msg{ID: i, Raw: fmt.Sprintf("raw-%d", i)})
		}
	}()

	var got []int
	for o := range p.Out() {
		if o.Err != nil {
			t.Fatalf("unexpected error for ID %d: %v", o.Msg.ID, o.Err)
		}
		got = append(got, o.Msg.ID)
		if len(got) == n {
			break
		}
	}

	for i, id := range got {
		if id != i {
			t.Fatalf("output order broken: position %d has ID %d, want %d (full: %v)", i, id, i, got)
		}
	}

	drained := make(chan struct{})
	go func() {
		p.Drain(context.Background())
		close(drained)
	}()
	select {
	case <-drained:
	case <-time.After(2 * time.Second):
		t.Fatal("Drain did not complete")
	}
}

// TestTimeoutDoesNotStallStream (rules 3 + the hint): a straggler that
// ignores its context must produce a TimeoutErr in its correct ordered
// position and must not block later messages from being emitted.
func TestTimeoutDoesNotStallStream(t *testing.T) {
	const n = 6
	const stragglerID = 2
	perMsgTimeout := 30 * time.Millisecond

	enricher := func(ctx context.Context, m Msg) (string, error) {
		if m.ID == stragglerID {
			// Deliberately ignore ctx and sleep well past perMsgTimeout.
			time.Sleep(400 * time.Millisecond)
			return "too-late", nil
		}
		return fmt.Sprintf("enriched-%d", m.ID), nil
	}

	p := NewPipeline(3, perMsgTimeout, enricher)

	go func() {
		for i := 0; i < n; i++ {
			p.Submit(Msg{ID: i, Raw: "x"})
		}
	}()

	start := time.Now()
	var got []Out
	for o := range p.Out() {
		got = append(got, o)
		if len(got) == n {
			break
		}
	}
	elapsed := time.Since(start)

	// The whole stream must complete close to perMsgTimeout, NOT close to
	// the straggler's real 400ms sleep -- proving it didn't stall the stream.
	if elapsed > 200*time.Millisecond {
		t.Fatalf("stream took %v, expected it to finish shortly after perMsgTimeout (%v) without waiting on the straggler", elapsed, perMsgTimeout)
	}

	for i, o := range got {
		if o.Msg.ID != i {
			t.Fatalf("order broken at position %d: got ID %d", i, o.Msg.ID)
		}
		if i == stragglerID {
			if o.Err != TimeoutErr {
				t.Fatalf("expected TimeoutErr at position %d, got %v", i, o.Err)
			}
		} else if o.Err != nil {
			t.Fatalf("unexpected error at position %d: %v", i, o.Err)
		}
	}

	drained := make(chan struct{})
	go func() {
		p.Drain(context.Background())
		close(drained)
	}()
	select {
	case <-drained:
	case <-time.After(2 * time.Second):
		t.Fatal("Drain did not complete (possibly blocked on straggler)")
	}
}

// TestDrainWithMessagesInFlight (rule 5): Drain with a short deadline while
// slow enrichments are still running must report the correct dropped count
// and cleanly close Out() without panicking, even though some background
// goroutines are still alive when Drain returns.
func TestDrainWithMessagesInFlight(t *testing.T) {
	const n = 5
	block := make(chan struct{}) // never closed: enrichments hang until the test ends
	started := make(chan struct{}, n)

	enricher := func(ctx context.Context, m Msg) (string, error) {
		started <- struct{}{}
		select {
		case <-block:
		case <-ctx.Done():
			// respects ctx in this test, but only after perMsgTimeout fires;
			// this still exercises the "give up waiting" path in Drain
			// because Drain's own ctx deadline is shorter than perMsgTimeout.
			<-block
		}
		return "unused", nil
	}

	p := NewPipeline(2, 10*time.Second, enricher) // perMsgTimeout intentionally long

	for i := 0; i < n; i++ {
		p.Submit(Msg{ID: i, Raw: "x"})
	}

	// Let at least the first couple of tasks actually start.
	<-started
	<-started

	// Drain a reader concurrently so trySend never blocks forever on a full
	// outCh (Drain's contract assumes an active/eventual consumer).
	var wg sync.WaitGroup
	wg.Add(1)
	var received int
	go func() {
		defer wg.Done()
		for range p.Out() {
			received++
		}
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()
	dropped := p.Drain(ctx)

	wg.Wait() // Out() closed; reader goroutine has exited

	if dropped+received != n {
		t.Fatalf("dropped(%d) + received(%d) should equal submitted(%d)", dropped, received, n)
	}
	if dropped == 0 {
		t.Fatalf("expected at least one dropped message given the short Drain deadline vs. hung enrichers")
	}
}

// TestConcurrentSubmitAndConsume exercises the pipeline under real
// concurrency (submitter + consumer running simultaneously) to catch data
// races; must pass under `go test -race`.
func TestConcurrentSubmitAndConsume(t *testing.T) {
	const n = 200
	enricher := func(ctx context.Context, m Msg) (string, error) {
		time.Sleep(time.Duration(m.ID%5) * time.Millisecond)
		return "ok", nil
	}
	p := NewPipeline(8, 200*time.Millisecond, enricher)

	go func() {
		for i := 0; i < n; i++ {
			p.Submit(Msg{ID: i, Raw: "x"})
		}
	}()

	next := 0
	for o := range p.Out() {
		if o.Msg.ID != next {
			t.Fatalf("out of order: got %d, want %d", o.Msg.ID, next)
		}
		next++
		if next == n {
			break
		}
	}

	done := make(chan struct{})
	go func() {
		p.Drain(context.Background())
		close(done)
	}()
	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("Drain did not complete")
	}
}
