---
title: Gỡ lỗi chương trình Go với GNU Debugger
date: 2011-10-30
by:
- Andrew Gerrand
tags:
- debug
- gdb
- technical
summary: Giới thiệu một bài viết mới về cách gỡ lỗi chương trình Go với GDB.
---


Năm ngoái chúng tôi đã [đưa tin](/blog/debugging-go-code-status-report)
rằng bộ công cụ [gc](/cmd/gc/)/[ld](/cmd/6l/)
của Go tạo ra thông tin gỡ lỗi DWARFv3 mà GNU Debugger (GDB) có thể đọc được.
Kể từ đó, công việc cải thiện khả năng hỗ trợ gỡ lỗi mã Go bằng GDB vẫn tiếp tục đều đặn.
Trong số các cải tiến có khả năng kiểm tra goroutine và in ra
các kiểu dữ liệu Go gốc,
bao gồm struct, slice, string, map,
interface và channel.

Để tìm hiểu thêm về Go và GDB, hãy xem bài viết [Gỡ lỗi với GDB](/doc/debugging_with_gdb.html).
