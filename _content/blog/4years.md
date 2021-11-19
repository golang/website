---
title: Four years of Go
date: 2013-11-10
by:
- Andrew Gerrand
tags:
- community
- birthday
summary: Happy 4th birthday, Go!
---


Today marks the fourth anniversary of Go as an open source project.

{{image "4years/4years-gopher.png"}}

Rather than talk about our technical progress (there'll be much to talk about
when we release Go 1.2 in a couple of weeks) we thought we would instead take
this occasion to look at how the Go community has grown.

Let's start with a chart:

{{image "4years/4years-graph.png"}}

This chart shows the growth of Google searches for the term
"[golang](http://www.google.com/trends/explore?hl=en-US#q=golang&date=10/2009+50m&cmpt=q)"
over the past four years.
Notice the knee in the curve around March 2012, when Go 1.0 was released.
If these searches are a decent proxy for interest, then it's clear that
interest in Go has grown remarkably since launch, and particularly so in the
last 2 years.

But where is the interest coming from?

The open source community has embraced Go,
with our community wiki listing [hundreds of Go projects](/wiki/Projects). Some popular ones:

  - [Docker](http://docker.io) is a tool for packaging and running applications
    in lightweight containers.
    Docker makes it easy to isolate, package,
    and deploy applications, and is beloved by system administrators.
    Its creator Solomon Hykes cited Go's standard library,
    concurrency primitives, and ease of deployment as key factors,
    and said "To put it simply, if Docker had not been written in Go,
    it would not have been as successful."

  - [Packer](http://packer.io) is a tool for automating the creation of
    machine images for deployment to virtual machines or cloud services.
    Its author, Mitchell Hashimoto, is now working on another Go project,
    [serf](http://www.serfdom.io/), a decentralized discovery service.
    Like Docker, these projects help with management of large-scale,
    cluster-based services.

  - [Bitly](http://bit.ly)'s [NSQ](http://bitly.github.io/nsq/) is a realtime
    distributed messaging platform designed for fault-tolerance and high-availability,
    and is used in production at bitly and a bunch of other companies.

  - [Canonical](http://canonical.com/)'s [JuJu](https://juju.ubuntu.com/)
    infrastructure automation system was rewritten in Go.
    Project lead Gustavo Niemeyer said "It's not a single aspect of Go that
    makes it a compelling choice,
    but rather the careful organization of well-crafted small pieces."

  - The [raft](https://github.com/goraft/raft) package provides an implementation
    of the [Raft](https://ramcloud.stanford.edu/wiki/download/attachments/11370504/raft.pdf)
    distributed consensus protocol.
    It is the basis of Go projects like [etcd](https://github.com/coreos/etcd)
    and [SkyDNS](https://github.com/skynetservices/skydns).

  - Other popular projects include [biogo](https://github.com/biogo/biogo),
    the [Gorilla Web Toolkit](http://www.gorillatoolkit.org/),
    [groupcache](https://github.com/golang/groupcache),
    Mozilla's [heka](https://github.com/mozilla-services/heka),
    the [kv](https://github.com/cznic/kv) and [ql](https://github.com/cznic/ql)
    lightweight storage systems,
    and the [Sky](http://skydb.io/) behavioral database.

But this is just the tip of the iceberg. The number of high-quality open
source Go projects is phenomenal.
Prolific Go hacker [Keith Rarick](http://xph.us/software/) put it well:
"The state of the Go ecosystem after only four years is astounding.
Compare Go in 2013 to Python in 1995 or Java in 1999. Or C++ in 1987!"

Businesses are enjoying Go, too. The [Go Users wiki page](/wiki/GoUsers)
lists dozens of success stories (and if you use Go,
please add yourself to it). Some examples:

  - [CloudFlare](https://blog.cloudflare.com/go-at-cloudflare) built their
    distributed DNS service entirely with Go,
    and are in the process of migrating their gigabytes-per-minute logging infrastructure to the language.
    Programmer John Graham-Cumming said "We've found Go to be the perfect match for our needs:
    the combination of familiar syntax, a powerful type system,
    a strong network library and built-in concurrency means that more and more
    projects are being built here in Go."

  - [SoundCloud](http://soundcloud.com) is an audio distribution service
    that has "dozens of [systems in Go](http://backstage.soundcloud.com/2012/07/go-at-soundcloud/),
    touching almost every part of the site, and in many cases powering features
    from top to bottom." Engineer Peter Bourgon said "Go demonstrates that the
    cruft that burdens other languages and ecosystems—stuff that developers
    have learned to deal with,
    often in anger—is simply not a necessary part of modern programming.
    With Go, I have a straightforward and non-adversarial relationship with my tools,
    from development to production."

  - The [ngrok](https://ngrok.com/) service allows web developers to provide
    remote access to their development environments.
    Its author Alan Shreve said that "ngrok's success as a project is due in
    no small part to choosing Go as the implementation language," citing Go's HTTP libraries,
    efficiency, cross-platform compatibility,
    and ease of deployment as the major benefits.

  - [Poptip](http://poptip.com) provides social analytics services,
    and product engineer Andy Bonventre said "What started as an experiment
    in writing a single service in Go turned into moving almost our entire infrastructure over to it.
    What I love about Go the most is not necessarily the features of the language,
    but the focus on tooling, testing, and other elements that make writing
    large applications much more manageable."

  - Music collaboration startup [Splice](http://splice.com) chose to build
    their service with Go.
    Co-founder Matt Aimonetti said "We seriously studied and considered many
    programming languages,
    but Go's simplicity, efficiency, philosophy and community won us over."

  - And, of course, engineering teams across Google are moving to Go.
    Engineer Matt Welsh recently [shared his experience](http://matt-welsh.blogspot.com.au/2013/08/rewriting-large-production-system-in-go.html)
    rewriting a large production service in Go.
    Other notable public examples include YouTube's [vitess project](https://github.com/youtube/vitess)
    and [dl.google.com](/talks/2013/oscon-dl.slide).
    We hope to share more stories like these soon.

In September 2012, [Apcera](http://apcera.com/) CEO Derek Collison [predicted](https://twitter.com/derekcollison/status/245522124666716160)
that "Go will become the dominant language for systems work in [Infastructure-as-a-Service],
Orchestration, and [Platform-as-a-Service] in 24 months." Looking at the list above,
it's easy to believe that prediction.

So how can you get involved? Whether you're a seasoned Go programmer or just Go-curious,
there are many ways to get started in the Go community:

  - [Join your nearest Go User Group](https://blog.golang.org/getthee-to-go-meetup),
    where your local gophers meet to share their knowledge and experience.
    These groups are popping up all over the world.
    I have personally spoken at Go groups in Amsterdam,
    Berlin, Gothenburg, London, Moscow, Munich,
    New York City, Paris, San Francisco, Seoul,
    Stockholm, Sydney, Tokyo, and Warsaw;
    but there are [many more](/wiki/GoUserGroups)!

  - Create or contribute to an open source Go project (or [to Go itself](/doc/contribute.html)).
    (And if you're building something, we'd love to hear from you on the [Go mailing list](http://groups.google.com/group/golang-nuts).)

  - If you're in Europe in February 2014, come along to the [Go Devroom](https://code.google.com/p/go-wiki/wiki/Fosdem2014)
    at [FOSDEM 2014](https://fosdem.org/2014/).

  - Attend [GopherCon](http://gophercon.com),
    the first major Go conference, in Denver in April 2014.
    The event is organized by the [Gopher Academy](http://www.gopheracademy.com),
    who also run a [Go job board](http://www.gopheracademy.com/jobs).

The Go team has been amazed by the growth of the Go community over the past
four years. We are thrilled to see so many great things being built with Go,
and deeply grateful to work with our wonderful and dedicated contributors.
Thank you, everyone.

Here's to four more years!
