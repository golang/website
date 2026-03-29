---
title: Chín năm của Go
date: 2018-11-10
by:
- Steve Francia
tags:
- community
- birthday
summary: Chúc mừng sinh nhật lần thứ 9 của Go!
template: true
---

## Giới thiệu

Hôm nay đánh dấu kỷ niệm chín năm kể từ ngày chúng tôi mã nguồn mở hóa bản phác thảo ban đầu của Go.
Mỗi dịp kỷ niệm, chúng tôi đều muốn dành thời gian để nhìn lại những gì đã diễn ra trong năm qua.
12 tháng vừa qua là một năm bứt phá đối với ngôn ngữ Go và cộng đồng của nó.

## Tình yêu dành cho Go và mức độ đón nhận

Nhờ tất cả các bạn, năm 2018 là một năm tuyệt vời với Go!
Trong nhiều khảo sát trong ngành, các gopher bày tỏ rằng họ hạnh phúc thế nào khi dùng Go,
và nhiều nhà phát triển không dùng Go cho biết họ dự định học Go trước bất kỳ ngôn ngữ nào khác.

Trong [Khảo sát Nhà phát triển 2018 của Stack Overflow](https://insights.stackoverflow.com/survey/2018#most-loved-dreaded-and-wanted),
Go vẫn giữ vị trí đáng mơ ước trong cả top 5 ngôn ngữ được yêu thích nhất và top 5 ngôn ngữ được muốn học nhất.
Người đang dùng Go thì yêu nó, còn người chưa dùng Go thì muốn dùng nó.

Trong [Khảo sát Nhà phát triển 2018 của ActiveState](https://www.activestate.com/developer-survey-2018-open-source-runtime-pains),
Go đứng đầu bảng với 36% người dùng trả lời rằng họ “Cực kỳ hài lòng” khi dùng Go
và 61% trả lời “Rất hài lòng” hoặc cao hơn.

[Khảo sát Nhà phát triển 2018 của JetBrains](https://www.jetbrains.com/research/devecosystem-2018/) đã trao cho Go
danh hiệu “Ngôn ngữ hứa hẹn nhất” với 12% người trả lời đang dùng Go và 16% có ý định dùng Go trong tương lai.

Trong [Khảo sát Nhà phát triển 2018 của HackerRank](https://research.hackerrank.com/developer-skills/2018/),
38% nhà phát triển trả lời rằng họ có ý định học Go tiếp theo.

Chúng tôi rất hào hứng với tất cả những gopher mới và đang tích cực cải thiện các tài nguyên giáo dục cũng như cộng đồng của mình.

## Cộng đồng Go

Thật khó tin rằng mới chỉ năm năm kể từ
những hội nghị Go và meetup Go đầu tiên.
Chúng tôi đã thấy sự tăng trưởng mạnh trong lĩnh vực lãnh đạo cộng đồng này trong năm qua.
Giờ đây đã có hơn 20 [hội nghị Go](/wiki/Conferences)
và hơn 300 [meetup liên quan tới Go](https://www.meetup.com/topics/golang/) trải rộng khắp toàn cầu.

Nhờ công sức bỏ ra cho rất nhiều hội nghị và meetup này,
đã có hàng trăm bài nói chuyện tuyệt vời trong năm nay.
Dưới đây là một vài bài nói chuyện chúng tôi yêu thích, đặc biệt thảo luận về sự phát triển của cộng đồng
và cách chúng ta có thể hỗ trợ các gopher trên toàn thế giới tốt hơn.

  - [Writing Accessible Go](https://www.youtube.com/watch?v=cVaDY0ChvOQ), của Julia Ferraioli tại GopherCon
  - [The Importance of Beginners](https://www.youtube.com/watch?v=7yMXs9TRvVI), của Natalie Pistunovich tại GopherCon
  - [The Legacy of Go, Part 2](https://www.youtube.com/watch?v=I_KcpgxcFyU), của Carmen Andoh tại GothamGo
  - [Growing a Community of Gophers](https://www.youtube.com/watch?v=dl1mCGKwlYY), của Cassandra Salisbury tại Gopherpalooza

Theo tinh thần đó, năm nay chúng tôi cũng đã [sửa đổi quy tắc ứng xử của mình](/blog/conduct-2018)
để hỗ trợ tốt hơn cho tính bao trùm trong cộng đồng Go.

Cộng đồng Go thực sự có tính toàn cầu.
Tại GopherCon Europe ở Iceland vào mùa hè vừa qua, các gopher đã theo đúng nghĩa đen đứng trên khoảng cách giữa hai lục địa.

{{image "9years/9years-iceland.jpg" 800}}

_(Ảnh của Winter Francia.)_

## Go 2

Sau năm năm kinh nghiệm với Go 1, chúng tôi đã bắt đầu xem xét
những gì cần thay đổi ở Go để hỗ trợ tốt hơn cho
[lập trình ở quy mô lớn](/talks/2012/splash.article).

Mùa xuân năm ngoái, chúng tôi đã công bố [bản thiết kế dự thảo cho Go modules](/blog/versioning-proposal),
cung cấp một cơ chế tích hợp cho việc quản lý phiên bản và phân phối package.
Bản phát hành Go gần đây nhất, Go 1.11, đã bao gồm
[hỗ trợ sơ bộ cho modules](/doc/go1.11#modules).

Mùa hè năm ngoái, chúng tôi đã công bố
[các bản thiết kế dự thảo ban đầu](/blog/go2draft)
cho cách Go 2 có thể hỗ trợ tốt hơn cho error values, xử lý lỗi và lập trình generic.

Chúng tôi rất hào hứng khi tiếp tục tinh chỉnh các thiết kế này cùng với sự trợ giúp từ cộng đồng
trong hành trình [hướng tới Go 2](/blog/toward-go2).

## Những người đóng góp cho Go

Dự án Go đã gia tăng số lượng đóng góp từ cộng đồng trong nhiều năm qua.
Dự án đạt một cột mốc lớn vào giữa năm 2018 khi, lần đầu tiên,
số đóng góp đến từ cộng đồng vượt quá số đóng góp từ nhóm Go.

{{image "9years/9years-graph.png" 600}}

## Xin cảm ơn

Từ góc độ cá nhân, thay mặt toàn bộ nhóm Go,
chúng tôi muốn chân thành cảm ơn tất cả các bạn.
Chúng tôi cảm thấy đặc biệt may mắn khi được làm việc trên dự án Go
và biết ơn rất nhiều gopher trên khắp thế giới đã đồng hành cùng chúng tôi.

Chúng tôi đặc biệt biết ơn hàng nghìn tình nguyện viên
đang giúp sức thông qua cố vấn, tổ chức, đóng góp
và hỗ trợ các gopher đồng hành cùng mình.
Chính các bạn đã làm nên Go như ngày hôm nay.
