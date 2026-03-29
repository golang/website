---
title: Mười hai năm của Go
date: 2021-11-10
by:
- Russ Cox, thay mặt nhóm Go
summary: Chúc mừng sinh nhật, Go!
---


Hôm nay chúng ta kỷ niệm sinh nhật lần thứ mười hai của bản phát hành mã nguồn mở Go.
Đây là một năm nhiều biến động và chúng tôi có rất nhiều điều để mong đợi trong năm tới.

Thay đổi dễ thấy nhất tại đây trên blog là
[ngôi nhà mới của chúng tôi trên go.dev](/blog/tidy-web),
là một phần trong nỗ lực hợp nhất tất cả các trang web về Go thành một trang thống nhất, mạch lạc.
Một phần khác của quá trình hợp nhất đó là
[thay thế godoc.org bằng pkg.go.dev](/blog/godoc.org-redirect).

Vào tháng 2, [bản phát hành Go 1.16](/blog/go1.16)
đã bổ sung [hỗ trợ macOS ARM64](/blog/ports),
thêm [giao diện hệ thống tệp](/pkg/io/fs) và [tệp nhúng](/pkg/embed),
đồng thời [bật modules theo mặc định](/blog/go116-module-changes),
cùng với vô số cải tiến và tối ưu hóa quen thuộc.

Vào tháng 8, [bản phát hành Go 1.17](/blog/go1.17)
đã bổ sung hỗ trợ Windows ARM64,
giúp [việc đưa ra quyết định về bộ mật mã TLS trở nên dễ dàng và an toàn hơn](/blog/tls-cipher-suites),
giới thiệu [đồ thị module đã được lược gọn](/doc/go1.17#go-command)
để modules còn hiệu quả hơn nữa trong các dự án lớn,
và bổ sung
[cú pháp ràng buộc build mới, dễ đọc hơn](https://pkg.go.dev/cmd/go#hdr-Build_constraints).
Ở tầng bên dưới, Go 1.17 cũng chuyển sang quy ước gọi hàm dựa trên thanh ghi cho các hàm Go
trên x86-64, cải thiện hiệu năng trong các ứng dụng nặng về CPU thêm 5-15%.

Trong suốt năm qua, chúng tôi đã
xuất bản [nhiều hướng dẫn mới](/doc/tutorial/),
một [hướng dẫn về cơ sở dữ liệu trong Go](/doc/database/),
một [hướng dẫn phát triển module](/doc/#developing-modules),
và một [tài liệu tham chiếu về Go modules](/ref/mod).
Một điểm nhấn là hướng dẫn mới
“[Xây dựng RESTful API với Go và Gin](/doc/tutorial/web-service-gin)”,
cũng có sẵn ở
[dạng tương tác bằng Google Cloud Shell](/s/cloud-shell-web-tutorial).

Chúng tôi cũng rất bận rộn ở phía IDE,
[bật gopls theo mặc định trong VS Code Go](/blog/gopls-vscode-go)
và mang đến vô số cải tiến cho cả `gopls` lẫn VS Code Go,
bao gồm [trải nghiệm gỡ lỗi mạnh mẽ](https://github.com/golang/vscode-go/blob/master/docs/debugging.md)
được hỗ trợ bởi Delve.

Chúng tôi cũng ra mắt [bản beta Go fuzzing](/blog/fuzz-beta)
và [chính thức đề xuất bổ sung generics vào Go](/blog/generics-proposal),
cả hai hiện đều được kỳ vọng sẽ xuất hiện trong Go 1.18.

Tiếp tục thích nghi với cách tiếp cận “ưu tiên trực tuyến”, nhóm Go đã tổ chức sự kiện thường niên lần thứ hai của mình:
[Ngày Go tại Google Open Source Live](https://opensourcelive.withgoogle.com/events/go-day-2021).
Bạn có thể xem các bài nói chuyện trên YouTube:

- “[Using Generics in Go](https://www.youtube.com/watch?v=nr8EpUO9jhw)”,
  của Ian Lance Taylor, giới thiệu generics và cách sử dụng chúng hiệu quả.

- “[Modern Enterprise Applications](https://www.youtube.com/watch?v=5fgG1qZaV4w)”,
  của Steve Francia, cho thấy Go đóng vai trò như thế nào trong quá trình hiện đại hóa doanh nghiệp.

- “[Building Better Projects with the Go Editor](https://www.youtube.com/watch?v=jMyzsp2E_0U)”,
  của Suzy Mueller, trình bày cách bộ công cụ tích hợp của VS Code Go
  giúp bạn điều hướng mã, gỡ lỗi kiểm thử, và nhiều hơn thế nữa.

- “[From Proof of Concept to Production](https://www.youtube.com/watch?v=e7PtBOsTpXE)”,
  của Benjamin Cane, Distinguished Engineer tại American Express,
  giải thích cách American Express bắt đầu sử dụng Go cho các nền tảng thanh toán và phần thưởng của họ.

## Hướng về phía trước

Chúng tôi vô cùng hào hứng với những gì đang chờ Go trong năm thứ 13.
Vào tháng tới, chúng tôi sẽ có hai bài nói chuyện tại [GopherCon 2021](https://www.gophercon.com/),
cùng với [rất nhiều diễn giả tài năng đến từ khắp cộng đồng Go](https://www.gophercon.com/agenda).
Hãy đăng ký miễn phí và ghi nhớ lịch nhé!

- “Why and How to Use Go Generics”,
  của Robert Griesemer và Ian Lance Taylor,
  những người đã dẫn dắt việc thiết kế và hiện thực tính năng mới này. \
  [Ngày 8 tháng 12, 11:00 sáng (miền Đông Hoa Kỳ)](https://www.gophercon.com/agenda/session/593015).

- “Debugging Go Code Using the Debug Adapter Protocol (DAP)”,
  của Suzy Mueller,
  trình bày cách sử dụng các tính năng gỡ lỗi nâng cao của VS Code Go với Delve. \
  [Ngày 9 tháng 12, 3:20 chiều (miền Đông Hoa Kỳ)](https://www.gophercon.com/agenda/session/593029).

Vào tháng 2, bản phát hành Go 1.18 sẽ mở rộng
quy ước gọi hàm mới dựa trên thanh ghi sang các kiến trúc không phải x86,
mang theo những cải tiến hiệu năng đáng kể.
Nó sẽ bao gồm hỗ trợ Go fuzzing mới.
Và đó sẽ là bản phát hành đầu tiên có hỗ trợ generics.

Generics sẽ là một trong những trọng tâm của chúng tôi trong năm 2022.
Bản phát hành đầu tiên trong Go 1.18 mới chỉ là sự khởi đầu.
Chúng tôi cần dành thời gian sử dụng generics và tìm hiểu điều gì hiệu quả
và điều gì không hiệu quả, để có thể viết ra các thực hành tốt nhất
và quyết định điều gì nên được bổ sung vào thư viện chuẩn
cũng như các thư viện khác.
Chúng tôi kỳ vọng rằng Go 1.19 (dự kiến vào tháng 8 năm 2022)
và các bản phát hành sau đó sẽ tiếp tục tinh chỉnh thiết kế và hiện thực của
generics cũng như tích hợp chúng sâu hơn vào toàn bộ trải nghiệm Go.

Một trọng tâm khác của năm 2022 là bảo mật chuỗi cung ứng.
Chúng tôi đã nói suốt nhiều năm về
[các vấn đề của dependency](https://research.swtch.com/deps).
Thiết kế của Go modules mang lại
[các bản build có thể tái tạo, có thể kiểm chứng và đã được kiểm chứng](https://research.swtch.com/vgo-repro),
nhưng vẫn còn nhiều việc phải làm.
Bắt đầu từ Go 1.18, lệnh `go` sẽ nhúng thêm nhiều thông tin vào tệp nhị phân
về cấu hình build của chúng, vừa để giúp việc tái tạo dễ hơn
vừa để hỗ trợ các dự án cần
[tạo SBOM](https://en.wikipedia.org/wiki/Software_bill_of_materials) cho các tệp nhị phân Go.
Chúng tôi cũng đã bắt đầu triển khai một
[cơ sở dữ liệu lỗ hổng bảo mật của Go](https://pkg.go.dev/golang.org/x/vuln)
và một công cụ liên quan để báo cáo lỗ hổng trong các dependency của chương trình.
Một trong những mục tiêu của công việc này là cải thiện đáng kể tỷ lệ tín hiệu trên nhiễu
của loại công cụ này:
nếu một chương trình không sử dụng hàm có lỗ hổng, chúng tôi không muốn báo cáo điều đó.
Trong suốt năm 2022, chúng tôi dự định cung cấp công cụ này dưới dạng công cụ độc lập
đồng thời tích hợp nó vào các công cụ hiện có, bao gồm `gopls`, VS Code Go, và [pkg.go.dev](https://pkg.go.dev).
Vẫn còn nhiều việc khác phải làm để cải thiện các khía cạnh khác trong tư thế bảo mật chuỗi cung ứng của Go.
Hãy chờ những chi tiết tiếp theo.

Nhìn chung, chúng tôi kỳ vọng năm 2022 sẽ là một năm nhiều biến động với Go,
và chúng tôi sẽ tiếp tục mang đến các bản phát hành và cải tiến đúng hạn
mà bạn vẫn luôn mong đợi.

## Xin cảm ơn!

Go còn hơn nhiều so với chỉ riêng chúng tôi, nhóm Go tại Google.
Xin cảm ơn sự giúp đỡ của các bạn trong việc đưa Go đến thành công
và đồng hành cùng chúng tôi trong cuộc phiêu lưu này.
Chúng tôi hy vọng tất cả các bạn luôn an toàn và nhận được những điều tốt đẹp nhất.
