---
title: Tương thích ngược, Go 1.21 và Go 2
date: 2023-08-14T12:00:00Z
by:
- Russ Cox
summary: Go 1.21 mở rộng cam kết của Go với tương thích ngược, để mỗi bộ công cụ Go mới đều là hiện thực tốt nhất có thể của ngữ nghĩa các bộ công cụ cũ hơn.
---

Go 1.21 bao gồm những tính năng mới nhằm cải thiện khả năng tương thích.
Trước khi bạn ngừng đọc, tôi biết điều đó nghe có vẻ nhàm chán.
Nhưng nhàm chán có thể là điều tốt.
Trong những ngày đầu của Go 1,
Go rất hào hứng và đầy bất ngờ.
Mỗi tuần chúng tôi cắt ra một bản snapshot mới
và mọi người lại tung xúc xắc
để xem chúng tôi đã thay đổi điều gì
và chương trình của họ sẽ hỏng ra sao.
Chúng tôi phát hành Go 1 cùng lời hứa tương thích
để loại bỏ sự phấn khích đó,
để các bản phát hành mới của Go trở nên nhàm chán.

Nhàm chán là tốt.
Nhàm chán là ổn định.
Nhàm chán nghĩa là có thể tập trung vào công việc của bạn,
chứ không phải vào việc Go khác đi thế nào.
Bài viết này nói về công việc quan trọng mà chúng tôi đã đưa vào Go 1.21
để giữ cho Go nhàm chán.

## Tương thích của Go 1 {#go1}

Chúng tôi đã tập trung vào khả năng tương thích trong hơn một thập kỷ.
Với Go 1, từ năm 2012, chúng tôi đã công bố tài liệu
“[Go 1 and the Future of Go Programs](/doc/go1compat)”
nêu lên một ý định rất rõ ràng:

> Chúng tôi dự định rằng các chương trình được viết theo đặc tả Go 1
> sẽ tiếp tục được biên dịch và chạy đúng, không thay đổi,
> trong suốt vòng đời của đặc tả đó. ...
> Các chương trình Go hoạt động hôm nay nên tiếp tục hoạt động
> ngay cả khi những bản phát hành Go 1 trong tương lai xuất hiện.

Có một vài ngoại lệ cho điều đó.
Thứ nhất, tương thích ở đây là tương thích mã nguồn.
Khi bạn cập nhật lên phiên bản Go mới,
bạn vẫn phải biên dịch lại mã của mình.
Thứ hai, chúng tôi có thể thêm API mới,
nhưng không theo cách làm hỏng mã hiện có.

Cuối tài liệu có cảnh báo rằng
“[Không thể bảo đảm rằng không thay đổi nào trong tương lai sẽ làm hỏng bất kỳ chương trình nào.]”
Sau đó nó liệt kê một số lý do khiến chương trình vẫn có thể bị hỏng.

Ví dụ, rõ ràng là nếu chương trình của bạn phụ thuộc vào một hành vi lỗi
và chúng tôi sửa lỗi đó, chương trình sẽ bị hỏng.
Nhưng chúng tôi cố gắng hết sức để làm hỏng càng ít càng tốt và giữ Go thật nhàm chán.
Cho tới nay, có hai cách tiếp cận chính mà chúng tôi đã dùng: kiểm tra API và kiểm thử.

## Kiểm tra API {#api}

Có lẽ sự thật rõ ràng nhất về khả năng tương thích
là chúng tôi không thể lấy đi API, nếu không các chương trình dùng nó sẽ hỏng.

Ví dụ, đây là một chương trình mà ai đó đã viết
và chúng tôi không được làm hỏng:

	package main

	import "os"

	func main() {
		os.Stdout.WriteString("hello, world\n")
	}

Chúng tôi không thể bỏ package `os`;
không thể bỏ biến toàn cục `os.Stdout`, vốn là một `*os.File`;
và cũng không thể bỏ phương thức `WriteString` của `os.File`.
Rõ ràng là loại bỏ bất kỳ thứ nào trong số đó cũng sẽ
làm hỏng chương trình này.

Có lẽ ít rõ ràng hơn là chúng tôi hoàn toàn không thể đổi kiểu của `os.Stdout`.
Giả sử chúng tôi muốn biến nó thành một interface có cùng các phương thức.
Chương trình vừa rồi sẽ không hỏng, nhưng chương trình này thì có:

	package main

	import "os"

	func main() {
		greet(os.Stdout)
	}

	func greet(f *os.File) {
		f.WriteString(“hello, world\n”)
	}

Chương trình này truyền `os.Stdout` vào một hàm tên là `greet`
yêu cầu đối số kiểu `*os.File`.
Vì vậy, đổi `os.Stdout` sang interface sẽ làm hỏng chương trình này.

Để hỗ trợ việc phát triển Go, chúng tôi dùng một công cụ duy trì
danh sách API được export của từng package
trong các tệp riêng biệt với package thực:

	% cat go/api/go1.21.txt
	pkg bytes, func ContainsFunc([]uint8, func(int32) bool) bool #54386
	pkg bytes, method (*Buffer) AvailableBuffer() []uint8 #53685
	pkg bytes, method (*Buffer) Available() int #53685
	pkg cmp, func Compare[$0 Ordered]($0, $0) int #59488
	pkg cmp, func Less[$0 Ordered]($0, $0) bool #59488
	pkg cmp, type Ordered interface {} #59488
	pkg context, func AfterFunc(Context, func()) func() bool #57928
	pkg context, func WithDeadlineCause(Context, time.Time, error) (Context, CancelFunc) #56661
	pkg context, func WithoutCancel(Context) Context #40221
	pkg context, func WithTimeoutCause(Context, time.Duration, error) (Context, CancelFunc) #56661

Một trong các bài kiểm thử chuẩn của chúng tôi kiểm tra rằng API thực tế của package khớp với các tệp đó.
Nếu chúng tôi thêm API mới vào một package, bài kiểm thử sẽ hỏng trừ khi chúng tôi thêm nó vào các tệp API.
Và nếu chúng tôi thay đổi hoặc loại bỏ API, bài kiểm thử cũng sẽ hỏng. Điều này giúp tránh sai sót.
Tuy nhiên, một công cụ như thế chỉ tìm được một lớp vấn đề nhất định, tức các thay đổi và việc gỡ bỏ API.
Còn có những cách khác để tạo ra thay đổi không tương thích cho Go.

Điều đó dẫn chúng ta tới cách tiếp cận thứ hai để giữ Go nhàm chán: kiểm thử.

## Kiểm thử {#testing}

Cách hiệu quả nhất để tìm ra những điểm không tương thích ngoài dự kiến
là chạy các bài kiểm thử hiện có trên phiên bản phát triển của bản phát hành Go tiếp theo.
Chúng tôi kiểm thử phiên bản phát triển của Go trên
toàn bộ mã Go nội bộ của Google theo hình thức luân phiên.
Khi các kiểm thử đều qua, chúng tôi cài commit đó làm
bộ công cụ Go sản xuất của Google.

Nếu một thay đổi làm hỏng các bài kiểm thử bên trong Google,
chúng tôi giả định rằng nó cũng sẽ làm hỏng kiểm thử bên ngoài Google,
và chúng tôi tìm cách giảm tác động.
Phần lớn thời gian, chúng tôi hoàn tác thay đổi hoàn toàn
hoặc tìm cách viết lại nó để không làm hỏng bất kỳ chương trình nào.
Tuy nhiên, đôi khi chúng tôi kết luận rằng thay đổi đó
quan trọng và vẫn “tương thích” dù nó
có làm hỏng một số chương trình.
Khi đó, chúng tôi vẫn cố giảm tác động nhiều nhất có thể,
rồi ghi lại vấn đề tiềm tàng trong ghi chú phát hành.

Dưới đây là hai ví dụ về kiểu vấn đề tương thích tinh vi như vậy
mà chúng tôi tìm thấy bằng việc kiểm thử Go bên trong Google nhưng vẫn đưa vào Go 1.1.

## Literal struct và trường mới {#struct}

Đây là đoạn mã chạy tốt trong Go 1:

	package main

	import "net"

	var myAddr = &net.TCPAddr{
		net.IPv4(18, 26, 4, 9),
		80,
	}

Package `main` khai báo biến toàn cục `myAddr`,
là composite literal của kiểu `net.TCPAddr`.
Trong Go 1, package `net` định nghĩa kiểu `TCPAddr`
là một struct với hai trường, `IP` và `Port`.
Chúng khớp với các trường trong composite literal,
nên chương trình biên dịch được.

Trong Go 1.1, chương trình ngừng biên dịch, với lỗi
“too few initializers in struct literal.”
Vấn đề là chúng tôi đã thêm trường thứ ba, `Zone`, vào `net.TCPAddr`,
và chương trình này không có giá trị cho trường thứ ba đó.
Cách sửa là viết lại chương trình bằng tagged literal,
để nó có thể build ở cả hai phiên bản Go:

	var myAddr = &net.TCPAddr{
		IP:   net.IPv4(18, 26, 4, 9),
		Port: 80,
	}

Vì literal này không chỉ định giá trị cho `Zone`, nó sẽ dùng
zero value (trong trường hợp này là chuỗi rỗng).

Yêu cầu phải dùng composite literal có tên trường cho các
struct trong thư viện chuẩn được nhắc đến rõ ràng trong [tài liệu tương thích](/doc/go1compat),
và `go vet` báo các literal cần gắn nhãn
để bảo đảm tương thích với các phiên bản sau.
Vấn đề này đủ mới trong Go 1.1
để đáng được nhắc ngắn trong ghi chú phát hành.
Ngày nay chúng tôi chỉ đơn giản đề cập đến trường mới.

## Độ chính xác của thời gian {#precision}

Vấn đề thứ hai mà chúng tôi tìm thấy trong quá trình kiểm thử Go 1.1
không liên quan gì đến API cả.
Nó liên quan đến thời gian.

Ngay sau khi Go 1 được phát hành, ai đó chỉ ra rằng
[`time.Now`](/pkg/time/#Now)
trả về thời điểm với độ chính xác micro giây,
nhưng với thêm một ít mã,
nó có thể trả về thời điểm với độ chính xác nano giây.
Nghe có vẻ tốt đúng không?
Chính xác hơn là tốt hơn.
Vì vậy chúng tôi đã thực hiện thay đổi đó.

Điều đó đã làm hỏng một số ít bài kiểm thử bên trong Google có dạng như sau:

	func TestSaveTime(t *testing.T) {
		t1 := time.Now()
		save(t1)
		if t2 := load(); t2 != t1 {
			t.Fatalf("load() = %v, want %v", t1, t2)
		}
	}

Đoạn mã này gọi `time.Now`
rồi truyền giá trị qua
`save` và `load`
và kỳ vọng sẽ lấy lại đúng cùng thời điểm.
Nếu `save` và `load` dùng một cách biểu diễn
chỉ lưu được độ chính xác micro giây,
thì điều đó hoạt động tốt trong Go 1 nhưng thất bại trong Go 1.1.

Để giúp sửa những kiểm thử như vậy,
chúng tôi thêm các phương thức [`Round`](/pkg/time/#Time.Round) và
[`Truncate`](/pkg/time/#Time.Truncate)
để bỏ đi độ chính xác không mong muốn,
và trong ghi chú phát hành,
chúng tôi ghi lại khả năng có vấn đề này
cùng các phương thức mới để hỗ trợ sửa lỗi.

Những ví dụ này cho thấy kiểm thử tìm ra
những dạng không tương thích khác với kiểm tra API.
Dĩ nhiên, kiểm thử cũng không phải là bảo đảm tuyệt đối
cho khả năng tương thích,
nhưng nó đầy đủ hơn việc chỉ kiểm tra API.
Có rất nhiều ví dụ về những vấn đề mà chúng tôi tìm được
khi kiểm thử rồi quyết định rằng chúng thực sự vi phạm
quy tắc tương thích và đã hoàn tác trước khi phát hành.

Ví dụ thay đổi độ chính xác thời gian là một trường hợp thú vị
của thứ đã làm hỏng chương trình nhưng chúng tôi vẫn phát hành
dù sao.
Chúng tôi thực hiện thay đổi đó vì độ chính xác tốt hơn
là điều tốt và được phép trong phạm vi hành vi đã được tài liệu hóa của hàm.

Ví dụ này cho thấy rằng đôi khi, dù đã nỗ lực rất nhiều
và chú ý rất kỹ, vẫn có những lúc thay đổi Go đồng nghĩa với
việc làm hỏng các chương trình Go.
Nói một cách nghiêm ngặt, các thay đổi đó “tương thích”
theo nghĩa của tài liệu Go 1, nhưng chúng vẫn làm hỏng chương trình.
Phần lớn các vấn đề tương thích loại này có thể được xếp
vào một trong ba nhóm:
thay đổi đầu ra,
thay đổi đầu vào,
và thay đổi giao thức.

## Thay đổi đầu ra {#output}

Thay đổi đầu ra xảy ra khi một hàm trả về đầu ra khác
so với trước đây, nhưng đầu ra mới cũng đúng như, hoặc thậm chí đúng hơn,
đầu ra cũ.
Nếu mã hiện có được viết chỉ để chấp nhận đầu ra cũ, nó sẽ hỏng.
Ta vừa thấy một ví dụ với `time.Now` tăng độ chính xác lên nano giây.

**Sort.** Một ví dụ khác xảy ra trong Go 1.6,
khi chúng tôi thay đổi hiện thực của sort
để chạy nhanh hơn khoảng 10%.
Đây là một chương trình ví dụ
sắp xếp danh sách màu theo độ dài tên:

	colors := strings.Fields(
		`black white red orange yellow green blue indigo violet`)
	sort.Sort(ByLen(colors))
	fmt.Println(colors)

	Go 1.5:  [red blue green white black yellow orange indigo violet]
	Go 1.6:  [red blue white green black orange yellow indigo violet]

Thay đổi thuật toán sort thường làm đổi
thứ tự của các phần tử bằng nhau,
và điều đó đã xảy ra ở đây.
Go 1.5 trả về green, white, black, theo thứ tự đó.
Go 1.6 trả về white, green, black.

Rõ ràng sort được phép trả về các kết quả bằng nhau theo bất kỳ thứ tự nào nó muốn,
và thay đổi này giúp nó nhanh hơn 10%, thật tuyệt.
Nhưng những chương trình kỳ vọng một đầu ra cụ thể sẽ bị hỏng.
Đây là ví dụ tốt cho thấy vì sao khả năng tương thích lại khó đến vậy.
Chúng tôi không muốn làm hỏng chương trình,
nhưng cũng không muốn bị khóa chặt
vào các chi tiết triển khai không được tài liệu hóa.

**Compress/flate.** Một ví dụ khác: trong Go 1.8, chúng tôi cải thiện
`compress/flate` để tạo đầu ra nhỏ hơn,
với chi phí CPU và bộ nhớ xấp xỉ như cũ.
Nghe như đôi bên cùng có lợi, nhưng nó lại làm hỏng một dự án bên trong Google
cần các bản build archive có thể tái lập:
giờ họ không thể tái tạo archive cũ nữa.
Họ đã fork `compress/flate` và `compress/gzip`
để giữ một bản sao của thuật toán cũ.

Chúng tôi cũng làm điều tương tự với compiler Go,
dùng một bản fork của package `sort` ([và các package khác](https://go.googlesource.com/go/+/go1.21.0/src/cmd/dist/buildtool.go#22))
để compiler tạo ra cùng kết quả
ngay cả khi nó được build bằng các phiên bản Go cũ hơn.

Đối với các điểm không tương thích kiểu thay đổi đầu ra như thế này,
câu trả lời tốt nhất là viết chương trình và kiểm thử
chấp nhận mọi đầu ra hợp lệ,
và dùng những lần hỏng như vậy
như cơ hội để thay đổi chiến lược kiểm thử,
không chỉ đơn thuần cập nhật kết quả mong đợi.
Nếu bạn cần đầu ra thật sự tái lập,
câu trả lời tốt thứ hai là fork mã
để tự cách ly khỏi các thay đổi,
nhưng hãy nhớ rằng
bạn cũng đang tự cách ly khỏi cả các bản sửa lỗi.

## Thay đổi đầu vào {#input}

Thay đổi đầu vào xảy ra khi một hàm thay đổi tập đầu vào nó chấp nhận
hoặc thay đổi cách nó xử lý chúng.

**ParseInt.** Ví dụ, Go 1.13 bổ sung hỗ trợ cho
dấu gạch dưới trong số lớn để tăng tính dễ đọc.
Cùng với thay đổi ngôn ngữ đó,
chúng tôi làm cho `strconv.ParseInt` chấp nhận cú pháp mới.
Thay đổi này không làm hỏng gì bên trong Google,
nhưng về sau chúng tôi nghe từ một người dùng bên ngoài
rằng mã của họ bị hỏng.
Chương trình của họ dùng các số được phân tách bằng dấu gạch dưới
như một định dạng dữ liệu.
Nó thử `ParseInt` trước và chỉ quay sang kiểm tra dấu gạch dưới nếu `ParseInt` thất bại.
Khi `ParseInt` không còn thất bại nữa, phần mã xử lý dấu gạch dưới ngừng chạy.

**ParseIP.** Một ví dụ khác là `net.ParseIP` của Go,
vốn đi theo các ví dụ trong RFC IP đời đầu,
thường hiển thị địa chỉ IP thập phân với các số 0 ở đầu.
Nó đọc địa chỉ IP 18.032.4.011 là 18.32.4.11, chỉ đơn giản có thêm vài số 0.
Về sau rất lâu, chúng tôi mới phát hiện ra rằng các thư viện C bắt nguồn từ BSD
diễn giải số 0 ở đầu trong địa chỉ IP như khởi đầu của số bát phân:
trong các thư viện đó, 18.032.4.011 nghĩa là 18.26.4.9!

Đây là một sự lệch nghiêm trọng
giữa Go và phần còn lại của thế giới,
nhưng thay đổi ý nghĩa của số 0 ở đầu
từ bản phát hành Go này sang bản phát hành kế tiếp
cũng là một sự lệch nghiêm trọng.
Nó sẽ là một điểm không tương thích rất lớn.
Cuối cùng, chúng tôi quyết định thay đổi `net.ParseIP` trong Go 1.17
để từ chối hoàn toàn số 0 ở đầu.
Cách phân tích chặt chẽ hơn này bảo đảm rằng khi Go và C
cùng phân tích thành công một địa chỉ IP,
hoặc khi Go phiên bản cũ và mới cùng làm được điều đó,
thì tất cả đều nhất trí về ý nghĩa của nó.

Thay đổi này không làm hỏng gì bên trong Google,
nhưng đội Kubernetes lo ngại về
các cấu hình đã lưu trước đây có thể phân tích được
nhưng sẽ ngừng phân tích với Go 1.17.
Địa chỉ có số 0 ở đầu
có lẽ nên bị loại khỏi những cấu hình đó,
vì Go diễn giải chúng khác với
gần như mọi ngôn ngữ khác,
nhưng điều đó nên diễn ra theo lộ trình của Kubernetes, không phải của Go.
Để tránh thay đổi ngữ nghĩa,
Kubernetes bắt đầu dùng bản fork riêng
của `net.ParseIP` nguyên bản.

Phản hồi tốt nhất trước các thay đổi đầu vào là xử lý đầu vào người dùng
bằng cách trước hết kiểm tra cú pháp mà bạn muốn chấp nhận
rồi mới phân tích giá trị,
nhưng đôi khi bạn vẫn cần fork mã.

## Thay đổi giao thức {#protocol}

Kiểu không tương thích phổ biến cuối cùng là thay đổi giao thức.
Thay đổi giao thức là một thay đổi được thực hiện trong một package
nhưng cuối cùng lại lộ ra bên ngoài
thông qua các giao thức mà chương trình dùng
để giao tiếp với thế giới bên ngoài.
Gần như bất kỳ thay đổi nào cũng có thể lộ ra bên ngoài
trong một số chương trình nhất định, như ta thấy với `ParseInt` và `ParseIP`,
nhưng thay đổi giao thức là thứ lộ ra bên ngoài
trong gần như mọi chương trình.

**HTTP/2.** Một ví dụ rõ ràng của thay đổi giao thức là khi
Go 1.6 thêm hỗ trợ tự động cho HTTP/2.
Giả sử một client Go 1.5 đang kết nối tới một
máy chủ hỗ trợ HTTP/2 qua một mạng có middlebox vô tình
làm hỏng HTTP/2.
Vì Go 1.5 chỉ dùng HTTP/1.1 nên chương trình hoạt động bình thường.
Nhưng rồi việc cập nhật lên Go 1.6 lại làm hỏng chương trình, bởi vì Go 1.6
bắt đầu dùng HTTP/2, và trong bối cảnh này thì HTTP/2 không hoạt động.

Go hướng tới hỗ trợ các giao thức hiện đại theo mặc định,
nhưng ví dụ này cho thấy bật HTTP/2 có thể làm hỏng chương trình
mà không phải lỗi của chương trình đó (cũng không phải lỗi của Go).
Nhà phát triển trong tình huống này có thể quay lại dùng Go 1.5,
nhưng điều đó không thật sự thỏa đáng.
Thay vào đó, Go 1.6 ghi nhận thay đổi này trong ghi chú phát hành
và làm cho việc tắt HTTP/2 trở nên đơn giản.

Thực tế, [Go 1.6 ghi lại hai cách](/doc/go1.6#http2) để tắt HTTP/2:
cấu hình tường minh trường `TLSNextProto` bằng API của package,
hoặc đặt biến môi trường GODEBUG:

	GODEBUG=http2client=0 ./myprog
	GODEBUG=http2server=0 ./myprog
	GODEBUG=http2client=0,http2server=0 ./myprog

Như ta sẽ thấy sau đây, Go 1.21 tổng quát hóa cơ chế GODEBUG này
để biến nó thành tiêu chuẩn cho mọi thay đổi có khả năng gây vỡ.

**SHA1.** Đây là một ví dụ tinh vi hơn về thay đổi giao thức.
Không ai nên còn dùng chứng chỉ HTTPS dựa trên SHA1 nữa.
Các nhà cấp phát chứng chỉ đã ngừng cấp chúng từ năm 2015,
và mọi trình duyệt lớn đều ngừng chấp nhận chúng từ năm 2017.
Đầu năm 2020, Go 1.18 tắt hỗ trợ cho chúng theo mặc định,
kèm theo một thiết lập GODEBUG để ghi đè thay đổi đó.
Chúng tôi cũng thông báo ý định gỡ bỏ thiết lập GODEBUG đó trong Go 1.19.

Đội Kubernetes cho chúng tôi biết rằng
một số cài đặt vẫn dùng chứng chỉ SHA1 riêng.
Tạm gác câu hỏi về bảo mật sang một bên,
Kubernetes không phải là nơi ép các doanh nghiệp đó
nâng cấp hạ tầng chứng chỉ của họ,
và việc fork `crypto/tls` cùng `net/http` để giữ hỗ trợ SHA1
sẽ vô cùng đau đớn.
Thay vào đó, chúng tôi đồng ý giữ thiết lập ghi đè thêm lâu hơn dự tính,
để tạo thêm thời gian cho một quá trình chuyển đổi có trật tự.
Dù sao thì chúng tôi cũng muốn làm hỏng càng ít chương trình càng tốt.

## Mở rộng hỗ trợ GODEBUG trong Go 1.21

Để cải thiện tương thích ngược ngay cả trong những trường hợp tinh vi
mà chúng ta vừa xem xét, Go 1.21 mở rộng và chính thức hóa việc dùng GODEBUG.

Trước hết, với bất kỳ thay đổi nào được Go 1 compatibility cho phép nhưng
vẫn có thể làm hỏng chương trình hiện có,
chúng tôi đều làm toàn bộ công việc như đã thấy bên trên để hiểu rõ
các vấn đề tương thích tiềm tàng, và chúng tôi thiết kế thay đổi để giữ được càng nhiều
chương trình hiện có hoạt động càng tốt.
Với những chương trình còn lại, cách tiếp cận mới là:

 1. Chúng tôi sẽ định nghĩa một thiết lập GODEBUG mới cho phép
    từng chương trình riêng lẻ chọn không dùng hành vi mới.
    Một thiết lập GODEBUG có thể không được thêm nếu điều đó là bất khả thi, nhưng trường hợp đó phải cực hiếm.

 2. Những thiết lập GODEBUG được thêm vào để tương thích sẽ được duy trì tối thiểu
    hai năm (bốn bản phát hành Go). Một số, như `http2client` và `http2server`,
    sẽ được duy trì lâu hơn rất nhiều, thậm chí vô thời hạn.

 3. Khi có thể, mỗi thiết lập GODEBUG có một bộ đếm [`runtime/metrics`](/pkg/runtime/metrics/)
    đi kèm, tên là
    `/godebug/non-default-behavior/<name>:events`,
    đếm số lần hành vi của một chương trình cụ thể
    đã thay đổi do giá trị không mặc định của thiết lập đó.
    Ví dụ, khi đặt `GODEBUG=http2client=0`,
    `/godebug/non-default-behavior/http2client:events` đếm
    số lần transport HTTP mà chương trình cấu hình không có hỗ trợ HTTP/2.

 4. Thiết lập GODEBUG của một chương trình được cấu hình để khớp với phiên bản Go
    được liệt kê trong tệp `go.mod` của package `main`.
    Nếu `go.mod` của chương trình ghi `go 1.20` và bạn cập nhật lên
    bộ công cụ Go 1.21, thì mọi hành vi được GODEBUG điều khiển đã thay đổi trong
    Go 1.21 sẽ giữ hành vi cũ của Go 1.20 cho tới khi bạn đổi
    `go.mod` thành `go 1.21`.

 5. Một chương trình có thể thay đổi từng thiết lập GODEBUG riêng lẻ bằng các dòng `//go:debug`
    trong package `main`.

 6. Mọi thiết lập GODEBUG đều được ghi lại trong [một danh sách trung tâm duy nhất](/doc/godebug#history)
    để dễ tra cứu.

Cách tiếp cận này có nghĩa là mỗi phiên bản Go mới phải trở thành hiện thực tốt nhất có thể
của các phiên bản Go cũ hơn, thậm chí vẫn bảo toàn cả các hành vi
đã bị thay đổi theo cách tương thích-nhưng-gây-vỡ ở các bản phát hành mới khi biên dịch mã cũ.

Ví dụ, trong Go 1.21, `panic(nil)` giờ gây ra một runtime panic (không nil),
để cho kết quả của [`recover`](/ref/spec/#Handling_panics) giờ đây báo cáo một cách đáng tin cậy
liệu goroutine hiện tại có đang panic hay không.
Hành vi mới này được điều khiển bởi một thiết lập GODEBUG và do đó phụ thuộc
vào dòng `go` trong `go.mod` của package `main`: nếu nó ghi `go 1.20` hoặc cũ hơn,
`panic(nil)` vẫn được cho phép.
Nếu nó ghi `go 1.21` hoặc mới hơn, `panic(nil)` sẽ trở thành panic với `runtime.PanicNilError`.
Và giá trị mặc định dựa trên phiên bản này có thể được ghi đè tường minh bằng cách thêm dòng sau vào package main:

	//go:debug panicnil=1

Sự kết hợp các tính năng này có nghĩa là chương trình có thể cập nhật lên bộ công cụ mới hơn
trong khi vẫn giữ được hành vi của các bộ công cụ cũ mà chúng từng dùng,
có thể áp dụng kiểm soát chi tiết hơn đối với từng thiết lập cụ thể khi cần,
và có thể dùng theo dõi sản xuất để hiểu những job nào
thực tế đang dùng những hành vi không mặc định này.
Kết hợp lại, chúng sẽ giúp việc triển khai bộ công cụ mới
mượt hơn cả trước đây.

Xem “[Go, Backwards Compatibility, and GODEBUG](/doc/godebug)” để biết thêm chi tiết.

## Cập nhật về Go 2 {#go2}

Trong đoạn trích từ “[Go 1 and the Future of Go Programs](/doc/go1compat)”
ở đầu bài viết, dấu ba chấm đã ẩn đi phần điều kiện sau:

> Tại một thời điểm không xác định trong tương lai, một đặc tả Go 2 có thể xuất hiện,
> nhưng cho tới lúc đó, [... mọi chi tiết về khả năng tương thích ...].

Điều đó làm dấy lên một câu hỏi hiển nhiên: khi nào thì ta nên kỳ vọng
đặc tả Go 2 làm hỏng các chương trình Go 1 cũ?

Câu trả lời là: không bao giờ.
Go 2, theo nghĩa đoạn tuyệt với quá khứ
và không còn biên dịch được các chương trình cũ,
sẽ không bao giờ xảy ra.
Go 2, theo nghĩa là cuộc đại tu lớn của Go 1
mà chúng tôi bắt đầu hướng tới từ năm 2017, thực ra đã xảy ra rồi.

Sẽ không có một Go 2 làm hỏng các chương trình Go 1.
Thay vào đó, chúng tôi sẽ càng đặt nặng tương thích hơn,
vì điều đó có giá trị hơn nhiều so với bất kỳ cuộc đoạn tuyệt nào với quá khứ.
Thực tế, chúng tôi tin rằng việc ưu tiên khả năng tương thích
là quyết định thiết kế quan trọng nhất mà chúng tôi đã đưa ra cho Go 1.

Vì vậy, điều bạn sẽ thấy trong vài năm tới
là rất nhiều công việc mới mẻ, thú vị, nhưng được thực hiện theo cách cẩn trọng,
tương thích, để chúng tôi có thể giữ cho việc nâng cấp của bạn từ
bộ công cụ này sang bộ công cụ khác nhàm chán nhất có thể.
