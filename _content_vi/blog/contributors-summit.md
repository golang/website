---
title: Contributors Summit
date: 2017-08-03
by:
- Sam Whited
tags:
- community
summary: Tường thuật từ Go Contributor Summit tại GopherCon 2017.
template: true
---

## Giới thiệu

Một ngày trước GopherCon, một nhóm thành viên đội Go và các contributor đã tụ họp tại Denver để thảo luận và lên kế hoạch cho tương lai của dự án Go.
Đây là sự kiện đầu tiên thuộc loại này, một cột mốc lớn đối với dự án Go.
Sự kiện gồm một phiên buổi sáng xoay quanh các cuộc thảo luận tập trung theo chủ đề, và một phiên buổi chiều gồm các cuộc thảo luận bàn tròn trong những nhóm nhỏ tách ra.

### Compiler và runtime

Phiên compiler và runtime bắt đầu bằng cuộc thảo luận về việc refactor `gc` và các công cụ liên quan thành các package có thể import.
Việc đó sẽ giảm overhead trong các công cụ lõi và trong các IDE vốn có thể tự nhúng compiler để kiểm tra cú pháp nhanh.
Mã cũng có thể được biên dịch hoàn toàn trong bộ nhớ, hữu ích trong những môi trường không cung cấp hệ thống tệp, hoặc để chạy kiểm thử liên tục trong quá trình phát triển nhằm có báo cáo sống về các lỗi vỡ.
Nhiều thảo luận hơn về việc có nên tiếp tục hướng công việc này hay không rất có thể sẽ còn xuất hiện trên các mailing list trong tương lai.

Cũng có rất nhiều thảo luận quanh việc thu hẹp khoảng cách giữa mã assembly tối ưu và Go.
Phần lớn mã crypto trong Go được viết bằng assembly vì lý do hiệu năng; điều này khiến nó khó gỡ lỗi, khó bảo trì, và khó đọc.
Hơn nữa, một khi đã dấn thân vào việc viết assembly, bạn thường không thể gọi ngược trở lại Go, làm hạn chế khả năng tái sử dụng mã.
Một bản viết lại bằng Go sẽ giúp việc bảo trì dễ hơn.
Việc thêm intrinsic của bộ xử lý và hỗ trợ tốt hơn cho phép toán 128-bit sẽ cải thiện hiệu năng crypto của Go.
Người ta đề xuất rằng package `math/bits` mới sẽ có trong 1.9 có thể được mở rộng cho mục đích này.

Vì không quá quen với sự phát triển của compiler và runtime, với tôi đây là một trong những phiên thú vị hơn cả trong ngày.
Tôi đã học được rất nhiều về tình trạng hiện tại, các vấn đề, và nơi mọi người muốn hướng tới.

### Quản lý dependency

Sau một cập nhật nhanh từ đội [dep](https://github.com/golang/dep) về tình trạng của dự án, phiên quản lý dependency chuyển dần sang cách thế giới Go sẽ hoạt động ra sao một khi dep (hoặc thứ gì đó tương tự dep) trở thành phương tiện quản lý package chính.
Công việc để làm cho Go dễ bắt đầu hơn và làm cho dep dễ dùng hơn đã được khởi động.
Trong Go 1.8, một giá trị mặc định cho `GOPATH` đã được đưa vào, nghĩa là người dùng chỉ cần thêm thư mục bin của Go vào `$PATH` trước khi có thể bắt đầu với dep.

Một cải tiến khả dụng khác trong tương lai mà dep có thể mở ra là cho phép Go làm việc từ bất kỳ thư mục nào (không chỉ từ workspace trong GOPATH), để mọi người có thể dùng cấu trúc thư mục và quy trình quen thuộc của họ với các ngôn ngữ khác.
Cũng có thể trong tương lai `go install` sẽ trở nên dễ dùng hơn bằng cách hướng dẫn người dùng thêm thư mục bin vào path, hoặc thậm chí tự động hóa quá trình đó.
Có rất nhiều lựa chọn tốt để làm bộ công cụ Go dễ dùng hơn, và thảo luận có lẽ sẽ còn tiếp diễn trên các mailing list.

### Thư viện chuẩn

Những cuộc thảo luận mà chúng tôi có quanh tương lai của ngôn ngữ Go phần lớn đã được Russ Cox đề cập trong bài blog của ông: [Toward Go 2](/blog//toward-go2), vậy nên ta hãy chuyển sang phiên về thư viện chuẩn.

Là một contributor của thư viện chuẩn và các subrepo, phiên này đặc biệt thú vị với tôi.
Điều gì nên đi vào thư viện chuẩn và subrepo, và chúng có thể thay đổi đến mức nào, là chủ đề chưa được xác định rõ.
Đội Go có thể gặp khó trong việc bảo trì một số lượng package khổng lồ khi họ có thể có hoặc không có người có chuyên môn cụ thể về từng lĩnh vực.
Để sửa những lỗi quan trọng trong các package của thư viện chuẩn, người ta phải đợi 6 tháng để phiên bản Go mới phát hành (hoặc phải có bản point release được phát hành trong trường hợp vấn đề bảo mật, điều làm tiêu tốn nguồn lực của đội).
Quản lý dependency tốt hơn có thể tạo điều kiện để di chuyển một số package ra khỏi thư viện chuẩn sang các dự án riêng với lịch phát hành riêng.

Cũng có thảo luận về những điều khó đạt được với các interface trong thư viện chuẩn.
Ví dụ, sẽ thật tuyệt nếu `io.Reader` nhận một context để các thao tác đọc bị chặn có thể bị hủy.

Cần có thêm [experience report](/wiki/experiencereports) trước khi ta có thể xác định điều gì sẽ thay đổi trong thư viện chuẩn.

### Công cụ và trình soạn thảo

Một language server để các trình soạn thảo dùng là chủ đề nóng trong phiên tooling, với nhiều người vận động để nhà phát triển IDE và công cụ cùng chấp nhận một “Go Language Server” chung nhằm lập chỉ mục và hiển thị thông tin về mã và package.
[Language Server Protocol](https://www.github.com/Microsoft/language-server-protocol) của Microsoft được gợi ý là điểm khởi đầu tốt vì được hỗ trợ rộng rãi trong các trình soạn thảo và IDE.

Jaana Burcu Dogan cũng nói về công việc của cô với distributed tracing và việc thông tin về các sự kiện runtime có thể được lấy dễ hơn và gắn vào trace như thế nào.
Một API “counter” chuẩn để báo số liệu thống kê đã được đề xuất, nhưng sẽ cần các experience report cụ thể từ cộng đồng trước khi một API như vậy có thể được thiết kế.

### Trải nghiệm contributor

Phiên cuối của ngày là về trải nghiệm contributor.
Cuộc thảo luận đầu tiên nói hoàn toàn về cách workflow Gerrit hiện tại có thể được làm dễ hơn cho contributor mới, điều đã dẫn đến các cải tiến trong tài liệu của nhiều repo và tác động đến workshop dành cho contributor mới vài ngày sau đó!

Việc giúp tìm ra tác vụ để làm dễ hơn, trao quyền cho người dùng thực hiện các công việc gardening trên issue tracker, và làm cho việc tìm reviewer dễ hơn cũng được cân nhắc.
Hy vọng chúng ta sẽ thấy những cải tiến ở các khía cạnh này và nhiều lĩnh vực khác của quy trình đóng góp trong các tuần và tháng tới!

### Các phiên breakout

Buổi chiều, người tham gia tách thành các nhóm nhỏ hơn để có những cuộc thảo luận sâu hơn về một số chủ đề từ phiên buổi sáng.
Những cuộc thảo luận này có các mục tiêu cụ thể hơn.
Ví dụ, một nhóm làm việc để xác định những phần hữu ích của một experience report và một danh sách tài liệu hiện có ghi lại trải nghiệm của người dùng Go, dẫn tới [trang wiki](/wiki/experiencereports) về experience report.

Một nhóm khác cân nhắc tương lai của lỗi trong Go.
Nhiều người dùng Go ban đầu bối rối hoặc không hiểu được sự thật rằng `error` là một interface, và việc gắn thêm thông tin vào lỗi mà không che mất các sentinel error như `io.EOF` có thể khó khăn.
Phiên breakout đã thảo luận các cách cụ thể để có thể sửa một số vấn đề đó trong những bản phát hành Go sắp tới, cũng như các cách xử lý lỗi có thể được cải thiện trong Go 2.

## Cộng đồng

Ngoài các thảo luận kỹ thuật, summit còn tạo cơ hội cho một nhóm người từ khắp nơi trên thế giới, những người thường xuyên nói chuyện và làm việc cùng nhau, được gặp trực tiếp, trong nhiều trường hợp là lần đầu tiên.
Không gì thay thế được một chút thời gian gặp mặt trực tiếp để xây dựng cảm giác tôn trọng lẫn nhau và tình đồng chí, điều vô cùng quan trọng khi một nhóm đa dạng với nền tảng và ý tưởng khác nhau cần cùng nhau làm việc trong một cộng đồng chung.
Trong các giờ nghỉ, thành viên đội Go tỏa ra giữa các contributor để thảo luận về Go cũng như giao lưu xã hội một chút, điều thật sự giúp gắn khuôn mặt với những cái tên vẫn duyệt mã của chúng ta hằng ngày.

Như Russ đã bàn trong [Toward Go 2](/blog//toward-go2), giao tiếp hiệu quả đòi hỏi phải hiểu rõ đối tượng của mình.
Việc có một mẫu rộng các contributor Go trong cùng một căn phòng đã giúp tất cả chúng tôi hiểu hơn về đối tượng của Go và khởi đầu nhiều cuộc thảo luận hiệu quả về tương lai của Go.
Trong tương lai, chúng tôi hy vọng sẽ có nhiều sự kiện kiểu này thường xuyên hơn để thúc đẩy đối thoại và cảm giác cộng đồng.

{{image "contributors-summit/IMG_20170712_145844.jpg"}}
{{image "contributors-summit/IMG_20170712_145854.jpg"}}
{{image "contributors-summit/IMG_20170712_145905.jpg"}}
{{image "contributors-summit/IMG_20170712_145911.jpg"}}
{{image "contributors-summit/IMG_20170712_145950.jpg"}}

Ảnh bởi Steve Francia

