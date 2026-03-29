---
title: Chương trình Go đầu tiên
date: 2013-07-18
by:
- Andrew Gerrand
tags:
- history
summary: Rob Pike đã đào lại chương trình Go đầu tiên từng được viết.
template: true
---


Brad Fitzpatrick và tôi (Andrew Gerrand) gần đây đã bắt đầu tái cấu trúc
[godoc](/cmd/godoc/), và tôi chợt nhận ra rằng đó là một trong
những chương trình Go cổ nhất.
Robert Griesemer đã bắt đầu viết nó từ đầu năm 2009,
và đến hôm nay chúng tôi vẫn còn dùng nó.

Khi tôi [tweet](https://twitter.com/enneff/status/357403054632484865) về điều này,
Dave Cheney đã trả lời bằng một [câu hỏi thú vị](https://twitter.com/davecheney/status/357406479415914497):
đâu là chương trình Go cổ nhất? Rob Pike đã lục lại thư cũ của mình và tìm thấy nó
trong một tin nhắn cũ gửi cho Robert và Ken Thompson.

Những gì dưới đây là chương trình Go đầu tiên. Nó được Rob viết vào tháng 2 năm 2008,
khi nhóm chỉ gồm Rob, Robert và Ken. Họ đã có một danh sách tính năng khá đầy đủ
(được nhắc tới trong [bài viết blog này](https://commandcenter.blogspot.com.au/2012/06/less-is-exponentially-more.html))
và một bản đặc tả ngôn ngữ sơ bộ. Ken vừa hoàn thành phiên bản đầu tiên hoạt động được của
trình biên dịch Go (nó chưa sinh mã máy bản địa mà chuyển mã Go sang C để
tạo mẫu nhanh) và đã đến lúc thử viết một chương trình bằng nó.

Rob gửi thư cho "nhóm Go":

	From: Rob 'Commander' Pike
	Date: Wed, Feb 6, 2008 at 3:42 PM
	To: Ken Thompson, Robert Griesemer
	Subject: slist

	it works now.

	roro=% a.out
	(defn foo (add 12 34))
	return: icounter = 4440
	roro=%

	here's the code.
	some ugly hackery to get around the lack of strings.

(Dòng `icounter` trong đầu ra của chương trình là số câu lệnh đã được thực thi,
được in ra để phục vụ việc gỡ lỗi.)

{{code "first-go-program/slist.go"}}

Chương trình này phân tích và in ra một
[S-expression](https://en.wikipedia.org/wiki/S-expression).
Nó không nhận dữ liệu đầu vào từ người dùng và không có import nào, chỉ dựa vào
khả năng `print` tích hợp sẵn để xuất kết quả.
Nó được viết đúng vào ngày đầu tiên có một
[trình biên dịch hoạt động được nhưng còn rất sơ khai](/change/8b8615138da3).
Phần lớn ngôn ngữ khi ấy chưa được hiện thực và một phần thậm chí còn chưa được đặc tả.

Dù vậy, sắc thái cơ bản của ngôn ngữ ngày nay vẫn có thể nhận ra trong chương trình này.
Khai báo kiểu và biến, điều khiển luồng, và câu lệnh package
vẫn chưa thay đổi nhiều.

Nhưng cũng có rất nhiều khác biệt và thiếu vắng.
Đáng chú ý nhất là sự thiếu vắng concurrency và interface, cả hai đều
được xem là thiết yếu ngay từ ngày đầu nhưng khi đó vẫn chưa được thiết kế.

`func` khi đó là `function`, và chữ ký của nó chỉ rõ giá trị trả về
_trước_ đối số, ngăn cách bằng {{raw "`<-`"}}, thứ mà hiện nay chúng ta dùng như toán tử
gửi/nhận trên channel. Ví dụ, hàm `WhiteSpace` nhận số nguyên
`c` và trả về một giá trị boolean.

{{raw `
	function WhiteSpace(bool <- c int)
`}}

Mũi tên này là một giải pháp tạm thời cho đến khi xuất hiện cú pháp tốt hơn để khai báo
nhiều giá trị trả về.

Method tách biệt khỏi function và có từ khóa riêng.

{{raw `
	method (this *Slist) Car(*Slist <-) {
		return this.list.car;
	}
`}}

Và method được khai báo trước trong định nghĩa struct, dù điều đó sớm thay đổi.

{{raw `
	type Slist struct {
		...
		Car method(*Slist <-);
	}
`}}

Không có string, dù chúng đã có trong đặc tả.
Để lách điều đó, Rob phải dựng chuỗi đầu vào như một mảng `uint8` với
một cách viết khá vụng về. (Mảng khi đó còn sơ khai và slice vẫn chưa được thiết kế,
chứ chưa nói đến hiện thực, dù đã có khái niệm chưa được hiện thực của
"open array".)

	input[i] = '('; i = i + 1;
	input[i] = 'd'; i = i + 1;
	input[i] = 'e'; i = i + 1;
	input[i] = 'f'; i = i + 1;
	input[i] = 'n'; i = i + 1;
	input[i] = ' '; i = i + 1;
	...

Cả `panic` và `print` đều là từ khóa tích hợp sẵn, không phải hàm được khai báo trước.

	print "parse error: expected ", c, "\n";
	panic "parse";

Và còn nhiều khác biệt nhỏ khác nữa; hãy thử xem bạn có thể nhận ra thêm điều gì.

Chưa đầy hai năm sau khi chương trình này được viết, Go đã được phát hành như một
dự án mã nguồn mở. Nhìn lại, thật ấn tượng khi thấy ngôn ngữ đã phát triển và trưởng thành đến mức nào.
(Thứ cuối cùng thay đổi giữa proto-Go này và Go mà ta biết ngày nay là việc loại bỏ dấu chấm phẩy.)

Nhưng điều còn ấn tượng hơn là chúng ta đã học được nhiều đến mức nào về việc _viết_ mã Go.
Ví dụ, Rob gọi receiver của method là `this`, nhưng giờ chúng ta dùng những tên ngắn hơn, phù hợp ngữ cảnh.
Còn hàng trăm ví dụ quan trọng hơn nữa
và cho đến tận hôm nay chúng ta vẫn tiếp tục khám phá ra những cách tốt hơn để viết mã Go.
(Hãy xem mẹo thông minh của [gói glog](https://github.com/golang/glog) để
[xử lý mức độ verbose](https://github.com/golang/glog/blob/c6f9652c7179652e2fd8ed7002330db089f4c9db/glog.go#L893).)

Tôi tự hỏi ngày mai chúng ta sẽ còn học được điều gì nữa.
