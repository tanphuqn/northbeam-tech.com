# Practice Questions — Senior Backend Engineer (Go)

[← Quay về note.md](../note.md)

Read the question, answer/code it yourself before checking `answers.md`. Time yourself to simulate live-coding conditions.

---

## Section A: Goroutines & Channels

**A1.** How is a goroutine different from an OS thread? Why can Go spawn hundreds of thousands of goroutines without crashing?

**A2.** What happens if `main()` returns while other goroutines are still running?

**A3.** What is a "goroutine leak"? Give an example of leaking code.

**A4.** What's the difference between a buffered and an unbuffered channel? When does a send block?

**A5.** What happens when:
   a. You send on a closed channel?
   b. You close an already-closed channel?
   c. You receive from a closed, empty channel?

**A6.** Who should close a channel — the sender or the receiver? Why?

**A7.** Explain the fan-in / fan-out pattern. Give a real-world scenario where you'd use it.

---

## Section B: sync package

**B1.** When do you use `sync.Mutex` vs `sync.RWMutex`? Give a concrete scenario for each.

**B2.** Why is a Go map unsafe for concurrent writes? What actually happens when you run that code?

**B3.** Find the bug in this code and fix it:
```go
func wrong() {
    var wg sync.WaitGroup
    for i := 0; i < 3; i++ {
        go func() {
            wg.Add(1)
            defer wg.Done()
            fmt.Println(i)
        }()
    }
    wg.Wait()
}
```

**B4.** What is `sync.Once` used for? Give a real-world example.

**B5.** Compare `sync.Mutex` vs `atomic.AddInt64` for a counter — which is faster, and what's the trade-off?

---

## Section C: Context & Race Conditions

**C1.** Why use `context` instead of hand-rolling a timer + channel for cancellation?

**C2.** What's the difference between `context.WithCancel`, `context.WithTimeout`, and `context.WithDeadline`?

**C3.** Why must you always call `cancel()` (even with `WithTimeout`)? What happens if you forget?

**C4.** What should `context.WithValue` be used for, and what should it NOT be used for? Why?

**C5.** What is a race condition? What tool does Go provide to detect it, and what command do you run?

**C6.** Explain the "loop variable capture" bug in Go before version 1.22. Give a buggy example and the fix.

**C7.** Does the following code have a race condition? If so, point out where and fix it:
```go
func count() int {
    counter := 0
    var wg sync.WaitGroup
    for i := 0; i < 100; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            counter++
        }()
    }
    wg.Wait()
    return counter
}
```

---

## Section D: Live-coding — code these yourself in 20-30 minutes each

**D1.** Write a **worker pool** that processes a list of jobs with at most K workers running concurrently. It must support cancellation via `context.Context` — when the context is cancelled, workers must stop instead of processing remaining jobs.

**D2.** Write a **rate limiter** (token bucket style) that allows at most N requests/second. Expose a method `Allow() bool`.

**D3.** Write a thread-safe **TTL cache** that automatically evicts expired entries via a periodic background goroutine.

**D4.** Write a function that takes multiple input channels (`...<-chan int`) and merges (fan-in) them into a single output channel.

**D5.** Write a **bounded queue** supporting multiple concurrent producers and consumers, safe from race conditions.

---

## Section E: Distributed Systems

**E1.** What is idempotency? Give an example of an API that needs it and how you'd implement it (e.g. idempotency key).

**E2.** What problem does the outbox pattern solve? Why can't you just "write to DB then publish a message"?

**E3.** How does a circuit breaker work? How is it different from retry?

**E4.** Explain the Saga pattern. When would you use it instead of a traditional distributed transaction (2PC)?

**E5.** What is exponential backoff and why is it needed when retrying calls to another service?

---

## Section F: System Design

**F1.** Design a system that processes millions of payment transactions per day for a food-delivery platform. Walk through: requirements clarification, high-level architecture, database choice, caching strategy, scaling approach, and failure handling.

**F2.** Describe a real project you worked on involving high-load/distributed systems: what was the bottleneck, how did you scale/optimize it, and what trade-off did you choose?

**F3.** What criteria would you use to shard PostgreSQL for a large-scale transaction system? What's the trade-off?

**F4.** How do cache-aside and write-through differ? When would you choose each?

---

## Section G: AWS & PostgreSQL

**G1.** Explain what a presigned URL in S3 is used for.

**G2.** At a conceptual level, how do ECS task definitions and service scaling work?

**G3.** What does `EXPLAIN ANALYZE` in PostgreSQL tell you? When do you use it?

**G4.** What problem does connection pooling (e.g. pgbouncer) solve under heavy traffic?

---

## Section H: Message Queue

**H1.** Compare RabbitMQ, Kafka, and Amazon SQS. When would you pick each one?

**H2.** What is message acknowledgement (ack), and how does it relate to at-least-once delivery?

**H3.** Describe two different ways to implement a priority queue on top of a message broker.

**H4.** How should a consumer distinguish between a transient error (worth retrying) and a permanent error (should go straight to DLQ)?

**H5.** Live-coding: implement a transaction processor that consumes messages `{id, value, timestamp}` and processes high-value transactions (value > 1000) before lower-value ones, even if they arrive later. Don't use an external priority-queue library — implement the ordering logic yourself.

---

## Section I: MongoDB

**I1.** What is the ESR rule for compound indexes, and why does field order matter?

**I2.** Why does a compound index `{a:1, b:1, c:1}` serve a query filtering only on `a`, or on `a+b`, but not one filtering only on `b`?

**I3.** What does `explain("executionStats")` tell you? What's the difference between `totalDocsExamined` and `nReturned`, and why does a big gap between them matter?

**I4.** When would you use a TTL index instead of a manual cleanup cron job? What's a limitation of TTL indexes?

**I5.** What makes a good shard key? What goes wrong with a badly chosen one?

---

## Section J: Logging & Monitoring

**J1.** List the standard log levels in order and explain when to use each.

**J2.** Why is structured (JSON) logging preferred over free-form text logging at scale?

**J3.** Compare ELK, Grafana+Loki, and CloudWatch Logs for centralized logging.

**J4.** What's the difference between logs, metrics, and distributed tracing? Why do you need all three?

---

## Section K: Code Quality

**K1.** Refactor this function to remove unnecessary nesting using early returns:
```go
func process(tx Transaction) error {
    if tx.Value > 0 {
        if tx.ID != "" {
            return nil
        } else {
            return errors.New("missing id")
        }
    } else {
        return errors.New("invalid value")
    }
}
```

**K2.** Explain each letter of SOLID with a concrete Go example (not just the textbook definition).

**K3.** Where should input validation happen in a service, and why don't you need to re-validate at every layer afterward?

**K4.** What does `%w` do in `fmt.Errorf("process transaction %s: %w", id, err)`, and why is it better than string-concatenating the error?

**K5.** When is it appropriate to use `panic`/`recover` in Go, versus returning a normal error?
