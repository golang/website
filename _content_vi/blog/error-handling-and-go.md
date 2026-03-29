---
title: Xử lý lỗi và Go
date: 2011-07-12
by:
- Andrew Gerrand
tags:
- error
- interface
- type
- technical
summary: Giới thiệu về lỗi trong Go.
template: true
---

## Giới thiệu

Nếu bạn đã từng viết bất kỳ đoạn mã Go nào thì hẳn bạn đã gặp kiểu tích hợp sẵn `error`.
Mã Go dùng các giá trị `error` để biểu thị một trạng thái bất thường.
Ví dụ, hàm `os.Open` trả về một giá trị `error` khác nil khi
nó không mở được tệp.

	func Open(name string) (file *File, err error)

Đoạn mã sau dùng `os.Open` để mở một tệp.
Nếu có lỗi xảy ra, nó gọi `log.Fatal` để in thông báo lỗi rồi dừng.

	f, err := os.Open("filename.ext")
	if err != nil {
	    log.Fatal(err)
	}
	// do something with the open *File f

Bạn có thể làm được rất nhiều việc trong Go chỉ với chừng đó hiểu biết về kiểu `error`,
nhưng trong bài viết này chúng ta sẽ nhìn kỹ hơn vào `error` và thảo luận một số
thực hành tốt khi xử lý lỗi trong Go.

## Kiểu error

Kiểu `error` là một kiểu interface. Một biến `error` đại diện cho bất kỳ
giá trị nào có thể tự mô tả chính nó dưới dạng chuỗi.
Đây là khai báo của interface đó:

	type error interface {
	    Error() string
	}

Kiểu `error`, giống như mọi kiểu tích hợp sẵn khác,
được [khai báo trước](/doc/go_spec.html#Predeclared_identifiers)
trong [universe block](/doc/go_spec.html#Blocks).

Hiện thực `error` được dùng phổ biến nhất là kiểu `errorString` không export
của gói [errors](/pkg/errors/).

	// errorString is a trivial implementation of error.
	type errorString struct {
	    s string
	}

	func (e *errorString) Error() string {
	    return e.s
	}

Bạn có thể tạo một giá trị như vậy bằng hàm `errors.New`.
Nó nhận một chuỗi, chuyển nó thành `errors.errorString`, rồi trả về
dưới dạng một giá trị `error`.

	// New returns an error that formats as the given text.
	func New(text string) error {
	    return &errorString{text}
	}

Đây là cách bạn có thể dùng `errors.New`:

{{raw `
	func Sqrt(f float64) (float64, error) {
	    if f < 0 {
	        return 0, errors.New("math: square root of negative number")
	    }
	    // implementation
	}
`}}

Người gọi truyền một đối số âm vào `Sqrt` sẽ nhận được một giá trị `error`
khác nil (mà biểu diễn cụ thể là một giá trị `errors.errorString`).
Người gọi có thể lấy chuỗi lỗi (“math:
square root of...”) bằng cách gọi phương thức `Error` của `error`,
hoặc đơn giản là in nó:

	f, err := Sqrt(-1)
	if err != nil {
	    fmt.Println(err)
	}

Gói [fmt](/pkg/fmt/) định dạng một giá trị `error` bằng cách gọi phương thức `Error() string` của nó.

Việc tóm tắt ngữ cảnh là trách nhiệm của hiện thực lỗi.
Lỗi do `os.Open` trả về sẽ có dạng “open /etc/passwd:
permission denied”, chứ không chỉ là “permission denied”. Lỗi do
`Sqrt` của ta trả về lại thiếu thông tin về đối số không hợp lệ.

Để thêm thông tin đó, một hàm hữu ích là `Errorf` của gói `fmt`.
Nó định dạng một chuỗi theo các quy tắc của `Printf` và trả nó về dưới dạng `error`
được tạo bởi `errors.New`.

{{raw `
	if f < 0 {
	    return 0, fmt.Errorf("math: square root of negative number %g", f)
	}
`}}

Trong nhiều trường hợp `fmt.Errorf` là đủ dùng,
nhưng vì `error` là một interface, bạn có thể dùng các cấu trúc dữ liệu tùy ý làm giá trị lỗi,
để cho phép người gọi kiểm tra chi tiết của lỗi.

Ví dụ, giả sử người gọi của chúng ta muốn khôi phục đối số không hợp lệ
đã truyền vào `Sqrt`.
Ta có thể hỗ trợ điều đó bằng cách định nghĩa một hiện thực lỗi mới thay vì dùng
`errors.errorString`:

	type NegativeSqrtError float64

	func (f NegativeSqrtError) Error() string {
	    return fmt.Sprintf("math: square root of negative number %g", float64(f))
	}

Một người gọi tinh vi hơn khi đó có thể dùng [type assertion](/doc/go_spec.html#Type_assertions)
để kiểm tra `NegativeSqrtError` và xử lý riêng,
trong khi người gọi chỉ đơn giản truyền lỗi vào `fmt.Println` hay `log.Fatal` sẽ
không thấy thay đổi gì về hành vi.

Làm thêm một ví dụ nữa, gói [json](/pkg/encoding/json/)
định nghĩa kiểu `SyntaxError` mà hàm `json.Decode` trả về
khi nó gặp lỗi cú pháp lúc phân tích một blob JSON.

	type SyntaxError struct {
	    msg    string // description of error
	    Offset int64  // error occurred after reading Offset bytes
	}

	func (e *SyntaxError) Error() string { return e.msg }

Trường `Offset` thậm chí không được hiển thị trong cách định dạng mặc định của lỗi,
nhưng người gọi có thể dùng nó để thêm thông tin tệp và dòng vào thông báo lỗi:

	if err := dec.Decode(&val); err != nil {
	    if serr, ok := err.(*json.SyntaxError); ok {
	        line, col := findLine(f, serr.Offset)
	        return fmt.Errorf("%s:%d:%d: %v", f.Name(), line, col, err)
	    }
	    return err
	}

(Đây là phiên bản đã được đơn giản hóa đôi chút của [một đoạn mã thực](https://github.com/camlistore/go4/blob/03efcb870d84809319ea509714dd6d19a1498483/jsonconfig/eval.go#L123-L135)
trong dự án [Camlistore](http://camlistore.org).)

Interface `error` chỉ yêu cầu một phương thức `Error`;
các hiện thực lỗi cụ thể có thể có thêm phương thức khác.
Ví dụ, gói [net](/pkg/net/) trả về lỗi kiểu `error`,
theo đúng thông lệ thường thấy, nhưng một số hiện thực lỗi của nó có
thêm các phương thức được định nghĩa bởi interface `net.Error`:

	package net

	type Error interface {
	    error
	    Timeout() bool   // Is the error a timeout?
	    Temporary() bool // Is the error temporary?
	}

Mã phía client có thể kiểm tra `net.Error` bằng type assertion rồi phân biệt
lỗi mạng tạm thời với lỗi vĩnh viễn.
Ví dụ, một web crawler có thể ngủ rồi thử lại khi gặp lỗi tạm thời
và bỏ cuộc trong các trường hợp khác.

	if nerr, ok := err.(net.Error); ok && nerr.Temporary() {
	    time.Sleep(1e9)
	    continue
	}
	if err != nil {
	    log.Fatal(err)
	}

## Đơn giản hóa việc xử lý lỗi lặp đi lặp lại

Trong Go, xử lý lỗi là quan trọng. Thiết kế và thông lệ của ngôn ngữ
khuyến khích bạn kiểm tra lỗi một cách tường minh ở đúng nơi chúng xảy ra (khác với
thông lệ ở các ngôn ngữ khác là ném exception rồi đôi khi bắt chúng).
Trong một số trường hợp điều này làm mã Go dài dòng,
nhưng may mắn là có vài kỹ thuật bạn có thể dùng để giảm thiểu việc xử lý lỗi lặp lại.

Hãy xét một ứng dụng [App Engine](https://cloud.google.com/appengine/docs/go/)
có một HTTP handler lấy một bản ghi từ datastore
rồi định dạng nó bằng template.

	func init() {
	    http.HandleFunc("/view", viewRecord)
	}

	func viewRecord(w http.ResponseWriter, r *http.Request) {
	    c := appengine.NewContext(r)
	    key := datastore.NewKey(c, "Record", r.FormValue("id"), 0, nil)
	    record := new(Record)
	    if err := datastore.Get(c, key, record); err != nil {
	        http.Error(w, err.Error(), 500)
	        return
	    }
	    if err := viewTemplate.Execute(w, record); err != nil {
	        http.Error(w, err.Error(), 500)
	    }
	}

Hàm này xử lý các lỗi trả về từ `datastore.Get` và
phương thức `Execute` của `viewTemplate`.
Trong cả hai trường hợp, nó hiển thị cho người dùng một thông báo lỗi đơn giản với mã HTTP
500 (“Internal Server Error”).
Lượng mã này trông còn có thể chấp nhận được,
nhưng chỉ cần thêm vài HTTP handler nữa là bạn sẽ nhanh chóng có rất nhiều
bản sao của cùng một đoạn xử lý lỗi.

Để giảm sự lặp lại này, ta có thể định nghĩa kiểu HTTP `appHandler` riêng, có thêm giá trị trả về `error`:

	type appHandler func(http.ResponseWriter, *http.Request) error

Sau đó ta có thể đổi `viewRecord` để trả về lỗi:

	func viewRecord(w http.ResponseWriter, r *http.Request) error {
	    c := appengine.NewContext(r)
	    key := datastore.NewKey(c, "Record", r.FormValue("id"), 0, nil)
	    record := new(Record)
	    if err := datastore.Get(c, key, record); err != nil {
	        return err
	    }
	    return viewTemplate.Execute(w, record)
	}

Cách này đơn giản hơn phiên bản gốc,
nhưng gói [http](/pkg/net/http/) không hiểu
các hàm trả về `error`.
Để khắc phục, ta có thể hiện thực phương thức `ServeHTTP`
của interface `http.Handler` trên `appHandler`:

	func (fn appHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	    if err := fn(w, r); err != nil {
	        http.Error(w, err.Error(), 500)
	    }
	}

Phương thức `ServeHTTP` gọi hàm `appHandler` và hiển thị
lỗi mà nó trả về (nếu có) cho người dùng.
Lưu ý rằng receiver của phương thức, `fn`, là một hàm.
(Go làm được điều đó!) Phương thức gọi hàm bằng cách gọi chính receiver
trong biểu thức `fn(w, r)`.

Giờ khi đăng ký `viewRecord` với gói http, ta dùng hàm `Handle`
(thay vì `HandleFunc`) vì `appHandler` là một `http.Handler`
(không phải `http.HandlerFunc`).

	func init() {
	    http.Handle("/view", appHandler(viewRecord))
	}

Khi hạ tầng xử lý lỗi cơ bản này đã có, ta có thể làm nó thân thiện hơn.
Thay vì chỉ hiển thị chuỗi lỗi,
tốt hơn là đưa cho người dùng một thông báo lỗi đơn giản với mã HTTP phù hợp,
đồng thời ghi đầy đủ lỗi đó vào App Engine developer console để phục vụ gỡ lỗi.

Để làm điều đó, ta tạo một struct `appError` chứa một `error` và một số trường khác:

	type appError struct {
	    Error   error
	    Message string
	    Code    int
	}

Tiếp theo ta sửa kiểu `appHandler` để trả về các giá trị `*appError`:

	type appHandler func(http.ResponseWriter, *http.Request) *appError

(Thông thường việc trả về kiểu cụ thể của một lỗi thay vì `error` là sai lầm,
vì những lý do được thảo luận trong [Go FAQ](/doc/go_faq.html#nil_error),
nhưng ở đây lại là lựa chọn đúng vì `ServeHTTP` là nơi duy nhất
nhìn thấy giá trị đó và dùng nội dung của nó.)

Và làm cho `ServeHTTP` của `appHandler` hiển thị `Message` của `appError`
cho người dùng với `Code` HTTP chính xác, đồng thời ghi `Error` đầy đủ
lên developer console:

	func (fn appHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	    if e := fn(w, r); e != nil { // e is *appError, not os.Error.
	        c := appengine.NewContext(r)
	        c.Errorf("%v", e.Error)
	        http.Error(w, e.Message, e.Code)
	    }
	}

Cuối cùng, ta cập nhật `viewRecord` theo chữ ký hàm mới và để nó
trả về nhiều ngữ cảnh hơn khi gặp lỗi:

	func viewRecord(w http.ResponseWriter, r *http.Request) *appError {
	    c := appengine.NewContext(r)
	    key := datastore.NewKey(c, "Record", r.FormValue("id"), 0, nil)
	    record := new(Record)
	    if err := datastore.Get(c, key, record); err != nil {
	        return &appError{err, "Record not found", 404}
	    }
	    if err := viewTemplate.Execute(w, record); err != nil {
	        return &appError{err, "Can't display record", 500}
	    }
	    return nil
	}

Phiên bản này của `viewRecord` dài đúng bằng phiên bản gốc,
nhưng giờ mỗi dòng đều có ý nghĩa cụ thể và chúng ta đang mang lại
một trải nghiệm thân thiện hơn cho người dùng.

Và chưa dừng ở đó; ta còn có thể tiếp tục cải thiện việc xử lý lỗi trong ứng dụng. Một vài ý tưởng:

  - cho bộ xử lý lỗi một mẫu HTML đẹp,

  - làm việc gỡ lỗi dễ hơn bằng cách ghi stack trace vào phản hồi HTTP khi người dùng là quản trị viên,

  - viết một hàm dựng cho `appError` lưu stack trace để việc gỡ lỗi dễ hơn,

  - recover từ panic bên trong `appHandler`,
    ghi lỗi lên console ở mức “Critical”, đồng thời nói với người dùng rằng “đã xảy ra
    một lỗi nghiêm trọng.” Đây là một điểm chạm hay để tránh phơi bày cho người dùng
    những thông báo lỗi khó hiểu do lỗi lập trình gây ra.
    Xem bài viết [Defer, Panic, and Recover](/doc/articles/defer_panic_recover.html)
    để biết thêm chi tiết.

## Kết luận

Xử lý lỗi đúng cách là một yêu cầu thiết yếu của phần mềm tốt.
Bằng cách áp dụng các kỹ thuật được mô tả trong bài viết này, bạn sẽ có thể
viết mã Go đáng tin cậy hơn và cô đọng hơn.
