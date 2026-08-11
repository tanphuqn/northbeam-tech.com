package ratelimiter

import "testing"

// run: go test ./...

func TestBasicAcceptUpToLimit(t *testing.T) {
	rl := New()
	for i := 0; i < 5; i++ {
		if !rl.Allow("alice", int64(i)) {
			t.Fatalf("request %d should be accepted", i)
		}
	}
	if rl.Allow("alice", 5) {
		t.Fatal("6th request within window should be rejected")
	}
}

func TestExactWindowBoundaryIsOutside(t *testing.T) {
	rl := New()
	// Fill the window at t=0..4 (5 accepted requests).
	for i := int64(0); i < 5; i++ {
		if !rl.Allow("alice", i) {
			t.Fatalf("request at t=%d should be accepted", i)
		}
	}
	// t=0 request ages out only when now - 0 >= 10000, i.e. now == 10000.
	// At now=9999 it is still inside (9999 < 10000): must still be full.
	if rl.Allow("alice", 9999) {
		t.Fatal("request at t=9999 should be rejected: t=0 entry is still inside the window")
	}
	// At now=10000, t=0 request is aged EXACTLY 10000ms -> outside -> evicted,
	// freeing one slot.
	if !rl.Allow("alice", 10000) {
		t.Fatal("request at t=10000 should be accepted: t=0 entry aged out exactly at boundary")
	}
}

func TestRejectedRequestsAreInvisible(t *testing.T) {
	rl := New()
	// Fill the window.
	for i := int64(0); i < 5; i++ {
		rl.Allow("alice", i)
	}
	// Hammer with many rejected requests; none should count toward the window.
	for i := int64(5); i < 9999; i++ {
		if rl.Allow("alice", i) {
			t.Fatalf("request at t=%d should be rejected (window full)", i)
		}
	}
	// The moment the original accepted t=0 request ages out (now=10000),
	// capacity must be regained immediately -- unaffected by the flood of
	// rejects in between.
	if !rl.Allow("alice", 10000) {
		t.Fatal("client should regain capacity the instant its oldest accepted request ages out")
	}
}

func TestPerClientIndependence(t *testing.T) {
	rl := New()
	for i := 0; i < 5; i++ {
		if !rl.Allow("alice", int64(i)) {
			t.Fatalf("alice request %d should be accepted", i)
		}
	}
	if rl.Allow("alice", 5) {
		t.Fatal("alice should be rate-limited")
	}
	// bob is a different client and must not be affected by alice's usage.
	if !rl.Allow("bob", 5) {
		t.Fatal("bob's first request should be accepted regardless of alice's state")
	}
}

func TestInvalidArgumentsChangeNothing(t *testing.T) {
	rl := New()
	if rl.Allow("", 100) {
		t.Fatal("empty clientID must be rejected")
	}
	if rl.Allow("alice", -1) {
		t.Fatal("negative now must be rejected")
	}
	// Neither invalid call should have created state / consumed capacity.
	for i := 0; i < 5; i++ {
		if !rl.Allow("alice", int64(i)) {
			t.Fatalf("alice request %d should be accepted; invalid calls must not have consumed capacity", i)
		}
	}
	if rl.ActiveClients(0) != 1 {
		t.Fatalf("expected 1 active client (alice), got %d", rl.ActiveClients(0))
	}
}

func TestActiveClients(t *testing.T) {
	rl := New()
	rl.Allow("alice", 0)
	rl.Allow("bob", 0)
	if got := rl.ActiveClients(0); got != 2 {
		t.Fatalf("expected 2 active clients, got %d", got)
	}
	// Once both clients' only accepted request ages out, they no longer count.
	if got := rl.ActiveClients(10000); got != 0 {
		t.Fatalf("expected 0 active clients after window expiry, got %d", got)
	}
}

func TestActiveClientsIgnoresRejectedOnlyClients(t *testing.T) {
	rl := New()
	for i := 0; i < 5; i++ {
		rl.Allow("alice", int64(i))
	}
	// This request is rejected (window full for alice); alice is still active
	// only because of her earlier *accepted* requests, not because of this one.
	rl.Allow("alice", 5)
	if got := rl.ActiveClients(5); got != 1 {
		t.Fatalf("expected 1 active client, got %d", got)
	}
}
