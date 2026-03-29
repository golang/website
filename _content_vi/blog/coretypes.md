---
title: Tạm biệt core types, chào Go như ta biết và yêu mến!
date: 2025-03-26
by:
- Robert Griesemer
summary: Go 1.25 đơn giản hóa đặc tả ngôn ngữ bằng cách loại bỏ khái niệm core type
---

Go 1.18 giới thiệu generics và cùng với đó là một số tính năng mới, bao gồm type parameter, type constraint, và các khái niệm mới như type set.
Nó cũng giới thiệu khái niệm _core type_.
Trong khi những thứ trước đó mang lại chức năng mới cụ thể, core type là một cấu trúc trừu tượng được đưa vào
vì lý do tiện dụng và để đơn giản hóa việc xử lý generic operand (các toán hạng có kiểu là type parameter).
Trong compiler Go, phần mã trước đây dựa vào [underlying type](/ref/spec/#Underlying_types) của một toán hạng
giờ phải gọi một hàm tính core type của toán hạng đó.
Trong đặc tả ngôn ngữ, ở nhiều chỗ chúng ta chỉ cần thay “underlying type” bằng “core type”.
Có gì mà không thích?

Hóa ra là khá nhiều!
Để hiểu vì sao chúng ta đi tới đây, sẽ hữu ích nếu xem lại ngắn gọn cách type parameter và type constraint hoạt động.

## Type parameter và type constraint

Type parameter là chỗ giữ chỗ cho một type argument trong tương lai;
nó hoạt động như một _biến kiểu_ có giá trị được biết ở thời gian biên dịch,
tương tự cách một hằng có tên đại diện cho một số, chuỗi, hay bool có giá trị được biết ở thời gian biên dịch.
Giống như biến thông thường, type parameter có kiểu.
Kiểu đó được mô tả bởi _type constraint_ của nó, thứ xác định
những phép toán nào được phép trên các toán hạng có kiểu là type parameter tương ứng.

Mọi kiểu cụ thể dùng để khởi tạo một type parameter đều phải thỏa mãn constraint của type parameter đó.
Điều này bảo đảm rằng một toán hạng có kiểu là type parameter sở hữu toàn bộ các thuộc tính của type constraint tương ứng,
bất kể kiểu cụ thể nào được dùng để khởi tạo type parameter.

Trong Go, type constraint được mô tả thông qua sự pha trộn giữa yêu cầu về phương thức và yêu cầu về kiểu, cùng nhau
định nghĩa nên một _type set_: đó là tập tất cả các kiểu thỏa mãn mọi yêu cầu. Go dùng một
dạng interface tổng quát hóa cho mục đích này. Một interface liệt kê tập phương thức và kiểu,
và tập kiểu do interface đó mô tả gồm tất cả các kiểu triển khai các phương thức đó
và nằm trong các kiểu đã được liệt kê.

Ví dụ, tập kiểu được mô tả bởi interface

```Go
type Constraint interface {
	~[]byte | ~string
	Hash() uint64
}
```

gồm tất cả các kiểu có biểu diễn là `[]byte` hoặc `string` và có method set bao gồm phương thức `Hash`.

Với điều này, giờ ta có thể viết các quy tắc chi phối phép toán trên generic operand.
Ví dụ, [quy tắc cho biểu thức chỉ số](/ref/spec#Index_expressions) nêu rằng (trong số các điều khác)
đối với một toán hạng `a` có kiểu type parameter `P`:

> Biểu thức chỉ số `a[x]` phải hợp lệ với giá trị của mọi kiểu trong type set của `P`.
> Kiểu phần tử của mọi kiểu trong type set của `P` phải đồng nhất.
  (Trong ngữ cảnh này, kiểu phần tử của kiểu chuỗi là `byte`.)

Những quy tắc này cho phép chỉ số hóa biến tổng quát `s` bên dưới ([playground](/play/p/M1LYKm3x3IB)):

```Go
func at[bytestring Constraint](s bytestring, i int) byte {
	return s[i]
}
```

Phép chỉ số `s[i]` được cho phép vì kiểu của `s` là `bytestring`, và type constraint (type set) của
`bytestring` chứa các kiểu `[]byte` và `string` mà việc chỉ số hóa với `i` là hợp lệ.

## Core type

Cách tiếp cận dựa trên type set này rất linh hoạt và phù hợp với chủ đích của
[đề xuất generics ban đầu](https://go.googlesource.com/proposal/+/refs/heads/master/design/43651-type-parameters.md):
một phép toán liên quan tới các toán hạng kiểu tổng quát phải hợp lệ nếu nó hợp lệ với mọi kiểu được type constraint tương ứng cho phép.
Để đơn giản hóa phần hiện thực, và biết rằng sau này chúng ta vẫn có thể nới lỏng quy tắc,
cách tiếp cận này đã _không_ được áp dụng ở mọi nơi.
Thay vào đó, ví dụ với [câu lệnh send](/ref/spec#Send_statements), đặc tả nói rằng

> _core type_ của biểu thức channel phải là một channel, hướng của channel phải cho phép thao tác gửi,
> và kiểu của giá trị được gửi phải có thể gán cho kiểu phần tử của channel.

Những quy tắc này dựa trên khái niệm core type được định nghĩa xấp xỉ như sau:

- Nếu một kiểu không phải type parameter, core type của nó đơn giản là [underlying type](/ref/spec#Underlying_types).
- Nếu kiểu đó là type parameter, core type là underlying type duy nhất của toàn bộ các kiểu trong type set của type parameter đó.
  Nếu type set có _các_ underlying type khác nhau, core type không tồn tại.

Ví dụ, `interface{ ~[]int }` có một core type (`[]int`), nhưng interface `Constraint` ở trên thì không có core type.
Để phức tạp hơn, đối với các phép toán trên channel và một số lời gọi built-in (`append`, `copy`) thì định nghĩa core type ở trên lại quá hạn chế.
Các quy tắc thực tế có các điều chỉnh cho phép hướng channel khác nhau và các type set chứa đồng thời kiểu `[]byte` và `string`.

Có nhiều vấn đề với cách tiếp cận này:

- Vì định nghĩa core type phải dẫn tới những quy tắc kiểu đúng đắn cho các đặc tính ngôn ngữ khác nhau,
nó trở nên quá hạn chế với một số phép toán cụ thể.
Ví dụ, các quy tắc [slice expression](/ref/spec#Slice_expressions) trong Go 1.24 có dựa vào core type,
và vì thế phép cắt một toán hạng kiểu `S` bị ràng buộc bởi `Constraint` lại không được phép, dù
nó có thể là hợp lệ.

- Khi cố hiểu một đặc tính ngôn ngữ cụ thể, người ta có thể phải học các chi tiết rắc rối của
core type ngay cả khi chỉ đang xét mã không tổng quát.
Một lần nữa, với slice expression, đặc tả ngôn ngữ nói về core type của toán hạng bị cắt,
thay vì chỉ nêu rằng toán hạng phải là một mảng, slice hoặc chuỗi.
Cách sau trực tiếp hơn, đơn giản hơn, rõ ràng hơn, và không đòi hỏi biết thêm một khái niệm khác
có thể hoàn toàn không liên quan trong trường hợp cụ thể.

- Vì khái niệm core type tồn tại, nên các quy tắc cho biểu thức chỉ số, và `len` và `cap` (cũng như các trường hợp khác),
vốn đều không dùng core type, lại trông như các ngoại lệ trong ngôn ngữ chứ không phải chuẩn mực.
Ngược lại, core type khiến những đề xuất như [issue #48522](/issue/48522), cho phép selector
`x.f` truy cập một trường `f` được mọi phần tử trong type set của `x` cùng chia sẻ, trông như đang thêm ngoại lệ nữa vào
ngôn ngữ.
Nếu không có core type, tính năng đó lại trở thành hệ quả tự nhiên và hữu ích của các quy tắc thông thường cho
truy cập trường trong mã không tổng quát.

## Go 1.25

Với bản phát hành Go 1.25 sắp tới (tháng 8 năm 2025), chúng tôi quyết định loại bỏ khái niệm core type khỏi
đặc tả ngôn ngữ để thay bằng văn xuôi tường minh (và tương đương!) ở những nơi cần thiết.
Điều này có nhiều lợi ích:

- Đặc tả Go trình bày ít khái niệm hơn, giúp việc học ngôn ngữ dễ hơn.
- Hành vi của mã không tổng quát có thể được hiểu mà không cần viện dẫn các khái niệm của generics.
- Cách tiếp cận cá thể hóa (quy tắc riêng cho từng phép toán cụ thể) mở ra khả năng cho những quy tắc linh hoạt hơn.
Chúng tôi đã nhắc đến [issue #48522](/issue/48522), nhưng còn có cả ý tưởng cho những
phép toán slice mạnh hơn, và [cải thiện suy luận kiểu](/issue/69153).

[Proposal issue #70128](/issue/70128) tương ứng gần đây đã được chấp thuận và các thay đổi liên quan đã được triển khai.
Cụ thể, điều này có nghĩa là rất nhiều đoạn văn trong đặc tả ngôn ngữ đã được trả về dạng nguyên bản của nó,
trước thời kỳ generics, và các đoạn mới được thêm vào ở nơi cần thiết để giải thích các quy tắc trong trường hợp
liên quan tới generic operand. Điều quan trọng là không có hành vi nào thay đổi.
Toàn bộ phần nói về core type đã bị loại bỏ.
Thông báo lỗi của compiler đã được cập nhật để không còn nhắc tới “core type” nữa, và trong nhiều
trường hợp, thông báo lỗi giờ cụ thể hơn khi chỉ ra chính xác kiểu nào trong một type set
đang gây ra vấn đề.

Dưới đây là một ví dụ về những thay đổi đã được thực hiện. Với hàm built-in `close`,
bắt đầu từ Go 1.18 đặc tả mở đầu như sau:

> Với một đối số `ch` có core type là channel,
> hàm built-in `close` ghi nhận rằng sẽ không còn giá trị nào được gửi trên channel đó nữa.

Một người đọc chỉ đơn giản muốn biết `close` hoạt động thế nào, trước tiên lại phải học về core type.
Bắt đầu từ Go 1.25, phần này sẽ lại mở đầu theo đúng cách từng mở đầu trước Go 1.18:

> Với một channel `ch`, hàm built-in `close(ch)`
> ghi nhận rằng sẽ không còn giá trị nào được gửi trên channel đó nữa.

Cách này ngắn hơn và dễ hiểu hơn.
Chỉ khi người đọc đang xử lý một generic operand thì họ mới cần suy ngẫm về
đoạn mới được thêm vào:

> Nếu kiểu của đối số truyền vào `close` là một type parameter
> thì mọi kiểu trong type set của nó đều phải là channel có cùng kiểu phần tử.
> Sẽ là lỗi nếu bất kỳ channel nào trong số đó là channel chỉ nhận.

Chúng tôi đã thực hiện các thay đổi tương tự ở từng nơi từng nhắc đến core type.
Tóm lại, dù việc cập nhật đặc tả này không ảnh hưởng tới bất kỳ chương trình Go hiện tại nào, nó vẫn mở ra
cánh cửa cho các cải tiến ngôn ngữ trong tương lai đồng thời làm cho ngôn ngữ hiện tại dễ học hơn
và đặc tả của nó đơn giản hơn.

