// Package cache is a teaching example for "Debug & Harden": a small
// service with realistic planted bugs, DIFFERENT from the leaderboard
// example in the practice paper, to broaden pattern recognition.
//
// This is the FIXED version. See NOTES.md for what was wrong in the
// original (buggy) version and why -- reproduced here in comments so you
// can see the before/after without needing two separate files.
package cache

import "sync"

// TTLCache is a simple in-memory cache where entries expire after a fixed
// duration measured in "ticks" (an abstract clock so tests don't need real
// time.Sleep -- also a common interview pattern: inject time instead of
// calling time.Now() directly).
type TTLCache struct {
	mu      sync.Mutex
	entries map[string]entry
	ttl     int64
}

type entry struct {
	value    string
	expireAt int64
}

func New(ttl int64) *TTLCache {
	return &TTLCache{entries: make(map[string]entry), ttl: ttl}
}

// Set stores value for key, expiring at now+ttl.
func (c *TTLCache) Set(key, value string, now int64) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entries[key] = entry{value: value, expireAt: now + c.ttl}
}

// Get returns the value for key if present and not expired as of now.
//
// BUG #1 (planted, now fixed): the original code checked
// `e.expireAt < now` to decide "expired", which meant an entry was still
// considered valid at the EXACT expiry tick (off-by-one: `<` should be
// `<=` -- an entry that expires "at" now should already be gone, mirroring
// the exact-boundary trap seen in the rate limiter problem). Fixed below.
func (c *TTLCache) Get(key string, now int64) (string, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	e, ok := c.entries[key]
	if !ok || e.expireAt <= now {
		return "", false
	}
	return e.value, true
}

// WarmAll seeds the cache with a batch of key/value pairs concurrently, one
// goroutine per pair, all sharing the same `now`.
//
// BUG #2 (planted, now fixed): the original code launched
//
//	for _, kv := range pairs {
//	    go func() { c.Set(kv.Key, kv.Value, now) }()
//	}
//
// which is the classic Go loop-variable-capture bug: `kv` is a single
// variable reused across iterations (pre-Go-1.22 semantics; Go 1.22+
// actually fixed this at the language level, but plenty of codebases still
// have go.mod pinned below 1.22, and interviewers love asking about it
// regardless). On older Go versions every goroutine could see the SAME
// (usually final) value of kv, silently writing the wrong key/value pairs.
// Fixed by passing kv explicitly into the closure, which is correct on
// every Go version.
func (c *TTLCache) WarmAll(pairs []KV, now int64) {
	var wg sync.WaitGroup
	for _, kv := range pairs {
		wg.Add(1)
		go func(kv KV) {
			defer wg.Done()
			c.Set(kv.Key, kv.Value, now)
		}(kv)
	}
	wg.Wait()
}

type KV struct {
	Key, Value string
}

// Len returns the number of entries currently stored, INCLUDING expired
// ones not yet evicted (Get is lazy-expiring; nothing proactively sweeps).
//
// BUG #3 (planted, now fixed): the original Len() did NOT take c.mu, while
// Set/Get did -- an unsynchronized read racing against synchronized writers
// is still a data race (Go's race detector flags ANY unsynchronized access
// to memory another goroutine writes, not just unsynchronized WRITES on
// both sides).
func (c *TTLCache) Len() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return len(c.entries)
}
