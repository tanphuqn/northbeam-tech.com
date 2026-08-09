# Code Quality: Clean Code, Modular Design, SOLID/OOP, Validation, Exception Handling, Maintainability

[← Quay về note.md](../note.md)

## Clean Code — nguyên tắc cốt lõi

- Đặt tên rõ ràng, nói lên được mục đích (`processTransaction` tốt hơn `doStuff`, `isHighValue` tốt hơn `flag`).
- Hàm nên làm **1 việc** và làm tốt việc đó (single responsibility ở mức hàm) — nếu tên hàm cần chữ "and" để mô tả, có thể nó đang làm quá nhiều việc.
- Tránh nesting quá sâu (if lồng if lồng if) — dùng early return để giảm độ sâu:
```go
// Khó đọc: nesting sâu
func process(tx Transaction) error {
    if tx.Value > 0 {
        if tx.ID != "" {
            // logic chính
            return nil
        } else {
            return errors.New("missing id")
        }
    } else {
        return errors.New("invalid value")
    }
}

// Rõ ràng hơn: early return, guard clause
func process(tx Transaction) error {
    if tx.Value <= 0 {
        return errors.New("invalid value")
    }
    if tx.ID == "" {
        return errors.New("missing id")
    }
    // logic chính
    return nil
}
```
- Comment chỉ nên giải thích **vì sao** (lý do không hiển nhiên), không giải thích **cái gì** (code rõ ràng đã tự nói lên điều đó).

## Modular Design

Chia code theo domain/trách nhiệm rõ ràng thay vì gom hết vào 1 file/hàm khổng lồ. Trong Go, thường tách theo layer: `controller` (nhận request, validate input) → `service` (business logic) → `repository/DB layer` (truy vấn dữ liệu) — đúng như pattern hiện có trong repo Verify365 (`controllers/` → `services/` → GORM models).

Lợi ích: dễ test từng phần độc lập (mock service khi test controller), dễ thay đổi 1 phần mà không ảnh hưởng phần khác (đổi DB không cần sửa controller), nhiều người có thể làm việc song song trên các module khác nhau.

## SOLID — áp dụng vào Go (Go không có class/inheritance như Java, nhưng vẫn áp dụng được qua interface)

- **S (Single Responsibility)**: 1 struct/hàm chỉ nên có 1 lý do để thay đổi. Vd: `PaymentService` chỉ lo logic thanh toán, không tự viết log ra file hay gọi email trực tiếp (nên inject `Logger`, `EmailService` riêng).
- **O (Open/Closed)**: mở rộng được mà không cần sửa code cũ — trong Go thường làm qua interface: thêm 1 provider thanh toán mới (Stripe → thêm PayPal) bằng cách implement cùng 1 interface `PaymentProvider`, không cần sửa code gọi hàm cũ.
- **L (Liskov Substitution)**: bất kỳ implementation nào của 1 interface đều phải dùng thay thế được cho nhau mà không phá vỡ logic gọi nó.
- **I (Interface Segregation)**: interface nhỏ, chỉ chứa method thực sự cần — thay vì 1 interface khổng lồ mà struct nào cũng phải implement thừa method không dùng tới.
- **D (Dependency Inversion)**: service nên phụ thuộc vào **interface** (abstraction) thay vì phụ thuộc trực tiếp vào 1 implementation cụ thể — giúp dễ mock khi test, dễ đổi implementation sau này.

```go
// Dependency inversion: PaymentService phụ thuộc vào interface, không phụ thuộc trực tiếp Stripe
type PaymentProvider interface {
    Charge(ctx context.Context, amount float64) error
}

type PaymentService struct {
    provider PaymentProvider // interface, không phải *StripeClient cụ thể
}

func (s *PaymentService) Pay(ctx context.Context, amount float64) error {
    return s.provider.Charge(ctx, amount)
}
```

## Validation

Luôn validate ở **boundary** của hệ thống (nơi dữ liệu đi vào: HTTP request body, message queue payload) — không tin tưởng dữ liệu đầu vào từ bên ngoài. Sau khi qua khỏi boundary và đã validate, phần code bên trong (service, business logic) có thể tin tưởng dữ liệu hợp lệ, không cần validate lại nhiều lần ở mọi hàm.

Ví dụ validate message từ queue trước khi xử lý (liên hệ ví dụ ở [11-Message-Queue.md](11-Message-Queue.md)):
```go
func validateTransaction(tx Transaction) error {
    if tx.ID == "" {
        return errors.New("id is required")
    }
    if tx.Value <= 0 {
        return errors.New("value must be positive")
    }
    return nil
}
```

## Exception Handling (error handling trong Go)

Go không có exception/try-catch — dùng error trả về tường minh (`if err != nil`). Nguyên tắc:
- Luôn kiểm tra error ngay sau khi gọi hàm có thể trả lỗi, không bỏ qua bằng `_`.
- Wrap error để giữ context khi trả ngược lên trên: `fmt.Errorf("process transaction %s: %w", tx.ID, err)` — `%w` giữ được error gốc để caller có thể `errors.Is`/`errors.As` kiểm tra loại lỗi cụ thể.
- Phân biệt lỗi **tạm thời** (nên retry — timeout, connection refused) với lỗi **nghiệp vụ/vĩnh viễn** (không nên retry — validation failed, not found) — xử lý khác nhau ở tầng gọi (liên hệ mục Retry ở [11-Message-Queue.md](11-Message-Queue.md)).
- `panic`/`recover` chỉ nên dùng cho lỗi thực sự không thể phục hồi (chương trình ở trạng thái không nhất quán) — không dùng thay cho error handling thông thường.

## Maintainability

- Viết test cho business logic quan trọng (đặc biệt là các nhánh lỗi/edge case), không chỉ test happy path.
- Đặt tên package/file theo domain, không theo loại kỹ thuật chung chung (vd: `payment/` tốt hơn `utils/` chứa lẫn lộn mọi thứ).
- Tránh premature abstraction: không tạo interface/abstraction cho thứ chỉ có 1 implementation và không có dấu hiệu sẽ có implementation thứ 2 — codebase đơn giản hơn thì dễ maintain hơn.
- Document (comment/README) tập trung vào lý do thiết kế (why), không lặp lại điều code đã tự nói (what).

## Checklist tự kiểm tra
- [ ] Refactor được 1 đoạn code nesting sâu thành early-return/guard clause
- [ ] Giải thích được từng chữ trong SOLID bằng ví dụ Go cụ thể (không chỉ định nghĩa suông)
- [ ] Giải thích được vì sao nên validate ở boundary, không cần validate lại nhiều lớp bên trong
- [ ] Giải thích được cách dùng `%w` để wrap error và vì sao nó tốt hơn nối chuỗi lỗi thủ công
- [ ] Phân biệt được khi nào dùng `panic` (hiếm, lỗi không thể phục hồi) vs error trả về bình thường
