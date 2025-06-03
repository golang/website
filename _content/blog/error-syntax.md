---
title: "[ On | No ] syntactic support for error handling"
date: 2025-06-03
by:
- Robert Griesemer
tags:
- error
- syntax
- technical
- proposal
summary: Go team plans around error handling support
---

One of the oldest and most persistent complaints about Go concerns the verbosity of error handling.
We are all intimately (some may say painfully) familiar with this code pattern:

```Go
x, err := call()
if err != nil {
        // handle err
}
```

The test `if err != nil` can be so pervasive that it drowns out the rest of the code.
This typically happens in programs that do a lot of API calls, and where handling errors
is rudimentary and they are simply returned.
Some programs end up with code that looks like this:

```Go
func printSum(a, b string) error {
	x, err := strconv.Atoi(a)
	if err != nil {
		return err
	}
	y, err := strconv.Atoi(b)
	if err != nil {
		return err
	}
	fmt.Println("result:", x + y)
	return nil
}
```

Of the ten lines of code in this function body, only four (the calls and the last two lines) appear to do real work.
The remaining six lines come across as noise.
The verbosity is real, and so it's no wonder that complaints about error handling have topped
our annual user surveys for years.
(For a while, the lack of generics surpassed complaints about error handling, but now that
Go supports generics, error handling is back on top.)

The Go team takes community feedback seriously, and so for many years now we have tried to
come up with a solution for this problem, together with input from the Go community.

The first explicit attempt by the Go team dates back to 2018, when Russ Cox
[formally described the problem](https://go.googlesource.com/proposal/+/master/design/go2draft-error-handling-overview.md)
as part of what we called the Go 2 effort at that time.
He outlined a possible solution based on a
[draft design](https://go.googlesource.com/proposal/+/master/design/go2draft-error-handling.md)
by Marcel van Lohuizen.
The design was based on a `check` and `handle` mechanism and was fairly comprehensive.
The draft includes a detailed analysis of alternative solutions, including comparisons with
approaches taken by other languages.
If you're wondering if your particular error handling idea was previously considered,
read this document!

```Go
// printSum implementation using the proposed check/handle mechanism.
func printSum(a, b string) error {
	handle err { return err }
	x := check strconv.Atoi(a)
	y := check strconv.Atoi(b)
	fmt.Println("result:", x + y)
	return nil
}
```

The `check` and `handle` approach was deemed too complicated and almost a year later, in 2019,
we followed up with the much simplified and by now
[infamous](/issue/32437#issuecomment-2278932700)
[`try` proposal](https://go.googlesource.com/proposal/+/master/design/32437-try-builtin.md).
It was based on the ideas of `check` and `handle`, but the `check` pseudo-keyword became
the `try` built-in function and the `handle` part was omitted.
To explore the impact of the `try` built-in, we wrote a simple tool
([tryhard](https://github.com/griesemer/tryhard))
that rewrites existing error handling code using `try`.
The proposal was argued over intensively, approaching 900 comments on the [GitHub issue](/issue/32437).

```Go
// printSum implementation using the proposed try mechanism.
func printSum(a, b string) error {
	// use a defer statement to augment errors before returning
	x := try(strconv.Atoi(a))
	y := try(strconv.Atoi(b))
	fmt.Println("result:", x + y)
	return nil
}
```

However, `try` affected control flow by returning from the enclosing function in case of an error,
and did so from potentially deeply nested expressions, thus hiding this control flow from view.
This made the proposal unpalatable to many, and despite significant investment
into this proposal we decided to abandon this effort too.
In retrospect it might have been better to introduce a new keyword,
something that we could do now since we have fine-grained control over the language version
via `go.mod` files and file-specific directives.
Restricting the use of `try` to assignments and statements might have alleviated some
of the other concerns. A [recent proposal](/issue/73376) by Jimmy Frasche, which essentially
goes back to the original `check` and `handle` design and addresses some of that design's
shortcomings, pursues that direction.

The repercussions of the `try` proposal led to much soul searching including a series of blog
posts by Russ Cox: ["Thinking about the Go Proposal Process"](https://research.swtch.com/proposals-intro).
One conclusion was that we likely diminished our chances for a better outcome by presenting an almost
fully baked proposal with little space for community feedback and a "threatening" implementation
timeline. Per ["Go Proposal Process: Large Changes"](https://research.swtch.com/proposals-large):
"in retrospect, `try` was a large enough change that the new design we published [...] should have
been a second draft design, not a proposal with an implementation timeline".
But irrespective of a possible process and communication failure in this case, the user sentiment towards
the proposal was very strongly not in favor.

We didn't have a better solution at that time and didn't pursue syntax changes for error handling for several years.
Plenty of people in the community were inspired, though, and we received a steady trickle
of error handling proposals, many very similar to each other, some interesting, some incomprehensible,
and some infeasible.
To keep track of the expanding landscape, another year later, Ian Lance Taylor created an
[umbrella issue](/issue/40432)
which summarizes the current state of proposed changes for improved error handling.
A [Go Wiki](/wiki/Go2ErrorHandlingFeedback) was created to collect related feedback, discussions, and articles.
Independently, other people have started tracking all the many error handling proposals
over the years.
It's amazing to see the sheer volume of them all, for instance in Sean K. H. Liao's blog post on
["go error handling proposals"](https://seankhliao.com/blog/12020-11-23-go-error-handling-proposals/).

The complaints about the verbosity of error handling persisted
(see [Go Developer Survey 2024 H1 Results](/blog/survey2024-h1-results)),
and so, after a series of increasingly refined Go team internal proposals, Ian Lance Taylor published
["reduce error handling boilerplate using `?`"](/issue/71203) in 2024.
This time the idea was to borrow from a construct implemented in
[Rust](https://www.rust-lang.org/), specifically the
[`?` operator](https://doc.rust-lang.org/std/result/index.html#the-question-mark-operator-).
The hope was that by leaning on an existing mechanism using an established notation, and taking into
account what we had learned over the years, we should be able to finally make some progress.
In small informal user studies where programmers were shown Go code using `?`, the vast majority
of participants correctly guessed the meaning of the code, which further convinced us to give it another
shot.
To be able to see the impact of the change, Ian wrote a tool that converts ordinary Go code
into code that uses the proposed new syntax, and we also prototyped the feature in the
compiler.

```Go
// printSum implementation using the proposed "?" statements.
func printSum(a, b string) error {
	x := strconv.Atoi(a) ?
	y := strconv.Atoi(b) ?
	fmt.Println("result:", x + y)
	return nil
}
```

Unfortunately, as with the other error handling ideas, this new proposal was also quickly overrun
with comments and many suggestions for minor tweaks, often based on individual preferences.
Ian closed the proposal and moved the content into a [discussion](/issue/71460)
to facilitate the conversation and to collect further feedback.
A slightly modified version was received
[a bit more positively](https://github.com/golang/go/discussions/71460#discussioncomment-12060294)
but broad support remained elusive.

After so many years of trying, with three full-fledged proposals by the Go team and
literally [hundreds](/issues?q=+is%3Aissue+label%3Aerror-handling) (!)
of community proposals, most of them variations on a theme,
all of which failed to attract sufficient (let alone overwhelming) support,
the question we now face is: how to proceed? Should we proceed at all?

_We think not._

To be more precise, we should stop trying to solve the _syntactic problem_, at least for the foreseeable
future.
The [proposal process](https://github.com/golang/proposal?tab=readme-ov-file#consensus-and-disagreement)
provides justification for this decision:

> The goal of the proposal process is to reach general consensus about the outcome in a timely manner.
> If proposal review cannot identify a general consensus in the discussion of the issue on the issue tracker,
> the usual result is that the proposal is declined.

Furthermore:

> It can happen that proposal review may not identify a general consensus and yet it is clear that the
> proposal should not be outright declined.
> [...]
> If the proposal review group cannot identify a consensus nor a next step for the proposal,
> the decision about the path forward passes to the Go architects [...], who review the discussion and
> aim to reach a consensus among themselves.

None of the error handling proposals reached anything close to a consensus,
so they were all declined.
Even the most senior members of the Go team at Google do not unanimously agree
on the best path forward _at this time_ (perhaps that will change at some point).
But without a strong consensus we cannot reasonably move forward.

There are valid arguments in favor of the status quo:

- If Go had introduced specific syntactic sugar for error handling early on, few would argue over it today.
But we are 15 years down the road, the opportunity has passed, and Go has
a perfectly fine way to handle errors, even if it may seem verbose at times.

- Looking from a different angle, let's assume we came across the perfect solution today.
Incorporating it into the language would simply lead from one unhappy group of users
(the one that roots for the change) to another (the one that prefers the status quo).
We were in a similar situation when we decided to add generics to the language, albeit with an
important difference:
today nobody is forced to use generics, and good generic libraries are written such that users
can mostly ignore the fact that they are generic, thanks to type inference.
On the contrary, if a new syntactic construct for error handling gets added to the language,
virtually everybody will need to start using it, lest their code become unidiomatic.

- Not adding extra syntax is in line with one of Go's design rules:
do not provide multiple ways of doing the same thing.
There are exceptions to this rule in areas with high "foot traffic": assignments come to mind.
Ironically, the ability to _redeclare_ a variable in
[short variable declarations](/ref/spec#Short_variable_declarations) (`:=`) was introduced to address a problem
that arose because of error handling:
without redeclarations, sequences of error checks require a differently named `err` variable for
each check (or additional separate variable declarations).
At that time, a better solution might have been to provide more syntactic support for error handling.
Then, the redeclaration rule may not have been needed, and with it gone, so would be various
associated [complications](/issue/377).

- Going back to actual error handling code, verbosity fades into the background if errors are
actually _handled_.
Good error handling often requires additional information added to an error.
For instance, a recurring comment in user surveys is about the lack of stack traces associated
with an error.
This could be addressed with support functions that produce and return an augmented
error.
In this (admittedly contrived) example, the relative amount of boilerplate is much smaller:

	```Go
	func printSum(a, b string) error {
		x, err := strconv.Atoi(a)
		if err != nil {
			return fmt.Errorf("invalid integer: %q", a)
		}
		y, err := strconv.Atoi(b)
		if err != nil {
			return fmt.Errorf("invalid integer: %q", b)
		}
		fmt.Println("result:", x + y)
		return nil
	}
	```

- New standard library functionality can help reduce error handling boilerplate as well,
very much in the vein of Rob Pike's 2015 blog post
["Errors are values"](/blog/errors-are-values).
For instance, in some cases [`cmp.Or`](/pkg/cmp#Or) may be used to deal with a
series of errors all at once:

	```Go
	func printSum(a, b string) error {
		x, err1 := strconv.Atoi(a)
		y, err2 := strconv.Atoi(b)
		if err := cmp.Or(err1, err2); err != nil {
			return err
		}
		fmt.Println("result:", x+y)
		return nil
	}
	```

- Writing, reading, and debugging code are all quite different activities.
Writing repeated error checks can be tedious, but today's IDEs provide powerful, even LLM-assisted
code completion.
Writing basic error checks is straightforward for these tools.
The verbosity is most obvious when reading code, but tools might help here as well;
for instance an IDE with a Go language setting could provide a toggle switch to hide error handling
code.
Such switches already exist for other code sections such as function bodies.

- When debugging error handling code, being able to quickly add a `println` or
have a dedicated line or source location for setting a breakpoint in a debugger is helpful.
This is easy when there is already a dedicated `if` statement.
But if all the error handling logic is hidden behind a `check`, `try`, or `?`, the code may have to
be changed into an ordinary `if` statement first, which complicates debugging
and may even introduce subtle bugs.

- There are also practical considerations:
Coming up with a new syntax idea for error handling is cheap;
hence the proliferation of a multitude of proposals from the community.
Coming up with a good solution that holds up to scrutiny: not so much.
It takes a concerted effort to properly design a language change and to actually implement it.
The real cost still comes afterwards:
all the code that needs to be changed, the documentation that needs to be updated,
the tools that need to be adjusted.
Taken all into account, language changes are very expensive, the Go team is relatively small,
and there are a lot of other priorities to address.
(These latter points may change: priorities can shift, team sizes can go up or down.)

- On a final note, some of us recently had the opportunity to attend
[Google Cloud Next 2025](https://cloud.withgoogle.com/next/25),
where the Go team had a booth and where we also hosted a small Go Meetup.
Every single Go user we had a chance to ask was adamant that we should not change the
language for better error handling.
Many mentioned that the lack of specific error handling support in Go is most apparent
when coming freshly from another language that has that support.
As one becomes more fluent and writes more idiomatic Go code, the issue becomes much less important.
This is of course not a sufficiently large set of people to be representative,
but it may be a different set of people than we see on GitHub, and their feedback serves as yet another data point.

Of course, there are also valid arguments in favor of change:

- Lack of better error handling support remains the top complaint in our user surveys.
If the Go team really does take user feedback seriously, we ought to do something about this eventually.
(Although there does not seem to be
[overwhelming support](https://github.com/golang/go/discussions/71460#discussioncomment-11977299)
for a language change either.)

- Perhaps the singular focus on reducing the character count is misguided.
A better approach might be to make default error handling highly visible with a keyword
while still removing boilerplate (`err != nil`).
Such an approach might make it easier for a reader (a code reviewer!) to see that an error
is handled, without "looking twice", resulting in improved code quality and safety.
This would bring us back to the beginnings of `check` and `handle`.

- We don't really know how much the issue is the straightforward syntactic verbosity of
error checking, versus the verbosity of good error handling:
constructing errors that are a useful part of an API and meaningful to developers and
end-users alike.
This is something we'd like to study in greater depth.

Still, no attempt to address error handling so far has gained sufficient traction.
If we are honestly taking stock of where we are, we can only admit that we
neither have a shared understanding of the problem,
nor do we all agree that there is a problem in the first place.
With this in mind, we are making the following pragmatic decision:

_For the foreseeable future, the Go team will stop pursuing syntactic language changes
for error handling.
We will also close all open and incoming proposals that concern themselves primarily
with the syntax of error handling, without further investigation._

The community has put tremendous effort into exploring, discussing, and debating these issues.
While this may not have resulted in any changes to error handling syntax, these efforts have
resulted in many other improvements to the Go language and our processes.
Maybe, at some point in the future, a clearer picture will emerge on error handling.
Until then, we look forward to focusing this incredible passion on new opportunities
to make Go better for everyone.

Thank you!
