---
title: Một đề xuất để thêm Generics vào Go
date: 2021-01-12
by:
- Ian Lance Taylor
tags:
- go2
- proposals
- generics
summary: Generics đang bước vào quy trình đề xuất thay đổi ngôn ngữ
---

## Đề xuất generics

Chúng tôi đã gửi [một đề xuất thay đổi ngôn ngữ Go](/issue/43651) để bổ sung hỗ trợ cho
type parameter trên kiểu và hàm, cho phép một dạng lập trình generic.

## Vì sao cần generics?

Generics có thể mang lại cho chúng ta những khối xây dựng mạnh mẽ giúp chia sẻ mã
và xây dựng chương trình dễ dàng hơn.
Lập trình generic nghĩa là viết các hàm và cấu trúc dữ liệu trong đó
một số kiểu sẽ được chỉ rõ sau.
Ví dụ, bạn có thể viết một hàm thao tác trên một slice của một kiểu dữ liệu tùy ý,
trong đó kiểu dữ liệu cụ thể chỉ được xác định khi hàm được gọi.
Hoặc bạn có thể định nghĩa một cấu trúc dữ liệu lưu trữ giá trị của bất kỳ kiểu nào,
trong đó kiểu cụ thể cần lưu chỉ được xác định khi bạn tạo một thể hiện của cấu trúc đó.

Kể từ khi Go được phát hành lần đầu vào năm 2009, hỗ trợ generics luôn là
một trong những tính năng ngôn ngữ được yêu cầu nhiều nhất.
Bạn có thể đọc thêm về lý do generics hữu ích trong
[một bài blog trước đó](/blog/why-generics).

Dù generics có những trường hợp sử dụng rõ ràng, việc đưa chúng vào một cách gọn gàng trong
một ngôn ngữ như Go là một nhiệm vụ khó.
Một trong [những nỗ lực đầu tiên (có khiếm khuyết) nhằm thêm generics vào
Go](/design/15292/2010-06-type-functions) có từ tận năm 2010.
Đã có nhiều nỗ lực khác trong suốt thập kỷ qua.

Trong vài năm gần đây, chúng tôi đã làm việc trên một chuỗi các bản thảo thiết kế
và culminate ở [một thiết kế dựa trên type
parameter](/design/go2draft-type-parameters).
Bản thảo thiết kế này đã nhận được rất nhiều ý kiến từ cộng đồng lập trình Go,
và nhiều người đã thử nghiệm với nó bằng [generics playground](https://go2goplay.golang.org) được mô tả trong [một
bài blog trước đó](/blog/generics-next-step).
Ian Lance Taylor đã có [một bài nói tại GopherCon
2019](https://www.youtube.com/watch?v=WzgLqE-3IhY)
về lý do nên thêm generics và chiến lược mà chúng tôi đang theo đuổi.
Robert Griesemer đã có [một bài nói tiếp theo về các thay đổi trong thiết kế
và trong phần hiện thực tại GopherCon
2020](https://www.youtube.com/watch?v=TborQFPY2IM).
Các thay đổi ngôn ngữ này hoàn toàn tương thích ngược, vì vậy các chương trình Go hiện có
sẽ tiếp tục hoạt động chính xác như ngày nay.
Chúng tôi đã đạt tới điểm mà chúng tôi tin rằng bản thảo thiết kế đã đủ tốt,
và đủ đơn giản, để đề xuất đưa nó vào Go.

## Bây giờ sẽ thế nào?

[Quy trình đề xuất thay đổi ngôn ngữ](/s/proposal)
là cách chúng tôi thực hiện các thay đổi đối với ngôn ngữ Go.
Giờ đây chúng tôi đã [bắt đầu quy trình này](/issue/43651)
để thêm generics vào một phiên bản Go trong tương lai.
Chúng tôi hoan nghênh các phê bình và bình luận thực chất, nhưng xin cố tránh
lặp lại những bình luận trước đó, và cũng xin [tránh các bình luận cộng trừ đơn giản](/wiki/NoPlusOne).
Thay vào đó, hãy thêm reaction thumbs-up/thumbs-down vào các bình luận
mà bạn đồng ý hoặc không đồng ý, hoặc vào chính đề xuất nói chung.

Như với mọi đề xuất thay đổi ngôn ngữ, mục tiêu của chúng tôi là hướng tới
một đồng thuận để либо thêm generics vào ngôn ngữ, либо để đề xuất này dừng lại.
Chúng tôi hiểu rằng với một thay đổi có quy mô như vậy, sẽ không thể làm cho tất cả mọi người
trong cộng đồng Go hài lòng, nhưng chúng tôi có ý định đi tới một quyết định
mà tất cả đều có thể chấp nhận.

Nếu đề xuất được chấp nhận, mục tiêu của chúng tôi là có một phần hiện thực hoàn chỉnh,
dù có thể chưa được tối ưu hoàn toàn, để mọi người dùng thử vào cuối năm,
có thể như một phần của các bản beta Go 1.18.
