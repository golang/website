---
title: "Go for Web Development"
linkTitle: "Web Development"
description: "With enhanced memory performance and support for several IDEs, Go powers fast and scalable web applications."
date: 2019-10-04T15:26:31-04:00
series: Use Cases
books:
icon:
  file: webdev-green.svg
  alt: web dev icon
iconDark:
  file: webdev-white.svg
  alt: web dev icon
---

## Overview {#overview .sectionHeading}

### Go delivers speed, security, and developer-friendly tools for Web Applications

Go is designed to enable developers to rapidly develop scalable and secure web applications. Go ships with an easy to use, secure and performant web server and includes it own web templating library. Go has excellent support for all of the latest technologies from [HTTP/2](https://pkg.go.dev/net/http), to databases like [MySQL](https://pkg.go.dev/mod/github.com/go-sql-driver/mysql), [MongoDB](https://pkg.go.dev/mod/go.mongodb.org/mongo-driver) and [Elasticsearch](https://pkg.go.dev/mod/github.com/elastic/go-elasticsearch/v8), to the latest encryption standards including [TLS 1.3](https://pkg.go.dev/crypto/tls). Go web applications run natively on [Google App Engine](https://cloud.google.com/appengine/) and [Google Cloud Run](https://cloud.google.com/run/) (for easy scaling) or on any environment, cloud, or operating system thanks to Go’s extreme portability.

## Key Benefits {#key-benefits .sectionHeading}

### Deploy across platforms in record speed

For enterprises, Go is preferred for providing rapid cross-platform deployment. With its goroutines, native compilation, and the URI-based package namespacing, Go code compiles to a single, small binary—with zero dependencies—making it very fast.

### Leverage Go’s out-of-the-box performance to scale with ease

Tigran Bayburtsyan, Co-Founder and CTO at Hexact Inc., summarizes five key reasons his company switched to Go:

-   **Compiles into a single binary** — “Using static linking, Go actually combining all dependency libraries and modules into one single binary file based on OS type and architecture.”

-   **Static type system** — “Type system is really important for large scale applications.”

-   **Performance** — “Go performed better because of its concurrency model and CPU scalability. Whenever we need to process some internal request, we are doing it with separate Goroutines which are 10x cheaper in resources than Python Threads.”

-   **No need for a web framework** — “In most of the cases you really don’t need any third-party library.”

-   **Great IDE support and debugging** — “After rewriting all projects to Go, we got 64 percent less code than we had earlier.”


{{projects `
  - company: Caddy
    url: https://caddyserver.com/
    logoSrc: caddy.svg
    logoSrcDark: caddy.svg
    desc: Caddy 2 is a powerful, enterprise-ready, open source web server with automatic HTTPS written in Go. Caddy offers greater memory safety than servers written in C. A hardened TLS stack powered by the Go standard library serves a significant portion of all Internet traffic.
    ctas:
      - text: Caddy 2
        url: https://caddyserver.com/
  - company: Cloudflare
    url: https://www.cloudflare.com/en-gb/
    logoSrc: cloudflare-icon.svg
    logoSrcDark: cloudflare-icon.svg
    desc: Cloudflare speeds up and protects millions of websites, APIs, SaaS services, and other properties connected to the Internet. “Go is at the heart of CloudFlare’s services including handling compression for high-latency HTTP connections, our entire DNS infrastructure, SSL, load testing and more.”
    ctas:
      - text: Cloudflare and Go
        url: https://blog.cloudflare.com/what-weve-been-doing-with-go/
  - company: gov.uk
    url: https://gov.uk/
    logoSrc: govuk_light.svg
    logoSrcDark: govuk_dark.svg
    desc: The simplicity and safety of the Go language were a good fit for the United Kingdom’s government’s HTTP infrastructure, and some brief experiments with the excellent net/http package convinced web developers they were on the right track. “In particular, Go’s concurrency model makes it absurdly easy to build performant I/O-bound applications.”
    ctas:
      - text: Building a new router for gov.uk
        url: https://technology.blog.gov.uk/2013/12/05/building-a-new-router-for-gov-uk/
      - text: Using Go in government
        url: https://technology.blog.gov.uk/2014/11/14/using-go-in-government/
  - company: Hugo
    url: http://gohugo.io/
    logoSrc: hugo.svg
    logoSrcDark: hugo.svg
    desc: Hugo is a fast and modern website engine written in Go, and designed to make website creation fun again. Websites built with Hugo are extremely fast and secure and can be hosted anywhere without any dependencies.
    ctas:
      - text: Hugo
        url: http://gohugo.io/
  - company: Mattermost
    url: https://mattermost.com/
    logoSrc: mattermost_light.svg
    logoSrcDark: mattermost_dark.svg
    desc: Mattermost is a flexible, open source messaging platform that enables secure team collaboration. It’s written in Go and React.
    ctas:
      - text: Mattermost
        url: https://mattermost.com/
  - company: Medium
    url: https://medium.org/
    logoSrc: medium_light.svg
    logoSrcDark: medium_dark.svg
    desc: Medium uses Go to power their social graph, their image server and several auxiliary services. “We’ve found Go very easy to build, package, and deploy. We like the type-safety without the verbosity and JVM tuning of Java.”
    ctas:
      - text: Medium's Go Services
        url: https://medium.engineering/how-medium-goes-social-b7dbefa6d413
  - company: The Economist
    url: https://economist.com/
    logoSrc: economist.svg
    logoSrcDark: economist.svg
    desc: The Economist needed more flexibility to deliver content to increasingly diverse digital channels. Services written in Go were a key component of the new system that would enable The Economist to deliver scalable, high performing services and quickly iterate new products. “Overall, it was determined that Go was the language best designed for usability and efficiency in a distributed, cloud-based system.”
    ctas:
      - text: The Economist's Go microservices
        url: https://www.infoq.com/articles/golang-the-economist/
`}}

## Get Started {#get-started .sectionHeading}

### Go books on web development

{{books `
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

### Courses
* [Learn to Create Web Applications using Go](https://www.usegolang.com), a paid online course

### Projects
*   {{pkg "github.com/gopherjs/gopherjs" "gopherjs"}}, a compiler from Go to JavaScript allowing developers to write front-end code in Go which will run in all browsers.
*   [Hugo](https://gohugo.io/), The world’s fastest framework for building websites
*   [Mattermost](https://mattermost.com/), a flexible, open source messaging platform
that enables secure team collaboration
*   [Caddy](https://caddyserver.com/), a powerful, enterprise-ready, open source web server with automatic HTTPS written in Go
