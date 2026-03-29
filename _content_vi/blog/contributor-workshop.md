---
title: Workshop đóng góp
date: 2017-08-09
by:
- Steve Francia
- Cassandra Salisbury
- Matt Broberg
- Dmitri Shuralyov
tags:
- community
summary: Workshop contributor của Go đã đào tạo contributor mới tại GopherCon.
template: true
---

## Tổng quan sự kiện

by [Steve](https://twitter.com/spf13)

Trong ngày cộng đồng tại GopherCon, đội Go đã tổ chức hai workshop nơi chúng tôi làm việc cùng mọi người để giúp họ thực hiện đóng góp đầu tiên cho dự án Go. Đây là lần đầu tiên dự án Go từng thử điều gì như thế này. Chúng tôi có khoảng 140 người tham gia và khoảng 35 người tình nguyện làm mentor. Mentor không chỉ nhận được cảm giác ấm áp vui vẻ khi giúp đỡ người khác, mà còn nhận được một chiếc mũ lưới Go Mentor cực kỳ sành điệu. Chúng tôi có contributor ở đủ mọi độ tuổi và mức kinh nghiệm từ Bắc và Nam Mỹ, châu Phi, châu Âu, châu Á và Úc. Đó thực sự là một nỗ lực toàn cầu của các Gopher cùng nhau hội tụ tại GopherCon.

Một trong những lý do chúng tôi tổ chức workshop là để nó đóng vai trò như một “forcing function”, buộc chúng tôi phải cải thiện trải nghiệm contributor. Để chuẩn bị cho workshop, chúng tôi đã viết lại hướng dẫn contributor, bao gồm thêm mục “troubleshooting” và xây dựng công cụ `go-contrib-init`, công cụ tự động hóa quá trình thiết lập môi trường phát triển để có thể đóng góp cho Go.

Đối với chính workshop, chúng tôi đã phát triển bài trình bày _“Contributing to Go,”_ và một bảng dashboard / scoreboard được trình chiếu trong sự kiện. Scoreboard được thiết kế để khuyến khích tất cả chúng tôi cùng làm việc hướng tới mục tiêu chung là thấy tổng điểm tập thể tăng lên. Người tham gia cộng thêm 1, 2 hoặc 3 điểm vào tổng điểm khi họ thực hiện các hành động như đăng ký tài khoản, tạo change list (còn gọi là CL, tương tự pull request), sửa đổi một CL, hoặc gửi một CL.

{{image "contributor-workshop/image17.png"}}

Brad Fitzpatrick, người năm nay ở nhà thay vì tới GopherCon, đã sẵn sàng và chờ để review mọi CL được gửi lên. Anh ấy review nhanh đến mức nhiều người tưởng rằng anh là bot tự động. Nội bộ đội chúng tôi giờ gọi anh ấy là “BradBot”, chủ yếu vì chúng tôi vừa ngưỡng mộ vừa hơi ghen tị.

{{image "contributor-workshop/image9.jpg"}}
{{image "contributor-workshop/image6.png"}}

### Tác động

Chúng tôi có tổng cộng 65 CL được gửi từ những người tham gia workshop (trong vòng một tuần kể từ workshop). Trong số đó, 44 CL đến từ các contributor trước đó chưa từng đóng góp cho bất kỳ repo nào trong dự án Go. Một nửa (22) của các đóng góp này đã được merge. Nhiều cái còn lại đang chờ codebase tan băng vì chúng tôi đang ở giữa giai đoạn đóng băng cho bản phát hành 1.9 sắp tới. Ngoài CL, nhiều người còn đóng góp cho dự án dưới dạng báo lỗi, [gardening task](/wiki/Gardening), và các kiểu đóng góp khác.

Loại đóng góp phổ biến nhất là các hàm ví dụ để dùng trong tài liệu. [Khảo sát người dùng Go](/blog/survey2016-results) xác định rằng tài liệu của chúng tôi thiếu ví dụ nghiêm trọng. Trong bài trình bày, chúng tôi đề nghị người dùng tìm một package mà họ yêu thích và thêm một ví dụ. Trong dự án Go, ví dụ được viết dưới dạng mã trong các tệp Go (theo quy ước đặt tên cụ thể) và công cụ `go doc` hiển thị chúng cùng với tài liệu. Đây là một đóng góp đầu tiên hoàn hảo vì nó là thứ có thể được merge trong lúc đóng băng, có tầm quan trọng lớn với người dùng, và là một bổ sung có phạm vi tương đối hẹp.

Một trong những ví dụ được thêm vào là ví dụ tạo Stringer, một trong những interface được dùng rộng rãi hơn trong Go.
[CL 49270](/cl/49270/)

Ngoài ví dụ, nhiều người còn đóng góp những bản sửa lỗi quan trọng bao gồm:

  - [CL 48988](/cl/48988/) sửa [issue #21029](/issue/21029)
  - [CL 49050](/cl/49050/) sửa [issue #20054](/issue/20054)
  - [CL 49031](/cl/49031/) sửa [issue #20166](/issue/20166)
  - [CL 49170](/cl/49170/) sửa [issue #20877](/issue/20877)

Một số người thậm chí còn làm chúng tôi ngạc nhiên khi tới với sẵn một lỗi mà họ muốn sửa. Nikhita đến nơi sẵn sàng xử lý [issue #20786](/issue/20786) và cô ấy đã gửi [CL 48871](/cl/48871/), sau đó tweet rằng:

{{image "contributor-workshop/image19.png"}}

Không chỉ có nhiều cải tiến tuyệt vời được thực hiện, mà quan trọng hơn, chúng tôi đã thu hẹp khoảng cách giữa đội Go cốt lõi và các thành viên cộng đồng rộng lớn hơn. Nhiều người trong đội Go nhận xét rằng các thành viên cộng đồng đang dạy cho họ những điều về dự án Go. Mọi người trong cộng đồng (trực tiếp và trên Twitter) nói rằng họ cảm thấy được chào đón để tham gia vào dự án.

{{image "contributor-workshop/image12.png"}}
{{image "contributor-workshop/image13.png"}}
{{image "contributor-workshop/image3.png"}}

### Tương lai

Sự kiện thành công vượt xa kỳ vọng của chúng tôi. Sameer Ajmani, quản lý đội Go, nói: “Workshop contributor cực kỳ vui và bổ ích cho đội Go. Chúng tôi nhăn mặt khi người dùng chạm vào những góc cạnh xấu xí trong quy trình của chúng tôi, và vui mừng khi họ xuất hiện trên dashboard. Tiếng hò reo khi điểm số cả nhóm chạm mốc 1000 thật tuyệt.”

Chúng tôi đang tìm cách làm workshop này dễ tổ chức hơn cho các sự kiện tương lai (như meetup và hội nghị). Thách thức lớn nhất của chúng tôi là cung cấp đủ mentoring để người dùng cảm thấy được hỗ trợ. Nếu bạn có ý tưởng hoặc muốn giúp quy trình này, vui lòng [cho tôi biết](mailto:spf@golang.org).

Tôi đã nhờ một vài người tham gia sự kiện chia sẻ trải nghiệm của họ bên dưới:

## Trải nghiệm đóng góp của tôi

by [Cassandra](https://twitter.com/cassandraoid)

Khi nghe về workshop go-contrib tôi đã rất hào hứng và sau đó là cực kỳ e ngại. Tôi được một thành viên của đội Go khuyến khích tham gia, nên tôi nghĩ sao không thử xem.

Khi bước vào phòng (nói thật đi, tôi chạy vào phòng vì bị muộn) tôi rất vui khi thấy căn phòng chật kín. Tôi nhìn quanh tìm những người đội mũ Gopher, dấu hiệu chính cho thấy họ là giáo viên. Tôi ngồi xuống một trong 16 bàn tròn có hai chiếc mũ và ba người không có mũ. Bật màn hình lên và sẵn sàng bắt đầu…

Jess Frazelle đứng lên bắt đầu bài trình bày và cung cấp cho cả nhóm [một liên kết](https://docs.google.com/presentation/d/1ap2fycBSgoo-jCswhK9lqgCIFroE1pYpsXC1ffYBCq4/edit#slide=id.p) để mọi người dễ theo dõi.

{{image "contributor-workshop/image16.png"}}

Những tiếng xì xào từ nền thấp dần biến thành một bản hòa âm vang vọng của tiếng nói, mọi người đang thiết lập máy tính với Go, họ bỏ qua một vài phần để bảo đảm GOPATH đã được thiết lập, và rồi… khoan, Gerrit là gì?

Đa số chúng tôi phải được giới thiệu sơ qua về Gerrit. Tôi không hề biết nó là gì, nhưng may thay đã có sẵn một slide hữu ích. Jess giải thích rằng đó là một lựa chọn khác thay cho GitHub với các công cụ code review nâng cao hơn một chút. Sau đó chúng tôi đi qua các thuật ngữ GitHub vs Geritt để hiểu hơn về quy trình.

{{image "contributor-workshop/image10.png"}}

Được rồi, giờ là lúc trở thành một **Go contributor đúng nghĩa**.

Để khiến chuyện này còn hào hứng hơn mức vốn dĩ của nó, đội Go đã dựng một trò chơi nơi chúng tôi có thể theo dõi cả nhóm kiếm được bao nhiêu điểm dựa trên hệ thống điểm Gerrit.

{{image "contributor-workshop/image7.png"}}

Nhìn thấy tên mình xuất hiện trên bảng và nghe sự phấn khích của mọi người là điều gây nghiện. Nó cũng gợi lên cảm giác đồng đội dẫn tới cảm giác hòa nhập và cảm giác rằng bạn thật sự là một phần của cộng đồng Go.

{{image "contributor-workshop/image11.png"}}

Chỉ trong 6 bước, một căn phòng khoảng 80 người đã có thể học cách đóng góp cho Go trong vòng một giờ. Đó là một kỳ tích!

Nó không hề khó như tôi dự đoán và cũng không nằm ngoài tầm với của một người hoàn toàn mới. Nó nuôi dưỡng cảm giác cộng đồng theo một cách chủ động và hữu hình cũng như cảm giác được bao gồm trong quy trình đóng góp lẫy lừng của Go.

Tôi thực sự muốn cảm ơn đội Go, những mentor Gopher đội mũ, và những người tham gia cùng tôi vì đã biến nó thành một trong những khoảnh khắc đáng nhớ nhất của tôi tại GopherCon.

## Trải nghiệm đóng góp của tôi

by [Matt](https://twitter.com/mbbroberg)

Tôi luôn thấy các ngôn ngữ lập trình có phần đáng sợ. Chúng là thứ cho phép cả thế giới viết mã. Với tầm ảnh hưởng đó, hẳn những người thông minh hơn tôi mới nên làm việc trên nó... nhưng nỗi sợ ấy là thứ cần vượt qua. Vậy nên khi có cơ hội tham gia workshop để đóng góp cho ngôn ngữ lập trình yêu thích mới của mình, tôi rất hào hứng muốn xem mình có thể giúp được gì. Một tháng sau, giờ đây tôi chắc chắn rằng bất kỳ ai và tất cả mọi người đều có thể (và nên) đóng góp lại cho Go.

Đây là các bước rất dài dòng của tôi để đi từ 0 tới 2 đóng góp cho Go:

### Thiết lập

Vì Go dùng Gerrit, tôi bắt đầu bằng việc thiết lập môi trường cho nó. [Hướng dẫn của Jess Frazzelle](https://docs.google.com/presentation/d/1ap2fycBSgoo-jCswhK9lqgCIFroE1pYpsXC1ffYBCq4/edit#slide=id.g1f953ef7df_0_9) là nơi bắt đầu rất tốt để không bỏ sót bước nào.

Phần thú vị thực sự bắt đầu khi bạn clone repo Go. Trớ trêu thay, bạn không hack vào Go dưới `$GOPATH`, vì vậy tôi đặt nó trong workspace khác của mình (là `~/Develop`).

	cd $DEV # Đó là thư mục mã nguồn của tôi nằm ngoài $GOPATH
	git clone --depth 1 https://go.googlesource.com/go

Sau đó cài công cụ trợ giúp tiện lợi `go-contrib-init`:

	go get -u golang.org/x/tools/cmd/go-contrib-init

Giờ bạn có thể chạy `go-contrib-init` từ thư mục `go/` vừa clone ở trên và xem liệu mình đã sẵn sàng đóng góp hay chưa. Nhưng khoan, nếu bạn đang làm theo, bạn chưa sẵn sàng ngay đâu.

Tiếp theo, cài `codereview` để có thể tham gia review mã trên Gerrit:

	go get -u golang.org/x/review/git-codereview

Package này bao gồm `git change` và `git mail`, sẽ thay thế workflow bình thường của bạn là `git commit` và `git push`.

Được rồi, phần cài đặt đã xong. Giờ hãy thiết lập [tài khoản Gerrit ở đây](https://go-review.googlesource.com/settings/#Profile), rồi [ký CLA](https://go-review.googlesource.com/settings#Agreements) phù hợp với bạn (tôi ký CLA cá nhân cho tất cả dự án Google, nhưng hãy chọn lựa chọn đúng cho bạn. Bạn có thể xem tất cả các CLA mình đã ký tại [cla.developers.google.com/clas](https://cla.developers.google.com/clas)).

VÀ BÙM. Bạn ổn rồi (để bắt đầu)! Nhưng đóng góp ở đâu?

### Đóng góp

Trong workshop, họ đưa chúng tôi vào repository `scratch`, một nơi an toàn để nghịch ngợm nhằm thành thạo quy trình:

	cd $(go env GOPATH)/src/golang.org/x
	git clone --depth 1 [[https://go.googlesource.com/scratch][go.googlesource.com/scratch]]

Điểm dừng đầu tiên là `cd` vào đó và chạy `go-contrib-init` để chắc rằng bạn đã sẵn sàng đóng góp:

	go-contrib-init
	All good. Happy hacking!

Từ đó, tôi tạo một thư mục mang tên tài khoản GitHub của mình, chạy `git add -u` rồi thử `git change`. Nó có một hash dùng để theo dõi công việc của bạn, đó là một dòng duy nhất bạn không nên đụng vào. Ngoài điều đó, nó có cảm giác giống `git commit`. Một khi tôi có được thông điệp commit đúng định dạng `package: description` (description bắt đầu bằng chữ thường), tôi dùng `git mail` để gửi nó lên Gerrit.

Hai ghi chú hay ở thời điểm này: `git change` cũng hoạt động như `git commit --amend`, nên nếu bạn cần cập nhật patch thì có thể `add` rồi `change` và mọi thứ sẽ gắn với cùng một patch. Thứ hai, bạn luôn có thể xem lại patch của mình từ [dashboard Gerrit cá nhân](https://go-review.googlesource.com/dashboard/).

Sau một vài lần qua lại, tôi chính thức có một đóng góp cho Go! Và nếu Jaana nói đúng thì đó có thể là cái đầu tiên có emoji ✌️.

{{image "contributor-workshop/image15.png"}}
{{image "contributor-workshop/image23.png"}}

### Đóng góp, cho thật

Repo scratch thì vui đấy, nhưng có vô vàn cách để lặn sâu vào các package của Go và đóng góp lại. Đến thời điểm này tôi đã đi dạo quanh rất nhiều package để xem thứ gì khả dụng và thú vị với mình. Và khi nói “đi dạo quanh”, ý tôi là cố tìm một danh sách package, rồi quay sang mã nguồn để xem bên dưới thư mục `go/src/` có gì:

{{image "contributor-workshop/image22.png"}}

Tôi quyết định xem mình có thể làm gì trong package `regexp`, có lẽ vì vừa yêu vừa sợ regex. Đây là lúc tôi chuyển sang [góc nhìn package trên website](https://godoc.org/regexp) (nên biết rằng mỗi package chuẩn đều có thể tìm thấy tại https://godoc.org/$PACKAGENAME). Ở đó tôi nhận thấy `QuoteMeta` thiếu cùng mức ví dụ chi tiết mà các hàm khác có (và tôi cũng cần luyện tập Gerrit).

{{image "contributor-workshop/image1.png"}}

Tôi bắt đầu xem `go/src/regexp` để tìm nơi thêm ví dụ và bị lạc rất nhanh. May cho tôi là [Francesc](https://twitter.com/francesc) ở quanh đó hôm ấy. Anh ấy hướng dẫn tôi rằng mọi ví dụ thực chất là các kiểm thử nội tuyến trong tệp `example_test.go`. Chúng đi theo định dạng các ca kiểm thử theo sau bởi phần “Output” được comment lại và rồi là đáp án của kiểm thử. Ví dụ:

	func ExampleRegexp_FindString() {
		re := regexp.MustCompile("fo.?")
		fmt.Printf("%q\n", re.FindString("seafood"))
		fmt.Printf("%q\n", re.FindString("meat"))
		// Output:
		// "foo"
		// ""
	}

Hay ho đúng không?? Tôi làm theo Francesc và thêm hàm `ExampleQuoteMeta` cùng một vài trường hợp mà tôi nghĩ sẽ hữu ích. Từ đó chỉ cần `git change` rồi `git mail` lên Gerrit!

Tôi phải nói rằng Steve Francia đã thách tôi “hãy tìm một thứ không phải open issue rồi sửa nó”, nên tôi đã đưa vài thay đổi tài liệu cho QuoteMeta vào patch của mình. Nó sẽ còn mở thêm chút nữa vì phạm vi tăng lên, nhưng tôi nghĩ lần này điều đó xứng đáng.

Tôi đoán tôi đã nghe thấy câu hỏi của bạn rồi: làm sao tôi kiểm tra rằng nó hoạt động? Thành thật mà nói thì không dễ. Chạy `go test example_test.go -run QuoteMeta -v` sẽ không làm được vì chúng tôi đang làm việc ngoài `$GOPATH`. Tôi loay hoay mãi để tìm ra cách cho đến khi [Kale Blakenship viết bài tuyệt vời này về testing trong Go](https://medium.com/@vCabbage/go-testing-standard-library-changes-1e9cbed11339). Hãy đánh dấu lại để đọc sau.

Bạn có thể xem [đóng góp hoàn chỉnh của tôi ở đây](https://go-review.googlesource.com/c/49130/). Điều tôi cũng hy vọng bạn sẽ thấy là việc đi vào guồng đóng góp thật sự đơn giản đến mức nào. Nếu bạn giống tôi, bạn sẽ giỏi tìm ra một lỗi gõ nhỏ hoặc ví dụ còn thiếu trong tài liệu để bắt đầu làm quen với workflow `git codereview`. Sau đó, bạn sẽ sẵn sàng tìm một open issue, lý tưởng là issue [được gắn nhãn cho bản phát hành sắp tới](https://github.com/golang/go/milestones), và thử sức với nó. Bất kể bạn chọn làm gì, hãy cứ tiến lên và làm. Đội Go đã chứng minh cho tôi thấy họ thực sự quan tâm đến việc giúp tất cả chúng ta đóng góp lại như thế nào. Tôi không thể chờ đến lần `git mail` tiếp theo của mình.

## Trải nghiệm làm mentor của tôi

by [Dmitri](https://twitter.com/dmitshur)

Tôi đã rất mong chờ tham gia sự kiện Contribution Workshop với vai trò mentor. Tôi có kỳ vọng cao cho sự kiện này, và đã nghĩ đó là một ý tưởng tuyệt vời ngay cả trước khi nó bắt đầu.

Tôi thực hiện đóng góp đầu tiên cho Go vào ngày 10 tháng 5 năm 2014. Tôi nhớ đó là khoảng bốn tháng kể từ thời điểm tôi muốn đóng góp, cho tới ngày tôi thực sự gửi CL đầu tiên. Tôi mất từng ấy thời gian để gom đủ can đảm và thật sự cam kết tìm hiểu quy trình. Lúc đó tôi đã là một kỹ sư phần mềm giàu kinh nghiệm. Dù vậy, quy trình đóng góp cho Go vẫn cho cảm giác xa lạ, khác với tất cả những quy trình tôi đã quen thuộc, nên trông có vẻ đáng sợ. Tuy nhiên nó được tài liệu hóa rất tốt, nên tôi biết đó chỉ là chuyện dành thời gian, ngồi xuống và làm. Yếu tố “chưa biết” đã ngăn tôi thử.

Sau vài tháng trôi qua, tôi nghĩ “thế là đủ rồi”, và quyết định dành trọn một ngày cuối tuần sắp tới để tìm hiểu quy trình. Tôi dành cả ngày thứ Bảy cho một việc duy nhất: gửi CL đầu tiên của mình cho Go. Tôi mở [Contribution Guide](/doc/contribute.html) và bắt đầu làm theo từng bước từ trên xuống. Trong vòng một giờ, tôi đã xong. Tôi đã gửi CL đầu tiên của mình. Tôi vừa kinh ngạc vừa sốc. Kinh ngạc, vì cuối cùng tôi đã gửi được một đóng góp cho Go, và nó còn được chấp nhận! Sốc, vì sao tôi lại chờ quá lâu mới làm chuyện này? Làm theo các bước trong [Contribution Guide](/doc/contribute.html) thật sự rất dễ, và toàn bộ quy trình diễn ra hoàn toàn trơn tru. Giá như ai đó đã nói với tôi rằng tôi sẽ xong trong một giờ và không có chuyện gì trục trặc, hẳn tôi đã làm sớm hơn rất nhiều!

Điều đó đưa tôi trở lại với sự kiện này và lý do tôi nghĩ nó là một ý tưởng hay đến vậy. Với bất kỳ ai từng muốn đóng góp cho Go nhưng cảm thấy e dè vì quy trình xa lạ và có vẻ dài dòng (giống như tôi trong bốn tháng đó), đây là cơ hội của họ! Không chỉ dễ cam kết sẽ tìm hiểu nó bằng cách tham dự sự kiện, mà đội Go và các mentor tình nguyện nhiệt tình cũng có mặt để giúp bạn đi qua từng bước.

Dù đã có kỳ vọng cao từ đầu, kỳ vọng của tôi về sự kiện vẫn còn bị vượt qua. Trước hết, đội Go đã chuẩn bị cực kỳ kỹ và đầu tư rất nhiều để khiến sự kiện trở nên vui hơn cho tất cả mọi người. Có một bài trình bày rất thú vị đi qua mọi bước đóng góp một cách nhanh gọn. Có một dashboard được làm riêng cho sự kiện, nơi mọi bước hoàn thành thành công của mỗi người đều được thưởng điểm cộng vào điểm chung toàn cầu. Điều đó biến nó thành một sự kiện rất cộng tác và đậm tính xã hội! Cuối cùng, và quan trọng nhất, là có các thành viên đội Go như Brad Fitzpatrick phía sau hậu trường, giúp review CL rất nhanh! Điều đó có nghĩa là các CL được gửi lên nhận review nhanh chóng, kèm các bước tiếp theo có thể hành động được, để mọi người có thể tiếp tục tiến lên và học thêm.

Ban đầu tôi dự đoán sự kiện sẽ phần nào tẻ nhạt, ở chỗ các bước đóng góp cực kỳ đơn giản để làm theo. Tuy nhiên tôi thấy điều đó không phải lúc nào cũng đúng, và tôi đã có thể dùng chuyên môn Go của mình để giúp những người bị mắc kẹt ở nhiều chỗ ngoài dự kiến. Hóa ra thế giới thực đầy những edge case. Ví dụ, có người có hai email git, một cá nhân và một cho công việc. Có sự chậm trễ trong việc ký CLA cho email công việc, nên họ thử dùng email cá nhân thay thế. Điều đó có nghĩa là mỗi commit phải được amend lại để dùng đúng email, điều mà các công cụ không tính tới. (May mắn là trong contribution guide có hẳn một phần troubleshooting nói chính xác về vấn đề này!) Có các lỗi tinh vi khác hoặc môi trường cấu hình sai mà một số người gặp phải, vì có hơn một bản cài Go là chuyện hơi bất thường. Đôi khi biến môi trường GOROOT phải được đặt tường minh tạm thời để godoc hiển thị các thay đổi trong đúng thư viện chuẩn (tôi còn đùa mà ngoái nhìn sau lưng xem Dave Cheney có nghe thấy khi mình nói những lời đó không).

Nhìn chung, tôi đã chứng kiến một vài gopher mới thực hiện những đóng góp Go đầu tiên của họ. Họ gửi CL, phản hồi feedback review, chỉnh sửa, lặp đi lặp lại cho tới khi mọi người đều hài lòng, và cuối cùng thấy đóng góp đầu tiên của mình cho Go được merge vào master! Rất đáng giá khi được nhìn thấy niềm vui trên khuôn mặt họ, bởi niềm vui khi thực hiện đóng góp đầu tiên là điều tôi hoàn toàn đồng cảm. Tôi cũng rất vui vì có thể giúp họ và giải thích những tình huống rắc rối mà đôi khi họ tự rơi vào. Từ những gì tôi thấy, rất nhiều gopher hạnh phúc đã rời khỏi sự kiện này, bao gồm cả tôi!

## Ảnh từ sự kiện

{{image "contributor-workshop/image2.jpg"}}
{{image "contributor-workshop/image4.jpg"}}
{{image "contributor-workshop/image5.jpg"}}
{{image "contributor-workshop/image8.jpg"}}
{{image "contributor-workshop/image14.jpg"}}
{{image "contributor-workshop/image18.jpg"}}
{{image "contributor-workshop/image20.jpg"}}
{{image "contributor-workshop/image21.jpg"}}

Ảnh bởi Sameer Ajmani & Steve Francia
