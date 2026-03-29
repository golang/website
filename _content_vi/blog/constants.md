---
title: Hằng số
date: 2014-08-25
by:
- Rob Pike
tags:
- constants
summary: Giới thiệu về hằng số trong Go.
template: true
---

## Giới thiệu

Go là ngôn ngữ kiểu tĩnh không cho phép các phép toán trộn lẫn các kiểu số.
Bạn không thể cộng `float64` với `int`, hay thậm chí `int32` với `int`.
Thế nhưng vẫn hợp lệ khi viết `1e6*time.Second` hay `math.Exp(1)` hoặc thậm chí {{raw "`1<<('\t'+2.0)`"}}.
Trong Go, hằng số, không giống biến, hành xử gần như các con số thông thường.
Bài viết này giải thích vì sao lại như vậy và điều đó có nghĩa gì.

## Bối cảnh: C

Trong những ngày đầu nghĩ về Go, chúng tôi đã nói nhiều về các vấn đề
do cách C và những ngôn ngữ hậu duệ của nó cho phép bạn trộn và ghép các kiểu số.
Nhiều lỗi bí ẩn, sự cố và vấn đề về tính khả chuyển được gây ra bởi các biểu thức
kết hợp số nguyên có kích thước và “signedness” khác nhau.
Mặc dù với một lập trình viên C lão luyện, kết quả của phép tính như

	unsigned int u = 1e9;
	long signed int i = -1;
	... i + u ...

có thể là quen thuộc, nhưng nó không hề hiển nhiên _a priori_.
Kết quả lớn đến đâu?
Giá trị của nó là gì?
Nó là signed hay unsigned?

Những lỗi khó chịu ẩn nấp ở đây.

C có một bộ quy tắc được gọi là “the usual arithmetic conversions” và
độ tinh vi của chúng thể hiện ở việc chúng đã thay đổi qua nhiều năm (đồng thời
đưa thêm nhiều lỗi mới, một cách hồi tố).

Khi thiết kế Go, chúng tôi quyết định tránh bãi mìn này bằng cách bắt buộc rằng _không được_ trộn các kiểu số.
Nếu bạn muốn cộng `i` và `u`, bạn phải tường minh về điều mình muốn kết quả trở thành.
Với

	var u uint
	var i int

bạn có thể viết `uint(i)+u` hoặc `i+int(u)`,
với cả ý nghĩa lẫn kiểu của phép cộng được diễn đạt rõ ràng,
nhưng không giống C, bạn không thể viết `i+u`.
Bạn thậm chí không thể trộn `int` và `int32`, ngay cả khi `int` là kiểu 32-bit.

Sự nghiêm ngặt này loại bỏ một nguyên nhân phổ biến gây lỗi và các thất bại khác.
Đây là một thuộc tính sống còn của Go.
Nhưng nó có cái giá: đôi khi nó buộc lập trình viên phải trang trí mã
của mình bằng các phép chuyển kiểu số vụng về để biểu đạt ý nghĩa rõ ràng.

Vậy còn hằng số thì sao?
Với các khai báo ở trên, điều gì sẽ làm cho việc viết `i` `=` `0` hoặc `u` `=` `0` trở nên hợp lệ?
_Kiểu_ của `0` là gì?
Sẽ là vô lý nếu bắt buộc hằng số phải có chuyển kiểu trong các ngữ cảnh đơn giản như `i` `=` `int(0)`.

Chúng tôi sớm nhận ra rằng câu trả lời nằm ở chỗ làm cho hằng số số học hoạt động khác
với cách chúng hành xử trong các ngôn ngữ giống C khác.
Sau rất nhiều suy nghĩ và thử nghiệm, chúng tôi nghĩ ra một thiết kế mà
chúng tôi tin rằng gần như luôn đem lại cảm giác đúng đắn,
giải phóng lập trình viên khỏi việc phải chuyển kiểu hằng số mọi lúc mà vẫn
có thể viết như `math.Sqrt(2)` mà không bị compiler trách móc.

Tóm lại, hằng số trong Go “cứ thế hoạt động”, ít nhất là phần lớn thời gian.
Hãy xem điều đó xảy ra như thế nào.

## Thuật ngữ

Trước hết, một định nghĩa nhanh.
Trong Go, `const` là từ khóa giới thiệu một tên cho một giá trị vô hướng như `2` hoặc `3.14159` hoặc `"scrumptious"`.
Những giá trị như thế, có tên hoặc không, được gọi là _hằng số_ trong Go.
Hằng số cũng có thể được tạo ra bởi các biểu thức xây dựng từ hằng số,
chẳng hạn `2+3` hoặc `2+3i` hoặc `math.Pi/2` hoặc `("go"+"pher")`.

Một số ngôn ngữ không có hằng số, và một số khác có định nghĩa khái quát hơn
về hằng số hoặc cách dùng từ `const`.
Ví dụ, trong C và C++, `const` là một định tính kiểu có thể mã hóa
những thuộc tính phức tạp hơn của những giá trị phức tạp hơn.

Nhưng trong Go, hằng số chỉ là một giá trị đơn giản, không thay đổi, và từ đây trở đi chúng ta chỉ nói về Go.

## Hằng chuỗi

Có nhiều loại hằng số số, số nguyên,
số thực, rune, signed, unsigned, imaginary,
complex, nên hãy bắt đầu với một dạng hằng số đơn giản hơn: chuỗi.
Hằng chuỗi rất dễ hiểu và tạo ra một không gian nhỏ hơn để
khảo sát các vấn đề kiểu của hằng số trong Go.

Một hằng chuỗi đặt một vài ký tự văn bản giữa cặp dấu nháy kép.
(Go cũng có literal chuỗi thô, đặt trong dấu backquote <code>``</code>,
nhưng cho mục đích thảo luận này, chúng có cùng tính chất.)
Đây là một hằng chuỗi:

	"Hello, 世界"

(Để xem chi tiết hơn nhiều về biểu diễn và diễn giải chuỗi,
xem [bài viết blog này](/blog/strings).)

Hằng chuỗi này có kiểu gì?
Câu trả lời hiển nhiên là `string`, nhưng điều đó _sai_.

Đây là một _hằng chuỗi không kiểu_, nghĩa là một giá trị văn bản hằng
chưa có kiểu cố định.
Đúng, nó là một chuỗi, nhưng nó chưa phải là một giá trị Go có kiểu `string`.
Nó vẫn là hằng chuỗi không kiểu ngay cả khi được đặt tên:

	const hello = "Hello, 世界"

Sau khai báo này, `hello` cũng là một hằng chuỗi không kiểu.
Một hằng số không kiểu chỉ đơn giản là một giá trị, một giá trị chưa được gán kiểu xác định
vốn sẽ buộc nó phải tuân theo các quy tắc nghiêm ngặt ngăn việc kết hợp các giá trị khác kiểu.

Chính khái niệm hằng số _không kiểu_ này làm cho chúng ta có thể dùng hằng số trong Go với độ tự do rất lớn.

Vậy thì một hằng chuỗi _có kiểu_ là gì?
Đó là một hằng đã được gán kiểu, như sau:

	const typedHello string = "Hello, 世界"

Hãy chú ý rằng khai báo của `typedHello` có kiểu `string` tường minh trước dấu bằng.
Điều này nghĩa là `typedHello` có kiểu Go là `string`, và không thể được gán cho một biến Go có kiểu khác.
Nói cách khác, đoạn mã này hoạt động:

{{play "constants/string1.go" `/START/` `/STOP/`}}

nhưng đoạn mã này thì không:

{{play "constants/string2.go" `/START/` `/STOP/`}}

Biến `m` có kiểu `MyString` và không thể được gán một giá trị có kiểu khác.
Nó chỉ có thể được gán những giá trị kiểu `MyString`, như sau:

{{play "constants/string3.go" `/START/` `/STOP/`}}

hoặc bằng cách ép chuyển kiểu, như sau:

{{play "constants/string4.go" `/START/` `/STOP/`}}

Quay lại hằng chuỗi _không kiểu_ của chúng ta,
nó có một thuộc tính hữu ích là, vì nó không có kiểu,
việc gán nó cho một biến có kiểu sẽ không gây lỗi kiểu.
Tức là, ta có thể viết

	m = "Hello, 世界"

hoặc

	m = hello

bởi vì, không giống các hằng có kiểu `typedHello` và `myStringHello`,
các hằng không kiểu `"Hello, 世界"` và `hello` _không có kiểu_.
Việc gán chúng cho một biến thuộc bất kỳ kiểu nào tương thích với chuỗi đều hoạt động không lỗi.

Những hằng chuỗi không kiểu này dĩ nhiên là chuỗi,
nên chúng chỉ có thể được dùng ở nơi chấp nhận chuỗi,
nhưng chúng không có _kiểu_ `string`.

## Kiểu mặc định

Là một lập trình viên Go, hẳn bạn đã thấy rất nhiều khai báo như

	str := "Hello, 世界"

và đến lúc này bạn có thể đang hỏi, “nếu hằng số không có kiểu, làm sao `str` có được kiểu trong khai báo biến này?”
Câu trả lời là một hằng số không kiểu có một kiểu mặc định,
một kiểu ngầm mà nó truyền cho một giá trị nếu cần kiểu mà không có kiểu nào được cung cấp.
Đối với hằng chuỗi không kiểu, kiểu mặc định dĩ nhiên là `string`, nên

	str := "Hello, 世界"

hoặc

	var str = "Hello, 世界"

nghĩa chính xác giống như

	var str string = "Hello, 世界"

Một cách để nghĩ về hằng số không kiểu là chúng sống trong một kiểu
không gian giá trị lý tưởng nào đó,
một không gian ít ràng buộc hơn hệ thống kiểu đầy đủ của Go.
Nhưng để làm bất cứ điều gì với chúng, ta phải gán chúng cho biến,
và khi điều đó xảy ra thì _biến_ (chứ không phải bản thân hằng số) cần một kiểu,
và hằng số có thể cho biến biết nó nên có kiểu gì.
Trong ví dụ này, `str` trở thành giá trị kiểu `string` vì hằng chuỗi không kiểu
đã cung cấp cho khai báo kiểu mặc định của nó là `string`.

Trong kiểu khai báo như vậy, một biến được khai báo cùng một kiểu và một giá trị khởi tạo.
Tuy nhiên, đôi khi khi dùng hằng số, nơi đến của giá trị lại không rõ ràng như thế.
Ví dụ hãy xét câu lệnh này:

{{play "constants/default1.go" `/START/` `/STOP/`}}

Chữ ký của `fmt.Printf` là

	func Printf(format string, a ...interface{}) (n int, err error)

tức là các đối số của nó (sau chuỗi định dạng) là các giá trị interface.
Điều xảy ra khi `fmt.Printf` được gọi với một hằng số không kiểu là một giá trị interface được tạo ra
để truyền làm đối số, và kiểu cụ thể được lưu cho đối số đó là kiểu mặc định của hằng số.
Quá trình này tương tự với điều ta đã thấy trước đó khi khai báo một giá trị khởi tạo bằng hằng chuỗi không kiểu.

Bạn có thể thấy kết quả trong ví dụ này, dùng định dạng `%v` để in
giá trị và `%T` để in kiểu của giá trị được truyền vào `fmt.Printf`:

{{play "constants/default2.go" `/START/` `/STOP/`}}

Nếu hằng số có kiểu, thì kiểu đó sẽ đi vào interface, như ví dụ này cho thấy:

{{play "constants/default3.go" `/START/` `/STOP/`}}

(Để biết thêm thông tin về cách giá trị interface hoạt động,
xem những phần đầu của [bài viết blog này](/blog/laws-of-reflection).)

Tóm lại, một hằng có kiểu tuân theo mọi quy tắc của các giá trị có kiểu trong Go.
Mặt khác, một hằng không kiểu không mang theo kiểu Go theo cùng cách đó
và có thể được trộn ghép một cách tự do hơn.
Tuy vậy, nó vẫn có một kiểu mặc định sẽ bộc lộ ra khi, và chỉ khi, không có thông tin kiểu nào khác sẵn có.

## Kiểu mặc định được quyết định bởi cú pháp

Kiểu mặc định của một hằng không kiểu được quyết định bởi cú pháp của nó.
Đối với hằng chuỗi, kiểu ngầm khả dĩ duy nhất là `string`.
Đối với [hằng số học](/ref/spec#Numeric_types), kiểu ngầm có sự đa dạng hơn.
Hằng số nguyên mặc định thành `int`, hằng số dấu phẩy động thành `float64`,
hằng rune thành `rune` (một bí danh của `int32`),
và hằng ảo thành `complex128`.
Đây là câu lệnh in chuẩn của chúng ta được dùng lặp lại để chỉ ra các kiểu mặc định đang hoạt động:

{{play "constants/syntax.go" `/START/` `/STOP/`}}

(Bài tập: Giải thích kết quả cho `'x'`.)

## Boolean

Mọi điều chúng ta nói về hằng chuỗi không kiểu đều có thể nói về hằng boolean không kiểu.
Các giá trị `true` và `false` là những hằng boolean không kiểu có thể được gán cho bất kỳ biến boolean nào,
nhưng một khi đã được cho kiểu, các biến boolean không thể bị trộn lẫn:

{{play "constants/bool.go" `/START/` `/STOP/`}}

Hãy chạy ví dụ và xem điều gì xảy ra, sau đó comment dòng “Bad” rồi chạy lại.
Mẫu hình ở đây hoàn toàn đi theo hằng chuỗi.

## Số thực

Hằng dấu phẩy động nhìn chung giống hằng boolean ở hầu hết các khía cạnh.
Ví dụ chuẩn của chúng ta hoạt động đúng như mong đợi khi chuyển sang số thực:

{{play "constants/float1.go" `/START/` `/STOP/`}}

Một điểm nhăn là Go có _hai_ kiểu dấu phẩy động: `float32` và `float64`.
Kiểu mặc định cho hằng dấu phẩy động là `float64`, dù một hằng dấu phẩy động không kiểu
vẫn có thể được gán cho giá trị `float32` mà không vấn đề gì:

{{play "constants/float2.go" `/START/` `/STOP/`}}

Giá trị dấu phẩy động là nơi tốt để giới thiệu khái niệm tràn, hay phạm vi giá trị.

Hằng số học sống trong một không gian số có độ chính xác tùy ý; chúng chỉ là các con số thông thường.
Nhưng khi được gán cho một biến, giá trị đó phải chứa vừa trong đích đến.
Ta có thể khai báo một hằng với giá trị rất lớn:

{{code "constants/float3.go" `/Huge/`}}

đó vẫn chỉ là một con số, nhưng ta không thể gán nó hay thậm chí in nó. Câu lệnh sau thậm chí không biên dịch:

{{play "constants/float3.go" `/Println/`}}

Lỗi là “constant 1.00000e+1000 overflows float64”, điều đó là đúng.
Nhưng `Huge` vẫn có thể hữu ích: ta có thể dùng nó trong các biểu thức với những hằng số khác
và dùng giá trị của các biểu thức đó nếu kết quả
có thể được biểu diễn trong phạm vi của `float64`.
Câu lệnh

{{play "constants/float4.go" `/Println/`}}

in ra `10`, đúng như mong đợi.

Theo cách liên quan, hằng dấu phẩy động có thể có độ chính xác rất cao,
để các phép toán liên quan tới chúng chính xác hơn.
Các hằng được định nghĩa trong package [math](/pkg/math) có nhiều chữ số hơn rất nhiều so với
những gì một `float64` có thể giữ. Đây là định nghĩa của `math.Pi`:

	Pi	= 3.14159265358979323846264338327950288419716939937510582097494459

Khi giá trị đó được gán cho một biến,
một phần độ chính xác sẽ mất đi;
việc gán sẽ tạo ra giá trị `float64` (hoặc `float32`)
gần nhất với giá trị có độ chính xác cao. Đoạn mã này

{{play "constants/float5.go" `/START/` `/STOP/`}}

in ra `3.141592653589793`.

Việc có nhiều chữ số khả dụng như vậy có nghĩa là các phép tính như `Pi/2` hoặc
những phép đánh giá phức tạp hơn có thể mang theo nhiều độ chính xác hơn
cho tới khi kết quả được gán, giúp việc viết các phép tính có hằng số dễ hơn mà không mất độ chính xác.
Nó cũng có nghĩa là không bao giờ có trường hợp mà các góc cạnh của số thực như vô cực,
soft underflow và `NaN` phát sinh trong biểu thức hằng.
(Chia cho hằng số không là lỗi biên dịch,
và khi mọi thứ đều là số thì không có chuyện “không phải là số”.)

## Số phức

Hằng số phức hành xử khá giống hằng dấu phẩy động.
Đây là một phiên bản của điệp khúc quen thuộc bây giờ của chúng ta, chuyển sang số phức:

{{play "constants/complex1.go" `/START/` `/STOP/`}}

Kiểu mặc định của một số phức là `complex128`, phiên bản có độ chính xác cao hơn được ghép từ hai giá trị `float64`.

Để rõ ràng trong ví dụ, chúng tôi viết đầy đủ biểu thức `(0.0+1.0i)`,
nhưng giá trị này có thể được rút gọn thành `0.0+1.0i`,
`1.0i` hoặc thậm chí chỉ `1i`.

Hãy chơi một mẹo nhỏ.
Ta biết rằng trong Go, một hằng số học đơn giản chỉ là một con số.
Điều gì sẽ xảy ra nếu con số đó là một số phức không có phần ảo, tức là một số thực?
Đây là một ví dụ:

{{code "constants/complex2.go" `/const Two/`}}

Đó là một hằng số phức không kiểu.
Dù nó không có phần ảo, _cú pháp_ của biểu thức vẫn định nghĩa nó có kiểu mặc định là `complex128`.
Vì thế, nếu ta dùng nó để khai báo biến, kiểu mặc định sẽ là `complex128`. Đoạn mã

{{play "constants/complex2.go" `/START/` `/STOP/`}}

in ra `complex128:` `(2+0i)`.
Nhưng về mặt số học, `Two` có thể được lưu trong một số dấu phẩy động vô hướng,
một `float64` hoặc `float32`, mà không mất thông tin.
Do đó ta có thể gán `Two` cho một `float64`, trong khởi tạo hoặc trong phép gán, mà không vấn đề gì:

{{play "constants/complex3.go" `/START/` `/STOP/`}}

Đầu ra là `2` `and` `2`.
Dù `Two` là một hằng số phức, nó vẫn có thể được gán cho các biến dấu phẩy động vô hướng.
Khả năng để một hằng “băng qua” các kiểu như vậy sẽ rất hữu ích.

## Số nguyên

Cuối cùng ta đến với số nguyên.
Chúng có nhiều bộ phận chuyển động hơn, [nhiều kích cỡ, signed hay unsigned, và hơn thế nữa](/ref/spec#Numeric_types), nhưng
chúng vẫn chơi theo cùng quy tắc.
Lần cuối cùng, đây là ví dụ quen thuộc của chúng ta, lần này chỉ dùng `int`:

{{play "constants/int1.go" `/START/` `/STOP/`}}

Cùng ví dụ này có thể dựng cho bất kỳ kiểu số nguyên nào, cụ thể là:

	int int8 int16 int32 int64
	uint uint8 uint16 uint32 uint64
	uintptr

(cộng với các bí danh `byte` cho `uint8` và `rune` cho `int32`).
Rất nhiều kiểu, nhưng mẫu hình trong cách hằng số hoạt động giờ đây hẳn đã đủ quen
để bạn có thể thấy mọi thứ sẽ diễn ra như thế nào.

Như đã nhắc ở trên, số nguyên có vài dạng và mỗi dạng
có kiểu mặc định riêng:
`int` cho các hằng đơn giản như `123` hoặc `0xFF` hoặc `-14`,
và `rune` cho ký tự có dấu nháy như `'a'`, `'世'` hoặc `'\r'`.

Không có dạng hằng nào có kiểu mặc định là một kiểu số nguyên unsigned.
Tuy nhiên, tính linh hoạt của hằng số không kiểu có nghĩa là ta vẫn có thể khởi tạo
các biến số nguyên unsigned bằng các hằng đơn giản miễn là ta làm rõ kiểu.
Điều này tương tự với việc ta có thể khởi tạo một `float64` bằng số phức có phần ảo bằng không.
Dưới đây là một vài cách khác nhau để khởi tạo một `uint`;
tất cả đều tương đương, nhưng tất cả đều phải nhắc đến kiểu một cách tường minh để kết quả là unsigned.

	var u uint = 17
	var u = uint(17)
	u := uint(17)

Tương tự với vấn đề phạm vi ở phần nói về giá trị dấu phẩy động,
không phải giá trị số nguyên nào cũng chứa vừa trong mọi kiểu số nguyên.
Có hai vấn đề có thể phát sinh: giá trị có thể quá lớn,
hoặc nó có thể là một giá trị âm đang được gán cho một kiểu unsigned.
Ví dụ, `int8` có phạm vi từ -128 đến 127,
nên các hằng nằm ngoài phạm vi đó không bao giờ có thể được gán cho biến kiểu `int8`:

{{play "constants/int2.go" `/var/`}}

Tương tự, `uint8`, còn được biết đến là `byte`,
có phạm vi từ 0 đến 255, nên một hằng lớn hoặc âm không thể được gán cho `uint8`:

{{play "constants/int3.go" `/var/`}}

Việc kiểm tra kiểu này có thể bắt được những sai sót như ví dụ sau:

{{play "constants/int4.go" `/START/` `/STOP/`}}

Nếu compiler phàn nàn về cách bạn dùng một hằng số, rất có thể đó là một lỗi thật như thế này.

## Bài tập: unsigned int lớn nhất

Đây là một bài tập nhỏ giàu thông tin.
Làm sao biểu diễn một hằng số biểu thị giá trị lớn nhất chứa vừa trong một `uint`?
Nếu ta đang nói về `uint32` thay vì `uint`, ta có thể viết

{{raw `
	const MaxUint32 = 1<<32 - 1
`}}

nhưng chúng ta muốn `uint`, không phải `uint32`.
Các kiểu `int` và `uint` có cùng số bit không xác định, hoặc là 32 hoặc là 64.
Vì số bit phụ thuộc kiến trúc, ta không thể chỉ viết ra một giá trị duy nhất.

Những người hâm mộ [số học bù hai](http://en.wikipedia.org/wiki/Two's_complement),
thứ mà các số nguyên của Go được định nghĩa để sử dụng, biết rằng biểu diễn của `-1` có tất cả các bit đặt là 1,
vì thế mẫu bit của `-1` về nội bộ giống với
số nguyên unsigned lớn nhất.
Do đó ta có thể nghĩ rằng mình có thể viết

{{play "constants/exercise1.go" `/const/`}}

nhưng điều đó là bất hợp pháp vì -1 không thể được biểu diễn bởi một biến unsigned;
`-1` không nằm trong phạm vi của giá trị unsigned.
Một phép chuyển kiểu cũng không giúp được gì, vì cùng lý do đó:

{{play "constants/exercise2.go" `/const/`}}

Dù tại thời gian chạy một giá trị -1 có thể được chuyển thành số nguyên unsigned, các quy tắc
cho [phép chuyển kiểu](/ref/spec#Conversions) của hằng số cấm loại ép buộc này ở thời gian biên dịch.
Nói cách khác, đoạn mã này hoạt động:

{{play "constants/exercise3.go" `/START/` `/STOP/`}}

nhưng chỉ vì `v` là một biến; nếu ta biến `v` thành hằng,
ngay cả hằng không kiểu, ta lại quay về vùng cấm:

{{play "constants/exercise4.go" `/START/` `/STOP/`}}

Ta quay lại cách tiếp cận trước đó, nhưng thay vì `-1` thì thử `^0`,
phép phủ định bit của một số lượng bit 0 tùy ý.
Nhưng điều đó cũng thất bại, vì lý do tương tự:
Trong không gian giá trị số,
`^0` biểu diễn vô hạn số 1, nên ta mất thông tin nếu gán nó cho bất kỳ số nguyên kích thước cố định nào:

{{play "constants/exercise5.go" `/const/`}}

Vậy thì làm sao biểu diễn số nguyên unsigned lớn nhất như một hằng số?

Chìa khóa là ràng buộc phép toán vào số lượng bit của một `uint` và tránh
những giá trị, như số âm, không thể biểu diễn trong một `uint`.
Giá trị `uint` đơn giản nhất là hằng có kiểu `uint(0)`.
Nếu `uint` có 32 hay 64 bit, `uint(0)` có tương ứng 32 hay 64 bit 0.
Nếu ta đảo từng bit trong số đó, ta sẽ được đúng số lượng bit 1 cần có, tức là giá trị `uint` lớn nhất.

Vì vậy ta không lật bit của hằng không kiểu `0`, mà lật bit của hằng có kiểu `uint(0)`.
Đây, vậy thì, là hằng số của chúng ta:

{{play "constants/exercise6.go" `/START/` `/STOP/`}}

Dù số bit cần để biểu diễn một `uint` trong môi trường thực thi hiện tại là bao nhiêu
(trên [playground](/blog/playground), đó là 32),
hằng số này vẫn biểu diễn chính xác giá trị lớn nhất mà một biến kiểu `uint` có thể nắm giữ.

Nếu bạn hiểu được phân tích dẫn chúng ta tới kết quả này,
thì bạn đã hiểu mọi điểm quan trọng về hằng số trong Go.

## Các con số

Khái niệm hằng số không kiểu trong Go có nghĩa là tất cả các hằng số số,
dù là số nguyên, dấu phẩy động, số phức,
hay thậm chí giá trị ký tự,
đều sống trong một loại không gian thống nhất nào đó.
Chỉ khi ta đưa chúng vào thế giới tính toán của biến,
phép gán và phép toán thì kiểu thực sự mới quan trọng.
Nhưng miễn là ta ở trong thế giới của các hằng số số, ta có thể trộn ghép giá trị tùy ý.
Tất cả các hằng sau đây đều có giá trị số bằng 1:

	1
	1.000
	1e3-99.0*10-9
	'\x01'
	'\u0001'
	'b' - 'a'
	1.0+3i-3.0i

Do đó, mặc dù chúng có các kiểu mặc định ngầm khác nhau,
khi được viết dưới dạng hằng số không kiểu, chúng có thể được gán cho biến thuộc bất kỳ kiểu số nào:

{{play "constants/numbers1.go" `/START/` `/STOP/`}}

Đầu ra từ đoạn mã này là: `1 1 1 1 1 (1+0i) 1`.

Bạn thậm chí có thể làm những thứ hơi điên rồ như

{{play "constants/numbers2.go" `/START/` `/STOP/`}}

và nhận được 145.5, điều vô nghĩa ngoại trừ việc chứng minh một điểm.

Nhưng điểm chính của những quy tắc này là tính linh hoạt.
Tính linh hoạt đó có nghĩa là, mặc dù trong Go là bất hợp pháp nếu trong
cùng một biểu thức trộn biến dấu phẩy động và biến số nguyên,
hoặc thậm chí biến `int` và `int32`, thì vẫn hoàn toàn ổn khi viết

	sqrt2 := math.Sqrt(2)

hoặc

	const millisecond = time.Second/1e3

hoặc

	bigBufferWithHeader := make([]byte, 512+1e6)

và để kết quả mang đúng ý nghĩa bạn mong đợi.

Bởi vì trong Go, hằng số số học hoạt động đúng như bạn kỳ vọng: như những con số.
