---
title: C? Go? Cgo!
date: 2011-03-17
by:
- Andrew Gerrand
tags:
- cgo
- technical
summary: Cách dùng cgo để cho package Go gọi mã C.
---

## Giới thiệu

Cgo cho phép package Go gọi mã C. Với một tệp mã nguồn Go được viết bằng một số tính năng đặc biệt,
cgo sẽ sinh ra các tệp Go và C có thể được kết hợp thành một package Go duy nhất.

Hãy bắt đầu bằng một ví dụ: đây là một package Go cung cấp hai hàm -
`Random` và `Seed` - bao bọc các hàm `random` và `srandom` của C.

	package rand

	/*
	#include <stdlib.h>
	*/
	import "C"

	func Random() int {
	    return int(C.random())
	}

	func Seed(i int) {
	    C.srandom(C.uint(i))
	}

Hãy xem điều gì đang diễn ra ở đây, bắt đầu từ câu lệnh import.

Package `rand` import `"C"`, nhưng bạn sẽ thấy không hề có package nào như vậy
trong thư viện chuẩn Go.
Đó là vì `C` là một "pseudo-package",
một tên đặc biệt được cgo diễn giải như một tham chiếu tới không gian tên của C.

Package `rand` có bốn tham chiếu tới package `C`:
các lệnh gọi `C.random` và `C.srandom`, phép chuyển đổi `C.uint(i)`,
và câu lệnh `import`.

Hàm `Random` gọi hàm `random` của thư viện chuẩn C và trả về kết quả.
Trong C, `random` trả về một giá trị kiểu C `long`,
kiểu mà cgo biểu diễn thành `C.long`.
Nó phải được chuyển thành kiểu Go trước khi có thể được dùng bởi mã Go ngoài package này,
bằng một phép chuyển kiểu Go thông thường:

	func Random() int {
	    return int(C.random())
	}

Đây là một hàm tương đương dùng biến tạm để minh họa phép chuyển kiểu rõ ràng hơn:

	func Random() int {
	    var r C.long = C.random()
	    return int(r)
	}

Hàm `Seed` làm điều ngược lại, theo một nghĩa nào đó.
Nó nhận một `int` Go thông thường, chuyển nó sang kiểu C `unsigned int`,
và truyền nó cho hàm C `srandom`.

	func Seed(i int) {
	    C.srandom(C.uint(i))
	}

Lưu ý rằng cgo biết kiểu `unsigned int` với tên `C.uint`;
xem [tài liệu cgo](/cmd/cgo) để có danh sách đầy đủ
các tên kiểu số này.

Chi tiết duy nhất trong ví dụ này mà chúng ta chưa xem là phần chú thích phía trên câu lệnh `import`.

	/*
	#include <stdlib.h>
	*/
	import "C"

Cgo nhận biết chú thích này. Mọi dòng bắt đầu bằng `#cgo` theo sau bởi
một ký tự khoảng trắng sẽ bị loại bỏ;
chúng trở thành directive cho cgo.
Những dòng còn lại được dùng làm phần header khi biên dịch các phần C của package.
Trong trường hợp này, các dòng đó chỉ là một câu lệnh `#include`,
nhưng chúng có thể là gần như bất kỳ đoạn mã C nào.
Các directive `#cgo` được dùng để cung cấp cờ cho compiler và linker
khi build các phần C của package.

Có một giới hạn: nếu chương trình của bạn dùng bất kỳ directive `//export` nào,
thì mã C trong phần chú thích chỉ có thể bao gồm khai báo (`extern int f();`),
không thể bao gồm định nghĩa (`int f() { return 1; }`).
Bạn có thể dùng directive `//export` để làm cho hàm Go có thể được mã C truy cập.

Các directive `#cgo` và `//export` được ghi lại trong [tài liệu cgo](/cmd/cgo/).

## Chuỗi và các thứ khác

Không giống Go, C không có kiểu chuỗi tường minh. Chuỗi trong C được biểu diễn bằng một mảng ký tự kết thúc bằng số không.

Việc chuyển đổi giữa chuỗi Go và chuỗi C được thực hiện bằng các hàm `C.CString`,
`C.GoString` và `C.GoStringN`.
Những phép chuyển đổi này tạo ra một bản sao của dữ liệu chuỗi.

Ví dụ tiếp theo này triển khai một hàm `Print` ghi một chuỗi ra
đầu ra chuẩn bằng hàm `fputs` của C từ thư viện `stdio`:

	package print

	// #include <stdio.h>
	// #include <stdlib.h>
	import "C"
	import "unsafe"

	func Print(s string) {
	    cs := C.CString(s)
	    C.fputs(cs, (*C.FILE)(C.stdout))
	    C.free(unsafe.Pointer(cs))
	}

Những cấp phát bộ nhớ do mã C thực hiện không được bộ quản lý bộ nhớ của Go biết đến.
Khi bạn tạo một chuỗi C bằng `C.CString` (hoặc bất kỳ cấp phát bộ nhớ C nào)
bạn phải nhớ giải phóng bộ nhớ khi dùng xong bằng cách gọi `C.free`.

Lệnh gọi tới `C.CString` trả về một con trỏ tới đầu mảng ký tự,
vì vậy trước khi hàm kết thúc, chúng ta chuyển nó thành một [`unsafe.Pointer`](/pkg/unsafe/#Pointer)
và giải phóng vùng cấp phát bộ nhớ bằng `C.free`.
Một thành ngữ phổ biến trong chương trình cgo là [`defer`](/doc/articles/defer_panic_recover.html)
lời gọi free ngay sau khi cấp phát (đặc biệt khi phần mã theo sau
phức tạp hơn một lệnh gọi hàm đơn lẻ),
như trong phiên bản viết lại của `Print` này:

	func Print(s string) {
	    cs := C.CString(s)
	    defer C.free(unsafe.Pointer(cs))
	    C.fputs(cs, (*C.FILE)(C.stdout))
	}

## Xây dựng package cgo

Để build package cgo, chỉ cần dùng [`go build`](/cmd/go/#hdr-Compile_packages_and_dependencies)
hoặc [`go install`](/cmd/go/#hdr-Compile_and_install_packages_and_dependencies) như bình thường.
Công cụ go nhận ra import đặc biệt `"C"` và tự động dùng cgo cho các tệp đó.

## Thêm tài nguyên về cgo

Tài liệu về [lệnh cgo](/cmd/cgo/) có nhiều
chi tiết hơn về pseudo-package C và quá trình build.
[Ví dụ cgo](/misc/cgo/) trong cây mã Go minh họa
các khái niệm nâng cao hơn.

Cuối cùng, nếu bạn tò mò về cách tất cả điều này hoạt động bên trong,
hãy xem chú thích giới thiệu của [cgocall.go](/src/runtime/cgocall.go) trong package runtime.
