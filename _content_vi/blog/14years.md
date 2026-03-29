---
title: Mười bốn năm của Go
date: 2023-11-10
by:
- Russ Cox, thay mặt nhóm Go
summary: Chúc mừng sinh nhật, Go!
---

<img src="/doc/gopher/gopherdrink.png" height="219" width="223" align="right" style="margin: 0 0 1em 1em;">

Hôm nay chúng ta kỷ niệm sinh nhật lần thứ mười bốn của bản phát hành mã nguồn mở Go!
Go đã có một năm tuyệt vời, với hai bản phát hành đầy tính năng và những cột mốc quan trọng khác.

Chúng tôi đã phát hành [Go 1.20 vào tháng 2](/blog/go1.20)
và [Go 1.21 vào tháng 8](/blog/go1.21),
tập trung nhiều hơn vào các cải tiến trong triển khai
so với các thay đổi mới của ngôn ngữ.

Tối ưu hóa theo hồ sơ thực thi (PGO),
[được giới thiệu trước trong Go 1.20](/blog/pgo-preview)
và
[được phát hành trong Go 1.21](/blog/pgo),
cho phép trình biên dịch Go đọc hồ sơ thực thi của chương trình
và sau đó dành nhiều thời gian hơn để tối ưu
những phần của chương trình chạy thường xuyên nhất.
Trong Go 1.21, các tải công việc thường nhận được
mức cải thiện sử dụng CPU từ 2% đến 7% khi bật PGO.
Hãy xem “[Tối ưu hóa theo hồ sơ thực thi trong Go 1.21](/blog/pgo)” để có cái nhìn tổng quan
và [hướng dẫn sử dụng tối ưu hóa theo hồ sơ thực thi](/doc/pgo)
để có tài liệu đầy đủ.

Go đã hỗ trợ thu thập hồ sơ độ bao phủ trong `go test`
[kể từ Go 1.2](/blog/cover).
Go 1.20 bổ sung hỗ trợ thu thập hồ sơ độ bao phủ trong các tệp nhị phân
được xây dựng bởi `go build`,
cho phép bạn thu thập độ bao phủ cả trong các bài kiểm thử tích hợp lớn hơn.
Xem “[Độ bao phủ mã cho các bài kiểm thử tích hợp Go](/blog/integration-test-coverage)” để biết chi tiết.

Tính tương thích đã là một phần quan trọng của Go kể từ
“[Go 1 và Tương lai của các chương trình Go](/doc/go1compat)”.
Go 1.21 tiếp tục cải thiện tính tương thích
bằng cách mở rộng các quy ước sử dụng GODEBUG
trong những tình huống chúng tôi cần thực hiện một thay đổi,
chẳng hạn như một bản sửa lỗi quan trọng,
mà thay đổi đó phải được cho phép nhưng vẫn có thể làm hỏng các chương trình hiện có.
Hãy xem bài blog
“[Tính tương thích ngược, Go 1.21 và Go 2](/blog/compat)”
để có cái nhìn tổng quan và
tài liệu
“[Go, Tính tương thích ngược và GODEBUG](/doc/godebug)” để biết chi tiết.

Go 1.21 cũng phát hành hỗ trợ quản lý toolchain tích hợp sẵn,
cho phép bạn thay đổi phiên bản của
bộ công cụ Go mà bạn sử dụng trong một module cụ thể
dễ dàng như khi thay đổi phiên bản của các dependency khác.
Hãy xem bài blog
“[Tính tương thích xuôi và Quản lý Toolchain trong Go 1.21](/blog/toolchain)”
để có cái nhìn tổng quan và tài liệu
“[Go Toolchains](/doc/toolchain)”
để biết chi tiết.

Một thành tựu quan trọng khác về công cụ là
việc tích hợp các chỉ mục trên đĩa vào
gopls, máy chủ LSP của Go.
Điều này đã giảm độ trễ khởi động và mức sử dụng bộ nhớ của gopls xuống 3-5 lần
trong các trường hợp sử dụng điển hình.
“[Mở rộng gopls cho hệ sinh thái Go đang phát triển](/blog/gopls-scalability)”
giải thích các chi tiết kỹ thuật.
Bạn có thể bảo đảm mình đang chạy gopls mới nhất bằng cách chạy:

```
go install golang.org/x/tools/gopls@latest
```

Go 1.21 giới thiệu các gói mới
[cmp](/pkg/cmp/),
[maps](/pkg/maps/),
và
[slices](/pkg/slices/)
đây là những thư viện chuẩn generic đầu tiên của Go,
đồng thời cũng mở rộng tập các kiểu có thể so sánh được.
Để biết chi tiết về điều đó, hãy xem bài blog
“[Tất cả các kiểu comparable của bạn](/blog/comparable)”.

Nhìn chung, chúng tôi tiếp tục tinh chỉnh generics
và viết các bài nói chuyện cùng bài blog giải thích
những chi tiết quan trọng.
Hai bài viết nổi bật trong năm nay là
“[Giải cấu trúc tham số kiểu](/blog/deconstructing-type-parameters)”,
và
“[Mọi điều bạn luôn muốn biết về suy luận kiểu - và hơn thế nữa một chút](/blog/type-inference)”.

Một gói mới quan trọng khác trong Go 1.21 là
[log/slog](/pkg/log/slog/),
gói này bổ sung một API chính thức cho
ghi log có cấu trúc vào thư viện chuẩn.
Xem “[Ghi log có cấu trúc với slog](/blog/slog)” để có cái nhìn tổng quan.

Đối với cổng WebAssembly (Wasm), Go 1.21 phát hành hỗ trợ
cho việc chạy trên WebAssembly System Interface (WASI) preview 1.
WASI preview 1 là một giao diện “hệ điều hành” mới dành cho Wasm
được hầu hết các môi trường Wasm phía máy chủ hỗ trợ.
Xem “[Hỗ trợ WASI trong Go](/blog/wasi)” để xem hướng dẫn chi tiết.

Về phía bảo mật, chúng tôi tiếp tục bảo đảm rằng
Go dẫn đầu trong việc giúp các nhà phát triển hiểu rõ
dependency và lỗ hổng của họ,
với [Govulncheck 1.0 ra mắt vào tháng 7](/blog/govulncheck).
Nếu bạn dùng VS Code, bạn có thể chạy govulncheck trực tiếp trong
trình soạn thảo bằng tiện ích mở rộng Go:
xem [hướng dẫn này](/doc/tutorial/govulncheck-ide) để bắt đầu.
Và nếu bạn dùng GitHub, bạn có thể chạy govulncheck như một phần của
CI/CD, với
[GitHub Action cho govulncheck](https://github.com/marketplace/actions/golang-govulncheck-action).
Để tìm hiểu thêm về việc kiểm tra dependency nhằm phát hiện vấn đề lỗ hổng,
xem bài nói chuyện tại Google I/O năm nay,
“[Xây dựng ứng dụng an toàn hơn với Go và Google](https://www.youtube.com/watch?v=HSt6FhsPT8c&ab_channel=TheGoProgrammingLanguage)”.

Một cột mốc bảo mật quan trọng khác là
các bản dựng toolchain có tính tái tạo rất cao trong Go 1.21.
Xem “[Các Go Toolchain được xác minh, tái tạo hoàn hảo](/blog/rebuild)” để biết chi tiết,
bao gồm cả phần trình diễn việc tái tạo một Go toolchain Ubuntu Linux
trên máy Mac mà hoàn toàn không dùng bất kỳ công cụ Linux nào.

Đó là một năm bận rộn!

Trong năm thứ 15 của Go, chúng tôi sẽ tiếp tục làm việc để biến Go thành môi trường tốt nhất
cho kỹ nghệ phần mềm ở quy mô lớn.
Một thay đổi mà chúng tôi đặc biệt hào hứng là
định nghĩa lại ngữ nghĩa `:=` trong vòng lặp `for` để loại bỏ
khả năng phát sinh lỗi aliasing ngoài ý muốn.
Xem “[Sửa các vòng lặp `for` trong Go 1.22](/blog/loopvar-preview)”
để biết chi tiết,
bao gồm hướng dẫn xem trước thay đổi này trong Go 1.21.

## Xin cảm ơn!

Dự án Go từ lâu đã luôn lớn hơn rất nhiều so với chỉ riêng chúng tôi trong nhóm Go tại Google.
Xin cảm ơn tất cả những người đóng góp và mọi người trong cộng đồng Go
đã giúp Go trở thành như ngày hôm nay.
Chúng tôi chúc tất cả các bạn những điều tốt đẹp nhất trong năm tới.
