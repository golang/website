---
title: "PayPal Taps Go to Modernize and Scale"
description: ""
date: 2019-10-03T16:26:31-04:00
series: Case Studies
resources:
- name: logo
  src: logo.svg
  params:
    alt: PayPal
---
## PayPal Taps Go to Modernize and Scale


### New Code Infrastructure Built on Go


PayPal was created to democratize financial services and empower people and businesses to join and thrive in the global economy. Central to this effort is PayPal's Payments Platform, which uses a combination of proprietary and third-party technologies to efficiently and securely facilitate transactions between millions of merchants and consumers worldwide. As the Payments Platform grew larger and more complicated, PayPal sought to modernize its systems and reduce time-to-market for new applications.

Go’s value in producing clean, efficient code that readily scales as software deployment scales made the language a strong fit to support PayPal’s goals.
 
PayPal had a NoSQL database in C++, but the complexity of the code was substantially decreasing its developers’ productivity. Go’s goroutines (lightweight threads of execution) and channels (which serve as the pipes that connect concurrent goroutines), made Go a natural choice in the NoSQL development team’s assessment for piloting system simplification and modernization.
 
As a proof of concept, a development team spent six months learning Go and reimplementing the NoSQL system from the ground up in Go, during which they also provided insights on how Go could be implemented more broadly at PayPal.


## **Using Go To Simplify For Scale**

As PayPal’s platform becomes more intricate, Go provides a way to readily simplify the complexity of creating and running software at scale. The language provides PayPal with great libraries and fast tools, plus concurrency, garbage collection, and type safety.
 
With Go, PayPal enables its developers to spend more time looking at code and thinking strategically, by freeing them from the noise of C++ development.
 
“Since our NoSQL and DB proxy used quite a bit of system details in a multi-threaded mode, the code got complex managing the different conditions,” said Bala Natarajan, Sr. Director of Engineering, Developer Experience and one of PayPal’s early Go evangelists. “Given that Go provides channels and routines to deal with complexity, we were able to structure the code to meet our requirements.”
 
After the success of this newly re-written NoSQL system, more platform and content teams within PayPal began adopting Go. Natarajan’s current team is responsible for PayPal’s build, test, and release pipelines—all built in Go. The company has a large build farm which is completely managed using Go infrastructure and supports builds-as-a-service for developers across the company.



### **Modernizing PayPal Systems With Go**

With the distributed computing required by PayPal, Go was the right language to refresh their systems. PayPal needed programming that is concurrent and parallel, compiled for high performance and highly portable, and that brings developers the benefits of a modular, composable open-source architecture—Go has delivered all that and more to help PayPal modernize its systems.

Security and supportability are key matters at PayPal, and the company’s operational pipelines are increasingly dominated by Go because the language’s cleanliness and modularity help them achieve these goals. PayPal’s deployment of Go engenders a platform of creativity for developers, allowing them to produce simple, efficient, and reliable software at scale for PayPal’s worldwide markets.
 
As PayPal continues to modernize their software-defined networking (SDN) infrastructure with Go, they are seeing performance benefits in addition to more maintainable code. For example, Go now powers routers, load balances, and an increasing number of production systems. **“In our tightly managed environments where we run Go code,”** says Natarajan, **“We have seen a CPU reduction of approximately 10% with cleaner and maintainable code.”**


### **Go Increases Developer Productivity**

As a global operation, PayPal needs its development teams to be effective at managing two kinds of scale: production scale, especially concurrent systems interacting with many other servers (such as cloud services); and development scale, especially large codebases developed by many programmers in coordination (such as open-source development).
 
PayPal leverages Go to address these issues of scale. The company’s developers benefit from Go’s ability to combine the ease of programming of an interpreted, dynamically typed language with the efficiency and safety of a statically typed, compiled language. As PayPal modernizes its system, support for networked and multicore computing is critical. Go not only delivers such support but delivers quickly—it takes at most a few seconds to compile a large executable on a single computer.
 
There are currently over 100 Go developers at PayPal, and future developers who choose to adopt Go will have an easier time getting the language approved thanks to the many successful implementations already in production at the company.
 
Most importantly, PayPal developers have increased their productivity with Go. Go’s concurrency mechanisms have made it easy to write programs that get the most out of PayPal’s multicore and networked machines. Developers using Go also benefit from the fact that it compiles quickly to machine code, and their apps gain the convenience of garbage collection and the power of run-time reflection.


## **Speeding PayPal's Time To Market**

The first-class languages at PayPal today are Java and Node, with Go primarily used as an infrastructure language. While Go may never replace Node.js for certain applications, Natarajan is pushing to make Go a first-class language at PayPal.
 
Through his efforts, PayPal is also evaluating moving to the Google Kubernetes Engine (GKE) to speed their new products’ time-to-market. The GKE is a managed, production-ready environment for deploying containerized applications, and brings Google's latest innovations in developer productivity, automated operations, and open source flexibility.

For PayPal, deploying the GKE would enable rapid development and iteration by making it easier for PayPal to deploy, update, and manage its applications and services. Plus PayPal will find it easier to run Machine Learning, General Purpose GPU, High-Performance Computing, and other workloads that benefit from specialized hardware accelerators supported by the GKE.
 
Most importantly for PayPal, the combination of Go development and the GKE allows the company to scale effortless to meet demand, as Kubernetes autoscaling will allow PayPal to handle increased user demand for services—keeping them available when it matters most, then scale back in the quiet periods to save money.


## **Getting Your Enterprise Started With Go**

PayPal’s story is not unique; dozens of other large enterprises are discovering how Go can help them ship reliable software faster. There are approximately two million developers using Go worldwide, spanning banking and commerce, gaming and media, technology, and other industries, at enterprises as diverse as Amazon, Apple, Dropbox, Epic Games, Facebook, IBM, Microsoft, Square, Uber, and of course Google.
 
To learn more about how Go can help your enterprise build reliable, scalable software as it does at PayPal,visit go.dev today.

