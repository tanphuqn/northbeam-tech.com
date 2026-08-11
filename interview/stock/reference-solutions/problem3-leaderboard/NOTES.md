=== NOTES ===

I found **four** defects — the spec's three reported symptoms map to four
distinct root causes because one symptom (growing background activity) was
actually caused by two independent bugs stacked together, and the test
requirement ("must pass with `go test -race`") was itself a strong hint at a
fourth, unreported one.

**Defect #1 — `Stop()` doesn't stop anything (symptom: snapshots keep
appearing after `stop()`, background activity keeps growing).**
Root cause: the original `StartSnapshots` did `for range b.ticker.C` — in Go,
the channel expression in a `range` clause is evaluated **once**, when the
loop starts. `Stop()` then did `b.ticker = nil`, but the already-running
goroutine had already captured the old ticker's channel value; reassigning
the struct field later is invisible to it. The ticker itself was also never
`.Stop()`'d (a resource leak independent of the logic bug), and `Stop()`
returned immediately with no way for a caller to know the goroutine was still
alive.
Fix: the goroutine now selects on the ticker's channel *and* a `stopCh` that
`Stop()` closes; `Stop()` also calls `ticker.Stop()` and blocks on a
`sync.WaitGroup` until the goroutine has actually exited, so "after `stop()`
returns, no further snapshots appear and no background work remains running"
is a real, synchronous guarantee rather than "eventually, probably."
Regression test: `TestStopFullyHaltsSnapshots` — starts snapshots, waits for
some to fire, calls `Stop()`, then waits again and asserts the log line count
didn't grow.

**Defect #2 — average trade size is low for fractional trades (two bugs
stacked).**
2a: `s.totalQty += int64(qty)` truncated the fractional part of *every trade*
at accumulation time, not just at the end — `1.5` became `1` before it was
ever added.
2b: `AvgSize` computed `float64(s.totalQty / s.count)` — since both were
`int64`, the division happened in integer arithmetic (truncating again)
*before* the result was converted to `float64`.
Fix: `totalQty` is now `float64`, accumulated without truncation; `AvgSize`
divides as `s.totalQty / float64(s.count)`, float division throughout.
Regression test: `TestAvgSizeExactWithFractionalTrades` — four trades of 1.5
must average to exactly 1.5 (old code gives 1.0).

**Defect #3 — leaderboard ties flip between runs.**
Root cause: `for sym := range b.stats` iterates a Go map in **randomized**
order (a language guarantee, not a bug in itself), and `sort.Slice`'s
comparator only compared `totalQty` — for equal values, `sort.Slice` is not
guaranteed stable, so ties kept whatever random relative order the map
produced going in.
Fix: the comparator itself now breaks ties alphabetically
(`if qi != qj { return qi > qj }; return symbols[i] < symbols[j]`), so the
result is deterministic regardless of map iteration order or sort stability.
Regression test: `TestTopNTieBreakIsDeterministic` — runs 50 fresh `Board`s
(new random map seed each time) with two equal-volume symbols and asserts the
alphabetical order every single time.

**Defect #4 — unsynchronized concurrent map access (not in the reported
symptoms, but implied by the `-race` requirement).**
Root cause: `AddTrade` (writer) and the snapshot goroutine's `TopN` call
(reader, from a second goroutine) touched the same `map[string]*stats` with
no synchronization at all. This is undefined behavior in Go and can crash
the process outright with "concurrent map read and map write" — not just a
theoretical race-detector nitpick.
Fix: added a `sync.Mutex` guarding `stats` (and the ticker/stopCh lifecycle
fields), held during `AddTrade`, `AvgSize`, and `TopN`.
Regression test: `TestConcurrentAddTradeAndSnapshotsNoRace` — runs
`AddTrade` in a loop concurrently with an active snapshot goroutine; passes
only under `-race` with the fix (old code either races or occasionally
panics with a fatal concurrent-map error).

**What I'd do with more time:**
- `StartSnapshots` called twice in a row now defensively stops the previous
  goroutine first — not in the reported symptoms, but the same bug class,
  so I hardened it while already in this code.
- Consider `sync.RWMutex` instead of `sync.Mutex` if reads (`AvgSize`,
  `TopN`) end up far more frequent than writes (`AddTrade`) in production —
  not changed here since correctness, not throughput, was the ask.
