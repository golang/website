---
title: Sáu năm của Go
date: 2015-11-10
by:
- Andrew Gerrand
summary: Chúc mừng sinh nhật lần thứ sáu của Go!
template: true
---


Hôm nay là tròn sáu năm kể từ khi ngôn ngữ Go được phát hành như một dự án mã nguồn mở.
Kể từ đó, hơn 780 người đóng góp đã thực hiện hơn 30.000 commit cho
22 kho lưu trữ của dự án. Hệ sinh thái vẫn đang tiếp tục phát triển,
với GitHub ghi nhận hơn 90.000 kho lưu trữ Go.
Và ở ngoài đời, chúng tôi vẫn đều đặn thấy các sự kiện Go và nhóm người dùng mới xuất hiện
[khắp](/blog/gophercon2015)
[thế](/blog/gouk15)
[giới](/blog/gopherchina).

{{image "6years/6years-gopher.png"}}

Vào tháng 8, chúng tôi [phát hành Go 1.5](/blog/go1.5), bản phát hành
quan trọng nhất kể từ Go 1. Nó có một
[bộ gom rác được thiết kế lại hoàn toàn](/doc/go1.5#gc)
khiến ngôn ngữ này phù hợp hơn với những ứng dụng nhạy cảm về độ trễ;
đánh dấu sự chuyển đổi từ toolchain trình biên dịch viết bằng C
sang một toolchain [được viết hoàn toàn bằng Go](/doc/go1.5#c);
và bao gồm các bản port tới [những kiến trúc mới](/doc/go1.5#ports), với khả năng hỗ trợ tốt hơn
cho bộ xử lý ARM (loại chip cung cấp sức mạnh cho phần lớn điện thoại thông minh).
Những cải tiến này khiến Go phù hợp hơn với phạm vi tác vụ rộng hơn,
và đó là xu hướng mà chúng tôi hy vọng sẽ tiếp tục trong những năm tới.

Những cải tiến về công cụ tiếp tục nâng cao năng suất cho nhà phát triển.
Chúng tôi đã giới thiệu [execution tracer](/cmd/trace/) và
lệnh "[go doc](/cmd/go/#hdr-Show_documentation_for_package_or_symbol)",
cùng với nhiều nâng cấp hơn nữa cho các
[công cụ phân tích tĩnh](/talks/2014/static-analysis.slide) khác nhau của mình.
Chúng tôi cũng đang làm việc trên một
[plugin Go chính thức cho Sublime Text](https://groups.google.com/forum/#!topic/Golang-nuts/8oCSjAiKXUQ),
và khả năng hỗ trợ tốt hơn cho các trình soạn thảo khác cũng đang được chuẩn bị.

Đầu năm tới, chúng tôi sẽ phát hành thêm nhiều cải tiến trong Go 1.6, bao gồm
hỗ trợ HTTP/2 cho server và client [net/http](/pkg/net/http/),
cơ chế vendoring package chính thức, hỗ trợ cho block trong template văn bản và HTML,
một memory sanitizer kiểm tra cả mã Go lẫn C/C++,
và bộ sưu tập quen thuộc gồm nhiều cải tiến và bản sửa lỗi khác.

Đây là lần thứ sáu chúng tôi có niềm vui được viết bài blog chúc mừng sinh nhật cho Go,
và điều đó sẽ không thể xảy ra nếu thiếu những con người tuyệt vời và đầy đam mê trong cộng đồng của chúng ta.
Nhóm Go muốn cảm ơn tất cả mọi người đã
đóng góp mã, viết thư viện mã nguồn mở, viết bài blog, giúp một gopher mới,
hoặc đơn giản là thử dùng Go. Không có các bạn, Go sẽ không thể hoàn thiện,
hữu ích hay thành công như hôm nay. Xin cảm ơn, và hãy cùng ăn mừng!
