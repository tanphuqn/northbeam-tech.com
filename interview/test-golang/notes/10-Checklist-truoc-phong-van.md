# Checklist ôn nhanh trước ngày phỏng vấn

[← Quay về note.md](../note.md)

## Concurrency & Go fundamentals
- [ ] Tự code lại 3 bài tập ở [7-Live-coding-Exercises.md](7-Live-coding-Exercises.md) mà không xem đáp án
- [ ] Giải thích được sự khác biệt Mutex vs RWMutex vs atomic — khi nào dùng cái nào ([3-sync-Mutex-RWMutex.md](3-sync-Mutex-RWMutex.md))
- [ ] Giải thích được vì sao cần `context` thay vì chỉ dùng channel + timer thủ công ([5-Context.md](5-Context.md))
- [ ] Biết chạy `go run -race` và đọc hiểu output khi có race ([6-Race-Condition.md](6-Race-Condition.md))

## Message Queue & Distributed Systems
- [ ] So sánh được RabbitMQ/Kafka/SQS và biết chọn cái nào cho tình huống cụ thể ([11-Message-Queue.md](11-Message-Queue.md))
- [ ] Tự code lại bài Priority Queue Transaction Processor không nhìn mẫu ([11-Message-Queue.md](11-Message-Queue.md), code tại `code/priority-queue-processor/main.go`)
- [ ] Giải thích được idempotency, retry+backoff, circuit breaker, Saga, CAP/eventual consistency ([8-Distributed-Systems.md](8-Distributed-Systems.md))

## System Design
- [ ] Kể được 1 câu chuyện dự án thực tế về high-load/distributed system (bottleneck → giải pháp → trade-off) ([9-System-Design.md](9-System-Design.md))
- [ ] Giải thích được microservices, REST vs gRPC, event-driven, load balancing, caching Redis, HA ([9-System-Design.md](9-System-Design.md))

## Database & Vận hành
- [ ] Ôn PostgreSQL: index, EXPLAIN ANALYZE, connection pooling
- [ ] Ôn MongoDB: compound index (ESR), explain, TTL index, sharding, archiving ([12-MongoDB.md](12-MongoDB.md))
- [ ] Giải thích được log level, structured logging, log aggregation (ELK/Grafana/CloudWatch) ([13-Logging-Monitoring.md](13-Logging-Monitoring.md))
- [ ] Ôn lại AWS ECS/S3/CloudWatch ở mức khái niệm (không cần chuyên sâu)

## Code Quality
- [ ] Refactor được code nesting sâu thành early-return, giải thích SOLID bằng ví dụ Go cụ thể ([14-Code-Quality.md](14-Code-Quality.md))
