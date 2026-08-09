# sync.WaitGroup

[← Quay về note.md](../note.md)

## Giải thích
Dùng để chờ một nhóm goroutine hoàn thành trước khi tiếp tục — giống như "đếm ngược".

Quy tắc:
- `Add(n)` phải được gọi **trước khi** goroutine bắt đầu (thường ở goroutine cha, không phải bên trong goroutine con) — tránh race giữa `Add` và `Wait`.
- `Done()` nên dùng `defer` ngay đầu goroutine để đảm bảo luôn được gọi kể cả khi panic.
- `Wait()` block tới khi counter về 0.

## Lỗi thường gặp
```go
// SAI: Add() gọi bên trong goroutine → có thể Wait() chạy trước khi Add() kịp gọi
func wrong() {
    var wg sync.WaitGroup
    for i := 0; i < 3; i++ {
        go func() {
            wg.Add(1) // race! Wait() có thể đã thấy counter = 0 và return sớm
            defer wg.Done()
            // ...
        }()
    }
    wg.Wait()
}

// ĐÚNG
func correct() {
    var wg sync.WaitGroup
    for i := 0; i < 3; i++ {
        wg.Add(1) // gọi ở goroutine cha, trước khi go func()
        go func() {
            defer wg.Done()
            // ...
        }()
    }
    wg.Wait()
}
```

## Checklist tự kiểm tra
- [ ] Giải thích được vì sao `Add()` phải gọi ở goroutine cha, trước `go func()`
- [ ] Giải thích được vì sao `Done()` nên đặt bằng `defer` ngay đầu goroutine
- [ ] Tự tìm ra lỗi trong đoạn code "SAI" ở trên mà không đọc chú thích trước
