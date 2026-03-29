---
title: Có gì trong một cái tên (bí danh)?
date: 2024-09-17
by:
- Robert Griesemer
tags:
- type aliases
- type parameters
- generics
summary: Mô tả về các kiểu bí danh tổng quát, một tính năng được lên kế hoạch cho Go 1.24
---

Bài viết này nói về các kiểu bí danh tổng quát, chúng là gì và tại sao chúng ta cần chúng.


## Bối cảnh

Go được thiết kế cho việc lập trình ở quy mô lớn.
Lập trình ở quy mô lớn có nghĩa là xử lý lượng dữ liệu lớn, nhưng
cũng là các codebase lớn, với nhiều kỹ sư cùng làm việc trên những codebase đó
trong thời gian dài.

Việc Go tổ chức mã nguồn thành các package cho phép lập trình ở quy mô lớn
bằng cách chia các codebase lớn thành những phần nhỏ hơn, dễ quản lý hơn,
thường do những người khác nhau viết, và được kết nối thông qua
API công khai.
Trong Go, các API này bao gồm các định danh được export bởi một package:
các hằng, kiểu, biến và hàm được export.
Điều này cũng bao gồm các trường được export của struct và các method của kiểu.

Khi các dự án phần mềm phát triển theo thời gian hoặc yêu cầu thay đổi,
cách tổ chức mã nguồn thành các package ban đầu có thể trở nên
không còn phù hợp và cần được _refactor_.
Việc refactor có thể liên quan đến việc di chuyển các định danh được export và
các khai báo tương ứng của chúng từ package cũ sang package mới.
Điều này cũng đòi hỏi mọi tham chiếu tới các khai báo đã di chuyển phải được
cập nhật để trỏ tới vị trí mới.
Trong các codebase lớn, việc thực hiện thay đổi như vậy một cách nguyên tử
có thể là không thực tế hoặc không khả thi; nói cách khác, thực hiện việc di chuyển
và cập nhật tất cả client trong một thay đổi duy nhất.
Thay vào đó, thay đổi phải diễn ra theo từng bước: chẳng hạn, để "di chuyển"
một hàm `F`, ta thêm khai báo của nó trong package mới mà không
xóa khai báo gốc trong package cũ.
Bằng cách đó, client có thể được cập nhật dần dần theo thời gian.
Khi mọi nơi gọi đều tham chiếu tới `F` trong package mới, khai báo gốc
của `F` có thể được xóa an toàn (trừ khi nó phải được giữ lại vô thời hạn vì
lý do tương thích ngược).
Russ Cox mô tả chi tiết việc refactor trong bài viết năm 2016 của ông về
[Codebase Refactoring (with help from Go)](/talks/2016/refactor.article).

Di chuyển một hàm `F` từ package này sang package khác trong khi vẫn giữ nó
ở package gốc là việc dễ dàng: chỉ cần một hàm wrapper.
Để di chuyển `F` từ `pkg1` sang `pkg2`, `pkg2` khai báo một hàm `F` mới
(hàm wrapper) với cùng chữ ký như `pkg1.F`, và `pkg2.F`
gọi `pkg1.F`.
Caller mới có thể gọi `pkg2.F`, caller cũ có thể gọi `pkg1.F`, nhưng trong cả hai
trường hợp, hàm cuối cùng được gọi vẫn là cùng một hàm.

Việc di chuyển hằng cũng đơn giản tương tự.
Biến cần thêm một chút công sức: có thể phải đưa vào một con trỏ tới
biến gốc trong package mới hoặc có lẽ dùng các hàm accessor.
Điều này kém lý tưởng hơn, nhưng ít nhất vẫn khả thi.
Điểm ở đây là đối với hằng, biến và hàm,
các tính năng ngôn ngữ hiện có cho phép refactor tăng dần như
mô tả ở trên.

Nhưng còn việc di chuyển một kiểu thì sao?

Trong Go, [(qualified) identifier](/ref/spec#Qualified_identifiers),
hay gọi ngắn gọn là _tên_, quyết định _định danh_
của kiểu:
một kiểu `T` được [định nghĩa](/ref/spec#Type_definitions) và export bởi package
`pkg1` là [khác](/ref/spec#Type_identity) với một định nghĩa kiểu `T`
_giống hệt về mặt khác_ được export bởi package `pkg2`.
Tính chất này làm phức tạp việc di chuyển `T` từ package này sang package khác
trong khi vẫn giữ một bản sao của nó ở package gốc.
Ví dụ, một giá trị có kiểu `pkg2.T` không thể [gán được](/ref/spec#Assignability)
cho một biến có kiểu `pkg1.T` vì tên kiểu và do đó định danh kiểu của chúng
là khác nhau.
Trong giai đoạn cập nhật tăng dần, client có thể có các giá trị và biến
của cả hai kiểu, dù ý định của lập trình viên là chúng phải có
cùng một kiểu.

Để giải quyết vấn đề này, [Go 1.9](/doc/go1.9) đã giới thiệu khái niệm
[_type alias_](/ref/spec#Alias_declarations).
Kiểu bí danh cung cấp một tên mới cho một kiểu hiện có mà không tạo ra
một kiểu mới với định danh khác.

Trái với một [định nghĩa kiểu](/ref/spec#Type_definitions) thông thường

```
type T T0
```

vốn khai báo một _kiểu mới_ không bao giờ đồng nhất với kiểu ở phía bên phải
của khai báo, một [khai báo bí danh](/ref/spec#Alias_declarations)

```
type A = T  // dấu "=" cho biết đây là một khai báo bí danh
```
chỉ khai báo một _tên mới_ `A` cho kiểu ở phía bên phải:
ở đây, `A` và `T` biểu thị cùng một kiểu `T` và do đó là đồng nhất.

Các khai báo bí danh cho phép cung cấp một tên mới (trong một package mới!)
cho một kiểu nhất định mà vẫn giữ nguyên định danh kiểu:

```
package pkg2

import "path/to/pkg1"

type T = pkg1.T
```

Tên kiểu đã thay đổi từ `pkg1.T` thành `pkg2.T` nhưng các giá trị
có kiểu `pkg2.T` có cùng kiểu với các biến có kiểu `pkg1.T`.


## Các kiểu bí danh tổng quát

[Go 1.18](/doc/go1.18) đã giới thiệu generics.
Kể từ bản phát hành đó, các định nghĩa kiểu và khai báo
hàm có thể được tùy biến thông qua các tham số kiểu.
Vì lý do kỹ thuật, các kiểu bí danh chưa có được khả năng tương tự vào thời điểm đó.
Hiển nhiên, khi ấy cũng chưa có các codebase lớn export các kiểu tổng quát
và cần refactor chúng.

Ngày nay, generics đã xuất hiện được vài năm, và các codebase lớn
đang sử dụng các tính năng tổng quát.
Sớm hay muộn nhu cầu refactor các codebase này sẽ xuất hiện, và cùng với đó là
nhu cầu di chuyển các kiểu tổng quát từ package này sang package khác.

Để hỗ trợ các đợt refactor tăng dần có liên quan đến kiểu tổng quát, bản phát hành Go 1.24 trong tương lai,
dự kiến vào đầu tháng 2 năm 2025, sẽ hỗ trợ đầy đủ tham số kiểu trên các kiểu bí danh
theo đề xuất [#46477](/issue/46477).
Cú pháp mới tuân theo cùng một mẫu như đối với định nghĩa kiểu và khai báo hàm,
với một danh sách tham số kiểu tùy chọn theo sau định danh (tên bí danh) ở phía bên trái.
Trước thay đổi này, ta chỉ có thể viết:

```
type Alias = someType
```

nhưng giờ đây ta cũng có thể khai báo tham số kiểu cùng với khai báo bí danh:

```
type Alias[P1 C1, P2 C2] = someType
```

Hãy xem lại ví dụ trước đó, nhưng lần này với các kiểu tổng quát.
Package gốc `pkg1` đã khai báo và export một kiểu tổng quát `G` với tham số kiểu `P`
được ràng buộc phù hợp:

```
package pkg1

type Constraint      someConstraint
type G[P Constraint] someType
```

Nếu phát sinh nhu cầu cung cấp quyền truy cập tới cùng kiểu `G` từ một package mới `pkg2`,
một kiểu bí danh tổng quát chính là thứ cần dùng [(playground)](/play/p/wKOf6NbVtdw?v=gotip):

```
package pkg2

import "path/to/pkg1"

type Constraint      = pkg1.Constraint  // cũng có thể dùng trực tiếp pkg1.Constraint trong G
type G[P Constraint] = pkg1.G[P]
```

Lưu ý rằng ta **không thể** chỉ viết đơn giản

```
type G = pkg1.G
```

vì một vài lý do:

1) Theo [các quy tắc hiện có của đặc tả](/ref/spec#Type_definitions), các kiểu tổng quát
phải được [khởi tạo](/ref/spec#Instantiations) khi chúng
được _sử dụng_.
Phía bên phải của khai báo bí danh sử dụng kiểu `pkg1.G` và
do đó các đối số kiểu phải được cung cấp.
Không làm như vậy sẽ đòi hỏi một ngoại lệ cho trường hợp này, khiến đặc tả
phức tạp hơn.
Không rõ sự tiện lợi nhỏ này có đáng để đánh đổi lấy sự phức tạp đó hay không.

2) Nếu khai báo bí danh không cần tự khai báo các tham số kiểu của riêng nó mà
thay vào đó chỉ đơn giản "kế thừa" chúng từ kiểu được đặt bí danh `pkg1.G`, thì khai báo của
`G` sẽ không cho thấy rằng nó là một kiểu tổng quát.
Các tham số kiểu và ràng buộc của nó sẽ phải được truy ra từ khai báo
của `pkg1.G` (bản thân nó cũng có thể là một bí danh).
Tính dễ đọc sẽ bị ảnh hưởng, trong khi mã dễ đọc là một trong những mục tiêu cốt lõi của dự án Go.

Việc viết ra một danh sách tham số kiểu tường minh thoạt nhìn có vẻ là một gánh nặng
không cần thiết, nhưng nó cũng mang lại thêm tính linh hoạt.
Thứ nhất, số lượng tham số kiểu do kiểu bí danh khai báo không nhất thiết phải
khớp với số lượng tham số kiểu của kiểu được đặt bí danh.
Hãy xem một kiểu map tổng quát:

```
type Map[K comparable, V any] mapImplementation
```

Nếu việc dùng `Map` như tập hợp là phổ biến, thì bí danh

```
type Set[K comparable] = Map[K, bool]
```

có thể hữu ích [(playground)](/play/p/IxeUPGCztqf?v=gotip).
Bởi vì nó là một bí danh, các kiểu như `Set[int]` và `Map[int, bool]` là
đồng nhất.
Điều này sẽ không đúng nếu `Set` là một kiểu [được định nghĩa](/ref/spec#Type_definitions)
(không phải bí danh).

Hơn nữa, các ràng buộc kiểu của một kiểu bí danh tổng quát không nhất thiết phải
khớp với các ràng buộc của kiểu được đặt bí danh, chúng chỉ cần
[thỏa mãn](/ref/spec#Satisfying_a_type_constraint) chúng.

Ví dụ, tiếp tục dùng ví dụ tập hợp ở trên, ta có thể định nghĩa
một `IntSet` như sau:

```
type integers interface{ ~int | ~int8 | ~int16 | ~int32 | ~int64 }
type IntSet[K integers] = Set[K]
```

Map này có thể được khởi tạo với bất kỳ kiểu khóa nào thỏa mãn ràng buộc
`integers` [(playground)](/play/p/0f7hOAALaFb?v=gotip).
Bởi vì `integers` thỏa mãn `comparable`, tham số kiểu `K` có thể được dùng
làm đối số kiểu cho tham số `K` của `Set`, theo các quy tắc khởi tạo
thông thường.

Cuối cùng, vì một bí danh cũng có thể biểu thị một type literal, các bí danh
có tham số cho phép tạo ra các type literal tổng quát
[(playground)](/play/p/wql3NJaUs0o?v=gotip):

```
type Point3D[E any] = struct{ x, y, z E }
```

Để nói rõ, không ví dụ nào trong số này là "trường hợp đặc biệt" hay bằng cách nào đó cần
thêm quy tắc mới vào đặc tả. Chúng tuân theo trực tiếp từ việc áp dụng
các quy tắc hiện có đã được đưa vào cho generics. Điều duy nhất thay đổi trong
đặc tả là khả năng khai báo tham số kiểu trong một khai báo bí danh.


## Một đoạn chuyển ý về tên kiểu

Trước khi có kiểu bí danh, Go chỉ có một dạng khai báo kiểu:

```
type TypeName existingType
```

Khai báo này tạo ra một kiểu mới và khác với một kiểu hiện có
và đặt cho kiểu mới đó một cái tên.
Khi đó việc gọi các kiểu như vậy là _named types_ là tự nhiên vì chúng có một _type name_
trái ngược với các [type literal](/ref/spec#Types) vô danh như
`struct{ x, y int }`.

Với việc đưa kiểu bí danh vào Go 1.9, ta cũng có thể đặt
tên (một bí danh) cho các type literal. Ví dụ, hãy xem:

```
type Point2D = struct{ x, y int }
```

Đột nhiên, khái niệm _named type_ để mô tả một thứ khác với
type literal không còn thật sự có nhiều ý nghĩa nữa, vì tên bí danh rõ ràng là
một tên cho một kiểu, và do đó kiểu được biểu thị (có thể là một type literal, không phải một type name!)
hợp lý mà nói cũng có thể được gọi là một "named type".

Vì các named type (đúng nghĩa) có những thuộc tính đặc biệt (có thể gắn method vào chúng,
chúng tuân theo các quy tắc gán khác nhau, v.v.), có vẻ thận trọng hơn nếu dùng một
thuật ngữ mới để tránh nhầm lẫn.
Vì vậy, kể từ Go 1.9, đặc tả gọi các kiểu trước đây được gọi là named types là _defined types_:
chỉ các defined types mới có những thuộc tính (method, giới hạn khả năng gán, v.v.) gắn
với tên của chúng.
Defined types được giới thiệu thông qua định nghĩa kiểu, còn kiểu bí danh
được giới thiệu thông qua khai báo bí danh.
Trong cả hai trường hợp, kiểu đều được gán tên.

Việc giới thiệu generics trong Go 1.18 khiến mọi thứ phức tạp hơn.
Tham số kiểu cũng là các kiểu, chúng có tên, và chúng chia sẻ các quy tắc
với defined types.
Ví dụ, giống như defined types, hai tham số kiểu có tên khác nhau
biểu thị hai kiểu khác nhau.
Nói cách khác, tham số kiểu là named types, và hơn nữa, chúng
hành xử tương tự các named types ban đầu của Go ở một vài khía cạnh.

Chưa hết, các kiểu được khai báo sẵn của Go (`int`, `string`, v.v.)
chỉ có thể được truy cập thông qua tên của chúng, và giống như defined types và
tham số kiểu, chúng khác nhau nếu tên của chúng khác nhau
(tạm bỏ qua các kiểu bí danh `byte` và `rune`).
Các kiểu được khai báo sẵn thật sự là named types.

Do đó, với Go 1.18, đặc tả đã quay một vòng trọn vẹn và chính thức
đưa trở lại khái niệm [named type](/ref/spec#Types), giờ đây bao gồm
"predeclared types, defined types, và type parameters".
Để hiệu chỉnh trường hợp kiểu bí danh biểu thị type literal, đặc tả nói:
"Một bí danh biểu thị một named type nếu kiểu được nêu trong khai báo bí danh
là một named type."

Lùi lại một bước và tạm rời khỏi hệ thuật ngữ của Go trong chốc lát, thuật ngữ
kỹ thuật chính xác cho một named type trong Go có lẽ là
[_nominal type_](https://en.wikipedia.org/wiki/Nominal_type_system).
Định danh của nominal type gắn tường minh với tên của nó, đúng chính xác với cách
named types của Go (theo thuật ngữ 1.18 hiện nay) vận hành.
Hành vi của nominal type đối lập với _structural type_, vốn có
hành vi chỉ phụ thuộc vào cấu trúc chứ không phụ thuộc vào tên
(nếu ngay từ đầu nó có tên).
Gộp tất cả lại, các kiểu khai báo sẵn, kiểu được định nghĩa và kiểu tham số của Go đều là
nominal types, còn type literals của Go và các bí danh biểu thị type literal
là structural types.
Cả nominal lẫn structural type đều có thể có tên, nhưng việc có tên
không có nghĩa kiểu đó là nominal, nó chỉ có nghĩa là nó được đặt tên.

Không điều nào trong số này thật sự quan trọng đối với việc dùng Go hằng ngày và trên thực tế các chi tiết
có thể được bỏ qua một cách an toàn.
Nhưng thuật ngữ chính xác lại quan trọng trong đặc tả vì nó giúp dễ hơn
trong việc mô tả các quy tắc chi phối ngôn ngữ.
Vậy đặc tả có nên đổi thuật ngữ thêm một lần nữa không?
Có lẽ không đáng với sự xáo trộn đó: không chỉ đặc tả cần
được cập nhật, mà còn rất nhiều tài liệu hỗ trợ khác.
Một số lượng không nhỏ sách viết về Go có thể trở nên thiếu chính xác.
Hơn nữa, "named", dù kém chính xác hơn, có lẽ lại trực quan dễ hiểu hơn
đối với đa số mọi người.
Nó cũng khớp với thuật ngữ gốc được dùng trong đặc tả, dù giờ đây
nó cần một ngoại lệ cho các kiểu bí danh biểu thị type literal.


## Khả dụng

Việc triển khai các kiểu bí danh tổng quát mất nhiều thời gian hơn dự kiến:
những thay đổi cần thiết đòi hỏi phải bổ sung một kiểu `Alias` mới được export
vào [`go/types`](/pkg/go/types) rồi sau đó thêm khả năng ghi nhận các tham số kiểu
cho kiểu đó.
Ở phía compiler, các thay đổi tương ứng cũng đòi hỏi sửa đổi
định dạng dữ liệu export, định dạng tệp mô tả các export
của một package, giờ đây cần có khả năng mô tả tham số kiểu cho bí danh.
Tác động của các thay đổi này không chỉ giới hạn ở compiler mà còn ảnh hưởng
tới các client của `go/types` và do đó là nhiều package của bên thứ ba.
Đây thật sự là một thay đổi tác động tới một codebase lớn; để tránh
làm hỏng mọi thứ, cần có một lộ trình triển khai tăng dần qua nhiều bản phát hành.

Sau tất cả công việc này, các kiểu bí danh tổng quát cuối cùng sẽ khả dụng mặc định trong Go 1.24.

Để các client bên thứ ba có thể chuẩn bị mã nguồn của họ, bắt đầu từ
Go 1.23, hỗ trợ cho kiểu bí danh tổng quát có thể được bật bằng cách đặt
`GOEXPERIMENT=aliastypeparams` khi gọi công cụ `go`.
Tuy nhiên, hãy lưu ý rằng hỗ trợ cho các bí danh tổng quát được export vẫn còn
thiếu trong phiên bản đó.

Hỗ trợ đầy đủ (bao gồm export) đã được triển khai ở tip, và thiết lập
mặc định cho `GOEXPERIMENT` sẽ sớm được chuyển sang để các kiểu bí danh
được bật theo mặc định.
Do đó, một lựa chọn khác là thử nghiệm với phiên bản Go mới nhất
ở tip.

Như thường lệ, vui lòng cho chúng tôi biết nếu bạn gặp bất kỳ vấn đề nào bằng cách tạo một
[issue](/issue/new);
chúng ta càng kiểm thử kỹ một tính năng mới, quá trình triển khai rộng rãi sẽ càng trơn tru.

Cảm ơn và chúc refactor vui vẻ!
