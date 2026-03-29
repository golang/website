---
title: Dấu vết thực thi Go mạnh hơn
date: 2024-03-14
by:
- Michael Knyszek
tags:
- debug
- technical
- tracing
summary: "Các tính năng và cải tiến mới cho execution trace trong năm qua."
template: true
---

Gói [runtime/trace](/pkg/runtime/trace) chứa một công cụ mạnh mẽ để hiểu và
xử lý sự cố trong các chương trình Go.
Chức năng bên trong cho phép tạo ra một dấu vết thực thi của từng goroutine trong một
khoảng thời gian.
Với [lệnh `go tool trace`](/pkg/cmd/trace) (hoặc công cụ mã nguồn mở tuyệt vời
[gotraceui](https://gotraceui.dev/)), bạn có thể trực quan hóa và khám phá dữ liệu trong
những execution trace này.

Điều kỳ diệu của trace là nó có thể dễ dàng hé lộ những thứ trong chương trình vốn rất khó nhìn thấy
bằng cách khác.
Ví dụ, một nút thắt cổ chai về đồng thời khi rất nhiều goroutine cùng chặn trên một channel
có thể khá khó nhìn thấy trong CPU profile, vì không có phần thực thi nào để lấy mẫu.
Nhưng trong execution trace, sự _thiếu vắng_ thực thi sẽ hiện ra rất rõ,
và stack trace của các goroutine đang bị chặn sẽ nhanh chóng chỉ ra thủ phạm.

{{image "execution-traces-2024/gotooltrace.png"}}

Lập trình viên Go thậm chí còn có thể instrument chính chương trình của mình bằng [task](/pkg/runtime/trace#Task),
[region](/pkg/runtime/trace#WithRegion), và [log](/pkg/runtime/trace#Log) để
liên hệ các mối quan tâm cấp cao hơn với chi tiết thực thi cấp thấp hơn.

## Vấn đề

Thật không may, lượng thông tin phong phú trong execution trace thường khó tiếp cận.
Bốn vấn đề lớn trong lịch sử đã cản trở việc dùng trace:

- Trace có overhead cao.
- Trace không scale tốt và có thể trở nên quá lớn để phân tích.
- Thường không rõ khi nào nên bắt đầu trace để chụp được một hành vi xấu cụ thể.
- Chỉ những gopher gan dạ nhất mới có thể phân tích trace bằng mã, vì không có
  gói công khai để parse và diễn giải execution trace.

Nếu bạn đã dùng trace trong vài năm gần đây, có lẽ bạn đã bực bội vì một hoặc nhiều
vấn đề trong số này.
Nhưng chúng tôi rất hào hứng chia sẻ rằng trong hai bản phát hành Go gần đây nhất, chúng tôi đã đạt tiến triển lớn ở cả bốn mặt.

## Trace overhead thấp

Trước Go 1.21, overhead khi chạy trace thường nằm đâu đó trong khoảng 10–20% CPU cho nhiều
ứng dụng, điều này khiến trace chỉ phù hợp cho việc dùng theo tình huống thay vì dùng liên tục như CPU
profiling.
Hóa ra phần lớn chi phí của trace đến từ traceback.
Nhiều sự kiện do runtime tạo ra có kèm stack trace, vốn vô cùng quý giá để thực sự
xác định các goroutine đang làm gì tại những thời điểm quan trọng.

Nhờ công trình của Felix Geisendörfer và Nick Ripley trong việc tối ưu hiệu quả của traceback,
overhead CPU lúc chạy của execution trace đã giảm mạnh, xuống còn 1–2% với nhiều
ứng dụng.
Bạn có thể đọc thêm về công việc này trong [bài viết rất hay của Felix](https://blog.felixge.de/reducing-gos-execution-tracer-overhead-with-frame-pointer-unwinding/)
về chủ đề đó.

## Trace có khả năng mở rộng

Định dạng trace và các sự kiện của nó được thiết kế xoay quanh việc phát ra tương đối hiệu quả, nhưng
công cụ lại phải parse và giữ trạng thái của toàn bộ trace.
Một trace vài trăm MiB có thể cần đến vài GiB RAM để phân tích!

Không may, vấn đề này mang tính nền tảng do cách trace được tạo ra.
Để giữ overhead runtime thấp, mọi sự kiện đều được ghi vào thứ tương đương với bộ đệm cục bộ theo luồng.
Nhưng điều đó có nghĩa các sự kiện xuất hiện sai thứ tự thật,
và gánh nặng tìm ra chuyện gì thực sự xảy ra được đẩy cho công cụ trace.

Ý tưởng then chốt để giúp trace scale mà vẫn giữ overhead thấp là thỉnh thoảng chia tách
trace đang được tạo ra.
Mỗi điểm tách sẽ hoạt động hơi giống như đồng thời tắt rồi bật lại trace
trong một bước.
Toàn bộ dữ liệu trace cho đến lúc đó sẽ đại diện cho một trace hoàn chỉnh và tự chứa,
trong khi dữ liệu trace mới tiếp tục trơn tru từ nơi nó dừng lại.

Như bạn có thể hình dung, để sửa điều này cần [xem xét lại và viết lại phần nền tảng của
hiện thực trace](/issue/60773) trong runtime.
Chúng tôi vui mừng cho biết công việc này đã có mặt trong Go 1.22 và hiện đã sẵn sàng dùng rộng rãi.
[Nhiều cải tiến hữu ích](/doc/go1.22#runtime/trace) đi kèm với lần viết lại này, bao gồm cả một số
cải tiến cho [`go tool trace`](/doc/go1.22#trace).
Nếu bạn tò mò, các chi tiết chuyên sâu đều nằm trong [tài liệu thiết kế](https://github.com/golang/proposal/blob/master/design/60773-execution-tracer-overhaul.md).

(Lưu ý: `go tool trace` vẫn nạp toàn bộ trace vào bộ nhớ, nhưng [loại bỏ giới hạn này](/issue/65315)
đối với trace do chương trình Go 1.22+ tạo ra giờ đây đã khả thi.)

## Ghi trace kiểu flight recorder

Giả sử bạn làm việc trên một dịch vụ web và một RPC mất rất lâu.
Bạn không thể bắt đầu trace tại thời điểm bạn biết RPC đã chạy quá lâu, vì
nguyên nhân gốc của request chậm đã xảy ra trước đó và không được ghi lại.

Có một kỹ thuật có thể giúp trong trường hợp này gọi là flight recording, có thể bạn đã quen
từ các môi trường lập trình khác.
Ý tưởng ở đây là luôn bật trace liên tục và luôn giữ lại dữ liệu trace gần đây nhất,
phòng khi cần dùng.
Sau đó, khi có điều gì đó đáng chú ý xảy ra, chương trình chỉ việc ghi ra những gì nó đang có.

Trước khi trace có thể được tách nhỏ, cách này gần như không khả thi.
Nhưng vì trace liên tục nay đã khả thi nhờ overhead thấp, và runtime
giờ có thể chia trace bất cứ lúc nào cần, hóa ra việc hiện thực flight
recording lại khá thẳng thắn.

Vì vậy, chúng tôi vui mừng công bố một thử nghiệm flight recorder, có trong
[gói golang.org/x/exp/trace](/pkg/golang.org/x/exp/trace#FlightRecorder).

Hãy thử nó!
Bên dưới là ví dụ thiết lập flight recording để chụp một HTTP request kéo dài nhằm giúp bạn bắt đầu.

{{code "execution-traces-2024/flightrecorder.go" `/START/` `/END/`}}

Nếu bạn có bất kỳ phản hồi nào, dù tích cực hay tiêu cực, xin hãy chia sẻ trong [issue đề xuất](/issue/63185)!

## API đọc trace

Cùng với việc viết lại hiện thực trace là một nỗ lực làm sạch các phần nội bộ khác của hệ trace,
như `go tool trace`.
Từ đó nảy sinh nỗ lực tạo ra một API đọc trace đủ tốt để chia sẻ và
giúp trace trở nên dễ tiếp cận hơn.

Giống như flight recorder, chúng tôi cũng vui mừng thông báo rằng chúng tôi có một API đọc trace
thử nghiệm muốn chia sẻ.
Nó có sẵn trong [chính gói chứa flight recorder,
golang.org/x/exp/trace](/pkg/golang.org/x/exp/trace#Reader).

Chúng tôi cho rằng nó đã đủ tốt để bắt đầu xây dựng những thứ trên đó, vì vậy hãy thử dùng!
Bên dưới là ví dụ đo tỷ lệ các sự kiện goroutine bị chặn do chờ
mạng.

{{code "execution-traces-2024/reader.go" `/START/` `/END/`}}

Và giống như flight recorder, [issue đề xuất](/issue/62627) là nơi rất phù hợp để để lại phản hồi.

Chúng tôi cũng muốn nhanh chóng nhắc đến Dominik Honnef, người đã thử nó từ sớm, đưa ra phản hồi tuyệt vời,
và đóng góp hỗ trợ cho các phiên bản trace cũ hơn vào API.

## Cảm ơn!

Công việc này được hoàn thành, không phần nhỏ nhờ sự giúp đỡ của [nhóm làm việc về chẩn đoán](/issue/57175),
được khởi xướng hơn một năm trước như một sự hợp tác giữa các bên liên quan từ khắp cộng đồng Go và mở cho công chúng.

Chúng tôi muốn dành một chút để cảm ơn các thành viên cộng đồng đã tham gia các buổi họp chẩn đoán
thường xuyên trong suốt năm qua: Felix Geisendörfer, Nick Ripley, Rhys Hiltner, Dominik
Honnef, Bryan Boreham, thepudds.

Các cuộc thảo luận, phản hồi và công sức của mọi người đã đóng vai trò quan trọng đưa chúng tôi tới
vị trí hôm nay.
Xin cảm ơn!
