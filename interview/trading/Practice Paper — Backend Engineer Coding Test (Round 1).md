# **Practice Paper — Backend Engineer Coding Test (Round 1\)**

**Purpose of this paper.** This is a **practice paper**, not the real test. It mirrors the real test's structure exactly — the same timing, rules, number and style of problems, submission format, and scoring convention. The problems themselves are different: **nothing in this paper appears in the real test.** We recommend attempting it once under real conditions (120 uninterrupted minutes, AI assistants disabled, Google allowed) before your scheduled test.

## **Structure of the real test (identical to this paper)**

**Total: 120 minutes · 120 points.** Each question's maximum score equals its allocated minutes.

| \# | Problem type | Language | Time / Points |
| :---- | :---- | :---- | :---- |
| 1 | Stateful correctness — implement a small component against an exact specification | Go or NodeJS/TypeScript — your choice | 35 |
| 2 | Concurrency — build a concurrent component with ordering, timeout and shutdown rules | Go only | 50 |
| 3 | Debug & harden — find and fix planted defects in code you didn't write (provided in Go and TypeScript; pick one) | Go or TypeScript — your choice | 35 |

## **Rules (identical to the real test)**

> * **Time limit: 120 minutes** from the moment you open the problems.  
> * **AI assistance is NOT allowed** (no Claude, ChatGPT, Copilot, Cursor, etc.). The real test session is monitored.  
> * **Google and documentation ARE allowed.** Look up any API or language detail you like. The problems are original — searching for solutions will only waste your time.  
> * Standard library only. No external packages except your language's built-in test tooling.

## **Submission format (identical to the real test)**

Work in your own coding environment (AI plugins disabled), then copy the final content of each file into the answer box of the corresponding problem. File names are entirely your choice — label every file with a header line:

`=== FILE: ratelimiter.go ===`  
`(file content)`

`=== NOTES ===`  
`(your notes)`

For each problem include: your implementation file(s); your test file(s) plus one line on how to run them (go test ./... / node \--test); and a NOTES section — approach, trade-offs, edge cases handled, what you'd do with more time. A few sentences is enough. **We read NOTES carefully** — they carry real points.

## **How grading works (identical convention)**

Within each problem: correctness and edge cases first, then your own tests and verification, then code clarity, organization and NOTES. Specifications are exact — expect edge cases to be deliberately buried in the precise wording. Read every rule twice.

# ---

**Practice Problem 1 — Sliding-Window Rate Limiter (\~35 min, Go or NodeJS/TypeScript)**

You are building an in-memory rate limiter used in front of a trading API. Clients send requests; each request carries a caller-supplied timestamp. The limiter must enforce the window rules exactly as specified.

### **Specification**

Implement a RateLimiter. Method names/shapes may follow your language's conventions; the *rules* are exact.

### **allow(clientID, now) → accepted (boolean)**

> * clientID: non-empty string. now: caller-supplied timestamp in milliseconds (do **not** read the system clock inside the limiter — the caller owns time). You may assume now values for a given client are non-decreasing.  
> * The request is **accepted** if the client has **fewer than 5** previously accepted requests inside the rolling window: an accepted request with timestamp t is inside the window while now − t \< 10000 (a request aged **exactly** 10000 ms is outside).  
> * If accepted, record the request's timestamp. If rejected, record **nothing**.  
> * Invalid arguments (empty clientID, negative now) must be rejected without changing any state.

### **The rule people miss**

**Rejected requests are invisible.** They never count toward the window — only *accepted* requests do. A client hammering the API with rejected requests must regain capacity the moment its old *accepted* requests age out.

### **activeClients(now) → integer**

> * Number of clients with at least one accepted request currently inside the window.

### **Requirements**

> * Single-threaded correctness is enough — no locking required. (Concurrency is Problem 2.)  
> * Standard library only.  
> * Include your own tests demonstrating the rules — including the ones you think we buried — plus one line on how to run them.  
> * Submit under \=== FILE: ... \=== headers plus a \=== NOTES \=== section: data-structure choice and why, and every edge case you covered.

# ---

**Practice Problem 2 — Ordered Results Pipeline (\~50 min, Go only)**

A market-data service receives a stream of messages with strictly increasing IDs. Each message must be enriched by a slow external function. Enrichment should run in parallel — but downstream consumers require results in **strict input order**. You are building the pipeline.

### **Specification**

`type Msg struct {`  
    `ID  int    // strictly increasing on input`  
    `Raw string`  
`}`

`type Out struct {`  
    `Msg      Msg`  
    `Enriched string`  
    `Err      error   // nil on success; TimeoutErr on timeout; enricher error otherwise`  
`}`

`// enricher may be slow, may return an error.`  
`// It is NOT guaranteed to respect ctx quickly.`  
`type Enricher func(ctx context.Context, m Msg) (string, error)`

`func NewPipeline(workers int, perMsgTimeout time.Duration, e Enricher) *Pipeline`

`// blocks if internal capacity is full (bounded memory)`  
`func (p *Pipeline) Submit(m Msg)`

`// results in strict ascending Msg.ID order`  
`func (p *Pipeline) Out() <-chan Out`

`func (p *Pipeline) Drain(ctx context.Context) (dropped int)`

### **Rules**

> 1. **Strict output order.** Results emerge from Out() in ascending Msg.ID order — even though enrichments complete out of order.  
> 2. **Parallel enrichment.** Up to workers messages may be inside the enricher simultaneously.  
> 3. **Timeout.** If enrichment of a message exceeds perMsgTimeout, emit Out{Msg, "", TimeoutErr} in its correct ordered position and move on. A timed-out straggler must not stall the output stream or corrupt ordering.  
> 4. **Errors don't stop the stream.** An enricher error is reported in that message's Out; later messages still flow.  
> 5. **Drain.** Drain(ctx): stop accepting new messages, let in-flight enrichments finish (up to ctx's deadline), emit everything already processed in order, report messages accepted-but-never-emitted as dropped (return the count), then close Out(). Document your choice for Submit-after-Drain.  
> 6. **Bounded memory.** The reorder buffer and queues must be bounded; Submit blocks when full. Pick a reasonable bound and document it.

### **Deliverables**

> * The implementation, standard library only.  
> * Tests that **prove** rules 1, 3 and 5 — ordering with an out-of-order-completing enricher, timeout with a deliberately slow enricher, and drain with messages in flight. Tests must pass with go test \-race.  
> * \=== NOTES \===: your design in a paragraph (how you reorder, how you bound memory), and what you'd harden with more time.

### **Hint that is not a trick**

The hard part is rule 1 and rule 3 interacting: a timed-out message still occupies a position in the output order. Prove to yourself that a straggler can neither stall the stream nor emit into the wrong position — or explain in NOTES exactly what guarantee you chose.

# ---

**Practice Problem 3 — Debug & Harden (\~35 min, Go or TypeScript — pick ONE)**

You have inherited a small service from a developer who has left the company. It tracks traded volume per symbol and serves a leaderboard. Two versions are provided — Go and TypeScript. **Choose one and work only on that version.** The same defects exist in both.

### **The specification the service is supposed to meet**

> * addTrade(symbol, qty) records a trade (qty is a positive decimal; fractional quantities are normal).  
> * avgSize(symbol) returns the **exact** average trade size (total quantity ÷ trade count), fractional part included.  
> * topN(n) returns the n symbols with highest total volume, **ties broken alphabetically** — the same input must always produce the same output.  
> * startSnapshots(everyMs) begins periodic snapshot logging; stop() must stop it: **after stop() returns, no further snapshots appear and no background work remains running.**

### **What production has reported**

> 1. In long-running processes, background activity keeps growing, and snapshot log lines keep appearing **after** stop() was called.  
> 2. Average trade sizes are consistently a little **lower** than the reference calculation, especially for symbols with many small fractional trades.  
> 3. When two symbols have equal volume, the leaderboard order **flips between runs**, breaking a downstream diff-based report.

### **Your tasks**

> 1. **Find the defects.** We believe there is more than one. Reproduce them if you can.  
> 2. **Fix them** with minimal, surgical changes — this is a live service, not a rewrite.  
> 3. **Explain each root cause** in \=== NOTES \===: what was wrong, why it produced the reported symptom, how your fix works.  
> 4. **Add a regression test for each defect** — the test that would have caught it before production did. (Go: tests must pass with go test \-race.)

### **GO VERSION (leaderboard.go)**

`// Package board — per-symbol volume tracker with a leaderboard.`  
`package board`

`import (`  
	`"fmt"`  
	`"sort"`  
	`"time"`  
`)`

`type stats struct {`  
	`totalQty int64 // accumulated quantity`  
	`count    int64`  
`}`

`type Board struct {`  
	`stats  map[string]*stats`  
	`ticker *time.Ticker`  
`}`

`func NewBoard() *Board {`  
	`return &Board{stats: make(map[string]*stats)}`  
`}`

`// AddTrade records a trade. qty is a positive decimal quantity.`  
`func (b *Board) AddTrade(symbol string, qty float64) {`  
	`s, ok := b.stats[symbol]`  
	`if !ok {`  
		`s = &stats{}`  
		`b.stats[symbol] = s`  
	`}`  
	`s.totalQty += int64(qty)`  
	`s.count++`  
`}`

`// AvgSize returns the exact average trade size for the symbol.`  
`func (b *Board) AvgSize(symbol string) float64 {`  
	`s, ok := b.stats[symbol]`  
	`if !ok || s.count == 0 {`  
		`return 0`  
	`}`  
	`return float64(s.totalQty / s.count)`  
`}`

`// TopN returns the n symbols with highest total volume, ties alphabetical.`  
`func (b *Board) TopN(n int) []string {`  
	`symbols := make([]string, 0, len(b.stats))`  
	`for sym := range b.stats {`  
		`symbols = append(symbols, sym)`  
	`}`  
	`sort.Slice(symbols, func(i, j int) bool {`  
		`return b.stats[symbols[i]].totalQty > b.stats[symbols[j]].totalQty`  
	`})`  
	`if n > len(symbols) {`  
		`n = len(symbols)`  
	`}`  
	`return symbols[:n]`  
`}`

`// StartSnapshots logs the top 3 periodically until Stop is called.`  
`func (b *Board) StartSnapshots(every time.Duration) {`  
	`b.ticker = time.NewTicker(every)`  
	`go func() {`  
		`for range b.ticker.C {`  
			`fmt.Println("snapshot:", b.TopN(3))`  
		`}`  
	`}()`  
`}`

`// Stop ends snapshot logging.`  
`func (b *Board) Stop() {`  
	`b.ticker = nil`  
`}`

### **TYPESCRIPT VERSION (leaderboard.ts)**

`// Per-symbol volume tracker with a leaderboard.`

`type Stats = { totalQty: number; count: number };`

`export class Board {`  
  `private stats = new Map<string, Stats>();`  
  `private timer: ReturnType<typeof setInterval> | null = null;`

  `// Records a trade. qty is a positive decimal quantity.`  
  `addTrade(symbol: string, qty: number): void {`  
    `const s = this.stats.get(symbol) ?? { totalQty: 0, count: 0 };`  
    `s.totalQty += Math.floor(qty);`  
    `s.count += 1;`  
    `this.stats.set(symbol, s);`  
  `}`

  `// Returns the exact average trade size for the symbol.`  
  `avgSize(symbol: string): number {`  
    `const s = this.stats.get(symbol);`  
    `if (!s || s.count === 0) return 0;`  
    `return Math.floor(s.totalQty / s.count);`  
  `}`

  `// Returns the n symbols with highest total volume, ties alphabetical.`  
  `topN(n: number): string[] {`  
    `const symbols = [...this.stats.keys()];`  
    `symbols.sort(`  
      `(a, b) => this.stats.get(b)!.totalQty - this.stats.get(a)!.totalQty`  
    `);`  
    `return symbols.slice(0, n);`  
  `}`

  `// Logs the top 3 periodically until stop() is called.`  
  `startSnapshots(everyMs: number): void {`  
    `this.timer = setInterval(() => {`  
      `console.log("snapshot:", this.topN(3));`  
    `}, everyMs);`  
  `}`

  `// Ends snapshot logging.`  
  `stop(): void {`  
    `this.timer = null;`  
  `}`  
`}`

## ---

**Self-check: what a strong submission looks like**

**Problem 1** — your tests hit the exact window boundary (a request exactly 10000 ms old), prove rejected requests never consume capacity, show per-client independence, and show invalid arguments changing nothing. NOTES names your data structure and why eviction is cheap.  
**Problem 2** — your ordering test uses an enricher that deliberately completes out of order and asserts strictly ascending IDs on output; your timeout test uses an enricher that ignores its context; your drain test has messages in flight and asserts the dropped count and channel closure. Everything passes go test \-race. NOTES states the reorder-buffer bound and the straggler guarantee you chose.  
**Problem 3** — you found more than one defect (hint: read the spec's exact words — "fractional", "always the same output", "no background work remains"), each fix is small, each root cause maps to one reported symptom, and each fix has a test that fails on the old code. A missed defect honestly flagged in NOTES beats a rewrite.  
*No solutions are provided for this paper — the real test provides none either. If you can convince yourself with your own tests, you are ready.*