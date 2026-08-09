# System Design — checklist khi được hỏi

[← Quay về note.md](../note.md)

Khi được hỏi thiết kế 1 hệ thống (vd: "thiết kế hệ thống xử lý transaction cho food delivery app"), nên đi theo trình tự:

1. **Clarify yêu cầu**: bao nhiêu request/giây, đọc nhiều hay ghi nhiều, cần strong consistency hay eventual consistency được, real-time hay có thể trễ vài giây.
2. **High-level architecture**: API Gateway → Service → Queue → Worker → DB. Vẽ được sơ đồ đơn giản.
3. **Data storage**: chọn DB phù hợp (PostgreSQL cho transactional, cân nhắc sharding/partitioning theo `region` hoặc `user_id` khi scale lớn).
4. **Caching**: cache-aside (đọc: check cache trước, miss thì đọc DB rồi set cache) vs write-through (ghi cache đồng thời với DB).
5. **Scaling**: horizontal scaling qua ECS/K8s, load balancer, stateless service để scale dễ dàng.
6. **Failure handling**: retry, circuit breaker, idempotency (liên hệ [8-Distributed-Systems.md](8-Distributed-Systems.md)).
7. **Trade-off**: luôn nói rõ bạn đánh đổi gì (vd: eventual consistency để đổi lấy throughput cao hơn).

**Chuẩn bị trước**: chọn sẵn 1 dự án thực tế bạn từng làm có liên quan tới xử lý transaction/high load, kể được: bottleneck cụ thể là gì, bạn đã scale/optimize ra sao, và trade-off nào bạn đã chọn. Họ chắc chắn sẽ hỏi câu này.

---

## Microservices

Chia hệ thống lớn thành nhiều service nhỏ, độc lập, mỗi service sở hữu 1 domain/data riêng (vd: `payment-service`, `order-service`, `user-service`). Mỗi service có thể deploy, scale, thay công nghệ độc lập.

Trade-off cần nói được:
- **Lợi**: scale độc lập từng phần, team làm việc song song không đụng nhau, fail của 1 service không nhất thiết sập toàn hệ thống.
- **Hại**: phức tạp hơn nhiều (network call thay vì function call, cần service discovery, distributed tracing, cần xử lý partial failure, data consistency giữa các service khó hơn monolith).

Câu hỏi hay gặp: "Khi nào KHÔNG nên dùng microservices?" — Trả lời: khi team nhỏ, domain chưa rõ ràng, hoặc chưa có nhu cầu scale riêng biệt — lúc đó monolith (module hóa tốt) đơn giản và nhanh hơn.

## REST vs gRPC

| | REST (JSON qua HTTP) | gRPC (protobuf qua HTTP/2) |
|---|---|---|
| Định dạng | JSON, dễ đọc, debug bằng tay được | Binary (protobuf), nhỏ gọn, nhanh hơn |
| Hợp đồng API | OpenAPI/Swagger (tùy chọn) | `.proto` file bắt buộc, generate code tự động 2 phía |
| Dùng khi | Public API, browser gọi trực tiếp, cần dễ debug | Giao tiếp internal giữa các service, cần performance cao, streaming |
| Streaming | Khó, phải dùng WebSocket/SSE riêng | Hỗ trợ sẵn (client streaming, server streaming, bidirectional) |

Trong hệ thống microservices thực tế: REST thường dùng ở tầng public-facing API (frontend gọi vào), gRPC dùng cho giao tiếp internal service-to-service vì nhanh và có strict schema.

## Event-driven Architecture

Thay vì service A gọi trực tiếp service B (request/response đồng bộ), service A **publish 1 event** (vd: `OrderCreated`) lên message broker, các service quan tâm (B, C, D) tự **subscribe** và xử lý độc lập.

Lợi ích: giảm coupling (A không cần biết B, C, D tồn tại), dễ thêm consumer mới mà không sửa code A, chịu tải tốt hơn (queue đóng vai trò buffer khi traffic tăng đột biến).
Đánh đổi: khó trace toàn bộ luồng xử lý hơn (cần distributed tracing), eventual consistency thay vì strong consistency, cần xử lý duplicate/out-of-order event.

## Load Balancing

Phân phối traffic đều cho nhiều instance của cùng 1 service, tránh 1 instance quá tải trong khi instance khác rảnh.

Thuật toán phổ biến:
- **Round robin**: chia đều lần lượt.
- **Least connections**: ưu tiên instance đang xử lý ít connection nhất.
- **Consistent hashing**: dùng khi cần cùng 1 client luôn được route tới cùng 1 instance (vd: sticky session, cache locality).

Có 2 tầng load balancing thường gặp: L4 (transport layer, dựa vào IP/port, nhanh) và L7 (application layer, dựa vào nội dung HTTP request như path/header, linh hoạt hơn nhưng chậm hơn 1 chút).

## Horizontal Scaling

Thêm nhiều instance chạy song song (horizontal = thêm máy, khác với vertical = nâng cấp máy mạnh hơn).

Điều kiện tiên quyết: service phải **stateless**.

### Stateless nghĩa là gì

Stateless nghĩa là mỗi request được xử lý **độc lập hoàn toàn**, không phụ thuộc vào dữ liệu được lưu lại từ request trước đó ngay trên chính instance đó — nói cách khác, state **không được giữ cục bộ trên 1 instance cụ thể**, mà phải đẩy ra một nơi lưu trữ dùng chung (Redis, DB, S3) mà mọi instance đều truy cập được.

```go
// Stateful — dễ gây bug khi scale ngang
var userSessions = map[string]string{} // sống trong RAM của CHÍNH instance này

func login(userID string) {
    token := generateToken()
    userSessions[userID] = token // chỉ instance này biết session này tồn tại
}

// Stateless — an toàn khi scale ngang
func login(userID string, redisClient *redis.Client) {
    token := generateToken()
    redisClient.Set(userID, token, ttl) // lưu ở Redis, instance nào cũng đọc được
}
```

Vì sao quan trọng: khi chạy nhiều instance sau load balancer, request được phân phối ngẫu nhiên/luân phiên tới bất kỳ instance nào — request 1 (login) có thể rơi vào instance A, request 2 (check auth, cùng user) có thể rơi vào instance B. Nếu service stateful (lưu session trong RAM), instance B sẽ không thấy session đó → user bị coi như chưa đăng nhập dù vừa login xong. Đây là bug kinh điển khi scale ngang mà quên chuyển state ra ngoài.

Trong AWS ECS: cấu hình Auto Scaling dựa trên metric (CPU %, memory %, hoặc custom metric như queue depth) để tự thêm/giảm task.

## Caching (Redis)

Dùng để giảm tải cho DB và giảm latency cho các dữ liệu đọc nhiều, ít thay đổi.

Các pattern chính (chi tiết cache-aside/write-through xem mục 4 ở trên):
- **Cache-aside**: phổ biến nhất, đơn giản, ứng dụng tự quản lý cache.
- **TTL (time-to-live)**: luôn đặt TTL cho cache entry để tránh dữ liệu cũ tồn tại mãi nếu invalidation bị sót.
- **Cache invalidation**: khó nhất trong caching ("There are only two hard things in Computer Science: cache invalidation and naming things"). Chiến lược: xóa cache khi update DB (thay vì update cache), hoặc dùng TTL ngắn nếu độ chính xác không cần tuyệt đối.

Redis cụ thể hay được hỏi:
- Redis là in-memory nên rất nhanh, nhưng dữ liệu có thể mất khi restart nếu không bật persistence (RDB snapshot / AOF log).
- Cấu trúc dữ liệu ngoài string: List, Set, Sorted Set (dùng cho leaderboard, priority theo score), Hash.
- Redis cũng dùng được làm message broker đơn giản (Pub/Sub) hoặc rate limiter (dùng lệnh `INCR` + `EXPIRE`).

## High Availability (HA)

Mục tiêu: hệ thống vẫn hoạt động khi 1 phần bị lỗi (1 instance chết, 1 AZ down...).

### Availability Zone (AZ) là gì

AZ là khái niệm hạ tầng của AWS (cloud khác cũng có tương tự): 1 **Region** (vd: `ap-southeast-1` Singapore) được chia thành nhiều **Availability Zone**, mỗi AZ là 1 hoặc nhiều data center vật lý **tách biệt nhau** (khác tòa nhà, khác nguồn điện, khác hệ thống làm mát), nhưng vẫn nằm trong cùng region nên kết nối mạng giữa các AZ rất nhanh.

```
Region: ap-southeast-1
├── AZ-1 (ap-southeast-1a)  ← data center A
├── AZ-2 (ap-southeast-1b)  ← data center B
└── AZ-3 (ap-southeast-1c)  ← data center C
```

"**AZ down**" nghĩa là toàn bộ 1 data center đó gặp sự cố (mất điện, lỗi mạng lớn...) khiến mọi máy chủ trong AZ đó ngừng hoạt động cùng lúc — đây là sự cố có thật đã từng xảy ra với AWS.

Nếu toàn bộ service (app + DB) chỉ chạy trên **1 AZ duy nhất**, AZ đó down → toàn bộ hệ thống sập, dù các AZ khác trong cùng region vẫn bình thường:

```
Setup KHÔNG an toàn:
Load Balancer → [Instance 1, 2, 3 đều ở AZ-1]  → AZ-1 down = 100% instance chết

Setup an toàn (multi-AZ):
Load Balancer → [Instance 1 (AZ-1), Instance 2 (AZ-2), Instance 3 (AZ-3)]
→ AZ-1 down = chỉ mất 1/3 instance, 2 AZ còn lại vẫn phục vụ traffic
```

### Kỹ thuật chính để đạt HA

- **Redundancy**: chạy ít nhất 2 instance của mỗi service, ở ít nhất 2 Availability Zone khác nhau — không chỉ nhiều instance mà phải nhiều instance ở **nhiều AZ khác nhau**.
- **Health check**: load balancer/orchestrator (ECS) tự động loại instance không phản hồi health check ra khỏi traffic.
- **Failover**: DB có replica (read replica cho PostgreSQL, hoặc **Multi-AZ RDS** — AWS tự duy trì 1 bản đồng bộ ở AZ khác, tự động failover khi AZ chứa primary down).
- **Graceful degradation**: khi 1 phần hệ thống lỗi, phần còn lại vẫn hoạt động ở mức giảm tính năng thay vì sập toàn bộ (vd: tắt tính năng gợi ý sản phẩm nhưng vẫn cho checkout được).

### SLA (Service Level Agreement) là gì

SLA là bản cam kết (thường mang tính hợp đồng) giữa nhà cung cấp dịch vụ và khách hàng về mức độ chất lượng dịch vụ, phổ biến nhất là cam kết **uptime** — hệ thống hoạt động bao nhiêu % thời gian. Biểu diễn bằng số "9", càng nhiều số 9 thì downtime cho phép càng ít:

| SLA | % uptime | Downtime tối đa/năm | Downtime tối đa/tháng |
|---|---|---|---|
| 2 số 9 | 99% | ~3.65 ngày | ~7.3 giờ |
| 3 số 9 | 99.9% | ~8.7 giờ | ~43.8 phút |
| 4 số 9 | 99.99% | ~52.6 phút | ~4.4 phút |
| 5 số 9 | 99.999% | ~5.26 phút | ~26 giây |

Ví dụ thực tế: AWS thường cam kết SLA 99.99% cho EC2/RDS trong 1 region — nếu không đạt, họ phải hoàn tiền (service credit) theo hợp đồng.

SLA quyết định kiến trúc bạn phải xây: muốn 99.9% chỉ cần vài kỹ thuật cơ bản (health check + auto-restart, 1 replica DB); muốn 99.99% trở lên bắt buộc phải multi-AZ, failover tự động, monitoring/alerting real-time — chi phí hạ tầng cao hơn nhiều vì phải chạy dư thừa (redundancy) ở nhiều nơi.

Phân biệt 3 khái niệm hay bị nhầm:
- **SLA**: cam kết chính thức, có ràng buộc (thường có phạt/bồi thường nếu không đạt).
- **SLO** (Service Level Objective): mục tiêu nội bộ team tự đặt ra (vd: "target 99.95%") — không nhất thiết có ràng buộc hợp đồng.
- **SLI** (Service Level Indicator): con số đo thực tế (vd: uptime thực đo được tháng này là 99.97%) — dùng để so sánh với SLO/SLA.

Câu trả lời tốt khi phỏng vấn hỏi "hệ thống bạn từng làm có SLA bao nhiêu": nói rõ con số cụ thể (99.9%/99.99%) VÀ kỹ thuật cụ thể bạn dùng để đạt được nó (multi-AZ, health check, retry, circuit breaker...), tránh chỉ nói con số suông.

## Checklist tự kiểm tra
- [ ] Đi được hết 7 bước ở phần đầu khi gặp 1 đề bài system design mới, không bỏ bước nào
- [ ] Chuẩn bị sẵn 1 câu chuyện dự án thực tế (bottleneck → giải pháp → trade-off)
- [ ] Vẽ được sơ đồ kiến trúc đơn giản trên giấy/whiteboard trong vài phút
- [ ] Luôn nói rõ trade-off thay vì chỉ đưa ra 1 giải pháp "hoàn hảo"
- [ ] Giải thích được khi nào NÊN và KHÔNG NÊN dùng microservices
- [ ] Phân biệt được REST vs gRPC và biết chọn cái nào cho internal vs public API
- [ ] Giải thích được event-driven architecture đánh đổi gì so với request/response đồng bộ
- [ ] Phân biệt được load balancing L4 vs L7
- [ ] Giải thích được stateless là gì và vì sao service phải stateless mới scale horizontal dễ dàng
- [ ] Giải thích được cache invalidation là bài toán khó và chiến lược bạn dùng
- [ ] Giải thích được Availability Zone (AZ) là gì và vì sao multi-AZ quan trọng cho HA
- [ ] Nói được 2-3 con số SLA/uptime phổ biến (99.9%, 99.99%) và phân biệt SLA/SLO/SLI
