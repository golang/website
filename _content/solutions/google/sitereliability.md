---
title: "Actuating Google Production: How Google’s Site Reliability Engineering Team Uses Go"
company: Google Site Reliability Engineering (SRE)
logoSrc: sitereliability.svg
logoSrcDark: sitereliability.svg
heroImgSrc: go_sitereliability_case_study.png
series: Case Studies
quote: |
  Google’s Site Reliability Engineering team has a mission to protect, provide for, and progress the software and systems behind all of Google’s public services — Google Search, Ads, Gmail, Android, YouTube, and App Engine, to name just a few — with an ever-watchful eye on their availability, latency, performance, and capacity.

  They shared their experience building core production management systems with Go, coming from experience with Python and C++.
authors:
  - Pierre Palatin, Site Reliability Engineer
---

Google runs a small number of very large services. Those services are powered
by a global infrastructure covering everything a developer needs: storage
systems, load balancers, network, logging, monitoring, and much more.
Nevertheless, it is not a static system—it cannot be. Architecture evolves,
new products and ideas are created, new versions must be rolled out, configs
pushed, database schema updated, and more. We end up deploying changes to our
systems dozens of times per second.

Because of this scale and critical need for reliability, Google pioneered Site
Reliability Engineering (SRE), a role that many other companies have since adopted.
“SRE is what you get when you treat operations as if it’s a software problem.
Our mission is to protect, provide for, and progress the software and systems
behind all of Google’s public services with an ever-watchful eye on their
availability, latency, performance, and capacity.”
— [Site Reliability Engineering (SRE)](https://sre.google/).

{{backgroundquote `
  quote: |
    Go promised a sweet spot between performance and readability that neither of
    the other languages [Python and C++] were able to offer.
`}}

In 2013-2014, Google’s SRE team realized that our approach to production
management was not cutting it anymore in many ways. We had advanced far beyond
shell scripts, but our scale had so many moving pieces and complexities that a
new approach was needed. We determined that we needed to move toward a
declarative model of our production, called "Prodspec", driving a dedicated
control plane, called "Annealing".

When we started those projects, Go was just becoming a viable option for
critical services at Google. Most engineers were more familiar with Python
and C++, either of which would have been valid choices. Nevertheless, Go
captured our interest. The appeal of novelty was certainly a factor of
course. But, more importantly, Go promised a sweet spot between performance
and readability that neither of the other languages were able to offer. We
started a small experiment with Go for some initial parts of Annealing and
Prodspec. As the projects progressed, those initial parts written in Go found
themselves at the core. We were happy with Go—its simplicity grew on us, the
performance was there, and concurrency primitives would have been hard to
replace.

{{backgroundquote `
  quote: |
    Now the majority of Google production is managed and maintained by our systems
    written in Go.
`}}

At no point was there ever a mandate or requirement to use Go, but we had no
desire to return to Python or C++. Go grew organically in Annealing and
Prodspec. It was the right choice, and thus is now our language of choice.
Now the majority of Google production is managed and maintained by our systems
written in Go.

The power of having a simple language in those projects is hard to overstate.
There have been cases where some feature was indeed missing, such as the
ability to enforce in the code that some complex structure should not be
mutated. But for each one of those cases, there have undoubtedly been tens or
hundred of cases where the simplicity helped.

{{backgroundquote `
  quote: |
    Go’s simplicity means that the code is easy to follow, whether it is to spot
    bugs during review or when trying to determine exactly what happened during a
    service disruption.
`}}

For example, Annealing impacts a wide variety of teams and services meaning
that we relied heavily on contributions across the company. The simplicity of
Go made it possible for people outside our team to see why some part or another
was not working for them, and often provide fixes or features themselves. This
allowed us to quickly grow.

Prodspec and Annealing are in charge of some quite critical components. Go’s
simplicity means that the code is easy to follow, whether it is to spot bugs
during review or when trying to determine exactly what happened during a
service disruption.

Go performance and concurrency support have also been key for our work. As our
model of production is declarative, we tend to manipulate a lot of structured
data, which describes what production is and what it should be. We have large
services so the data can grow large, often making purely sequential processing
not efficient enough.

We are manipulating this data in many ways and many places. It is not a matter
of having a smart person come up with a parallel version of our algorithm. It
is a matter of casual parallelism, finding the next bottleneck and
parallelising that code section. And Go enables exactly that.

As a result of our success with Go, we now use Go for every new development for
Prodspec and Annealing.

In addition to the Site Reliability Engineering team, engineering teams across
Google have adopted Go in their development process. Read about how the
[Core Data Solutions](/solutions/google/coredata/),
[Firebase Hosting](/solutions/google/firebase/), and
[Chrome](/solutions/google/chrome/) teams use Go to build fast, reliable,
and efficient software at scale.