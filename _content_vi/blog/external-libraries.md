---
title: Điểm sáng về các thư viện Go bên ngoài
date: 2011-06-03
by:
- Andrew Gerrand
tags:
- community
- libraries
summary: Một số thư viện Go phổ biến và cách sử dụng chúng.
---


Trong khi các tác giả của Go đã làm việc tích cực để cải thiện thư viện chuẩn của Go,
cộng đồng rộng lớn hơn đã tạo ra một hệ sinh thái ngày càng phong phú của các thư viện bên ngoài.
Trong bài viết này, chúng ta sẽ xem một số thư viện Go phổ biến và cách chúng có thể được sử dụng.

[Mgo](http://labix.org/mgo) (đọc là "mango") là một trình điều khiển cơ sở dữ liệu MongoDB.
[MongoDB](http://www.mongodb.org/) là một [cơ sở dữ liệu hướng tài liệu](http://en.wikipedia.org/wiki/Document-oriented_database)
với một danh sách dài các tính năng phù hợp cho [nhiều trường hợp sử dụng](http://www.mongodb.org/display/DOCS/Use%2BCases).
Gói `mgo` cung cấp một API Go giàu tính biểu cảm và đúng phong cách Go để làm việc với MongoDB,
từ các thao tác cơ bản như chèn và cập nhật bản ghi cho đến các tính năng nâng cao hơn như
[MapReduce](http://www.mongodb.org/display/DOCS/MapReduce) và
[GridFS](http://www.mongodb.org/display/DOCS/GridFS).
Mgo có nhiều tính năng thú vị bao gồm tự động khám phá cụm và
nạp trước kết quả, hãy xem [trang chủ mgo](http://labix.org/mgo) để biết
chi tiết và mã ví dụ.
Khi làm việc với những tập dữ liệu lớn, Go, MongoDB
và mgo là một sự kết hợp mạnh mẽ.

[Authcookie](https://github.com/dchest/authcookie) là một thư viện web dùng để
tạo và xác minh cookie xác thực người dùng.
Nó cho phép máy chủ web phát ra các token an toàn về mặt mật mã gắn với
một người dùng cụ thể và sẽ hết hạn sau một khoảng thời gian xác định.
Nó có một API đơn giản giúp việc thêm cơ chế xác thực
vào các ứng dụng web hiện có trở nên dễ dàng.
Xem [tệp README](https://github.com/dchest/authcookie/blob/master/README.md)
để biết chi tiết và mã ví dụ.

[Go-charset](http://code.google.com/p/go-charset) cung cấp hỗ trợ cho
việc chuyển đổi giữa mã hóa UTF-8 chuẩn của Go và nhiều bộ ký tự khác nhau.
Gói `go-charset` hiện thực một `io.Reader` và `io.Writer` có khả năng chuyển đổi,
vì vậy bạn có thể bao bọc các Reader và Writer hiện có (như kết nối mạng
hoặc file descriptor),
giúp việc giao tiếp với các hệ thống dùng mã hóa ký tự khác trở nên dễ dàng.

[Go-socket.io](https://github.com/madari/go-socket.io) là một hiện thực của Go
cho [Socket.IO](http://socket.io/),
một API client/server cho phép máy chủ web đẩy thông điệp tới trình duyệt web.
Tùy thuộc vào khả năng của trình duyệt người dùng,
Socket.IO sử dụng kiểu truyền tải tốt nhất cho kết nối,
có thể là websocket hiện đại, AJAX long polling,
hoặc [một cơ chế khác](http://socket.io/#transports).
Go-socket.io bắc cầu giữa các máy chủ Go và các client JavaScript giàu tính năng
trên nhiều loại trình duyệt.
Để cảm nhận rõ hơn về go-socket.io hãy xem [ví dụ máy chủ chat](https://github.com/madari/go-socket.io/blob/master/example/example.go).

Điều đáng nhắc tới là các gói này đều [có thể cài bằng goinstall](/cmd/goinstall/).
Với một bản [cài đặt](/doc/install.html) Go được cập nhật,
bạn có thể cài tất cả chúng bằng một lệnh duy nhất:

	goinstall launchpad.net/mgo \
	    github.com/dchest/authcookie \
	    go-charset.googlecode.com/hg/charset \
	    github.com/madari/go-socket.io

Sau khi được cài bằng goinstall, các gói có thể được import bằng chính các đường dẫn đó:

	import (
	    "launchpad.net/mgo"
	    "github.com/dchest/authcookie"
	    "go-charset.googlecode.com/hg/charset"
	    "github.com/madari/go-socket.io"
	)

Ngoài ra, vì giờ chúng là một phần của hệ thống Go cục bộ,
chúng ta có thể xem tài liệu của chúng bằng [godoc](/cmd/godoc/):

	godoc launchpad.net/mgo Database # xem tài liệu cho kiểu Database

Dĩ nhiên, đây mới chỉ là phần nổi của tảng băng;
còn rất nhiều thư viện Go tuyệt vời khác được liệt kê trên [bảng điều khiển gói](http://godashboard.appspot.com/package)
và sẽ còn nhiều thư viện nữa trong tương lai.
