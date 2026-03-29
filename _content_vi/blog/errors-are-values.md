---
title: Lỗi là giá trị
date: 2015-01-12
by:
- Rob Pike
summary: Các thành ngữ và mẫu xử lý lỗi trong Go.
---


Một chủ đề thường được bàn tới giữa các lập trình viên Go,
đặc biệt là những người mới với ngôn ngữ này, là cách xử lý lỗi.
Cuộc trò chuyện thường chuyển thành lời than phiền về số lần mà chuỗi mã

	if err != nil {
		return err
	}

xuất hiện.
Gần đây chúng tôi đã quét tất cả các dự án mã nguồn mở mà chúng tôi có thể tìm thấy và
phát hiện rằng đoạn mã này chỉ xuất hiện khoảng một lần trên mỗi một hoặc hai trang,
ít thường xuyên hơn nhiều so với những gì người ta hay nghĩ.
Dù vậy, nếu ấn tượng vẫn tồn tại rằng người ta phải gõ

	if err != nil

mọi lúc, thì hẳn phải có gì đó không ổn, và mục tiêu hiển nhiên sẽ là chính Go.

Điều này vừa đáng tiếc, vừa gây hiểu lầm, và rất dễ sửa.
Có lẽ điều đang xảy ra là lập trình viên mới với Go hỏi rằng,
"Người ta xử lý lỗi như thế nào?", học mẫu này rồi dừng lại ở đó.
Trong những ngôn ngữ khác, người ta có thể dùng khối try-catch hay cơ chế tương tự để xử lý lỗi.
Vì vậy lập trình viên nghĩ rằng, khi tôi vốn sẽ dùng try-catch
trong ngôn ngữ cũ, thì trong Go tôi chỉ việc gõ `if` `err` `!=` `nil`.
Theo thời gian, mã Go tích lũy rất nhiều đoạn như vậy, và kết quả tạo cảm giác cồng kềnh.

Bất kể lời giải thích này có đúng hay không,
rõ ràng là những lập trình viên Go đó đã bỏ lỡ một điểm cốt lõi về lỗi:
_Lỗi là giá trị._

Giá trị có thể được lập trình, và vì lỗi là giá trị, lỗi cũng có thể được lập trình.

Dĩ nhiên, một câu lệnh phổ biến khi làm việc với giá trị lỗi là kiểm tra nó có nil hay không,
nhưng còn vô số điều khác người ta có thể làm với một giá trị lỗi,
và việc áp dụng một số điều đó có thể làm chương trình của bạn tốt hơn,
loại bỏ phần lớn boilerplate phát sinh khi mọi lỗi đều được kiểm tra bằng một câu lệnh if máy móc.

Đây là một ví dụ đơn giản từ kiểu [`Scanner`](/pkg/bufio/#Scanner) của
gói `bufio`.
Phương thức [`Scan`](/pkg/bufio/#Scanner.Scan) của nó thực hiện I/O bên dưới,
điều hiển nhiên có thể dẫn tới lỗi.
Thế nhưng bản thân phương thức `Scan` lại không hề để lộ lỗi.
Thay vào đó, nó trả về một giá trị boolean, và một phương thức riêng, được gọi ở cuối vòng quét,
sẽ báo xem có lỗi xảy ra hay không.
Mã phía người dùng trông như sau:

	scanner := bufio.NewScanner(input)
	for scanner.Scan() {
		token := scanner.Text()
		// process token
	}
	if err := scanner.Err(); err != nil {
		// process the error
	}

Đúng là vẫn có một lần kiểm tra nil cho lỗi, nhưng nó chỉ xuất hiện và chạy một lần.
Phương thức `Scan` thay vào đó hoàn toàn có thể được định nghĩa là

	func (s *Scanner) Scan() (token []byte, error)

và khi đó mã ví dụ phía người dùng có thể là (tùy vào cách lấy token),

	scanner := bufio.NewScanner(input)
	for {
		token, err := scanner.Scan()
		if err != nil {
			return err // hoặc có thể break
		}
		// process token
	}

Điều này không khác quá nhiều, nhưng có một điểm phân biệt quan trọng.
Trong đoạn mã này, phía người dùng phải kiểm tra lỗi ở mỗi vòng lặp,
nhưng trong API `Scanner` thật, việc xử lý lỗi được tách khỏi yếu tố cốt lõi của API,
tức việc lặp qua các token.
Vì thế với API thật, mã phía người dùng có cảm giác tự nhiên hơn:
lặp cho tới khi xong, rồi mới lo đến lỗi.
Xử lý lỗi không làm mờ luồng điều khiển.

Tất nhiên, điều đang diễn ra bên dưới là,
ngay khi `Scan` gặp lỗi I/O, nó ghi lại lỗi đó và trả về `false`.
Một phương thức khác, [`Err`](/pkg/bufio/#Scanner.Err),
sẽ báo giá trị lỗi khi phía người dùng hỏi.
Dù đơn giản, điều này không giống với việc rải

	if err != nil

khắp nơi hoặc buộc phía người dùng phải kiểm tra lỗi sau mỗi token.
Đó là lập trình với giá trị lỗi.
Lập trình đơn giản, đúng vậy, nhưng vẫn là lập trình.

Điều đáng nhấn mạnh là dù thiết kế thế nào,
điều tối quan trọng vẫn là chương trình phải kiểm tra lỗi theo cách mà chúng được phơi bày.
Phần thảo luận ở đây không nói về việc tránh kiểm tra lỗi,
mà là về việc dùng ngôn ngữ để xử lý lỗi một cách duyên dáng.

Chủ đề mã kiểm tra lỗi lặp đi lặp lại xuất hiện khi tôi tham dự GoCon mùa thu 2014 ở Tokyo.
Một gopher đầy nhiệt huyết, dùng tên [`@jxck_`](https://twitter.com/jxck_) trên Twitter,
lặp lại lời than phiền quen thuộc về việc kiểm tra lỗi.
Anh ấy có một đoạn mã đại khái như thế này:

	_, err = fd.Write(p0[a:b])
	if err != nil {
		return err
	}
	_, err = fd.Write(p1[c:d])
	if err != nil {
		return err
	}
	_, err = fd.Write(p2[e:f])
	if err != nil {
		return err
	}
	// and so on

Nó rất lặp lại.
Trong mã thật, vốn dài hơn,
còn có nhiều thứ khác đang diễn ra nên không dễ chỉ việc refactor bằng một hàm trợ giúp,
nhưng ở dạng lý tưởng hóa này, một function literal đóng trên biến lỗi sẽ giúp:

	var err error
	write := func(buf []byte) {
		if err != nil {
			return
		}
		_, err = w.Write(buf)
	}
	write(p0[a:b])
	write(p1[c:d])
	write(p2[e:f])
	// and so on
	if err != nil {
		return err
	}

Mẫu này hoạt động tốt, nhưng lại đòi hỏi một closure trong mỗi hàm thực hiện ghi;
một hàm trợ giúp riêng thì vụng về hơn vì biến `err`
cần được duy trì qua các lần gọi (hãy thử xem).

Ta có thể làm điều này sạch hơn, tổng quát hơn và tái sử dụng được bằng cách mượn ý tưởng từ
phương thức `Scan` ở trên.
Tôi đã nhắc tới kỹ thuật này trong cuộc trò chuyện nhưng `@jxck_` không thấy rõ cách áp dụng.
Sau một hồi trao đổi dài, phần nào bị cản trở bởi rào cản ngôn ngữ,
tôi hỏi liệu tôi có thể mượn laptop của anh ấy để gõ thử vài dòng mã hay không.

Tôi định nghĩa một đối tượng tên là `errWriter`, đại loại như sau:

	type errWriter struct {
		w   io.Writer
		err error
	}

và cho nó một phương thức tên là `write`.
Nó không cần phải có chữ ký `Write` chuẩn,
và việc viết thường một phần là để nhấn mạnh sự khác biệt.
Phương thức `write` gọi phương thức `Write` của `Writer` bên dưới
và ghi lại lỗi đầu tiên để tham chiếu về sau:

	func (ew *errWriter) write(buf []byte) {
		if ew.err != nil {
			return
		}
		_, ew.err = ew.w.Write(buf)
	}

Ngay khi có lỗi xảy ra, phương thức `write` trở thành no-op nhưng giá trị lỗi vẫn được giữ lại.

Với kiểu `errWriter` và phương thức `write` của nó, đoạn mã ở trên có thể được refactor:

	ew := &errWriter{w: fd}
	ew.write(p0[a:b])
	ew.write(p1[c:d])
	ew.write(p2[e:f])
	// and so on
	if ew.err != nil {
		return ew.err
	}

Cách này sạch hơn, kể cả so với việc dùng closure,
và cũng làm cho chuỗi thao tác ghi thực sự đang được thực hiện dễ nhìn thấy hơn trên trang.
Không còn sự bừa bộn nữa.
Việc lập trình với giá trị lỗi (và interface) đã khiến mã đẹp hơn.

Rất có thể một đoạn mã khác trong cùng package có thể xây tiếp từ ý tưởng này,
hoặc thậm chí dùng trực tiếp `errWriter`.

Ngoài ra, một khi `errWriter` đã tồn tại, nó còn có thể làm nhiều thứ hơn để hỗ trợ,
đặc biệt trong những ví dụ bớt giả tạo hơn.
Nó có thể tích lũy số byte.
Nó có thể gộp các lần ghi vào một bộ đệm duy nhất để sau đó truyền đi theo cách nguyên tử.
Và còn nhiều nữa.

Thực tế, mẫu này xuất hiện thường xuyên trong thư viện chuẩn.
Các gói [`archive/zip`](/pkg/archive/zip/) và
[`net/http`](/pkg/net/http/) đều dùng nó.
Liên quan trực tiếp hơn tới chủ đề này, [`Writer` của gói `bufio`](/pkg/bufio/)
thực ra chính là một hiện thực của ý tưởng `errWriter`.
Dù `bufio.Writer.Write` có trả về lỗi,
điều đó chủ yếu là để tuân theo interface [`io.Writer`](/pkg/io/#Writer).
Phương thức `Write` của `bufio.Writer` hành xử hệt như phương thức `errWriter.write`
ở trên, với `Flush` là nơi báo lỗi, vì vậy ví dụ của ta có thể được viết như sau:

	b := bufio.NewWriter(fd)
	b.Write(p0[a:b])
	b.Write(p1[c:d])
	b.Write(p2[e:f])
	// and so on
	if b.Flush() != nil {
		return b.Flush()
	}

Cách tiếp cận này có một nhược điểm đáng kể, ít nhất với một số ứng dụng:
không có cách nào biết được đã xử lý xong bao nhiêu trước khi lỗi xảy ra.
Nếu thông tin đó quan trọng, thì cần một cách tiếp cận tinh vi hơn.
Tuy nhiên trong nhiều trường hợp, một lần kiểm tra tất cả hoặc không gì cả ở cuối là đủ.

Chúng ta mới chỉ xem một kỹ thuật để tránh mã xử lý lỗi lặp đi lặp lại.
Hãy nhớ rằng việc dùng `errWriter` hay `bufio.Writer` không phải là cách duy nhất để đơn giản hóa xử lý lỗi,
và cách tiếp cận này cũng không phù hợp cho mọi tình huống.
Dù vậy, bài học then chốt là lỗi là giá trị và toàn bộ sức mạnh của
ngôn ngữ Go đều sẵn có để xử lý chúng.

Hãy dùng ngôn ngữ để đơn giản hóa việc xử lý lỗi của bạn.

Nhưng hãy nhớ: Dù bạn làm gì, luôn luôn kiểm tra lỗi!

Cuối cùng, để đọc trọn vẹn câu chuyện về cuộc tương tác của tôi với @jxck_, bao gồm cả đoạn video nhỏ anh ấy đã ghi,
hãy ghé [blog của anh ấy](http://jxck.hatenablog.com/entry/golang-error-handling-lesson-by-rob-pike).
