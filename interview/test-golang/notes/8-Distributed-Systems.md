# Distributed Systems — khái niệm cần nói được khi phỏng vấn

[← Quay về note.md](../note.md)

| Khái niệm | Giải thích ngắn | Khi nào dùng / ví dụ |
|---|---|---|
| **Idempotency** | Gọi API/thao tác nhiều lần cho cùng 1 kết quả, không tạo tác dụng phụ lặp lại | Payment API: dùng `idempotency-key` để tránh charge tiền 2 lần khi client retry |
| **Outbox pattern** | Ghi event vào bảng "outbox" cùng transaction với business data, rồi 1 process riêng đọc và publish lên message queue | Tránh mất event khi DB commit thành công nhưng publish message thất bại |
| **Retry + Exponential backoff** | Khi gọi service khác thất bại, thử lại sau khoảng thời gian tăng dần (1s, 2s, 4s...) thay vì retry liên tục | Gọi API bên thứ 3 (CreditSafe, Stripe...) bị timeout |
| **Circuit breaker** | Khi 1 dependency lỗi liên tục, "ngắt mạch" tạm thời — trả lỗi ngay thay vì tiếp tục gọi, để tránh cascading failure | Service downstream đang down, tránh làm nghẽn toàn hệ thống |
| **Dead-letter queue (DLQ)** | Message xử lý thất bại nhiều lần được đẩy sang queue riêng để xử lý thủ công/sau | Xử lý transaction lỗi liên tục do dữ liệu sai |
| **Saga pattern** | Chia 1 transaction phân tán thành chuỗi bước nhỏ, mỗi bước có "compensating action" để rollback nếu bước sau thất bại | Đặt hàng + trừ kho + thanh toán trên nhiều service khác nhau |
| **Optimistic locking** | Thêm cột `version`, khi update thì check version có đổi không — nếu có người khác đã update thì fail và retry | Tránh lock DB lâu khi update tần suất cao |
| **Data consistency (CAP)** | Trade-off giữa Consistency, Availability, Partition tolerance — không thể có cả 3 cùng lúc khi có network partition | Chọn CP (ưu tiên đúng dữ liệu, có thể tạm từ chối request) hay AP (luôn trả lời, chấp nhận dữ liệu có thể cũ) tùy nghiệp vụ |
| **Eventual consistency** | Dữ liệu ở các node/service có thể tạm thời khác nhau, nhưng cuối cùng sẽ hội tụ về cùng 1 giá trị nếu không có write mới | Đếm số like trên mạng xã hội — hiển thị hơi trễ vài giây không sao, miễn cuối cùng đúng |

## Giải thích sâu hơn: Data Consistency & Eventual Consistency

Trong hệ thống phân tán, dữ liệu thường được replicate ra nhiều node để chịu tải và chịu lỗi tốt hơn. Câu hỏi đặt ra: khi ghi vào 1 node, các node khác thấy giá trị mới **ngay lập tức** (strong consistency) hay **sau một khoảng trễ** (eventual consistency)?

- **Strong consistency**: mọi lần đọc sau khi ghi đều thấy giá trị mới nhất, ở bất kỳ node nào. Đổi lại: chậm hơn (phải đợi đồng bộ xong mới trả response), sẵn sàng kém hơn khi có network partition (theo định lý CAP).
- **Eventual consistency**: ghi xong trả response ngay, việc đồng bộ ra các node khác diễn ra "sau đó". Đổi lại: nhanh hơn, sẵn sàng cao hơn, nhưng có khoảng thời gian ngắn dữ liệu đọc được có thể chưa cập nhật (stale read).

Khi trả lời phỏng vấn, nên nói rõ: bạn chọn cái nào cho bài toán cụ thể, và vì sao. Ví dụ: số dư tài khoản ngân hàng cần strong consistency (không thể sai dù chỉ 1 giây); số lượt xem video có thể eventual consistency (chấp nhận đếm hơi trễ để đổi lấy throughput cao).

## Checklist tự kiểm tra
- [ ] Giải thích được từng khái niệm trong bảng trên bằng lời của chính mình (không đọc lại)
- [ ] Cho được ví dụ thực tế khác (ngoài ví dụ trong bảng) cho mỗi khái niệm
- [ ] Phân biệt được circuit breaker vs retry — hai cái không thay thế nhau
- [ ] Giải thích được vì sao Saga phù hợp hơn 2PC trong hệ thống phân tán lớn
- [ ] Giải thích được định lý CAP và cho ví dụ 1 hệ thống nên chọn CP, 1 hệ thống nên chọn AP
- [ ] Phân biệt được strong consistency vs eventual consistency, cho ví dụ mỗi loại
