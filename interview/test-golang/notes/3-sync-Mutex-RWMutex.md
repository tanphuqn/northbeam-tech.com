# sync.Mutex / sync.RWMutex

[← Quay về note.md](../note.md)

## Giải thích
Khi nhiều goroutine cùng đọc/ghi một biến chia sẻ (shared state) mà không dùng channel, bạn cần khóa (lock) để tránh race condition.

- `sync.Mutex`: chỉ 1 goroutine được vào critical section tại một thời điểm (dù đọc hay ghi).
- `sync.RWMutex`: cho phép **nhiều reader đồng thời**, nhưng **writer thì độc quyền** (không ai được đọc/ghi khi có writer). Dùng khi đọc nhiều hơn ghi rất nhiều (read-heavy workload).

Nguyên tắc: luôn `defer mu.Unlock()` ngay sau `mu.Lock()` để tránh quên unlock khi có panic/early return.

## Ví dụ: Counter an toàn với Mutex
```go
type SafeCounter struct {
    mu    sync.Mutex
    count map[string]int
}

func NewSafeCounter() *SafeCounter {
    return &SafeCounter{count: make(map[string]int)}
}

func (c *SafeCounter) Inc(key string) {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.count[key]++
}

func (c *SafeCounter) Value(key string) int {
    c.mu.Lock()
    defer c.mu.Unlock()
    return c.count[key]
}
```
Lưu ý: map trong Go **không an toàn** khi nhiều goroutine ghi đồng thời — sẽ panic "concurrent map writes". Đây là câu hỏi phỏng vấn rất phổ biến.

## Ví dụ: RWMutex cho cache read-heavy
```go
type Cache struct {
    mu   sync.RWMutex
    data map[string]string
}

func (c *Cache) Get(key string) (string, bool) {
    c.mu.RLock()         // nhiều goroutine có thể RLock cùng lúc
    defer c.mu.RUnlock()
    v, ok := c.data[key]
    return v, ok
}

func (c *Cache) Set(key, value string) {
    c.mu.Lock()          // độc quyền khi ghi
    defer c.mu.Unlock()
    c.data[key] = value
}
```

## Checklist tự kiểm tra
- [ ] Giải thích được khi nào chọn Mutex, khi nào chọn RWMutex
- [ ] Giải thích được vì sao map cần Mutex bảo vệ khi ghi đồng thời
- [ ] Tự viết lại `SafeCounter` và `Cache` ở trên không nhìn mẫu
- [ ] Nhớ luôn dùng `defer mu.Unlock()` ngay sau `mu.Lock()`
