package main

import (
	"fmt"
	"time"
)

// Demo: "ch <- i" is a SEND operation, not a plain assignment.
// This prints timestamps so you can see the sender blocking until a receiver is ready.

func unbufferedDemo() {
	fmt.Println("\n--- unbuffered channel: make(chan int) ---")
	ch := make(chan int)

	go func() {
		for i := 1; i <= 3; i++ {
			fmt.Printf("[sender]   about to send %d at %s\n", i, time.Now().Format("15:04:05.000"))
			ch <- i // blocks here until someone runs <-ch
			fmt.Printf("[sender]   sent %d (receiver picked it up) at %s\n", i, time.Now().Format("15:04:05.000"))
		}
		close(ch)
	}()

	// receiver deliberately sleeps before each receive, so you can see the sender
	// print "about to send" first, then wait, then print "sent" only once we receive.
	for v := range ch {
		time.Sleep(500 * time.Millisecond) // simulate receiver being "busy"
		fmt.Printf("[receiver] received %d at %s\n", v, time.Now().Format("15:04:05.000"))
	}
}

func bufferedDemo() {
	fmt.Println("\n--- buffered channel: make(chan int, 2) ---")
	ch := make(chan int, 2) // capacity 2

	go func() {
		for i := 1; i <= 3; i++ {
			fmt.Printf("[sender]   about to send %d at %s\n", i, time.Now().Format("15:04:05.000"))
			ch <- i // only blocks once the buffer (2 slots) is full
			fmt.Printf("[sender]   sent %d (buffer accepted it) at %s\n", i, time.Now().Format("15:04:05.000"))
		}
		close(ch)
	}()

	time.Sleep(1 * time.Second) // let sender fill the buffer first, uninterrupted
	for v := range ch {
		fmt.Printf("[receiver] received %d at %s\n", v, time.Now().Format("15:04:05.000"))
	}
}

func main() {
	unbufferedDemo()
	bufferedDemo()
}
