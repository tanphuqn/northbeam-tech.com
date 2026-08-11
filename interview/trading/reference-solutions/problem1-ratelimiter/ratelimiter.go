// Package ratelimiter implements an in-memory sliding-window rate limiter.
package ratelimiter

// windowMs is the rolling window size in milliseconds.
const windowMs = 10000

// maxRequests is the max number of accepted requests allowed inside the window.
const maxRequests = 5

// RateLimiter enforces a per-client sliding-window limit of maxRequests
// accepted requests per windowMs milliseconds.
//
// Not safe for concurrent use (single-threaded correctness only, per spec).
type RateLimiter struct {
	// clients maps clientID -> accepted timestamps, oldest first.
	// Since callers guarantee `now` is non-decreasing per client, and we only
	// ever append accepted timestamps in call order, this slice is always
	// sorted ascending. That lets us evict expired entries from the front
	// in O(k) instead of scanning/sorting on every call.
	clients map[string][]int64
}

// New creates an empty RateLimiter.
func New() *RateLimiter {
	return &RateLimiter{clients: make(map[string][]int64)}
}

// Allow reports whether the request from clientID at time now (ms) is accepted.
// Invalid arguments (empty clientID, negative now) are rejected without
// changing any state.
func (r *RateLimiter) Allow(clientID string, now int64) bool {
	if clientID == "" || now < 0 {
		return false
	}

	ts := r.clients[clientID]
	ts = evict(ts, now)

	if len(ts) >= maxRequests {
		// Rejected: record nothing, but DO persist the eviction we just did
		// above, so a client that's been idle regains capacity even while
		// being hammered with requests that end up rejected.
		r.clients[clientID] = ts
		return false
	}

	ts = append(ts, now)
	r.clients[clientID] = ts
	return true
}

// ActiveClients returns the number of clients with at least one accepted
// request currently inside the window as of now.
func (r *RateLimiter) ActiveClients(now int64) int {
	count := 0
	for id, ts := range r.clients {
		ts = evict(ts, now)
		r.clients[id] = ts // opportunistic cleanup; keeps future scans cheap
		if len(ts) > 0 {
			count++
		}
	}
	return count
}

// evict drops timestamps that have aged out of the window: a request with
// timestamp t is inside the window while now - t < windowMs; a request aged
// EXACTLY windowMs is outside (per spec) and must be evicted.
func evict(ts []int64, now int64) []int64 {
	i := 0
	for i < len(ts) && now-ts[i] >= windowMs {
		i++
	}
	if i == 0 {
		return ts
	}
	// Reslice to drop the expired prefix. We deliberately don't reuse the
	// backing array via copy() here to keep this simple; at max 5 elements
	// per client this is not a real cost.
	out := make([]int64, len(ts)-i)
	copy(out, ts[i:])
	return out
}
