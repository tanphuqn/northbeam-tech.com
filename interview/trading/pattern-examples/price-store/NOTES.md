**Pattern this teaches:** "stateful correctness" problems aren't always about
sliding windows (Problem 1 in the practice paper). A **price store with
point-in-time lookup** is a different flavor of the same category: you need
full history per key, kept sorted, and a **binary search** (`sort.Search`)
to answer "what was true as of time T" — not just "what's the latest value."

**Key edge case:** `PriceAsOf` at a timestamp *before* the first recorded
observation must return "not found," not the first price — a classic
off-by-one trap if you get the binary search boundary wrong (`i == 0` means
every entry is after `ts`).

**Compare to the rate limiter (Problem 1):** that one only needed a *rolling
window* (evict from the front, append at the back — a queue). This one needs
*permanent history with random point-in-time access* (a sorted slice +
binary search). Recognizing which of these two shapes a "store" problem
actually needs is the real skill being tested.
