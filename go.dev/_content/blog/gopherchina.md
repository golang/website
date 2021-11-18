---
title: GopherChina Trip Report
date: 2015-07-01
by:
- Robert Griesemer
tags:
- community
- china
summary: Reporting from GopherChina 2015, the first Go conference in China.
---


We have known for some time that Go is more popular in China than in any other
country.
According to Google Trends, most [searches for the term “golang”](https://www.google.com/trends/explore#q=golang) come from The People’s Republic than anywhere else.
[Others](http://herman.asia/why-is-go-popular-in-china) have speculated on
the same observation, yet so far we have had
[sparse concrete information](https://news.ycombinator.com/item?id=8872400)
about the phenomenon.

The first Go conference in China, [GopherChina](http://gopherchina.org/),
seemed like an excellent opportunity to explore the situation by putting some
Western Gopher feet on Chinese ground. An actual invitation made it real and I
decided to accept and give a presentation about gofmt’s impact on software
development.

{{image "gopherchina/image04.jpg"}}

_Hello, Shanghai!_

The conference took place over an April weekend in Shanghai, in the
[Puruan Building](https://www.google.com/maps/place/Puruan+Bldg,+Pudong,+Shanghai,+China)
of the Shanghai Pudong Software Park, easily reachable by subway within an hour
or less from Shanghai’s more central parts.
Modelled after [GopherCon](http://www.gophercon.com), the conference was
single-track, with all talks presented in a conference room that fit about 400
attendees.
It was organized by volunteers, lead by [Asta Xie](https://github.com/astaxie),
and with robust sponsorship from major industry names. According to the
organizers, many more people were hoping to attend than could be accommodated
due to space constraints.

{{image "gopherchina/image01.jpg"}}

_The welcoming committee with Asta Xie (2nd from left), the primary organizer._

Each attendee received a bag filled with the obligatory GopherChina t-shirt,
various sponsor-related informational brochures, stickers, and the occasional
stuffed “something” (no fluffy Gophers, though). At least one 3rd party vendor
was advertising technical books, including several original (not translated
from English) Go books.

{{image "gopherchina/image05.jpg"}}

_Go books!_

On first impression, the average attendee seemed pretty young, which made for
an enthusiastic crowd, and the event appeared well run.

With the exception of my talk, all presentations were given in Mandarin and
thus were incomprehensible to me. Asta Xie, the primary organizer, assisted
with a few simultaneous translations whispered into my ear, and the occasional
English slide provided additional clues: “69GB” stands out even without any
Mandarin knowledge (more on that below). Consequently, I ended up listening to
a handful of presentations only, and instead spent much of my time talking with
attendees outside the main conference room. Yet judging from the slides, the
quality of most presentations seemed high, comparable with our experience at
GopherCon in Denver last year. Each talk got a one hour time slot which allowed
for plenty of technical detail, and many (dozens) of questions from an
enthusiastic audience.

As expected, many of the presentations were about web services, backends for
mobile applications, and so on. Some of the systems appear to be huge by any
measure.
For instance, a talk by [Yang Zhou](http://gopherchina.org/user/zhouyang)
described a large-scale internal messaging system, used by
[Qihoo 360](http://www.360.cn/), a major Chinese software firm, all written
in Go. The presentation discussed how his team managed to reduce an original
heap size of 69GB (!) and the resulting long GC pauses of 3-6s to more
manageable numbers, and how they run millions of goroutines per machine, on a
fleet of thousands of machines. A future guest blog post is planned describing
this system in more detail.

{{image "gopherchina/image03.jpg"}}

_Packed conference room on Saturday._

In another presentation, [Feng Guo](http://gopherchina.org/user/guofeng) from
[DaoCloud](https://www.daocloud.io/) talked about how they use Go in their
company for what they call the “continuous delivery” of applications. DaoCloud
takes care of automatically moving software hosted on GitHub (and Chinese
equivalents) to the cloud. A software developer simply pushes a new version on
GitHub and DaoCloud takes care of the rest: running tests,
[Dockerizing](https://www.docker.com/) it, and shipping it using your
preferred cloud service provider.

Several speakers were from well-recognized major software firms (I showed the
conference program to non-technical people and they easily recognized several
of the firm’s names). Much more so than in the US, it seems Go is not just
hugely popular with newcomers and startups, but has very much found its way
into larger organizations and is employed at a scale that we are only starting
to see elsewhere.

Not being an expert in web services myself, in my presentation I veered off the
general conference theme a bit by talking about
[gofmt](/cmd/gofmt/) and how its widespread use has started
to shape expectations not just for Go but other languages as well.
I presented in English but had my slides translated to Mandarin beforehand. Due
to the significant language barrier I wasn’t expecting too many questions on my
talk itself.
Instead I decided the keep it short and leave plenty of time for general
questions on Go, which the audience appreciated.

{{image "gopherchina/image06.jpg"}}

_No social event in China is complete without fantastic food._

A couple of days after the conference I visited the 4-year-old startup company
[Qiniu](http://www.qiniu.com/) (“Seven Bulls”), at the invitation of its
[CEO](http://gopherchina.org/user/xushiwei) Wei Hsu, facilitated and
translated with the help of Asta Xie. Qiniu is a cloud-based storage provider
for mobile applications; Wei Hsu presented at the conference and also happens
to be the author of one of the first Chinese books on Go (the leftmost one in
the picture above).

{{image "gopherchina/image02.jpg"}}
{{image "gopherchina/image00.jpg"}}

_Qiniu lobby, engineering._

Qiniu is an extremely successful all-Go shop, with about 160 employees, serving
over 150,000 companies and developers, storing over 50 Billion files, and
growing by over 500 Million files per day. When asked about the reasons for
Go’s success in China, Wei Hsu is quick to answer: PHP is extremely popular in
China, but relatively slow and not well-suited for large systems. Like in the
US, universities teach C++ and Java as primary languages, but for many
applications C++ is too complex a tool and Java too bulky. In his opinion, Go
now plays the role that traditionally belonged to PHP, but Go runs much faster,
is type safe, and scales more easily. He loves the fact that Go is simple and
applications are easy to deploy. He thought the language to be “perfect” for
them and his primary request was for a recommended or even standardized package
to easily access database systems. He did mention that they had GC problems in
the past but were able to work around them. Hopefully our upcoming 1.5 release
will address this. For Qiniu, Go appeared just at the right time and the right
(open source) place.

According to Asta Xie, Qiniu is just one of many Go shops in the PRC. Large
companies such as Alibaba, Baidu, Tencent, and Weibo, are now all using Go in
one form or another. He pointed out that while Shanghai and neighboring cities
like [Suzhou](https://www.google.com/maps/place/Suzhou,+Jiangsu,+China) are
high-tech centres, even more software developers are found in the Beijing area.
For 2016,  Asta hopes to organize a larger (1000, perhaps 1500 people)
successor conference in Beijing.

It appears that we have found the Go users in China: They are everywhere!

_Some of the GopherChina materials, including videos, are now available alongside Go coursework on a_ [_3rd party site_](http://www.imooc.com/view/407).
