---
title: Go 2018 Survey Results
date: 2019-03-28
by:
- Todd Kulesza, Steve Francia
tags:
- survey
- community
summary: What we learned from the December 2018 Go User Survey.
---

## Thank you

<style>
    p.note {
        font-size: 0.80em;
        font-family: "Helvetica Neue", Arial, sans-serif; /* Helvetica on Mac aka sans-serif has broken U+2007 */
    }
</style>

This post summarizes the results of our 2018 user survey and draws comparisons
between the results of our prior surveys from [2016](https://blog.golang.org/survey2016-results)
and [2017](https://blog.golang.org/survey2017-results).

This year we had 5,883 survey respondents from 103 different countries.
We are grateful to everyone who provided their feedback through this survey
to help shape the future of Go. Thank you!

## Summary of findings

  - For the first time, **half of survey respondents are now using Go as part of their daily routine**.
    This year also saw significant increases in the number of respondents who
    develop in Go as part of their jobs and use Go outside of work responsibilities.
  - The most **common uses for Go remain API/RPC services and CLI tools**.
    Automation tasks, while not as common as CLI tools and API services,
    are a fast-growing area for Go.
  - **Web development remains the most common domain** that survey respondents work in,
    but **DevOps showed the highest year-over-year growth** and is now the second most common domain.
  - A large majority of survey respondents said **Go is their most-preferred programming language**,
    despite generally feeling less proficient with it than at least one other language.
  - **VS Code and GoLand are surging in popularity** and are now the most popular code editors among survey respondents.
  - Highlighting the portable nature of Go,
    **many Go developers use more than one primary OS** for development.
    Linux and macOS are particularly popular,
    with a large majority of survey respondents using one or both of these operating
    systems to write Go code.
  - Survey respondents appear to be **shifting away from on-prem Go deployments**
    and moving towards containers and serverless cloud deployments.
  - The majority of respondents said they feel welcome in the Go community,
    and most ideas for improving the Go community specifically focus on **improving the experience of newcomers**.

Read on for all of the details.

## Programming background

This year's results show a significant increase in the number of survey
respondents who are paid to write Go as part of their jobs (68% → 72%),
continuing a year-over-year trend that has been growing since our first survey in 2016.
We also see an increase in the number of respondents who program in Go outside
of work (64% → 70%).
For the first time, the number of survey respondents who write in Go as
part of their daily routine reached 50% (up from 44% in 2016).
These findings suggest companies are continuing to embrace Go for professional
software development at a consistent pace,
and that Go's general popularity with developers remains strong.

{{image "survey2018/fig1.svg" 600}}
{{image "survey2018/fig2.svg" 600}}

To better understand where developers use Go,
we broke responses down into three groups:
1) people who are using Go both in and outside of work,
2) people who use Go professionally but not outside of work,
and 3) people who only write Go outside of their job responsibilities.
Nearly half (46%) of respondents write Go code both professionally and on
their own time (a 10-point increase since 2017),
while the remaining respondents are closely split between either only writing Go at work,
or only writing Go outside of work.
The large percentage of respondents who both use Go at work and choose to
use it outside of work suggests that the language appeals to developers
who do not view software engineering as a day job:
they also choose to hack on code outside of work responsibilities,
and (as evidenced by 85% of respondents saying they'd prefer Go for their next project,
see section _Attitudes towards Go_ below) Go is the top language they'd
prefer to use for these non-work-related projects.

{{image "survey2018/fig4.svg" 600}}

When asked how long they've been using Go,
participants' answers are strongly trending upward over time,
with a higher percentage of responses in the 2-4 and 4+ year buckets each year.
This is expected for a newer programming language,
and we're glad to see that the percentage of respondents who are new to
Go is dropping more slowly than the percentage of respondents who have been
using Go for 2+ years is increasing,
as this suggests that developers are not dropping out of the ecosystem after
initially learning the language.

{{image "survey2018/fig5.svg" 600}}

As in prior years, Go ranks at the top of respondents' preferred languages
and languages in which they have expertise.
A majority of respondents (69%) claimed expertise in 5 different languages,
highlighting that their attitudes towards Go are influenced by experiences
with other programming stacks.
The charts below are sorted by the number of respondents who ranked each
language as their most preferred/understood (the darkest blue bars),
which highlights three interesting bits:

  - While about ⅓ of respondents consider Go to be the language in which
    they have the most expertise,
    twice that many respondents consider it their most preferred programming language.
    So even though many respondents feel they haven't become as proficient with
    Go as with some other language,
    they still frequently prefer to develop with Go.
  - Few survey respondents rank Rust as a language in which they have expertise (6.8%),
    yet 19% rank it as a top preferred language,
    indicating a high level of interest in Rust among this audience.
  - Only three languages have more respondents who say they prefer the language
    than say they have expertise with it:
    Rust (2.41:1 ratio of preference:expertise),
    Kotlin (1.95:1), and Go (1.02:1).
    Higher preference than expertise implies interest—but little direct experience—in a language,
    while lower preference than expertise numbers suggests barriers to proficient use.
    Ratios near 1.0 suggest that most developers are able to work effectively
    _and_ enjoyably with a given language.
    This data is corroborated by [Stack Overflow's 2018 developer survey](https://insights.stackoverflow.com/survey/2018/#most-loved-dreaded-and-wanted),
    which also found Rust, Kotlin, and Go to be among the most-preferred programming languages.

{{image "survey2018/fig6.svg" 600}}
{{image "survey2018/fig7.svg" 600}}

<p class="note">
    <i>Reading the data</i>: Participants could rank their top 5 languages. The color coding starts with dark blue for the top rank and lightens for each successive rank. These charts are sorted by the percentage of participants who ranked each language as their top choice.
</p>

## Development domains

Survey respondents reported working on a median of three different domains,
with a large majority (72%) working in 2-5 different areas.
Web development is the most prevalent at 65%,
and it increased its dominance as the primary area survey respondents work
in (up from 61% last year):
web development has been the most common domain for Go development since 2016.
This year DevOps noticeably increased, from 36% to 41% of respondents,
taking over the number two spot from Systems Programming.
We did not find any domains with lower usage in 2018 than in 2017,
suggesting that respondents are adopting Go for a wider variety of projects,
rather than shifting usage from one domain to another.

{{image "survey2018/fig8.svg" 600}}

Since 2016, the top two uses of Go have been writing API/RPC services and
developing CLI applications.
While CLI usage has remained stable at 63% for three years,
API/RPC usage has increased from 60% in 2016 to 65% in 2017 to 73% today.
These domains play to core strengths of Go and are both central to cloud-native
software development,
so we expect them to remain two of the primary scenarios for Go developers into the future.
The percentage of respondents who write web services that directly return
HTML has steadily dropped while API/RPC usage has increased,
suggesting some migration to the API/RPC model for web services.
Another year-over-year trend suggests that automation is also a growing area for Go,
with 38% of respondents now using Go for scripts and automation tasks (up from 31% in 2016).

{{image "survey2018/fig9.svg" 600}}

To better understand the contexts in which developers are using Go,
we added a question about Go adoption across different industries.
Perhaps unsurprisingly for a relatively new language,
over half of survey respondents work in companies in the _Internet/web services_
and _Software_ categories (i.e., tech companies).
The only other industries with >3% responses were _Finance, banking, or insurance_
and _Media, advertising, publishing, or entertainment_.
(In the chart below, we've condensed all of the categories with response
rates below 3% into the "Other" category.) We'll continue tracking Go's
adoption across industries to better understand developer needs outside
of technology companies.

{{image "survey2018/fig10.svg" 600}}

## Attitudes towards Go

This year we added a question asking "How likely are you to recommend Go
to a friend or colleague?" to calculate our [Net Promoter Score](https://en.wikipedia.org/wiki/Net_Promoter).
This score attempts to measure how many more "promoters" a product has than
"detractors" and ranges from -100 to 100;
a positive value suggests most people are likely to recommend using a product,
while negative values suggest most people are likely to recommend against using it.
Our 2018 score is 61 (68% promoters - 7% detractors) and will serve as a
baseline to help us gauge community sentiment towards the Go ecosystem over time.

{{image "survey2018/fig11.svg" 600}}

In addition to NPS, we asked several questions about developer satisfaction with Go.
Overall, survey respondents indicated a high level of satisfaction,
consistent with prior years.
Large majorities say they are happy with Go (89%),
would prefer to use Go for their next project (85%),
and feel that it is working well for their team (66%),
while a plurality feel that Go is at least somewhat critical to their company's success (44%).
While all of these metrics showed an increase in 2017,
they remained mostly stable this year.
(The wording of the first question changed in 2018 from "_I would recommend using Go to others_"
to "_Overall, I'm happy with Go_",
so those results are not directly comparable.)

{{image "survey2018/fig12.svg" 600}}

Given the strong sentiment towards preferring Go for future development,
we want to understand what prevents developers from doing so.
These remained largely unchanged since last year:
about ½ of survey respondents work on existing projects written in other languages,
and ⅓ work on a team or project that prefers to use a different language.
Missing language features and libraries round out the most common reasons
respondents did not use Go more.
We also asked about the biggest challenges developers face while using Go;
unlike most of our survey questions, respondents could type in anything
they wished to answer this question.
We analyzed the results via machine learning to identify common themes and
counting the number of responses that supported each theme.
The top three major challenges we identified are:

  - Package management (e.g., "Keeping up with vendoring",
    "dependency / packet [sic] management / vendoring not unified")
  - Differences from more familiar programming languages (e.g.,
    "syntax close to C-languages with slightly different semantics makes me
    look up references somewhat more than I'd like",
    "coworkers who come from non-Go backgrounds trying to use Go as a version
    of their previous language but with channels and Goroutines")
  - Lack of generics (e.g., "Lack of generics makes it difficult to persuade
    people who have not tried Go that they would find it efficient.",
    "Hard to build richer abstractions (want generics)")

{{image "survey2018/fig13.svg" 600}}
{{image "survey2018/fig14.svg" 600}}

This year we added several questions about developer satisfaction with different aspects of Go.
Survey respondents were very satisfied with Go applications' CPU performance (46:1,
meaning 46 respondents said they were satisfied for every 1 respondent who
said they were not satisfied),
build speed (37:1), and application memory utilization (32:1).
Responses for application debuggability (3.2:1)  and binary size (6.4:1),
however, suggest room for improvement.

The dissatisfaction with binary size largely comes from developers building CLIs,
only 30% of whom are satisfied with the size of Go's generated binaries.
For all other types of applications, however,
developer satisfaction was > 50%, and binary size was consistently ranked
at the bottom of the list of important factors.

Debuggability, conversely, stands out when we look at how respondents ranked
the importance of each aspect;
44% of respondents ranked debuggability as their most or second-most important aspect,
but only 36% were satisfied with the current state of Go debugging.
Debuggability was consistently rated about as important as memory usage
and build speed but with significantly lower satisfaction levels,
and this pattern held true regardless of the type of software respondents were building.
The two most recent Go releases, Go 1.11 and 1.12,
both contained significant improvements to debuggability.
We plan to investigate how developers debug Go applications in more depth this year,
with a goal of improving the overall debugging experience for Go developers.

{{image "survey2018/fig15.svg" 600}}
{{image "survey2018/fig29.svg" 600}}

## Development environments

We asked respondents which operating systems they primarily use when writing Go code.
A majority (65%) of respondents said they use Linux,
50% use macOS, and 18% use Windows, consistent with last year.
This year we also looked at how many respondents develop on multiple OSes vs. a single OS.
Linux and macOS remain the clear leaders,
with 81% of respondents developing on some mix of these two systems.
Only 3% of respondents evenly split their time between all three OSes.
Overall, 41% of respondents use multiple operating systems for Go development,
highlighting the cross-platform nature of Go.

{{image "survey2018/fig16.svg" 600}}

Last year, VS Code edged out Vim as the most popular Go editor among survey respondents.
This year it significantly expanded its lead to become the preferred editor
for over ⅓ of our survey respondents (up from 27% last year).
GoLand also experienced strong growth and is now the second most-preferred editor at 22%,
swapping places with Vim (down to 17%).
The surging popularity of VS Code and GoLand appear to be coming at the
expense of Sublime Text and Atom.
Vim also saw the number of respondents ranking it their top choice drop,
but it remains the most popular second-choice editor at 14%.
Interestingly, we found no differences in the level of satisfaction respondents
reported for their editor(s) of choice.

We also asked respondents what would most improve Go support in their preferred editor.
Like the "biggest challenge" question above,
participants could write in their own response rather than select from a
multiple-choice list.
A thematic analysis on the responses revealed that _improved debugging support_ (e.g.,
"Live debugging", "Integrated debugging",
"Even better debugging") was the most-common request,
followed by _improved code completion_ (e.g.,
"autocomplete performance and quality", "smarter autocomplete").
Other requests include better integration with Go's CLI toolchain,
better support for modules/packages, and general performance improvements.

{{image "survey2018/fig17.svg" 600}}
{{image "survey2018/fig18.svg" 600}}

This year we also added a question asking which deployment architectures
are most important to Go developers.
Unsurprisingly, survey respondents overwhelmingly view x86/x86-64 as their
top deployment platform (76% of respondents listed it as their most important
deployment architecture,
and 84% had it in their top 3).
The ranking of the second- and third-choice architectures,
however, is informative:
there is significant interest in ARM64 (45%),
WebAssembly (30%), and ARM (22%), but very little interest in other platforms.

{{image "survey2018/fig19.svg" 600}}

## Deployments and services

For 2018 we see a continuation of the trend from on-prem to cloud hosting
for both Go and non-Go deployments.
The percentage of survey respondents who deploy Go applications to on-prem
servers dropped from 43% → 32%,
mirroring the 46% → 36% drop reported for non-Go deployments.
The cloud services which saw the highest year-over-year growth include AWS
Lambda (4% → 11% for Go,
10% → 15% non-Go) and Google Kubernetes Engine (8% → 12% for Go,
5% → 10% non-Go), suggesting that serverless and containers are becoming
increasingly popular deployment platforms.
This service growth appears to be driven by respondents who had already
adopted cloud services,
however, as we found no meaningful growth in the percentage of respondents
who deploy to at least one cloud service this year (55% → 56%).
We also see steady growth in Go deployments to GCP since 2016,
increasing from 12% → 19% of respondents.

{{image "survey2018/fig20.svg" 600}}

Perhaps correlated with the decrease in on-prem deployments,
this year we saw cloud storage become the second-most used service by survey respondents,
increasing from 32% → 44%.
Authentication & federation services also saw a significant increase (26% → 33%).
The primary service survey respondents access from Go remains open-source
relational databases,
which ticked up from 61% → 65% of respondents.
As the below chart shows, service usage increased across the board.

{{image "survey2018/fig21.svg" 600}}

## Go community

The top community sources for finding answers to Go questions continue to
be Stack Overflow (23% of respondents marked it as their top source),
Go web sites (18% for godoc.org, 14% for golang.org),
and reading source code (8% for source code generally,
4% for GitHub specifically).
The order remains largely consistent with prior years.
The primary sources for Go news remain the Go blog,
Reddit's r/golang, Twitter, and Hacker News.
These were also the primary distribution methods for this survey,
however, so there is likely some bias in this result.
In the two charts below, we've grouped sources used by less than < 5% of
respondents into the "Other" category.

{{image "survey2018/fig24.svg" 600}}
{{image "survey2018/fig25.svg" 600}}

This year, 55% of survey respondents said they have or are interested in
contributing to the Go community,
slightly down from 59% last year.
Because the two most common areas for contribution (the standard library
and official Go tools) require interacting with the core Go team,
we suspect this decrease may be related to a dip in the percentage of participants
who agreed with the statements "I feel comfortable approaching the Go project
leadership with questions and feedback" (30% → 25%) and "I am confident
in the leadership of Go (54% → 46%).

{{image "survey2018/fig26.svg" 600}}
{{image "survey2018/fig27.svg" 600}}

An important aspect of community is helping everyone feel welcome,
especially people from traditionally under-represented demographics.
To better understand this, we asked an optional question about identification
across several under-represented groups.
In 2017 we saw year-over-year increases across the board.
For 2018, we saw a similar percentage of respondents (12%) identify as part
of an under-represented group,
and this was paired with a significant decrease in the percentage of respondents
who do **not** identify as part of an under-represented group.
In 2017, for every person who identified as part of an under-represented group,
3.5 people identified as not part of an under-represented group (3.5:1 ratio).
In 2018 that ratio improved to 3.08:1. This suggests that the Go community
is at least retaining the same proportions of under-represented members,
and may even be increasing.

{{image "survey2018/fig28.svg" 600}}

Maintaining a healthy community is extremely important to the Go project,
so for the past three years we've been measuring the extent to which developers
feel welcome in the Go community.
This year we saw a drop in the percentage of survey respondents who agree
with the statement "I feel welcome in the Go community", from 66% → 59%.

To better understand this decrease, we looked more closely at who reported
feeling less welcome.
Among traditionally under-represented groups,
fewer people reported feeling unwelcome in 2018,
suggesting that outreach in that area has been helpful.
Instead, we found a linear relationship between the length of time someone
has used Go and how welcome they feel:
newer Go developers felt significantly less welcome (at 50%) than developers
with 1-2 years of experience (62%),
who in turn felt less welcome than developers with a few years of experience (73%).
This interpretation of the data is supported by responses to the question
"What changes would make the Go community more welcoming?".
Respondents' comments can be broadly grouped into four categories:

  - Reduce a perception of elitism, especially for newcomers to Go (e.g.,
    "less dismissiveness", "Less defensiveness and hubris")
  - Increase transparency at the leadership level (e.g.,
    "Future direction and planning discussions",
    "Less top down leadership", "More democratic")
  - Increase introductory resources (e.g., "A more clear introduction for contributors",
    "Fun challenges to learn best practices")
  - More events and meetups, with a focus on covering a larger geographic area (e.g.,
    "More meetups & social events", "Events in more cities")

This feedback is very helpful and gives us concrete areas we can focus on
to improve the experience of being a Go developer.
While it doesn't represent a large percentage of our user base,
we take this feedback very seriously and are working on improving each area.

{{image "survey2018/fig22.svg" 600}}
{{image "survey2018/fig23.svg" 600}}

## Conclusion

We hope you've enjoyed seeing the results of our 2018 developer survey.
These results are impacting our 2019 planning,
and in the coming months we'll share some ideas with you to address specific
issues and needs the community has highlighted for us.
Once again, thank you to everyone who contributed to this survey!
