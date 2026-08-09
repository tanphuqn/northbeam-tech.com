# Race Condition

[← Quay về note.md](../note.md)

## Giải thích
Race condition xảy ra khi 2+ goroutine truy cập cùng 1 biến, ít nhất 1 bên ghi, mà không có đồng bộ hóa (mutex/channel). Kết quả chương trình phụ thuộc vào thứ tự thực thi ngẫu nhiên của scheduler → bug khó tái hiện, "hôm nay chạy đúng mai chạy sai".

**Công cụ bắt buộc phải biết**: chạy `go run -race main.go` hoặc `go test -race ./...` để Go runtime tự phát hiện race (dùng ThreadSanitizer).

## Ví dụ bug race + fix

```go
// BUG: race condition trên biến "counter"
func raceyCounter() {
    counter := 0
    var wg sync.WaitGroup
    for i := 0; i < 1000; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            counter++ // đọc-sửa-ghi không atomic → race
        }()
    }
    wg.Wait()
    fmt.Println(counter) // kết quả không ổn định, thường < 1000
}

// FIX 1: dùng Mutex
func safeCounterMutex() {
    var mu sync.Mutex
    counter := 0
    var wg sync.WaitGroup
    for i := 0; i < 1000; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            mu.Lock()
            counter++
            mu.Unlock()
        }()
    }
    wg.Wait()
    fmt.Println(counter) // luôn = 1000
}

// FIX 2: dùng atomic (nhanh hơn Mutex cho phép toán đơn giản)
func safeCounterAtomic() {
    var counter int64
    var wg sync.WaitGroup
    for i := 0; i < 1000; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            atomic.AddInt64(&counter, 1)
        }()
    }
    wg.Wait()
    fmt.Println(atomic.LoadInt64(&counter)) // luôn = 1000
}
```

## Lỗi kinh điển: capture biến vòng lặp (loop variable capture)
Trước Go 1.22, biến loop được **tái sử dụng** qua mỗi vòng lặp — nếu goroutine tham chiếu trực tiếp biến đó mà không truyền qua tham số, tất cả goroutine có thể thấy cùng 1 giá trị (giá trị cuối cùng của loop).

```go
// BUG (trước Go 1.22): tất cả goroutine có thể in ra "3" thay vì 0,1,2
for i := 0; i < 3; i++ {
    go func() {
        fmt.Println(i) // "i" bị capture theo reference, dùng chung cho mọi goroutine
    }()
}

// FIX: truyền i vào làm tham số → mỗi goroutine có bản copy riêng
for i := 0; i < 3; i++ {
    go func(i int) {
        fmt.Println(i)
    }(i)
}
```
Từ **Go 1.22 trở lên**, mỗi vòng lặp `for` tạo một biến mới → bug này không còn xảy ra mặc định. Nhưng phỏng vấn viên rất hay hỏi câu này để test hiểu biết nền tảng, nên vẫn cần biết rõ.

## Checklist tự kiểm tra
- [ ] Giải thích được race condition là gì và vì sao nó khó tái hiện
- [ ] Biết chạy `go run -race` / `go test -race ./...` và đọc hiểu output
- [ ] Tự viết lại `raceyCounter` và cả 2 cách fix (Mutex, atomic) không nhìn mẫu
- [ ] Giải thích được bug loop variable capture trước Go 1.22 và cách fix
