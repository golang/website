---
title: "Go for Cloud Computing"
description: "Go Helps Enterprises Build and Scale Cloud Computing Systems"
date: 2019-10-03T16:26:31-04:00
series: Use Cases
---

### _Go Helps Enterprises Build and Scale Cloud Computing Systems_

## **Why Use Go for Cloud Computing**

As applications and processing move to the cloud, concurrency becomes a very big issue. Cloud computing systems, by
their very nature, share and scale resources. Coordinating access to shared resources is an issue that impacts every
application processing in the cloud, and requires programming languages “[explicitly geared to develop highly reliable
concurrent applications](https://tomassetti.me/best-programming-languages/).”

Go was created to address exactly these concurrency needs for scaled applications, microservices, and cloud development.
In fact, over 75 percent of projects in the Cloud Native Computing Foundation ([CNCF](https://www.cncf.io/projects/))
are written in Go.

## **Who Uses Go for Cloud Computing**

### **American Express**

American Express uses Go to improve microservices and speed cloud development productivity. Go began at American Express
with the efforts of the payment processing platform team, focused on microservices, transaction routing, and load
balancing. By leveraging goroutines, the payment processing team has seen improved performance numbers in its real-time
transaction processing, and validated Go’s garbage collection as a huge improvement over other languages.

Today, American Express also uses Go via:

*   Docker—a SaaS product, written in Go, that uses operating system level virtualization to develop and deliver
*   software in containers hosted on a Docker Engine Kubernetes—an open-source container-orchestration system, written
*   in Go, that follows a primary/replica architecture across clusters of American Express hosts Prometheus—an
*   open-source software application written in Go used for real-time event monitoring and alerting

This triumvirate of Go solutions has helped modernize American Express’s infrastructure and opened the door for Go as a
key player in the American Express payment ecosystem.

### **AT&T**

Within AT&T's DirectTV division, a microservices development team oversees VUD monitoring and analytics for the video
ingestion pipeline as it comes in from content providers. The team builds small microservices in Go as monitoring
points, checking when video content goes from one state to another throughout their system. Being able to re-write old
microservices in a cloud-friendly language like Go delivered a tremendous development cost-savings to AT&T.  The team
also developed a Go SDK to support future Go development on AT&T's platform.

The DirecTV division also hosts Kubernetes themselves and are looking at Knative (a Kubernetes-based platform to build,
deploy, and manage modern serverless workloads).

### **Dropbox**

Dropbox decided to migrate its performance-critical backends from Python to Go to leverage better concurrency support
and faster execution speed.  Go delivers better performance for Dropbox engineering teams, making them more productive
with a standard library, debugging tools that work, and proving it is easier to both write and consume services in Go.

Today, Dropbox has over 500 million users and most of the company's infrastructure is written in Go—over 1.3 million
lines of Go and every Dropbox engineer hired goes through Go training during onboarding.

### **MercadoLibre**

MercadoLibre uses Go to scale its eCommerce empire.  Go produces efficient code that readily scales as MercadoLibre’s
online commerce grows and supports the company as a boon for developers—improving their productivity while streamlining
and expanding MercadoLibre services.

Go started with the core APIs team, which builds and maintains the largest APIs at the center of MercadoLibre’s
microservices solutions. The API team converted their architecture to Go to great performance benefits, and one large Go
program is now able to run 100,000 requests per machine with just 24 megabytes of memory.

MercadoLibre leverages Go’s expressive and clean syntax to make it easier for developers to write programs that get the
most out of the company’s multicore and networked machines. And while speed in development yields cost efficiency for
the company, developers individually benefit from the swift learning curve Go delivers. Not only are MercadoLibre's very
experienced developers able to build highly critical applications very, very quickly with Go, but even new programmers
have been able to produce significant solutions.

## **How to Use Go for Cloud Computing**

Go's strengths shine when it comes to building services. Its speed and built-in support for concurrency result in fast
and efficient services, while static typing, robust tooling, and emphasis on simplicity and readability help build
reliable and maintainable code.

Go has a strong ecosystem supporting service development. The[ standard library](https://golang.org/pkg/) includes
packages for common needs like HTTP servers and clients, JSON/XML parsing, SQL databases, and a range of
security/encryption functionality, while the Go runtime includes tools for[ race
detection](https://golang.org/doc/articles/race_detector.html),[
benchmarking](https://golang.org/pkg/testing/#hdr-Benchmarks)/profiling, code generation, and static code analysis.

The major Cloud providers ([GCP](https://cloud.google.com/go/home),[ AWS](https://aws.amazon.com/sdk-for-go/),[
Azure](https://docs.microsoft.com/en-us/azure/go/)) have Go APIs for their services, and popular open source libraries
provide support for API tooling ([Swagger](https://github.com/go-swagger/go-swagger)), transport ([protocol
buffers](https://github.com/golang/protobuf),[ gRPC](https://grpc.io/docs/quickstart/go/)), monitoring
([OpenCensus](https://godoc.org/go.opencensus.io)), Object-Relational Mapping ([gORM](https://gorm.io/)), and
authentication ([JWT](https://github.com/dgrijalva/jwt-go)). The open source community has also provided several service
frameworks, including[ Go Kit](https://gokit.io/),[ Go Micro](https://micro.mu/docs/go-micro.html), and[
Gizmo](https://github.com/nytimes/gizmo), which can be a great way to get started quickly.

Service developers often make a tradeoff between development cycle time and server performance. Go's fast build times
make iterative development possible, while still yielding the benefits of fast compiled code. Plus, Go servers tend to
have lower memory and CPU utilization, making them cheaper to run in pay-as-you-go deployments

Overall, Go's mission to make it easy to write simple, reliable, and efficient software means it is a great choice for
developing services.

## **Go Solutions to Legacy Challenges**

Historically, challenges facing cloud computing systems have included the need for highly concurrent and distributed
processing, multi-nodes and multi-cores, the lack of shared memory, and the very major bottleneck of single-threaded
applications.

Further, cloud engineering teams want to be able to develop cloud applications locally, they want to be able to develop
cross-cloud applications with a simple idiomatic interface, and they want to be able to run both on-premises and on the
cloud.

"Go was designed to be scalable to large systems and usable without an IDE, but also productive and being especially
good at networking and concurrency," notes Federico Tomassetti in his[
blog](https://tomassetti.me/best-programming-languages/). "Other than a well-thought design, it has some specific
features for concurrency like a type of light-weight processes called goroutines."

Goroutines do not have names; they are just anonymous workers. They expose no unique identifier, name, or data structure
to the programmer. Some people are surprised by this, expecting the go statement to return some item that can be used to
access and control the goroutine later.

The fundamental reason goroutines are anonymous is so that the full Go language is available when programming concurrent
code. By contrast, the usage patterns that develop when threads and goroutines are named can restrict what a library
using them can do.

For example, once one names a goroutine and constructs a model around it, it becomes special, and one would be tempted
to associate _all_ computation with that goroutine—ignoring the possibility of using multiple, possibly shared
goroutines for the processing. If the net/http package associated per-request state with a goroutine, clients would be
unable to use more goroutines when serving a request.

Go solves the problems of modern cloud development, delivering a standard idiomatic APIs designed around user needs,
plus out-of-the-box support for multiple cloud environments (including on-premises), the ability to write and test
locally (run in production), and open cloud development . . . granting development teams the power to choose and the
power to move.

## **Resources for Learning More**

*   Go.dev (when it’s live)
*   http://linuxgizmos.com/latest-edgex-iot-middleware-release-gets-smaller-faster-and-more-secure/ or similar?
*   https://ewanvalentine.io/microservices-in-golang-part-1/ (through part 10) is a pretty nice walk-through of Go with
*   microservices (using go-micro). Do we want to cite? Similarly,[ https://awesome-go.com/](https://awesome-go.com/)
