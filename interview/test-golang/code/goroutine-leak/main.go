package main

import (
	"context"
	"fmt"
	"runtime"
	"time"
)

// --- BUGGY: unbuffered channel, nobody ever receives -> goroutine leaks forever ---
func leaky() {
	ch := make(chan int)
	go func() {
		ch <- 1 // blocks forever, no one reads
	}()
}

// --- FIX 1: buffered channel -> send doesn't need a receiver waiting, goroutine exits ---
func fixedBuffered() {
	ch := make(chan int, 1) // capacity 1, send succeeds immediately
	go func() {
		ch <- 1 // does not block, goroutine finishes right away
	}()
}

// --- FIX 2: context cancellation -> goroutine gives up instead of blocking forever ---
func fixedContext(ctx context.Context) {
	ch := make(chan int)
	go func() {
		select {
		case ch <- 1:
			// sent successfully
		case <-ctx.Done():
			// caller cancelled, stop waiting
			return
		}
	}()
}

// --- FIX 3: always have a receiver matching the sender ---
func fixedReceiver() {
	ch := make(chan int)
	go func() {
		ch <- 1
	}()
	<-ch // receive it, so the sender is unblocked and both goroutines finish
}

func printGoroutines(label string) {
	fmt.Printf("%-25s goroutines = %d\n", label, runtime.NumGoroutine())
}

func main() {
	printGoroutines("start")

	// 1. Leaky version — goroutines pile up and NEVER get cleaned
	for i := 0; i < 5; i++ {
		leaky()
	}
	time.Sleep(100 * time.Millisecond)
	printGoroutines("after 5x leaky()")

	// 2. Fix 1: buffered channel — no leak
	for i := 0; i < 5; i++ {
		fixedBuffered()
	}
	time.Sleep(100 * time.Millisecond)
	printGoroutines("after 5x fixedBuffered()")

	// 3. Fix 2: context cancellation — goroutine exits when we cancel
	ctx, cancel := context.WithCancel(context.Background())
	for i := 0; i < 5; i++ {
		fixedContext(ctx)
	}
	cancel() // nobody will ever receive, so cancel tells the goroutines to give up
	time.Sleep(100 * time.Millisecond)
	printGoroutines("after 5x fixedContext() + cancel")

	// 4. Fix 3: matching receiver — no leak
	for i := 0; i < 5; i++ {
		fixedReceiver()
	}
	time.Sleep(100 * time.Millisecond)
	printGoroutines("after 5x fixedReceiver()")
}
