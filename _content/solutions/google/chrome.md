---
title: "Chrome Content Optimization Service Runs on Go"
company: Chrome
logoSrc: chrome.svg
logoSrcDark: chrome.svg
heroImgSrc: go_chrome_case_study.png
series: Case Studies
quote: |
  Google Chrome is a more simple, secure, and faster web browser than ever,
  with Google's smarts built-in.

  In this case study, the Chrome Optimization Guide
  team shared how they experimented with Go, ramped up quickly, and their plans to
  use Go going forward.
---

When the product Chrome comes to mind, you probably think solely of the
user-installed browser. But behind the scenes, Chrome has an extensive fleet of
backends. Among these is the Chrome Optimization Guide service. This service
forms an important basis for Chrome's user experience strategy, operating in the
critical path for users, and is implemented in Go.

The Chrome Optimization Guide service is designed to bring the power of Google
to Chrome by providing hints to the installed browser about what optimizations
may be performed on a page load, as well as when they can be applied most
effectively. It comprises a conjunction of real-time servers and batch logs
analysis.

All Lite mode users of Chrome receive data via the service through the following
mechanisms: a data blob push that provides hints for well-known sites in their
geography, a check-in to Google servers to retrieve hints for hosts that the
specific user visits often, and on demand for page loads for which a hint is not
already on the device. Were the Chrome Optimization Guide service to suddenly
disappear, users might notice a dramatic change in the speed of their page loads
and the amount of data consumed while browsing the web.

{{backgroundquote `
  author: Sophie Chang
  title: Software Engineer
  quote: |
    Given that Go was a success for us, we plan to continue to use
    it where appropriate
`}}

When the Chrome engineering team started building the service, only a few
members had comfort with Go. Most of the team was more familiar with C++, but
they found the complex boilerplate required to stand up a C++ server to be too
much. The team shared that “[they] were pretty motivated to learn Go due to its
simplicity, fast ramp-up, and ecosystem.” and that “[their] sense of adventure
was rewarded.” Millions of users rely on this service to make their Chrome
experience better, and choosing Go was no small decision. After their experience
so far, the team also shared that “given that Go was a success for us, we plan
to continue to use it where appropriate.”

In addition to the Chrome Optimization Guide team, engineering teams across
Google have adopted Go in their development process. Read about how the [Core
Data Solutions](/solutions/google/coredata/) and [Firebase
Hosting](/solutions/google/firebase/) teams use Go to build fast, reliable,
and efficient software at scale.

*Editorial note: The Go team would like to thank Sophie Chang for her
contributions to this story.*
