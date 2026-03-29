---
title: Go tròn 10 tuổi
date: 2019-11-08
by:
- Russ Cox, thay mặt nhóm Go
summary: Chúc mừng sinh nhật 10 tuổi, Go!
template: true
---


Chúc mừng sinh nhật, Go!

Cuối tuần này, chúng ta kỷ niệm 10 năm kể từ
[bản phát hành Go](https://opensource.googleblog.com/2009/11/hey-ho-lets-go.html),
đánh dấu cột mốc Go tròn 10 tuổi với tư cách là một ngôn ngữ lập trình mã nguồn mở
và một hệ sinh thái để xây dựng phần mềm mạng hiện đại.

Nhân dịp này,
[Renee French](https://twitter.com/reneefrench),
người tạo ra
[Go gopher](/blog/gopher),
đã vẽ bức tranh đầy thú vị này:

<a href="10years/gopher10th-large.jpg">
{{image "10years/gopher10th-small.jpg" 850}}
</a>

Kỷ niệm 10 năm của Go khiến tôi nhớ lại đầu tháng 11 năm 2009,
khi chúng tôi đang chuẩn bị giới thiệu Go với thế giới.
Chúng tôi không biết sẽ nhận được phản ứng như thế nào,
liệu có ai quan tâm đến ngôn ngữ nhỏ bé này hay không.
Tôi hy vọng rằng ngay cả khi không ai sử dụng Go,
thì ít nhất chúng tôi cũng sẽ làm nổi bật một số ý tưởng hay,
đặc biệt là cách Go tiếp cận tính đồng thời và interface,
để có thể ảnh hưởng đến các ngôn ngữ ra đời sau đó.

Khi đã rõ ràng là mọi người thực sự hào hứng với Go,
tôi tìm hiểu lịch sử của những ngôn ngữ phổ biến
như C, C++, Perl, Python và Ruby,
để xem mất bao lâu mỗi ngôn ngữ mới đạt được sự phổ biến rộng rãi.
Ví dụ, trong suy nghĩ của tôi, Perl có vẻ như xuất hiện với hình hài hoàn chỉnh
vào giữa đến cuối những năm 1990, cùng các CGI script và web,
nhưng thực ra nó được phát hành lần đầu vào năm 1987.
Mẫu hình này lặp lại ở gần như mọi ngôn ngữ tôi xem qua:
có vẻ như cần xấp xỉ một thập kỷ cải tiến bền bỉ, âm thầm
và phổ biến rộng rãi trước khi một ngôn ngữ mới thật sự cất cánh.

Tôi tự hỏi: sau một thập kỷ, Go sẽ đứng ở đâu?

Hôm nay, chúng ta đã có câu trả lời cho câu hỏi đó:
Go có mặt ở khắp nơi, được ít nhất [một triệu lập trình viên trên toàn thế giới](https://research.swtch.com/gophercount) sử dụng.

Mục tiêu ban đầu của Go là hạ tầng hệ thống mạng,
điều mà ngày nay chúng ta gọi là phần mềm đám mây.
Ngày nay, mọi nhà cung cấp đám mây lớn đều sử dụng hạ tầng đám mây cốt lõi được viết bằng Go,
như Docker, Etcd, Istio, Kubernetes, Prometheus và Terraform;
phần lớn các
[dự án của Cloud Native Computing Foundation](https://www.cncf.io/projects/)
được viết bằng Go.
Vô số công ty cũng đang sử dụng Go để đưa công việc của họ lên đám mây,
từ các startup xây dựng mới từ đầu
đến các doanh nghiệp đang hiện đại hóa ngăn xếp phần mềm.
Go cũng đã được ứng dụng vượt xa mục tiêu đám mây ban đầu,
với các trường hợp sử dụng trải dài
từ việc điều khiển những hệ thống nhúng siêu nhỏ bằng
[GoBot](https://gobot.io) và [TinyGo](https://tinygo.org/)
đến việc phát hiện ung thư bằng
[phân tích dữ liệu lớn quy mô lớn và học máy tại GRAIL](https://medium.com/grail-eng/bigslice-a-cluster-computing-system-for-go-7e03acd2419b),
và mọi thứ ở giữa hai đầu đó.

Tất cả những điều này cho thấy Go đã thành công vượt xa mọi mong đợi táo bạo nhất của chúng tôi.
Và thành công của Go không chỉ nằm ở ngôn ngữ.
Nó nằm ở ngôn ngữ, hệ sinh thái, và đặc biệt là cộng đồng cùng nhau hợp tác.

Năm 2009, ngôn ngữ này là một ý tưởng hay cùng một bản phác thảo đang hoạt động.
Lệnh `go` khi đó chưa tồn tại:
chúng tôi chạy những lệnh như `6g` để biên dịch và `6l` để liên kết tệp nhị phân,
được tự động hóa bằng makefile.
Chúng tôi gõ dấu chấm phẩy ở cuối câu lệnh.
Toàn bộ chương trình dừng lại trong quá trình garbage collection,
và bộ gom rác khi đó cũng khó tận dụng tốt hai lõi xử lý.
Go chỉ chạy trên Linux và Mac, trên x86 32 bit và 64 bit và ARM 32 bit.

Trong suốt thập kỷ qua, với sự giúp đỡ của các lập trình viên Go trên khắp thế giới,
chúng tôi đã phát triển ý tưởng và bản phác thảo này thành một ngôn ngữ năng suất cao
với hệ thống công cụ tuyệt vời,
một trình triển khai đạt chất lượng sản xuất,
một
[bộ gom rác hiện đại hàng đầu](/blog/ismmkeynote),
và [các bản port tới 12 hệ điều hành và 10 kiến trúc](/doc/install/source#introduction).

Bất kỳ ngôn ngữ lập trình nào cũng cần sự hỗ trợ của một hệ sinh thái phát triển mạnh.
Bản phát hành mã nguồn mở là hạt giống cho hệ sinh thái đó,
nhưng kể từ đó, rất nhiều người đã đóng góp thời gian và tài năng của mình
để bồi đắp hệ sinh thái Go bằng những hướng dẫn, sách, khóa học, bài viết blog,
podcast, công cụ, tích hợp, và tất nhiên là các gói Go có thể tải bằng `go get`.
Go không thể nào thành công nếu không có sự hỗ trợ của hệ sinh thái này.

Tất nhiên, hệ sinh thái lại cần sự hỗ trợ của một cộng đồng phát triển mạnh.
Trong năm 2019 có hàng chục hội nghị Go trên khắp thế giới,
cùng với
[hơn 150 nhóm meetup Go với hơn 90.000 thành viên](https://www.meetup.com/pro/go).
[GoBridge](https://golangbridge.org)
và
[Women Who Go](https://medium.com/@carolynvs/www-loves-gobridge-ccb26309f667)
giúp đưa những tiếng nói mới vào cộng đồng Go,
thông qua cố vấn, đào tạo, và học bổng hội nghị.
Chỉ riêng năm nay, họ đã hướng dẫn
hàng trăm người thuộc những nhóm trước đây ít được đại diện
thông qua các buổi workshop, nơi thành viên cộng đồng giảng dạy và cố vấn cho những người mới đến với Go.

Hiện có
[hơn một triệu lập trình viên Go](https://research.swtch.com/gophercount)
trên toàn thế giới,
và các công ty trên khắp địa cầu đang tìm cách tuyển thêm.
Thực tế, mọi người thường nói với chúng tôi rằng việc học Go
đã giúp họ có được những công việc đầu tiên trong ngành công nghệ.
Sau cùng, điều chúng tôi tự hào nhất về Go
không phải là một tính năng được thiết kế tốt hay một đoạn mã thông minh
mà là tác động tích cực mà Go đã mang lại cho cuộc sống của rất nhiều người.
Chúng tôi hướng đến việc tạo ra một ngôn ngữ giúp chúng tôi trở thành những lập trình viên giỏi hơn,
và chúng tôi vô cùng phấn khởi khi Go đã giúp được rất nhiều người khác.

Khi
[\#GoTurns10](https://twitter.com/search?q=%23GoTurns10),
tôi hy vọng mọi người sẽ dành ít phút để cùng kỷ niệm
cộng đồng Go và tất cả những gì chúng ta đã đạt được.
Thay mặt toàn bộ nhóm Go tại Google,
xin cảm ơn tất cả mọi người đã đồng hành cùng chúng tôi trong suốt thập kỷ vừa qua.
Hãy biến thập kỷ tiếp theo trở nên tuyệt vời hơn nữa!

<div>
<center>
<a href="10years/gopher10th-pin-large.jpg">
{{image "10years/gopher10th-pin-small.jpg" 150}}
</center>
</div>
