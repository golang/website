---
title: Bước tiếp theo cho Generics
date: 2020-06-16
by:
- Ian Lance Taylor
- Robert Griesemer
tags:
- go2
- proposals
- generics
summary: Một bản thảo thiết kế generics đã được cập nhật, cùng một công cụ chuyển đổi để thử nghiệm
---

## Giới thiệu

Đã gần một năm kể từ lần [chúng tôi viết gần nhất về khả năng
thêm generics vào Go](/blog/why-generics).
Đã đến lúc cập nhật.

## Thiết kế cập nhật

Chúng tôi tiếp tục tinh chỉnh [bản thảo thiết kế
generics](https://go.googlesource.com/proposal/+/refs/heads/master/design/go2draft-contracts.md).
Chúng tôi đã viết một trình kiểm tra kiểu cho nó: một chương trình có thể parse mã Go
dùng generics như được mô tả trong bản thảo thiết kế và báo mọi lỗi kiểu.
Chúng tôi đã viết mã ví dụ.
Và chúng tôi đã thu thập phản hồi từ rất, rất nhiều người, cảm ơn
vì điều đó!

Dựa trên những gì đã học được, chúng tôi công bố một [bản thảo thiết kế
cập nhật](https://go.googlesource.com/proposal/+/refs/heads/master/design/go2draft-type-parameters.md).
Thay đổi lớn nhất là chúng tôi bỏ ý tưởng về contracts.
Sự khác biệt giữa contracts và interface types gây bối rối, nên
chúng tôi loại bỏ sự khác biệt đó.
Type parameters giờ được ràng buộc bởi interface types.
Interface types giờ được phép bao gồm type list, dù chỉ
khi dùng làm ràng buộc; trong bản thảo thiết kế trước đó, type list là
một tính năng của contracts.
Những trường hợp phức tạp hơn sẽ dùng parameterized interface type.

Chúng tôi hy vọng mọi người sẽ thấy bản thảo thiết kế này đơn giản hơn và dễ
hiểu hơn.

## Công cụ thử nghiệm

Để giúp quyết định nên tinh chỉnh bản thảo thiết kế дальше như thế nào, chúng tôi
công bố một công cụ chuyển đổi.
Đây là một công cụ cho phép mọi người kiểm tra kiểu và chạy mã được viết
theo phiên bản generics mô tả trong bản thảo thiết kế.
Nó hoạt động bằng cách chuyển mã generic thành mã Go thông thường.
Quá trình chuyển đổi này áp đặt một số giới hạn, nhưng chúng tôi hy vọng rằng
nó đủ tốt để mọi người cảm nhận được mã Go generic
có thể trông như thế nào.
Phần hiện thực generics thật, nếu chúng được chấp nhận vào
ngôn ngữ, sẽ hoạt động theo cách khác.
(Chúng tôi mới chỉ bắt đầu phác thảo xem một
phần hiện thực trực tiếp trong trình biên dịch sẽ trông ra sao.)

Công cụ có sẵn trên một biến thể của Go playground tại
[https://go2goplay.golang.org](https://go2goplay.golang.org).
Playground này hoạt động giống hệt Go playground thông thường, nhưng
hỗ trợ mã generic.

Bạn cũng có thể tự xây và dùng công cụ.
Nó có sẵn trong một nhánh của repo Go master.
Hãy làm theo [hướng dẫn cài đặt Go từ
mã nguồn](/doc/install/source).
Ở chỗ các hướng dẫn đó bảo bạn checkout tag bản phát hành mới nhất,
hãy chạy `git checkout dev.go2go` thay vào đó.
Sau đó build Go toolchain như hướng dẫn.

Công cụ chuyển đổi được tài liệu hóa trong
[README.go2go](https://go.googlesource.com/go/+/refs/heads/dev.go2go/README.go2go.md).

## Các bước tiếp theo

Chúng tôi hy vọng công cụ này sẽ cho cộng đồng Go cơ hội
thử nghiệm với generics.
Có hai điều chính mà chúng tôi hy vọng sẽ học được.

Thứ nhất, mã generic có hợp lý không?
Nó có mang cảm giác Go không?
Mọi người gặp phải những điều bất ngờ nào?
Thông báo lỗi có hữu ích không?

Thứ hai, chúng tôi biết rằng nhiều người đã nói Go cần generics, nhưng
không nhất thiết biết chính xác điều đó có nghĩa là gì.
Bản thiết kế này có giải quyết vấn đề theo cách hữu ích không?
Nếu có một vấn đề khiến bạn nghĩ “Tôi có thể giải quyết chuyện này nếu Go
có generics,” liệu bạn có giải quyết được nó khi dùng công cụ này không?

Chúng tôi sẽ dùng phản hồi thu thập được từ cộng đồng Go để quyết định nên
tiến lên như thế nào.
Nếu bản thảo thiết kế được đón nhận tốt và không cần những thay đổi
đáng kể, bước tiếp theo sẽ là một [đề xuất thay đổi ngôn ngữ
chính thức](/s/proposal).
Để đặt kỳ vọng, nếu mọi người hoàn toàn hài lòng với bản thảo thiết kế
và nó không cần điều chỉnh thêm, thời điểm sớm nhất generics có thể được thêm
vào Go sẽ là bản Go 1.17,
dự kiến vào tháng 8 năm 2021.
Dĩ nhiên trên thực tế có thể xuất hiện những vấn đề không lường trước được, nên đây là
một mốc thời gian lạc quan; chúng tôi không thể đưa ra dự đoán chắc chắn.

## Phản hồi

Cách tốt nhất để đưa phản hồi cho các thay đổi ngôn ngữ sẽ là qua
danh sách thư `golang-nuts@googlegroups.com`.
Danh sách thư không hoàn hảo, nhưng có vẻ đây là lựa chọn tốt nhất cho
thảo luận ban đầu.
Khi viết về bản thảo thiết kế, xin hãy đặt `[generics]` ở
đầu dòng Subject và bắt đầu các luồng khác nhau cho các chủ đề cụ thể khác nhau.

Nếu bạn phát hiện lỗi trong trình kiểm tra kiểu generics hoặc công cụ chuyển đổi,
hãy tạo issue trong Go issue tracker tiêu chuẩn tại
[go.dev/issue](/issue).
Xin hãy bắt đầu tiêu đề issue bằng `cmd/go2go:`.
Lưu ý rằng issue tracker không phải nơi tốt nhất để thảo luận thay đổi
của ngôn ngữ, vì nó không cung cấp luồng trao đổi và cũng không
phù hợp với các cuộc trò chuyện dài.

Chúng tôi mong nhận được phản hồi của bạn.

## Lời cảm ơn

Chúng tôi chưa xong, nhưng đã đi được một chặng đường dài.
Chúng tôi sẽ không thể tới đây nếu không có rất nhiều sự giúp đỡ.

Chúng tôi muốn cảm ơn Philip Wadler và các cộng sự vì đã suy nghĩ
một cách hình thức về generics trong Go và giúp chúng tôi làm rõ các khía cạnh lý thuyết
của thiết kế.
Bài báo của họ [Featherweight Go](https://arxiv.org/abs/2005.11710)
phân tích generics trong một phiên bản Go bị giới hạn, và họ cũng đã
phát triển một nguyên mẫu [trên GitHub](https://github.com/rhu1/fgg).

Chúng tôi cũng muốn cảm ơn [những
người](https://go.googlesource.com/proposal/+/refs/heads/master/design/go2draft-type-parameters.md#acknowledgements)
đã cung cấp phản hồi chi tiết cho một phiên bản trước đó của bản thảo thiết kế.

Và cuối cùng nhưng chắc chắn không kém phần quan trọng, chúng tôi muốn cảm ơn rất nhiều người trong
nhóm Go, nhiều người đóng góp cho Go issue tracker, và tất cả những ai khác
đã chia sẻ ý tưởng và phản hồi cho các bản thảo thiết kế trước.
Chúng tôi đã đọc tất cả, và rất biết ơn. Chúng tôi sẽ không thể có mặt ở đây nếu không có
mọi người.
