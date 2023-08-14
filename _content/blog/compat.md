---
title: Backward Compatibility, Go 1.21, and Go 2
date: 2023-08-14T12:00:00Z
by:
- Russ Cox
summary: Go 1.21 expands Go's commitment to backward compatibility, so that every new Go toolchain is the best possible implementation of older toolchain semantics as well.
---

Go 1.21 includes new features to improve compatibility.
Before you stop reading, I know that sounds boring.
But boring can be good.
Back in the early days of Go 1,
Go was exciting and full of surprises.
Each week we cut a new snapshot release
and everyone got to roll the dice
to see what we’d changed
and how their programs would break.
We released Go 1 and its compatibility promise
to remove the excitement,
so that new releases of Go would be boring.

Boring is good.
Boring is stable.
Boring means being able to focus on your work,
not on what’s different about Go.
This post is about the important work we shipped in Go 1.21
to keep Go boring.

## Go 1 Compatibility {#go1}

We’ve been focused on compatibility for over a decade.
For Go 1, back in 2012, we published a document titled
“[Go 1 and the Future of Go Programs](/doc/go1compat)”
that lays out a very clear intention:

> It is intended that programs written to the Go 1 specification
> will continue to compile and run correctly, unchanged,
> over the lifetime of that specification. ...
> Go programs that work today should continue to work
> even as future releases of Go 1 arise.

There are a few qualifications to that.
First, compatibility means source compatibility.
When you update to a new version of Go,
you do have to recompile your code.
Second, we can add new APIs,
but not in a way that breaks existing code.

The end of the document warns,
 “[It] is impossible to guarantee that no future change will break any program.”
Then it lays out a number of reasons why programs might still break.

For example, it makes sense that if your program depends on a buggy behavior
and we fix the bug, your program will break.
But we try very hard to break as little as possible and keep Go boring.
There are two main approaches we’ve used so far: API checking and testing.

## API Checking {#api}

Perhaps the clearest fact about compatibility
is that we can’t take away API, or else programs using it will break.

For example, here’s a program someone has written
that we can’t break:

	package main

	import "os"

	func main() {
		os.Stdout.WriteString("hello, world\n")
	}

We can’t remove the package `os`;
we can’t remove the global variable `os.Stdout`, which is an `*os.File`;
and we also can’t remove the `os.File` method `WriteString`.
It should be clear that removing any of those would
break this program.

It’s perhaps less clear that we can’t change the type of `os.Stdout` at all.
Suppose we want to make it an interface with the same methods.
The program we just saw wouldn’t break, but this one would:

	package main

	import "os"

	func main() {
		greet(os.Stdout)
	}

	func greet(f *os.File) {
		f.WriteString(“hello, world\n”)
	}

This program passes `os.Stdout` to a function named `greet`
that requires an argument of type `*os.File`.
So changing `os.Stdout` to an interface will break this program.

To help us as we develop Go, we use a tool that maintains
a list of each package’s exported API
in files separate from the actual packages:

	% cat go/api/go1.21.txt
	pkg bytes, func ContainsFunc([]uint8, func(int32) bool) bool #54386
	pkg bytes, method (*Buffer) AvailableBuffer() []uint8 #53685
	pkg bytes, method (*Buffer) Available() int #53685
	pkg cmp, func Compare[$0 Ordered]($0, $0) int #59488
	pkg cmp, func Less[$0 Ordered]($0, $0) bool #59488
	pkg cmp, type Ordered interface {} #59488
	pkg context, func AfterFunc(Context, func()) func() bool #57928
	pkg context, func WithDeadlineCause(Context, time.Time, error) (Context, CancelFunc) #56661
	pkg context, func WithoutCancel(Context) Context #40221
	pkg context, func WithTimeoutCause(Context, time.Duration, error) (Context, CancelFunc) #56661

One of our standard tests checks that the actual package APIs match those files.
If we add new API to a package, the test breaks unless we add it to the API files.
And if we change or remove API, the test breaks too. This helps us avoid mistakes.
However, a tool like this only finds a certain class of problems, namely API changes and removals.
There are other ways to make incompatible changes to Go.

That leads us to the second approach we use to keep Go boring: testing.

## Testing {#testing}

The most effective way to find unexpected incompatibilities
is to run existing tests against the development version of the next Go release.
We test the development version of Go against
all of Google’s internal Go code on a rolling basis.
When tests are passing, we install that commit as
Google’s production Go toolchain.

If a change breaks tests inside Google,
we assume it will also break tests outside Google,
and we look for ways to reduce the impact.
Most of the time, we roll back the change entirely
or find a way to rewrite it so that it doesn’t break any programs.
Sometimes, however, we conclude that the change is
important to make and “compatible” even though it does
break some programs.
In that case, we still work to reduce the impact as much as possible,
and then we document the potential problem in the release notes.

Here are two examples of that kind of subtle compatibility problems
we found by testing Go inside Google but still included in Go 1.1.

## Struct Literals and New Fields {#struct}

Here is some code that runs fine in Go 1:

	package main

	import "net"

	var myAddr = &net.TCPAddr{
		net.IPv4(18, 26, 4, 9),
		80,
	}

Package `main` declares a global variable `myAddr`,
which is a composite literal of type `net.TCPAddr`.
In Go 1, package `net` defines the type `TCPAddr`
as a struct with two fields, `IP` and `Port`.
Those match the fields in the composite literal,
so the program compiles.

In Go 1.1, the program stopped compiling, with a compiler error
that said “too few initializers in struct literal.”
The problem is that we added a third field, `Zone`, to `net.TCPAddr`,
and this program is missing the value for that third field.
The fix is to rewrite the program using tagged literals,
so that it builds in both versions of Go:

	var myAddr = &net.TCPAddr{
		IP:   net.IPv4(18, 26, 4, 9),
		Port: 80,
	}

Since this literal doesn’t specify a value for `Zone`, it will use the
zero value (an empty string in this case).

This requirement to use composite literals for standard library
structs is explicitly called out in the [compatibility document](/doc/go1compat),
and `go vet` reports literals that need tags
to ensure compatibility with later versions of Go.
This problem was new enough in Go 1.1
to merit a short comment in the release notes.
Nowadays we just mention the new field.

## Time Precision {#precision}

The second problem we found while testing Go 1.1
had nothing to do with APIs at all.
It had to do with time.

Shortly after Go 1 was released, someone pointed
out that [`time.Now`](/pkg/time/#Now)
returned times with microsecond precision,
but with some extra code,
it could return times with nanosecond precision instead.
That sounds good, right?
More precision is better.
So we made that change.

That broke a handful of tests inside Google that
were schematically like this one:

	func TestSaveTime(t *testing.T) {
		t1 := time.Now()
		save(t1)
		if t2 := load(); t2 != t1 {
			t.Fatalf("load() = %v, want %v", t1, t2)
		}
	}

This code calls `time.Now`
and then round-trips the result
through `save` and `load`
and expects to get the same time back.
If `save` and `load` use a representation
that only stores microsecond precision,
that will work fine in Go 1 but fail in Go 1.1.

To help fix tests like this,
we added [`Round`](/pkg/time/#Time.Round) and
[`Truncate`](/pkg/time/#Time.Truncate) methods
to discard unwanted precision,
and in the release notes,
we documented the possible problem
and the new methods to help fix it.

These examples show how
testing finds different kinds of incompatibility
than the API checks do.
Of course, testing is not a complete guarantee
of compatibility either,
but it’s more complete than just API checks.
There are many examples of problems we’ve found
while testing that we decided did break the compatibility
rules and rolled back before the release.
The time precision change is an interesting example
of something that broke programs but that we released
anyway.
We made the change because the improved precision
was better and was allowed within the documented behavior of the function.

This example shows that sometimes, despite significant effort
and attention, there are times when changing Go means
breaking Go programs.
The changes are, strictly speaking, “compatible” in the sense
of the Go 1 document, but they still break programs.
Most of these compatibility issues can be placed
in one of three categories:
output changes,
input changes,
and protocol changes.

## Output Changes {#output}

An output change happens when a function gives a different output
than it used to, but the new output is just as correct as, or even more correct than,
the old output.
If existing code is written to expect only the old output, it will break.
We just saw an example of this, with `time.Now` adding nanosecond precision.

**Sort.** Another example happened in Go 1.6,
when we changed the implementation of sort
to run about 10% faster.
Here's an example program
that sorts a list of colors by length of name:

	colors := strings.Fields(
		`black white red orange yellow green blue indigo violet`)
	sort.Sort(ByLen(colors))
	fmt.Println(colors)

	Go 1.5:  [red blue green white black yellow orange indigo violet]
	Go 1.6:  [red blue white green black orange yellow indigo violet]

Changing sort algorithms often changes
how equal elements are ordered,
and that happened here.
Go 1.5 returned green, white, black, in that order.
Go 1.6 returned white, green, black.

Sort is clearly allowed to return equal results in any order it likes,
and this change made it 10% faster, which is great.
But programs that expect a specific output will break.
This is a good example of why compatibility is so hard.
We don’t want to break programs,
but we also don’t want to be locked in
to undocumented implementation details.

**Compress/flate.** As another example, in Go 1.8 we improved
`compress/flate` to produce smaller outputs,
with roughly the same CPU and memory overheads.
That sounds like a win-win, but it broke a project inside Google
that needed reproducible archive builds:
now they couldn’t reproduce their old archives.
They forked `compress/flate` and `compress/gzip`
to keep a copy of the old algorithm.

We do a similar thing with the Go compiler,
using a fork of the `sort` package ([and others](https://go.googlesource.com/go/+/go1.21.0/src/cmd/dist/buildtool.go#22))
so that the compiler produces the same results
even when it is built using earlier versions of Go.

For output change incompatibilities like these,
the best answer is to write programs and tests
that accept any valid output,
and to use these kinds of breakages
as an opportunity to change your testing strategy,
not just update the expected answers.
If you need truly reproducible outputs,
the next best answer is to fork the code
to insulate yourself from changes,
but remember that
you’re also insulating yourself from bug fixes.

## Input Changes {#input}

An input change happens when a function changes which inputs it accepts
or how it processes them.

**ParseInt.** For example, Go 1.13 added support for
underscores in large numbers for readability.
Along with the language change,
we made `strconv.ParseInt` accept the new syntax.
This change didn't break anything inside Google,
but much later we heard from an external user
whose code did break.
Their program used numbers separated by underscores
as a data format.
It tried `ParseInt` first and only fell back to checking for underscores if `ParseInt` failed.
When `ParseInt` stopped failing, the underscore-handling code stopped running.

**ParseIP.** As another example, Go’s `net.ParseIP`,
followed the examples in early IP RFCs,
which often showed decimal IP addresses with leading zeros.
It read the IP address 18.032.4.011 address as 18.32.4.11, just with a few extra zeros.
We found out much later that BSD-derived C libraries
interpret leading zeros in IP addresses as starting an octal number:
in those libraries, 18.032.4.011 means 18.26.4.9!

This was a serious mismatch
between Go and the rest of the world,
but changing the meaning of leading zeros
from one Go release to the next
would be a serious mismatch too.
It would be a huge incompatibility.
In the end, we decided to change `net.ParseIP` in Go 1.17
to reject leading zeros entirely.
This stricter parsing ensures that when Go and C
both parse an IP address successfully,
or when old and new Go versions do,
they all agree about what it means.

This change didn't break anything inside Google,
but the Kubernetes team was concerned about
saved configurations that might have parsed before
but would stop parsing with Go 1.17.
Addresses with leading zeros
probably should be removed from those configs,
since Go interprets them differently from
essentially every other language,
but that should happen on Kubernetes’s timeline, not Go’s.
To avoid the semantic change,
Kubernetes started using its own forked copy
of the original `net.ParseIP`.

The best response to input changes is to process user input
by first validating the syntax you want to accept
before parsing the values,
but sometimes you need to fork the code instead.

## Protocol Changes {#protocol}

The final common kind of incompatibility is protocol changes.
A protocol change is a change made to a package
that ends up externally visible
in the protocols a program uses
to communicate with the external world.
Almost any change can become externally visible
in certain programs, as we saw with `ParseInt` and `ParseIP`,
but protocol changes are externally visible
in essentially all programs.

**HTTP/2.** A clear example of a protocol change is when
Go 1.6 added automatic support for HTTP/2.
Suppose a Go 1.5 client is connecting to an HTTP/2-capable
server over a network with middleboxes that happen to
break HTTP/2.
Since Go 1.5 only uses HTTP/1.1, the program works fine.
But then updating to Go 1.6 breaks the program, because Go 1.6
starts using HTTP/2, and in this context, HTTP/2 doesn’t work.

Go aims to support modern protocols by default,
but this example shows that enabling HTTP/2 can break programs
through no fault of their own (nor any fault of Go's).
Developers in this situation could go back to using Go 1.5,
but that’s not very satisfying.
Instead, Go 1.6 documented the change in the release notes
and made it straightforward to disable HTTP/2.

In fact, [Go 1.6 documented two ways](/doc/go1.6#http2) to disable HTTP/2:
configure the `TLSNextProto` field explicitly using the package API,
or set the GODEBUG environment variable:

	GODEBUG=http2client=0 ./myprog
	GODEBUG=http2server=0 ./myprog
	GODEBUG=http2client=0,http2server=0 ./myprog

As we’ll see later, Go 1.21 generalizes this GODEBUG mechanism
to make it a standard for all potentially breaking changes.

**SHA1.** Here's a subtler example of a protocol change.
No one should be using SHA1-based certificates for HTTPS anymore.
Certificate authorities stopped issuing them in 2015,
and all the major browsers stopped accepting them in 2017.
In early 2020, Go 1.18 disabled support for them by default,
with a GODEBUG setting to override that change.
We also announced our intent to remove the GODEBUG setting in Go 1.19.

The Kubernetes team let us know that
some installations still use private SHA1 certificates.
Putting aside the security questions,
it's not Kubernetes’s place to force these enterprises to
upgrade their certificate infrastructure,
and it would be extremely painful to fork `crypto/tls`
and `net/http` to keep SHA1 support.
Instead, we agreed to keep the override in place longer than we had planned,
to create more time for an orderly transition.
After all, we want to break as few programs as possible.

## Expanded GODEBUG Support in Go 1.21

To improve backwards compatibility even in these subtle cases
we’ve been examining, Go 1.21 expands and formalizes the use of GODEBUG.

To begin with, for any change that is permitted by Go 1 compatibility but
still might break existing programs,
we do all the work we just saw to understand potential
compatibility problems, and we engineer the change to keep as many existing
programs working as possible.
For the remaining programs, the new approach is:

 1. We will define a new GODEBUG setting that allows
    individual programs to opt out of the new behavior.
    A GODEBUG setting may not be added if doing so is infeasible, but that should be extremely rare.

 2. GODEBUG settings added for compatibility will be maintained for a minimum
    of two years (four Go releases). Some, such as `http2client` and `http2server`,
    will be maintained much longer, even indefinitely.

 3. When possible, each GODEBUG setting has an associated
    [`runtime/metrics`](/pkg/runtime/metrics/) counter
    named `/godebug/non-default-behavior/<name>:events`
    that counts the number of times a particular program’s behavior
    has changed based on a non-default value for that setting.
    For example, when `GODEBUG=http2client=0` is set,
    `/godebug/non-default-behavior/http2client:events` counts the
    number of HTTP transports that the program has configured without HTTP/2 support.

 4. A program’s GODEBUG settings are configured to match the Go version
   listed in the main package’s `go.mod` file.
   If your program’s `go.mod` file says `go 1.20` and you update to
   a Go 1.21 toolchain, any GODEBUG-controlled behaviors changed in
   Go 1.21 will retain their old Go 1.20 behavior until you change the
   `go.mod` to say `go 1.21`.

 5. A program can change individual GODEBUG settings by using `//go:debug` lines
    in package `main`.

 6. All GODEBUG settings are documented in a [single, central list](/doc/godebug#history)
   for easy reference.

This approach means that each new version of Go should be the best possible
implementation of older versions of Go, even preserving behaviors that were
changed in compatible-but-breaking ways in later releases when compiling old code.

For example, in Go 1.21, `panic(nil)` now causes a (non-nil) runtime panic,
so that the result of [`recover`](/ref/spec/#Handling_panics) now reliably
reports whether the current goroutine is panicking.
This new behavior is controlled by a GODEBUG setting and therefore dependent
on the main package’s `go.mod`’s `go` line: if it says `go 1.20` or earlier,
`panic(nil)` is still allowed.
If it says `go 1.21` or later, `panic(nil)` turns into a panic with a `runtime.PanicNilError`.
And the version-based default can be overridden explicitly by adding a line like this to package main:

	//go:debug panicnil=1

This combination of features means that programs can update to newer toolchains
while preserving the behaviors of the earlier toolchains they used,
can apply finer-grained control over specific settings as needed,
and can use production monitoring to understand which jobs make
use of these non-default behaviors in practice.
Combined, those should make rolling out new toolchains even
smoother than in the past.

See “[Go, Backwards Compatibility, and GODEBUG](/doc/godebug)” for more details.

## An Update on Go 2 {#go2}

In the quoted text from “[Go 1 and the Future of Go Programs](/doc/go1compat)”
at the top of this post, the ellipsis hid the following qualifier:

> At some indefinite point, a Go 2 specification may arise,
> but until that time, [... all the compatibility details ...].

That raises an obvious question: when should we expect the
Go 2 specification that breaks old Go 1 programs?

The answer is never.
Go 2, in the sense of breaking with the past
and no longer compiling old programs,
is never going to happen.
Go 2 in the sense of being the major revision of Go 1
we started toward in 2017 has already happened.

There will not be a Go 2 that breaks Go 1 programs.
Instead, we are going to double down on compatibility,
which is far more valuable than any possible break with the past.
In fact, we believe that prioritizing compatibility
was the most important design decision we made for Go 1.

So what you will see over the next few years
is plenty of new, exciting work, but done in a careful,
compatible way, so that we can keep your upgrades from
one toolchain to the next as boring as possible.
