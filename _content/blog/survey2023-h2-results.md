---
title: Go Developer Survey 2023 H2 Results
date: 2023-12-05
by:
- Todd Kulesza
tags:
- survey
- community
- developer experience research
summary: What we learned from our 2023 H2 developer survey
---

<style type="text/css" scoped>
  .chart {
    margin-left: 1.5rem;
    margin-right: 1.5rem;
    width: 800px;
  }
  blockquote p {
    color: var(--color-text-subtle) !important;
  }

  .quote_source {
    font-style: italic;
  }

  @media (prefers-color-scheme: dark) {
    .chart {
      border-radius: 8px;
    }
  }
</style>

## Background

In August 2023, the Go team at Google conducted our bi-annual survey of Go
developers. We recruited participants via a public post on the Go blog and a
randomized prompt in VS Code, resulting in 4,005 responses. We primarily
focused survey questions around a few topics: general sentiment and feedback
about developing with Go, technology stacks used alongside Go, how developers
start new Go projects, recent experiences with toolchain error messages, and
understanding developer interest around ML/AI.

Thank you to everyone who participated in this survey! This report shares what
we learned from your feedback.

## tl;dr

1. Go developers said they are **more interested in AI/ML tooling that
   improves the quality, reliability, and performance of code they write**,
   rather than writing code for them. An always-awake, never-busy expert
   "reviewer" might be one of the more helpful forms of AI developer
   assistance.
1. The top requests for improving toolchain warnings and errors were to **make
   the messages more comprehensible and actionable**; this sentiment was
   shared by developers of all experience levels, but was particularly strong
   among newer Go developers.
1. Our experiment with project templates (`gonew`) appears to solve critical
   problems for Go developers (especially developers new to Go) and does so in
   a way that matches their existing workflows for starting a new project.
   Based on these findings, we believe **`gonew` can substantially reduce
   onboarding barriers for new Go developers and ease adoption of Go in
   organizations**.
1. Three out of every four respondents work on Go software that also uses
   cloud services; this is evidence that **developers see Go as a language for
   modern, cloud-based development**.
1. **Developer sentiment towards Go remains extremely positive**, with 90% of
   survey respondents saying they felt satisfied while working with Go during
   the prior year.

## Contents

- <a href="#sentiment">Developer sentiment</a>
- <a href="#devenv">Developer environments</a>
- <a href="#stacks">Tech stacks</a>
- <a href="#gonew">How developers start new Go projects</a>
- <a href="#err_handling">Developer goals for error handling</a>
- <a href="#mlai">Understanding ML/AI use cases</a>
- <a href="#err_msgs">Toolchain error messages</a>
- <a href="#microservices">Microservices</a>
- <a href="#modules">Module authorship and maintenance</a>
- <a href="#demographics">Demographics</a>
- <a href="#firmographics">Firmographics</a>
- <a href="#methodology">Methodology</a>
- <a href="#closing">Closing</a>

## Developer sentiment {#sentiment}

Go developers continue to report high levels of satisfaction with the Go
ecosystem. A large majority of respondents said they felt satisfied while
working with Go over the past year (90% satisfied, 6% dissatisfied), and a
majority (52%) went further and said they were "very satisfied", the highest
rating. Longtime readers have likely noticed that this number doesn't change
much from year to year. This is expected for a large, stable project like Go;
we view this metric as a [lagging
indicator](https://en.wikipedia.org/wiki/Economic_indicator#Lagging_indicators)
that can help confirm widespread issues in the Go ecosystem, but isn't where
we expect to first learn about potential problems.

We typically find that the longer someone has worked with Go, the more likely
they are to report being satisfied with it. This trend continued in 2023;
among respondents with less than one year of Go experience, 82% reported
satisfaction with the Go development experience, compared to the 94% of Go
developers with five or more years of experience. There are likely a mix of
factors contributing to this, such as some respondents developing an
appreciation for Go's design choices over time, or deciding Go isn't a good
fit for their work and so not returning to this survey in following years
(i.e., [survivorship bias](https://en.wikipedia.org/wiki/Survivorship_bias)).
Still, this data helps us quantify the current getting started experience for
Go developers, and it seems clear we could do more to help emerging Gophers
find their footing and enjoy early successes developing with Go.

The key takeaway is that a large majority of people who chose to work with Go
during the past year were happy with their experience. Further, the number of
people working with Go continues to increase; we see evidence of this from
external research like [Stack Overflow's Developer
Survey](https://survey.stackoverflow.co/2023/#most-popular-technologies-language-prof)
(which found 14% of professional developers worked with Go during the past
year, a roughly 15% year-over-year increase), as well as analytics for
[go.dev](/) (which show an 8% rise in visitors year-over-year).
Combining this growth with a high satisfaction score is evidence that Go
continues to appeal to developers, and suggests that many developers who
choose to learn the language feel good about their decision long afterwards.
In their own words:

> "After 30+ years of development in C, C++, Java, and now seven years of
> programming in Go, it is still the most productive language by far. It's not
> perfect (no language is), but it has the best balance of productivity,
> complexity, and performance." <span class="quote_source">--- Professional Go
> developer w/ 5 -- 9 years of experience</span>

> "This is currently the best language I know, and I've tried many. The
> tooling is awesome, compile times are great, and I can be really productive.
> I'm glad I have Go as a tool, and I don't need to use TypeScript
> server-side. Thanks." <span class="quote_source">--- Open source Go
> developer w/ 3 -- 4 years of experience</span>

<img src="survey2023h2/csat.svg" alt="Chart of developer satisfaction with Go"
class="chart" />

## Developer environments {#devenv}

As in prior years, the majority of survey respondents told us they work with
Go on Linux (63%) and macOS (58%) systems. Small variations in these numbers
from year to year are most likely dependent upon who finds and responds to
this survey (particularly on the Go blog), as we don't see consistent
year-over-year trends in the random sample coming from VS Code.

We do continue to see that newer members of the Go community are more likely
to be working with Windows than more experienced Go developers. We interpret
this as a signal that Windows-based development is important for onboarding
new developers to the Go ecosystem, and is a topic our team hopes to focus on
more in 2024.

<img src="survey2023h2/os_dev.svg" alt="Chart of operating systems respondents
use when developing Go software" class="chart" /> <img
src="survey2023h2/os_dev_exp.svg" alt="Chart of operating systems respondents
use when developing Go software, split by duration of experience"
class="chart" />

Respondents continue to be heavily focused on Linux deployments. Given the
prevalence of Go for cloud development and containerized workloads, this is
not surprising but is still an important confirmation. We found few meaningful
differences based on factors such as organization size or experience level;
indeed, while novice Go developers appear more likely to *develop* on Windows,
92% still *deploy* to Linux systems. Perhaps the most interesting finding from
this breakdown is that more experienced Go developers said they deploy to a
wider variety of systems (most notably WebAssembly and IoT), though it's
unclear if this is because such deployments are challenging for newer Go
developers or the result of experienced Go developers using Go in a broader
range of contexts. We also observed that both IoT and WebAssembly have
steadily increased in recent years, with each rising from 3% in 2021 to 6% and
5% in 2023, respectively.

<img src="survey2023h2/os_deploy.svg" alt="Chart of platforms respondents
deploy Go software to" class="chart" />

The computing architecture landscape has changed over the past few years, and
we see that reflected in the current architectures Go developers say they work
with. While x86-compatible systems still account for the majority of
development (89%), ARM64 is also now used by a majority of respondents (56%).
This adoption appears to be partly driven by Apple Silicon; macOS developers
are now more likely to say they develop for ARM64 than for x86-based
architectures (76% vs. 71%). However, Apple hardware isn't the only factor
driving ARM64 adoption: among respondents who don't develop on macOS at all,
29% still say they develop for ARM64.

<img src="survey2023h2/arch.svg" alt="Chart of architectures respondents use
with Go" class="chart" />

The most common code editors among Go Developer Survey respondents continue to
be [VS Code](https://code.visualstudio.com/) (44%) and
[GoLand](https://www.jetbrains.com/go/) (31%). Both of these proportions
ticked down slightly from 2023 H1 (46% and 33%, respectively), but remain
within this survey's margin of error. Among the "Other" category,
[Helix](https://helix-editor.com/) accounted for the majority of responses.
Similar to the results for operating systems above, we don't believe this
represents a meaningful shift in code editor usage, but rather shows some of
the variability we expect to see in a community survey such as this. In
particular, we exclude the randomly sampled respondents from VS Code for this
question, as we know that group is heavily biased towards VS Code. However,
that has the side effect of making these results more susceptible to variation
each year.

We also looked at respondents' level of satisfaction with Go based on the
editor they prefer using. After controlling for length of experience, we found
no differences: we don't believe people enjoy working with Go more or less
based on which code editor they use. That doesn't necessarily mean all Go
editors are equal, but may reflect that people find the editor that is best
for their own needs. This would suggest the Go ecosystem has a healthy
diversity of different editors geared towards different use cases and
developer preferences.

<img src="survey2023h2/editor_self_select.svg" alt="Chart of code editors
respondents prefer to use with Go" class="chart" />

## Tech stacks {#stacks}

To better understand the web of software and services that Go developers
interact with, we asked several questions about tech stacks. We're sharing
these results with the community to show which tools and platforms are in
common use today, but we believe everyone should consider their own needs and
use cases when selecting a tech stack. More plainly: we neither intend for
readers to use this data to select components of their tech stack because they
are popular, nor to avoid components because they are not commonly used.

First, we can say with confidence that Go is a language for modern cloud-based
development. Indeed, 75% of respondents work on Go software that integrates
with cloud services. For nearly half of respondents, this involved AWS (48%),
and almost one-third used GCP (29%) for their Go development and deployments.
For both AWS and GCP, usage is equally balanced among large enterprises and
smaller organizations. Microsoft Azure is the only cloud provider that is
significantly more likely to be used in large organizations (companies with >
1,000 employees) than smaller shops; other providers show no meaningful
differences in usage based on the size of the organization.

<img src="survey2023h2/cloud.svg" alt="Chart of cloud platforms respondents
use with Go" class="chart" />

Databases are extremely common components of software systems, and we found
that 91% of respondents said the Go services they work on use at least one.
Most frequently this was PostgreSQL (59%), but with double digits of
respondents reporting use of six additional databases, it's safe to say there
are not just a couple of standard DBs for Go developers to consider. We again
see differences based on organization size, with respondents from smaller
organizations more likely to report using PostgreSQL and Redis, while
developers from large organizations are somewhat more likely to use a database
specific to their cloud provider.

<img src="survey2023h2/db.svg" alt="Chart of databases respondents use with
Go" class="chart" />

Another common component respondents reported using were caches or key-value
stores; 68% of respondents said they work on Go software incorporating at
least one of these. Redis was clearly the most common (57%), followed at a
distance by etcd (10%) and memcached (7%).

<img src="survey2023h2/cache.svg" alt="Chart of caches respondents use with
Go" class="chart" />

Similar to databases, survey respondents told us they use a range of different
observability systems. Prometheus and Grafana were the most commonly cited
(both at 43%), but Open Telemetry, Datadog, and Sentry were all in double
digits.

<img src="survey2023h2/metrics.svg" alt="Chart of metric systems respondents
use with Go" class="chart" />

Lest anyone wonder "Have we JSON'd all the things?"... yes, yes we have.
Nearly every respondent (96%!) said their Go software uses the JSON data
format; that's about as close to universal as you'll see with self-reported
data. YAML, CSV, and protocol buffers are also all used by roughly half of
respondents, and double-digit proportions work with TOML and XML as well.

<img src="survey2023h2/data.svg" alt="Chart of data formats respondents use
with Go" class="chart" />

For authentication and authorization services, we found most respondents are
building upon the foundations provided by standards such as
[JWT](https://jwt.io/introduction) and [OAuth2](https://oauth.net/2/). This
also appears to be an area where an organization's cloud provider's solution
is about as likely to be used as most turn-key alternatives.

<img src="survey2023h2/auth.svg" alt="Chart of authentication systems
respondents use with Go" class="chart" />

Finally, we have a bit of a grab bag of other services that don't neatly fit
into the above categories. We found that nearly half of respondents work with
gRPC in their Go software (47%). For infrastructure-as-code needs, Terraform
was the tool of choice for about ¼ of respondents. Other fairly common
technologies used alongside Go included Apache Kafka, ElasticSearch, GraphQL,
and RabbitMQ.

<img src="survey2023h2/other_tech.svg" alt="Chart of authentication systems
respondents use with Go" class="chart" />

We also looked at which technologies tended to be used together. While nothing
clearly analogous to the classic [LAMP
stack](https://en.wikipedia.org/wiki/LAMP_(software_bundle)) emerged from this
analysis, we did identify some interesting patterns:

- All or nothing: Every category (except data formats) showed a strong
  correlation where if a respondent answered “None” to one category, they
  likely answered “None” for all of the others. We interpret this as evidence
  that a minority of use cases require none of these tech stack components,
  but once the use case requires any one of them, it likely requires (or is at
  least simplified by) more than just one.
- A bias towards cross-platform technologies: Provider-specific solutions
  (i.e., services that are unique to a single cloud platform) were not
  commonly adopted. However, if respondents used one provider-specific
  solution (e.g., for metrics), they were substantially more likely to also
  say they used cloud-specific solutions in order areas (e.g., databases,
  authentication, caching, etc.).
- Multicloud: The three biggest cloud platforms were most likely to be
  involved in multicloud setups. For example, if an organization is using any
  non-AWS cloud provider, they’re probably also using AWS. This pattern was
  clearest for Amazon Web Services, but was also apparent (to a lesser extent)
  for Google Cloud Platform and Microsoft Azure.

## How developers start new Go projects {#gonew}

As part of our [experimentation with project
templates](/blog/gonew), we wanted to understand how Go
developers get started with new projects today. Respondents told us their
biggest challenges were choosing an appropriate way to structure their project
(54%) and learning how to write idiomatic Go (47%). As two respondents phrased
it:

> "Finding an appropriate structure and the right abstraction levels for a new
> project can be quite tedious; looking at high-profile community and
> enterprise projects for inspiration can be quite confusing as everyone
> structures their project differently" <span class="quote_source">---
> Professional Go developer w/ 5 -- 9 years of Go experience</span>

> "It would be great if [Go had a] toolchain to create [a project's] basic
> structure for web or CLI like \`go init \<project name\>\`" <span
> class="quote_source">--- Professional Go developer w/ 3 -- 4 years of
> experience</span>

Newer Go developers were even more likely to encounter these challenges: the
proportions increased to 59% and 53% for respondents with less than two years
of experience with Go, respectively. These are both areas we hope to improve
via our `gonew` prototype: templates can provide new Go developers with
well-tested project structures and design patterns, with initial
implementations written in idiomatic Go. These survey results have helped our
team to keep the purpose of `gonew` focused on tasks the Go community most
struggle with.

<img src="survey2023h2/new_challenge.svg" alt="Chart of challenges respondents
faced when starting new Go projects" class="chart" />

A majority of respondents told us they either use templates or copy+paste code
from existing projects when starting a new Go project (58%). Among respondents
with less than five years of Go experience, this proportion increased to
nearly ⅔ (63%). This was an important confirmation that the template-based
approach in `gonew` seems to meet developers where they already are, aligning
a common, informal approach with `go` command-style tooling. This is further
supported by the common feature requests for project templates: a majority of
respondents requested 1) a pre-configured directory structure to organize
their project and 2) sample code for common tasks in the project domain. These
results are well-aligned with the challenges developers said they faced in the
previous section. The responses to this question also help tease apart the
difference between project structure and design patterns, with nearly twice as
many respondents saying they want Go project templates to provide the former
than the latter.

<img src="survey2023h2/new_approach.svg" alt="Chart of approaches respondents
used when starting new Go projects" class="chart" />

<img src="survey2023h2/templates.svg" alt="Chart of functionality respondents
requested when starting new Go projects" class="chart" />

A majority of respondents told us the ability to make changes to a template
*and* have those changes propagate to projects based on that template was of
at least moderate importance. Anecdotally, we haven't spoken with any
developers who *currently* have this functionality with home-grown template
approaches, but it suggests this is an interesting avenue for future
development.

<img src="survey2023h2/template_updates.svg" alt="Chart of respondent interest
in updatable templates" class="chart" />

## Developer goals for error handling {#err_handling}

A perennial topic of discussion among Go developers is potential improvements
to error handling. As one respondent summarized:

> "Error handling adds too much boilerplate (I know, you probably heard this
> before)" <span class="quote_source">--- Open source Go developer w/ 1 -- 2
> years of experience</span>

But, we also hear from numerous developers that they appreciate Go's approach
to error handling:

> "Go error handling is simple and effective. As I have backends in Java and
> C# and exploring Rust and Zig now, I am always pleased to go back to write
> Go code. And one of the reasons is, believe it or not, error handling. It is
> really simple, plain and effective. Please leave it that way." <span
> class="quote_source">--- Open source Go developer w/ 5 -- 9 years of
> experience</span>

Rather than ask about specific modifications to error handling in Go, we
wanted to better understand developers' higher-level goals and whether Go's
current approach has proven useful and usable. We found that a majority of
respondents appreciate Go's approach to error handling (55%) and say it helps
them know when to check for errors (50%). Both of these outcomes were stronger
for respondents with more Go experience, suggesting that either developers
grow to appreciate Go's approach to error handling over time, or that this is
one factor leading developers to eventually leave the Go ecosystem (or at
least stop responding to Go-related surveys). Many survey respondents also
felt that Go requires a lot of tedious, boilerplate code to check for errors
(43%); this remained true regardless of how much prior Go experience
respondents had. Interestingly, when respondents said they appreciate Go's
error handling, they were unlikely to say it also results in lots of
boilerplate code---our team had a hypothesis that Go developers can both
appreciate the language's approach to error handling and feel it's too
verbose, but only 14% of respondents agreed with both statements.

Specific issues that respondents cited include challenges knowing which error
types to check for (28%), wanting to easily show a stack trace along with the
error message (28%), and the ease with which errors can be entirely ignored
(19%). About ⅓ of respondents were also interested in adopting concepts from
other languages, such as Rust's `?` operator (31%).

The Go team has no plans to add exceptions to the language, but since this is
anecdotally a common request, we included it as a response choice. Only 1 in
10 respondents said they wished they could use exceptions in Go, and this was
inversely related to experience---more veteran Go developers were less likely
to be interested in exceptions than respondents newer to the Go community.

<img src="survey2023h2/error_handling.svg" alt="Chart of respondents' thoughts
about Go's error handling approach" class="chart" />

## Understanding ML/AI use cases {#mlai}

The Go team is considering how the unfolding landscape of new ML/AI
technologies may impact software development in two distinct veins: 1) how
might ML/AI tooling help engineers write better software, and 2) how might Go
help engineers bring ML/AI support to their applications and services? Below,
we delve into each of these areas.

### Helping engineers write better software

There's little denying we're in [a hype cycle around the possibilities for
AI/ML](https://www.gartner.com/en/articles/what-s-new-in-artificial-intelligence-from-the-2023-gartner-hype-cycle).
We wanted to take a step back to focus on the broader challenges developers
face and where they think AI might prove useful in their regular work. The
answers were a bit surprising, especially given the industry's current focus
on coding assistants.

First, we see a few AI use cases that about half of respondents thought could
be helpful: generating tests (49%), suggesting best practices in-situ (47%),
and catching likely mistakes early in the development process (46%). A
unifying theme of these top use cases is that each could help improve the
quality and reliability of code an engineer is writing. A fourth use case
(help writing documentation) garnered interest from about ⅓ of respondents.
The remaining cases comprise a long tail of potentially fruitful ideas, but
these are of significantly less general interest than the top four.

When we look at developers' duration of experience with Go, we find that
novice respondents are interested in help resolving compiler errors and
explaining what a piece of Go code does more than veteran Go developers. These
might be areas where AI could help improve the getting started experience for
new Gophers; for example, an AI assistant could help explain in natural
language what an undocumented block of code does, or suggest common solutions
to specific error messages. Conversely, we see no differences between
experience levels for topics like "catch common mistakes"---both novice and
veteran Go developers say they would appreciate tooling to help with this.

One can squint at this data and see three broad trends:

1. Respondents voiced interest in getting feedback from "expert reviewers" in
   real-time, not just during review time.
1. Generally, respondents appeared most interested in tooling that saves them
   from potentially less-enjoyable tasks (e.g., writing tests or documenting
   code).
1. Wholesale writing or translating of code was of fairly low interest,
   especially to developers with more than a year or two of experience.

Taken together, it appears that today, developers are less excited by the
prospect of machines doing the fun (e.g., creative, enjoyable, appropriately
challenging) parts of software development, but do see value in another set of
"eyes" reviewing their code and potentially handling dull or repetitive tasks
for them. As one respondent phrased it:

> "I'm specifically interested in using AI/ML to improve my productivity with
> Go. Having a system that is trained in Go best practices, can catch
> anti-patterns, bugs, generate tests, with a low rate of hallucination, would
> be killer." <span class="quote_source">--- Professional Go developer w/ 5 --
> 9 years of experience</span>

This survey, however, is just one data point in a quickly-evolving research
field, so it's best to keep these results in context.

<img src="survey2023h2/ml_use_cases.svg" alt="Chart of respondents' interest
in AI/ML support for development tasks" class="chart" />

### Bringing AI features to applications and services

In addition to looking at how Go developers might benefit from AI/ML-powered
tooling, we explored their plans for building AI-powered applications and
services (or supporting infrastructure) with Go. We found that we're still
early in [the adoption
curve](https://en.wikipedia.org/wiki/Technology_adoption_life_cycle): most
respondents have not yet tried to use Go in these areas, though every topic
saw some level of interest from roughly half of respondents. For example, a
majority of respondents reported interest in integrating the Go services they
work on with LLMs (49%), but only 13% have already done so or are currently
evaluating this use case. At the time of this survey, responses gently suggest
that developers may be most interested in using Go to call LLMs directly,
build the data pipelines needed to power ML/AI systems, and for creating API
endpoints other services can call to interact with ML/AI models. As one
example, this respondent described the benefits they hoped to gain by using Go
in their data pipelines:

> "I want to integrate the ETL [extract, transform, and load] part using Go,
> to keep a consistent, robust, reliable codebase." <span
> class="quote_source">--- Professional Go developer w/ 3 -- 4 years of
> experience</span>

<img src="survey2023h2/ml_adoption.svg" alt="Chart of respondents' current use
of (and interest in) Go for AI/ML systems" class="chart" />

## Toolchain error messages {#err_msgs}

Many developers can relate to the frustrating experience of seeing an error
message, thinking they know what it means and how to resolve it, but after
hours of fruitless debugging realize it meant something else entirely. One
respondent explained their frustration as follows:

> "So often the printed complaints wind up having nothing to do with the
> problem, but it can take an hour before I discover that that's the case. The
> error messages are unnervingly terse, and don't seem to go out of their way
> to guess as to what the user might be trying to do or [explain what they're]
> doing wrong." <span class="quote_source">--- Professional Go developer w/ 10+
> years of experience</span>

We believe the warnings and errors emitted by developer tooling should be
brief, understandable, and actionable: the human reading them should be able
to accurately understand what went wrong and what they can do to resolve the
issue. This is an admittedly high bar to strive for, and with this survey we
took some measurements to understand how developers perceive Go's current
warning and error messages.

When thinking about the most recent Go error message they worked through,
respondents told us there was much room for improvement. Only a small majority
understood what the problem was from the error message alone (54%), and even
fewer knew what to do next to resolve the issue (41%). It appears a relatively
small amount of additional information could meaningfully increase these
proportions, as ¼ of respondents said they mostly knew how to fix the problem,
but needed to see an example first. Further, with 11% of respondents saying
they couldn't make sense of the error message, we now have a baseline for
current understandability of the Go toolchain's error messages.

Improvements to Go's toolchain error messages would especially benefit
less-experienced Gophers. Respondents with up to two years of experience were
less likely than veteran Gophers to say they understood the problem (47% vs.
61%) or knew how to fix it (29% vs. 52%), and were twice as likely to need to
search online to fix the issue (21% vs. 9%) or even make sense of what the
error meant (15% vs. 7%).

We hope to focus on improving toolchain error messages during 2024. These
survey results suggest this is an area of frustration for developers of all
experience levels, and will particularly help newer developers get started
with Go.

<img src="survey2023h2/err_exp.svg" alt="Chart of error handling experiences"
class="chart" />

<img src="survey2023h2/err_exp_exp.svg" alt="Chart of error handling
experiences, split by duration of Go experience" class="chart" />

To understand *how* these messages might be improved, we asked survey
respondents an open-ended question: "If you could make a wish and improve one
thing about error messages in the Go toolchain, what would you change?". The
responses largely align with our hypothesis that good error messages are both
understandable and actionable. The most common response was some form of "Help
me understand what led to this error" (36%), 21% of respondents explicitly
asked for guidance to fix the problem, and 14% of respondents called out
languages such as Rust or Elm as exemplars which strive to do both of these
things. In the words of one respondent:

> "For compilation errors, Elm or Rust-style output pinpointing exact issue in
> the source code. Errors should include suggestions to fix them where
> possible... I think a general policy of 'optimize error output to be read by
> humans' with 'provide suggestions where possible' would be very welcome
> here." <span class="quote_source">--- Professional Go developer w/ 5 -- 9
> years of experience</span>

Understandably, there is a fuzzy conceptual boundary between toolchain error
messages and runtime error messages. For example, one of the top requests
involved improved stack traces or other approaches to assist debugging runtime
crashes (22%). Similarly, a surprising theme in 4% of the feedback was about
challenges with getting help from the `go` command itself. These are great
examples of the Go community helping us identify related pain points that
weren't otherwise on our radar. We started this investigation focused on
improving compile-time errors, but one of the core areas Go developers would
like to see improved actually relates to run-time errors, while another was
about the `go` command's help system.

> "When an error is thrown, the call stack can be huge and includes a bunch of
> files I don't care about. I just want to know where the problem is in MY
> code, not the library I'm using, or how the panic was handled." <span
> class="quote_source">--- Professional Go developer w/ 1 -- 2 years of
> experience</span>

> "Getting help via \`go help run\` dumps a wall of text, with links to
> further readings to find the available command-line flags. Or the fact that
> it understands \`go run --help\` but instead of showing the help, it says
> 'please run go help run instead'. Just show me list of flags in \`go run
> --help\`." <span class="quote_source">--- Professional Go developer w/ 3 --
> 4 years of experience</span>

<img src="survey2023h2/text_err_wish.svg" alt="Chart of potential improvements
for Go's error messages" class="chart" />

## Microservices {#microservices}

We commonly hear that developers find Go to be a great fit for microservices,
but we have never tried to quantify how many Go developers have adopted this
type of service architecture, understand how those services communicate with
one another, or the challenges developers encounter when working on them. This
year we added a few questions to better understand this space.

A plurality of respondents said they work mostly on microservices (43%), with
another ¼ saying they work on a mix of both microservices and monoliths. Only
about ⅕ of respondents work mostly on monolithic Go applications. This is one
of the few areas where we see differences based on the size of organization
respondents work at---large organizations seem more likely to have adopted a
microservice architecture than smaller companies. Respondents from large
organizations (>1,000 employees) were most likely to say they work on
microservices (55%), with only 11% of these respondents working primarily on
monoliths.

<img src="survey2023h2/service_arch.svg" alt="Chart of respondents' primary
service architecture" class="chart" />

We see some bifurcation in the number of microservices comprising Go
platforms. One group is composed of a handful (2 to 5) of services (40%),
while the other consists of larger collections, with a minimum of 10 component
services (37%). The number of microservices involved does not appear to be
correlated with organization size.

<img src="survey2023h2/service_num.svg" alt="Chart of the number of
microservices respondents' systems involve" class="chart" />

A large majority of respondents use some form of direct response request
(e.g., RPC, HTTP, etc.) for microservice communication (72%). A smaller
proportion use message queues (14%) or a pub/sub approach (9%); again, we see
no differences here based on organization size.

<img src="survey2023h2/service_comm.svg" alt="Chart of how microservices
communicate with one another" class="chart" />

A majority of respondents build microservices in a polyglot of languages, with
only about ¼ exclusively using Go. Python is the most common companion
language (33%), alongside Node.js (28%) and Java (26%). We again see
differences based on organization size, with larger organizations more likely
to be integrating Python (43%) and Java (36%) microservices, while smaller
organizations are a bit more likely to only use Go (30%). Other languages
appeared to be used equally based on organization size.

<img src="survey2023h2/service_lang.svg" alt="Chart of other languages that Go
microservices interact with" class="chart" />

Overall, respondents told us testing and debugging were their biggest
challenge when writing microservice-based applications, followed by
operational complexity. Many other challenges occupy the long tail on this
graph, though "portability" stands out as a non-issue for most respondents. We
interpret this to mean that such services aren't intended to be portable
(beyond basic containerization); for example, if an organization's
microservices are initially powered by PostgreSQL databases, developers aren't
concerned with potentially porting this to an Oracle database in the near
future.

<img src="survey2023h2/service_challenge.svg" alt="Chart of challenges
respondents face when writing microservice-based applications" class="chart" />

## Module authorship and maintenance {#modules}

Go has a vibrant ecosystem of community-driven modules, and we want to
understand the motivations and challenges faced by developers who maintain
these modules. We found that about ⅕ of respondents maintain (or used to
maintain) an open-source Go module. This was a surprisingly high proportion,
and may be biased due to how we share this survey: module maintainers may be
more likely to closely follow the Go blog (where this survey is announced)
than other Go developers.

<img src="survey2023h2/mod_maintainer.svg" alt="Chart of how many respondents
have served as a maintainer for a public Go module" class="chart" />

Module maintainers appear to be largely self-motivated---they report working
on modules that they need for personal (58%) or work (56%) projects, that they
do so because they enjoy working on these modules (63%) and being part of the
public Go community (44%), and that they learn useful skills from their module
maintainership (44%). More external motivations, such as receiving recognition
(15%), career advancement (36%), or cash money (20%) are towards the bottom of
the list.

<img src="survey2023h2/mod_motivation.svg" alt="Chart of the motivations of
public module maintainers" class="chart" />

Given the forms of [intrinsic
motivation](https://en.wikipedia.org/wiki/Motivation#Intrinsic_and_extrinsic) identified above, it
follows that a key challenge for module maintainers is finding time to devote
to their module (41%). While this might not seem like an actionable finding in
itself (we can't give Go developers an extra hour or two each day, right?),
it's a helpful lens through which to view module tooling and
development---these tasks are most likely occurring while the developer is
already pressed for time, and perhaps it's been weeks or months since they
last had an opportunity to work on it, so things aren't fresh in their memory.
Thus, aspects like understandable and actionable error messages can be
particularly helpful: rather than require someone to once again search for
specific `go` command syntax, perhaps the error output could provide the
solution they need right in their terminal.

<img src="survey2023h2/mod_challenge.svg" alt="Chart of challenges respondents
face when maintaining public Go modules" class="chart" />

## Demographics {#demographics}

Most survey respondents reported using Go for their primary job (78%), and a
majority (59%) said they use it for personal or open-source projects. In fact,
it's common for respondents to use Go for *both* work and personal/OSS
projects, with 43% of respondents saying they use Go in each of these
situations.

<img src="survey2023h2/where.svg" alt="Chart of situations in which
respondents recently used Go" class="chart" />

The majority of respondents have been working with Go for under five years
(68%). As we've seen in [prior
years](/blog/survey2023-q1-results#novice-respondents-are-more-likely-to-prefer-windows-than-more-experienced-respondents),
people who found this survey via VS Code tended to be less experienced than
people who found the survey via other channels.

When we break down where people use Go by their experience level, two findings
stand out. First, a majority of respondents from all experience levels said
they're using Go professionally; indeed, for people with over two years of
experience, the vast majority use Go at work (85% -- 91%). A similar trend
exists for open-source development. The second finding is that developers with
less Go experience are more likely to be using Go to expand their skill set
(38%) or to evaluate it for use at work (13%) than more experienced Go
developers. We interpret this to mean that many Gophers initially view Go as
part of "upskilling" or expanding their understanding of software development,
but that within a year or two, they look to Go as more of a tool for doing
than learning.

<img src="survey2023h2/go_exp.svg" alt="Chart of how long respondents have
been working with Go" class="chart" />

<img src="survey2023h2/where_exp.svg" alt="Chart of situations in which
respondents recently used Go, split by their level of Go experience"
class="chart" />

The most common use cases for Go continue to be API/RPC services (74%) and
command line tools (62%). People tell us Go is a great choice for these types
of software for several reasons, including its built-in HTTP server and
concurrency primitives, ease of cross-compilation, and single-binary
deployments.

The intended audience for much of this tooling is in business settings (62%),
with 17% of respondents reporting that they develop primarily for more
consumer-oriented applications. This isn't surprising given the low use of Go
for consumer-focused applications such as desktop, mobile, or gaming, vs. its
very high use for backend services, CLI tooling, and cloud development, but it
is a useful confirmation of how heavily Go is used in B2B settings.

We also looked for differences based on respondents' level of experience with
Go and organization size. More experienced Go developers reported building a
wider variety of different things in Go; this trend was consistent across
every category of app or service. We did not find any notable differences in
what respondents are building based on their organization size.

<img src="survey2023h2/what.svg" alt="Chart of the types of things respondents
are building with Go" class="chart" />

<img src="survey2023h2/enduser.svg" alt="Chart of the audience using the
software respondents build" class="chart" />

Respondents were about equally likely to say this was the first time they've
responded to the Go Developer Survey vs. saying they had taken this survey
before. There is a meaningful difference between people who learned about this
survey via the Go blog, where 61% reported taking this survey previously, vs.
people who learned about this survey via a notification in VS Code, where only
31% said they've previously taken this survey. We don't expect people to
perfectly recall every survey they've responded to on the internet, but this
gives us some confidence that we're hearing from a balanced mix of new and
repeat respondents with each survey. Further, this tells us our combination of
social media posts and random in-editor sampling are both necessary for
hearing from a diverse set of Go developers.

<img src="survey2023h2/return_respondent.svg" alt="Chart of how many
respondents said they have taken this survey before" class="chart" />

## Firmographics {#firmographics}

Respondents to this survey reported working at a mix of different
organizations, from thousand-person-plus enterprises (27%), to midsize
businesses (25%) and smaller organizations with < 100 employees (44%). About
half of respondents work in the technology industry (50%), a large increase
over the next most-common industry---financial services---at 13%.

This is statistically unchanged from the past few Go Developer Surveys---we
continue to hear from people in different countries and in organizations of
different sizes and industries at consistent rates year after year.

<img src="survey2023h2/org_size.svg" alt="Chart of the different organization
sizes where respondents use Go" class="chart" />

<img src="survey2023h2/industry.svg" alt="Chart of the different industries
where respondents use Go" class="chart" />

<img src="survey2023h2/location.svg" alt="Chart of countries or regions where
respondents are located" class="chart" />

## Methodology {#methodology}

Most survey respondents "self-selected" to take this survey, meaning they
found it on the Go blog or other social Go channels. A potential problem with
this approach is that people who don't follow these channels are less likely
to learn about the survey, and might respond differently than people who do
closely follow them. About 40% of respondents were randomly sampled, meaning
they responded to the survey after seeing a prompt for it in VS Code (everyone
using the VS Code Go plugin between mid-July -- mid-August 2023 had a 10% of
receiving this random prompt). This randomly sampled group helps us generalize
these findings to the larger community of Go developers.

### How to read these results

Throughout this report we use charts of survey responses to provide supporting
evidence for our findings. All of these charts use a similar format. The title
is the exact question that survey respondents saw. Unless otherwise noted,
questions were multiple choice and participants could only select a single
response choice; each chart's subtitle will tell the reader if the question
allowed multiple response choices or was an open-ended text box instead of a
multiple choice question. For charts of open-ended text responses, a Go team
member read and manually categorized the responses. Many open-ended questions
elicited a wide variety of responses; to keep the chart sizes reasonable, we
condensed them to a maximum of the top 10 themes, with additional themes all
grouped under "Other". The percentage labels shown in charts are rounded to
the nearest integer (e.g., 1.4% and 0.8% will both be displayed as 1%), but
the length of each bar and row ordering are based on the unrounded values.

To help readers understand the weight of evidence underlying each finding, we
included error bars showing the 95% [confidence
interval](https://en.wikipedia.org/wiki/Confidence_interval) for responses;
narrower bars indicate increased confidence. Sometimes two or more responses
have overlapping error bars, which means the relative order of those responses
is not statistically meaningful (i.e., the responses are effectively tied).
The lower right of each chart shows the number of people whose responses are
included in the chart, in the form "n = [number of respondents]".

We include select quotes from respondents to help clarify many of our
findings. These quotes include the length of times the respondent has used Go.
If the respondent said they use Go at work, we refer to them as a
"professional Go developer"; if they don't use Go at work but do use Go for
open-source development, we refer to them as an "open-source Go developer".

## Closing {#closing}

The final question on our survey always asks respondents whether there's
anything else they'd like to share with us about Go. The most common piece of
feedback people provide is "thanks!", and this year was no different (33%). In
terms of requested language improvements, we see a three-way statistical tie
between improved expressivity (12%), improved error handling (12%), and
improved type safety or reliability (9%). Respondents had a variety of ideas
for improving expressivity, with the general trend of this feedback being
"Here's a specific thing I write frequently, and I wish it were easier to
express this in Go". The issues with error handling continue to be complaints
about the verbosity of this code today, while feedback about type safety most
commonly touched on [sum types](https://en.wikipedia.org/wiki/Tagged_union).
This type of high-level feedback is extremely useful when the Go team tries to
plan focus areas for the coming year, as it tells us general directions in
which the community is hoping to steer the ecosystem.

> "I know about Go's attitude towards simplicity and I appreciate it. I just
> wish there [were] slightly more features. For me it would be better error
> handling (not exceptions though), and maybe some common creature comforts
> like map/reduce/filter and ternary operators. Anything not too obscure
> that'll save me some 'if' statements." <span class="quote_source">---
> Professional Go developer w/ 1 -- 2 years of experience</span>

> "Please keep Go in line with the long term values Go established so long ago
> &mdash; language and library stability. [...] It is an environment I can
> rely on to not break my code after 2 or 3 years. For that, thank you very
> much." <span class="quote_source">--- Professional Go developer w/ 10+ years
> of experience</span>

<img src="survey2023h2/text_anything_else.svg" alt="Chart of other topics
respondents shared with us" class="chart" />

That's all for this bi-annual iteration of the Go Developer Survey. Thanks to
everyone who shared their feedback about Go---we have immense gratitude for
taking your time to help shape Go's future, and we hope you see some of your
own feedback reflected in this report. 🩵

--- Todd (on behalf of the Go team at Google)
