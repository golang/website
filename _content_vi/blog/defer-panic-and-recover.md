---
title: Defer, Panic và Recover
date: 2010-08-04
by:
- Andrew Gerrand
tags:
- defer
- panic
- recover
- technical
- function
summary: Giới thiệu các cơ chế điều khiển luồng defer, panic và recover của Go.
template: true
---


Go có các cơ chế điều khiển luồng quen thuộc:
if, for, switch, goto.
Nó cũng có câu lệnh `go` để chạy mã trong một goroutine riêng.
Ở đây tôi muốn bàn về một vài cơ chế ít phổ biến hơn:
defer, panic và recover.

Một **câu lệnh defer** đẩy một lời gọi hàm vào một danh sách.
Danh sách các lời gọi đã lưu sẽ được thực thi sau khi hàm bao quanh trả về.
Defer thường được dùng để đơn giản hóa các hàm thực hiện nhiều thao tác dọn dẹp khác nhau.

Ví dụ, hãy xem một hàm mở hai tệp rồi sao chép nội dung từ tệp này sang tệp kia:

	func CopyFile(dstName, srcName string) (written int64, err error) {
	    src, err := os.Open(srcName)
	    if err != nil {
	        return
	    }

	    dst, err := os.Create(dstName)
	    if err != nil {
	        return
	    }

	    written, err = io.Copy(dst, src)
	    dst.Close()
	    src.Close()
	    return
	}

Cách này chạy được, nhưng có một lỗi. Nếu lời gọi `os.Create` thất bại,
hàm sẽ trả về mà không đóng tệp nguồn.
Điều này có thể dễ dàng khắc phục bằng cách đặt lời gọi `src.Close` trước câu lệnh return thứ hai,
nhưng nếu hàm phức tạp hơn thì vấn đề có thể sẽ không dễ
nhận ra và xử lý như vậy.
Bằng cách thêm các câu lệnh defer, ta có thể bảo đảm rằng các tệp luôn được đóng:

	func CopyFile(dstName, srcName string) (written int64, err error) {
	    src, err := os.Open(srcName)
	    if err != nil {
	        return
	    }
	    defer src.Close()

	    dst, err := os.Create(dstName)
	    if err != nil {
	        return
	    }
	    defer dst.Close()

	    return io.Copy(dst, src)
	}

Câu lệnh defer cho phép ta nghĩ tới việc đóng từng tệp ngay sau khi mở nó,
và bảo đảm rằng, bất kể số lượng câu lệnh return trong hàm là bao nhiêu,
các tệp _sẽ_ được đóng.

Hành vi của câu lệnh defer rất đơn giản và dễ đoán. Có ba quy tắc cơ bản:

1. _Đối số của hàm bị defer sẽ được đánh giá tại thời điểm câu lệnh defer được đánh giá._

Trong ví dụ này, biểu thức "i" được đánh giá khi lời gọi `Println` được defer.
Lời gọi bị defer sẽ in ra "0" sau khi hàm trả về.

	func a() {
	    i := 0
	    defer fmt.Println(i)
	    i++
	    return
	}

2. _Các lời gọi hàm bị defer được thực thi theo thứ tự Last In First Out sau khi hàm bao quanh trả về._

Hàm này sẽ in ra "3210":

{{raw `
	func b() {
	    for i := 0; i < 4; i++ {
	        defer fmt.Print(i)
	    }
	}
`}}

3. _Các hàm bị defer có thể đọc và gán vào các named return value của hàm trả về._

Trong ví dụ này, một hàm bị defer tăng giá trị trả về `i` _sau khi_
hàm bao quanh trả về.
Vì thế, hàm này trả về 2:

	func c() (i int) {
	    defer func() { i++ }()
	    return 1
	}

Điều này rất tiện để chỉnh sửa giá trị lỗi được trả về của một hàm; lát nữa ta sẽ thấy một ví dụ về điều đó.

**Panic** là một hàm dựng sẵn dừng luồng điều khiển thông thường và bắt đầu _panicking_.
Khi hàm F gọi `panic`, việc thực thi của F dừng lại,
mọi hàm bị defer trong F vẫn được thực thi bình thường,
rồi F trả về cho caller của nó.
Đối với caller, lúc đó F hành xử như một lời gọi tới `panic`.
Quá trình này tiếp tục đi ngược lên stack cho đến khi mọi hàm trong goroutine hiện tại đều đã trả về,
lúc đó chương trình bị crash.
Panic có thể được khởi tạo bằng cách gọi trực tiếp `panic`.
Nó cũng có thể do lỗi runtime gây ra,
chẳng hạn truy cập mảng vượt chỉ số.

**Recover** là một hàm dựng sẵn giúp giành lại quyền điều khiển từ một goroutine đang panicking.
Recover chỉ hữu ích bên trong các hàm bị defer.
Trong quá trình thực thi bình thường, một lời gọi `recover` sẽ trả về `nil` và không có tác dụng nào khác.
Nếu goroutine hiện tại đang panicking, một lời gọi `recover` sẽ bắt lấy
giá trị đã truyền vào `panic` và khôi phục việc thực thi bình thường.

Dưới đây là một chương trình ví dụ minh họa cơ chế của panic và defer:

	package main

	import "fmt"

	func main() {
	    f()
	    fmt.Println("Returned normally from f.")
	}

	func f() {
	    defer func() {
	        if r := recover(); r != nil {
	            fmt.Println("Recovered in f", r)
	        }
	    }()
	    fmt.Println("Calling g.")
	    g(0)
	    fmt.Println("Returned normally from g.")
	}

	func g(i int) {
	    if i > 3 {
	        fmt.Println("Panicking!")
	        panic(fmt.Sprintf("%v", i))
	    }
	    defer fmt.Println("Defer in g", i)
	    fmt.Println("Printing in g", i)
	    g(i + 1)
	}

Hàm `g` nhận số nguyên `i`, và panic nếu `i` lớn hơn 3,
nếu không thì nó tự gọi lại chính nó với đối số `i+1`.
Hàm `f` defer một hàm gọi `recover` và in ra giá trị đã được recover
(nếu nó khác `nil`).
Hãy thử hình dung đầu ra của chương trình này sẽ như thế nào trước khi đọc tiếp.

Chương trình sẽ in ra:

	Calling g.
	Printing in g 0
	Printing in g 1
	Printing in g 2
	Printing in g 3
	Panicking!
	Defer in g 3
	Defer in g 2
	Defer in g 1
	Defer in g 0
	Recovered in f 4
	Returned normally from f.

Nếu ta bỏ hàm bị defer khỏi `f`, panic sẽ không được recover và
đi tới đỉnh call stack của goroutine,
kết thúc chương trình.
Chương trình đã chỉnh sửa này sẽ in ra:

	Calling g.
	Printing in g 0
	Printing in g 1
	Printing in g 2
	Printing in g 3
	Panicking!
	Defer in g 3
	Defer in g 2
	Defer in g 1
	Defer in g 0
	panic: 4

	panic PC=0x2a9cd8
	[stack trace omitted]

Để xem một ví dụ thực tế về **panic** và **recover**,
xem [gói json](/pkg/encoding/json/) trong
thư viện chuẩn Go.
Nó mã hóa một interface bằng một tập các hàm đệ quy.
Nếu có lỗi khi duyệt qua giá trị,
`panic` sẽ được gọi để tháo ngược stack lên lời gọi hàm cấp cao nhất,
nơi nó sẽ recover khỏi panic và trả về một giá trị lỗi thích hợp (xem
các phương thức 'error' và 'marshal' của kiểu `encodeState` trong [encode.go](/src/pkg/encoding/json/encode.go)).

Thông lệ trong các thư viện Go là ngay cả khi một package dùng panic nội bộ,
API bên ngoài của nó vẫn trình bày các giá trị lỗi trả về một cách tường minh.

Những cách dùng khác của **defer** (ngoài ví dụ `file.Close` ở trên) bao gồm giải phóng mutex:

	mu.Lock()
	defer mu.Unlock()

in phần chân trang:

	printHeader()
	defer printFooter()

và nhiều hơn nữa.

Tóm lại, câu lệnh defer (có hoặc không đi kèm panic và recover) cung cấp
một cơ chế điều khiển luồng khác thường nhưng mạnh mẽ.
Nó có thể được dùng để mô hình hóa nhiều tính năng vốn được triển khai bằng các cấu trúc
chuyên dụng trong những ngôn ngữ lập trình khác. Hãy thử dùng nó.
