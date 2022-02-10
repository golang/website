---
title: 'Using Go at Google'
date: 2020-08-27
company: Google
logoSrc: google.svg
logoSrcDark: google.svg
heroImgSrc: go_core_data_case_study.png
carouselImgSrc: go_google_case_study_carousel.png
series: Case Studies
type: solutions
description: |-
  Google is a technology company whose mission is to organize the world’s
  information and make it universally accessible and useful.

  Go was created at Google in 2007 to improve programming productivity in an
  era of multi-core networked machines and large codebases. Today, over 10
  years since its public announcement in 2009, Go's use inside Google has grown
  tremendously.
quote: Go was created at Google in 2007, and since then, engineering teams
  across Google have adopted Go to build products and services at massive scale.

---

{{pullquote `
  author: Rob Pike
  quote: |
    Go started in September 2007 when Robert Griesemer, Ken Thompson, and I began
    discussing a new language to address the engineering challenges we and our
    colleagues at Google were facing in our daily work.

    When we first released Go to the public in November 2009, we didn’t know if the
    language would be widely adopted or if it might influence future languages.
    Looking back from 2020, Go has succeeded in both ways: it is widely used both
    inside and outside Google, and its approaches to network concurrency and
    software engineering have had a noticeable effect on other languages and their
    tools.

    Go has turned out to have a much broader reach than we had ever expected. Its
    growth in the industry has been phenomenal, and it has powered many projects at
    Google.
`}}

The following stories are a small sample of the many ways that Go is used at Google.

### How Google's Core Data Solutions Team Uses Go

Google's mission is “to organize the world's information and make it universally
accessible and useful.”  One of the teams responsible for organizing that
information is Google’s Core Data Solutions team. The team, among other things,
maintains services to index web pages across the globe. These web indexing
services help support products like Google Search by keeping search results
updated and comprehensive, and they’re written in Go.

[Learn more](/solutions/google/coredata/)

---

### Chrome Content Optimization Service Runs on Go

When the product Chrome comes to mind, you probably think solely of the user-installed browser. But behind the scenes, Chrome has an extensive fleet of backends. Among these is the Chrome Optimization Guide service. This service forms an important basis for Chrome’s user experience strategy, operating in the critical path for users, and is implemented in Go.

[Learn more](/solutions/google/chrome/)

---

### How the Firebase Hosting Team Scaled With Go

The Firebase Hosting team provides static web hosting services for Google Cloud customers. They provide a static web host that sits behind a global content delivery network, and offer users tools that are easy to use. The team also develops features that range from uploading site files to registering domains to tracking usage.

[Learn more](/solutions/google/firebase/)

---

### Actuating Google Production: How Google’s Site Reliability Engineering Team Uses Go

Google runs a small number of very large services. Those services are powered by a global infrastructure covering everything one needs: storage systems, load balancers, network, logging, monitoring, and many more. Nevertheless, it is not a static system - it cannot be. Architecture evolves, new products and ideas are created, new versions must be rolled out, configs pushed, database schema updated, and more. We end up deploying changes to our systems dozens of times per second.

[Learn more](/solutions/google/sitereliability/)