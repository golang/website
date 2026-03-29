---
title: Ngôn ngữ lập trình Go tròn hai tuổi
date: 2011-11-10
by:
- Andrew Gerrand
tags:
- appengine
- community
- gopher
summary: Chúc mừng sinh nhật lần thứ hai của Go!
template: true
---


Hai năm trước, một nhóm nhỏ tại Google đã công khai dự án non trẻ của mình:
Ngôn ngữ lập trình Go.
Họ giới thiệu một đặc tả ngôn ngữ, hai trình biên dịch,
một thư viện chuẩn khiêm tốn, một vài công cụ mới lạ,
và rất nhiều tài liệu chính xác (dù ngắn gọn).
Họ dõi theo với sự hào hứng khi các lập trình viên trên khắp thế giới bắt đầu thử nghiệm với Go.
Nhóm tiếp tục lặp lại, cải tiến những gì họ đã xây dựng,
và dần dần có thêm hàng chục, rồi hàng trăm lập trình viên
từ cộng đồng mã nguồn mở tham gia cùng.
Các tác giả của Go tiếp tục tạo ra rất nhiều thư viện,
công cụ mới, và hàng loạt [tài liệu](/doc/docs.html).
Họ kỷ niệm một năm thành công trước công chúng bằng một [bài blog](/blog/go-one-year-ago-today)
vào tháng 11 năm ngoái với câu kết rằng “Go chắc chắn đã sẵn sàng cho production,
nhưng vẫn còn chỗ để cải thiện.
Trọng tâm trước mắt của chúng tôi là làm cho chương trình Go nhanh hơn và
hiệu quả hơn trong bối cảnh các hệ thống hiệu năng cao.”

Hôm nay là dịp kỷ niệm hai năm kể từ ngày Go ra mắt,
và Go giờ đây nhanh hơn và ổn định hơn bao giờ hết.
Việc tinh chỉnh cẩn trọng các bộ sinh mã, primitive đồng thời,
bộ gom rác, và các thư viện cốt lõi của Go đã làm tăng hiệu năng của chương trình Go,
và hỗ trợ gốc cho [profiling](/blog/profiling-go-programs)
và [debugging](/blog/debugging-go-programs-with-gnu-debugger)
khiến việc phát hiện và loại bỏ các vấn đề hiệu năng trong mã người dùng trở nên dễ hơn.
Go cũng giờ đây dễ học hơn với [A Tour of Go](/tour/),
một hướng dẫn tương tác mà bạn có thể học ngay trong trình duyệt web của mình.

Trong năm nay, chúng tôi đã giới thiệu [Go runtime](http://code.google.com/appengine/docs/go/) ở trạng thái thử nghiệm
cho nền tảng App Engine của Google,
và đã đều đặn mở rộng mức độ hỗ trợ của Go runtime đối với các API của App Engine.
Ngay trong tuần này, chúng tôi đã phát hành [phiên bản 1.6.0](http://code.google.com/appengine/downloads.html)
của Go App Engine SDK,
bao gồm hỗ trợ cho [backends](http://code.google.com/appengine/docs/go/backends/overview.html)
(các tiến trình chạy lâu),
khả năng kiểm soát chi tiết hơn với datastore indexes, và nhiều cải tiến khác.
Ngày nay, Go runtime đã gần đạt mức ngang bằng tính năng và là một lựa chọn khả thi
so với Python runtime và Java runtime.
Trên thực tế, hiện nay chúng tôi phục vụ [golang.org](/) bằng cách chạy một phiên bản
của [godoc](/cmd/godoc/) trên dịch vụ App Engine.

Nếu như năm 2010 là năm của khám phá và thử nghiệm,
thì năm 2011 là năm của tinh chỉnh và lập kế hoạch cho tương lai.
Trong năm nay, chúng tôi đã phát hành nhiều phiên bản Go dạng "[release](/doc/devel/release.html)"
ổn định và được hỗ trợ tốt hơn so với các bản chụp hằng tuần.
Chúng tôi cũng giới thiệu [gofix](/cmd/gofix/) để giảm bớt
đau đớn khi chuyển sang bản phát hành mới hơn.
Hơn nữa, tháng trước chúng tôi đã công bố [kế hoạch cho Go version 1](/blog/preview-of-go-version-1),
một bản phát hành sẽ được hỗ trợ trong nhiều năm tới.
Công việc hướng tới Go 1 hiện đã bắt đầu và bạn có thể theo dõi tiến độ
qua bản chụp hằng tuần mới nhất tại [weekly.golang.org](http://weekly.golang.org/pkg/).

Kế hoạch là phát hành Go 1 vào đầu năm 2012.
Chúng tôi cũng hy vọng có thể đưa Go App Engine runtime ra khỏi trạng thái “thử nghiệm” vào cùng thời điểm.

Nhưng đó chưa phải là tất cả. Năm 2011 cũng là một năm thú vị với chú gopher.
Nó đã xuất hiện dưới dạng thú nhồi bông (một món quà rất được săn đón tại Google I/O
và các buổi nói chuyện về Go khác) và dưới dạng tượng vinyl
(được tặng cho mọi người tham dự OSCON và hiện có bán tại [Google Store](http://www.googlestore.com/Fun/Go+Gopher+Figurine.axd)).

{{image "2years/2years-gophers.jpg"}}

Và bất ngờ nhất, vào dịp Halloween, chú ấy còn xuất hiện cùng cô bạn gopher của mình!

{{image "2years/2years-costume.jpg"}}

Ảnh do [Chris Nokleberg](https://plus.google.com/106640494112897458359/posts) chụp.
