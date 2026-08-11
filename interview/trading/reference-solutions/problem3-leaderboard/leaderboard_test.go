package board

import (
	"bytes"
	"strings"
	"sync"
	"testing"
	"time"
)

// run: go test -race ./...

// Regression test for defect #1: after Stop() returns, no further snapshot
// lines may appear, and no background goroutine may remain running.
func TestStopFullyHaltsSnapshots(t *testing.T) {
	b := NewBoard()
	var buf bytes.Buffer
	var mu sync.Mutex // guards buf, since the snapshot goroutine writes concurrently
	b.out = syncWriter{&buf, &mu}

	b.AddTrade("AAA", 10)
	b.StartSnapshots(5 * time.Millisecond)

	time.Sleep(30 * time.Millisecond) // let a few ticks fire
	b.Stop()                          // must block until the goroutine has actually exited

	mu.Lock()
	countAtStop := strings.Count(buf.String(), "snapshot:")
	mu.Unlock()
	if countAtStop == 0 {
		t.Fatal("expected at least one snapshot line before Stop")
	}

	time.Sleep(50 * time.Millisecond) // give a leaking goroutine time to misbehave

	mu.Lock()
	countAfter := strings.Count(buf.String(), "snapshot:")
	mu.Unlock()
	if countAfter != countAtStop {
		t.Fatalf("snapshot lines kept appearing after Stop() returned: %d -> %d", countAtStop, countAfter)
	}
}

// Regression test for defect #2: fractional trade sizes must be preserved
// exactly, not truncated per-trade or via integer division in AvgSize.
func TestAvgSizeExactWithFractionalTrades(t *testing.T) {
	b := NewBoard()
	b.AddTrade("BBB", 1.5)
	b.AddTrade("BBB", 1.5)
	b.AddTrade("BBB", 1.5)
	b.AddTrade("BBB", 1.5)

	got := b.AvgSize("BBB")
	want := 1.5
	if got != want {
		t.Fatalf("AvgSize = %v, want %v (old code truncated each trade to 1, giving 1.0)", got, want)
	}

	b2 := NewBoard()
	b2.AddTrade("CCC", 0.1)
	b2.AddTrade("CCC", 0.2)
	got2 := b2.AvgSize("CCC")
	want2 := 0.15
	if diff := got2 - want2; diff > 1e-9 || diff < -1e-9 {
		t.Fatalf("AvgSize = %v, want ~%v", got2, want2)
	}
}

// Regression test for defect #3: ties must break alphabetically, every time,
// regardless of map iteration order.
func TestTopNTieBreakIsDeterministic(t *testing.T) {
	for i := 0; i < 50; i++ {
		b := NewBoard()
		// Equal volume for both -- must always resolve to [AAA, BBB] (alphabetical),
		// never flipping, regardless of Go's randomized map iteration order.
		b.AddTrade("BBB", 100)
		b.AddTrade("AAA", 100)

		got := b.TopN(2)
		want := []string{"AAA", "BBB"}
		if len(got) != 2 || got[0] != want[0] || got[1] != want[1] {
			t.Fatalf("iteration %d: TopN tie order = %v, want %v", i, got, want)
		}
	}
}

func TestTopNOrdersByVolumeThenAlpha(t *testing.T) {
	b := NewBoard()
	b.AddTrade("ZZZ", 300)
	b.AddTrade("AAA", 100)
	b.AddTrade("MMM", 100) // ties with AAA, must come after alphabetically
	b.AddTrade("BBB", 200)

	got := b.TopN(4)
	want := []string{"ZZZ", "BBB", "AAA", "MMM"}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("TopN = %v, want %v", got, want)
		}
	}
}

// Regression test for defect #4: AddTrade and the snapshot goroutine's TopN
// call must not race on the underlying map. Must pass under `go test -race`.
func TestConcurrentAddTradeAndSnapshotsNoRace(t *testing.T) {
	b := NewBoard()
	var discard bytes.Buffer
	var mu sync.Mutex
	b.out = syncWriter{&discard, &mu}

	b.StartSnapshots(1 * time.Millisecond)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < 500; i++ {
			b.AddTrade("AAA", 1.23)
		}
	}()
	wg.Wait()

	b.Stop()

	// Compare with a tolerance: summing 1.23 five hundred times via
	// float64 addition naturally accumulates a tiny epsilon of drift --
	// that's expected float behavior, not a defect.
	got := b.AvgSize("AAA")
	if diff := got - 1.23; diff > 1e-9 || diff < -1e-9 {
		t.Fatalf("AvgSize = %v, want ~1.23", got)
	}
}

// syncWriter serializes writes so the test's own reads of the buffer (via
// strings.Count) don't race with the snapshot goroutine's writes.
type syncWriter struct {
	buf *bytes.Buffer
	mu  *sync.Mutex
}

func (w syncWriter) Write(p []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.buf.Write(p)
}
