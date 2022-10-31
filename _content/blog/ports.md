---
title: Go on ARM and Beyond
date: 2020-12-17
by:
- Russ Cox
summary: Go's support for ARM64 and other architectures
---


The industry is abuzz about non-x86 processors recently,
so we thought it would be worth a brief post about Go’s support for them.

It has always been important to us for Go to be portable,
not overfitting to any particular operating system or architecture.
The [initial open source release of Go](https://opensource.googleblog.com/2009/11/hey-ho-lets-go.html)
included support for two operating systems (Linux and Mac OS X) and three
architectures (64-bit x86,
32-bit x86, and 32-bit ARM).

Over the years, we’ve added support for many more operating systems and architecture combinations:

- Go 1 (March 2012) supported the original systems as well as FreeBSD,
  NetBSD, and OpenBSD on 64-bit and 32-bit x86,
  and Plan 9 on 32-bit x86.
- Go 1.3 (June 2014) added support for Solaris on 64-bit x86.
- Go 1.4 (December 2014) added support for Android on 32-bit ARM and Plan 9 on 64-bit x86.
- Go 1.5 (August 2015) added support for Linux on 64-bit ARM and 64-bit PowerPC,
  as well as iOS on 32-bit and 64-bit ARM.
- Go 1.6 (February 2016) added support for Linux on 64-bit MIPS,
  as well as Android on 32-bit x86.
  It also added an official binary download for Linux on 32-bit ARM,
  primarily for Raspberry Pi systems.
- Go 1.7 (August 2016) added support for Linux on z Systems (S390x) and Plan 9 on 32-bit ARM.
- Go 1.8 (February 2017) added support for Linux on 32-bit MIPS,
  and it added official binary downloads for Linux on 64-bit PowerPC and z Systems.
- Go 1.9 (August 2017) added official binary downloads for Linux on 64-bit ARM.
- Go 1.12 (February 2018) added support for Windows 10 IoT Core on 32-bit ARM,
  such as the Raspberry Pi 3.
  It also added support for AIX on 64-bit PowerPC.
- Go 1.14 (February 2019) added support for Linux on 64-bit RISC-V.

Although the x86-64 port got most of the attention in the early days of Go,
today all our target architectures are well supported by our [SSA-based compiler back end](https://www.youtube.com/watch?v=uTMvKVma5ms)
and produce excellent code.
We’ve been helped along the way by many contributors,
including engineers from Amazon, ARM, Atos,
IBM, Intel, and MIPS.

Go supports cross-compiling for all these systems out of the box with minimal effort.
For example, to build an app for 32-bit x86-based Windows from a 64-bit Linux system:

	GOARCH=386 GOOS=windows go build myapp  # writes myapp.exe

In the past year, several major vendors have made announcements of new ARM64
hardware for servers,
laptops and developer machines.
Go was well-positioned for this. For years now,
Go has been powering Docker, Kubernetes, and the rest of the Go ecosystem
on ARM64 Linux servers,
as well as mobile apps on ARM64 Android and iOS devices.

Since Apple’s announcement of the Mac transitioning to Apple Silicon this summer,
Apple and Google have been working together to ensure that Go and the broader
Go ecosystem work well on them,
both running Go x86 binaries under Rosetta 2 and running native Go ARM64 binaries.
Earlier this week, we released the first Go 1.16 beta,
which includes native support for Macs using the M1 chip.
You can download and try the Go 1.16 beta for M1 Macs and all your other
systems on [the Go download page](/dl/#go1.16beta1).
(Of course, this is a beta release and, like all betas,
it certainly has bugs we don’t know about.
If you run into any problems, please report them at [golang.org/issue/new](/issue/new).)

It’s always nice to use the same CPU architecture for local development as in production,
to remove one variation between the two environments.
If you deploy to ARM64 production servers,
Go makes it easy to develop on ARM64 Linux and Mac systems too.
But of course, it’s still as easy as ever to work on one system and cross-compile
for deployment to another,
whether you’re working on an x86 system and deploying to ARM,
working on Windows and deploying to Linux,
or some other combination.

The next target we’d like to add support for is ARM64 Windows 10 systems.
If you have expertise and would like to help,
we’re coordinating work on [golang.org/issue/36439](https://github.com/golang/go/issues/36439).

