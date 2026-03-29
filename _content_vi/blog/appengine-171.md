---
title: Các cập nhật Go trong App Engine 1.7.1
date: 2012-08-22
by:
- Andrew Gerrand
tags:
- appengine
- release
summary: App Engine SDK 1.7.1 bổ sung memcache và các chức năng khác cho Go.
---


Tuần này chúng tôi đã phát hành App Engine SDK phiên bản 1.7.1.
Nó bao gồm một số cập nhật đáng kể dành riêng cho môi trường chạy App Engine của Go.

[Package memcache](https://developers.google.com/appengine/docs/go/memcache/reference)
đã có thêm một số bổ sung cho kiểu tiện ích [Codec](https://developers.google.com/appengine/docs/go/memcache/reference#Codec) của nó.
Các phương thức SetMulti, AddMulti, CompareAndSwap và CompareAndSwapMulti
giúp việc lưu trữ và cập nhật dữ liệu đã mã hóa trong [Dịch vụ Memcache](https://developers.google.com/appengine/docs/go/memcache/overview) trở nên dễ dàng hơn.

[Công cụ bulkloader](https://developers.google.com/appengine/docs/go/tools/uploadingdata)
giờ đây có thể được dùng với các ứng dụng Go,
cho phép người dùng tải lên và tải xuống hàng loạt các bản ghi datastore.
Điều này hữu ích cho sao lưu và xử lý ngoại tuyến,
đồng thời là một trợ giúp lớn khi di chuyển các ứng dụng Python hoặc Java sang môi trường chạy Go.

[Dịch vụ Images](https://developers.google.com/appengine/docs/go/images/overview)
giờ đã sẵn sàng cho người dùng Go.
[Package appengine/image](https://developers.google.com/appengine/docs/go/images/reference) mới hỗ trợ
phục vụ hình ảnh trực tiếp từ Blobstore và thay đổi kích thước hoặc cắt ảnh ngay tức thì.
Lưu ý rằng đây chưa phải là toàn bộ dịch vụ ảnh như trong SDK Python và Java,
vì phần lớn chức năng tương đương đã có sẵn trong [package image chuẩn của Go](/pkg/image/)
và các package bên ngoài như [graphics-go](http://code.google.com/p/graphics-go/).

Hàm [runtime.RunInBackground](https://developers.google.com/appengine/docs/go/backends/runtime#RunInBackground) mới
cho phép các yêu cầu backend tạo ra một yêu cầu mới độc lập với yêu cầu ban đầu.
Chúng có thể chạy ở chế độ nền miễn là backend còn sống.

Cuối cùng, chúng tôi đã bổ sung một số chức năng còn thiếu:
[package xmpp](https://developers.google.com/appengine/docs/go/xmpp/reference) giờ đây
hỗ trợ gửi cập nhật hiện diện và lời mời chat, cũng như truy xuất
trạng thái hiện diện của người dùng khác,
và [package user](https://developers.google.com/appengine/docs/go/users/reference) hỗ trợ
xác thực client bằng OAuth.

Bạn có thể lấy SDK mới từ [trang tải xuống App Engine](https://developers.google.com/appengine/downloads#Google_App_Engine_SDK_for_Go)
và xem [tài liệu đã cập nhật](https://developers.google.com/appengine/docs/go).
