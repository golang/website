---
title: Thirteen Years of Go
date: 2022-11-10
by:
- Russ Cox, for the Go team
summary: Happy Birthday, Go!
---

<img src="../doc/gopher/gopherbelly300.jpg" height="300" width="300" align="right" style="margin: 0 0 1em 1em;">

Today we celebrate the thirteenth birthday of the Go open source release.
[The Gopher](/doc/gopher) is a teenager!

It's been an eventful year for Go.
The most significant event was the release of
[Go 1.18 in March](/blog/go1.18),
which brought many improvements but most notably
Go workspaces, fuzzing, and generics.

Workspaces make it easy to work on multiple modules simultaneously,
which is most helpful when you are maintaining a set of related modules with
module dependencies between them.
To learn about workspaces, see Beth Brown's blog post
“[Get familiar with workspaces](/blog/get-familiar-with-workspaces)”
and the [workspace reference](/ref/mod#workspaces).

Fuzzing is a new feature of `go` `test` that
helps you find inputs that your code doesn't handle properly:
you define a fuzz test that should pass for any input at all,
and then fuzzing tries different random inputs, guided by code coverage,
to try to make the fuzz test fail.
Fuzzing is particularly useful when developing code that must be
robust against arbitrary (even attacker-controlled) inputs.
To learn more about fuzzing, see the tutorial
“[Getting started with fuzzing](/doc/tutorial/fuzz)”
and the [fuzzing reference](/security/fuzz/),
and keep an eye out for Katie Hockman's GopherCon 2022 talk
“Fuzz Testing Made Easy”,
which should be online soon.

Generics, quite possibly Go's most requested feature,
adds parametric polymorphism to Go, to allow writing
code that works with a variety of different types but is still
statically checked at compile time.
To learn more about generics, see the tutorial
“[Getting started with generics](/doc/tutorial/generics)”.
For more detail see
the blog posts
“[An Introduction to Generics](/blog/intro-generics)”
and
“[When to Use Generics](/blog/when-generics)”,
or the talks
“[Using Generics in Go](https://www.youtube.com/watch?v=nr8EpUO9jhw)”
from Go Day on Google Open Source Live 2021,
and
“[Generics!](https://www.youtube.com/watch?v=Pa_e9EeCdy8)” from GopherCon 2021,
by Robert Griesemer and Ian Lance Taylor.

Compared to Go 1.18, the [Go 1.19 release in August](/blog/go1.19) was a relatively quiet one:
it focused on refining and improving the features that Go 1.18 introduced
as well as internal stability improvements and optimizations.
One visible change in Go 1.19 was the addition of
support for [links, lists, and headings in Go doc comments](/doc/comment).
Another was the addition of a [soft memory limit](/doc/go1.19#runtime)
for the garbage collector, which is particularly useful in container workloads.
For more about recent garbage collector improvements,
see Michael Knyszek's blog post “[Go runtime: 4 years later](/blog/go119runtime)”,
his talk “[Respecting Memory Limits in Go](https://www.youtube.com/watch?v=07wduWyWx8M&list=PLtoVuM73AmsJjj5tnZ7BodjN_zIvpULSx)”,
and the new “[Guide to the Go Garbage Collector](/doc/gc-guide)”.

We've continued to work on making Go development scale gracefully to ever larger code bases,
especially in our work on VS Code Go and the Gopls language server.
This year, Gopls releases focused on improving stability and performance,
while delivering support for generics as well as new analyses and code lenses.
If you aren't using VS Code Go or Gopls yet, give them a try.
See Suzy Mueller's talk
“[Building Better Projects with the Go Editor](https://www.youtube.com/watch?v=jMyzsp2E_0U)”
for an overview.
And as a bonus,
[Debugging Go in VS Code](/s/vscode-go-debug)
got more reliable and powerful with Delve's native
[Debug Adapter Protocol](https://microsoft.github.io/debug-adapter-protocol/) support.
Try Suzy's “[Debugging Treasure Hunt](https://www.youtube.com/watch?v=ZPIPPRjwg7Q)”!

Another part of development scale is the number of dependencies in a project.
A month or so after Go's 12th birthday,
the [Log4shell vulnerability](https://en.wikipedia.org/wiki/Log4Shell) served
as a wake-up call for the industry
about the importance of supply chain security.
Go's module system was designed specifically for this purpose,
to help you understand and track your dependencies,
identify which specific ones you are using,
and determine whether any of them have known vulnerabilities.
Filippo Valsorda's blog post
“[How Go Mitigates Supply Chain Attacks](/blog/supply-chain)”
gives an overview of our approach.
In September, we previewed
Go's approach to vulnerability management
in Julie Qiu's blog post “[Vulnerability Management for Go](/blog/vuln)”.
The core of that work is a new, curated vulnerability database
and a new [govulncheck command](https://pkg.go.dev/golang.org/x/vuln/cmd/govulncheck),
which uses advanced static analysis to eliminate most of the false positives
that would result from using module requirements alone.

Part of our effort to understand Go users is our annual end-of-year Go survey.
This year, our user experience researchers added a lightweight mid-year Go survey as well.
We aim to gather enough responses to be statistically significant
without being a burden on the Go community as a whole.
For the results, see Alice Merrick's blog post
“[Go Developer Survey 2021 Results](/blog/survey2021-results)”
and Todd Kulesza's post
“[Go Developer Survey 2022 Q2 Results](/blog/survey2022-q2-results)”.

As the world starts traveling more,
we've also been happy to meet many of you in person at Go conferences in 2022,
particularly at GopherCon Europe in Berlin in July and at GopherCon in Chicago in October.
Last week we held our annual virtual event,
[Go Day on Google Open Source Live](https://opensourcelive.withgoogle.com/events/go-day-2022).
Here are some of the talks we've given at those events:

 - “[How Go Became its Best Self](https://www.youtube.com/watch?v=vQm_whJZelc)”,
   by Cameron Balahan, at GopherCon Europe.
 - “[Go team Q&A](https://www.youtube.com/watch?v=KbOTTU9yEpI)”,
   with Cameron Balahan, Michael Knyszek, and Than McIntosh, at GopherCon Europe.
 - “[Compatibility: How Go Programs Keep Working](https://www.youtube.com/watch?v=v24wrd3RwGo)”,
   by Russ Cox at GopherCon.
 - “[A Holistic Go Experience](https://www.gophercon.com/agenda/session/998660)”,
   by Cameron Balahan at GopherCon (video not yet posted)
 - “[Structured Logging for Go](https://opensourcelive.withgoogle.com/events/go-day-2022/watch?talk=talk2)”,
   by Jonathan Amsterdam at Go Day on Google Open Source Live
 - “[Writing your Applications Faster and More Securely with Go](https://opensourcelive.withgoogle.com/events/go-day-2022/watch?talk=talk3)”,
   by Cody Oss at Go Day on Google Open Source Live
 - “[Respecting Memory Limits in Go](https://opensourcelive.withgoogle.com/events/go-day-2022/watch?talk=talk4),
   by Michael Knyszek at Go Day on Google Open Source Live

One other milestone for this year was the publication of
“[The Go Programming Language and Environment](https://cacm.acm.org/magazines/2022/5/260357-the-go-programming-language-and-environment/fulltext)”,
by Russ Cox, Robert Griesemer, Rob Pike, Ian Lance Taylor, and Ken Thompson,
in _Communications of the ACM_.
The article, by the original designers and implementers of Go,
explains what we believe makes Go so popular and productive.
In short, it is that Go effort focuses on delivering a full development environment
targeting the entire software development process,
with a focus on scaling both to large software engineering efforts
and large deployments.

In Go's 14th year, we'll keep working to make Go the best environment
for software engineering at scale.
We plan to focus particularly on supply chain security, improved compatibility,
and structured logging, all of which have been linked already in this post.
And there will be plenty of other improvements as well,
including profile-guided optimization.

## Thank You!

Go has always been far more than what the Go team at Google does.
Thanks to all of you—our contributors and everyone in the Go community—for
your help making Go the successful programming environment that it is today.
We wish you all the best in the coming year.

