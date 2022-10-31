---
title: Go 2016 Survey Results
date: 2017-03-06
by:
- Steve Francia, for the Go team
tags:
- survey
- community
summary: What we learned from the December 2017 Go User Survey.
---

## Thank you

This post summarizes the result of our December 2016 user survey along with our commentary and insights.
We are grateful to everyone who provided their feedback through the survey to help shape the future of Go.

## Programming background

Of the 3,595 survey respondents, 89% said they program in Go at work or outside of work,
with 39% using Go both at home and at work, 27% using Go only at home, and 23% using Go only at work.

We asked about the areas in which people work.
63% said they work in web development, but only 9% listed web development alone.
In fact, 77% chose two or more areas, and 53% chose three or more.

We also asked about the kinds of programs people write in Go.
63% of respondents write command-line programs, 60% write API or RPC services, and 52% write web services.
Like in the previous question, most made multiple choices,
with 85% choosing two or more and 72% choosing three or more.

We asked about people’s expertise and preference among programming languages.
Unsurprisingly, Go ranked highest among respondents’ first choices in both expertise (26%) and preference (62%).
With Go excluded, the top five first choices for language expertise were
Python (18%), Java (17%), JavaScript (13%), C (11%), and PHP (8%);
and the top five first choices for language preference were
Python (22%), JavaScript (10%), C (9%), Java (9%), and Ruby (7%).
Go is clearly attracting many programmers from dynamic languages.

{{raw (file "survey2016/background.html")}}

## Go usage

Users are overwhelmingly happy with Go:
they agree that they would recommend Go to others by a ratio of 19:1,
that they’d prefer to use Go for their next project (14:1),
and that Go is working well for their teams (18:1).
Fewer users agree that Go is critical to their company’s success (2.5:1).

When asked what they like most about Go, users most commonly mentioned
Go’s simplicity, ease of use, concurrency features, and performance.
When asked what changes would most improve Go,
users most commonly mentioned generics, package versioning, and dependency management.
Other popular responses were GUIs, debugging, and error handling.

When asked about the biggest challenges to their own personal use of Go,
users mentioned many of the technical changes suggested in the previous question.
The most common themes in the non-technical challenges were convincing others to use Go
and communicating the value of Go to others, including management.
Another common theme was learning Go or helping others learn,
including finding documentation like getting-started walkthroughs,
tutorials, examples, and best practices.

Some representative common feedback, paraphrased for confidentiality:

{{raw (file "survey2016/quotes.html")}}

We appreciate the feedback given to identify these challenges faced by our users and community.
In 2017 we are focusing on addressing these issues and hope to make as many significant improvements as we can.
We welcome suggestions and contributions from the community in making these challenges into strengths for Go.

{{raw (file "survey2016/usage.html")}}

## Development and deployment

When asked which operating systems they develop Go on,
63% of respondents say they use Linux, 44% use MacOS, and 19% use Windows,
with multiple choices allowed and 49% of respondents developing on multiple systems.
The 51% of responses choosing a single system split into
29% on Linux, 17% on MacOS, 5% on Windows, and 0.2% on other systems.

Go deployment is roughly evenly split between privately managed servers and hosted cloud servers.

{{raw (file "survey2016/dev.html")}}

## Working Effectively

We asked how strongly people agreed or disagreed with various statements about Go.
Users most agreed that Go’s performance meets their needs (57:1 ratio agree versus disagree),
that they are able to quickly find answers to their questions (20:1),
and that they are able to effectively use Go’s concurrency features (14:1).
On the other hand, users least agreed that they are able to effectively
debug uses of Go’s concurrency features (2.7:1).

Users mostly agreed that they were able to quickly find libraries they need (7.5:1).
When asked what libraries are still missing, the most common request by far was a library for writing GUIs.
Another popular topic was requests around data processing, analytics, and numerical and scientific computing.

Of the 30% of users who suggested ways to improve Go’s documentation,
the most common suggestion by far was more examples.

The primary sources for Go news are the Go blog,
Reddit’s /r/golang and Twitter;
there may be some bias here since these are also how the survey was announced.

The primary sources for finding answers to Go questions are the Go web site,
Stack Overflow, and reading source code directly.

{{raw (file "survey2016/effective.html")}}

## The Go Project

55% of respondents expressed interest in contributing in some way to the Go community and projects.
Unfortunately, relatively few agreed that they felt welcome to do so (3.3:1)
and even fewer felt that the process was clear (1.3:1).
In 2017, we intend to work on improving the contribution process and to
continue to work to make all contributors feel welcome.

Respondents agree that they are confident in the leadership of the Go project (9:1),
but they agree much less that the project leadership understands their needs (2.6:1),
and they agree even less that they feel comfortable approaching project leadership with questions and feedback (2.2:1).
In fact, these were the only questions in the survey for which more than half of respondents
did not mark “somewhat agree”, “agree”, or “strongly agree” (many were neutral or did not answer).

We hope that the survey and this blog post convey to those of you
who aren’t comfortable reaching out that the Go project leadership is listening.
Throughout 2017 we will be exploring new ways to engage with users to better understand their needs.

{{raw (file "survey2016/project.html")}}

## Community

At the end of the survey, we asked some demographic questions.
The country distribution of responses roughly matches the country distribution of site visits to golang.org,
but the responses under-represent some Asian countries.
In particular, India, China, and Japan each accounted for about 5% of the site visits to golang.org in 2016
but only 3%, 2%, and 1% of survey responses.

An important part of a community is making everyone feel welcome,
especially people from under-represented demographics.
We asked an optional question about identification across a few diversity groups.
37% of respondents left the question blank and 12% of respondents chose “I prefer not to answer”,
so we cannot make many broad conclusions from the data.
However, one comparison stands out: the 9% who identified as underrepresented agreed
with the statement “I feel welcome in the Go community” by a ratio of 7.5:1,
compared with 15:1 in the survey as a whole.
We aim to make the Go community even more welcoming.
We support and are encouraged by the efforts of organizations like GoBridge and Women Who Go.

The final question on the survey was just for fun: what’s your favorite Go keyword?
Perhaps unsurprisingly, the most popular response was `go`, followed by `defer`, `func`, `interface`, and `select`.

{{raw (file "survey2016/community.html")}}
