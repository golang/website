---
title: Làm quen với workspaces
date: 2022-04-05
by:
- Beth Brown, đại diện cho nhóm Go
tags:
- go
- workspaces
- go1.18
summary: Tìm hiểu về Go workspaces và một số quy trình làm việc mà chúng cho phép.
---

Go 1.18 bổ sung workspace mode cho Go, cho phép bạn làm việc trên nhiều module
đồng thời.

Bạn có thể lấy Go 1.18 tại trang [download](/dl/). [Release notes](/doc/go1.18)
có thêm chi tiết về tất cả các thay đổi.

## Workspaces

[Workspaces](/ref/mod#workspaces) trong Go 1.18 cho phép bạn làm việc trên
nhiều module cùng lúc mà không cần sửa tệp `go.mod` của từng
module. Mỗi module trong một workspace được coi là một main module khi
giải quyết phụ thuộc.

Trước đây, để thêm một tính năng vào một module rồi dùng nó trong module khác, bạn
cần либо phát hành các thay đổi lên module thứ nhất, либо [sửa
tệp go.mod](/doc/tutorial/call-module-code) của module phụ thuộc
với một chỉ thị `replace` trỏ tới các thay đổi cục bộ chưa phát hành đó. Để
phát hành mà không có lỗi, bạn lại phải xóa chỉ thị `replace` khỏi
tệp `go.mod` của module phụ thuộc sau khi đã phát hành thay đổi cục bộ trên
module thứ nhất.

Với Go workspaces, bạn kiểm soát mọi phụ thuộc bằng một tệp `go.work` ở
gốc thư mục workspace. Tệp `go.work` có các chỉ thị `use` và
`replace` ghi đè các tệp `go.mod` riêng lẻ, nên không còn
cần sửa từng tệp `go.mod` một.

Bạn tạo một workspace bằng cách chạy `go work init` với danh sách các thư mục
module làm đối số cách nhau bằng dấu cách. Workspace không nhất thiết phải chứa
các module mà bạn đang làm việc cùng. Lệnh `init` tạo một tệp `go.work`
liệt kê các module trong workspace. Nếu chạy `go work init` mà không có
đối số, lệnh sẽ tạo một workspace trống.

Để thêm module vào workspace, hãy chạy `go work use [moddir]` hoặc tự sửa
tệp `go.work`. Chạy `go work use -r .` để thêm đệ quy các thư mục trong
thư mục đối số mà có tệp `go.mod` vào workspace. Nếu một thư mục
không có tệp `go.mod`, hoặc không còn tồn tại nữa, thì chỉ thị `use` cho
thư mục đó sẽ bị xóa khỏi tệp `go.work`.

Cú pháp của tệp `go.work` tương tự tệp `go.mod` và có các
chỉ thị sau:

- `go`: phiên bản go toolchain, ví dụ `go 1.18`
- `use`: thêm một module trên đĩa vào tập main modules của workspace.
  Đối số của nó là một đường dẫn tương đối tới thư mục chứa
  tệp `go.mod` của module. Chỉ thị `use` không thêm các module trong thư mục con
  của thư mục đã chỉ định.
- `replace`: Tương tự chỉ thị `replace` trong tệp `go.mod`, một
  chỉ thị `replace` trong tệp `go.work` thay thế nội dung của
  _một phiên bản cụ thể_ của một module, hoặc _mọi phiên bản_ của một module, bằng
  nội dung nằm ở nơi khác.

## Quy trình làm việc

Workspaces rất linh hoạt và hỗ trợ nhiều kiểu quy trình làm việc. Các phần sau
là cái nhìn tổng quan ngắn gọn về những gì chúng tôi nghĩ sẽ là phổ biến nhất.

### Thêm một tính năng vào module upstream và dùng nó trong chính module của bạn

1. Tạo một thư mục cho workspace.
2. Clone module upstream mà bạn muốn sửa.
3. Thêm tính năng của bạn vào phiên bản cục bộ của module upstream.
4. Chạy `go work init [path-to-upstream-mod-dir]` trong thư mục workspace.
5. Sửa module của chính bạn để hiện thực tính năng vừa thêm
   vào module upstream.
6. Chạy `go work use [path-to-your-module]` trong thư mục workspace.

   Lệnh `go work use` thêm đường dẫn đến module của bạn vào tệp `go.work`:

   ```
   go 1.18

   use (
          ./path-to-upstream-mod-dir
          ./path-to-your-module
   )
   ```

7. Chạy và kiểm thử module của bạn với tính năng mới vừa thêm vào module upstream.
8. Phát hành module upstream với tính năng mới.
9. Phát hành module của bạn dùng tính năng mới đó.

### Làm việc với nhiều module phụ thuộc lẫn nhau trong cùng một kho mã

Khi làm việc với nhiều module trong cùng một kho mã, tệp `go.work`
xác định workspace thay vì dùng chỉ thị `replace` trong từng
tệp `go.mod`.

1. Tạo một thư mục cho workspace.
2. Clone kho mã chứa các module bạn muốn sửa. Các module không nhất thiết
   phải nằm trong thư mục workspace vì bạn sẽ chỉ ra đường dẫn tương đối tới
   từng module bằng chỉ thị `use`.
3. Chạy `go work init [path-to-module-one] [path-to-module-two]` trong
   thư mục workspace.

   Ví dụ: Bạn đang làm việc trên `example.com/x/tools/groundhog`, vốn phụ thuộc
   vào các package khác trong module `example.com/x/tools`.

   Bạn clone kho mã rồi chạy `go work init tools tools/groundhog` trong
   thư mục workspace.

   Nội dung tệp `go.work` sẽ tương tự như sau:

   ```
   go 1.18

   use (
           ./tools
           ./tools/groundhog
   )
   ```

   Mọi thay đổi cục bộ trong module `tools` sẽ được
   `tools/groundhog` sử dụng trong workspace của bạn.

### Chuyển đổi giữa các cấu hình phụ thuộc

Để kiểm thử module của bạn với các cấu hình phụ thuộc khác nhau, bạn có thể либо
tạo nhiều workspace với các tệp `go.work` riêng, либо giữ một workspace
duy nhất rồi comment các chỉ thị `use` mà bạn không muốn trong một tệp `go.work`.

Để tạo nhiều workspace:

1. Tạo các thư mục riêng cho những nhu cầu phụ thuộc khác nhau.
2. Chạy `go work init` trong từng thư mục workspace.
3. Thêm các phụ thuộc bạn muốn trong mỗi thư mục bằng `go work use
   [path-to-dependency]`.
4. Chạy `go run [path-to-your-module]` trong từng thư mục workspace để dùng
   các phụ thuộc do tệp `go.work` của thư mục đó chỉ định.

Để thử các phụ thuộc khác nhau trong cùng một workspace, mở tệp `go.work`
và thêm hoặc comment các phụ thuộc mong muốn.

### Vẫn đang dùng GOPATH?

Có thể workspaces sẽ làm bạn đổi ý. Người dùng `GOPATH` có thể giải quyết
phụ thuộc của mình bằng một tệp `go.work` đặt ở gốc thư mục `GOPATH`.
Workspaces không nhắm tới việc tái tạo hoàn toàn mọi quy trình `GOPATH`,
nhưng chúng có thể tạo ra một cấu hình chia sẻ được phần nào sự tiện lợi của `GOPATH`
trong khi vẫn giữ được lợi ích của modules.

Để tạo một workspace cho GOPATH:

1. Chạy `go work init` ở gốc thư mục `GOPATH`.
2. Để dùng một module cục bộ hoặc một phiên bản cụ thể làm phụ thuộc trong
   workspace, hãy chạy `go work use [path-to-module]`.
3. Để thay thế các phụ thuộc hiện có trong tệp `go.mod` của các module, hãy dùng
   `go work replace [path-to-module]`.
4. Để thêm tất cả module trong GOPATH hoặc bất kỳ thư mục nào, hãy chạy `go work use
   -r` để thêm đệ quy các thư mục có tệp `go.mod` vào workspace.
   Nếu một thư mục không có tệp `go.mod`, hoặc không còn tồn tại, thì chỉ thị `use`
   cho thư mục đó sẽ bị xóa khỏi tệp `go.work`.

> Lưu ý: Nếu bạn có các dự án không có tệp `go.mod` mà muốn thêm vào
workspace, hãy chuyển vào thư mục dự án của chúng và chạy `go mod init`,
rồi thêm module mới vào workspace bằng `go work use [path-to-module].`

## Các lệnh workspace

Bên cạnh `go work init` và `go work use`, Go 1.18 còn giới thiệu các
lệnh sau cho workspaces:

- `go work sync`: đẩy các phụ thuộc trong tệp `go.work` ngược lại vào
  các tệp `go.mod` của từng module trong workspace.
- `go work edit`: cung cấp giao diện dòng lệnh để sửa `go.work`,
  chủ yếu dùng cho công cụ hoặc script.

Các lệnh build nhận biết module và một số subcommand của `go mod` sẽ xem xét biến môi trường `GOWORK`
để xác định xem chúng có đang ở trong ngữ cảnh workspace hay không.

Workspace mode được bật nếu biến `GOWORK` chỉ tới một đường dẫn tệp kết thúc bằng
`.work`. Để xác định tệp `go.work` nào đang được dùng, hãy chạy
`go env GOWORK`. Đầu ra sẽ rỗng nếu lệnh `go` không ở trong workspace
mode.

Khi workspace mode được bật, tệp `go.work` sẽ được parse để xác định ba
tham số của workspace mode: một phiên bản Go, một danh sách thư mục và một
danh sách thay thế.

Một vài lệnh có thể thử trong workspace mode (miễn là bạn đã biết chúng
làm gì!):

```
go work init
go work sync
go work use
go list
go build
go test
go run
go vet
```

## Cải tiến trải nghiệm trình soạn thảo

Chúng tôi đặc biệt hào hứng với những nâng cấp cho language server của Go là
[gopls](https://pkg.go.dev/golang.org/x/tools/gopls) và
[phần mở rộng VSCode Go](https://marketplace.visualstudio.com/items?itemName=golang.go),
giúp việc làm việc với nhiều module trong một trình soạn thảo tương thích LSP trở nên mượt mà
và đáng giá.

Tìm tham chiếu, tự động hoàn thành mã và đi tới định nghĩa hoạt động xuyên các module
trong workspace. Phiên bản [0.8.1](https://github.com/golang/tools/releases/tag/gopls%2Fv0.8.1)
của `gopls` bổ sung diagnostics, completion, formatting và hover cho
tệp `go.work`. Bạn có thể tận dụng các tính năng này của gopls với bất kỳ trình soạn thảo nào
tương thích [LSP](https://microsoft.github.io/language-server-protocol/).

#### Ghi chú riêng cho từng trình soạn thảo

- [Bản phát hành mới nhất của vscode-go](https://github.com/golang/vscode-go/releases/tag/v0.32.0)
  cho phép truy cập nhanh đến tệp `go.work` của workspace qua menu
  Quick Pick trên thanh trạng thái Go.

![Truy cập tệp go.work qua menu Quick Pick trên thanh trạng thái Go](https://user-images.githubusercontent.com/4999471/157268414-fba63843-5a14-44ba-be82-d42765568856.gif)

- [GoLand](https://www.jetbrains.com/go/) hỗ trợ workspaces và có kế hoạch bổ sung
  tô sáng cú pháp cùng tự động hoàn thành mã cho tệp `go.work`.

Để biết thêm thông tin về việc dùng `gopls` với các trình soạn thảo khác nhau, hãy xem `gopls` [documentation](https://pkg.go.dev/golang.org/x/tools/gopls#readme-editors).

## Tiếp theo là gì?

- Tải xuống và cài đặt [Go 1.18](/dl/).
- Thử dùng [workspaces](/ref/mod#workspaces) với [Go
  workspaces Tutorial](/doc/tutorial/workspaces).
- Nếu bạn gặp bất kỳ vấn đề nào với workspaces, hoặc muốn đề xuất điều gì,
  hãy tạo [issue](/issue/new).
- Đọc [tài liệu bảo trì workspace](https://pkg.go.dev/cmd/go#hdr-Workspace_maintenance).
- Khám phá các lệnh module để [làm việc bên ngoài một
  module đơn lẻ](/ref/mod#commands-outside), bao gồm `go work init`,
  `go work sync` và nhiều lệnh khác.
