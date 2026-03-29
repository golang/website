---
title: Contributors Summit 2019
date: 2019-08-15
by:
- Carmen Andoh and contributors
tags:
- community
summary: Tường thuật từ Go Contributor Summit tại GopherCon 2019.
template: true
---

## Giới thiệu

Năm thứ ba liên tiếp, đội Go và các contributor đã tụ họp một ngày trước GopherCon để thảo luận và lên kế hoạch cho tương lai của dự án Go.
Sự kiện bao gồm việc tự tổ chức thành các nhóm breakout, một cuộc thảo luận kiểu town-hall về quy trình proposal vào buổi sáng, và các cuộc thảo luận bàn tròn breakout vào buổi chiều dựa trên các chủ đề do contributor lựa chọn.
Chúng tôi đã mời năm contributor viết về trải nghiệm của họ trong nhiều cuộc thảo luận tại summit năm nay.

{{image "contributors-summit-2019/group.jpg" 800}}

_(Ảnh bởi Steve Francia.)_

## Compiler và Runtime (tường thuật bởi Lynn Boger)

Go contributors summit là một cơ hội tuyệt vời để gặp gỡ và thảo luận các chủ đề và ý tưởng với những người khác cũng đang đóng góp cho Go.

Ngày bắt đầu bằng khoảng thời gian để mọi người trong phòng làm quen với nhau.
Có sự pha trộn tốt giữa đội Go cốt lõi và những người khác đang tích cực đóng góp cho Go.
Từ đó, chúng tôi quyết định những chủ đề nào được quan tâm và chia nhóm lớn thành các nhóm nhỏ hơn.
Lĩnh vực quan tâm của tôi là compiler, nên tôi tham gia nhóm đó và ở lại với họ gần như suốt thời gian.

Ở cuộc họp đầu tiên, một danh sách dài các chủ đề được đưa ra và vì vậy nhóm compiler quyết định tiếp tục gặp nhau trong cả ngày.
Tôi có vài chủ đề quan tâm mà mình chia sẻ và cũng có nhiều chủ đề người khác nêu ra mà tôi thấy hứng thú.
Không phải mọi mục trong danh sách đều được bàn sâu; dưới đây là danh sách những chủ đề có nhiều quan tâm và thảo luận nhất, tiếp theo là một số bình luận ngắn về các chủ đề khác.

**Kích thước tệp nhị phân**.
Mọi người bày tỏ lo ngại về kích thước tệp nhị phân, đặc biệt là việc nó tiếp tục tăng theo mỗi bản phát hành.
Một số nguyên nhân khả dĩ được xác định như tăng mức inlining và các tối ưu hóa khác.
Có lẽ tồn tại một nhóm người dùng muốn tệp nhị phân nhỏ, một nhóm khác muốn hiệu năng tốt nhất có thể, và cũng có thể có những người không quá quan tâm.
Điều này dẫn sang chủ đề TinyGo, và mọi người ghi nhận rằng TinyGo không phải là hiện thực đầy đủ của Go và điều quan trọng là tránh để TinyGo tách khỏi Go và chia cắt người dùng.
Cần thêm nghiên cứu để hiểu nhu cầu của người dùng và các nguyên nhân chính xác góp phần vào kích thước hiện tại.
Nếu có cơ hội giảm kích thước mà không ảnh hưởng hiệu năng, những thay đổi đó có thể được thực hiện, nhưng nếu hiệu năng bị ảnh hưởng thì một số người dùng vẫn sẽ ưu tiên hiệu năng hơn.

**Assembly vector**.
Cách tận dụng assembly vector trong Go được thảo luận khá lâu và đã là chủ đề được quan tâm từ trước.
Tôi tách nó thành ba khả năng riêng, vì tất cả đều liên quan tới việc dùng chỉ lệnh vector nhưng khác nhau ở cách dùng, bắt đầu với chủ đề assembly vector.
Đây lại là một trường hợp đánh đổi của compiler.

Với hầu hết các đích, có những hàm quan trọng trong các package chuẩn như crypto, hash, math và các package khác, nơi việc dùng assembly là cần thiết để có hiệu năng tốt nhất; tuy nhiên việc có những hàm lớn viết bằng assembly khiến chúng khó hỗ trợ và bảo trì, và có thể đòi hỏi các hiện thực khác nhau cho từng nền tảng đích.
Một giải pháp là dùng macro assembly hoặc các kỹ thuật sinh mã cấp cao khác để làm assembly vector dễ đọc và dễ hiểu hơn.

Một mặt khác của câu hỏi này là liệu compiler Go có thể trực tiếp sinh ra các chỉ lệnh SIMD vector khi biên dịch tệp mã nguồn Go hay không, bằng cách nâng cấp compiler Go để biến đổi chuỗi mã nhằm “simdize” mã và tận dụng các chỉ lệnh vector.
Việc triển khai SIMD trong compiler Go sẽ làm tăng độ phức tạp và thời gian biên dịch, và không phải lúc nào cũng cho ra mã chạy nhanh hơn.
Cách mã được biến đổi trong một số trường hợp còn có thể phụ thuộc vào nền tảng đích nên điều đó không lý tưởng.

Một cách khác để tận dụng chỉ lệnh vector trong Go là cung cấp cơ chế giúp dễ dùng các chỉ lệnh vector trực tiếp từ mã nguồn Go.
Những chủ đề được thảo luận gồm intrinsics, hoặc các hiện thực vốn đã tồn tại trong những compiler khác như Rust.
Trong gcc, một số nền tảng cung cấp inline asm, và Go có thể cũng có thể cung cấp khả năng đó, nhưng theo kinh nghiệm của tôi, việc trộn inline asm với mã Go làm tăng độ phức tạp của compiler về mặt theo dõi việc dùng thanh ghi và gỡ lỗi.
Nó cho phép người dùng làm những điều mà compiler có thể không mong đợi hay không muốn, và thực sự làm tăng thêm một lớp phức tạp.
Nó cũng có thể bị chèn ở những chỗ không lý tưởng.

Tóm lại, điều quan trọng là phải cung cấp cách tận dụng các chỉ lệnh vector sẵn có, và làm cho việc viết chúng dễ hơn và an toàn hơn.
Ở nơi có thể, các hàm nên dùng nhiều mã Go nhất có thể, và có thể tìm ra một cách dùng assembly cấp cao.
Cũng có thảo luận về việc thiết kế một package vector thử nghiệm để triển khai thử một số ý tưởng này.

**Quy ước gọi mới**.
Một số người quan tâm đến chủ đề [các thay đổi ABI để cung cấp quy ước gọi dựa trên thanh ghi](/issue/18597).
Tình trạng hiện tại đã được báo cáo cùng với các chi tiết.
Mọi người thảo luận những gì còn phải làm trước khi nó có thể được sử dụng.
Đặc tả ABI cần được viết trước và chưa rõ khi nào điều đó sẽ xong.
Tôi biết điều này sẽ có lợi cho một số nền tảng nhiều hơn những nền tảng khác và quy ước gọi bằng thanh ghi được dùng trong hầu hết compiler cho các nền tảng khác.

**Tối ưu hóa tổng quát**.
Một số tối ưu hóa có lợi hơn cho vài nền tảng ngoài x86 đã được thảo luận.
Đặc biệt, các tối ưu hóa vòng lặp như hoisting các invariant và strength reduction có thể được thực hiện và đem lại lợi ích lớn hơn trên một số nền tảng.
Các giải pháp khả dĩ đã được thảo luận, và việc triển khai có lẽ sẽ tùy thuộc vào các đích coi những cải tiến này là quan trọng.

**Tối ưu hóa dựa trên phản hồi**.
Chủ đề này được thảo luận và tranh luận như một hướng cải tiến tương lai.
Theo kinh nghiệm của tôi, rất khó tìm ra các chương trình có ý nghĩa để dùng cho việc thu thập dữ liệu hiệu năng rồi sau đó dùng dữ liệu đó để tối ưu mã.
Nó làm tăng thời gian biên dịch và cần nhiều không gian để lưu dữ liệu, trong khi dữ liệu đó có thể chỉ có ý nghĩa với một tập nhỏ chương trình.

**Các thay đổi đang chuẩn bị gửi**.
Một vài thành viên trong nhóm nhắc đến những thay đổi mà họ đang thực hiện và dự định sớm gửi lên, bao gồm các cải tiến cho makeslice và viết lại rulegen.

**Mối lo về thời gian biên dịch**.
Thời gian biên dịch được thảo luận ngắn gọn. Mọi người lưu ý rằng phase timing đã được thêm vào đầu ra GOSSAFUNC.

**Giao tiếp giữa các contributor compiler**.
Có người hỏi liệu có cần một mailing list riêng cho compiler Go hay không.
Người ta gợi ý dùng golang-dev cho mục đích đó, thêm từ compiler vào tiêu đề thư để nhận diện.
Nếu lưu lượng trên golang-dev quá lớn thì khi đó có thể cân nhắc một mailing list riêng cho compiler ở thời điểm sau.

**Cộng đồng**.
Tôi thấy ngày hôm đó rất có ích trong việc kết nối với những người đã hoạt động tích cực trong cộng đồng và có cùng lĩnh vực quan tâm.
Tôi có thể gặp nhiều người mà trước đó chỉ biết qua tên tài khoản xuất hiện trong issue, mailing list hay CL.
Tôi có thể thảo luận một số chủ đề và issue hiện có, và nhận được phản hồi tương tác trực tiếp thay vì phải chờ phản hồi trực tuyến.
Tôi được khuyến khích viết issue về những vấn đề mình từng thấy.
Các kết nối này không chỉ diễn ra trong ngày hôm đó mà còn khi gặp lại những người khác trong suốt hội nghị, sau khi đã được giới thiệu trong ngày đầu tiên, dẫn đến nhiều cuộc thảo luận thú vị.
Hy vọng những kết nối này sẽ dẫn đến giao tiếp hiệu quả hơn và cách xử lý issue cùng thay đổi mã tốt hơn trong tương lai.

## Công cụ (tường thuật bởi Paul Jolly)

Phiên breakout về công cụ tại contributor summit có dạng mở rộng, với hai phiên tiếp theo trong các ngày hội nghị chính được nhóm [golang-tools](/wiki/golang-tools) tổ chức.
Bản tóm tắt này được chia thành hai phần: phiên công cụ tại contributor workshop, và báo cáo tổng hợp từ các phiên golang-tools trong những ngày hội nghị chính.

**Contributor summit**.
Phiên công cụ bắt đầu bằng phần giới thiệu của khoảng 25 người tham dự, tiếp theo là một màn brainstorm các chủ đề, bao gồm: gopls, ARM 32-bit, eval, signal, analysis, API go/packages, refactoring, pprof, trải nghiệm module, phân tích mono repo, go mobile, dependency, tích hợp editor, quyết định tối ưu của compiler, debugging, visualization, documentation.
Rất nhiều người với rất nhiều mối quan tâm về rất nhiều công cụ!

Phiên này tập trung vào hai lĩnh vực (toàn bộ thời lượng chỉ cho phép vậy): gopls và trực quan hóa.
[Gopls](/wiki/gopls) (đọc là “go please”) là một hiện thực máy chủ [Language Server Protocol (LSP)](https://langserver.org) cho Go.
Rebecca Stambler, tác giả chính của gopls, cùng phần còn lại của đội Go tools muốn nghe trải nghiệm của mọi người với gopls: độ ổn định, tính năng còn thiếu, tích hợp với editor có hoạt động không, v.v.?
Cảm nhận chung là gopls đang ở trạng thái rất tốt và hoạt động cực kỳ tốt với đa số trường hợp sử dụng.
Độ bao phủ kiểm thử tích hợp cần được cải thiện, nhưng đó là một bài toán khó để làm “đúng” trên mọi editor.
Chúng tôi đã thảo luận về một cách tốt hơn để người dùng báo lỗi gopls mà họ gặp phải qua editor, telemetry/diagnostics, các chỉ số hiệu năng của gopls, tất cả đều là những chủ đề được thảo luận kỹ hơn ở các phiên golang-tools diễn ra trong các ngày hội nghị chính (xem bên dưới).
Một lĩnh vực thảo luận then chốt là cách mở rộng gopls, ví dụ dưới dạng các kiểm tra kiểu vet dựa trên go/analysis bổ sung, kiểm tra lint, refactoring, v.v.
Hiện chưa có giải pháp tốt, nhưng vấn đề này đang được điều tra tích cực.
Cuộc trò chuyện chuyển sang chủ đề rất rộng là trực quan hóa, với màn giới thiệu dựa trên demo của Anthony Starks (người nhân tiện cũng đã có một bài nói xuất sắc về [Go cho các màn hình hiển thị thông tin](https://www.youtube.com/watch?v=NyDNJnioWhI) tại GopherCon 2018).

**Các ngày hội nghị**.
Những phiên golang-tools trong các ngày hội nghị chính là sự tiếp nối của [các cuộc gọi hằng tháng](/wiki/golang-tools) đã diễn ra kể từ khi nhóm được hình thành ở GopherCon 2018.
Biên bản đầy đủ có tại phiên [ngày 1](https://docs.google.com/document/d/1-RVyttQ0ncjCpR_sRwizf-Ubedkr0Emwmk2LhnsUOmE/edit) và [ngày 2](https://docs.google.com/document/d/1ZI_WqpLCB8DO6teJ3aBuXTeYD2iZZZlkDptmcY6Ja60/edit#heading=h.x9lkytc2gxmg).
Những phiên này cũng được tham dự rất đông, với 25-30 người mỗi phiên.
Đội Go tools xuất hiện rất mạnh (một dấu hiệu tốt cho mức hỗ trợ đang được đầu tư cho lĩnh vực này), và đội nền tảng của Uber cũng có mặt.
Khác với contributor summit, mục tiêu của các phiên này là ra về với những action item cụ thể.

**Gopls**.
“Mức sẵn sàng” của gopls là trọng tâm lớn của cả hai phiên.
Câu hỏi này về cơ bản quy về việc xác định thời điểm thích hợp để nói với những người tích hợp editor rằng “chúng tôi đã có phiên bản cắt đầu tiên tốt của gopls”, rồi biên soạn danh sách các tích hợp/plugin editor “được chứng nhận” là hoạt động tốt với gopls.
Trung tâm của quá trình “chứng nhận” đó là một quy trình được định nghĩa rõ ràng để người dùng có thể báo những vấn đề họ gặp với gopls.
Hiệu năng và bộ nhớ không phải là blockers cho “bản phát hành” đầu tiên này.
Cuộc trò chuyện về cách mở rộng gopls, bắt đầu từ contributor summit hôm trước, tiếp tục được đào sâu.
Dù có rất nhiều lợi ích rõ ràng và sức hấp dẫn trong việc mở rộng gopls (kiểm tra go/analysis tùy biến, hỗ trợ linter, refactoring, sinh mã...), vẫn chưa có câu trả lời rõ ràng cho cách hiện thực hóa điều này theo cách có khả năng mở rộng.
Những người tham dự đồng ý rằng điều này không nên bị xem là blocker cho “bản phát hành” đầu tiên, nhưng vẫn cần tiếp tục được nghiên cứu.
Theo tinh thần của gopls và tích hợp editor, Heschi Kreinick từ đội Go tools nêu thêm chủ đề hỗ trợ debugging.
Delve đã trở thành debugger thực tế của Go và đang ở trạng thái tốt; giờ cần xác định tình trạng tích hợp debugger-editor, theo quy trình tương tự như với gopls và các tích hợp “được chứng nhận”.

**Go Discovery Site**.
Phiên golang-tools thứ hai bắt đầu với màn giới thiệu rất tốt về Go Discovery Site của Julie Qiu từ đội Go tools, cùng một demo nhanh.
Julie nói về kế hoạch cho Discovery Site: open-source dự án, những tín hiệu được dùng trong xếp hạng tìm kiếm, cách [godoc.org](http://godoc.org/) cuối cùng sẽ được thay thế, cách submodule nên hoạt động, và cách người dùng có thể khám phá các major version mới.

**Build Tag**.
Sau đó cuộc trò chuyện chuyển sang hỗ trợ build tag trong gopls.
Đây là một lĩnh vực rõ ràng cần được hiểu tốt hơn (các trường hợp sử dụng hiện đang được thu thập trong [issue 33389](/issue/33389)).
Sau cuộc thảo luận này, phiên họp khép lại với đề nghị từ Alexander Zolotov thuộc đội JetBrains GoLand rằng đội gopls và đội GoLand nên chia sẻ kinh nghiệm ở lĩnh vực này cũng như nhiều lĩnh vực khác, vì GoLand đã có khá nhiều kinh nghiệm trước đó.

**Hãy tham gia cùng chúng tôi!**
Chúng tôi hoàn toàn có thể nói về các chủ đề liên quan tới công cụ suốt nhiều ngày!
Tin tốt là các cuộc gọi golang-tools sẽ tiếp tục diễn ra trong tương lai gần.
Mọi người quan tâm tới công cụ Go đều rất được khuyến khích tham gia: [wiki](/wiki/golang-tools) có thêm chi tiết.

## Sử dụng trong doanh nghiệp (tường thuật bởi Daniel Theophanes)

Việc chủ động tìm hiểu nhu cầu của những nhà phát triển ít lên tiếng hơn sẽ là thách thức lớn nhất, và cũng là chiến thắng lớn nhất, của ngôn ngữ Go. Có một bộ phận lớn lập trình viên không chủ động tham gia vào cộng đồng Go.
Một số là đối tác kinh doanh, người làm marketing, hay kiểm thử chất lượng nhưng cũng làm phát triển.
Một số khác lại đội mũ quản lý và đưa ra quyết định tuyển dụng hoặc công nghệ.
Những người khác chỉ đơn giản làm công việc của họ rồi trở về với gia đình.
Và cuối cùng, nhiều khi những nhà phát triển này làm việc trong các doanh nghiệp có hợp đồng bảo vệ IP rất nghiêm ngặt.
Dù phần lớn trong số họ cuối cùng có thể sẽ không trực tiếp tham gia vào mã nguồn mở hay các proposal của cộng đồng Go, khả năng họ dùng Go vẫn phụ thuộc vào cả hai điều đó.

Cộng đồng Go và các proposal của Go cần hiểu nhu cầu của những nhà phát triển ít lên tiếng đó.
Proposal Go có thể tạo ra tác động lớn lên thứ được chấp nhận và sử dụng.
Ví dụ, vendor folder và sau này là proxy của Go modules là cực kỳ quan trọng cho các doanh nghiệp kiểm soát chặt chẽ mã nguồn và thường có ít cuộc trò chuyện trực tiếp hơn với cộng đồng Go.
Chính nhờ có các cơ chế đó mà các tổ chức này mới có thể dùng Go.
Do đó, chúng ta không chỉ phải chú ý đến người dùng Go hiện tại, mà còn đến các nhà phát triển và tổ chức từng cân nhắc Go nhưng đã quyết định không dùng.
Chúng ta cần hiểu những lý do đó.

