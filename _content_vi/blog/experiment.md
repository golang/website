---
title: Thử nghiệm, Đơn giản hóa, Phát hành
date: 2019-08-01
by:
- Russ Cox
tags:
- community
- go2
- proposals
summary: Phiên bản bài blog của bài nói chuyện của tôi tại GopherCon 2019.
template: true
---

## Giới thiệu

Đây là phiên bản bài blog của bài nói chuyện của tôi vào tuần trước tại GopherCon 2019.

{{video "https://www.youtube.com/embed/kNHo788oO5Y?rel=0"}}

Tất cả chúng ta đều đang cùng nhau bước trên con đường đến Go 2,
nhưng không ai trong chúng ta biết chính xác con đường đó dẫn tới đâu
hay thậm chí đôi khi nó đang đi theo hướng nào.
Bài viết này bàn về cách chúng ta thực sự
tìm ra và đi theo con đường đến Go 2.
Đây là hình dạng của quy trình đó.

<div style="margin-left: 2em;">
{{image "experiment/expsimp1.png" 179}}
</div>

Chúng ta thử nghiệm với Go như nó đang tồn tại hiện nay,
để hiểu nó rõ hơn,
học xem điều gì hoạt động tốt và điều gì không.
Sau đó chúng ta thử nghiệm với những thay đổi khả dĩ,
để hiểu chúng rõ hơn,
và một lần nữa học xem điều gì hoạt động tốt và điều gì không.
Dựa trên những gì học được từ các thử nghiệm đó,
chúng ta đơn giản hóa.
Rồi chúng ta lại thử nghiệm.
Rồi lại đơn giản hóa.
Cứ thế tiếp diễn.
Và cứ thế tiếp diễn.

## Bốn chữ R của việc Đơn giản hóa

Trong quá trình này, có bốn cách chính để chúng ta có thể đơn giản hóa
trải nghiệm tổng thể của việc viết chương trình Go:
định hình lại, định nghĩa lại, loại bỏ và giới hạn.

**Đơn giản hóa bằng cách Định hình lại**

Cách đầu tiên để đơn giản hóa là định hình lại những gì đang có sang một dạng mới,
một dạng rốt cuộc đơn giản hơn về tổng thể.

Mỗi chương trình Go chúng ta viết đều là một thử nghiệm để kiểm tra chính Go.
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

Chúng tôi sẽ viết cùng một đoạn mã cho slice của byte,
và slice của string, v.v.
Các chương trình của chúng tôi quá phức tạp, vì Go khi đó lại quá đơn giản.

Vì vậy chúng tôi đã lấy rất nhiều hàm như `addToList` trong các chương trình
và định hình lại chúng thành một hàm duy nhất do chính Go cung cấp.
Việc thêm `append` khiến ngôn ngữ Go phức tạp hơn một chút,
nhưng xét trên cán cân tổng thể
thì nó làm cho trải nghiệm viết chương trình Go trở nên đơn giản hơn,
ngay cả khi đã tính đến chi phí học về `append`.

Đây là một ví dụ khác.
Đối với Go 1, chúng tôi xem xét rất nhiều công cụ phát triển
trong bộ phân phối Go, và định hình lại chúng thành một lệnh mới duy nhất.

	5a      8g
	5g      8l
	5l      cgo
	6a      gobuild
	6cov    gofix         →     go
	6g      goinstall
	6l      gomake
	6nm     gopack
	8a      govet

Lệnh `go` giờ đây đã trở nên quá trung tâm đến mức
người ta dễ quên rằng đã từng có một thời gian dài chúng ta không có nó
và điều đó kéo theo nhiều công sức phụ trội như thế nào.

Chúng tôi đã thêm mã và độ phức tạp vào bộ phân phối Go,
nhưng xét trên tổng thể chúng tôi đã đơn giản hóa trải nghiệm viết chương trình Go.
Cấu trúc mới cũng tạo ra không gian cho những thử nghiệm thú vị khác,
mà chúng ta sẽ thấy sau đây.

**Đơn giản hóa bằng cách Định nghĩa lại**

Cách thứ hai để đơn giản hóa là định nghĩa lại
chức năng mà chúng ta đã có,
cho phép nó làm được nhiều hơn.
Giống như đơn giản hóa bằng cách định hình lại,
đơn giản hóa bằng cách định nghĩa lại làm cho chương trình dễ viết hơn,
nhưng lần này không cần học thêm điều gì mới.

Ví dụ, ban đầu `append` được định nghĩa là chỉ đọc từ slice.
Khi append vào một byte slice, bạn có thể append các byte từ một byte slice khác,
nhưng không thể append các byte từ một string.
Chúng tôi đã định nghĩa lại append để cho phép append từ string,
mà không cần thêm gì mới vào ngôn ngữ.

	var b []byte
	var more []byte
	b = append(b, more...) // ok

	var b []byte
	var more string
	b = append(b, more...) // ok later

**Đơn giản hóa bằng cách Loại bỏ**

Cách thứ ba để đơn giản hóa là loại bỏ chức năng
khi nó hóa ra kém hữu ích
hoặc kém quan trọng hơn chúng tôi từng kỳ vọng.
Loại bỏ chức năng có nghĩa là bớt đi một thứ phải học,
bớt đi một thứ phải sửa lỗi,
bớt đi một thứ gây xao nhãng hoặc bị dùng sai.
Dĩ nhiên, việc loại bỏ cũng
buộc người dùng phải cập nhật các chương trình hiện có,
có thể làm chúng phức tạp hơn,
để bù lại phần đã bị loại bỏ.
Nhưng kết quả tổng thể vẫn có thể là
quá trình viết chương trình Go trở nên đơn giản hơn.

Một ví dụ của điều này là khi chúng tôi loại bỏ
các dạng boolean của thao tác kênh không chặn khỏi ngôn ngữ:

{{raw `
	ok := c <- x  // trước Go 1, là gửi không chặn
	x, ok := <-c  // trước Go 1, là nhận không chặn
`}}

Những thao tác này cũng có thể được thực hiện bằng `select`,
khiến việc phải quyết định dùng dạng nào trở nên rối rắm.
Loại bỏ chúng đã đơn giản hóa ngôn ngữ mà không làm giảm sức mạnh của nó.

**Đơn giản hóa bằng cách Giới hạn**

Chúng ta cũng có thể đơn giản hóa bằng cách giới hạn những gì được phép.
Ngay từ ngày đầu tiên, Go đã giới hạn cách mã hóa của tệp mã nguồn Go:
chúng phải là UTF-8.
Giới hạn này làm cho mọi chương trình cố đọc tệp mã nguồn Go trở nên đơn giản hơn.
Những chương trình đó không phải lo về các tệp mã nguồn Go
được mã hóa theo Latin-1 hay UTF-16 hay UTF-7 hay bất kỳ thứ gì khác.

Một giới hạn quan trọng khác là `gofmt` cho việc định dạng chương trình.
Không có gì từ chối mã Go không được định dạng bằng `gofmt`,
nhưng chúng tôi đã xây dựng một quy ước rằng các công cụ viết lại chương trình Go
sẽ để chúng ở dạng `gofmt`.
Nếu bạn cũng giữ các chương trình của mình ở dạng `gofmt`,
thì các công cụ viết lại này sẽ không tạo ra thay đổi định dạng nào.
Khi bạn so sánh trước và sau,
những khác biệt duy nhất bạn thấy sẽ là các thay đổi thực sự.
Giới hạn này đã đơn giản hóa các công cụ viết lại chương trình
và dẫn tới những thử nghiệm thành công như
`goimports`, `gorename`, và nhiều công cụ khác.

## Quy trình Phát triển Go

Chu trình thử nghiệm và đơn giản hóa này là một mô hình tốt cho những gì chúng tôi đã làm trong mười năm qua.
nhưng nó có một vấn đề:
nó quá đơn giản.
Chúng ta không thể chỉ thử nghiệm và đơn giản hóa.

Chúng ta phải phát hành kết quả.
Chúng ta phải cung cấp nó để mọi người sử dụng.
Dĩ nhiên, việc sử dụng nó cho phép có thêm các thử nghiệm,
và có thể có thêm sự đơn giản hóa,
và quy trình cứ thế lặp đi lặp lại.

<div style="margin-left: 2em;">
{{image "experiment/expsimp2.png" 326}}
</div>

Chúng tôi đã phát hành Go tới tất cả các bạn lần đầu tiên
vào ngày 10 tháng 11 năm 2009.
Sau đó, cùng với sự giúp đỡ của các bạn, chúng tôi đã phát hành Go 1 cùng nhau vào tháng 3 năm 2012.
Và từ đó tới nay chúng tôi đã phát hành thêm mười hai phiên bản Go.
Tất cả những mốc này đều rất quan trọng,
để mở đường cho thêm nhiều thử nghiệm,
để giúp chúng tôi hiểu rõ hơn về Go,
và dĩ nhiên để đưa Go vào sử dụng trong môi trường production.

Khi chúng tôi phát hành Go 1,
chúng tôi đã chủ động chuyển trọng tâm sang việc sử dụng Go,
để hiểu phiên bản này của ngôn ngữ sâu hơn nhiều
trước khi thử thêm bất kỳ sự đơn giản hóa nào liên quan tới
thay đổi ngôn ngữ.
Chúng tôi cần dành thời gian để thử nghiệm,
để thực sự hiểu điều gì hoạt động và điều gì không.

Dĩ nhiên, chúng tôi đã có mười hai bản phát hành kể từ Go 1,
vì vậy chúng tôi vẫn tiếp tục thử nghiệm, đơn giản hóa và phát hành.
Nhưng chúng tôi tập trung vào những cách đơn giản hóa việc phát triển Go
mà không có thay đổi ngôn ngữ đáng kể và không phá vỡ
các chương trình Go hiện có.
Ví dụ, Go 1.5 phát hành bộ gom rác đồng thời đầu tiên
và các bản phát hành sau đó tiếp tục cải thiện nó,
đơn giản hóa việc phát triển Go bằng cách loại bỏ thời gian dừng như một mối lo thường trực.

Tại Gophercon năm 2017, chúng tôi đã thông báo rằng sau năm năm
thử nghiệm, đã đến lúc
nghĩ về
những thay đổi đáng kể có thể đơn giản hóa việc phát triển Go.
Con đường tới Go 2 của chúng ta thực chất vẫn giống con đường tới Go 1:
thử nghiệm rồi đơn giản hóa rồi phát hành,
hướng đến mục tiêu tổng thể là làm cho việc phát triển Go trở nên đơn giản hơn.

Đối với Go 2, các chủ đề cụ thể mà chúng tôi tin là
quan trọng nhất cần giải quyết là
xử lý lỗi, generics và dependencies.
Kể từ đó chúng tôi nhận ra rằng còn có một
chủ đề quan trọng khác là công cụ dành cho lập trình viên.

Phần còn lại của bài viết này sẽ bàn về cách
công việc của chúng tôi trong từng lĩnh vực này
đi theo con đường đó.
Trên đường đi,
chúng ta sẽ có một đoạn rẽ nhỏ,
dừng lại để xem xét chi tiết kỹ thuật
của những gì sắp được phát hành trong Go 1.13
cho việc xử lý lỗi.

## Lỗi

Viết một chương trình
hoạt động đúng trong mọi trường hợp
đã đủ khó
khi mọi đầu vào đều hợp lệ và chính xác
và không có gì mà chương trình phụ thuộc vào bị lỗi.
Khi bạn đưa lỗi vào bài toán,
việc viết một chương trình hoạt động đúng
bất kể chuyện gì xảy ra lại càng khó hơn.

Trong quá trình suy nghĩ về Go 2,
chúng tôi muốn hiểu rõ hơn
liệu Go có thể giúp làm cho công việc đó đơn giản hơn phần nào hay không.

Có hai khía cạnh khác nhau có thể
được đơn giản hóa:
giá trị lỗi và cú pháp lỗi.
Chúng ta sẽ lần lượt xem từng khía cạnh,
với đoạn rẽ kỹ thuật mà tôi đã hứa tập trung
vào các thay đổi giá trị lỗi trong Go 1.13.

**Giá trị lỗi**

Giá trị lỗi phải bắt đầu từ đâu đó.
Đây là hàm `Read` từ phiên bản đầu tiên của gói `os`:

	export func Read(fd int64, b *[]byte) (ret int64, errno int64) {
		r, e := syscall.read(fd, &b[0], int64(len(b)));
		return r, e
	}

Khi đó chưa có kiểu `File`, và cũng chưa có kiểu lỗi.
`Read` cùng các hàm khác trong gói
trả về trực tiếp một `errno int64` từ lời gọi hệ thống Unix bên dưới.

Đoạn mã này được commit vào ngày 10 tháng 9 năm 2008 lúc 12:14 chiều.
Giống như mọi thứ khi đó, nó là một thử nghiệm,
và mã thay đổi rất nhanh.
Hai giờ năm phút sau, API đã thay đổi:

	export type Error struct { s string }

	func (e *Error) Print() { … } // ra standard error!
	func (e *Error) String() string { … }

	export func Read(fd int64, b *[]byte) (ret int64, err *Error) {
		r, e := syscall.read(fd, &b[0], int64(len(b)));
		return r, ErrnoToError(e)
	}

API mới này giới thiệu kiểu `Error` đầu tiên.
Một lỗi giữ một chuỗi và có thể trả về chuỗi đó
đồng thời cũng có thể in nó ra standard error.

Mục đích ở đây là khái quát hóa vượt ra ngoài mã số nguyên.
Chúng tôi biết từ kinh nghiệm trước đó
rằng các số lỗi của hệ điều hành là một biểu diễn quá hạn chế,
rằng việc không phải nhét
mọi chi tiết về một lỗi vào 64 bit sẽ làm chương trình đơn giản hơn.
Việc dùng chuỗi lỗi đã hoạt động khá ổn
cho chúng tôi trong quá khứ, nên chúng tôi làm điều tương tự ở đây.
API mới này tồn tại trong bảy tháng.

Vào tháng Tư năm sau, sau khi có thêm kinh nghiệm với interfaces,
chúng tôi quyết định khái quát hóa thêm nữa
và cho phép người dùng tự định nghĩa cách cài đặt lỗi,
bằng cách biến chính kiểu `os.Error` thành một interface.
Chúng tôi đơn giản hóa bằng cách loại bỏ phương thức `Print`.

Đến Go 1, hai năm sau đó,
dựa trên một gợi ý của Roger Peppe,
`os.Error` trở thành kiểu dựng sẵn `error`,
và phương thức `String` được đổi tên thành `Error`.
Từ đó đến nay không có gì thay đổi.
Nhưng chúng tôi đã viết rất nhiều chương trình Go,
và nhờ đó đã thử nghiệm rất nhiều về cách
tốt nhất để cài đặt và sử dụng lỗi.

**Lỗi là giá trị**

Việc biến `error` thành một interface đơn giản
và cho phép có nhiều cách cài đặt khác nhau
có nghĩa là toàn bộ ngôn ngữ Go
đều sẵn sàng để định nghĩa và kiểm tra lỗi.
Chúng tôi thích nói rằng [lỗi là giá trị](/blog/errors-are-values),
giống như bất kỳ giá trị Go nào khác.

Đây là một ví dụ.
Trên Unix,
một nỗ lực kết nối mạng
rốt cuộc sẽ dùng lời gọi hệ thống `connect`.
Lời gọi hệ thống đó trả về một `syscall.Errno`,
là một kiểu số nguyên có tên đại diện cho
mã lỗi lời gọi hệ thống
và cài đặt interface `error`:

	package syscall

	type Errno int64

	func (e Errno) Error() string { ... }

	const ECONNREFUSED = Errno(61)

	    ... err == ECONNREFUSED ...

Gói `syscall` cũng định nghĩa các hằng có tên
cho những số lỗi do hệ điều hành máy chủ định nghĩa.
Trong trường hợp này, trên hệ thống này, `ECONNREFUSED` là số 61.
Mã nhận một lỗi từ một hàm
có thể kiểm tra xem lỗi có phải là `ECONNREFUSED`
bằng [phép so sánh giá trị](/ref/spec#Comparison_operators) thông thường.

Tiến lên một mức,
trong gói `os`,
bất kỳ lỗi lời gọi hệ thống nào cũng được báo cáo bằng
một cấu trúc lỗi lớn hơn ghi lại
thao tác đã được thử thực hiện ngoài chính lỗi đó.
Có một vài cấu trúc như vậy.
Cấu trúc này, `SyscallError`, mô tả một lỗi
khi gọi một lời gọi hệ thống cụ thể
không ghi nhận thêm thông tin nào khác:

	package os

	type SyscallError struct {
		Syscall string
		Err     error
	}

	func (e *SyscallError) Error() string {
		return e.Syscall + ": " + e.Err.Error()
	}

Tiến lên thêm một mức nữa,
trong gói `net`,
bất kỳ lỗi mạng nào cũng được báo cáo bằng một
cấu trúc lỗi còn lớn hơn ghi lại chi tiết
về thao tác mạng xung quanh,
chẳng hạn như dial hay listen,
và mạng cùng các địa chỉ liên quan:

	package net

	type OpError struct {
		Op     string
		Net    string
		Source Addr
		Addr   Addr
		Err    error
	}

	func (e *OpError) Error() string { ... }

Ghép tất cả lại với nhau,
các lỗi do những thao tác như `net.Dial` trả về có thể được định dạng thành chuỗi,
nhưng chúng cũng là các giá trị dữ liệu Go có cấu trúc.
Trong trường hợp này, lỗi là một `net.OpError`, bổ sung ngữ cảnh
cho một `os.SyscallError`, thứ này lại bổ sung ngữ cảnh cho một `syscall.Errno`:

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

Khi chúng tôi nói lỗi là giá trị, chúng tôi muốn nói cả hai điều:
toàn bộ ngôn ngữ Go đều có thể được dùng để định nghĩa chúng
và
toàn bộ ngôn ngữ Go đều có thể được dùng để kiểm tra chúng.

Đây là một ví dụ từ gói net.
Hóa ra là khi bạn thử kết nối socket,
phần lớn thời gian bạn hoặc sẽ kết nối được hoặc sẽ nhận connection refused,
nhưng đôi khi bạn có thể gặp một `EADDRNOTAVAIL` giả,
không vì lý do chính đáng nào cả.
Go bảo vệ chương trình người dùng khỏi kiểu lỗi này bằng cách thử lại.
Để làm vậy, nó phải kiểm tra cấu trúc lỗi để tìm xem
liệu `syscall.Errno` nằm sâu bên trong có phải là `EADDRNOTAVAIL` hay không.

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

Một [type assertion](/ref/spec#Type_assertions) bóc đi lớp bọc `net.OpError`.
Sau đó một type assertion thứ hai bóc đi lớp bọc `os.SyscallError`.
Rồi hàm kiểm tra lỗi đã được bóc để so sánh bằng với `EADDRNOTAVAIL`.

Điều mà chúng tôi học được sau nhiều năm kinh nghiệm,
từ việc thử nghiệm với lỗi trong Go,
là việc có thể định nghĩa
các cách cài đặt tùy ý của interface `error` là cực kỳ mạnh mẽ,
rằng toàn bộ ngôn ngữ Go đều sẵn sàng
cho cả việc tạo ra lẫn tháo gỡ lỗi,
và không cần bắt buộc phải dùng
bất kỳ một cách cài đặt duy nhất nào.

Những thuộc tính này, rằng lỗi là giá trị,
và rằng không có một cách cài đặt lỗi bắt buộc duy nhất, là
những điều quan trọng cần được giữ lại.

Việc không áp đặt một cách cài đặt lỗi duy nhất
đã cho phép mọi người thử nghiệm
những chức năng bổ sung mà một lỗi có thể cung cấp,
dẫn đến nhiều gói,
chẳng hạn như
[github.com/pkg/errors](https://godoc.org/github.com/pkg/errors),
[gopkg.in/errgo.v2](https://godoc.org/gopkg.in/errgo.v2),
[github.com/hashicorp/errwrap](https://godoc.org/github.com/hashicorp/errwrap),
[upspin.io/errors](https://godoc.org/upspin.io/errors),
[github.com/spacemonkeygo/errors](https://godoc.org/github.com/spacemonkeygo/errors),
và nhiều hơn nữa.

Tuy vậy, một vấn đề của việc thử nghiệm không bị ràng buộc
là với tư cách người dùng thư viện
bạn phải lập trình dựa trên hợp của
tất cả những cách cài đặt có thể gặp phải.
Một sự đơn giản hóa có vẻ đáng để khám phá cho Go 2
là định nghĩa một phiên bản chuẩn của chức năng thường được thêm vào,
dưới dạng các interface tùy chọn đã được thống nhất,
để những cách cài đặt khác nhau có thể tương tác với nhau.

**Unwrap**

Chức năng thường được thêm vào nhiều nhất
trong các gói này là một phương thức nào đó có thể được
gọi để loại bỏ ngữ cảnh khỏi một lỗi,
trả về lỗi bên trong.
Các gói dùng những tên và ý nghĩa khác nhau
cho thao tác này, và đôi khi nó loại bỏ một mức ngữ cảnh,
trong khi đôi khi nó loại bỏ nhiều mức nhất có thể.

Đối với Go 1.13, chúng tôi đã đưa ra quy ước rằng một
cách cài đặt lỗi bổ sung ngữ cảnh có thể tháo bỏ cho lỗi bên trong
nên cài đặt một phương thức `Unwrap` trả về lỗi bên trong,
gỡ bỏ lớp ngữ cảnh đó.
Nếu không có lỗi bên trong nào phù hợp để lộ ra cho caller,
thì hoặc lỗi đó không nên có phương thức `Unwrap`,
hoặc phương thức `Unwrap` nên trả về nil.

	// Phương thức tùy chọn trong Go 1.13 cho các cách cài đặt lỗi.

	interface {
		// Unwrap loại bỏ một lớp ngữ cảnh,
		// trả về lỗi bên trong nếu có, nếu không thì nil.
		Unwrap() error
	}

Cách để gọi phương thức tùy chọn này là dùng hàm trợ giúp `errors.Unwrap`,
hàm này xử lý các trường hợp như chính lỗi là nil hoặc không hề có phương thức `Unwrap`.

	package errors

	// Unwrap trả về kết quả của việc gọi
	// phương thức Unwrap trên err,
	// nếu kiểu của err có định nghĩa phương thức Unwrap.
	// Nếu không, Unwrap trả về nil.
	func Unwrap(err error) error

Chúng ta có thể dùng phương thức `Unwrap`
để viết một phiên bản đơn giản hơn, tổng quát hơn của `spuriousENOTAVAIL`.
Thay vì tìm các cách cài đặt lớp bọc lỗi cụ thể
như `net.OpError` hay `os.SyscallError`,
phiên bản tổng quát có thể lặp, gọi `Unwrap` để loại bỏ ngữ cảnh,
cho đến khi hoặc nó chạm tới `EADDRNOTAVAIL` hoặc không còn lỗi nào nữa:

	func spuriousENOTAVAIL(err error) bool {
		for err != nil {
			if err == syscall.EADDRNOTAVAIL {
				return true
			}
			err = errors.Unwrap(err)
		}
		return false
	}

Tuy vậy, vòng lặp này quá phổ biến, nên Go 1.13 định nghĩa thêm một hàm thứ hai, `errors.Is`,
hàm này liên tục unwrap một lỗi để tìm một target cụ thể.
Vì vậy chúng ta có thể thay toàn bộ vòng lặp bằng một lời gọi `errors.Is` duy nhất:

	func spuriousENOTAVAIL(err error) bool {
		return errors.Is(err, syscall.EADDRNOTAVAIL)
	}

Đến lúc này có lẽ chúng ta thậm chí sẽ không định nghĩa hàm đó nữa;
gọi thẳng `errors.Is` tại các call site cũng rõ ràng tương đương, và đơn giản hơn.

Go 1.13 cũng giới thiệu hàm `errors.As`
unwrap cho tới khi tìm thấy một kiểu cài đặt cụ thể.

Nếu bạn muốn viết mã làm việc với
các lỗi được bọc theo cách tùy ý,
`errors.Is` là phiên bản nhận biết lớp bọc
của phép kiểm tra bằng nhau giữa các lỗi:

	err == target

	    →

	errors.Is(err, target)

Và `errors.As` là phiên bản nhận biết lớp bọc
của một type assertion đối với lỗi:

	target, ok := err.(*Type)
	if ok {
	    ...
	}

	    →

	var target *Type
	if errors.As(err, &target) {
	   ...
	}

**Có nên Unwrap hay không?**

Việc có nên cho phép unwrap một lỗi hay không là một quyết định API,
giống như việc có nên export một trường của struct hay không là một quyết định API.
Đôi khi việc lộ chi tiết đó cho mã gọi là phù hợp,
và đôi khi thì không.
Khi phù hợp, hãy cài đặt Unwrap.
Khi không phù hợp, đừng cài đặt Unwrap.

Cho đến giờ, `fmt.Errorf` chưa từng lộ
một lỗi bên dưới được định dạng bằng `%v` cho caller kiểm tra.
Tức là, kết quả của `fmt.Errorf` trước đây không thể unwrap được.
Hãy xem ví dụ này:

	// errors.Unwrap(err2) == nil
	// err1 không khả dụng (giống các phiên bản Go trước đây)
	err2 := fmt.Errorf("connect: %v", err1)

Nếu `err2` được trả về cho
một caller, caller đó chưa từng có cách nào để mở `err2` và truy cập `err1`.
Chúng tôi đã giữ nguyên thuộc tính đó trong Go 1.13.

Cho những lúc bạn thực sự muốn cho phép unwrap kết quả của `fmt.Errorf`,
chúng tôi cũng thêm một verb định dạng mới là `%w`, định dạng giống `%v`,
yêu cầu một đối số là giá trị lỗi,
và khiến phương thức `Unwrap` của lỗi kết quả trả về chính đối số đó.
Trong ví dụ của chúng ta, giả sử ta thay `%v` bằng `%w`:

	// errors.Unwrap(err4) == err3
	// (%w là mới trong Go 1.13)
	err4 := fmt.Errorf("connect: %w", err3)

Bây giờ, nếu `err4` được trả về cho caller,
caller có thể dùng `Unwrap` để lấy lại `err3`.

Điều quan trọng cần lưu ý là những quy tắc tuyệt đối kiểu như
“luôn dùng `%v` (hoặc không bao giờ cài đặt `Unwrap`)” hay “luôn dùng `%w` (hoặc luôn luôn cài đặt `Unwrap`)”
đều sai giống như những quy tắc tuyệt đối kiểu “không bao giờ export trường struct” hay “luôn luôn export trường struct”.
Thay vào đó, quyết định đúng phụ thuộc vào
việc caller có nên có thể kiểm tra và phụ thuộc vào
thông tin bổ sung mà việc dùng `%w` hoặc cài đặt `Unwrap` làm lộ ra hay không.

Để minh họa cho điểm này,
mọi kiểu bọc lỗi trong standard library
vốn đã có trường `Err` được export
giờ đều có thêm phương thức `Unwrap` trả về trường đó,
nhưng các cách cài đặt có trường lỗi không export thì không,
và các chỗ dùng `fmt.Errorf` với `%v` hiện có vẫn tiếp tục dùng `%v`, không phải `%w`.

**In giá trị lỗi (Đã từ bỏ)**

Cùng với bản thảo thiết kế cho Unwrap,
chúng tôi cũng công bố một
[bản thảo thiết kế cho một phương thức tùy chọn để in lỗi phong phú hơn](/design/go2draft-error-printing),
bao gồm thông tin stack frame
và hỗ trợ lỗi được bản địa hóa, dịch thuật.

	// Phương thức tùy chọn cho các cách cài đặt lỗi
	type Formatter interface {
		Format(p Printer) (next error)
	}

	// Interface truyền vào Format
	type Printer interface {
		Print(args ...interface{})
		Printf(format string, args ...interface{})
		Detail() bool
	}

Phần này không đơn giản như `Unwrap`,
và tôi sẽ không đi vào chi tiết ở đây.
Khi chúng tôi thảo luận thiết kế này với cộng đồng Go trong mùa đông vừa rồi,
chúng tôi nhận ra rằng thiết kế đó chưa đủ đơn giản.
Nó quá khó để từng kiểu lỗi riêng lẻ cài đặt,
và nó cũng không giúp đủ nhiều cho các chương trình hiện có.
Xét trên tổng thể, nó không đơn giản hóa việc phát triển Go.

Kết quả của cuộc thảo luận với cộng đồng này là
chúng tôi đã từ bỏ thiết kế in ấn đó.

**Cú pháp lỗi**

Đó là về giá trị lỗi.
Giờ hãy xem nhanh về cú pháp lỗi,
một thử nghiệm khác cũng đã bị từ bỏ.

Đây là một đoạn mã từ
[`compress/lzw/writer.go`](https://go.googlesource.com/go/+/go1.12/src/compress/lzw/writer.go#209) trong standard library:

{{raw `
	// Ghi savedCode nếu hợp lệ.
	if e.savedCode != invalidCode {
		if err := e.write(e, e.savedCode); err != nil {
			return err
		}
		if err := e.incHi(); err != nil && err != errOutOfCodes {
			return err
		}
	}

	// Ghi mã eof.
	eof := uint32(1)<<e.litWidth + 1
	if err := e.write(e, eof); err != nil {
		return err
	}
`}}

Thoáng nhìn qua, khoảng một nửa đoạn mã này là kiểm tra lỗi.
Mắt tôi như bị lướt qua khi đọc nó.
Và chúng ta biết rằng mã vừa khó viết vừa khó đọc thì rất dễ bị đọc sai,
khiến nó trở thành nơi trú ngụ tốt cho những lỗi khó tìm.
Ví dụ, một trong ba lần kiểm tra lỗi này không giống hai lần còn lại,
điều rất dễ bỏ sót nếu chỉ lướt nhanh.
Nếu bạn đang debug đoạn mã này, bạn sẽ mất bao lâu để nhận ra điều đó?

Tại Gophercon năm ngoái, chúng tôi đã
[trình bày một bản thảo thiết kế](/design/go2draft-error-handling)
cho một cấu trúc điều khiển luồng mới được đánh dấu bằng từ khóa `check`.
`Check` tiêu thụ kết quả lỗi từ một lời gọi hàm hoặc biểu thức.
Nếu lỗi khác nil, `check` sẽ trả về lỗi đó.
Nếu không, `check` sẽ đánh giá thành các kết quả khác
từ lời gọi đó. Chúng ta có thể dùng `check` để đơn giản hóa mã lzw:

{{raw `
	// Ghi savedCode nếu hợp lệ.
	if e.savedCode != invalidCode {
		check e.write(e, e.savedCode)
		if err := e.incHi(); err != errOutOfCodes {
			check err
		}
	}

	// Ghi mã eof.
	eof := uint32(1)<<e.litWidth + 1
	check e.write(e, eof)
`}}

Phiên bản này của cùng đoạn mã dùng `check`,
giúp loại bỏ bốn dòng mã và
quan trọng hơn là làm nổi bật rằng
lời gọi tới `e.incHi` được phép trả về `errOutOfCodes`.

Có lẽ quan trọng nhất,
thiết kế này cũng cho phép định nghĩa các khối xử lý lỗi
sẽ được chạy khi các `check` phía sau thất bại.
Điều đó sẽ cho phép bạn chỉ viết một lần đoạn mã thêm ngữ cảnh dùng chung,
như trong đoạn sau:

{{raw `
	handle err {
		err = fmt.Errorf("closing writer: %w", err)
	}

	// Ghi savedCode nếu hợp lệ.
	if e.savedCode != invalidCode {
		check e.write(e, e.savedCode)
		if err := e.incHi(); err != errOutOfCodes {
			check err
		}
	}

	// Ghi mã eof.
	eof := uint32(1)<<e.litWidth + 1
	check e.write(e, eof)
`}}

Về bản chất, `check` là cách viết ngắn gọn cho câu lệnh `if`,
còn `handle` thì giống
[`defer`](/ref/spec#Defer_statements) nhưng chỉ dành cho đường trả về lỗi.
Khác với exception trong các ngôn ngữ khác,
thiết kế này vẫn giữ được thuộc tính quan trọng của Go rằng
mọi lời gọi có thể thất bại đều được đánh dấu rõ ràng trong mã,
giờ dùng từ khóa `check` thay cho `if err != nil`.

Vấn đề lớn của thiết kế này
là `handle` chồng lấn quá nhiều,
và theo những cách gây bối rối, với `defer`.

Vào tháng Năm, chúng tôi đăng
[một thiết kế mới với ba sự đơn giản hóa](/design/32437-try-builtin):
để tránh sự nhầm lẫn với `defer`, thiết kế bỏ `handle` để chỉ dùng `defer`;
để khớp với một ý tưởng tương tự trong Rust và Swift, thiết kế đổi tên `check` thành `try`;
và để cho phép thử nghiệm theo cách mà các parser hiện có như `gofmt` vẫn nhận ra được,
nó biến `check` (giờ là `try`) từ một từ khóa thành một hàm dựng sẵn.

Giờ cùng đoạn mã đó sẽ trông như sau:

{{raw `
	defer errd.Wrapf(&err, "closing writer")

	// Ghi savedCode nếu hợp lệ.
	if e.savedCode != invalidCode {
		try(e.write(e, e.savedCode))
		if err := e.incHi(); err != errOutOfCodes {
			try(err)
		}
	}

	// Ghi mã eof.
	eof := uint32(1)<<e.litWidth + 1
	try(e.write(e, eof))
`}}

Chúng tôi đã dành phần lớn tháng Sáu để thảo luận công khai đề xuất này trên GitHub.

Ý tưởng cốt lõi của `check` hay `try` là rút ngắn
lượng cú pháp bị lặp lại ở mỗi lần kiểm tra lỗi,
và đặc biệt là loại bỏ câu lệnh `return` khỏi tầm nhìn,
vẫn giữ việc kiểm tra lỗi là rõ ràng và làm nổi bật hơn những biến thể thú vị.
Tuy vậy, một điểm đáng chú ý được nêu ra trong phần phản hồi công khai
là nếu không có câu lệnh `if` và `return` rõ ràng,
sẽ không có chỗ để chèn một dòng in phục vụ debug,
sẽ không có chỗ để đặt breakpoint,
và cũng không có đoạn mã nào để hiển thị là chưa được chạy trong kết quả đo độ bao phủ mã.
Những lợi ích mà chúng tôi hướng tới
đã phải trả giá bằng việc làm những tình huống đó trở nên phức tạp hơn.
Xét trên tổng thể, từ điều này cũng như những cân nhắc khác,
hoàn toàn không rõ liệu kết quả cuối cùng có
đơn giản hóa việc phát triển Go hay không,
nên chúng tôi đã từ bỏ thử nghiệm này.

Đó là toàn bộ về xử lý lỗi,
vốn là một trong những trọng tâm chính của năm nay.

## Generics

Giờ đến một chủ đề ít gây tranh cãi hơn một chút: generics.

Chủ đề lớn thứ hai mà chúng tôi xác định cho Go 2 là
một cách nào đó để viết mã với
type parameters.
Điều này sẽ cho phép viết các cấu trúc dữ liệu generic
và cả những hàm generic
làm việc với bất kỳ kiểu slice nào,
hay bất kỳ kiểu channel nào,
hoặc bất kỳ kiểu map nào.
Ví dụ, đây là một bộ lọc channel generic:

{{raw `
	// Filter sao chép các giá trị từ c sang channel trả về,
	// chỉ chuyển tiếp những giá trị thỏa mãn f.
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

Chúng tôi đã nghĩ về generics từ khi công việc về Go bắt đầu,
và đã viết rồi bác bỏ thiết kế cụ thể đầu tiên vào năm 2010.
Chúng tôi viết rồi bác bỏ thêm ba thiết kế nữa trước cuối năm 2013.
Bốn thử nghiệm bị từ bỏ,
nhưng không phải những thử nghiệm thất bại,
Chúng tôi đã học được từ chúng,
giống như đã học từ `check` và `try`.
Mỗi lần như vậy, chúng tôi đều học được rằng con đường tới Go 2 không nằm đúng theo hướng đó,
và nhận ra những hướng khác có thể đáng để khám phá.
Nhưng đến năm 2013, chúng tôi quyết định rằng cần tập trung vào những mối quan tâm khác,
vì vậy chúng tôi gác toàn bộ chủ đề này lại trong vài năm.

Năm ngoái chúng tôi bắt đầu khám phá và thử nghiệm trở lại,
và đã trình bày một
[thiết kế mới](https://github.com/golang/proposal/blob/master/design/go2draft-contracts.md),
dựa trên ý tưởng về contract,
tại Gophercon mùa hè năm ngoái.
Chúng tôi tiếp tục thử nghiệm và đơn giản hóa,
đồng thời làm việc
với các chuyên gia lý thuyết ngôn ngữ lập trình
để hiểu rõ hơn về thiết kế này.

Nhìn chung, tôi hy vọng rằng chúng ta đang đi đúng hướng,
hướng tới một thiết kế sẽ đơn giản hóa việc phát triển Go.
Dù vậy, cũng có thể chúng ta sẽ nhận ra rằng thiết kế này cũng không hiệu quả.
Chúng ta có thể sẽ phải từ bỏ thử nghiệm này
và điều chỉnh con đường của mình dựa trên những gì đã học được.
Rồi chúng ta sẽ biết.

Tại Gophercon 2019, Ian Lance Taylor đã nói về
lý do chúng ta có thể muốn thêm generics vào Go
và xem trước ngắn gọn bản thảo thiết kế mới nhất.
Để biết chi tiết, hãy xem bài blog của ông ấy “[Why Generics?](/blog/why-generics)”

## Dependencies

Chủ đề lớn thứ ba mà chúng tôi xác định cho Go 2 là quản lý dependencies.

Năm 2010, chúng tôi công bố một công cụ tên là `goinstall`,
mà chúng tôi gọi là
“[một thử nghiệm về cài đặt package](https://groups.google.com/forum/#!msg/golang-nuts/8JFwR3ESjjI/cy7qZzN7Lw4J).”
Nó tải dependencies xuống và lưu chúng trong cây thư mục
của bộ phân phối Go của bạn, trong GOROOT.

Khi thử nghiệm với `goinstall`,
chúng tôi học được rằng bộ phân phối Go và các package đã cài đặt
nên được giữ tách biệt,
để có thể chuyển sang một bộ phân phối Go mới
mà không làm mất tất cả các package Go của bạn.
Vì vậy vào năm 2011 chúng tôi giới thiệu `GOPATH`,
một biến môi trường chỉ định
nơi tìm các package không có trong bộ phân phối Go chính.

Việc thêm GOPATH tạo ra thêm nhiều nơi chứa package Go
nhưng lại đơn giản hóa tổng thể việc phát triển Go,
bằng cách tách bộ phân phối Go của bạn khỏi các thư viện Go của bạn.

**Tương thích**

Thử nghiệm `goinstall` chủ ý bỏ qua
khái niệm phiên bản package một cách tường minh.
Thay vào đó, `goinstall` luôn tải bản sao mới nhất.
Chúng tôi làm vậy để có thể tập trung vào các
vấn đề thiết kế khác của việc cài đặt package.

`Goinstall` trở thành `go get` như một phần của Go 1.
Khi mọi người hỏi về phiên bản,
chúng tôi khuyến khích họ thử nghiệm bằng cách
tạo ra các công cụ bổ sung, và họ đã làm vậy.
Chúng tôi cũng khuyến khích các tác giả package
cung cấp cho người dùng của họ
khả năng tương thích ngược giống như
chúng tôi đã làm cho các thư viện Go 1.
Trích [Go FAQ](/doc/faq#get_version):

<div style="margin-left: 2em; font-style: italic;">

“Các package được dự định cho mục đích công khai nên cố gắng duy trì khả năng tương thích ngược khi chúng phát triển.

Nếu cần chức năng khác,
hãy thêm một tên mới thay vì thay đổi tên cũ.

Nếu cần một sự phá vỡ hoàn toàn,
hãy tạo một package mới với một import path mới.”

</div>

Quy ước này
đơn giản hóa trải nghiệm tổng thể của việc sử dụng một package
bằng cách giới hạn những gì tác giả có thể làm:
tránh các thay đổi phá vỡ API;
đặt tên mới cho chức năng mới;
và
đặt import path mới cho một thiết kế package hoàn toàn mới.

Dĩ nhiên, mọi người vẫn tiếp tục thử nghiệm.
Một trong những thử nghiệm thú vị nhất
được bắt đầu bởi Gustavo Niemeyer.
Ông đã tạo ra một trình chuyển hướng Git tên là
[`gopkg.in`](https://gopkg.in),
nơi cung cấp các import path khác nhau
cho các phiên bản API khác nhau,
để giúp các tác giả package
tuân theo quy ước
đặt import path mới
cho một thiết kế package mới.

Ví dụ,
mã nguồn Go trong kho GitHub
[go-yaml/yaml](https://github.com/go-yaml/yaml)
có các API khác nhau
ở các semantic version tag v1 và v2.
Máy chủ `gopkg.in` cung cấp chúng với
các import path khác nhau
[gopkg.in/yaml.v1](https://godoc.org/gopkg.in/yaml.v1)
và
[gopkg.in/yaml.v2](https://godoc.org/gopkg.in/yaml.v2).

Quy ước cung cấp khả năng tương thích ngược,
để một phiên bản mới hơn của package có thể được dùng
thay cho một phiên bản cũ hơn,
chính là điều giúp quy tắc rất đơn giản của `go get` là “luôn tải bản sao mới nhất”
vẫn hoạt động tốt cho đến tận ngày nay.

**Versioning và Vendoring**

Nhưng trong môi trường production, bạn cần chính xác hơn
về phiên bản dependency, để build có thể tái lập.

Nhiều người đã thử nghiệm xem điều đó nên trông như thế nào,
xây dựng những công cụ phục vụ nhu cầu của họ,
bao gồm `goven` của Keith Rarick (2012) và `godep` (2013),
`glide` của Matt Butcher (2014), và `gb` của Dave Cheney (2015).
Tất cả những công cụ này đều dùng mô hình sao chép các package dependency
vào chính kho mã nguồn của bạn.
Cơ chế chính xác được dùng
để làm cho các package đó sẵn sàng để import thì khác nhau,
nhưng tất cả đều phức tạp hơn mức lẽ ra phải có.

Sau một cuộc thảo luận trên toàn cộng đồng,
chúng tôi đã thông qua một đề xuất của Keith Rarick
để thêm hỗ trợ tường minh cho việc tham chiếu các dependency đã được sao chép
mà không cần mánh khóe GOPATH.
Đây là đơn giản hóa bằng cách định hình lại:
giống như với `addToList` và `append`,
các công cụ này vốn đã triển khai khái niệm đó,
nhưng lúng túng hơn mức cần thiết.
Việc thêm hỗ trợ tường minh cho thư mục vendor
đã làm cho các cách dùng này đơn giản hơn về tổng thể.

Việc đưa thư mục vendor vào lệnh `go`
dẫn đến nhiều thử nghiệm hơn với chính mô hình vendoring,
và chúng tôi nhận ra rằng mình đã đưa vào một vài vấn đề.
Nghiêm trọng nhất là chúng tôi đã đánh mất _tính duy nhất của package_.
Trước đây, trong bất kỳ lần build nào,
một import path
có thể xuất hiện trong rất nhiều package khác nhau,
và tất cả các import đều trỏ tới cùng một đích.
Giờ với vendoring, cùng một import path ở các
package khác nhau có thể trỏ tới các bản sao vendored khác nhau của package đó,
tất cả đều sẽ xuất hiện trong binary kết quả cuối cùng.

Thời điểm đó, chúng tôi chưa có tên cho thuộc tính này:
tính duy nhất của package.
Nó chỉ đơn giản là cách mô hình GOPATH hoạt động.
Chúng tôi không thật sự đánh giá hết tầm quan trọng của nó cho tới khi nó biến mất.

Ở đây có một sự song song với các đề xuất cú pháp lỗi `check` và `try`.
Trong trường hợp đó, chúng tôi đã dựa vào
cách mà câu lệnh `return` hiển thị hoạt động
theo những cách mà chúng tôi không nhận ra
cho đến khi cân nhắc loại bỏ nó.

Khi chúng tôi thêm hỗ trợ thư mục vendor,
có rất nhiều công cụ khác nhau để quản lý dependencies.
Chúng tôi nghĩ rằng một sự thống nhất rõ ràng
về định dạng của thư mục vendor
và metadata vendoring
sẽ cho phép các công cụ khác nhau tương tác với nhau,
giống như việc thống nhất về
cách các chương trình Go được lưu trong tệp văn bản
giúp có sự tương tác
giữa trình biên dịch Go, các trình soạn thảo văn bản,
và các công cụ như `goimports` và `gorename`.

Điều này hóa ra là lạc quan một cách ngây thơ.
Các công cụ vendoring đều khác nhau ở những sắc thái ngữ nghĩa tinh vi.
Khả năng tương tác sẽ đòi hỏi phải thay đổi tất cả chúng
để thống nhất về ngữ nghĩa,
có khả năng làm hỏng người dùng tương ứng của từng công cụ.
Sự hội tụ đã không xảy ra.

**Dep**

Tại Gophercon năm 2016, chúng tôi bắt đầu một nỗ lực
nhằm định nghĩa một công cụ duy nhất để quản lý dependencies.
Là một phần của nỗ lực đó, chúng tôi đã thực hiện khảo sát
với nhiều kiểu người dùng khác nhau
để hiểu họ cần gì
về mặt quản lý dependencies,
và một nhóm bắt đầu làm việc trên công cụ mới,
sau này trở thành `dep`.

`Dep` nhắm tới khả năng thay thế tất cả các
công cụ quản lý dependency hiện có.
Mục tiêu là đơn giản hóa bằng cách định hình lại
nhiều công cụ khác nhau hiện có thành một công cụ duy nhất.
Nó đã đạt được điều đó một phần.
`Dep` cũng khôi phục tính duy nhất của package cho người dùng của nó,
bằng cách chỉ có một thư mục vendor
ở đỉnh cây dự án.

Nhưng `dep` cũng đưa vào một vấn đề nghiêm trọng
mà phải mất một thời gian chúng tôi mới nhận ra đầy đủ.
Vấn đề là `dep` chấp nhận một lựa chọn thiết kế từ `glide`,
nhằm hỗ trợ và khuyến khích các thay đổi không tương thích với một package nhất định
mà không thay đổi import path.

Đây là một ví dụ.
Giả sử bạn đang xây dựng chương trình của riêng mình,
và bạn cần có một tệp cấu hình,
vì vậy bạn dùng phiên bản 2 của một package YAML phổ biến trong Go:

<div style="margin-left: 2em;">
{{image "experiment/yamldeps1.png" 214}}
</div>

Giờ giả sử chương trình của bạn
import Kubernetes client.
Hóa ra Kubernetes dùng YAML rất nhiều,
và nó dùng phiên bản 1 của cùng package phổ biến đó:

<div style="margin-left: 2em;">
{{image "experiment/yamldeps2.png" 557}}
</div>

Phiên bản 1 và phiên bản 2 có API không tương thích,
nhưng chúng cũng có import path khác nhau,
vì vậy không có sự mơ hồ nào về việc một import nhất định đang nói tới cái nào.
Kubernetes nhận phiên bản 1,
trình phân tích cấu hình của bạn nhận phiên bản 2,
và mọi thứ đều hoạt động.

`Dep` đã từ bỏ mô hình này.
Phiên bản 1 và phiên bản 2 của package yaml giờ sẽ
có cùng import path,
tạo ra xung đột.
Việc dùng cùng một import path cho hai phiên bản không tương thích,
kết hợp với tính duy nhất của package,
khiến không thể build chương trình này,
trong khi trước đó bạn vẫn build được:

<div style="margin-left: 2em;">
{{image "experiment/yamldeps3.png" 450}}
</div>

Phải mất một thời gian chúng tôi mới hiểu rõ vấn đề này,
bởi vì chúng tôi đã áp dụng quy ước
“API mới thì import path mới”
quá lâu tới mức xem đó là điều hiển nhiên.
Thử nghiệm dep đã giúp chúng tôi
đánh giá đúng hơn quy ước đó,
và chúng tôi đã đặt tên cho nó:
_quy tắc tương thích import_:

<div style="margin-left: 2em; font-style: italic;">

“Nếu một package cũ và một package mới có cùng import path,
thì package mới phải tương thích ngược với package cũ.”

</div>

**Go Modules**

Chúng tôi lấy những gì hoạt động tốt trong thử nghiệm dep
và những gì đã học được về điều gì không hoạt động tốt,
rồi thử nghiệm với một thiết kế mới tên là `vgo`.
Trong `vgo`, các package tuân theo quy tắc tương thích import,
để chúng ta có thể cung cấp tính duy nhất của package
nhưng vẫn không làm hỏng các bản build như bản ta vừa xem.
Điều đó cũng cho phép chúng tôi đơn giản hóa các phần khác của thiết kế.

Ngoài việc khôi phục quy tắc tương thích import,
một phần quan trọng khác của thiết kế `vgo`
là đặt tên cho khái niệm một nhóm package
và cho phép nhóm đó được tách biệt
khỏi ranh giới kho mã nguồn.
Tên của một nhóm package Go là module,
vì vậy giờ chúng tôi gọi hệ thống này là Go modules.

Go modules hiện đã được tích hợp với lệnh `go`,
giúp tránh nhu cầu phải sao chép các thư mục vendor qua lại.

**Thay thế GOPATH**

Cùng với Go modules là sự kết thúc của GOPATH như một
không gian tên toàn cục.
Gần như toàn bộ công việc khó khăn trong việc chuyển đổi cách dùng Go hiện có
và các công cụ sang modules là do thay đổi này gây ra,
do việc rời bỏ GOPATH.

Ý tưởng cốt lõi của GOPATH
là cây thư mục GOPATH
là nguồn chân lý toàn cục
về các phiên bản đang được dùng,
và các phiên bản đang được dùng sẽ không thay đổi
khi bạn di chuyển giữa các thư mục.
Nhưng chế độ GOPATH toàn cục lại mâu thuẫn trực tiếp
với yêu cầu production về build tái lập theo từng dự án,
mà bản thân điều đó lại đơn giản hóa trải nghiệm
phát triển và triển khai Go theo nhiều cách quan trọng.

Build tái lập theo từng dự án có nghĩa là
khi bạn đang làm việc trong checkout của dự án A,
bạn sẽ có cùng tập phiên bản dependency với các lập trình viên khác của dự án A
tại commit đó,
như được xác định bởi tệp `go.mod`.
Khi bạn chuyển sang làm việc trong checkout của dự án B,
giờ bạn sẽ có tập phiên bản dependency mà dự án đó chọn,
cũng là tập mà các lập trình viên khác của dự án B nhận được.
Nhưng chúng có khả năng khác với dự án A.
Việc tập phiên bản dependency
thay đổi khi bạn chuyển từ dự án A sang dự án B
là cần thiết để giữ cho quá trình phát triển của bạn đồng bộ
với quá trình phát triển của các lập trình viên khác ở A và ở B.
Không thể tiếp tục tồn tại một GOPATH toàn cục duy nhất nữa.

Phần lớn độ phức tạp của việc áp dụng modules
phát sinh trực tiếp từ việc mất đi GOPATH toàn cục duy nhất.
Mã nguồn của một package nằm ở đâu?
Trước đây, câu trả lời chỉ phụ thuộc vào biến môi trường GOPATH của bạn,
mà đa số mọi người hiếm khi thay đổi.
Giờ câu trả lời phụ thuộc vào dự án bạn đang làm việc,
thứ có thể thay đổi thường xuyên.
Mọi thứ đều cần được cập nhật cho quy ước mới này.

Hầu hết các công cụ phát triển đều dùng gói
[`go/build`](https://godoc.org/go/build) để tìm và nạp mã nguồn Go.
Chúng tôi vẫn giữ cho gói đó hoạt động,
nhưng API của nó không lường trước modules,
và các cách giải quyết tạm mà chúng tôi thêm vào để tránh thay đổi API
chậm hơn mức chúng tôi mong muốn.
Chúng tôi đã công bố một gói thay thế,
[`golang.org/x/tools/go/packages`](https://godoc.org/golang.org/x/tools/go/packages).
Các công cụ phát triển giờ nên dùng gói đó thay thế.
Nó hỗ trợ cả GOPATH lẫn Go modules,
và nhanh hơn, dễ dùng hơn.
Trong một hoặc hai bản phát hành tới chúng tôi có thể chuyển nó vào standard library,
nhưng hiện tại [`golang.org/x/tools/go/packages`](https://godoc.org/golang.org/x/tools/go/packages)
đã ổn định và sẵn sàng để dùng.

**Go Module Proxies**

Một trong những cách mà modules đơn giản hóa việc phát triển Go
là bằng cách tách khái niệm một nhóm package
khỏi kho mã nguồn quản lý phiên bản bên dưới
nơi chúng được lưu trữ.

Khi chúng tôi nói chuyện với người dùng Go về dependencies,
hầu như tất cả những người dùng Go trong công ty của họ
đều hỏi làm thế nào để định tuyến các lần tải package của `go get`
qua máy chủ riêng của họ,
để kiểm soát tốt hơn mã nào có thể được dùng.
Và ngay cả các lập trình viên mã nguồn mở cũng lo ngại
về việc dependencies biến mất
hoặc thay đổi bất ngờ,
làm hỏng các bản build của họ.
Trước modules, người dùng đã thử
những giải pháp phức tạp cho các vấn đề này,
bao gồm cả việc chặn các lệnh version control
mà lệnh `go` chạy.

Thiết kế Go modules giúp việc
đưa vào ý tưởng về một module proxy
có thể được hỏi xin một phiên bản module cụ thể trở nên dễ dàng.

Các công ty giờ có thể dễ dàng chạy module proxy của riêng mình,
với các quy tắc tùy chỉnh về những gì được phép
và nơi lưu các bản sao đã được cache.
Dự án mã nguồn mở [Athens](https://docs.gomods.io) đã xây dựng đúng một proxy như vậy,
và Aaron Schlesinger đã có một bài nói về nó tại Gophercon 2019.
(Chúng tôi sẽ thêm liên kết ở đây khi video có sẵn.)

Và đối với các lập trình viên cá nhân cũng như các nhóm mã nguồn mở,
nhóm Go tại Google đã [ra mắt một proxy](https://groups.google.com/forum/#!topic/golang-announce/0wo8cOhGuAI) đóng vai trò
là bản sao gương công khai của mọi package Go mã nguồn mở,
và Go 1.13 sẽ dùng proxy đó theo mặc định khi ở chế độ module.
Katie Hockman đã có một [bài nói về hệ thống này tại Gophercon 2019](https://youtu.be/KqTySYYhPUE).

**Trạng thái Go Modules**

Go 1.11 giới thiệu modules như một bản xem trước thử nghiệm, phải bật bằng opt-in.
Chúng tôi tiếp tục thử nghiệm và đơn giản hóa.
Go 1.12 phát hành các cải tiến,
và Go 1.13 sẽ phát hành thêm nhiều cải tiến nữa.

Modules giờ đã đi đến điểm
mà chúng tôi tin rằng chúng sẽ phục vụ được đa số người dùng,
nhưng chúng tôi chưa sẵn sàng tắt GOPATH ngay.
Chúng tôi sẽ tiếp tục thử nghiệm, đơn giản hóa và điều chỉnh.

Chúng tôi hoàn toàn nhận thức được rằng
cộng đồng người dùng Go
đã xây dựng gần một thập kỷ kinh nghiệm,
công cụ và quy trình làm việc xung quanh GOPATH,
và sẽ mất một thời gian để chuyển tất cả những thứ đó sang Go modules.

Nhưng một lần nữa,
chúng tôi nghĩ rằng giờ đây modules sẽ
hoạt động rất tốt với đa số người dùng,
và tôi khuyến khích bạn hãy xem thử
khi Go 1.13 được phát hành.

Như một điểm dữ liệu,
dự án Kubernetes có rất nhiều dependencies,
và họ đã chuyển sang dùng Go modules
để quản lý chúng.
Có lẽ bạn cũng làm được.
Và nếu bạn không thể,
hãy cho chúng tôi biết điều gì không hoạt động với bạn
hoặc điều gì quá phức tạp,
bằng cách [gửi báo cáo lỗi](/issue/new),
và chúng tôi sẽ thử nghiệm rồi đơn giản hóa.

## Công cụ

Xử lý lỗi, generics và quản lý dependencies
sẽ còn cần ít nhất thêm vài năm nữa,
và hiện tại chúng tôi sẽ tập trung vào chúng.
Xử lý lỗi gần như đã xong,
modules sẽ là phần tiếp theo sau đó,
và có thể là generics sau đó.

Nhưng giả sử ta nhìn xa hơn vài năm,
tới lúc chúng ta đã thử nghiệm và đơn giản hóa xong
và đã phát hành xử lý lỗi, modules và generics.
Sau đó thì sao?
Rất khó để dự đoán tương lai,
nhưng tôi nghĩ rằng khi cả ba thứ này đã được phát hành,
đó có thể đánh dấu sự khởi đầu của một giai đoạn yên ắng mới đối với các thay đổi lớn.
Trọng tâm của chúng tôi lúc đó có thể sẽ chuyển sang
đơn giản hóa việc phát triển Go bằng các công cụ tốt hơn.

Một phần công việc về công cụ đã bắt đầu,
vì vậy bài viết này kết thúc bằng cách nhìn vào điều đó.

Trong khi chúng tôi giúp cập nhật tất cả các
công cụ hiện có của cộng đồng Go để hiểu Go modules,
chúng tôi nhận ra rằng việc có quá nhiều công cụ trợ giúp phát triển
mỗi công cụ chỉ làm một việc nhỏ không phục vụ người dùng tốt.
Từng công cụ riêng lẻ quá khó kết hợp,
quá chậm khi gọi, và quá khác nhau để sử dụng.

Chúng tôi bắt đầu một nỗ lực nhằm hợp nhất những
công cụ trợ giúp phát triển thường cần nhất vào một công cụ duy nhất,
giờ được gọi là `gopls` (đọc là “go, please”).
`Gopls` nói
[Language Server Protocol, LSP](https://langserver.org/),
và làm việc với bất kỳ môi trường phát triển tích hợp nào
hoặc trình soạn thảo văn bản nào có hỗ trợ LSP,
nghĩa là về cơ bản là mọi thứ ở thời điểm này.

`Gopls` đánh dấu một sự mở rộng trọng tâm của dự án Go,
từ việc cung cấp các công cụ độc lập kiểu trình biên dịch, dòng lệnh
như go vet hay gorename
sang cả việc cung cấp một dịch vụ IDE hoàn chỉnh.
Rebecca Stambler đã có một bài nói chi tiết hơn về `gopls` và IDE tại Gophercon 2019.
(Chúng tôi sẽ thêm liên kết ở đây khi video có sẵn.)

Sau `gopls`, chúng tôi cũng có những ý tưởng để hồi sinh `go fix` theo một
cách mở rộng được và làm cho `go vet` còn hữu ích hơn nữa.

## Kết

<div style="margin-left: 2em;">
{{image "experiment/expsimp2.png" 326}}
</div>

Vậy đó là con đường đến Go 2.
Chúng ta sẽ thử nghiệm rồi đơn giản hóa.
Rồi thử nghiệm rồi đơn giản hóa.
Rồi phát hành.
Rồi thử nghiệm rồi đơn giản hóa.
Rồi làm lại tất cả.
Trông có thể giống hoặc thậm chí khiến bạn cảm thấy như con đường đó đi vòng quanh.
Nhưng mỗi lần chúng ta thử nghiệm và đơn giản hóa
chúng ta lại học thêm được một chút về hình hài của Go 2
và tiến thêm một bước gần hơn tới nó.
Ngay cả những thử nghiệm bị từ bỏ như `try`
hay bốn thiết kế generics đầu tiên của chúng tôi
hay `dep` cũng không phải là thời gian lãng phí.
Chúng giúp chúng tôi hiểu điều gì cần được
đơn giản hóa trước khi có thể phát hành,
và trong một số trường hợp chúng giúp chúng tôi hiểu rõ hơn
về điều mà chúng tôi từng xem là hiển nhiên.

Tại một thời điểm nào đó, chúng ta sẽ nhận ra rằng mình đã
thử nghiệm đủ, đơn giản hóa đủ,
và phát hành đủ,
và khi đó chúng ta sẽ có Go 2.

Cảm ơn tất cả các bạn trong cộng đồng Go
đã giúp chúng tôi thử nghiệm
và đơn giản hóa
và phát hành
và tìm đường đi trên con đường này.
