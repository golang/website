---
title: "How the Firebase Hosting Team Scaled With Go"
company: Firebase
logoSrc: firebase.svg
logoSrcDark: firebase.svg
heroImgSrc: go_firebase_case_study.png
series: Case Studies
quote: |
  Firebase is Google’s mobile platform that helps you quickly develop high-quality
  apps and grow your business.

  The Firebase Hosting team shared their journey with Go, including their
  backend migration from Node.js, the ease of onboarding new Go developers, and
  how Go has helped them scale.
---

The Firebase Hosting team provides static web hosting services for Google Cloud
customers. They provide a static web host that sits behind a global content
delivery network, and offer users tools that are easy to use. The team also
develops features that range from uploading site files to registering domains to
tracking usage.

Before joining Google, Firebase Hosting’s tech stack was written in Node.js. The
team started to use Go when they needed to interoperate with several other
Google services. They decided to use Go to help them scale easily and
efficiently, knowing that “concurrency would continue to be a big need.” They
“were confident Go would be more performant,” and “liked that Go is more terse”
than other languages they were considering, said Michael Bleigh, a software
engineer on the team.

Starting with one small service written in Go, the team migrated their entire
backend in a series of moves. The team progressively identified large features
they wanted to implement and, in the process, rewrote them in Go and moved to
Google Cloud and Google’s internal cluster management system. **Now the Firebase
Hosting team has replaced 100% of backend Node.js code with Go.**

The team’s experience writing in Go began with one engineer. “Through
peer-to-peer learning and Go being generally easy to get started with, everyone
on the team now has Go dev experience,” said Bleigh. They’ve found that while a
majority of people who are new to the team haven’t had any experience with Go,
“most of them are productive within a couple weeks.”

"Using Go, it's easy to see how the code is organized and what the code does,"
said Bleigh, speaking for the team. “Go is generally very readable and
understandable. The language’s error handling, receivers, and interfaces are all
easy to understand due to the idioms in the language.”

Concurrency continues to be a focus for the team as they scale. Robert Rossney,
a software engineer, shared that “Go makes it very easy to put all of the hard
concurrency stuff in one place, and everywhere else it's abstracted.” Rossney
also spoke to the benefits of using a language built with concurrency in mind,
saying that “there are also a lot of ways to do concurrency in Go. We’ve had to
learn when each route is best, how to determine when a problem is a concurrency
problem, how to debug–but that comes out of the fact that you actually can write
these patterns in Go code.”

{{backgroundquote `
  author: Robert Rossney
  title: Software Engineer
  quote: |
    Generally speaking, there’s not a time on the team where we’re feeling
    frustrated with Go, it just kind of gets out of the way and lets you do work.
`}}

Hundreds of thousands of customers host their websites with Firebase Hosting,
which means Go code is used to serve billions of requests per day. “Our customer
base and traffic have doubled multiple times since migrating to Go without ever
requiring fine-tuned optimizations” shared Bleigh.  With Go, the team has seen
performance improvements both in the software and on the team, with excellent
productivity gains. “Generally speaking,” Rossney mentioned, “...there’s not a
time on the team where we’re feeling frustrated with Go, it just kind of gets
out of the way and lets you do work.”

In addition to the Firebase Hosting team, engineering teams across Google have
adopted Go in their development process. Read about how the [Core Data
Solutions](/solutions/google/coredata/) and [Chrome](/solutions/google/chrome/)
teams use Go to build fast, reliable, and efficient software at scale.
