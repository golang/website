---
title: "Go: một năm trước vào ngày này"
date: 2010-11-10
by:
- Andrew Gerrand
tags:
- birthday
summary: Chúc mừng sinh nhật lần thứ nhất của Go!
---


Vào ngày 10 tháng 11 năm 2009, chúng tôi ra mắt dự án Go:
một ngôn ngữ lập trình mã nguồn mở tập trung vào sự đơn giản và hiệu quả.
Một năm trôi qua kể từ đó đã chứng kiến rất nhiều bước phát triển
cả trong chính dự án Go lẫn trong cộng đồng của nó.

Chúng tôi bắt đầu với mục tiêu xây dựng một ngôn ngữ cho lập trình hệ thống,
tức là loại chương trình mà người ta thường viết bằng C hoặc C++,
và chúng tôi đã ngạc nhiên trước tính hữu dụng của Go như một ngôn ngữ đa dụng.
Chúng tôi từng dự đoán sẽ nhận được sự quan tâm từ các lập trình viên C, C++ và Java,
nhưng làn sóng quan tâm từ người dùng các ngôn ngữ định kiểu động như
Python và JavaScript lại là điều ngoài dự liệu.
Sự kết hợp giữa biên dịch gốc, kiểu tĩnh, quản lý bộ nhớ
và cú pháp gọn nhẹ của Go dường như đã chạm đúng nhu cầu
của một bộ phận rất rộng trong cộng đồng lập trình.

Bộ phận ấy sau đó lớn dần thành một cộng đồng tận tâm gồm những người viết Go đầy nhiệt huyết.
[Danh sách thư điện tử](http://groups.google.com/group/golang-nuts) của chúng tôi hiện có hơn 3.800 thành viên,
với khoảng 1.500 bài viết mỗi tháng.
Dự án đã có hơn 130 người đóng góp
(tức những người đã gửi mã hoặc tài liệu),
và trong số 2.800 commit kể từ ngày ra mắt, gần một phần ba được đóng góp
bởi các lập trình viên ngoài nhóm cốt lõi.
Để đưa toàn bộ số mã đó vào trạng thái tốt, gần 14.000 email đã được trao đổi trên
[danh sách thư phát triển](http://groups.google.com/group/golang-dev) của chúng tôi.

Những con số đó phản ánh một khối lao động mà thành quả của nó thể hiện rõ trong codebase của dự án.
Các trình biên dịch đã được cải thiện đáng kể,
với việc sinh mã nhanh hơn và hiệu quả hơn,
hơn một trăm lỗi đã được báo cáo đã được sửa,
và hỗ trợ cho ngày càng nhiều hệ điều hành và kiến trúc.
Bản port sang Windows đang tiến rất gần tới hoàn thiện nhờ một nhóm
người đóng góp tận tụy (một trong số đó đã trở thành committer đầu tiên của dự án không thuộc Google).
Bản port ARM cũng đã tiến rất xa,
gần đây đã đạt cột mốc vượt qua toàn bộ bài kiểm thử.

Bộ công cụ Go cũng đã được mở rộng và cải thiện.
Công cụ tài liệu của Go, [godoc](/cmd/godoc/),
giờ đây hỗ trợ tài liệu cho các cây mã nguồn khác
(bạn có thể duyệt và tìm kiếm mã của chính mình)
và cung cấp giao diện ["code walk"](/doc/codewalk/)
để trình bày tài liệu hướng dẫn (cùng rất nhiều cải tiến khác).
[Goinstall](/cmd/goinstall/),
một công cụ quản lý package mới, cho phép người dùng cài đặt và cập nhật
các package bên ngoài chỉ với một lệnh.
[Gofmt](/cmd/gofmt/),
công cụ định dạng mã của Go, giờ đây thực hiện các phép đơn giản hóa cú pháp khi có thể.
Goplay,
một công cụ “biên dịch trong lúc gõ” trên nền web,
là cách thuận tiện để thử nghiệm Go trong những lúc bạn không thể truy cập
[Go Playground](/doc/play/).

Thư viện chuẩn đã tăng thêm hơn 42.000 dòng mã và hiện bao gồm
20 [package](/pkg/) mới.
Trong số đó có các package [jpeg](/pkg/image/jpeg/),
[jsonrpc](/pkg/rpc/jsonrpc/),
[mime](/pkg/mime/), [netchan](/pkg/netchan/),
và [smtp](/pkg/smtp/),
cũng như hàng loạt package [mật mã học](/pkg/crypto/) mới.
Nói rộng hơn, thư viện chuẩn đã liên tục được tinh chỉnh và sửa đổi
khi hiểu biết của chúng tôi về các thành ngữ của Go ngày càng sâu hơn.

Câu chuyện về gỡ lỗi cũng đã tốt hơn.
Những cải tiến gần đây đối với đầu ra DWARF của các trình biên dịch gc giúp
GNU debugger, GDB, trở nên hữu ích cho các tệp nhị phân Go, và chúng tôi đang tích cực làm việc
để phần thông tin gỡ lỗi đó đầy đủ hơn.
(Xem [bài blog gần đây](/blog/debugging-go-code-status-report) để biết chi tiết.)

Giờ đây, việc liên kết với các thư viện hiện có được viết bằng
những ngôn ngữ khác ngoài Go cũng trở nên dễ dàng hơn bao giờ hết.
Hỗ trợ Go đã có trong bản phát hành [SWIG](http://www.swig.org/) mới nhất,
phiên bản 2.0.1, giúp việc liên kết với mã C và C++ dễ dàng hơn,
và công cụ [cgo](/cmd/cgo/) của chúng tôi cũng đã có rất nhiều bản sửa lỗi và cải tiến.

[Gccgo](/doc/install/gccgo),
front end của Go cho GNU C Compiler, vẫn bắt kịp gc compiler
như một hiện thực song song của Go.
Giờ đây nó đã có một bộ gom rác hoạt động tốt, và đã được chấp nhận vào lõi GCC.
Chúng tôi hiện đang hướng tới việc đưa [gofrontend](http://code.google.com/p/gofrontend/)
trở thành một front end của trình biên dịch Go theo giấy phép BSD,
tách rời hoàn toàn khỏi GCC.

Bên ngoài chính dự án Go, Go cũng đang bắt đầu được dùng để xây dựng phần mềm thực tế.
Hiện có hơn 200 chương trình và thư viện Go được liệt kê trên [Project dashboard](http://godashboard.appspot.com/project),
và hàng trăm dự án khác trên [Google Code](http://code.google.com/hosting/search?q=label:Go)
và [GitHub](https://github.com/search?q=language:Go).
Trên danh sách thư và kênh IRC của chúng tôi, bạn có thể gặp các lập trình viên từ khắp nơi trên thế giới
đang sử dụng Go cho các dự án lập trình của họ.
(Xem [bài blog khách mời](/blog/real-go-projects-smarttwitter-and-webgo)
tháng trước của chúng tôi để thấy một ví dụ thực tế.) Bên trong Google cũng có
nhiều nhóm chọn Go để xây dựng phần mềm production,
và chúng tôi cũng đã nhận được báo cáo từ các công ty khác đang phát triển những hệ thống quy mô đáng kể bằng Go.
Chúng tôi cũng đã liên lạc với một số nhà giáo dục đang dùng Go như một ngôn ngữ giảng dạy.

Bản thân ngôn ngữ cũng đã trưởng thành và chín chắn hơn.
Trong năm qua chúng tôi đã nhận được nhiều yêu cầu tính năng.
Nhưng Go là một ngôn ngữ nhỏ, và chúng tôi đã nỗ lực rất nhiều để bảo đảm rằng
mỗi tính năng mới đều đạt được sự cân bằng phù hợp giữa tính đơn giản và tính hữu dụng.
Kể từ khi ra mắt, chúng tôi đã thực hiện một số thay đổi với ngôn ngữ,
trong đó nhiều thay đổi được thúc đẩy bởi phản hồi từ cộng đồng.

- Dấu chấm phẩy giờ đây là tùy chọn trong gần như mọi trường hợp. [spec](/doc/go_spec.html#Semicolons)
- Hai hàm dựng sẵn mới `copy` và `append` khiến việc quản lý slice
  hiệu quả và trực quan hơn.
  [spec](/doc/go_spec.html#Appending_and_copying_slices)
- Có thể bỏ qua cận trên và cận dưới khi tạo sub-slice.
  Điều này có nghĩa `s[:]` là dạng rút gọn của `s[0:len(s)]`.
  [spec](/doc/go_spec.html#Slices)
- Hàm dựng sẵn mới `recover` bổ sung cho `panic` và `defer`
  như một cơ chế xử lý lỗi.
  [blog](/blog/defer-panic-and-recover),
  [spec](/doc/go_spec.html#Handling_panics)
- Các kiểu số phức mới (`complex`,
  `complex64`, và `complex128`) giúp đơn giản hóa một số phép toán.
  [spec](/doc/go_spec.html#Complex_numbers),
  [spec](/doc/go_spec.html#Imaginary_literals)
- Cú pháp composite literal cho phép lược bỏ thông tin kiểu dư thừa
  (ví dụ khi khai báo mảng hai chiều).
  [release.2010-10-27](/doc/devel/release.html#2010-10-27),
  [spec](/doc/go_spec.html#Composite_literals)
- Cú pháp tổng quát cho tham số hàm biến số (`...T`) và cách truyền tiếp chúng
  (`v...`) giờ đã được đặc tả.
  [spec](/doc/go_spec.html#Function_Types),
  [spec](/doc/go_spec.html#Passing_arguments_to_..._parameters),
  [release.2010-09-29](/doc/devel/release.html#2010-09-29)

Go chắc chắn đã sẵn sàng cho production,
nhưng vẫn còn chỗ để cải thiện.
Trọng tâm trước mắt của chúng tôi là làm cho chương trình Go nhanh hơn và
hiệu quả hơn trong bối cảnh các hệ thống hiệu năng cao.
Điều đó có nghĩa là cải thiện bộ gom rác,
tối ưu mã được sinh ra, và cải thiện các thư viện cốt lõi.
Chúng tôi cũng đang khám phá một số bổ sung khác cho hệ thống kiểu
để việc lập trình generic trở nên dễ dàng hơn.
Rất nhiều điều đã diễn ra trong một năm; đó vừa là khoảng thời gian đầy phấn khích vừa rất thỏa mãn.
Chúng tôi hy vọng năm tới sẽ còn nhiều thành quả hơn cả năm vừa qua.

_Nếu bạn vẫn đang định quay trở lại với Go, bây giờ là thời điểm rất tốt để làm điều đó! Hãy xem_
[_Documentation_](/doc/docs.html) _và_ [_Getting Started_](/doc/install.html)
_để biết thêm thông tin, hoặc cứ thoải mái tung hoành trong_ [_Go Playground_](/doc/play/).
