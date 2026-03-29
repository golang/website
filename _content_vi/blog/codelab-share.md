---
title: Chia sẻ bộ nhớ bằng giao tiếp
date: 2010-07-13
by:
- Andrew Gerrand
tags:
- concurrency
- technical
summary: Bản xem trước của codelab mới cho Go, Share Memory by Communicating.
template: true
---


Các mô hình luồng truyền thống (thường dùng khi viết chương trình Java,
C++ và Python, chẳng hạn) yêu cầu lập trình viên giao tiếp
giữa các luồng bằng cách dùng bộ nhớ dùng chung.
Thông thường, các cấu trúc dữ liệu dùng chung được bảo vệ bằng khóa,
và các luồng sẽ tranh chấp những khóa đó để truy cập dữ liệu.
Trong một số trường hợp, điều này trở nên dễ dàng hơn nhờ sử dụng các cấu trúc dữ liệu an toàn cho luồng
như Queue của Python.

Các nguyên thủy đồng thời của Go, goroutine và channel, cung cấp một
cách thanh nhã và khác biệt để cấu trúc phần mềm đồng thời.
(Các khái niệm này có một [lịch sử thú vị](https://swtch.com/~rsc/thread/) bắt đầu từ C.
[Communicating Sequential Processes](http://www.usingcsp.com/) của A. R. Hoare.)
Thay vì dùng khóa một cách tường minh để điều phối quyền truy cập vào dữ liệu dùng chung,
Go khuyến khích dùng channel để truyền tham chiếu tới dữ liệu giữa các goroutine.
Cách tiếp cận này bảo đảm rằng chỉ có một goroutine có quyền truy cập dữ liệu tại một thời điểm nhất định.
Khái niệm này được tóm tắt trong tài liệu [Effective Go](/doc/effective_go.html)
(bắt buộc phải đọc với mọi lập trình viên Go):

_Đừng giao tiếp bằng cách chia sẻ bộ nhớ; thay vào đó, hãy chia sẻ bộ nhớ bằng cách giao tiếp._

Hãy xét một chương trình thăm dò một danh sách URL.
Trong môi trường luồng truyền thống, người ta có thể cấu trúc dữ liệu như sau:

	type Resource struct {
	    url        string
	    polling    bool
	    lastPolled int64
	}

	type Resources struct {
	    data []*Resource
	    lock *sync.Mutex
	}

Và rồi một hàm Poller (nhiều bản sao của nó sẽ chạy trong các luồng riêng) có thể trông như sau:

{{raw `
	func Poller(res *Resources) {
	    for {
	        // get the least recently-polled Resource
	        // and mark it as being polled
	        res.lock.Lock()
	        var r *Resource
	        for _, v := range res.data {
	            if v.polling {
	                continue
	            }
	            if r == nil || v.lastPolled < r.lastPolled {
	                r = v
	            }
	        }
	        if r != nil {
	            r.polling = true
	        }
	        res.lock.Unlock()
	        if r == nil {
	            continue
	        }

	        // poll the URL

	        // update the Resource's polling and lastPolled
	        res.lock.Lock()
	        r.polling = false
	        r.lastPolled = time.Nanoseconds()
	        res.lock.Unlock()
	    }
	}
`}}

Hàm này dài khoảng một trang và cần thêm chi tiết để hoàn chỉnh.
Nó thậm chí còn chưa bao gồm logic thăm dò URL (thứ tự thân nó
chỉ là vài dòng), và cũng không xử lý một cách nhã nhặn trường hợp
cạn kiệt pool Resource.

Hãy xem cùng chức năng đó được triển khai theo thành ngữ Go.
Trong ví dụ này, Poller là một hàm nhận Resource cần được thăm dò
từ một channel đầu vào,
và gửi chúng sang một channel đầu ra khi hoàn tất.

{{raw `
	type Resource string

	func Poller(in, out chan *Resource) {
	    for r := range in {
	        // poll the URL

	        // send the processed Resource to out
	        out <- r
	    }
	}
`}}

Logic tinh vi từ ví dụ trước đã biến mất một cách dễ thấy,
và cấu trúc dữ liệu Resource của chúng ta cũng không còn chứa dữ liệu sổ sách nữa.
Thực tế, tất cả những gì còn lại chính là các phần quan trọng.
Điều này sẽ cho bạn cảm nhận ban đầu về sức mạnh của những tính năng ngôn ngữ đơn giản ấy.

Có nhiều phần bị lược bỏ khỏi các đoạn mã trên.
Để xem phần diễn giải của một chương trình Go hoàn chỉnh, đúng thành ngữ và dùng các ý tưởng này,
hãy xem Codewalk [_Share Memory By Communicating_](/doc/codewalk/sharemem/).
