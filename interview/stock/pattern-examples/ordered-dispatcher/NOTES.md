**Pattern this teaches:** the core idea underneath the practice paper's much
more complex streaming `Pipeline` (Problem 2), stripped down to its essence:

> Give each task its ORIGINAL INDEX. Run tasks concurrently. Each goroutine
> writes its result into a pre-sized slice at ITS OWN index. Because every
> goroutine owns a different index, there's no data race on the writes
> themselves — you only need a `sync.WaitGroup` to know when everyone's done.

This is the pattern to reach for FIRST, before adding any of: streaming
output (a channel instead of a returned slice), timeouts, `Drain`/shutdown
semantics, or bounded memory. Get this batch version right, understand *why*
index-isolation avoids needing a mutex for the writes, and the streaming
version becomes "add a reorder buffer because now results need to leave one
at a time, in order, without waiting for the whole batch."

**Where this version stops short of the real Problem 2:**
- No timeout per task — a single slow `Work` call blocks the whole `Dispatch`
  call from returning, since we wait for everyone via `wg.Wait()`.
- No streaming — callers get nothing until ALL tasks finish, not incrementally.
- No graceful shutdown / drain — there's no way to say "stop early, give me
  what's done so far."
- A panic inside `Work` crashes the whole program (see the comment in
  `dispatcher.go`) — a hardened version would `recover()` per-task.

Each of those gaps is exactly one of the practice paper's Problem 2 rules
(1=order, 3=timeout, 5=drain, 6=bounded memory) — useful to trace explicitly:
which rule would you add, and what would it force you to change here?
