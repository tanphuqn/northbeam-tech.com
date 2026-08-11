// Package pricestore is a teaching example for the "stateful correctness"
// category: implement a small in-memory component exactly to spec, including
// point-in-time historical lookups.
//
// This is a DIFFERENT problem from the rate limiter (which only ever needed
// "recent" state within a rolling window). A price store needs to answer
// "what was the price at time T", which means keeping full history per
// symbol and doing a point-in-time lookup -- a different core skill
// (binary search over time-ordered data) than the rate limiter's sliding
// window eviction.
package pricestore

import "sort"

type pricePoint struct {
	ts    int64
	price float64
}

// Store tracks a price history per symbol. Not safe for concurrent use
// (matches the "single-threaded correctness is enough" pattern seen in
// Problem 1 of the practice paper -- concurrency, if required, is usually a
// separate problem).
type Store struct {
	// prices[symbol] is kept sorted ascending by ts. Callers are assumed to
	// report non-decreasing timestamps per symbol (same assumption the
	// practice paper's rate limiter made) so SetPrice can just append.
	prices map[string][]pricePoint
}

func New() *Store {
	return &Store{prices: make(map[string][]pricePoint)}
}

// SetPrice records a new price observation for symbol at ts (ms).
// Invalid input (empty symbol, negative ts, non-positive price) is a no-op.
func (s *Store) SetPrice(symbol string, ts int64, price float64) {
	if symbol == "" || ts < 0 || price <= 0 {
		return
	}
	history := s.prices[symbol]
	// Defensive: if ts ever arrives out of order despite the stated
	// assumption, don't silently corrupt the sort invariant -- insert at the
	// correct position instead of blindly appending.
	i := sort.Search(len(history), func(i int) bool { return history[i].ts > ts })
	history = append(history, pricePoint{})
	copy(history[i+1:], history[i:])
	history[i] = pricePoint{ts: ts, price: price}
	s.prices[symbol] = history
}

// LatestPrice returns the most recent price for symbol and whether one exists.
func (s *Store) LatestPrice(symbol string) (float64, bool) {
	history := s.prices[symbol]
	if len(history) == 0 {
		return 0, false
	}
	return history[len(history)-1].price, true
}

// PriceAsOf returns the price that was in effect at time ts: the most recent
// observation with timestamp <= ts. Returns false if no observation exists
// at or before ts (e.g. ts is before the symbol's first recorded price).
func (s *Store) PriceAsOf(symbol string, ts int64) (float64, bool) {
	history := s.prices[symbol]
	// Find the first index with ts_i > ts; the answer is the entry just
	// before that (the last one that's <= ts).
	i := sort.Search(len(history), func(i int) bool { return history[i].ts > ts })
	if i == 0 {
		return 0, false // every observation is after ts
	}
	return history[i-1].price, true
}
