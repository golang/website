---
title: "Go for Cloud & Network Services"
linkTitle: "Cloud & Network Services"
description: "With a strong ecosystem of tools and APIs on major cloud providers, it is easier than ever to build services with Go."
date: 2019-10-04T15:26:31-04:00
series: Use Cases
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

### _Go helps enterprises build and scale Cloud Computing systems_

## **Why use Go for Cloud Computing**

As applications and processing move to the cloud, concurrency becomes a very big issue. Cloud computing systems, by their very nature, share and scale resources. Coordinating access to shared resources is an issue that impacts every application processing in the cloud, and requires programming languages “[explicitly geared to develop highly reliable concurrent applications](https://tomassetti.me/best-programming-languages/).”
 
Go was created to address exactly these concurrency needs for scaled applications, microservices, and cloud development. In fact, over 75 percent of projects in the Cloud Native Computing Foundation ([CNCF](https://www.cncf.io/projects/)) are written in Go.

“Go makes it very easy to scale as a company,” [says developer Ruchi Malik](https://builtin.com/software-engineering-perspectives/golang-advantages) at Choozle. “This is very important because, as our engineering team grows, each service can be managed by a different unit.”

## **Who uses Go for Cloud Computing**

### **American Express**

American Express uses Go to improve microservices and speed cloud development productivity. Go began at American Express with the efforts of the payment processing platform team, focused on microservices, transaction routing, and load balancing. By leveraging goroutines, the payment processing team has seen improved performance numbers in its real-time transaction processing, and validated Go’s garbage collection as a huge improvement over other languages.

Today, American Express also uses Go via:

*   [Docker](https://www.docker.com/)—a SaaS product, written in Go, that uses operating system level virtualization to develop and deliver software in containers hosted on a Docker Engine
*   [Kubernetes](https://kubernetes.io/)—an open-source container-orchestration system, written in Go, that follows a primary/replica architecture across clusters of
*   [Prometheus](https://prometheus.io/)—an open-source software application written in Go used for real-time event monitoring and alerting

This triumvirate of Go solutions has helped modernize American Express’s infrastructure and opened the door for Go as a key player in the American Express payment ecosystem.

### **AT&T**

Within AT&T's DirectTV division, a microservices development team oversees VUD monitoring and analytics for the video ingestion pipeline as it comes in from content providers. The team builds small microservices in Go as monitoring points, checking when video content goes from one state to another throughout their system. Being able to rewrite old microservices in a cloud-friendly language like Go delivered a tremendous development cost-savings to AT&T.  The team also developed a Go SDK to support future Go development on AT&T's platform.
 
The DirecTV division also hosts Kubernetes themselves and are looking at Knative (a Kubernetes-based platform to build, deploy, and manage modern serverless workloads).

### **Dropbox**

Dropbox was built on Python, [but in 2013 decided to migrate](https://blogs.dropbox.com/tech/2014/07/open-sourcing-our-go-libraries/) their performance-critical backends to Go. Dropbox developers wanted better concurrency support and faster execution speeds, and were willing to write around 200,000 lines of new Go code. Dropbox has since built many open-source libraries to support large-scale production systems, including libraries for caching, errors, database/sqlbuilder, and hash2. 
 
Today, Dropbox has over 500 million users and most of the company's infrastructure is written in Go—over 1.3 million lines of Go and every Dropbox engineer hired goes through Go training during onboarding. Dropbox libraries can be found at [Dropbox's Go github](https://github.com/dropbox/godropbox).

### **MercadoLibre**

MercadoLibre uses Go to scale its eCommerce empire.  Go produces efficient code that readily scales as MercadoLibre’s online commerce grows and supports the company as a boon for developers—improving their productivity while streamlining and expanding MercadoLibre services.
 
Go started with the core APIs team, which builds and maintains the largest APIs at the center of MercadoLibre’s microservices solutions. With Go, MercadoLibre’s build times are three times (3x) faster and their test suite runs an amazing 24 times faster. This means the company’s developers can make a change, then build and test that change much faster than they could before. And dropping MercadoLibre’s test suite runtimes from 90-seconds to just 3-seconds with Go was a huge boon for its developers—allowing them to keep focus (and context) while the much faster tests complete.
 
MercadoLibre leverages Go’s expressive and clean syntax to make it easier for developers to write programs that get the most out of the company’s multicore and networked machines. And while speed in development yields cost efficiency for the company, developers individually benefit from the swift learning curve Go delivers. Not only are MercadoLibre's very experienced developers able to build highly critical applications very, very quickly with Go, but even new programmers have been able to produce significant solutions. 

## **How to use Go for Cloud Computing**

Go's strengths shine when it comes to building services. Its speed and built-in support for concurrency results in fast and efficient services, while static typing, robust tooling, and emphasis on simplicity and readability help build reliable and maintainable code.
 
Go has a strong ecosystem supporting service development. The [standard library](https://golang.org/pkg/) includes packages for common needs like HTTP servers and clients, JSON/XML parsing, SQL databases, and a range of security/encryption functionality, while the Go runtime includes tools for [race detection](https://golang.org/doc/articles/race_detector.html), [benchmarking](https://golang.org/pkg/testing/#hdr-Benchmarks)/profiling, code generation, and static code analysis.

Two popular Go tools for cloud computing include Docker and Kubernetes:

**Docker is a platform-as-a-service that delivers software in containers.** Containers bundle software, libraries, and config files, are hosted by a [Docker Engine](https://www.docker.com/), and are run by a single operating-system kernel (utilizing less system resources than virtual machines).

Cloud developers use Docker to manage their Go code and support multiple platforms, as Docker supports the development workflow and deployment process. “Docker can help us to compile Go code in a clean, isolated environment,” [writes Jérôme Petazzoni](https://www.docker.com/blog/docker-golang/), founder of Tiny Shell Script. “To use different versions of the Go toolchain; and to cross-compile between different operating systems and platforms.”

**Kubernetes is an open-source container-orchestration system, written in Go, for automating web app deployment.** Web apps are often built using containers (as noted above) packaged with their dependencies and configurations. Kubernetes helps deploying and managing those containers at scale. Cloud programmers use Kubernetes to build, deliver, and scale containerized apps quickly—managing the growing complexity via APIs that controls how the containers will run.

The major Cloud providers ([GCP](https://cloud.google.com/go/home), [AWS](https://aws.amazon.com/sdk-for-go/), [Azure](https://docs.microsoft.com/en-us/azure/go/)) have Go APIs for their services, and popular open source libraries provide support for API tooling ([Swagger](https://github.com/go-swagger/go-swagger)), transport ([protocol buffers](https://github.com/golang/protobuf), [gRPC](https://grpc.io/docs/quickstart/go/)), monitoring ([OpenCensus](https://godoc.org/go.opencensus.io)), Object-Relational Mapping ([gORM](https://gorm.io/)), and authentication ([JWT](https://github.com/dgrijalva/jwt-go)). The open source community has also provided several service frameworks, including [Go Kit](https://gokit.io/[), [Go Micro](https://micro.mu/docs/go-micro.html), and [Gizmo](https://github.com/nytimes/gizmo), which can be a great way to get started quickly.
 
Service developers often make a tradeoff between development cycle time and server performance. Go's fast build times make iterative development possible, while still yielding the benefits of fast compiled code. Plus, Go servers tend to have lower memory and CPU utilization, making them cheaper to run in pay-as-you-go deployments.
 
Overall, Go's mission to make it easy to write simple, reliable, and efficient software means it is a great choice for developing services.

## **Go solutions to legacy challenges**

Historically, challenges facing cloud computing systems have included the need for highly concurrent and distributed processing, multi-nodes and multi-cores, the lack of shared memory, and the very major bottleneck of single-threaded applications.
 
Further, cloud engineering teams want to be able to develop cloud applications locally, they want to be able to develop cross-cloud applications with a simple idiomatic interface, and they want to be able to run both on-premise and on the cloud.
 
"Go was designed to be scalable to large systems and usable without an IDE, but also productive and being especially good at networking and concurrency," notes Federico Tomassetti in his [blog](https://tomassetti.me/best-programming-languages/). "Other than a well-thought design, it has some specific features for concurrency like a type of light-weight processes called goroutines."
 
Goroutines do not have names; they are just anonymous workers. They expose no unique identifier, name, or data structure to the programmer. Some people are surprised by this, expecting the go statement to return some item that can be used to access and control the goroutine later.
 
The fundamental reason goroutines are anonymous is so that the full Go language is available when programming concurrent code. By contrast, the usage patterns that develop when threads and goroutines are named can restrict what a library using them can do.
 
For example, once one names a goroutine and constructs a model around it, it becomes special, and one would be tempted to associate _all_ computation with that goroutine—ignoring the possibility of using multiple, possibly shared goroutines for the processing. If the net/http package associated per-request state with a goroutine, clients would be unable to use more goroutines when serving a request.
 
Go solves the problems of modern cloud development, delivering a standard idiomatic APIs designed around user needs, plus out-of-the-box support for multiple cloud environments (including on-premises), the ability to write and test locally (run in production), and open cloud development - granting development teams the power to choose and the power to move.

## **Resources for Learning More**

*   [EdgeX](http://linuxgizmos.com/latest-edgex-iot-middleware-release-gets-smaller-faster-and-more-secure/) - IoT Middleware
*   [Microservices in Golang](https://ewanvalentine.io/microservices-in-golang-part-1/) - walkthrough of microservices with Go-micro
*    [Awesome Go](https://awesome-go.com/) - curated list of Go frameworks, libraries, and software
