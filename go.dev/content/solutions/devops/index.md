---
title: "Go for Development Operations & Site Reliability Engineering (SRE)"
linkTitle: "DevOps & Site Reliability"
description: "With fast build times, lean syntax, an automatic formatter and doc generator, Go is built to support both DevOps and SRE."
date: 2019-10-03T17:16:43-04:00
series: Use Cases
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

### _Go helps enterprises automate and scale for CI/CD_


## **Why use Go for DevOps & SRE**

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


## **Who uses Go for DevOps & SRE**


### **Docker**

Docker is a software-as-a-service (SaaS) product, written in Go, that uses operating-system level virtualization to
develop and deliver software in containers hosted on a Docker Engine. DevOps/SRE teams leverage Docker to “[drive secure
automation and deployment at massive scale](https://www.docker.com/solutions/cicd),” supporting their CI/CD efforts.

 


### **Google**

Google leverages Go for the Google Cloud Platform (GCP), as well as at the heart of Kubernetes—an open-source
container-orchestration system, written in Go, for automating application deployment, scaling, and management. At[
Google](https://landing.google.com/sre/), SRE's "protect, provide for, and progress the software and systems behind all
of Google’s public services—Google Search, Ads, Gmail, Android, YouTube, and App Engine, to name just a few—with an
ever-watchful eye on their availability, latency, performance, and capacity… they keep important, revenue-critical
systems up and running.”


### **IBM**

IBM’s DevOps teams are heavily invested in Docker and Kubernetes, plus other DevOps and CI/CD tools written in Go as
found on the company’s [GitHub](https://github.com/IBM?utf8=%E2%9C%93&q=&type=&language=go). IBM engineering
organizations leverage Red Hat's cloud platform, [OpenShift](https://www.openshift.com), written primarily in Go, and
Red Hat's new addition, [CoreOS](https://coreos.com). CoreOS, also written in Go, delivers one of the best enterprise
Kubernetes distributions available in Tectonic—bringing automated operations, Open Cloud Services, Prometheus
monitoring, and more to simplify Kubernetes deployments, reduce engineering operating costs, and speed time to
production.


### **Microsoft**

Microsoft DevOps includes the company's fully managed Azure Kubernetes Service
([AKS](https://azure.microsoft.com/en-us/services/kubernetes-service/)). AKS was designed to make deploying and managing
containerized applications easy by offering serverless Kubernetes, an integrated CI/CD experience, and enterprise-grade
security and governance.

Like IBM, Microsoft is also leveraging Red Hat's OpenShift, written in Go, via [Azure Red Hat
OpenShift](https://azure.microsoft.com/en-us/services/openshift/) services. This Microsoft solution provides DevOps
teams with OpenShift clusters to maintain regulatory compliance and focus on application development.

[summary of how Microsoft uses Go for DevOps/SRE?  Does Robert van Gent have details?]


### **Terraform**

Terraform is a[ tool for building, changing, and versioning infrastructure](https://www.terraform.io/intro/index.html)
safely and efficiently. It can manage existing and popular service providers as well as custom in-house solutions.
Written in Go, Terraform supports a number of cloud infrastructure providers such as AWS, IBM Cloud, GCP, and Microsoft
Azure. From a DevOps/SRE perspective, Terraform describes infrastructure as code using a high-level configuration
syntax. It leverages execution plans and resource graphs to automate changes to infrastructure with minimal human
interaction.


## **How to use Go for DevOps & SRE**

Go has been enthusiastically adopted by the DevOps and SRE communities. As previously noted, many underpinnings of the
modern cloud environment are themselves written in Go—including Docker, Etcd, Istio, Kubernetes, Prometheus, Terraform,
and many others.

[do we want text for featured users/projects, or just logos? links?]

[can we say that Google has switched to recommend Go for all new SRE code?]

DevOps/SRE teams write software ranging from small one-time scripts, to command-line interfaces (CLI), to complex
automation and services where Go excels for all of them.

**For small scripts:** Go's fast build times and automatic formatter ([gofmt](https://golang.org/cmd/gofmt/)) enable rapid iteration. Go’s extensive standard library—including packages for common needs like HTTP, file I/O, time, regular expressions, exec, and JSON/CSV formats—lets DevOps/SREs get right into their business logic. Plus, Go's static type system and explicit error handling make even small scripts more robust. 

**For CLIs:** every site reliability engineer has written “one-time use” scripts that turned into CLIs used by dozens of other engineers every day. And small deployment automation scripts turn into rollout management services. With Go, DevOps/SREs are in a great position to be successful when software scope inevitably creeps. 

**For larger applications:** Go's garbage collector means DevOps/SRE teams don't have to worry about memory management. And Go’s automatic documentation generator (godoc) makes code self-documenting.

With Go, DevOps/SREs seek to “balance the risk of unavailability with the goals of rapid innovation and efficient
service operations,"[ says Marc Alvidrez](https://landing.google.com/sre/), engineer at Google. "So that users’ overall
happiness—with features, service, and performance—is optimized."

## **Go solutions to legacy challenges**

Traditionally, “DevOps has been more about collaboration between developer and operations. It has also focused more on
deployments,"[ says Matt Watson](https://stackify.com/site-reliability-engineering/), founder and CEO of Stackify. "Site
reliability engineering is more focused on operations and monitoring. Depending on how you define DevOps, it could be
related or not."

Across deployment, operations, and monitoring, DevOps/SRE teams strive to achieve simplicity, reliability, and speed
with their systems. But in complex development environments, such disparate goals are hard to unite. 

**Go helps by
allowing engineers to focus on building**, even as they optimize for deployment and support.

For simplicity, Go delivers code readability, built in testing/profiling/benchmarking, a standard library, and a
homogenous environment—statically linked—[ meaning](https://blog.gopheracademy.com/advent-2018/go-devops/) “there’s no
need for external libraries, copy dependencies or worry for imports. All the code and its dependencies are in the
binary, so that’s all you need to distribute.”

 

For reliability, open source Go delivers pointers, error handling, and safe Type, meaning string operations on an int
cannot happen, because it will be caught by the compiler.

 

For speed, Go delivers fast compilation and machine-code execution, small binary sizes, superior garbage collection, and
import-defined dependences, meaning all dependencies are included in the binary. For a list of practical Go benchmarks,
visit this[ list of performance benchmarks](https://stackimpact.com/blog/practical-golang-benchmarks/) in various
functionalities.

 

"With systems becoming distributed and more complex—spread over a group of services (or microservices),"[ writes Natalie
Pistunovich](https://blog.gopheracademy.com/advent-2018/go-devops/), engineering manager at Fraugster. "Observability is
becoming a trade that helps keep you on track with the system’s health."

 

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
