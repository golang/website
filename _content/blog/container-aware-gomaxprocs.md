---
title: Container-aware GOMAXPROCS
date: 2025-08-20
by:
- Michael Pratt
- Carlos Amedee
summary:  New GOMAXPROCS defaults in Go 1.25 improve behavior in containers.
---

Go 1.25 includes new container-aware `GOMAXPROCS` defaults, providing more sensible default behavior for many container workloads, avoiding throttling that can impact tail latency, and improving Go's out-of-the-box production-readiness.
In this post, we will dive into how Go schedules goroutines, how that scheduling interacts with container-level CPU controls, and how Go can perform better with awareness of container CPU controls.

## `GOMAXPROCS`

One of Go's strengths is its built-in and easy-to-use concurrency via goroutines.
From a semantic perspective, goroutines appear very similar to operating system threads, enabling us to write simple, blocking code.
On the other hand, goroutines are more lightweight than operating system threads, making it much cheaper to create and destroy them on the fly.

While a Go implementation could map each goroutine to a dedicated operating system thread, Go keeps goroutines lightweight with a runtime scheduler that makes threads fungible.
Any Go-managed thread can run any goroutine, so creating a new goroutine doesn't require creating a new thread, and waking a goroutine doesn't necessarily require waking another thread.

That said, along with a scheduler comes scheduling questions.
For example, exactly how many threads should we use to run goroutines?
If 1,000 goroutines are runnable, should we schedule them on 1,000 different threads?

This is where [`GOMAXPROCS`](/pkg/runtime#GOMAXPROCS) comes in.
Semantically, `GOMAXPROCS` tells the Go runtime the "available parallelism" that Go should use.
In more concrete terms, `GOMAXPROCS` is the maximum number of threads to use for running goroutines at once.

So, if `GOMAXPROCS=8` and there are 1,000 runnable goroutines, Go will use 8 threads to run 8 goroutines at a time.
Often, goroutines run for a very short time and then block, at which point Go will switch to running another goroutine on that same thread.
Go will also preempt goroutines that don't block on their own, ensuring all goroutines get a chance to run.

From Go 1.5 through Go 1.24, `GOMAXPROCS` defaulted to the total number of CPU cores on the machine.
Note that in this post, "core" more precisely means "logical CPU."
For example, a machine with 4 physical CPUs with hyperthreading has 8 logical CPUs.

This typically makes a good default for "available parallelism" because it naturally matches the available parallelism of the hardware.
That is, if there are 8 cores and Go runs more than 8 threads at a time, the operating system will have to multiplex these threads onto the 8 cores, much like how Go multiplexes goroutines onto threads.
This extra layer of scheduling is not always a problem, but it is unnecessary overhead.

## Container Orchestration

Another of Go's core strengths is the convenience of deploying applications via a container, and managing the number of cores Go uses is especially important when deploying an application within a container orchestration platform.
Container orchestration platforms like [Kubernetes](https://kubernetes.io/) take a set of machine resources and schedule containers within the available resources based on requested resources.
Packing as many containers as possible within a cluster's resources requires the platform to be able to predict the resource usage of each scheduled container.
We want Go to adhere to the resource utilization constraints that the container orchestration platform sets.

Let's explore the effects of the `GOMAXPROCS` setting in the context of Kubernetes, as an example.
Platforms like Kubernetes provide a mechanism to limit the resources consumed by a container.
Kubernetes has the concept of CPU resource limits, which signal to the underlying operating system how many core resources a specific container or set of containers will be allocated.
Setting a CPU limit translates to the creation of a Linux [control group](https://docs.kernel.org/admin-guide/cgroup-v2.html#cpu) CPU bandwidth limit.

Before Go 1.25, Go was unaware of CPU limits set by orchestration platforms.
Instead, it would set `GOMAXPROCS` to the number of cores on the machine it was deployed to.
If there was a CPU limit in place, the application may try to use far more CPU than allowed by the limit.
To prevent an application from exceeding its limit, the Linux kernel will [throttle](https://kubernetes.io/docs/concepts/configuration/manage-resources-containers/#how-pods-with-resource-limits-are-run) the application.

Throttling is a blunt mechanism for restricting containers that would otherwise exceed their CPU limit: it completely pauses application execution for the remainder of the throttling period.
The throttling period is typically 100ms, so throttling can cause substantial tail latency impact compared to the softer scheduling multiplexing effects of a lower `GOMAXPROCS` setting.
Even if the application never has much parallelism, tasks performed by the Go runtime—such as garbage collection—can still cause CPU spikes that trigger throttling.

## New default

We want Go to provide efficient and reliable defaults when possible, so in Go 1.25, we have made `GOMAXPROCS` take into account its container environment by default.
If a Go process is running inside a container with a CPU limit, `GOMAXPROCS` will default to the CPU limit if it is less than the core count.

Container orchestration systems may adjust container CPU limits on the fly, so Go 1.25 will also periodically check the CPU limit and adjust `GOMAXPROCS` automatically if it changes.

Both of these defaults only apply if `GOMAXPROCS` is otherwise unspecified.
Setting the `GOMAXPROCS` environment variable or calling `runtime.GOMAXPROCS` continues to behave as before.
The [`runtime.GOMAXPROCS`](/pkg/runtime#GOMAXPROCS) documentation covers the details of the new behavior.

## Slightly different models

Both `GOMAXPROCS` and a container CPU limit place a limit on the maximum amount of CPU the process can use, but their models are subtly different.

`GOMAXPROCS` is a parallelism limit.
If `GOMAXPROCS=8` Go will never run more than 8 goroutines at a time.

By contrast, CPU limits are a throughput limit.
That is, they limit the total CPU time used in some period of wall time.
The default period is 100ms.
So an "8 CPU limit" is actually a limit of 800ms of CPU time every 100ms of wall time.

This limit could be filled by running 8 threads continuously for the entire 100ms, which is equivalent to `GOMAXPROCS=8`.
On the other hand, the limit could also be filled by running 16 threads for 50ms each, with each thread being idle or blocked for the other 50ms.

In other words, a CPU limit doesn't limit the total number of CPUs the container can run on.
It only limits total CPU time.

Most applications have fairly consistent CPU usage across 100ms periods, so the new `GOMAXPROCS` default is a pretty good match to the CPU limit, and certainly better than the total core count!
However, it is worth noting that particularly spiky workloads may see a latency increase from this change due to `GOMAXPROCS` preventing short-lived spikes of additional threads beyond the CPU limit average.

In addition, since CPU limits are a throughput limit, they can have a fractional component (e.g., 2.5 CPU).
On the other hand, `GOMAXPROCS` must be a positive integer.
Thus, Go must round the limit to a valid `GOMAXPROCS` value.
Go always rounds up to enable use of the full CPU limit.

## CPU Requests

Go's new `GOMAXPROCS` default is based on the container's CPU limit, but container orchestration systems also provide a "CPU request" control.
While the CPU limit specifies the maximum CPU a container may use, the CPU request specifies the minimum CPU guaranteed to be available to the container at all times.

It is common to create containers with a CPU request but no CPU limit, as this allows containers to utilize machine CPU resources beyond the CPU request that would otherwise be idle due to lack of load from other containers.
Unfortunately, this means that Go cannot set `GOMAXPROCS` based on the CPU request, which would prevent utilization of additional idle resources.

Containers with a CPU request are still [constrained](https://kubernetes.io/docs/concepts/configuration/manage-resources-containers/#how-pods-with-resource-limits-are-run) when exceeding their request if the machine is busy.
The weight-based constraint of exceeding requests is "softer" than the hard period-based throttling of CPU limits, but CPU spikes from high `GOMAXPROCS` can still have an adverse impact on application behavior.

## Should I set a CPU limit?

We have learned about the problems caused by having `GOMAXPROCS` too high, and that setting a container CPU limit allows Go to automatically set an appropriate `GOMAXPROCS`, so an obvious next step is to wonder whether all containers should set a CPU limit.

While that may be good advice to automatically get a reasonable `GOMAXPROCS` defaults, there are many other factors to consider when deciding whether to set a CPU limit, such as prioritizing utilization of idle resources by avoiding limits vs prioritizing predictable latency by setting limits.

The worst behaviors from a mismatch between `GOMAXPROCS` and effective CPU limits occur when `GOMAXPROCS` is significantly higher than the effective CPU limit.
For example, a small container receiving 2 CPUs running on a 128 core machine.
These are the cases where it is most valuable to consider setting an explicit CPU limit, or, alternatively, explicitly setting `GOMAXPROCS`.

## Conclusion

Go 1.25 provides more sensible default behavior for many container workloads by setting `GOMAXPROCS` based on container CPU limits.
Doing so avoids throttling that can impact tail latency, improves efficiency, and generally tries to ensure Go is production-ready out-of-the-box.
You can get the new defaults simply by setting the Go version to 1.25.0 or higher in your `go.mod`.

Thanks to everyone in the community that contributed to the [long](/issue/33803) [discussions](/issue/73193) that made this a reality, and in particular to feedback from the maintainers of [`go.uber.org/automaxprocs`](https://pkg.go.dev/go.uber.org/automaxprocs) from Uber, which has long provided similar behavior to its users.
