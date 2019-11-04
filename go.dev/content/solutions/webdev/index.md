---
title: "Go for Web Development"
linkTitle: "Web Development"
description: "With enhanced memory performance and support for several IDEs, Go powers fast and scalable web applications."
date: 2019-10-04T15:26:31-04:00
series: Use Cases
books:
  - title: Web Development with Go
    url: https://www.amazon.com/Web-Development-Go-Building-Scalable-ebook/dp/B01JCOC6Z6
    thumbnail: /images/books/web-development-with-go.jpg
  - title: Go Web Programming
    url: https://www.amazon.com/Web-Programming-Sau-Sheong-Chang/dp/1617292567
    thumbnail: /images/books/go-web-programming.jpg
  - title: "Web Development Cookbook: Build full-stack web applications with Go"
    url: https://www.amazon.com/Web-Development-Cookbook-full-stack-applications-ebook/dp/B077TVQ28W
    thumbnail: /images/books/go-web-development-cookbook.jpg
  - title: Building RESTful Web services with Go
    url: https://www.amazon.com/Building-RESTful-Web-services-gracefully-ebook/dp/B072QB8KL1
    thumbnail: /images/books/building-restful-web-services-with-go.jpg
  - title: Mastering Go Web Services
    url: https://www.amazon.com/Mastering-Web-Services-Nathan-Kozyra-ebook/dp/B00W5GUKL6
    thumbnail: /images/books/mastering-go-web-services.jpg
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

## Go delivers speed, security, and developer-friendly tools for Web Applications

Go is designed to enable developers to rapidly develop scalable and secure web applications. Go ships with an easy to use, secure and performant web server and includes it own web templating library. Go has excellent support for all of the latest technologies from [HTTP/2](https://pkg.go.dev/net/http), to databases like [MySQL](https://pkg.go.dev/mod/github.com/go-sql-driver/mysql), [MongoDB](https://pkg.go.dev/mod/go.mongodb.org/mongo-driver) and [ElasticSearch](https://pkg.go.dev/mod/github.com/elastic/go-elasticsearch/v7), to the latest encryption standards including [TLS 1.3](https://pkg.go.dev/crypto/tls). Go web applications run natively on [Google App Engine](https://cloud.google.com/appengine/) (for easy scaling) or on any environment, cloud, or operating system thanks to Go’s extreme portability. 

{{% pullquote author="Tigran Bayburtsyan, Co-Founder and CTO at Hexact Inc." link="https://hackernoon.com/5-reasons-why-we-switched-from-python-to-go-4414d5f42690" %}}
Go Language is the easiest language that I’ve ever seen and used... For me, Go is easier to learn than even JavaScript.
{{% /pullquote %}}

Bayburtsyan summarizes the five key reasons his company switched to Go:

1.   **Compiles into a single binary** — “Using static linking, Go actually combining all dependency libraries and modules into one single binary file based on OS type and architecture.”

2.   **Static type system** — “Type system is really important for large scale applications.”

3.   **Performance** — “Go performed better because of its concurrency model and CPU scalability. Whenever we need to process some internal request, we are doing it with separate Goroutines which are 10x cheaper in resources than Python Threads.”

4.   **No need for a web framework** — “In most of the cases you really don’t need any third-party library.”

5.   **Great IDE support and debugging** — “After rewriting all projects to Go, we got 64 percent less code than we had earlier.”

{{% gopher gopher=front align=right %}}
For enterprises, Go is preferred for providing rapid cross-platform deployment. With its goroutines, native compilation, and the URI-based package namespacing, Go code executes to a single, small binary—with zero dependencies—making it very fast.

“If you are looking for powerful tools for web programming, mobile development, microservices, and ERP systems,” [writes Andrew Smith](https://dzone.com/articles/golang-web-development-better-than-python), marketing manager at QArea. “Go web development has proved to be faster than using Python for the same kind of tasks in many use cases.”

## Featured Go web development users and projects 

{{% mediaList %}}
    {{% mediaListBox img-src="/images/logos/hugo.svg" img-alt="Hugo Logo"  img-link="https://gohugo.io" title="Hugo" align=top  %}}
[Hugo](https://gohugo.io) is a fast and modern static site generator written in Go, and designed to make website creation fun again. Hugo is one of the most popular website engines available today. Websites built with Hugo are extremely fast and secure and can be hosted anywhere without any dependencies. 
    {{% /mediaListBox %}}
    {{% mediaListBox img-src="/images/logos/caddy.svg" img-alt="Caddy Logo"  img-link="https://caddyserver.com" title="Caddy" align=top  %}}
[Caddy 2](https://caddyserver.com) is a powerful, enterprise-ready, open source web server with automatic HTTPS written in Go. Written in Go, Caddy offers greater memory safety than servers written in C. A hardened TLS stack powered by the Go standard library serves a significant portion of all Internet traffic. 
    {{% /mediaListBox %}}
    {{% mediaListBox img-src="/images/logos/lets-encrypt-icon.svg" img-alt="Let's Encrypt Logo"  img-link="https://letsencrypt.org" title="Let's Encrypt" align=top  %}}
Let’s Encrypt is a free, automated, and open Certificate Authority run by the non-profit Internet Security Research Group (ISRG). The CA's [website](https://letsencrypt.org/) is powered by Hugo—a static site generator written in Go. Plus, for web services written in Go (using the built-in Go web server), Let's Encrypt allows web developers to serve pages over HTTPS with quick and minimal setup.
    {{% /mediaListBox %}}
    {{% mediaListBox img-src="Cloudflare Logo" img-alt="" title="hold" align="top" %}}

    {{% /mediaListBox %}}
    {{% mediaListBox img-src="/images/logos/govuk.svg" img-alt="Gov.UK Logo" img-link="https://gov.uk" title="Gov.UK" align=top %}}
GOV.UK made the decision to develop in Go back in 2013. The simplicity and safety of the Go language were a good fit for a core component of their HTTP infrastructure, and some brief experiments with the excellent net/http package convinced web developers they were on the right track. “In particular, Go’s concurrency model makes it absurdly easy to build performant I/O-bound applications,” note Nick Stenning and James Stewart in their [Technology in Government blog](https://technology.blog.gov.uk/2013/12/05/building-a-new-router-for-gov-uk/).
    {{% /mediaListBox %}}
    {{% mediaListBox img-src="/images/logos/uber-app-icon.svg" img-alt="Uber Logo" title="Uber" align=top %}}
[In early 2015](https://eng.uber.com/go-geofence/), Uber leveraged Go to build a microservice for their geofence lookups. At Uber, a geofence refers to a human-defined geographic area on the Earth’s surface, and Uber uses geofences to show its users which products are available at a given location, to define areas with specific requirements such as airports, and to implement dynamic pricing in areas where many people are requesting rides at the same time. By 2016, this service was Uber’s highest queries per second (QPS) service out of the hundreds they run in production—thanks to Go helping them build and scale it so fast.
    {{% /mediaListBox %}}

{{% /mediaList %}}

## Key Solutions

### Go books on web development 

{{% books %}}

{{< headerWithLink header="Web Frameworks" link="https://pkg.go.dev/search?q=web+framework" level=3 >}} 

*   [Buffalo](https://gobuffalo.io/en/), for rapid web development in Go. While Buffalo can be considered as a framework, it's mostly an ecosystem of Go and Javascript libraries curated to fit together.
*   [Echo](https://echo.labstack.com), a high performance, extensible, and minimalist Go web framework providing optimized HTTP router, group APIs, data binding for JSON and XML, HTTP/2 support, and much more.
*   [Flamingo](https://www.flamingo.me), a fast open-source framework based on Go with clean and scalable architecture designed to build custom, fast and flexible frontend interfaces. Includes Flamingo Core and Flamingo Commerce.
*   [Gin](https://gin-gonic.com/), a web framework written in Go. It features a martini-like API with much better performance, up to 40 times faster. If you need performance and good productivity, you will love Gin.
*   [Gorilla](http://www.gorillatoolkit.org/), a web toolkit for the Go programming language. Packages include a powerful URL router and dispatcher, context, RPC, schema, sessions, websocket, and more.


{{< headerWithLink header="Routers" search="http router" level=3 >}} 

{{< headerWithLink header="Template Engines" search="templates" level=3 >}} 

{{< headerWithLink header="Databases" search="database" level=3 >}} 

*   [Bleve](http://blevesearch.com/), full-text search and indexing for Go, providing a simple top-level API to index any object in the data model and override default mapping to customize behavior.
*   [CockroachDB](https://www.cockroachlabs.com/), an evolution of the database—architected for the cloud to deliver resilient, consistent, distributed SQL at scale.


### Courses
* https://www.usegolang.com

### Projects
*   [GopherJS](https://github.com/gopherjs/gopherjs), a compiler from Go to JavaScript allowing developers to write front-end code in Go which will run in all browsers.
*   [Hugo](https://gohugo.io/)
*   [Mattermost](https://mattermost.com/)
*   [Caddy](https://caddyserver.com/)
