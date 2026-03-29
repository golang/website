---
title: "Các mẫu đồng thời trong Go: Context"
date: 2014-07-29
by:
- Sameer Ajmani
tags:
- concurrency
- cancellation
- context
summary: Giới thiệu về package context của Go.
template: true
---

## Giới thiệu

Trong các máy chủ Go, mỗi request đến được xử lý trong goroutine riêng của nó.
Các handler thường khởi động thêm goroutine để truy cập các backend như cơ sở dữ liệu và dịch vụ RPC.
Tập goroutine làm việc trên một request thường cần truy cập các giá trị dành riêng cho request như danh tính của người dùng cuối, token ủy quyền và deadline của request.
Khi một request bị hủy hoặc quá thời gian, toàn bộ goroutine đang làm việc trên request đó phải thoát nhanh để hệ thống có thể thu hồi mọi tài nguyên mà chúng đang dùng.

Tại Google, chúng tôi đã phát triển package `context` giúp dễ dàng truyền các giá trị theo phạm vi request, tín hiệu hủy và deadline qua các ranh giới API tới mọi goroutine tham gia xử lý request.
Package này được công khai tại [context](/pkg/context).
Bài viết này mô tả cách sử dụng package đó và cung cấp một ví dụ chạy hoàn chỉnh.

## Context

Cốt lõi của package `context` là kiểu `Context`:

{{code "context/interface.go" `/A Context/` `/^}/`}}

(Mô tả này đã được rút gọn; [godoc](/pkg/context) là nguồn chính xác nhất.)

Phương thức `Done` trả về một channel đóng vai trò là tín hiệu hủy cho các hàm chạy thay mặt cho `Context`: khi channel bị đóng, các hàm nên từ bỏ công việc và trả về.
Phương thức `Err` trả về lỗi cho biết vì sao `Context` bị hủy.
Bài viết [Pipelines and Cancellation](/blog/pipelines) thảo luận chi tiết hơn về thành ngữ dùng channel `Done`.

`Context` _không_ có phương thức `Cancel` vì cùng một lý do mà channel `Done` là chỉ-nhận: hàm nhận tín hiệu hủy thường không phải là hàm phát tín hiệu.
Đặc biệt, khi một thao tác cha khởi động các goroutine cho các thao tác con, các thao tác con đó không nên có khả năng hủy thao tác cha.
Thay vào đó, hàm `WithCancel` (mô tả bên dưới) cung cấp cách hủy một giá trị `Context` mới.

`Context` an toàn khi được dùng đồng thời bởi nhiều goroutine.
Mã có thể truyền một `Context` duy nhất vào bất kỳ số lượng goroutine nào và hủy `Context` đó để phát tín hiệu cho tất cả.

Phương thức `Deadline` cho phép các hàm xác định liệu chúng có nên bắt đầu công việc ngay từ đầu hay không; nếu thời gian còn lại quá ít, công việc đó có thể không đáng làm.
Mã cũng có thể dùng deadline để đặt timeout cho các thao tác I/O.

`Value` cho phép `Context` mang dữ liệu gắn với request.
Dữ liệu đó phải an toàn cho việc dùng đồng thời bởi nhiều goroutine.

### Context phát sinh

Package `context` cung cấp các hàm để _phát sinh_ các giá trị `Context` mới từ các context hiện có.
Các giá trị này tạo thành một cây: khi một `Context` bị hủy, mọi `Context` phát sinh từ nó cũng bị hủy theo.

`Background` là gốc của mọi cây `Context`; nó không bao giờ bị hủy:

{{code "context/interface.go" `/Background returns/` `/func Background/`}}

`WithCancel` và `WithTimeout` trả về các giá trị `Context` phát sinh có thể bị hủy sớm hơn `Context` cha.
`Context` gắn với một request đến thường bị hủy khi handler của request trả về.
`WithCancel` cũng hữu ích để hủy các request dư thừa khi dùng nhiều bản sao.
`WithTimeout` hữu ích để đặt deadline cho các request tới các máy chủ backend:

{{code "context/interface.go" `/WithCancel/` `/func WithTimeout/`}}

`WithValue` cung cấp cách gắn các giá trị gắn với request vào một `Context`:

{{code "context/interface.go" `/WithValue/` `/func WithValue/`}}

Cách tốt nhất để thấy package `context` được dùng thế nào là thông qua một ví dụ hoàn chỉnh.

## Ví dụ: Tìm kiếm web Google

Ví dụ của chúng ta là một máy chủ HTTP xử lý các URL như `/search?q=golang&timeout=1s` bằng cách chuyển tiếp truy vấn "golang" tới [Google Web Search API](https://developers.google.com/web-search/docs/) và hiển thị kết quả.
Tham số `timeout` cho máy chủ biết rằng nó phải hủy request sau khi khoảng thời gian đó trôi qua.

Mã được chia thành ba package:

  - [server](context/server/server.go) cung cấp hàm `main` và handler cho `/search`.
  - [userip](context/userip/userip.go) cung cấp các hàm để trích xuất địa chỉ IP của người dùng từ request và gắn nó với một `Context`.
  - [google](context/google/google.go) cung cấp hàm `Search` để gửi truy vấn tới Google.

### Chương trình server

Chương trình [server](context/server/server.go) xử lý các request như `/search?q=golang` bằng cách trả về một vài kết quả tìm kiếm đầu tiên của Google cho `golang`.
Nó đăng ký `handleSearch` để xử lý endpoint `/search`.
Handler tạo một `Context` ban đầu tên là `ctx` và sắp xếp để nó bị hủy khi handler trả về.
Nếu request bao gồm tham số URL `timeout`, `Context` sẽ bị hủy tự động khi timeout hết hạn:

{{code "context/server/server.go" `/func handleSearch/` `/defer cancel/`}}

Handler trích truy vấn từ request và trích địa chỉ IP của client bằng cách gọi package `userip`.
Địa chỉ IP của client là cần thiết cho các request backend, vì vậy `handleSearch` gắn nó vào `ctx`:

{{code "context/server/server.go" `/Check the search query/` `/userip.NewContext/`}}

Handler gọi `google.Search` với `ctx` và `query`:

{{code "context/server/server.go" `/Run the Google search/` `/elapsed/`}}

Nếu tìm kiếm thành công, handler sẽ hiển thị kết quả:

{{code "context/server/server.go" `/resultsTemplate/` `/(?m)}$/`}}

### Package userip

Package [userip](context/userip/userip.go) cung cấp các hàm để trích xuất địa chỉ IP của người dùng từ một request và gắn nó với một `Context`.
`Context` cung cấp một ánh xạ khóa-giá trị, trong đó cả khóa lẫn giá trị đều có kiểu `interface{}`.
Kiểu khóa phải hỗ trợ so sánh bằng, và giá trị phải an toàn cho việc dùng đồng thời bởi nhiều goroutine.
Những package như `userip` che giấu chi tiết của ánh xạ này và cung cấp cách truy cập được định kiểu mạnh vào một giá trị `Context` cụ thể.

Để tránh va chạm khóa, `userip` định nghĩa một kiểu không export tên là `key` và dùng một giá trị của kiểu này làm khóa context:

{{code "context/userip/userip.go" `/The key type/` `/const userIPKey/`}}

`FromRequest` trích xuất một giá trị `userIP` từ `http.Request`:

{{code "context/userip/userip.go" `/func FromRequest/` `/}/`}}

`NewContext` trả về một `Context` mới mang theo giá trị `userIP` được cung cấp:

{{code "context/userip/userip.go" `/func NewContext/` `/}/`}}

`FromContext` trích xuất một `userIP` từ `Context`:

{{code "context/userip/userip.go" `/func FromContext/` `/}/`}}

### Package google

Hàm [google.Search](context/google/google.go) thực hiện một request HTTP tới [Google Web Search API](https://developers.google.com/web-search/docs/) và phân tích kết quả mã hóa JSON.
Nó nhận tham số `Context` tên là `ctx` và trả về ngay nếu `ctx.Done` bị đóng trong lúc request đang bay.

Request tới Google Web Search API bao gồm truy vấn tìm kiếm và IP của người dùng dưới dạng các tham số truy vấn:

{{code "context/google/google.go" `/func Search/` `/q.Encode/`}}

`Search` dùng một hàm phụ trợ tên `httpDo` để phát request HTTP và hủy nó nếu `ctx.Done` bị đóng trong lúc request hoặc response đang được xử lý.
`Search` truyền một closure vào `httpDo` để xử lý response HTTP:

{{code "context/google/google.go" `/var results/` `/return results/`}}

Hàm `httpDo` chạy request HTTP và xử lý response của nó trong một goroutine mới.
Nó hủy request nếu `ctx.Done` bị đóng trước khi goroutine kết thúc:

{{code "context/google/google.go" `/func httpDo/` `/^}/`}}

## Điều chỉnh mã cho Context

Nhiều framework máy chủ cung cấp package và kiểu để mang các giá trị theo phạm vi request.
Chúng ta có thể định nghĩa các hiện thực mới của interface `Context` để làm cầu nối giữa mã dùng các framework hiện có và mã kỳ vọng một tham số `Context`.

Ví dụ, package [github.com/gorilla/context](http://www.gorillatoolkit.org/pkg/context) của Gorilla cho phép handler gắn dữ liệu với các request đến bằng cách cung cấp một ánh xạ từ request HTTP sang cặp khóa-giá trị.
Trong [gorilla.go](context/gorilla/gorilla.go), chúng tôi cung cấp một hiện thực `Context` mà phương thức `Value` trả về các giá trị gắn với một request HTTP cụ thể trong package Gorilla.

Các package khác cũng đã cung cấp khả năng hủy tương tự `Context`.
Ví dụ, [Tomb](https://godoc.org/gopkg.in/tomb.v2) cung cấp phương thức `Kill` phát tín hiệu hủy bằng cách đóng một channel `Dying`.
`Tomb` cũng cung cấp các phương thức để chờ các goroutine đó thoát, tương tự `sync.WaitGroup`.
Trong [tomb.go](context/tomb/tomb.go), chúng tôi cung cấp một hiện thực `Context` bị hủy khi `Context` cha của nó bị hủy hoặc khi `Tomb` được cung cấp bị kill.

## Kết luận

Tại Google, chúng tôi yêu cầu lập trình viên Go truyền một tham số `Context` làm đối số đầu tiên cho mọi hàm trên đường đi lời gọi giữa request đi vào và request đi ra.
Điều này cho phép mã Go do nhiều đội khác nhau phát triển tương tác tốt với nhau.
Nó cung cấp khả năng điều khiển đơn giản đối với timeout và việc hủy, đồng thời bảo đảm rằng những giá trị quan trọng như thông tin xác thực bảo mật được truyền đúng cách qua các chương trình Go.

Những framework máy chủ muốn xây dựng dựa trên `Context` nên cung cấp các hiện thực của `Context` để làm cầu nối giữa package của chính chúng và những package kỳ vọng một tham số `Context`.
Các thư viện client của chúng khi đó sẽ nhận một `Context` từ mã gọi.
Bằng cách thiết lập một giao diện chung cho dữ liệu theo phạm vi request và việc hủy, `Context` giúp các nhà phát triển package dễ dàng chia sẻ mã nhằm tạo ra các dịch vụ có khả năng mở rộng.

## Đọc thêm

- [Go Concurrency Patterns: Pipelines and cancellation (bài blog năm 2014)](pipelines.md)
- [Contexts and structs (bài blog năm 2021)](context-and-structs.md)

