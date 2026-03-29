---
title: Cập nhật Bộ Quy tắc Ứng xử
date: 2021-09-16
by:
 - Carmen Andoh
 - Russ Cox
 - Steve Francia
summary: Một cập nhật nhỏ cho, và một cập nhật về việc thực thi, Bộ Quy tắc Ứng xử của Go
---

Mặc dù các chi tiết của [Bộ Quy tắc Ứng xử](/conduct) của chúng tôi đã được
[điều chỉnh](/blog/conduct-2018) theo thời gian, [mục tiêu của chúng tôi](/blog/open-source#code-of-conduct) thì không đổi.
Chúng tôi muốn cộng đồng Go trở nên bao trùm, chào đón, hữu ích và tôn trọng nhất có thể.
Nếu bạn muốn dùng hoặc thảo luận về Go, chúng tôi muốn bạn cảm thấy được chào đón ở đây.

Cộng đồng đã đủ lớn để,
thay vì giả định rằng mọi người đều biết điều gì được kỳ vọng ở họ,
Bộ Quy tắc Ứng xử đóng vai trò như một thỏa thuận,
đặt ra những kỳ vọng rõ ràng cho hành vi của chúng ta trong cả tương tác trực tuyến lẫn ngoại tuyến.
Nếu chúng ta không đáp ứng thỏa thuận đó,
mọi người có thể chỉ ra điều ấy và chúng ta có thể sửa hành vi của mình.

Trong bài viết này, chúng tôi muốn cung cấp hai cập nhật:
thứ nhất, cập nhật về cách chúng tôi tiếp cận việc thực thi Bộ Quy tắc Ứng xử,
và thứ hai, cập nhật đối với chính các Gopher Value.

## Thực thi

Chúng tôi muốn mọi người cảm thấy được chào đón ở đây.
Điều gì xảy ra khi thành viên trong cộng đồng khiến người khác cảm thấy không được chào đón?
Những hành vi đó có thể được báo cáo cho Project Steward,
người làm việc với một ủy ban từ Open Source Programs Office của Google
để xác định nên xử lý từng báo cáo như thế nào.

Kể từ [lần sửa đổi Bộ Quy tắc Ứng xử tháng 5 năm 2018](/blog/conduct-2018), các thành viên cộng đồng đã gửi hơn 300 báo cáo về hành vi,
trung bình khoảng một đến hai báo cáo mỗi tuần.
Kết quả điển hình là gặp trực tiếp người có hành vi bị báo cáo
và giúp họ hiểu cách chịu trách nhiệm về hành động của mình và sửa chúng trong tương lai.

Nhưng còn những người làm điều tệ hơn là chỉ khiến người khác cảm thấy không được chào đón,
hoặc từ chối sửa đổi hành vi của họ thì sao?
[Nghịch lý của sự khoan dung](https://en.wikipedia.org/wiki/Paradox_of_tolerance)
là nhóm người duy nhất mà chúng ta không thể chào đón là những người khiến người khác cảm thấy không được chào đón.
Chúng ta buộc phải lựa chọn giữa họ và những người mà họ sẽ đẩy đi.
Chúng tôi chọn đứng về phía những người thân thiện, chào đón, kiên nhẫn, chu đáo,
tôn trọng, rộng lượng và mang tính xây dựng, những người hành xử theo Gopher Values
và khiến cộng đồng của chúng tôi trở nên tốt đẹp hơn.

Khi một người khiến người khác cảm thấy không được chào đón buộc phải bị loại trừ,
các kênh hoặc không gian cụ thể có thể chặn cá nhân đó dựa trên quan sát của chính họ,
mà không cần chờ báo cáo hành vi.
Ví dụ, đội phát hành Go,
nhóm chịu trách nhiệm chính trong việc chăm sóc issue tracker của Go,
có thể nhanh chóng chặn người dùng hành xử không phù hợp (thường là lăng mạ bằng lời nói)
mà không cần liên hệ với bất kỳ ai trong chúng tôi. Tính đến nay, họ đã chặn khoảng một chục tài khoản.

Hành vi được báo cáo cũng có thể (nhưng rất hiếm khi) đạt đến mức
khiến chúng tôi cân nhắc bước nghiêm trọng là trục xuất một thành viên cộng đồng,
tạm thời hoặc vĩnh viễn, khỏi mọi không gian do chúng tôi tổ chức:
danh sách thư, issue tracker, các sự kiện có thư mời, v.v.

Ví dụ về kiểu sai phạm có thể dẫn đến việc trục xuất toàn cộng đồng bao gồm:

1. Đe dọa người khác.
2. Lạm dụng hoặc hành hung người khác.
3. Brigading hoặc khuyến khích hay điều phối hành vi ngược đãi trực tuyến theo cách khác.
4. Cố tình lập báo cáo hành vi sai sự thật về người khác.
5. Quấy rối người khác.
   Quấy rối có thể là một sự cố nghiêm trọng đơn lẻ hoặc một chuỗi dài các sự cố nhỏ hơn.
6. Hành vi giáp ranh kéo dài.
   Từng vi phạm riêng lẻ có thể trông không đáng kể,
   nhưng lặp lại theo thời gian sẽ tạo thành một mẫu hành vi
   không phù hợp với Gopher Values của chúng tôi
   và cộng dồn thành tổn hại đáng kể.

Trục xuất không phải điều nên cân nhắc một cách nhẹ tay.
Cho tới nay, chỉ có một số lượng nhỏ (một chữ số) cá nhân
bị trục xuất hoàn toàn khỏi các không gian của Go.

## Một Gopher Value mới

Một chủ đề lặp đi lặp lại mà chúng tôi thấy trong các báo cáo về những vấn đề nhỏ
là mọi người không chấp nhận rằng lời nói và hành động của họ ảnh hưởng tới người khác.
Trong những trường hợp cực đoan, có người nói những điều như “nhưng đây là internet mà.”
Chúng tôi mong muốn trở nên chào đón hơn internet nói chung rất nhiều.
Vì mục tiêu đó, chúng tôi đang bổ sung thêm một Gopher value vào Bộ Quy tắc Ứng xử:
“Hãy có trách nhiệm.”

Toàn bộ phần “[Gopher Values](/conduct#values)”
của Bộ Quy tắc Ứng xử giờ như sau:

  - **Hãy thân thiện và chào đón.**

  - **Hãy kiên nhẫn.**
      - Hãy nhớ rằng mọi người có phong cách giao tiếp khác nhau và không phải
        ai cũng dùng ngôn ngữ mẹ đẻ của mình.
        (Ý nghĩa và sắc thái có thể bị mất đi trong quá trình dịch.)

  - **Hãy chu đáo.**
      - Giao tiếp hiệu quả đòi hỏi nỗ lực.
        Hãy nghĩ về cách lời nói của bạn sẽ được diễn giải.
      - Hãy nhớ rằng đôi khi tốt nhất là đừng bình luận gì cả.

  - **Hãy tôn trọng.**
      - Đặc biệt, hãy tôn trọng sự khác biệt về quan điểm.

  - **Hãy rộng lượng.**
      - Hãy diễn giải lập luận của người khác trong thiện chí, đừng cố tìm cách bất đồng.
      - Khi chúng ta bất đồng, hãy cố hiểu tại sao.

  - **Hãy mang tính xây dựng.**
      - Tránh làm chệch hướng: bám sát chủ đề; nếu bạn muốn nói về điều gì khác,
        hãy bắt đầu một cuộc trò chuyện mới.
      - Tránh chỉ trích thiếu xây dựng: đừng chỉ than phiền về hiện trạng;
        hãy đưa ra, hoặc ít nhất là mời gọi, các đề xuất về cách cải thiện.
      - Tránh cà khịa (những bình luận ngắn gọn, thiếu hiệu quả, mang tính chĩa mũi dùi)
      - Tránh thảo luận các vấn đề có khả năng xúc phạm hoặc nhạy cảm;
        điều này quá thường xuyên dẫn tới xung đột không cần thiết.
      - Tránh microaggression (những sự xúc phạm nhỏ, ngắn và thường nhật bằng lời nói, hành vi hoặc
        môi trường truyền đi các sự hạ thấp, xúc phạm và ác ý
        đối với một người hoặc một nhóm).

  - **Hãy có trách nhiệm.**
      - Điều bạn nói và làm đều quan trọng.
        Hãy chịu trách nhiệm cho lời nói và hành động của mình, bao gồm cả hệ quả của chúng,
        dù là có chủ ý hay không.

Hai năm qua đã rất khó khăn:
thế giới cực kỳ bất định và có thể sẽ còn như vậy trong tương lai gần.
Mọi người đang căng thẳng và kiệt sức.
Trong những thời điểm như thế này, việc ghi nhớ Gopher Values lại càng quan trọng hơn bao giờ hết.

Khi chúng ta tiếp tục phát triển Go, điều tối quan trọng là nuôi dưỡng và duy trì cảm giác về _trách nhiệm tập thể_.
Là thành viên của cộng đồng Go, tất cả chúng ta phải ý thức
về việc hành động và cách cư xử của mình ảnh hưởng thế nào tới những nhóm và không gian nơi chúng ta cộng tác.
Chúng ta phải chịu trách nhiệm cho hành vi của mình, tác động của nó,
và việc nó có thể khuyến khích người khác tiến tới (hoặc rời xa) hợp tác mang tính xây dựng như thế nào.

Đó cũng là trách nhiệm tập thể của chúng ta trong việc lên tiếng để bảo đảm một động lực nhóm
cho phép sự trao đổi ý tưởng và quan điểm một cách hiệu quả.
Bộ Quy tắc Ứng xử của chúng ta nêu rõ rằng chúng ta phải “tôn trọng các quan điểm và trải nghiệm khác biệt.”
Khi bất đồng trở thành bất hòa hay thiếu tôn trọng, chúng ta phải lên tiếng.
Bất kể quan điểm của mình là gì, chúng ta có trách nhiệm đối với hành động và tác động của chúng.

Khi cam kết với những giá trị này,
chúng ta tạo ra một môi trường an toàn và chào đón cho mọi người, đáp ứng các mục tiêu cộng đồng chung:
hợp tác và giao tiếp mang tính xây dựng để đưa Go tới thành công.

Cảm ơn tất cả mọi người đã tham gia cùng chúng tôi để trở thành một phần của cộng đồng Go.
Chúng tôi hy vọng bạn cảm thấy được chào đón, và chúng tôi sẽ tiếp tục làm việc
để bao gồm nhiều người nhất có thể.
Nếu bạn có bất kỳ câu hỏi hoặc lo ngại nào,
xin cứ thoải mái liên hệ với bất kỳ ai trong chúng tôi qua email:
_carmen@golang.org_ (Carmen),
_rsc@golang.org_ (Russ),
và
_spf@golang.org_ (Steve).
