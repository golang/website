---
title: "Các mẫu đồng thời trong Go: Hết thời gian, chuyển bước"
date: 2010-09-23
by:
- Andrew Gerrand
tags:
- concurrency
- technical
summary: Cách triển khai timeout bằng hỗ trợ đồng thời của Go.
template: true
---


Lập trình đồng thời có những thành ngữ riêng của nó.
Một ví dụ điển hình là timeout. Mặc dù channel của Go không hỗ trợ trực tiếp,
chúng rất dễ triển khai.
Giả sử ta muốn nhận từ channel `ch`,
nhưng muốn chờ nhiều nhất một giây để giá trị đến nơi.
Ta bắt đầu bằng cách tạo một channel báo hiệu và khởi động một goroutine
ngủ trước khi gửi vào channel đó:

{{raw `
	timeout := make(chan bool, 1)
	go func() {
	    time.Sleep(1 * time.Second)
	    timeout <- true
	}()
`}}

Sau đó ta có thể dùng câu lệnh `select` để nhận từ `ch` hoặc `timeout`.
Nếu không có gì đến trên `ch` sau một giây,
nhánh timeout sẽ được chọn và nỗ lực đọc từ `ch` bị hủy bỏ.

{{raw `
	select {
	case <-ch:
	    // a read from ch has occurred
	case <-timeout:
	    // the read from ch has timed out
	}
`}}

Channel `timeout` được đệm với chỗ cho 1 giá trị,
cho phép goroutine timeout gửi vào channel rồi thoát ra.
Goroutine không biết (và cũng không quan tâm) liệu giá trị đó có được nhận hay không.
Điều này có nghĩa là goroutine sẽ không bị treo mãi nếu việc nhận từ `ch` xảy ra
trước khi hết thời gian.
Channel `timeout` cuối cùng sẽ được bộ gom rác giải phóng.

(Trong ví dụ này, chúng tôi dùng `time.Sleep` để minh họa cơ chế của goroutine và channel.
Trong chương trình thực tế, bạn nên dùng [`time.After`](/pkg/time/#After),
một hàm trả về một channel và gửi vào channel đó sau khoảng thời gian được chỉ định.)

Hãy xem một biến thể khác của mẫu này.
Trong ví dụ này, ta có một chương trình đọc đồng thời từ nhiều cơ sở dữ liệu nhân bản.
Chương trình chỉ cần một câu trả lời,
và nó nên chấp nhận câu trả lời đến trước.

Hàm `Query` nhận một slice các kết nối cơ sở dữ liệu và một chuỗi `query`.
Nó truy vấn từng cơ sở dữ liệu song song và trả về phản hồi đầu tiên nó nhận được:

{{raw `
	func Query(conns []Conn, query string) Result {
	    ch := make(chan Result)
	    for _, conn := range conns {
	        go func(c Conn) {
	            select {
	            case ch <- c.DoQuery(query):
	            default:
	            }
	        }(conn)
	    }
	    return <-ch
	}
`}}

Trong ví dụ này, closure thực hiện một thao tác gửi không chặn,
điều mà nó đạt được bằng cách dùng thao tác gửi trong câu lệnh `select` với một nhánh `default`.
Nếu việc gửi không thể hoàn thành ngay thì nhánh mặc định sẽ được chọn.
Việc làm cho thao tác gửi không chặn bảo đảm rằng không goroutine nào được khởi động
trong vòng lặp sẽ bị treo lại.
Tuy nhiên, nếu kết quả đến trước khi hàm chính kịp thực hiện lệnh nhận,
thì việc gửi có thể thất bại vì chưa ai sẵn sàng nhận.

Vấn đề này là ví dụ kinh điển về thứ được gọi là [race condition](https://en.wikipedia.org/wiki/Race_condition),
nhưng cách sửa thì rất đơn giản.
Ta chỉ cần bảo đảm channel `ch` có bộ đệm (bằng cách thêm độ dài bộ đệm
làm đối số thứ hai cho [make](/pkg/builtin/#make)),
để chắc chắn rằng lần gửi đầu tiên có chỗ để đặt giá trị.
Điều này bảo đảm việc gửi sẽ luôn thành công,
và giá trị đầu tiên tới nơi sẽ được lấy ra bất kể thứ tự thực thi.

Hai ví dụ này cho thấy sự đơn giản mà Go dùng để biểu đạt
những tương tác phức tạp giữa các goroutine.
