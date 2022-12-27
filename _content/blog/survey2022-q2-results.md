---
title: Go Developer Survey 2022 Q2 Results
date: 2022-09-08
by:
- Todd Kulesza
tags:
- survey
- community
summary: An analysis of the results from the 2022 Q2 Go Developer Survey.
---

<style type="text/css" scoped>
  .chart {
    margin-left: 1.5rem;
    margin-right: 1.5rem;
    width: 800px;
  }
  .quote {
    margin-left: 2rem;
    margin-right: 2rem;
    color: #999;
    font-style: italic;
    font-size: 120%;
  }
  @media (prefers-color-scheme: dark) {
    .chart {
      border-radius: 8px;
    }
  }
</style>

## Overview

This article shares the results of the 2022 June edition of the Go Developer
Survey. On behalf of the Go team, thank you to the 5,752 people who told us
about their experience working with new features introduced in Go 1.18,
including generics, security tooling, and workspaces. You've helped us better
understand how developers are discovering and using this functionality, and as
this article will discuss, provided useful insights for additional
improvements. Thank you! 💙

### Key findings

- __Generics has seen quick adoption__. A large majority of respondents were
  aware generics had been included in the Go 1.18 release, and about 1 in 4
  respondents said they've already started using generics in their Go code.
  The most common single piece of generics-related feedback was "thank you!",
  but it is clear developers are already running into some limitations of the
  initial generics implementation.
- __Fuzzing is new to most Go developers__. Awareness of Go's built-in fuzz
  testing was much lower than generics, and respondents had much more
  uncertainty around why or when they might consider using fuzz testing.
- __Third-party dependencies are a top security concern__. Avoiding
  dependencies with known vulnerabilities was the top security-related
  challenge for respondents. More broadly, security work can often be
  unplanned and unrewarded, implying that tooling needs to be respectful of
  developers' time and attention.
- __We can do better when announcing new functionality__. Randomly sampled
  participants were less likely to know about recent Go tooling releases than
  people who found the survey via the Go blog. This suggests we should either
  look beyond blog articles to communicate changes in the Go ecosystem, or
  expand efforts to share these articles more widely.
- __Error handling remains a challenge__. Following the release of generics,
  respondents' top challenge when working with Go shifted to error handling.
  Overall, however, satisfaction with Go remains very high, and we found no
  significant changes in how respondents said they were using Go.


### How to read these results

Throughout this post, we use charts of survey responses to provide supporting
evidence for our findings. All of these charts use a similar format. The title
is the exact question that survey respondents saw. Unless otherwise noted,
questions were multiple choice and participants could only select a single
response choice; each chart's subtitle will tell you if the question allowed
multiple response choices or was an open-ended text box instead of a multiple
choice question. For charts of open-ended text responses, a Go team member
read and manually categorized all of the responses. Many open-ended questions
elicited a wide variety of responses; to keep the chart sizes reasonable, we
condensed them to a maximum of the top 10 themes, with additional themes all
grouped under "Other".

To help readers understand the weight of evidence underlying each finding, we
include error bars showing the 95% confidence interval for responses; narrower
bars indicate increased confidence. Sometimes two or more responses have
overlapping error bars, which means the relative order of those responses is
not statistically meaningful (i.e., the responses are effectively tied). The
lower right of each chart shows the number of people whose responses are
included in the chart, in the form "_n = [number of respondents]_".

### A note on methodology

Most survey respondents "self-selected" to take the survey, meaning they found
it on [the Go blog](https://go.dev/blog),   [@golang on
Twitter](https://twitter.com/golang), or other social Go channels. A potential
problem with this approach is that people who don't follow these channels are
less likely to learn about the survey, and might respond differently than
people who _do_ closely follow them. About one third of respondents were
randomly sampled, meaning they responded to the survey after seeing a prompt
for it in VS Code (everyone using the VS Code Go plugin between June 1 - June
21st 2022 had a 10% of receiving this random prompt). This randomly sampled
group helps us generalize these findings to the larger community of Go
developers. Most survey questions showed no meaningful difference between
these groups, but in the few cases with important differences, readers will
see charts that break down responses into "Random sample" and "Self-selected"
groups.

## Generics

<div class="quote">"[Generics] seemed like the only obvious missing feature from the first time I used the language. Has helped reduce code duplication a lot." &mdash; A survey respondent discussing generics</div>

After Go 1.18 was released with support for type parameters (more commonly
referred to as _generics_), we wanted to understand what the initial awareness
and adoption of generics looked like, as well as identify common challenges or
blockers for using generics.

The vast majority of survey respondents (86%) were already aware generics
shipped as part of the Go 1.18 release. We had hoped to see a simple majority
here, so this was much more awareness than we'd been expecting. We also found
that about a quarter of respondents had begun using generics in Go code (26%),
including 14% who said they are already using generics in production or
released code. A majority of respondents (54%) were not opposed to using
generics, but didn't have a need for them today. We also found that 8% of
respondents _wanted_ to use generics in Go, but were currently blocked by
something.

<img src="survey2022q2/generics_awareness.svg" alt="Chart showing most
respondents were aware Go 1.18 included generics" class="chart" /> <img
src="survey2022q2/generics_use.svg" alt="Chart showing 26% of respondents are
already using Go generics" class="chart" />

What was blocking some developers from using generics? A majority of
respondents fell into one of two categories. First, 30% of respondents said
they hit a limit of the current implementation of generics, such as wanting
parameterized methods, improved type inference, or switching on types.
Respondents said these issues limited the potential use cases for generics or
felt they made generic code unnecessarily verbose. The second category
involved depending on something that didn't (yet) support generics---linters
were the most common tool preventing adoption, but this list also included
things like organizations remaining on an earlier Go release or depending on a
Linux distribution that did not yet provide Go 1.18 packages (26%). A steep
learning curve or lack of helpful documentation was cited by 12% of
respondents. Beyond these top issues, respondents told us about a wide range
of less-common (though still meaningful) challenges, as shown in the chart
below. To avoid focusing on hypotheticals, this analysis only includes people
who said they were already using generics, or who tried to use generics but
were blocked by something.

<img src="survey2022q2/text_gen_challenge.svg" alt="Chart showing the top
generic challenges" class="chart" />

We also asked survey respondents who had tried using generics to share any
additional feedback. Encouragingly, one in ten respondents said generics had
already simplified their code, or resulted in less code duplication. The most
common response was some variation of "thank you!" or a general positive
sentiment (43%); for comparison, only 6% of respondents evinced a negative
reaction or sentiment. Mirroring the findings from the "biggest challenge"
question above, nearly one third of respondents discussed hitting a limitation
of Go's implementation of generics. The Go team is using this set of results
to help decide if or how some of these limitations could be relaxed.

<img src="survey2022q2/text_gen_feedback.svg" alt="Chart showing most generics
feedback was positive or referenced a limitation of the current
implementation" class="chart" />

## Security

<div class="quote">"[The biggest challenge is] finding time given competing priorities; business customers want their features over security." &mdash; A survey respondent discussing security challenges</div>

Following the [2020 SolarWinds
breach](https://en.wikipedia.org/wiki/2020_United_States_federal_government_data_breach#SolarWinds_exploit),
the practice of developing software securely has received renewed attention.
The Go team has prioritized work in this area, including tools for creating [a
software bill of materials (SBOM)](https://pkg.go.dev/debug/buildinfo), [fuzz
testing](https://go.dev/doc/fuzz/), and most recently, [vulnerability
scanning](https://go.dev/blog/vuln/). To support these efforts, this survey
asked several questions about software development security practices and
challenges. We specifically wanted to understand:

- What types of security tools are Go developers using today?
- How do Go developers find and resolve vulnerabilities?
- What are the biggest challenges to writing secure Go software?

Our results suggest that while static analysis tooling is in widespread use
(65% of respondents), a minority of respondents currently use it to find
vulnerabilities (35%) or otherwise improve code security (33%). Respondents
said that security tooling is most commonly run during CI/CD time (84%), with
a minority saying developers run these tools locally during development (22%).
This aligns with additional security research our team has conducted, which
found that security scanning at CI/CD time is a desired backstop, but
developers often considered this too late for a first notification: they would
prefer to know a dependency may be vulnerable _before_ building upon it, or to
verify that a version update resolved a vulnerability without waiting for CI
to run a full battery of additional tests against their PR.

<img src="survey2022q2/dev_techniques.svg" alt="Chart showing prevalence of 9
different development techniques" class="chart" /> <img
src="survey2022q2/security_sa_when.svg" alt="Chart showing most respondents
run security tools during CI" class="chart" />

We also asked respondents about their biggest challenges around developing
secure software. The most wide-spread difficulty was evaluating the security
of third-party libraries (57% of respondents), a topic vulnerability scanners
(such as [GitHub's dependabot](https://github.com/dependabot) or the Go team's
[govulncheck](https://pkg.go.dev/golang.org/x/vuln/cmd/govulncheck)) can help
address. The other top challenges suggest opportunities for additional
security tooling: respondents said it’s hard to consistently apply best
practices while writing code, and validating that the resulting code doesn’t
have vulnerabilities.

<img src="survey2022q2/security_challenges.svg" alt="Chart showing the most
common security challenge is evaluating the security of third-party libraries"
class="chart" />

Fuzz testing, another approach for increasing application security, was still
quite new to most respondents. Only 12% said they use it at work, and 5% said
they've adopted Go's built-in fuzzing tools. An open-ended follow-up question
asking what made fuzzing difficult to use found that the main reasons were not
technical problems: the top three responses discussed not understanding how to
use fuzz testing (23%), a lack of time to devote to fuzzing or security more
broadly (22%), and understanding why and when developers might want to use
fuzz testing (14%). These findings indicate that we still have work to do in
terms of communicating the value of fuzz testing, what should be fuzz tested,
and how to apply it to a variety of different code bases.

<img src="survey2022q2/fuzz_use.svg" alt="Chart showing most respondents have
not tried fuzz testing yet" class="chart" /> <img
src="survey2022q2/text_fuzz_challenge.svg" alt="Chart showing the biggest fuzz
testing challenges relate to understanding, rather than technical issues"
class="chart" />

To better understand common tasks around vulnerability detection and
resolution, we asked respondents whether they'd learned of any vulnerabilities
in their Go code or its dependencies during the past year. For those who did,
we followed up with questions asking how the most recent vulnerability was
discovered, how they investigated and/or resolved it, and what was most
challenging about the whole process.

First, we found evidence that vulnerability scanning is effective. One quarter
of respondents said they'd learned of a vulnerability in one of their
third-party dependencies. Recall, however, that only about ⅓ of respondents
were using vulnerability scanning at all---when we look at responses from
people who said they ran some sort of vulnerability scanner, this proportion
nearly doubles, from 25% → 46%. Besides vulnerabilities in dependencies or in
Go itself, 12% of respondents said they learned about vulnerabilities in their
own code.

A majority of respondents said they learned of vulnerabilities via security
scanners (65%). The single most common tool respondents cited was [GitHub's
dependabot](https://github.com/dependabot) (38%), making it more frequently
referenced than all other vulnerability scanners combined (27%). After
scanning tools, the most common method for learning about vulnerabilities were
public reports, such as release notes and CVEs (22%).

<img src="survey2022q2/security_found_vuln.svg" alt="Chart showing that most
respondents have not found security vulnerabilities during the past year"
class="chart" /> <img src="survey2022q2/text_vuln_find.svg" alt="Chart showing
that vulnerability scanners are the most common way to learn about security
vulnerabilities" class="chart" />

Once respondents learned about a vulnerability, the most common resolution was
to upgrade the vulnerable dependency (67%). Among respondents who also
discussed using a vulnerability scanner (a proxy for participants who were
discussing a vulnerability in a third-party dependency), this increased to
85%. Less than one third of respondents discussed reading the CVE or
vulnerability report (31%), and only 12% mentioned a deeper investigation to
understand whether (and how) their software was impacted by the vulnerability.

That only 12% of respondents said they performed an investigation into whether
a vulnerability was reachable in their code, or the potential impact it may
have had on their service, was surprising. To understand this better, we also
looked at what respondents said was most challenging about responding to
security vulnerabilities. They described several different topics in roughly
equal proportions, from ensuring that dependency updates didn't break
anything, to understanding how to update indirect dependencies via go.mod
files. Also in this list is the type of investigation needed to understand a
vulnerability's impact or root cause. When we focus on only the respondents
who said they performed these investigations, however, we see a clear
correlation: 70% of respondents who said they performed an investigation into
the vulnerability's potential impact cited it as the most challenging part of
this process. Reasons included not just the difficulty of the task, but the
fact that it was often both unplanned and unrewarded work.

The Go team believes these deeper investigations, which require an
understanding of _how_ an application uses a vulnerable dependency, are
crucial for understanding the risk the vulnerability may present to an
organization, as well as understanding whether a data breach or other security
compromise occurred. Thus, [we designed
`govulncheck`](https://go.dev/blog/vuln) to only alert developers when a
vulnerability is invoked, and point developers to the exact places in their
code using the vulnerable functions. Our hope is that this will make it easier
for developers to quickly investigate the vulnerabilities that truly matter to
their application, thus reducing the overall amount of unplanned work in this
space.

<img src="survey2022q2/text_vuln_resolve.svg" alt="Chart showing most
respondents resolved vulnerabilities by upgrading dependencies" class="chart" />
<img src="survey2022q2/text_vuln_challenge.svg" alt="Chart showing a 6-way
tie for tasks that were most challenging when investigating and resolving
security vulnerabilities" class="chart" />

## Tooling

Next, we investigated three questions focused on tooling:

- Has the editor landscape shifted since our last survey?
- Are developers using workspaces? If so, what challenges have they
  encountered while getting started?
- How do developers handle internal package documentation?

VS Code appears to be continuing to grow in popularity among survey
respondents, with the proportion of respondents saying it's their preferred
editor for Go code increasing from 42% → 45% since 2021. VS Code and GoLand,
the two most popular editors, showed no differences in popularity between
small and large organizations, though hobbyist developers were more likely to
prefer VS Code to GoLand. This analysis excludes the randomly sampled VS Code
respondents---we'd expect people we invited to the survey to show a preference
for the tool used to distribute the invitation, which is exactly what we saw
(91% of the randomly sampled respondents preferred VS Code).

Following the 2021 switch to [power VS Code's Go support via the gopls
language server](https://go.dev/blog/gopls-vscode-go), the Go team has been
interested in understanding developer pain points related to gopls. While we
receive a healthy amount of feedback from developers currently using gopls, we
wondered whether a large proportion of developers had disabled it shortly
after release, which could mean we weren't hearing feedback about particularly
problematic use cases. To answer this question, we asked respondents who said
they preferred an editor which supports gopls whether or not they _used_
gopls, finding that only 2% said they had disabled it; for VS Code
specifically, this dropped to 1%. This increases our confidence that we're
hearing feedback from a representative group of developers. For readers who
still have unresolved issues with gopls, please let us know by <a
href="https://github.com/golang/go/issues">filing an issue on GitHub</a>.

<img src="survey2022q2/editor_self_select.svg" alt="Chart showing the top
preferred editors for Go are VS Code, GoLand, and Vim / Neovim" class="chart" />
<img src="survey2022q2/use_gopls.svg" alt="Chart showing only 2% of
respondents disabled gopls" class="chart"/>

Regarding workspaces, it seems many people first learned about Go's support
for multi-module workspaces via this survey. Respondents who learned of the
survey through VS Code's randomized prompt were especially likely to say they
had not heard of workspaces before (53% of randomly sampled respondents vs.
33% of self-selecting respondents), a trend we also observed with awareness of
generics (though this was higher for both groups, with 93% of self-selecting
respondents aware that generics landed in Go 1.18 vs. 68% of randomly sampled
respondents). One interpretation is that there is a large audience of Go
developers we do not currently reach through the Go blog or existing social
media channels, which has traditionally been our primary mechanism for sharing
new functionality.

We found that 9% of respondents said they had tried workspaces, and an
additional 5% would like to but are blocked by something. Respondents
discussed a variety of challenges when trying to use Go workspaces. A lack of
documentation and helpful error message from the `go work` command top the
list (21%), followed by technical challenges such as refactoring existing
repositories (13%). Similar to challenges discussed in the security section,
we again see "lack of time / not a priority" in this list---we interpret this
to mean the bar to understand and setup workspaces is still a bit too high
compared to the benefits they provide, potentially because developers already
had workarounds in place.

<img src="survey2022q2/workspaces_use_s.svg" alt="Chart showing a majority of
randomly sampled respondents were not aware of workspaces prior to this
survey" class="chart" /> <img src="survey2022q2/text_workspace_challenge.svg"
alt="Chart showing that documentation and error messages were the top
challenge when trying to use Go workspaces" class="chart" />

Prior to the release of Go modules, organizations were able to run internal
documentation servers (such as [the one that powered
godoc.org](https://github.com/golang/gddo)) to provide employees with
documentation for private, internal Go packages. This remains true with
[pkg.go.dev](https://pkg.go.dev), but setting up such a server is more complex
than it used to be. To understand if we should invest in making this process
easier, we asked respondents how they view documentation for internal Go
modules today, and whether that's their preferred way of working.

The results show the most common way to view internal Go documentation today
is by reading the code (81%), and while about half of the respondents were
happy with this, a large proportion would prefer to have an internal
documentation server (39%). We also asked who might be most likely to
configure and maintain such a server: by a 2-to-1 margin, respondents thought
it would be a software engineer rather than someone from a dedicated IT
support or operations team. This strongly suggests that a documentation server
should be a turn-key solution, or at least easy for a single developer to get
running quickly (over, say, a lunch break), on the theory that this type of
work is yet one more responsibility on developers' already full plates.

<img src="survey2022q2/doc_viewing_today.svg" alt="Chart showing most
respondents use source code directly for internal package documentation"
class="chart" /> <img src="survey2022q2/doc_viewing_ideal.svg" alt="Chart
showing 39% of respondents would prefer to use a documentation server instead
of viewing source for docs" class="chart" /> <img
src="survey2022q2/doc_server_owner.svg" alt="Chart showing most respondents
expect a software engineer to be responsible for such a documentation server"
class="chart" />

## Who we heard from

Overall, the demographics and firmographics of respondents did not
meaningfully shift since [our 2021
survey](https://go.dev/blog/survey2021-results). A small majority of
respondents (53%) have at least two years of experience using Go, while the
rest are newer to the Go community. About ⅓ of respondents work at small
businesses (< 100 employees), ¼ work at medium-sized businesses (100 -- 1,000
employees), and ¼ work at enterprises (> 1,000 employees). Similar to last
year, we found that our VS Code prompt helped encourage survey participation
outside of North America and Europe.

<img src="survey2022q2/go_exp.svg" alt="Chart showing distribution of
respondents' Go experience" class="chart" /> <img src="survey2022q2/where.svg"
alt="Chart showing distribution of where respondents' use Go" class="chart" />
<img src="survey2022q2/org_size.svg" alt="Chart showing distribution of
organization sizes for survey respondents" class="chart" /> <img
src="survey2022q2/industry.svg" alt="Chart showing distribution of industry
classifications for survey respondents" class="chart" /> <img
src="survey2022q2/location_s.svg" alt="Chart showing where in the world survey
respondents live" class="chart" />

## How respondents use Go

Similar to the previous section, we did not find any statistically significant
year-over-year changes in how respondents are using Go. The two most common
use cases remain building API/RPC services (73%) and writing CLIs (60%). We
used linear models to investigate whether there was a relationship between how
long a respondent had been using Go and the types of things they were building
with it. We found that respondents with < 1 year of Go experience are more
likely to be building something in the bottom half of this chart (GUIs, IoT,
games, ML/AI, or mobile apps), suggesting that there is interest in using Go
in these domains, but the drop-off after one year of experience also implies
that developers hit significant barriers when working with Go in these areas.

A majority of respondents use either Linux (59%) or macOS (52%) when
developing with Go, and the vast majority deploy to Linux systems (93%). This
cycle we added a response choice for developing on Windows Subsystem for Linux
(WSL), finding that 13% of respondents use this when working with Go.

<img src="survey2022q2/go_app.svg" alt="Chart showing distribution of what
respondents build with Go" class="chart" /> <img src="survey2022q2/os_dev.svg"
alt="Chart showing Linux and macOS are the most common development systems"
class="chart" /> <img src="survey2022q2/os_deploy.svg" alt="Chart showing
Linux is the most common deployment platform" class="chart" />

## Sentiment and challenges

Finally, we asked respondents about their overall level of satisfaction or
dissatisfaction with Go during that past year, as well as the biggest
challenge they face when using Go. We found that 93% of respondents said they
were "somewhat" (30%) or "very" (63%) satisfied, which is not statistically
different from the 92% of respondents who said they were satisfied during the
2021 Go Developer Survey.

After years of generics consistently being the most commonly discussed
challenge when using Go, the support for type parameters in Go 1.18 finally
resulted in a new top challenge: our old friend, error handling. To be sure,
error handling is statistically tied with several other challenges, including
missing or immature libraries for certain domains, helping developers learn
and implement best practices, and other revisions to the type system, such as
support for enums or more functional programming syntax. Post-generics, there
appears to be a very long tail of challenges facing Go developers.

<img src="survey2022q2/csat.svg" alt="Chart showing 93% of survey respondents
are satisfied using Go, with 4% dissatisfied" class="chart" /> <img
src="survey2022q2/text_biggest_challenge.svg" alt="Chart showing a long tail
of challenges reported by survey respondents" class="chart" />

## Survey methodology

We publicly announced this survey on June 1st, 2022 via
[go.dev/blog](https://go.dev/blog) and [@golang](https://twitter.com/golang)
on Twitter. We also randomly prompted 10% of [VS
Code](https://code.visualstudio.com/) users via the Go plugin between June 1st
-- 21st. The survey closed on June 22nd, and partial responses (i.e., people
who started but did not finish the survey) were also recorded. We filtered out
data from respondents who completed the survey especially quickly (< 30
seconds) or tended to check all of the response choices for multi-select
questions. This left 5,752 responses.

About ⅓ of respondents came from the randomized VS Code prompt, and this group
tended to have less experience with Go than people who found the survey via
the Go blog or Go's social media channels. We used linear and logistic models
to investigate whether apparent differences between these groups were better
explained by this difference in experience, which was usually the case. The
exceptions are noted in the text.

This year we very much hoped to also share the raw dataset with the community,
similar to developer surveys from [Stack
Overflow](https://insights.stackoverflow.com/survey),
[JetBrains](https://www.jetbrains.com/lp/devecosystem-2021/), and others.
Recent legal guidance unfortunately prevents us from doing that right now, but
we're working on this and expect to be able to share the raw dataset for our
next Go Developer Survey.

## Conclusion

This iteration of the Go Developer Survey focused on new functionality from
the Go 1.18 release. We found that generics adoption is well under way, with
developers already hitting some limitations of the current implementation.
Fuzz testing and workspaces have seen slower adoption, though largely not for
technical reasons: the primary challenge with both was understanding when and
how to use them. A lack of developer time to focus on these topics was another
challenge, and this theme carried into security tooling as well. These
findings are helping the Go team prioritize our next efforts and will
influence how we approach the design of future tooling.

Thank you for joining us in the tour of Go developer research---we hope it's
been insightful and interesting. Most importantly, thank you to everyone who
has responded to our surveys over the years. Your feedback helps us understand
the constraints Go developers work under and identify challenges they face. By
sharing these experiences, you're helping to improve the Go ecosystem for
everyone. On behalf of Gophers everywhere, we appreciate you!
