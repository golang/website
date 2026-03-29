---
title: Fuzzing đã sẵn sàng ở mức beta
date: 2021-06-03
by:
- Katie Hockman
- Jay Conrod
tags:
- fuzz
- testing
summary: Fuzzing gốc của Go hiện đã sẵn sàng cho thử nghiệm beta trên tip.
---


Chúng tôi rất vui mừng thông báo rằng fuzzing gốc đã sẵn sàng cho thử nghiệm beta trên tip!

Fuzzing là một dạng kiểm thử tự động liên tục biến đổi đầu vào của
một chương trình để tìm ra các vấn đề như panic hoặc bug. Những phép biến đổi dữ liệu
bán ngẫu nhiên này có thể phát hiện vùng bao phủ mã mới mà các unit test hiện có có thể bỏ lỡ, và
làm lộ ra những lỗi ở các trường hợp biên vốn có thể không bao giờ bị phát hiện.
Vì fuzzing có thể chạm tới những trường hợp biên này, fuzz testing đặc biệt có giá trị
trong việc tìm ra các khai thác và lỗ hổng bảo mật.

Xem
[golang.org/s/draft-fuzzing-design](/s/draft-fuzzing-design)
để biết thêm chi tiết về tính năng này.


## Bắt đầu

Để bắt đầu, bạn có thể chạy:

	$ go install golang.org/dl/gotip@latest
	$ gotip download

Lệnh này sẽ xây dựng bộ công cụ Go từ nhánh master. Sau khi chạy xong, `gotip`
có thể đóng vai trò thay thế trực tiếp cho lệnh `go`. Giờ bạn có thể chạy các lệnh
như

	$ gotip test -fuzz=Fuzz

## Viết một fuzz test

Một fuzz test phải nằm trong một tệp `*_test.go` dưới dạng một hàm có mẫu `FuzzXxx`.
Hàm này phải nhận một đối số `*testing.F`, tương tự như cách `*testing.T`
được truyền vào một hàm `TestXxx`.

Dưới đây là ví dụ về một fuzz test đang kiểm tra hành vi của [gói
net/url](https://pkg.go.dev/net/url#ParseQuery).

	//go:build go1.18
	// +build go1.18

	package fuzz

	import (
		"net/url"
		"reflect"
		"testing"
	)

	func FuzzParseQuery(f *testing.F) {
		f.Add("x=1&y=2")
		f.Fuzz(func(t *testing.T, queryStr string) {
			query, err := url.ParseQuery(queryStr)
			if err != nil {
				t.Skip()
			}
			queryStr2 := query.Encode()
			query2, err := url.ParseQuery(queryStr2)
			if err != nil {
				t.Fatalf("ParseQuery failed to decode a valid encoded query %s: %v", queryStr2, err)
			}
			if !reflect.DeepEqual(query, query2) {
				t.Errorf("ParseQuery gave different query after being encoded\nbefore: %v\nafter: %v", query, query2)
			}
		})
	}

Bạn có thể đọc thêm về fuzzing trên pkg.go.dev, bao gồm [tổng quan
về fuzzing với Go](https://pkg.go.dev/testing@master#hdr-Fuzzing) và
[godoc cho kiểu `testing.F` mới](https://pkg.go.dev/testing@master#F).

## Kỳ vọng

Đây là một tính năng mới vẫn còn ở giai đoạn beta, vì vậy bạn nên kỳ vọng sẽ còn có lỗi
và bộ tính năng chưa hoàn chỉnh. Hãy kiểm tra [issue tracker với các issue gắn nhãn
“fuzz”](https://github.com/golang/go/issues?q=is%3Aopen+is%3Aissue+label%3Afuzz)
để luôn cập nhật về các lỗi hiện có và các tính năng còn thiếu.

Xin lưu ý rằng fuzzing có thể tiêu tốn nhiều bộ nhớ và có thể ảnh hưởng đến
hiệu năng của máy trong lúc chạy. `go test -fuzz` mặc định chạy fuzzing
song song trong `$GOMAXPROCS` tiến trình. Bạn có thể giảm số tiến trình được
dùng trong khi fuzz bằng cách đặt rõ cờ `-parallel` với `go test`.
Hãy đọc tài liệu cho lệnh `go test` bằng cách chạy `gotip help
testflag` nếu bạn muốn biết thêm.

Cũng xin lưu ý rằng bộ máy fuzzing ghi các giá trị giúp mở rộng độ bao phủ kiểm thử vào
một thư mục fuzz cache trong `$GOCACHE/fuzz` trong lúc chạy. Hiện chưa có
giới hạn cho số lượng tệp hoặc tổng số byte có thể được ghi vào fuzz
cache, vì vậy nó có thể chiếm khá nhiều dung lượng lưu trữ (ví dụ: vài GB). Bạn có thể
xóa fuzz cache bằng cách chạy `gotip clean -fuzzcache`.

## Tiếp theo là gì?

Tính năng này sẽ có mặt bắt đầu từ Go 1.18.

Nếu bạn gặp bất kỳ vấn đề nào hoặc có ý tưởng cho một tính năng, vui lòng [tạo issue](/issue/new/?&labels=fuzz).

Để thảo luận và đưa phản hồi chung về tính năng này, bạn cũng có thể tham gia
[kênh #fuzzing](https://gophers.slack.com/archives/CH5KV1AKE) trên
Gophers Slack.

Chúc fuzzing vui vẻ!
