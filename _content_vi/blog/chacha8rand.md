---
title: "Tính ngẫu nhiên an toàn trong Go 1.22"
date: 2024-05-02
by:
- Russ Cox
- Filippo Valsorda
summary: ChaCha8Rand là một bộ sinh số giả ngẫu nhiên an toàn về mặt mật mã mới được dùng trong Go 1.22.
template: true
---

Máy tính không ngẫu nhiên.
Ngược lại, các nhà thiết kế phần cứng phải làm việc rất vất vả để bảo đảm máy tính chạy mọi chương trình theo cùng một cách mỗi lần.
Vì vậy khi một chương trình cần số ngẫu nhiên, điều đó đòi hỏi thêm công sức.
Theo truyền thống, các nhà khoa học máy tính và các ngôn ngữ lập trình
phân biệt hai loại số ngẫu nhiên khác nhau:
ngẫu nhiên thống kê và ngẫu nhiên mật mã.
Trong Go, hai loại đó lần lượt được cung cấp bởi [`math/rand`](/pkg/math/rand/)
và [`crypto/rand`](/pkg/crypto/rand).
Bài viết này nói về cách Go 1.22 đưa hai loại này lại gần nhau hơn,
bằng cách dùng một nguồn ngẫu nhiên mật mã trong `math/rand`
(cũng như `math/rand/v2`, như đã nhắc trong [bài viết trước](/blog/randv2)).
Kết quả là tính ngẫu nhiên tốt hơn và ít thiệt hại hơn rất nhiều khi
lập trình viên vô tình dùng `math/rand` thay vì `crypto/rand`.

Trước khi giải thích Go 1.22 đã làm gì, hãy xem kỹ hơn
sự khác biệt giữa ngẫu nhiên thống kê và ngẫu nhiên mật mã.

## Tính ngẫu nhiên thống kê

Các số ngẫu nhiên vượt qua những kiểm thử thống kê cơ bản
thường phù hợp cho các trường hợp như mô phỏng, lấy mẫu,
phân tích số, thuật toán ngẫu nhiên không liên quan mật mã,
[kiểm thử ngẫu nhiên](/doc/security/fuzz/),
[xáo trộn đầu vào](https://en.wikipedia.org/wiki/Fisher%E2%80%93Yates_shuffle),
và
[backoff lũy thừa có ngẫu nhiên](https://en.wikipedia.org/wiki/Exponential_backoff#Collision_avoidance).
Những công thức toán học rất cơ bản, rất dễ tính
lại hóa ra hoạt động đủ tốt cho những trường hợp này.
Tuy nhiên, vì phương pháp quá đơn giản, một người quan sát
biết thuật toán đang dùng thường có thể dự đoán phần còn lại
của chuỗi sau khi thấy đủ số lượng giá trị.

Gần như mọi môi trường lập trình đều cung cấp một cơ chế sinh
số ngẫu nhiên thống kê có nguồn gốc ngược qua C tới
Research Unix Third Edition (V3), bản đã bổ sung cặp hàm `srand` và `rand`.
Trang hướng dẫn của chúng có ghi chú:

> _CẢNH BÁO   Tác giả của thủ tục này đã viết
bộ sinh số ngẫu nhiên trong nhiều năm và chưa từng
được biết đến là đã viết ra cái nào hoạt động đúng._

Ghi chú này một phần là câu đùa nhưng cũng là sự thừa nhận rằng các bộ sinh như thế
[về bản chất không hề ngẫu nhiên](https://www.tuhs.org/pipermail/tuhs/2024-March/029587.html).

Mã nguồn của bộ sinh cho thấy rõ nó đơn giản đến mức nào.
Khi dịch từ assembly PDP-11 sang C hiện đại, nó như sau:

	uint16 ranx;

	void
	srand(uint16 seed)
	{
	    ranx = seed;
	}

	int16
	rand(void)
	{
	    ranx = 13077*ranx + 6925;
	    return ranx & ~0x8000;
	}

Việc gọi `srand` sẽ khởi tạo bộ sinh bằng một số nguyên duy nhất,
và `rand` trả về số tiếp theo từ bộ sinh.
Phép AND trong câu lệnh return xóa bit dấu để chắc chắn kết quả là số dương.

Hàm này là một trường hợp cụ thể của lớp
[linear congruential generator (LCG)](https://en.wikipedia.org/wiki/Linear_congruential_generator),
được Knuth phân tích trong _The Art of Computer Programming_, Tập 2, mục 3.2.1.
Lợi ích chính của LCG là có thể chọn hằng số sao cho chúng
phát ra mọi giá trị đầu ra có thể đúng một lần trước khi lặp lại,
như cách hiện thực Unix đã làm với đầu ra 15 bit.
Tuy vậy, một vấn đề nghiêm trọng của LCG là các bit cao của trạng thái hoàn toàn không ảnh hưởng đến bit thấp,
vì vậy mọi phép cắt ngắn chuỗi xuống còn _k_ bit đều tất yếu lặp lại với chu kỳ ngắn hơn.
Bit thấp nhất buộc phải đổi qua lại: 0, 1, 0, 1, 0, 1.
Hai bit thấp nhất phải đếm lên hoặc xuống: 0, 1, 2, 3, 0, 1, 2, 3, hoặc 0, 3, 2, 1, 0, 3, 2, 1.
Có bốn chuỗi ba bit có thể có; hiện thực Unix gốc lặp lại 0, 5, 6, 3, 4, 1, 2, 7.
(Có thể tránh các vấn đề này bằng cách giảm giá trị theo modulo một số nguyên tố,
nhưng điều đó vào thời điểm ấy khá tốn kém.
Xem bài báo CACM năm 1988 của S. K. Park và K. W. Miller,
“[Random number generators: good ones are hard to find](https://dl.acm.org/doi/10.1145/63039.63042)”
để có một phân tích ngắn,
và chương đầu của Knuth Tập 2 để có phân tích dài hơn.)

Ngay cả với những vấn đề đã biết này,
các hàm `srand` và `rand` vẫn được đưa vào tiêu chuẩn C đầu tiên,
và chức năng tương đương cũng xuất hiện trong hầu như mọi ngôn ngữ kể từ đó.
LCG từng là chiến lược triển khai thống trị,
dù chúng dần mất đi sự ưa chuộng vì một số nhược điểm quan trọng.
Một trường hợp sử dụng còn đáng kể là [`java.util.Random`](https://github.com/openjdk/jdk8u-dev/blob/master/jdk/src/share/classes/java/util/Random.java),
thứ đứng sau [`java.lang.Math.random`](https://github.com/openjdk/jdk8u-dev/blob/master/jdk/src/share/classes/java/util/Random.java).

Một điều khác bạn có thể thấy từ hiện thực ở trên
là trạng thái nội bộ bị lộ hoàn toàn qua kết quả của `rand`.
Người quan sát biết thuật toán và nhìn thấy chỉ một kết quả
cũng có thể dễ dàng tính được mọi kết quả tương lai.
Nếu bạn đang vận hành một máy chủ tính ra một số giá trị ngẫu nhiên
trở thành công khai và một số giá trị ngẫu nhiên khác phải được giữ bí mật,
thì dùng kiểu bộ sinh này sẽ là thảm họa:
các bí mật sẽ không còn là bí mật.

Các bộ sinh ngẫu nhiên hiện đại hơn không tệ như bộ sinh Unix ban đầu,
nhưng chúng vẫn không hoàn toàn không thể đoán trước.
Để minh họa điều đó, tiếp theo chúng ta sẽ xem bộ sinh `math/rand` gốc của Go 1
và bộ sinh PCG mà chúng tôi thêm vào `math/rand/v2`.

## Bộ sinh Go 1

Bộ sinh dùng trong `math/rand` của Go 1 là một biến thể của
[linear-feedback shift register](https://en.wikipedia.org/wiki/Linear-feedback_shift_register).
Thuật toán này dựa trên ý tưởng của George Marsaglia,
được Don Mitchell và Jim Reeds điều chỉnh,
và sau đó được Ken Thompson tùy biến thêm cho Plan 9 rồi Go.
Nó không có tên chính thức, nên bài viết này gọi nó là bộ sinh Go 1.

Trạng thái nội bộ của bộ sinh Go 1 là một slice `vec` gồm 607 giá trị uint64.
Trong slice đó có hai phần tử đặc biệt: `vec[606]`, phần tử cuối, được gọi là “tap”,
và `vec[334]` được gọi là “feed”.
Để sinh số ngẫu nhiên tiếp theo,
bộ sinh cộng tap và feed
để tạo ra một giá trị `x`,
ghi `x` trở lại feed,
dịch toàn bộ slice sang phải một vị trí
(tap chuyển đến `vec[0]` và `vec[i]` chuyển tới `vec[i+1]`),
và trả về `x`.
Bộ sinh được gọi là “linear feedback” vì tap được _cộng_ vào feed;
toàn bộ trạng thái là một “shift register” vì mỗi bước đều dịch các phần tử trong slice.

Dĩ nhiên, việc thật sự di chuyển mọi phần tử của slice ở mỗi bước sẽ quá tốn kém,
nên hiện thực để nguyên dữ liệu slice tại chỗ
và dịch vị trí tap và feed lùi lại
sau mỗi bước. Mã trông như sau:

{{raw `
	func (r *rngSource) Uint64() uint64 {
		r.tap--
		if r.tap < 0 {
			r.tap += len(r.vec)
		}

		r.feed--
		if r.feed < 0 {
			r.feed += len(r.vec)
		}

		x := r.vec[r.feed] + r.vec[r.tap]
		r.vec[r.feed] = x
		return uint64(x)
	}
`}}

Việc sinh số tiếp theo khá rẻ: hai phép trừ, hai phép cộng có điều kiện, hai lần tải, một phép cộng và một lần ghi.

Thật không may, vì bộ sinh trả trực tiếp một phần tử slice từ vector trạng thái nội bộ,
đọc 607 giá trị từ bộ sinh sẽ bộc lộ toàn bộ trạng thái của nó.
Với các giá trị đó, bạn có thể dự đoán tất cả giá trị tương lai, bằng cách điền `vec`
trong bản sao của riêng bạn rồi chạy thuật toán.
Bạn cũng có thể khôi phục tất cả giá trị trước đó, bằng cách chạy thuật toán ngược lại
(lấy feed trừ tap rồi dịch slice sang trái).

Để minh họa đầy đủ, đây là một [chương trình không an toàn](/play/p/v0QdGjUAtzC)
sinh token xác thực giả ngẫu nhiên
cùng với mã dự đoán token tiếp theo dựa trên một chuỗi token trước đó.
Như bạn có thể thấy, bộ sinh Go 1 hoàn toàn không cung cấp bảo mật nào (và cũng không nhằm mục đích đó).
Chất lượng của các số được sinh ra cũng phụ thuộc vào cách thiết lập ban đầu của `vec`.

## Bộ sinh PCG

Đối với `math/rand/v2`, chúng tôi muốn cung cấp một bộ sinh ngẫu nhiên thống kê hiện đại hơn
và đã chọn thuật toán PCG của Melissa O'Neill, được công bố năm 2014 trong bài báo của bà
“[PCG: A Family of Simple Fast Space-Efficient Statistically Good Algorithms for Random Number Generation](https://www.pcg-random.org/pdf/hmc-cs-2014-0905.pdf)”.
Phân tích rất đầy đủ trong bài báo có thể khiến người đọc khó nhận ra ngay từ cái nhìn đầu tiên việc các bộ sinh này đơn giản đến mức nào:
PCG là một LCG 128 bit có hậu xử lý.

Nếu trạng thái `p.x` là một `uint128` (giả định), mã để tính giá trị tiếp theo sẽ là:

	const (
		pcgM = 0x2360ed051fc65da44385df649fccf645
		pcgA = 0x5851f42d4c957f2d14057b7ef767814f
	)

	type PCG struct {
		x uint128
	}

	func (p *PCG) Uint64() uint64 {
		p.x = p.x * pcgM + pcgA
		return scramble(p.x)
	}

Toàn bộ trạng thái chỉ là một số 128 bit duy nhất,
và bước cập nhật là một phép nhân rồi cộng 128 bit.
Trong câu lệnh return, hàm `scramble` rút trạng thái 128 bit
xuống còn trạng thái 64 bit.
PCG nguyên bản dùng (một lần nữa, với kiểu `uint128` giả định):

	func scramble(x uint128) uint64 {
		return bits.RotateLeft(uint64(x>>64) ^ uint64(x), -int(x>>122))
	}

Mã này XOR hai nửa của trạng thái 128 bit
rồi xoay kết quả theo sáu bit cao nhất của trạng thái.
Phiên bản này được gọi là PCG-XSL-RR, viết tắt của “xor shift low, right rotate”.

Dựa trên [một gợi ý từ O'Neill trong quá trình thảo luận đề xuất](/issue/21835#issuecomment-739065688),
PCG của Go dùng một hàm scramble mới dựa trên phép nhân,
trộn bit mạnh tay hơn:

	func scramble(x uint128) uint64 {
		hi, lo := uint64(x>>64), uint64(x)
		hi ^= hi >> 32
		hi *= 0xda942042e4dd58b5
		hi ^= hi >> 48
		hi *= lo | 1
	}

O'Neill gọi PCG với bộ trộn này là PCG-DXSM, viết tắt của “double xorshift multiply.”
Numpy cũng dùng dạng PCG này.

Dù PCG cần nhiều phép tính hơn để sinh mỗi giá trị,
nó dùng trạng thái ít hơn đáng kể: hai uint64 thay vì 607.
Nó cũng ít nhạy cảm hơn nhiều với các giá trị khởi tạo của trạng thái đó,
và [vượt qua nhiều kiểm thử thống kê mà các bộ sinh khác không vượt qua được](https://www.pcg-random.org/statistical-tests.html).
Ở nhiều khía cạnh, đây là một bộ sinh thống kê lý tưởng.

Dẫu vậy, PCG không phải là không thể đoán trước.
Mặc dù việc xáo trộn bit để chuẩn bị kết quả không làm lộ
trạng thái trực tiếp như ở LCG và bộ sinh Go 1,
[PCG-XSL-RR vẫn có thể bị đảo ngược](https://pdfs.semanticscholar.org/4c5e/4a263d92787850edd011d38521966751a179.pdf),
và sẽ không có gì bất ngờ nếu PCG-DXSM cũng vậy.
Đối với bí mật, chúng ta cần thứ gì đó khác.

## Tính ngẫu nhiên mật mã

_Số ngẫu nhiên mật mã_ trong thực tế cần phải hoàn toàn không thể đoán trước,
ngay cả với người quan sát biết cách chúng được tạo ra
và đã quan sát bất kỳ số lượng giá trị sinh trước đó nào.
Sự an toàn của các giao thức mật mã, khóa bí mật, thương mại hiện đại,
quyền riêng tư trực tuyến và nhiều thứ khác đều phụ thuộc thiết yếu vào khả năng truy cập
tính ngẫu nhiên mật mã.

Việc cung cấp tính ngẫu nhiên mật mã cuối cùng là công việc của
hệ điều hành, hệ thống có thể thu thập sự ngẫu nhiên thực từ thiết bị vật lý, từ thời gian
của chuột, bàn phím, đĩa và mạng, và gần đây hơn là
[nhiễu điện được CPU đo trực tiếp](https://web.archive.org/web/20141230024150/http://www.cryptography.com/public/pdf/Intel_TRNG_Report_20120312.pdf).
Khi hệ điều hành đã thu thập được một lượng ngẫu nhiên đủ ý nghĩa,
ví dụ ít nhất 256 bit, nó có thể dùng các thuật toán
băm hoặc mã hóa mật mã để kéo dãn hạt giống đó thành
một chuỗi số ngẫu nhiên dài tùy ý.
(Trên thực tế, hệ điều hành cũng liên tục thu thập và
bổ sung thêm ngẫu nhiên vào chuỗi này.)

Các giao diện hệ điều hành chính xác đã thay đổi theo thời gian.
Mười năm trước, hầu hết hệ thống cung cấp một tệp thiết bị có tên
`/dev/random` hoặc thứ tương tự.
Ngày nay, vì nhận thức được tính nền tảng của ngẫu nhiên,
hệ điều hành cung cấp trực tiếp một system call thay thế.
(Điều này cũng cho phép chương trình đọc dữ liệu ngẫu nhiên ngay cả
khi bị cắt khỏi hệ thống tệp.)
Trong Go, package [`crypto/rand`](/pkg/crypto/rand/) trừu tượng hóa các chi tiết đó,
cung cấp cùng một giao diện trên mọi hệ điều hành: [`rand.Read`](/pkg/crypto/rand/#Read).

Sẽ không thực tế nếu `math/rand` phải hỏi hệ điều hành để lấy
ngẫu nhiên mỗi khi cần một `uint64`.
Nhưng chúng ta có thể dùng kỹ thuật mật mã để định nghĩa một bộ sinh ngẫu nhiên trong tiến trình
vượt trội hơn LCG, bộ sinh Go 1, và thậm chí cả PCG.

## Bộ sinh ChaCha8Rand

Bộ sinh mới của chúng tôi, được đặt cái tên khá đơn điệu là ChaCha8Rand cho mục đích đặc tả
và được hiện thực thành [`rand.ChaCha8`](/pkg/math/rand/v2/#ChaCha8) của `math/rand/v2`,
là một phiên bản chỉnh sửa nhẹ của [mã dòng ChaCha](https://cr.yp.to/chacha.html) của Daniel J. Bernstein.
ChaCha được dùng rộng rãi ở dạng 20 vòng gọi là ChaCha20, bao gồm trong TLS và SSH.
Bài báo “[Too Much Crypto](https://eprint.iacr.org/2019/1492.pdf)” của Jean-Philippe Aumasson
lập luận khá thuyết phục rằng dạng 8 vòng ChaCha8 cũng an toàn (và nhanh hơn khoảng 2.5 lần).
Chúng tôi dùng ChaCha8 làm lõi cho ChaCha8Rand.

Hầu hết các mã dòng, bao gồm ChaCha8, hoạt động bằng cách định nghĩa một hàm nhận
một khóa và một số hiệu khối rồi tạo ra một khối dữ liệu trông có vẻ ngẫu nhiên với kích thước cố định.
Tiêu chuẩn mật mã mà chúng hướng tới (và thường đáp ứng) là đầu ra này không thể phân biệt
với dữ liệu ngẫu nhiên thực khi không có một cuộc tìm kiếm brute force đắt đỏ theo cấp số mũ nào đó.
Một thông điệp được mã hóa hoặc giải mã bằng cách XOR các khối dữ liệu đầu vào liên tiếp
với các khối được sinh ngẫu nhiên liên tiếp.
Để dùng ChaCha8 như một `rand.Source`,
chúng tôi dùng trực tiếp các khối đã sinh thay vì XOR chúng với dữ liệu đầu vào
(điều này tương đương với việc mã hóa hoặc giải mã toàn số không).

Chúng tôi đã thay đổi một vài chi tiết để ChaCha8Rand phù hợp hơn với việc sinh số ngẫu nhiên. Tóm tắt:

 - ChaCha8Rand nhận một hạt giống dài 32 byte, được dùng làm khóa ChaCha8.
 - ChaCha8 sinh ra các khối 64 byte, với phép tính coi mỗi khối là 16 `uint32`.
   Một hiện thực phổ biến là tính bốn khối cùng lúc bằng [chỉ thị SIMD](https://en.wikipedia.org/wiki/Single_instruction,_multiple_data)
   trên 16 thanh ghi véc-tơ gồm bốn `uint32` mỗi thanh ghi.
   Điều này sinh ra bốn khối đan xen, phải được tháo trộn lại trước khi XOR với dữ liệu vào.
   ChaCha8Rand định nghĩa chính các khối đan xen đó là luồng dữ liệu ngẫu nhiên,
   loại bỏ chi phí của bước tháo trộn.
   (Về mặt bảo mật, có thể xem đây là ChaCha8 chuẩn rồi tới một bước trộn lại.)
 - ChaCha8 kết thúc một khối bằng cách cộng một số giá trị vào từng `uint32` trong khối.
   Một nửa các giá trị là vật liệu khóa, nửa còn lại là các hằng số đã biết.
   ChaCha8Rand định nghĩa rằng các hằng số đã biết sẽ không được cộng lại,
   loại bỏ một nửa số phép cộng cuối cùng.
   (Về mặt bảo mật, có thể xem đây là ChaCha8 chuẩn rồi trừ đi các hằng số đã biết.)
 - Mỗi khối sinh ra thứ 16, ChaCha8Rand lấy 32 byte cuối của khối cho riêng nó,
   biến chúng thành khóa cho 16 khối tiếp theo.
   Điều này cung cấp một dạng [forward secrecy](https://en.wikipedia.org/wiki/Forward_secrecy):
   nếu một hệ thống bị xâm phạm bởi cuộc tấn công có thể
   khôi phục toàn bộ trạng thái bộ nhớ của bộ sinh, chỉ các giá trị được sinh
   kể từ lần thay khóa gần nhất mới có thể bị khôi phục. Quá khứ thì không thể tiếp cận.
   ChaCha8Rand theo định nghĩa đến đây buộc phải sinh 4 khối một lúc,
   nhưng chúng tôi chọn xoay khóa mỗi 16 khối để vẫn mở khả năng cho
   các hiện thực nhanh hơn dùng véc-tơ 256 bit hoặc 512 bit,
   có thể sinh 8 hoặc 16 khối một lúc.

Chúng tôi đã viết và công bố [đặc tả C2SP cho ChaCha8Rand](https://c2sp.org/chacha8rand),
cùng với các ca kiểm thử.
Điều này sẽ cho phép các hiện thực khác chia sẻ tính lặp lại với hiện thực Go
cho cùng một hạt giống.

Go runtime hiện duy trì một trạng thái ChaCha8Rand cho mỗi lõi (300 byte),
được khởi tạo bằng tính ngẫu nhiên mật mã do hệ điều hành cung cấp,
để có thể sinh số ngẫu nhiên nhanh mà không có tranh chấp khóa.
Việc dành riêng 300 byte cho mỗi lõi nghe có vẻ đắt,
nhưng trên hệ thống 16 lõi, nó xấp xỉ việc lưu một trạng thái bộ sinh Go 1 dùng chung duy nhất (4.872 byte).
Tốc độ xứng đáng với chi phí bộ nhớ đó.
Bộ sinh ChaCha8Rand theo mỗi lõi này hiện được dùng ở ba nơi khác nhau trong thư viện chuẩn Go:

 1. Các hàm của package `math/rand/v2`, như
   [`rand.Float64`](/pkg/math/rand/v2/#Float64) và
   [`rand.N`](/pkg/math/rand/v2/#N), luôn dùng ChaCha8Rand.

 2. Các hàm của package `math/rand`, như
   [`rand.Float64`](/pkg/math/rand/#Float64) và
   [`rand.Intn`](/pkg/math/rand/#Intn),
   dùng ChaCha8Rand khi
   [`rand.Seed`](/pkg/math/rand/#Seed) chưa được gọi.
   Việc áp dụng ChaCha8Rand vào `math/rand` cải thiện độ an toàn cho chương trình
   ngay cả trước khi chúng cập nhật lên `math/rand/v2`,
   miễn là chúng không gọi `rand.Seed`.
   (Nếu `rand.Seed` được gọi, hiện thực bắt buộc phải quay về bộ sinh Go 1 để tương thích.)

 3. Runtime chọn hạt giống băm cho từng map mới
    bằng ChaCha8Rand thay vì một [bộ sinh dựa trên wyrand](https://github.com/wangyi-fudan/wyhash)
    kém an toàn hơn mà trước đây nó sử dụng.
    Hạt giống ngẫu nhiên là cần thiết bởi vì nếu
    kẻ tấn công biết chính xác hàm băm được hiện thực map dùng,
    họ có thể chuẩn bị đầu vào khiến map rơi vào hành vi bậc hai
    (xem Crosby và Wallach, “[Denial of Service via Algorithmic Complexity Attacks](https://www.usenix.org/conference/12th-usenix-security-symposium/denial-service-algorithmic-complexity-attacks)”).
    Việc dùng hạt giống cho từng map, thay vì một hạt giống toàn cục cho mọi map,
    cũng tránh được [các hành vi suy biến khác](https://accidentallyquadratic.tumblr.com/post/153545455987/rust-hash-iteration-reinsertion).
    Không hoàn toàn rõ ràng rằng map có cần hạt giống ngẫu nhiên theo chuẩn mật mã hay không,
    nhưng cũng không rõ là không cần. Chúng tôi thấy thận trọng hơn và việc đổi sang đó thì rất đơn giản.

Mã cần các thể hiện ChaCha8Rand riêng của mình có thể tạo trực tiếp [`rand.ChaCha8`](/pkg/math/rand/v2/#ChaCha8).

## Sửa những sai lầm bảo mật

Go hướng tới việc giúp lập trình viên viết mã an toàn theo mặc định.
Khi chúng tôi quan sát thấy một lỗi phổ biến có hệ quả bảo mật,
chúng tôi tìm cách giảm rủi ro của lỗi đó
hoặc loại bỏ nó hoàn toàn.
Trong trường hợp này, bộ sinh toàn cục của `math/rand` quá dễ đoán,
dẫn tới những vấn đề nghiêm trọng trong nhiều bối cảnh.

Ví dụ, khi Go 1.20 ngừng khuyến khích dùng [`Read` của `math/rand`](/pkg/math/rand/#Read),
chúng tôi nhận được phản hồi từ các lập trình viên phát hiện ra (nhờ công cụ chỉ ra
việc dùng chức năng đã bị deprecate) rằng họ đã
dùng nó ở những nơi mà [`Read` của `crypto/rand`](/pkg/crypto/rand/#Read)
rõ ràng mới là thứ cần thiết, như tạo vật liệu khóa.
Trong Go 1.20, sai lầm đó
là một vấn đề bảo mật nghiêm trọng cần điều tra kỹ
để hiểu thiệt hại.
Các khóa đã được dùng ở đâu?
Chúng bị lộ như thế nào?
Những đầu ra ngẫu nhiên khác có bị lộ khiến kẻ tấn công suy ra được khóa không?
Và vân vân.
Trong Go 1.22, sai lầm đó chỉ còn là một sai lầm.
Tất nhiên vẫn tốt hơn nếu dùng `crypto/rand`,
vì kernel hệ điều hành có thể làm tốt hơn trong việc giữ giá trị ngẫu nhiên
bí mật khỏi nhiều kiểu con mắt tò mò,
kernel liên tục thêm entropy mới vào bộ sinh của nó,
và kernel đã trải qua nhiều kiểm chứng hơn.
Nhưng việc vô tình dùng `math/rand` không còn là thảm họa bảo mật nữa.

Cũng có nhiều trường hợp sử dụng trông không giống “crypto”
nhưng vẫn cần tính ngẫu nhiên không thể đoán trước.
Những trường hợp này trở nên vững vàng hơn khi dùng ChaCha8Rand thay cho bộ sinh Go 1.

Ví dụ, hãy xét việc sinh một
[UUID ngẫu nhiên](https://en.wikipedia.org/wiki/Universally_unique_identifier#Version_4_(random)).
Vì UUID không phải bí mật, việc dùng `math/rand` có vẻ ổn.
Nhưng nếu `math/rand` được khởi tạo bằng thời gian hiện tại,
thì việc chạy nó cùng một lúc trên các máy tính khác nhau
sẽ cho cùng một giá trị, khiến chúng không còn “duy nhất toàn cục”.
Điều này đặc biệt dễ xảy ra trên hệ thống mà thời gian hiện tại
chỉ có độ chính xác tới mili giây.
Ngay cả với việc tự động khởi tạo bằng entropy do hệ điều hành cung cấp,
như được đưa vào từ Go 1.20,
hạt giống của bộ sinh Go 1 vẫn chỉ là số nguyên 63 bit,
vì thế chương trình sinh UUID khi khởi động
chỉ có thể sinh ra 2⁶³ UUID đầu tiên có thể có và
dễ thấy va chạm sau khoảng 2³¹ UUID.
Trong Go 1.22, bộ sinh ChaCha8Rand mới
được khởi tạo từ 256 bit entropy và có thể sinh
2²⁵⁶ UUID đầu tiên có thể có.
Nó không phải lo về va chạm.

Ví dụ khác, hãy xét việc cân bằng tải trong một front-end server
ngẫu nhiên gán các request đến cho các back-end server.
Nếu kẻ tấn công có thể quan sát các lần gán và biết
thuật toán dễ đoán tạo ra chúng,
thì kẻ tấn công có thể gửi một dòng
request phần lớn là rẻ nhưng sắp xếp để mọi request đắt đỏ
đều rơi vào một back-end server duy nhất.
Đây là một vấn đề không quá thường gặp nhưng vẫn khả dĩ nếu dùng bộ sinh Go 1.
Trong Go 1.22, nó hoàn toàn không còn là vấn đề.

Trong tất cả ví dụ này, Go 1.22 đã loại bỏ hoặc giảm mạnh
các vấn đề bảo mật.

## Hiệu năng

Lợi ích bảo mật của ChaCha8Rand đúng là có một chi phí nhỏ,
nhưng ChaCha8Rand vẫn nằm trong cùng mặt bằng với cả bộ sinh Go 1 lẫn PCG.
Các biểu đồ sau so sánh hiệu năng của ba bộ sinh,
trên nhiều loại phần cứng khác nhau, khi chạy hai phép toán:
phép toán nguyên thủy “Uint64”, trả về `uint64` tiếp theo trong luồng ngẫu nhiên,
và phép toán cấp cao hơn “N(1000)”, trả về một giá trị ngẫu nhiên trong đoạn [0, 1000).

<div style="background-color: white;">
<img src="chacha8rand/amd.svg">
<img src="chacha8rand/intel.svg">
<img src="chacha8rand/amd32.svg">
<img src="chacha8rand/intel32.svg">
<img src="chacha8rand/m1.svg">
<img src="chacha8rand/m3.svg">
<img src="chacha8rand/taut2a.svg">
</div>

Các biểu đồ “running 32-bit code” cho thấy chip x86 64-bit hiện đại
đang thực thi mã được build với `GOARCH=386`, nghĩa là chúng
đang chạy ở chế độ 32 bit.
Trong trường hợp đó, việc PCG đòi hỏi phép nhân 128 bit
khiến nó chậm hơn ChaCha8Rand, thứ chỉ dùng số học SIMD 32 bit.
Các hệ thống 32 bit thật ngày càng ít quan trọng theo từng năm,
nhưng việc ChaCha8Rand nhanh hơn PCG
trên những hệ thống đó vẫn rất đáng chú ý.

Trên một số hệ thống, “Go 1: Uint64” nhanh hơn “PCG: Uint64”,
nhưng “Go 1: N(1000)” lại chậm hơn “PCG: N(1000)”.
Điều này xảy ra vì “Go 1: N(1000)” đang dùng thuật toán của `math/rand`
để rút một `int64` ngẫu nhiên xuống giá trị trong đoạn [0, 1000),
và thuật toán đó thực hiện hai phép chia số nguyên 64 bit.
Ngược lại, “PCG: N(1000)” và “ChaCha8: N(1000)” dùng [thuật toán `math/rand/v2` nhanh hơn](/blog/randv2#problem.rand),
gần như luôn tránh được các phép chia đó.
Việc loại bỏ phép chia 64 bit chi phối phần thay đổi thuật toán
trong thực thi 32 bit và trên Ampere.

Nhìn chung, ChaCha8Rand chậm hơn bộ sinh Go 1,
nhưng nó không bao giờ chậm quá gấp đôi, và trên các máy chủ điển hình,
chênh lệch không bao giờ vượt quá 3ns.
Rất ít chương trình bị nghẽn cổ chai bởi khác biệt này,
và nhiều chương trình sẽ được hưởng lợi từ bảo mật tốt hơn.

## Kết luận

Go 1.22 làm cho chương trình của bạn an toàn hơn mà không cần thay đổi mã.
Chúng tôi làm được điều đó bằng cách xác định sai lầm phổ biến là vô tình dùng `math/rand`
thay vì `crypto/rand`, rồi tăng cường `math/rand`.
Đây là một bước nhỏ trong hành trình không ngừng của Go nhằm giữ cho chương trình
an toàn theo mặc định.

Những kiểu sai lầm này không chỉ xuất hiện trong Go.
Ví dụ, package `keypair` của npm cố gắng sinh cặp khóa RSA
bằng Web Crypto API, nhưng nếu chúng không sẵn có, nó quay về dùng `Math.random` của JavaScript.
Đây không hề là một trường hợp cá biệt,
và độ an toàn của hệ thống không thể phụ thuộc vào việc lập trình viên không bao giờ mắc lỗi.
Thay vào đó, chúng tôi hy vọng rồi mọi ngôn ngữ lập trình
cuối cùng đều sẽ chuyển sang dùng bộ sinh giả ngẫu nhiên mạnh về mặt mật mã
ngay cả cho tính ngẫu nhiên “toán học”,
loại bỏ kiểu sai lầm này, hoặc ít nhất giảm mạnh bán kính ảnh hưởng của nó.
Hiện thực [ChaCha8Rand](https://c2sp.org/chacha8rand) của Go 1.22
chứng minh rằng cách tiếp cận này có tính cạnh tranh với các bộ sinh khác.
