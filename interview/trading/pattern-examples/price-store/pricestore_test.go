package pricestore

import "testing"

// run: go test ./...

func TestLatestPrice(t *testing.T) {
	s := New()
	if _, ok := s.LatestPrice("AAPL"); ok {
		t.Fatal("expected no price for unknown symbol")
	}
	s.SetPrice("AAPL", 100, 150.0)
	s.SetPrice("AAPL", 200, 151.5)
	got, ok := s.LatestPrice("AAPL")
	if !ok || got != 151.5 {
		t.Fatalf("LatestPrice = %v, %v; want 151.5, true", got, ok)
	}
}

func TestPriceAsOfExactAndBetween(t *testing.T) {
	s := New()
	s.SetPrice("AAPL", 100, 150.0)
	s.SetPrice("AAPL", 200, 151.5)
	s.SetPrice("AAPL", 300, 149.0)

	cases := []struct {
		ts   int64
		want float64
		ok   bool
	}{
		{50, 0, false},     // before any observation
		{100, 150.0, true}, // exact match on first point
		{150, 150.0, true}, // between points -> most recent one <= ts
		{200, 151.5, true}, // exact match on second point
		{250, 151.5, true},
		{300, 149.0, true},
		{9999, 149.0, true}, // after last point -> still the latest known
	}
	for _, c := range cases {
		got, ok := s.PriceAsOf("AAPL", c.ts)
		if ok != c.ok || (ok && got != c.want) {
			t.Errorf("PriceAsOf(%d) = %v, %v; want %v, %v", c.ts, got, ok, c.want, c.ok)
		}
	}
}

func TestSetPriceOutOfOrderStillSortsCorrectly(t *testing.T) {
	s := New()
	// Even though the spec assumes non-decreasing ts, a robust store
	// shouldn't silently corrupt itself if that assumption is ever violated.
	s.SetPrice("AAPL", 300, 149.0)
	s.SetPrice("AAPL", 100, 150.0)
	s.SetPrice("AAPL", 200, 151.5)

	got, ok := s.PriceAsOf("AAPL", 250)
	if !ok || got != 151.5 {
		t.Fatalf("PriceAsOf(250) = %v, %v; want 151.5, true", got, ok)
	}
}

func TestInvalidInputsAreNoOps(t *testing.T) {
	s := New()
	s.SetPrice("", 100, 150.0)
	s.SetPrice("AAPL", -1, 150.0)
	s.SetPrice("AAPL", 100, 0)
	s.SetPrice("AAPL", 100, -5)
	if _, ok := s.LatestPrice("AAPL"); ok {
		t.Fatal("invalid SetPrice calls should not have created state")
	}
}
