---
title: Cập nhật Bộ Quy tắc Ứng xử của Go
date: 2018-05-23
by:
- Steve Francia
tags:
- conduct
summary: Sửa đổi Bộ Quy tắc Ứng xử của Go.
---


Vào tháng 11 năm 2015, chúng tôi đã giới thiệu Bộ Quy tắc Ứng xử của Go.
Nó được xây dựng từ sự hợp tác giữa
các thành viên đội Go tại Google và cộng đồng Go.
Tôi may mắn là một trong những thành viên cộng đồng
được mời tham gia cả quá trình soạn thảo lẫn sau đó là thực thi
Bộ Quy tắc Ứng xử của Go.
Kể từ đó, chúng tôi đã rút ra hai bài học về
những giới hạn trong bộ quy tắc ứng xử đã khiến chúng tôi
không thể nuôi dưỡng nền văn hóa an toàn
cần thiết cho thành công của Go.

Bài học đầu tiên chúng tôi học được là các hành vi độc hại của
người tham gia dự án ở ngoài các không gian của dự án có thể gây
tác động tiêu cực lên dự án, ảnh hưởng đến sự an toàn và an ninh của
các thành viên cộng đồng. Đã có một vài
sự cố được báo cáo mà hành động diễn ra bên ngoài không gian của dự án
nhưng tác động lại được cảm nhận trong cộng đồng của chúng tôi. Phần ngôn ngữ
cụ thể trong bộ quy tắc ứng xử của chúng tôi đã giới hạn khả năng
phản hồi của chúng tôi chỉ với những hành động xảy ra “trong các
diễn đàn chính thức do dự án Go vận hành”. Chúng tôi cần một cách
để bảo vệ các thành viên cộng đồng ở bất cứ nơi đâu họ hiện diện.

Bài học thứ hai chúng tôi học được là những yêu cầu cần thiết
để thực thi bộ quy tắc
ứng xử đã đặt gánh nặng quá lớn lên các tình nguyện viên.
Phiên bản đầu tiên của bộ quy tắc ứng xử trình bày
nhóm công tác như những người kỷ luật. Chẳng bao lâu sau,
rõ ràng điều này là quá sức, nên vào đầu năm 2017 [chúng tôi đã thay đổi vai trò của nhóm](/cl/37014)
thành cố vấn và hòa giải viên.
Dù vậy, các thành viên cộng đồng trong nhóm công tác
vẫn cho biết họ cảm thấy quá tải, thiếu đào tạo và dễ tổn thương.
Sự thay đổi đầy thiện chí này khiến chúng tôi không còn cơ chế thực thi
mà vẫn không giải quyết được vấn đề gánh nặng trên các tình nguyện viên.

Vào giữa năm 2017, tôi đại diện dự án Go trong một cuộc họp với
Open Source Programs Office và Open Source Strategy Team của Google
để giải quyết những thiếu sót trong
bộ quy tắc ứng xử tương ứng của chúng tôi, đặc biệt là trong phần thực thi.
Rất nhanh chóng, chúng tôi nhận ra rằng các vấn đề của mình có rất nhiều điểm chung,
và rằng cùng nhau xây dựng một bộ quy tắc ứng xử chung cho tất cả
các dự án mã nguồn mở của Google là điều hợp lý.
Chúng tôi bắt đầu với văn bản của
Contributor Covenant Code of Conduct v1.4
rồi thực hiện các thay đổi, chịu ảnh hưởng từ
kinh nghiệm của chúng tôi trong cộng đồng Go và những kinh nghiệm tập thể của chúng tôi với mã nguồn mở.
Kết quả là [mẫu bộ quy tắc ứng xử](https://opensource.google.com/docs/releasing/template/CODE_OF_CONDUCT/) của Google.

Hôm nay, dự án Go đang áp dụng bộ quy tắc ứng xử mới này,
và chúng tôi đã cập nhật [golang.org/conduct](/conduct).
Bộ quy tắc ứng xử đã sửa đổi này giữ lại phần lớn ý định, cấu trúc và
ngôn ngữ của bộ quy tắc ứng xử Go nguyên bản, đồng thời đưa ra hai
thay đổi cốt lõi giải quyết những thiếu sót đã được xác định ở trên.

Thứ nhất, [bộ quy tắc ứng xử mới làm rõ](/conduct/#scope) rằng những người
tham gia vào bất kỳ hình thức quấy rối hoặc hành vi không phù hợp nào,
ngay cả bên ngoài không gian của dự án, đều không được chào đón trong không gian dự án của chúng tôi.
Điều này có nghĩa là Bộ Quy tắc Ứng xử được áp dụng cả bên ngoài
không gian dự án khi có cơ sở hợp lý để tin rằng
hành vi của một cá nhân có thể gây tác động tiêu cực
lên dự án hoặc cộng đồng của nó.

Thứ hai, thay cho nhóm công tác,
[bộ quy tắc ứng xử mới giới thiệu một Project Steward duy nhất](/conduct/#reporting)
người sẽ được đào tạo và hỗ trợ rõ ràng cho vai trò này.
Project Steward sẽ nhận các báo cáo vi phạm
rồi làm việc với một ủy ban,
gồm đại diện từ Open Source Programs Office
và Google Open Source Strategy team,
để tìm ra cách giải quyết.

Project Steward đầu tiên của chúng tôi sẽ là [Cassandra Salisbury](https://twitter.com/cassandraoid).
Cô ấy được cộng đồng Go biết đến rộng rãi với vai trò thành viên của Go Bridge,
người tổ chức nhiều buổi meetup và hội nghị về Go,
và là người dẫn dắt nhóm công tác tiếp cận cộng đồng Go.
Cassandra hiện làm việc trong đội Go tại Google
với trọng tâm là vận động và hỗ trợ cộng đồng Go.

Chúng tôi biết ơn tất cả những ai đã phục vụ trong Nhóm Công tác Bộ Quy tắc Ứng xử
ban đầu. Những nỗ lực của các bạn là nền tảng để tạo nên một cộng đồng
bao trùm và an toàn.

Chúng tôi tin rằng bộ quy tắc ứng xử đã góp phần giúp
dự án Go trở nên chào đón hơn hiện nay so với năm 2015,
và tất cả chúng ta nên tự hào vì điều đó.

Chúng tôi hy vọng bộ quy tắc ứng xử mới sẽ giúp bảo vệ các thành viên cộng đồng
của chúng tôi hiệu quả hơn nữa.
