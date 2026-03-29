---
title: Các bài nói về Go tại FOSDEM 2014
date: 2014-02-24
by:
- Andrew Gerrand
tags:
- fosdem
- youtube
- talk
summary: Tường thuật từ Go Devroom tại FOSDEM 2014.
template: true
---

## Giới thiệu

Tại [FOSDEM](http://fosdem.org/) vào ngày 2 tháng 2 năm 2014, các thành viên của cộng đồng Go
đã trình bày một loạt bài nói trong Go Devroom. Ngày hôm đó là một thành công lớn,
với 13 bài nói xuất sắc được trình bày trước một khán phòng luôn chật kín người.

Các bản ghi hình video của những bài nói này giờ đã có sẵn, và một số video tiêu biểu
được giới thiệu bên dưới.

Toàn bộ loạt bài nói có tại
[danh sách phát trên YouTube](http://www.youtube.com/playlist?list=PLtLJO5JKE5YDKG4WcaNts3IVZqhDmmuBH).
(Bạn cũng có thể xem trực tiếp tại
[kho video của FOSDEM](http://video.fosdem.org/2014/K4601/Sunday/).)

## Mở rộng quy mô với Go: Vitess của YouTube

Kỹ sư Google Sugu Sougoumarane đã mô tả cách anh và nhóm của mình xây dựng
[Vitess](https://github.com/youtube/vitess) bằng Go để giúp mở rộng quy mô
[YouTube](https://youtube.com).

Vitess là một tập hợp các máy chủ và công cụ chủ yếu được phát triển bằng Go.
Nó giúp mở rộng các cơ sở dữ liệu MySQL cho web, và hiện đang được dùng như
một thành phần nền tảng trong hạ tầng MySQL của YouTube.

Bài nói trình bày một số lịch sử về việc vì sao nhóm chọn Go, và cách lựa chọn đó
đã mang lại hiệu quả.
Sugu cũng nói về các mẹo và kỹ thuật dùng để mở rộng Vitess bằng Go.

{{video "https://www.youtube.com/embed/qATTTSg6zXk"}}

Slide của bài nói [có tại đây](https://github.com/youtube/vitess/blob/master/doc/Vitess2014.pdf?raw=true).

## Camlistore

[Camlistore](http://camlistore.org/) được thiết kế để trở thành "hệ thống lưu trữ cá nhân cho cả đời,
đặt bạn vào vị trí kiểm soát và được thiết kế để tồn tại lâu dài." Nó là mã nguồn mở,
đã có gần 4 năm phát triển tích cực, và cực kỳ linh hoạt. Trong bài nói này,
Brad Fitzpatrick và Mathieu Lonjaret giải thích vì sao họ xây dựng nó,
nó làm gì, và nói về thiết kế của nó.

{{video "https://www.youtube.com/embed/yvjeIZgykiA"}}

## Tự viết trình biên dịch Go của riêng bạn

Elliot Stoneham giải thích tiềm năng của Go như một ngôn ngữ di động và
điểm qua các công cụ Go khiến điều đó trở nên hấp dẫn như vậy.

Anh ấy nói: "Dựa trên kinh nghiệm của tôi khi viết một trình biên dịch thử nghiệm từ Go sang Haxe,
tôi sẽ nói về các vấn đề thực tế của việc sinh mã và mô phỏng runtime cần thiết.
Tôi sẽ so sánh một số quyết định thiết kế của mình với những quyết định của hai trình biên dịch/trình chuyển đổi Go khác
xây dựng trên thư viện go.tools. Mục tiêu của tôi là khuyến khích bạn thử một trong những trình biên dịch Go 'đột biến' mới này.
Tôi hy vọng một số bạn sẽ được truyền cảm hứng để đóng góp cho một trong số chúng hoặc thậm chí viết một trình biên dịch mới của riêng mình."

{{video "https://www.youtube.com/embed/Qe8Dq7V3hXY"}}

## Thêm nữa

Còn rất nhiều bài nói tuyệt vời khác, vì vậy hãy xem toàn bộ loạt bài
[trong danh sách phát YouTube](http://www.youtube.com/playlist?list=PLtLJO5JKE5YDKG4WcaNts3IVZqhDmmuBH).
Đặc biệt, các [lightning talk](http://www.youtube.com/watch?v=cwpI5ONWGxc&list=PLtLJO5JKE5YDKG4WcaNts3IVZqhDmmuBH&index=7) rất thú vị.

Tôi muốn gửi lời cảm ơn cá nhân tới các diễn giả xuất sắc, Mathieu
Lonjaret vì đã quản lý thiết bị ghi hình, và đội ngũ FOSDEM vì đã biến tất cả điều này thành hiện thực.
