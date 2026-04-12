---
title: Go 1.26 is released
date: 2026-02-10
by:
- Carlos Amedee, on behalf of the Go team
summary: Go 1.26 adds a new garbage collector, cgo overhead reduction, experimental simd/archsimd package, experimental runtime/secret package, and more.
---

Today the Go team is pleased to release Go 1.26.
You can find its binary archives and installers on the [download page](/dl/).

## Language changes

Go 1.26 introduces two significant refinements to the language
[syntax and type system](/doc/go1.26#language).

First, the built-in `new` function, which creates a new variable, now allows its operand to be an
expression, specifying the initial value of the variable.

A simple example of this change means that code such as this:

```go
x := int64(300)
ptr := &x
```

Can be simplified to:

```go
ptr := new(int64(300))
```

Second, generic types may now refer to themselves in their own type parameter list. This change
simplifies the implementation of complex data structures and interfaces.

## Performance improvements

The previously experimental [Green Tea garbage collector](/doc/go1.26#new-garbage-collector)
is now enabled by default.

The baseline [cgo overhead has been reduced](/doc/go1.26#faster-cgo-calls)
by approximately 30%.

The compiler can now [allocate the backing store](/doc/go1.26#compiler) for
slices on the stack in more situations, which improves performance.

## Tool improvements

The `go fix` command has been completely rewritten to use the
[Go analysis framework](/pkg/golang.org/x/tools/go/analysis), and now includes a
couple dozen "[modernizers](/pkg/golang.org/x/tools/go/analysis/passes/modernize)", analyzers
that suggest safe fixes to help your code take advantage of newer features of the language
and standard library. It also includes the
[`inline` analyzer](/pkg/golang.org/x/tools/go/analysis/passes/inline#hdr-Analyzer_inline), which
attempts to inline all calls to each function annotated with a `//go:fix inline` directive.
Two upcoming blog posts will address these features in more detail.

## More improvements and changes

Go 1.26 introduces many improvements over Go 1.25 across
its [tools](/doc/go1.26#tools), the [runtime](/doc/go1.26#runtime),
[compiler](/doc/go1.26#compiler), [linker](/doc/go1.26#linker),
and the [standard library](/doc/go1.26#library).
This includes the addition of three new packages: [`crypto/hpke`](/doc/go1.26#new-cryptohpke-package),
[`crypto/mlkem/mlkemtest`](/doc/go1.26#cryptomlkempkgcryptomlkem), and
[`testing/cryptotest`](/doc/go1.26#testingcryptotestpkgtestingcryptotest).
There are [port-specific](/doc/go1.26#ports) changes
and [`GODEBUG` settings](/doc/godebug#go-126) updates.

Some of the additions in Go 1.26 are in an experimental stage
and become exposed only when you explicitly opt in. Notably:

- An [experimental `simd/archsimd` package](/doc/go1.26#simd) provides access to "single instruction,
multiple data" (SIMD) operations.

- An [experimental `runtime/secret` package](/doc/go1.26#new-experimental-runtimesecret-package) provides
a facility for securely erasing temporaries used in code that manipulates secret
information, typically cryptographic in nature.

- An [experimental `goroutineleak` profile](/doc/go1.26#goroutineleak-profiles)
in the `runtime/pprof` package that reports leaked goroutines.

These experiments are all expected to be generally available in a
future version of Go. We encourage you to try them out ahead of time.
We really value your feedback!

Please refer to the [Go 1.26 Release Notes](/doc/go1.26) for the complete list
of additions, changes, and improvements in Go 1.26.

Over the next few weeks, follow-up blog posts will cover some of the topics
relevant to Go 1.26 in more detail. Check back later to read those posts.

Thanks to everyone who contributed to this release by writing code, filing bugs,
trying out experimental additions, sharing feedback, and testing the release candidates.
Your efforts helped make Go 1.26 as stable as possible.
As always, if you notice any problems, please [file an issue](/issue/new).

We hope you enjoy using the new release!
