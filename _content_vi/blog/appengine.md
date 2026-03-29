---
title: Go và Google App Engine
date: 2011-05-10
by:
- David Symonds
- Nigel Tao
- Andrew Gerrand
tags:
- appengine
- release
summary: Công bố hỗ trợ Go trong Google App Engine.
---


App Engine của Google cung cấp một cách đáng tin cậy,
có khả năng mở rộng và dễ dàng để xây dựng và triển khai ứng dụng cho web.
Hơn một trăm nghìn ứng dụng được lưu trữ tại appspot.com và các tên miền tùy chỉnh
sử dụng hạ tầng App Engine.
Ban đầu được viết cho ứng dụng Python, đến năm 2009 hệ thống đã thêm môi trường chạy Java.
Và hôm nay, tại Google I/O, chúng tôi rất vui mừng thông báo rằng Go sẽ là cái tên tiếp theo.
Hiện tại nó được đánh dấu là một tính năng App Engine mang tính thử nghiệm,
vì đây vẫn còn là những ngày đầu, nhưng cả đội App Engine lẫn đội Go đều rất
hào hứng với cột mốc này.

Với những ngày đầu, ý chúng tôi là việc triển khai vẫn đang được mở dần.
Tính tới hôm nay, App Engine SDK cho Go đã [sẵn sàng để tải về](http://code.google.com/p/googleappengine/downloads/list),
và chúng tôi sẽ sớm cho phép triển khai ứng dụng Go vào hạ tầng lưu trữ App Engine.
Hôm nay, thông qua SDK, bạn sẽ có thể viết ứng dụng web,
tìm hiểu các API (và ngôn ngữ này, nếu nó còn mới với bạn),
và chạy ứng dụng web cục bộ.
Khi khả năng triển khai đầy đủ được bật, việc đẩy ứng dụng của bạn lên đám mây của Google sẽ rất dễ dàng.

Một trong những điều thú vị nhưng ít hiển nhiên hơn của thông báo này là nó mang lại
một cách rất dễ dàng để thử Go.
Bạn thậm chí không cần cài Go từ trước vì SDK
hoàn toàn tự chứa.
Chỉ cần tải SDK xuống, giải nén và bắt đầu viết mã.
Hơn nữa, “dev app server” của SDK có nghĩa là bạn thậm chí không cần
tự chạy compiler;
mọi thứ đều tự động một cách tuyệt vời.

Điều bạn sẽ tìm thấy trong SDK là nhiều API App Engine chuẩn,
được thiết kế riêng theo phong cách Go tốt, bao gồm Datastore,
Blobstore, URL Fetch, Mail, Users, v.v.
Nhiều API khác sẽ được bổ sung khi môi trường phát triển.
Môi trường chạy cung cấp toàn bộ ngôn ngữ Go và gần như mọi thư viện chuẩn,
ngoại trừ một vài thứ không hợp lý trong môi trường App Engine.
Ví dụ, không có package `unsafe` và package `syscall` đã được tinh giản.
(Phần triển khai dùng một phiên bản mở rộng của thiết lập trong [Go Playground](/doc/play/)
trên [golang.org](/).)

Ngoài ra, dù goroutine và channel đều hiện diện,
khi một ứng dụng Go chạy trên App Engine thì chỉ có một luồng được chạy trong một instance nhất định.
Tức là, mọi goroutine đều chạy trong một luồng hệ điều hành duy nhất,
vì vậy sẽ không có song song CPU cho một yêu cầu client nhất định.
Chúng tôi kỳ vọng giới hạn này sẽ được gỡ bỏ vào một thời điểm nào đó.

Bất chấp các hạn chế nhỏ này, đây vẫn là ngôn ngữ thật:
Mã được triển khai dưới dạng mã nguồn và được biên dịch trên đám mây bằng compiler x86 64-bit (6g),
khiến nó trở thành ngôn ngữ biên dịch thực thụ đầu tiên chạy trên App Engine.
Go trên App Engine cho phép triển khai các ứng dụng web hiệu quả,
tốn nhiều CPU.

Nếu bạn muốn biết thêm, hãy đọc [tài liệu](http://code.google.com/appengine/docs/go/)
(hãy bắt đầu với “[Getting Started](http://code.google.com/appengine/docs/go/gettingstarted/)”).
Các thư viện và SDK là mã nguồn mở, được lưu trữ tại [http://code.google.com/p/appengine-go/](http://code.google.com/p/appengine-go/).
Chúng tôi đã tạo một danh sách thư điện tử [google-appengine-go](http://groups.google.com/group/google-appengine-go) mới;
bạn cứ thoải mái liên hệ tại đó với các câu hỏi riêng về App Engine.
[Trình theo dõi lỗi của App Engine](http://code.google.com/p/googleappengine/issues/list)
là nơi để báo cáo các vấn đề liên quan tới Go SDK mới.

Go App Engine SDK [đã có sẵn](http://code.google.com/p/googleappengine/downloads/list)
cho Linux và Mac OS X (10.5 trở lên);
chúng tôi hy vọng phiên bản Windows cũng sẽ sớm xuất hiện.

Chúng tôi muốn gửi lời cảm ơn vì mọi sự giúp đỡ và nhiệt tình mà chúng tôi nhận được
từ đội App Engine của Google trong quá trình biến điều này thành hiện thực.
