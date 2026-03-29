---
title: Mười ba năm của Go
date: 2022-11-10
by:
- Russ Cox, thay mặt nhóm Go
summary: Chúc mừng sinh nhật, Go!
---

<img src="../doc/gopher/gopherbelly300.jpg" height="300" width="300" align="right" style="margin: 0 0 1em 1em;">

Hôm nay chúng ta kỷ niệm sinh nhật lần thứ mười ba của bản phát hành mã nguồn mở Go.
[Gopher](/doc/gopher) giờ đã là một thiếu niên!

Đây là một năm nhiều sự kiện đối với Go.
Sự kiện quan trọng nhất là
[bản phát hành Go 1.18 vào tháng 3](/blog/go1.18),
mang đến nhiều cải tiến nhưng nổi bật nhất là
workspaces của Go, fuzzing và generics.

Workspaces giúp dễ dàng làm việc với nhiều module cùng lúc,
điều này đặc biệt hữu ích khi bạn đang duy trì một nhóm module có liên quan với nhau
và có dependency giữa các module đó.
Để tìm hiểu về workspaces, hãy xem bài blog của Beth Brown
“[Làm quen với workspaces](/blog/get-familiar-with-workspaces)”
và [tài liệu tham chiếu về workspaces](/ref/mod#workspaces).

Fuzzing là một tính năng mới của `go test`
giúp bạn tìm ra những đầu vào mà mã của bạn xử lý chưa đúng:
bạn định nghĩa một fuzz test phải thành công với mọi đầu vào,
rồi fuzzing sẽ thử nhiều đầu vào ngẫu nhiên khác nhau, được dẫn hướng bởi độ bao phủ mã,
để tìm cách làm cho fuzz test thất bại.
Fuzzing đặc biệt hữu ích khi phát triển những đoạn mã cần phải
vững chắc trước các đầu vào tùy ý (thậm chí do kẻ tấn công kiểm soát).
Để tìm hiểu thêm về fuzzing, hãy xem hướng dẫn
“[Bắt đầu với fuzzing](/doc/tutorial/fuzz)”
và [tài liệu tham chiếu về fuzzing](/security/fuzz/),
đồng thời hãy để ý bài nói chuyện tại GopherCon 2022 của Katie Hockman
“Fuzz Testing Made Easy”,
có lẽ sẽ sớm được đưa lên mạng.

Generics, có lẽ là tính năng được mong chờ nhất của Go,
bổ sung tính đa hình tham số cho Go, cho phép viết
mã hoạt động với nhiều kiểu dữ liệu khác nhau nhưng vẫn
được kiểm tra tĩnh tại thời điểm biên dịch.
Để tìm hiểu thêm về generics, hãy xem hướng dẫn
“[Bắt đầu với generics](/doc/tutorial/generics)”.
Để xem chi tiết hơn, hãy đọc
các bài blog
“[Giới thiệu về Generics](/blog/intro-generics)”
và
“[Khi nào nên dùng Generics](/blog/when-generics)”,
hoặc các bài nói chuyện
“[Using Generics in Go](https://www.youtube.com/watch?v=nr8EpUO9jhw)”
từ Go Day on Google Open Source Live 2021,
và
“[Generics!](https://www.youtube.com/watch?v=Pa_e9EeCdy8)” từ GopherCon 2021,
của Robert Griesemer và Ian Lance Taylor.

So với Go 1.18, [bản phát hành Go 1.19 vào tháng 8](/blog/go1.19) tương đối yên ắng:
nó tập trung vào việc tinh chỉnh và cải thiện những tính năng mà Go 1.18 đã giới thiệu,
cũng như cải thiện độ ổn định nội bộ và tối ưu hóa.
Một thay đổi dễ thấy trong Go 1.19 là việc bổ sung
hỗ trợ cho [liên kết, danh sách và tiêu đề trong chú thích tài liệu Go](/doc/comment).
Một thay đổi khác là việc bổ sung [giới hạn bộ nhớ mềm](/doc/go1.19#runtime)
cho bộ gom rác, đặc biệt hữu ích trong các tải công việc chạy trong container.
Để biết thêm về những cải tiến gần đây của bộ gom rác,
xem bài blog của Michael Knyszek “[Go runtime: 4 years later](/blog/go119runtime)”,
bài nói chuyện của anh ấy “[Respecting Memory Limits in Go](https://www.youtube.com/watch?v=07wduWyWx8M&list=PLtoVuM73AmsJjj5tnZ7BodjN_zIvpULSx)”,
và tài liệu mới “[Hướng dẫn về Go Garbage Collector](/doc/gc-guide)”.

Chúng tôi vẫn tiếp tục làm việc để giúp việc phát triển với Go mở rộng một cách trơn tru
trên những codebase ngày càng lớn,
đặc biệt là trong công việc của chúng tôi với VS Code Go và language server Gopls.
Trong năm nay, các bản phát hành Gopls tập trung vào việc cải thiện độ ổn định và hiệu năng,
đồng thời hỗ trợ generics cũng như các phân tích và code lens mới.
Nếu bạn vẫn chưa dùng VS Code Go hoặc Gopls, hãy thử chúng.
Hãy xem bài nói chuyện của Suzy Mueller
“[Building Better Projects with the Go Editor](https://www.youtube.com/watch?v=jMyzsp2E_0U)”
để có cái nhìn tổng quan.
Và như một phần thưởng thêm,
[Debugging Go in VS Code](/s/vscode-go-debug)
đã trở nên đáng tin cậy và mạnh mẽ hơn nhờ hỗ trợ gốc cho
[Debug Adapter Protocol](https://microsoft.github.io/debug-adapter-protocol/) của Delve.
Hãy thử “[Debugging Treasure Hunt](https://www.youtube.com/watch?v=ZPIPPRjwg7Q)” của Suzy!

Một khía cạnh khác của việc mở rộng phát triển là số lượng dependency trong một dự án.
Khoảng một tháng sau sinh nhật lần thứ 12 của Go,
[lỗ hổng Log4Shell](https://en.wikipedia.org/wiki/Log4Shell) đã trở thành
hồi chuông cảnh tỉnh cho cả ngành
về tầm quan trọng của bảo mật chuỗi cung ứng.
Hệ thống module của Go được thiết kế chính xác cho mục đích này,
giúp bạn hiểu và theo dõi các dependency của mình,
xác định chính xác những dependency nào bạn đang dùng,
và xác định liệu trong số đó có dependency nào có lỗ hổng đã biết hay không.
Bài blog của Filippo Valsorda
“[How Go Mitigates Supply Chain Attacks](/blog/supply-chain)”
đưa ra cái nhìn tổng quan về cách tiếp cận của chúng tôi.
Vào tháng 9, chúng tôi đã giới thiệu trước
cách tiếp cận của Go đối với quản lý lỗ hổng
trong bài blog của Julie Qiu “[Vulnerability Management for Go](/blog/vuln)”.
Cốt lõi của công việc đó là một cơ sở dữ liệu lỗ hổng mới, được tuyển chọn cẩn thận,
và lệnh mới [govulncheck](https://pkg.go.dev/golang.org/x/vuln/cmd/govulncheck),
sử dụng phân tích tĩnh nâng cao để loại bỏ phần lớn các cảnh báo dương tính giả
vốn sẽ xuất hiện nếu chỉ dựa vào riêng yêu cầu module.

Một phần trong nỗ lực tìm hiểu người dùng Go của chúng tôi là khảo sát Go thường niên vào cuối năm.
Năm nay, các nhà nghiên cứu trải nghiệm người dùng của chúng tôi cũng bổ sung thêm một khảo sát Go nhẹ hơn vào giữa năm.
Mục tiêu của chúng tôi là thu thập đủ phản hồi để có ý nghĩa thống kê
mà không trở thành gánh nặng cho cộng đồng Go nói chung.
Để xem kết quả, hãy đọc bài blog của Alice Merrick
“[Go Developer Survey 2021 Results](/blog/survey2021-results)”
và bài của Todd Kulesza
“[Go Developer Survey 2022 Q2 Results](/blog/survey2022-q2-results)”.

Khi thế giới bắt đầu đi lại nhiều hơn,
chúng tôi cũng rất vui khi được gặp nhiều bạn trực tiếp tại các hội nghị Go trong năm 2022,
đặc biệt là tại GopherCon Europe ở Berlin vào tháng 7 và GopherCon ở Chicago vào tháng 10.
Tuần trước, chúng tôi tổ chức sự kiện trực tuyến thường niên của mình,
[Go Day on Google Open Source Live](https://opensourcelive.withgoogle.com/events/go-day-2022).
Dưới đây là một số bài nói chuyện mà chúng tôi đã trình bày tại các sự kiện đó:

- “[How Go Became its Best Self](https://www.youtube.com/watch?v=vQm_whJZelc)”,
  của Cameron Balahan, tại GopherCon Europe.
- “[Go team Q&A](https://www.youtube.com/watch?v=KbOTTU9yEpI)”,
  với Cameron Balahan, Michael Knyszek và Than McIntosh, tại GopherCon Europe.
- “[Compatibility: How Go Programs Keep Working](https://www.youtube.com/watch?v=v24wrd3RwGo)”,
  của Russ Cox tại GopherCon.
- “[A Holistic Go Experience](https://www.gophercon.com/agenda/session/998660)”,
  của Cameron Balahan tại GopherCon (video vẫn chưa được đăng).
- “[Structured Logging for Go](https://opensourcelive.withgoogle.com/events/go-day-2022/watch?talk=talk2)”,
  của Jonathan Amsterdam tại Go Day on Google Open Source Live.
- “[Writing your Applications Faster and More Securely with Go](https://opensourcelive.withgoogle.com/events/go-day-2022/watch?talk=talk3)”,
  của Cody Oss tại Go Day on Google Open Source Live.
- “[Respecting Memory Limits in Go](https://opensourcelive.withgoogle.com/events/go-day-2022/watch?talk=talk4)”,
  của Michael Knyszek tại Go Day on Google Open Source Live.

Một cột mốc khác của năm nay là việc công bố
“[The Go Programming Language and Environment](https://cacm.acm.org/magazines/2022/5/260357-the-go-programming-language-and-environment/fulltext)”,
của Russ Cox, Robert Griesemer, Rob Pike, Ian Lance Taylor và Ken Thompson,
trên _Communications of the ACM_.
Bài viết, do chính những người thiết kế và hiện thực ban đầu của Go chấp bút,
giải thích điều gì theo chúng tôi đã làm cho Go trở nên phổ biến và hiệu quả đến vậy.
Nói ngắn gọn, điểm cốt lõi là Go tập trung vào việc cung cấp một môi trường phát triển đầy đủ
nhắm đến toàn bộ quy trình phát triển phần mềm,
với trọng tâm là khả năng mở rộng cho cả các nỗ lực kỹ nghệ phần mềm quy mô lớn
lẫn các triển khai quy mô lớn.

Trong năm thứ 14 của Go, chúng tôi sẽ tiếp tục làm việc để biến Go thành môi trường tốt nhất
cho kỹ nghệ phần mềm ở quy mô lớn.
Chúng tôi dự định tập trung đặc biệt vào bảo mật chuỗi cung ứng, cải thiện tính tương thích,
và structured logging, tất cả đều đã được liên kết trong bài viết này.
Và sẽ còn rất nhiều cải tiến khác nữa,
bao gồm tối ưu hóa dựa trên hồ sơ thực thi.

## Xin cảm ơn!

Go từ lâu đã luôn lớn hơn rất nhiều so với riêng những gì nhóm Go tại Google làm.
Xin cảm ơn tất cả các bạn, những người đóng góp và mọi người trong cộng đồng Go,
vì đã giúp Go trở thành môi trường lập trình thành công như ngày hôm nay.
Chúng tôi chúc tất cả các bạn những điều tốt đẹp nhất trong năm tới.
