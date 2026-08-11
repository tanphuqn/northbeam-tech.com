=== NOTES ===

**Data structure.** `map[clientID][]int64` of accepted timestamps, oldest first.
Since the caller guarantees `now` is non-decreasing per client, and we only
ever append at the end, the slice is always sorted ascending — eviction is
just trimming a prefix, no sorting or heap needed. Per-client slice length is
capped at 5 (maxRequests), so this is O(1)-ish in practice regardless of
traffic volume.

**Why not a real deque / ring buffer?** At a cap of 5 elements, a plain slice
with `append` + reslice is simpler and fast enough; a ring buffer would only
pay off at much higher per-client limits.

**Edge cases covered:**
- Exact window boundary: a request aged *exactly* 10000ms is outside (tested
  explicitly — this is the "buried" edge case the spec calls out).
- Rejected requests never consume capacity, including a flood of thousands of
  rejects between two accepted requests.
- Per-client isolation — one client's state never affects another's.
- Invalid arguments (empty clientID, negative now) change no state.
- `activeClients` only counts clients with at least one *accepted* request
  still inside the window, not clients who only ever got rejected.

**What I'd do with more time:**
- `activeClients` currently does an opportunistic full-map eviction sweep,
  which is O(clients). For very large client counts, a lazy strategy (evict
  per-client only on that client's next `Allow` call, plus a periodic sweep)
  would avoid the full scan on every `activeClients` call.
- The map never removes a client key once created, even after all its
  timestamps evict to empty — for a long-lived process with a huge number of
  distinct one-off clients, this is an unbounded-map-keys leak. Worth
  addressing if `clientID` cardinality is attacker-controlled.
