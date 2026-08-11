package dispatcher

import (
	"context"
	"fmt"
	"sync/atomic"
	"testing"
	"time"
)

// run: go test -race ./...

func TestResultsPreserveInputOrderDespiteOutOfOrderCompletion(t *testing.T) {
	const n = 10
	tasks := make([]Task[string], n)
	for i := 0; i < n; i++ {
		i := i
		tasks[i] = Task[string]{Work: func(ctx context.Context) string {
			// Earlier-index tasks deliberately finish LAST.
			time.Sleep(time.Duration(n-i) * 3 * time.Millisecond)
			return fmt.Sprintf("result-%d", i)
		}}
	}

	got := Dispatch(context.Background(), tasks, 5)

	for i, r := range got {
		want := fmt.Sprintf("result-%d", i)
		if r != want {
			t.Fatalf("results[%d] = %q, want %q (order not preserved)", i, r, want)
		}
	}
}

func TestConcurrencyBound(t *testing.T) {
	const n = 20
	const workers = 3

	var current atomic.Int32
	var maxSeen atomic.Int32

	tasks := make([]Task[int], n)
	for i := 0; i < n; i++ {
		tasks[i] = Task[int]{Work: func(ctx context.Context) int {
			c := current.Add(1)
			for {
				m := maxSeen.Load()
				if c <= m || maxSeen.CompareAndSwap(m, c) {
					break
				}
			}
			time.Sleep(5 * time.Millisecond)
			current.Add(-1)
			return 0
		}}
	}

	Dispatch(context.Background(), tasks, workers)

	if got := maxSeen.Load(); got > workers {
		t.Fatalf("observed %d concurrent tasks, want <= %d", got, workers)
	}
}

func TestEmptyTaskList(t *testing.T) {
	got := Dispatch[int](context.Background(), nil, 4)
	if len(got) != 0 {
		t.Fatalf("expected empty results, got %v", got)
	}
}
