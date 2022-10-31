---
title: Contribution Workshop
date: 2017-08-09
by:
- Steve Francia
- Cassandra Salisbury
- Matt Broberg
- Dmitri Shuralyov
tags:
- community
summary: The Go contributor workshop trained new contributors at GopherCon.
---

## Event Overview

by [Steve](https://twitter.com/spf13)

During the community day at GopherCon, the Go team held two workshops
where we worked with people to help them make their first contribution to the
Go project. This was the first time the Go project has ever attempted anything
like this. We had about 140 participants and about 35 people who volunteered as
mentors. Mentors not only received warm fuzzy feelings for helping others, but
also a very stylish Go Mentor trucker hat. We had contributors of all
ages and experience levels coming from North and South America, Africa, Europe,
Asia, and Australia. It was truly a worldwide effort of Gophers coming together
at GopherCon.

One of our reasons for running the workshop was for it to act as a forcing
function to make us improve our contributor experience. In preparation for the
workshop, we rewrote our contributor guide, including adding a "troubleshooting"
section and built a tool `go-contrib-init`, which automated the process of
setting up a development environment to be able to contribute to Go.

For the workshop itself, we developed a presentation _"Contributing to Go,"_
and a dashboard / scoreboard that was presented during the event. The
scoreboard was designed to encourage us all to work together towards a common
goal of seeing our collective score increase. Participants added 1, 2 or 3 points to
the total score when they performed actions like registering an account, making
a change list (also known as a CL, similar to a pull request),
amending a CL, or submitting a CL.

{{image "contributor-workshop/image17.png"}}

Brad Fitzpatrick, who stayed home from GopherCon this year, was ready and
waiting to review all CLs submitted. He was so quick to review that many people
thought he was an automated bot. Internally our team is now calling him
"BradBot" mostly because we are in awe and a bit jealous.

{{image "contributor-workshop/image9.jpg"}}
{{image "contributor-workshop/image6.png"}}

### Impact

We had a total of 65 CLs submitted from the people who participated in the
workshop (within a week of the workshop). Of these, 44 were from contributors
who had never previously contributed to any of the repos in the Go project.
Half (22) of these contributions were already merged. Many of the others are
waiting on the codebase to thaw as we are in the middle of a freeze for the
upcoming 1.9 release. In addition to CLs, many contributed to the project in
the form of bug reports,
[gardening tasks](/wiki/Gardening), and other types
of contributions.

The most common type of contribution was an example function to be used in the
documentation. The [Go User survey](https://blog.golang.org/survey2016-results)
identified that our documentation was significantly lacking examples. In the
presentation, we asked users to find a package they loved and to add an example. In
the Go project, examples are written as code in Go files
(with specific naming) and the `go doc` tool displays them alongside the documentation.
This is a perfect first contribution as it's something that can be merged
during a freeze, it's of critical importance to our users, and it's an addition
that has a relatively narrow scope.

One of the examples added is that of creating a Stringer, one of the more
widely used interfaces in Go.
[CL 49270](/cl/49270/)

In addition to examples, many people contributed critical bug fixes including:

  - [CL 48988](/cl/48988/) fixing [issue #21029](/issue/21029)
  - [CL 49050](/cl/49050/) fixing [issue #20054](/issue/20054)
  - [CL 49031](/cl/49031/) fixing [issue #20166](/issue/20166)
  - [CL 49170](/cl/49170/) fixing [issue #20877](/issue/20877)

Some people even surprised us by arriving with a bug in mind that they wanted
to fix. Nikhita arrived ready to tackle
[issue #20786](/issue/20786)
and she did submit
[CL 48871](/cl/48871/),
after which she tweeted:

{{image "contributor-workshop/image19.png"}}

Not only were some great improvements made, but most importantly, we narrowed
the gap between the core Go team and the broader community members. Many people
on the Go team remarked that the community members were teaching them things
about the Go project. People in the community (in person, and on Twitter)
remarked that felt welcome to participate in the project.

{{image "contributor-workshop/image12.png"}}
{{image "contributor-workshop/image13.png"}}
{{image "contributor-workshop/image3.png"}}

### Future

The event was successful well beyond our expectations. Sameer Ajmani, Go team
manager said, "The contributor workshop was incredibly fun and educational–for
the Go team.  We cringed as users hit the rough edges in our process, and
celebrated when they got up on the dashboard. The cheer when the group score
hit 1000 was awesome."

We are looking into ways to make this workshop easier to run for future events
(like meetups and conferences). Our biggest challenge is providing enough
mentorship so that users feel supported. If you have any ideas or would like to
help with this process please [let me know](mailto:spf@golang.org).

I've asked a few participants of the event to share their experiences below:

## My Contribution Experience

by [Cassandra](https://twitter.com/cassandraoid)

When I heard about the go-contrib workshop I was very excited and then I was
extremely intimidated. I was encouraged by a member of the Go team to
participate, so I thought what the heck.

As I walked into the room (let's be real, I ran into the room because I was
running late) I was pleased to see the room was jam-packed. I looked around for
people in Gopher caps, which was the main indicator they were teachers. I sat
down at one of the 16 round tables that had two hats and three non-hats.
Brought up my screen and was ready to roll…

Jess Frazelle stood up and started the presentation and provided the group with
[a link](https://docs.google.com/presentation/d/1ap2fycBSgoo-jCswhK9lqgCIFroE1pYpsXC1ffYBCq4/edit#slide=id.p)
to make it easy to follow.

{{image "contributor-workshop/image16.png"}}

The murmurs grew from a deep undercurrent to a resounding melody of voices,
people were getting their computers set up with Go, they were skipping ahead to
make sure their GOPATH was set, and were… wait what's Gerrit?

Most of us had to get a little intro to Gerrit. I had no clue what it was, but
luckily there was a handy slide. Jess explained that it was an alternative to
GitHub with slightly more advanced code review tools. We then went through
GitHub vs Geritt terminology, so we had better understanding of the process.

{{image "contributor-workshop/image10.png"}}

Ok, now it was time to become a **freaking Go contributor**.

To make this more exciting than it already is, the Go team set up a game where
we could track as a group how many points we could rack up based on the Gerrit
score system.

{{image "contributor-workshop/image7.png"}}

Seeing your name pop up on the board and listening to everyone's excitement was
intoxicating. It also invoked a sense of teamwork that lead to a feeling of
inclusion and feeling like you were truly a part of the Go community.

{{image "contributor-workshop/image11.png"}}

In 6 steps a room of around 80 people were able to learn how to contribute to
go within an hour. That's a feat!

It wasn't nearly as difficult as I anticipated and it wasn't out of scope for a
total newbie. It fostered a sense of community in an active and tangible way as
well as a sense of inclusion in the illustrious process of Go contributions.

I'd personally like to thank the Go Team, the Gopher mentors in hats, and my
fellow participants for making it one of my most memorable moments at
GopherCon.

## My Contribution Experience

by [Matt](https://twitter.com/mbbroberg)

I've always found programming languages to be intimidating. It's the code that
enables the world to write code. Given the impact, surely smarter people than
me should be working on it... but that fear was something to overcome. So when
the opportunity to join a workshop to contribute to my new favorite programming
language came up, I was excited to see how I could help. A month
later, I'm now certain that anyone and everyone can (and should) contribute back to Go.

Here are my very verbose steps to go from 0 to 2 contributions to Go:

### The Setup

Given Go's use of Gerrit, I started by setting up my environment for it. [Jess Frazzelle's guide](https://docs.google.com/presentation/d/1ap2fycBSgoo-jCswhK9lqgCIFroE1pYpsXC1ffYBCq4/edit#slide=id.g1f953ef7df_0_9)
is a great place to start to not miss a step.

The real fun starts when you clone the Go repo. Ironically, you don't hack on
Go under `$GOPATH`, so I put it in my other workspace (which is `~/Develop`).

	cd $DEV # That's my source code folder outside of $GOPATH
	git clone --depth 1 https://go.googlesource.com/go

Then install the handy dandy helper tool, `go-contrib-init`:

	go get -u golang.org/x/tools/cmd/go-contrib-init

Now you can run `go-contrib-init` from the `go/` folder we cloned above and see
whether or not we're ready to contribute. But hold on if you're following along,
you're not ready just yet.

Next, install `codereview` so you can participate in a Gerrit code review:

	go get -u golang.org/x/review/git-codereview

This package includes `git change` and `git mail` which will replace your
normal workflow of `git commit` and `git push` respectively.

Okay, installations are out of the way. Now set up your [Gerrit account here](https://go-review.googlesource.com/settings/#Profile),
then [sign the CLA](https://go-review.googlesource.com/settings#Agreements) appropriate for
you (I signed a personal one for all Google projects, but choose the right option for you.
You can see all CLAs you've signed at [cla.developers.google.com/clas](https://cla.developers.google.com/clas)).

AND BAM. You're good (to go)! But where to contribute?

### Contributing

In the workshop, they sent us into the `scratch` repository, which is a safe place to
fool around in order to master the workflow:

	cd $(go env GOPATH)/src/golang.org/x
	git clone --depth 1 [[https://go.googlesource.com/scratch][go.googlesource.com/scratch]]

First stop is to `cd` in and run `go-contrib-init` to make sure you're ready to contribute:

	go-contrib-init
	All good. Happy hacking!

From there, I made a folder named after my GitHub account, did a `git add -u`
then took `git change` for a spin. It has a hash that keeps track of your work,
which is the one line you shouldn't touch. Other than that, it feels just like
`git commit`. Once I got the commit message matching the format of
`package: description` (description begins with a lowercase), I used
`git mail` to send it over to Gerrit.

Two good notes to take at this point: `git change` also works like `git commit --amend`, so
if you need to update your patch you can `add` then `change` and it will all
link to the same patch. Secondly, you can always review your patch from your
[personal Gerrit dashboard](https://go-review.googlesource.com/dashboard/).

After a few back and forths, I officially had a contribution to Go! And if Jaana
is right, it might be the first with emojis ✌️.

{{image "contributor-workshop/image15.png"}}
{{image "contributor-workshop/image23.png"}}

### Contributing, For Real

The scratch repo is fun and all, but there's a ton of ways to get into the
depths of Go's packages and give back. It's at this point where I cruised
around the many packages available to see what was available and interesting to
me. And by "cruised around" I mean attempted to find a list of packages, then
went to my source code to see what's around under the `go/src/` folder:

{{image "contributor-workshop/image22.png"}}

I decided to see what I can do in the `regexp` package, maybe out of love and
fear of regex. Here's where I switched to the
[website's view of the package](https://godoc.org/regexp) (it's good to know
that each standard package can be found at https://godoc.org/$PACKAGENAME). In
there I noticed that `QuoteMeta` was missing the same level of detailed examples
other functions have (and I could use the practice using Gerrit).

{{image "contributor-workshop/image1.png"}}

I started looking at `go/src/regexp` to try to find where to add examples and I
got lost pretty quickly. Lucky for me, [Francesc](https://twitter.com/francesc) was around that day. He walked
me through how all examples are actually in-line tests in a `example_test.go`
file. They follow the format of test cases followed by "Output" commented out
and then the answers to the tests. For example:

	func ExampleRegexp_FindString() {
		re := regexp.MustCompile("fo.?")
		fmt.Printf("%q\n", re.FindString("seafood"))
		fmt.Printf("%q\n", re.FindString("meat"))
		// Output:
		// "foo"
		// ""
	}

Kind of cool, right?? I followed Francesc's lead and added a function
`ExampleQuoteMeta` and added a few I thought would be helpful. From there it's
a `git change` and `git mail` to Gerrit!

I have to say that Steve Francia challenged me to "find something that isn't an
open issue and fix it," so I included some documentation changes for QuoteMeta
in my patch. It's going to be open for a bit longer given the additional scope,
but I think it's worth it on this one.

I can hear your question already: how did I verify it worked? Well it wasn't
easy to be honest. Running `go test example_test.go -run QuoteMeta -v` won't do
it since we're working outside of our $GOPATH. I struggled to figure it out
until [Kale Blakenship wrote this awesome post on testing in Go](https://medium.com/@vCabbage/go-testing-standard-library-changes-1e9cbed11339).
Bookmark this one for later.

You can see my completed [contribution here](https://go-review.googlesource.com/c/49130/). What I also hope you see is
how simple it is to get into the flow of contributing. If you're like me,
you'll be good at finding a small typo or missing example in the docs to start to
get used to the `git codereview` workflow. After that, you'll be ready to find
an open issue, ideally one [tagged for an upcoming release](https://github.com/golang/go/milestones), and give it a go. No matter
what you choose to do, definitely go forth and do it. The Go team proved to me
just how much they care about helping us all contribute back. I can't wait for
my next `git mail`.

## My Mentorship Experience

by [Dmitri](https://twitter.com/dmitshur)

I was looking forward to participating in the Contribution Workshop event as a
mentor. I had high expectations for the event, and thought it was a great idea
before it started.

I made my first contribution to Go on May 10th, 2014. I remember it was about
four months from the moment I wanted to contribute, until that day, when I
actually sent my first CL. It took that long to build up the courage and fully
commit to figuring out the process. I was an experienced software engineer at
the time. Despite that, the Go contribution process felt alien—being unlike all
other processes I was already familiar with—and therefore seemed intimidating.
It was well documented though, so I knew it would be just a matter of finding
the time, sitting down, and doing it. The "unknown" factor kept me from giving
it a shot.

After a few months passed, I thought "enough is enough," and decided to
dedicate an entire day of an upcoming weekend to figuring out the process. I
set aside all of Saturday for doing one thing: sending my first CL to Go. I
opened up [the Contribution Guide](/doc/contribute.html)
and started following all the steps, from the very top. Within an hour, I was
done. I had sent my first CL. I was both in awe and shock. In awe, because I
had finally sent a contribution to Go, and it was accepted! In shock, because,
why did I wait so long to finally do this? Following the steps in
[the Contribution Guide](/doc/contribute.html) was very
easy, and the entire process went completely smoothly. If only someone had told
me that I'd be done within an hour and nothing would go wrong, I would've done
it much sooner!

Which brings me to this event and why I thought it was such a good idea. For
anyone who ever wanted to contribute to Go, but felt daunted by the unfamiliar
and seemingly lengthy process (like I was during those four months), this was
their chance! Not only is it easy to commit to figuring it out by attending the
event, but also the Go team and helpful volunteer mentors would be there to
help you along the way.

Despite the already high expectations I had for the event, my expectations were
exceeded. For one, the Go team had prepared really well and invested a lot in
making the event that much more enjoyable for everyone. There was a very fun
presentation that went over all the contributing steps quickly. There was a
dashboard made for the event, where everyone's successfully completed steps
were rewarded with points towards a global score. That made it into a very
collaborative and social event! Finally, and most importantly, they were Go
team members like Brad Fitzpatrick behind the scenes, helping review CLs
promptly! That meant the CLs that were submitted received reviews quickly, with
actionable next steps, so everyone could move forward and learn more.

I originally anticipated the event to be somewhat dull, in that the
contribution steps are extremely simple to follow. However, I found that wasn't
always the case, and I was able to use my expertise in Go to help out people
who got stuck in various unexpected places. It turns out, the real world is
filled with edge cases. For instance, someone had two git emails, one personal
and another for work. There was a delay with signing the CLA for the work
email, so they tried to use their personal email instead. That meant each
commit had to be amended to use the right email, something the tools didn't
take into account. (Luckily, there is a troubleshooting section in the
contribution guide covering this exact issue!) There were other subtle mistakes
or environment misconfiguration that some people ran into, because having more
than one Go installation was a bit unusual. Sometimes, the GOROOT environment
variable had to be explicitly set, temporarily, to get godoc to show changes in
the right standard library (I was tongue-in-cheek looking over my shoulder to
check for Dave Cheney as I uttered those words).

Overall, I oversaw a few new gophers make their first Go contributions. They
sent the CLs, responded to review feedback, made edits, iterated until everyone
was happy, and eventually saw their first Go contributions get merged to
master! It was very rewarding to see the happiness on their faces, because the
joy of making one's first contribution is something I can relate to myself. It
was also great to be able to help them out, and explain tricky situations that
they sometimes found themselves. From what I can tell, many happy gophers
walked away from the event, myself included!

## Photos from the event

{{image "contributor-workshop/image2.jpg"}}
{{image "contributor-workshop/image4.jpg"}}
{{image "contributor-workshop/image5.jpg"}}
{{image "contributor-workshop/image8.jpg"}}
{{image "contributor-workshop/image14.jpg"}}
{{image "contributor-workshop/image18.jpg"}}
{{image "contributor-workshop/image20.jpg"}}
{{image "contributor-workshop/image21.jpg"}}

Photos by Sameer Ajmani & Steve Francia
