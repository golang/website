---
title: Bảy năm của Go
date: 2016-11-10
by:
- Nhóm Go
summary: Chúc mừng sinh nhật lần thứ bảy của Go!
---


<img src="7years/gopherbelly300.jpg" align="right">

Hôm nay đánh dấu tròn bảy năm kể từ khi chúng tôi phát hành mã nguồn mở bản phác thảo ban đầu của Go.
Với sự giúp đỡ của cộng đồng mã nguồn mở, bao gồm hơn một nghìn
người đóng góp cá nhân cho các kho mã nguồn của Go,
Go đã trưởng thành thành một ngôn ngữ được sử dụng trên khắp thế giới.

Những thay đổi đáng kể nhất mà người dùng có thể thấy trong Go trong năm qua là
việc bổ sung hỗ trợ tích hợp sẵn cho
[HTTP/2](https://www.youtube.com/watch?v=FARQMJndUn0#t=0m0s) trong
[Go 1.6](/doc/go1.6) và việc tích hợp
[package context](/blog/context) vào thư viện chuẩn trong [Go 1.7](/doc/go1.7).
Nhưng chúng tôi còn thực hiện rất nhiều cải tiến ít thấy hơn.
Go 1.7 đã thay đổi trình biên dịch x86-64 để dùng back end mới dựa trên SSA,
qua đó cải thiện hiệu năng của hầu hết chương trình Go thêm 10-20%.
Đối với Go 1.8, dự kiến phát hành vào tháng 2 năm sau,
chúng tôi cũng đã chuyển các trình biên dịch của những kiến trúc khác sang dùng back end mới đó.
Chúng tôi còn bổ sung các bản port mới, tới Android trên x86 32 bit, Linux trên MIPS 64 bit,
và Linux trên IBM z Systems.
Chúng tôi cũng phát triển các kỹ thuật garbage collection mới giúp giảm
thời gian dừng “stop the world” điển hình xuống [dưới 100 micro giây](/design/17503-eliminate-rescan).
(Hãy so điều đó với tin lớn của Go 1.5 là [10 mili giây trở xuống](/blog/go15gc).)

Năm nay mở đầu bằng một cuộc hackathon Go toàn cầu,
[Gopher Gala](/blog/gophergala), vào tháng 1.
Sau đó là các [hội nghị Go](/wiki/Conferences) ở Ấn Độ và Dubai vào tháng 2,
Trung Quốc và Nhật Bản vào tháng 4, San Francisco vào tháng 5, Denver vào tháng 7,
London vào tháng 8, Paris vào tháng trước, và Brazil vào cuối tuần vừa rồi.
Và GothamGo ở New York sẽ diễn ra vào tuần tới.
Năm nay cũng chứng kiến hơn 30 [nhóm người dùng Go](/wiki/GoUserGroups) mới,
tám chapter mới của [Women Who Go](http://www.womenwhogo.org/),
và bốn workshop [GoBridge](https://golangbridge.org/) trên khắp thế giới.

Chúng tôi tiếp tục vừa choáng ngợp vừa biết ơn
trước sự nhiệt tình và ủng hộ của cộng đồng Go.
Dù bạn tham gia bằng cách gửi thay đổi, báo lỗi,
chia sẻ chuyên môn của mình trong các cuộc thảo luận thiết kế, viết blog hay sách,
tổ chức meetup, giúp người khác học hoặc tiến bộ hơn,
mã nguồn mở các package Go bạn viết, hay đơn giản chỉ là trở thành một phần của cộng đồng Go,
nhóm Go xin cảm ơn sự giúp đỡ, thời gian và năng lượng của bạn.
Go sẽ không thể trở thành thành công như ngày hôm nay nếu thiếu bạn.

Xin cảm ơn, và chúc thêm một năm nữa thật vui và thành công cùng Go!
