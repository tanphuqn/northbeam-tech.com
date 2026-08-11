// Package board — per-symbol volume tracker with a leaderboard.
package board

import (
	"fmt"
	"io"
	"os"
	"sort"
	"sync"
	"time"
)

type stats struct {
	// totalQty was int64 in the original code, which truncated the
	// fractional part of every trade via int64(qty) before it was ever
	// accumulated (defect #2a), and AvgSize then did integer division on
	// top of that (defect #2b). Fractional quantities are normal per spec,
	// so this must be float64 end to end.
	totalQty float64
	count    int64
}

type Board struct {
	// mu protects stats (and the snapshot goroutine's lifecycle fields).
	// The original code had NO synchronization at all: AddTrade (writer)
	// and the snapshot goroutine's TopN call (reader, via a second
	// goroutine) touched the same map concurrently with no lock. That is a
	// genuine data race (undefined behavior, and Go maps can fatally crash
	// the process on concurrent read/write) -- defect #4, caught by
	// `go test -race`.
	mu    sync.Mutex
	stats map[string]*stats

	ticker *time.Ticker
	stopCh chan struct{}
	doneWG sync.WaitGroup
	out    io.Writer // defaults to os.Stdout; overridable in tests
}

func NewBoard() *Board {
	return &Board{
		stats: make(map[string]*stats),
		out:   os.Stdout,
	}
}

// AddTrade records a trade. qty is a positive decimal quantity.
func (b *Board) AddTrade(symbol string, qty float64) {
	b.mu.Lock()
	defer b.mu.Unlock()

	s, ok := b.stats[symbol]
	if !ok {
		s = &stats{}
		b.stats[symbol] = s
	}
	s.totalQty += qty
	s.count++
}

// AvgSize returns the exact average trade size for the symbol.
func (b *Board) AvgSize(symbol string) float64 {
	b.mu.Lock()
	defer b.mu.Unlock()

	s, ok := b.stats[symbol]
	if !ok || s.count == 0 {
		return 0
	}
	// Both operands must be float64 BEFORE dividing. The original code did
	// `float64(s.totalQty / s.count)` with both fields as int64 -- integer
	// division truncated the result before the float64 conversion even ran.
	return s.totalQty / float64(s.count)
}

// TopN returns the n symbols with highest total volume, ties broken
// alphabetically. The same input must always produce the same output.
func (b *Board) TopN(n int) []string {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.topNLocked(n)
}

// topNLocked assumes b.mu is already held.
func (b *Board) topNLocked(n int) []string {
	symbols := make([]string, 0, len(b.stats))
	for sym := range b.stats {
		symbols = append(symbols, sym)
	}
	// Go map iteration order is randomized, and the original comparator only
	// compared by totalQty -- sort.Slice is NOT stable, so ties kept
	// whatever random relative order the map produced, flipping between
	// runs (defect #3). Fix: make the alphabetical tiebreak part of the
	// comparator itself, so the result is deterministic regardless of
	// starting order or sort stability.
	sort.Slice(symbols, func(i, j int) bool {
		qi, qj := b.stats[symbols[i]].totalQty, b.stats[symbols[j]].totalQty
		if qi != qj {
			return qi > qj
		}
		return symbols[i] < symbols[j]
	})
	if n > len(symbols) {
		n = len(symbols)
	}
	if n < 0 {
		n = 0
	}
	return symbols[:n]
}

// StartSnapshots logs the top 3 periodically until Stop is called.
func (b *Board) StartSnapshots(every time.Duration) {
	b.mu.Lock()
	defer b.mu.Unlock()

	// Guard against a second StartSnapshots leaking a previous goroutine.
	// Not explicitly reported, but the same class of bug as the Stop defect,
	// so hardened defensively while we're in here.
	if b.ticker != nil {
		b.stopLocked()
	}

	b.ticker = time.NewTicker(every)
	b.stopCh = make(chan struct{})
	ticker := b.ticker
	stopCh := b.stopCh

	b.doneWG.Add(1)
	go func() {
		defer b.doneWG.Done()
		for {
			select {
			case <-ticker.C:
				b.mu.Lock()
				top := b.topNLocked(3)
				b.mu.Unlock()
				fmt.Fprintln(b.out, "snapshot:", top)
			case <-stopCh:
				return
			}
		}
	}()
}

// Stop ends snapshot logging. After Stop returns, no further snapshots
// appear and no background work remains running.
//
// The original code just set b.ticker = nil. That did nothing to the
// running goroutine: `for range b.ticker.C` (as originally written)
// evaluates b.ticker.C ONCE when the loop starts and keeps reading that same
// channel forever -- reassigning the field later is invisible to it, the
// underlying *time.Ticker is never Stop()'d (a resource leak on top of the
// logic bug), and Stop() returned immediately without waiting for anything,
// so callers had no way to know the goroutine was still running (defect #1).
func (b *Board) Stop() {
	b.mu.Lock()
	b.stopLocked()
	b.mu.Unlock()

	// Wait outside the lock: the goroutine takes b.mu itself when logging a
	// snapshot, so waiting while still holding it would deadlock against a
	// snapshot currently in flight.
	b.doneWG.Wait()
}

// stopLocked assumes b.mu is already held.
func (b *Board) stopLocked() {
	if b.ticker == nil {
		return // already stopped / never started; idempotent
	}
	b.ticker.Stop()
	close(b.stopCh)
	b.ticker = nil
	b.stopCh = nil
}
