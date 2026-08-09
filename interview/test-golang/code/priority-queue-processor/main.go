package main

import (
	"container/heap"
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

// Transaction matches the message structure from the queue.
type Transaction struct {
	ID        string    `json:"id"`
	Value     float64   `json:"value"`
	Timestamp time.Time `json:"timestamp"`
}

// --- Priority queue implementation using container/heap ---

type txItem struct {
	tx       Transaction
	priority float64 // higher value = higher priority
	index    int     // required by heap.Interface
}

type priorityQueue []*txItem

func (pq priorityQueue) Len() int { return len(pq) }

func (pq priorityQueue) Less(i, j int) bool {
	// max-heap: higher priority value comes out first
	return pq[i].priority > pq[j].priority
}

func (pq priorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *priorityQueue) Push(x any) {
	item := x.(*txItem)
	item.index = len(*pq)
	*pq = append(*pq, item)
}

func (pq *priorityQueue) Pop() any {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil
	*pq = old[:n-1]
	return item
}

// --- Processor: consumes messages (simulating RabbitMQ), buffers into a
//     priority queue, and processes high-value transactions first ---

const highValueThreshold = 1000

type Processor struct {
	mu     sync.Mutex
	pq     priorityQueue
	notify chan struct{} // signals the worker that a new item is available
}

func NewProcessor() *Processor {
	return &Processor{
		pq:     make(priorityQueue, 0),
		notify: make(chan struct{}, 1),
	}
}

// Enqueue is called for every message received from the queue.
func (p *Processor) Enqueue(tx Transaction) {
	priority := tx.Value // simplest priority function: higher value = higher priority

	p.mu.Lock()
	heap.Push(&p.pq, &txItem{tx: tx, priority: priority})
	p.mu.Unlock()

	select {
	case p.notify <- struct{}{}:
	default: // worker already has a pending notification, no need to send another
	}
}

// Run starts the worker loop; call in a goroutine. Stops when ctx is cancelled.
func (p *Processor) Run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-p.notify:
			for {
				p.mu.Lock()
				if p.pq.Len() == 0 {
					p.mu.Unlock()
					break
				}
				item := heap.Pop(&p.pq).(*txItem)
				p.mu.Unlock()

				p.process(item.tx)
			}
		}
	}
}

func (p *Processor) process(tx Transaction) {
	label := "normal"
	if tx.Value > highValueThreshold {
		label = "HIGH-VALUE"
	}
	fmt.Printf("[processed:%s] id=%s value=%.2f timestamp=%s\n",
		label, tx.ID, tx.Value, tx.Timestamp.Format(time.RFC3339))
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	processor := NewProcessor()
	go processor.Run(ctx)

	// Simulate messages arriving from the queue out of priority order:
	// low value first, high value arrives later but must still jump the queue.
	raw := []string{
		`{"id":"tx1","value":100,"timestamp":"2026-01-01T10:00:00Z"}`,
		`{"id":"tx2","value":50,"timestamp":"2026-01-01T10:00:01Z"}`,
		`{"id":"tx3","value":5000,"timestamp":"2026-01-01T10:00:02Z"}`, // high-value, arrives 3rd
		`{"id":"tx4","value":20,"timestamp":"2026-01-01T10:00:03Z"}`,
	}

	for _, r := range raw {
		var tx Transaction
		if err := json.Unmarshal([]byte(r), &tx); err != nil {
			fmt.Println("skip invalid message:", err) // permanent error, do not retry
			continue
		}
		processor.Enqueue(tx)
	}

	time.Sleep(200 * time.Millisecond) // give the worker time to drain the queue
}
