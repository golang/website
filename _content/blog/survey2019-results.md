---
title: Go Developer Survey 2019 Results
date: 2020-04-20
by:
- Todd Kulesza
tags:
- survey
- community
summary: An analysis of the results from the 2019 Go Developer Survey.
---

## What a response!

I want to start with an enormous **thank you** to the thousands of Go developers
who participated in this year’s survey.
For 2019, we saw 10,975 responses, nearly [twice as many as last year](https://blog.golang.org/survey2018-results)!
On behalf of the rest of the team, I cannot adequately stress how much we
appreciate you taking the time and effort to tell us about your experiences with Go. Thank you!

## A note about prior years

Sharp-eyed readers may notice that our year-over-year comparisons don’t
quite square with numbers we’ve shared in the past.
The reason is that from 2016–2018, we calculated percentages for each
question using the total number of people who started the survey as the denominator.
While that’s nice and consistent, it ignores the fact that not everyone
finishes the survey—up to 40% of participants stop before reaching the final page,
which meant questions that occurred later in the survey appeared to perform
worse solely because they were later.
Thus, this year we’ve recalculated all of our results (including the 2016–2018
responses shown in this post) to use the number of people who responded
to a given question as the denominator for that question.
We’ve included the number of 2019 responses for each chart—in the form
of "n=[number of respondents]" on the x-axis or in the chart’s legend—to
give readers a better understanding of the weight of evidence underlying each finding.

Similarly, we learned that in prior surveys options that appeared earlier
in response lists had a disproportionate response rate.
To address this, we added an element of randomization into the survey.
Some of our multiple-choice questions have lists of choices with no logical ordering,
such as "I write the following in Go:
[list of types of applications]".
Previously these choices had been alphabetized,
but for 2019 they were presented in a random order to each participant.
This means year-over-year comparison for certain questions are invalid for 2018 → 2019,
but trends from 2016–2018 are not invalidated.
You can think of this as setting a more accurate baseline for 2019.
We retained alphabetical ordering in cases where respondents are likely
to scan for a particular name,
such as their preferred editor.
We explicitly call out which questions this applies to below.

A third major change was to improve our analysis of questions with open-ended,
free-text responses.
Last year we used machine learning to roughly—but quickly—categorize these responses.
This year two researchers manually analyzed and categorized these responses,
allowing for a more granular analysis but preventing valid comparisons with
last year’s numbers.
Like the randomization discussed above, the purpose of this change is to
give us a reliable baseline for 2019 onward.

## Without further ado…

This is a long post. Here’s the tl;dr of our major findings:

- The demographics of our respondents are similar to Stack Overflow’s survey respondents,
  which increases our confidence that these results are representative of
  the larger Go developer audience.
- A majority of respondents use Go every day, and this number has been trending up each year.
- Go’s use is still concentrated in technology companies,
  but Go is increasingly found in a wider variety of industries,
  such as finance and media.
- Methodology changes showed us that most of our year-over-year metrics
  are stable and higher than we previously realized.
- Respondents are using Go to solve similar problems,
  particularly building API/RPC services and CLIs,
  regardless of the size of organization they work at.
- Most teams try to update to the latest Go release quickly;
  when third-party providers are late to support the current Go release,
  this creates an adoption blocker for developers.
- Almost everyone in the Go ecosystem is now using modules, but some confusion around package management remains.
- High-priority areas for improvement include improving the developer experience for debugging,
  working with modules, and working with cloud services.
- VS Code and GoLand have continued to see increased use; they’re now preferred by 3 out of 4 respondents.

## Who did we hear from?

This year we asked some new demographic questions to help us better understand
the people who’ve responded to this survey.
In particular, we asked about the duration of professional programming experience
and the size of the organizations where people work.
These were modeled on questions that StackOverflow asks in their annual survey,
and the distribution of responses we saw is very close to StackOverflow’s 2019 results.
Our take-away is the respondents to this survey have similar levels of professional
experience and proportional representation of different sizes of organizations
as the StackOverflow survey audience (with the obvious difference that we’re
primarily hearing from developers working with Go).
That increases our confidence when generalizing these findings to the estimated
1 million Go developers worldwide.
These demographic questions will also help us in the future to identify
which year-over-year changes may be the result of a shift in who responded to the survey,
rather than changes in sentiment or behavior.

{{image "survey2019/fig1.svg" 700}}
{{image "survey2019/fig2.svg" 700}}

Looking at Go experience, we see that a majority of respondents (56%) are
relatively new to Go,
having used it for less than two years.
Majorities also said they use Go at work (72%) and outside of work (62%).
The percentage of respondents using Go professionally appears to be trending up each year.

As you can see in the chart below, in 2018 we saw a spike in these numbers,
but that increase disappeared this year.
This is one of many signals suggesting that the audience who answered the
survey in 2018 was significantly different than in the other three years.
In this case they were significantly more likely to be using Go outside
of work and a different language while at work,
but we see similar outliers across multiple survey questions.

{{image "survey2019/fig3.svg" 700}}
{{image "survey2019/fig4.svg" 700}}

Respondents who have been using Go the longest have different backgrounds
than newer Go developers.
These Go veterans were more likely to claim expertise in C/C++ and less
likely to claim expertise in JavaScript,
TypeScript, and PHP.
One caveat is that this is self-reported "expertise";
it may be more helpful to think of it instead as "familiarity".
Python appears to be the language (other than Go) familiar to the most respondents,
regardless of how long they’ve been working with Go.

{{image "survey2019/fig5.svg" 750}}

Last year we asked about which industries respondents work in,
finding that a majority reported working in software,
internet, or web service companies.
This year it appears respondents represent a broader range of industries.
However, we also simplified the list of industries to reduce confusion from
potentially overlapping categories (e.g.,
the separate categories for "Software" and "Internet / web services" from
2018 were combined into "Technology" for 2019).
Thus, this isn’t strictly an apples-to-apples comparison.
For example, it’s possible that one effect of simplifying the category list
was to reduce the use of the "Software" category as a catch-all for respondents
writing Go software for an industry that wasn’t explicitly listed.

{{image "survey2019/fig6.svg" 700}}

Go is a successful open-source project, but that doesn’t mean the developers
working with it are also writing free or open-source software.
As in prior years, we found that most respondents are not frequent contributors
to Go open-source projects,
with 75% saying they do so "infrequently" or "never".
As the Go community expands, we see the proportion of respondents who’ve
never contributed to Go open-source projects slowly trending up.

{{image "survey2019/fig7.svg" 700}}

## Developer tools

As in prior years, the vast majority of survey respondents reported working
with Go on Linux and macOS systems.
This is one area of strong divergence between our respondents and StackOverflow’s 2019 results:
in our survey, only 20% of respondents use Windows as a primary development platform,
while for StackOverflow it was 45% of respondents.
Linux is used by 66% and macOS by 53%—both much higher than the StackOverflow audience,
which reported 25% and 30%, respectively.

{{image "survey2019/fig8.svg" 700}}
{{image "survey2019/fig9.svg" 700}}

The trend in editor consolidation has continued this year.
GoLand saw the sharpest increase in use this year,
rising from 24% → 34%.
VS Code’s growth slowed, but it remains the most popular editor among respondents at 41%.
Combined, these two editors are now preferred by 3 out of 4 respondents.

Every other editor saw a small decrease. This doesn’t mean those editors
aren’t being used at all,
but they’re not what respondents say they _prefer_ to use for writing Go code.

{{image "survey2019/fig10.svg" 700}}

This year we added a question about internal Go documentation tooling,
such as [gddo](https://github.com/golang/gddo).
A small minority of respondents (6%) reported that their organization runs
its own Go documentation server,
though this proportion nearly doubles (to 11%) when we look at respondents
at large organizations (those with at least 5,000 employees).
A follow-up asked of respondents who said their organization had stopped
running its own documentation server suggests that the top reason to retire
their server was a combination of low perceived benefits (23%) versus the
amount of effort required to initially set it up and maintain it (38%).

{{image "survey2019/fig11.svg" 700}}

## Sentiments towards Go

Large majorities of respondents agreed that Go is working well for their
teams (86%) and that they’d prefer to use it for their next project (89%).
We also found that over half of respondents (59%) believe Go is critical
to the success of their companies.
All of these metrics have remained stable since 2016.

Normalizing the results changed most of these numbers for prior years.
For example, the percentage of respondents who agreed with the statement
"Go is working well for my team" was previously in the 50’s and 60’s because
of participant drop-off;
when we remove participants who never saw the question,
we see it’s been fairly stable since 2016.

{{image "survey2019/fig12.svg" 700}}

Looking at sentiments toward problem solving in the Go ecosystem,
we see similar results.
Large percentages of respondents agreed with each statement (82%–88%),
and these rates have been largely stable over the past four years.

{{image "survey2019/fig13.svg" 700}}

This year we took a more nuanced look at satisfaction across industries
to establish a baseline.
Overall, respondents were positive about using Go at work,
regardless of industry sector.
We do see small variations in dissatisfaction in a few areas,
most notably manufacturing, which we plan to investigate with follow-up research.
Similarly, we asked about satisfaction with—and the importance of—various
aspects of Go development.
Pairing these measures together highlighted three topics of particular focus:
debugging (including debugging concurrency),
using modules, and using cloud services.
Each of these topics was rated "very" or "critically" important by a majority
of respondents but had significantly lower satisfaction scores compared to other topics.

{{image "survey2019/fig14.svg" 800}}
{{image "survey2019/fig15.svg" 750}}

Turning to sentiments toward the Go community,
we see some differences from prior years.
First, there is a dip in the percentage of respondents who agreed with the
statement "I feel welcome in the Go community", from 82% to 75%.
Digging deeper revealed that the proportion of respondents who "slightly"
or "moderately agreed" decreased,
while the proportions who "neither agree nor disagree" and "strongly agree"
both increased (up 5 and 7 points, respectively).
This polarizing split suggests two or more groups whose experiences in the
Go community are diverging,
and is thus another area we plan to further investigate.

The other big differences are a clear upward trend in responses to the statement
"I feel welcome to contribute to the Go project" and a large year-over-year
increase in the proportion of respondents who feel Go’s project leadership
understands their needs.

All of these results show a pattern of higher agreement correlated with
increased Go experience,
beginning at about two years.
In other words, the longer a respondent has been using Go,
the more likely they were to agree with each of these statements.

{{image "survey2019/fig16.svg" 700}}

This likely comes as no surprise, but people who responded to the Go Developer
Survey tended to like Go.
However, we also wanted to understand which _other_ languages respondents enjoy working with.
Most of these numbers have not significantly changed from prior years,
with two exceptions:
TypeScript (which has increased 10 points),
and Rust (up 7 points).
When we break these results down by duration of Go experience,
we see the same pattern as we found for language expertise.
In particular, Python is the language and ecosystem that Go developers are
most likely to also enjoy building with.

{{image "survey2019/fig17.svg" 700}}

In 2018 we first asked the "Would you recommend…" [Net Promoter Score](https://en.wikipedia.org/wiki/Net_Promoter) (NPS) question,
yielding a score of 61.
This year our NPS result is a statistically unchanged 60 (67% "promoters"
minus 7% "detractors").

{{image "survey2019/fig18.svg" 700}}

## Working with Go

Building API/RPC services (71%) and CLIs (62%) remain the most common uses of Go.
The chart below appears to show major changes from 2018,
but these are most likely the result of randomizing the order of choices,
which used to be listed alphabetically:
3 of the 4 choices beginning with ’A’ decreased,
while everything else remained stable or increased.
Thus, this chart is best interpreted as a more accurate baseline for 2019
with trends from 2016–2018.
For example, we believe that the proportion of respondents building web
services which return HTML has been decreasing since 2016 but were likely
undercounted because this response was always at the bottom of a long list of choices.
We also broke this out by organization size and industry but found no significant differences:
it appears respondents use Go in roughly similar ways whether they work
at a small tech start-up or a large retail enterprise.

A related question asked about the larger areas in which respondents work with Go.
The most common area by far was web development (66%),
but other common areas included databases (45%),
network programming (42%), systems programming (38%),
and DevOps tasks (37%).

{{image "survey2019/fig19.svg" 700}}
{{image "survey2019/fig20.svg" 700}}

In addition to what respondents are building,
we also asked about some of the development techniques they use.
A large majority of respondents said they depend upon text logs for debugging (88%),
and their free-text responses suggest this is because alternative tooling
is challenging to use effectively.
However, local stepwise debugging (e.g., with Delve),
profiling, and testing with the race detector were not uncommon,
with ~50% of respondents depending upon at least one of these techniques.

{{image "survey2019/fig21.svg" 700}}

Regarding package management, we found that the vast majority of respondents
have adopted modules for Go (89%).
This has been a big shift for developers,
and nearly the entire community appears to be going through it simultaneously.

{{image "survey2019/fig22.svg" 700}}

We also found that 75% of respondents evaluate the current Go release for production use,
with an additional 12% waiting one release cycle.
This suggests a large majority of Go developers are using (or at the least,
trying to use) the current or previous stable release,
highlighting the importance for platform-as-a-service providers to quickly
support new stable releases of Go.


{{image "survey2019/fig23.svg" 700}}

## Go in the clouds

Go was designed with modern distributed computing in mind,
and we want to continue to improve the developer experience of building
cloud services with Go.
This year we expanded the questions we asked about cloud development to
better understand how respondents are working with cloud providers,
what they like about the current developer experience,
and what can be improved.
As mentioned earlier, some of the 2018 results appear to be outliers,
such as an unexpectedly low result for self-owned servers,
and an unexpectedly high result for GCP deployments.

We see two clear trends:

1. The three largest global cloud providers (Amazon Web Services,
Google Cloud Platform, and Microsoft Azure) all appear to be trending up
in usage among survey respondents,
while most other providers are used by a smaller proportion of respondents each year.
2. On-prem deployments to self-owned or company-owned servers continue to
decrease and are now statistically tied with AWS (44% vs.
42%) as the most common deployment targets.

Looking at which types of cloud platforms respondents are using,
we see differences between the major providers.
Respondents deploying to AWS and Azure were most likely to be using VMs
directly (65% and 51%,
respectively), while those deploying to GCP were almost twice as likely
to be using the managed Kubernetes platform (GKE,
64%) than VMs (35%).
We also found that respondents deploying to AWS were equally likely to be
using a managed Kubernetes platform (32%) as they were to be using a managed
serverless platform (AWS Lambda, 33%).
Both GCP (17%) and Azure (7%) had lower proportions of respondents using
serverless platforms,
and free-text responses suggest a primary reason was delayed support for
the latest Go runtime on these platforms.

Overall, a majority of respondents were satisfied with using Go on all three
major cloud providers.
Respondents reported similar satisfaction levels with Go development for
AWS (80% satisfied) and GCP (78%).
Azure received a lower satisfaction score (57% satisfied),
and free-text responses suggest that the main driver was a perception that
Go lacks first-class support on this platform (25% of free-text responses).
Here, "first-class support" refers to always staying up-to-date with the latest Go release,
and ensuring new features are available to Go developers at time of launch.
This was the same top pain-point reported by respondents using GCP (14%),
and particularly focused on support for the latest Go runtime in serverless deployments.
Respondents deploying to AWS, in contrast,
were most likely to say the SDK could use improvements,
such as being more idiomatic (21%).
SDK improvements were also the second most common request for both GCP (9%)
and Azure (18%) developers.

{{image "survey2019/fig24.svg" 700}}
{{image "survey2019/fig25.svg" 700}}
{{image "survey2019/fig26.svg" 700}}

## Pain points

The top reasons respondents say they are unable to use Go more remain working
on a project in another language (56%),
working on a team that prefers to use another language (37%),
and the lack of a critical feature in Go itself (25%).

This was one of the questions where we randomized the choice list,
so year-over-year comparisons aren’t valid,
though 2016–2018 trends are.
For example, we are confident that the number of developers unable to use
Go more frequently because their team prefers a different language is decreasing each year,
but we don’t know whether that decrease dramatically accelerated this year,
or was always a bit lower than our 2016–2018 numbers estimated.

{{image "survey2019/fig27.svg" 700}}

The top two adoption blockers (working on an existing non-Go project and
working on a team that prefers a different language) don’t have direct technical solutions,
but the remaining blockers might.
Thus, this year we asked for more details,
to better understand how we might help developers increase their use of Go.
The charts in the remainder of this section are based on free-text responses
which were manually categorized,
so they have _very_ long tails;
categories totalling less than 3% of the total responses have been grouped
into the "Other" category for each chart.
A single response may mention multiple topics,
thus charts do not sum to 100%.

Among the 25% of respondents who said Go lacks language features they need,
79% pointed to generics as a critical missing feature.
Continued improvements to error handling (in addition to the Go 1.13 changes) was cited by 22%,
while 13% requested more functional programming features,
particularly built-in map/filter/reduce functionality.
To be clear, these numbers are from the subset of respondents who said they
would be able to use Go more were it not missing one or more critical features they need,
not the entire population of survey respondents.


{{image "survey2019/fig28.svg" 700}}

Respondents who said Go "isn’t an appropriate language" for what they work
on had a wide variety of reasons and use-cases.
The most common was that they work on some form of front-end development (22%),
such as GUIs for web, desktop, or mobile.
Another common response was that the respondent said they worked in a domain
with an already-dominant language (9%),
making it a challenge to use something different.
Some respondents also told us which domain they were referring to (or simply
mentioned a domain without mentioning another language being more common),
which we show via the "I work on [domain]" rows below.
An additional top reason cited by respondents was a need for better performance (9%),
particularly for real-time computing.

{{image "survey2019/fig29.svg" 700}}

The biggest challenges respondents reported remain largely consistent with last year.
Go’s lack of generics and modules/package management still top the list
(15% and 12% of responses,
respectively), and the proportion of respondents highlighting tooling problems increased.
These numbers are different from the above charts because this question
was asked of _all_ respondents,
regardless of what they said their biggest Go adoption blockers were.
All three of these are areas of focus for the Go team this year,
and we hope to greatly improve the developer experience,
particularly around modules, tooling, and the getting started experience,
in the coming months.

{{image "survey2019/fig30.svg" 700}}

Diagnosing faults and performance issues can be challenging in any language.
Respondents told us their top challenge for both of these was not something
specific to Go’s implementation or tooling,
but a more fundamental issue:
a self-reported lack of knowledge, experience, or best practices.
We hope to help address these knowledge gaps via documentation and other
educational materials later this year.
The other major problems do involve tooling,
specifically a perceived unfavorable cost/benefit trade-off to learning/using
Go’s debugging and profiling tooling,
and challenges making the tooling work in various environments (e.g.,
debugging in containers, or getting performance profiles from production systems).

{{image "survey2019/fig31.svg" 700}}
{{image "survey2019/fig32.svg" 700}}

Finally, when we asked what would most improve Go support in respondents’
editing environment,
the most common response was for general improvements or better support
for the language server (gopls, 19%).
This was expected, as gopls replaces about 80 extant tools and is still in beta.
When respondents were more specific about what they’d like to see improved,
they were most likely to report the debugging experience (14%) and faster
or more reliable code completion (13%).
A number of participants also explicitly referenced the need to frequently
restart VS Code when using gopls (8%);
in the time since this survey was in the field (late November – early December 2019),
many of these gopls improvements have already landed,
and this continues to be a high-priority area for the team.

{{image "survey2019/fig33.svg" 700}}

## The Go community

Roughly two thirds of respondents used Stack Overflow to answer their Go-related questions (64%).
The other top sources of answers were godoc.org (47%),
directly reading source code (42%), and golang.org (33%).

{{image "survey2019/fig34.svg" 700}}

The long tail on the previous chart highlights the large variety of different
sources (nearly all of them community-driven) and modalities that respondents
rely on to overcome challenges while developing with Go.
Indeed, for many Gophers, this may be one of their main points of interaction
with the larger community:
as our community expands, we’ve seen higher and higher proportions of respondents
who do not attend any Go-related events.
For 2019, that proportion nearly reached two thirds of respondents (62%).

{{image "survey2019/fig35.svg" 700}}

Due to updated Google-wide privacy guidelines,
we can no longer ask about which countries respondents live in.
Instead we asked about preferred spoken/written language as a very rough
proxy for Go’s worldwide usage,
with the benefit of providing data for potential localization efforts.

Because this survey is in English, there is likely a strong bias toward
English speakers and people from areas where English is a common second or third language.
Thus, the non-English numbers should be interpreted as likely minimums rather
than an approximation of Go’s global audience.

{{image "survey2019/fig36.svg" 700}}

We found 12% of respondents identify with a traditionally underrepresented group (e.g.,
ethnicity, gender identity, et al.) and 3% identify as female.
(This question should have said "woman" instead of "female".
The mistake has been corrected in our draft survey for 2020,
and we apologize for it.)
We strongly suspect this 3% is undercounting women in the Go community.
For example, we know women software developers in the US respond to the
StackOverflow Developer Survey at [about half the rate we’d expect based on US employment figures](https://insights.stackoverflow.com/survey/2019#developer-profile-_-developer-type) (11% vs 20%).
Since we don’t know the proportion of responses in the US,
we can’t safely extrapolate from these numbers beyond saying the actual
proportion is likely higher than 3%.
Furthermore, GDPR required us to change how we ask about sensitive information,
which includes gender and traditionally underrepresented groups.
Unfortunately these changes prevent us from being able to make valid comparisons
of these numbers with prior years.

Respondents who identified with underrepresented groups or preferred not
to answer this question showed higher rates of disagreement with the statement
"I feel welcome in the Go community" (8% vs.
4%) than those who do not identify with an underrepresented group,
highlighting the importance of our continued outreach efforts.

{{image "survey2019/fig37.svg" 700}}
{{image "survey2019/fig38.svg" 700}}
{{image "survey2019/fig39.svg" 800}}

## Conclusion

We hope you’ve enjoyed seeing the results of our 2019 developer survey.
Understanding developers’ experiences and challenges helps us plan and prioritize work for 2020.
Once again, an enormous thank you to everyone who contributed to this survey—your
feedback is helping to steer Go’s direction in the coming year and beyond.
