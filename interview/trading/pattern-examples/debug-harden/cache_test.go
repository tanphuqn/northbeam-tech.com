package cache

import (
	"fmt"
	"testing"
)

// run: go test -race ./...

// Regression test for BUG #1: an entry must be gone at the EXACT expiry tick.
func TestExpiryExactBoundary(t *testing.T) {
	c := New(10)
	c.Set("k", "v", 0) // expires at tick 10

	if _, ok := c.Get("k", 9); !ok {
		t.Fatal("entry should still be present at tick 9 (not yet expired)")
	}
	if _, ok := c.Get("k", 10); ok {
		t.Fatal("entry should be expired at tick 10 (expireAt <= now)")
	}
}

// Regression test for BUG #2: every pair in a WarmAll batch must be stored
// under its OWN key/value, not all collapsing onto the loop variable's final
// value.
func TestWarmAllStoresEachPairIndependently(t *testing.T) {
	c := New(100)
	const n = 50
	pairs := make([]KV, n)
	for i := 0; i < n; i++ {
		pairs[i] = KV{Key: fmt.Sprintf("k%d", i), Value: fmt.Sprintf("v%d", i)}
	}

	c.WarmAll(pairs, 0)

	for i := 0; i < n; i++ {
		want := fmt.Sprintf("v%d", i)
		got, ok := c.Get(fmt.Sprintf("k%d", i), 0)
		if !ok || got != want {
			t.Fatalf("key k%d = %q, %v; want %q, true (loop-capture bug would collapse these onto the last pair)", i, got, ok, want)
		}
	}
}

// Regression test for BUG #3: Len() must be synchronized against concurrent
// Set() calls. Must pass under `go test -race`.
func TestLenNoRaceWithConcurrentSet(t *testing.T) {
	c := New(100)
	done := make(chan struct{})
	go func() {
		defer close(done)
		for i := 0; i < 500; i++ {
			c.Set(fmt.Sprintf("k%d", i), "v", 0)
		}
	}()
	for i := 0; i < 500; i++ {
		_ = c.Len()
	}
	<-done
}
