---
title: Gỡ lỗi những gì bạn triển khai trong Go 1.12
date: 2019-03-21
by:
- David Chase
tags:
- debug
- technical
summary: Go 1.12 cải thiện khả năng hỗ trợ gỡ lỗi các tệp nhị phân đã tối ưu hóa.
template: true
---

## Giới thiệu

Go 1.11 và Go 1.12 đánh dấu bước tiến đáng kể hướng tới việc cho phép
lập trình viên gỡ lỗi chính những tệp nhị phân đã tối ưu hóa mà họ triển khai lên môi trường production.

Khi trình biên dịch Go ngày càng quyết liệt hơn trong việc tạo ra các tệp nhị phân nhanh hơn,
chúng ta đã đánh đổi khả năng gỡ lỗi.
Trong Go 1.10, người dùng cần tắt hoàn toàn các tối ưu hóa để có
trải nghiệm gỡ lỗi tốt với các công cụ tương tác như Delve.
Nhưng người dùng không nên phải đánh đổi hiệu năng để lấy khả năng gỡ lỗi,
đặc biệt khi vận hành các dịch vụ production.
Nếu vấn đề xảy ra trong production,
bạn cần gỡ lỗi ngay trong production, và điều đó không nên đòi hỏi phải triển khai
các tệp nhị phân không được tối ưu hóa.

Trong Go 1.11 và 1.12, chúng tôi tập trung cải thiện trải nghiệm gỡ lỗi trên
các tệp nhị phân đã tối ưu hóa (thiết lập mặc định của trình biên dịch Go).
Các cải tiến gồm có

  - Kiểm tra giá trị chính xác hơn, đặc biệt với đối số tại thời điểm vào hàm;
  - Xác định ranh giới câu lệnh chính xác hơn để việc step bớt giật cục
    và điểm dừng thường rơi đúng chỗ lập trình viên mong đợi hơn;
  - Và hỗ trợ sơ bộ để Delve gọi các hàm Go (goroutine và
    garbage collection khiến việc này khó hơn trong C và C++).

## Gỡ lỗi mã đã tối ưu hóa với Delve

[Delve](https://github.com/go-delve/delve) là trình gỡ lỗi cho Go trên x86,
hỗ trợ cả Linux và macOS.
Delve hiểu goroutine và các đặc điểm khác của Go, đồng thời mang lại một trong những
trải nghiệm gỡ lỗi Go tốt nhất.
Delve cũng là bộ máy gỡ lỗi phía sau [GoLand](https://www.jetbrains.com/go/),
[VS Code](https://code.visualstudio.com/),
và [Vim](https://github.com/fatih/vim-go).

Thông thường Delve sẽ biên dịch lại mã đang gỡ lỗi với `-gcflags "all=-N -l"`,
tham số này tắt inline và hầu hết các tối ưu hóa.
Để gỡ lỗi mã đã tối ưu hóa bằng Delve, trước tiên hãy biên dịch tệp nhị phân đã tối ưu hóa,
sau đó dùng `dlv exec your_program` để gỡ lỗi.
Hoặc, nếu bạn có tệp core từ một lần crash,
bạn có thể kiểm tra nó bằng `dlv core your_program your_core`.
Với Go 1.12 và các bản Delve mới nhất, bạn sẽ có thể xem được nhiều biến,
ngay cả trong các tệp nhị phân đã tối ưu hóa.

## Cải thiện việc kiểm tra giá trị

Khi gỡ lỗi các tệp nhị phân đã tối ưu hóa do Go 1.10 tạo ra,
giá trị biến thường hoàn toàn không truy cập được.
Ngược lại, bắt đầu từ Go 1.11, các biến thường có thể được xem
ngay cả trong các tệp nhị phân đã tối ưu hóa,
trừ khi chúng đã bị tối ưu loại bỏ hoàn toàn.
Trong Go 1.11, trình biên dịch bắt đầu phát ra danh sách vị trí DWARF để trình gỡ lỗi
có thể theo dõi biến khi chúng di chuyển ra vào các thanh ghi và tái dựng
các đối tượng phức tạp bị tách ra trên nhiều thanh ghi và vị trí trong stack.

## Cải thiện việc step

Đây là ví dụ về việc step qua một hàm đơn giản trong trình gỡ lỗi ở phiên bản 1.10,
với các lỗi (bỏ qua và lặp lại dòng) được tô nổi bật bằng các mũi tên đỏ.

{{image "debug-opt/stepping.svg" 450}}

Những lỗi như vậy khiến bạn dễ mất dấu vị trí hiện tại khi step
qua chương trình và làm cản trở việc dừng đúng breakpoint.

Go 1.11 và 1.12 ghi lại thông tin ranh giới câu lệnh và làm tốt hơn
trong việc theo dõi số dòng mã nguồn qua các tối ưu hóa và inline.
Kết quả là trong Go 1.12, khi step qua đoạn mã này, trình gỡ lỗi dừng ở mọi dòng
và theo đúng thứ tự bạn mong đợi.

## Gọi hàm

Hỗ trợ gọi hàm trong Delve vẫn đang được phát triển, nhưng những trường hợp đơn giản đã hoạt động. Ví dụ:

	(dlv) call fib(6)
	> main.main() ./hello.go:15 (PC: 0x49d648)
	Values returned:
		~r1: 8

## Hướng phát triển tiếp theo

Go 1.12 là một bước tiến hướng tới trải nghiệm gỡ lỗi tốt hơn cho các tệp nhị phân đã tối ưu hóa
và chúng tôi có kế hoạch cải thiện xa hơn nữa.

Có những đánh đổi căn bản giữa khả năng gỡ lỗi và hiệu năng,
vì vậy chúng tôi đang tập trung vào các lỗi gỡ lỗi có mức ưu tiên cao nhất,
đồng thời xây dựng các chỉ số tự động để theo dõi tiến độ và phát hiện hồi quy.

Chúng tôi tập trung vào việc tạo ra thông tin chính xác cho trình gỡ lỗi về vị trí biến,
để nếu một biến có thể được in ra, nó sẽ được in đúng.
Chúng tôi cũng đang xem xét cách làm cho giá trị biến xuất hiện thường xuyên hơn,
đặc biệt tại những điểm quan trọng như vị trí gọi hàm,
mặc dù trong nhiều trường hợp, cải thiện điều này sẽ đòi hỏi phải làm chậm việc thực thi chương trình.
Cuối cùng, chúng tôi đang cải thiện việc step:
chúng tôi tập trung vào thứ tự step khi có panic,
thứ tự step quanh các vòng lặp, và nói chung cố gắng bám theo thứ tự
mã nguồn khi có thể.

## Ghi chú về hỗ trợ macOS

Go 1.11 bắt đầu nén thông tin gỡ lỗi để giảm kích thước tệp nhị phân.
Delve hỗ trợ điều này nguyên bản, nhưng cả LLDB lẫn GDB đều không hỗ trợ
thông tin gỡ lỗi đã nén trên macOS.
Nếu bạn đang dùng LLDB hoặc GDB, có hai cách khắc phục:
biên dịch tệp nhị phân với `-ldflags=-compressdwarf=false`,
hoặc dùng [splitdwarf](https://godoc.org/golang.org/x/tools/cmd/splitdwarf)
(`go get golang.org/x/tools/cmd/splitdwarf`) để giải nén thông tin gỡ lỗi
trong một tệp nhị phân hiện có.
