=== NOTES ===

**Design in a paragraph.** Each `Submit` is assigned a monotonically
increasing internal `seq` (not `Msg.ID` itself — ordering is by *submission
order*, which the spec guarantees matches ID order, so this sidesteps having
to assume IDs are contiguous). A single dispatch loop reads from a bounded
input channel and, per task, acquires a semaphore slot (capacity =
`workers`) before spawning a goroutine that races the enricher against
`perMsgTimeout` via `context.WithTimeout` + `select`. Completed results
(success or timeout) are stored in a `map[seq]Out` "reorder buffer" guarded
by a mutex; the same critical section walks forward from `nextSeq`, popping
and sending any now-contiguous run to `Out()`. **The send happens while still
holding the mutex** — an earlier version sent after unlocking, which let a
later-seq goroutine race ahead and send before an earlier-seq goroutine got
scheduled, corrupting order under real concurrency. This only surfaced under
`-race` with many workers and variable delays; worth flagging because it's
exactly the kind of interaction rule 1 + rule 2 create together.

**How memory is bounded.** Input queue and output queue are each sized
`2*workers`; `Submit` blocks once the input queue is full, and the
worker-count semaphore caps concurrent enrichment attempts. Since the reorder
buffer only ever holds entries for tasks that have already been accepted
(bounded by the input queue) but not yet fully enriched, it can never exceed
roughly `2*workers` entries either — it's a consequence of the other two
bounds, not a separately-bounded structure.

**The straggler guarantee I chose.** A timed-out enrichment's own goroutine
is *not* forcibly killed (Go has no such mechanism, and the spec says the
enricher may ignore `ctx`) — it's abandoned to finish or leak in the
background, while the pipeline immediately emits `TimeoutErr` in the
straggler's correct ordered position and frees its worker-semaphore slot so
later messages aren't blocked. This means, briefly, more than `workers`
enrichments can be physically running at once (the straggler + a newly
started one) — `workers` bounds *awaited* concurrency, not literal OS-level
concurrency. I considered forcing strict `workers` concurrency (not starting
a new task until the straggler's goroutine actually returns), but that would
let one badly-behaved enricher call permanently throttle the whole pipeline
to fewer effective workers — worse than the trade-off I chose, and directly
against rule 3 ("must not stall the output stream").

**Drain semantics.** `Drain` stops new `Submit`s (documented: calling
`Submit` after `Drain` has started panics — fail loud rather than silently
drop, since it almost always means a shutdown-ordering bug in the caller),
waits up to `ctx`'s deadline for everything already accepted to resolve, then
closes an internal `stopEmit` gate so any sends still in flight (or arriving
late, from a straggler that finishes after the deadline) drop cleanly instead
of touching a channel that might already be closed. `dropped` is computed as
`submitted - emitted`, which correctly counts both "never started" and
"finished but stuck behind a still-unresolved earlier message" as dropped —
I don't try to distinguish those two cases since the spec only asks for a
count.

**Caveat / contract worth stating explicitly:** `Drain` assumes a consumer is
reading (or will read) `Out()` around the same time — if the output buffer is
full and nothing drains it, a pending send inside the mutex-held critical
section can block `Drain`'s own wait indefinitely. This mirrors ordinary
backpressure everywhere else in the pipeline; I'd flag it as a real risk in
a production system and would consider giving `Drain` its own bypass path
that force-closes `Out()` after the deadline regardless of pending sends, at
the cost of a possible partial/incomplete final read.

**What I'd harden with more time:**
- Replace the raw `int64`/`int32` + `atomic.*` fields with `atomic.Int64` /
  `atomic.Bool` for clarity (functionally identical, just newer API).
- Add a metrics/logging hook for dropped and timed-out messages — useful in
  production, out of scope for the exercise.
- Stress-test with much higher worker counts and adversarial enrichers
  (e.g., one that panics) to confirm panic-safety inside `runTask`'s inner
  goroutine doesn't take down the whole pipeline.
