---
title: "Bộ giải mã GIF: một bài tập về interface trong Go"
date: 2011-05-25
by:
- Rob Pike
tags:
- gif
- gopher
- image
- interface
- lagomorph
- lzw
- moustache
- rodent
- technical
summary: Cách các interface của Go hoạt động rất đẹp trong bộ giải mã GIF của Go.
template: true
---

## Giới thiệu

Tại hội nghị Google I/O ở San Francisco vào ngày 10 tháng 5 năm 2011,
chúng tôi đã thông báo rằng ngôn ngữ Go hiện đã có mặt trên Google App Engine.
Go là ngôn ngữ đầu tiên được cung cấp trên App Engine mà biên dịch
trực tiếp ra mã máy,
điều này khiến nó trở thành lựa chọn tốt cho các tác vụ nặng CPU như xử lý ảnh.

Theo hướng đó, chúng tôi đã trình diễn một chương trình tên là [Moustachio](http://moustach-io.appspot.com/)
giúp việc cải thiện một bức ảnh như thế này trở nên dễ dàng:

{{image "gif-decoder/image00.jpg"}}

bằng cách thêm ria mép và chia sẻ kết quả:

{{image "gif-decoder/image02.jpg"}}

Toàn bộ xử lý đồ họa, bao gồm việc vẽ ria mép anti-aliased,
đều được thực hiện bởi một chương trình Go chạy trên App Engine.
(Mã nguồn có tại [dự án appengine-go](http://code.google.com/p/appengine-go/source/browse/example/moustachio/).)

Dù hầu hết ảnh trên web, ít nhất là những ảnh có khả năng được thêm ria mép, là JPEG,
vẫn có vô số định dạng khác đang hiện diện,
và có vẻ hợp lý khi Moustachio chấp nhận ảnh tải lên ở một vài định dạng đó.
Bộ giải mã JPEG và PNG đã có sẵn trong thư viện ảnh của Go,
nhưng định dạng GIF đáng kính thì chưa có,
nên chúng tôi quyết định viết một bộ giải mã GIF kịp cho buổi công bố.
Bộ giải mã đó có một vài phần thể hiện cách các interface của Go
giúp giải một số bài toán dễ dàng hơn.
Phần còn lại của bài viết này mô tả một vài trường hợp như vậy.

## Định dạng GIF

Trước tiên, hãy đi nhanh một vòng qua định dạng GIF. Một tệp ảnh GIF là _paletted_,
tức mỗi giá trị điểm ảnh là một chỉ số vào một bảng màu cố định được chứa trong tệp.
Định dạng GIF ra đời từ thời mà màn hình thường chỉ có tối đa 8
bit mỗi pixel,
và bảng màu được dùng để chuyển tập giá trị hữu hạn đó thành các bộ ba RGB (đỏ,
lục, lam) cần để phát sáng màn hình.
(Điều này trái ngược với JPEG, chẳng hạn,
vốn không có bảng màu vì mã hóa của nó biểu diễn riêng biệt các tín hiệu màu
riêng rẽ.)

Một ảnh GIF có thể có từ 1 đến 8 bit mỗi pixel, tính cả hai đầu, nhưng 8 bit mỗi pixel là phổ biến nhất.

Nói đơn giản hóa đi một chút, một tệp GIF chứa phần đầu định nghĩa độ sâu pixel
và kích thước ảnh,
một bảng màu (256 bộ ba RGB cho ảnh 8-bit),
và sau đó là dữ liệu pixel.
Dữ liệu pixel được lưu dưới dạng một dòng bit một chiều,
được nén bằng thuật toán LZW, rất hiệu quả cho đồ họa do máy tính tạo ra
dù không tốt lắm cho ảnh chụp.
Dữ liệu nén sau đó được chia thành các khối có độ dài kèm trước bằng một byte
đếm (0-255), theo sau là đúng số byte đó:

{{image "gif-decoder/image03.gif"}}

## Gỡ khối dữ liệu pixel

Để giải mã dữ liệu pixel GIF trong Go, ta có thể dùng bộ giải nén LZW từ
gói `compress/lzw`.
Nó có một hàm `NewReader` trả về một đối tượng mà,
như [tài liệu](/pkg/compress/lzw/#NewReader) nói,
“thỏa mãn các lần đọc bằng cách giải nén dữ liệu đọc từ r”:

	func NewReader(r io.Reader, order Order, litWidth int) io.ReadCloser

Ở đây `order` xác định thứ tự đóng gói bit và `litWidth` là độ rộng từ theo bit,
với tệp GIF thì tương ứng với độ sâu pixel, thường là 8.

Nhưng ta không thể cứ đưa trực tiếp tệp đầu vào cho `NewReader` làm đối số đầu tiên
vì bộ giải nén cần một dòng byte trong khi dữ liệu GIF là một dòng các khối
phải được mở gói ra.
Để xử lý điều này, ta có thể bọc `io.Reader` đầu vào bằng một đoạn mã để gỡ khối nó,
và làm cho đoạn mã đó một lần nữa hiện thực `Reader`.
Nói cách khác, ta đặt mã gỡ khối vào phương thức `Read` của một kiểu mới,
mà ta gọi là `blockReader`.

Đây là cấu trúc dữ liệu cho một `blockReader`.

	type blockReader struct {
	   r     reader    // Input source; implements io.Reader and io.ByteReader.
	   slice []byte    // Buffer of unread data.
	   tmp   [256]byte // Storage for slice.
	}

Reader `r` sẽ là nguồn dữ liệu ảnh,
có thể là một tệp hoặc kết nối HTTP.
Các trường `slice` và `tmp` sẽ được dùng để quản lý việc gỡ khối.
Đây là toàn bộ phương thức `Read`.
Nó là một ví dụ hay về cách dùng slice và array trong Go.

	1  func (b *blockReader) Read(p []byte) (int, os.Error) {
	2      if len(p) == 0 {
	3          return 0, nil
	4      }
	5      if len(b.slice) == 0 {
	6          blockLen, err := b.r.ReadByte()
	7          if err != nil {
	8              return 0, err
	9          }
	10          if blockLen == 0 {
	11              return 0, os.EOF
	12          }
	13          b.slice = b.tmp[0:blockLen]
	14          if _, err = io.ReadFull(b.r, b.slice); err != nil {
	15              return 0, err
	16          }
	17      }
	18      n := copy(p, b.slice)
	19      b.slice = b.slice[n:]
	20      return n, nil
	21  }

Dòng 2-4 chỉ là một kiểm tra an toàn: nếu không có chỗ nào để đặt dữ liệu, hãy trả về zero.
Điều đó đáng lẽ không bao giờ xảy ra, nhưng cẩn thận vẫn hơn.

Dòng 5 hỏi xem có dữ liệu còn sót lại từ lần gọi trước hay không bằng cách kiểm tra
độ dài của `b.slice`.
Nếu không có, slice sẽ có độ dài zero và ta cần đọc
khối tiếp theo từ `r`.

Một khối GIF bắt đầu bằng một byte đếm, được đọc ở dòng 6.
Nếu giá trị đếm là zero, GIF định nghĩa đây là khối kết thúc,
nên ta trả về `EOF` ở dòng 11.

Giờ ta biết cần đọc `blockLen` byte,
nên ta trỏ `b.slice` tới `blockLen` byte đầu của `b.tmp` rồi
dùng hàm trợ giúp `io.ReadFull` để đọc đúng số byte đó.
Hàm đó sẽ trả về lỗi nếu không thể đọc chính xác столько byte,
điều đáng lẽ không bao giờ xảy ra.
Nếu không, ta có `blockLen` byte sẵn sàng để đọc.

Dòng 18-19 sao chép dữ liệu từ `b.slice` sang bộ đệm của bên gọi.
Ta đang hiện thực `Read`, không phải `ReadFull`,
nên ta được phép trả về ít byte hơn số lượng yêu cầu.
Điều đó khiến việc hiện thực rất dễ: ta chỉ cần chép dữ liệu từ `b.slice` sang bộ đệm của bên gọi (`p`),
và giá trị trả về từ `copy` là số byte đã chuyển.
Sau đó ta reslice `b.slice` để bỏ đi `n` byte đầu,
sẵn sàng cho lần gọi tiếp theo.

Trong lập trình Go, đây là một kỹ thuật hay: ghép một slice (`b.slice`) với một array (`b.tmp`).
Trong trường hợp này, điều đó có nghĩa là phương thức `Read` của kiểu `blockReader` không bao giờ phải cấp phát.
Nó cũng có nghĩa là ta không cần giữ một biến đếm riêng (nó đã hàm ý trong độ dài slice),
và hàm dựng sẵn `copy` bảo đảm rằng ta không bao giờ chép quá mức cần thiết.
(Để tìm hiểu thêm về slice, xem [bài viết này trên Go Blog](/blog/go-slices-usage-and-internals).)

Với kiểu `blockReader`, ta có thể mở khối dòng dữ liệu ảnh
chỉ bằng cách bọc reader đầu vào,
giả sử là một tệp, như sau:

	deblockingReader := &blockReader{r: imageFile}

Việc bọc này biến một dòng ảnh GIF được phân khối thành một dòng byte đơn giản
có thể truy cập qua các lần gọi tới phương thức `Read` của `blockReader`.

## Nối các mảnh lại

Khi `blockReader` đã được hiện thực và bộ nén LZW đã có trong thư viện,
ta đã có mọi mảnh cần để giải mã dòng dữ liệu ảnh.
Ta khâu chúng lại với cú sấm này,
lấy thẳng từ mã:

	lzwr := lzw.NewReader(&blockReader{r: d.r}, lzw.LSB, int(litWidth))
	if _, err = io.ReadFull(lzwr, m.Pix); err != nil {
	   break
	}

Chỉ vậy thôi.

Dòng đầu tiên tạo một `blockReader` rồi truyền nó cho `lzw.NewReader`
để tạo một bộ giải nén.
Ở đây `d.r` là `io.Reader` giữ dữ liệu ảnh,
`lzw.LSB` định nghĩa thứ tự byte trong bộ giải nén LZW,
và `litWidth` là độ sâu pixel.

Với bộ giải nén đó, dòng thứ hai gọi `io.ReadFull` để giải nén
dữ liệu và lưu nó vào ảnh, `m.Pix`.
Khi `ReadFull` trả về, dữ liệu ảnh đã được giải nén và được lưu trong ảnh
`m`, sẵn sàng để hiển thị.

Đoạn mã này chạy đúng ngay lần đầu tiên. Thật sự đấy.

Ta có thể tránh biến tạm `lzwr` bằng cách đặt lời gọi `NewReader`
vào ngay danh sách đối số của `ReadFull`,
giống như cách ta tạo `blockReader` bên trong lời gọi `NewReader`,
nhưng như vậy có lẽ nhồi quá nhiều thứ vào một dòng.

## Kết luận

Các interface của Go giúp việc xây phần mềm bằng cách ghép các mảnh nhỏ
như thế này để tái cấu trúc dữ liệu trở nên dễ dàng.
Trong ví dụ này, ta hiện thực giải mã GIF bằng cách móc nối một bộ gỡ khối
và một bộ giải nén qua interface `io.Reader`,
tương tự một đường ống Unix có kiểm tra kiểu.
Ngoài ra, ta đã viết bộ gỡ khối như một hiện thực (ngầm) của interface `Reader`,
vì thế không cần bất kỳ khai báo hay boilerplate nào thêm để đưa nó vào
pipeline xử lý.
Thật khó để hiện thực bộ giải mã này vừa gọn vừa sạch và an toàn đến vậy trong hầu hết ngôn ngữ,
nhưng cơ chế interface cộng với một vài quy ước khiến điều đó gần như tự nhiên trong Go.

Điều đó xứng đáng với thêm một bức ảnh nữa, lần này là GIF:

{{image "gif-decoder/image01.gif"}}

Định dạng GIF được định nghĩa tại [http://www.w3.org/Graphics/GIF/spec-gif89a.txt](http://www.w3.org/Graphics/GIF/spec-gif89a.txt).
