---
title: "MercadoLibre Grows with Go"
company: MercadoLibre
logoSrc: mercadolibre_light.svg
logoSrcDark: mercadolibre_dark.svg
heroImgSrc: go_mercadolibre_case_study_logo.png
carouselImgSrc: go_mercadolibre_case_study.png
date: 2019-11-10T16:26:31-04:00
series: Case Studies
quote: Go provides clean, efficient code that readily scales as MercadoLibre’s online commerce grows, and increases developer productivity by allowing their engineers to serve their ever-increasing audience while writing less code.
---

{{pullquote `
  author: Eric Kohan
  title: Software Engineering Manager
  company: MercadoLibre
  quote: |
    I think that **the tour of Go is by far the best introduction to a language that I’ve seen**, It’s really simple and it gives you a fair overview of probably 80 percent of the language. When we want to get developers to learn Go, and to get to production fast, we tell them to start with the tour of Go.
`}}

## Go helps integrated ecosystem attract developers and scale eCommerce

MercadoLibre, Inc. hosts the largest online commerce ecosystem in Latin America and is present in 18 countries. Founded
in 1999 and headquartered in Argentina, the company has turned to Go to help it scale and modernize its ecosystem. Go
provides clean, efficient code that readily scales as MercadoLibre’s online commerce grows, and increases developer
productivity by allowing their engineers to serve their ever-increasing audience while writing less code.

### MercadoLibre taps Go for scale

Back in 2015, there was a growing sense within MercadoLibre that their existing API framework, on Groovy and Grails, was
reaching its limits and the company needed a different platform to continue scaling. MercadoLibre’s platform was (and
continues) to expand exponentially, which created a lot of extra work for its developers: Both Groovy and Grails require
a lot of decisions from developers and Groovy is a dynamic programming language. This was not a good combination for
quickly scaling growth, as MercadoLibre needed very experienced developers in this very resource intensive environment
to develop and tune to achieve desired performance. Test execution times were slow, and build and deploy times were
slow. Thus, the need for code efficiency and scalability became as important as the need for speed in code development.


### Go improves system efficiency

As one example of Go’s contributions to network efficiency, the core API team builds and maintains the largest APIs at
the center of the company’s microservices solutions. This team creates user APIs, which in turn are used by the
MercadoLibre Marketplace, by the MercadoPago FinTech platform, by MercadoLibre’s shipping and logistics solutions, and
other hosted solutions. With the high service levels demanded by these solutions—the average user API has between eight
and ten million requests per minute—the team employs Go to serve them at less than ten milliseconds per request.

The API team also deploys Docker containers—a software-as-a-service (SaaS) product, also written in Go—to virtualize
their development and readily deploy their microservices via the Docker Engine. This system supports larger,
mission-critical APIs that handle **more than 20 million requests per minute in Go.**

One API made important use of Go’s concurrency primitives to efficiently multiplex IDs from several services. The team
was able to accomplish this with just a few lines of Go code, and the success of this API convinced the core API team to
migrate more and more microservices to Go. The end result for MercadoLibre has been improved cost-efficiencies and
system response times.

### Go for scalability

Historically, much of the company’s stack was based on Grails and Groovy backed by relational  databases. However this
big framework with multiple layers was soon found encountering scalability issues.

Converting that legacy architecture to Go as a new, very thin framework for building APIs streamlined those intermediate
layers and yielded great performance benefits. For example, one large Go service is now able to **run 70,000 requests
per machine with just 20 MB of RAM.**

{{backgroundquote `
  author: Eric Kohan
  title: Software Engineering Manager
  company: MercadoLibre
  quote: |
    Go was just marvelous for us. It’s very powerful
    and very easy to learn, and with backend infrastructure, has been great for us in terms of scalability.
`}}

Using **Go allowed MercadoLibre to cut the number of servers** they use for this service to one-eighth the original
number (from 32 servers down to four), plus each server can operate with less power (originally four CPU cores, now down
to two CPU cores). With Go, the company **obviated 88 percent of their servers and cut CPU on the remaining ones in
half**—producing a tremendous cost-savings.

Sitting between developers and the cloud providers, MercadoLibre uses a platform called Fury—a platform-as-a-service
tool for building, deploying, monitoring, and managing services in a cloud-agnostic way. As a result, any team that
wants to create a new service in Go has access to proven templates for a variety of service types, and can quickly spin
up a repository in GitHub with starter code, a Docker image for the service, and a deployment pipeline. The end result
is a system that allows engineers to focus on building innovative services while avoiding the tedious stages of setting
up a new project—all while effectively standardizing the build and deployment pipelines.

Today, **roughly half of Mercadolibre's traffic is handled by Go applications.**


### MercadoLibre uses Go for developers

The programming _lingua francas_ for MercadoLibre’s infrastructure are currently Go and Java. Every app, every program,
every microservice is hosted on its own GitHub repository, plus the company uses an additional GitHub repository of
toolkits to solve new problems and allow clients to interact with its services.

These extensive and well-curated Go and Java toolkits allow programmers to develop new apps quickly and with great
support. Plus, in a community of more than 2,800 developers, MercadoLibre has multiple internal groups available for
chat and guidance on deploying Go, whether across different development centers or different countries. The company also
fosters internal working groups to provide training sessions for new MercadoLibre Go developers, and hosts Go meetups
for external developers to help build a broader community of Latin American Go developers.


### Go as a recruiting tool

MercadoLibre’s Go advocacy has also become a strong recruiting tool for the company. MercadoLibre was among the first
companies using Go in Argentina, and is perhaps the largest in Latin America using the language so widely in production.
Headquartered in Buenos Aires, with many start-ups and emerging technology companies nearby, MercadoLibre's adoption of
Go has shaped the market for developers across the Pampas.

{{backgroundquote `
  author: Eric Kohan
  title: Software Engineering Manager
  company: MercadoLibre
  quote: |
    We really see eye-to-eye with the larger philosophy of the language. We love Go's simplicity, and we find that having its very explicit error handling has been a gain for developers because it results in safer, more stable code in production.
`}}

Buenos Aires is today a very competitive market for programmers, offering computer programmers many employment options,
and the high demand for technology in the region drives great salaries, great benefits, and the ability to be selective
when choosing an employer. As such, MercadoLibre—like all employers of engineers and programmers in the region—strives
to provide an exciting workplace and strong career path. Go has proven to be a key differentiator for MercadoLibre: the
company organizes Go workshops for external developers so they can come and learn Go, and when they enjoy what they are
doing and the people they talk to, they quickly recognize MercadoLibre as an enticing place to work.

### Go enabling developers

MercadoLibre employs Go for its simplicity with systems at scale, but that simplicity is also why the company's
developers love Go.

The company also uses web pages like[ Go by Example](https://gobyexample.com/) and[ Effective
Go](/doc/effective_go.html) to educate new programmers, and shares representative internal APIs
written in Go to speed understanding and proficiency. MercadoLibre developers get the resources they need to embrace the
language, then leverage their own skills and enthusiasm to start programming.

{{backgroundquote `
  author: Federico Martin Roasio
  title: Technical Project Lead
  company: MercadoLibre
  quote: |
    Go has been great for writing business logic, and we are the team that writes those APIs.
`}}

MercadoLibre leverages Go’s expressive and clean syntax to make it easier for developers to write programs that run
efficiently on modern cloud platforms. And while speed in development yields cost efficiency for the company, developers
individually benefit from the swift learning curve Go delivers. Not only are MercadoLibre's experienced engineers able
to build highly critical applications very quickly with Go, but even entry-level engineers have been able to write
services that, in other languages, MercadoLibre would only trust to more senior developers. For example, a key set of
user APIs—handling almost ten million requests per minute—were developed by entry-level software engineers, many of whom
only knew about programming from recent courses at university. Similarly, MercadoLibre has seen developers already
proficient with other programming languages (such as Java or .NET or Ruby) learn Go fast enough start writing production
services in just a few weeks.

With Go, MercadoLibre’s **build times are three times (3x) faster** and their **test suite runs an amazing 24 times
faster**. This means the company’s developers can make a change, then build and test that change much faster than they
could before.

And dropping MercadoLibre’s test suite runtimes from 90-seconds to **just 3-seconds with Go** was a huge boon for its
developers—allowing them to keep focus (and context) while the much faster tests complete.

Leveraging this success, MercadoLibre is committed not only to ongoing education for its programmers, but ongoing Go
education. The company sends key engineering leaders to GopherCon and other Go events each year, MercadoLibre’s
infrastructure and security teams encourage all the development teams to keep Go versions up to date, and the company
has a team developing a _Go-meli-toolkit_: A complete Go library to interface all the services provided by Fury.

### Getting your enterprise started with Go

Just as MercadoLibre started with a proof-of-concept project to implement Go, dozens of other large enterprises are
adopting Go as well.

There are over one million developers using Go worldwide—spanning banking and commerce, gaming and media, technology, and other industries, at enterprises as diverse as [American Express](/solutions/americanexpress), [PayPal](/solutions/paypal), Capital One, Dropbox, IBM, Monzo, New York Times, Salesforce, Square, Target, Twitch, Uber, and of course Google.

To learn more about how Go can help your enterprise build reliable, scalable software as it does at MercadoLibre, visit [go.dev](/) today.
