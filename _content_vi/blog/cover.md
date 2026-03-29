---
title: Câu chuyện về cover
date: 2013-12-02
by:
- Rob Pike
tags:
- tools
- coverage
- testing
summary: Giới thiệu công cụ đo độ bao phủ mã của Go 1.12.
template: true
---

## Giới thiệu

Ngay từ đầu dự án, Go đã được thiết kế với công cụ trong đầu.
Những công cụ đó bao gồm một số mảnh công nghệ mang tính biểu tượng nhất của Go như công cụ trình bày tài liệu [godoc](/cmd/godoc), công cụ định dạng mã [gofmt](/cmd/gofmt), và công cụ viết lại API [gofix](/cmd/fix).
Có lẽ quan trọng nhất là [`go` command](/cmd/go), chương trình tự động cài đặt, xây dựng và kiểm thử các chương trình Go chỉ bằng chính mã nguồn như đặc tả build.

Việc phát hành Go 1.2 giới thiệu một công cụ mới cho test coverage có cách tiếp cận khác thường với cách nó tạo ra thống kê coverage, một cách tiếp cận xây trên nền công nghệ do godoc và các công cụ bạn bè của nó đặt ra.

## Hỗ trợ cho công cụ

Trước tiên, một chút nền: điều gì có nghĩa là một [ngôn ngữ hỗ trợ công cụ tốt](/talks/2012/splash.article#TOC_17.)?
Điều đó có nghĩa là ngôn ngữ khiến việc viết công cụ tốt trở nên dễ dàng và hệ sinh thái của nó hỗ trợ xây dựng công cụ đủ mọi loại.

Có một số đặc tính của Go khiến nó phù hợp với tooling.
Trước hết, Go có cú pháp đều đặn, dễ phân tích cú pháp.
Ngữ pháp của nó hướng tới việc không có các trường hợp đặc biệt đòi hỏi bộ máy phân tích phức tạp.

Khi có thể, Go dùng các cấu trúc từ vựng và cú pháp để làm cho các thuộc tính ngữ nghĩa trở nên dễ hiểu.
Ví dụ như việc dùng chữ cái viết hoa để xác định các tên được export và các quy tắc phạm vi được đơn giản hóa mạnh mẽ so với các ngôn ngữ khác trong truyền thống C.

Cuối cùng, thư viện chuẩn đi kèm các package chất lượng sản xuất để lex và parse mã nguồn Go.
Chúng cũng bao gồm, một cách ít phổ biến hơn, một package chất lượng sản xuất để pretty-print các cây cú pháp Go.

Kết hợp lại, các package này tạo thành lõi của công cụ gofmt, nhưng pretty-printer xứng đáng được nhấn mạnh riêng.
Vì nó có thể nhận một cây cú pháp Go tùy ý rồi xuất ra mã đúng định dạng chuẩn, dễ đọc với con người, và chính xác, nó tạo ra khả năng xây dựng những công cụ biến đổi cây parse rồi xuất ra mã đã chỉnh sửa nhưng vẫn đúng và dễ đọc.

Một ví dụ là công cụ gofix, tự động hóa việc viết lại mã để dùng các tính năng ngôn ngữ mới hoặc thư viện đã cập nhật.
Gofix cho phép chúng tôi thực hiện những thay đổi căn bản đối với ngôn ngữ và thư viện trong [giai đoạn chuẩn bị cho Go 1.0](/blog/the-path-to-go-1), với sự tự tin rằng người dùng chỉ cần chạy công cụ để cập nhật mã nguồn của họ lên phiên bản mới nhất.

Bên trong Google, chúng tôi đã dùng gofix để tạo ra những thay đổi quét ngang trong một kho mã khổng lồ mà gần như không tưởng với các ngôn ngữ khác mà chúng tôi dùng.
Không còn cần hỗ trợ nhiều phiên bản của cùng một API nữa; chúng tôi có thể dùng gofix để cập nhật cả công ty chỉ trong một thao tác.

Dĩ nhiên không chỉ các công cụ lớn như vậy được những package này cho phép.
Chúng cũng làm cho việc viết các chương trình khiêm tốn hơn như plugin cho IDE trở nên dễ dàng.
Tất cả những mảnh ghép này xây dựng lẫn nhau, khiến môi trường Go năng suất hơn bằng cách tự động hóa nhiều tác vụ.

## Test coverage

Test coverage là thuật ngữ mô tả mức độ mã của một package được chạy bởi việc thực thi các kiểm thử của package đó.
Nếu việc chạy bộ kiểm thử khiến 80% số câu lệnh nguồn của package được thực thi, ta nói rằng test coverage là 80%.

Chương trình cung cấp test coverage trong Go 1.2 là công cụ mới nhất tận dụng sự hỗ trợ về tooling trong hệ sinh thái Go.

Cách thông thường để tính test coverage là instrument tệp nhị phân.
Ví dụ, chương trình GNU [gcov](http://gcc.gnu.org/onlinedocs/gcc/Gcov.html) đặt breakpoint tại các nhánh được tệp nhị phân thực thi.
Khi mỗi nhánh được thực thi, breakpoint sẽ bị xóa và các câu lệnh đích của nhánh được đánh dấu là “covered”.

Cách tiếp cận này thành công và được dùng rộng rãi. Một công cụ test coverage ban đầu cho Go cũng từng hoạt động như vậy.
Nhưng nó có vấn đề.
Nó khó triển khai, vì việc phân tích thực thi của tệp nhị phân là một việc đầy thách thức.
Nó cũng đòi hỏi một cách đáng tin cậy để gắn vết thực thi trở lại với mã nguồn, điều cũng có thể khó,
đúng như bất kỳ ai từng dùng debugger ở mức mã nguồn đều có thể chứng thực.
Các vấn đề ở đây gồm thông tin debug không chính xác và những thứ như các hàm bị inline làm phức tạp việc phân tích.
Quan trọng hơn cả, cách tiếp cận này rất không khả chuyển.
Nó cần được thực hiện lại cho từng kiến trúc, và ở một mức nào đó cho từng hệ điều hành,
vì hỗ trợ debug thay đổi rất nhiều giữa các hệ thống.

Tuy vậy, nó vẫn hoạt động, và ví dụ nếu bạn là người dùng gccgo thì công cụ gcov có thể cung cấp thông tin test coverage cho bạn.
Tuy nhiên, nếu bạn là người dùng gc, bộ compiler Go được dùng phổ biến hơn, thì cho đến Go 1.2 bạn hoàn toàn không có may mắn đó.

## Test coverage cho Go

Đối với công cụ test coverage mới cho Go, chúng tôi chọn một cách tiếp cận khác tránh phải debug động.
Ý tưởng rất đơn giản: viết lại mã nguồn của package trước khi biên dịch để thêm instrumentation, biên dịch và chạy mã nguồn đã chỉnh sửa đó, rồi dump ra thống kê.
Việc viết lại rất dễ sắp xếp vì lệnh `go` kiểm soát toàn bộ luồng từ nguồn tới kiểm thử rồi tới thực thi.

Đây là một ví dụ. Giả sử ta có một package đơn giản chỉ gồm một tệp như sau:

{{code "cover/pkg.go"}}

và kiểm thử này:

{{code "cover/pkg_test.go"}}

Để lấy test coverage cho package,
ta chạy kiểm thử với coverage được bật bằng cách thêm cờ `-cover` cho `go test`:

	% go test -cover
	PASS
	coverage: 42.9% of statements
	ok  	size	0.026s
	%

Hãy lưu ý coverage là 42.9%, không được tốt lắm.
Trước khi hỏi làm sao tăng con số này, hãy xem nó được tính như thế nào.

Khi test coverage được bật, `go test` chạy công cụ “cover”, một chương trình riêng đi kèm bản phân phối, để viết lại mã nguồn trước khi biên dịch. Đây là phiên bản đã viết lại của hàm `Size`:

{{code "cover/pkg.cover" `/func/` `/^}/`}}

Mỗi phần có thể thực thi của chương trình đều được chú thích bằng một câu lệnh gán mà khi được thực thi sẽ ghi nhận rằng phần đó đã chạy.
Counter đó được gắn với vị trí nguồn gốc của các câu lệnh mà nó đếm thông qua một cấu trúc dữ liệu chỉ đọc thứ hai cũng do công cụ cover sinh ra.
Khi lần chạy kiểm thử hoàn tất, các counter được thu thập và phần trăm được tính bằng cách xem bao nhiêu cái đã được đặt.

Dù câu lệnh gán để chú thích trông có vẻ đắt đỏ, nó biên dịch xuống chỉ một lệnh “move”.
Do đó overhead lúc chạy của nó là vừa phải, chỉ thêm khoảng 3% khi chạy một kiểm thử điển hình (thực tế hơn).
Điều đó khiến việc đưa test coverage vào pipeline phát triển chuẩn trở thành điều hợp lý.

## Xem kết quả

Test coverage cho ví dụ của chúng ta là kém.
Để biết tại sao, ta yêu cầu `go test` ghi ra một “coverage profile”, một tệp giữ các thống kê đã thu thập để ta có thể nghiên cứu chúng kỹ hơn.
Làm điều này rất dễ: dùng cờ `-coverprofile` để chỉ định tệp đầu ra:

	% go test -coverprofile=coverage.out
	PASS
	coverage: 42.9% of statements
	ok  	size	0.030s
	%

(Cờ `-coverprofile` tự động bật `-cover` để kích hoạt phân tích coverage.)
Kiểm thử chạy giống như trước, nhưng kết quả được lưu vào tệp.
Để nghiên cứu chúng, ta tự chạy công cụ test coverage, không thông qua `go test`.
Đầu tiên, ta có thể yêu cầu phân rã coverage theo hàm, dù trong trường hợp này điều đó không soi sáng nhiều vì chỉ có một hàm:

	% go tool cover -func=coverage.out
	size.go:	Size          42.9%
	total:      (statements)  42.9%
	%

Một cách xem dữ liệu thú vị hơn nhiều là lấy bản trình bày HTML của mã nguồn được trang trí bằng thông tin coverage.
Cách hiển thị này được kích hoạt bằng cờ `-html`:

	$ go tool cover -html=coverage.out

Khi chạy lệnh này, một cửa sổ trình duyệt bật lên, hiển thị mã nguồn được phủ (xanh lá), không được phủ (đỏ), và không được instrument (xám).
Đây là ảnh chụp màn hình:

{{image "cover/set.png"}}

Với cách trình bày này, điều sai trở nên hiển nhiên: chúng ta đã quên kiểm thử một vài trường hợp!
Và chúng ta có thể thấy chính xác những trường hợp nào, điều khiến việc cải thiện test coverage trở nên dễ dàng.

## Heat map

Một lợi thế lớn của cách tiếp cận ở mức mã nguồn này với test coverage là rất dễ instrument mã theo những cách khác nhau.
Ví dụ, ta có thể hỏi không chỉ một câu lệnh đã được thực thi hay chưa, mà còn nó đã chạy bao nhiêu lần.

Lệnh `go test` chấp nhận cờ `-covermode` để đặt chế độ coverage thành một trong ba thiết lập:

  - set:    mỗi câu lệnh có chạy không?
  - count:  mỗi câu lệnh chạy bao nhiêu lần?
  - atomic: giống count, nhưng đếm chính xác trong các chương trình song song

Mặc định là `set`, tức là điều ta đã thấy.
Thiết lập `atomic` chỉ cần khi cần đếm chính xác trong khi chạy các thuật toán song song. Nó dùng các phép toán nguyên tử từ package [sync/atomic](/pkg/sync/atomic/), vốn có thể khá đắt đỏ.
Tuy nhiên, với đa số mục đích, chế độ `count` hoạt động tốt và, giống như chế độ mặc định `set`, cũng rất rẻ.

Hãy thử đếm số lần thực thi câu lệnh cho một package chuẩn, package định dạng `fmt`.
Ta chạy kiểm thử và ghi ra coverage profile để sau đó có thể trình bày thông tin một cách đẹp hơn.

	% go test -covermode=count -coverprofile=count.out fmt
	ok  	fmt	0.056s	coverage: 91.7% of statements
	%

Đó là tỷ lệ test coverage tốt hơn rất nhiều so với ví dụ trước.
(Tỷ lệ coverage không bị ảnh hưởng bởi chế độ coverage.)
Ta có thể hiển thị bản phân rã theo hàm:

	% go tool cover -func=count.out
	fmt/format.go: init              100.0%
	fmt/format.go: clearflags        100.0%
	fmt/format.go: init              100.0%
	fmt/format.go: computePadding     84.6%
	fmt/format.go: writePadding      100.0%
	fmt/format.go: pad               100.0%
	...
	fmt/scan.go:   advance            96.2%
	fmt/scan.go:   doScanf            96.8%
	total:         (statements)       91.7%

Phần thưởng lớn thực sự nằm ở đầu ra HTML:

	% go tool cover -html=count.out

Đây là cách hàm `pad` trông như thế nào trong bản trình bày đó:

{{image "cover/count.png"}}

Hãy để ý cường độ màu xanh thay đổi ra sao. Những câu lệnh xanh sáng hơn có số lần thực thi cao hơn; những màu xanh ít bão hòa hơn biểu diễn số lần thực thi thấp hơn.
Bạn thậm chí có thể rê chuột lên câu lệnh để xem số đếm thực tế bật lên trong tool tip.
Tại thời điểm bài viết này được viết, số đếm trông như sau
(chúng tôi đã chuyển số đếm từ tooltip sang đầu dòng để dễ hiển thị hơn):

	2933    if !f.widPresent || f.wid == 0 {
	2985        f.buf.Write(b)
	2985        return
	2985    }
	  56    padding, left, right := f.computePadding(len(b))
	  56    if left > 0 {
	  37        f.writePadding(left, padding)
	  37    }
	  56    f.buf.Write(b)
	  56    if right > 0 {
	  13        f.writePadding(right, padding)
	  13    }

Đó là rất nhiều thông tin về việc thực thi của hàm, thông tin có thể hữu ích trong profiling.

## Basic block

Có thể bạn đã để ý rằng số đếm trong ví dụ trước không đúng như bạn mong đợi trên các dòng có dấu ngoặc đóng.
Điều đó là vì, như thường lệ, test coverage là một khoa học không chính xác.

Điều đang xảy ra ở đây đáng để giải thích. Ta muốn các chú thích coverage được phân định bởi các nhánh trong chương trình, giống như cách chúng được làm khi tệp nhị phân được instrument theo phương pháp truyền thống.
Tuy nhiên làm điều đó bằng cách viết lại mã nguồn là chuyện khó, bởi các nhánh không xuất hiện tường minh trong mã nguồn.

Điều mà chú thích coverage thực hiện là instrument các block, thường được ranh giới bởi dấu ngoặc nhọn.
Làm điều này đúng trong trường hợp tổng quát là rất khó.
Một hệ quả của thuật toán được dùng là dấu ngoặc đóng trông như thuộc về block mà nó đóng, trong khi dấu ngoặc mở lại trông như thuộc về bên ngoài block.
Một hệ quả thú vị hơn là trong biểu thức như

	f() && g()

không có nỗ lực nào nhằm instrument riêng biệt các lời gọi tới `f` và `g`. Bất kể thực tế ra sao, nó sẽ luôn trông như cả hai chạy cùng số lần, tức số lần `f` đã chạy.

Công bằng mà nói, ngay cả `gcov` cũng gặp khó ở đây. Công cụ đó làm đúng phần instrumentation nhưng cách trình bày của nó dựa trên dòng và vì vậy có thể bỏ lỡ một số sắc thái.

## Bức tranh lớn

Đó là câu chuyện về test coverage trong Go 1.2.
Một công cụ mới với cách hiện thực thú vị không chỉ cho phép thống kê test coverage, mà còn cả những cách trình bày dễ diễn giải và thậm chí khả năng trích xuất thông tin profiling.

Kiểm thử là một phần quan trọng của phát triển phần mềm và test coverage là một cách đơn giản để thêm tính kỷ luật vào chiến lược kiểm thử của bạn.
Hãy đi và kiểm thử, rồi phủ nó.

