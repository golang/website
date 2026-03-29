---
title: Viết ứng dụng App Engine có khả năng mở rộng
date: 2011-11-01
by:
- David Symonds
tags:
- appengine
- optimization
summary: Cách xây dựng ứng dụng web có khả năng mở rộng bằng Go với Google App Engine.
---


Trở lại hồi tháng 5, chúng tôi đã [công bố](/blog/go-and-google-app-engine)
môi trường chạy Go cho App Engine.
Kể từ đó, chúng tôi đã mở nó cho mọi người sử dụng,
thêm nhiều API mới và cải thiện hiệu năng.
Chúng tôi rất vui khi thấy tất cả những cách thú vị mà mọi người đang dùng Go trên App Engine.
Một trong những lợi ích chính của môi trường chạy Go,
ngoài việc được làm việc với một ngôn ngữ tuyệt vời,
là hiệu năng cao.
Ứng dụng Go biên dịch ra mã máy gốc, không có trình thông dịch hay máy ảo
nằm giữa chương trình của bạn và máy tính.

Làm cho ứng dụng web của bạn chạy nhanh là điều quan trọng vì ai cũng biết rằng
độ trễ của một website có ảnh hưởng đo được đến sự hài lòng của người dùng,
và [Google web search dùng nó như một yếu tố xếp hạng](https://googlewebmastercentral.blogspot.com/2010/04/using-site-speed-in-web-search-ranking.html).
Được công bố hồi tháng 5 còn có việc App Engine sẽ [rời trạng thái Preview](http://googleappengine.blogspot.com/2011/05/year-ahead-for-google-app-engine.html)
và chuyển sang [mô hình định giá mới](https://www.google.com/enterprise/cloud/appengine/pricing.html),
tạo thêm một lý do nữa để viết các ứng dụng App Engine hiệu quả.

Để giúp các nhà phát triển Go dùng App Engine dễ dàng viết các ứng dụng
hiệu quả và có khả năng mở rộng cao, gần đây chúng tôi đã cập nhật một số bài viết App Engine hiện có
để thêm các đoạn mã nguồn Go và liên kết tới tài liệu Go liên quan.

  - [Best practices for writing scalable applications](http://code.google.com/appengine/articles/scaling/overview.html)

  - [Managing Your App's Resource Usage](http://code.google.com/appengine/articles/managing-resources.html)
