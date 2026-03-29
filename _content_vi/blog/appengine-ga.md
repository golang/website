---
title: Go cho App Engine nay đã chính thức sẵn sàng
date: 2011-07-21
by:
- Andrew Gerrand
tags:
- appengine
- release
summary: Giờ đây bạn có thể dùng Go trên App Engine!
---


Các nhóm Go và App Engine rất hào hứng thông báo rằng môi trường chạy Go
cho App Engine nay đã chính thức sẵn sàng.
Điều này có nghĩa là bạn có thể lấy ứng dụng Go mà bạn đang làm (hoặc định
làm) và triển khai nó lên App Engine ngay bây giờ với [SDK 1.5.2](http://code.google.com/appengine/downloads.html) mới.

Kể từ khi chúng tôi công bố môi trường chạy Go tại Google I/O, chúng tôi đã tiếp tục [cải thiện và mở rộng](http://code.google.com/p/googleappengine/wiki/SdkForGoReleaseNotes)
hỗ trợ Go cho các API App Engine và đã bổ sung Channels API.
Go Datastore API giờ cũng hỗ trợ transaction và ancestor query.
Xem [tài liệu Go App Engine](https://code.google.com/appengine/docs/go/)
để biết toàn bộ chi tiết.

Đối với những ai đã dùng Go SDK,
xin lưu ý rằng bản phát hành 1.5.2 giới thiệu `api_version` 2.
Lý do là SDK mới dựa trên Go `release.r58.1` (phiên bản ổn định hiện tại
của Go) và không tương thích ngược với bản phát hành trước đó.
Các ứng dụng hiện có có thể cần thay đổi theo [ghi chú phát hành r58](/doc/devel/release.html#r58).
Sau khi cập nhật mã, bạn nên triển khai lại ứng dụng với dòng
`api_version: 2` trong tệp `app.yaml`.
Các ứng dụng được viết dựa trên `api_version` 1 sẽ ngừng hoạt động sau ngày 18 tháng 8.

Cuối cùng, chúng tôi xin gửi lời cảm ơn sâu sắc tới những người kiểm thử đáng tin cậy và vô số báo cáo lỗi của họ.
Sự giúp đỡ của họ là vô giá trong việc đạt tới cột mốc quan trọng này.

_Cách nhanh nhất để bắt đầu với Go trên App Engine là_ [_hướng dẫn Bắt đầu_](http://code.google.com/appengine/docs/go/gettingstarted/).

_Lưu ý rằng môi trường chạy Go vẫn được xem là mang tính thử nghiệm; nó chưa được hỗ trợ tốt như các môi trường chạy Python và Java._
