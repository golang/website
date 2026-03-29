---
title: Phân tích type parameter
date: 2023-09-26
by:
- Ian Lance Taylor
summary: Vì sao chữ ký hàm trong các gói slices lại phức tạp đến vậy.
---

## Chữ ký hàm của gói slices

Hàm [`slices.Clone`](https://pkg.go.dev/slices#Clone) khá đơn giản:
nó tạo một bản sao của một slice thuộc bất kỳ kiểu nào.

```Go
func Clone[S ~[]E, E any](s S) S {
	return append(s[:0:0], s...)
}
```

Điều này hoạt động vì việc append vào một slice có zero capacity sẽ
cấp phát một backing array mới.
Phần thân hàm rốt cuộc còn ngắn hơn cả chữ ký hàm,
một phần vì phần thân ngắn, nhưng cũng vì
chữ ký thì dài.
Trong bài viết này, chúng ta sẽ giải thích vì sao chữ ký lại được viết theo cách đó.

## Clone đơn giản

Trước tiên, ta hãy viết một hàm generic `Clone` đơn giản.
Đây không phải phiên bản trong gói `slices`.
Ta muốn nhận vào một slice của bất kỳ kiểu phần tử nào, và trả về một slice mới.

```Go
func Clone1[E any](s []E) []E {
	// body omitted
}
```

Hàm generic `Clone1` có một type parameter duy nhất là `E`.
Nó nhận một đối số `s` là một slice có kiểu `E`, và
trả về một slice cùng kiểu.
Chữ ký này khá trực quan với bất kỳ ai đã quen với generic
trong Go.

Tuy nhiên, có một vấn đề.
Named slice type không phổ biến trong Go, nhưng vẫn có người dùng chúng.

```Go
// MySlice is a slice of strings with a special String method.
type MySlice []string

// String returns the printable version of a MySlice value.
func (s MySlice) String() string {
	return strings.Join(s, "+")
}
```

Giả sử ta muốn sao chép một `MySlice` rồi lấy phiên bản có thể in ra,
nhưng với các string được sắp xếp theo thứ tự.

```Go
func PrintSorted(ms MySlice) string {
	c := Clone1(ms)
	slices.Sort(c)
	return c.String() // FAILS TO COMPILE
}
```

Đáng tiếc là cách này không hoạt động.
Trình biên dịch báo lỗi:

```
c.String undefined (type []string has no field or method String)
```

Ta có thể thấy vấn đề nếu tự tay khởi tạo `Clone1` bằng cách
thay type parameter bằng type argument.

```Go
func InstantiatedClone1(s []string) []string
```

[Quy tắc gán giá trị của Go](/ref/spec#Assignability) cho phép
ta truyền một giá trị kiểu `MySlice` vào tham số kiểu
`[]string`, nên việc gọi `Clone1` là hợp lệ.
Nhưng `Clone1` sẽ trả về một giá trị kiểu `[]string`, chứ không phải giá trị
kiểu `MySlice`.
Kiểu `[]string` không có phương thức `String`, nên trình biên dịch
báo lỗi.

## Clone linh hoạt hơn

Để khắc phục vấn đề này, ta phải viết một phiên bản `Clone` mà
trả về cùng kiểu với đối số đầu vào.
Nếu làm được như vậy, thì khi gọi `Clone` với một giá trị kiểu
`MySlice`, kết quả trả về cũng sẽ có kiểu `MySlice`.

Ta biết nó phải trông đại loại như thế này.

```Go
func Clone2[S ?](s S) S // INVALID
```

Hàm `Clone2` này trả về một giá trị có cùng kiểu với
đối số của nó.

Ở đây tôi viết ràng buộc là `?`, nhưng đó chỉ là
một ký hiệu giữ chỗ.
Để làm điều này hoạt động, ta cần viết một ràng buộc cho phép ta viết
phần thân của hàm.
Với `Clone1`, ta chỉ cần dùng ràng buộc `any` cho kiểu phần tử.
Với `Clone2` thì không được: ta muốn bắt buộc `s` phải là
một kiểu slice.

Vì ta biết mình muốn một slice, ràng buộc của `S` phải là một
slice.
Ta không quan tâm kiểu phần tử của slice là gì, nên cứ gọi nó là
`E`, như ta đã làm với `Clone1`.

```Go
func Clone3[S []E](s S) S // INVALID
```

Cách này vẫn chưa hợp lệ, vì ta chưa khai báo `E`.
Type argument cho `E` có thể là bất kỳ kiểu nào, điều đó có nghĩa
nó cũng phải là một type parameter.
Vì nó có thể là bất kỳ kiểu nào, ràng buộc của nó là `any`.

```Go
func Clone4[S []E, E any](s S) S
```

Mọi thứ đã khá gần, và ít nhất nó sẽ biên dịch, nhưng ta vẫn
chưa đến đích.
Nếu biên dịch phiên bản này, ta sẽ gặp lỗi khi gọi `Clone4(ms)`.

```
MySlice does not satisfy []string (possibly missing ~ for []string in []string)
```

Trình biên dịch đang nói rằng ta không thể dùng type argument
`MySlice` cho type parameter `S`, vì `MySlice` không
thỏa mãn ràng buộc `[]E`.
Đó là vì `[]E` khi dùng làm ràng buộc chỉ cho phép một slice type
literal, như `[]string`.
Nó không cho phép một named type như `MySlice`.

## Ràng buộc theo underlying type

Như thông báo lỗi gợi ý, câu trả lời là thêm một dấu `~`.

```Go
func Clone5[S ~[]E, E any](s S) S
```

Nói lại cho rõ, viết type parameter và ràng buộc là `[S []E, E any]`
có nghĩa type argument cho `S` có thể là bất kỳ unnamed slice type nào,
nhưng không thể là một named type được định nghĩa từ một slice literal.
Viết `[S ~[]E, E any]`, có dấu `~`, có nghĩa type argument
cho `S` có thể là bất kỳ kiểu nào có underlying type là một kiểu slice.

Với bất kỳ named type nào `type T1 T2`, underlying type của `T1` là
underlying type của `T2`.
Underlying type của một predeclared type như `int` hoặc một type literal
như `[]string` đơn giản chính là bản thân kiểu đó.
Về chi tiết chính xác, [xem language
spec](/ref/spec#Underlying_types).
Trong ví dụ của chúng ta, underlying type của `MySlice` là `[]string`.

Vì underlying type của `MySlice` là một slice, ta có thể truyền
một đối số kiểu `MySlice` vào `Clone5`.
Như bạn có thể đã nhận ra, chữ ký của `Clone5` cũng chính là
chữ ký của `slices.Clone`.
Cuối cùng ta đã đi đến nơi mình muốn.

Trước khi tiếp tục, hãy bàn một chút vì sao cú pháp Go lại yêu cầu phải có `~`.
Thoạt nhìn, có vẻ như ta lúc nào cũng muốn cho phép truyền `MySlice`,
vậy tại sao không biến đó thành mặc định?
Hoặc, nếu cần hỗ trợ so khớp chính xác, tại sao không đảo ngược lại,
để ràng buộc `[]E` cho phép named type còn ràng buộc kiểu như
`=[]E` mới chỉ cho phép slice type literal?

Để giải thích điều này, trước hết hãy quan sát rằng một danh sách type parameter như
`[T ~MySlice]` là vô nghĩa.
Đó là vì `MySlice` không phải underlying type của bất kỳ kiểu nào khác.
Ví dụ, nếu ta có định nghĩa `type MySlice2 MySlice`,
thì underlying type của `MySlice2` là `[]string`, chứ không phải `MySlice`.
Vậy nên hoặc `[T ~MySlice]` sẽ không cho phép kiểu nào cả, hoặc nó sẽ
giống hệt `[T MySlice]` và chỉ khớp với `MySlice`.
Dù theo cách nào, `[T ~MySlice]` cũng không hữu ích.
Để tránh sự mơ hồ đó, ngôn ngữ cấm `[T ~MySlice]`, và
trình biên dịch sẽ báo một lỗi dạng

```
invalid use of ~ (underlying type of MySlice is []string)
```

Nếu Go không yêu cầu dấu ngã, khiến `[S []E]` khớp với mọi kiểu
có underlying type là `[]E`, thì ta sẽ phải định nghĩa
ý nghĩa của `[S MySlice]`.

Ta có thể cấm `[S MySlice]`, hoặc nói rằng `[S MySlice]`
chỉ khớp với `MySlice`, nhưng cách nào cũng gặp rắc rối với
predeclared type.
Một predeclared type, như `int`, có underlying type chính là nó.
Ta muốn cho phép mọi người viết ràng buộc chấp nhận
bất kỳ type argument nào có underlying type là `int`.
Trong ngôn ngữ hiện tại, họ làm điều đó bằng cách viết `[T ~int]`.
Nếu không yêu cầu dấu ngã, ta vẫn cần một cách để nói "mọi
kiểu có underlying type là `int`".
Cách tự nhiên để nói điều đó là `[T int]`.
Điều đó sẽ có nghĩa `[T MySlice]` và `[T int]` sẽ hành xử
khác nhau, dù trông rất giống nhau.

Ta có lẽ cũng có thể nói rằng `[S MySlice]` khớp với mọi kiểu có
underlying type là underlying type của `MySlice`, nhưng điều đó khiến
`[S MySlice]` trở nên không cần thiết và gây bối rối.

Chúng tôi cho rằng tốt hơn hết là yêu cầu `~` và nói thật rõ ràng khi nào
ta đang khớp theo underlying type thay vì chính kiểu đó.

## Suy luận kiểu

Giờ khi đã giải thích chữ ký của `slices.Clone`, ta hãy xem việc sử dụng `slices.Clone`
thực tế được đơn giản hóa như thế nào nhờ type inference.
Hãy nhớ rằng chữ ký của `Clone` là

```Go
func Clone[S ~[]E, E any](s S) S
```

Một lời gọi `slices.Clone` sẽ truyền một slice vào tham số `s`.
Simple type inference sẽ cho phép trình biên dịch suy ra rằng type
argument cho type parameter `S` là kiểu của slice được truyền vào
`Clone`.
Sau đó type inference đủ mạnh để nhận ra type argument
cho `E` chính là kiểu phần tử của type argument đã truyền cho `S`.

Điều này có nghĩa là ta có thể viết

```Go
	c := Clone(ms)
```

mà không cần phải viết

```Go
	c := Clone[MySlice, string](ms)
```

Nếu ta tham chiếu tới `Clone` mà không gọi nó, ta vẫn phải chỉ rõ một
type argument cho `S`, vì trình biên dịch không có gì để suy ra
nó.
May mắn là trong trường hợp đó, type inference vẫn có thể suy ra type
argument cho `E` từ đối số của `S`, và ta không cần
phải chỉ rõ riêng.

Nói cách khác, ta có thể viết

```Go
	myClone := Clone[MySlice]
```

mà không cần phải viết

```Go
	myClone := Clone[MySlice, string]
```

## Phân tích type parameter

Kỹ thuật tổng quát mà ta dùng ở đây, trong đó ta định nghĩa một type
parameter `S` bằng cách sử dụng một type parameter khác là `E`, là một cách
để phân tích kiểu trong chữ ký hàm generic.
Bằng cách phân tích một kiểu, ta có thể đặt tên, và áp ràng buộc lên, mọi khía cạnh
của kiểu đó.

Ví dụ, đây là chữ ký của `maps.Clone`.

```Go
func Clone[M ~map[K]V, K comparable, V any](m M) M
```

Giống hệt `slices.Clone`, ta dùng một type parameter cho kiểu của
tham số `m`, rồi phân tích kiểu đó bằng hai type parameter khác là
`K` và `V`.

Trong `maps.Clone`, ta ràng buộc `K` phải là comparable, như yêu cầu đối với
kiểu khóa của map.
Ta có thể ràng buộc các kiểu thành phần theo bất kỳ cách nào mình muốn.

```Go
func WithStrings[S ~[]E, E interface { String() string }](s S) (S, []string)
```

Điều này nói rằng đối số của `WithStrings` phải là một kiểu slice mà
kiểu phần tử của nó có phương thức `String`.

Vì mọi kiểu trong Go đều có thể được xây dựng từ các kiểu thành phần, ta luôn
có thể dùng type parameter để phân tích các kiểu đó và áp ràng buộc theo ý muốn.
