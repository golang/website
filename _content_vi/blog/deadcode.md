---
title: Tìm các hàm không thể với tới bằng deadcode
date: 2023-12-12
by:
- Alan Donovan
summary: deadcode là một lệnh mới giúp xác định những hàm không thể được gọi.
---

Những hàm là một phần của mã nguồn dự án nhưng không bao giờ có thể được chạm tới trong bất kỳ lần thực thi nào được gọi là “dead code”, và chúng kéo lùi nỗ lực bảo trì codebase.
Hôm nay chúng tôi rất vui được chia sẻ một công cụ tên là `deadcode` để giúp bạn xác định chúng.

```
$ go install golang.org/x/tools/cmd/deadcode@latest
$ deadcode -help
The deadcode command reports unreachable functions in Go programs.

Usage: deadcode [flags] package...
```

## Ví dụ

Trong khoảng hơn một năm trở lại đây, chúng tôi đã thực hiện nhiều thay đổi đối với cấu trúc của [gopls](/blog/gopls-scalability), language server cho Go vận hành VS Code và các trình soạn thảo khác.
Một thay đổi điển hình có thể viết lại một hàm hiện có, cẩn thận bảo đảm hành vi mới của nó đáp ứng nhu cầu của toàn bộ các caller hiện tại.
Đôi khi, sau khi đã bỏ ra rất nhiều công sức như vậy, chúng tôi lại bực bội phát hiện rằng một trong các caller thực ra chưa từng được chạm tới trong bất kỳ lần thực thi nào, nên đáng ra có thể xóa an toàn từ đầu.
Nếu biết điều này sớm hơn thì công việc refactor đã dễ hơn.

Chương trình Go đơn giản dưới đây minh họa vấn đề:

```
module example.com/greet
go 1.21
```

```
package main

import "fmt"

func main() {
	var g Greeter
	g = Helloer{}
	g.Greet()
}

type Greeter interface{ Greet() }

type Helloer struct{}
type Goodbyer struct{}

var _ Greeter = Helloer{}  // Helloer  implements Greeter
var _ Greeter = Goodbyer{} // Goodbyer implements Greeter

func (Helloer) Greet()  { hello() }
func (Goodbyer) Greet() { goodbye() }

func hello()   { fmt.Println("hello") }
func goodbye() { fmt.Println("goodbye") }
```

Khi ta thực thi nó, nó in ra hello:

```
$ go run .
hello
```

Từ đầu ra của nó, dễ thấy chương trình này thực thi hàm `hello` nhưng không thực thi hàm `goodbye`.
Điều ít rõ ràng hơn khi chỉ nhìn thoáng qua là hàm `goodbye` không bao giờ có thể được gọi.
Tuy nhiên, ta không thể đơn giản xóa `goodbye`, vì nó cần thiết cho phương thức `Goodbyer.Greet`, phương thức này đến lượt mình lại cần để triển khai interface `Greeter` mà ta thấy phương thức `Greet` của nó được gọi từ `main`.
Nhưng nếu ta lần theo tiến trình đi về phía trước từ main, ta có thể thấy rằng không có giá trị `Goodbyer` nào từng được tạo ra, nên lời gọi `Greet` trong `main` chỉ có thể đi tới `Helloer.Greet`.
Đó là ý tưởng đằng sau thuật toán mà công cụ `deadcode` dùng.

Khi chạy deadcode trên chương trình này, công cụ báo cho ta rằng hàm `goodbye` và phương thức `Goodbyer.Greet` đều không thể chạm tới:

```
$ deadcode .
greet.go:23: unreachable func: goodbye
greet.go:20: unreachable func: Goodbyer.Greet
```

Với kiến thức đó, ta có thể xóa an toàn cả hai hàm, cùng với chính kiểu `Goodbyer`.

Công cụ cũng có thể giải thích vì sao hàm `hello` đang còn sống. Nó phản hồi bằng một chuỗi lời gọi hàm dẫn tới `hello`, bắt đầu từ main:

```
$ deadcode -whylive=example.com/greet.hello .
                  example.com/greet.main
dynamic@L0008 --> example.com/greet.Helloer.Greet
 static@L0019 --> example.com/greet.hello
```

Đầu ra được thiết kế để dễ đọc trên terminal, nhưng bạn có thể dùng cờ `-json` hoặc `-f=template` để chỉ định các định dạng đầu ra phong phú hơn cho công cụ khác tiêu thụ.

## Cách hoạt động

Lệnh `deadcode` [nạp](https://pkg.go.dev/golang.org/x/tools/go/packages), [phân tích cú pháp](https://pkg.go.dev/go/parser), và [kiểm tra kiểu](https://pkg.go.dev/go/types) các package được chỉ định, sau đó chuyển chúng thành một [biểu diễn trung gian](https://pkg.go.dev/golang.org/x/tools/go/ssa) tương tự compiler thông thường.

Sau đó nó dùng một thuật toán tên là [Rapid Type Analysis](https://pkg.go.dev/golang.org/x/tools/go/callgraph/rta) (RTA) để xây dựng tập các hàm có thể chạm tới, ban đầu chỉ gồm các điểm vào của từng package `main`: hàm `main`, và hàm khởi tạo package, nơi gán các biến toàn cục và gọi các hàm tên là `init`.

RTA xem xét các câu lệnh trong thân của từng hàm có thể chạm tới để thu thập ba loại thông tin: tập các hàm mà nó gọi trực tiếp; tập các lời gọi động mà nó thực hiện thông qua phương thức interface; và tập các kiểu mà nó chuyển sang interface.

Các lời gọi hàm trực tiếp thì dễ: ta chỉ thêm callee vào tập các hàm có thể chạm tới, và nếu đó là lần đầu ta gặp callee đó, ta sẽ kiểm tra thân hàm của nó giống như đã làm với main.

Các lời gọi động qua phương thức interface thì khó hơn, vì ta không biết tập các kiểu triển khai interface đó. Ta không muốn giả định rằng mọi phương thức có thể có trong chương trình mà kiểu của nó khớp đều là mục tiêu khả dĩ cho lời gọi, vì một số kiểu trong số đó có thể chỉ được khởi tạo từ dead code! Đó là lý do ta thu thập tập các kiểu được chuyển sang interface: việc chuyển này làm mỗi kiểu đó có thể chạm tới từ `main`, khiến các phương thức của nó giờ đây trở thành các mục tiêu khả dĩ của lời gọi động.

Điều này dẫn đến một tình huống con gà và quả trứng. Khi gặp từng hàm có thể chạm tới mới, ta phát hiện thêm lời gọi phương thức interface và thêm các phép chuyển kiểu cụ thể sang kiểu interface.
Nhưng khi tích Descartes của hai tập này (lời gọi phương thức interface × kiểu cụ thể) lớn dần lên, ta lại phát hiện thêm các hàm có thể chạm tới.
Lớp bài toán này, gọi là “dynamic programming”, có thể được giải bằng cách (về mặt ý niệm) đánh dấu kiểm vào một bảng hai chiều lớn, thêm hàng và cột khi cần, cho tới khi không còn dấu kiểm nào để thêm.
Các dấu kiểm trong bảng cuối cùng cho biết cái gì có thể chạm tới; các ô trống là dead code.

<!--
  Source:
  https://docs.google.com/presentation/d/1DH6Ycdqpt-Zel88lINAuudA6cp0e64ILfHOJq8hJ3v8
  Exported using "File > Download > SVG"
  Cropped using Inkscape "Edit > Resize Page to Selection"
-->  
<div class="image">
<center>
  <img src="deadcode-rta.svg" alt="illustration of Rapid Type Analysis"/><br/>  <i>
   Hàm <code>main</code> làm cho <code>Helloer</code> được
   khởi tạo, và lời gọi <code>g.Greet</code><br/>
   sẽ điều phối đến phương thức <code>Greet</code> của từng kiểu đã được khởi tạo cho tới thời điểm đó.
  </i>
</center>
</div>

Các lời gọi động đến các hàm (không phải phương thức) được xử lý tương tự như interface chỉ có một phương thức.
Và các lời gọi được thực hiện [bằng reflection](https://pkg.go.dev/reflect#Value.Call) được xem như chạm tới bất kỳ phương thức nào của bất kỳ kiểu nào được dùng trong một phép chuyển đổi interface, hoặc bất kỳ kiểu nào có thể suy ra từ đó thông qua package `reflect`.
Nhưng nguyên tắc ở mọi trường hợp đều là như nhau.

## Kiểm thử

RTA là một phân tích toàn chương trình. Điều đó có nghĩa là nó luôn bắt đầu từ một hàm main và tiến về phía trước: bạn không thể bắt đầu từ một package thư viện như `encoding/json`.

Tuy nhiên, đa số package thư viện đều có kiểm thử, và kiểm thử có các hàm main.
Ta không thấy chúng vì chúng được `go test` sinh ra phía sau hậu trường, nhưng ta có thể đưa chúng vào phân tích bằng cờ `-test`.

Nếu việc này báo rằng một hàm trong package thư viện là dead, đó là dấu hiệu cho thấy độ bao phủ kiểm thử của bạn có thể cần được cải thiện.
Ví dụ, lệnh sau liệt kê toàn bộ các hàm trong `encoding/json` không được bất kỳ kiểm thử nào của nó chạm tới:

```
$ deadcode -test -filter=encoding/json encoding/json
encoding/json/decode.go:150:31: unreachable func: UnmarshalFieldError.Error
encoding/json/encode.go:225:28: unreachable func: InvalidUTF8Error.Error
```

(Cờ `-filter` giới hạn đầu ra vào những package khớp với biểu thức chính quy. Mặc định, công cụ báo mọi package trong module ban đầu.)

## Tính đúng đắn

Mọi công cụ phân tích tĩnh [tất yếu](https://en.wikipedia.org/wiki/Rice%27s_theorem) đều tạo ra những xấp xỉ không hoàn hảo của các hành vi động khả dĩ của chương trình đích.
Các giả định và suy luận của công cụ có thể là “sound”, tức bảo thủ nhưng có lẽ quá thận trọng, hoặc “unsound”, tức lạc quan nhưng không phải lúc nào cũng đúng.

Công cụ deadcode cũng không ngoại lệ: nó phải xấp xỉ tập các mục tiêu của lời gọi động thông qua giá trị hàm và interface hoặc thông qua reflection.
Ở khía cạnh này, công cụ là sound. Nói cách khác, nếu nó báo một hàm là dead code, điều đó có nghĩa là hàm đó không thể được gọi ngay cả qua các cơ chế động đó. Tuy nhiên, công cụ có thể không báo ra một số hàm mà trên thực tế cũng không bao giờ được thực thi.

Công cụ deadcode cũng phải xấp xỉ tập các lời gọi được thực hiện từ những hàm không viết bằng Go, tức các hàm mà nó không nhìn thấy.
Ở khía cạnh này, công cụ không sound.
Phân tích của nó không biết đến các hàm chỉ được gọi từ mã assembly, hay việc alias hàm phát sinh từ [`go:linkname` directive](https://pkg.go.dev/cmd/compile#hdr-Compiler_Directives).
May mắn là cả hai tính năng này hiếm khi được dùng bên ngoài runtime Go.

## Hãy thử nó

Chúng tôi chạy `deadcode` định kỳ trên các dự án của mình, đặc biệt sau công việc refactor, để giúp xác định những phần của chương trình không còn cần nữa.

Khi dead code đã được an táng, bạn có thể tập trung loại bỏ những đoạn mã đáng ra đã đến lúc kết thúc nhưng vẫn ngoan cố còn sống, tiếp tục rút cạn sinh lực của bạn. Chúng tôi gọi các hàm undead như vậy là “vampire code”!

Hãy thử xem:

```
$ go install golang.org/x/tools/cmd/deadcode@latest
```

Chúng tôi thấy nó hữu ích, và hy vọng bạn cũng vậy.

