---
title: Mọi kiểu comparable của bạn
date: 2023-02-17
by:
- Robert Griesemer
summary: type parameters, type sets, comparable types, constraint satisfaction
---

Ngày 1 tháng 2, chúng tôi đã phát hành phiên bản Go mới nhất, 1.20,
trong đó có một vài thay đổi ngôn ngữ.
Ở đây chúng ta sẽ bàn về một trong những thay đổi đó: ràng buộc kiểu `comparable`
được khai báo sẵn giờ đây được thỏa mãn bởi mọi [kiểu comparable](/ref/spec#Comparison_operators).
Điều đáng ngạc nhiên là trước Go 1.20, một số kiểu comparable lại không thỏa mãn `comparable`!

Nếu bạn thấy bối rối, bạn đã đến đúng chỗ.
Hãy xét khai báo map hợp lệ sau

```Go
var lookupTable map[any]string
```

trong đó kiểu khóa của map là `any` (vốn là một
[kiểu comparable](/ref/spec#Comparison_operators)).
Điều này hoạt động hoàn toàn ổn trong Go.
Mặt khác, trước Go 1.20, kiểu map tổng quát có vẻ tương đương

```Go
type genericLookupTable[K comparable, V any] map[K]V
```

có thể được dùng giống như kiểu map thông thường, nhưng lại sinh lỗi biên dịch khi
`any` được dùng làm kiểu khóa:

```Go
var lookupTable genericLookupTable[any, string] // ERROR: any does not implement comparable (Go 1.18 and Go 1.19)
```

Bắt đầu từ Go 1.20, đoạn mã này biên dịch bình thường.

Hành vi của `comparable` trước Go 1.20 đặc biệt khó chịu vì nó
ngăn chúng tôi viết loại thư viện tổng quát mà chúng tôi hy vọng có thể viết
với generics ngay từ đầu.
Hàm [`maps.Clone`](/issue/57436) được đề xuất

```Go
func Clone[M ~map[K]V, K comparable, V any](m M) M { … }
```

có thể được viết ra nhưng lại không thể được dùng cho một map như `lookupTable` vì cùng lý do
mà `genericLookupTable` không thể dùng với `any` làm khóa.

Trong bài viết này, chúng tôi hy vọng làm sáng tỏ đôi chút cơ chế ngôn ngữ phía sau tất cả điều này.
Để làm vậy, trước tiên ta cần một chút thông tin nền.

## Tham số kiểu và ràng buộc

Go 1.18 giới thiệu generics và cùng với đó,
[_tham số kiểu_](/ref/spec#Type_parameter_declarations)
như một cấu trúc ngôn ngữ mới.

Trong một hàm thông thường, một tham số trải trên một tập giá trị bị giới hạn bởi kiểu của nó.
Tương tự, trong một hàm (hoặc kiểu) tổng quát, một tham số kiểu trải trên một tập kiểu bị giới hạn
bởi [_ràng buộc kiểu_](/ref/spec#Type_constraints) của nó.
Vì vậy, một ràng buộc kiểu định nghĩa _tập kiểu_ được phép
dùng làm đối số kiểu.

Go 1.18 cũng thay đổi cách ta nhìn nhận interface: trong quá khứ, một
interface định nghĩa một tập phương thức; giờ đây một interface định nghĩa một tập kiểu.
Cách nhìn mới này hoàn toàn tương thích ngược:
với bất kỳ tập phương thức nào do một interface định nghĩa, ta có thể hình dung ra tập (vô hạn)
của mọi kiểu triển khai những phương thức đó.
Ví dụ, với interface [`io.Writer`](/pkg/io#Writer),
ta có thể hình dung tập vô hạn của mọi kiểu có phương thức `Write`
với chữ ký thích hợp.
Tất cả các kiểu đó đều _triển khai_ interface vì chúng đều có
phương thức `Write` được yêu cầu.

Nhưng cách nhìn theo tập kiểu mới mạnh hơn cách nhìn theo tập phương thức cũ:
ta có thể mô tả tập kiểu một cách tường minh, chứ không chỉ gián tiếp qua phương thức.
Điều này cho ta những cách mới để kiểm soát tập kiểu.
Bắt đầu từ Go 1.18, một interface có thể nhúng không chỉ các interface khác,
mà còn là bất kỳ kiểu nào, một hợp của các kiểu, hoặc một tập vô hạn kiểu có cùng
[underlying type](/ref/spec#Underlying_types). Những kiểu này sau đó được đưa vào
[quá trình tính tập kiểu](/ref/spec#General_interfaces):
ký hiệu hợp `A|B` nghĩa là "kiểu `A` hoặc kiểu `B`",
còn ký hiệu `~T` biểu thị "mọi kiểu có underlying type là `T`".
Ví dụ, interface

```Go
interface {
	~int | ~string
	io.Writer
}
```

định nghĩa tập mọi kiểu mà underlying type của chúng là `int` hoặc `string`
và đồng thời cũng triển khai phương thức `Write` của `io.Writer`.

Những generalized interface như vậy không thể dùng làm kiểu biến.
Nhưng vì chúng mô tả tập kiểu nên chúng được dùng làm ràng buộc kiểu, vốn
là các tập kiểu.
Ví dụ, ta có thể viết một hàm `min` tổng quát

```Go
func min[P interface{ ~int64 | ~float64 }](x, y P) P
```

nhận bất kỳ đối số `int64` hoặc `float64` nào.
(Dĩ nhiên, một hiện thực thực tế hơn sẽ dùng một ràng buộc liệt kê
mọi kiểu cơ bản có toán tử <code>&lt;</code>.)

Nhân tiện, vì việc liệt kê tường minh các kiểu không có phương thức là chuyện phổ biến,
một chút [syntactic sugar](https://en.wikipedia.org/wiki/Syntactic_sugar)
cho phép chúng ta [bỏ phần `interface{}` bao quanh](/ref/spec#General_interfaces),
dẫn tới cách viết ngắn gọn và đúng thành ngữ hơn

```Go
func min[P ~int64 | ~float64](x, y P) P { … }
```

Với cách nhìn theo tập kiểu, ta cũng cần một cách mới để giải thích điều gì có nghĩa là
[_triển khai_](/ref/spec#Implementing_an_interface) một interface.
Ta nói một kiểu (không phải interface) `T` triển khai
interface `I` nếu `T` là một phần tử của tập kiểu của interface đó.
Nếu `T` bản thân là một interface, nó mô tả một tập kiểu. Mọi kiểu đơn lẻ trong tập đó
cũng phải có trong tập kiểu của `I`, nếu không `T` sẽ chứa các kiểu không triển khai `I`.
Vì vậy, nếu `T` là một interface, nó triển khai interface `I` nếu tập kiểu
của `T` là tập con của tập kiểu của `I`.

Giờ ta đã có đủ thành phần để hiểu việc thỏa mãn ràng buộc.
Như đã thấy trước đó, một ràng buộc kiểu mô tả tập các kiểu đối số được chấp nhận
cho một tham số kiểu. Một đối số kiểu thỏa mãn ràng buộc của tham số kiểu tương ứng
nếu đối số kiểu đó nằm trong tập do interface ràng buộc mô tả.
Đây là một cách khác để nói rằng đối số kiểu triển khai ràng buộc.
Trong Go 1.18 và Go 1.19, việc thỏa mãn ràng buộc đồng nghĩa với việc triển khai ràng buộc.
Như ta sẽ thấy ngay sau đây, trong Go 1.20 việc thỏa mãn ràng buộc không còn hoàn toàn giống
việc triển khai ràng buộc nữa.

## Các phép toán trên giá trị của tham số kiểu

Một ràng buộc kiểu không chỉ xác định những đối số kiểu nào được chấp nhận cho tham số kiểu,
nó còn quyết định các phép toán nào có thể thực hiện trên giá trị của tham số kiểu đó.
Như ta mong đợi, nếu một ràng buộc định nghĩa một phương thức như `Write`,
thì phương thức `Write` có thể được gọi trên giá trị của tham số kiểu tương ứng.
Nói rộng hơn, một phép toán như `+` hay `*` mà được mọi kiểu trong tập kiểu
do ràng buộc định nghĩa hỗ trợ thì sẽ được phép với giá trị của tham số kiểu tương ứng.

Ví dụ, với ví dụ `min`, trong thân hàm mọi phép toán được
kiểu `int64` và `float64` hỗ trợ đều được phép trên giá trị của tham số kiểu `P`.
Điều đó bao gồm mọi phép toán số học cơ bản, nhưng cũng bao gồm so sánh như <code>&lt;</code>.
Tuy vậy, nó không bao gồm các phép toán bit như `&` hoặc `|`
vì các phép toán đó không được định nghĩa trên giá trị `float64`.

## Các kiểu comparable

Trái với các phép toán một ngôi và hai ngôi khác, `==` không chỉ được định nghĩa trên
một tập giới hạn các [kiểu được khai báo sẵn](/ref/spec#Types), mà còn trên vô số loại kiểu khác nhau,
bao gồm mảng, struct và interface.
Không thể liệt kê hết các kiểu này trong một ràng buộc.
Ta cần một cơ chế khác để biểu đạt rằng một tham số kiểu phải hỗ trợ `==`
(và dĩ nhiên cả `!=`) nếu ta quan tâm tới nhiều hơn các kiểu được khai báo sẵn.

Chúng tôi giải bài toán này thông qua kiểu được khai báo sẵn
[`comparable`](/ref/spec#Predeclared_identifiers), được đưa vào cùng Go 1.18.
`comparable` là
một kiểu interface có tập kiểu là tập vô hạn của các kiểu comparable, và
có thể được dùng làm ràng buộc bất cứ khi nào ta cần một đối số kiểu hỗ trợ `==`.

Tuy nhiên, tập các kiểu do `comparable` bao gồm lại không hoàn toàn giống
tập tất cả [kiểu comparable](/ref/spec#Comparison_operators) được đặc tả Go định nghĩa.
[Theo cách xây dựng](/ref/spec#Interface_types), tập kiểu do một interface chỉ định
(bao gồm cả `comparable`) không chứa chính interface đó (hay bất kỳ interface nào khác).
Vì vậy, một interface như `any` không nằm trong `comparable`,
dù mọi interface đều hỗ trợ `==`.
Vậy chuyện gì xảy ra?

Việc so sánh các interface (và các kiểu hợp thành chứa chúng) có thể panic lúc chạy:
điều này xảy ra khi dynamic type, tức kiểu của giá trị thực được lưu trong
biến interface, lại không comparable.
Hãy xem lại ví dụ `lookupTable` ban đầu: nó chấp nhận giá trị khóa tùy ý.
Nhưng nếu ta thử đưa vào một giá trị có khóa không hỗ trợ `==`, chẳng hạn
một giá trị slice, ta sẽ gặp panic lúc chạy:

```Go
lookupTable[[]int{}] = "slice"  // PANIC: runtime error: hash of unhashable type []int
```

Ngược lại, `comparable` chỉ chứa các kiểu mà compiler bảo đảm sẽ không panic khi dùng `==`.
Chúng tôi gọi những kiểu đó là _strictly comparable_.

Hầu hết thời gian, đây chính xác là điều ta muốn: thật yên tâm khi biết rằng `==` trong một hàm tổng quát
sẽ không panic nếu các toán hạng bị ràng buộc bởi `comparable`, và đó cũng là điều
ta trực giác mong đợi.

Thật không may, định nghĩa này của `comparable` cùng với các quy tắc
thỏa mãn ràng buộc lại ngăn chúng tôi viết mã tổng quát hữu ích, như kiểu
`genericLookupTable` được trình bày trước đó:
để `any` là kiểu đối số chấp nhận được, `any` phải thỏa mãn (và vì vậy triển khai) `comparable`.
Nhưng tập kiểu của `any` lớn hơn (không phải tập con của) tập kiểu của `comparable`
nên nó không triển khai `comparable`.

```Go
var lookupTable GenericLookupTable[any, string] // ERROR: any does not implement comparable (Go 1.18 and Go 1.19)
```

Người dùng nhận ra vấn đề này từ sớm và đã gửi hàng loạt issue và proposal trong thời gian ngắn
([#51338](/issue/51338),
[#52474](/issue/52474),
[#52531](/issue/52531),
[#52614](/issue/52614),
[#52624](/issue/52624),
[#53734](/issue/53734),
v.v.).
Rõ ràng đây là vấn đề mà chúng tôi cần giải quyết.

Giải pháp “hiển nhiên” là đơn giản đưa cả những kiểu không strictly comparable vào
tập kiểu `comparable`.
Nhưng điều đó dẫn tới những bất nhất với mô hình tập kiểu.
Hãy xét ví dụ sau:

```Go
func f[Q comparable]() { … }

func g[P any]() {
        _ = f[int] // (1) ok: int implements comparable
        _ = f[P]   // (2) error: type parameter P does not implement comparable
        _ = f[any] // (3) error: any does not implement comparable (Go 1.18, Go.19)
}
```

Hàm `f` yêu cầu một đối số kiểu strictly comparable.
Rõ ràng việc khởi tạo `f` với `int` là ổn: giá trị `int` không bao giờ panic với `==`
nên `int` triển khai `comparable` (trường hợp 1).
Mặt khác, khởi tạo `f` với `P` là không được phép: tập kiểu của `P` được định nghĩa
bởi ràng buộc `any`, và `any` là viết tắt cho tập mọi kiểu có thể có.
Tập này bao gồm cả những kiểu hoàn toàn không comparable.
Vì vậy, `P` không triển khai `comparable` và không thể được dùng để khởi tạo `f`
(trường hợp 2).
Và cuối cùng, dùng chính kiểu `any` (thay vì tham số kiểu bị ràng buộc bởi `any`)
cũng không được, vì đúng cùng lý do đó (trường hợp 3).

Tuy nhiên, trong trường hợp này chúng ta lại thật sự muốn có thể dùng kiểu `any` làm đối số kiểu.
Con đường duy nhất để thoát khỏi thế tiến thoái lưỡng nan này là thay đổi ngôn ngữ theo cách nào đó.
Nhưng thay đổi như thế nào?

## Triển khai interface so với thỏa mãn ràng buộc

Như đã nhắc trước đó, thỏa mãn ràng buộc chính là triển khai interface:
một đối số kiểu `T` thỏa mãn ràng buộc `C` nếu `T` triển khai `C`.
Điều này hợp lý: `T` phải nằm trong tập kiểu mà `C` kỳ vọng, vốn
chính là định nghĩa của việc triển khai interface.

Nhưng đó cũng chính là vấn đề, vì nó ngăn chúng ta dùng những
kiểu không strictly comparable làm đối số kiểu cho `comparable`.

Vì vậy, cho Go 1.20, sau gần một năm thảo luận công khai về vô số lựa chọn thay thế
(xem các issue nêu ở trên), chúng tôi quyết định đưa vào một ngoại lệ chỉ cho riêng trường hợp này.
Để tránh sự bất nhất, thay vì thay đổi ý nghĩa của `comparable`,
chúng tôi phân biệt giữa _triển khai interface_,
vốn liên quan tới việc truyền giá trị cho biến, và _thỏa mãn ràng buộc_,
vốn liên quan tới việc truyền đối số kiểu cho tham số kiểu.
Khi đã tách riêng, chúng tôi có thể cho mỗi khái niệm này những quy tắc
hơi khác nhau, và đó chính là điều chúng tôi đã làm với proposal [#56548](/issue/56548).

Tin tốt là ngoại lệ này khá cục bộ trong [đặc tả](/ref/spec#Satisfying_a_type_constraint).
Việc thỏa mãn ràng buộc gần như vẫn giống hệt triển khai interface, với một lưu ý:

> Một kiểu `T` thỏa mãn ràng buộc `C` nếu
>
> - `T` triển khai `C`; hoặc
> - `C` có thể được viết dưới dạng `interface{ comparable; E }`, trong đó `E` là một basic interface
>   và `T` là [comparable](/ref/spec#Comparison_operators) và triển khai `E`.

Gạch đầu dòng thứ hai là ngoại lệ.
Không đi quá sâu vào hình thức của đặc tả, điều ngoại lệ này nói là:
một ràng buộc `C` kỳ vọng các kiểu strictly comparable (và có thể còn có các yêu cầu khác
như phương thức `E`) sẽ được thỏa mãn bởi bất kỳ đối số kiểu `T` nào hỗ trợ `==`
(và cũng triển khai các phương thức trong `E`, nếu có).
Hoặc nói ngắn gọn hơn: một kiểu hỗ trợ `==` cũng thỏa mãn `comparable`
(dù nó có thể không triển khai nó).

Ta có thể thấy ngay sự thay đổi này là tương thích ngược:
trước Go 1.20, thỏa mãn ràng buộc cũng giống như triển khai interface, và ta vẫn
còn quy tắc đó (gạch đầu dòng thứ nhất).
Mọi mã dựa vào quy tắc đó vẫn tiếp tục hoạt động như trước.
Chỉ khi quy tắc đó thất bại thì ta mới cần xem ngoại lệ.

Hãy quay lại ví dụ trước:

```Go
func f[Q comparable]() { … }

func g[P any]() {
        _ = f[int] // (1) ok: int satisfies comparable
        _ = f[P]   // (2) error: type parameter P does not satisfy comparable
        _ = f[any] // (3) ok: satisfies comparable (Go 1.20)
}
```

Giờ đây `any` quả thực thỏa mãn (nhưng không triển khai!) `comparable`.
Tại sao?
Bởi vì Go cho phép dùng `==` với giá trị kiểu `any`
(ứng với kiểu `T` trong quy tắc của đặc tả),
và vì ràng buộc `comparable` (ứng với ràng buộc `C` trong quy tắc)
có thể được viết dưới dạng `interface{ comparable; E }` trong đó `E` đơn giản là interface rỗng
trong ví dụ này (trường hợp 3).

Điều thú vị là `P` vẫn không thỏa mãn `comparable` (trường hợp 2).
Lý do là `P` là một tham số kiểu bị ràng buộc bởi `any` (nó _không phải_ là `any`).
Phép toán `==` _không_ khả dụng với mọi kiểu trong tập kiểu của `P`
và vì thế không khả dụng trên `P`;
nó không phải là một [kiểu comparable](/ref/spec#Comparison_operators).
Vì vậy ngoại lệ không áp dụng.
Nhưng điều này ổn: ta vẫn muốn biết rằng `comparable`, yêu cầu strict comparability,
được cưỡng chế trong phần lớn trường hợp. Ta chỉ cần ngoại lệ cho
các kiểu Go hỗ trợ `==`, về cơ bản là vì lý do lịch sử:
chúng ta vốn luôn có khả năng so sánh các kiểu không strictly comparable.

## Hệ quả và cách khắc phục

Chúng ta, những gopher, luôn tự hào rằng hành vi đặc thù của ngôn ngữ
có thể được giải thích và quy về một bộ quy tắc khá gọn, được viết rõ
trong đặc tả ngôn ngữ.
Qua nhiều năm, chúng tôi đã tinh chỉnh các quy tắc này, và khi có thể thì làm cho chúng
đơn giản hơn, thường cũng tổng quát hơn.
Chúng tôi cũng cẩn thận giữ các quy tắc độc lập,
luôn cảnh giác trước những hệ quả ngoài ý muốn và không may.
Tranh luận được giải quyết bằng cách tra đặc tả, chứ không phải bằng mệnh lệnh.
Đó là điều chúng tôi theo đuổi từ khi Go ra đời.

_Không ai có thể đơn giản thêm một ngoại lệ vào một hệ thống kiểu đã được chế tác cẩn thận
mà không có hệ quả!_

Vậy cái giá phải trả là gì?
Có một nhược điểm khá rõ (dù nhẹ), và một nhược điểm ít rõ hơn (nhưng nghiêm trọng hơn).
Hiển nhiên, giờ đây ta có một quy tắc phức tạp hơn cho việc thỏa mãn ràng buộc,
và có thể nói rằng nó kém thanh nhã hơn trước.
Điều này có lẽ không ảnh hưởng đáng kể tới công việc hằng ngày của chúng ta.

Nhưng ta thực sự phải trả giá cho ngoại lệ này: trong Go 1.20, các hàm tổng quát
dựa vào `comparable` không còn hoàn toàn an toàn kiểu một cách tĩnh nữa.
Các phép toán `==` và `!=` có thể panic nếu áp dụng lên toán hạng thuộc
tham số kiểu `comparable`, dù khai báo nói rằng
chúng strictly comparable.
Chỉ một giá trị không comparable cũng có thể len lỏi qua
nhiều hàm hoặc kiểu tổng quát bằng đường của một đối số kiểu không strictly
comparable duy nhất và gây panic.
Trong Go 1.20, giờ đây ta có thể khai báo

```Go
var lookupTable genericLookupTable[any, string]
```

mà không có lỗi biên dịch, nhưng sẽ gặp panic lúc chạy
nếu ta từng dùng một kiểu khóa không strictly comparable trong trường hợp này, giống hệt như ta sẽ
gặp với kiểu `map` dựng sẵn.
Chúng ta đã từ bỏ an toàn kiểu tĩnh để đổi lấy kiểm tra lúc chạy.

Có thể sẽ có những tình huống mà điều này là chưa đủ,
và ta muốn cưỡng chế strict comparability.
Quan sát sau cho phép ta làm chính xác điều đó, ít nhất ở dạng giới hạn:
tham số kiểu không được hưởng ngoại lệ mà ta đã thêm vào
quy tắc thỏa mãn ràng buộc.
Ví dụ, trong ví dụ trước đó, tham số kiểu `P` trong hàm
`g` bị ràng buộc bởi `any` (bản thân nó comparable nhưng không strictly comparable)
nên `P` không thỏa mãn `comparable`.
Ta có thể dùng kiến thức này để tạo ra một kiểu “assert compile time”
cho một kiểu `T` cho trước:

```Go
type T struct { … }
```

Ta muốn khẳng định rằng `T` strictly comparable.
Thật hấp dẫn khi viết như sau:

```Go
// isComparable may be instantiated with any type that supports ==
// including types that are not strictly comparable because of the
// exception for constraint satisfaction.
func isComparable[_ comparable]() {}

// Tempting but not quite what we want: this declaration is also
// valid for types T that are not strictly comparable.
var _ = isComparable[T] // compile-time error if T does not support ==
```

Khai báo biến giả (blank) này đóng vai trò như “assertion” của chúng ta.
Nhưng vì ngoại lệ trong quy tắc thỏa mãn ràng buộc,
`isComparable[T]` chỉ thất bại nếu `T` hoàn toàn không comparable;
nó vẫn thành công nếu `T` hỗ trợ `==`.
Ta có thể khắc phục vấn đề này bằng cách dùng `T` không phải làm đối số kiểu,
mà làm ràng buộc kiểu:

```Go
func _[P T]() {
	_ = isComparable[P] // P supports == only if T is strictly comparable
}
```

Đây là ví dụ trên playground [thành công](/play/p/9i9iEto3TgE) và [thất bại](/play/p/5d4BeKLevPB)
minh họa cơ chế này.

## Quan sát cuối cùng

Điều thú vị là cho tới hai tháng trước khi Go 1.18
được phát hành, compiler đã hiện thực việc thỏa mãn ràng buộc đúng như cách chúng ta làm
giờ đây trong Go 1.20.
Nhưng vì vào thời điểm đó thỏa mãn ràng buộc có nghĩa là triển khai interface,
nên chúng ta đã có một hiện thực không nhất quán với đặc tả ngôn ngữ.
Chúng tôi được cảnh báo về sự việc đó qua [issue #50646](/issue/50646).
Chúng tôi đang ở rất gần thời điểm phát hành và phải ra quyết định nhanh.
Trong bối cảnh chưa có giải pháp thuyết phục, có vẻ an toàn hơn khi làm cho
hiện thực nhất quán với đặc tả.
Một năm sau, với rất nhiều thời gian để cân nhắc các cách tiếp cận khác nhau,
có vẻ hiện thực mà chúng tôi từng có lại chính là hiện thực mà chúng tôi muốn ngay từ đầu.
Chúng tôi đã đi trọn một vòng.

Như thường lệ, vui lòng cho chúng tôi biết nếu có gì đó không hoạt động như mong đợi
bằng cách gửi issue tại [https://go.dev/issue/new](/issue/new).

Xin cảm ơn!
