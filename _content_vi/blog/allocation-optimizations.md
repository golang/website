---
title: Cấp phát trên stack
date: 2026-02-27
by:
- Keith Randall
summary: Mô tả một số thay đổi gần đây nhằm thực hiện cấp phát trên stack thay vì heap.
template: true
---

Chúng tôi luôn tìm cách để các chương trình Go chạy nhanh hơn. Trong
2 bản phát hành gần đây, chúng tôi đã tập trung vào việc giảm thiểu một nguồn gây
chậm cụ thể, đó là cấp phát trên heap. Mỗi khi một chương trình Go cấp phát bộ nhớ
từ heap, sẽ có một đoạn mã khá lớn cần chạy
để đáp ứng việc cấp phát đó. Ngoài ra, cấp phát trên heap còn tạo thêm
tải cho bộ gom rác. Ngay cả với các cải tiến
gần đây như [Green Tea](/blog/greenteagc), bộ gom rác
vẫn phát sinh chi phí đáng kể.

Vì vậy, chúng tôi đã nghiên cứu những cách để thực hiện nhiều lần cấp phát hơn trên stack
thay vì heap. Cấp phát trên stack rẻ hơn đáng kể
khi thực hiện (đôi khi gần như miễn phí). Hơn nữa, chúng không tạo tải
cho bộ gom rác, vì các cấp phát trên stack có thể được thu hồi
tự động cùng với chính stack frame. Cấp phát trên stack
cũng cho phép tái sử dụng nhanh chóng, điều này rất thân thiện với bộ nhớ đệm.

## Cấp phát trên stack cho slice có kích thước hằng

Hãy xem xét tác vụ xây dựng một slice các task để xử lý:
{{raw `
	func process(c chan task) {
		var tasks []task
		for t := range c {
			tasks = append(tasks, t)
		}
		processAll(tasks)
	}
`}}

Hãy cùng xem điều gì xảy ra lúc chạy khi lấy task từ
channel `c` và thêm chúng vào slice `tasks`.

Ở vòng lặp đầu tiên, `tasks` chưa có vùng lưu trữ nền, vì vậy
`append` phải cấp phát một vùng. Vì nó không biết slice cuối cùng
sẽ lớn đến mức nào, nên nó không thể quá mạnh tay. Hiện tại, nó
cấp phát một vùng lưu trữ nền có kích thước 1.

Ở vòng lặp thứ hai, vùng lưu trữ nền đã tồn tại, nhưng nó đã
đầy. `append` lại phải cấp phát một vùng lưu trữ nền mới, lần này có
kích thước 2. Vùng lưu trữ nền cũ có kích thước 1 giờ trở thành rác.

Ở vòng lặp thứ ba, vùng lưu trữ nền có kích thước 2 đã
đầy. `append` *lại tiếp tục* phải cấp phát một vùng lưu trữ nền mới, lần này
có kích thước 4. Vùng lưu trữ nền cũ có kích thước 2 giờ trở thành rác.

Ở vòng lặp thứ tư, vùng lưu trữ nền có kích thước 4 mới chỉ có 3
phần tử. `append` chỉ cần đặt phần tử mới vào vùng lưu trữ nền hiện có
và tăng độ dài của slice. Tuyệt! Không cần gọi bộ cấp phát cho
vòng lặp này.

Ở vòng lặp thứ năm, vùng lưu trữ nền có kích thước 4 đã đầy, và
`append` lại phải cấp phát một vùng lưu trữ nền mới, lần này có kích thước
8.

Và cứ tiếp tục như vậy. Chúng ta thường tăng gấp đôi kích thước cấp phát mỗi lần nó
đầy, để cuối cùng có thể thêm hầu hết task mới vào slice
mà không cần cấp phát. Nhưng có một lượng chi phí khá đáng kể ở
"giai đoạn khởi động" khi slice còn nhỏ. Trong giai đoạn khởi động này, chúng ta
tốn nhiều thời gian trong bộ cấp phát, và tạo ra một đống rác,
điều này có vẻ khá lãng phí. Và cũng có thể trong chương trình của bạn,
slice đó không bao giờ thật sự lớn. Giai đoạn khởi động này có thể là tất cả
những gì bạn từng gặp.

Nếu đoạn mã này là một phần cực kỳ nóng trong chương trình của bạn, bạn có thể
bị cám dỗ bắt đầu slice với kích thước lớn hơn để tránh tất cả các
lần cấp phát đó.

{{raw `
	func process2(c chan task) {
		tasks := make([]task, 0, 10) // có lẽ nhiều nhất là 10 task
		for t := range c {
			tasks = append(tasks, t)
		}
		processAll(tasks)
	}
`}}

Đây là một tối ưu hóa hợp lý. Nó không bao giờ sai; chương trình của
bạn vẫn chạy đúng. Nếu ước lượng quá nhỏ, bạn sẽ nhận các
lần cấp phát từ `append` như trước. Nếu ước lượng quá lớn, bạn sẽ
lãng phí một ít bộ nhớ.

Nếu ước lượng của bạn về số lượng task là tốt, thì chỉ còn
một vị trí cấp phát trong chương trình này. Lệnh gọi `make` cấp phát một
vùng lưu trữ nền cho slice với kích thước chính xác, và `append` sẽ không bao giờ phải
tái cấp phát.

Điều bất ngờ là nếu bạn benchmark đoạn mã này với 10
phần tử trong channel, bạn sẽ thấy mình không giảm số lần
cấp phát xuống 1, mà giảm xuống 0!

Lý do là compiler đã quyết định cấp phát vùng lưu trữ nền
trên stack. Vì nó biết cần kích thước bao nhiêu (10 lần
kích thước của một task) nên nó có thể cấp phát bộ nhớ cho nó trong stack frame của
`process2` thay vì trên heap[<sup>1</sup>](#footnotes). Lưu ý
rằng điều này phụ thuộc vào việc vùng lưu trữ nền không [thoát
ra heap](/doc/gc-guide#Escape_analysis) bên trong `processAll`.

## Cấp phát trên stack cho slice có kích thước biến thiên

Nhưng tất nhiên, việc gán cứng một ước lượng kích thước là khá cứng nhắc.
Có lẽ chúng ta có thể truyền vào một độ dài ước lượng?

{{raw `
	func process3(c chan task, lengthGuess int) {
		tasks := make([]task, 0, lengthGuess)
		for t := range c {
			tasks = append(tasks, t)
		}
		processAll(tasks)
	}
`}}

Điều này cho phép caller chọn kích thước tốt cho slice `tasks`, có thể
thay đổi tùy theo nơi đoạn mã này được gọi.

Thật không may, trong Go 1.24, kích thước không hằng của vùng lưu trữ nền
khiến compiler không còn có thể cấp phát vùng lưu trữ nền trên
stack. Nó sẽ nằm trên heap, biến đoạn mã không cấp phát của chúng ta
thành đoạn mã cấp phát 1 lần. Vẫn tốt hơn so với việc để `append` thực hiện tất cả
các lần cấp phát trung gian, nhưng vẫn đáng tiếc.

Nhưng đừng lo, Go 1.25 đã đến!

Hãy tưởng tượng bạn quyết định làm như sau, để chỉ nhận được cấp phát trên stack
trong những trường hợp ước lượng nhỏ:

{{raw `
	func process4(c chan task, lengthGuess int) {
		var tasks []task
		if lengthGuess <= 10 {
			tasks = make([]task, 0, 10)
		} else {
			tasks = make([]task, 0, lengthGuess)
		}
		for t := range c {
			tasks = append(tasks, t)
		}
		processAll(tasks)
	}
`}}

Khá xấu, nhưng nó sẽ hoạt động. Khi ước lượng nhỏ, bạn dùng
`make` với kích thước hằng và vì vậy có một vùng lưu trữ nền được cấp phát trên stack, còn
khi ước lượng lớn hơn bạn dùng `make` với kích thước biến thiên và cấp phát
vùng lưu trữ nền từ heap.

Nhưng trong Go 1.25, bạn không cần đi theo con đường xấu xí này. Compiler Go
1.25 sẽ tự thực hiện phép biến đổi này cho bạn! Với một số vị trí
cấp phát slice nhất định, compiler tự động cấp phát một vùng lưu trữ nền
nhỏ (hiện tại là 32 byte) trên stack cho slice, và dùng vùng lưu trữ đó
làm kết quả của `make` nếu kích thước được yêu cầu đủ
nhỏ. Nếu không, nó sẽ dùng cấp phát trên heap như bình thường.

Trong Go 1.25, `process3` thực hiện không cấp phát heap,
nếu `lengthGuess` đủ nhỏ để một slice có độ dài đó vừa trong 32
byte. (Và dĩ nhiên `lengthGuess` là một ước lượng đúng về số
phần tử trong `c`.)

Chúng tôi luôn cải thiện hiệu năng của Go, vì vậy hãy nâng cấp lên bản phát hành Go
mới nhất và [ngạc
nhiên](https://youtu.be/FUm0pfgWehI?si=QRTt_JYwr-cRHDNJ&t=960) trước việc
chương trình của bạn trở nên nhanh hơn và tiết kiệm bộ nhớ hơn bao nhiêu!

## Cấp phát trên stack cho slice do append cấp phát

Được rồi, nhưng bạn vẫn không muốn phải thay đổi API để thêm
ước lượng độ dài kỳ quặc này. Có cách nào khác không?

Nâng cấp lên Go 1.26!

{{raw `
	func process(c chan task) {
		var tasks []task
		for t := range c {
			tasks = append(tasks, t)
		}
		processAll(tasks)
	}
`}}

Trong Go 1.26, chúng tôi cấp phát cùng kiểu vùng lưu trữ nền nhỏ, có tính suy đoán
trên stack, nhưng giờ đây chúng tôi có thể dùng nó trực tiếp tại vị trí `append`.

Ở vòng lặp đầu tiên, chưa có vùng lưu trữ nền cho `tasks`, vì vậy
`append` dùng một vùng lưu trữ nền nhỏ được cấp phát trên stack làm
lần cấp phát đầu tiên. Ví dụ, nếu chúng ta có thể chứa 4 `task` trong vùng lưu trữ nền đó,
lần `append` đầu tiên sẽ cấp phát một vùng lưu trữ nền có độ dài 4 từ stack.

3 vòng lặp tiếp theo thêm trực tiếp vào vùng lưu trữ nền trên stack,
không yêu cầu cấp phát.

Ở vòng lặp thứ 4, vùng lưu trữ nền trên stack cuối cùng đã đầy và chúng ta
phải chuyển sang heap để lấy thêm vùng lưu trữ nền. Nhưng chúng ta đã tránh được
gần như toàn bộ chi phí khởi động được mô tả trước đó trong bài viết này.
Không có cấp phát heap với kích thước 1, 2 và 4, và cũng không có số rác
mà cuối cùng chúng trở thành. Nếu slice của bạn nhỏ, có lẽ bạn sẽ không bao giờ
phải cấp phát heap.

## Cấp phát trên stack cho slice do append cấp phát nhưng thoát ra ngoài

Được rồi, tất cả điều này đều ổn khi slice `tasks` không thoát ra ngoài. Nhưng nếu
tôi trả về slice thì sao? Khi đó nó không thể được cấp phát trên stack, đúng không?

Đúng vậy! Vùng lưu trữ nền cho slice được trả về bởi `extract` bên dưới
không thể được cấp phát trên stack, vì stack frame của `extract`
biến mất khi `extract` trả về.

{{raw `
	func extract(c chan task) []task {
		var tasks []task
		for t := range c {
			tasks = append(tasks, t)
		}
		return tasks
	}
`}}

Nhưng bạn có thể nghĩ rằng, slice *được trả về* không thể được cấp phát trên
stack. Nhưng còn tất cả những slice trung gian chỉ trở thành
rác thì sao? Có lẽ chúng ta có thể cấp phát chúng trên stack?

{{raw `
	func extract2(c chan task) []task {
		var tasks []task
		for t := range c {
			tasks = append(tasks, t)
		}
		tasks2 := make([]task, len(tasks))
		copy(tasks2, tasks)
		return tasks2
	}
`}}

Khi đó slice `tasks` không bao giờ thoát ra khỏi `extract2`. Nó có thể hưởng lợi từ
tất cả các tối ưu hóa được mô tả ở trên. Sau đó, ở cuối
`extract2`, khi chúng ta biết kích thước cuối cùng của slice, chúng ta thực hiện một lần cấp phát heap
với kích thước cần thiết, sao chép các `task` vào đó, và trả về
bản sao.

Nhưng bạn có thật sự muốn viết thêm toàn bộ đoạn mã đó không? Có vẻ
dễ sinh lỗi. Có lẽ compiler có thể thực hiện phép biến đổi này cho chúng ta?

Trong Go 1.26, có thể!

Với các slice thoát ra ngoài, compiler sẽ biến đổi mã `extract`
ban đầu thành một thứ như thế này:

{{raw `
	func extract3(c chan task) []task {
		var tasks []task
		for t := range c {
			tasks = append(tasks, t)
		}
		tasks = runtime.move2heap(tasks)
		return tasks
	}
`}}

`runtime.move2heap` là một hàm đặc biệt của compiler+runtime, nó là
hàm đồng nhất đối với các slice đã được cấp phát trên heap.
Đối với các slice đang ở trên stack, nó cấp phát một slice mới trên
heap, sao chép slice được cấp phát trên stack sang bản sao trên heap, rồi trả về
bản sao trên heap.

Điều này đảm bảo rằng với mã `extract` gốc của chúng ta, nếu số lượng
phần tử vừa trong bộ đệm nhỏ được cấp phát trên stack, chúng ta thực hiện đúng 1
lần cấp phát với chính xác kích thước phù hợp. Nếu số phần tử vượt quá
dung lượng của bộ đệm nhỏ được cấp phát trên stack, chúng ta sẽ thực hiện kiểu cấp phát gấp đôi
thông thường một khi bộ đệm trên stack bị tràn.

Tối ưu hóa mà Go 1.26 thực hiện thực ra còn tốt hơn đoạn mã
được tối ưu hóa thủ công, vì nó không yêu cầu thêm
lần cấp phát+sao chép mà đoạn mã tối ưu hóa thủ công luôn thực hiện ở cuối.
Nó chỉ yêu cầu cấp phát+sao chép trong trường hợp chúng ta đã chỉ
thao tác trên một slice có nền là stack cho đến thời điểm trả về.

Chúng ta vẫn phải trả chi phí cho một lần sao chép, nhưng chi phí đó gần như hoàn toàn
được bù lại bởi các lần sao chép trong giai đoạn khởi động mà chúng ta không còn phải
thực hiện nữa. (Thực tế, sơ đồ mới trong trường hợp xấu nhất chỉ phải sao chép nhiều hơn sơ đồ cũ
một phần tử.)

## Kết luận

Tối ưu hóa thủ công vẫn có thể hữu ích, đặc biệt nếu bạn có một
ước lượng tốt về kích thước của slice từ trước. Nhưng hy vọng là giờ đây compiler
sẽ tự bắt được nhiều trường hợp đơn giản cho bạn và cho phép bạn tập trung
vào những trường hợp còn lại thật sự quan trọng.

Có rất nhiều chi tiết mà compiler cần đảm bảo để thực hiện đúng
tất cả các tối ưu hóa này. Nếu bạn cho rằng một trong các tối ưu hóa đó
đang gây ra vấn đề về tính đúng đắn hoặc hiệu năng (theo hướng tiêu cực)
cho bạn, bạn có thể tắt chúng bằng
`-gcflags=all=-d=variablemakehash=n`. Nếu việc tắt các tối ưu hóa này
giúp ích, vui lòng [tạo issue](/issue/new) để chúng tôi điều tra.

## Chú thích

<sup>1</sup> Stack của Go không có bất kỳ cơ chế kiểu `alloca` nào cho
stack frame có kích thước động. Mọi stack frame của Go đều có kích thước
hằng.
