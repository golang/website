---
title: Ví dụ có thể kiểm thử trong Go
date: 2015-05-07
by:
- Andrew Gerrand
tags:
- godoc
- testing
summary: Cách thêm ví dụ, đồng thời cũng là kiểm thử, vào các package của bạn.
template: true
---

## Giới thiệu

Các [ví dụ](/pkg/testing/#hdr-Examples) trong Godoc là những đoạn mã Go
được hiển thị như tài liệu package và cũng được xác minh bằng cách
chạy chúng như các bài kiểm thử.
Chúng cũng có thể được chạy bởi người dùng truy cập trang web godoc của package
và nhấn nút "Run" tương ứng.

Việc có tài liệu có thể thực thi cho một package bảo đảm rằng thông tin
sẽ không bị lỗi thời khi API thay đổi.

Thư viện chuẩn có rất nhiều ví dụ như vậy
(chẳng hạn xem [`strings` package](/pkg/strings/#Contains)).

Bài viết này giải thích cách viết các hàm ví dụ của riêng bạn.

## Ví dụ là kiểm thử

Các ví dụ được biên dịch (và có thể được thực thi) như một phần trong
bộ kiểm thử của package.

Giống như các bài kiểm thử thông thường, ví dụ là những hàm nằm trong các tệp
`_test.go` của package.
Tuy nhiên, khác với các hàm kiểm thử thông thường, các hàm ví dụ không nhận tham số
và bắt đầu bằng từ `Example` thay vì `Test`.

Gói [`reverse`](https://pkg.go.dev/golang.org/x/example/hello/reverse/)
là một phần của [kho ví dụ Go](https://cs.opensource.google/go/x/example).
Dưới đây là một ví dụ minh họa hàm `String` của nó:

	package reverse_test

	import (
		"fmt"

		"golang.org/x/example/hello/reverse"
	)

	func ExampleString() {
		fmt.Println(reverse.String("hello"))
		// Output: olleh
	}

Đoạn mã này có thể nằm trong `example_test.go` trong thư mục `reverse`.

Máy chủ tài liệu package Go _pkg.go.dev_ hiển thị ví dụ này
bên cạnh [tài liệu của hàm `String`](https://pkg.go.dev/golang.org/x/example/hello/reverse/#String):

{{image "examples/pkgdoc.png" 517}}

Khi chạy bộ kiểm thử của package, ta có thể thấy hàm ví dụ được thực thi
mà không cần thêm cấu hình nào:

	$ go test -v
	=== RUN   TestString
	--- PASS: TestString (0.00s)
	=== RUN   ExampleString
	--- PASS: ExampleString (0.00s)
	PASS
	ok  	golang.org/x/example/hello/reverse	0.209s

## Chú thích Output

Điều đó có nghĩa là hàm `ExampleString` "pass" như thế nào?

Khi thực thi ví dụ,
framework kiểm thử sẽ thu thập dữ liệu được ghi ra standard output
và sau đó so sánh đầu ra đó với chú thích "Output:" của ví dụ.
Bài kiểm thử pass nếu đầu ra của nó khớp với chú thích đầu ra.

Để thấy một ví dụ thất bại, ta có thể đổi nội dung chú thích đầu ra thành một giá trị
rõ ràng là sai

	func ExampleString() {
		fmt.Println(reverse.String("hello"))
		// Output: golly
	}

rồi chạy kiểm thử lại:

	$ go test
	--- FAIL: ExampleString (0.00s)
	got:
	olleh
	want:
	golly
	FAIL

Nếu ta bỏ hẳn chú thích output

	func ExampleString() {
		fmt.Println(reverse.String("hello"))
	}

thì hàm ví dụ sẽ được biên dịch nhưng không được thực thi:

	$ go test -v
	=== RUN   TestString
	--- PASS: TestString (0.00s)
	PASS
	ok  	golang.org/x/example/hello/reverse	0.110s

Những ví dụ không có chú thích output hữu ích để minh họa mã không thể
chạy như unit test, chẳng hạn mã truy cập mạng,
đồng thời vẫn bảo đảm ít nhất ví dụ đó biên dịch được.

## Tên hàm ví dụ

Godoc dùng một quy ước đặt tên để gắn một hàm ví dụ với
một định danh cấp package.

	func ExampleFoo()     // tài liệu cho hàm hoặc kiểu Foo
	func ExampleBar_Qux() // tài liệu cho phương thức Qux của kiểu Bar
	func Example()        // tài liệu cho toàn bộ package

Theo quy ước này, godoc hiển thị ví dụ `ExampleString`
bên cạnh tài liệu cho hàm `String`.

Có thể cung cấp nhiều ví dụ cho cùng một định danh bằng cách dùng hậu tố
bắt đầu bằng dấu gạch dưới theo sau là một chữ cái thường.
Mỗi ví dụ sau đây đều tài liệu cho hàm `String`:

	func ExampleString()
	func ExampleString_second()
	func ExampleString_third()

## Ví dụ lớn hơn

Đôi khi ta cần nhiều hơn chỉ một hàm để viết một ví dụ tốt.

Ví dụ, để minh họa [`sort` package](/pkg/sort/)
ta nên cho thấy một hiện thực của `sort.Interface`.
Vì phương thức không thể được khai báo bên trong thân hàm, ví dụ phải
bao gồm thêm ngữ cảnh bên cạnh chính hàm ví dụ.

Để làm điều này, ta có thể dùng một "ví dụ toàn tệp".
Ví dụ toàn tệp là một tệp kết thúc bằng `_test.go` và chứa đúng một
hàm ví dụ, không có hàm test hay benchmark nào, và có ít nhất một khai báo
cấp package khác.
Khi hiển thị kiểu ví dụ này, godoc sẽ cho thấy toàn bộ tệp.

Đây là một ví dụ toàn tệp từ package `sort`:

{{raw `
	package sort_test

	import (
		"fmt"
		"sort"
	)

	type Person struct {
		Name string
		Age  int
	}

	func (p Person) String() string {
		return fmt.Sprintf("%s: %d", p.Name, p.Age)
	}

	// ByAge implements sort.Interface for []Person based on
	// the Age field.
	type ByAge []Person

	func (a ByAge) Len() int           { return len(a) }
	func (a ByAge) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
	func (a ByAge) Less(i, j int) bool { return a[i].Age < a[j].Age }

	func Example() {
		people := []Person{
			{"Bob", 31},
			{"John", 42},
			{"Michael", 17},
			{"Jenny", 26},
		}

		fmt.Println(people)
		sort.Sort(ByAge(people))
		fmt.Println(people)

		// Output:
		// [Bob: 31 John: 42 Michael: 17 Jenny: 26]
		// [Michael: 17 Jenny: 26 Bob: 31 John: 42]
	}
`}}

Một package có thể chứa nhiều ví dụ toàn tệp; mỗi tệp một ví dụ.
Hãy xem [`sort` package's source code](/src/sort/)
để thấy điều này trong thực tế.

## Kết luận

Ví dụ Godoc là một cách tuyệt vời để viết và duy trì mã như tài liệu.
Chúng cũng cung cấp những ví dụ có thể chỉnh sửa, hoạt động thực sự và có thể chạy
để người dùng của bạn xây dựng tiếp trên đó.
Hãy dùng chúng!
