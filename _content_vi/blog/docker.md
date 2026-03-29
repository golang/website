---
title: Triển khai máy chủ Go với Docker
date: 2014-09-26
by:
- Andrew Gerrand
summary: Cách dùng các ảnh nền chính thức mới của Docker cho Go.
template: true
---

## Giới thiệu

Tuần này Docker đã [công bố](https://blog.docker.com/2014/09/docker-hub-official-repos-announcing-language-stacks/)
các ảnh nền chính thức cho Go và các ngôn ngữ lớn khác,
giúp lập trình viên có một cách đáng tin cậy và dễ dàng để xây dựng container cho các chương trình Go của họ.

Trong bài viết này, chúng ta sẽ đi qua một công thức để tạo một Docker container cho
một ứng dụng web Go đơn giản và triển khai container đó lên Google Compute Engine.
Nếu bạn chưa quen với Docker, bạn nên đọc
[Understanding Docker](https://docs.docker.com/engine/understanding-docker/)
trước khi tiếp tục.

## Ứng dụng minh họa

Để minh họa, chúng ta sẽ dùng chương trình
[outyet](https://pkg.go.dev/golang.org/x/example/outyet) từ
[kho ví dụ của Go](https://cs.opensource.google/go/x/example),
một máy chủ web đơn giản báo xem phiên bản Go tiếp theo đã được phát hành hay chưa
(được thiết kế để vận hành các trang như [isgo1point4.outyet.org](http://isgo1point4.outyet.org/)).
Nó không có phụ thuộc nào ngoài thư viện chuẩn và không cần thêm
tệp dữ liệu nào lúc chạy; với một máy chủ web, nó gần như đơn giản hết mức có thể.

Dùng `go get` để tải và cài đặt outyet vào
[workspace](/doc/code.html#Workspaces) của bạn:

	$ go get golang.org/x/example/outyet

## Viết Dockerfile

Hãy thay nội dung tệp tên `Dockerfile` trong thư mục `outyet` bằng nội dung sau:

	# Start from a Debian image with the latest version of Go installed
	# and a workspace (GOPATH) configured at /go.
	FROM golang

	# Copy the local package files to the container's workspace.
	ADD . /go/src/golang.org/x/example/outyet

	# Build the outyet command inside the container.
	# (You may fetch or manage dependencies here,
	# either manually or with a tool like "godep".)
	RUN go install golang.org/x/example/outyet

	# Run the outyet command by default when the container starts.
	ENTRYPOINT /go/bin/outyet

	# Document that the service listens on port 8080.
	EXPOSE 8080

`Dockerfile` này chỉ rõ cách dựng một container chạy `outyet`,
bắt đầu từ các phụ thuộc cơ bản (một hệ Debian có cài Go;
[ảnh Docker `golang` chính thức](https://registry.hub.docker.com/_/golang/)),
thêm mã nguồn của package `outyet`, biên dịch nó, rồi cuối cùng chạy nó.

Các bước `ADD`, `RUN` và `ENTRYPOINT` là các tác vụ phổ biến với bất kỳ dự án Go nào.
Để đơn giản hóa điều này, có một
[`onbuild` variant](https://github.com/docker-library/golang/blob/9ff2ccca569f9525b023080540f1bb55f6b59d7f/1.3.1/onbuild/Dockerfile)
của ảnh `golang` tự động sao chép mã nguồn package, tải các phụ thuộc của ứng dụng,
biên dịch chương trình và cấu hình để nó chạy khi khởi động.

Với biến thể `onbuild`, `Dockerfile` đơn giản hơn nhiều:

	FROM golang:onbuild
	EXPOSE 8080

## Xây dựng và chạy ảnh

Gọi Docker từ thư mục package `outyet` để xây dựng một ảnh dùng `Dockerfile`:

	$ docker build -t outyet .

Lệnh này sẽ tải ảnh nền `golang` từ Docker Hub, sao chép mã nguồn package
vào đó, biên dịch package bên trong nó, rồi gắn nhãn ảnh kết quả là `outyet`.

Để chạy một container từ ảnh vừa tạo:

	$ docker run --publish 6060:8080 --name test --rm outyet

Cờ `--publish` bảo Docker công bố cổng `8080` của container trên
cổng ngoài `6060`.

Cờ `--name` đặt cho container một tên dễ đoán để việc thao tác với nó thuận tiện hơn.

Cờ `--rm` bảo Docker xóa ảnh container khi máy chủ outyet thoát.

Khi container đang chạy, mở `http://localhost:6060/` trong trình duyệt web
và bạn sẽ thấy đại loại như thế này:

{{image "docker/outyet.png"}}

(Nếu Docker daemon của bạn đang chạy trên một máy khác, hoặc trong máy ảo,
bạn nên thay `localhost` bằng địa chỉ của máy đó. Nếu bạn đang
dùng [boot2docker](http://boot2docker.io/) trên OS X hoặc Windows, bạn có thể tìm
địa chỉ đó bằng `boot2docker ip`.)

Sau khi đã xác minh rằng ảnh hoạt động, hãy tắt container đang chạy
từ một cửa sổ terminal khác:

	$ docker stop test

## Tạo kho trên Docker Hub

[Docker Hub](https://hub.docker.com/), kho lưu trữ container mà từ đó chúng ta
đã kéo ảnh `golang` trước đó, cung cấp một tính năng gọi là
[Automated Builds](http://docs.docker.com/docker-hub/builds/) để xây dựng
ảnh từ một kho GitHub hoặc BitBucket.

Bằng cách commit [Dockerfile](https://go.googlesource.com/example/+/refs/heads/master/outyet/)
vào kho mã và tạo một
[automated build](https://registry.hub.docker.com/u/adg1/outyet/)
cho nó, bất kỳ ai cài Docker đều có thể tải và chạy ảnh của chúng ta chỉ bằng một
lệnh duy nhất. (Ta sẽ thấy ích lợi của điều này ở phần tiếp theo.)

Để thiết lập Automated Build, hãy commit Dockerfile vào kho mã của bạn trên
[GitHub](https://github.com/) hoặc [BitBucket](https://bitbucket.org/),
tạo một tài khoản trên Docker Hub và làm theo hướng dẫn
[tạo Automated Build](http://docs.docker.com/docker-hub/builds/).

Khi hoàn tất, bạn có thể chạy container của mình bằng tên của automated build:

	$ docker run goexample/outyet

(Hãy thay `goexample/outyet` bằng tên automated build mà bạn đã tạo.)

## Triển khai container lên Google Compute Engine

Google cung cấp các
[ảnh Google Compute Engine tối ưu cho container](https://developers.google.com/compute/docs/containers/container_vms)
giúp việc khởi tạo một máy ảo chạy một Docker container bất kỳ trở nên dễ dàng.
Khi khởi động, một chương trình chạy trên instance sẽ đọc một tệp cấu hình
chỉ rõ container nào cần chạy, tải ảnh container và chạy nó.

Hãy tạo một tệp [containers.yaml](https://cloud.google.com/compute/docs/containers/container_vms#container_manifest)
chỉ rõ ảnh docker cần chạy và các cổng cần mở:

	version: v1beta2
	containers:
	- name: outyet
	  image: goexample/outyet
	  ports:
	  - name: http
	    hostPort: 80
	    containerPort: 8080

(Lưu ý rằng chúng ta đang công bố cổng `8080` của container ra cổng ngoài `80`,
là cổng mặc định để phục vụ lưu lượng HTTP. Và một lần nữa, bạn nên thay
`goexample/outyet` bằng tên Automated Build của mình.)

Dùng [công cụ gcloud](https://cloud.google.com/sdk/#Quick_Start)
để tạo một VM instance chạy container:

	$ gcloud compute instances create outyet \
		--image container-vm-v20140925 \
		--image-project google-containers \
		--metadata-from-file google-container-manifest=containers.yaml \
		--tags http-server \
		--zone us-central1-a \
		--machine-type f1-micro

Đối số đầu tiên (`outyet`) chỉ rõ tên instance, một nhãn tiện lợi
cho mục đích quản trị.

Cờ `--image` và `--image-project` chỉ rõ ảnh hệ thống đặc biệt
tối ưu cho container cần dùng (hãy sao chép nguyên văn các cờ này).

Cờ `--metadata-from-file` cung cấp tệp `containers.yaml` của bạn cho VM.

Cờ `--tags` gắn nhãn VM instance của bạn là một máy chủ HTTP, điều chỉnh
firewall để lộ cổng 80 trên giao diện mạng công khai.

Các cờ `--zone` và `--machine-type` chỉ rõ vùng chạy VM
và loại máy cần dùng. (Để xem danh sách loại máy và vùng,
hãy chạy `gcloud compute machine-types list`.)

Sau khi hoàn tất, lệnh gcloud sẽ in ra một số thông tin về
instance. Trong đầu ra đó, hãy tìm phần `networkInterfaces` để biết
địa chỉ IP ngoài của instance. Trong vòng vài phút, bạn sẽ có thể
truy cập địa chỉ IP đó bằng trình duyệt web và thấy trang “Has Go 1.4 been released
yet?”.

(Để xem chuyện gì đang diễn ra trên VM instance mới, bạn có thể ssh vào nó bằng
`gcloud compute ssh outyet`. Từ đó, thử `sudo docker ps` để xem
những Docker container nào đang chạy.)

## Tìm hiểu thêm

Đây mới chỉ là phần nổi của tảng băng, còn rất nhiều điều bạn có thể làm với Go, Docker và Google Compute Engine.

Để tìm hiểu thêm về Docker, hãy xem [bộ tài liệu rất đầy đủ của họ](https://docs.docker.com/).

Để tìm hiểu thêm về Docker và Go, hãy xem [kho Docker Hub `golang` chính thức](https://registry.hub.docker.com/_/golang/) và bài viết của Kelsey Hightower [Optimizing Docker Images for Static Go Binaries](https://medium.com/@kelseyhightower/optimizing-docker-images-for-static-binaries-b5696e26eb07).

Để tìm hiểu thêm về Docker và [Google Compute Engine](http://cloud.google.com/compute),
hãy xem [trang Container-optimized VMs](https://cloud.google.com/compute/docs/containers/container_vms)
và [kho Docker Hub google/docker-registry](https://registry.hub.docker.com/u/google/docker-registry/).
