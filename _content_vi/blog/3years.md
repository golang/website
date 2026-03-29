---
title: Go tròn ba tuổi
date: 2012-11-10
by:
- Russ Cox
tags:
- community
- birthday
summary: Chúc mừng sinh nhật lần thứ ba của Go!
---


Dự án mã nguồn mở Go hôm nay [tròn ba tuổi](http://google-opensource.blogspot.com/2009/11/hey-ho-lets-go.html).

Thật tuyệt khi nhìn lại chặng đường mà Go đã đi được trong ba năm ấy.
Khi chúng tôi ra mắt, Go là một ý tưởng được hậu thuẫn bởi hai bản hiện thực chạy trên Linux và OS X.
Cú pháp, ngữ nghĩa và các thư viện thay đổi thường xuyên khi chúng tôi phản ứng với phản hồi từ người dùng
và kinh nghiệm thực tế với ngôn ngữ.

Kể từ khi Go được phát hành mã nguồn mở,
chúng tôi đã có thêm
hàng trăm người đóng góp bên ngoài,
những người đã mở rộng và cải thiện Go theo vô số cách khác nhau,
bao gồm cả việc viết một bản port cho Windows từ đầu.
Chúng tôi bổ sung hệ thống quản lý package
[goinstall](https://groups.google.com/d/msg/golang-nuts/8JFwR3ESjjI/cy7qZzN7Lw4J),
về sau trở thành
[lệnh `go`](/cmd/go/).
Chúng tôi cũng bổ sung
[hỗ trợ Go trên App Engine](/blog/go-for-app-engine-is-now-generally).
Trong năm qua, chúng tôi còn thực hiện [nhiều bài nói chuyện](/doc/#talks), tạo ra [tour nhập môn tương tác](/tour/)
và gần đây đã bổ sung hỗ trợ cho [ví dụ thực thi được trong tài liệu package](/pkg/strings/#pkg-examples).

Có lẽ sự phát triển quan trọng nhất trong năm qua
là việc ra mắt phiên bản ổn định đầu tiên,
[Go 1](/blog/go1).
Những người viết chương trình bằng Go 1 giờ đây có thể yên tâm rằng chương trình của họ
sẽ tiếp tục biên dịch và chạy mà không cần thay đổi, trong nhiều môi trường khác nhau,
trên thang thời gian tính bằng nhiều năm.
Trong quá trình chuẩn bị phát hành Go 1, chúng tôi đã dành nhiều tháng để dọn dẹp
[ngôn ngữ và thư viện](/doc/go1.html)
để biến nó thành thứ có thể trường tồn theo thời gian.

Hiện tại chúng tôi đang hướng đến việc phát hành Go 1.1 trong năm 2013.
Sẽ có một số chức năng mới, nhưng bản phát hành đó sẽ chủ yếu tập trung
vào việc giúp Go đạt hiệu năng còn tốt hơn hiện nay.

Điều khiến chúng tôi đặc biệt vui mừng là cộng đồng đã phát triển xung quanh Go:
danh sách thư và các kênh IRC dường như luôn tràn ngập thảo luận,
và năm nay cũng đã có một số cuốn sách về Go được xuất bản. Cộng đồng đang phát triển mạnh.
Việc dùng Go trong môi trường production cũng đã bùng nổ, đặc biệt là kể từ Go 1.

Chúng tôi sử dụng Go tại Google theo nhiều cách khác nhau, phần lớn trong số đó không nhìn thấy từ bên ngoài.
Một vài ví dụ công khai bao gồm
[phục vụ Chrome và các bản tải xuống khác](https://groups.google.com/d/msg/golang-nuts/BNUNbKSypE0/E4qSfpx9qI8J),
[mở rộng cơ sở dữ liệu MySQL tại YouTube](http://code.google.com/p/vitess/),
và tất nhiên là vận hành
[trang chủ Go](/)
trên [App Engine](https://developers.google.com/appengine/docs/go/overview).
[Thanksgiving Doodle](/blog/from-zero-to-go-launching-on-google) của năm ngoái
và trang
[Jam with Chrome](http://www.jamwithchrome.com/technology)
gần đây cũng được phục vụ bởi các chương trình Go.

Các công ty và dự án khác cũng đang sử dụng Go, bao gồm
[BBC Worldwide](http://www.quora.com/Go-programming-language/Is-Google-Go-ready-for-production-use/answer/Kunal-Anand),
[Canonical](http://dave.cheney.net/wp-content/uploads/2012/08/august-go-meetup.pdf),
[CloudFlare](http://blog.cloudflare.com/go-at-cloudflare),
[Heroku](/blog/go-at-heroku),
[Novartis](https://plus.google.com/114945221884326152379/posts/d1SVaqkRyTL),
[SoundCloud](http://backstage.soundcloud.com/2012/07/go-at-soundcloud/),
[SmugMug](http://sorcery.smugmug.com/2012/04/06/deriving-json-types-in-go/),
[StatHat](/blog/building-stathat-with-go),
[Tinkercad](https://tinkercad.com/about/jobs),
và
[nhiều đơn vị khác](/wiki/GoUsers).

Xin chúc cho còn nhiều năm nữa của việc lập trình hiệu quả với Go.
