# Channels

[← Quay về note.md](../note.md)

## Giải thích
Channel là "ống dẫn" để goroutine giao tiếp và đồng bộ với nhau — theo triết lý Go: *"Don't communicate by sharing memory; share memory by communicating."*

Hai loại:
- **Unbuffered channel** (`make(chan int)`): gửi và nhận phải xảy ra đồng thời (synchronous rendezvous). Gửi sẽ block tới khi có người nhận.
- **Buffered channel** (`make(chan int, 3)`): gửi không block cho tới khi buffer đầy.

Quy tắc quan trọng:
- Chỉ **người gửi (sender)** nên close channel, không bao giờ để receiver close.
- Gửi vào channel đã đóng → **panic**.
- Đóng channel đã đóng → **panic**.
- Đọc từ channel đã đóng → nhận về zero value ngay lập tức, không block (dùng `v, ok := <-ch` để check `ok == false` nghĩa là channel đã đóng và rỗng).

## Ví dụ: Producer-Consumer cơ bản
```go
func producerConsumer() {
    ch := make(chan int, 5) // buffered
    done := make(chan struct{})

    // Producer
    go func() {
        for i := 0; i < 10; i++ {
            ch <- i
        }
        close(ch) // báo hiệu không còn dữ liệu
    }()

    // Consumer
    go func() {
        for v := range ch { // tự động thoát khi channel đóng và rỗng
            fmt.Println("received:", v)
        }
        close(done)
    }()

    <-done
}
```

## Ví dụ: `select` — chờ nhiều channel cùng lúc
```go
func selectExample(ch1, ch2 <-chan int, timeout <-chan time.Time) {
    for {
        select {
        case v := <-ch1:
            fmt.Println("from ch1:", v)
        case v := <-ch2:
            fmt.Println("from ch2:", v)
        case <-timeout:
            fmt.Println("timeout, exiting")
            return
        default:
            // optional: chạy khi không channel nào sẵn sàng (non-blocking)
        }
    }
}
```

## Pattern quan trọng: Fan-in / Fan-out
- **Fan-out**: nhiều goroutine cùng đọc từ 1 channel để xử lý song song (worker pool).
- **Fan-in**: gộp kết quả từ nhiều channel vào 1 channel duy nhất.

```go
// Fan-out: nhiều worker xử lý cùng 1 input channel
func fanOut(input <-chan int, workerCount int) <-chan int {
    output := make(chan int)
    var wg sync.WaitGroup

    for i := 0; i < workerCount; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            for v := range input {
                output <- v * v // xử lý (ví dụ: bình phương)
            }
        }()
    }

    go func() {
        wg.Wait()
        close(output) // đóng khi tất cả worker xong
    }()

    return output
}
```

## Checklist tự kiểm tra
- [ ] Giải thích được khác biệt buffered vs unbuffered channel
- [ ] Biết ai nên close channel và vì sao
- [ ] Nhớ được 3 hành vi: gửi vào channel đóng, đóng channel đã đóng, đọc từ channel đóng+rỗng
- [ ] Tự viết được `select` với nhiều case + timeout
- [ ] Giải thích được fan-in/fan-out và viết lại được code fan-out ở trên không nhìn mẫu
