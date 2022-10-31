---
title: Go 2017 Survey Results
date: 2018-02-26
by:
- Steve Francia
tags:
- survey
- community
summary: What we learned from the December 2017 Go User Survey.
---

## Thank you

This post summarizes the result of our 2017 user survey along with commentary
and insights. It also draws key comparisons between the results of the 2016 and
2017 survey.

This year we had 6,173 survey respondents, 70% more than the 3,595 we had in the
[Go 2016 User Survey](https://blog.golang.org/survey2016-results). In
addition, it also had a slightly higher completion rate (84% → 87%) and a
higher response rate to most of the questions. We believe that survey length is
the main cause of this improvement as the 2017 survey was shortened in response
to feedback that the 2016 survey was too long.

We are grateful to everyone who provided their feedback through the survey to
help shape the future of Go.

## Programming background

For the first time, more survey respondents say they are paid to write Go
than say they write it outside work. This indicates a significant shift in
Go's user base and in its acceptance by companies for professional software
development.

The areas people who responded to the survey work in is mostly consistent with
last year, however, mobile and desktop applications have fallen significantly.

Another important shift: the #1 use of Go is now writing API/RPC services (65%,
up 5% over 2016), taking over the top spot from writing CLI tools in Go (63%).
Both take full advantage of Go's distinguishing features and are key elements of
modern cloud computing. As more companies adopt Go, we expect these two uses
of Go to continue to thrive.

Most of the metrics reaffirm things we have learned in prior years. Go
programmers still overwhelmingly prefer Go. As more time passes Go users are
deepening their experience in Go. While Go has increased its lead among Go
developers, the order of language rankings remains quite consistent with last
year.

{{raw (file "survey2017/background.html")}}

## Go usage

In nearly every question around the usage and perception of Go, Go has
demonstrated improvement over our prior survey. Users are happier using Go, and
a greater percentage prefer using Go for their next project.

When asked about the biggest challenges to their own personal use of Go, users
clearly conveyed that lack of dependency management and lack of generics were
their two biggest issues, consistent with 2016. In 2017 we laid a foundation to
be able to address these issues. We improved our proposal and development
process with the addition of
[Experience Reports](/wiki/ExperienceReports) which is
enabling the project to gather and obtain feedback critical to making these
significant changes. We also made
[significant changes](/doc/go1.10#build) under the hood in
how Go obtains, and builds packages. This is foundational work essential to
addressing our dependency management needs.

These two issues will continue to be a major focus of the project through 2018.

In this section we asked two new questions. Both center around what
developers are doing with Go in a more granular way than we've previously asked.
We hope this data will provide insights for the Go project and ecosystem.

Since last year there has been an increase of the percentage of people who
identified "Go lacks critical features" as the reason they don't use Go more and
a decreased percentage who identified "Go not being an appropriate fit". Other
than these changes, the list remains consistent with last year.

{{raw (file "survey2017/usage.html")}}

## Development and deployment

We asked programmers which operating systems they develop Go on; the ratios of
their responses remain consistent with last year. 64% of respondents say
they use Linux, 49% use MacOS, and 18% use Windows, with multiple choices
allowed.

Continuing its explosive growth, VSCode is now the most popular editor among
Gophers. IntelliJ/GoLand also saw significant increase in usage. These largely
came at the expense of Atom and Sublime Text which saw relative usage drops.
This question had a 6% higher response rate from last year.

Survey respondents demonstrated significantly higher satisfaction with Go
support in their editors over 2016 with the ratio of satisfied to dissatisfied
doubling (9:1 → 18:1). Thank you to everyone who worked on Go editor support
for all your hard work.

Go deployment is roughly evenly split between privately managed servers and
hosted cloud servers. For Go applications, Google Cloud services saw significant
increase over 2016. For Non-Go applications, AWS Lambda saw the largest increase in use.

{{raw (file "survey2017/dev.html")}}

## Working Effectively

We asked how strongly people agreed or disagreed with various statements about
Go. All questions are repeated from last year with the addition of one new
question which we introduced to add further clarification around how users are
able to both find and **use** Go libraries.

All responses either indicated a small improvement or are comparable to 2016.

As in 2016, the most commonly requested missing library for Go is one for
writing GUIs though the demand is not as pronounced as last year. No other
missing library registered a significant number of responses.

The primary sources for finding answers to Go questions are the Go web site,
Stack Overflow, and reading source code directly. Stack Overflow showed a small
increase from usage over last year.

The primary sources for Go news are still the Go blog, Reddit’s /r/golang and
Twitter; like last year, there may be some bias here since these are also how
the survey was announced.

{{raw (file "survey2017/effective.html")}}

## The Go Project

59% of respondents expressed interest in contributing in some way to the Go
community and projects, up from 55% last year. Respondents also indicated that
they felt much more welcome to contribute than in 2016. Unfortunately,
respondents indicated only a very tiny improvement in understanding how to
contribute. We will be actively working with the community and its leaders
to make this a more accessible process.

Respondents showed an increase in agreement that they are confident in the
leadership of the Go project (9:1 → 11:1). They also showed a small increase in
agreement that the project leadership understands their needs (2.6:1 → 2.8:1)
and in agreement that they feel comfortable approaching project leadership with
questions and feedback (2.2:1 → 2.4:1). While improvements were made, this
continues to be an area of focus for the project and its leadership going
forward. We will continue to work to improve our understanding of user needs and
approachability.

We tried some [new ways](https://blog.golang.org/8years#TOC_1.3.) to engage
with users in 2017 and while progress was made, we are still working on making these
solutions scalable for our growing community.

{{raw (file "survey2017/project.html")}}

## Community

At the end of the survey, we asked some demographic questions.

The country distribution of responses is largely similar to last year with minor
fluctuations. Like last year, the distribution of countries is similar to the
visits to golang.org, though some Asian countries remain under-represented in
the survey.

Perhaps the most significant improvement over 2016 came from the question which
asked to what degree do respondents agreed with the statement, "I feel welcome
in the Go community". Last year the agreement to disagreement ratio was 15:1. In
2017 this ratio nearly doubled to 25:1.

An important part of a community is making everyone feel welcome, especially
people from under-represented demographics. We asked an optional question about
identification across a few underrepresented groups. We had a 4% increase in
response rate over last year. The percentage of each underrepresented group
increased over 2016, some quite significantly.

Like last year, we took the results of the statement “I feel welcome in the Go
community” and broke them down by responses to the various underrepresented
categories. Like the whole, most of the respondents who identified as
underrepresented also felt significantly more welcome in the Go community than
in 2016. Respondents who identified as a woman showed the most significant
improvement with an increase of over 400% in the ratio of agree:disagree to this
statement (3:1 → 13:1). People who identified as ethnically or racially
underrepresented had an increase of over 250% (7:1 → 18:1). Like last year,
those who identified as not underrepresented still had a much higher percentage
of agreement to this statement than those identifying from underrepresented
groups.

We are encouraged by this progress and hope that the momentum continues.

The final question on the survey was just for fun: what’s your favorite Go
keyword? Perhaps unsurprisingly, the most popular response was `go`, followed by
`defer`, `func`, `interface`, and `select`, unchanged from last year.

{{raw (file "survey2017/community.html")}}

Finally, on behalf of the entire Go project, we are grateful for everyone who
has contributed to our project, whether by being a part of our great community,
by taking this survey or by taking an interest in Go.
