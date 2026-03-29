---
title: Thử nghiệm, Đơn giản hóa, Phát hành
date: 2019-08-01
by:
- Russ Cox
tags:
- community
- go2
- proposals
summary: Cách chúng tôi phát triển Go, từ một bài nói chuyện tại GopherCon 2019.
template: true
---

## Giới thiệu

Đây là phiên bản bài blog của bài nói chuyện tôi đã trình bày tuần trước tại GopherCon 2019.

{{video "https://www.youtube.com/embed/kNHo788oO5Y?rel=0"}}

Tất cả chúng ta đều đang cùng nhau đi trên con đường hướng tới Go 2,
nhưng không ai trong chúng ta biết chính xác con đường đó dẫn tới đâu,
hay đôi khi thậm chí nó đang đi theo hướng nào.
Bài viết này thảo luận về cách chúng ta thực sự
tìm ra và đi theo con đường tới Go 2.
Đây là hình dạng của quá trình đó.

<div style="margin-left: 2em;">
{{image "experiment/expsimp1.png" 179}}
</div>

Chúng ta thử nghiệm với Go như nó đang tồn tại hiện nay
để hiểu nó rõ hơn,
học xem điều gì vận hành tốt và điều gì không.
Sau đó chúng ta thử nghiệm với những thay đổi khả dĩ
để hiểu chúng rõ hơn,
và tiếp tục học xem điều gì hiệu quả, điều gì không.
Dựa trên những gì học được từ các thử nghiệm đó,
chúng ta đơn giản hóa.
Rồi lại thử nghiệm tiếp.
Rồi lại đơn giản hóa tiếp.
Và cứ thế.
Và cứ thế.

## Bốn chữ R của việc đơn giản hóa

Trong quá trình này, có bốn cách chính để chúng ta có thể đơn giản hóa
trải nghiệm tổng thể khi viết chương trình Go:
tái định hình, tái định nghĩa, loại bỏ và giới hạn.

**Đơn giản hóa bằng cách tái định hình**

Cách đầu tiên là tái định hình những gì đang có sang một dạng mới,
một dạng mà xét tổng thể sẽ đơn giản hơn.

Mỗi chương trình Go chúng ta viết đều là một thí nghiệm để kiểm nghiệm chính Go.
Trong những ngày đầu của Go, chúng tôi nhanh chóng nhận ra rằng
người ta thường viết mã như hàm `addToList` này:

	func addToList(list []int, x int) []int {
		n := len(list)
		if n+1 > cap(list) {
			big := make([]int, n, (n+5)*2)
			copy(big, list)
			list = big
		}
		list = list[:n+1]
		list[n] = x
		return list
	}

Chúng tôi cũng viết cùng kiểu mã đó cho slice byte,
slice string, v.v.
Chương trình của chúng tôi quá phức tạp,
vì Go khi ấy lại quá đơn giản.

Vì thế chúng tôi lấy rất nhiều hàm kiểu `addToList` trong các chương trình
và tái định hình chúng thành một hàm do chính Go cung cấp.
Việc thêm `append` làm ngôn ngữ Go phức tạp hơn một chút,
nhưng xét toàn cục,
nó làm trải nghiệm viết chương trình Go đơn giản hơn,
ngay cả khi tính cả chi phí học về `append`.

Đây là một ví dụ khác.
Đối với Go 1, chúng tôi xem xét rất nhiều công cụ phát triển
trong bộ phân phối Go và tái định hình chúng thành một lệnh mới.

	5a      8g
	5g      8l
	5l      cgo
	6a      gobuild
	6cov    gofix         →     go
	6g      goinstall
	6l      gomake
	6nm     gopack
	8a      govet

Lệnh `go` nay đã trở nên quá trung tâm đến mức
người ta dễ quên rằng chúng ta từng sống rất lâu mà không có nó
và việc đó kéo theo bao nhiêu công sức thừa.

Chúng tôi thêm mã và độ phức tạp vào bộ phân phối Go,
nhưng xét tổng thể thì đã đơn giản hóa trải nghiệm viết chương trình Go.
Cấu trúc mới cũng tạo ra khoảng trống cho những thử nghiệm thú vị khác,
mà ta sẽ thấy ở phần sau.

**Đơn giản hóa bằng cách tái định nghĩa**

Cách thứ hai là tái định nghĩa
chức năng mà ta đã có,
để nó làm được nhiều hơn.
Giống như đơn giản hóa bằng tái định hình,
đơn giản hóa bằng tái định nghĩa giúp chương trình dễ viết hơn,
nhưng lần này không cần học thêm gì mới.

Ví dụ, ban đầu `append` chỉ được định nghĩa để đọc từ slice.
Khi nối thêm vào một `[]byte`, bạn có thể nối các byte từ một `[]byte` khác,
nhưng không thể nối các byte từ một string.
Chúng tôi tái định nghĩa `append` để cho phép nối từ string,
mà không thêm bất kỳ điều gì mới vào ngôn ngữ.

	var b []byte
	var more []byte
	b = append(b, more...) // ok

	var b []byte
	var more string
	b = append(b, more...) // ok later

**Đơn giản hóa bằng cách loại bỏ**

Cách thứ ba là đơn giản hóa bằng cách loại bỏ chức năng
khi hóa ra nó kém hữu ích
hoặc kém quan trọng hơn chúng tôi từng kỳ vọng.
Loại bỏ chức năng đồng nghĩa bớt đi một thứ phải học,
bớt đi một thứ phải sửa lỗi,
bớt đi một thứ gây xao nhãng hoặc dễ bị dùng sai.
Dĩ nhiên, việc loại bỏ cũng
buộc người dùng phải cập nhật chương trình hiện có của họ,
có thể làm chúng phức tạp hơn,
để bù cho phần bị bỏ đi.
Nhưng kết quả chung vẫn có thể là
quá trình viết chương trình Go trở nên đơn giản hơn.

Một ví dụ là khi chúng tôi loại bỏ
các dạng boolean của thao tác channel không chặn khỏi ngôn ngữ:

{{raw `
	ok := c <- x  // before Go 1, was non-blocking send
	x, ok := <-c  // before Go 1, was non-blocking receive
`}}

Các thao tác này cũng có thể thực hiện bằng `select`,
khiến việc phải quyết định dùng dạng nào trở nên rối rắm.
Loại bỏ chúng đã đơn giản hóa ngôn ngữ mà không làm giảm sức mạnh của nó.

**Đơn giản hóa bằng cách giới hạn**

Chúng ta cũng có thể đơn giản hóa bằng cách giới hạn những gì được phép.
Ngay từ ngày đầu, Go đã giới hạn cách mã hóa của tệp mã nguồn Go:
chúng phải là UTF-8.
Giới hạn này làm cho mọi chương trình cố gắng đọc tệp mã nguồn Go trở nên đơn giản hơn.
Những chương trình đó không phải bận tâm về
các tệp mã nguồn Go được mã hóa bằng Latin-1, UTF-16, UTF-7 hay thứ gì khác.

Một giới hạn quan trọng khác là `gofmt` cho định dạng chương trình.
Không gì từ chối mã Go không được định dạng bằng `gofmt`,
nhưng chúng tôi đã xác lập một quy ước rằng các công cụ viết lại chương trình Go
sẽ để chúng ở dạng `gofmt`.
Nếu bạn cũng giữ chương trình của mình ở dạng `gofmt`,
thì các bộ viết lại đó sẽ không tạo ra thay đổi định dạng nào.
Khi bạn so sánh trước và sau,
mọi khác biệt bạn nhìn thấy đều là thay đổi thực sự.
Giới hạn này đã đơn giản hóa các bộ viết lại chương trình
và dẫn tới những thử nghiệm thành công như
`goimports`, `gorename`, và nhiều công cụ khác.

## Quy trình phát triển Go

Chu kỳ thử nghiệm rồi đơn giản hóa này là một mô hình tốt
cho những gì chúng tôi đã làm trong mười năm qua.
Nhưng nó có một vấn đề:
nó quá đơn giản.
Chúng ta không thể chỉ thử nghiệm và đơn giản hóa.

Chúng ta phải phát hành kết quả.
Chúng ta phải đưa nó ra để mọi người sử dụng.
Dĩ nhiên, việc sử dụng lại tạo điều kiện cho nhiều thử nghiệm hơn,
và có thể là nhiều lần đơn giản hóa hơn,
và quá trình cứ thế lặp đi lặp lại.

<div style="margin-left: 2em;">
{{image "experiment/expsimp2.png" 326}}
</div>

Chúng tôi lần đầu phát hành Go cho tất cả các bạn
vào ngày 10 tháng 11 năm 2009.
Sau đó, với sự giúp đỡ của các bạn, chúng ta cùng phát hành Go 1 vào tháng 3 năm 2012.
Và từ đó đến nay chúng tôi đã phát hành thêm mười hai bản Go nữa.
Tất cả đều là những cột mốc quan trọng,
để cho phép nhiều thử nghiệm hơn,
để giúp chúng tôi hiểu thêm về Go,
và tất nhiên để đưa Go vào sử dụng trong sản xuất.

Khi phát hành Go 1,
chúng tôi đã chủ động chuyển trọng tâm sang việc dùng Go,
để hiểu phiên bản này của ngôn ngữ rõ hơn nhiều
trước khi thử thêm bất kỳ sự đơn giản hóa nào
liên quan tới thay đổi ngôn ngữ.
Chúng tôi cần dành thời gian để thử nghiệm,
để thực sự hiểu điều gì hiệu quả và điều gì không.

Dĩ nhiên, từ Go 1 đến nay chúng tôi đã có mười hai bản phát hành,
vì vậy chúng tôi vẫn tiếp tục thử nghiệm, đơn giản hóa và phát hành.
Nhưng chúng tôi tập trung vào những cách đơn giản hóa việc phát triển Go
mà không đòi hỏi thay đổi ngôn ngữ lớn và không làm hỏng
các chương trình Go hiện có.
Ví dụ, Go 1.5 phát hành bộ gom rác đồng thời đầu tiên,
và các bản phát hành sau đó tiếp tục cải thiện nó,
đơn giản hóa việc phát triển Go bằng cách loại bỏ thời gian dừng
khỏi danh sách mối lo thường trực.

Tại Gophercon 2017, chúng tôi thông báo rằng sau năm năm
thử nghiệm, đã đến lúc
suy nghĩ lại về
những thay đổi đáng kể có thể đơn giản hóa việc phát triển Go.
Con đường đến Go 2 thực ra cũng chính là con đường đến Go 1:
thử nghiệm, đơn giản hóa và phát hành,
hướng tới mục tiêu tổng thể là làm cho việc phát triển Go đơn giản hơn.

Đối với Go 2, những chủ đề cụ thể mà chúng tôi tin là
quan trọng nhất cần giải quyết là
xử lý lỗi, generics và dependencies.
Kể từ đó chúng tôi nhận ra rằng
một chủ đề quan trọng khác là công cụ cho lập trình viên.

Phần còn lại của bài viết này bàn về cách
công việc của chúng tôi trong từng lĩnh vực
đi theo con đường đó.
Trên đường đi,
chúng ta sẽ có một chặng rẽ,
dừng lại để xem xét chi tiết kỹ thuật
của những gì sắp được phát hành trong Go 1.13
cho xử lý lỗi.

## Lỗi

Đã khó để viết một chương trình
hoạt động đúng trong mọi trường hợp
khi mọi đầu vào đều hợp lệ và chính xác
và không thứ gì mà chương trình phụ thuộc vào bị hỏng.
Khi đưa lỗi vào bài toán,
việc viết một chương trình vẫn hoạt động đúng
bất kể điều gì xảy ra lại còn khó hơn.

Trong quá trình suy nghĩ về Go 2,
chúng tôi muốn hiểu rõ hơn
xem Go có thể giúp làm công việc đó đơn giản hơn hay không.

Có hai khía cạnh khác nhau
có thể được đơn giản hóa:
error values và error syntax.
Ta sẽ lần lượt xem từng phần,
với chặng rẽ kỹ thuật tôi đã hứa tập trung vào
những thay đổi của Go 1.13 đối với error values.

**Giá trị lỗi**

Error value phải bắt đầu từ đâu đó.
Đây là hàm `Read` từ phiên bản đầu tiên của gói `os`:

	export func Read(fd int64, b *[]byte) (ret int64, errno int64) {
		r, e := syscall.read(fd, &b[0], int64(len(b)));
		return r, e
	}

Khi đó chưa có kiểu `File`, và cũng chưa có kiểu error.
`Read` cùng các hàm khác trong gói
trả về trực tiếp `errno int64` từ lời gọi hệ thống Unix bên dưới.

Đoạn mã này được commit vào ngày 10 tháng 9 năm 2008 lúc 12:14 trưa.
Khi đó, giống như mọi thứ khác, nó là một thí nghiệm,
và mã thay đổi rất nhanh.
Hai giờ năm phút sau, API đã thay đổi:

	export type Error struct { s string }

	func (e *Error) Print() { … } // to standard error!
	func (e *Error) String() string { … }

	export func Read(fd int64, b *[]byte) (ret int64, err *Error) {
		r, e := syscall.read(fd, &b[0], int64(len(b)));
		return r, ErrnoToError(e)
	}

API mới này giới thiệu kiểu `Error` đầu tiên.
Một error giữ một string và có thể trả về string đó
đồng thời cũng có thể in nó ra stderr.

Mục đích ở đây là khái quát hóa vượt ra ngoài mã số nguyên.
Chúng tôi biết từ kinh nghiệm trước đó
rằng các mã lỗi của hệ điều hành là một biểu diễn quá hạn chế,
rằng việc không phải nhồi nhét toàn bộ chi tiết về một lỗi vào 64 bit
sẽ giúp chương trình đơn giản hơn.
Việc dùng chuỗi lỗi từng hoạt động khá ổn với chúng tôi trước đây,
nên chúng tôi cũng làm như vậy ở đây.
API mới này tồn tại trong bảy tháng.

Tháng 4 năm sau, sau khi có thêm kinh nghiệm với interface,
chúng tôi quyết định khái quát hóa hơn nữa
và cho phép người dùng tự định nghĩa hiện thực lỗi,
bằng cách biến chính kiểu `os.Error` thành một interface.
Chúng tôi đơn giản hóa bằng cách bỏ phương thức `Print`.

Đến Go 1, hai năm sau đó,
dựa trên một đề xuất của Roger Peppe,
`os.Error` trở thành kiểu dựng sẵn `error`,
và phương thức `String` được đổi tên thành `Error`.
Từ đó đến nay không có gì thay đổi.
Nhưng chúng tôi đã viết rất nhiều chương trình Go,
và kết quả là đã thử nghiệm rất nhiều về cách
hiện thực và sử dụng lỗi sao cho tốt nhất.

**Errors Are Values**

Việc biến `error` thành một interface đơn giản
và cho phép nhiều hiện thực khác nhau
có nghĩa là toàn bộ ngôn ngữ Go
đều có sẵn để định nghĩa và kiểm tra lỗi.
Chúng tôi thích nói rằng [errors are values](/blog/errors-are-values),
giống như mọi giá trị Go khác.

Đây là một ví dụ.
Trên Unix,
một lần thử quay số kết nối mạng
cuối cùng sẽ dùng lời gọi hệ thống `connect`.
Lời gọi đó trả về một `syscall.Errno`,
là một kiểu số nguyên có tên đại diện cho
mã lỗi system call
và hiện thực interface `error`:

	package syscall

	type Errno int64

	func (e Errno) Error() string { ... }

	const ECONNREFUSED = Errno(61)

	    ... err == ECONNREFUSED ...

Gói `syscall` cũng định nghĩa các hằng có tên
cho các mã lỗi do hệ điều hành máy chủ quy định.
Trong trường hợp này, trên hệ thống này, `ECONNREFUSED` là số 61.
Mã nhận được error từ một hàm
có thể kiểm tra xem error đó có phải là `ECONNREFUSED` hay không
bằng [so sánh giá trị](/ref/spec#Comparison_operators) thông thường.

Đi lên một tầng,
trong gói `os`,
bất kỳ lỗi system call nào cũng được báo cáo bằng
một cấu trúc lỗi lớn hơn, ghi lại thao tác nào đã được thử
ngoài chính lỗi.
Có một vài cấu trúc như vậy.
Cấu trúc này, `SyscallError`, mô tả lỗi
khi gọi một system call cụ thể
mà không ghi thêm thông tin nào khác:

	package os

	type SyscallError struct {
		Syscall string
		Err     error
	}

	func (e *SyscallError) Error() string {
		return e.Syscall + ": " + e.Err.Error()
	}

Đi lên thêm một tầng nữa,
trong gói `net`,
bất kỳ lỗi mạng nào cũng được báo cáo bằng một cấu trúc lỗi
còn lớn hơn nữa, ghi lại chi tiết
của thao tác mạng bao quanh,
chẳng hạn dial hay listen,
và network cùng địa chỉ liên quan:

	package net

	type OpError struct {
		Op     string
		Net    string
		Source Addr
		Addr   Addr
		Err    error
	}

	func (e *OpError) Error() string { ... }

Kết hợp tất cả những điều đó,
các lỗi được trả về từ các thao tác như `net.Dial` có thể được định dạng thành chuỗi,
nhưng chúng cũng là các giá trị dữ liệu Go có cấu trúc.
Trong trường hợp này, lỗi là một `net.OpError`, thêm ngữ cảnh
vào một `os.SyscallError`, vốn lại thêm ngữ cảnh vào một `syscall.Errno`:

	c, err := net.Dial("tcp", "localhost:50001")

	// "dial tcp [::1]:50001: connect: connection refused"

	err is &net.OpError{
		Op:   "dial",
		Net:  "tcp",
		Addr: &net.TCPAddr{IP: ParseIP("::1"), Port: 50001},
		Err: &os.SyscallError{
			Syscall: "connect",
			Err:     syscall.Errno(61), // == ECONNREFUSED
		},
	}

Khi chúng tôi nói errors are values, ý là cả
toàn bộ ngôn ngữ Go đều có sẵn để định nghĩa chúng
và cũng là
toàn bộ ngôn ngữ Go đều có sẵn để kiểm tra chúng.

Đây là một ví dụ từ gói net.
Hóa ra khi bạn thử thiết lập kết nối socket,
phần lớn thời gian bạn sẽ kết nối được hoặc nhận connection refused,
nhưng đôi khi bạn có thể gặp `EADDRNOTAVAIL` giả,
không có lý do rõ ràng.
Go bảo vệ chương trình người dùng khỏi kiểu lỗi này bằng cách thử lại.
Để làm vậy, nó phải kiểm tra cấu trúc lỗi để tìm xem
`syscall.Errno` nằm sâu bên trong có phải là `EADDRNOTAVAIL` hay không.

Đây là đoạn mã:

	func spuriousENOTAVAIL(err error) bool {
		if op, ok := err.(*OpError); ok {
			err = op.Err
		}
		if sys, ok := err.(*os.SyscallError); ok {
			err = sys.Err
		}
		return err == syscall.EADDRNOTAVAIL
	}

Một [type assertion](/ref/spec#Type_assertions) bóc đi lớp bọc `net.OpError` nếu có.
Và sau đó type assertion thứ hai bóc đi lớp bọc `os.SyscallError` nếu có.
Rồi hàm kiểm tra lỗi đã được gỡ bọc xem có bằng `EADDRNOTAVAIL` hay không.

Điều chúng tôi học được sau nhiều năm kinh nghiệm,
từ quá trình thử nghiệm với lỗi trong Go,
là việc có thể định nghĩa
các hiện thực tùy ý của interface `error`,
có toàn bộ ngôn ngữ Go sẵn có
để vừa dựng vừa gỡ cấu trúc lỗi,
và không bắt buộc dùng một hiện thực duy nhất,
là rất mạnh mẽ.

Những thuộc tính này, rằng lỗi là giá trị
và rằng không có một hiện thực lỗi bắt buộc duy nhất,
là những điều quan trọng cần giữ lại.

Việc không áp đặt một hiện thực lỗi duy nhất
đã giúp mọi người thử nghiệm
những chức năng bổ sung mà một error có thể cung cấp,
dẫn tới nhiều gói,
chẳng hạn
[github.com/pkg/errors](https://godoc.org/github.com/pkg/errors),
[gopkg.in/errgo.v2](https://godoc.org/gopkg.in/errgo.v2),
[github.com/hashicorp/errwrap](https://godoc.org/github.com/hashicorp/errwrap),
[upspin.io/errors](https://godoc.org/upspin.io/errors),
[github.com/spacemonkeygo/errors](https://godoc.org/github.com/spacemonkeygo/errors),
và nhiều nữa.

Tuy vậy, một vấn đề của việc thử nghiệm không ràng buộc là
ở phía người dùng,
bạn phải lập trình theo hợp của
mọi hiện thực có thể gặp.
Một sự đơn giản hóa mà có vẻ đáng để khám phá cho Go 2
là định nghĩa một phiên bản chuẩn của
chức năng thường được thêm vào,
dưới dạng các interface tùy chọn đã được thống nhất,
để các hiện thực khác nhau có thể tương tác với nhau.

**Unwrap**

Chức năng được thêm vào phổ biến nhất
trong các gói này là một phương thức nào đó
có thể được gọi để bỏ ngữ cảnh khỏi một lỗi,
trả về error nằm bên trong.
Các gói dùng những tên và ý nghĩa khác nhau
cho thao tác này, và đôi khi nó bỏ một lớp ngữ cảnh,
còn đôi khi bỏ được nhiều lớp nhất có thể.

Đối với Go 1.13, chúng tôi đã đưa ra một quy ước rằng một hiện thực error
thêm ngữ cảnh có thể gỡ bỏ vào một lỗi bên trong
nên hiện thực phương thức `Unwrap` trả về lỗi bên trong đó,
tức là gỡ lớp ngữ cảnh.
Nếu không có lỗi bên trong nào phù hợp để lộ cho phía gọi,
thì hoặc error đó không nên có `Unwrap`,
hoặc `Unwrap` nên trả về nil.

	// Go 1.13 optional method for error implementations.

	interface {
		// Unwrap removes one layer of context,
		// returning the inner error if any, or else nil.
		Unwrap() error
	}

Cách gọi phương thức tùy chọn này là dùng hàm trợ giúp `errors.Unwrap`,
vốn xử lý các trường hợp như bản thân error là nil
hoặc hoàn toàn không có phương thức `Unwrap`.

	package errors

	// Unwrap returns the result of calling
	// the Unwrap method on err,
	// if err’s type defines an Unwrap method.
	// Otherwise, Unwrap returns nil.
	func Unwrap(err error) error

Ta có thể dùng phương thức `Unwrap`
để viết một phiên bản đơn giản hơn và tổng quát hơn của `spuriousENOTAVAIL`.
Thay vì tìm các hiện thực wrapper cụ thể
như `net.OpError` hay `os.SyscallError`,
phiên bản tổng quát có thể lặp, gọi `Unwrap` để gỡ ngữ cảnh,
cho đến khi chạm tới `EADDRNOTAVAIL` hoặc không còn lỗi nào:

	func spuriousENOTAVAIL(err error) bool {
		for err != nil {
			if err == syscall.EADDRNOTAVAIL {
				return true
			}
			err = errors.Unwrap(err)
		}
		return false
	}

Vòng lặp này quá phổ biến,
nên Go 1.13 còn định nghĩa thêm hàm thứ hai là `errors.Is`,
hàm này liên tục unwrap một error để tìm một đích cụ thể.
Vì vậy ta có thể thay toàn bộ vòng lặp bằng một lời gọi `errors.Is`:

	func spuriousENOTAVAIL(err error) bool {
		return errors.Is(err, syscall.EADDRNOTAVAIL)
	}

Đến đây có lẽ ta thậm chí sẽ không định nghĩa riêng hàm nữa;
gọi trực tiếp `errors.Is` ngay tại call site cũng rõ ràng tương đương
mà lại đơn giản hơn.

Go 1.13 cũng giới thiệu hàm `errors.As`
unwrap cho tới khi tìm thấy một kiểu hiện thực cụ thể.

Nếu bạn muốn viết mã hoạt động với
các lỗi được bọc tùy ý,
thì `errors.Is` là phiên bản có nhận thức wrapper
của phép kiểm tra bằng nhau giữa các lỗi:

	err == target

	    →

	errors.Is(err, target)

Và `errors.As` là phiên bản có nhận thức wrapper
của type assertion trên lỗi:

	target, ok := err.(*Type)
	if ok {
	    ...
	}

	    →

	var target *Type
	if errors.As(err, &target) {
	   ...
	}

**Có Unwrap Hay Không?**

Việc có cho phép unwrap một error hay không là một quyết định API,
giống như việc có export một field của struct hay không cũng là một quyết định API.
Đôi khi việc lộ chi tiết đó cho mã gọi là phù hợp,
đôi khi thì không.
Khi phù hợp, hãy hiện thực `Unwrap`.
Khi không phù hợp, đừng hiện thực `Unwrap`.

Cho tới nay, `fmt.Errorf` chưa từng lộ ra
lỗi bên dưới được định dạng với `%v` cho phía gọi kiểm tra.
Tức là, kết quả của `fmt.Errorf` trước nay không thể bị unwrap.
Hãy xét ví dụ này:

	// errors.Unwrap(err2) == nil
	// err1 is not available (same as earlier Go versions)
	err2 := fmt.Errorf("connect: %v", err1)

Nếu `err2` được trả về cho
phía gọi, phía đó trước nay chưa từng có cách nào mở `err2` ra để truy cập `err1`.
Chúng tôi giữ nguyên thuộc tính đó trong Go 1.13.

Đối với những lúc bạn thật sự muốn cho phép unwrap kết quả của `fmt.Errorf`,
chúng tôi cũng thêm một format verb mới là `%w`, định dạng giống `%v`,
yêu cầu đối số phải là giá trị error,
và làm cho phương thức `Unwrap` của lỗi kết quả trả về chính đối số đó.
Trong ví dụ trên, giả sử ta thay `%v` bằng `%w`:

	// errors.Unwrap(err4) == err3
	// (%w is new in Go 1.13)
	err4 := fmt.Errorf("connect: %w", err3)

Giờ đây, nếu `err4` được trả về cho phía gọi,
phía gọi có thể dùng `Unwrap` để lấy ra `err3`.

Điều quan trọng cần lưu ý là các quy tắc tuyệt đối kiểu
“luôn dùng `%v` (hoặc đừng bao giờ hiện thực `Unwrap`)” hay “luôn dùng `%w` (hoặc luôn hiện thực `Unwrap`)”
đều sai như các quy tắc tuyệt đối kiểu “đừng bao giờ export field của struct” hay “luôn export field của struct”.
Thay vào đó, quyết định đúng phụ thuộc vào việc
liệu phía gọi có nên được phép kiểm tra và phụ thuộc vào
thông tin bổ sung mà việc dùng `%w` hay hiện thực `Unwrap` phơi bày hay không.

Để minh họa cho điểm này,
mọi kiểu bọc lỗi trong thư viện chuẩn
vốn đã có trường `Err` được export
giờ cũng có thêm `Unwrap` trả về trường đó,
nhưng những hiện thực có trường lỗi không export thì không,
và các chỗ dùng `fmt.Errorf` với `%v` từ trước
vẫn tiếp tục dùng `%v`, không đổi sang `%w`.

**In giá trị lỗi (đã bỏ)**

Cùng với bản phác thảo thiết kế cho Unwrap,
chúng tôi cũng công bố một
[bản phác thảo thiết kế cho một phương thức tùy chọn để in lỗi phong phú hơn](/design/go2draft-error-printing),
bao gồm thông tin stack frame
và hỗ trợ lỗi đã được bản địa hóa, dịch thuật.

	// Optional method for error implementations
	type Formatter interface {
		Format(p Printer) (next error)
	}

	// Interface passed to Format
	type Printer interface {
		Print(args ...interface{})
		Printf(format string, args ...interface{})
		Detail() bool
	}

Thiết kế này không đơn giản như `Unwrap`,
và tôi sẽ không đi vào chi tiết ở đây.
Khi chúng tôi thảo luận thiết kế này với cộng đồng Go trong mùa đông vừa qua,
chúng tôi học được rằng thiết kế đó chưa đủ đơn giản.
Nó quá khó để từng kiểu lỗi riêng lẻ hiện thực,
và nó cũng không giúp được đủ nhiều cho các chương trình hiện có.
Xét tổng thể, nó không làm cho việc phát triển Go đơn giản hơn.

Kết quả của cuộc thảo luận với cộng đồng này
là chúng tôi từ bỏ thiết kế in ấn đó.

**Cú pháp lỗi**

Đó là phần giá trị lỗi.
Giờ ta hãy nhìn nhanh sang cú pháp lỗi,
một thử nghiệm khác cũng đã bị bỏ.

Đây là một đoạn mã từ
[`compress/lzw/writer.go`](https://go.googlesource.com/go/+/go1.12/src/compress/lzw/writer.go#209) trong thư viện chuẩn:

{{raw `
	// Write the savedCode if valid.
	if e.savedCode != invalidCode {
		if err := e.write(e, e.savedCode); err != nil {
			return err
		}
		if err := e.incHi(); err != nil && err != errOutOfCodes {
			return err
		}
	}

	// Write the eof code.
	eof := uint32(1)<<e.litWidth + 1
	if err := e.write(e, eof); err != nil {
		return err
	}
`}}

Nhìn thoáng qua, khoảng một nửa đoạn mã này là kiểm tra lỗi.
Mắt tôi như mờ đi khi đọc nó.
Và chúng ta biết rằng mã vừa nhàm chán khi viết vừa nhàm chán khi đọc thì rất dễ bị đọc sót,
khiến nó trở thành nơi trú ẩn lý tưởng cho những lỗi khó tìm.
Ví dụ, một trong ba lần kiểm tra lỗi này không giống hai lần còn lại,
điều rất dễ bỏ lỡ khi lướt nhanh.
Nếu bạn đang debug đoạn mã này, bạn sẽ mất bao lâu để nhận ra điều đó?

Tại Gophercon năm ngoái, chúng tôi
[trình bày một bản thiết kế nháp](/design/go2draft-error-handling)
cho một cấu trúc điều khiển luồng mới được đánh dấu bằng từ khóa `check`.
`Check` tiêu thụ kết quả lỗi từ một lời gọi hàm hoặc biểu thức.
Nếu lỗi khác nil, `check` trả về lỗi đó.
Nếu không, `check` đánh giá thành các kết quả còn lại
từ lời gọi đó. Ta có thể dùng `check` để đơn giản hóa đoạn mã lzw:

{{raw `
	// Write the savedCode if valid.
	if e.savedCode != invalidCode {
		check e.write(e, e.savedCode)
		if err := e.incHi(); err != errOutOfCodes {
			check err
		}
	}

	// Write the eof code.
	eof := uint32(1)<<e.litWidth + 1
	check e.write(e, eof)
`}}

Phiên bản này của cùng đoạn mã dùng `check`,
giúp loại bỏ bốn dòng mã và
quan trọng hơn là làm nổi bật việc
lời gọi `e.incHi` được phép trả về `errOutOfCodes`.

Có lẽ quan trọng nhất là
thiết kế này cũng cho phép định nghĩa các khối xử lý lỗi
sẽ chạy khi những lần `check` sau đó thất bại.
Điều đó cho phép bạn viết mã thêm ngữ cảnh dùng chung chỉ một lần,
như trong đoạn sau:

{{raw `
	handle err {
		err = fmt.Errorf("closing writer: %w", err)
	}

	// Write the savedCode if valid.
	if e.savedCode != invalidCode {
		check e.write(e, e.savedCode)
		if err := e.incHi(); err != errOutOfCodes {
			check err
		}
	}

	// Write the eof code.
	eof := uint32(1)<<e.litWidth + 1
	check e.write(e, eof)
`}}

Về bản chất, `check` là một cách viết ngắn của câu lệnh `if`,
còn `handle` thì giống như
[`defer`](/ref/spec#Defer_statements) nhưng chỉ dành cho đường trả về do lỗi.
Khác với exception trong các ngôn ngữ khác,
thiết kế này vẫn giữ lại thuộc tính quan trọng của Go rằng
mọi lời gọi có thể thất bại đều được đánh dấu tường minh trong mã,
lúc này dùng từ khóa `check` thay cho `if err != nil`.

Vấn đề lớn của thiết kế này
là `handle` chồng lấn quá nhiều,
và theo những cách gây bối rối, với `defer`.

Tháng 5 vừa rồi chúng tôi đăng
[một thiết kế mới với ba sự đơn giản hóa](/design/32437-try-builtin):
để tránh sự nhầm lẫn với `defer`, thiết kế bỏ `handle` để dùng thẳng `defer`;
để gần với một ý tưởng tương tự trong Rust và Swift, thiết kế đổi tên `check` thành `try`;
và để cho phép thử nghiệm theo cách mà các trình phân tích sẵn có như `gofmt` có thể nhận ra,
nó đổi `check` (giờ là `try`) từ một từ khóa thành một hàm dựng sẵn.

Giờ cùng đoạn mã đó sẽ như sau:

{{raw `
	defer errd.Wrapf(&err, "closing writer")

	// Write the savedCode if valid.
	if e.savedCode != invalidCode {
		try(e.write(e, e.savedCode))
		if err := e.incHi(); err != errOutOfCodes {
			try(err)
		}
	}

	// Write the eof code.
	eof := uint32(1)<<e.litWidth + 1
	try(e.write(e, eof))
`}}

Chúng tôi dành phần lớn tháng 6 để thảo luận công khai đề xuất này trên GitHub.

Ý tưởng cốt lõi của `check` hay `try` là rút ngắn
lượng cú pháp lặp lại ở mỗi lần kiểm tra lỗi,
đặc biệt là loại bỏ câu lệnh `return` khỏi tầm nhìn,
giữ cho việc kiểm tra lỗi vẫn tường minh và làm nổi bật hơn các biến thể thú vị.
Tuy nhiên, một điểm đáng chú ý được nêu ra trong cuộc thảo luận phản hồi công khai là,
nếu không có câu lệnh `if` và `return` tường minh,
thì không còn chỗ để đặt một lệnh in phục vụ debug,
không còn chỗ để đặt breakpoint,
và cũng không còn đoạn mã nào để công cụ coverage hiển thị là chưa được chạy.
Những lợi ích mà chúng tôi hướng tới
đi kèm cái giá là làm các tình huống này phức tạp hơn.
Xét tổng thể, từ điểm này cũng như các cân nhắc khác,
hoàn toàn không rõ kết quả cuối cùng có thật sự
làm việc phát triển Go đơn giản hơn hay không,
nên chúng tôi đã từ bỏ thử nghiệm này.

Đó là toàn bộ về xử lý lỗi,
vốn là một trong những trọng tâm chính của năm nay.

## Generics

Giờ đến một chủ đề ít gây tranh cãi hơn một chút: generics.

Chủ đề lớn thứ hai mà chúng tôi xác định cho Go 2 là
một cách nào đó để viết mã với
tham số kiểu.
Điều này sẽ cho phép viết các cấu trúc dữ liệu tổng quát
và cả các hàm tổng quát
làm việc với mọi loại slice,
mọi loại channel,
hay mọi loại map.
Ví dụ, đây là một bộ lọc channel tổng quát:

{{raw `
	// Filter copies values from c to the returned channel,
	// passing along only those values satisfying f.
	func Filter(type value)(f func(value) bool, c <-chan value) <-chan value {
		out := make(chan value)
		go func() {
			for v := range c {
				if f(v) {
					out <- v
				}
			}
			close(out)
		}()
		return out
	}
`}}

Chúng tôi đã nghĩ về generics từ khi bắt đầu làm Go,
và từng viết rồi bác bỏ thiết kế cụ thể đầu tiên vào năm 2010.
Chúng tôi còn viết rồi bác bỏ thêm ba thiết kế nữa trước cuối năm 2013.
Bốn thử nghiệm bị bỏ,
nhưng không phải là những thử nghiệm thất bại,
chúng tôi đã học được từ chúng,
giống như đã học được từ `check` và `try`.
Mỗi lần như vậy, chúng tôi học được rằng con đường tới Go 2 không nằm chính xác theo hướng đó,
và chúng tôi nhận ra những hướng khác có thể đáng để khám phá.
Nhưng đến năm 2013, chúng tôi quyết định rằng cần tập trung vào những mối quan tâm khác,
vì vậy toàn bộ chủ đề này được gác sang một bên trong vài năm.

Năm ngoái chúng tôi bắt đầu khám phá và thử nghiệm trở lại,
và chúng tôi đã trình bày
[một thiết kế mới](https://github.com/golang/proposal/blob/master/design/go2draft-contracts.md),
dựa trên ý tưởng về contract,
tại Gophercon mùa hè năm ngoái.
Chúng tôi tiếp tục thử nghiệm và đơn giản hóa,
đồng thời làm việc
với các chuyên gia về lý thuyết ngôn ngữ lập trình
để hiểu thiết kế này rõ hơn.

Nhìn chung, tôi lạc quan rằng chúng ta đang đi đúng hướng,
hướng tới một thiết kế có thể đơn giản hóa việc phát triển Go.
Ngay cả vậy, chúng ta vẫn có thể phát hiện ra rằng thiết kế này cũng không ổn.
Chúng ta có thể lại phải từ bỏ thử nghiệm này
và điều chỉnh con đường dựa trên những gì đã học.
Rồi chúng ta sẽ biết.

Tại Gophercon 2019, Ian Lance Taylor đã nói về
lý do vì sao chúng ta có thể muốn thêm generics vào Go
và đồng thời điểm qua ngắn gọn bản thiết kế nháp mới nhất.
Để biết chi tiết, hãy xem bài blog của anh ấy: “[Why Generics?](/blog/why-generics)”

## Dependencies

Chủ đề lớn thứ ba mà chúng tôi xác định cho Go 2 là quản lý dependency.

Năm 2010 chúng tôi công bố một công cụ tên là `goinstall`,
và gọi nó là
“[một thử nghiệm về cài đặt package](https://groups.google.com/forum/#!msg/golang-nuts/8JFwR3ESjjI/cy7qZzN7Lw4J).”
Nó tải dependency về và lưu chúng
trong cây bộ phân phối Go của bạn, trong GOROOT.

Khi thử nghiệm với `goinstall`,
chúng tôi học được rằng bộ phân phối Go và các package đã cài
nên được giữ tách biệt,
để có thể chuyển sang bộ phân phối Go mới
mà không làm mất toàn bộ package Go của bạn.
Vì vậy đến năm 2011 chúng tôi giới thiệu `GOPATH`,
một biến môi trường chỉ định
nơi tìm các package không có trong bộ phân phối Go chính.

Việc thêm GOPATH tạo ra thêm nhiều nơi chứa package Go
nhưng xét tổng thể đã đơn giản hóa việc phát triển Go,
bằng cách tách bộ phân phối Go của bạn khỏi các thư viện Go của bạn.

**Tính tương thích**

Thử nghiệm `goinstall` cố ý bỏ qua
khái niệm version package tường minh.
Thay vào đó, `goinstall` luôn tải bản mới nhất.
Chúng tôi làm vậy để có thể tập trung vào các
bài toán thiết kế khác của việc cài đặt package.

`Goinstall` trở thành `go get` như một phần của Go 1.
Khi mọi người hỏi về version,
chúng tôi khuyến khích họ thử nghiệm bằng cách
tạo ra thêm công cụ, và họ đã làm vậy.
Và chúng tôi khuyến khích các tác giả package
cung cấp cho người dùng của họ
cùng mức tương thích ngược
mà chúng tôi áp dụng cho các thư viện Go 1.
Trích [Go FAQ](/doc/faq#get_version):

<div style="margin-left: 2em; font-style: italic;">

“Các package được tạo ra để dùng công khai nên cố gắng duy trì tính tương thích ngược trong quá trình phát triển.

Nếu cần chức năng khác,
hãy thêm một tên mới thay vì thay đổi tên cũ.

Nếu cần một sự phá vỡ hoàn toàn,
hãy tạo một package mới với một import path mới.”

</div>

Quy ước này
đơn giản hóa trải nghiệm tổng thể khi sử dụng một package
bằng cách giới hạn những gì tác giả có thể làm:
tránh thay đổi phá vỡ API;
đặt tên mới cho chức năng mới;
và
đặt đường dẫn import mới cho một thiết kế package hoàn toàn mới.

Dĩ nhiên, mọi người vẫn tiếp tục thử nghiệm.
Một trong những thử nghiệm thú vị nhất
do Gustavo Niemeyer khởi xướng.
Anh ấy tạo ra một Git redirector tên là
[`gopkg.in`](https://gopkg.in),
cung cấp các đường dẫn import khác nhau
cho các phiên bản API khác nhau,
để giúp tác giả package
tuân theo quy ước
rằng một thiết kế package mới
phải có một đường dẫn import mới.

Ví dụ,
mã nguồn Go trong repository GitHub
[go-yaml/yaml](https://github.com/go-yaml/yaml)
có các API khác nhau
ở semantic version tag v1 và v2.
Máy chủ `gopkg.in` cung cấp chúng với
các đường dẫn import khác nhau là
[gopkg.in/yaml.v1](https://godoc.org/gopkg.in/yaml.v1)
và
[gopkg.in/yaml.v2](https://godoc.org/gopkg.in/yaml.v2).

Quy ước cung cấp tương thích ngược,
để một phiên bản package mới hơn có thể được dùng
thay cho phiên bản cũ hơn,
chính là điều làm cho quy tắc rất đơn giản của `go get` là “luôn tải bản mới nhất”
vẫn hoạt động tốt cho đến tận ngày nay.

**Version hóa và vendoring**

Nhưng trong bối cảnh production, bạn cần chính xác hơn
về version dependency để bản build có thể tái lập.

Nhiều người đã thử nghiệm xem điều đó nên trông như thế nào,
xây dựng các công cụ phục vụ nhu cầu của họ,
bao gồm `goven` của Keith Rarick (2012) và `godep` (2013),
`glide` của Matt Butcher (2014), và `gb` của Dave Cheney (2015).
Tất cả các công cụ này đều dùng mô hình sao chép package dependency
vào chính repository mã nguồn của bạn.
Cơ chế chính xác được dùng
để đưa các package đó vào tầm import thì khác nhau,
nhưng tất cả đều phức tạp hơn mức tưởng là cần thiết.

Sau một cuộc thảo luận trên toàn cộng đồng,
chúng tôi chấp nhận đề xuất của Keith Rarick
để thêm hỗ trợ tường minh cho việc tham chiếu tới dependency được sao chép
mà không cần mẹo GOPATH.
Đây là đơn giản hóa bằng tái định hình:
giống như với `addToList` và `append`,
các công cụ này vốn đã hiện thực khái niệm đó,
nhưng vụng về hơn mức cần thiết.
Việc thêm hỗ trợ tường minh cho thư mục vendor
đã làm những trường hợp này trở nên đơn giản hơn xét tổng thể.

Việc phát hành hỗ trợ thư mục vendor trong lệnh `go`
dẫn tới nhiều thử nghiệm hơn với chính việc vendoring,
và chúng tôi nhận ra rằng mình đã đưa vào một vài vấn đề.
Nghiêm trọng nhất là chúng tôi làm mất _package uniqueness_.
Trước đó, trong một bản build bất kỳ,
một import path
có thể xuất hiện trong rất nhiều package,
và tất cả các import đó đều trỏ tới cùng một đích.
Giờ với vendoring, cùng một import path trong các
package khác nhau có thể trỏ tới các bản sao vendored khác nhau của package,
tất cả đều sẽ xuất hiện trong binary tạo ra cuối cùng.

Khi đó, chúng tôi chưa có tên cho thuộc tính này:
package uniqueness.
Nó chỉ đơn giản là cách mô hình GOPATH hoạt động.
Chúng tôi không hoàn toàn đánh giá đúng tầm quan trọng của nó cho đến khi nó biến mất.

Ở đây có một sự tương đồng với các đề xuất cú pháp lỗi `check` và `try`.
Trong trường hợp đó, chúng tôi đang dựa vào
cách câu lệnh `return` hiển thị trước mắt hoạt động
theo những cách mà chúng tôi không thật sự nhận ra
cho đến khi cân nhắc việc bỏ nó đi.

Khi chúng tôi thêm hỗ trợ cho thư mục vendor,
đã tồn tại rất nhiều công cụ quản lý dependency khác nhau.
Chúng tôi nghĩ rằng nếu có một thỏa thuận rõ ràng
về định dạng thư mục vendor
và metadata của vendoring,
thì các công cụ khác nhau sẽ có thể tương tác với nhau,
giống như sự thống nhất về
cách chương trình Go được lưu trong tệp văn bản
cho phép trình biên dịch Go, trình soạn thảo văn bản
và các công cụ như `goimports` và `gorename`
tương tác với nhau.

Điều này hóa ra đã quá lạc quan.
Các công cụ vendoring đều khác nhau ở những sắc thái ngữ nghĩa tinh vi.
Muốn tương tác được sẽ phải thay đổi tất cả chúng
để cùng thống nhất ngữ nghĩa,
rất có thể làm hỏng người dùng hiện tại của chúng.
Sự hội tụ đã không xảy ra.

**Dep**

Tại Gophercon 2016, chúng tôi bắt đầu nỗ lực
định nghĩa một công cụ duy nhất để quản lý dependency.
Là một phần của nỗ lực đó, chúng tôi tiến hành khảo sát
với nhiều kiểu người dùng khác nhau
để hiểu họ cần gì
trong quản lý dependency,
và một nhóm bắt tay vào làm một công cụ mới,
sau này trở thành `dep`.

`Dep` nhắm tới việc có thể thay thế mọi công cụ
quản lý dependency hiện có.
Mục tiêu là đơn giản hóa bằng tái định hình các
công cụ khác nhau hiện hữu thành một công cụ duy nhất.
Nó đã phần nào làm được điều đó.
`Dep` cũng khôi phục package uniqueness cho người dùng của nó,
bằng cách chỉ có một thư mục vendor
ở đỉnh cây dự án.

Nhưng `dep` cũng đưa vào một vấn đề nghiêm trọng
mà chúng tôi mất một thời gian mới đánh giá hết được.
Vấn đề là `dep` chấp nhận một lựa chọn thiết kế từ `glide`,
đó là hỗ trợ và khuyến khích những thay đổi không tương thích đối với một package
mà không đổi import path.

Đây là một ví dụ.
Giả sử bạn đang xây dựng chương trình của riêng mình,
và bạn cần có một tệp cấu hình,
nên bạn dùng phiên bản 2 của một gói Go YAML phổ biến:

<div style="margin-left: 2em;">
{{image "experiment/yamldeps1.png" 214}}
</div>

Giờ giả sử chương trình của bạn
import Kubernetes client.
Hóa ra Kubernetes dùng YAML rất nhiều,
và nó dùng phiên bản 1 của cùng gói phổ biến đó:

<div style="margin-left: 2em;">
{{image "experiment/yamldeps2.png" 557}}
</div>

Phiên bản 1 và phiên bản 2 có API không tương thích,
nhưng chúng cũng có import path khác nhau,
nên không có sự mơ hồ nào về việc một import cụ thể ám chỉ cái nào.
Kubernetes lấy phiên bản 1,
bộ phân tích cấu hình của bạn lấy phiên bản 2,
và mọi thứ đều hoạt động.

`Dep` đã từ bỏ mô hình này.
Phiên bản 1 và phiên bản 2 của gói yaml giờ đây
sẽ có cùng import path,
gây ra xung đột.
Việc dùng cùng import path cho hai phiên bản không tương thích,
kết hợp với package uniqueness,
khiến việc build chương trình này trở nên bất khả thi,
trong khi trước đó bạn có thể build được:

<div style="margin-left: 2em;">
{{image "experiment/yamldeps3.png" 450}}
</div>

Chúng tôi mất một thời gian mới hiểu được vấn đề này,
vì đã áp dụng quy ước
“API mới nghĩa là import path mới”
quá lâu nên chúng tôi mặc nhiên coi đó là điều hiển nhiên.
Thử nghiệm dep đã giúp chúng tôi
đánh giá đúng hơn quy ước đó,
và chúng tôi đặt cho nó một cái tên:
_import compatibility rule_:

<div style="margin-left: 2em; font-style: italic;">

“Nếu một package cũ và một package mới có cùng import path,
thì package mới phải tương thích ngược với package cũ.”

</div>

**Go Modules**

Chúng tôi lấy những gì hoạt động tốt trong thử nghiệm dep
và những gì học được về những điểm không hoạt động tốt,
rồi thử nghiệm một thiết kế mới tên là `vgo`.
Trong `vgo`, các package tuân theo import compatibility rule,
để chúng ta có thể vừa đảm bảo package uniqueness
vừa không làm hỏng các bản build như ví dụ vừa rồi.
Điều này cũng cho phép chúng tôi đơn giản hóa
các phần khác của thiết kế.

Ngoài việc khôi phục import compatibility rule,
một phần quan trọng khác của thiết kế `vgo`
là đặt tên cho khái niệm một nhóm package
và cho phép nhóm đó được tách ra
khỏi ranh giới của repository mã nguồn.
Tên của một nhóm package Go là module,
vì vậy giờ chúng tôi gọi hệ thống này là Go modules.

Go modules hiện đã được tích hợp với lệnh `go`,
nhờ đó không cần phải sao chép thư mục vendor qua lại nữa.

**Thay thế GOPATH**

Go modules kéo theo sự kết thúc của GOPATH với vai trò
một không gian tên toàn cục.
Gần như toàn bộ công việc khó khăn khi chuyển cách dùng Go hiện có
và công cụ hiện có sang modules đều do thay đổi này gây ra,
tức là do rời bỏ GOPATH.

Ý tưởng nền tảng của GOPATH
là cây thư mục GOPATH
là nguồn chân lý toàn cục
cho việc đang dùng những phiên bản nào,
và những phiên bản đó không đổi
khi bạn di chuyển giữa các thư mục.
Nhưng chế độ GOPATH toàn cục xung đột trực tiếp
với yêu cầu trong production là
bản build tái lập theo từng dự án,
bản thân điều này lại đơn giản hóa
trải nghiệm phát triển và triển khai Go theo nhiều cách quan trọng.

Bản build tái lập theo từng dự án có nghĩa là
khi bạn đang làm việc trong bản checkout của dự án A,
bạn nhận đúng bộ phiên bản dependency mà các lập trình viên khác của dự án A nhận được
tại commit đó,
như được định nghĩa bởi tệp `go.mod`.
Khi bạn chuyển sang làm việc trong bản checkout của dự án B,
lúc này bạn sẽ nhận bộ phiên bản dependency do dự án đó chọn,
cũng chính là bộ mà những người khác của dự án B nhận được.
Nhưng bộ đó rất có thể khác với dự án A.
Việc bộ phiên bản dependency
thay đổi khi bạn chuyển từ dự án A sang dự án B
là điều cần thiết để giữ cho quá trình phát triển của bạn đồng bộ
với quá trình phát triển của những người khác ở A và ở B.
Không thể còn một GOPATH toàn cục duy nhất nữa.

Phần lớn độ phức tạp của việc chấp nhận modules
phát sinh trực tiếp từ việc mất đi GOPATH toàn cục duy nhất đó.
Mã nguồn của một package nằm ở đâu?
Trước đây, câu trả lời chỉ phụ thuộc vào biến môi trường GOPATH của bạn,
thứ mà phần lớn mọi người hiếm khi thay đổi.
Giờ đây, câu trả lời phụ thuộc vào dự án bạn đang làm,
thứ có thể thay đổi thường xuyên.
Mọi thứ đều cần được cập nhật cho quy ước mới này.

Hầu hết công cụ phát triển dùng gói
[`go/build`](https://godoc.org/go/build) để tìm và tải mã nguồn Go.
Chúng tôi đã giữ cho gói đó tiếp tục hoạt động,
nhưng API của nó không được thiết kế với modules trong đầu,
và những cách vòng mà chúng tôi thêm vào để tránh thay đổi API
chạy chậm hơn mức chúng tôi mong muốn.
Chúng tôi đã phát hành một gói thay thế,
[`golang.org/x/tools/go/packages`](https://godoc.org/golang.org/x/tools/go/packages).
Công cụ phát triển giờ nên dùng gói đó thay thế.
Nó hỗ trợ cả GOPATH lẫn Go modules,
đồng thời nhanh hơn và dễ dùng hơn.
Trong một hoặc hai bản phát hành tới chúng tôi có thể chuyển nó vào thư viện chuẩn,
nhưng hiện tại [`golang.org/x/tools/go/packages`](https://godoc.org/golang.org/x/tools/go/packages)
đã ổn định và sẵn sàng để sử dụng.

**Module proxy của Go**

Một trong những cách mà modules đơn giản hóa việc phát triển Go
là tách khái niệm một nhóm package
khỏi repository quản lý mã nguồn bên dưới
nơi chúng được lưu trữ.

Khi chúng tôi trò chuyện với người dùng Go về dependency,
gần như mọi người dùng Go trong công ty của họ
đều hỏi làm thế nào để định tuyến các lần `go get` tải package
qua máy chủ riêng của họ,
để kiểm soát tốt hơn loại mã nào được phép dùng.
Và ngay cả các lập trình viên nguồn mở cũng lo ngại
về việc dependency biến mất
hoặc thay đổi bất ngờ,
làm hỏng các bản build của họ.
Trước modules, người dùng từng thử
những giải pháp phức tạp cho các vấn đề này,
bao gồm cả việc chặn các lệnh của hệ quản lý phiên bản
mà lệnh `go` chạy.

Thiết kế Go modules giúp việc
đưa vào khái niệm module proxy trở nên dễ dàng,
một proxy có thể được yêu cầu cung cấp một phiên bản module cụ thể.

Giờ đây các công ty có thể dễ dàng chạy module proxy của riêng mình,
với các quy tắc tùy biến về những gì được phép
và nơi lưu các bản sao cache.
Dự án nguồn mở [Athens](https://docs.gomods.io) đã xây dựng một proxy như vậy,
và Aaron Schlesinger đã có một bài nói về nó tại Gophercon 2019.
(Chúng tôi sẽ thêm liên kết ở đây khi video có sẵn.)

Và đối với các lập trình viên cá nhân cùng các nhóm nguồn mở,
nhóm Go tại Google đã [ra mắt một proxy](https://groups.google.com/forum/#!topic/golang-announce/0wo8cOhGuAI)
đóng vai trò là bản phản chiếu công khai của mọi package Go nguồn mở,
và Go 1.13 sẽ dùng proxy đó mặc định khi ở module mode.
Katie Hockman đã có [một bài nói về hệ thống này tại Gophercon 2019](https://youtu.be/KqTySYYhPUE).

**Tình trạng của Go modules**

Go 1.11 giới thiệu modules như một bản xem trước thử nghiệm, phải chủ động bật.
Chúng tôi tiếp tục thử nghiệm và đơn giản hóa.
Go 1.12 phát hành các cải tiến,
và Go 1.13 sẽ phát hành thêm nhiều cải tiến nữa.

Modules giờ đã tới điểm
mà chúng tôi tin rằng chúng sẽ phục vụ được phần lớn người dùng,
nhưng chúng tôi chưa sẵn sàng tắt GOPATH ngay lúc này.
Chúng tôi sẽ tiếp tục thử nghiệm, đơn giản hóa và hiệu chỉnh.

Chúng tôi hoàn toàn nhận thức được rằng
cộng đồng người dùng Go
đã tích lũy gần một thập kỷ kinh nghiệm
cũng như công cụ và workflow xoay quanh GOPATH,
và sẽ cần thời gian để chuyển tất cả những thứ đó sang Go modules.

Nhưng một lần nữa,
chúng tôi nghĩ rằng modules giờ đây
sẽ hoạt động rất tốt cho đa số người dùng,
và tôi khuyến khích bạn xem thử
khi Go 1.13 được phát hành.

Như một điểm dữ liệu,
dự án Kubernetes có rất nhiều dependency,
và họ đã chuyển sang dùng Go modules
để quản lý chúng.
Có lẽ bạn cũng có thể.
Và nếu bạn không thể,
hãy cho chúng tôi biết điều gì không hoạt động với bạn
hoặc điều gì quá phức tạp,
bằng cách [gửi bug report](/issue/new),
và chúng tôi sẽ thử nghiệm rồi đơn giản hóa.

## Công cụ

Xử lý lỗi, generics và quản lý dependency
sẽ còn cần ít nhất vài năm nữa,
và hiện tại chúng tôi sẽ tập trung vào chúng.
Xử lý lỗi đã gần xong,
tiếp theo sẽ là modules,
và có thể sau đó là generics.

Nhưng giả sử chúng ta nhìn xa hơn vài năm,
tới lúc chúng ta đã thử nghiệm, đơn giản hóa
và phát hành xong xử lý lỗi, modules và generics.
Rồi thì sao?
Rất khó để dự đoán tương lai,
nhưng tôi nghĩ rằng một khi ba mảng này đã được phát hành,
đó có thể đánh dấu sự bắt đầu của một giai đoạn yên ắng mới đối với các thay đổi lớn.
Khi đó trọng tâm của chúng tôi có lẽ sẽ chuyển sang
đơn giản hóa việc phát triển Go bằng các công cụ tốt hơn.

Một phần công việc về công cụ đã bắt đầu rồi,
nên bài viết này kết thúc bằng việc nhìn vào hướng đó.

Trong lúc chúng tôi giúp cập nhật tất cả các công cụ hiện có của cộng đồng Go
để hiểu Go modules,
chúng tôi nhận ra rằng việc có quá nhiều công cụ trợ giúp phát triển
mỗi công cụ chỉ làm một việc nhỏ
không còn phục vụ người dùng tốt nữa.
Từng công cụ riêng lẻ quá khó kết hợp,
quá chậm để gọi,
và cách sử dụng thì quá khác nhau.

Chúng tôi bắt đầu một nỗ lực nhằm hợp nhất những công cụ trợ giúp phát triển
được cần đến thường xuyên nhất vào một công cụ duy nhất,
giờ được gọi là `gopls` (phát âm như “go, please”).
`Gopls` nói
[Language Server Protocol, LSP](https://langserver.org/),
và làm việc với bất kỳ môi trường phát triển tích hợp
hoặc trình soạn thảo văn bản nào có hỗ trợ LSP,
về cơ bản là hầu như mọi thứ vào thời điểm này.

`Gopls` đánh dấu sự mở rộng trọng tâm của dự án Go,
từ chỗ chỉ cung cấp các công cụ kiểu compiler chạy dòng lệnh độc lập
như go vet hoặc gorename
sang chỗ cũng cung cấp một dịch vụ IDE hoàn chỉnh.
Rebecca Stambler đã có một bài nói với nhiều chi tiết hơn về `gopls` và IDE tại Gophercon 2019.
(Chúng tôi sẽ thêm liên kết ở đây khi video có sẵn.)

Sau `gopls`, chúng tôi cũng có các ý tưởng để hồi sinh `go fix`
theo một cách có thể mở rộng
và để làm cho `go vet` hữu ích hơn nữa.

## Kết

<div style="margin-left: 2em;">
{{image "experiment/expsimp2.png" 326}}
</div>

Vậy đó là con đường tới Go 2.
Chúng ta sẽ thử nghiệm rồi đơn giản hóa.
Rồi lại thử nghiệm và đơn giản hóa.
Rồi phát hành.
Rồi lại thử nghiệm và đơn giản hóa.
Rồi lại làm tất cả một lần nữa.
Nó có thể trông, hoặc thậm chí có cảm giác, như thể con đường đang đi vòng tròn.
Nhưng mỗi lần chúng ta thử nghiệm và đơn giản hóa
chúng ta lại học thêm một chút về hình dạng mà Go 2 nên có
và tiến thêm một bước nữa đến gần nó.
Ngay cả những thử nghiệm bị bỏ như `try`
hay bốn thiết kế generics đầu tiên của chúng tôi
hay `dep` cũng không phải là thời gian lãng phí.
Chúng giúp chúng tôi học được những gì cần được
đơn giản hóa trước khi có thể phát hành,
và trong một số trường hợp còn giúp chúng tôi hiểu rõ hơn
những điều mà trước đây chúng tôi mặc nhiên coi là hiển nhiên.

Đến một lúc nào đó, chúng tôi sẽ nhận ra rằng mình đã
thử nghiệm đủ,
đơn giản hóa đủ,
và phát hành đủ,
và khi đó chúng ta sẽ có Go 2.

Cảm ơn tất cả các bạn trong cộng đồng Go
đã giúp chúng tôi thử nghiệm
và đơn giản hóa
và phát hành
và tìm đường đi trên hành trình này.
