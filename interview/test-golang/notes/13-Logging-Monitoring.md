# Logging & Monitoring

[← Quay về note.md](../note.md)

## Log levels

Thứ tự phổ biến từ chi tiết nhất đến nghiêm trọng nhất:
`DEBUG` → `INFO` → `WARN` → `ERROR` → `FATAL`

- `DEBUG`: chi tiết kỹ thuật, chỉ bật khi cần điều tra sâu (thường tắt ở production để tránh log quá nhiều).
- `INFO`: sự kiện bình thường đáng ghi lại (request tới, job hoàn thành).
- `WARN`: điều bất thường nhưng chưa làm hỏng luồng xử lý (retry lần 1, dùng giá trị mặc định vì thiếu config).
- `ERROR`: 1 thao tác cụ thể thất bại, cần chú ý nhưng hệ thống vẫn chạy tiếp được.
- `FATAL`: lỗi khiến chương trình phải dừng hẳn.

Nguyên tắc: production nên chạy ở mức `INFO` trở lên (trừ khi đang debug tạm thời), vì `DEBUG` tạo quá nhiều log gây tốn chi phí lưu trữ và khó tìm log quan trọng.

## Structured logging (JSON)

Thay vì log dạng text tự do (`fmt.Println("user login failed: " + err.Error())`), log dưới dạng JSON có field rõ ràng:
```json
{"level":"error","msg":"user login failed","user_id":"123","reason":"invalid_password","trace_id":"abc-xyz","timestamp":"2026-07-16T10:00:00Z"}
```

Lý do quan trọng: log dạng JSON **query/filter được** dễ dàng ở hệ thống log tập trung (tìm tất cả log có `user_id=123`, hoặc đếm số lỗi theo `reason` trong 1 giờ) — text tự do rất khó làm việc này ở quy mô lớn.

Repo này dùng `zerolog` — đã là structured logging sẵn (mỗi log entry là JSON với field key-value), đúng theo khuyến nghị trên.

## Log aggregation (ELK / Grafana / CloudWatch)

Khi hệ thống chạy nhiều instance/service, log rải rác trên từng máy không dùng được — cần gom log tập trung về 1 nơi để tìm kiếm/dashboard/alert.

Các hệ sinh thái phổ biến:
- **ELK stack** (Elasticsearch + Logstash + Kibana): log được index vào Elasticsearch, Kibana dùng để search/visualize. Mạnh về full-text search log.
- **Grafana + Loki**: Loki lưu log (chỉ index metadata, không index full text như Elasticsearch nên rẻ hơn), Grafana dùng để query/dashboard — thường đi kèm Prometheus cho metrics.
- **CloudWatch Logs** (AWS): tích hợp sẵn với ECS/Lambda, không cần tự vận hành hạ tầng, nhưng khả năng query/dashboard hạn chế hơn ELK/Grafana.

Câu hỏi hay gặp: "Log tập trung giải quyết vấn đề gì mà log file trên từng server không giải quyết được?" — Trả lời: khả năng search/filter across nhiều instance cùng lúc, alert tự động khi có pattern lỗi bất thường, giữ log lại được sau khi container/instance đã bị xóa (container thường ephemeral, log trong container mất khi container chết).

## Monitoring & Alerting — khái niệm liên quan

- **Metrics** (khác với log): số liệu đo theo thời gian (request/s, latency p95/p99, error rate, CPU/memory) — dùng để phát hiện xu hướng và trigger alert, thường query nhanh hơn nhiều so với log.
- **Distributed tracing**: theo dõi 1 request đi qua nhiều service (mỗi service thêm 1 "span" vào cùng 1 `trace_id`) — giúp tìm ra service nào gây chậm trong 1 luồng request phức tạp.
- **Alerting**: đặt ngưỡng trên metric (vd: error rate > 5% trong 5 phút) để tự động báo cho on-call thay vì phải ngồi xem dashboard liên tục.

Trong context 1 service xử lý transaction: nên log ít nhất — thời điểm nhận message, thời điểm xử lý xong/thất bại, số lần retry, và luôn kèm `trace_id`/`transaction_id` để nối các log của cùng 1 giao dịch lại với nhau khi debug.

## Checklist tự kiểm tra
- [ ] Kể đúng thứ tự log level và giải thích khi nào dùng mức nào
- [ ] Giải thích được vì sao structured logging (JSON) tốt hơn log text tự do ở quy mô lớn
- [ ] So sánh được ELK vs Grafana/Loki vs CloudWatch — biết chọn cái nào tùy tình huống
- [ ] Phân biệt được log vs metrics vs tracing — 3 cái bổ trợ nhau, không thay thế nhau
- [ ] Giải thích được vì sao cần `trace_id`/`transaction_id` xuyên suốt log của 1 giao dịch
