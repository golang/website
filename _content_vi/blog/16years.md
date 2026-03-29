---
title: Sinh nhật 16 ngọt ngào của Go
date: 2025-11-14
by:
- Austin Clements, thay mặt nhóm Go
tags:
- community
- birthday
summary: Chúc mừng sinh nhật, Go!
---

Thứ Hai vừa rồi, ngày 10 tháng 11, chúng tôi đã kỷ niệm 16 năm kể từ [bản phát hành mã nguồn mở của Go](https://opensource.googleblog.com/2009/11/hey-ho-lets-go.html)\!

Chúng tôi đã phát hành [Go 1.24 vào tháng 2](/blog/go1.24) và [Go 1.25 vào tháng 8](/blog/go1.25),
theo đúng nhịp phát hành đều đặn và đáng tin cậy mà nay đã trở nên quen thuộc.
Tiếp tục sứ mệnh xây dựng nền tảng ngôn ngữ năng suất nhất cho việc phát triển các hệ thống production,
các bản phát hành này đã mang đến những API mới để xây dựng phần mềm vững chắc và đáng tin cậy,
những bước tiến đáng kể trong thành tích của Go về phát triển phần mềm an toàn,
và một số cải tiến rất đáng kể ở tầng bên dưới.
Trong khi đó, không ai có thể phớt lờ những chuyển động mang tính địa chấn trong ngành của chúng ta do AI tạo sinh tạo ra.
Nhóm Go đang mang tư duy cẩn trọng và không thỏa hiệp của mình vào những vấn đề và cơ hội của không gian năng động này,
nhằm đưa cách tiếp cận sẵn sàng cho production của Go vào việc xây dựng tích hợp AI, sản phẩm AI, agent và hạ tầng AI vững chắc.

# Các cải tiến cốt lõi của ngôn ngữ và thư viện

Được phát hành lần đầu như một thử nghiệm trong Go 1.24 và trở thành tính năng chính thức trong Go 1.25,
package mới [`testing/synctest`](https://pkg.go.dev/testing/synctest)
giúp việc viết kiểm thử cho [mã đồng thời, bất đồng bộ](/blog/testing-time)
trở nên đơn giản hơn rất nhiều.
Loại mã này đặc biệt phổ biến trong các dịch vụ mạng,
và từ trước đến nay vốn rất khó để kiểm thử cho tốt.
`synctest` hoạt động bằng cách ảo hóa chính thời gian.
Nó biến những bài kiểm thử vốn chậm, flaky, hoặc tệ hơn là cả hai,
thành những bài kiểm thử đáng tin cậy và gần như tức thời,
thường chỉ cần thêm vài dòng mã.
Đây cũng là một ví dụ điển hình cho cách tiếp cận tích hợp của Go với phát triển phần mềm:
đằng sau một API gần như tối giản,
`synctest` che giấu sự tích hợp sâu với runtime Go và các phần khác của thư viện chuẩn.

Đây không phải là cải tiến duy nhất mà package `testing` nhận được trong năm qua.
API mới [`testing.B.Loop`](https://pkg.go.dev/testing#B.Loop) vừa dễ dùng hơn API `testing.B.N` ban đầu,
vừa xử lý nhiều [cạm bẫy](/blog/testing-b-loop) truyền thống và thường vô hình khi viết benchmark trong Go.
Package `testing` cũng có các API mới giúp [dọn dẹp](https://pkg.go.dev/testing#T.Context)
trong những kiểm thử dùng [`Context`](https://pkg.go.dev/context#Context),
và [giúp dễ dàng](https://pkg.go.dev/testing#T.Output) ghi vào log của bài kiểm thử.

Go và containerization đã trưởng thành cùng nhau và phối hợp rất tốt với nhau.
Go 1.25 phát hành [container-aware scheduling](/blog/container-aware-gomaxprocs),
khiến cặp đôi này còn mạnh hơn nữa.
Không cần nhà phát triển phải động tay,
tính năng này tự động điều chỉnh mức độ song song của các tải công việc Go chạy trong container,
tránh tình trạng CPU throttling có thể ảnh hưởng đến tail latency và cải thiện mức sẵn sàng cho production của Go ngay khi cài đặt mặc định.

Tính năng mới [flight recorder](/blog/flight-recorder) trong Go 1.25
dựa trên execution tracer vốn đã rất mạnh của chúng tôi,
cho phép quan sát sâu vào hành vi động của các hệ thống production.
Trong khi execution tracer thường thu thập *quá nhiều* thông tin để có thể áp dụng thực tế
trong các dịch vụ production chạy lâu dài,
flight recorder giống như một cỗ máy thời gian nhỏ,
cho phép một dịch vụ chụp lại các sự kiện gần đây với độ chi tiết cao *sau khi* có điều gì đó xảy ra.

## Phát triển phần mềm an toàn

Go tiếp tục củng cố cam kết của mình đối với phát triển phần mềm an toàn,
đạt được những bước tiến đáng kể trong các package mật mã gốc
và tiếp tục phát triển thư viện chuẩn theo hướng an toàn hơn.

Go cung cấp trọn bộ các package mật mã gốc trong thư viện chuẩn,
và trong năm qua chúng đã đạt hai cột mốc lớn.
Một cuộc kiểm toán bảo mật do công ty an ninh độc lập [Trail of Bits](https://www.trailofbits.com/) thực hiện
đã cho [kết quả xuất sắc](/blog/tob-crypto-audit),
chỉ có một phát hiện mức độ nghiêm trọng thấp.
Hơn nữa, thông qua nỗ lực hợp tác giữa Nhóm Bảo mật Go và [Geomys](https://geomys.org/),
các package này đã đạt chứng nhận CAVP,
mở đường cho [chứng nhận FIPS 140-3 đầy đủ](/blog/fips140).
Đây là bước phát triển rất quan trọng đối với người dùng Go trong một số môi trường chịu quản lý chặt chẽ.
Việc tuân thủ FIPS 140, vốn trước đây gây nhiều vướng mắc do phải dựa vào những giải pháp không được hỗ trợ,
giờ sẽ được tích hợp liền mạch, giải quyết đồng thời các mối quan tâm liên quan đến an toàn,
trải nghiệm nhà phát triển, chức năng, tốc độ phát hành và tuân thủ.

Thư viện chuẩn của Go tiếp tục phát triển theo hướng *an toàn theo mặc định*
và *an toàn ngay từ thiết kế*.
Ví dụ, API [`os.Root`](https://pkg.go.dev/os#Root), được thêm vào Go 1.24,
cho phép [truy cập hệ thống tệp chống traversal](/blog/osroot),
qua đó chống lại hiệu quả một lớp lỗ hổng trong đó kẻ tấn công có thể lừa chương trình truy cập
những tệp vốn không được phép truy cập.
Những lỗ hổng như vậy từ lâu vốn rất khó giải quyết
nếu không có sự hỗ trợ từ nền tảng và hệ điều hành.
API [`os.Root`](https://pkg.go.dev/os#Root) mới mang lại một giải pháp
đơn giản, nhất quán và có tính di động.

## Những cải tiến ở tầng bên dưới

Bên cạnh các thay đổi dễ thấy đối với người dùng, Go cũng đạt được nhiều cải tiến đáng kể ở tầng bên dưới trong năm qua.

Đối với Go 1.24, chúng tôi đã [thiết kế lại hoàn toàn cách cài đặt `map`](/blog/swisstable),
dựa trên những ý tưởng mới và tốt nhất trong thiết kế bảng băm.
Thay đổi này hoàn toàn trong suốt đối với người dùng,
nhưng mang lại cải thiện đáng kể cho hiệu năng của `map`,
giảm tail latency của các thao tác `map`,
và trong một số trường hợp còn tiết kiệm bộ nhớ đáng kể.

Go 1.25 bao gồm một bước tiến lớn, mang tính thử nghiệm,
trong bộ gom rác của Go mang tên [Green Tea](/blog/greenteagc).
Green Tea giảm overhead của garbage collection trong nhiều ứng dụng ít nhất 10%
và đôi khi lên đến 40%.
Nó sử dụng một thuật toán mới, được thiết kế cho khả năng và ràng buộc của phần cứng ngày nay,
đồng thời mở ra một không gian thiết kế mới mà chúng tôi đang rất hào hứng khám phá.
Ví dụ, trong bản phát hành Go 1.26 sắp tới,
Green Tea sẽ đạt thêm 10% giảm overhead của garbage collector
trên phần cứng hỗ trợ các lệnh vector AVX-512,
điều gần như không thể tận dụng được với thuật toán cũ.
Green Tea sẽ được bật mặc định trong Go 1.26;
người dùng chỉ cần nâng cấp phiên bản Go là đã có thể hưởng lợi.

# Tiếp tục phát triển toàn bộ ngăn xếp phần mềm

Go còn hơn nhiều so với ngôn ngữ và thư viện chuẩn.
Nó là một nền tảng phát triển phần mềm, và trong năm qua,
chúng tôi cũng đã thực hiện bốn bản phát hành định kỳ của [gopls language server](/gopls),
đồng thời thiết lập quan hệ hợp tác để hỗ trợ các framework mới nổi dành cho ứng dụng agentic.

Gopls cung cấp hỗ trợ Go cho VS Code và các trình soạn thảo, IDE dựa trên LSP khác.
Mỗi bản phát hành đều mang đến hàng loạt tính năng và cải tiến cho trải nghiệm
đọc và viết mã Go (hãy xem ghi chú phát hành của [v0.17.0](/gopls/release/v0.17.0),
[v0.18.0](/gopls/release/v0.18.0), [v0.19.0](/gopls/release/v0.19.0),
và [v0.20.0](/gopls/release/v0.20.0) để biết đầy đủ chi tiết,
hoặc xem [tài liệu tính năng gopls mới](/gopls/features)\!).
Một số điểm nổi bật gồm nhiều analyzer mới và được nâng cấp
để giúp nhà phát triển viết mã Go đúng phong cách hơn và vững chắc hơn;
hỗ trợ refactor để tách biến, inline biến, và JSON struct tags;
và một [máy chủ tích hợp dạng thử nghiệm](/gopls/features/mcp) cho
Model Context Protocol (MCP), cung cấp một phần chức năng của gopls cho trợ lý AI
dưới dạng các công cụ MCP.

Từ gopls v0.18.0, chúng tôi bắt đầu khám phá *automatic code modernizers*.
Khi Go phát triển, mỗi bản phát hành lại mang đến năng lực mới và thành ngữ mới;
những cách mới và tốt hơn để làm những điều mà lập trình viên Go từ trước đến nay vẫn làm theo cách khác.
Go vẫn giữ [cam kết tương thích](/doc/go1compat) của mình:
cách cũ sẽ tiếp tục hoạt động mãi mãi.
Tuy nhiên, điều này cũng tạo ra sự phân nhánh giữa thành ngữ cũ và thành ngữ mới.
Modernizer là các công cụ phân tích tĩnh có thể nhận diện những thành ngữ cũ
và gợi ý các cách thay thế nhanh hơn, dễ đọc hơn, an toàn hơn và *hiện đại* hơn,
với độ tin cậy gần như chỉ cần một nút bấm.
Điều mà `gofmt` từng làm cho [tính nhất quán về phong cách](/blog/gofmt),
chúng tôi hy vọng modernizer có thể làm cho tính nhất quán về thành ngữ.
Chúng tôi đã tích hợp modernizer như các gợi ý trong IDE,
nơi chúng không chỉ giúp nhà phát triển duy trì tiêu chuẩn mã nhất quán hơn,
mà còn giúp họ khám phá các tính năng mới và bắt kịp trình độ hiện tại.
Chúng tôi cũng tin rằng modernizer có thể giúp các trợ lý lập trình AI
bắt kịp trình độ hiện tại và chống lại xu hướng củng cố kiến thức lỗi thời
về ngôn ngữ Go, API và thành ngữ của nó.
Bản phát hành Go 1.26 sắp tới sẽ bao gồm một cuộc đại tu toàn diện cho lệnh `go fix`,
vốn đã ngủ quên từ lâu, để nó có thể áp dụng toàn bộ bộ modernizer hàng loạt,
như một sự quay về với [gốc rễ từ trước Go 1.0](/blog/introducing-gofix).

Cuối tháng 9, với sự hợp tác của [Anthropic](https://www.anthropic.com/) và cộng đồng Go,
chúng tôi đã phát hành [v1.0.0](https://github.com/modelcontextprotocol/go-sdk/releases/tag/v1.0.0)
của [Go SDK chính thức](https://github.com/modelcontextprotocol/go-sdk) cho
[Model Context Protocol (MCP)](https://modelcontextprotocol.io/).
SDK này hỗ trợ cả client MCP lẫn server MCP,
và cũng là nền tảng cho chức năng MCP mới trong gopls.
Việc đóng góp công việc này theo mô hình mã nguồn mở giúp trao quyền cho các lĩnh vực khác
trong hệ sinh thái agentic mã nguồn mở đang phát triển dựa trên Go,
chẳng hạn như [Agent Development Kit (ADK) for Go](https://github.com/google/adk-go)
mới được phát hành bởi [Google](https://www.google.com/).
ADK Go xây dựng trên Go MCP SDK để cung cấp một framework đúng phong cách Go
cho việc xây dựng các ứng dụng và hệ thống đa agent theo mô-đun.
Go MCP SDK và ADK Go cho thấy thế mạnh riêng của Go
về đồng thời, hiệu năng và độ tin cậy
đang tạo ra sự khác biệt cho Go trong phát triển AI production,
và chúng tôi kỳ vọng sẽ có nhiều tải công việc AI hơn được viết bằng Go trong những năm tới.

# Nhìn về phía trước

Go còn một năm rất đáng mong đợi phía trước.

Chúng tôi đang thúc đẩy năng suất của nhà phát triển thông qua lệnh `go fix` hoàn toàn mới,
hỗ trợ sâu hơn cho các trợ lý lập trình AI,
và các cải tiến liên tục cho gopls cũng như VS Code Go.
Việc đưa Green Tea garbage collector vào trạng thái phát hành chính thức,
hỗ trợ gốc cho các tính năng phần cứng Single Instruction Multiple Data (SIMD),
và hỗ trợ trong runtime cùng thư viện chuẩn để viết mã mở rộng tốt hơn nữa
trên phần cứng đa lõi quy mô lớn sẽ tiếp tục giúp Go đồng bộ với phần cứng hiện đại
và cải thiện hiệu quả production.
Chúng tôi đang tập trung vào các thư viện và công cụ chẩn đoán trong “ngăn xếp production” của Go,
bao gồm một [bản nâng cấp lớn cho `encoding/json`](/issue/71497),
được thúc đẩy bởi Joe Tsai và nhiều người trong cộng đồng Go;
[profiling goroutine bị rò rỉ](/design/74609-goroutine-leak-detection-gc),
do nhóm Programming Systems của [Uber](https://www.uber.com/us/en/about/) đóng góp;
và nhiều cải tiến khác cho `net/http`, `unicode`, và những package nền tảng khác.
Chúng tôi đang làm việc để tạo ra những con đường rõ ràng, dễ tiếp cận cho việc xây dựng với Go và AI,
phát triển nền tảng ngôn ngữ một cách cẩn trọng trước nhu cầu đang thay đổi của nhà phát triển ngày nay,
đồng thời xây dựng công cụ và năng lực hỗ trợ cho cả nhà phát triển con người lẫn các trợ lý và hệ thống AI.

Nhân dịp kỷ niệm 16 năm bản phát hành mã nguồn mở của Go, chúng tôi cũng đang nhìn vào
tương lai của chính dự án mã nguồn mở Go.
Từ [khởi đầu khiêm tốn](https://www.youtube.com/watch?v=wwoWei-GAPo),
Go đã hình thành nên một cộng đồng đóng góp đầy sức sống.
Để tiếp tục đáp ứng tốt nhất nhu cầu của tập người dùng ngày càng mở rộng,
đặc biệt trong giai đoạn ngành phần mềm đang có nhiều biến động,
chúng tôi đang tìm cách mở rộng tốt hơn quy trình phát triển của Go
mà không đánh mất những nguyên tắc nền tảng của Go,
đồng thời đưa cộng đồng đóng góp tuyệt vời của chúng ta tham gia sâu hơn.

Go sẽ không thể có được vị thế ngày hôm nay nếu thiếu cộng đồng người dùng và người đóng góp tuyệt vời của mình.
Chúng tôi chúc tất cả các bạn những điều tốt đẹp nhất trong năm tới\!
