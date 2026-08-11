**Pattern this teaches:** a different set of 3 classic Go bugs than the
leaderboard exercise, chosen to broaden pattern recognition rather than
repeat it:

1. **Off-by-one on a time boundary** (`expireAt < now` vs `<= now`) — the
   exact same shape of bug as the rate limiter's exact-10000ms case and the
   leaderboard's fractional-truncation case: read the spec's comparison
   operators literally, don't assume.
2. **Loop-variable capture in a goroutine closure** — `for _, kv := range
   pairs { go func(){ use(kv) }() }` is one of THE most common real-world Go
   bugs (pre-Go-1.22 semantics: `kv` is one shared variable across every
   iteration). Fixed by passing it as a parameter: `go func(kv KV){...}(kv)`.
   Go 1.22+ changed the language so each loop iteration gets its own `kv`,
   but plenty of production codebases are pinned below 1.22, and interviewers
   ask about this specifically because it's such a common historical trap.
3. **Unsynchronized READ racing a synchronized WRITE** — `Len()` didn't take
   the mutex while `Set()`/`Get()` did. A common misconception is "only
   writes need locking, or only if both sides are unlocked" — false: Go's
   race detector (and the memory model) flags ANY unsynchronized access to
   memory another goroutine might be writing, read or write, doesn't matter.

**How to practice with this:** open `buggy_original.go.txt`, try to spot all
3 bugs from the "what production reported" hints at the top BEFORE reading
`cache.go`'s inline fix comments. That mirrors the real exercise format
(reported symptoms → find root cause → minimal fix → regression test) more
closely than reading the fixed version first.
