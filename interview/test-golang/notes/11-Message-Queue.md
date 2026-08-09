# Message Queue: RabbitMQ / Kafka / SQS, Producer–Consumer, Priority Queue, Retry, DLQ

[← Quay về note.md](../note.md)

## So sánh nhanh 3 hệ thống phổ biến

| | RabbitMQ | Kafka | Amazon SQS |
|---|---|---|---|
| Mô hình | Message broker truyền thống, push tới consumer | Distributed log, consumer tự kéo (pull) và giữ offset | Managed queue của AWS, fully managed |
| Ordering | Đảm bảo thứ tự trong 1 queue | Đảm bảo thứ tự trong 1 partition | FIFO queue đảm bảo thứ tự, Standard queue thì không |
| Throughput | Trung bình-cao | Rất cao (thiết kế cho streaming, hàng triệu event/s) | Cao, tự scale, không cần quản lý infra |
| Dùng khi | Task queue, RPC, cần routing linh hoạt (exchange/topic) | Event streaming, log aggregation, cần replay lại lịch sử event | Muốn giải pháp managed, không muốn tự vận hành broker |
| Retry/DLQ | Có sẵn qua plugin/policy | Phải tự implement (dùng thêm topic riêng) | Có sẵn qua Redrive Policy |

Câu trả lời phỏng vấn tốt: không có cái nào "tốt nhất" tuyệt đối — chọn dựa trên yêu cầu (cần replay event lịch sử → Kafka; cần routing phức tạp/RPC → RabbitMQ; muốn zero-ops, tích hợp sẵn AWS → SQS).

## Producer–Consumer

Producer: bên tạo và đẩy message vào queue (không cần biết ai xử lý). Consumer: bên lấy message ra và xử lý (không cần biết ai đã gửi). Đây chính là cách message queue **decouple** 2 phía — producer và consumer có thể scale độc lập, chạy tốc độ khác nhau.

Điểm quan trọng khi thiết kế consumer:
- **Acknowledgement (ack)**: chỉ báo "đã xử lý xong" sau khi thực sự xử lý thành công — nếu consumer crash giữa chừng mà chưa ack, message quay lại queue để consumer khác xử lý (at-least-once delivery).
- **Prefetch/concurrency**: giới hạn số message 1 consumer xử lý đồng thời, tránh 1 consumer bị quá tải trong khi consumer khác rảnh.

## Priority Queue

Một số hệ thống cần xử lý message quan trọng hơn trước (vd: giao dịch giá trị cao xử lý trước giao dịch giá trị thấp), dù chúng tới queue sau.

Có 2 cách phổ biến để làm priority queue với message queue:
1. **Nhiều queue theo mức ưu tiên** (vd: `queue.high`, `queue.low`) — consumer luôn ưu tiên đọc hết `queue.high` trước khi đọc `queue.low`. Đơn giản, dễ implement, hầu hết broker hỗ trợ tốt.
2. **Priority field trên message + priority queue nội bộ ở consumer** — RabbitMQ hỗ trợ `x-max-priority` trên queue để tự sắp xếp, hoặc consumer tự dùng cấu trúc heap (min-heap/max-heap) để sắp xếp trước khi xử lý.

## Retry

Khi xử lý 1 message thất bại (lỗi tạm thời: timeout, service phụ thuộc down...), không nên bỏ luôn — nên retry với **exponential backoff** (xem thêm ở [8-Distributed-Systems.md](8-Distributed-Systems.md)).

Best practice:
- Giới hạn số lần retry tối đa (vd: 3-5 lần), tránh retry vô hạn.
- Phân biệt lỗi **tạm thời** (nên retry: timeout, connection refused) với lỗi **vĩnh viễn** (không nên retry: dữ liệu JSON sai định dạng — retry mãi cũng không tự sửa được).
- Track số lần đã retry trên chính message (header hoặc field riêng) để biết khi nào nên đẩy sang DLQ.

## Dead-Letter Queue (DLQ)

Khi 1 message thất bại quá số lần retry cho phép, đẩy nó sang 1 queue riêng (DLQ) thay vì xử lý mãi hoặc mất luôn. Mục đích: không chặn các message khác phía sau (head-of-line blocking), đồng thời không mất dữ liệu để sau này debug/xử lý thủ công.

Quy trình chuẩn: `main queue` → xử lý thất bại → tăng retry count → nếu còn dưới ngưỡng thì requeue (có backoff) → nếu vượt ngưỡng thì đẩy sang `dlq` → có alert/monitor riêng theo dõi DLQ để người vận hành xử lý.

---

## Ví dụ thực hành: Transaction Processor có Priority Queue (Go)

Đây là dạng bài thường gặp trong test thực tế (tương tự đề: "xử lý transaction từ message queue, ưu tiên giao dịch giá trị cao, tự implement priority queue, không dùng thư viện ngoài cho phần priority"). Ví dụ dưới đây dùng `container/heap` (built-in Go, không phải thư viện ngoài) để tự viết priority queue, mô phỏng consumer đọc từ channel (đóng vai trò message queue).

```go
package main

import (
	"container/heap"
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

// Transaction matches the message structure from the queue.
type Transaction struct {
	ID        string    `json:"id"`
	Value     float64   `json:"value"`
	Timestamp time.Time `json:"timestamp"`
}

// --- Priority queue implementation using container/heap ---

type txItem struct {
	tx       Transaction
	priority float64 // higher value = higher priority
	index    int     // required by heap.Interface
}

type priorityQueue []*txItem

func (pq priorityQueue) Len() int { return len(pq) }

func (pq priorityQueue) Less(i, j int) bool {
	// max-heap: higher priority value comes out first
	return pq[i].priority > pq[j].priority
}

func (pq priorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *priorityQueue) Push(x any) {
	item := x.(*txItem)
	item.index = len(*pq)
	*pq = append(*pq, item)
}

func (pq *priorityQueue) Pop() any {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil
	*pq = old[:n-1]
	return item
}

// --- Processor: consumes from an input channel (simulating RabbitMQ),
//     buffers into a priority queue, and processes high-value transactions first ---

const highValueThreshold = 1000

type Processor struct {
	mu    sync.Mutex
	pq    priorityQueue
	notify chan struct{} // signals the worker that a new item is available
}

func NewProcessor() *Processor {
	return &Processor{
		pq:     make(priorityQueue, 0),
		notify: make(chan struct{}, 1),
	}
}

// Enqueue is called for every message received from the queue.
func (p *Processor) Enqueue(tx Transaction) {
	priority := tx.Value // simplest priority function: higher value = higher priority

	p.mu.Lock()
	heap.Push(&p.pq, &txItem{tx: tx, priority: priority})
	p.mu.Unlock()

	select {
	case p.notify <- struct{}{}:
	default: // worker already has a pending notification, no need to send another
	}
}

// Run starts the worker loop; call in a goroutine. Stops when ctx is cancelled.
func (p *Processor) Run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-p.notify:
			for {
				p.mu.Lock()
				if p.pq.Len() == 0 {
					p.mu.Unlock()
					break
				}
				item := heap.Pop(&p.pq).(*txItem)
				p.mu.Unlock()

				p.process(item.tx)
			}
		}
	}
}

func (p *Processor) process(tx Transaction) {
	label := "normal"
	if tx.Value > highValueThreshold {
		label = "HIGH-VALUE"
	}
	fmt.Printf("[processed:%s] id=%s value=%.2f timestamp=%s\n",
		label, tx.ID, tx.Value, tx.Timestamp.Format(time.RFC3339))
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	processor := NewProcessor()
	go processor.Run(ctx)

	// Simulate messages arriving from the queue out of priority order:
	// low value first, high value arrives later but must still jump the queue.
	raw := []string{
		`{"id":"tx1","value":100,"timestamp":"2026-01-01T10:00:00Z"}`,
		`{"id":"tx2","value":50,"timestamp":"2026-01-01T10:00:01Z"}`,
		`{"id":"tx3","value":5000,"timestamp":"2026-01-01T10:00:02Z"}`, // high-value, arrives 3rd
		`{"id":"tx4","value":20,"timestamp":"2026-01-01T10:00:03Z"}`,
	}

	for _, r := range raw {
		var tx Transaction
		if err := json.Unmarshal([]byte(r), &tx); err != nil {
			fmt.Println("skip invalid message:", err) // permanent error, do not retry
			continue
		}
		processor.Enqueue(tx)
	}

	time.Sleep(200 * time.Millisecond) // give the worker time to drain the queue
}
```

Code chạy được: [`code/priority-queue-processor/main.go`](../code/priority-queue-processor/main.go)
```bash
cd test/code/priority-queue-processor
go run main.go
```

Kết quả mẫu:
```
[processed:HIGH-VALUE] id=tx3 value=5000.00 timestamp=2026-01-01T10:00:02Z
[processed:normal] id=tx1 value=100.00 timestamp=2026-01-01T10:00:00Z
[processed:normal] id=tx2 value=50.00 timestamp=2026-01-01T10:00:01Z
[processed:normal] id=tx4 value=20.00 timestamp=2026-01-01T10:00:03Z
```
Đúng như kỳ vọng: `tx3` (value 5000, tới thứ 3 trong danh sách) được xử lý **trước** `tx1`, `tx2`, `tx4` dù chúng tới trước nó — vì processor luôn lấy phần tử ưu tiên cao nhất ra trước từ heap, thay vì xử lý theo đúng thứ tự tới (FIFO).

### Giải thích thiết kế
- `container/heap` là built-in Go, không phải "external library" — đúng yêu cầu đề bài dạng "tự implement priority logic".
- Channel `notify` (buffered 1) dùng làm tín hiệu "có việc mới" thay vì polling liên tục — worker chỉ thức dậy khi cần.
- `sync.Mutex` bảo vệ heap vì `Enqueue` (producer side) và `Run` (consumer side) chạy trên các goroutine khác nhau.
- Tách `Enqueue` (nhận message) và `Run` (xử lý) cho phép mô phỏng đúng producer-consumer: message có thể đến bất cứ lúc nào, xử lý diễn ra bất đồng bộ.
- Lỗi parse JSON là lỗi vĩnh viễn (permanent) → bỏ qua (hoặc gửi thẳng sang DLQ), không nên retry.

## Checklist tự kiểm tra
- [ ] So sánh được RabbitMQ/Kafka/SQS và biết chọn cái nào cho tình huống cụ thể
- [ ] Giải thích được ack/nack và tại sao nó liên quan tới at-least-once delivery
- [ ] Giải thích được 2 cách làm priority queue (nhiều queue theo mức, hoặc heap ở consumer)
- [ ] Phân biệt được lỗi tạm thời (nên retry) vs lỗi vĩnh viễn (nên bỏ/đẩy DLQ ngay)
- [ ] Tự code lại ví dụ Transaction Processor ở trên không nhìn mẫu, thử thêm test case retry khi `process()` trả lỗi giả lập
