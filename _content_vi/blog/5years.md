---
title: Nửa thập kỷ cùng Go
date: 2014-11-10
by:
- Andrew Gerrand
summary: Chúc mừng sinh nhật lần thứ năm của Go!
template: true
---


Năm năm trước, chúng tôi ra mắt dự án Go. Cảm giác như mới chỉ hôm qua
khi chúng tôi còn chuẩn bị cho bản phát hành công khai đầu tiên: [website](https://web.archive.org/web/20091112094121/http://golang.org/) của chúng tôi
có một màu vàng rất đáng yêu, chúng tôi gọi Go là một “ngôn ngữ hệ thống”,
và bạn phải kết thúc câu lệnh bằng dấu chấm phẩy cũng như viết Makefile để build mã của mình.
Khi đó, chúng tôi hoàn toàn không biết Go sẽ được đón nhận ra sao.
Liệu mọi người có chia sẻ tầm nhìn và mục tiêu của chúng tôi không? Liệu họ có thấy Go hữu ích không?

Khi ra mắt, Go nhận được một làn sóng chú ý mạnh mẽ. Google vừa tạo ra một
ngôn ngữ lập trình mới, và ai cũng muốn thử xem nó là gì.
Một số lập trình viên không hứng thú với tập tính năng có phần bảo thủ của Go,
thoạt nhìn họ thấy “chẳng có gì đáng xem ở đây”,
nhưng một nhóm nhỏ hơn đã nhìn thấy mầm mống của một hệ sinh thái
được may đo cho nhu cầu của họ với tư cách là những kỹ sư phần mềm đang làm việc thực tế.
Chính nhóm người này đã tạo nên hạt nhân của cộng đồng Go.

{{image "5years/gophers5th.jpg" 850}}

[_Hình minh họa Gopher_](/blog/gopher) _của_ [_Renee French_](http://reneefrench.blogspot.com.au/)

Sau bản phát hành đầu tiên, chúng tôi mất một thời gian để truyền đạt đúng
mục tiêu và tinh thần thiết kế phía sau Go.
Rob Pike đã làm điều đó rất xuất sắc trong bài luận năm 2012
[_Go at Google: Language Design in the Service of Software Engineering_](/talks/2012/splash.article)
và theo cách gần gũi hơn trong bài blog
[_Less is exponentially more_](https://commandcenter.blogspot.com.au/2012/06/less-is-exponentially-more.html).
Bài nói chuyện của Andrew Gerrand
[_Code that grows with grace_](http://vimeo.com/53221560)
([slides](/talks/2012/chat.slide)) và
[_Go for Gophers_](https://www.youtube.com/watch?v=dKGmK_Z1Zl0)
([slides](/talks/2014/go4gophers.slide)) đưa ra góc nhìn kỹ thuật,
đi sâu hơn vào triết lý thiết kế của Go.

Theo thời gian, từ số ít đã thành số đông.
Bước ngoặt của dự án là việc phát hành Go 1 vào tháng 3 năm 2012,
phiên bản mang lại một ngôn ngữ và thư viện chuẩn ổn định mà các lập trình viên có thể tin cậy.
Đến năm 2014, dự án đã có hàng trăm người đóng góp cốt lõi,
hệ sinh thái có vô số [thư viện và công cụ](https://godoc.org/)
được duy trì bởi hàng nghìn nhà phát triển,
và cộng đồng rộng lớn hơn có rất nhiều thành viên nhiệt huyết
(hay như chúng tôi gọi họ là “gophers”).
Theo các chỉ số hiện tại của chúng tôi, cộng đồng Go đang phát triển nhanh hơn nhiều
so với điều mà chúng tôi từng nghĩ là có thể.

Vậy những gopher đó có thể được tìm thấy ở đâu?
Họ có mặt tại rất nhiều sự kiện Go đang mọc lên trên khắp thế giới.
Năm nay, chúng tôi chứng kiến một số hội nghị Go chuyên biệt:
[GopherCon](/blog/gophercon) đầu tiên và
[dotGo](http://www.dotgo.eu/) tại Denver và Paris,
[Go DevRoom tại FOSDEM](/blog/fosdem14) và thêm hai kỳ nữa của hội nghị [GoCon](https://github.com/GoCon/GoCon) tổ chức hai năm một lần
tại Tokyo.
Ở mỗi sự kiện, các gopher từ khắp nơi trên thế giới đã hào hứng trình bày các dự án Go của mình.
Đối với nhóm Go, việc được gặp nhiều lập trình viên cùng chia sẻ tầm nhìn và sự phấn khích với chúng tôi
là một điều vô cùng mãn nguyện.

{{image "5years/conferences.jpg"}}

_Hơn 1.200 gopher đã tham dự GopherCon ở Denver và dotGo ở Paris._

Ngoài ra còn có hàng chục [Go User Group](/wiki/GoUserGroups) do cộng đồng tự vận hành
ở nhiều thành phố trên khắp thế giới.
Nếu bạn chưa từng đến nhóm gần mình, hãy cân nhắc tham gia.
Và nếu nơi bạn ở chưa có nhóm nào, có lẽ bạn nên [bắt đầu một nhóm](/blog/getthee-to-go-meetup)?

Ngày nay, Go đã tìm được vị trí của mình trong đám mây.
Go xuất hiện đúng vào lúc ngành công nghiệp đang trải qua một sự dịch chuyển mang tính địa chấn
về phía điện toán đám mây, và chúng tôi rất vui khi thấy nó nhanh chóng trở thành
một phần quan trọng của làn sóng đó.
Sự đơn giản, hiệu quả, các primitive đồng thời tích hợp,
và thư viện chuẩn hiện đại khiến Go rất phù hợp cho phát triển phần mềm đám mây
(suy cho cùng, đó chính là điều nó được thiết kế để phục vụ).
Những dự án đám mây mã nguồn mở quan trọng như
[Docker](https://www.docker.com/) và
[Kubernetes](https://github.com/GoogleCloudPlatform/kubernetes)
đã được viết bằng Go, và những công ty hạ tầng như Google, CloudFlare, Canonical,
Digital Ocean, GitHub, Heroku và Microsoft hiện đều đang dùng Go
để gánh vác những phần việc nặng.

Vậy tương lai sẽ như thế nào? Chúng tôi tin rằng năm 2015 sẽ là năm lớn nhất của Go từ trước tới nay.

Go 1.4, ngoài [những tính năng và bản sửa lỗi mới](/doc/go1.4),
đã đặt nền móng cho một bộ gom rác độ trễ thấp mới
và hỗ trợ chạy Go trên thiết bị di động.
Phiên bản này dự kiến phát hành vào ngày 1 tháng 12 năm 2014.
Chúng tôi kỳ vọng GC mới sẽ có mặt trong Go 1.5, dự kiến vào ngày 1 tháng 6 năm 2015,
qua đó giúp Go trở nên hấp dẫn hơn với một phạm vi ứng dụng rộng hơn.
Chúng tôi rất nóng lòng muốn thấy mọi người sẽ đưa Go đi đến đâu.

Và sẽ còn nhiều sự kiện tuyệt vời khác nữa, với [GothamGo](http://gothamgo.com/) tại
New York (15/11), thêm một Go DevRoom tại FOSDEM ở Brussels (31/1 và 1/2;
[hãy tham gia!](https://groups.google.com/d/msg/golang-nuts/1xgBazQzs1I/hwrZ5ni8cTEJ)),
[GopherCon India](http://www.gophercon.in/) tại Bengaluru (19-21/2),
[GopherCon](http://gophercon.com/) quay lại Denver vào tháng 7,
và [dotGo](http://www.dotgo.eu/) một lần nữa ở Paris vào tháng 11.

Nhóm Go muốn gửi lời cảm ơn tới tất cả các gopher ngoài kia.
Xin chúc cho năm năm tiếp theo.

_Để kỷ niệm 5 năm của Go, trong suốt tháng tới,_
[_Gopher Academy_](http://blog.gopheracademy.com/)
_sẽ đăng một loạt bài viết của những người dùng Go nổi bật. Nhớ ghé xem_
[_blog của họ_](http://blog.gopheracademy.com/)
_để theo dõi thêm những nội dung thú vị về Go._
