---
title: "Go for Cloud & Network Services"
linkTitle: "Cloud & Network Services"
description: "With a strong ecosystem of tools and APIs on major cloud providers, it is easier than ever to build services with Go."
date: 2019-10-04T15:26:31-04:00
series: Use Cases
icon:
  file: cloud-green.svg
  alt: cloud icon
iconDark:
  file: cloud-white.svg
  alt: cloud icon
---

## Overview {#overview .sectionHeading}

<div class="UseCase-halfColumn">
    <h3>Go helps enterprises build and scale cloud computing systems</h3>
    <p>As applications and processing move to the cloud, concurrency becomes a very big issue. Cloud computing systems, by their very nature, share and scale resources. Coordinating access to shared resources is an issue that impacts every application processing in the cloud, and requires programming languages “explicitly geared to develop highly reliable concurrent applications.”</p>
  </div>

{{quote `
  author: Ruchi Malik
  title: developer at Choozle
  link: https://builtin.com/software-engineering-perspectives/golang-advantages
  quote: |
    Go makes it very easy to scale as a company. This is very important because, as our engineering team grows, each service can be managed by a different unit.
`}}

## Key Benefits {#key-benefits .sectionHeading}

### Address tradeoff between development cycle time and server performance

Go was created to address exactly these concurrency needs for scaled applications, microservices, and cloud development. In fact, over 75 percent of projects in the Cloud Native Computing Foundation are written in Go.

Go helps reduce the need to make this tradeoff, with its fast build times that enable iterative development, lower memory and CPU utilization. Servers built with Go experience instant start up times and are cheaper to run in pay-as-you-go and serverless deployments.

### Address challenges with the modern cloud, delivering standard idiomatic APIs

Go addresses many challenges developers face with the modern cloud, delivering standard idiomatic APIs, and built in concurrency to take advantage of multicore processors. Go’s low-latency and “no knob” tuning make Go a great balance between performance and productivity - granting engineering teams the power to choose and the power to move.

## Use Case {#use-case .sectionHeading}

### Use Go for Cloud Computing

Go's strengths shine when it comes to building services. Its speed and built-in support for concurrency results in fast and efficient services, while static typing, robust tooling, and emphasis on simplicity and readability help build reliable and maintainable code.

Go has a strong ecosystem supporting service development. The [standard library](/pkg/) includes packages for common needs like HTTP servers and clients, JSON/XML parsing, SQL databases, and a range of security/encryption functionality, while the Go runtime includes tools for [race detection](/doc/articles/race_detector.html), [benchmarking](/pkg/testing/#hdr-Benchmarks)/profiling, code generation, and static code analysis.

The major Cloud providers ([GCP](https://cloud.google.com/go/home), [AWS](https://aws.amazon.com/sdk-for-go/), [Azure](https://docs.microsoft.com/en-us/azure/go/)) have Go APIs for their services, and popular open source libraries provide support for API tooling ([Swagger](https://github.com/go-swagger/go-swagger)), transport ([protocol buffers](https://github.com/golang/protobuf), [gRPC](https://grpc.io/docs/quickstart/go/)), monitoring ([OpenCensus](https://godoc.org/go.opencensus.io)), Object-Relational Mapping ([gORM](https://gorm.io/)), and authentication ([JWT](https://github.com/dgrijalva/jwt-go)). The open source community has also provided several service frameworks, including [Go Kit](https://gokit.io/), [Go Micro](https://micro.mu/docs/go-micro.html), and [Gizmo](https://github.com/nytimes/gizmo), which can be a great way to get started quickly.

### Go tools for Cloud Computing

{{toolsblurbs `
  - title: Docker
    url: https://www.docker.com/
    iconSrc: /images/logos/docker.svg
    paragraphs:
      - Docker is a platform-as-a-service that delivers software in containers. Containers bundle software, libraries, and config files, are hosted by a Docker Engine, and are run by a single operating-system kernel (utilizing less system resources than virtual machines).
      - Cloud developers use Docker to manage their Go code and support multiple platforms, as Docker supports the development workflow and deployment process.
  - title: Kubernetes
    url: https://kubernetes.io/
    iconSrc: /images/logos/kubernetes.svg
    paragraphs:
      - Kubernetes is an open-source container-orchestration system, written in Go, for automating web app deployment. Web apps are often built using containers (as noted above) packaged with their dependencies and configurations. Kubernetes helps deploying and managing those containers at scale. Cloud programmers use Kubernetes to build, deliver, and scale containerized apps quickly—managing the growing complexity via APIs that controls how the containers will run.
`}}

{{projects `
  - company: Google
    url: http://cloud.google.com/go
    logoSrc: google-cloud.svg
    logoSrcDark: google-cloud.svg
    desc: Google Cloud uses Go across its ecosystem of products and tools, including Kubernetes, gVisor, Knative, Istio, and Anthos. Go is fully supported on Google Cloud across all APIs and runtimes.
    ctas:
      - text: Go on Google Cloud Platform
        url: http://cloud.google.com/go
  - company: Capital One
    url: https://www.capitalone.com/
    logoSrc: capitalone_light.svg
    logoSrcDark: capitalone_dark.svg
    desc: Capital One uses Go to power the Credit Offers API, a critical service. The engineering team is also building their serverless architecture with Go, citing Go’s speed and simplicity, and mentioning that “[they] didn’t want to go serverless without Go.”
    ctas:
      - text: Credit Offers API
        url: https://medium.com/capital-one-tech/a-serverless-and-go-journey-credit-offers-api-74ef1f9fde7f
  - company: Dropbox
    url: https://www.dropbox.com/
    logoSrc: dropbox.svg
    logoSrcDark: dropbox.svg
    desc: Dropbox was built on Python, but in 2013 decided to migrate their performance-critical backends to Go. Today, most of the company’s infrastructure is written in Go.
    ctas:
      - text: Dropbox libraries
        url: https://blogs.dropbox.com/tech/2014/07/open-sourcing-our-go-libraries/
  - company: Mercado Libre
    url: https://www.mercadolibre.com.ar/
    logoSrc: mercadolibre_light.svg
    logoSrcDark: mercadolibre_dark.svg
    desc: MercadoLibre uses Go to scale its eCommerce platform. Go produces efficient code that readily scales as MercadoLibre’s online commerce grows. Go improves their productivity while streamlining and expanding MercadoLibre services.
    ctas:
      - text: MercadoLibre & Go
        url: http://go.dev/solutions/mercadolibre
  - company: The New York Times
    url: https://www.nytimes.com/
    logoSrc: the-new-york-times-icon.svg
    logoSrcDark: the-new-york-times-icon.svg
    desc: The New York Times adopted Go “to build better back-end services”. As the usage of Go expanded with in the company they felt the need to create a toolkit to “to help developers quickly configure and build microservice APIs and pubsub daemons”, which they have open sourced.
    ctas:
      - text: NYTimes - Gizmo
        url: https://open.nytimes.com/introducing-gizmo-aa7ea463b208
      - text: Gizmo GitHub
        url: https://github.com/nytimes/gizmo
  - company: Twitch
    url: https://www.twitch.tv/
    logoSrc: twitch.svg
    logoSrcDark: twitch.svg
    desc: Twitch uses Go to power many of its busiest systems that serve live video and chat to millions of users.
    ctas:
      - text: Go’s march to low-latency GC
        url: https://blog.twitch.tv/en/2016/07/05/gos-march-to-low-latency-gc-a6fa96f06eb7/
  - company: Uber
    url: https://www.uber.com/
    logoSrc: uber_light.svg
    logoSrcDark: uber_dark.svg
    desc: Uber uses Go to power several of its critical services that impact the experience of millions of drivers and passengers around the world. From their real-time analytics engine, AresDB, to their microservice for Geo-querying, Geofence, and their resource scheduler, Peloton.
    ctas:
      - text: AresDB
        url: https://eng.uber.com/aresdb/
      - text: Geofence
        url: https://eng.uber.com/go-geofence/
      - text: Peloton
        url:  https://eng.uber.com/open-sourcing-peloton/
`}}

## Get Started {#get-started .sectionHeading}

### Go books for cloud computing

{{books `
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
    url: https://www.amazon.com/Mastering-Web-Services-Nathan-Kozyra/dp/178398130X
    thumbnail: /images/books/mastering-go-web-services.jpg
`}}

{{libraries `
  - title: Web frameworks
    viewMoreUrl: https://pkg.go.dev/search?q=web+framework
    items:
      - text: Buffalo
        url: https://gobuffalo.io/en/
        desc: A framework for rapid web development in Go, curating Go and JS libraries together.
      - text: Echo
        url: https://echo.labstack.com/
        desc: A high performance, extensible, and minimalist Go web framework
      - text: Flamingo
        url: https://www.flamingo.me/
        desc: A fast open-source framework based on Go with clean and scalable architecture
      - text: Gin
        url: https://gin-gonic.com/
        desc: A web framework written in Go, with a martini-like API.
      - text: Gorilla
        url: http://www.gorillatoolkit.org/
        desc: A web toolkit for the Go programming language.
  - title: Routers
    viewMoreUrl: https://pkg.go.dev/search?q=http%20router
    items:
      - text: julienschmidt/httprouter
        url: https://pkg.go.dev/github.com/julienschmidt/httprouter?tab=overview
        desc: A lightweight high performance HTTP request router
      - text: gorilla/mux
        url: https://pkg.go.dev/github.com/gorilla/mux?tab=overview
        desc: A powerful HTTP router and URL matcher for building Go web servers with 🦍
      - text: Chi
        url: https://pkg.go.dev/github.com/go-chi/chi?tab=overview
        desc: A lightweight, idiomatic and composable router for building Go HTTP services.
      - text: net/http
        url: https://pkg.go.dev/net/http
        desc: A standard library HTTP package
  - title: Template Engines
    viewMoreUrl: https://pkg.go.dev/search?q=templates
    items:
      - text: html/template
        url: https://pkg.go.dev/html/template
        desc: A standard library HTML template engine
      - text: flosch/pongo2
        url: https://pkg.go.dev/github.com/flosch/pongo2?tab=overview
        desc: A Django-syntax like templating-language
  - title: Databases & Drivers
    viewMoreUrl: https://pkg.go.dev/search?q=database%20OR%20sql
    items:
      - text: database/sql
        url: https://pkg.go.dev/database/sql
        desc: A standard library interface with driver support for MySQL, Postgres, Oracle, MS SQL, BigQuery and most SQL databases
      - text: mongo-driver/mongo
        url: https://pkg.go.dev/go.mongodb.org/mongo-driver/mongo?tab=overview
        desc: The MongoDB supported driver for Go
      - text: elastic/go-elasticsearch
        url: https://pkg.go.dev/github.com/elastic/go-elasticsearch/v8?tab=overview
        desc: An Elasticsearch client for Go
      - text: GORM
        url: https://gorm.io/
        desc: An ORM library for Go
      - text: Bleve
        url: http://blevesearch.com/
        desc: Full-text search and indexing for Go
      - text: CockroachDB
        url: https://www.cockroachlabs.com/
        desc: An evolution of the database—architected for the cloud to deliver resilient, consistent, distributed SQL at scale
  - title: Web Libraries
    viewMoreUrl: https://pkg.go.dev/search?q=web
    items:
      - text: markbates/goth
        url: https://pkg.go.dev/github.com/markbates/goth?tab=overview
        desc: Authentication for web apps
      - text: jinzhu/gorm
        url: https://pkg.go.dev/github.com/jinzhu/gorm?tab=overview
        desc: An ORM library for Go
      - text: dgrijalva/jwt-go
        url: https://pkg.go.dev/github.com/dgrijalva/jwt-go?tab=overview
        desc: A Go implementation of json web tokens
  - title: Other Projects
    items:
      - text: gopherjs
        url: https://pkg.go.dev/github.com/gopherjs/gopherjs?tab=overview
        desc: A compiler from Go to JavaScript allowing developers to write front-end code in Go which will run in all browsers.
`}}
