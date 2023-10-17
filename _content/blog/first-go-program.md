---
title: The first Go program
date: 2013-07-18
by:
- Andrew Gerrand
tags:
- history
summary: Rob Pike dug up the first Go program ever written.
---


Brad Fitzpatrick and I (Andrew Gerrand) recently started restructuring
[godoc](/cmd/godoc/), and it occurred to me that it is one
of the oldest Go programs.
Robert Griesemer started writing it back in early 2009,
and we're still using it today.

When I [tweeted](https://twitter.com/enneff/status/357403054632484865) about
this, Dave Cheney replied with an [interesting question](https://twitter.com/davecheney/status/357406479415914497):
what is the oldest Go program? Rob Pike dug into his mail and found it
in an old message to Robert and Ken Thompson.

What follows is the first Go program. It was written by Rob in February 2008,
when the team was just Rob, Robert, and Ken. They had a solid feature list
(mentioned in [this blog post](https://commandcenter.blogspot.com.au/2012/06/less-is-exponentially-more.html))
and a rough language specification. Ken had just finished the first working version of
a Go compiler (it didn't produce native code, but rather transliterated Go code
to C for fast prototyping) and it was time to try writing a program with it.

Rob sent mail to the "Go team":

	From: Rob 'Commander' Pike
	Date: Wed, Feb 6, 2008 at 3:42 PM
	To: Ken Thompson, Robert Griesemer
	Subject: slist

	it works now.

	roro=% a.out
	(defn foo (add 12 34))
	return: icounter = 4440
	roro=%

	here's the code.
	some ugly hackery to get around the lack of strings.

(The `icounter` line in the program output is the number of executed
statements, printed for debugging.)

{{code "first-go-program/slist.go"}}

The program parses and prints an
[S-expression](https://en.wikipedia.org/wiki/S-expression).
It takes no user input and has no imports, relying only on the built-in
`print` facility for output.
It was written literally the first day there was a
[working but rudimentary compiler](/change/8b8615138da3).
Much of the language wasn't implemented and some of it wasn't even specified.

Still, the basic flavor of the language today is recognizable in this program.
Type and variable declarations, control flow, and package statements haven't
changed much.

But there are many differences and absences.
Most significant are the lack of concurrency and interfaces—both
considered essential since day 1 but not yet designed.

A `func` was a `function`, and its signature specified return values
_before_ arguments, separating them with {{raw "`<-`"}}, which we now use as the channel
send/receive operator. For example, the `WhiteSpace` function takes the integer
`c` and returns a boolean.

{{raw `
	function WhiteSpace(bool <- c int)
`}}

This arrow was a stop-gap measure until a better syntax arose for declaring
multiple return values.

Methods were distinct from functions and had their own keyword.

{{raw `
	method (this *Slist) Car(*Slist <-) {
		return this.list.car;
	}
`}}

And methods were pre-declared in the struct definition, although that changed soon.

{{raw `
	type Slist struct {
		...
		Car method(*Slist <-);
	}
`}}

There were no strings, although they were in the spec.
To work around this, Rob had to build the input string as an `uint8` array with
a clumsy construction. (Arrays were rudimentary and slices hadn't been designed
yet, let alone implemented, although there was the unimplemented concept of an
"open array".)

	input[i] = '('; i = i + 1;
	input[i] = 'd'; i = i + 1;
	input[i] = 'e'; i = i + 1;
	input[i] = 'f'; i = i + 1;
	input[i] = 'n'; i = i + 1;
	input[i] = ' '; i = i + 1;
	...

Both `panic` and `print` were built-in keywords, not pre-declared functions.

	print "parse error: expected ", c, "\n";
	panic "parse";

And there are many other little differences; see if you can identify some others.

Less than two years after this program was written, Go was released as an
open source project. Looking back, it is striking how much the language has
grown and matured. (The last thing to change between this proto-Go and the Go
we know today was the elimination of semicolons.)

But even more striking is how much we have learned about _writing_ Go code.
For instance, Rob called his method receivers `this`, but now we use shorter
context-specific names. There are hundreds of more significant examples
and to this day we're still discovering better ways to write Go code.
(Check out the [glog package](https://github.com/golang/glog)'s clever trick for
[handling verbosity levels](https://github.com/golang/glog/blob/c6f9652c7179652e2fd8ed7002330db089f4c9db/glog.go#L893).)

I wonder what we'll learn tomorrow.
