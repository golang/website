---
title: "American Express Uses Go for Payments & Rewards"
company: American Express
logoSrc: american-express.svg
logoSrcDark: american-express.svg
heroImgSrc: go_amex_case_study_logo.png
carouselImgSrc: go_amex_case_study.png
date: 2019-12-19
series: Case Studies
quote: Go provides American Express with the speed and scalability it needs for both its payment and rewards networks.
---

{{pullquote `
  author: Glen Balliet
  title: Engineering Director of loyalty platforms
  company: American Express
  quote: |
    What makes Go different from other programming languages is cognitive load. You can do more with less code, which makes it easier to reason about and understand the code that you do end up writing.

    The majority of Go code ends up looking quite similar, so, even if you’re working with a completely new codebase, you can get up and running pretty quickly.
`}}

## Go Improves Microservices and Speeds Productivity

Founded in 1850, American Express is a globally integrated payments company offering charge and credit card products, merchant acquisition and processing services, network services, and travel-related services.

American Express’ payment processing systems have been developed over its long history and have been updated across multiple architectural evolutions. Foremost in any update, payment processing needs to be fast, especially at very large transaction volumes, with resilience built across systems that must all be compliant with security and regulatory standards. With Go, American Express gains the speed and scalability it needs for both its payment and rewards networks.

### Modernizing American Express systems

American Express understands that the programming language landscape is changing drastically. The company's existing systems were purpose-built for high concurrency and low latency, but knowing that those systems would be re-platformed in the near future. The payments platform team decided to take the time to identify what languages were ideal for American Express's evolving needs.

The payments and rewards platform teams at American Express were among the first to start evaluating Go. These teams were focused on microservices, transaction routing, and load-balancing use cases, and they needed to modernize their architecture. Many American Express developers were familiar with the language’s capabilities and wanted to pilot Go for their high concurrency and low latency applications (such as custom transactional load balancers). With this goal in mind, the teams began lobbying senior leadership to deploy Go on the American Express payment platform.

"We wanted to find the optimal language for writing fast and efficient applications for payment processing," says Benjamin Cane, vice president and principal engineer at American Express. "To do so, we started an internal programming language showdown with the goal of seeing which language best fit our design and performance needs."

### Comparing languages

For their assessment, Cane's team chose to build a microservice in four different programming languages. They then compared the four languages for speed/performance, tooling, testing, and ease of development.

For the service, they decided on an ISO8583 to JSON converter. ISO8583 is an international standard for financial transactions, and it’s commonly used within American Express's payment network. For the programming languages, they chose to compare C++, Go, Java and Node.js. With the exception of Go, all of these languages were already in use within American Express.

From a speed perspective, Go achieved the second-best performance at 140,000 requests per second. Go showed that it excels when used for backend microservices.

While Go may not have been the fastest language tested, its powerful tooling helped bolster its overall results. Go's built-in testing framework, profiling capabilities, and benchmarking tools impressed the team. "It is easy to write effective tests in Go," says Cane. "The benchmarking and profiling features make it simple to tune our application. Coupled with its fast build times, Go makes it easy to write well-tested and optimized code."

Ultimately, Go was selected by the team as the preferred language for building high-performance microservices. The tooling, testing frameworks, performance, and language simplicity were all key contributors.

### Go for infrastructure

"Many of our services are running in Docker containers within our Kubernetes-based internal cloud platform" says Cane. Kubernetes is an open-source container-orchestration system written in Go. It provides clusters of hosts to run container based workloads, most notably Docker containers. Docker is a software product, also written in Go, that uses operating system level virtualization to provide portable software runtimes called containers.

American Express also collects application metrics via Prometheus, an open-source monitoring and alerting toolkit written in Go. Prometheus collects and aggregates real-time events and metrics for monitoring and alerts.

This triumvirate of Go solutions—Kubernetes, Docker, and Prometheus—has helped modernize American Express infrastructure.

### Improving performance with Go

Today, scores of developers are programming with Go at American Express, with most working on platforms designed for high availability and performance.

"Tooling has always been a critical area of need for our legacy code base," says Cane. "We have found that Go has excellent tooling, plus built-in testing, benchmarking, and profiling frameworks. It is easy to write efficient and resilient applications."

{{backgroundquote `
  author: Benjamin Cane
  title: Vice President and Principal Engineer
  company: American Express
  quote: |
    After working on Go, most of our developers don't want to go back to other languages.
`}}

American Express is just beginning to see the benefits of Go. For example, Go was designed from the ground up with concurrency in mind – using lightweight “goroutines” rather than heavier-weight operating system threads – making it practical to create hundreds of thousands of goroutines in the same address space. Using goroutines, American Express has seen improved performance numbers in its real-time transaction processing.

Go’s garbage collection is also a major improvement over other languages, both in terms of performance and ease of development. “We saw far better results of garbage collection in Go than we did in other languages, and garbage collection for real time transaction processing is a big deal.” says Cane. “Tuning garbage collection in other languages can be very complicated. With Go you don’t tune anything.”

To learn more, read ["Choosing Go at American Express"](https://americanexpress.io/choosing-go/) which goes into more depth about American Express's Go adoption.

### Getting your enterprise started with Go

Just as American Express is using Go to modernize its payment and rewards networks, dozens of other large enterprises are adopting Go as well.

There are over one million developers using Go worldwide—spanning banking and commerce, gaming and media, technology, and other industries, at enterprises as diverse as [PayPal](/solutions/paypal), [Mercado Libre](/solutions/mercadolibre), Capital One, Dropbox, IBM, Mercado Libre, Monzo, New York Times, Salesforce, Square, Target, Twitch, Uber, and of course Google.

To learn more about how Go can help your enterprise build reliable, scalable software as it does at American Express, visit [go.dev](/) today.
