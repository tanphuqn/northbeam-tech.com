# Database (MongoDB): Query Optimization, Index, Explain, TTL Index, Sharding, Data Archiving

[← Quay về note.md](../note.md)

Lưu ý: repo Verify365 hiện dùng PostgreSQL/GORM, nhưng JD/đề test có nhắc MongoDB — phần này để bạn nắm khái niệm chung, không phụ thuộc dự án cụ thể.

## Query Optimization cơ bản

- Luôn query theo field đã có index, tránh full collection scan (`COLLSCAN`).
- Tránh `$where` (chạy JavaScript, rất chậm) và regex không có anchor đầu chuỗi (`/abc/` chậm hơn `/^abc/` vì không tận dụng được index).
- Chỉ select field cần thiết (projection) thay vì lấy cả document, đặc biệt khi document có field lớn (mảng dài, blob).
- Dùng `limit()` sớm khi chỉ cần vài kết quả, kết hợp index để tránh sort toàn bộ collection trong bộ nhớ.

## Index & Compound Index

- **Single-field index**: tăng tốc query lọc theo 1 field (`db.collection.createIndex({ email: 1 })`).
- **Compound index**: index trên nhiều field cùng lúc (`{ firmId: 1, status: 1, createdAt: -1 }`) — thứ tự field trong compound index rất quan trọng, theo nguyên tắc **ESR** (Equality → Sort → Range): field dùng để so sánh bằng (`=`) đặt trước, field dùng để sort đặt giữa, field dùng cho range (`$gt`, `$lt`) đặt cuối.
- 1 compound index có thể phục vụ nhiều query khác nhau nếu chúng dùng **tiền tố (prefix)** của index đó — vd: index `{a:1, b:1, c:1}` phục vụ được query lọc theo `a`, hoặc `a+b`, hoặc `a+b+c`, nhưng KHÔNG phục vụ tốt cho query chỉ lọc theo `b` hoặc `c` một mình.
- Đánh đổi: mỗi index thêm vào giúp đọc nhanh hơn nhưng làm chậm ghi (insert/update phải cập nhật thêm index) và tốn thêm dung lượng lưu trữ — không nên tạo index tùy tiện cho mọi field.

## Explain

`db.collection.find({...}).explain("executionStats")` cho biết:
- Query có dùng index không (`IXSCAN`) hay quét toàn bộ (`COLLSCAN`).
- Số document thực sự được kiểm tra (`totalDocsExamined`) so với số document trả về (`nReturned`) — nếu 2 số này lệch nhau quá nhiều, nghĩa là index chưa hiệu quả (đang phải quét nhiều hơn cần thiết).
- Thời gian thực thi (`executionTimeMillis`).

Dùng khi: debug query chậm, hoặc trước khi merge 1 query mới vào production để chắc chắn nó dùng đúng index.

## TTL Index

Tự động xóa document sau một khoảng thời gian, dựa trên 1 field kiểu Date:
```js
db.sessions.createIndex({ createdAt: 1 }, { expireAfterSeconds: 3600 }) // tự xóa sau 1 giờ
```
Dùng cho: session tạm thời, OTP/verification code, log ngắn hạn, cache document — bất cứ dữ liệu nào có "hạn dùng" rõ ràng, tránh phải viết cron job dọn dẹp thủ công.

Lưu ý: MongoDB chạy background task dọn TTL mỗi ~60 giây, nên document có thể tồn tại thêm một chút sau khi hết hạn — không dùng TTL index cho yêu cầu xóa chính xác tức thời.

## Sharding

Chia 1 collection lớn ra nhiều **shard** (mỗi shard là 1 replica set riêng), dựa trên **shard key** — MongoDB tự route query/write tới đúng shard dựa vào giá trị shard key.

Điểm quan trọng khi chọn shard key:
- Cần có **độ phân tán tốt** (high cardinality) để dữ liệu chia đều các shard, tránh "hot shard" (1 shard nhận hầu hết traffic).
- Query không có shard key trong điều kiện lọc sẽ phải **scatter-gather** (hỏi tất cả shard rồi gộp kết quả) — chậm hơn nhiều so với query trúng đúng 1 shard.
- Ví dụ tốt: shard theo `userId` hoặc `tenantId` nếu phần lớn query đều lọc theo field đó. Ví dụ xấu: shard theo `status` (chỉ vài giá trị cố định như "active"/"inactive") → dữ liệu dồn hết vào 1-2 shard.

## Data Archiving

Khi collection quá lớn, dữ liệu cũ ít truy cập nên được archive (chuyển sang collection/storage khác rẻ hơn) thay vì giữ mãi trong collection chính — giúp giảm kích thước index, tăng tốc query trên dữ liệu "nóng" (hay truy cập).

Chiến lược phổ biến:
- Archive theo thời gian (vd: chuyển document cũ hơn 1 năm sang collection `orders_archive` hoặc sang S3/cold storage).
- Kết hợp TTL index để tự xóa dữ liệu ở collection chính sau khi đã archive xong.
- Cân nhắc archive bất đồng bộ qua batch job định kỳ, tránh chạy trong giờ cao điểm vì archive thường phải đọc/ghi lượng lớn dữ liệu.

## Checklist tự kiểm tra
- [ ] Giải thích được nguyên tắc ESR khi thiết kế compound index
- [ ] Giải thích được vì sao 1 compound index chỉ phục vụ tốt cho query dùng đúng "tiền tố" của nó
- [ ] Đọc hiểu được output `explain()` cơ bản: IXSCAN vs COLLSCAN, totalDocsExamined vs nReturned
- [ ] Biết khi nào dùng TTL index thay vì cron job dọn dẹp thủ công
- [ ] Giải thích được tiêu chí chọn shard key tốt, và hậu quả khi chọn sai (hot shard, scatter-gather)
- [ ] Giải thích được lý do và chiến lược archive dữ liệu cũ
