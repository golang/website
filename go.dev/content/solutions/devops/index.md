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
  - title: 
    url: 
    thumbnail: 
  - title: 
    url: 
    thumbnail: 
  - title: 
    url: 
    thumbnail: 
  - title: 
    url: 
    thumbnail: 
resources:
- name: icon
  src: cog.png
  params:
    alt: cog
- name: icon-white
  src: cog-white.png
  params:
    alt: cog
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

Go serves both siblings, DevOps and SRE, with its fast build times and lean syntax—readily supporting automation while
scaling for speed and code maintainability as development infrastructure grows over time.

## Featured Go users & projects

{{% mediaList %}}
    {{% mediaListBox img-src="/images/logos/paypal.svg" img-alt="Paypal Logo" title="" align=top %}}
PayPal [uses Go across its payment ecosystem](/solutions/paypal), from build, test, and release pipelines, to NoSQL databases, to a large build farm completely managed in Go. The team originally explored using Go to decrease the complexity of a NoSQL databases' code (which was written in C++). The NoSQL team rebuilt one of their databases in Go and immediately noticed the benefits of having a language built for concurrency, simplicity, and speed.
    {{% /mediaListBox %}}
    {{% mediaListBox img-src="/images/logos/netflix.svg" img-alt="Netflix Logo" title="" align=top %}}
Netflix uses Go to [handle large scale data caching](https://medium.com/netflix-techblog/application-data-caching-using-ssds-5bf25df851ef), with a service called [Rend](https://github.com/netflix/rend), which manages globally replicated storage for personalization data. It's a high-performance server written in Go that acts as a proxy in front of other processes that actually store the data - and plays a critical role in member personalization for over 150 million users.
    {{% /mediaListBox %}}
    {{% mediaListBox img-src="/images/logos/docker.svg" img-alt="Docker Logo" img-link="https://docker.com" title="" align=top %}}
Docker is a software-as-a-service (SaaS) product, written in Go, that DevOps/SRE teams leverage to “[drive secure
automation and deployment at massive scale](https://www.docker.com/solutions/cicd),” supporting their CI/CD efforts.
    {{% /mediaList %}}
    {{% mediaListBox img-src="/images/logos/google.svg" img-alt="Google Logo" img-link="https://google.com" title="" align=top %}}
[Google SREs use Go](https://landing.google.com/sre/) to protect, provide for, and progress the software and systems behind all
of Google’s public services. Google now recommends that Go be used for all new SRE code. 
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
safely and efficiently. It supports a number of cloud providers such as AWS, IBM Cloud, GCP, and Microsoft
Azure - and its written in Go.
    {{% /mediaList %}}

{{% /mediaList %}}


## **How to use Go for DevOps & SRE**

Go has been enthusiastically adopted by the DevOps and SRE communities. As previously noted, many underpinnings of the
modern cloud environment are themselves written in Go—including Docker, Etcd, Istio, Kubernetes, Prometheus, Terraform,
and many others.

DevOps/SRE teams write software ranging from small one-time scripts, to command-line interfaces (CLI), to complex
automation and services, and Go's feature set has benefits for every situation.

**For small scripts:** Go's fast build times and automatic formatter ([gofmt](https://golang.org/cmd/gofmt/)) enable rapid iteration. Go’s extensive standard library—including packages for common needs like HTTP, file I/O, time, regular expressions, exec, and JSON/CSV formats—lets DevOps/SREs get right into their business logic. Plus, Go's static type system and explicit error handling make even small scripts more robust. 

**For CLIs:** every site reliability engineer has written “one-time use” scripts that turned into CLIs used by dozens of other engineers every day. And small deployment automation scripts turn into rollout management services. With Go, DevOps/SREs are in a great position to be successful when software scope inevitably creeps. 

**For larger applications:** Go's garbage collector means DevOps/SRE teams don't have to worry about memory management. And Go’s automatic documentation generator (godoc) makes code self-documenting.

With Go, DevOps/SREs seek to “balance the risk of unavailability with the goals of rapid innovation and efficient
service operations,"[ says Marc Alvidrez](https://landing.google.com/sre/), engineer at Google. "So that users’ overall
happiness—with features, service, and performance—is optimized."

## Key Solutions

### Go books on DevOps & SRE

{{% books %}}

{{< headerWithLink header="Frameworks" link="https://pkg.go.dev/search?q=framework" level=3 >}} 

## **Go solutions to legacy challenges**

Traditionally, “DevOps has been more about collaboration between developer and operations. It has also focused more on
deployments,"[ says Matt Watson](https://stackify.com/site-reliability-engineering/), founder and CEO of Stackify. "Site
reliability engineering is more focused on operations and monitoring. Depending on how you define DevOps, it could be
related or not."

{{% gopher gopher=machine align=right %}}

Across deployment, operations, and monitoring, DevOps/SRE teams strive to achieve simplicity, reliability, and speed
with their systems. But in complex development environments, such disparate goals are hard to unite. 

**Go helps by
allowing engineers to focus on building**, even as they optimize for deployment and support.

For simplicity, Go delivers code readability, built in testing/profiling/benchmarking, a standard library, and a
homogenous environment that is statically linked.

{{% pullquote author="Natalie Pistunovich, Engineering Manager at Fraugster." link="https://blog.gopheracademy.com/advent-2018/go-devops/" %}}
[With Go] there’s no need for external libraries, copy dependencies or worry for imports. All the code and its dependencies are in the
binary, so that’s all you need to distribute.
{{% /pullquote%}}

For reliability, open source Go delivers pointers, error handling, and safe Type, meaning string operations on an int
cannot happen, because it will be caught by the compiler.

For speed, Go delivers fast compilation and machine-code execution, small binary sizes, superior garbage collection, and
import-defined dependences, meaning all dependencies are included in the binary. For a list of practical Go benchmarks,
visit this[ list of performance benchmarks](https://stackimpact.com/blog/practical-golang-benchmarks/) in various
functionalities.

Many of the modern tooling apps, for DevOps/SRE and for observability, are written in Go. For example:

*   [Grafana](https://grafana.com/)
*   [Helm](https://helm.sh/)
*   [Istio](https://istio.io/)
*   [Jaeger](https://www.jaegertracing.io/)
*   [The Open Tracing Project](https://opentracing.io/)

As DevOps/SRE teams automate the processes between software development and IT teams, Go can help them build, test, and
release software faster and more reliably. Scaling infrastructure and development for CI/CD is critical to many large
technology firms today, and Go is the right language for enterprises looking to scale successfully.


## **Resources for learning more**

*   [Migrate](https://github.com/golang-migrate/migrate) - database migration tool written in Go
