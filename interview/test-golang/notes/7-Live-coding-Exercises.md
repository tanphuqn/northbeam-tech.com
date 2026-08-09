# Bài tập luyện live-coding (tự làm trước khi xem gợi ý)

[← Quay về note.md](../note.md)

## Bài 1: Worker Pool có giới hạn concurrency + hỗ trợ cancel
Yêu cầu: xử lý N job với tối đa K worker chạy song song, hỗ trợ hủy giữa chừng qua context.

```go
func WorkerPool(ctx context.Context, jobs []int, workerCount int, process func(int) int) []int {
    jobCh := make(chan int)
    resultCh := make(chan int)
    var wg sync.WaitGroup

    // gửi job
    go func() {
        defer close(jobCh)
        for _, j := range jobs {
            select {
            case jobCh <- j:
            case <-ctx.Done():
                return
            }
        }
    }()

    // worker
    for i := 0; i < workerCount; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            for {
                select {
                case j, ok := <-jobCh:
                    if !ok {
                        return
                    }
                    select {
                    case resultCh <- process(j):
                    case <-ctx.Done():
                        return
                    }
                case <-ctx.Done():
                    return
                }
            }
        }()
    }

    go func() {
        wg.Wait()
        close(resultCh)
    }()

    var results []int
    for r := range resultCh {
        results = append(results, r)
    }
    return results
}
```

## Bài 2: Rate Limiter (token bucket) dùng ticker
```go
type RateLimiter struct {
    tokens chan struct{}
}

func NewRateLimiter(rps int) *RateLimiter {
    rl := &RateLimiter{tokens: make(chan struct{}, rps)}
    for i := 0; i < rps; i++ {
        rl.tokens <- struct{}{} // đổ đầy bucket lúc khởi tạo
    }

    go func() {
        ticker := time.NewTicker(time.Second / time.Duration(rps))
        defer ticker.Stop()
        for range ticker.C {
            select {
            case rl.tokens <- struct{}{}: // thêm token mới, nếu bucket chưa đầy
            default: // bucket đầy, bỏ qua
            }
        }
    }()
    return rl
}

func (rl *RateLimiter) Allow() bool {
    select {
    case <-rl.tokens:
        return true
    default:
        return false
    }
}
```

## Bài 3: Cache có TTL, tự dọn dẹp bằng background goroutine
```go
type Item struct {
    value      string
    expiresAt  time.Time
}

type TTLCache struct {
    mu    sync.RWMutex
    items map[string]Item
}

func NewTTLCache(cleanupInterval time.Duration) *TTLCache {
    c := &TTLCache{items: make(map[string]Item)}
    go c.cleanupLoop(cleanupInterval)
    return c
}

func (c *TTLCache) Set(key, value string, ttl time.Duration) {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.items[key] = Item{value: value, expiresAt: time.Now().Add(ttl)}
}

func (c *TTLCache) Get(key string) (string, bool) {
    c.mu.RLock()
    defer c.mu.RUnlock()
    item, ok := c.items[key]
    if !ok || time.Now().After(item.expiresAt) {
        return "", false
    }
    return item.value, true
}

func (c *TTLCache) cleanupLoop(interval time.Duration) {
    ticker := time.NewTicker(interval)
    defer ticker.Stop()
    for range ticker.C {
        now := time.Now()
        c.mu.Lock()
        for k, v := range c.items {
            if now.After(v.expiresAt) {
                delete(c.items, k)
            }
        }
        c.mu.Unlock()
    }
}
```

## Checklist tự kiểm tra
- [ ] Tự code lại cả 3 bài trong 20-30 phút mỗi bài, không xem đáp án
- [ ] Với Bài 1: thử test case cancel giữa chừng, kiểm tra worker dừng đúng
- [ ] Với Bài 2: thử gọi `Allow()` liên tục và quan sát hành vi giới hạn rate
- [ ] Với Bài 3: thử set TTL ngắn và quan sát cleanup loop tự xóa entry hết hạn
