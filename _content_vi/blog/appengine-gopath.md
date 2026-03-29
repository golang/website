---
title: App Engine SDK và workspace (GOPATH)
date: 2013-01-09
by:
- Andrew Gerrand
tags:
- appengine
- tools
- gopath
summary: App Engine SDK 1.7.4 bổ sung hỗ trợ cho workspace kiểu GOPATH.
---

## Giới thiệu

Khi chúng tôi phát hành Go 1, chúng tôi đã giới thiệu [công cụ go](/cmd/go/) và,
cùng với nó, khái niệm workspace.
Workspace (được chỉ định bởi biến môi trường GOPATH) là một quy ước
để tổ chức mã nguồn, giúp đơn giản hóa việc lấy,
xây dựng và cài đặt các package Go.
Nếu bạn chưa quen với workspace, vui lòng đọc [bài viết này](/doc/code.html)
hoặc xem [video screencast này](http://www.youtube.com/watch?v=XCsL89YtqCs) trước khi đọc tiếp.

Cho tới gần đây, các công cụ trong App Engine SDK chưa nhận biết workspace.
Nếu không có workspace, lệnh
"[go get](/cmd/go/#hdr-Download_and_install_packages_and_dependencies)"
không thể hoạt động,
và vì vậy tác giả ứng dụng phải tự cài đặt và cập nhật dependency của ứng dụng thủ công. Điều đó rất phiền.

Tất cả đã thay đổi với phiên bản 1.7.4 của App Engine SDK.
Công cụ [dev\_appserver](https://developers.google.com/appengine/docs/go/tools/devserver)
và [appcfg](https://developers.google.com/appengine/docs/go/tools/uploadinganapp)
giờ đã nhận biết workspace.
Khi chạy cục bộ hoặc tải ứng dụng lên,
các công cụ này giờ sẽ tìm dependency trong những workspace được chỉ định bởi
biến môi trường GOPATH.
Điều đó có nghĩa là giờ đây bạn có thể dùng "go get" khi xây dựng ứng dụng App Engine,
và chuyển đổi giữa chương trình Go thông thường và ứng dụng App Engine mà không cần thay đổi
môi trường hoặc thói quen của mình.

Ví dụ, giả sử bạn muốn xây dựng một ứng dụng dùng OAuth 2.0 để xác thực
với một dịch vụ từ xa.
Một thư viện OAuth 2.0 phổ biến cho Go là package [oauth2](https://godoc.org/golang.org/x/oauth2),
gói mà bạn có thể cài vào workspace bằng lệnh:

	go get golang.org/x/oauth2

Khi viết ứng dụng App Engine của bạn, hãy import package oauth giống như trong một chương trình Go thông thường:

	import "golang.org/x/oauth2"

Giờ đây, dù chạy ứng dụng bằng dev\_appserver hay triển khai nó bằng appcfg,
các công cụ sẽ tìm thấy package oauth trong workspace của bạn. Nó cứ thế hoạt động.

## Ứng dụng lai chạy độc lập/App Engine

Go App Engine SDK được xây dựng dựa trên package [net/http](/pkg/net/http/)
chuẩn của Go để phục vụ các yêu cầu web và,
do đó, nhiều web server Go có thể chạy trên App Engine chỉ với một vài thay đổi.
Ví dụ, [godoc](/cmd/godoc/) được đưa vào
bản phân phối Go như một chương trình độc lập,
nhưng nó cũng có thể chạy như một ứng dụng App Engine (godoc phục vụ [golang.org](/) từ App Engine).

Nhưng sẽ thật tuyệt nếu bạn có thể viết một chương trình vừa là
web server độc lập vừa là ứng dụng App Engine? Bằng cách dùng [ràng buộc xây dựng](/pkg/go/build/#hdr-Build_Constraints), bạn có thể làm vậy.

Ràng buộc xây dựng là các dòng chú thích xác định xem một tệp có nên
được đưa vào một package hay không.
Chúng thường được dùng nhất trong mã xử lý nhiều hệ điều hành
hoặc kiến trúc bộ xử lý khác nhau.
Ví dụ, package [path/filepath](/pkg/path/filepath/)
bao gồm tệp [symlink.go](/src/pkg/path/filepath/symlink.go),
trong đó chỉ định một ràng buộc xây dựng để đảm bảo tệp đó không được build trên Windows
(hệ thống không có symbolic link):

	// +build !windows

App Engine SDK giới thiệu một từ khóa ràng buộc xây dựng mới: "appengine". Các tệp chỉ định

	// +build appengine

sẽ được App Engine SDK build và bị công cụ go bỏ qua. Ngược lại, các tệp chỉ định

	// +build !appengine

sẽ bị App Engine SDK bỏ qua, trong khi công cụ go vẫn build chúng bình thường.

Thư viện [goprotobuf](http://code.google.com/p/goprotobuf/) dùng
cơ chế này để cung cấp hai cách triển khai cho một phần quan trọng của cơ chế mã hóa/giải mã:
[pointer\_unsafe.go](http://code.google.com/p/goprotobuf/source/browse/proto/pointer_unsafe.go)
là phiên bản nhanh hơn nhưng không thể dùng trên App Engine vì nó dùng
[package unsafe](/pkg/unsafe/),
trong khi [pointer\_reflect.go](http://code.google.com/p/goprotobuf/source/browse/proto/pointer_reflect.go)
là phiên bản chậm hơn tránh dùng unsafe bằng cách dùng [package reflect](/pkg/reflect/) thay thế.

Hãy lấy một web server Go đơn giản và biến nó thành một ứng dụng lai. Đây là main.go:

	package main

	import (
	    "fmt"
	    "net/http"
	)

	func main() {
	    http.HandleFunc("/", handler)
	    http.ListenAndServe("localhost:8080", nil)
	}

	func handler(w http.ResponseWriter, r *http.Request) {
	    fmt.Fprint(w, "Hello!")
	}

Hãy build bằng công cụ go và bạn sẽ có một tệp thực thi web server độc lập.

Hạ tầng App Engine cung cấp hàm main riêng để chạy phần tương đương
với ListenAndServe.
Để chuyển main.go thành một ứng dụng App Engine, hãy bỏ lệnh gọi tới ListenAndServe
và đăng ký handler trong một hàm init (được chạy trước main). Đây là app.go:

	package main

	import (
	    "fmt"
	    "net/http"
	)

	func init() {
	    http.HandleFunc("/", handler)
	}

	func handler(w http.ResponseWriter, r *http.Request) {
	    fmt.Fprint(w, "Hello!")
	}

Để biến đây thành một ứng dụng lai, ta cần tách nó thành một phần dành riêng cho App Engine,
một phần dành riêng cho tệp nhị phân độc lập, và các phần dùng chung cho cả hai phiên bản.
Trong trường hợp này, không có phần nào dành riêng cho App Engine,
vì vậy ta chỉ tách thành hai tệp:

app.go chỉ định và đăng ký hàm handler.
Nó giống hệt đoạn mã bên trên,
và không cần ràng buộc build nào vì nó phải được đưa vào mọi phiên bản của chương trình.

main.go chạy web server. Nó bao gồm ràng buộc build "!appengine",
vì nó chỉ được đưa vào khi build tệp nhị phân độc lập.

	// +build !appengine

	package main

	import "net/http"

	func main() {
	    http.ListenAndServe("localhost:8080", nil)
	}

Để xem một ứng dụng lai phức tạp hơn, hãy xem [công cụ present](https://godoc.org/golang.org/x/tools/present).

## Kết luận

Chúng tôi hy vọng những thay đổi này sẽ giúp việc làm việc trên các ứng dụng có dependency bên ngoài
trở nên dễ dàng hơn, và giúp duy trì các codebase chứa cả chương trình độc lập lẫn ứng dụng App Engine.
