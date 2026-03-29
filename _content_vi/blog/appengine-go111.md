---
title: Công bố môi trường chạy Go 1.11 mới của App Engine
date: 2018-10-16
by:
- Eno Compton
- Tyler Bui-Palsulich
tags:
- appengine
summary: Google Cloud công bố môi trường chạy Go 1.11 mới cho App Engine, với ít giới hạn hơn về cấu trúc ứng dụng.
template: true
---


[App Engine](https://cloud.google.com/appengine/) đã ra mắt
[hỗ trợ thử nghiệm cho Go](/blog/go-and-google-app-engine)
vào năm 2011. Trong những năm tiếp theo, cộng đồng Go đã phát triển đáng kể và
đã định hình các mẫu
thành ngữ cho ứng dụng trên đám mây. Hôm nay, Google Cloud
[công bố môi trường chạy Go 1.11 mới](https://cloud.google.com/blog/products/application-development/go-1-11-is-now-available-on-app-engine)
cho môi trường chuẩn của App Engine, cung cấp toàn bộ
sức mạnh của App Engine, những thứ như chỉ trả tiền cho những gì bạn dùng, tự động mở rộng,
và hạ tầng được quản lý, trong khi vẫn hỗ trợ Go theo đúng thành ngữ.

Bắt đầu từ Go 1.11, Go trên App Engine không còn giới hạn về cấu trúc ứng dụng,
các package được hỗ trợ, giá trị `context.Context`, hay HTTP client. Hãy viết ứng dụng Go
của bạn theo cách bạn muốn, thêm một tệp `app.yaml`, và ứng dụng đã sẵn sàng
để triển khai trên App Engine.
[Specifying Dependencies](https://cloud.google.com/appengine/docs/standard/go111/specifying-dependencies)
mô tả cách môi trường chạy mới
hỗ trợ [vendoring](/cmd/go/#hdr-Vendor_Directories) và
[modules](/doc/go1.11#modules) (thử nghiệm) cho việc
quản lý dependency.

Cùng với [Cloud Functions hỗ trợ Go](https://twitter.com/kelseyhightower/status/1035278586754813952)
(sẽ nói thêm trong một bài viết sau), App Engine cung cấp một cách hấp dẫn để chạy
mã Go trên Google Cloud Platform (GCP) mà không cần bận tâm tới hạ tầng
bên dưới.

Hãy cùng xem cách tạo một ứng dụng nhỏ cho App Engine. Trong ví dụ
này, chúng tôi giả định một quy trình làm việc dựa trên `GOPATH`, dù Go modules
cũng đã có [hỗ trợ thử nghiệm](https://cloud.google.com/appengine/docs/standard/go111/specifying-dependencies).

Đầu tiên, bạn tạo ứng dụng trong `GOPATH` của mình:

{{code "appengine/main.go"}}

Đoạn mã chứa một cách thiết lập theo đúng thành ngữ cho một HTTP server nhỏ phản hồi
“Hello, 世界.” Nếu bạn có kinh nghiệm với App Engine trước đây, bạn sẽ nhận ra sự
vắng mặt của bất kỳ lệnh gọi nào tới `appengine.Main()`, thứ nay hoàn toàn là tùy chọn.
Hơn nữa, mã ứng dụng hoàn toàn có tính di động, không có sự ràng buộc nào với
hạ tầng mà ứng dụng của bạn được triển khai lên.

Nếu bạn cần dùng dependency bên ngoài, bạn có thể thêm các dependency đó vào
một thư mục `vendor` hoặc vào một tệp `go.mod`, cả hai đều được môi trường chạy mới
hỗ trợ.

Sau khi hoàn tất mã ứng dụng, hãy tạo một tệp `app.yaml` để chỉ định
môi trường chạy:

	runtime: go111

Cuối cùng, hãy chuẩn bị máy của bạn với một tài khoản Google Cloud Platform:

  - Tạo một tài khoản với [GCP](https://cloud.google.com).
  - [Tạo một project](https://cloud.google.com/resource-manager/docs/creating-managing-projects).
  - Cài đặt [Cloud SDK](https://cloud.google.com/sdk/) trên hệ thống của bạn.

Khi mọi thiết lập đã hoàn tất, bạn có thể triển khai chỉ với một lệnh:

	gcloud app deploy

Chúng tôi tin rằng các nhà phát triển Go sẽ thấy môi trường chạy Go 1.11 mới cho App Engine
là một bổ sung thú vị cho các lựa chọn hiện có để chạy ứng dụng Go. Có một
[gói miễn phí](https://cloud.google.com/free/). Hãy xem
[hướng dẫn bắt đầu](https://cloud.google.com/appengine/docs/standard/go111/building-app/)
hoặc
[hướng dẫn di chuyển](https://cloud.google.com/appengine/docs/standard/go111/go-differences)
và triển khai một ứng dụng lên môi trường chạy mới ngay hôm nay!
