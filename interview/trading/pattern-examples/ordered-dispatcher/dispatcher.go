// Package dispatcher is a teaching example for the "ordered parallel
// concurrency" category, showing the core pattern in its simplest form:
// batch fan-out/fan-in with index-preserved results.
//
// This is deliberately SIMPLER than the practice paper's streaming Pipeline
// (no incremental Out() channel, no Drain, no bounded-memory backpressure) --
// it isolates the one idea those problems all build on: give each task its
// original index, run them concurrently, and write each result into a
// pre-sized results slice at that index. Because each goroutine writes to a
// *different* slice index, no locking is needed for the writes themselves --
// only a WaitGroup to know when everyone's done. This is the pattern to
// reach for first before adding streaming/timeout/backpressure complexity.
package dispatcher

import (
	"context"
	"sync"
)

// Task is anything dispatchable; Work does the (possibly slow) processing.
type Task[T any] struct {
	Work func(ctx context.Context) T
}

// Dispatch runs up to `workers` tasks concurrently and returns their results
// in the SAME order as the input tasks, regardless of completion order.
func Dispatch[T any](ctx context.Context, tasks []Task[T], workers int) []T {
	if workers < 1 {
		workers = 1
	}
	results := make([]T, len(tasks)) // pre-sized: each goroutine owns a distinct index
	sem := make(chan struct{}, workers)
	var wg sync.WaitGroup

	for i, task := range tasks {
		wg.Add(1)
		sem <- struct{}{} // acquire a slot; blocks once `workers` are busy
		go func(i int, task Task[T]) {
			defer wg.Done()
			defer func() { <-sem }()
			results[i] = task.Work(ctx)
			// NOTE: a panic inside Work is NOT recovered here -- it will
			// crash the whole program, since a goroutine panic can only be
			// recovered from within that same goroutine. A production
			// version would wrap this call in its own recover() and turn a
			// panic into an error result instead. Left out here to keep the
			// core pattern visible; see NOTES.md.
		}(i, task)
	}

	wg.Wait()
	return results
}
