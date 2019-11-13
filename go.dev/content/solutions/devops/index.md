---
title: "Go for Development Operations & Site Reliability Engineering (SRE)"
linkTitle: "DevOps & Site Reliability"
description: "With fast build times, lean syntax, an automatic formatter and doc generator, Go is built to support both DevOps and SRE."
date: 2019-10-03T17:16:43-04:00
series: Use Cases
books:
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
resources:
- name: icon
  src: ops-green.svg
  params:
    alt: ops icon
- name: icon-white
  src: ops-white.svg
  params:
    alt: ops icon
---

## Go helps enterprises automate and scale

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

## Featured Go users & projects

{{% mediaList %}}
    {{% mediaListBox img-src="/images/logos/netflix.svg" img-alt="Netflix Logo" title="" align=top %}}
Netflix uses Go to [handle large scale data caching](https://medium.com/netflix-techblog/application-data-caching-using-ssds-5bf25df851ef), with a service called [Rend](https://github.com/netflix/rend), which manages globally replicated storage for personalization data.
    {{% /mediaListBox %}}
    {{% mediaListBox img-src="/images/logos/docker.svg" img-alt="Docker Logo" img-link="https://docker.com" title="" align=top %}}
Docker is a software-as-a-service (SaaS) product, written in Go, that DevOps/SRE teams leverage to “[drive secure
automation and deployment at massive scale](https://www.docker.com/solutions/cicd),” supporting their CI/CD efforts.
    {{% /mediaList %}}
    {{% mediaListBox img-src="/images/logos/drone.svg" img-alt="Drone Logo" img-link="https://github.com/drone" title="" align=top %}}
Drone is a [Continuous Delivery system built on container technology](https://github.com/drone), written in Go, that uses a simple YAML configuration file, a superset of docker-compose, to define and execute Pipelines inside Docker containers.
    {{% /mediaList %}}
    {{% mediaListBox img-src="/images/logos/etcd.svg" img-alt="etcd Logo" img-link="https://github.com/etcd-io/etcd" title="" align=top %}}
etcd is a [strongly consistent, distributed key-value store](https://github.com/etcd-io/etcd) that provides a reliable way to store data that needs to be accessed by a distributed system or cluster of machines, and its written in Go.
    {{% /mediaList %}}
    {{% mediaListBox img-src="/images/logos/ibm.svg" img-alt="IBM Logo" img-link="https://ibm.com" title="" align=top %}}
IBM’s DevOps teams use Go through Docker and Kubernetes, plus other DevOps and CI/CD tools written in Go. The company also supports connection to it's messaging middleware through a [Go-specific API](https://developer.ibm.com/messaging/2019/02/05/simplified-ibm-mq-applications-golang/).
    {{% /mediaList %}}
    {{% mediaListBox img-src="/images/logos/microsoft.svg" img-alt="Microsoft Logo" img-link="https://microsoft.com" title="" align=top %}}
Microsoft uses Go in [Azure Red Hat
OpenShift](https://azure.microsoft.com/en-us/services/openshift/) services. This Microsoft solution provides DevOps
teams with OpenShift clusters to maintain regulatory compliance and focus on application development.
    {{% /mediaList %}}
    {{% mediaListBox img-src="/images/logos/terraform.png" img-alt="Terraform Logo" img-link="https://terraform.io" title="" align=top %}}
Terraform is a [tool for building, changing, and versioning infrastructure](https://www.terraform.io/intro/index.html)
safely and efficiently. It supports a number of cloud providers such as AWS, IBM Cloud, GCP, and Microsoft Azure - and its written in Go.
    {{% /mediaList %}}
    {{% mediaListBox img-src="/images/logos/prometheus.svg" img-alt="Prometheus Logo" img-link="https://github.com/prometheus/prometheus" title="" align=top %}}
Prometheus is an [open-source systems monitoring and alerting toolkit](https://github.com/prometheus/prometheus) originally built at SoundCloud. Most Prometheus components are written in Go, making them easy to build and deploy as static binaries.
    {{% /mediaList %}}
    {{% mediaListBox img-src="/images/logos/youtube.svg" img-alt="YouTube Logo" img-link="https://youtube.com" title="" align=top %}}
YouTube uses Go with Vitess (now part of [PlanetScale](https://planetscale.com/)), its [database clustering system](https://github.com/vitessio/vitess) for horizontal scaling of MySQL through generalized sharding. Since 2011 it's been a core component of YouTube's database infrastructure, and has grown to encompass tens of thousands of MySQL nodes.
    {{% /mediaList %}}

{{% /mediaList %}}


## **How to use Go for DevOps & SRE**

Go has been enthusiastically adopted by the DevOps and SRE communities. As previously noted, many underpinnings of the
modern cloud environment are themselves written in Go—including Docker, Etcd, Istio, Kubernetes, Prometheus, Terraform,
and many others.

DevOps/SRE teams write software ranging from small scripts, to command-line interfaces (CLI), to complex
automation and services, and Go's feature set has benefits for every situation.

**For small scripts:** Go's fast build and startup times. Go’s extensive standard library—including packages for common needs like HTTP, file I/O, time, regular expressions, exec, and JSON/CSV formats—lets DevOps/SREs get right into their business logic. Plus, Go's static type system and explicit error handling make even small scripts more robust. 

**For CLIs:** every site reliability engineer has written “one-time use” scripts that turned into CLIs used by dozens of other engineers every day. And small deployment automation scripts turn into rollout management services. With Go, DevOps/SREs are in a great position to be successful when software scope inevitably creeps. Starting with Go puts you in a great position to be successful when that happens.

**For larger applications:** Go's garbage collector means DevOps/SRE teams don't have to worry about memory management. And Go’s automatic documentation generator ([godoc](https://godoc.org/golang.org/x/tools/cmd/godoc)) makes code self-documenting.

## Key solutions

### Go books on DevOps & SRE

{{% books %}}

{{< headerWithLink header="Monitoring and tracing" search="tracing" level=3 >}} 

*   {{< pkg "https://github.com/opentracing/opentracing-go">}}, vendor-neutral APIs and instrumentation for distributed tracing
*   {{< pkg "https://github.com/jaegertracing/jaeger-client-go">}}, and open source distributed tracing system developed by Uber
*   {{< pkg "https://github.com/grafana/grafana">}}, an open-source platform for monitoring and observability
*   {{< pkg "https://github.com/istio/istio">}}, an open-source service mesh and integratable platform  

{{< headerWithLink header="CLI libraries" search="command line OR CLI" level=3 >}} 

*   {{< pkg "github.com/spf13/cobra" >}}, a library for creating powerful modern CLI applications and a program to generate applications and CLI applications in Go
*   {{< pkg "github.com/spf13/viper" >}}, a complete configuration solution for Go applications, designed to work within an app to handle configuration needs and formats
*   {{< pkg "github.com/urfave/cli" >}}, a minimal framework for creating and organizing command line Go applications

### Other

*   {{< pkg "https://github.com/golang-migrate/migrate">}}, database migration tool written in Go
