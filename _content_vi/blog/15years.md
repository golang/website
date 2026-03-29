---
title: Go tròn 15 tuổi
date: 2024-11-11
by:
- Austin Clements, thay mặt nhóm Go
tags:
- community
- birthday
summary: Chúc mừng sinh nhật lần thứ 15 của Go!
---

<div style="float:right; margin: 0 0 1em 1em; width: 245px">
<img src="/doc/gopher/fifteen.gif" height="245" width="245"><br/>
<i>Cảm ơn Renee French đã vẽ và tạo hoạt hình cho chú gopher chơi trò “15 puzzle”.</i>
</div>

Chúc mừng sinh nhật, Go\!

Vào Chủ nhật, chúng tôi đã kỷ niệm 15 năm kể từ [bản phát hành mã nguồn mở của Go](https://opensource.googleblog.com/2009/11/hey-ho-lets-go.html)\!

Rất nhiều điều đã thay đổi kể từ [mốc Go tròn 10 tuổi](/blog/10years),
cả trong Go lẫn trong thế giới xung quanh. Nhưng ở những khía cạnh khác, nhiều điều vẫn giữ nguyên:
Go vẫn cam kết với tính ổn định, an toàn, và hỗ trợ kỹ nghệ phần mềm cũng như vận hành production ở quy mô lớn.

Và Go vẫn đang phát triển mạnh mẽ\! Số người dùng Go đã tăng hơn gấp ba trong
năm năm qua, đưa nó trở thành một trong những ngôn ngữ tăng trưởng nhanh nhất.
Từ khởi đầu chỉ mới mười lăm năm trước, Go đã trở thành một ngôn ngữ top 10
và là ngôn ngữ của đám mây hiện đại.

Với các bản phát hành [Go 1.22 vào tháng 2](/blog/go1.22) và [Go 1.23 vào tháng 8](/blog/go1.23),
đây có thể xem là năm của vòng lặp `for`. Go 1.22 khiến các biến được giới thiệu bởi
vòng lặp `for` [có phạm vi theo từng lần lặp](/blog/loopvar-preview), thay vì theo toàn bộ vòng lặp,
qua đó xử lý một “bẫy” lâu năm của ngôn ngữ. Hơn mười năm trước, khi chuẩn bị cho
việc phát hành Go 1, nhóm Go đã phải đưa ra quyết định về nhiều chi tiết của ngôn ngữ;
trong đó có việc liệu vòng lặp `for` có nên tạo một biến lặp mới ở mỗi lần lặp hay không.
Thú vị là cuộc thảo luận lúc đó khá ngắn và gần như không mang màu sắc quan điểm.
Rob Pike đã chốt nó đúng theo phong cách của Rob Pike bằng một từ duy nhất:
“stet” (cứ để nguyên). Và mọi chuyện đã diễn ra như vậy. Dù lúc đó có vẻ không đáng kể,
nhiều năm kinh nghiệm trong môi trường production đã cho thấy hệ quả của quyết định này.
Nhưng trong quãng thời gian đó, chúng tôi cũng xây dựng được những công cụ mạnh để hiểu
tác động của thay đổi đối với Go, nổi bật là khả năng phân tích và kiểm thử trên toàn hệ sinh thái
trong toàn bộ codebase của Google, đồng thời thiết lập được quy trình làm việc với cộng đồng
và tiếp nhận phản hồi. Sau quá trình kiểm thử, phân tích và thảo luận rộng rãi với cộng đồng,
chúng tôi triển khai thay đổi này cùng với một
[công cụ hash bisection](https://go.googlesource.com/proposal/+/master/design/60078-loopvar.md#transition-support-tooling)
để hỗ trợ nhà phát triển xác định chính xác đoạn mã bị ảnh hưởng ở quy mô lớn.

Thay đổi đối với vòng lặp `for` là một phần của quỹ đạo thay đổi có kiểm soát trong suốt năm năm.
Điều đó sẽ không thể xảy ra nếu không có [tính tương thích xuôi của ngôn ngữ](/blog/toolchain)
được giới thiệu trong Go 1.21. Và bản thân điều này lại được xây dựng trên nền tảng do Go modules đặt ra,
vốn được giới thiệu trong Go 1.14 cách đây bốn năm rưỡi.

Go 1.23 tiếp tục phát triển trên thay đổi đó để giới thiệu iterator và
[vòng lặp for-range do người dùng định nghĩa](/blog/range-functions).
Kết hợp với generics, được giới thiệu trong Go 1.18 chỉ mới hai năm rưỡi trước\!,
điều này tạo ra một nền tảng mạnh mẽ và thuận tiện cho các collection tùy biến
và nhiều mẫu lập trình khác.

Các bản phát hành này cũng mang đến nhiều cải tiến về mức sẵn sàng cho production,
bao gồm [những nâng cấp rất được mong đợi cho bộ định tuyến HTTP trong thư viện chuẩn](/blog/routing-enhancements),
[một cuộc đại tu toàn diện của execution traces](/blog/execution-traces-2024),
và [nguồn ngẫu nhiên mạnh hơn](/blog/chacha8rand) cho mọi ứng dụng Go.
Ngoài ra, việc giới thiệu [package thư viện chuẩn v2 đầu tiên](/blog/randv2)
cũng đặt ra khuôn mẫu cho quá trình tiến hóa và hiện đại hóa thư viện trong tương lai.

Trong năm qua, chúng tôi cũng đã thận trọng triển khai [telemetry theo cơ chế tự nguyện tham gia](/blog/gotelemetry)
cho các công cụ của Go. Hệ thống này sẽ cung cấp cho các nhà phát triển Go dữ liệu để ra quyết định tốt hơn,
trong khi vẫn hoàn toàn [mở](https://telemetry.go.dev/) và ẩn danh.
Telemetry của Go lần đầu xuất hiện trong
[gopls](https://github.com/golang/tools/blob/master/gopls/README.md),
máy chủ ngôn ngữ Go, nơi nó đã dẫn tới một [loạt cải tiến](https://github.com/golang/go/issues?q=is%3Aissue+label%3Agopls%2Ftelemetry-wins).
Nỗ lực này mở đường để việc lập trình bằng Go trở thành một trải nghiệm còn tốt hơn nữa cho mọi người.

Nhìn về phía trước, chúng tôi đang phát triển Go để tận dụng tốt hơn khả năng của phần cứng hiện tại
và tương lai. Phần cứng đã thay đổi rất nhiều trong 15 năm qua. Để bảo đảm Go tiếp tục hỗ trợ
các tải công việc production hiệu năng cao, quy mô lớn trong *15 năm tới*,
chúng tôi cần thích nghi với các hệ thống đa lõi lớn, các tập lệnh tiên tiến,
và tầm quan trọng ngày càng tăng của locality trong các phân cấp bộ nhớ ngày càng không đồng nhất.
Một số cải tiến trong số này sẽ hoàn toàn trong suốt với người dùng.
Go 1.24 sẽ có một triển khai `map` hoàn toàn mới ở tầng bên dưới, hiệu quả hơn
trên các CPU hiện đại. Và chúng tôi đang thử nghiệm các thuật toán garbage collection mới
được thiết kế xoay quanh khả năng và ràng buộc của phần cứng hiện đại.
Một số cải tiến khác sẽ xuất hiện dưới dạng API và công cụ mới để các nhà phát triển Go
có thể tận dụng phần cứng hiện đại tốt hơn. Chúng tôi đang nghiên cứu cách hỗ trợ
các lệnh vector và ma trận mới nhất của phần cứng, cũng như nhiều cách để ứng dụng
có thể xây dựng locality cho CPU và bộ nhớ. Một nguyên tắc cốt lõi định hướng nỗ lực của chúng tôi
là *tối ưu hóa có tính kết hợp*: tác động của một tối ưu hóa lên codebase nên càng cục bộ càng tốt,
để sự dễ dàng trong phát triển ở phần còn lại của codebase không bị ảnh hưởng.

Chúng tôi tiếp tục bảo đảm thư viện chuẩn của Go an toàn theo mặc định
và an toàn ngay từ thiết kế. Điều này bao gồm các nỗ lực đang diễn ra nhằm đưa vào
hỗ trợ gốc, tích hợp sẵn cho mật mã học đạt chứng nhận FIPS, để các ứng dụng cần FIPS
chỉ còn cách một lần bật cờ cấu hình. Hơn nữa, chúng tôi đang phát triển các package trong thư viện chuẩn
của Go ở những nơi có thể, và theo gương của `math/rand/v2`, đang cân nhắc những nơi mà
API mới có thể cải thiện đáng kể sự dễ dàng khi viết mã Go an toàn và bảo mật.

Chúng tôi đang làm cho Go tốt hơn cho AI và làm cho AI tốt hơn cho Go bằng cách
nâng cao năng lực của Go trong hạ tầng AI, ứng dụng AI và trợ giúp cho nhà phát triển.
Go là một ngôn ngữ tuyệt vời để xây dựng các hệ thống production,
và chúng tôi cũng muốn nó trở thành một ngôn ngữ tuyệt vời để [xây dựng các hệ thống AI production](/blog/llmpowered).
Sự đáng tin cậy của Go như một ngôn ngữ dành cho hạ tầng đám mây đã khiến nó trở thành lựa chọn tự nhiên
cho hạ tầng [LLM](https://ollama.com/), [vector database](https://weaviate.io/),
[AI cục bộ](https://localai.io/) và [các hệ thống tương tự](https://zilliz.com/what-is-milvus).
Đối với các ứng dụng AI, chúng tôi sẽ tiếp tục xây dựng hỗ trợ hạng nhất cho Go
trong các SDK AI phổ biến, bao gồm
[LangChainGo](https://pkg.go.dev/github.com/tmc/langchaingo) và
[Genkit](https://developers.googleblog.com/en/introducing-genkit-for-go-build-scalable-ai-powered-apps-in-go/).
Và ngay từ buổi đầu, Go đã hướng tới việc cải thiện toàn bộ quy trình kỹ nghệ phần mềm,
vì vậy một cách tự nhiên, chúng tôi cũng đang tìm cách đưa những công cụ và kỹ thuật AI mới nhất
vào việc giảm bớt công việc nặng nhọc cho nhà phát triển, để dành nhiều thời gian hơn cho phần thú vị nhất
như… thực sự lập trình\!

## Xin cảm ơn

Tất cả những điều này chỉ có thể xảy ra nhờ những người đóng góp tuyệt vời
và cộng đồng đầy sức sống của Go. Mười lăm năm trước, chúng tôi chỉ có thể mơ về
thành công mà Go đã đạt được và cộng đồng đã phát triển quanh Go.
Xin cảm ơn tất cả mọi người đã góp một phần, dù lớn hay nhỏ.
Chúng tôi chúc tất cả các bạn những điều tốt đẹp nhất trong năm tới.
