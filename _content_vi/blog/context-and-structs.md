---
title: Context và struct
date: 2021-02-24
by:
- Jean Barkhuysen, Matt T. Proud
tags:
- context
- cancellation
---

## Giới thiệu

Trong nhiều API Go, đặc biệt là các API hiện đại, đối số đầu tiên của hàm và phương thức thường là [`context.Context`](/pkg/context/). Context cung cấp một cách truyền deadline, việc caller hủy, và các giá trị gắn với phạm vi request qua các ranh giới API và giữa các tiến trình. Nó thường được dùng khi một thư viện tương tác, trực tiếp hoặc gián tiếp, với các máy chủ từ xa như cơ sở dữ liệu, API, và những thứ tương tự.

[Tài liệu về context](/pkg/context/) nêu rằng:

> Không nên lưu Context bên trong một kiểu struct; thay vào đó hãy truyền nó vào từng hàm cần dùng.

Bài viết này mở rộng lời khuyên đó bằng các lý do và ví dụ mô tả tại sao việc truyền Context lại quan trọng hơn việc lưu nó trong một kiểu khác. Nó cũng nhấn mạnh một trường hợp hiếm hoi mà việc lưu Context trong kiểu struct có thể hợp lý, và cách làm điều đó một cách an toàn.

## Ưu tiên truyền context như đối số

Để hiểu lời khuyên không lưu context trong struct, hãy xem cách tiếp cận ưu tiên là truyền context như một đối số:

```
// Worker fetches and adds works to a remote work orchestration server.
type Worker struct { /* … */ }

type Work struct { /* … */ }

func New() *Worker {
  return &Worker{}
}

func (w *Worker) Fetch(ctx context.Context) (*Work, error) {
  _ = ctx // A per-call ctx is used for cancellation, deadlines, and metadata.
}

func (w *Worker) Process(ctx context.Context, work *Work) error {
  _ = ctx // A per-call ctx is used for cancellation, deadlines, and metadata.
}
```

Ở đây, các phương thức `(*Worker).Fetch` và `(*Worker).Process` đều nhận context trực tiếp. Với thiết kế truyền như đối số này, người dùng có thể đặt deadline, việc hủy, và metadata cho từng lần gọi. Và cũng rõ ràng `context.Context` truyền vào mỗi phương thức sẽ được dùng như thế nào: không có kỳ vọng rằng một `context.Context` truyền vào một phương thức sẽ được dùng bởi phương thức khác. Điều này là vì context được giới hạn vào đúng mức nhỏ nhất mà thao tác cần, giúp tăng đáng kể tính hữu ích và độ rõ ràng của `context` trong package này.

## Lưu context trong struct dẫn tới nhầm lẫn

Hãy nhìn lại ví dụ `Worker` ở trên với cách tiếp cận kém được ưa chuộng hơn là để context trong struct. Vấn đề của cách này là khi bạn lưu context trong struct, bạn làm mờ vòng đời đối với caller, hoặc tệ hơn là trộn lẫn hai phạm vi lại với nhau theo những cách khó đoán:

```
type Worker struct {
  ctx context.Context
}

func New(ctx context.Context) *Worker {
  return &Worker{ctx: ctx}
}

func (w *Worker) Fetch() (*Work, error) {
  _ = w.ctx // A shared w.ctx is used for cancellation, deadlines, and metadata.
}

func (w *Worker) Process(work *Work) error {
  _ = w.ctx // A shared w.ctx is used for cancellation, deadlines, and metadata.
}
```

Các phương thức `(*Worker).Fetch` và `(*Worker).Process` đều dùng một context được lưu trong Worker. Điều này ngăn caller của Fetch và Process, vốn có thể tự mang các context khác nhau, chỉ định deadline, yêu cầu hủy, và gắn metadata cho từng lần gọi. Ví dụ: người dùng không thể cung cấp deadline chỉ cho `(*Worker).Fetch`, hoặc chỉ hủy lời gọi `(*Worker).Process`. Vòng đời của caller bị trộn lẫn với một context dùng chung, và context bị ràng buộc vào vòng đời khi `Worker` được tạo.

API này cũng gây nhầm lẫn cho người dùng hơn nhiều so với cách truyền như đối số. Người dùng có thể tự hỏi:

- Vì `New` nhận một `context.Context`, liệu constructor có đang làm việc gì cần hủy hoặc deadline không?
- `context.Context` được truyền vào `New` có áp dụng cho công việc trong `(*Worker).Fetch` và `(*Worker).Process` không? Không cái nào? Hay chỉ một trong hai?

API sẽ cần một lượng tài liệu khá lớn để nói tường minh cho người dùng biết chính xác `context.Context` được dùng vào việc gì. Người dùng thậm chí có thể phải đọc mã nguồn thay vì có thể tin cậy vào việc cấu trúc API tự truyền tải điều đó.

Và cuối cùng, sẽ khá nguy hiểm nếu thiết kế một máy chủ sản xuất mà các request của nó không có context riêng và do đó không thể tôn trọng việc hủy một cách đầy đủ. Nếu không có khả năng đặt deadline cho từng lần gọi, [tiến trình của bạn có thể bị backlog](https://sre.google/sre-book/handling-overload/) và cạn kiệt tài nguyên (như bộ nhớ)!

## Ngoại lệ của quy tắc: giữ tương thích ngược

Khi Go 1.7, phiên bản [giới thiệu context.Context](/doc/go1.7), được phát hành, rất nhiều API đã phải bổ sung hỗ trợ context theo cách tương thích ngược. Ví dụ, [các phương thức của `Client` trong `net/http`](/pkg/net/http/), như `Get` và `Do`, là những ứng viên rất phù hợp cho context. Mỗi request gửi ra ngoài bằng các phương thức này sẽ hưởng lợi từ deadline, khả năng hủy, và hỗ trợ metadata mà `context.Context` mang lại.

Có hai cách để thêm hỗ trợ `context.Context` theo cách tương thích ngược: đưa context vào một struct, như ta sắp thấy, hoặc nhân đôi hàm, trong đó bản sao nhận `context.Context` và có hậu tố `Context` trong tên hàm. Cách nhân đôi nên được ưu tiên hơn cách đưa context vào struct, và được bàn thêm trong [Keeping your modules compatible](/blog/module-compatibility). Tuy nhiên, trong một số trường hợp nó không thực tế: ví dụ nếu API của bạn lộ ra số lượng lớn hàm, thì việc nhân đôi tất cả chúng có thể không khả thi.

Package `net/http` đã chọn cách đưa context vào struct, tạo thành một nghiên cứu tình huống hữu ích. Hãy xem `Do` của `net/http`. Trước khi có `context.Context`, `Do` được định nghĩa như sau:

```
// Do sends an HTTP request and returns an HTTP response [...]
func (c *Client) Do(req *Request) (*Response, error)
```

Sau Go 1.7, `Do` có thể đã trông như sau, nếu không vì thực tế rằng nó sẽ phá vỡ tương thích ngược:

```
// Do sends an HTTP request and returns an HTTP response [...]
func (c *Client) Do(ctx context.Context, req *Request) (*Response, error)
```

Nhưng việc giữ tương thích ngược và tuân theo [cam kết tương thích Go 1](/doc/go1compat) là tối quan trọng đối với thư viện chuẩn. Vì vậy, thay vào đó, những người bảo trì đã chọn thêm `context.Context` vào struct `http.Request` để hỗ trợ `context.Context` mà không phá vỡ tương thích ngược:

```
// A Request represents an HTTP request received by a server or to be sent by a client.
// ...
type Request struct {
  ctx context.Context

  // ...
}

// NewRequestWithContext returns a new Request given a method, URL, and optional
// body.
// [...]
// The given ctx is used for the lifetime of the Request.
func NewRequestWithContext(ctx context.Context, method, url string, body io.Reader) (*Request, error) {
  // Simplified for brevity of this article.
  return &Request{
    ctx: ctx,
    // ...
  }
}

// Do sends an HTTP request and returns an HTTP response [...]
func (c *Client) Do(req *Request) (*Response, error)
```

Khi cải tạo API của bạn để hỗ trợ context, có thể sẽ hợp lý khi thêm `context.Context` vào một struct như trên. Tuy nhiên, hãy nhớ cân nhắc trước tiên việc nhân đôi các hàm của bạn, điều này cho phép cải tạo `context.Context` theo cách tương thích ngược mà không hy sinh tính hữu ích và tính dễ hiểu. Ví dụ:

```
// Call uses context.Background internally; to specify the context, use
// CallContext.
func (c *Client) Call() error {
  return c.CallContext(context.Background())
}

func (c *Client) CallContext(ctx context.Context) error {
  // ...
}
```

## Kết luận

Context giúp việc lan truyền các thông tin quan trọng xuyên thư viện và xuyên API xuống ngăn xếp lời gọi trở nên dễ dàng. Nhưng nó phải được dùng một cách nhất quán và rõ ràng thì mới còn dễ hiểu, dễ gỡ lỗi và hiệu quả.

Khi được truyền như đối số đầu tiên của một phương thức thay vì được lưu trong một kiểu struct, người dùng có thể tận dụng đầy đủ khả năng mở rộng của nó để xây dựng một cây mạnh mẽ gồm thông tin hủy, deadline và metadata xuyên suốt ngăn xếp lời gọi. Và tuyệt nhất là phạm vi của nó được hiểu rất rõ khi nó được truyền như một đối số, dẫn tới khả năng hiểu và gỡ lỗi tốt ở mọi tầng của ngăn xếp.

Khi thiết kế API có context, hãy nhớ lời khuyên: truyền `context.Context` như một đối số; đừng lưu nó trong struct.

## Đọc thêm

- [Go Concurrency Patterns: Context (bài blog năm 2014)](context.md)

