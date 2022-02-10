---
title: "How Google's Core Data Solutions Team Uses Go"
company: Core Data
logoSrc: google.svg
logoSrcDark: google.svg
heroImgSrc: go_core_data_case_study.png
series: Case Studies
quote: |
  Google is a technology company whose mission is to organize the world’s
  information and make it universally accessible and useful.

  In this case study, Google’s Core Data Solutions team shares their journey
  with Go, including their decision to rewrite web indexing services in Go,
  taking advantage of Go’s built-in concurrency, and observing how Go helps to
  improve the development process.
authors:
  - Prasanna Meda, Software Engineer, Core Data Solutions
---

Google's mission is “to organize the world's information and make it universally
accessible and useful.”  One of the teams responsible for organizing that
information is Google’s Core Data Solutions team. The team, among other things,
maintains services to index web pages across the globe. These web indexing
services help support products like Google Search by keeping search results
updated and comprehensive, and they’re written in Go.

In 2015, to keep up with Google’s scale, our team needed to rewrite our indexing
stack from a single monolithic binary written in C++ to multiple components in a
microservices architecture. We decided to rewrite many indexing services in Go,
which we now use to power the majority of our architecture.

{{backgroundquote `
  author: Minjae Hwang
  title: Software Engineer
  quote: |
    Go’s built-in concurrency is a natural fit because engineers on the team are
    encouraged to use concurrency and parallel algorithms.
`}}

When choosing a language, our team found that several of Go’s features made it
particularly suitable. For instance, Go’s built-in concurrency is a natural fit
because engineers on the team are encouraged to use concurrency and parallel
algorithms. Engineers have also found that “Go code is more natural,” allowing
them to spend their time focusing on business logic and analysis rather than on
managing memory and optimizing performance.

Writing code is much simpler when writing in Go, as it helps lessen cognitive
burden during the development process. For example, when working with C++,
sophisticated IDEs might, “show that the source code has no compile error when
there actually is one” whereas “in Go, [the code] will always compile when [the
IDE] says the code has no compile error,” said MinJae Hwang, a software engineer
on the Core Data Solutions team. Reducing small friction points along the
development process, such as shortening the cycle of fixing compile errors,
helped our team ship faster during the original rewrite, and has helped keep our
maintenance costs low.

“When I’m in C++ and I want to use more packages, I have to write pieces such as
headers. When I'm writing in Go, **built-in tools allow me to use packages more
easily. My development velocity is much faster,**” Hwang also shared.

With simple language syntax and support of Go tools, several members of our team
find it much easier to write in Go code. We’ve also found that Go does a really
good job of static type checking and that certain Go fundamentals, such as the
godoc command, have helped the team build a more disciplined culture around
writing documentation.

{{backgroundquote `
  author: Prasanna Meda
  title: Software Engineer
  quote: |
    ...Google’s web indexing was re-architected within a year. More impressively,
    most developers on the team were rewriting in Go while also learning it.
`}}

Working on a product used so heavily around the world is no small task and our
team’s decision to use Go wasn’t a simple one, but doing so helped us move
faster. As a result, Google’s web indexing was re-architected within a year.
More impressively, most developers on the team were rewriting in Go while also
learning it.

In addition to the Core Data Solutions team, engineering teams across Google
have adopted Go in their development process. Read about how the
[Chrome](/solutions/google/chrome/) and [Firebase
Hosting](/solutions/google/firebase/) teams use Go to build fast, reliable, and
efficient software at scale.
