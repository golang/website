---
title: Bốn năm của Go
date: 2013-11-10
by:
- Andrew Gerrand
tags:
- community
- birthday
summary: Chúc mừng sinh nhật lần thứ tư của Go!
template: true
---


Hôm nay đánh dấu kỷ niệm bốn năm của Go với tư cách là một dự án mã nguồn mở.

{{image "4years/4years-gopher.png"}}

Thay vì nói về những tiến bộ kỹ thuật của chúng tôi (sẽ còn rất nhiều điều để nói
khi chúng tôi phát hành Go 1.2 trong vài tuần tới), chúng tôi nghĩ nhân dịp này
sẽ phù hợp hơn để nhìn lại cách cộng đồng Go đã phát triển.

Hãy bắt đầu bằng một biểu đồ:

{{image "4years/4years-graph.png"}}

Biểu đồ này cho thấy mức tăng trưởng của số lượt tìm kiếm trên Google cho từ khóa
"[golang](http://www.google.com/trends/explore?hl=en-US#q=golang&date=10/2009+50m&cmpt=q)"
trong bốn năm qua.
Hãy chú ý đến điểm gấp của đường cong vào khoảng tháng 3 năm 2012, khi Go 1.0 được phát hành.
Nếu những lượt tìm kiếm này là một đại diện tương đối hợp lý cho mức độ quan tâm,
thì có thể thấy rõ rằng
sự quan tâm đến Go đã tăng lên đáng kể kể từ khi ra mắt, đặc biệt là trong
2 năm gần đây.

Nhưng sự quan tâm đó đến từ đâu?

Cộng đồng mã nguồn mở đã đón nhận Go,
với wiki cộng đồng của chúng tôi liệt kê [hàng trăm dự án Go](/wiki/Projects). Một vài dự án nổi bật:

- [Docker](http://docker.io) là một công cụ để đóng gói và chạy ứng dụng
  trong các container nhẹ.
  Docker giúp việc cô lập, đóng gói
  và triển khai ứng dụng trở nên dễ dàng, và được các quản trị viên hệ thống rất yêu thích.
  Người tạo ra nó, Solomon Hykes, cho biết thư viện chuẩn của Go,
  các primitive đồng thời, và khả năng triển khai dễ dàng là những yếu tố then chốt,
  và nói rằng “Nói đơn giản, nếu Docker không được viết bằng Go,
  thì nó đã không thể thành công đến vậy.”

- [Packer](http://packer.io) là công cụ tự động hóa việc tạo
  machine image để triển khai lên máy ảo hoặc dịch vụ đám mây.
  Tác giả của nó, Mitchell Hashimoto, hiện cũng đang phát triển một dự án Go khác,
  [serf](http://www.serfdom.io/), một dịch vụ khám phá phi tập trung.
  Giống như Docker, các dự án này giúp việc quản lý các dịch vụ quy mô lớn,
  vận hành theo cụm, trở nên dễ dàng hơn.

- [Bitly](http://bit.ly) với [NSQ](http://bitly.github.io/nsq/) là một nền tảng nhắn tin phân tán theo thời gian thực
  được thiết kế cho khả năng chịu lỗi và tính sẵn sàng cao,
  đồng thời đang được dùng trong production tại bitly và nhiều công ty khác.

- Hệ thống tự động hóa hạ tầng [JuJu](https://juju.ubuntu.com/) của [Canonical](http://canonical.com/)
  đã được viết lại bằng Go.
  Trưởng dự án Gustavo Niemeyer nói rằng “Không phải một khía cạnh đơn lẻ nào của Go
  khiến nó trở thành lựa chọn hấp dẫn,
  mà chính là cách tổ chức cẩn trọng của những mảnh ghép nhỏ được chế tác tốt.”

- Package [raft](https://github.com/goraft/raft) cung cấp một bản hiện thực
  của giao thức đồng thuận phân tán [Raft](https://ramcloud.stanford.edu/wiki/download/attachments/11370504/raft.pdf).
  Đây là nền tảng của các dự án Go như [etcd](https://github.com/coreos/etcd)
  và [SkyDNS](https://github.com/skynetservices/skydns).

- Các dự án phổ biến khác bao gồm [biogo](https://github.com/biogo/biogo),
  [Gorilla Web Toolkit](http://www.gorillatoolkit.org/),
  [groupcache](https://github.com/golang/groupcache),
  [heka](https://github.com/mozilla-services/heka) của Mozilla,
  các hệ thống lưu trữ gọn nhẹ [kv](https://github.com/cznic/kv) và [ql](https://github.com/cznic/ql),
  cùng cơ sở dữ liệu hành vi [Sky](http://skydb.io/).

Nhưng đó mới chỉ là phần nổi của tảng băng.
Số lượng dự án Go mã nguồn mở chất lượng cao là điều đáng kinh ngạc.
Người viết Go năng suất cao [Keith Rarick](http://xph.us/software/) đã nói rất đúng:
“Tình trạng của hệ sinh thái Go chỉ sau bốn năm thật đáng kinh ngạc.
Hãy so Go năm 2013 với Python năm 1995, hoặc Java năm 1999. Hay C++ năm 1987!”

Doanh nghiệp cũng đang tận hưởng lợi ích từ Go. [Trang wiki Go Users](/wiki/GoUsers)
liệt kê hàng chục câu chuyện thành công (và nếu bạn dùng Go,
hãy thêm tên mình vào đó). Một vài ví dụ:

- [CloudFlare](https://blog.cloudflare.com/go-at-cloudflare) đã xây dựng
  toàn bộ dịch vụ DNS phân tán của họ bằng Go,
  và hiện đang trong quá trình chuyển hạ tầng ghi log hàng gigabyte mỗi phút sang ngôn ngữ này.
  Lập trình viên John Graham-Cumming nói rằng “Chúng tôi thấy Go là sự kết hợp hoàn hảo cho nhu cầu của mình:
  sự kết hợp giữa cú pháp quen thuộc, hệ thống kiểu mạnh,
  thư viện mạng mạnh mẽ và tính đồng thời tích hợp khiến ngày càng nhiều
  dự án ở đây được xây dựng bằng Go.”

- [SoundCloud](http://soundcloud.com) là một dịch vụ phân phối âm thanh
  có “hàng chục [hệ thống viết bằng Go](http://backstage.soundcloud.com/2012/07/go-at-soundcloud/),
  chạm tới gần như mọi phần của website, và trong nhiều trường hợp còn vận hành tính năng
  từ đầu đến cuối.” Kỹ sư Peter Bourgon nói rằng “Go chứng minh rằng
  những gánh nặng rườm rà của các ngôn ngữ và hệ sinh thái khác, những thứ mà lập trình viên
  đã quen phải chịu đựng, thường là trong bực tức, thực ra không phải là phần tất yếu của lập trình hiện đại.
  Với Go, tôi có một mối quan hệ thẳng thắn và không đối đầu với công cụ của mình,
  từ giai đoạn phát triển cho đến production.”

- Dịch vụ [ngrok](https://ngrok.com/) cho phép các nhà phát triển web
  cung cấp truy cập từ xa vào môi trường phát triển của họ.
  Tác giả của nó, Alan Shreve, nói rằng “thành công của ngrok như một dự án
  phần lớn là nhờ việc chọn Go làm ngôn ngữ hiện thực,” viện dẫn thư viện HTTP của Go,
  hiệu quả, tính tương thích đa nền tảng,
  và khả năng triển khai dễ dàng là những lợi ích chính.

- [Poptip](http://poptip.com) cung cấp dịch vụ phân tích mạng xã hội,
  và kỹ sư sản phẩm Andy Bonventre nói rằng “Điều bắt đầu như một thử nghiệm
  viết một dịch vụ đơn lẻ bằng Go cuối cùng đã biến thành việc chuyển gần như toàn bộ hạ tầng của chúng tôi sang nó.
  Điều tôi yêu thích nhất ở Go không hẳn là những tính năng của ngôn ngữ,
  mà là sự tập trung vào công cụ, kiểm thử, và những yếu tố khác khiến việc viết
  các ứng dụng lớn trở nên dễ quản lý hơn nhiều.”

- Startup cộng tác âm nhạc [Splice](http://splice.com) đã chọn xây dựng
  dịch vụ của họ bằng Go.
  Đồng sáng lập Matt Aimonetti nói rằng “Chúng tôi đã nghiên cứu nghiêm túc và cân nhắc nhiều
  ngôn ngữ lập trình,
  nhưng sự đơn giản, hiệu quả, triết lý và cộng đồng của Go đã chinh phục chúng tôi.”

- Và tất nhiên, các nhóm kỹ sư khắp Google cũng đang chuyển sang Go.
  Kỹ sư Matt Welsh gần đây đã [chia sẻ trải nghiệm](http://matt-welsh.blogspot.com.au/2013/08/rewriting-large-production-system-in-go.html)
  của mình khi viết lại một dịch vụ production lớn bằng Go.
  Những ví dụ công khai đáng chú ý khác bao gồm [vitess](https://github.com/youtube/vitess) của YouTube
  và [dl.google.com](/talks/2013/oscon-dl.slide).
  Chúng tôi hy vọng sẽ sớm chia sẻ thêm nhiều câu chuyện như vậy.

Vào tháng 9 năm 2012, CEO Derek Collison của [Apcera](http://apcera.com/) đã [dự đoán](https://twitter.com/derekcollison/status/245522124666716160)
rằng “Go sẽ trở thành ngôn ngữ thống trị cho công việc hệ thống trong [Infrastructure-as-a-Service],
Orchestration, và [Platform-as-a-Service] trong vòng 24 tháng.” Nhìn vào danh sách ở trên,
thật dễ để tin vào dự đoán đó.

Vậy làm thế nào để bạn tham gia? Dù bạn là một lập trình viên Go dày dạn kinh nghiệm
hay chỉ mới tò mò về Go,
vẫn có rất nhiều cách để bắt đầu với cộng đồng Go:

- [Tham gia Go User Group gần bạn nhất](/blog/getthee-to-go-meetup),
  nơi các gopher địa phương gặp nhau để chia sẻ kiến thức và kinh nghiệm.
  Những nhóm này đang mọc lên ở khắp nơi trên thế giới.
  Cá nhân tôi đã từng nói chuyện tại các nhóm Go ở Amsterdam,
  Berlin, Gothenburg, London, Moscow, Munich,
  New York City, Paris, San Francisco, Seoul,
  Stockholm, Sydney, Tokyo và Warsaw;
  nhưng còn [nhiều nhóm khác nữa](/wiki/GoUserGroups)!

- Tạo hoặc đóng góp cho một dự án Go mã nguồn mở (hoặc [cho chính Go](/doc/contribute.html)).
  (Và nếu bạn đang xây dựng thứ gì đó, chúng tôi rất muốn được nghe từ bạn trên [danh sách thư Go](http://groups.google.com/group/golang-nuts).)

- Nếu bạn ở châu Âu vào tháng 2 năm 2014, hãy tới [Go Devroom](https://code.google.com/p/go-wiki/wiki/Fosdem2014)
  tại [FOSDEM 2014](https://fosdem.org/2014/).

- Tham dự [GopherCon](http://gophercon.com),
  hội nghị lớn đầu tiên về Go, tại Denver vào tháng 4 năm 2014.
  Sự kiện được tổ chức bởi [Gopher Academy](http://www.gopheracademy.com),
  nơi cũng vận hành [bảng việc làm Go](http://www.gopheracademy.com/jobs).

Nhóm Go đã vô cùng kinh ngạc trước sự phát triển của cộng đồng Go trong bốn năm qua.
Chúng tôi rất vui khi thấy quá nhiều điều tuyệt vời đang được xây dựng bằng Go,
và vô cùng biết ơn khi được làm việc cùng những người đóng góp tuyệt vời và tận tâm của chúng tôi.
Xin cảm ơn tất cả mọi người.

Xin chúc cho bốn năm tiếp theo nữa!
