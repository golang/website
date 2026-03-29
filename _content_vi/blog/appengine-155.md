---
title: Đã phát hành Go App Engine SDK 1.5.5
date: 2011-10-11
by:
- Andrew Gerrand
tags:
- appengine
- gofix
- release
summary: Go App Engine SDK 1.5.5 bao gồm Go release.r60.2.
---


Hôm nay chúng tôi đã phát hành Go App Engine SDK phiên bản 1.5.5.
Bạn có thể tải nó từ [trang tải xuống của App Engine](http://code.google.com/appengine/downloads.html).

Bản phát hành này bao gồm các thay đổi và cải tiến cho API App Engine, đồng thời
nâng chuỗi công cụ Go đi kèm lên [release.r60.2](/doc/devel/release.html#r60)
(bản phát hành ổn định hiện tại).
Cũng được bao gồm trong bản phát hành này là các công cụ
[godoc](/cmd/godoc/),
[gofmt](/cmd/gofmt/),
và [gofix](/cmd/gofix/) của chuỗi công cụ Go.
Bạn có thể tìm thấy chúng trong thư mục gốc của SDK.

Một số thay đổi trong bản phát hành này không tương thích ngược,
vì vậy chúng tôi đã tăng `api_version` của SDK lên 3.
Các ứng dụng hiện có sẽ cần thay đổi mã khi chuyển sang `api_version` 3.

Công cụ gofix đi kèm với SDK đã được tùy biến với các mô-đun dành riêng cho App Engine.
Nó có thể được dùng để tự động cập nhật các ứng dụng Go nhằm hoạt động với các package
appengine mới nhất và thư viện chuẩn Go đã được cập nhật.
Để cập nhật ứng dụng của bạn, hãy chạy:

	/path/to/sdk/gofix /path/to/your/app

SDK hiện cũng bao gồm mã nguồn package appengine,
vì vậy bạn có thể dùng godoc cục bộ để đọc tài liệu API App Engine:

	/path/to/sdk/godoc appengine/datastore Get

**Lưu ý quan trọng:** Chúng tôi đã ngừng hỗ trợ `api_version` 2.
Các ứng dụng Go dùng `api_version` 2 sẽ ngừng hoạt động sau ngày 16 tháng 12 năm 2011.
Vui lòng cập nhật ứng dụng của bạn để dùng `api_version` 3 trước thời điểm đó.

Xem [ghi chú phát hành](http://code.google.com/p/googleappengine/wiki/SdkForGoReleaseNotes)
để biết danh sách đầy đủ các thay đổi.
Vui lòng gửi mọi câu hỏi về SDK mới tới [nhóm thảo luận Go App Engine](http://groups.google.com/group/google-appengine-go).
