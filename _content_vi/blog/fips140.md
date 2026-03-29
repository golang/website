---
title: Mô-đun mật mã Go FIPS 140-3
date: 2025-07-15
by:
- Filippo Valsorda (Geomys)
- Daniel McCarney (Geomys)
- Roland Shoemaker (Google)
summary: Go giờ đã có chế độ tuân thủ FIPS 140-3 gốc, tích hợp sẵn.
---

FIPS 140 là một tiêu chuẩn dành cho các hiện thực mật mã, và mặc dù nó không nhất thiết
làm tăng mức độ an toàn, việc tuân thủ FIPS 140 lại là yêu cầu trong một số
môi trường được quản lý đang ngày càng dùng Go nhiều hơn. Cho đến nay, tuân thủ FIPS 140
là một nguồn ma sát lớn với người dùng Go, buộc họ phải dùng những giải pháp không được hỗ trợ,
vướng các vấn đề về an toàn, trải nghiệm phát triển, chức năng, tốc độ phát hành
và tuân thủ.

Go đang giải quyết nhu cầu ngày càng tăng này bằng hỗ trợ FIPS 140 gốc được tích hợp trực tiếp
vào thư viện chuẩn và lệnh `go`, giúp Go trở thành cách dễ nhất, an toàn nhất
để tuân thủ FIPS 140. Mô-đun mật mã Go đã được xác thực FIPS 140-3
giờ là nền tảng của các thư viện crypto tích hợp sẵn của Go, bắt đầu từ
Go Cryptographic Module v1.0.0 có trong Go 1.24, phát hành vào tháng 2 vừa qua.

Mô-đun v1.0.0 đã được cấp [chứng chỉ Cryptographic Algorithm Validation Program
(CAVP) A6650][], đã được nộp lên Cryptographic Module
Validation Program (CMVP), và đã có mặt trong [Modules In Process List][] vào tháng 5.
Các mô-đun nằm trong danh sách MIP đang chờ NIST xét duyệt và đã có thể được triển khai trong
một số môi trường được quản lý nhất định.

[Geomys][] dẫn dắt việc hiện thực mô-đun trong hợp tác với Go Security
Team, và đang theo đuổi một chứng nhận FIPS 140-3 có phạm vi áp dụng rộng cho lợi ích
của cộng đồng Go. Google và các bên liên quan khác trong ngành có quan hệ hợp đồng
với Geomys để đưa những Operating Environment cụ thể vào chứng chỉ.

Thông tin chi tiết hơn về mô-đun có trong
[tài liệu](/doc/security/fips140).

Hiện nay một số người dùng Go vẫn dựa vào GOEXPERIMENT [Go+BoringCrypto][] hoặc một
fork của nó như một phần trong chiến lược tuân thủ FIPS 140 của họ. Khác với Go Cryptographic Module
FIPS 140-3, Go+BoringCrypto chưa bao giờ được hỗ trợ chính thức
và gặp nhiều vấn đề lớn về trải nghiệm phát triển, do nó được tạo ra
chỉ để phục vụ nhu cầu nội bộ của Google. Nó sẽ bị loại bỏ trong một bản phát hành tương lai khi Google
chuyển sang mô-đun gốc.

## Trải nghiệm phát triển gốc

Mô-đun này tích hợp hoàn toàn trong suốt vào các ứng dụng Go. Trên thực tế,
mọi chương trình Go được xây bằng Go 1.24 đã dùng nó cho mọi thuật toán
được FIPS 140-3 phê duyệt! Mô-đun này thực chất chỉ là một tên khác của các gói
`crypto/internal/fips140/...` trong thư viện chuẩn, vốn cung cấp
phần hiện thực cho các phép toán được phơi bày bởi những gói như `crypto/ecdsa` và
`crypto/rand`.

Những gói này không dùng cgo, nghĩa là chúng cross-compile như mọi chương trình Go khác,
không phải trả overhead hiệu năng FFI, và không gặp phải
[các vấn đề an toàn về quản lý bộ nhớ][] như Go+BoringCrypto và các fork của nó.

Khi khởi động một tệp nhị phân Go, mô-đun có thể được đưa vào chế độ FIPS 140-3 bằng
tùy chọn GODEBUG `fips140=on`, có thể được đặt như biến môi trường hoặc qua
tệp `go.mod`. Nếu chế độ FIPS 140-3 được bật, mô-đun sẽ dùng NIST DRBG cho tính ngẫu nhiên,
`crypto/tls` sẽ tự động chỉ đàm phán các phiên bản và thuật toán TLS được FIPS 140-3 phê duyệt,
và nó sẽ thực hiện các self-test bắt buộc trong lúc khởi tạo và khi sinh khóa. Chỉ vậy thôi;
không có khác biệt hành vi nào khác.

Cũng có một chế độ nghiêm ngặt thử nghiệm là `fips140=only`, khiến mọi
thuật toán không được phê duyệt trả về lỗi hoặc panic. Chúng tôi hiểu rằng điều này có thể
quá cứng nhắc cho hầu hết triển khai và đang [tìm kiếm phản hồi](/issue/74630)
về hình hài của một khung ép buộc chính sách phù hợp.

Cuối cùng, ứng dụng có thể dùng [`GOFIPS140` environment
variable](/doc/security/fips140#the-gofips140-environment-variable)
để build dựa trên các phiên bản cũ hơn đã được xác thực của các gói `crypto/internal/fips140/...`.
`GOFIPS140` hoạt động như `GOOS` và `GOARCH`, và nếu đặt là
`GOFIPS140=v1.0.0` thì chương trình sẽ được build dựa trên snapshot v1.0.0 của
các gói đúng như khi chúng được nộp cho CMVP để xác thực. Snapshot này được phát hành cùng
phần còn lại của thư viện chuẩn Go dưới dạng `lib/fips140/v1.0.0.zip`.

Khi dùng `GOFIPS140`, giá trị mặc định của GODEBUG `fips140` sẽ là `on`, nên gộp lại,
tất cả những gì cần để build dựa trên mô-đun FIPS 140-3 và chạy ở chế độ FIPS 140-3
là `GOFIPS140=v1.0.0 go build`. Chỉ vậy thôi.

Nếu một toolchain được build với `GOFIPS140` đã đặt, mọi bản build mà nó sinh ra sẽ
mặc định dùng giá trị đó.

Phiên bản `GOFIPS140` dùng để build một tệp nhị phân có thể được kiểm tra bằng
`go version -m`.

Các phiên bản Go tương lai sẽ tiếp tục phát hành và hoạt động với v1.0.0 của Go
Cryptographic Module cho tới khi phiên bản tiếp theo được Geomys chứng nhận đầy đủ, nhưng
một số tính năng mật mã mới có thể không sẵn có khi build dựa trên các mô-đun cũ.
Bắt đầu từ Go 1.24.3, bạn có thể dùng `GOFIPS140=inprocess` để
chọn động mô-đun mới nhất mà chứng nhận của Geomys đã đạt tới trạng thái In Process.
Geomys dự định xác thực các phiên bản mô-đun mới ít nhất mỗi năm một lần để tránh
việc các bản build FIPS 140 bị tụt lại quá xa và cả mỗi khi có một lỗ hổng
trong mô-đun mà không thể giảm thiểu ở mã thư viện chuẩn phía gọi tới.

## Bảo mật không thỏa hiệp

Ưu tiên hàng đầu của chúng tôi khi phát triển mô-đun là đạt hoặc vượt
mức độ an toàn của các gói mật mã hiện có trong thư viện chuẩn Go. Có thể điều này nghe bất ngờ,
nhưng đôi khi cách dễ nhất để đạt và chứng minh việc tuân thủ các yêu cầu an ninh của FIPS 140
là không vượt quá chúng. Chúng tôi từ chối chấp nhận điều đó.

Ví dụ, `crypto/ecdsa` [luôn sinh chữ ký hedged][]. Chữ ký hedged
sinh nonce bằng cách kết hợp khóa riêng, thông điệp và các byte ngẫu nhiên.
Giống [ECDSA xác định][RFC 6979], chúng bảo vệ khỏi việc bộ sinh số ngẫu nhiên gặp lỗi,
điều vốn có thể làm lộ khóa riêng(!). Khác với ECDSA xác định, chúng cũng chống lại [vấn đề API][]
và [tấn công lỗi][], đồng thời không làm lộ việc hai thông điệp có bằng nhau hay không.
FIPS 186-5 đã bổ sung hỗ trợ cho ECDSA xác định theo [RFC 6979][], nhưng không cho hedged ECDSA.

Thay vì hạ cấp xuống chữ ký ECDSA ngẫu nhiên thông thường hoặc xác định
trong chế độ FIPS 140-3 (hoặc tệ hơn, giữa các chế độ), chúng tôi đã [thay đổi
thuật toán hedging][] và xâu chuỗi lập luận qua nửa tá tài liệu để [chứng minh thuật toán mới
là một phép tổ hợp tuân thủ của DRBG và ECDSA truyền thống][]. Trong lúc đó,
chúng tôi cũng [bổ sung hỗ trợ chọn dùng cho chữ ký xác định][].

Một ví dụ khác là sinh số ngẫu nhiên. FIPS 140-3 có các quy tắc chặt chẽ về cách
sinh ngẫu nhiên mật mã, về cơ bản ép buộc dùng một [CSPRNG][] chạy trong userspace.
Ngược lại, chúng tôi tin rằng kernel phù hợp hơn để sinh ra các byte ngẫu nhiên an toàn,
vì nó ở vị trí tốt nhất để thu thập entropy từ hệ thống, và để phát hiện khi tiến trình
hoặc thậm chí máy ảo bị nhân bản (điều có thể dẫn tới việc tái sử dụng các byte lẽ ra là ngẫu nhiên).
Vì thế, [crypto/rand][] chuyển mọi thao tác đọc xuống kernel.

Để dung hòa điều này, trong chế độ FIPS 140-3 chúng tôi duy trì một NIST DRBG trong userspace
tuân thủ tiêu chuẩn, dựa trên AES-256-CTR, rồi tiêm thêm vào đó 128 bit lấy từ kernel ở mỗi lần đọc.
Entropy bổ sung này được xem là dữ liệu bổ sung “không được ghi công”
theo mục đích FIPS 140-3, nhưng trên thực tế khiến nó mạnh ngang với việc đọc trực tiếp từ kernel,
dù chậm hơn.

Cuối cùng, toàn bộ Go Cryptographic Module v1.0.0 đều nằm trong phạm vi của [đợt
đánh giá bảo mật gần đây bởi Trail of Bits](/blog/tob-crypto-audit), và
không bị ảnh hưởng bởi phát hiện duy nhất không phải loại thông tin.

Kết hợp với các bảo đảm an toàn bộ nhớ do trình biên dịch và runtime Go cung cấp,
chúng tôi tin điều này hiện thực hóa mục tiêu biến Go thành một trong những giải pháp
dễ nhất và an toàn nhất cho việc tuân thủ FIPS 140.

## Hỗ trợ nền tảng rộng

Một mô-đun FIPS 140-3 chỉ tuân thủ nếu được vận hành trên một Operating Environment
đã được thử nghiệm hoặc “Vendor Affirmed”, về cơ bản là sự kết hợp của hệ điều hành
và nền tảng phần cứng. Để hỗ trợ nhiều trường hợp sử dụng Go nhất có thể, xác thực của Geomys
được thử nghiệm trên [một trong những tập Operating Environment toàn diện nhất][] trong ngành.

Phòng thí nghiệm của Geomys đã thử nghiệm nhiều biến thể Linux (Alpine Linux trên Podman, Amazon
Linux, Google Prodimage, Oracle Linux, Red Hat Enterprise Linux và SUSE Linux
Enterprise Server), macOS, Windows và FreeBSD trên nhiều nền tảng x86-64 (AMD và
Intel), ARMv8/9 (Ampere Altra, Apple M, AWS Graviton và Qualcomm Snapdragon),
ARMv7, MIPS, z/ Architecture và POWER, tổng cộng 23 môi trường đã thử nghiệm.

Một số môi trường trong số này được các bên liên quan chi trả, số khác do Geomys tài trợ vì lợi ích
của cộng đồng Go.

Ngoài ra, chứng nhận của Geomys còn liệt kê một tập rộng các nền tảng tổng quát dưới dạng
Vendor Affirmed Operating Environments:
* Linux 3.10+ trên x86-64 và ARMv7/8/9,
* macOS 11–15 trên bộ xử lý Apple M,
* FreeBSD 12–14 trên x86-64,
* Windows 10 và Windows Server 2016–2022 trên x86-64, và
* Windows 11 cùng Windows Server 2025 trên x86-64 và ARMv8/9.

## Bao phủ thuật toán toàn diện

Có thể gây ngạc nhiên, nhưng ngay cả khi dùng một thuật toán được FIPS 140-3 phê duyệt,
được hiện thực bởi một mô-đun FIPS 140-3 trên một Operating Environment được hỗ trợ,
vẫn chưa chắc là đủ để tuân thủ; thuật toán đó phải được chính thức bao phủ bởi
quy trình thử nghiệm như một phần của chứng nhận. Vì thế, để giúp việc xây dựng ứng dụng
tuân thủ FIPS 140 bằng Go trở nên dễ dàng nhất có thể, mọi thuật toán được FIPS 140-3 phê duyệt
trong thư viện chuẩn đều được hiện thực bởi Go Cryptographic Module và được thử nghiệm
trong quá trình xác thực, từ chữ ký số cho đến TLS key schedule.

Trao đổi khóa hậu lượng tử ML-KEM (FIPS 203), [được đưa vào Go 1.24][mlkem
relnote], cũng đã được xác thực, nghĩa là `crypto/tls` có thể thiết lập các kết nối
an toàn hậu lượng tử tuân thủ FIPS 140-3 với X25519MLKEM768.

Trong một số trường hợp, chúng tôi đã xác thực cùng một thuật toán dưới nhiều định danh NIST khác nhau,
để có thể sử dụng chúng một cách hoàn toàn tuân thủ cho nhiều mục đích khác nhau.
Ví dụ, [HKDF được thử nghiệm và xác thực dưới *bốn* tên][hkdf]:
SP 800-108 Feedback KDF, SP 800-56C two-step KDF, Implementation Guidance D.P
OneStepNoCounter KDF, và SP 800-133 Section 6.3 KDF.

Cuối cùng, chúng tôi cũng xác thực một số thuật toán nội bộ như CMAC Counter KDF,
để có thể phơi bày các chức năng tương lai như [XAES-256-GCM][].

Nhìn chung, mô-đun FIPS 140-3 gốc mang lại hồ sơ tuân thủ tốt hơn Go+BoringCrypto,
đồng thời cung cấp nhiều thuật toán hơn cho các ứng dụng bị giới hạn theo FIPS 140-3.

Chúng tôi mong đợi mô-đun mật mã Go gốc mới sẽ giúp các lập trình viên Go chạy các workload
tuân thủ FIPS 140 dễ dàng hơn và an toàn hơn.

[Geomys]: https://geomys.org
[Cryptographic Algorithm Validation Program (CAVP) certificate A6650]: https://csrc.nist.gov/projects/cryptographic-algorithm-validation-program/details?validation=39260
[Modules In Process List]: https://csrc.nist.gov/Projects/cryptographic-module-validation-program/modules-in-process/modules-in-process-list
[Go+BoringCrypto]: /doc/security/fips140#goboringcrypto
[memory management security issues]: /blog/tob-crypto-audit#cgo-memory-management
[GODEBUG option]: /doc/godebug
[always produced hedged signatures]: https://cs.opensource.google/go/go/+/refs/tags/go1.23.0:src/crypto/ecdsa/ecdsa.go;l=417
[API issues]: https://github.com/MystenLabs/ed25519-unsafe-libs
[fault attacks]: https://en.wikipedia.org/wiki/Differential_fault_analysis
[RFC 6979]: https://www.rfc-editor.org/rfc/rfc6979
[switched the hedging algorithm]: https://github.com/golang/go/commit/9776d028f4b99b9a935dae9f63f32871b77c49af
[prove the new one is a compliant composition of a DRBG and traditional ECDSA]: https://github.com/cfrg/draft-irtf-cfrg-det-sigs-with-noise/issues/6#issuecomment-2067819904
[added opt-in support for deterministic signatures]: /doc/go1.24#cryptoecdsapkgcryptoecdsa
[CSPRNG]: https://en.wikipedia.org/wiki/Cryptographically_secure_pseudorandom_number_generator
[crypto/rand]: https://pkg.go.dev/crypto/rand
[one of the most comprehensive sets of Operating Environments]: https://csrc.nist.gov/projects/cryptographic-algorithm-validation-program/details?product=19371&displayMode=Aggregated
[mlkem relnote]: /doc/go1.24#crypto-mlkem
[hkdf]: https://words.filippo.io/dispatches/fips-hkdf/
[XAES-256-GCM]: https://c2sp.org/XAES-256-GCM
