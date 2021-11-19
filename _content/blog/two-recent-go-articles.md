---
title: Two recent Go articles
date: 2013-03-06
by:
- Andrew Gerrand
tags:
- google
- talk
- ethos
summary: "Two Go articles: “Go at Google: Language Design in the Service of Software Engineering” and “Getting Started with Go, App Engine and Google+ API”"
---

## Introduction

In today's blog post I'd like to highlight a couple of recent articles about Go.

## Go at Google

In October last year, Rob Pike presented a keynote at the ACM [SPLASH](http://splashcon.org/2012/) conference in Tucson.
The talk, titled [Go at Google](/talks/2012/splash.slide),
was a comprehensive discussion of the motivations behind Go.
Rob later expanded on his talk to produce an essay titled [Go at Google: Language Design in the Service of Software Engineering](http://go.dev/talks/2012/splash.article).
Here is the abstract:

	The Go programming language was conceived in late 2007 as an
	answer to some of the problems we were seeing developing
	software infrastructure at Google. The computing landscape
	today is almost unrelated to the environment in which the
	languages being used, mostly C++, Java, and Python, had been
	created. The problems introduced by multicore processors,
	networked systems, massive computation clusters, and the web
	programming model were being worked around rather than
	addressed head-on. Moreover, the scale has changed: today's
	server programs comprise tens of millions of lines of code,
	are worked on by hundreds or even thousands of programmers,
	and are updated literally every day.  To make matters worse,
	build times, even on large compilation clusters, have
	stretched to many minutes, even hours.

	Go was designed and developed to make working in this
	environment more productive. Besides its better-known
	aspects such as built-in concurrency and garbage collection,
	Go's design considerations include rigorous dependency
	management, the adaptability of software architecture as
	systems grow, and robustness across the boundaries between
	components.

This article explains how these issues were addressed while building an efficient,
compiled programming language that feels lightweight and pleasant.
Examples and explanations will be taken from the real-world problems faced at Google.

If you have wondered about the design decisions behind Go,
you may find your questions answered by [the essay](/talks/2012/splash.article).
It is recommended reading for both new and experienced Go programmers.

## Go at the Google Developers Academy

At Google I/O 2012 the Google Developers team [launched](http://googledevelopers.blogspot.com.au/2012/06/google-launches-new-developer-education.html) the [Google Developers Academy](https://developers.google.com/academy/),
a program that provides training materials on Google technologies.
Go is one of those technologies and we're pleased to announce the first
GDA article featuring Go front and center:

[Getting Started with Go, App Engine and Google+ API](https://developers.google.com/appengine/training/go-plus-appengine/) is
an introduction to writing web applications in Go.
It demonstrates how to build and deploy App Engine applications and make
calls to the Google+ API using the Google APIs Go Client.
This is a great entry point for Go programmers eager to get started with
Google's developer ecosystem.
