---
title: "Go for Cloud & Network Services"
linkTitle: "Cloud & Network Services"
description: "With a strong ecosystem of tools and APIs on major cloud providers, it is easier than ever to build services with Go."
date: 2019-10-04T15:26:31-04:00
series: Use Cases
books:
  - title: Building Microservices with Go 
    url: https://www.amazon.com/Building-Microservices-Go-efficient-microservices/dp/1786468662/
    thumbnail: /images/books/building-microservices-with-go.jpg
  - title: Hands-On Software Architecture with Golang
    url: https://www.amazon.com/dp/1788622596/ref=cm_sw_r_tw_dp_U_x_-aZWDbS8PD7R4
    thumbnail: /images/books/hands-on-software-architecture-with-golang.jpg
  - title: Building RESTful Web services with Go
    url: https://www.amazon.com/Building-RESTful-Web-services-gracefully-ebook/dp/B072QB8KL1
    thumbnail: /images/books/building-restful-web-services-with-go.jpg
  - title: Mastering Go Web Services
    url: https://www.amazon.com/Mastering-Web-Services-Nathan-Kozyra-ebook/dp/B00W5GUKL6
    thumbnail: /images/books/mastering-go-web-services.jpg
  - title: 
    url: 
    thumbnail: 
resources:
- name: icon
  src: globe.png
  params:
    alt: globe
- name: icon-white
  src: globe-white.png
  params:
    alt: globe
---

## Go helps enterprises build and scale Cloud Computing systems

As applications and processing move to the cloud, concurrency becomes a very big issue. Cloud computing systems, by their very nature, share and scale resources. Coordinating access to shared resources is an issue that impacts every application processing in the cloud, and requires programming languages “[explicitly geared to develop highly reliable concurrent applications](https://tomassetti.me/best-programming-languages/).”
 
Go was created to address exactly these concurrency needs for scaled applications, microservices, and cloud development. In fact, **over 75 percent of projects in the Cloud Native Computing Foundation ([CNCF](https://www.cncf.io/projects/)) are written in Go**.

{{< pullquote author="Ruchi Malik, developer at Choozle" link="https://builtin.com/software-engineering-perspectives/golang-advantages" >}}
Go makes it very easy to scale as a company. This is very important because, as our engineering team grows, each service can be managed by a different unit.
{{< /pullquote >}}

Service developers often make a tradeoff between development cycle time and server performance. Go's helps reduce the need to make this tradeoff, with its fast build times that enable iterative development, lower memory and CPU utilization, making servers built with Go cheaper to run in pay-as-you-go deployments.

Go addresses many challenges developers face with the modern cloud, delivering standard idiomatic APIs, out-of-the-box support for multiple cloud environments (including on-prem), and a great balance between speed and efficiency - granting engineering teams the power to choose and the power to move.

## Featured Go users & projects

{{% mediaList %}}
    {{% mediaListBox img-src="/images/logos/google-cloud.svg" img-alt="Google Cloud Logo" title="" align=top %}}
Google Cloud is built on Go. Many critical cloud projects used across the industry like [Kubernetes](https://kubernetes.io/), [Istio](https://istio.io/) and [gVisor](https://gvisor.dev/) were created in Go at Google Cloud. Go is fully supported on [Google Cloud](https://cloud.google.com) across all APIs and runtimes including serverless [App Engine](https://cloud.google.com/appengine/) and [Google Cloud Functions](https://cloud.google.com/functions/). 
    {{% /mediaListBox %}}
    {{% mediaListBox img-src="/images/logos/capital-one.svg" img-alt="Capital One Logo" title="" align=top %}}
Capital One uses Go to power the [Credit Offers API, a critical service](https://medium.com/capital-one-tech/a-serverless-and-go-journey-credit-offers-api-74ef1f9fde7f). The engineering team is also building their serverless architecture with Go, citing Go's speed and simplicity, and mentioning that "[they] didn't want to go serverless without Go."
    {{% /mediaListBox %}}
    {{% mediaListBox img-src="/images/logos/dropbox.svg" img-alt="Dropbox Logo" title="" align=top %}}
Dropbox was built on Python, [but in 2013 decided to migrate](https://blogs.dropbox.com/tech/2014/07/open-sourcing-our-go-libraries/) their performance-critical backends to Go. Today, most of the company's infrastructure is written in Go.  Dropbox libraries can be found at [Dropbox's Go github](https://github.com/dropbox/godropbox).
    {{% /mediaListBox %}}
    {{% mediaListBox img-src="/images/logos/mercadoLibre.svg" img-alt="MercadoLibre Logo" title="" align=top %}}
MercadoLibre uses Go to [scale its eCommerce platform](/solutions/mercadolibre).  Go produces efficient code that readily scales as MercadoLibre’s online commerce grows and supports the company as a boon for developers—improving their productivity while streamlining and expanding MercadoLibre services.
    {{% /mediaListBox %}}
    {{% mediaListBox img-src="/images/logos/twitch.svg" img-alt="Twitch Logo" title="" align=top %}}
Twitch [uses Go to power many of its busiest systems](https://blog.twitch.tv/en/2016/07/05/gos-march-to-low-latency-gc-a6fa96f06eb7/) that serve live video and chat to millions of users. 
    {{% /mediaListBox %}}
    {{% mediaListBox img-src="/images/logos/the-new-york-times-icon.svg" img-alt="The New York Times Logo" title="" align=top %}}
The New York Times adopted Go ["to build better back-end services"](https://open.nytimes.com/introducing-gizmo-aa7ea463b208). As the usage of Go expanded with in the company they felt the need to create a toolkit to "to help developers quickly configure and build microservice APIs and pubsub daemons", which they have [open sourced](https://github.com/nytimes/gizmo). 
    {{% /mediaListBox %}}
    {{% mediaListBox img-src="/images/logos/uber-app-icon.svg" img-alt="Uber Logo" title="" align=top %}}
Uber uses Go to power several of its critical services that impact the experience of millions of drivers and passengers around the world. From their [real-time analytics engine](https://eng.uber.com/aresdb/), AresDB, to their [microservice for Geo-querying](https://eng.uber.com/go-geofence/), Geofence, and [their resource scheduler](https://eng.uber.com/open-sourcing-peloton/).
    {{% /mediaListBox %}}
    {{% mediaListBox img-src="/images/logos/gokit.png" img-alt="Go Kit Logo" title="" align=top %}}
[Go kit](https://gokit.io) is an open source toolkit for creating microservices. Described as "fill[ing] in the gaps left by the otherwise excellent standard library, giving your team the confidence to adopt Go throughout your stack."
    {{% /mediaListBox %}}
{{% /mediaList %}}

## How to use Go for Cloud Computing

Go's strengths shine when it comes to building services. Its speed and built-in support for concurrency results in fast and efficient services, while static typing, robust tooling, and emphasis on simplicity and readability help build reliable and maintainable code.
 
Go has a strong ecosystem supporting service development. The [standard library](https://golang.org/pkg/) includes packages for common needs like HTTP servers and clients, JSON/XML parsing, SQL databases, and a range of security/encryption functionality, while the Go runtime includes tools for [race detection](https://golang.org/doc/articles/race_detector.html), [benchmarking](https://golang.org/pkg/testing/#hdr-Benchmarks)/profiling, code generation, and static code analysis.

Two popular Go tools for cloud computing include [Docker](https://docker.com) and [Kubernetes](https://kubernetes.io):

**Docker is a platform-as-a-service that delivers software in containers.** Containers bundle software, libraries, and config files, are hosted by a [Docker Engine](https://www.docker.com/), and are run by a single operating-system kernel (utilizing less system resources than virtual machines).

Cloud developers use Docker to manage their Go code and support multiple platforms, as Docker supports the development workflow and deployment process. “Docker can help us to compile Go code in a clean, isolated environment,” [writes Jérôme Petazzoni](https://www.docker.com/blog/docker-golang/), founder of Tiny Shell Script. “To use different versions of the Go toolchain; and to cross-compile between different operating systems and platforms.”

**Kubernetes is an open-source container-orchestration system, written in Go, for automating web app deployment.** Web apps are often built using containers (as noted above) packaged with their dependencies and configurations. Kubernetes helps deploying and managing those containers at scale. Cloud programmers use Kubernetes to build, deliver, and scale containerized apps quickly—managing the growing complexity via APIs that controls how the containers will run.

The major Cloud providers ([GCP](https://cloud.google.com/go/home), [AWS](https://aws.amazon.com/sdk-for-go/), [Azure](https://docs.microsoft.com/en-us/azure/go/)) have Go APIs for their services, and popular open source libraries provide support for API tooling ([Swagger](https://github.com/go-swagger/go-swagger)), transport ([protocol buffers](https://github.com/golang/protobuf), [gRPC](https://grpc.io/docs/quickstart/go/)), monitoring ([OpenCensus](https://godoc.org/go.opencensus.io)), Object-Relational Mapping ([gORM](https://gorm.io/)), and authentication ([JWT](https://github.com/dgrijalva/jwt-go)). The open source community has also provided several service frameworks, including [Go Kit](https://gokit.io/[), [Go Micro](https://micro.mu/docs/go-micro.html), and [Gizmo](https://github.com/nytimes/gizmo), which can be a great way to get started quickly.

## Key solutions

### Go books on web development 

{{% books %}}

{{< headerWithLink header="Service frameworks" search="service framework" level=3 >}} 

*   {{< pkg "https://github.com/go-kit/kit" go-kit >}}, a programming toolkit for building microservices (or elegant monoliths) in Go
*   {{< pkg "https://github.com/micro/go-micro" go-micro >}}, a framework for microservice development
*   {{< pkg "https://github.com/nytimes/gizmo" >}}, a microservice framework from The New York Times
*   {{< pkg "go.uber.org/yarpc" >}}, a microservice framework from Uber 

{{< headerWithLink header="Cloud client Libraries" search="cloud client OR cloud SDK" level=3 >}} 

*   {{< pkg "cloud.google.com/go" "cloud.google.com/go">}}, Google Cloud Client Libraries for Go 
*   {{< pkg "github.com/aws/aws-sdk-go/aws/client" "aws/client">}}, AWS SDK for Go 
*   {{< pkg "github.com/aliyun/alibaba-cloud-sdk-go/sdk" >}}, Alibaba Cloud SDK for Go 
*   {{< pkg "github.com/Azure/azure-sdk-for-go" >}}, Azure SDK for Go
*   {{< pkg "github.com/IBM-Cloud/bluemix-go" >}}, IBM Cloud SDK for Go
*   {{< pkg "github.com/heroku/heroku-go/v5" >}}, Heroku client library 


{{< headerWithLink header="Libraries" search="REST OR API OR web OR cloud" level=3 >}} 

*   {{< pkg "https://github.com/go-swagger/go-swagger" go-swagger >}}, a simple and powerful representation of RESTful APIs
*   {{< pkg "github.com/emicklei/go-restful" >}}, create REST-style services without magic 
*   {{< pkg "https://github.com/golang/protobuf" >}}, Go support for Google's protocol buffer
*   {{< pkg "https://github.com/grpc/grpc-go" grpc-go >}}, a high performance, open source, general RPC framework 
*   {{< pkg "https://github.com/census-instrumentation/opencensus-go/blob/master/opencensus.go" >}}, a set of libraries for collecting application metrics and distributed traces
*   {{< pkg "https://github.com/jinzhu/gorm" >}}, an ORM library for Go
*   {{< pkg "https://github.com/dgrijalva/jwt-go" >}}, a Go implementation of json web tokens 
*   {{< pkg "github.com/lileio/lile" >}}, quickly create RPC base services

### Other resources

*   [Microservices in Golang](https://ewanvalentine.io/microservices-in-golang-part-1/), a walkthrough of microservices with Go-micro
*   [Awesome Go](https://awesome-go.com/), a curated list of Go frameworks, libraries, and software
