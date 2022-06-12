---
title: Go Developer Survey 2021 Results
date: 2022-04-19
by:
- Alice Merrick
tags:
- survey
- community
summary: An analysis of the results from the 2021 Go Developer Survey.
---


## A huge thank you to the community for the amazing response!

In 2021, we ran the Go Developer Survey from Oct 26th to Nov 16th and had 11,840 responses—the largest turnout so far in the 6 years we’ve been running the survey! Thank you for putting in the time to provide the community with these insights on your experiences using Go. 

## Highlights {#highlights}

- Most responses were consistent with prior years. For example, [satisfaction with Go is still very high at 92%](#satisfaction) and 75% of respondents use Go at work.
- This year we [randomly sampled](#changes) some participants using the Go VS Code plugin, which resulted in some shifts in who responded to the survey.
- Missing critical libraries, language features and infrastructure were the most [common barriers to using Go](#adoption). (note: this survey ran prior to the release of Go 1.18 with generics, the #1 most reported missing feature)
- Respondents want to [prioritize improvements](#prioritization) to debugging and dependency management.
- The [biggest challenges when using modules](#modules) concerned versioning, using private repos and multi-module workflows. (note: this survey was run prior to Go 1.18 which introduced workspaces addressing many of these concerns).
- 81% of respondents are [confident in the long term direction of the Go project](#satisfaction).

## Who did we hear from? {#demographics}

Our demographics have been pretty stable year over year ([See 2020 results](https://go.dev/blog/survey2020-results)). Consistent with previous years, Go is primarily used in the tech industry. 70% of respondents were software developers, with a few in IT or DevOps and 76% of respondents said they program in Go at work. 
<img src="survey2021/industry_yoy.svg" alt="Bar chart of industries where respondents work" width="700"/>
<img src="survey2021/where_yoy.svg" alt="Bar chart showing Go used more at work than outide of work" width="700"/>
<img src="survey2021/app_yoy.svg" alt="Bar chart of uses for Go where API/RPC services and CLI apps are most common" width="700"/>

Some new demographics from 2021:
* Most respondents describe their organization as an enterprise or small to medium business, with about a quarter describing their organization as a startup. Consultancies and public institutions were much less common. 
* The vast majority of respondents work on teams of less than ten people. 
* Over half (55%) of respondents use Go at work on a daily basis. Respondents use Go less frequently outside of work.

<img src="survey2021/orgtype.svg" alt="Bar chart of organization type where enterprise is the most common response" width="700"/>

<img src="survey2021/teamsize.svg" alt="Bar chart of team size where 2 to 5 is the most common size" width="700"/>

<img src="survey2021/gofreq_comparison.svg" alt="Frequency of using Go at work versus outside of work where using Go at work is most often on a daily basis and outside of work is less common and most often on a weekly basis" width="700"/>

### Gender identity {#gender}
We ask about gender identity on the survey because it gives us an idea of who is being represented in the results and adds another dimension to measure the inclusivity of the community. The Go team values diversity and inclusion, not only because it’s the right thing to do, but because diverse voices help us make better decisions. This year we rephrased the gender identity question to be more inclusive of other gender identities. About the same proportion identified as women as previous years (2%). This was true in the [randomly sampled group](#changes) as well, suggesting this is not just due to sampling. 
<img src="survey2021/gender.svg" alt="Bar chart showing gender identity of respondents where 92% of respondents identify as male" width="700"/>

### Assistive technology

This year we again found that about 8% of respondents are using some form of assistive technology.  Most challenges concerned a need for higher contrast themes and increased font sizes on Go-related websites or in their code editors; we’re planning to act on the website feedback later this year. These accessibility needs are something we should all keep in mind when contributing to the Go ecosystem.

## A closer look at challenges to Go adoption {#adoption}
This year we revised our questions to target actual cases where Go wasn’t adopted and why. First, we asked whether or not respondents had evaluated using another language against Go in the last year. 43% of respondents said they had either evaluated switching to Go, from Go, or adopting Go when there wasn’t a previously established language. 80% of these evaluations were primarily for business reasons.

<img src="survey2021/evaluated.svg" alt="Chart showing proportion of respondents who evaluated Go against another language in the last year" width="700"/>

We expected the most common use cases for Go would be the most common intended uses for those evaluating Go. API/RPC services was by far the most common use, but surprisingly, data processing was the second most common intended use case.

<img src="survey2021/intended_app.svg" alt="Chart showing the kind application they considered using Go" width="700"/>

Of those respondents who evaluated Go, 75% ended up using Go. (Of course, since nearly all survey respondents report using Go, we likely are not hearing from developers who evaluated Go and decided against using it.) 

<img src="survey2021/adopted.svg" alt="Chart showing proportion who used Go compared to those who stayed with the current language or chose another language" width="700"/>

To those who evaluated Go and didn’t use it, we then asked what challenges prevented them from using Go and which of those was the primary barrier. 
<img src="survey2021/blockers.svg" alt="Chart showing barriers to using Go" width="700"/>

The picture we get from these results supports previous findings that missing features and lack of ecosystem / library support are the most significant technical barriers to Go adoption. 

We asked for more details on what features or libraries respondents were missing and found that generics was the most common critical missing feature—we expect this to be a less significant barrier after the introduction of generics in Go 1.18. The next most common missing features had to do with Go’s type system. We would like to see how introducing generics may influence or resolve underlying needs around Go’s type system before making additional changes. For now, we will gather more information on the contexts for these needs and may in the future explore different ways to meet those needs such as through tooling, libraries or changes to the type system.

As for missing libraries, there was no clear consensus on what addition would unblock the largest proportion of those wanting to adopt Go. That will require additional exploration.

So what did respondents use instead when they didn’t choose Go?

<img src="survey2021/lang_instead.svg" alt="Chart of which languages respondents used instead of Go" width="700"/>

 Rust, Python, and Java are the most common choices. [Rust and Go have complementary feature sets](https://thenewstack.io/rust-vs-go-why-theyre-better-together/), so Rust may be a good option for when Go doesn’t meet feature needs for a project. The primary reasons for using Python were missing libraries and existing infrastructure support, so Python’s large package ecosystem may make it difficult to switch to Go. Similarly, the most common reason for using Java instead was because of Go’s missing features, which may be alleviated by the introduction of generics in the 1.18 release.

## Go satisfaction and prioritization {#satisfaction}
Let’s look at areas where Go is doing well and where things can be improved. 

Consistent with last year, 92% of respondents said they were very or somewhat satisfied using Go during the past year. 

<img src="survey2021/csat.svg" alt="Overall satisfaction on a 5 points scale from very dissatisfied to very satisfied" width="700"/>

Year over year trends in community attitudes have seen minor fluctuations. Those using Go for less than 3 months tend to be less likely to agree with these statements. Respondents are increasingly finding Go critical for their company's success. 

<img src="survey2021/attitudes_yoy.svg" alt="Attitudes around using Go at work" width="700"/>
<img src="survey2021/attitudes_community_yoy.svg" alt="Community attitudes around welcomeness and confidence in direction of the Go project" width="700"/>

### Prioritization {#prioritization}
The last few years we’ve asked respondents to rate specific areas on how satisfied they are and how important those areas are to them; we use this information to identify areas that are important to respondents, but with which they are unsatisfied. However, most of these areas have shown only minor differences in both importance and satisfaction.

<img src="survey2021/imp_vs_sat2.svg" alt="Scatter plot of importance compared to satisfaction showing most areas have high satisfaction and where binary size is less important than other areas" width="700"/>

This year we introduced a new question to explore alternative ways to prioritize work on specific areas. “Let’s say you have 10 GopherCoins to spend on improving the following aspects of working with Go. How would you distribute your coins?” Two areas that stood out as receiving significantly more GopherCoins were dependency management (using modules) and diagnosing bugs, areas that we’ll be dedicating resources during 2022.

<img src="survey2021/improvements_sums.svg" alt="Overall sum of coins spent on each area for improvement" width="700"/>

### Challenges when working with modules {#modules}
The most common module-related challenge was working across multiple modules (19% of respondents), followed by comments about versioning (including trepidation around committing to a stable v1 API). Related to versioning, 9% of responses discussed version management or updating dependencies. Rounding out the top 5 were challenges around private repos (including authentication with GitLab in particular) and remembering the different `go mod` commands plus understanding their error messages.

## Learning Go {#learning}

This year we adopted a new construct to explore relative productivity among different levels of experience with Go. The vast majority of respondents (88%) agree that they regularly reach a high level of productivity and 85% agree they’re often able to achieve a flow state when writing in Go. The proportion of agreement increases as experience with Go increases.

<img src="survey2021/productivity.svg" alt="Charts showing proportion of respondents who agree they feel productive using Go and can achieve a state of flow while writing in Go" width="700"/>

### In which areas should we invest in best practice docs?

Half of respondents wanted more guidance on best practices on performance optimization and project directory structure. Unsurprisingly, new Gophers (using Go for less than 1 year) need more guidance than more experienced Gophers, though the top areas were consistent across both groups. Notably, new Gophers asked for more guidance in concurrency than more experienced Gophers.

<img src="survey2021/best_practices.svg" alt="Chart showing which areas respondents want more guidance on best practices" width="700"/>

### How do developers learn a new language?
About half of respondents learned a new language at work, but almost as many (45%) learn outside of school or work. Respondents most often (90%) reported learning alone. Of those who said they learned at work, where there may be opportunities to learn as a group, 84% learned alone rather than as a group.

<img src="survey2021/learn_where.svg" alt="Chart showing half of respondents learned a new language at work while 45% learned a new language outside of school or work" width="700"/>
<img src="survey2021/learn_with.svg" alt="Chart showing 90% of respondents learned their last new language alone" width="700"/>

Many of the top resources highlight the importance of good documentation, but live instruction stands out as a particularly useful resource for language learning as well. 

<img src="survey2021/learning_resources.svg" alt="Chart showing which resources are most helpful for learning a new programming language where reading reference docs and written tutorials are most useful" width="700"/>

## Developer tools and practices {#devtools}
As in prior years, the vast majority of survey respondents reported working
with Go on Linux (63%) and macOS (55%) systems. The proportion of respondents who primarily develop on Linux appears to be slightly trending down over time.

<img src="survey2021/os_yoy.svg" alt="Primary operating system from 2019 to 2021" width="700"/>

### Targeted platforms
Over 90% of respondents target Linux! Even though more respondents develop on macOS than Windows, they more often deploy to Windows than macOS. 

<img src="survey2021/os_deploy.svg" alt="Chart showing which platforms respondents deploy their Go code on " width="700"/>

### Fuzzing
Most respondents are unfamiliar with fuzzing or still consider themselves new to fuzzing. Based on this finding, we plan to 1) ensure Go’s fuzzing documentation explains fuzzing concepts in addition to the specifics of fuzzing in Go, and 2) design output and error messages to be actionable, so as to help developers who are new to fuzzing apply it successfully. 

<img src="survey2021/fuzz.svg" alt="Chart showing proportion of respondents who have used fuzzing" width="700"/>

## Cloud computing {#cloud}

Go was designed with modern distributed computing in mind, and we want to continue to improve the developer experience of building cloud services with Go. The proportion of respondents deploying Go programs to three largest global cloud providers (Amazon Web Services, Google Cloud Platform, and Microsoft Azure) remained about the same this year and on-prem deployments to self-owned or company-owned servers continue to decrease.

<img src="survey2021/cloud_yoy.svg" alt="Bar chart of cloud providers used to deploy Go programs where AWS is the most common at 44%" width="700"/>

Respondents deploying to AWS saw increases in deploying to a managed Kubernetes platform, now at 35% of those who deploy to any of the three largest cloud providers. All of these cloud providers saw a drop in the proportion of users deploying Go programs to VMs. 


<img src="survey2021/cloud_services_yoy.svg" alt="Bar charts of proportion of services being used with each provider" width="700"/>

## Changes this year {#changes}

Last year we introduced a [modular survey design](https://go.dev/blog/survey2020-results) so that we could ask more questions without lengthening the survey. We continued the modular design this year, although some questions were discontinued and others were added or modified. No respondents saw all the questions on the survey. Additionally, some questions may have much smaller sample sizes because they were asked selectively based on a previous question. 

The most significant change to the survey this year was in how we recruited participants. In previous years, we announced the survey through the Go Blog, where it was picked up on various social channels like Twitter, Reddit, or Hacker News. This year, in addition to the traditional channels, we used the VS Code Go plugin to randomly select users  to be shown a prompt asking if they’d like to participate in the survey. This created a random sample that we used to compare the self-selected respondents from our traditional channels and helped identify potential effects of [self-selection bias](https://en.wikipedia.org/wiki/Self-selection_bias). 

<img src="survey2021/rsamp.svg" alt="Proportion of respondents from each source" width="700"/>

Almost a third of our respondents were sourced this way so their responses had the potential to significantly impact the responses we saw this year. Some of the key differences we see between these two groups are:

### More new Gophers
The randomly selected sample had a higher proportion of new Gophers (those using Go for less than a year). It could be that new Gophers are less plugged into the Go ecosystem or the social channels,  so they were more likely to see the survey advertised in their IDE than find it through other means. Regardless of the reason, it's great to hear from a wider slice of the Go community.

<img src="survey2021/goex_s.svg" alt="Comparison of proportion of respondents with each level of experience for randomly sampled versus self-selected groups" width="700"/>

### More VS Code users
It’s unsurprising 91% of respondents who came to the survey from the VS Code plugin prefer to use VS Code when using Go. As a result, we saw much higher editor preferences for VS Code this year. When we exclude the random sample, the results are not statistically different from last year, so we know this is the result of the change in our sample and not overall preference. Similarly, VS Code users are also more likely to develop on Windows than other respondents so we see a slight increase in the preference for Windows this year. We also saw slight shifts in the usage of certain developer techniques that are common with VS Code editor usage. 

<img src="survey2021/editor_s.svg" alt="Grouped bar chart of which editor respondents prefer from each sample group" width="700"/>

<img src="survey2021/os_s.svg" alt="Grouped bar chart of primary operating system respondents use to develop go on" width="700"/>

<img src="survey2021/devtech_s.svg" alt="Grouped bar chart showing which techniques respondents use when writing in Go" width="700"/>

### Different resources
The randomly selected sample was less likely to rate social channels like the Go Blog as among their top resources for answering Go-related questions, so they may have been less likely to see the survey advertised on those channels.

 <img src="survey2021/resources_s.svg" alt="Grouped bar chart showing the top resources respondents use when writing in Go" width="700"/>

## Conclusion {#conclusion}

Thank you for joining us in reviewing the results of our 2021 developer survey! To reiterate, some key takeaways:

* Most of our year over year metrics remained stable with most changes owing to our change in sample.
* Satisfaction with Go remains high!
* Three-quarters of respondents use Go at work and many use Go on a daily basis so helping you get work done is a top priority.
* We will prioritize improvements to debugging and dependency management workflows.
* We will continue to work towards making Go an inclusive community for all kinds of Gophers.

Understanding developers’ experiences and challenges helps us measure our progress and directs the future of Go. Thanks again to everyone who contributed to this survey—we couldn't have done it without you. We hope to see you next year!


