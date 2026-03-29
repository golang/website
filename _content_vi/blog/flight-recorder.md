---
title: Flight Recorder trong Go 1.25
date: 2025-09-26
by:
- Carlos Amedee
- Michael Knyszek
tags:
- debug
- technical
- tracing
- flight recorder
summary: Go 1.25 giới thiệu một công cụ mới trong bộ công cụ chẩn đoán: flight recording.
template: true
---

Trong năm 2024, chúng tôi đã giới thiệu
[execution trace mạnh hơn của Go](/blog/execution-traces-2024). Trong bài viết đó,
chúng tôi đã hé lộ một số chức năng mới có thể mở khóa với execution tracer mới,
bao gồm cả *flight recording*. Chúng tôi vui mừng thông báo rằng flight recording
giờ đã có trong Go 1.25, và đó là một công cụ mạnh mẽ mới trong hộp đồ nghề chẩn đoán của Go.

## Execution trace

Trước tiên, hãy điểm lại nhanh execution trace của Go.

Runtime Go có thể được yêu cầu ghi ra một nhật ký lưu lại nhiều sự kiện xảy ra trong
quá trình thực thi của một ứng dụng Go. Nhật ký đó được gọi là runtime execution trace.
Execution trace của Go chứa rất nhiều thông tin về cách các goroutine tương tác với nhau
và với hệ thống bên dưới. Điều này khiến chúng đặc biệt hữu ích khi gỡ lỗi các vấn đề về độ trễ, vì
chúng cho bạn biết cả lúc goroutine đang chạy, và quan trọng không kém, lúc chúng không chạy.

Gói [runtime/trace](/pkg/runtime/trace) cung cấp API để thu thập
một execution trace trong một cửa sổ thời gian xác định bằng cách gọi `runtime/trace.Start` và `runtime/trace.Stop`.
Cách này hoạt động tốt nếu đoạn mã bạn muốn trace chỉ là test, microbenchmark hoặc công cụ
dòng lệnh. Bạn có thể thu một trace cho toàn bộ quá trình thực thi đầu-cuối, hoặc chỉ những phần bạn quan tâm.

Tuy nhiên, với những dịch vụ web chạy lâu dài, tức loại ứng dụng mà Go nổi tiếng,
thì như vậy là chưa đủ. Máy chủ web có thể chạy hàng ngày hoặc thậm chí hàng tuần, và việc thu trace cho
toàn bộ quá trình thực thi sẽ tạo ra quá nhiều dữ liệu để sàng lọc. Thường chỉ một phần
trong quá trình chạy của chương trình bị lỗi, chẳng hạn một request timeout hoặc
một lần health check thất bại. Đến lúc nó xảy ra thì đã quá muộn để gọi `Start`!

Một cách tiếp cận cho bài toán này là lấy mẫu ngẫu nhiên execution trace trên toàn bộ fleet.
Dù cách này mạnh và có thể giúp phát hiện vấn đề trước khi chúng trở thành outage, nó
đòi hỏi rất nhiều hạ tầng để bắt đầu. Lượng lớn execution trace data
sẽ cần được lưu trữ, phân loại và xử lý, trong đó phần lớn thậm chí chẳng chứa điều gì
thú vị. Và khi bạn đang cố đào đến cùng một sự cố cụ thể,
nó gần như không dùng được.

## Flight recording

Và đó là lúc flight recorder xuất hiện.

Một chương trình thường biết khi nào có gì đó trục trặc, nhưng nguyên nhân gốc có thể đã xảy ra
từ lâu. Flight recorder cho phép bạn thu một trace của vài giây cuối cùng
dẫn tới thời điểm chương trình phát hiện có vấn đề.

Flight recorder thu execution trace như bình thường, nhưng thay vì ghi ra
socket hoặc tệp, nó đệm vài giây trace gần nhất trong bộ nhớ. Bất kỳ lúc nào,
chương trình có thể yêu cầu nội dung của bộ đệm và chụp lại chính xác
cửa sổ thời gian có vấn đề. Flight recorder giống như một lưỡi dao mổ cắt thẳng vào vùng bệnh.

## Ví dụ

Hãy học cách dùng flight recorder qua một ví dụ. Cụ thể, ta sẽ dùng nó để
chẩn đoán một vấn đề hiệu năng với một máy chủ HTTP hiện thực trò chơi "đoán số".
Nó phơi ra endpoint `/guess-number` nhận vào một số nguyên và phản hồi cho phía gọi
biết họ có đoán đúng số hay không. Ngoài ra còn có một goroutine, cứ mỗi phút một lần,
gửi một báo cáo về tất cả các số đã được đoán tới một dịch vụ khác thông qua một HTTP request.

{{raw `
	// bucket is a simple mutex-protected counter.
	type bucket struct {
		mu      sync.Mutex
		guesses int
	}

	func main() {
		// Make one bucket for each valid number a client could guess.
		// The HTTP handler will look up the guessed number in buckets by
		// using the number as an index into the slice.
		buckets := make([]bucket, 100)

		// Every minute, we send a report of how many times each number was guessed.
		go func() {
			for range time.Tick(1 * time.Minute) {
				sendReport(buckets)
			}
		}()

		// Choose the number to be guessed.
		answer := rand.Intn(len(buckets))

		http.HandleFunc("/guess-number", func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// Fetch the number from the URL query variable "guess" and convert it
			// to an integer. Then, validate it.
			guess, err := strconv.Atoi(r.URL.Query().Get("guess"))
			if err != nil || !(0 <= guess && guess < len(buckets)) {
				http.Error(w, "invalid 'guess' value", http.StatusBadRequest)
				return
			}

			// Select the appropriate bucket and safely increment its value.
			b := &buckets[guess]
			b.mu.Lock()
			b.guesses++
			b.mu.Unlock()

			// Respond to the client with the guess and whether it was correct.
			fmt.Fprintf(w, "guess: %d, correct: %t", guess, guess == answer)

			log.Printf("HTTP request: endpoint=/guess-number guess=%d duration=%s", guess, time.Since(start))
		})
		log.Fatal(http.ListenAndServe(":8090", nil))
	}

	// sendReport posts the current state of buckets to a remote service.
	func sendReport(buckets []bucket) {
		counts := make([]int, len(buckets))

		for index := range buckets {
			b := &buckets[index]
			b.mu.Lock()
			defer b.mu.Unlock()

			counts[index] = b.guesses
		}

		// Marshal the report data into a JSON payload.
		b, err := json.Marshal(counts)
		if err != nil {
			log.Printf("failed to marshal report data: error=%s", err)
			return
		}
		url := "http://localhost:8091/guess-number-report"
		if _, err := http.Post(url, "application/json", bytes.NewReader(b)); err != nil {
			log.Printf("failed to send report: %s", err)
		}
	}
`}}

Đây là toàn bộ mã cho máy chủ:
[https://go.dev/play/p/rX1eyKtVglF](/play/p/rX1eyKtVglF), và cho một client đơn giản:
[https://go.dev/play/p/2PjQ-1ORPiw](/play/p/2PjQ-1ORPiw). Để tránh phải chạy thêm tiến trình thứ ba,
"client" cũng hiện thực máy chủ nhận report, dù trong hệ thống thật phần này sẽ
tách biệt.

Giả sử sau khi triển khai ứng dụng lên production, chúng tôi nhận được phàn nàn từ
người dùng rằng một số lời gọi `/guess-number` mất nhiều thời gian hơn mong đợi. Khi xem
log, chúng tôi thấy đôi khi thời gian phản hồi vượt quá 100 mili giây, trong khi đa số
lời gọi chỉ mất cỡ micro giây.

```
2025/09/19 16:52:02 HTTP request: endpoint=/guess-number guess=69 duration=625ns
2025/09/19 16:52:02 HTTP request: endpoint=/guess-number guess=62 duration=458ns
2025/09/19 16:52:02 HTTP request: endpoint=/guess-number guess=42 duration=1.417µs
2025/09/19 16:52:02 HTTP request: endpoint=/guess-number guess=86 duration=115.186167ms
2025/09/19 16:52:02 HTTP request: endpoint=/guess-number guess=0 duration=127.993375ms
```

Trước khi tiếp tục, hãy dành một phút xem bạn có nhận ra vấn đề gì không!

Bất kể bạn có tìm ra hay không, giờ ta hãy đào sâu hơn và xem cách tìm vấn đề
từ những nguyên lý đầu tiên. Đặc biệt, sẽ rất tuyệt nếu ta có thể thấy
ứng dụng đã làm gì trong khoảng thời gian trước phản hồi chậm. Đây chính là
điều mà flight recorder được xây ra để làm! Ta sẽ dùng nó để chụp execution trace khi
ta thấy phản hồi đầu tiên vượt quá 100 mili giây.

Đầu tiên, trong `main`, ta sẽ cấu hình và khởi động flight recorder:

{{raw `
	// Set up the flight recorder
	fr := trace.NewFlightRecorder(trace.FlightRecorderConfig{
		MinAge:   200 * time.Millisecond,
		MaxBytes: 1 << 20, // 1 MiB
	})
	fr.Start()
`}}

`MinAge` cấu hình khoảng thời gian dữ liệu trace được giữ lại một cách đáng tin cậy, và chúng tôi
khuyên nên đặt nó khoảng 2x cửa sổ thời gian của sự kiện. Ví dụ, nếu bạn
đang gỡ lỗi một timeout 5 giây, hãy đặt nó là 10 giây. `MaxBytes` cấu hình
kích thước bộ đệm trace để bạn không làm bùng nổ mức dùng bộ nhớ. Trung bình,
bạn có thể kỳ vọng trace tạo ra vài MB dữ liệu mỗi giây thực thi,
hoặc 10 MB/s với một dịch vụ bận rộn.

Tiếp theo, ta thêm một hàm trợ giúp để chụp snapshot và ghi nó ra tệp:

{{raw `
	var once sync.Once

	// captureSnapshot captures a flight recorder snapshot.
	func captureSnapshot(fr *trace.FlightRecorder) {
		// once.Do ensures that the provided function is executed only once.
		once.Do(func() {
			f, err := os.Create("snapshot.trace")
			if err != nil {
				log.Printf("opening snapshot file %s failed: %s", f.Name(), err)
				return
			}
			defer f.Close() // ignore error

			// WriteTo writes the flight recorder data to the provided io.Writer.
			_, err = fr.WriteTo(f)
			if err != nil {
				log.Printf("writing snapshot to file %s failed: %s", f.Name(), err)
				return
			}

			// Stop the flight recorder after the snapshot has been taken.
			fr.Stop()
			log.Printf("captured a flight recorder snapshot to %s", f.Name())
		})
	}
`}}

Và cuối cùng, ngay trước khi ghi log một request đã hoàn tất, ta sẽ kích hoạt snapshot nếu request
mất hơn 100 mili giây:

```go
// Capture a snapshot if the response takes more than 100ms.
// Only the first call has any effect.
if fr.Enabled() && time.Since(start) > 100*time.Millisecond {
	go captureSnapshot(fr)
}
```

Đây là toàn bộ mã của máy chủ, giờ đã được instrument bằng flight recorder:
[https://go.dev/play/p/3V33gfIpmjG](/play/p/3V33gfIpmjG)

Giờ ta chạy máy chủ lại và gửi request cho đến khi gặp một request chậm đủ để kích hoạt
snapshot.

Khi đã có trace, ta cần một công cụ giúp kiểm tra nó. Toolchain Go
cung cấp sẵn một công cụ phân tích execution trace thông qua
[lệnh `go tool trace`](https://pkg.go.dev/cmd/trace). Hãy chạy `go tool trace snapshot.trace`
để khởi chạy công cụ, nó sẽ tạo một web server cục bộ, rồi mở URL được hiển thị trong trình duyệt
(nếu công cụ không tự mở trình duyệt cho bạn).

Công cụ này cho ta một vài cách để nhìn vào trace, nhưng hãy tập trung vào việc trực quan hóa trace
để có cảm giác chuyện gì đang diễn ra. Hãy bấm “View trace by proc”.

Trong góc nhìn này, trace được trình bày dưới dạng một dòng thời gian của các sự kiện. Ở đầu trang, trong
phần “STATS”, ta có thể thấy bản tóm tắt trạng thái của ứng dụng, bao gồm
số luồng, kích thước heap và số goroutine.

Bên dưới, trong phần “PROCS”, ta có thể thấy việc thực thi goroutine được ánh xạ
lên `GOMAXPROCS` (số luồng hệ điều hành mà ứng dụng Go tạo ra) như thế nào. Ta
có thể thấy khi nào và bằng cách nào mỗi goroutine bắt đầu, chạy và cuối cùng dừng thực thi.

Tạm thời, hãy tập trung vào khoảng trống khổng lồ về hoạt động ở phía bên phải của
trình xem. Trong một khoảng thời gian khoảng 100ms, chẳng có gì xảy ra cả!

<a href="flight-recorder/flight_recorder_1.png"><img src="flight-recorder/flight_recorder_1.png" width=100%></a>

Bằng cách chọn công cụ `zoom` (hoặc nhấn `3`), ta có thể quan sát phần trace ngay
sau khoảng trống đó với nhiều chi tiết hơn.

<a href="flight-recorder/flight_recorder_2.png"><img src="flight-recorder/flight_recorder_2.png" width=100%></a>

Ngoài hoạt động của từng goroutine riêng lẻ, ta còn có thể thấy cách các goroutine tương tác
thông qua “flow event”. Một incoming flow event chỉ ra điều gì đã khiến một goroutine
bắt đầu chạy. Một outgoing flow edge chỉ ra một goroutine đã tác động lên goroutine khác
ra sao. Việc bật hiển thị tất cả flow event thường cung cấp những manh mối gợi ý về nguồn gốc vấn đề.

<a href="flight-recorder/flight_recorder_3.png"><img src="flight-recorder/flight_recorder_3.png" width=100%></a>

Trong trường hợp này, ta có thể thấy nhiều goroutine nối trực tiếp tới một
goroutine duy nhất ngay sau quãng dừng hoạt động.

Bấm vào goroutine duy nhất đó sẽ hiện ra một bảng sự kiện đầy các outgoing flow event,
khớp với điều ta thấy khi bật chế độ xem flow.

Goroutine đó đã làm gì khi nó chạy? Một phần thông tin được lưu trong trace là ảnh chụp
stack trace tại các thời điểm khác nhau. Khi nhìn vào goroutine này, ta thấy stack trace lúc bắt đầu cho biết
nó đang chờ HTTP request hoàn tất khi goroutine được lên lịch chạy.
Và stack trace lúc kết thúc cho thấy hàm `sendReport`
đã trả về và goroutine đang chờ ticker cho thời điểm gửi report tiếp theo.

<a href="flight-recorder/flight_recorder_4.png"><img src="flight-recorder/flight_recorder_4.png" style="padding: inherit;margin:auto;display: block;"></a>

Giữa lúc bắt đầu và lúc kết thúc goroutine chạy, ta thấy có một số lượng rất lớn
“outgoing flows”, nơi nó tương tác với các goroutine khác. Bấm vào một trong các
mục `Outgoing flow` sẽ đưa ta tới phần xem tương tác đó.

<a href="flight-recorder/flight_recorder_5.png"><img src="flight-recorder/flight_recorder_5.png" width=100%></a>

Luồng tương tác này chỉ ra `Unlock` trong `sendReport`:

```go
for index := range buckets {
	b := &buckets[index]
	b.mu.Lock()
	defer b.mu.Unlock()

	counts[index] = b.guesses
}
```

Trong `sendReport`, ý định của chúng ta là khóa từng bucket rồi mở khóa ngay sau
khi sao chép giá trị.

Nhưng đây là vấn đề: thực ra ta không hề mở khóa ngay sau khi sao chép
giá trị chứa trong `bucket.guesses`. Bởi vì ta dùng một câu lệnh `defer` để mở khóa,
việc mở khóa đó không xảy ra cho đến khi hàm trả về. Ta giữ khóa không chỉ
qua hết vòng lặp, mà còn cho tới tận sau khi HTTP request hoàn tất. Đây là một lỗi tinh vi
có thể rất khó lần ra trong một hệ thống production lớn.

May mắn là execution tracing đã giúp chúng ta chỉ ra chính xác vấn đề. Tuy nhiên,
nếu cố dùng execution tracer trên một máy chủ chạy lâu dài mà không có chế độ
flight-recording mới, nó có thể tích lũy một lượng cực lớn execution trace data,
và người vận hành sẽ phải lưu, truyền và sàng lọc nó. Flight recorder trao cho ta
sức mạnh của việc nhìn lại quá khứ. Nó cho phép chụp đúng thứ đã xảy ra,
sau khi sự cố đã diễn ra, rồi nhanh chóng khoanh vùng nguyên nhân.

Flight recorder chỉ là bổ sung mới nhất vào hộp công cụ của lập trình viên Go để
chẩn đoán hoạt động bên trong của các ứng dụng đang chạy. Chúng tôi đã liên tục cải thiện
tracing trong vài bản phát hành gần đây. Go 1.21 đã giảm rất mạnh overhead lúc chạy
của tracing. Định dạng trace trở nên mạnh mẽ hơn và cũng có thể chia tách được trong Go 1.22,
từ đó dẫn đến những tính năng như flight recorder. Các công cụ mã nguồn mở như
[gotraceui](https://gotraceui.dev/), và [khả năng sắp có để parse execution trace
bằng mã](/issue/62627) là những cách khác để tận dụng sức mạnh của
execution trace. [Trang Diagnostics](/doc/diagnostics) liệt kê thêm nhiều
công cụ khác theo ý bạn. Chúng tôi hy vọng bạn sẽ dùng chúng trong lúc viết và tinh chỉnh
các ứng dụng Go của mình.

## Lời cảm ơn

Chúng tôi muốn dành một chút để cảm ơn các thành viên cộng đồng đã tích cực tham gia các cuộc họp
diagnostics, đóng góp cho các thiết kế và cung cấp phản hồi qua nhiều năm:
Felix Geisendörfer ([@felixge.de](https://bsky.app/profile/felixge.de)),
Nick Ripley ([@nsrip-dd](https://github.com/nsrip-dd)),
Rhys Hiltner ([@rhysh](https://github.com/rhysh)),
Dominik Honnef ([@dominikh](https://github.com/dominikh)),
Bryan Boreham ([@bboreham](https://github.com/bboreham)),
và PJ Malloy ([@thepudds](https://github.com/thepudds)).

Những cuộc thảo luận, phản hồi và công sức của mọi người đã đóng vai trò quan trọng trong việc đưa
chúng tôi tới một tương lai diagnostics tốt hơn. Xin cảm ơn!
