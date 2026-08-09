# Context (context.Context)

[← Quay về note.md](../note.md)

## Giải thích
`context` dùng để:
1. **Truyền deadline/timeout** xuống các goroutine con để họ tự hủy công việc khi hết giờ.
2. **Cancel theo chuỗi** (cha hủy thì tất cả con cũng phải dừng) — rất quan trọng trong hệ thống xử lý request (vd: HTTP request bị client hủy giữa chừng).
3. Truyền giá trị request-scoped (như request ID, user ID) — dùng hạn chế, không dùng để truyền config/dependency.

Các hàm tạo context:
- `context.Background()`: gốc, dùng ở main/entrypoint.
- `context.WithCancel(parent)`: trả về context + hàm `cancel()` để hủy thủ công.
- `context.WithTimeout(parent, duration)`: tự động cancel sau khoảng thời gian.
- `context.WithDeadline(parent, time)`: tự động cancel tại một thời điểm cụ thể.
- `context.WithValue(parent, key, value)`: gắn giá trị.

Quy tắc quan trọng: **luôn gọi `cancel()`** (thường bằng `defer cancel()`) để giải phóng tài nguyên, kể cả khi context tự hết hạn.

## Ví dụ: Timeout cho một tác vụ chậm
```go
func doWorkWithTimeout() error {
    ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
    defer cancel() // luôn gọi, tránh leak

    resultCh := make(chan int)
    go func() {
        time.Sleep(3 * time.Second) // giả lập việc chậm (vd: gọi DB/API)
        resultCh <- 42
    }()

    select {
    case res := <-resultCh:
        fmt.Println("result:", res)
        return nil
    case <-ctx.Done():
        // ctx.Err() trả về context.DeadlineExceeded hoặc context.Canceled
        return fmt.Errorf("timed out: %w", ctx.Err())
    }
}
```

## Ví dụ: Cancel theo chuỗi (parent hủy → con dừng theo)
```go
func processRequest(ctx context.Context, jobs <-chan int) {
    for {
        select {
        case <-ctx.Done():
            fmt.Println("cancelled:", ctx.Err())
            return
        case job, ok := <-jobs:
            if !ok {
                return
            }
            fmt.Println("processing job:", job)
        }
    }
}
```

## Checklist tự kiểm tra
- [ ] Phân biệt được `WithCancel`, `WithTimeout`, `WithDeadline`
- [ ] Giải thích được vì sao luôn phải gọi `cancel()` kể cả khi dùng `WithTimeout`
- [ ] Giải thích được `context.WithValue` nên/không nên dùng cho việc gì
- [ ] Tự viết lại ví dụ timeout ở trên không nhìn mẫu
