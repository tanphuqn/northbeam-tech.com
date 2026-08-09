# Senior Backend Engineer (Go) — Interview Prep Notes

Mục tiêu: nắm vững Go concurrency (trọng tâm live-coding), message queue, distributed systems, system design, database, logging/monitoring, và code quality — tư duy thiết kế và giải quyết vấn đề là chính, không phải cú pháp 1 ngôn ngữ cụ thể.

## Bối cảnh

Bộ tài liệu này được soạn để ôn phỏng vấn vị trí **Senior Backend Engineer (Go/Python, Distributed Systems)** — role thực tế là **Go-first**, Python chỉ là secondary stack. Vài điểm cần nhớ về JD:

- Go là bắt buộc, không thương lượng — live coding dùng Golang, tập trung vào **Concurrency**: Goroutines / Channels / sync.Mutex / WaitGroup / Context / Race condition.
- Team làm 1 module backend cho hệ thống phân tán quy mô lớn, xử lý hàng triệu transaction/ngày (FoodTech platform, 50+ quốc gia).
- Quy trình: phỏng vấn nội bộ → endorse CV cho client Đức → 2 vòng kỹ thuật (1. Tech stack + live coding, 2. System Design).
- Nhà tuyển dụng sau đó gửi thêm danh sách chủ đề mở rộng (không phụ thuộc 1 ngôn ngữ cụ thể):
  Backend Fundamentals (Concurrency, Async, HTTP, Error Handling) · Message Queue (RabbitMQ/Kafka/SQS, Producer-Consumer, Priority Queue, Retry, DLQ) · System Design (Microservices, REST/gRPC, Event-driven, Load Balancing, Horizontal Scaling, Redis Caching, HA) · Distributed Systems (Idempotency, Data Consistency, Retry, Circuit Breaker, Eventual Consistency, Saga) · Database MongoDB (Query optimization, Index/Compound Index, Explain, TTL Index, Sharding, Archiving) · Logging & Monitoring (log level, structured logging, ELK/Grafana/CloudWatch) · Code Quality (Clean Code, SOLID/OOP, Validation, Exception Handling, Maintainability).
- Nhà tuyển dụng cũng gửi 1 đề bài mẫu thực tế (Question 1.1, 30 phút, code implementation): viết service Node.js xử lý transaction từ RabbitMQ, ưu tiên giao dịch giá trị cao (value > 1000), tự implement priority logic (không dùng thư viện priority queue ngoài). Vì mọi ví dụ trong bộ tài liệu này đều dùng Go cho nhất quán (quyết định đã chốt khi làm phần Message Queue), bài này được giải lại bằng Go ở [notes/11-Message-Queue.md](notes/11-Message-Queue.md) và [code/priority-queue-processor/](code/priority-queue-processor/) — logic/tư duy thiết kế giữ nguyên, chỉ đổi ngôn ngữ hiện thực.

## Cách bộ tài liệu này được xây dựng (để hiểu nhanh khi quay lại sau)

- Bắt đầu từ 1 file `note.md` duy nhất, sau đó theo yêu cầu được tách thành nhiều file nhỏ đánh số trong `notes/` để dễ đọc từng chủ đề riêng, mỗi file có mục **Checklist tự kiểm tra** ở cuối.
- `practice/questions.md` + `practice/answers.md` được tách riêng theo yêu cầu "chỉ đọc câu hỏi và tự làm" — cố tình để câu hỏi và đáp án ở 2 file khác nhau, tránh nhìn thấy đáp án ngay khi luyện tập. Bộ câu hỏi viết bằng **tiếng Anh** vì phỏng vấn thực tế sẽ trả lời bằng tiếng Anh.
- `code/` chứa các đoạn code runnable minh họa cho phần lý thuyết dễ hiểu nhầm (goroutine leak, `ch <- i` là gửi/block chứ không phải phép gán, priority queue transaction processor) — mỗi demo nằm trong 1 thư mục con riêng vì mỗi file đều là `package main` với `func main()` riêng, không thể gộp chung 1 thư mục.
- Mỗi file trong `notes/` và `practice/` có link "← Quay về note.md" ngay dưới tiêu đề để điều hướng nhanh về mục lục.
- Ban đầu thư mục này nằm trong repo `Verify365` (dùng làm scratch folder), sau đó được chuyển ra thư mục riêng độc lập — nội dung bên trong **không phụ thuộc vào code hay context của Verify365**, tất cả các đường link nội bộ đều là **đường dẫn tương đối** (`notes/...`, `../note.md`, `../code/...`) nên di chuyển cả thư mục `test/` sang vị trí mới vẫn hoạt động bình thường, không cần sửa gì thêm.

## Mục lục

### Backend Fundamentals (Concurrency, Async, Error Handling)
1. [Goroutines](notes/1-Goroutines.md)
2. [Channels](notes/2-Channels.md)
3. [sync.Mutex / sync.RWMutex](notes/3-sync-Mutex-RWMutex.md)
4. [sync.WaitGroup](notes/4-sync-WaitGroup.md)
5. [Context (context.Context)](notes/5-Context.md)
6. [Race Condition](notes/6-Race-Condition.md)
7. [Bài tập luyện live-coding](notes/7-Live-coding-Exercises.md)

### Distributed Systems & System Design
8. [Distributed Systems — Idempotency, Retry, Circuit Breaker, Saga, CAP, Eventual Consistency](notes/8-Distributed-Systems.md)
9. [System Design — Microservices, REST/gRPC, Event-driven, Load Balancing, Caching, HA](notes/9-System-Design.md)
10. [Checklist ôn nhanh trước ngày phỏng vấn](notes/10-Checklist-truoc-phong-van.md)

### Message Queue & Database
11. [Message Queue: RabbitMQ/Kafka/SQS, Producer-Consumer, Priority Queue, Retry, DLQ](notes/11-Message-Queue.md)
12. [Database (MongoDB): Index, Explain, TTL Index, Sharding, Archiving](notes/12-MongoDB.md)

### Vận hành & Chất lượng code
13. [Logging & Monitoring](notes/13-Logging-Monitoring.md)
14. [Code Quality: Clean Code, SOLID, Validation, Error Handling](notes/14-Code-Quality.md)

## Cấu trúc thư mục

- `notes/` — 14 file lý thuyết theo mục lục trên
- `practice/` — bộ câu hỏi luyện tập bằng tiếng Anh: [`questions.md`](practice/questions.md) (đọc và tự trả lời trước) + [`answers.md`](practice/answers.md) (đối chiếu sau)
- `code/` — code demo chạy được (mỗi demo trong 1 thư mục con riêng, chạy bằng `go run main.go`):
  - [`goroutine-leak/main.go`](code/goroutine-leak/main.go) — minh họa goroutine leak và 3 cách fix (dùng trong mục 1)
  - [`channel-send/main.go`](code/channel-send/main.go) — minh họa `ch <- i` là hành động gửi/block, không phải phép gán (dùng trong mục 2)
  - [`priority-queue-processor/main.go`](code/priority-queue-processor/main.go) — transaction processor ưu tiên giao dịch giá trị cao, tự implement priority queue bằng `container/heap` (dùng trong mục 11, sát với dạng bài test thực tế)

## Việc còn dang dở / có thể làm tiếp

- `notes/1-Goroutines.md` và `notes/2-Channels.md` đã có phần giải thích rất chi tiết (giữ nguyên transcript hỏi-đáp) — các file `notes/3` đến `notes/10` súc tích hơn; nếu cần đào sâu thêm phần nào, có thể hỏi lại theo đúng cách đã làm với Goroutines/Channels để mở rộng.
- Có thể cân nhắc thêm 1 demo runnable cho Worker Pool / Rate Limiter / TTL Cache (hiện mới chỉ có code mẫu trong `notes/7-Live-coding-Exercises.md`, chưa tách thành file chạy được trong `code/` như 3 demo kia).
- Sau khi chuyển sang thư mục mới `/Users/tanphuqn/Projects/me/northbeam-tech.com/interview`, có thể xóa dòng "ban đầu nằm trong repo Verify365" ở trên nếu muốn tài liệu gọn hơn — giữ lại chỉ để biết nguồn gốc.
