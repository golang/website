---
title: Sinh mã
date: 2014-12-22
by:
- Rob Pike
tags:
- programming
- technical
summary: Cách dùng go generate.
template: true
---


Một thuộc tính của tính toán phổ dụng, tức tính đầy đủ Turing, là một chương trình máy tính có thể viết ra một chương trình máy tính khác.
Đây là một ý tưởng mạnh mẽ nhưng không phải lúc nào cũng được đánh giá đúng mức, dù nó xảy ra rất thường xuyên.
Chẳng hạn, nó là một phần lớn trong định nghĩa của một trình biên dịch.
Đó cũng là cách lệnh `go test` hoạt động: nó quét các package cần kiểm thử,
viết ra một chương trình Go chứa bộ khung kiểm thử được tùy biến cho package,
rồi biên dịch và chạy nó.
Máy tính hiện đại nhanh đến mức chuỗi thao tác nghe có vẻ tốn kém này có thể hoàn thành chỉ trong một phần của giây.

Còn rất nhiều ví dụ khác về các chương trình viết ra chương trình.
Ví dụ [Yacc](https://godoc.org/golang.org/x/tools/cmd/goyacc), đọc mô tả của một văn phạm rồi viết ra một chương trình để phân tích văn phạm đó.
"Trình biên dịch" protocol buffer đọc một mô tả giao diện và sinh ra các định nghĩa cấu trúc,
phương thức và mã hỗ trợ khác.
Các công cụ cấu hình đủ loại cũng hoạt động như vậy, xem xét metadata hoặc môi trường
rồi sinh ra phần khung được tùy biến theo trạng thái cục bộ.

Vì thế, các chương trình viết ra chương trình là những thành phần quan trọng trong kỹ nghệ phần mềm,
nhưng những chương trình như Yacc, vốn sinh ra mã nguồn, cần được tích hợp vào quy trình build
để đầu ra của chúng có thể được biên dịch.
Khi dùng một công cụ build bên ngoài như Make, điều này thường dễ thực hiện.
Nhưng trong Go, nơi công cụ go lấy mọi thông tin build cần thiết từ chính mã nguồn Go, lại có một vấn đề.
Đơn giản là không có cơ chế nào để chạy Yacc chỉ bằng công cụ go.

Cho đến bây giờ thì có.

[Bản phát hành Go mới nhất](/blog/go1.4), 1.4,
bao gồm một lệnh mới giúp việc chạy các công cụ như thế trở nên dễ dàng hơn.
Nó được gọi là `go generate`, và hoạt động bằng cách quét tìm các chú thích đặc biệt trong mã nguồn Go
xác định các lệnh tổng quát cần chạy.
Điều quan trọng cần hiểu là `go generate` không phải là một phần của `go build`.
Nó không hề có phân tích phụ thuộc và phải được chạy tường minh trước khi chạy `go build`.
Nó được thiết kế cho tác giả của package Go dùng, chứ không phải người dùng package đó.

Lệnh `go generate` rất dễ dùng.
Để khởi động, đây là cách dùng nó để sinh một văn phạm Yacc.

Trước hết, cài công cụ Yacc của Go:

	go get golang.org/x/tools/cmd/goyacc

Giả sử bạn có một tệp đầu vào Yacc tên là `gopher.y` định nghĩa văn phạm cho ngôn ngữ mới của bạn.
Để sinh tệp mã nguồn Go hiện thực văn phạm đó,
bình thường bạn sẽ gọi lệnh như sau:

	goyacc -o gopher.go -p parser gopher.y

Tùy chọn `-o` đặt tên tệp đầu ra, còn `-p` chỉ định tên package.

Để `go generate` điều khiển quy trình này, trong một trong các tệp `.go` thông thường (không phải tệp sinh ra)
trong cùng thư mục, thêm chú thích này vào bất cứ đâu trong tệp:

	//go:generate goyacc -o gopher.go -p parser gopher.y

Đoạn văn bản này chính là lệnh ở trên với tiền tố là một chú thích đặc biệt mà `go generate` nhận ra.
Chú thích phải bắt đầu ngay đầu dòng và không có khoảng trắng giữa `//` và `go:generate`.
Sau mốc đó, phần còn lại của dòng chỉ định một lệnh để `go generate` chạy.

Giờ hãy chạy nó. Chuyển vào thư mục mã nguồn và chạy `go generate`, rồi `go build`, v.v.:

	$ cd $GOPATH/myrepo/gopher
	$ go generate
	$ go build
	$ go test

Vậy là xong.
Giả sử không có lỗi nào, lệnh `go generate` sẽ gọi `yacc` để tạo `gopher.go`,
lúc đó thư mục sẽ chứa đủ bộ tệp mã nguồn Go, nên ta có thể build, test và làm việc như bình thường.
Mỗi lần `gopher.y` được sửa, chỉ cần chạy lại `go generate` để sinh lại parser.

Để biết thêm chi tiết về cách `go generate` hoạt động, bao gồm tùy chọn, biến môi trường,
v.v., hãy xem [tài liệu thiết kế](/s/go1.4-generate).

Go generate không làm gì mà Make hay cơ chế build khác không làm được,
nhưng nó đi kèm với công cụ `go`, không cần cài đặt thêm, và phù hợp tự nhiên với hệ sinh thái Go.
Chỉ cần nhớ rằng nó dành cho tác giả package, không phải người dùng package,
nếu chỉ vì lý do là chương trình mà nó gọi có thể không có trên máy đích.
Ngoài ra, nếu package chứa nó được thiết kế để có thể được import bằng `go get`,
thì sau khi tệp được sinh ra (và đã được kiểm thử!), nó phải được commit vào
kho mã nguồn để sẵn sàng cho người dùng package.

Giờ ta đã có nó, hãy dùng nó cho một việc mới.
Là một ví dụ rất khác về cách `go generate` có thể giúp ích, có một chương trình mới trong
kho `golang.org/x/tools` tên là `stringer`.
Nó tự động viết các phương thức string cho các tập hằng số nguyên.
Nó không phải là một phần của bản phân phối phát hành, nhưng rất dễ cài:

	$ go get golang.org/x/tools/cmd/stringer

Đây là ví dụ trong tài liệu của
[`stringer`](https://godoc.org/golang.org/x/tools/cmd/stringer).
Giả sử ta có đoạn mã chứa một tập hằng số nguyên xác định các loại thuốc viên khác nhau:

	package painkiller

	type Pill int

	const (
		Placebo Pill = iota
		Aspirin
		Ibuprofen
		Paracetamol
		Acetaminophen = Paracetamol
	)

Để gỡ lỗi, ta muốn các hằng số này có thể tự in đẹp, nghĩa là ta cần một phương thức có chữ ký:

	func (p Pill) String() string

Viết bằng tay cũng dễ, có thể như thế này:

	func (p Pill) String() string {
		switch p {
		case Placebo:
			return "Placebo"
		case Aspirin:
			return "Aspirin"
		case Ibuprofen:
			return "Ibuprofen"
		case Paracetamol: // == Acetaminophen
			return "Paracetamol"
		}
		return fmt.Sprintf("Pill(%d)", p)
	}

Dĩ nhiên còn nhiều cách khác để viết hàm này.
Ta có thể dùng một slice string được lập chỉ mục bởi `Pill`, hoặc map, hoặc kỹ thuật khác.
Dù làm thế nào đi nữa, ta vẫn phải bảo trì nó nếu thay đổi tập thuốc viên, và phải chắc rằng nó đúng.
(Hai tên cho paracetamol khiến việc này khó hơn bình thường.)
Ngoài ra, câu hỏi chọn cách tiếp cận nào còn phụ thuộc vào kiểu và giá trị:
có dấu hay không dấu, dày hay thưa, bắt đầu từ 0 hay không, v.v.

Chương trình `stringer` xử lý tất cả những chi tiết này.
Dù có thể chạy độc lập, nó được thiết kế để được điều khiển bởi `go generate`.
Để dùng nó, thêm một chú thích generate vào mã nguồn, có thể đặt gần định nghĩa kiểu:

	//go:generate stringer -type=Pill

Quy tắc này chỉ rõ rằng `go generate` nên chạy công cụ `stringer` để sinh một phương thức `String` cho kiểu `Pill`.
Đầu ra sẽ tự động được ghi vào `pill_string.go` (mặc định này có thể ghi đè bằng cờ `-output`).

Hãy chạy nó:

{{raw `
	$ go generate
	$ cat pill_string.go
	// Code generated by stringer -type Pill pill.go; DO NOT EDIT.

	package painkiller

	import "fmt"

	const _Pill_name = "PlaceboAspirinIbuprofenParacetamol"

	var _Pill_index = [...]uint8{0, 7, 14, 23, 34}

	func (i Pill) String() string {
		if i < 0 || i+1 >= Pill(len(_Pill_index)) {
			return fmt.Sprintf("Pill(%d)", i)
		}
		return _Pill_name[_Pill_index[i]:_Pill_index[i+1]]
	}
	$
`}}

Mỗi lần ta thay đổi định nghĩa của `Pill` hoặc các hằng số, tất cả những gì cần làm là chạy

	$ go generate

để cập nhật phương thức `String`.
Và dĩ nhiên nếu trong cùng package ta thiết lập nhiều kiểu theo cách này,
thì chỉ một lệnh duy nhất đó sẽ cập nhật toàn bộ phương thức `String` của chúng.

Không có gì phải nghi ngờ rằng phương thức sinh ra trông xấu.
Nhưng điều đó không sao, bởi con người không cần làm việc trên nó; mã sinh tự động thường xấu.
Nó đang nỗ lực để đạt hiệu năng cao.
Tất cả tên bị ép lại thành một chuỗi duy nhất,
giúp tiết kiệm bộ nhớ (chỉ một string header cho tất cả tên, kể cả khi có rất nhiều tên).
Sau đó một mảng `_Pill_index` ánh xạ từ giá trị sang tên bằng một kỹ thuật đơn giản, hiệu quả.
Cũng lưu ý rằng `_Pill_index` là một mảng (không phải slice; bớt thêm một header) gồm `uint8`,
kiểu số nguyên nhỏ nhất đủ để bao phủ không gian giá trị.
Nếu có nhiều giá trị hơn, hoặc có giá trị âm,
kiểu sinh ra cho `_Pill_index` có thể đổi thành `uint16` hoặc `int8`: bất cứ thứ gì phù hợp nhất.

Cách tiếp cận mà các phương thức do `stringer` in ra sử dụng sẽ thay đổi tùy theo đặc tính của tập hằng số.
Ví dụ, nếu các hằng số thưa, nó có thể dùng map.
Đây là một ví dụ nhỏ dựa trên một tập hằng số biểu diễn lũy thừa của hai:

	const _Power_name = "p0p1p2p3p4p5..."

	var _Power_map = map[Power]string{
		1:    _Power_name[0:2],
		2:    _Power_name[2:4],
		4:    _Power_name[4:6],
		8:    _Power_name[6:8],
		16:   _Power_name[8:10],
		32:   _Power_name[10:12],
		...,
	}

	func (i Power) String() string {
		if str, ok := _Power_map[i]; ok {
			return str
		}
		return fmt.Sprintf("Power(%d)", i)
	}

Tóm lại, việc sinh phương thức tự động cho phép ta làm tốt hơn điều mà ta kỳ vọng ở con người.

Đã có nhiều cách dùng khác của `go generate` được cài sẵn trong cây mã Go.
Ví dụ bao gồm sinh bảng Unicode trong package `unicode`,
tạo các phương thức hiệu quả để mã hóa và giải mã mảng trong `encoding/gob`,
sinh dữ liệu múi giờ trong package `time`, v.v.

Hãy dùng `go generate` một cách sáng tạo.
Nó ở đó để khuyến khích thử nghiệm.

Và ngay cả khi bạn không làm vậy, hãy dùng công cụ `stringer` mới để viết các phương thức `String` cho các hằng số nguyên của bạn.
Hãy để máy làm việc đó.
