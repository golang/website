---
title: Mười một năm của Go
date: 2020-11-10T12:01:00Z
by:
- Russ Cox, thay mặt nhóm Go
summary: Chúc mừng sinh nhật, Go!
---


Hôm nay chúng ta kỷ niệm sinh nhật lần thứ mười một của bản phát hành mã nguồn mở Go.
Những buổi tiệc mà chúng ta từng có để
[mừng Go tròn 10 tuổi](/blog/10years)
giờ dường như đã là một ký ức xa xôi.
Đây là một năm đầy khó khăn, nhưng
chúng tôi vẫn giữ cho quá trình phát triển Go tiếp tục tiến lên
và tích lũy được khá nhiều dấu mốc nổi bật.

Vào tháng 11, chúng tôi ra mắt [go.dev và pkg.go.dev](/blog/go.dev)
không lâu sau sinh nhật lần thứ 10 của Go.

Vào tháng 2, [bản phát hành Go 1.14](/blog/go1.14)
đã mang đến bản triển khai Go modules đầu tiên chính thức được xem là “sẵn sàng cho môi trường production”,
cùng với nhiều cải tiến về hiệu năng,
bao gồm
[defer nhanh hơn](/design/34481-opencoded-defers)
và
[goroutine preemption không hợp tác](/design/24543/conservative-inner-frame)
để giảm độ trễ của việc lập lịch
và garbage collection.

Vào đầu tháng 3, chúng tôi ra mắt
[API mới cho protocol buffers](/blog/protobuf-apiv2),
[google.golang.org/protobuf](https://pkg.go.dev/google.golang.org/protobuf),
với khả năng hỗ trợ được cải thiện đáng kể cho reflection của protocol buffer và các custom message.

<img src="11years/gophermask.jpg" height="450" width="300" align="right" style="border: 2px solid black; margin: 0 0 1em 1em;">

Khi đại dịch bùng phát, chúng tôi quyết định tạm dừng mọi thông báo
hoặc đợt ra mắt công khai trong mùa xuân,
vì nhận thấy rằng sự chú ý của mọi người khi đó nên dành cho những việc khác quan trọng hơn.
Nhưng chúng tôi vẫn tiếp tục làm việc, và một thành viên trong nhóm đã tham gia vào
sự hợp tác giữa Apple và Google về
[thông báo phơi nhiễm bảo toàn quyền riêng tư](https://www.google.com/covid19/exposurenotifications/)
để hỗ trợ các nỗ lực truy vết tiếp xúc trên toàn thế giới.
Vào tháng 5, nhóm đó đã ra mắt
[máy chủ backend tham chiếu](https://github.com/google/exposure-notifications-server),
được viết bằng Go.

Chúng tôi tiếp tục cải thiện [gopls](https://www.youtube.com/watch?v=EFJfdWzBHwE),
công cụ cung cấp
[khả năng hỗ trợ hiểu Go](https://github.com/golang/tools/blob/master/gopls/doc/user.md)
ở nhiều trình soạn thảo.
Vào tháng 6,
[tiện ích mở rộng Go cho VSCode chính thức gia nhập dự án Go](/blog/vscode-go)
và hiện được duy trì bởi chính những nhà phát triển đang làm việc trên gopls.

Cũng trong tháng 6, nhờ phản hồi của các bạn, chúng tôi đã công khai mã nguồn của
[phần đứng sau pkg.go.dev](/blog/pkgsite)
như một phần của dự án Go.

Cuối tháng 6, chúng tôi
[phát hành bản nháp thiết kế mới nhất cho generics](/blog/generics-next-step),
kèm theo một công cụ nguyên mẫu và [sân chơi generics](https://go2goplay.golang.org/).

Vào tháng 7, chúng tôi công bố và thảo luận ba bản nháp thiết kế mới cho các thay đổi trong tương lai:
[dòng `//go:build` mới để chọn tệp](/design/draft-gobuild),
[giao diện hệ thống tệp](/design/draft-iofs),
và
[nhúng tệp tại thời điểm build](/design/draft-embed).
(Như ghi chú bên dưới, chúng ta sẽ thấy tất cả những thay đổi đó trong năm 2021.)

Vào tháng 8, [bản phát hành Go 1.15](/blog/go1.15)
chủ yếu mang đến các tối ưu hóa và sửa lỗi thay vì tính năng mới.
Điểm đáng kể nhất là việc bắt đầu viết lại linker,
giúp nó chạy nhanh hơn 20% và dùng ít hơn 30% bộ nhớ
trung bình đối với các bản build lớn.

Tháng trước, chúng tôi thực hiện [khảo sát người dùng Go thường niên](/blog/survey2020).
Chúng tôi sẽ đăng kết quả lên blog sau khi phân tích xong.

Cộng đồng Go cũng đã thích nghi với cách tiếp cận “ưu tiên trực tuyến” giống như mọi người,
và chúng tôi đã chứng kiến nhiều buổi meetup trực tuyến cùng hơn một chục hội nghị Go trực tuyến trong năm nay.
Tuần trước, nhóm Go đã tổ chức
[Ngày Go tại Google Open Source Live](https://opensourcelive.withgoogle.com/events/go)
(video có tại liên kết).

## Hướng về phía trước

Chúng tôi cũng vô cùng hào hứng với những gì đang chờ Go trong năm thứ 12.
Trước mắt nhất, trong tuần này các thành viên nhóm Go sẽ
tham gia trình bày tám sự kiện tại
[GopherCon 2020](https://www.gophercon.com/).
Hãy ghi lại lịch nhé!

- “Typing [Generic] Go”,
  bài nói chuyện của Robert Griesemer,\
  [Ngày 11 tháng 11, 10:00 sáng (miền Đông Hoa Kỳ)](https://www.gophercon.com/agenda/session/233094);
  [Hỏi đáp lúc 10:30 sáng](https://www.gophercon.com/agenda/session/417935).
- “What to Expect When You’re NOT Expecting”,
  một buổi ghi hình trực tiếp podcast Go time với hội đồng các chuyên gia gỡ lỗi,
  trong đó có Hana Kim,\
  [Ngày 11 tháng 11, 12:00 trưa](https://www.gophercon.com/agenda/session/2334490).
- “Evolving the Go Memory Manager's RAM and CPU Efficiency”,
  bài nói chuyện của Michael Knyszek,\
  [Ngày 11 tháng 11, 1:00 chiều](https://www.gophercon.com/agenda/session/233086);
  [Hỏi đáp lúc 1:50 chiều](https://www.gophercon.com/agenda/session/417940).
- “Implementing Faster Defers”,
  bài nói chuyện của Dan Scales,\
  [Ngày 11 tháng 11, 5:10 chiều](https://www.gophercon.com/agenda/session/233397);
  [Hỏi đáp lúc 5:40 chiều](https://www.gophercon.com/agenda/session/417941).
- “Go Team - Ask Me Anything”,
  buổi hỏi đáp trực tiếp với Julie Qiu, Rebecca Stambler, Russ Cox, Sameer Ajmani và Van Riper,\
  [Ngày 12 tháng 11, 3:00 chiều](https://www.gophercon.com/agenda/session/420539).
- “Pardon the Interruption: Loop Preemption in Go 1.14”,
  bài nói chuyện của Austin Clements,\
  [Ngày 12 tháng 11, 4:45 chiều](https://www.gophercon.com/agenda/session/233441);
  [Hỏi đáp lúc 5:15 chiều](https://www.gophercon.com/agenda/session/417943).
- “Working with Errors”,
  bài nói chuyện của Jonathan Amsterdam,\
  [Ngày 13 tháng 11, 1:00 chiều](https://www.gophercon.com/agenda/session/233432);
  [Hỏi đáp lúc 1:50 chiều](https://www.gophercon.com/agenda/session/417945).
- “Crossing the Chasm for Go: Two Million Users and Growing”,
  bài nói chuyện của Carmen Andoh,\
  [Ngày 13 tháng 11, 5:55 chiều](https://www.gophercon.com/agenda/session/233426).

## Các bản phát hành Go

Vào tháng 2, bản phát hành Go 1.16 sẽ bao gồm
[giao diện hệ thống tệp](https://tip.golang.org/pkg/io/fs/)
và
[nhúng tệp tại thời điểm build](https://tip.golang.org/pkg/embed/).
Nó sẽ hoàn tất việc viết lại linker, mang đến thêm các cải tiến về hiệu năng.
Và nó sẽ bao gồm hỗ trợ cho các máy Mac Apple Silicon (`GOARCH=arm64`) mới.

Vào tháng 8, bản phát hành Go 1.17 chắc chắn sẽ mang đến thêm nhiều tính năng và cải tiến,
mặc dù hiện tại còn khá xa nên các chi tiết cụ thể vẫn chưa hoàn toàn rõ ràng.
Nó sẽ bao gồm quy ước gọi hàm mới dựa trên thanh ghi cho x86-64
(mà không làm hỏng mã assembly hiện có!),
giúp chương trình chạy nhanh hơn trên diện rộng.
(Các kiến trúc khác sẽ theo sau trong những bản phát hành sau.)
Một tính năng hay chắc chắn sẽ được đưa vào là
[dòng `//go:build` mới](/design/draft-gobuild),
ít gây lỗi hơn nhiều so với
[dòng `// +build` hiện tại](/cmd/go/#hdr-Build_constraints).
Một tính năng rất được mong đợi khác mà chúng tôi hy vọng sẽ sẵn sàng để thử nghiệm beta trong năm tới
là
[hỗ trợ fuzzing trong lệnh `go test`](/design/draft-fuzzing).

## Go Modules

Trong năm tới, chúng tôi sẽ tiếp tục phát triển hỗ trợ cho Go modules
và tích hợp chúng thật tốt vào toàn bộ hệ sinh thái Go.
Go 1.16 sẽ mang đến trải nghiệm Go modules mượt mà nhất từ trước đến nay.
Một kết quả sơ bộ từ cuộc khảo sát gần đây của chúng tôi là 96% người dùng
đã chuyển sang sử dụng Go modules (tăng từ 90% của một năm trước).

Cuối cùng, chúng tôi cũng sẽ dần khép lại hỗ trợ cho cách phát triển dựa trên GOPATH:
mọi chương trình sử dụng dependency ngoài thư viện chuẩn đều sẽ cần `go.mod`.
(Nếu bạn vẫn chưa chuyển sang modules, hãy xem
[trang wiki về GOPATH](/wiki/GOPATH)
để biết chi tiết về bước cuối cùng trong hành trình từ GOPATH sang modules.)

Ngay từ đầu, [mục tiêu của Go modules](https://research.swtch.com/vgo-intro)
là “đưa khái niệm phiên bản package vào vốn từ làm việc
của cả lập trình viên Go lẫn các công cụ của chúng ta,”
nhằm tạo điều kiện cho việc hỗ trợ modules và phiên bản ở mức sâu trong toàn bộ hệ sinh thái Go.
[Mirror, cơ sở dữ liệu checksum và chỉ mục của Go module](/blog/modules2019)
được tạo ra nhờ cách hiểu trên phạm vi toàn hệ sinh thái về việc một phiên bản package là gì.
Trong năm tới, chúng ta sẽ thấy hỗ trợ module phong phú hơn được thêm vào nhiều công cụ và hệ thống.
Ví dụ, chúng tôi dự định nghiên cứu các công cụ mới để giúp tác giả module phát hành phiên bản mới
(`go release`)
cũng như giúp người dùng module cập nhật mã của họ để rời xa
những API đã không còn được khuyến nghị dùng nữa (một `go fix` mới).

Lấy một ví dụ lớn hơn,
[chúng tôi đã tạo ra gopls](https://github.com/golang/tools/blob/master/gopls/README.md)
để gom nhiều công cụ mà các trình soạn thảo dùng cho hỗ trợ Go,
trong đó không công cụ nào hỗ trợ modules,
thành một công cụ duy nhất có hỗ trợ.
Trong năm tới,
chúng tôi sẽ sẵn sàng để tiện ích mở rộng Go cho VSCode dùng `gopls` theo mặc định,
nhằm mang lại trải nghiệm module xuất sắc ngay từ đầu,
và chúng tôi sẽ phát hành gopls 1.0.
Tất nhiên, một trong những điều tuyệt vời nhất về gopls là nó trung lập với trình soạn thảo:
bất kỳ trình soạn thảo nào hiểu
[language server protocol](https://langserver.org/)
đều có thể sử dụng nó.

Một ứng dụng quan trọng khác của thông tin phiên bản là theo dõi xem
bất kỳ package nào trong một bản build có lỗ hổng bảo mật đã biết hay không.
Trong năm tới, chúng tôi dự định phát triển một cơ sở dữ liệu về các lỗ hổng đã biết
cũng như các công cụ để kiểm tra chương trình của bạn đối chiếu với cơ sở dữ liệu đó.

Trang khám phá package Go
[pkg.go.dev](https://pkg.go.dev/)
là một ví dụ khác của hệ thống nhận biết phiên bản được Go modules tạo điều kiện.
Chúng tôi đã tập trung vào việc hoàn thiện chức năng cốt lõi và trải nghiệm người dùng,
bao gồm cả
[bản thiết kế lại ra mắt hôm nay](/blog/pkgsite-redesign).
Trong năm tới,
chúng tôi sẽ hợp nhất godoc.org vào pkg.go.dev.
Chúng tôi cũng sẽ mở rộng dòng thời gian phiên bản cho mỗi package,
hiển thị các thay đổi quan trọng trong từng phiên bản,
các lỗ hổng đã biết và nhiều thông tin khác,
theo đúng mục tiêu tổng thể là làm nổi bật những gì bạn cần để đưa ra
[quyết định sáng suốt khi thêm dependency](https://research.swtch.com/deps).

Chúng tôi rất hào hứng khi thấy hành trình từ GOPATH đến Go modules
đang tiến gần đến hồi kết và tất cả những công cụ nhận biết dependency tuyệt vời
mà Go modules đang tạo điều kiện.

## Generics

Tính năng tiếp theo mà mọi người đều nghĩ tới dĩ nhiên là generics.
Như đã nhắc ở trên, vào tháng 6 chúng tôi đã công bố
[bản nháp thiết kế mới nhất cho generics](/blog/generics-next-step).
Kể từ đó, chúng tôi tiếp tục tinh chỉnh những chỗ còn thô
và chuyển sự chú ý sang các chi tiết để hiện thực một phiên bản sẵn sàng cho môi trường production.
Chúng tôi sẽ làm việc đó trong suốt năm 2021, với mục tiêu có
một thứ gì đó để mọi người có thể thử vào cuối năm,
có lẽ là một phần của các bản beta Go 1.18.

## Xin cảm ơn!

Go còn hơn nhiều so với chỉ riêng chúng tôi, nhóm Go tại Google.
Chúng tôi mang ơn những người đóng góp đang cùng làm việc với chúng tôi trong các bản phát hành và công cụ của Go.
Xa hơn nữa, Go chỉ có thể thành công nhờ tất cả các bạn, những người đang làm việc trong
và đóng góp cho hệ sinh thái Go đầy sức sống.
Đây là một năm khó khăn với thế giới bên ngoài Go.
Hơn bao giờ hết, chúng tôi trân trọng việc các bạn dành thời gian
đồng hành cùng chúng tôi và giúp Go thành công đến vậy.
Xin cảm ơn.
Chúng tôi hy vọng tất cả các bạn luôn an toàn và nhận được những điều tốt đẹp nhất.
