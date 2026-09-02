---
title: Go 1.27 được phát hành
date: 2026-08-19
by:
- Nicholas Husin, on behalf of the Go team
summary: Go 1.27 bổ sung các phương thức generic, gói encoding/json/v2, gói uuid, phân bổ bộ nhớ nhanh hơn, hồ sơ rò rỉ goroutine và nhiều nội dung khác.
---

Hôm nay nhóm Go vui mừng bản phát hành Go 1.27. Bạn có thể tìm thấy các kho lưu trữ binary và trình cài đặt của bản phát hành này trên [trang tải xuống](/dl/).

Go 1.27 mang đến những cải tiến lớn trên ngôn ngữ, toolchain, runtime và thư viện chuẩn. Dưới đây là một số điểm nổi bật chính.

## Các thay đổi về ngôn ngữ

Go 1.27 giới thiệu ba cập nhật đáng chú ý cho [đặc tả ngôn ngữ](/doc/go1.27#language).

Đầu tiên, các phương thức generic hiện đã được hỗ trợ. Ví dụ, xem [`math/rand/v2.Rand`](/pkg/math/rand/v2#Rand):

```go
// Trước Go 1.27, phải thêm một phương thức riêng trên Rand cho từng kiểu
// (lược bỏ các phương thức cho số nguyên không dấu để ngắn gọn).
func (r *Rand) Int32N(n int32) int32
func (r *Rand) Int64N(n int64) int64
func (r *Rand) IntN(n int) int

// Go 1.27 thêm một phương thức generic mới hoạt động với mọi kiểu số nguyên.
func (r *Rand) N[Int intType](n Int) Int
```

Thứ hai, khóa trong một [struct literal](/ref/spec#Composite_literals) giờ đây có thể là bất kỳ [field selector](/ref/spec#Selectors) hợp lệ nào của kiểu struct, cho phép khởi tạo trực tiếp các trường trong các struct lồng nhau hoặc được nhúng:

```go
type Habitat struct {
    Burrow string
}

type Gopher struct {
    Name    string
    Habitat // Struct được nhúng.
}

// Go 1.27 cho phép sử dụng Burrow trực tiếp làm khóa.
g := Gopher{
    Name:   "Gopher",
    Burrow: "Burrow #42",
}
```

Cuối cùng, suy luận kiểu của hàm đã được tổng quát hóa để áp dụng trong mọi ngữ cảnh gán. Các hàm generic giờ đây có thể được sử dụng mà không cần đối số kiểu tường minh trong các struct literal, phép chuyển đổi kiểu và thao tác gửi qua channel:

```go
func GenericFormatter[T any](v T) string {
    return fmt.Sprintf("value: %v", v)
}

type IntFormatter func(int) string

// Go 1.27 suy luận T = int trong struct literal, phép chuyển đổi và thao tác gửi qua channel.
formatters := []IntFormatter{GenericFormatter}
fn := IntFormatter(GenericFormatter)
ch := make(chan IntFormatter, 1)
ch <- GenericFormatter
```

## Cải tiến công cụ

- [`go fix`](/doc/go1.27#go-fix) bao gồm một số [modernizer](/pkg/golang.org/x/tools/go/analysis/passes/modernize) mới:
  `atomictypes`, `embedlit`, `slicesbackward` và `unsafefuncs`.
- [`go doc`](/doc/go1.27#go-doc) hiện hỗ trợ các truy vấn `package@version` như `go doc example.com/pkg@v1.2.3`.
- [`go mod tidy`](/doc/go1.27#go-mod-tidy) hiện tự động hợp nhất nhiều khối `require` trong `go.mod` thành cấu trúc hai khối chuẩn gồm direct và indirect.

## Hiệu năng và runtime

- [Cấp phát bộ nhớ chuyên biệt theo kích thước](/doc/go1.27#faster-memory-allocation)

giảm chi phí cấp phát các đối tượng nhỏ (<80B) tới 30%, cải thiện hiệu năng

tổng thể khoảng ~1% đối với các chương trình có nhiều hoạt động cấp phát.
- Hồ sơ [`goroutineleak`](/doc/go1.27#goroutineleak-profiles) trong

[`runtime/pprof`](/pkg/runtime/pprof) hiện đã được cung cấp rộng rãi, cho phép

tự động phát hiện các goroutine bị chặn vĩnh viễn.

## Các bổ sung cho thư viện chuẩn

- [`encoding/json/v2`](/doc/go1.27#jsonv2) cung cấp khả năng xử lý JSON cấp cao

với các tùy chọn có thể cấu hình và các giá trị mặc định nghiêm ngặt hơn,

đồng thời có [`encoding/json/jsontext`](/doc/go1.27#jsonv2) cho việc truyền

dữ liệu streaming cấp thấp. Gói [`encoding/json`](/pkg/encoding/json) hiện có

được hỗ trợ bởi triển khai v2 để giải mã nhanh hơn trong khi vẫn duy trì khả

năng tương thích ngược.
- [`crypto/mldsa`](/doc/go1.27#crypto_mldsa) triển khai lược đồ chữ ký

ML-DSA hậu lượng tử (FIPS 204), được tích hợp vào

[`crypto/x509`](/pkg/crypto/x509) và [`crypto/tls`](/pkg/crypto/tls).
- [`uuid`](/doc/go1.27#uuid) cung cấp hỗ trợ gốc cho việc tạo và phân tích

UUID.
- [`simd`](/doc/go1.27#simd) và

[`simd/archsimd`](/doc/go1.27#archsimd) dành riêng cho từng kiến trúc cung cấp

hỗ trợ SIMD thử nghiệm.
- [`net/http/httptest`](/doc/go1.27#nethttphttptestpkgnethttphttptest) bổ sung

[`NewTestServer`](/pkg/net/http/httptest#NewTestServer), cung cấp một mạng

giả lập trong bộ nhớ phù hợp để sử dụng với gói

[`testing/synctest`](/pkg/testing/synctest).

Vui lòng đọc [ghi chú bản phát hành Go 1.27](/doc/go1.27) để xem danh sách đầy đủ

các thay đổi và chi tiết.

Trong vài tuần tới, các bài đăng blog tiếp theo sẽ trình bày chi tiết hơn một số

chủ đề liên quan đến Go 1.27. Hãy quay lại sau để đọc các bài viết đó.

Cảm ơn tất cả mọi người đã đóng góp cho bản phát hành này bằng cách viết mã,

báo lỗi, dùng thử các bổ sung thử nghiệm và kiểm thử các bản ứng viên phát hành.

Như thường lệ, nếu bạn nhận thấy bất kỳ vấn đề nào, vui lòng [báo cáo một issue](/issue/new).

Chúng tôi hy vọng bạn thích sử dụng Go 1.27!
