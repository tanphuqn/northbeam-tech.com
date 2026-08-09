# Answers — see questions.md

[← Quay về note.md](../note.md)

Only open this AFTER attempting `questions.md` yourself. See `note.md` for fuller explanations and more examples.

---

## Section A: Goroutines & Channels

**A1.** A goroutine is managed by the Go runtime, not mapped 1:1 to an OS thread. Its stack starts at ~2KB and grows/shrinks dynamically, and the Go scheduler (M:N scheduling) multiplexes many goroutines onto a small number of real OS threads. OS threads typically reserve 1-8MB of fixed stack plus kernel-managed context-switch overhead, making them far more expensive.

**A2.** The program exits immediately; any goroutines still running are killed with no cleanup — deferred calls don't run. This is why you need `sync.WaitGroup` or a channel to wait for completion before `main()` returns.

**A3.** A goroutine leak is a goroutine blocked forever (waiting on a channel nobody sends to/receives from) that never exits, living until the program terminates.
```go
func leaky() {
    ch := make(chan int) // unbuffered
    go func() {
        ch <- 1 // blocks forever, nobody reads
    }()
} // function returns, that goroutine leaks
```

**A4.** Unbuffered channel: send and receive must happen at the same time (rendezvous) — send blocks until a receiver is ready. Buffered channel: send doesn't block until the buffer is full; receive doesn't block if the buffer has an item.

**A5.**
a. Panic: "send on closed channel"
b. Panic: "close of closed channel"
c. Returns the zero value immediately, doesn't block — use `v, ok := <-ch`; `ok == false` means the channel is closed and drained.

**A6.** The sender should close the channel, never the receiver. Reason: the receiver has no way to know whether another sender still intends to send; if the receiver closes it and a sender sends afterward, it panics.

**A7.** Fan-out: multiple goroutines read from one input channel to process work in parallel (e.g. a worker pool). Fan-in: merge results from multiple channels into a single channel. Example: parallel image resizing — fan-out to split work across workers, fan-in to collect results.

---

## Section B: sync package

**B1.** `Mutex`: use when reads and writes happen at similar rates, or writes dominate — it locks exclusively for both. `RWMutex`: use for read-heavy workloads where writes are rare — multiple readers can hold `RLock()` concurrently, only a writer blocks everyone. Example: a config cache that's read constantly but rarely updated fits `RWMutex` well.

**B2.** Maps have no built-in synchronization. When 2+ goroutines write to the same map concurrently, the Go runtime detects it and **panics**: "fatal error: concurrent map writes" — this is a fatal runtime error, not a recoverable one in most cases; the program crashes.

**B3.** Bug: `wg.Add(1)` is called **inside** the goroutine, creating a race between `Add()` and `Wait()`: if `Wait()` runs before any goroutine calls `Add(1)`, it sees counter = 0 and returns immediately without waiting for anything.
```go
func correct() {
    var wg sync.WaitGroup
    for i := 0; i < 3; i++ {
        wg.Add(1) // called in the parent, BEFORE go func()
        go func(i int) {
            defer wg.Done()
            fmt.Println(i)
        }(i)
    }
    wg.Wait()
}
```

**B4.** `sync.Once` guarantees a block of code runs exactly once even when called from multiple goroutines concurrently — commonly used for singleton initialization (e.g. setting up a connection pool, loading config once).
```go
var once sync.Once
var instance *Config

func GetConfig() *Config {
    once.Do(func() {
        instance = loadConfig()
    })
    return instance
}
```

**B5.** `atomic` is faster than `Mutex` for simple operations (increment/compare-and-swap) because it uses low-level, lock-free CPU instructions instead of goroutine context-switching under contention. But `atomic` only works for simple single-variable operations; if you need to protect a more complex block of logic (multiple variables, multiple steps), you need `Mutex`.

---

## Section C: Context & Race Conditions

**C1.** `context` propagates cancellation **down a chain** (parent cancels → every child observes it via `ctx.Done()`), has built-in deadline/timeout support, and is Go's standard convention — most libraries (net/http, database/sql, ...) accept a `context` to auto-cancel when needed. A hand-rolled timer/channel doesn't propagate across multiple call layers and is easy to get wrong.

**C2.**
- `WithCancel`: returns a `cancel()` function you call manually whenever you decide to cancel.
- `WithTimeout`: auto-cancels after a given duration from when it's called.
- `WithDeadline`: auto-cancels at a specific absolute `time.Time`.
(`WithTimeout` is actually implemented as `WithDeadline(parent, time.Now().Add(duration))` internally.)

**C3.** If you don't call `cancel()`, the context and its internal resources (the timer inside `WithTimeout`/`WithDeadline`) aren't released promptly, causing a leak — especially bad if contexts are created in a loop or on every request. `defer cancel()` guarantees cleanup regardless of which return path the function takes.

**C4.** Use it for **request-scoped** values that travel with a single request as metadata (request ID, trace ID, user ID for logging/tracing). Do NOT use it to pass config, dependencies, or optional function parameters — `WithValue` isn't type-safe (it uses `interface{}`), is hard to trace, and hides a function's real dependencies (those should be explicit parameters).

**C5.** A race condition is when 2+ goroutines access the same variable, with at least one write, without synchronization — the result depends on non-deterministic execution order. Tool: Go's race detector (ThreadSanitizer), run with `go run -race main.go` or `go test -race ./...`.

**C6.** Before Go 1.22, the loop variable in `for i := ...` was **reused/shared** across all iterations — if a goroutine closure captured that variable directly (instead of passing it as a parameter), all goroutines could end up reading the same final value once they actually ran (since the loop may have already finished by then).
```go
// BUG: may print "3 3 3" instead of "0 1 2"
for i := 0; i < 3; i++ {
    go func() { fmt.Println(i) }()
}

// FIX: pass i as a parameter so each goroutine gets its own copy
for i := 0; i < 3; i++ {
    go func(i int) { fmt.Println(i) }(i)
}
```
Since Go 1.22, each loop iteration creates a fresh `i`, so this bug no longer happens by default — but it's still a classic interview question to test fundamentals.

**C7.** Yes, there's a race condition at `counter++` — it's a non-atomic read-modify-write, and concurrent goroutines doing this simultaneously will lose updates.
```go
func count() int64 {
    var counter int64
    var wg sync.WaitGroup
    for i := 0; i < 100; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            atomic.AddInt64(&counter, 1) // or use a Mutex
        }()
    }
    wg.Wait()
    return atomic.LoadInt64(&counter)
}
```

---

## Section D: Live-coding

See full sample code in `note.md` section 7 (Exercises 1, 2, 3). For D4 and D5, here's the approach:

**D4 (merge/fan-in multiple channels):**
```go
func merge(channels ...<-chan int) <-chan int {
    out := make(chan int)
    var wg sync.WaitGroup
    wg.Add(len(channels))

    for _, ch := range channels {
        go func(c <-chan int) {
            defer wg.Done()
            for v := range c {
                out <- v
            }
        }(ch)
    }

    go func() {
        wg.Wait()
        close(out)
    }()

    return out
}
```

**D5 (bounded queue, multiple producers/consumers):**
```go
type BoundedQueue struct {
    ch chan int
}

func NewBoundedQueue(capacity int) *BoundedQueue {
    return &BoundedQueue{ch: make(chan int, capacity)}
}

func (q *BoundedQueue) Push(ctx context.Context, v int) error {
    select {
    case q.ch <- v:
        return nil
    case <-ctx.Done():
        return ctx.Err()
    }
}

func (q *BoundedQueue) Pop(ctx context.Context) (int, error) {
    select {
    case v := <-q.ch:
        return v, nil
    case <-ctx.Done():
        return 0, ctx.Err()
    }
}
```
A buffered channel already is a thread-safe bounded queue by design — no extra Mutex needed since channels self-synchronize.

---

## Section E: Distributed Systems

**E1.** Idempotency: calling the same request multiple times produces the same result with no repeated side effects. Example: Payment API — the client sends an `Idempotency-Key` (UUID) header; the server stores this key with its result; if the same key arrives again, it returns the stored result instead of reprocessing (prevents double-charging when a client retries after a timeout).

**E2.** The outbox pattern solves the "dual write" problem: if you write to the DB and then publish a message separately, there's a gap between the two operations — if the process crashes in between (DB write succeeded but publish didn't happen), the message is silently lost. Outbox pattern: write business data + event into an "outbox" table within the **same DB transaction**, then a separate process (polling or CDC) reads the outbox table and publishes to the message queue, guaranteeing atomicity.

**E3.** A circuit breaker tracks the error rate of calls to a dependency; once the error rate crosses a threshold, it "opens" — all subsequent requests fail immediately (fail fast) without actually calling the dependency, for some cooldown period. It then tries a "half-open" probe to see if the dependency has recovered. Difference from retry: retry re-attempts a single failed request; a circuit breaker stops an entire stream of requests once it knows a dependency is unhealthy, preventing cascading failure.

**E4.** Saga pattern: break a distributed transaction into a sequence of local transactions across services, where each step has a "compensating transaction" to undo it if a later step fails. Used instead of 2PC (two-phase commit) because 2PC requires locking resources across multiple services simultaneously, which doesn't scale well in large distributed systems. Example: place order → deduct inventory → charge payment; if payment fails, run a compensating action to restore inventory.

**E5.** Exponential backoff: increase the wait time between retries progressively (e.g. 1s, 2s, 4s, 8s...) instead of retrying immediately every time. Needed because if all clients retry immediately while a downstream service is already overloaded, it makes the overload worse (thundering herd). Usually combined with jitter (small randomization) so multiple clients don't retry at exactly the same moment.

---

## Section F: System Design

**F1-F2.** These are open-ended questions with no single "correct" answer — interviewers evaluate how systematically you think. See the framework in `note.md` section 9. For F2, prepare one concrete story from your own project experience ahead of time — don't improvise on the spot.

**F3.** Typically shard by `user_id` or `tenant_id`/`firm_id` (depending on the domain) so that queries related to one user/tenant stay on the same shard, avoiding cross-shard joins. Trade-off: system-wide aggregate queries (cross-shard aggregation) become harder and slower; you need extra application-layer logic to route queries to the right shard; rebalancing when adding a new shard is complex.

**F4.** Cache-aside: the application checks the cache first, and on a miss reads the DB and populates the cache itself — cache and DB can be briefly out of sync (eventual consistency); this is the simplest and most common pattern. Write-through: every write goes to the cache and DB at the same time — the cache stays consistent with the DB, but writes are slower (must wait for both). Choose cache-aside when prioritizing simplicity/read performance; choose write-through when the cache must always be accurate immediately after a write.

---

## Section G: AWS & PostgreSQL

**G1.** A presigned URL is a temporary, signed URL (with an expiry) that lets a client upload/download an S3 object directly without AWS credentials and without the request passing through your backend — reducing server load and improving transfer speed.

**G2.** Task definition: describes the configuration of a single container (image, CPU/memory, ports, env vars...) — like a blueprint. Service: runs N instances of that task definition, automatically replaces failed tasks (health check failures), and can auto-scale the task count based on metrics (CPU, request count...) via Application Auto Scaling.

**G3.** `EXPLAIN ANALYZE` actually executes the query and shows: the execution plan (whether it uses an index or does a full table scan), actual time spent at each step, and actual row counts processed at each step — helping you find slow queries caused by missing indexes or a suboptimal plan. Use it when debugging a slow query or before deploying a new query to high-traffic production.

**G4.** Each PostgreSQL connection consumes significant server-side resources (memory, process); with thousands of clients connecting directly, the DB can easily become overloaded on connections. Connection pooling (pgbouncer) sits between the app and the DB, maintaining a small pool of real DB connections and multiplexing many virtual client connections through them — significantly reducing DB load under heavy traffic.

---

## Section H: Message Queue

**H1.** RabbitMQ: traditional broker, push-based delivery, strong routing flexibility (exchanges/topics), good for task queues and RPC-style messaging. Kafka: distributed log, consumers pull and track their own offset, extremely high throughput, supports replaying historical events — best for event streaming/log aggregation. SQS: AWS-managed queue, zero ops, built-in DLQ via redrive policy — best when you want a managed solution and are already on AWS. Pick based on need: replay history → Kafka; complex routing/RPC → RabbitMQ; no-ops managed queue → SQS.

**H2.** Ack (acknowledgement) tells the broker "I've successfully processed this message" — only after ack does the broker remove it from the queue. If the consumer crashes before acking, the message is redelivered (to the same or another consumer). This redelivery-on-failure behavior is exactly what "at-least-once delivery" means — a message may be processed more than once (e.g. if it was processed but the ack was lost), so consumers must be idempotent.

**H3.** (1) Multiple queues by priority level (e.g. `queue.high`, `queue.low`) — the consumer always drains the high queue before touching the low queue; simple, well-supported by most brokers. (2) A priority field on the message combined with an in-memory priority queue (heap) on the consumer side — messages get buffered and popped in priority order regardless of arrival order (some brokers, like RabbitMQ, also support native priority via `x-max-priority`).

**H4.** Transient errors (timeout, connection refused, temporary unavailability) are worth retrying because the same request will likely succeed shortly after. Permanent errors (malformed JSON, validation failure, missing required field) will never succeed no matter how many times you retry — these should be sent straight to the DLQ (or dropped with a logged reason) instead of wasting retry attempts and delaying other messages behind them.

**H5.** See the full runnable implementation in `notes/11-Message-Queue.md` and `code/priority-queue-processor/main.go`. Core idea: use `container/heap` to build a max-heap keyed on `Value`; a `mutex`-protected `Enqueue` pushes onto the heap and signals a worker via a buffered notify channel; the worker drains the heap fully (highest priority first) each time it wakes up. Test it by feeding a high-value transaction after a couple of low-value ones and confirming it gets processed first.

---

## Section I: MongoDB

**I1.** ESR = Equality, Sort, Range. In a compound index, put fields used for equality (`=`) filters first, fields used for sorting next, and fields used for range filters (`$gt`/`$lt`) last. This ordering lets MongoDB use the index to narrow down matches as much as possible before falling back to scanning a range, minimizing the documents it has to examine.

**I2.** A compound index can only be used efficiently by queries that filter on a **prefix** of its fields, left to right. `{a:1, b:1, c:1}` supports `a` alone, `a+b` together, or `a+b+c` together — because the index is physically sorted by `a` first, then `b` within each `a`, then `c` within each `b`. A query filtering only on `b` (skipping `a`) can't use this sort order efficiently, since matching `b` values are scattered across the whole index, not grouped together.

**I3.** `explain("executionStats")` shows whether the query used an index (`IXSCAN`) or scanned the whole collection (`COLLSCAN`), plus `totalDocsExamined` (how many documents MongoDB actually looked at) versus `nReturned` (how many it actually returned). A large gap between these two numbers means the index isn't narrowing things down well — the query is examining far more documents than it needs to, signaling a missing or poorly designed index.

**I4.** Use a TTL index when data has a natural, fixed expiry and you don't need exact-to-the-second deletion — sessions, OTP codes, short-lived cache documents. It avoids writing and maintaining a separate cleanup cron job. Limitation: MongoDB's background TTL sweep runs roughly every 60 seconds, so expired documents can briefly linger past their exact expiry time — not suitable when you need precise, immediate deletion.

**I5.** A good shard key has high cardinality (many distinct values) so data spreads evenly across shards, and matches the fields most queries filter on (e.g. `userId` or `tenantId`) so most queries hit a single shard. A bad shard key (e.g. a `status` field with only 2-3 possible values) creates a "hot shard" where most data and traffic pile onto one or two shards, and queries that don't include the shard key must scatter-gather across all shards — much slower than a single-shard query.

---

## Section J: Logging & Monitoring

**J1.** DEBUG → INFO → WARN → ERROR → FATAL. DEBUG: fine-grained detail, usually disabled in production. INFO: normal noteworthy events (request received, job completed). WARN: something unusual but the flow still continues (first retry attempt, falling back to a default). ERROR: a specific operation failed, needs attention, but the system keeps running. FATAL: an error severe enough that the program must stop.

**J2.** Structured (JSON) logs are queryable and filterable at scale — you can search "all logs where `user_id=123`" or aggregate "count of errors grouped by `reason` in the last hour" directly in a log platform. Free-form text logs require fragile regex parsing to extract the same information and don't scale well once you have many services producing large volumes of logs.

**J3.** ELK (Elasticsearch + Logstash + Kibana): full-text indexing in Elasticsearch, powerful search via Kibana, but relatively expensive to run at scale since everything is indexed. Grafana + Loki: Loki only indexes metadata (labels), not full log text, making it cheaper to operate; pairs naturally with Prometheus metrics in the same Grafana dashboards. CloudWatch Logs: fully managed, integrates natively with ECS/Lambda with zero extra infrastructure, but has more limited query/dashboard capability than ELK or Grafana.

**J4.** Logs are discrete event records (what happened, with detail) — good for debugging a specific incident. Metrics are numeric measurements over time (request rate, latency percentiles, error rate) — good for spotting trends and triggering alerts quickly since they're cheaper to query than logs. Distributed tracing follows a single request as it moves through multiple services (each service adds a "span" under a shared `trace_id`) — good for finding which specific service in a chain is causing latency. You need all three because each answers a different question: logs = what exactly happened, metrics = is something trending wrong, tracing = where in the chain is it happening.

---

## Section K: Code Quality

**K1.**
```go
func process(tx Transaction) error {
    if tx.Value <= 0 {
        return errors.New("invalid value")
    }
    if tx.ID == "" {
        return errors.New("missing id")
    }
    return nil
}
```
Guard clauses handle the failure cases first and return early, so the "main path" isn't nested inside multiple `if` blocks — easier to scan and extend.

**K2. S**ingle Responsibility: a `PaymentService` should only handle payment logic, not also format emails or write logs directly — inject a `Logger` and `EmailService` instead. **O**pen/Closed: add a new payment provider (e.g. PayPal alongside Stripe) by implementing the same `PaymentProvider` interface, without touching the code that calls it. **L**iskov Substitution: any type implementing `PaymentProvider` must be swappable for another without breaking the caller's expectations. **I**nterface Segregation: keep interfaces small and focused — a `PaymentProvider` interface should only have `Charge`, not also unrelated methods like `SendReceipt` that not every implementation needs. **D**ependency Inversion: `PaymentService` should depend on the `PaymentProvider` interface, not directly on a concrete `*StripeClient` — makes it mockable in tests and swappable later.

**K3.** Validate at the boundary where external data enters the system — HTTP request bodies, message queue payloads. Once data has passed that check, everything downstream (services, business logic) can trust it's well-formed and doesn't need to re-validate the same fields again — re-validating at every layer is redundant work and clutters the code with defensive checks for situations that can no longer occur.

**K4.** `%w` wraps the original error inside the new one while preserving a reference to it, so callers further up the stack can use `errors.Is(err, someSentinelErr)` or `errors.As(err, &someErrType)` to check for a specific underlying error even though it's now wrapped in a more descriptive message. Plain string concatenation (`"process transaction: " + err.Error()`) throws away the original error value — you can only inspect the resulting text, not programmatically check what kind of error it was.

**K5.** Use `panic`/`recover` only for truly unrecoverable situations — programmer bugs or a state the program genuinely cannot continue from (e.g. a required invariant is violated at startup). For anything expected to happen during normal operation (invalid input, a record not found, a network call failing), return a normal `error` — that's the idiomatic Go way to signal failure, and it keeps control flow explicit and easy to reason about.
