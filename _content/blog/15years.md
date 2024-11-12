---
title: Go Turns 15
date: 2024-11-11
by:
- Austin Clements, for the Go team
tags:
- community
- birthday
summary: Happy 15th birthday, Go!
---

<div style="float:right; margin: 0 0 1em 1em; width: 245px">
<img src="/doc/gopher/fifteen.gif" height="245" width="245"><br/>
<i>Thanks to Renee French for drawing and animating the gopher doing the “15 puzzle”.</i>
</div>

Happy birthday, Go\!

On Sunday, we celebrated the 15th anniversary of [the Go open source
release](https://opensource.googleblog.com/2009/11/hey-ho-lets-go.html)\!

So much has changed since [Go's 10 year anniversary](/blog/10years),
both in Go and in the world. In other ways, so much has stayed the
same: Go remains committed to stability, safety, and supporting
software engineering and production at scale.

And Go is going strong\! Go's user base has more than tripled in the
past five years, making it one of the fastest growing languages. From
its beginnings just fifteen years ago, Go has become a top 10 language
and the language of the modern cloud.

With the releases of [Go 1.22 in February](/blog/go1.22) and [Go 1.23
in August](/blog/go1.23), it's been the year of `for` loops. Go 1.22
made variables introduced by `for` loops [scoped per
iteration](/blog/loopvar-preview), rather than per loop, addressing a
long-standing language "gotcha". Over ten years ago, leading up to the
release of Go 1, the Go team made decisions about several language
details; among them whether `for` loops should create a new loop
variable on each iteration. Amusingly, the discussion was quite brief
and distinctly unopinionated. Rob Pike closed it out in true Rob Pike
fashion with a single word: “stet” (leave it be). And so it was. While
seemingly insignificant at the time, years of production experience
highlighted the implications of this decision. But in that time, we
also built robust tools for understanding the effects of changes to
Go—notably, ecosystem-wide analysis and testing across the entire
Google codebase—and established processes for working with the
community and getting feedback. Following extensive testing, analysis,
and community discussion, we rolled out the change, accompanied by a
[hash bisection
tool](https://go.googlesource.com/proposal/+/master/design/60078-loopvar.md#transition-support-tooling)
to assist developers in pinpointing code affected by the change at
scale.

The change to `for` loops was part of a five year trajectory of
measured changes. It would not have been possible without [forward
language compatibility](/blog/toolchain) introduced in Go 1.21. This,
in turn, built upon the foundation laid by Go modules, which were
introduced in Go 1.14 four and a half years ago.

Go 1.23 further built on this change to introduce iterators and
[user-defined for-range loops](/blog/range-functions). Combined with
generics—introduced in Go 1.18, just two and a half years ago\!—this
creates a powerful and ergonomic foundation for custom collections and
many other programming patterns.

These releases have also brought many improvements in production
readiness, including [much-anticipated enhancements to the standard
library's HTTP router](/blog/routing-enhancements), a [total overhaul
of execution traces](/blog/execution-traces-2024), and [stronger
randomness](/blog/chacha8rand) for all Go applications. Additionally,
the introduction of our [first v2 standard library
package](/blog/randv2) establishes a template for future library
evolution and modernization.

Over the past year we've also been cautiously rolling out [opt-in
telemetry](/blog/gotelemetry) for Go tools. This system will give Go's
developers data to make better decisions, while remaining completely
[open](https://telemetry.go.dev/) and anonymous. Go telemetry first
appeared in
[gopls](https://github.com/golang/tools/blob/master/gopls/README.md),
the Go language server, where it has already led to a [litany of
improvements](https://github.com/golang/go/issues?q=is%3Aissue+label%3Agopls%2Ftelemetry-wins).
This effort paves the way to make programming in Go an even better
experience for everyone.

Looking forward, we're evolving Go to better leverage the capabilities
of current and future hardware. Hardware has changed a lot in the past
15 years. In order to ensure Go continues to support high-performance,
large-scale production workloads for the *next* 15 years, we need to
adapt to large multicores, advanced instruction sets, and the growing
importance of locality in increasingly non-uniform memory hierarchies.
Some of these improvements will be transparent. Go 1.24 will have a
totally new `map` implementation under the hood that's more efficient
on modern CPUs. And we're prototyping new garbage collection
algorithms designed around the capabilities and constraints of modern
hardware. Some improvements will be in the form of new APIs and tools
so Go developers can better leverage modern hardware. We're looking at
how to support the latest vector and matrix hardware instructions, and
multiple ways that applications can build in CPU and memory locality.
A core principle guiding our efforts is *composable optimization*: the
impact of an optimization on a codebase should be as localized as
possible, ensuring that the ease of development across the rest of the
codebase is not compromised.

We're continuing to ensure Go's standard library is safe by default
and safe by design. This includes ongoing efforts to incorporate
built-in, native support for FIPS-certified cryptography, so that FIPS
crypto will be just a flag flip away for applications that need it.
Furthermore, we're evolving Go's standard library packages where we
can and, following the example of `math/rand/v2`, considering where
new APIs can significantly enhance the ease of writing safe and secure
Go code.

We're working on making Go better for AI—and AI better for Go—by
enhancing Go's capabilities in AI infrastructure, applications, and
developer assistance. Go is a great language for building production
systems, and we want it to be a great language for [building
production *AI* systems](/blog/llmpowered), too.
Go's dependability as a language
for Cloud infrastructure has made it a natural choice for
[LLM](https://ollama.com/) [infrastructure](https://weaviate.io/)
[as](https://localai.io/) [well](https://zilliz.com/what-is-milvus).
For AI applications, we will continue building out first-class support
for Go in popular AI SDKs, including
[LangChainGo](https://pkg.go.dev/github.com/tmc/langchaingo) and
[Genkit](https://developers.googleblog.com/en/introducing-genkit-for-go-build-scalable-ai-powered-apps-in-go/).
And from its very beginning, Go aimed to improve the end-to-end
software engineering process, so naturally we're looking at bringing
the latest tools and techniques from AI to bear on reducing developer
toil, leaving more time for the fun stuff—like actually programming\!

## Thank you

All of this is only possible because of Go's incredible contributors
and thriving community. Fifteen years ago we could only dream of the
success that Go has become and the community that has developed around
Go. Thank you to everyone who has played a part, large and small. We
wish you all the best in the coming year.

