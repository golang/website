---
title: Gccgo trong GCC 4.7.1
date: 2012-07-11
by:
- Ian Lance Taylor
tags:
- release
summary: GCC 4.7.1 bổ sung hỗ trợ cho Go 1.
---


Ngôn ngữ Go từ trước đến nay luôn được xác định bởi [đặc tả](/ref/spec),
chứ không phải bởi một hiện thực cụ thể.
Nhóm Go đã viết hai trình biên dịch khác nhau hiện thực đặc tả đó: gc và gccgo.
Việc có hai hiện thực khác nhau giúp bảo đảm rằng đặc tả là đầy đủ và chính xác:
khi các trình biên dịch bất đồng, chúng tôi sửa đặc tả,
và thay đổi một hoặc cả hai trình biên dịch cho phù hợp.
Gc là trình biên dịch gốc, và công cụ `go` mặc định dùng nó.
Gccgo là một hiện thực khác với trọng tâm khác,
và trong bài viết này chúng ta sẽ xem kỹ hơn về nó.

Gccgo được phân phối như một phần của GCC, GNU Compiler Collection.
GCC hỗ trợ nhiều frontend khác nhau cho các ngôn ngữ khác nhau;
gccgo là một frontend Go kết nối với backend GCC.
Frontend Go tách biệt khỏi dự án GCC và được thiết kế để có thể
kết nối với các backend trình biên dịch khác,
nhưng hiện tại chỉ hỗ trợ GCC.

So với gc, gccgo biên dịch chậm hơn nhưng hỗ trợ các tối ưu hóa mạnh hơn,
vì vậy một chương trình bị giới hạn bởi CPU được xây dựng bằng gccgo thường sẽ chạy nhanh hơn.
Mọi tối ưu hóa đã được hiện thực trong GCC qua nhiều năm đều có sẵn,
bao gồm inlining, tối ưu vòng lặp, vectorization,
lập lịch lệnh và nhiều hơn nữa.
Dù không phải lúc nào cũng sinh ra mã tốt hơn,
trong một số trường hợp chương trình biên dịch bằng gccgo có thể chạy nhanh hơn 30%.

Trình biên dịch gc chỉ hỗ trợ những bộ xử lý phổ biến nhất:
x86 (32-bit và 64-bit) và ARM.
Tuy nhiên, gccgo hỗ trợ mọi bộ xử lý mà GCC hỗ trợ.
Không phải tất cả các bộ xử lý đó đều đã được kiểm thử kỹ cho gccgo,
nhưng nhiều bộ xử lý đã được, bao gồm x86 (32-bit và 64-bit),
SPARC, MIPS, PowerPC và thậm chí cả Alpha.
Gccgo cũng đã được kiểm thử trên những hệ điều hành mà gc không hỗ trợ,
đáng chú ý là Solaris.

Gccgo cung cấp thư viện Go chuẩn đầy đủ.
Nhiều tính năng cốt lõi của runtime Go là giống nhau ở cả gccgo và gc,
bao gồm bộ lập lịch goroutine, channel,
bộ cấp phát bộ nhớ và bộ gom rác.
Gccgo hỗ trợ chia nhỏ stack của goroutine như gc,
nhưng hiện tại chỉ trên x86 (32-bit hoặc 64-bit) và chỉ khi dùng gold
linker (trên các bộ xử lý khác,
mỗi goroutine sẽ có một stack lớn, và một chuỗi lời gọi hàm sâu
có thể vượt quá cuối stack và làm chương trình bị crash).

Các bản phân phối gccgo hiện chưa bao gồm một phiên bản của lệnh `go`.
Tuy nhiên, nếu bạn cài lệnh `go` từ một bản phát hành Go tiêu chuẩn,
nó đã hỗ trợ gccgo thông qua tùy chọn `-compiler`:
`go build -compiler gccgo myprog`.
Các công cụ dùng cho lời gọi qua lại giữa Go và C/C++,
`cgo` và SWIG, cũng hỗ trợ gccgo.

Chúng tôi đã đặt frontend Go dưới cùng giấy phép BSD như phần còn lại của các công cụ Go.
Bạn có thể tải mã nguồn của frontend tại
[dự án gofrontend](https://github.com/golang/gofrontend).
Lưu ý rằng khi frontend Go được liên kết với backend GCC để tạo gccgo,
giấy phép GPL của GCC sẽ có hiệu lực ưu tiên.

Bản phát hành GCC mới nhất, 4.7.1, bao gồm gccgo với hỗ trợ cho Go 1.
Nếu bạn cần hiệu năng tốt hơn cho các chương trình Go bị giới hạn bởi CPU,
hoặc cần hỗ trợ những bộ xử lý hay hệ điều hành mà gc không hỗ trợ,
gccgo có thể là câu trả lời.
