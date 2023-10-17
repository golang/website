---
title: "Development Operations & Site Reliability Engineering"
linkTitle: "Development Operations & Site Reliability Engineering"
description: "With fast build times, lean syntax, an automatic formatter and doc generator, Go is built to support both DevOps and SRE."
date: 2019-10-03T17:16:43-04:00
series: Use Cases
books:
icon:
  file: devops-green.svg
  alt: ops icon
iconDark:
  file: devops-white.svg
  alt: ops icon
---

## Overview {#overview .sectionHeading}

### Go helps enterprises automate and scale

Development Operations (DevOps) teams help engineering organizations automate tasks and improve their continuous
integration and continuous delivery and deployment (CI/CD) process. DevOps can topple developmental silos and implement
tooling and automation to enhance software development, deployment, and support.

Site Reliability Engineering (SRE) was born at Google to make the company’s “large-scale sites more reliable, efficient,
and scalable,”[ writes Silvia Fressard](https://opensource.com/article/18/10/what-site-reliability-engineer), an
independent DevOps consultant. “And the practices they developed responded so well to Google’s needs that other big tech
companies, such as Amazon and Netflix, also adopted them.” SRE requires a mix of development and operations skills, and
“[empowers software developers](https://stackify.com/site-reliability-engineering/) to own the ongoing daily operation
of their applications in production.”

Go serves both siblings, DevOps and SRE, from its fast build times and lean syntax to its security and reliability support. Go's concurrency and networking features also make it ideal for tools that manage cloud deployment—readily supporting automation while
scaling for speed and code maintainability as development infrastructure grows over time.

DevOps/SRE teams write software ranging from small scripts, to command-line interfaces (CLI), to complex automation and services, and Go’s feature set has benefits for every situation.

## Key Benefits {#key-benefits .sectionHeading}

### Easily build small scripts with Go’s robust standard library and static typing
Go’s fast build and startup times. Go’s extensive standard library—including packages for
common needs like HTTP, file I/O, time, regular expressions, exec, and JSON/CSV formats—lets DevOps/SREs get right into their business logic. Plus, Go’s static type system and explicit error handling make even small scripts more robust.

### Quickly deploy CLIs with Go’s fast build times
Every site reliability engineer has written “one-time use” scripts that turned into CLIs used by dozens of other engineers every day. And small deployment automation scripts turn into rollout management services. With Go, DevOps/SREs are in a great position to be successful when software scope inevitably creeps. Starting with Go puts you in a great position to be successful when that happens.

### Scale and maintain larger applications with Go’s low memory footprint and doc generator
Go’s garbage collector means DevOps/SRE teams don’t have to worry about memory management. And Go’s automatic documentation generator (godoc) makes code self-documenting–lowering maintenance overhead and establishing best practices from the get-go.

{{projects `
  - company: Docker
    url: https://docker.com/
    logoSrc: docker.svg
    logoSrcDark: docker.svg
    desc: Docker is a software-as-a-service (SaaS) product, written in Go, that DevOps/SRE teams leverage to “drive secure automation and deployment at massive scale,” supporting their CI/CD efforts.
    ctas:
      - text: Docker CI/CD
        url: https://www.docker.com/solutions/cicd
  - company: Drone
    url: https://github.com/drone
    logoSrc: drone.svg
    logoSrcDark: drone.svg
    desc: Drone is a Continuous Delivery system built on container technology, written in Go, that uses a simple YAML configuration file, a superset of docker-compose, to define and execute Pipelines inside Docker containers.
    ctas:
      - text: Drone
        url: https://github.com/drone
  - company: etcd
    url: https://github.com/etcd-io/etcd
    logoSrc: etcd.svg
    logoSrcDark: etcd.svg
    desc: etcd is a strongly consistent, distributed key-value store that provides a reliable way to store data that needs to be accessed by a distributed system or cluster of machines, and it's written in Go.
    ctas:
      - text: etcd
        url: https://github.com/etcd-io/etcd
  - company: IBM
    url: https://ibm.com/
    logoSrc: ibm.svg
    logoSrcDark: ibm.svg
    desc: IBM’s DevOps teams use Go through Docker and Kubernetes, plus other DevOps and CI/CD tools written in Go. The company also supports connection to it’s messaging middleware through a Go-specific API.
    ctas:
      - text: IBM Applications in Golang
        url: https://developer.ibm.com/messaging/2019/02/05/simplified-ibm-mq-applications-golang/
  - company: Netflix
    url: http://netflix.com/
    logoSrc: netflix.svg
    logoSrcDark: netflix.svg
    desc: Netflix uses Go to handle large scale data caching, with a service called Rend, which manages globally replicated storage for personalization data.
    ctas:
      - text: Application Data Caching
        url: https://medium.com/netflix-techblog/application-data-caching-using-ssds-5bf25df851ef
      - text: Rend
        url: https://github.com/netflix/rend
  - company: Microsoft
    url: https://microsoft.com/
    logoSrc: microsoft_light.svg
    logoSrcDark: microsoft_dark.svg
    desc: Microsoft uses Go in Azure Red Hat OpenShift services. This Microsoft solution provides DevOps teams with OpenShift clusters to maintain regulatory compliance and focus on application development.
    ctas:
      - text: OpenShift
        url: https://azure.microsoft.com/en-us/services/openshift/
  - company: Terraform
    url: https://terraform.io/
    logoSrc: terraform-icon.svg
    logoSrcDark: terraform-icon.svg
    desc: Terraform is a tool for building, changing, and versioning infrastructure safely and efficiently. It supports a number of cloud providers such as AWS, IBM Cloud, GCP, and Microsoft Azure - and it’s written in Go.
    ctas:
      - text: Terraform
        url: https://www.terraform.io/intro/index.html
  - company: Prometheus
    url: https://github.com/prometheus/prometheus
    logoSrc: prometheus.svg
    logoSrcDark: prometheus.svg
    desc: Prometheus is an open-source systems monitoring and alerting toolkit originally built at SoundCloud. Most Prometheus components are written in Go, making them easy to build and deploy as static binaries.
    ctas:
      - text: Prometheus
        url: https://github.com/prometheus/prometheus
  - company: YouTube
    url: https://youtube.com/
    logoSrc: youtube.svg
    logoSrcDark: youtube.svg
    desc: YouTube uses Go with Vitess (now part of PlanetScale), its database clustering system for horizontal scaling of MySQL through generalized sharding. Since 2011 it’s been a core component of YouTube’s database infrastructure, and has grown to encompass tens of thousands of MySQL nodes.
    ctas:
      - text: Vitess
        url: https://github.com/vitessio/vitess
`}}

## Get Started {#get-started .sectionHeading}

### Go books on DevOps & SRE

{{books `
  - title: Go Programming for Network Operations
    url: https://www.amazon.com/Go-Programming-Network-Operations-Automation-ebook/dp/B07JKKN34L/ref=sr_1_16
    thumbnail: /images/books/go-programming-for-network-operations.jpg
  - title: Go Programming Blueprints
    url: https://github.com/matryer/goblueprints
    thumbnail: /images/learn/go-programming-blueprints.png
  - title: Go in Action
    url: https://www.amazon.com/Go-Action-William-Kennedy/dp/1617291781
    thumbnail: /images/books/go-in-action.jpg
  - title: The Go Programming Language
    url: https://www.gopl.io/
    thumbnail: /images/learn/go-programming-language-book.png
`}}

{{libraries `
  - title: Monitoring and tracing
    viewMoreUrl: https://pkg.go.dev/search?q=tracing
    items:
      - text: open-telemetry/opentelemetry-go
        url: https://pkg.go.dev/go.opentelemetry.io/otel
        desc: Vendor-neutral APIs and instrumentation for monitoring and distributed tracing
      - text: jaegertracing/jaeger-client-go
        url: https://pkg.go.dev/github.com/jaegertracing/jaeger-client-go?tab=overview
        desc: An open source distributed tracing system developed by Uber formats
      - text: grafana/grafana
        url: https://pkg.go.dev/github.com/grafana/grafana?tab=overview
        desc: An open-source platform for monitoring and observability
      - text: istio/istio
        url: https://pkg.go.dev/github.com/istio/istio?tab=overview
        desc: An open-source service mesh and integratable platform
  - title: CLI Libraries
    viewMoreUrl: https://pkg.go.dev/search?q=command%20line%20OR%20CLI
    items:
      - text: spf13/cobra
        url: https://pkg.go.dev/github.com/spf13/cobra?tab=overview
        desc: A library for creating powerful modern CLI applications and a program to generate applications and CLI applications in Go
      - text: spf13/viper
        url: https://pkg.go.dev/github.com/spf13/viper?tab=overview
        desc: A complete configuration solution for Go applications, designed to work within an app to handle configuration needs and formats
      - text: urfave/cli
        url: https://pkg.go.dev/github.com/urfave/cli?tab=overview
        desc: A minimal framework for creating and organizing command line Go applications
  - title: Other projects
    items:
      - text: golang-migrate/migrate
        url: https://pkg.go.dev/github.com/golang-migrate/migrate?tab=overview
        desc: A database migration tool written in Go
`}}
