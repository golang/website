---
title: Eight years of Go
date: 2017-11-10
by:
- Steve Francia
tags:
- community
- birthday
summary: Happy 8th birthday, Go!
---


Today we celebrate 8 years since Go was released as an open source project.
During [Go’s 4th anniversary](https://blog.golang.org/4years), Andrew
finished the post with “Here's to four more years!”. Now that we have reached
that milestone, I cannot help but reflect on how much the project and
ecosystem has grown since then. In our post 4 years ago we included a chart
demonstrating Go's rising popularity on Google Trends with the search term
"golang". Today, we’re including an updated chart. In this relative scale of
popularity, what was 100 four years ago is now a mere 17. Go’s popularity has
increased exponentially over the last 8 years and continues to grow.

{{image "8years/image1.png"}}

Source: [trends.google.com](https://trends.google.com/trends/explore?date=2009-10-01%202017-10-30&q=golang&hl=en-US)

## Developers love Go

Go has been embraced by developers all over the world with approximately one
million users worldwide. In the [freshly published 2017 Octoverse](https://octoverse.github.com/)
by GitHub, **Go has become the #9 most popular language**, surpassing C.
**Go is the fastest growing language on GitHub in 2017** in the top 10 with
**52% growth over the previous year**. In growth, Go swapped places with
Javascript, which fell to the second spot with 44%.

{{image "8years/image2.png"}}

Source: [octoverse.github.com](https://octoverse.github.com/)

In [Stack Overflow's 2017 developer survey](https://insights.stackoverflow.com/survey/2017#most-loved-dreaded-and-wanted)
, Go was the only language that was both on the **top 5 most loved and top 5 most wanted** languages.
People who use Go, love it, and the people who aren’t using Go, want to be.

{{image "8years/image3.png"}}
{{image "8years/image4.png"}}

Source: [insights.stackoverflow.com/survey/2017](https://insights.stackoverflow.com/survey/2017#most-loved-dreaded-and-wanted)

## Go: The language of Cloud Infrastructure

In 2014, analyst Donnie Berkholz called Go
[the emerging language of cloud infrastructure](http://redmonk.com/dberkholz/2014/03/18/go-the-emerging-language-of-cloud-infrastructure/).
**By 2017, Go has emerged as the language of cloud infrastructure**.
Today, **every single cloud company has critical components of their cloud infrastructure implemented in Go**
including Google Cloud, AWS, Microsoft Azure, Digital Ocean, Heroku and many others. Go
is a key part of cloud companies like Alibaba, Cloudflare, and Dropbox. Go is
a critical part of open infrastructure including Kubernetes, Cloud Foundry,
Openshift, NATS, Docker, Istio, Etcd, Consul, Juju and many more. Companies
are increasingly choosing Go to build cloud infrastructure solutions.

## Go’s Great Community

It may be hard to imagine that only four years ago the Go community was
transitioning from online-only to include in-person community with its first
conference. Now the Go community has had over 30 conferences all around the
world with hundreds of presentations and tens of thousands of attendees.
There are hundreds of Go meetups meeting monthly covering much of the globe.
Wherever you live, you are likely to find a Go meetup nearby.

Two different organizations have been established to help with inclusivity in
the Go community, Go Bridge and Women Who Go; the latter has grown to over 25
chapters. Both have been instrumental in offering free trainings. In 2017
alone over 50 scholarships to conferences have been given through efforts of
Go Bridge and Women Who Go.

This year we had two significant firsts for the Go project. We had our first
[contributor summit](https://blog.golang.org/contributors-summit) where
people from across the Go community came together to
discuss the needs and future of the Go project. Shortly after, we had the
first [Go contributor workshop](https://blog.golang.org/contributor-workshop)
where hundreds of people came to make their first Go contribution.

{{image "8years/photo.jpg"}}

Photo by Sameer Ajmani

## Go’s impact on open source

Go has become a major force in the world of open source powering some of the
most popular projects and enabling innovations across many industries. Find
thousands of additional applications and libraries at [awesome-go](https://github.com/avelino/awesome-go). Here are
just a handful of the most popular:

  - [Moby](https://mobyproject.org/) (formerly Docker) is a tool for packaging
    and running applications in lightweight containers.
    Its creator Solomon Hykes cited Go's standard library,
    concurrency primitives, and ease of deployment as key factors,
    and said "To put it simply, if Docker had not been written in Go,
    it would not have been as successful."

  - [Kubernetes](https://kubernetes.io/) is a system for automating deployment,
    scaling and management of containerized applications.
    Initially designed by Google and used in the Google cloud,
    Kubernetes now is a critical part of every major cloud offering.

  - [Hugo](https://gohugo.io/) is now the most popular open-source static website engine.
    With its amazing speed and flexibility, Hugo makes building websites fun again.
    According to [w3techs](https://w3techs.com/technologies/overview/content_management/all),
    Hugo now has nearly 3x the usage of Jekyll, the former leader.

  - [Prometheus](https://prometheus.io/) is an open source monitoring solution
    and time series database that powers metrics and alerting designed to be
    the system you go to during an outage to allow you to quickly diagnose problems.

  - [Grafana](https://grafana.com/) is an open source,
    feature-rich metrics dashboard and graph editor for Graphite,
    Elasticsearch, OpenTSDB, Prometheus and InfluxDB.

  - [Lantern](https://getlantern.org/) delivers fast, reliable and secure access to blocked websites and apps.

  - [Syncthing](https://syncthing.net/) is an open-source cross platform
    peer-to-peer continuous file synchronization application

  - [Keybase](https://keybase.io/) is a new and free security app for mobile
    phones and computers.
    Think of it as an open source Dropbox & Slack with end-to-end encryption
    public-key cryptography.

  - [Fzf](https://github.com/junegunn/fzf) is an interactive Unix filter
    for command-line that can be used with any list;
    files, command history, processes, hostnames,
    bookmarks, git commits, etc.
    Fzf supports Unix, macOS and has beta support for Windows.
    It also can operate as a vim plugin.

Many of these authors have said that their projects would not exist without
Go. Some like Kubernetes and Docker created entirely new solutions. Others
like Hugo, Syncthing and Fzf created more refined experiences where many
solutions already existed. The popularity of these applications alone is
proof that Go is an ideal language for a broad set of use cases.

## Thank You

This is the eighth time we have had the pleasure of writing a birthday blog
post for Go and we continue to be overwhelmed by and grateful for the
enthusiasm and support of the Go community.

Since Go was first open sourced we have had 10 releases of the language,
libraries and tooling with more than 1680 contributors making over 50,000
commits to the project's 34 repositories; More than double the number of
contributors and nearly double the number of commits from only [two years ago](https://blog.golang.org/6years).
This year we announced that we have begun planning [Go 2](https://blog.golang.org/toward-go2), our first major
revision of the language and tooling.

The Go team would like to thank everyone who has contributed to the project,
whether you participate by contributing changes, reporting bugs, sharing your
expertise in design discussions, writing blog posts or books, running events,
attending or speaking at events, helping others learn or improve, open
sourcing Go packages you wrote, contributing artwork, introducing Go to
someone, or being part of the Go community. Without you, Go would not be as
complete, useful, or successful as it is today.

Thank you, and here’s to eight more years!
