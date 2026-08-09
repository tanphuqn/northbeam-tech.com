# Goroutines

[← Quay về note.md](../note.md)

## Giải thích cơ bản
Goroutine là một "lightweight thread" được Go runtime quản lý (không phải OS thread trực tiếp). Chi phí khởi tạo cực nhỏ (~2KB stack, tự grow), nên bạn có thể tạo hàng trăm nghìn goroutine mà không sập máy — khác hẳn thread OS truyền thống (thường ~1-8MB/thread).

Điểm quan trọng cần hiểu:
- `go func(){...}()` chỉ **lên lịch chạy**, không đảm bảo chạy ngay hay chạy trước goroutine khác.
- Nếu `main()` kết thúc, mọi goroutine đang chạy dở sẽ bị "giết" ngay lập tức — không có cleanup, không defer nào chạy.
- Goroutine leak: nếu một goroutine bị block mãi mãi (chờ channel không ai gửi/nhận), nó sẽ tồn tại tới khi chương trình kết thúc → memory leak.

## Ví dụ: chờ nhiều goroutine hoàn thành bằng WaitGroup
```go
package main

import (
    "fmt"
    "sync"
)

func main() {
    var wg sync.WaitGroup

    for i := 0; i < 5; i++ {
        wg.Add(1)
        go func(id int) {
            defer wg.Done()
            fmt.Printf("goroutine %d done\n", id)
        }(i) // truyền i vào để tránh bug capture biến vòng lặp
    }

    wg.Wait() // chờ tất cả goroutine hoàn thành
    fmt.Println("all done")
}
```
Chạy nhiều lần sẽ thấy thứ tự in ra khác nhau mỗi lần — vì các goroutine chạy song song, không có thứ tự đảm bảo, chỉ có `wg.Wait()` đảm bảo tất cả đã xong trước khi in "all done".

---

## Goroutine leak — lỗi thường gặp nhất

### Vấn đề cốt lõi
Channel không buffer (`make(chan int)`) giống như một ống hẹp: người gửi và người nhận phải "bắt tay" cùng lúc. Nếu không có ai ở đầu kia nhận, người gửi phải **đứng chờ mãi mãi**.

```go
func leaky() {
    ch := make(chan int)     // (1) tạo channel, dung lượng = 0

    go func() {
        ch <- 1              // (3) goroutine con cố gửi giá trị 1 vào ch
    }()                      // (2) khởi động goroutine con, cha chạy tiếp ngay, KHÔNG chờ

    // (4) hàm leaky() không có dòng nào đọc từ ch → return luôn
}
```

Chạy từng bước:
- Dòng (2): `go func(){...}()` chỉ là "nhờ" Go chạy đoạn code ở một luồng riêng — goroutine cha không đợi nó xong.
- Dòng (3): `ch <- 1` nghĩa là "tôi muốn gửi số 1 vào channel". Vì `ch` không buffer, câu lệnh này **block** cho tới khi có ai chạy `<-ch` để nhận.
- Dòng (4): `leaky()` không hề có `<-ch` ở đâu cả → không ai bao giờ đến nhận → goroutine con đứng chờ **vĩnh viễn**.

### Hậu quả
Khi `leaky()` return, goroutine cha tiếp tục làm việc khác, nhưng goroutine con vẫn còn sống, đứng im chờ mãi trong bộ nhớ — không tự biến mất, không ai dọn. Gọi `leaky()` 1000 lần → 1000 goroutine "ma" kẹt mãi mãi, ngốn bộ nhớ dần → đây gọi là **goroutine leak**.

---

## 3 cách fix + demo runnable

Đã có sẵn code demo đầy đủ tại [`code/goroutine-leak/main.go`](../code/goroutine-leak/main.go) — chạy bằng:
```bash
cd test/code/goroutine-leak
go run main.go
```

Kết quả mẫu:
```
start                             goroutines = 1
after 5x leaky()                  goroutines = 6   <- tăng 5, KHÔNG BAO GIỜ giảm (leak thật)
after 5x fixedBuffered()          goroutines = 6   <- không tăng thêm, goroutine tự thoát xong
after 5x fixedContext() + cancel  goroutines = 6   <- không tăng thêm, cancel() giúp thoát
after 5x fixedReceiver()          goroutines = 6   <- không tăng thêm, có người nhận nên thoát
```

Điểm mấu chốt: sau bước `leaky()`, con số **6 giữ nguyên mãi mãi** ở các bước sau — đó là 5 goroutine bị kẹt vĩnh viễn. Trong khi 3 hàm fix, dù gọi thêm 5 lần mỗi hàm, số goroutine sống không tăng thêm vì chúng hoàn thành và biến mất ngay.

### Fix 1: Buffered channel
```go
func fixedBuffered() {
    ch := make(chan int, 1) // capacity 1, gửi không cần chờ người nhận
    go func() {
        ch <- 1 // không block, goroutine thoát ngay
    }()
}
```
Lý do: channel có "chỗ chứa" tạm (capacity 1), nên gửi xong là xong ngay, không cần ai nhận ngay lúc đó.

### Fix 2: Context cancellation
```go
func fixedContext(ctx context.Context) {
    ch := make(chan int)
    go func() {
        select {
        case ch <- 1:
            // gửi thành công
        case <-ctx.Done():
            // caller đã cancel, dừng chờ
            return
        }
    }()
}
```
Lý do: dùng `select` để vừa chờ gửi, vừa chờ tín hiệu hủy — nếu không ai nhận và context bị cancel, goroutine tự thoát thay vì chờ mãi.

### Fix 3: Đảm bảo luôn có người nhận khớp
```go
func fixedReceiver() {
    ch := make(chan int)
    go func() {
        ch <- 1
    }()
    <-ch // nhận giá trị, sender được giải phóng, cả hai bên đều hoàn thành
}
```
Lý do: đơn giản nhất — miễn là có `<-ch` khớp với `ch <- 1` ở đâu đó, không có ai bị kẹt.

---

## Checklist tự kiểm tra
- [ ] Giải thích được vì sao unbuffered channel gây block khi không có người nhận
- [ ] Tự viết được ví dụ goroutine leak từ đầu (không nhìn code mẫu)
- [ ] Giải thích được cả 3 cách fix và khi nào nên dùng cách nào
- [ ] Chạy thử `code/goroutine-leak/main.go`, quan sát `runtime.NumGoroutine()` tăng/không tăng để tự chứng minh
