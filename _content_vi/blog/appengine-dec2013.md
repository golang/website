---
title: "Go trên App Engine: công cụ, kiểm thử và đồng thời"
date: 2013-12-13
by:
- Andrew Gerrand
- Johan Euphrosine
tags:
- appengine
summary: Công bố các cải tiến cho Go trên App Engine.
---

## Bối cảnh

Khi chúng tôi [ra mắt Go cho App Engine](/blog/go-and-google-app-engine)
vào tháng 5 năm 2011, SDK khi đó chỉ là một phiên bản đã được chỉnh sửa của SDK Python.
Vào thời điểm đó, chưa có cách chính thống nào để xây dựng hoặc tổ chức chương trình Go, vì thế
cách tiếp cận kiểu Python là hợp lý. Kể từ đó, Go 1.0 đã được phát hành,
bao gồm [công cụ go](/cmd/go/) và một
[quy ước](/doc/code.html) để tổ chức chương trình Go.

Vào tháng 1 năm 2013, chúng tôi đã công bố
[sự tích hợp tốt hơn](/blog/the-app-engine-sdk-and-workspaces-gopath)
giữa Go App Engine SDK và công cụ go, khuyến khích việc dùng
các import path thông thường trong ứng dụng App Engine và cho phép dùng "go
get" để lấy các dependency của ứng dụng.

Với bản phát hành App Engine 1.8.8 gần đây, chúng tôi vui mừng công bố thêm
nhiều cải tiến cho trải nghiệm phát triển Go trên App Engine.

## Công cụ goapp

Go App Engine SDK hiện bao gồm công cụ "goapp", một phiên bản dành riêng cho App Engine
của công cụ "go". Tên mới này cho phép người dùng giữ cả công cụ
"go" thông thường lẫn công cụ "goapp" trong PATH hệ thống của họ.

Ngoài các [lệnh](/cmd/go/) hiện có của công cụ "go",
công cụ "goapp" còn cung cấp các lệnh mới để làm việc với ứng dụng App Engine.
Lệnh "[goapp serve](https://developers.google.com/appengine/docs/go/tools/devserver)"
khởi động máy chủ phát triển cục bộ và lệnh
"[goapp deploy](https://developers.google.com/appengine/docs/go/tools/uploadinganapp)"
tải ứng dụng lên App Engine.

Các ưu điểm chính của các lệnh "goapp serve" và "goapp deploy"
là giao diện người dùng đơn giản hơn và tính nhất quán với các lệnh hiện có như
"go get" và "go fmt".
Ví dụ, để chạy một phiên bản cục bộ của ứng dụng trong thư mục hiện tại, hãy chạy:

	$ goapp serve

Để tải nó lên App Engine:

	$ goapp deploy

Bạn cũng có thể chỉ định import path Go để phục vụ hoặc triển khai:

	$ goapp serve github.com/user/myapp

Bạn thậm chí có thể chỉ định một tệp YAML để phục vụ hoặc triển khai một
[module](https://developers.google.com/appengine/docs/go/modules/) cụ thể:

	$ goapp deploy mymodule.yaml

Các lệnh này có thể thay thế phần lớn các trường hợp dùng `dev_appserver.py` và `appcfg.py`,
mặc dù các công cụ Python vẫn còn sẵn cho những trường hợp ít gặp hơn.

## Kiểm thử đơn vị cục bộ

Go App Engine SDK giờ hỗ trợ kiểm thử đơn vị cục bộ, sử dụng [package testing](https://developers.google.com/appengine/docs/go/tools/localunittesting)
gốc của Go và lệnh "[go test](/cmd/go/#hdr-Test_packages)"
(được SDK cung cấp dưới dạng "goapp test").

Hơn nữa, giờ đây bạn có thể viết các bài kiểm thử dùng dịch vụ App Engine.
[Package aetest](https://developers.google.com/appengine/docs/go/tools/localunittesting#Go_Introducing_the_aetest_package)
cung cấp một giá trị appengine.Context ủy quyền các request tới một
phiên bản tạm thời của máy chủ phát triển.

Để biết thêm thông tin về cách dùng "goapp test" và package aetest, hãy xem
[tài liệu Local Unit Testing for Go](https://developers.google.com/appengine/docs/go/tools/localunittesting).
Lưu ý rằng package aetest vẫn đang ở giai đoạn đầu;
chúng tôi hy vọng sẽ bổ sung thêm nhiều tính năng theo thời gian.

## Hỗ trợ đồng thời tốt hơn

Giờ đây có thể cấu hình số lượng request đồng thời được phục vụ bởi
mỗi instance động của ứng dụng bằng cách đặt tùy chọn
[`max_concurrent_requests`](https://developers.google.com/appengine/docs/go/modules/#max_concurrent_requests)
(chỉ khả dụng với [module Automatic Scaling](https://developers.google.com/appengine/docs/go/modules/#automatic_scaling)).

Đây là ví dụ về tệp `app.yaml`:

	application: maxigopher
	version: 1
	runtime: go
	api_version: go1
	automatic_scaling:
	  max_concurrent_requests: 100

Điều này cấu hình mỗi instance của ứng dụng phục vụ tối đa 100 request
đồng thời (tăng từ mặc định là 10). Bạn có thể cấu hình các instance Go
phục vụ tối đa 500 request đồng thời.

Thiết lập này cho phép các instance của bạn xử lý nhiều request đồng thời hơn bằng cách
tận dụng khả năng xử lý đồng thời hiệu quả của Go, điều này sẽ mang lại
khả năng tận dụng instance tốt hơn và cuối cùng là ít giờ instance bị tính phí hơn.

## Kết luận

Với những thay đổi này, Go trên App Engine tiện lợi và hiệu quả hơn bao giờ hết,
và chúng tôi hy vọng bạn sẽ thích những cải tiến này. Vui lòng tham gia
[nhóm google-appengine-go](http://groups.google.com/group/google-appengine-go/)
để nêu câu hỏi hoặc thảo luận về các thay đổi này với đội ngũ kỹ sư và
phần còn lại của cộng đồng.
