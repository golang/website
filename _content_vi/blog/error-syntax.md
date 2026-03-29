---
title: "[ Có | Không ] hỗ trợ cú pháp cho xử lý lỗi"
date: 2025-06-03
by:
- Robert Griesemer
tags:
- error
- syntax
- technical
- proposal
summary: Kế hoạch của nhóm Go xoay quanh hỗ trợ xử lý lỗi
---

Một trong những lời phàn nàn lâu đời và dai dẳng nhất về Go liên quan đến sự dài dòng của xử lý lỗi.
Tất cả chúng ta đều rất quen thuộc, có người nói là khá đau đớn, với mẫu mã này:

```Go
x, err := call()
if err != nil {
        // handle err
}
```

Phép kiểm tra `if err != nil` có thể xuất hiện dày đặc đến mức nhấn chìm phần còn lại của đoạn mã.
Điều này thường xảy ra trong những chương trình thực hiện nhiều lời gọi API, nơi việc xử lý lỗi
còn đơn giản và chúng chỉ được trả thẳng ra ngoài.
Một số chương trình rốt cuộc có dạng như sau:

```Go
func printSum(a, b string) error {
	x, err := strconv.Atoi(a)
	if err != nil {
		return err
	}
	y, err := strconv.Atoi(b)
	if err != nil {
		return err
	}
	fmt.Println("result:", x + y)
	return nil
}
```

Trong mười dòng mã của thân hàm này, chỉ có bốn dòng (các lời gọi và hai dòng cuối) có vẻ như làm việc thật sự.
Sáu dòng còn lại giống như nhiễu.
Sự dài dòng là có thật, nên chẳng có gì ngạc nhiên khi những lời than phiền về xử lý lỗi
đứng đầu khảo sát người dùng thường niên của chúng tôi suốt nhiều năm.
(Trong một thời gian, việc thiếu generics đã vượt lên trên lời phàn nàn về xử lý lỗi, nhưng bây giờ
Go đã hỗ trợ generics, xử lý lỗi lại quay về vị trí số một.)

Nhóm Go coi trọng phản hồi từ cộng đồng, và vì thế trong nhiều năm qua chúng tôi đã cố gắng
tìm ra lời giải cho bài toán này cùng với ý kiến đóng góp từ cộng đồng Go.

Nỗ lực rõ ràng đầu tiên của nhóm Go bắt đầu từ năm 2018, khi Russ Cox
[mô tả bài toán một cách chính thức](https://go.googlesource.com/proposal/+/master/design/go2draft-error-handling-overview.md)
như một phần của cái mà khi đó chúng tôi gọi là nỗ lực Go 2.
Ông phác ra một lời giải khả dĩ dựa trên một
[bản thảo thiết kế](https://go.googlesource.com/proposal/+/master/design/go2draft-error-handling.md)
của Marcel van Lohuizen.
Thiết kế đó dựa trên cơ chế `check` và `handle` và khá toàn diện.
Bản thảo có một phân tích chi tiết về các giải pháp thay thế, bao gồm cả so sánh với
các cách tiếp cận ở những ngôn ngữ khác.
Nếu bạn thắc mắc liệu ý tưởng xử lý lỗi cụ thể của mình đã từng được xem xét hay chưa,
hãy đọc tài liệu này!

```Go
// printSum implementation using the proposed check/handle mechanism.
func printSum(a, b string) error {
	handle err { return err }
	x := check strconv.Atoi(a)
	y := check strconv.Atoi(b)
	fmt.Println("result:", x + y)
	return nil
}
```

Cách tiếp cận `check` và `handle` bị đánh giá là quá phức tạp và gần một năm sau, vào năm 2019,
chúng tôi tiếp nối bằng
[đề xuất `try`](https://go.googlesource.com/proposal/+/master/design/32437-try-builtin.md)
đã trở nên [khá tai tiếng](/issue/32437#issuecomment-2278932700).
Nó dựa trên ý tưởng của `check` và `handle`, nhưng pseudo-keyword `check` trở thành
hàm dựng sẵn `try` và phần `handle` bị bỏ đi.
Để khám phá tác động của `try`, chúng tôi đã viết một công cụ đơn giản
([tryhard](https://github.com/griesemer/tryhard))
để chuyển mã xử lý lỗi hiện có sang dùng `try`.
Đề xuất này đã bị tranh luận rất dữ dội, với gần 900 bình luận trên [issue GitHub](/issue/32437).

```Go
// printSum implementation using the proposed try mechanism.
func printSum(a, b string) error {
	// use a defer statement to augment errors before returning
	x := try(strconv.Atoi(a))
	y := try(strconv.Atoi(b))
	fmt.Println("result:", x + y)
	return nil
}
```

Tuy nhiên, `try` tác động đến luồng điều khiển bằng cách trả về khỏi hàm bao quanh khi có lỗi,
và còn làm điều đó từ những biểu thức có thể lồng rất sâu, nên che giấu luồng điều khiển này khỏi tầm mắt.
Điều này khiến đề xuất trở nên khó chấp nhận với nhiều người, và dù đã đầu tư đáng kể,
chúng tôi cũng quyết định từ bỏ nỗ lực này.
Nhìn lại, có lẽ tốt hơn nếu giới thiệu một từ khóa mới,
điều mà giờ đây chúng ta có thể làm vì đã có khả năng kiểm soát chi tiết phiên bản ngôn ngữ
qua tệp `go.mod` và các chỉ thị riêng theo tệp.
Việc giới hạn cách dùng `try` cho các phép gán và câu lệnh có thể đã xoa dịu một số
quan ngại khác. [Một đề xuất gần đây](/issue/73376) của Jimmy Frasche, về cơ bản
quay trở lại thiết kế `check` và `handle` ban đầu và xử lý một số thiếu sót của nó,
đang đi theo hướng đó.

Hệ quả của đề xuất `try` đã dẫn đến rất nhiều suy ngẫm, bao gồm một loạt bài blog
của Russ Cox: ["Thinking about the Go Proposal Process"](https://research.swtch.com/proposals-intro).
Một kết luận là có lẽ chúng tôi đã tự làm giảm cơ hội có được một kết quả tốt hơn khi đưa ra
một đề xuất gần như đã hoàn thiện, ít không gian cho phản hồi cộng đồng, và một mốc thời gian
triển khai gây cảm giác đe dọa. Theo ["Go Proposal Process: Large Changes"](https://research.swtch.com/proposals-large):
“nhìn lại, `try` là một thay đổi đủ lớn để bản thiết kế mới mà chúng tôi công bố [...] đáng lẽ phải
là bản thảo thiết kế thứ hai chứ không phải một đề xuất kèm mốc triển khai”.
Nhưng bất kể trong trường hợp này có lỗi về quy trình hay truyền thông hay không, cảm nhận của người dùng
đối với đề xuất này là cực kỳ không đồng tình.

Khi đó chúng tôi chưa có lời giải tốt hơn và đã không theo đuổi thay đổi cú pháp cho xử lý lỗi trong vài năm.
Dù vậy, rất nhiều người trong cộng đồng đã được truyền cảm hứng, và chúng tôi liên tục nhận được
những đề xuất xử lý lỗi mới, nhiều đề xuất rất giống nhau, có cái thú vị, có cái khó hiểu,
và có cái không khả thi.
Để theo dõi bức tranh đang mở rộng này, một năm sau đó Ian Lance Taylor đã tạo ra một
[umbrella issue](/issue/40432)
tóm tắt tình trạng hiện tại của các thay đổi được đề xuất nhằm cải thiện xử lý lỗi.
Một [Go Wiki](/wiki/Go2ErrorHandlingFeedback) cũng được tạo để thu thập phản hồi, thảo luận và bài viết liên quan.
Độc lập với chúng tôi, những người khác cũng đã bắt đầu theo dõi tất cả các đề xuất xử lý lỗi
qua nhiều năm.
Thật đáng kinh ngạc khi nhìn thấy số lượng khổng lồ của chúng, chẳng hạn trong bài viết của Sean K. H. Liao
về ["go error handling proposals"](https://seankhliao.com/blog/12020-11-23-go-error-handling-proposals/).

Những lời than phiền về sự dài dòng của xử lý lỗi vẫn tiếp diễn
(xem [Go Developer Survey 2024 H1 Results](/blog/survey2024-h1-results)),
vì vậy, sau một loạt đề xuất nội bộ ngày càng được tinh chỉnh của nhóm Go, Ian Lance Taylor đã công bố
["reduce error handling boilerplate using `?`"](/issue/71203) vào năm 2024.
Lần này ý tưởng là vay mượn từ một cấu trúc đã được hiện thực trong
[Rust](https://www.rust-lang.org/), cụ thể là
[toán tử `?`](https://doc.rust-lang.org/std/result/index.html#the-question-mark-operator-).
Hy vọng là bằng cách dựa vào một cơ chế đã tồn tại với một ký hiệu đã được công nhận, và tính đến
những gì chúng tôi đã học được qua các năm, cuối cùng chúng tôi sẽ tạo ra được tiến triển thật sự.
Trong các nghiên cứu người dùng nhỏ, không chính thức, khi lập trình viên được xem mã Go dùng `?`, đa số áp đảo
đã đoán đúng ý nghĩa của đoạn mã, điều càng thuyết phục chúng tôi cho nó thêm một
cơ hội.
Để nhìn rõ tác động của thay đổi này, Ian viết một công cụ chuyển mã Go thông thường
sang mã dùng cú pháp mới được đề xuất, và chúng tôi cũng làm một nguyên mẫu của tính năng đó trong
trình biên dịch.

```Go
// printSum implementation using the proposed "?" statements.
func printSum(a, b string) error {
	x := strconv.Atoi(a) ?
	y := strconv.Atoi(b) ?
	fmt.Println("result:", x + y)
	return nil
}
```

Thật không may, cũng như các ý tưởng xử lý lỗi khác, đề xuất mới này cũng nhanh chóng bị lấn át
bởi các bình luận và nhiều đề xuất tinh chỉnh nhỏ, thường dựa trên sở thích cá nhân.
Ian đã đóng đề xuất và chuyển nội dung sang một [cuộc thảo luận](/issue/71460)
để tạo điều kiện cho trao đổi và thu thập thêm phản hồi.
Một phiên bản có điều chỉnh nhẹ được đón nhận
[tích cực hơn đôi chút](https://github.com/golang/go/discussions/71460#discussioncomment-12060294)
nhưng sự ủng hộ rộng rãi vẫn rất khó đạt được.

Sau ngần ấy năm thử nghiệm, với ba đề xuất đầy đủ từ nhóm Go và
thực sự là [hàng trăm](/issues?q=+is%3Aissue+label%3Aerror-handling) (!)
đề xuất từ cộng đồng, đa số chỉ là các biến thể của cùng một chủ đề,
tất cả đều không thu hút được sự ủng hộ đủ mạnh, huống hồ là áp đảo,
câu hỏi giờ đây là: nên đi tiếp thế nào? Liệu có nên đi tiếp không?

_Chúng tôi nghĩ là không._

Nói chính xác hơn, chúng ta nên ngừng cố giải bài toán _cú pháp_, ít nhất là trong tương lai gần.
[Quy trình đề xuất](https://github.com/golang/proposal?tab=readme-ov-file#consensus-and-disagreement)
cung cấp cơ sở cho quyết định này:

> Mục tiêu của quy trình đề xuất là đạt được đồng thuận chung về kết quả trong một khoảng thời gian hợp lý.
> Nếu việc rà soát đề xuất không xác định được một đồng thuận chung trong phần thảo luận của issue trên issue tracker,
> kết quả thông thường là đề xuất bị từ chối.

Ngoài ra:

> Có thể xảy ra trường hợp quá trình rà soát đề xuất không xác định được đồng thuận chung nhưng lại thấy rõ rằng
> đề xuất không nên bị bác bỏ hoàn toàn.
> [...]
> Nếu nhóm rà soát đề xuất không xác định được đồng thuận hay bước tiếp theo cho đề xuất,
> quyết định về hướng đi tiếp sẽ được chuyển cho các Go architects [...], những người sẽ xem xét thảo luận và
> cố gắng đạt đồng thuận giữa chính họ.

Không có đề xuất xử lý lỗi nào đạt đến mức gần được đồng thuận,
nên tất cả đều bị từ chối.
Ngay cả những thành viên kỳ cựu nhất của nhóm Go tại Google cũng không hoàn toàn nhất trí
về con đường tốt nhất _ở thời điểm hiện tại_ (có thể điều đó sẽ thay đổi vào lúc nào đó).
Nhưng nếu không có đồng thuận mạnh, chúng ta không thể hợp lý mà tiến về phía trước.

Có những lập luận hợp lệ ủng hộ nguyên trạng:

- Nếu Go đã giới thiệu cú pháp rút gọn riêng cho xử lý lỗi từ sớm, ngày nay hẳn ít ai còn tranh luận về nó.
Nhưng chúng ta đã đi được 15 năm, cơ hội đó đã qua, và Go vẫn có
một cách hoàn toàn ổn để xử lý lỗi, dù đôi lúc có thể hơi dài dòng.

- Nhìn từ góc độ khác, giả sử hôm nay chúng ta tình cờ gặp được lời giải hoàn hảo.
Việc đưa nó vào ngôn ngữ chỉ đơn giản chuyển từ một nhóm người dùng không hài lòng
(nhóm muốn thay đổi) sang một nhóm khác (nhóm thích nguyên trạng).
Chúng tôi đã ở trong tình huống tương tự khi quyết định thêm generics vào ngôn ngữ, dù có
một khác biệt quan trọng:
ngày nay không ai bị buộc phải dùng generics, và các thư viện generic tốt được viết sao cho người dùng
hầu như có thể bỏ qua việc chúng là generic nhờ type inference.
Ngược lại, nếu một cấu trúc cú pháp mới cho xử lý lỗi được thêm vào ngôn ngữ,
hầu như ai cũng sẽ phải bắt đầu dùng nó, nếu không mã của họ sẽ trở nên không còn đúng phong cách.

- Việc không thêm cú pháp mới phù hợp với một trong các quy tắc thiết kế của Go:
đừng cung cấp nhiều cách để làm cùng một việc.
Có những ngoại lệ cho quy tắc này ở các khu vực có “lưu lượng cao”: phép gán là một ví dụ.
Trớ trêu thay, khả năng _khai báo lại_ biến trong
[short variable declarations](/ref/spec#Short_variable_declarations) (`:=`) được đưa vào để giải quyết
một vấn đề phát sinh từ xử lý lỗi:
nếu không có redeclaration, các chuỗi kiểm tra lỗi sẽ cần một biến `err` với tên khác nhau cho
mỗi lần kiểm tra (hoặc cần thêm các khai báo biến riêng biệt).
Vào thời điểm đó, một lời giải tốt hơn có thể đã là cung cấp thêm hỗ trợ cú pháp cho xử lý lỗi.
Khi ấy, quy tắc redeclaration có lẽ không còn cần thiết, và cùng với nó, nhiều [rắc rối liên quan](/issue/377)
cũng đã không tồn tại.

- Quay lại với mã xử lý lỗi thực tế, sự dài dòng mờ dần vào nền nếu lỗi thực sự được
_xử lý_.
Xử lý lỗi tốt thường đòi hỏi bổ sung thêm thông tin cho lỗi.
Ví dụ, một bình luận thường lặp lại trong khảo sát người dùng là việc thiếu stack trace gắn với lỗi.
Điều này có thể được giải quyết bằng các hàm hỗ trợ sinh ra và trả về lỗi đã được tăng cường.
Trong ví dụ sau (dĩ nhiên có phần gượng ép), tỷ lệ boilerplate tương đối nhỏ hơn nhiều:

	```Go
	func printSum(a, b string) error {
		x, err := strconv.Atoi(a)
		if err != nil {
			return fmt.Errorf("invalid integer: %q", a)
		}
		y, err := strconv.Atoi(b)
		if err != nil {
			return fmt.Errorf("invalid integer: %q", b)
		}
		fmt.Println("result:", x + y)
		return nil
	}
	```

- Chức năng mới trong thư viện chuẩn cũng có thể giúp giảm boilerplate xử lý lỗi,
rất đúng tinh thần bài viết năm 2015 của Rob Pike
["Errors are values"](/blog/errors-are-values).
Ví dụ, trong một số trường hợp [`cmp.Or`](/pkg/cmp#Or) có thể được dùng để xử lý
một chuỗi lỗi cùng lúc:

	```Go
	func printSum(a, b string) error {
		x, err1 := strconv.Atoi(a)
		y, err2 := strconv.Atoi(b)
		if err := cmp.Or(err1, err2); err != nil {
			return err
		}
		fmt.Println("result:", x+y)
		return nil
	}
	```

- Viết mã, đọc mã và gỡ lỗi mã là những hoạt động rất khác nhau.
Việc viết các đoạn kiểm tra lỗi lặp đi lặp lại có thể khá tẻ nhạt, nhưng IDE ngày nay cung cấp
khả năng hoàn thành mã mạnh mẽ, thậm chí có hỗ trợ LLM.
Việc viết các kiểm tra lỗi cơ bản là chuyện đơn giản với những công cụ này.
Sự dài dòng thường dễ thấy nhất khi đọc mã, nhưng công cụ cũng có thể giúp ở đây;
ví dụ một IDE với chế độ dành cho Go có thể cung cấp công tắc để ẩn mã xử lý lỗi.
Những công tắc như vậy đã tồn tại cho các phần khác của mã, chẳng hạn thân hàm.

- Khi gỡ lỗi mã xử lý lỗi, việc có thể nhanh chóng thêm một `println` hoặc
có một dòng hay vị trí mã nguồn riêng để đặt breakpoint trong debugger là điều rất hữu ích.
Điều đó rất dễ khi đã có sẵn một câu lệnh `if`.
Nhưng nếu toàn bộ logic xử lý lỗi bị che sau `check`, `try` hoặc `?`, mã có thể phải
được đổi ngược về câu lệnh `if` thông thường trước, điều đó làm việc gỡ lỗi phức tạp hơn
và thậm chí có thể đưa vào những lỗi tinh vi.

- Cũng còn những cân nhắc thực tế khác:
Nghĩ ra một ý tưởng cú pháp mới cho xử lý lỗi là việc rẻ;
vì thế mới có sự bùng nổ của vô số đề xuất từ cộng đồng.
Còn nghĩ ra một lời giải tốt, chịu được sự soi xét kỹ lưỡng, thì không hề rẻ.
Việc thiết kế đúng đắn một thay đổi ngôn ngữ và hiện thực nó thật sự
đòi hỏi một nỗ lực phối hợp đáng kể.
Chi phí thực sự còn đến sau đó:
tất cả đoạn mã cần thay đổi, tài liệu cần cập nhật,
các công cụ cần điều chỉnh.
Tính mọi thứ lại, thay đổi ngôn ngữ là cực kỳ đắt đỏ, nhóm Go lại tương đối nhỏ,
và còn rất nhiều ưu tiên khác cần được giải quyết.
(Những điểm này sau này có thể thay đổi: ưu tiên có thể dịch chuyển, quy mô nhóm có thể tăng hoặc giảm.)

- Cuối cùng, gần đây một số người trong chúng tôi có cơ hội tham dự
[Google Cloud Next 2025](https://cloud.withgoogle.com/next/25),
nơi nhóm Go có gian hàng và cũng tổ chức một buổi Go Meetup nhỏ.
Mọi người dùng Go mà chúng tôi kịp hỏi đều khẳng định rằng chúng tôi không nên thay đổi
ngôn ngữ chỉ để xử lý lỗi tốt hơn.
Nhiều người nhắc rằng việc Go thiếu hỗ trợ xử lý lỗi chuyên biệt là điều rõ rệt nhất
khi vừa chuyển sang từ một ngôn ngữ có hỗ trợ đó.
Khi người ta trở nên thành thạo hơn và viết mã Go đúng phong cách hơn, vấn đề này trở nên ít quan trọng hơn nhiều.
Dĩ nhiên đây không phải một tập người đủ lớn để đại diện,
nhưng có thể lại là một nhóm khác so với những người xuất hiện trên GitHub, và phản hồi của họ lại thêm một điểm dữ liệu nữa.

Dĩ nhiên cũng có những lập luận hợp lệ ủng hộ thay đổi:

- Việc thiếu hỗ trợ xử lý lỗi tốt hơn vẫn là lời phàn nàn số một trong khảo sát người dùng.
Nếu nhóm Go thật sự nghiêm túc với phản hồi của người dùng, sớm muộn gì chúng tôi cũng nên làm gì đó.
(Mặc dù dường như cũng không có
[sự ủng hộ áp đảo](https://github.com/golang/go/discussions/71460#discussioncomment-11977299)
cho một thay đổi ngôn ngữ.)

- Có lẽ việc chỉ chăm chăm giảm số ký tự là sai hướng.
Một cách tiếp cận tốt hơn có thể là làm cho việc xử lý lỗi mặc định trở nên rất dễ nhìn với một từ khóa
nhưng vẫn bỏ được boilerplate (`err != nil`).
Cách này có thể giúp người đọc (một code reviewer!) dễ thấy rằng một lỗi
đang được xử lý mà không phải “nhìn lại lần hai”, từ đó cải thiện chất lượng và độ an toàn của mã.
Điều này sẽ đưa ta quay về điểm khởi đầu của `check` và `handle`.

- Chúng tôi thực ra không biết bao nhiêu phần của vấn đề đến từ sự dài dòng cú pháp đơn thuần của
việc kiểm tra lỗi, so với sự dài dòng của xử lý lỗi tốt:
xây dựng những lỗi hữu ích như một phần của API, có ý nghĩa với lập trình viên và
người dùng cuối.
Đây là điều mà chúng tôi muốn nghiên cứu sâu hơn.

Dù vậy, chưa nỗ lực nào nhằm cải thiện xử lý lỗi cho tới nay giành được đủ lực kéo.
Nếu thành thật nhìn nhận vị trí hiện tại, chúng ta chỉ có thể thừa nhận rằng
chúng ta không những chưa có một cách hiểu chung về bài toán,
mà còn không đồng ý rằng liệu có thật sự có bài toán ở đây hay không.
Với điều này trong đầu, chúng tôi đưa ra quyết định thực dụng sau:

_Trong tương lai gần, nhóm Go sẽ ngừng theo đuổi các thay đổi cú pháp của ngôn ngữ
cho xử lý lỗi.
Chúng tôi cũng sẽ đóng mọi đề xuất đang mở và mới gửi liên quan chủ yếu
đến cú pháp xử lý lỗi, mà không điều tra thêm._

Cộng đồng đã đổ rất nhiều công sức vào việc khám phá, thảo luận và tranh luận các vấn đề này.
Dù điều đó có thể chưa dẫn tới thay đổi cú pháp xử lý lỗi nào, những nỗ lực ấy
đã dẫn tới nhiều cải tiến khác cho ngôn ngữ Go và quy trình của chúng tôi.
Có thể vào một thời điểm nào đó trong tương lai, bức tranh về xử lý lỗi sẽ trở nên rõ ràng hơn.
Cho đến lúc đó, chúng tôi mong muốn tập trung niềm đam mê đáng kinh ngạc này vào những cơ hội mới
để làm cho Go tốt hơn cho mọi người.

Xin cảm ơn!
