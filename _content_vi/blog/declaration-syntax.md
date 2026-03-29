---
title: Cú pháp khai báo của Go
date: 2010-07-07
by:
- Rob Pike
tags:
- c
- syntax
- ethos
summary: Vì sao cú pháp khai báo của Go không giống và cũng đơn giản hơn nhiều so với C.
---

## Giới thiệu

Những người mới đến với Go thường thắc mắc vì sao cú pháp khai báo lại khác với
truyền thống đã được thiết lập trong họ ngôn ngữ C.
Trong bài viết này, chúng ta sẽ so sánh hai cách tiếp cận và giải thích
vì sao khai báo trong Go lại có hình thức như vậy.

## Cú pháp của C

Trước tiên, hãy nói về cú pháp của C. C đã chọn một cách tiếp cận khác thường nhưng thông minh
cho cú pháp khai báo.
Thay vì mô tả kiểu bằng cú pháp đặc biệt,
người ta viết một biểu thức có liên quan đến mục đang được khai báo,
và phát biểu kiểu của biểu thức đó sẽ là gì. Vì vậy

	int x;

khai báo x là một `int`: biểu thức 'x' sẽ có kiểu `int`.
Nói chung, để tìm ra cách viết kiểu của một biến mới,
hãy viết một biểu thức có chứa biến đó sao cho biểu thức được đánh giá ra một kiểu cơ bản,
rồi đặt kiểu cơ bản ở bên trái và biểu thức ở bên phải.

Do đó, các khai báo

	int *p;
	int a[3];

cho biết rằng p là con trỏ tới `int` vì '\*p' có kiểu `int`,
và a là mảng các `int` vì a[3] (bỏ qua giá trị chỉ số cụ thể,
ở đây nó được mượn để biểu diễn kích thước mảng) có kiểu `int`.

Thế còn hàm thì sao? Ban đầu, khai báo hàm trong C viết kiểu
của các đối số bên ngoài dấu ngoặc, như sau:

	int main(argc, argv)
	    int argc;
	    char *argv[];
	{ /* ... */ }

Một lần nữa, ta thấy rằng `main` là một hàm vì biểu thức main(argc,
argv) trả về một `int`.
Trong ký hiệu hiện đại, ta sẽ viết

	int main(int argc, char *argv[]) { /* ... */ }

nhưng cấu trúc cơ bản vẫn vậy.

Đây là một ý tưởng cú pháp thông minh, hoạt động tốt với các kiểu đơn giản
nhưng có thể trở nên rối rất nhanh.
Ví dụ nổi tiếng là khai báo con trỏ hàm.
Làm theo quy tắc và bạn sẽ có:

	int (*fp)(int a, int b);

Ở đây, `fp` là con trỏ tới một hàm vì nếu bạn viết biểu thức (\*fp)(a,
b) thì bạn sẽ gọi một hàm trả về `int`.
Điều gì xảy ra nếu một trong các đối số của `fp` bản thân nó cũng là một hàm?

	int (*fp)(int (*ff)(int x, int y), int b)

Đến đây thì đã bắt đầu khó đọc.

Dĩ nhiên, khi khai báo hàm ta có thể bỏ tên của các tham số, vì vậy `main` có thể được khai báo là

	int main(int, char *[])

Hãy nhớ rằng `argv` được khai báo như sau,

	char *argv[]

vì vậy bạn bỏ tên ở giữa khai báo của nó để tạo ra kiểu.
Tuy nhiên, không hề hiển nhiên rằng bạn khai báo một thứ có kiểu `char *[]` bằng cách
đặt tên của nó vào giữa.

Và hãy xem điều gì xảy ra với khai báo của `fp` nếu bạn không đặt tên cho tham số:

	int (*fp)(int (*)(int, int), int)

Không chỉ khó nhận ra phải đặt tên vào đâu trong

	int (*)(int, int)

mà thậm chí còn không thật sự rõ đây là khai báo con trỏ hàm.
Và nếu kiểu trả về cũng là một con trỏ hàm thì sao?

	int (*(*fp)(int (*)(int, int), int))(int, int)

Thật khó để nhìn ra khai báo này đang nói về `fp`.

Bạn có thể dựng những ví dụ còn phức tạp hơn, nhưng chỉ bấy nhiêu cũng đủ minh họa
một số khó khăn mà cú pháp khai báo của C có thể tạo ra.

Tuy vậy, còn một điểm nữa cần nhắc tới.
Vì cú pháp kiểu và cú pháp khai báo là một,
nên việc phân tích các biểu thức có chứa kiểu ở giữa có thể trở nên khó khăn.
Đó là lý do, chẳng hạn, ép kiểu trong C luôn đặt kiểu trong ngoặc đơn, như

	(int)M_PI

## Cú pháp của Go

Những ngôn ngữ ngoài họ C thường dùng một cú pháp kiểu riêng trong khai báo.
Mặc dù đây là một điểm khác, tên thường được đặt trước,
thường theo sau bởi dấu hai chấm.
Do đó các ví dụ ở trên sẽ trở thành đại loại như thế này (trong một ngôn ngữ giả tưởng nhưng dễ hình dung)

	x: int
	p: pointer to int
	a: array[3] of int

Những khai báo này rất rõ ràng, dù dài dòng, bạn chỉ cần đọc từ trái sang phải.
Go lấy cảm hứng từ đó, nhưng vì lợi ích của sự ngắn gọn nó bỏ dấu hai chấm
và loại bỏ một số từ khóa:

	x int
	p *int
	a [3]int

Không có sự tương ứng trực tiếp nào giữa hình thức `[3]int` và cách
sử dụng `a` trong một biểu thức.
(Chúng ta sẽ quay lại với con trỏ ở phần tiếp theo.) Bạn có được sự rõ ràng
đổi lấy một cú pháp tách biệt.

Bây giờ hãy xét hàm. Hãy chép lại khai báo của `main` theo cách nó sẽ xuất hiện trong Go,
mặc dù hàm `main` thực tế trong Go không nhận đối số nào:

	func main(argc int, argv []string) int

Bề ngoài nó không khác C là mấy,
ngoài việc đổi từ mảng `char` sang string,
nhưng nó đọc rất tự nhiên từ trái sang phải:

hàm `main` nhận một `int` và một slice các string rồi trả về một `int`.

Bỏ tên tham số đi mà nó vẫn rõ ràng, chúng luôn đứng trước nên không gây nhầm lẫn.

	func main(int, []string) int

Một ưu điểm của kiểu trình bày từ trái sang phải này là nó hoạt động tốt thế nào
khi các kiểu trở nên phức tạp hơn.
Đây là khai báo của một biến hàm (tương tự con trỏ hàm trong C):

	f func(func(int,int) int, int) int

Hoặc nếu `f` trả về một hàm:

	f func(func(int,int) int, int) func(int, int) int

Nó vẫn đọc rõ ràng, từ trái sang phải,
và luôn luôn rõ tên nào đang được khai báo, tên luôn đứng trước.

Sự phân biệt giữa cú pháp kiểu và cú pháp biểu thức khiến việc viết và gọi closure trong Go trở nên dễ dàng:

	sum := func(a, b int) int { return a+b } (3, 4)

## Con trỏ

Con trỏ là ngoại lệ chứng minh quy tắc.
Chẳng hạn, hãy lưu ý rằng với mảng và slice,
cú pháp kiểu của Go đặt dấu ngoặc vuông ở bên trái kiểu nhưng cú pháp biểu thức
lại đặt chúng ở bên phải biểu thức:

	var a []int
	x = a[1]

Vì lý do quen thuộc, con trỏ trong Go sử dụng ký hiệu `*` từ C,
nhưng chúng tôi không thể tự thuyết phục mình thực hiện một sự đảo ngược tương tự cho kiểu con trỏ.
Do đó con trỏ hoạt động như sau

	var p *int
	x = *p

Chúng tôi không thể viết

	var p *int
	x = p*

vì dấu `*` hậu tố đó sẽ bị nhập nhằng với phép nhân. Chúng tôi có thể dùng ký hiệu `^` của Pascal, chẳng hạn:

	var p ^int
	x = p^

và có lẽ đáng ra nên làm như vậy (rồi chọn toán tử khác cho xor),
vì dấu sao tiền tố ở cả kiểu lẫn biểu thức làm mọi thứ phức tạp lên
theo nhiều cách.
Ví dụ, dù bạn có thể viết

	[]int("hi")

như một phép chuyển kiểu, bạn vẫn phải đặt kiểu trong ngoặc nếu nó bắt đầu bằng `*`:

	(*int)(nil)

Nếu chúng tôi chấp nhận từ bỏ `*` làm cú pháp con trỏ, những dấu ngoặc đó đã không cần thiết.

Vì vậy cú pháp con trỏ của Go gắn liền với dạng quen thuộc của C,
nhưng sự gắn bó đó cũng có nghĩa là chúng ta không thể hoàn toàn thoát khỏi việc dùng dấu ngoặc
để phân biệt kiểu và biểu thức trong văn phạm.

Tuy nhiên, nhìn chung, chúng tôi tin rằng cú pháp kiểu của Go dễ hiểu hơn cú pháp của C, đặc biệt khi mọi thứ trở nên phức tạp.

## Ghi chú

Khai báo của Go được đọc từ trái sang phải. Người ta cũng từng chỉ ra rằng khai báo của C được đọc theo hình xoắn ốc!
Xem [ The "Clockwise/Spiral Rule"](http://c-faq.com/decl/spiral.anderson.html) của David Anderson.
