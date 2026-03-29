Thư mục này chứa các bài viết blog của Go, ở định dạng *.article.
Xem https://pkg.go.dev/golang.org/x/tools/present?tab=doc
để biết tài liệu về định dạng tệp hoặc xem bất kỳ
bài viết nào để có ví dụ.

Tên tệp bài viết nên ngắn, thuận tiện để gõ tay trong URL.
Thông thường tên tệp bài viết có tối đa ba từ, được phân tách bằng dấu gạch ngang.
Một năm ở cuối tên thường không có dấu gạch ngang đứng trước.

Mọi tệp hỗ trợ cho một bài viết, ngay cả khi chỉ có một tệp, nên được
đặt trong một thư mục được đặt tên theo bài viết (bỏ phần hậu tố .article).

Nếu bài viết của bạn có mã được dự định là một chương trình chạy được, vui lòng dùng
.code, hoặc tốt hơn là .play, để nạp các dòng từ một tệp .go hỗ trợ.
Bằng cách đó bạn có thể dễ dàng kiểm tra rằng tệp .go, và do đó là đoạn mã
trong bài viết của bạn, vẫn hoạt động.

Vui lòng dùng .image và .video để nhúng hình ảnh và video,
thay vì dùng các thẻ HTML thô. Các lệnh .image và .video
cung cấp một cách để điều chỉnh phần triển khai của các phần nhúng đó
ở cùng một nơi.

Vui lòng dùng .html khi cần thêm các khối HTML lớn,
giữ phần văn bản của bài viết trong tệp .article chính. Một cách dùng quan trọng khác của
.html là tách riêng một đoạn HTML xuất hiện nhiều lần.
Các chuỗi HTML ngắn, như <div><center> hoặc </div></center>,
thì có thể đặt trực tiếp trong các tệp bài viết.
