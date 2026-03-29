---
title: GOMAXPROCS nhận biết container
date: 2025-08-20
by:
- Michael Pratt
- Carlos Amedee
summary: Giá trị mặc định GOMAXPROCS mới trong Go 1.25 cải thiện hành vi trong container.
---

Go 1.25 bao gồm các giá trị mặc định `GOMAXPROCS` mới có nhận biết container, mang lại hành vi mặc định hợp lý hơn cho nhiều tải công việc trong container, tránh hiện tượng throttling có thể ảnh hưởng tới độ trễ tail, và cải thiện mức sẵn sàng cho sản xuất của Go ngay từ khi cài đặt.
Trong bài viết này, chúng ta sẽ đi sâu vào cách Go lập lịch goroutine, cách việc lập lịch đó tương tác với các cơ chế kiểm soát CPU ở cấp container, và cách Go có thể hoạt động tốt hơn khi nhận biết các cơ chế kiểm soát CPU của container.

## `GOMAXPROCS`

Một trong những điểm mạnh của Go là khả năng đồng thời tích hợp sẵn và dễ dùng thông qua goroutine.
Từ góc nhìn ngữ nghĩa, goroutine trông rất giống các luồng của hệ điều hành, cho phép chúng ta viết mã chặn đơn giản.
Mặt khác, goroutine nhẹ hơn luồng hệ điều hành nhiều, khiến việc tạo và hủy chúng tức thời rẻ hơn đáng kể.

Mặc dù một hiện thực Go có thể ánh xạ từng goroutine sang một luồng hệ điều hành riêng, Go giữ cho goroutine nhẹ bằng một bộ lập lịch trong runtime biến các luồng thành tài nguyên có thể thay thế lẫn nhau.
Bất kỳ luồng nào do Go quản lý cũng có thể chạy bất kỳ goroutine nào, nên tạo một goroutine mới không đòi hỏi tạo một luồng mới, và đánh thức một goroutine cũng không nhất thiết phải đánh thức thêm một luồng khác.

Tuy nhiên, đi cùng bộ lập lịch là các câu hỏi về lập lịch.
Ví dụ, chính xác thì ta nên dùng bao nhiêu luồng để chạy goroutine?
Nếu có 1.000 goroutine sẵn sàng chạy, ta có nên lập lịch chúng trên 1.000 luồng khác nhau không?

Đó là lúc [`GOMAXPROCS`](/pkg/runtime#GOMAXPROCS) xuất hiện.
Về mặt ngữ nghĩa, `GOMAXPROCS` cho runtime Go biết “mức song song khả dụng” mà Go nên dùng.
Nói cụ thể hơn, `GOMAXPROCS` là số luồng tối đa dùng để chạy goroutine cùng lúc.

Vì vậy, nếu `GOMAXPROCS=8` và có 1.000 goroutine sẵn sàng chạy, Go sẽ dùng 8 luồng để chạy đồng thời 8 goroutine.
Thông thường, goroutine chạy trong thời gian rất ngắn rồi bị chặn, lúc đó Go sẽ chuyển sang chạy một goroutine khác trên chính luồng đó.
Go cũng sẽ preempt các goroutine không tự chặn, bảo đảm mọi goroutine đều có cơ hội được chạy.

Từ Go 1.5 đến Go 1.24, `GOMAXPROCS` mặc định bằng tổng số lõi CPU trên máy.
Lưu ý rằng trong bài viết này, “lõi” chính xác hơn có nghĩa là “CPU logic”.
Ví dụ, một máy có 4 CPU vật lý có hyperthreading sẽ có 8 CPU logic.

Đây thường là một mặc định tốt cho “mức song song khả dụng” vì nó khớp tự nhiên với mức song song của phần cứng.
Tức là, nếu có 8 lõi mà Go chạy hơn 8 luồng cùng lúc, hệ điều hành sẽ phải multiplex các luồng đó lên 8 lõi, giống như cách Go multiplex goroutine lên luồng.
Lớp lập lịch bổ sung này không phải lúc nào cũng là vấn đề, nhưng nó là chi phí không cần thiết.

## Điều phối container

Một điểm mạnh cốt lõi khác của Go là sự thuận tiện khi triển khai ứng dụng bằng container, và việc quản lý số lõi mà Go sử dụng đặc biệt quan trọng khi triển khai ứng dụng trong nền tảng điều phối container.
Những nền tảng điều phối container như [Kubernetes](https://kubernetes.io/) lấy một tập tài nguyên máy và lập lịch các container trong tài nguyên khả dụng dựa trên tài nguyên được yêu cầu.
Để nhồi được nhiều container nhất có thể vào tài nguyên của cụm, nền tảng cần dự đoán được mức sử dụng tài nguyên của từng container được lập lịch.
Chúng ta muốn Go tuân theo các ràng buộc sử dụng tài nguyên mà nền tảng điều phối container đặt ra.

Hãy xem xét tác động của thiết lập `GOMAXPROCS` trong bối cảnh Kubernetes như một ví dụ.
Các nền tảng như Kubernetes cung cấp cơ chế giới hạn tài nguyên mà một container tiêu thụ.
Kubernetes có khái niệm giới hạn tài nguyên CPU, báo cho hệ điều hành bên dưới biết một container cụ thể hoặc một tập container sẽ được cấp bao nhiêu tài nguyên lõi.
Việc đặt giới hạn CPU sẽ được chuyển thành việc tạo ra giới hạn băng thông CPU của Linux [control group](https://docs.kernel.org/admin-guide/cgroup-v2.html#cpu).

Trước Go 1.25, Go không biết đến các giới hạn CPU do nền tảng điều phối đặt ra.
Thay vào đó, nó đặt `GOMAXPROCS` bằng số lõi trên máy nơi nó được triển khai.
Nếu có giới hạn CPU đang có hiệu lực, ứng dụng có thể cố dùng nhiều CPU hơn rất nhiều so với mức được phép.
Để ngăn ứng dụng vượt quá giới hạn, kernel Linux sẽ [throttle](https://kubernetes.io/docs/concepts/configuration/manage-resources-containers/#how-pods-with-resource-limits-are-run) ứng dụng.

Throttling là một cơ chế thô bạo để hạn chế các container nếu không chúng sẽ vượt quá giới hạn CPU: nó dừng hẳn việc thực thi ứng dụng trong phần thời gian còn lại của chu kỳ throttling.
Chu kỳ throttling thường là 100ms, nên throttling có thể gây tác động rất lớn đến độ trễ tail so với hiệu ứng multiplex lập lịch mềm hơn của một thiết lập `GOMAXPROCS` thấp hơn.
Ngay cả khi ứng dụng chưa bao giờ có nhiều song song, các tác vụ do runtime Go thực hiện, như garbage collection, vẫn có thể gây các đỉnh CPU làm kích hoạt throttling.

## Mặc định mới

Chúng tôi muốn Go cung cấp các mặc định hiệu quả và đáng tin cậy khi có thể, vì vậy trong Go 1.25, chúng tôi đã làm cho `GOMAXPROCS` mặc định tính đến môi trường container của nó.
Nếu một tiến trình Go đang chạy bên trong container có giới hạn CPU, `GOMAXPROCS` sẽ mặc định bằng giới hạn CPU nếu nó nhỏ hơn số lõi.

Các hệ thống điều phối container có thể điều chỉnh giới hạn CPU của container trong lúc chạy, nên Go 1.25 cũng sẽ định kỳ kiểm tra giới hạn CPU và tự động điều chỉnh `GOMAXPROCS` nếu giới hạn thay đổi.

Cả hai mặc định này chỉ áp dụng nếu `GOMAXPROCS` không được chỉ định theo cách khác.
Việc đặt biến môi trường `GOMAXPROCS` hoặc gọi `runtime.GOMAXPROCS` vẫn tiếp tục hoạt động như trước.
Tài liệu của [`runtime.GOMAXPROCS`](/pkg/runtime#GOMAXPROCS) trình bày chi tiết hành vi mới.

## Các mô hình hơi khác nhau

Cả `GOMAXPROCS` và giới hạn CPU của container đều đặt giới hạn lên lượng CPU tối đa mà tiến trình có thể dùng, nhưng mô hình của chúng khác nhau đôi chút.

`GOMAXPROCS` là giới hạn về mức song song.
Nếu `GOMAXPROCS=8` thì Go sẽ không bao giờ chạy quá 8 goroutine cùng một lúc.

Ngược lại, giới hạn CPU là giới hạn về throughput.
Tức là chúng giới hạn tổng thời gian CPU được sử dụng trong một khoảng thời gian thực.
Chu kỳ mặc định là 100ms.
Vì vậy, “giới hạn 8 CPU” thực ra là giới hạn 800ms thời gian CPU cho mỗi 100ms thời gian thực.

Giới hạn này có thể được lấp đầy bằng cách chạy liên tục 8 luồng trong toàn bộ 100ms, tương đương với `GOMAXPROCS=8`.
Mặt khác, giới hạn cũng có thể được lấp đầy bằng cách chạy 16 luồng trong 50ms mỗi luồng, trong khi mỗi luồng rảnh hoặc bị chặn trong 50ms còn lại.

Nói cách khác, giới hạn CPU không giới hạn tổng số CPU mà container có thể chạy trên đó.
Nó chỉ giới hạn tổng thời gian CPU.

Phần lớn ứng dụng có mức sử dụng CPU khá ổn định giữa các chu kỳ 100ms, nên mặc định `GOMAXPROCS` mới khớp khá tốt với giới hạn CPU, và chắc chắn tốt hơn là số lõi tổng!
Tuy nhiên, cũng cần lưu ý rằng những tải công việc đặc biệt giật cục có thể thấy độ trễ tăng lên do thay đổi này vì `GOMAXPROCS` ngăn các đợt bùng phát ngắn của số luồng vượt quá mức trung bình của giới hạn CPU.

Ngoài ra, vì giới hạn CPU là giới hạn throughput nên nó có thể có thành phần thập phân (ví dụ 2.5 CPU).
Trong khi đó, `GOMAXPROCS` phải là một số nguyên dương.
Vì thế, Go phải làm tròn giới hạn đó thành một giá trị `GOMAXPROCS` hợp lệ.
Go luôn làm tròn lên để cho phép sử dụng toàn bộ giới hạn CPU.

## CPU request

Mặc định `GOMAXPROCS` mới của Go dựa trên giới hạn CPU của container, nhưng các hệ thống điều phối container cũng cung cấp cơ chế “CPU request”.
Trong khi giới hạn CPU chỉ định lượng CPU tối đa mà container được phép dùng, CPU request chỉ định lượng CPU tối thiểu luôn được bảo đảm cho container.

Việc tạo container có CPU request nhưng không có giới hạn CPU là khá phổ biến, vì điều đó cho phép container tận dụng tài nguyên CPU của máy vượt quá mức request nếu những tài nguyên đó nhàn rỗi do các container khác thiếu tải.
Thật không may, điều đó cũng có nghĩa là Go không thể đặt `GOMAXPROCS` dựa trên CPU request, bởi như thế sẽ ngăn việc tận dụng thêm tài nguyên nhàn rỗi.

Các container có CPU request vẫn bị [ràng buộc](https://kubernetes.io/docs/concepts/configuration/manage-resources-containers/#how-pods-with-resource-limits-are-run) khi vượt quá request nếu máy đang bận.
Ràng buộc dựa trên trọng số khi vượt request “mềm” hơn việc throttling cứng theo chu kỳ của giới hạn CPU, nhưng các đỉnh CPU do `GOMAXPROCS` cao vẫn có thể ảnh hưởng xấu tới hành vi ứng dụng.

## Tôi có nên đặt giới hạn CPU không?

Chúng ta đã tìm hiểu những vấn đề do `GOMAXPROCS` quá cao gây ra, và rằng việc đặt giới hạn CPU cho container cho phép Go tự động đặt `GOMAXPROCS` phù hợp, vì vậy bước tiếp theo hiển nhiên là tự hỏi liệu mọi container có nên đặt giới hạn CPU hay không.

Dù đó có thể là lời khuyên hay để tự động có được mặc định `GOMAXPROCS` hợp lý, còn có nhiều yếu tố khác cần cân nhắc khi quyết định có đặt giới hạn CPU hay không, chẳng hạn ưu tiên tận dụng tài nguyên nhàn rỗi bằng cách tránh giới hạn so với ưu tiên độ trễ có thể dự đoán được bằng cách đặt giới hạn.

Những hành vi tệ nhất do không khớp giữa `GOMAXPROCS` và giới hạn CPU hiệu dụng xảy ra khi `GOMAXPROCS` cao hơn đáng kể so với giới hạn CPU hiệu dụng.
Ví dụ, một container nhỏ được cấp 2 CPU nhưng chạy trên máy 128 lõi.
Đây là những trường hợp có giá trị nhất để cân nhắc đặt giới hạn CPU tường minh, hoặc thay vào đó đặt `GOMAXPROCS` một cách tường minh.

## Kết luận

Go 1.25 mang lại hành vi mặc định hợp lý hơn cho nhiều tải công việc container bằng cách đặt `GOMAXPROCS` dựa trên giới hạn CPU của container.
Điều này tránh hiện tượng throttling có thể ảnh hưởng tới độ trễ tail, cải thiện hiệu quả, và nói chung cố gắng bảo đảm rằng Go sẵn sàng cho môi trường sản xuất ngay từ đầu.
Bạn có thể nhận được các mặc định mới chỉ bằng cách đặt phiên bản Go là 1.25.0 hoặc cao hơn trong `go.mod` của mình.

Xin cảm ơn mọi người trong cộng đồng đã đóng góp vào các [cuộc thảo luận](/issue/33803) [kéo dài](/issue/73193) giúp biến điều này thành hiện thực, và đặc biệt là phản hồi từ những người bảo trì [`go.uber.org/automaxprocs`](https://pkg.go.dev/go.uber.org/automaxprocs) tại Uber, gói đã cung cấp từ lâu hành vi tương tự cho người dùng của họ.

