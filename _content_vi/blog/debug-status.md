---
title: Gỡ lỗi mã Go (báo cáo hiện trạng)
date: 2010-11-02
by:
- Luuk van Dijk
tags:
- debug
- gdb
summary: Điều gì hoạt động và điều gì chưa hoạt động khi gỡ lỗi chương trình Go bằng GDB.
---


Khi nói đến gỡ lỗi, không gì hiệu quả hơn vài câu lệnh print được đặt chiến lược
để kiểm tra biến hoặc một lệnh panic đúng chỗ để lấy stack trace.
Tuy nhiên, đôi khi bạn thiếu hoặc sự kiên nhẫn hoặc mã nguồn,
và trong những trường hợp đó một trình gỡ lỗi tốt có thể vô cùng giá trị.
Đó là lý do trong vài bản phát hành gần đây chúng tôi đã cải thiện hỗ trợ
cho GDB, GNU debugger, trong linker gc của Go (6l,
8l).

Trong bản phát hành mới nhất (2010-11-02), linker 6l và 8l phát ra thông tin gỡ lỗi DWARF3
khi ghi tệp nhị phân ELF (Linux,
FreeBSD) hoặc Mach-O (Mac OS X).
Mã DWARF đủ phong phú để cho phép bạn làm những việc sau:

  - nạp một chương trình Go vào GDB phiên bản 7.x,
  - liệt kê mọi tệp mã nguồn Go, C và assembly theo dòng (một số phần của runtime Go được viết bằng C và assembly),
  - đặt breakpoint theo dòng và step qua mã,
  - in stack trace và kiểm tra các stack frame, và
  - tìm địa chỉ và in nội dung của hầu hết các biến.

Vẫn còn một số bất tiện:

  - Mã DWARF được phát ra không thể được đọc bởi GDB phiên bản 6.x đi kèm với Mac OS X.
    Chúng tôi rất sẵn lòng nhận các bản vá để làm đầu ra DWARF tương thích với
    GDB tiêu chuẩn trên OS X,
    nhưng cho đến khi điều đó được sửa, bạn sẽ cần tải xuống,
    biên dịch và cài đặt GDB 7.x để dùng trên OS X.
    Mã nguồn có thể tìm thấy tại [http://sourceware.org/gdb/download/](http://sourceware.org/gdb/download/).
    Do đặc thù của OS X bạn sẽ cần cài đặt tệp nhị phân trên một
    hệ thống tệp cục bộ với `chgrp procmod` và `chmod g+s`.
  - Tên được định danh kèm theo tên gói và,
    vì GDB không hiểu các gói Go, bạn phải tham chiếu từng mục bằng tên đầy đủ của nó.
    Ví dụ, biến tên `v` trong gói `main` phải được tham chiếu
    là `'main.v'`, đặt trong dấu nháy đơn.
    Hệ quả là việc tab completion tên biến và tên hàm không hoạt động.
  - Thông tin lexical scoping có phần bị che mờ.
    Nếu có nhiều biến trùng tên,
    thực thể thứ n sẽ có hậu tố dạng ‘#n’.
    Chúng tôi dự định sửa điều này, nhưng nó sẽ đòi hỏi một số thay đổi trong dữ liệu trao đổi
    giữa trình biên dịch và linker.
  - Biến slice và string được biểu diễn dưới dạng cấu trúc nền
    trong thư viện runtime.
    Chúng sẽ trông giống như `{data = 0x2aaaaab3e320, len = 1, cap = 1}.` Với slice,
    bạn phải dereference con trỏ dữ liệu để kiểm tra các phần tử.

Một số thứ vẫn chưa hoạt động:

  - Biến channel, function, interface và map chưa thể kiểm tra được.
  - Chỉ các biến Go mới được chú thích thông tin kiểu; các biến C của runtime thì không.
  - Tệp nhị phân Windows và ARM không chứa thông tin gỡ lỗi DWARF và vì vậy không thể kiểm tra bằng GDB.

Trong những tháng tới chúng tôi dự định xử lý các vấn đề này,
hoặc bằng cách thay đổi trình biên dịch và linker hoặc bằng cách dùng các tiện ích mở rộng Python của GDB.
Trong thời gian đó, chúng tôi hy vọng các lập trình viên Go sẽ được hưởng lợi từ việc có quyền truy cập tốt hơn
vào công cụ gỡ lỗi quen thuộc này.

P.S. Thông tin DWARF cũng có thể được đọc bởi các công cụ khác ngoài GDB.
Ví dụ, trên Linux bạn có thể dùng nó với profiler toàn hệ thống sysprof.
