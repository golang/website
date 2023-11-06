---
template: true
title: Go 1 Release Notes
---

## Introduction to Go 1 {#introduction}

Go version 1, Go 1 for short, defines a language and a set of core libraries
that provide a stable foundation for creating reliable products, projects, and
publications.

The driving motivation for Go 1 is stability for its users. People should be able to
write Go programs and expect that they will continue to compile and run without
change, on a time scale of years, including in production environments such as
Google App Engine. Similarly, people should be able to write books about Go, be
able to say which version of Go the book is describing, and have that version
number still be meaningful much later.

Code that compiles in Go 1 should, with few exceptions, continue to compile and
run throughout the lifetime of that version, even as we issue updates and bug
fixes such as Go version 1.1, 1.2, and so on. Other than critical fixes, changes
made to the language and library for subsequent releases of Go 1 may
add functionality but will not break existing Go 1 programs.
[The Go 1 compatibility document](go1compat.html)
explains the compatibility guidelines in more detail.

Go 1 is a representation of Go as it used today, not a wholesale rethinking of
the language. We avoided designing new features and instead focused on cleaning
up problems and inconsistencies and improving portability. There are a number
changes to the Go language and packages that we had considered for some time and
prototyped but not released primarily because they are significant and
backwards-incompatible. Go 1 was an opportunity to get them out, which is
helpful for the long term, but also means that Go 1 introduces incompatibilities
for old programs. Fortunately, the `go` `fix` tool can
automate much of the work needed to bring programs up to the Go 1 standard.

This document outlines the major changes in Go 1 that will affect programmers
updating existing code; its reference point is the prior release, r60 (tagged as
r60.3). It also explains how to update code from r60 to run under Go 1.

## Changes to the language {#language}

### Append {#append}

The `append` predeclared variadic function makes it easy to grow a slice
by adding elements to the end.
A common use is to add bytes to the end of a byte slice when generating output.
However, `append` did not provide a way to append a string to a `[]byte`,
which is another common case.

{{code "/doc/progs/go1.go" `/greeting := ..byte/` `/append.*hello/`}}

By analogy with the similar property of `copy`, Go 1
permits a string to be appended (byte-wise) directly to a byte
slice, reducing the friction between strings and byte slices.
The conversion is no longer necessary:

{{code "/doc/progs/go1.go" `/append.*world/`}}

_Updating_:
This is a new feature, so existing code needs no changes.

### Close {#close}

The `close` predeclared function provides a mechanism
for a sender to signal that no more values will be sent.
It is important to the implementation of `for` `range`
loops over channels and is helpful in other situations.
Partly by design and partly because of race conditions that can occur otherwise,
it is intended for use only by the goroutine sending on the channel,
not by the goroutine receiving data.
However, before Go 1 there was no compile-time checking that `close`
was being used correctly.

To close this gap, at least in part, Go 1 disallows `close` on receive-only channels.
Attempting to close such a channel is a compile-time error.

{{code "/doc/progs/go1.go" `/STARTCLOSE/` `/ENDCLOSE/`}}

_Updating_:
Existing code that attempts to close a receive-only channel was
erroneous even before Go 1 and should be fixed. The compiler will
now reject such code.

### Composite literals {#literals}

In Go 1, a composite literal of array, slice, or map type can elide the
type specification for the elements' initializers if they are of pointer type.
All four of the initializations in this example are legal; the last one was illegal before Go 1.

{{code "/doc/progs/go1.go" `/type Date struct/` `/STOP/`}}

_Updating_:
This change has no effect on existing code, but the command
`gofmt` `-s` applied to existing source
will, among other things, elide explicit element types wherever permitted.

### Goroutines during init {#init}

The old language defined that `go` statements executed during initialization created goroutines but that they did not begin to run until initialization of the entire program was complete.
This introduced clumsiness in many places and, in effect, limited the utility
of the `init` construct:
if it was possible for another package to use the library during initialization, the library
was forced to avoid goroutines.
This design was done for reasons of simplicity and safety but,
as our confidence in the language grew, it seemed unnecessary.
Running goroutines during initialization is no more complex or unsafe than running them during normal execution.

In Go 1, code that uses goroutines can be called from
`init` routines and global initialization expressions
without introducing a deadlock.

{{code "/doc/progs/go1.go" `/PackageGlobal/` `/^}/`}}

_Updating_:
This is a new feature, so existing code needs no changes,
although it's possible that code that depends on goroutines not starting before `main` will break.
There was no such code in the standard repository.

### The rune type {#rune}

The language spec allows the `int` type to be 32 or 64 bits wide, but current implementations set `int` to 32 bits even on 64-bit platforms.
It would be preferable to have `int` be 64 bits on 64-bit platforms.
(There are important consequences for indexing large slices.)
However, this change would waste space when processing Unicode characters with
the old language because the `int` type was also used to hold Unicode code points: each code point would waste an extra 32 bits of storage if `int` grew from 32 bits to 64.

To make changing to 64-bit `int` feasible,
Go 1 introduces a new basic type, `rune`, to represent
individual Unicode code points.
It is an alias for `int32`, analogous to `byte`
as an alias for `uint8`.

Character literals such as `'a'`, `'語'`, and `'\u0345'`
now have default type `rune`,
analogous to `1.0` having default type `float64`.
A variable initialized to a character constant will therefore
have type `rune` unless otherwise specified.

Libraries have been updated to use `rune` rather than `int`
when appropriate. For instance, the functions `unicode.ToLower` and
relatives now take and return a `rune`.

{{code "/doc/progs/go1.go" `/STARTRUNE/` `/ENDRUNE/`}}

_Updating_:
Most source code will be unaffected by this because the type inference from
`:=` initializers introduces the new type silently, and it propagates
from there.
Some code may get type errors that a trivial conversion will resolve.

### The error type {#error}

Go 1 introduces a new built-in type, `error`, which has the following definition:

	    type error interface {
	        Error() string
	    }

Since the consequences of this type are all in the package library,
it is discussed [below](#errors).

### Deleting from maps {#delete}

In the old language, to delete the entry with key `k` from map `m`, one wrote the statement,

	    m[k] = value, false

This syntax was a peculiar special case, the only two-to-one assignment.
It required passing a value (usually ignored) that is evaluated but discarded,
plus a boolean that was nearly always the constant `false`.
It did the job but was odd and a point of contention.

In Go 1, that syntax has gone; instead there is a new built-in
function, `delete`. The call

{{code "/doc/progs/go1.go" `/delete\(m, k\)/`}}

will delete the map entry retrieved by the expression `m[k]`.
There is no return value. Deleting a non-existent entry is a no-op.

_Updating_:
Running `go` `fix` will convert expressions of the form `m[k] = value,
false` into `delete(m, k)` when it is clear that
the ignored value can be safely discarded from the program and
`false` refers to the predefined boolean constant.
The fix tool
will flag other uses of the syntax for inspection by the programmer.

### Iterating in maps {#iteration}

The old language specification did not define the order of iteration for maps,
and in practice it differed across hardware platforms.
This caused tests that iterated over maps to be fragile and non-portable, with the
unpleasant property that a test might always pass on one machine but break on another.

In Go 1, the order in which elements are visited when iterating
over a map using a `for` `range` statement
is defined to be unpredictable, even if the same loop is run multiple
times with the same map.
Code should not assume that the elements are visited in any particular order.

This change means that code that depends on iteration order is very likely to break early and be fixed long before it becomes a problem.
Just as important, it allows the map implementation to ensure better map balancing even when programs are using range loops to select an element from a mapl.

{{code "/doc/progs/go1.go" `/Sunday/` `/^	}/`}}

_Updating_:
This is one change where tools cannot help. Most existing code
will be unaffected, but some programs may break or misbehave; we
recommend manual checking of all range statements over maps to
verify they do not depend on iteration order. There were a few such
examples in the standard repository; they have been fixed.
Note that it was already incorrect to depend on the iteration order, which
was unspecified. This change codifies the unpredictability.

### Multiple assignment {#multiple_assignment}

The language specification has long guaranteed that in assignments
the right-hand-side expressions are all evaluated before any left-hand-side expressions are assigned.
To guarantee predictable behavior,
Go 1 refines the specification further.

If the left-hand side of the assignment
statement contains expressions that require evaluation, such as
function calls or array indexing operations, these will all be done
using the usual left-to-right rule before any variables are assigned
their value. Once everything is evaluated, the actual assignments
proceed in left-to-right order.

These examples illustrate the behavior.

{{code "/doc/progs/go1.go" `/sa :=/` `/then sc.0. = 2/`}}

_Updating_:
This is one change where tools cannot help, but breakage is unlikely.
No code in the standard repository was broken by this change, and code
that depended on the previous unspecified behavior was already incorrect.

### Returns and shadowed variables {#shadowing}

A common mistake is to use `return` (without arguments) after an assignment to a variable that has the same name as a result variable but is not the same variable.
This situation is called _shadowing_: the result variable has been shadowed by another variable with the same name declared in an inner scope.

In functions with named return values,
the Go 1 compilers disallow return statements without arguments if any of the named return values is shadowed at the point of the return statement.
(It isn't part of the specification, because this is one area we are still exploring;
the situation is analogous to the compilers rejecting functions that do not end with an explicit return statement.)

This function implicitly returns a shadowed return value and will be rejected by the compiler:

{{code "/doc/progs/go1.go" `/^func Bug/` `/^}$/`}}

_Updating_:
Code that shadows return values in this way will be rejected by the compiler and will need to be fixed by hand.
The few cases that arose in the standard repository were mostly bugs.

### Copying structs with unexported fields {#unexported}

The old language did not allow a package to make a copy of a struct value containing unexported fields belonging to a different package.
There was, however, a required exception for a method receiver;
also, the implementations of `copy` and `append` have never honored the restriction.

Go 1 will allow packages to copy struct values containing unexported fields from other packages.
Besides resolving the inconsistency,
this change admits a new kind of API: a package can return an opaque value without resorting to a pointer or interface.
The new implementations of `time.Time` and
`reflect.Value` are examples of types taking advantage of this new property.

As an example, if package `p` includes the definitions,

	    type Struct struct {
	        Public int
	        secret int
	    }
	    func NewStruct(a int) Struct {  // Note: not a pointer.
	        return Struct{a, f(a)}
	    }
	    func (s Struct) String() string {
	        return fmt.Sprintf("{%d (secret %d)}", s.Public, s.secret)
	    }

a package that imports `p` can assign and copy values of type
`p.Struct` at will.
Behind the scenes the unexported fields will be assigned and copied just
as if they were exported,
but the client code will never be aware of them. The code

	    import "p"

	    myStruct := p.NewStruct(23)
	    copyOfMyStruct := myStruct
	    fmt.Println(myStruct, copyOfMyStruct)

will show that the secret field of the struct has been copied to the new value.

_Updating_:
This is a new feature, so existing code needs no changes.

### Equality {#equality}

Before Go 1, the language did not define equality on struct and array values.
This meant,
among other things, that structs and arrays could not be used as map keys.
On the other hand, Go did define equality on function and map values.
Function equality was problematic in the presence of closures
(when are two closures equal?)
while map equality compared pointers, not the maps' content, which was usually
not what the user would want.

Go 1 addressed these issues.
First, structs and arrays can be compared for equality and inequality
(`==` and `!=`),
and therefore be used as map keys,
provided they are composed from elements for which equality is also defined,
using element-wise comparison.

{{code "/doc/progs/go1.go" `/type Day struct/` `/Printf/`}}

Second, Go 1 removes the definition of equality for function values,
except for comparison with `nil`.
Finally, map equality is gone too, also except for comparison with `nil`.

Note that equality is still undefined for slices, for which the
calculation is in general infeasible. Also note that the ordered
comparison operators (< <=
`>` `>=`) are still undefined for
structs and arrays.

_Updating_:
Struct and array equality is a new feature, so existing code needs no changes.
Existing code that depends on function or map equality will be
rejected by the compiler and will need to be fixed by hand.
Few programs will be affected, but the fix may require some
redesign.

## The package hierarchy {#packages}

Go 1 addresses many deficiencies in the old standard library and
cleans up a number of packages, making them more internally consistent
and portable.

This section describes how the packages have been rearranged in Go 1.
Some have moved, some have been renamed, some have been deleted.
New packages are described in later sections.

### The package hierarchy {#hierarchy}

Go 1 has a rearranged package hierarchy that groups related items
into subdirectories. For instance, `utf8` and
`utf16` now occupy subdirectories of `unicode`.
Also, [some packages](#subrepo) have moved into
subrepositories of
[`code.google.com/p/go`](https://code.google.com/p/go)
while [others](#deleted) have been deleted outright.

<table class="codetable" frame="border" summary="Moved packages">
<colgroup align="left" width="60%"></colgroup>
<colgroup align="left" width="40%"></colgroup>
<tbody><tr>
<th align="left">Old path</th>
<th align="left">New path</th>
</tr>
<tr>
<td colspan="2"><hr></hr></td>
</tr>
<tr><td>asn1</td> <td>encoding/asn1</td></tr>
<tr><td>csv</td> <td>encoding/csv</td></tr>
<tr><td>gob</td> <td>encoding/gob</td></tr>
<tr><td>json</td> <td>encoding/json</td></tr>
<tr><td>xml</td> <td>encoding/xml</td></tr>
<tr>
<td colspan="2"><hr></hr></td>
</tr>
<tr><td>exp/template/html</td> <td>html/template</td></tr>
<tr>
<td colspan="2"><hr></hr></td>
</tr>
<tr><td>big</td> <td>math/big</td></tr>
<tr><td>cmath</td> <td>math/cmplx</td></tr>
<tr><td>rand</td> <td>math/rand</td></tr>
<tr>
<td colspan="2"><hr></hr></td>
</tr>
<tr><td>http</td> <td>net/http</td></tr>
<tr><td>http/cgi</td> <td>net/http/cgi</td></tr>
<tr><td>http/fcgi</td> <td>net/http/fcgi</td></tr>
<tr><td>http/httptest</td> <td>net/http/httptest</td></tr>
<tr><td>http/pprof</td> <td>net/http/pprof</td></tr>
<tr><td>mail</td> <td>net/mail</td></tr>
<tr><td>rpc</td> <td>net/rpc</td></tr>
<tr><td>rpc/jsonrpc</td> <td>net/rpc/jsonrpc</td></tr>
<tr><td>smtp</td> <td>net/smtp</td></tr>
<tr><td>url</td> <td>net/url</td></tr>
<tr>
<td colspan="2"><hr></hr></td>
</tr>
<tr><td>exec</td> <td>os/exec</td></tr>
<tr>
<td colspan="2"><hr></hr></td>
</tr>
<tr><td>scanner</td> <td>text/scanner</td></tr>
<tr><td>tabwriter</td> <td>text/tabwriter</td></tr>
<tr><td>template</td> <td>text/template</td></tr>
<tr><td>template/parse</td> <td>text/template/parse</td></tr>
<tr>
<td colspan="2"><hr></hr></td>
</tr>
<tr><td>utf8</td> <td>unicode/utf8</td></tr>
<tr><td>utf16</td> <td>unicode/utf16</td></tr>
</tbody></table>

Note that the package names for the old `cmath` and
`exp/template/html` packages have changed to `cmplx`
and `template`.

_Updating_:
Running `go` `fix` will update all imports and package renames for packages that
remain inside the standard repository. Programs that import packages
that are no longer in the standard repository will need to be edited
by hand.

### The package tree exp {#exp}

Because they are not standardized, the packages under the `exp` directory will not be available in the
standard Go 1 release distributions, although they will be available in source code form
in [the repository](https://code.google.com/p/go/) for
developers who wish to use them.

Several packages have moved under `exp` at the time of Go 1's release:

  - `ebnf`
  - `html`<sup>†</sup>
  - `go/types`

(<sup>†</sup>The `EscapeString` and `UnescapeString` types remain
in package `html`.)

All these packages are available under the same names, with the prefix `exp/`: `exp/ebnf` etc.

Also, the `utf8.String` type has been moved to its own package, `exp/utf8string`.

Finally, the `gotype` command now resides in `exp/gotype`, while
`ebnflint` is now in `exp/ebnflint`.
If they are installed, they now reside in `$GOROOT/bin/tool`.

_Updating_:
Code that uses packages in `exp` will need to be updated by hand,
or else compiled from an installation that has `exp` available.
The `go` `fix` tool or the compiler will complain about such uses.

### The package tree old {#old}

Because they are deprecated, the packages under the `old` directory will not be available in the
standard Go 1 release distributions, although they will be available in source code form for
developers who wish to use them.

The packages in their new locations are:

  - `old/netchan`

_Updating_:
Code that uses packages now in `old` will need to be updated by hand,
or else compiled from an installation that has `old` available.
The `go` `fix` tool will warn about such uses.

### Deleted packages {#deleted}

Go 1 deletes several packages outright:

  - `container/vector`
  - `exp/datafmt`
  - `go/typechecker`
  - `old/regexp`
  - `old/template`
  - `try`

and also the command `gotry`.

_Updating_:
Code that uses `container/vector` should be updated to use
slices directly. See
[the Go
Language Community Wiki](https://code.google.com/p/go-wiki/wiki/SliceTricks) for some suggestions.
Code that uses the other packages (there should be almost zero) will need to be rethought.

### Packages moving to subrepositories {#subrepo}

Go 1 has moved a number of packages into other repositories, usually sub-repositories of
[the main Go repository](https://code.google.com/p/go/).
This table lists the old and new import paths:

<table class="codetable" frame="border" summary="Sub-repositories">
<colgroup align="left" width="40%"></colgroup>
<colgroup align="left" width="60%"></colgroup>
<tbody><tr>
<th align="left">Old</th>
<th align="left">New</th>
</tr>
<tr>
<td colspan="2"><hr></hr></td>
</tr>
<tr><td>crypto/bcrypt</td> <td>code.google.com/p/go.crypto/bcrypt</td></tr>
<tr><td>crypto/blowfish</td> <td>code.google.com/p/go.crypto/blowfish</td></tr>
<tr><td>crypto/cast5</td> <td>code.google.com/p/go.crypto/cast5</td></tr>
<tr><td>crypto/md4</td> <td>code.google.com/p/go.crypto/md4</td></tr>
<tr><td>crypto/ocsp</td> <td>code.google.com/p/go.crypto/ocsp</td></tr>
<tr><td>crypto/openpgp</td> <td>code.google.com/p/go.crypto/openpgp</td></tr>
<tr><td>crypto/openpgp/armor</td> <td>code.google.com/p/go.crypto/openpgp/armor</td></tr>
<tr><td>crypto/openpgp/elgamal</td> <td>code.google.com/p/go.crypto/openpgp/elgamal</td></tr>
<tr><td>crypto/openpgp/errors</td> <td>code.google.com/p/go.crypto/openpgp/errors</td></tr>
<tr><td>crypto/openpgp/packet</td> <td>code.google.com/p/go.crypto/openpgp/packet</td></tr>
<tr><td>crypto/openpgp/s2k</td> <td>code.google.com/p/go.crypto/openpgp/s2k</td></tr>
<tr><td>crypto/ripemd160</td> <td>code.google.com/p/go.crypto/ripemd160</td></tr>
<tr><td>crypto/twofish</td> <td>code.google.com/p/go.crypto/twofish</td></tr>
<tr><td>crypto/xtea</td> <td>code.google.com/p/go.crypto/xtea</td></tr>
<tr><td>exp/ssh</td> <td>code.google.com/p/go.crypto/ssh</td></tr>
<tr>
<td colspan="2"><hr></hr></td>
</tr>
<tr><td>image/bmp</td> <td>code.google.com/p/go.image/bmp</td></tr>
<tr><td>image/tiff</td> <td>code.google.com/p/go.image/tiff</td></tr>
<tr>
<td colspan="2"><hr></hr></td>
</tr>
<tr><td>net/dict</td> <td>code.google.com/p/go.net/dict</td></tr>
<tr><td>net/websocket</td> <td>code.google.com/p/go.net/websocket</td></tr>
<tr><td>exp/spdy</td> <td>code.google.com/p/go.net/spdy</td></tr>
<tr>
<td colspan="2"><hr></hr></td>
</tr>
<tr><td>encoding/git85</td> <td>code.google.com/p/go.codereview/git85</td></tr>
<tr><td>patch</td> <td>code.google.com/p/go.codereview/patch</td></tr>
<tr>
<td colspan="2"><hr></hr></td>
</tr>
<tr><td>exp/wingui</td> <td>code.google.com/p/gowingui</td></tr>
</tbody></table>

_Updating_:
Running `go` `fix` will update imports of these packages to use the new import paths.
Installations that depend on these packages will need to install them using
a `go get` command.

## Major changes to the library {#major}

This section describes significant changes to the core libraries, the ones that
affect the most programs.

### The error type and errors package {#errors}

The placement of `os.Error` in package `os` is mostly historical: errors first came up when implementing package `os`, and they seemed system-related at the time.
Since then it has become clear that errors are more fundamental than the operating system. For example, it would be nice to use `Errors` in packages that `os` depends on, like `syscall`.
Also, having `Error` in `os` introduces many dependencies on `os` that would otherwise not exist.

Go 1 solves these problems by introducing a built-in `error` interface type and a separate `errors` package (analogous to `bytes` and `strings`) that contains utility functions.
It replaces `os.NewError` with
[`errors.New`](/pkg/errors/#New),
giving errors a more central place in the environment.

So the widely-used `String` method does not cause accidental satisfaction
of the `error` interface, the `error` interface uses instead
the name `Error` for that method:

	    type error interface {
	        Error() string
	    }

The `fmt` library automatically invokes `Error`, as it already
does for `String`, for easy printing of error values.

{{code "/doc/progs/go1.go" `/START ERROR EXAMPLE/` `/END ERROR EXAMPLE/`}}

All standard packages have been updated to use the new interface; the old `os.Error` is gone.

A new package, [`errors`](/pkg/errors/), contains the function

	func New(text string) error

to turn a string into an error. It replaces the old `os.NewError`.

{{code "/doc/progs/go1.go" `/ErrSyntax/`}}

_Updating_:
Running `go` `fix` will update almost all code affected by the change.
Code that defines error types with a `String` method will need to be updated
by hand to rename the methods to `Error`.

### System call errors {#errno}

The old `syscall` package, which predated `os.Error`
(and just about everything else),
returned errors as `int` values.
In turn, the `os` package forwarded many of these errors, such
as `EINVAL`, but using a different set of errors on each platform.
This behavior was unpleasant and unportable.

In Go 1, the
[`syscall`](/pkg/syscall/)
package instead returns an `error` for system call errors.
On Unix, the implementation is done by a
[`syscall.Errno`](/pkg/syscall/#Errno) type
that satisfies `error` and replaces the old `os.Errno`.

The changes affecting `os.EINVAL` and relatives are
described [elsewhere](#os).

_Updating_:
Running `go` `fix` will update almost all code affected by the change.
Regardless, most code should use the `os` package
rather than `syscall` and so will be unaffected.

### Time {#time}

Time is always a challenge to support well in a programming language.
The old Go `time` package had `int64` units, no
real type safety,
and no distinction between absolute times and durations.

One of the most sweeping changes in the Go 1 library is therefore a
complete redesign of the
[`time`](/pkg/time/) package.
Instead of an integer number of nanoseconds as an `int64`,
and a separate `*time.Time` type to deal with human
units such as hours and years,
there are now two fundamental types:
[`time.Time`](/pkg/time/#Time)
(a value, so the `*` is gone), which represents a moment in time;
and [`time.Duration`](/pkg/time/#Duration),
which represents an interval.
Both have nanosecond resolution.
A `Time` can represent any time into the ancient
past and remote future, while a `Duration` can
span plus or minus only about 290 years.
There are methods on these types, plus a number of helpful
predefined constant durations such as `time.Second`.

Among the new methods are things like
[`Time.Add`](/pkg/time/#Time.Add),
which adds a `Duration` to a `Time`, and
[`Time.Sub`](/pkg/time/#Time.Sub),
which subtracts two `Times` to yield a `Duration`.

The most important semantic change is that the Unix epoch (Jan 1, 1970) is now
relevant only for those functions and methods that mention Unix:
[`time.Unix`](/pkg/time/#Unix)
and the [`Unix`](/pkg/time/#Time.Unix)
and [`UnixNano`](/pkg/time/#Time.UnixNano) methods
of the `Time` type.
In particular,
[`time.Now`](/pkg/time/#Now)
returns a `time.Time` value rather than, in the old
API, an integer nanosecond count since the Unix epoch.

{{code "/doc/progs/go1.go" `/sleepUntil/` `/^}/`}}

The new types, methods, and constants have been propagated through
all the standard packages that use time, such as `os` and
its representation of file time stamps.

_Updating_:
The `go` `fix` tool will update many uses of the old `time` package to use the new
types and methods, although it does not replace values such as `1e9`
representing nanoseconds per second.
Also, because of type changes in some of the values that arise,
some of the expressions rewritten by the fix tool may require
further hand editing; in such cases the rewrite will include
the correct function or method for the old functionality, but
may have the wrong type or require further analysis.

## Minor changes to the library {#minor}

This section describes smaller changes, such as those to less commonly
used packages or that affect
few programs beyond the need to run `go` `fix`.
This category includes packages that are new in Go 1.
Collectively they improve portability, regularize behavior, and
make the interfaces more modern and Go-like.

### The archive/zip package {#archive_zip}

In Go 1, [`*zip.Writer`](/pkg/archive/zip/#Writer) no
longer has a `Write` method. Its presence was a mistake.

_Updating_:
What little code is affected will be caught by the compiler and must be updated by hand.

### The bufio package {#bufio}

In Go 1, [`bufio.NewReaderSize`](/pkg/bufio/#NewReaderSize)
and
[`bufio.NewWriterSize`](/pkg/bufio/#NewWriterSize)
functions no longer return an error for invalid sizes.
If the argument size is too small or invalid, it is adjusted.

_Updating_:
Running `go` `fix` will update calls that assign the error to \_.
Calls that aren't fixed will be caught by the compiler and must be updated by hand.

### The compress/flate, compress/gzip and compress/zlib packages {#compress}

In Go 1, the `NewWriterXxx` functions in
[`compress/flate`](/pkg/compress/flate),
[`compress/gzip`](/pkg/compress/gzip) and
[`compress/zlib`](/pkg/compress/zlib)
all return `(*Writer, error)` if they take a compression level,
and `*Writer` otherwise. Package `gzip`'s
`Compressor` and `Decompressor` types have been renamed
to `Writer` and `Reader`. Package `flate`'s
`WrongValueError` type has been removed.

_Updating_
Running `go` `fix` will update old names and calls that assign the error to \_.
Calls that aren't fixed will be caught by the compiler and must be updated by hand.

### The crypto/aes and crypto/des packages {#crypto_aes_des}

In Go 1, the `Reset` method has been removed. Go does not guarantee
that memory is not copied and therefore this method was misleading.

The cipher-specific types `*aes.Cipher`, `*des.Cipher`,
and `*des.TripleDESCipher` have been removed in favor of
`cipher.Block`.

_Updating_:
Remove the calls to Reset. Replace uses of the specific cipher types with
cipher.Block.

### The crypto/elliptic package {#crypto_elliptic}

In Go 1, [`elliptic.Curve`](/pkg/crypto/elliptic/#Curve)
has been made an interface to permit alternative implementations. The curve
parameters have been moved to the
[`elliptic.CurveParams`](/pkg/crypto/elliptic/#CurveParams)
structure.

_Updating_:
Existing users of `*elliptic.Curve` will need to change to
simply `elliptic.Curve`. Calls to `Marshal`,
`Unmarshal` and `GenerateKey` are now functions
in `crypto/elliptic` that take an `elliptic.Curve`
as their first argument.

### The crypto/hmac package {#crypto_hmac}

In Go 1, the hash-specific functions, such as `hmac.NewMD5`, have
been removed from `crypto/hmac`. Instead, `hmac.New` takes
a function that returns a `hash.Hash`, such as `md5.New`.

_Updating_:
Running `go` `fix` will perform the needed changes.

### The crypto/x509 package {#crypto_x509}

In Go 1, the
[`CreateCertificate`](/pkg/crypto/x509/#CreateCertificate)
function and
[`CreateCRL`](/pkg/crypto/x509/#Certificate.CreateCRL)
method in `crypto/x509` have been altered to take an
`interface{}` where they previously took a `*rsa.PublicKey`
or `*rsa.PrivateKey`. This will allow other public key algorithms
to be implemented in the future.

_Updating_:
No changes will be needed.

### The encoding/binary package {#encoding_binary}

In Go 1, the `binary.TotalSize` function has been replaced by
[`Size`](/pkg/encoding/binary/#Size),
which takes an `interface{}` argument rather than
a `reflect.Value`.

_Updating_:
What little code is affected will be caught by the compiler and must be updated by hand.

### The encoding/xml package {#encoding_xml}

In Go 1, the [`xml`](/pkg/encoding/xml/) package
has been brought closer in design to the other marshaling packages such
as [`encoding/gob`](/pkg/encoding/gob/).

The old `Parser` type is renamed
[`Decoder`](/pkg/encoding/xml/#Decoder) and has a new
[`Decode`](/pkg/encoding/xml/#Decoder.Decode) method. An
[`Encoder`](/pkg/encoding/xml/#Encoder) type was also introduced.

The functions [`Marshal`](/pkg/encoding/xml/#Marshal)
and [`Unmarshal`](/pkg/encoding/xml/#Unmarshal)
work with `[]byte` values now. To work with streams,
use the new [`Encoder`](/pkg/encoding/xml/#Encoder)
and [`Decoder`](/pkg/encoding/xml/#Decoder) types.

When marshaling or unmarshaling values, the format of supported flags in
field tags has changed to be closer to the
[`json`](/pkg/encoding/json) package
(`` `xml:"name,flag"` ``). The matching done between field tags, field
names, and the XML attribute and element names is now case-sensitive.
The `XMLName` field tag, if present, must also match the name
of the XML element being marshaled.

_Updating_:
Running `go` `fix` will update most uses of the package except for some calls to
`Unmarshal`. Special care must be taken with field tags,
since the fix tool will not update them and if not fixed by hand they will
misbehave silently in some cases. For example, the old
`"attr"` is now written `",attr"` while plain
`"attr"` remains valid but with a different meaning.

### The expvar package {#expvar}

In Go 1, the `RemoveAll` function has been removed.
The `Iter` function and Iter method on `*Map` have
been replaced by
[`Do`](/pkg/expvar/#Do)
and
[`(*Map).Do`](/pkg/expvar/#Map.Do).

_Updating_:
Most code using `expvar` will not need changing. The rare code that used
`Iter` can be updated to pass a closure to `Do` to achieve the same effect.

### The flag package {#flag}

In Go 1, the interface [`flag.Value`](/pkg/flag/#Value) has changed slightly.
The `Set` method now returns an `error` instead of
a `bool` to indicate success or failure.

There is also a new kind of flag, `Duration`, to support argument
values specifying time intervals.
Values for such flags must be given units, just as `time.Duration`
formats them: `10s`, `1h30m`, etc.

{{code "/doc/progs/go1.go" `/timeout/`}}

_Updating_:
Programs that implement their own flags will need minor manual fixes to update their
`Set` methods.
The `Duration` flag is new and affects no existing code.

### The go/\* packages {#go}

Several packages under `go` have slightly revised APIs.

A concrete `Mode` type was introduced for configuration mode flags
in the packages
[`go/scanner`](/pkg/go/scanner/),
[`go/parser`](/pkg/go/parser/),
[`go/printer`](/pkg/go/printer/), and
[`go/doc`](/pkg/go/doc/).

The modes `AllowIllegalChars` and `InsertSemis` have been removed
from the [`go/scanner`](/pkg/go/scanner/) package. They were mostly
useful for scanning text other then Go source files. Instead, the
[`text/scanner`](/pkg/text/scanner/) package should be used
for that purpose.

The [`ErrorHandler`](/pkg/go/scanner/#ErrorHandler) provided
to the scanner's [`Init`](/pkg/go/scanner/#Scanner.Init) method is
now simply a function rather than an interface. The `ErrorVector` type has
been removed in favor of the (existing) [`ErrorList`](/pkg/go/scanner/#ErrorList)
type, and the `ErrorVector` methods have been migrated. Instead of embedding
an `ErrorVector` in a client of the scanner, now a client should maintain
an `ErrorList`.

The set of parse functions provided by the [`go/parser`](/pkg/go/parser/)
package has been reduced to the primary parse function
[`ParseFile`](/pkg/go/parser/#ParseFile), and a couple of
convenience functions [`ParseDir`](/pkg/go/parser/#ParseDir)
and [`ParseExpr`](/pkg/go/parser/#ParseExpr).

The [`go/printer`](/pkg/go/printer/) package supports an additional
configuration mode [`SourcePos`](/pkg/go/printer/#Mode);
if set, the printer will emit `//line` comments such that the generated
output contains the original source code position information. The new type
[`CommentedNode`](/pkg/go/printer/#CommentedNode) can be
used to provide comments associated with an arbitrary
[`ast.Node`](/pkg/go/ast/#Node) (until now only
[`ast.File`](/pkg/go/ast/#File) carried comment information).

The type names of the [`go/doc`](/pkg/go/doc/) package have been
streamlined by removing the `Doc` suffix: `PackageDoc`
is now `Package`, `ValueDoc` is `Value`, etc.
Also, all types now consistently have a `Name` field (or `Names`,
in the case of type `Value`) and `Type.Factories` has become
`Type.Funcs`.
Instead of calling `doc.NewPackageDoc(pkg, importpath)`,
documentation for a package is created with:

	    doc.New(pkg, importpath, mode)

where the new `mode` parameter specifies the operation mode:
if set to [`AllDecls`](/pkg/go/doc/#AllDecls), all declarations
(not just exported ones) are considered.
The function `NewFileDoc` was removed, and the function
`CommentText` has become the method
[`Text`](/pkg/go/ast/#CommentGroup.Text) of
[`ast.CommentGroup`](/pkg/go/ast/#CommentGroup).

In package [`go/token`](/pkg/go/token/), the
[`token.FileSet`](/pkg/go/token/#FileSet) method `Files`
(which originally returned a channel of `*token.File`s) has been replaced
with the iterator [`Iterate`](/pkg/go/token/#FileSet.Iterate) that
accepts a function argument instead.

In package [`go/build`](/pkg/go/build/), the API
has been nearly completely replaced.
The package still computes Go package information
but it does not run the build: the `Cmd` and `Script`
types are gone.
(To build code, use the new
[`go`](/cmd/go/) command instead.)
The `DirInfo` type is now named
[`Package`](/pkg/go/build/#Package).
`FindTree` and `ScanDir` are replaced by
[`Import`](/pkg/go/build/#Import)
and
[`ImportDir`](/pkg/go/build/#ImportDir).

_Updating_:
Code that uses packages in `go` will have to be updated by hand; the
compiler will reject incorrect uses. Templates used in conjunction with any of the
`go/doc` types may need manual fixes; the renamed fields will lead
to run-time errors.

### The hash package {#hash}

In Go 1, the definition of [`hash.Hash`](/pkg/hash/#Hash) includes
a new method, `BlockSize`. This new method is used primarily in the
cryptographic libraries.

The `Sum` method of the
[`hash.Hash`](/pkg/hash/#Hash) interface now takes a
`[]byte` argument, to which the hash value will be appended.
The previous behavior can be recreated by adding a `nil` argument to the call.

_Updating_:
Existing implementations of `hash.Hash` will need to add a
`BlockSize` method. Hashes that process the input one byte at
a time can implement `BlockSize` to return 1.
Running `go` `fix` will update calls to the `Sum` methods of the various
implementations of `hash.Hash`.

_Updating_:
Since the package's functionality is new, no updating is necessary.

### The http package {#http}

In Go 1 the [`http`](/pkg/net/http/) package is refactored,
putting some of the utilities into a
[`httputil`](/pkg/net/http/httputil/) subdirectory.
These pieces are only rarely needed by HTTP clients.
The affected items are:

  - ClientConn
  - DumpRequest
  - DumpRequestOut
  - DumpResponse
  - NewChunkedReader
  - NewChunkedWriter
  - NewClientConn
  - NewProxyClientConn
  - NewServerConn
  - NewSingleHostReverseProxy
  - ReverseProxy
  - ServerConn

The `Request.RawURL` field has been removed; it was a
historical artifact.

The `Handle` and `HandleFunc`
functions, and the similarly-named methods of `ServeMux`,
now panic if an attempt is made to register the same pattern twice.

_Updating_:
Running `go` `fix` will update the few programs that are affected except for
uses of `RawURL`, which must be fixed by hand.

### The image package {#image}

The [`image`](/pkg/image/) package has had a number of
minor changes, rearrangements and renamings.

Most of the color handling code has been moved into its own package,
[`image/color`](/pkg/image/color/).
For the elements that moved, a symmetry arises; for instance,
each pixel of an
[`image.RGBA`](/pkg/image/#RGBA)
is a
[`color.RGBA`](/pkg/image/color/#RGBA).

The old `image/ycbcr` package has been folded, with some
renamings, into the
[`image`](/pkg/image/)
and
[`image/color`](/pkg/image/color/)
packages.

The old `image.ColorImage` type is still in the `image`
package but has been renamed
[`image.Uniform`](/pkg/image/#Uniform),
while `image.Tiled` has been removed.

This table lists the renamings.

<table class="codetable" frame="border" summary="image renames">
<colgroup align="left" width="50%"></colgroup>
<colgroup align="left" width="50%"></colgroup>
<tbody><tr>
<th align="left">Old</th>
<th align="left">New</th>
</tr>
<tr>
<td colspan="2"><hr></hr></td>
</tr>
<tr><td>image.Color</td> <td>color.Color</td></tr>
<tr><td>image.ColorModel</td> <td>color.Model</td></tr>
<tr><td>image.ColorModelFunc</td> <td>color.ModelFunc</td></tr>
<tr><td>image.PalettedColorModel</td> <td>color.Palette</td></tr>
<tr>
<td colspan="2"><hr></hr></td>
</tr>
<tr><td>image.RGBAColor</td> <td>color.RGBA</td></tr>
<tr><td>image.RGBA64Color</td> <td>color.RGBA64</td></tr>
<tr><td>image.NRGBAColor</td> <td>color.NRGBA</td></tr>
<tr><td>image.NRGBA64Color</td> <td>color.NRGBA64</td></tr>
<tr><td>image.AlphaColor</td> <td>color.Alpha</td></tr>
<tr><td>image.Alpha16Color</td> <td>color.Alpha16</td></tr>
<tr><td>image.GrayColor</td> <td>color.Gray</td></tr>
<tr><td>image.Gray16Color</td> <td>color.Gray16</td></tr>
<tr>
<td colspan="2"><hr></hr></td>
</tr>
<tr><td>image.RGBAColorModel</td> <td>color.RGBAModel</td></tr>
<tr><td>image.RGBA64ColorModel</td> <td>color.RGBA64Model</td></tr>
<tr><td>image.NRGBAColorModel</td> <td>color.NRGBAModel</td></tr>
<tr><td>image.NRGBA64ColorModel</td> <td>color.NRGBA64Model</td></tr>
<tr><td>image.AlphaColorModel</td> <td>color.AlphaModel</td></tr>
<tr><td>image.Alpha16ColorModel</td> <td>color.Alpha16Model</td></tr>
<tr><td>image.GrayColorModel</td> <td>color.GrayModel</td></tr>
<tr><td>image.Gray16ColorModel</td> <td>color.Gray16Model</td></tr>
<tr>
<td colspan="2"><hr></hr></td>
</tr>
<tr><td>ycbcr.RGBToYCbCr</td> <td>color.RGBToYCbCr</td></tr>
<tr><td>ycbcr.YCbCrToRGB</td> <td>color.YCbCrToRGB</td></tr>
<tr><td>ycbcr.YCbCrColorModel</td> <td>color.YCbCrModel</td></tr>
<tr><td>ycbcr.YCbCrColor</td> <td>color.YCbCr</td></tr>
<tr><td>ycbcr.YCbCr</td> <td>image.YCbCr</td></tr>
<tr>
<td colspan="2"><hr></hr></td>
</tr>
<tr><td>ycbcr.SubsampleRatio444</td> <td>image.YCbCrSubsampleRatio444</td></tr>
<tr><td>ycbcr.SubsampleRatio422</td> <td>image.YCbCrSubsampleRatio422</td></tr>
<tr><td>ycbcr.SubsampleRatio420</td> <td>image.YCbCrSubsampleRatio420</td></tr>
<tr>
<td colspan="2"><hr></hr></td>
</tr>
<tr><td>image.ColorImage</td> <td>image.Uniform</td></tr>
</tbody></table>

The image package's `New` functions
([`NewRGBA`](/pkg/image/#NewRGBA),
[`NewRGBA64`](/pkg/image/#NewRGBA64), etc.)
take an [`image.Rectangle`](/pkg/image/#Rectangle) as an argument
instead of four integers.

Finally, there are new predefined `color.Color` variables
[`color.Black`](/pkg/image/color/#Black),
[`color.White`](/pkg/image/color/#White),
[`color.Opaque`](/pkg/image/color/#Opaque)
and
[`color.Transparent`](/pkg/image/color/#Transparent).

_Updating_:
Running `go` `fix` will update almost all code affected by the change.

### The log/syslog package {#log_syslog}

In Go 1, the [`syslog.NewLogger`](/pkg/log/syslog/#NewLogger)
function returns an error as well as a `log.Logger`.

_Updating_:
What little code is affected will be caught by the compiler and must be updated by hand.

### The mime package {#mime}

In Go 1, the [`FormatMediaType`](/pkg/mime/#FormatMediaType) function
of the `mime` package has been simplified to make it
consistent with
[`ParseMediaType`](/pkg/mime/#ParseMediaType).
It now takes `"text/html"` rather than `"text"` and `"html"`.

_Updating_:
What little code is affected will be caught by the compiler and must be updated by hand.

### The net package {#net}

In Go 1, the various `SetTimeout`,
`SetReadTimeout`, and `SetWriteTimeout` methods
have been replaced with
[`SetDeadline`](/pkg/net/#IPConn.SetDeadline),
[`SetReadDeadline`](/pkg/net/#IPConn.SetReadDeadline), and
[`SetWriteDeadline`](/pkg/net/#IPConn.SetWriteDeadline),
respectively. Rather than taking a timeout value in nanoseconds that
apply to any activity on the connection, the new methods set an
absolute deadline (as a `time.Time` value) after which
reads and writes will time out and no longer block.

There are also new functions
[`net.DialTimeout`](/pkg/net/#DialTimeout)
to simplify timing out dialing a network address and
[`net.ListenMulticastUDP`](/pkg/net/#ListenMulticastUDP)
to allow multicast UDP to listen concurrently across multiple listeners.
The `net.ListenMulticastUDP` function replaces the old
`JoinGroup` and `LeaveGroup` methods.

_Updating_:
Code that uses the old methods will fail to compile and must be updated by hand.
The semantic change makes it difficult for the fix tool to update automatically.

### The os package {#os}

The `Time` function has been removed; callers should use
the [`Time`](/pkg/time/#Time) type from the
`time` package.

The `Exec` function has been removed; callers should use
`Exec` from the `syscall` package, where available.

The `ShellExpand` function has been renamed to [`ExpandEnv`](/pkg/os/#ExpandEnv).

The [`NewFile`](/pkg/os/#NewFile) function
now takes a `uintptr` fd, instead of an `int`.
The [`Fd`](/pkg/os/#File.Fd) method on files now
also returns a `uintptr`.

There are no longer error constants such as `EINVAL`
in the `os` package, since the set of values varied with
the underlying operating system. There are new portable functions like
[`IsPermission`](/pkg/os/#IsPermission)
to test common error properties, plus a few new error values
with more Go-like names, such as
[`ErrPermission`](/pkg/os/#ErrPermission)
and
[`ErrNotExist`](/pkg/os/#ErrNotExist).

The `Getenverror` function has been removed. To distinguish
between a non-existent environment variable and an empty string,
use [`os.Environ`](/pkg/os/#Environ) or
[`syscall.Getenv`](/pkg/syscall/#Getenv).

The [`Process.Wait`](/pkg/os/#Process.Wait) method has
dropped its option argument and the associated constants are gone
from the package.
Also, the function `Wait` is gone; only the method of
the `Process` type persists.

The `Waitmsg` type returned by
[`Process.Wait`](/pkg/os/#Process.Wait)
has been replaced with a more portable
[`ProcessState`](/pkg/os/#ProcessState)
type with accessor methods to recover information about the
process.
Because of changes to `Wait`, the `ProcessState`
value always describes an exited process.
Portability concerns simplified the interface in other ways, but the values returned by the
[`ProcessState.Sys`](/pkg/os/#ProcessState.Sys) and
[`ProcessState.SysUsage`](/pkg/os/#ProcessState.SysUsage)
methods can be type-asserted to underlying system-specific data structures such as
[`syscall.WaitStatus`](/pkg/syscall/#WaitStatus) and
[`syscall.Rusage`](/pkg/syscall/#Rusage) on Unix.

_Updating_:
Running `go` `fix` will drop a zero argument to `Process.Wait`.
All other changes will be caught by the compiler and must be updated by hand.

#### The os.FileInfo type {#os_fileinfo}

Go 1 redefines the [`os.FileInfo`](/pkg/os/#FileInfo) type,
changing it from a struct to an interface:

	    type FileInfo interface {
	        Name() string       // base name of the file
	        Size() int64        // length in bytes
	        Mode() FileMode     // file mode bits
	        ModTime() time.Time // modification time
	        IsDir() bool        // abbreviation for Mode().IsDir()
	        Sys() interface{}   // underlying data source (can return nil)
	    }

The file mode information has been moved into a subtype called
[`os.FileMode`](/pkg/os/#FileMode),
a simple integer type with `IsDir`, `Perm`, and `String`
methods.

The system-specific details of file modes and properties such as (on Unix)
i-number have been removed from `FileInfo` altogether.
Instead, each operating system's `os` package provides an
implementation of the `FileInfo` interface, which
has a `Sys` method that returns the
system-specific representation of file metadata.
For instance, to discover the i-number of a file on a Unix system, unpack
the `FileInfo` like this:

	    fi, err := os.Stat("hello.go")
	    if err != nil {
	        log.Fatal(err)
	    }
	    // Check that it's a Unix file.
	    unixStat, ok := fi.Sys().(*syscall.Stat_t)
	    if !ok {
	        log.Fatal("hello.go: not a Unix file")
	    }
	    fmt.Printf("file i-number: %d\n", unixStat.Ino)

Assuming (which is unwise) that `"hello.go"` is a Unix file,
the i-number expression could be contracted to

	    fi.Sys().(*syscall.Stat_t).Ino

The vast majority of uses of `FileInfo` need only the methods
of the standard interface.

The `os` package no longer contains wrappers for the POSIX errors
such as `ENOENT`.
For the few programs that need to verify particular error conditions, there are
now the boolean functions
[`IsExist`](/pkg/os/#IsExist),
[`IsNotExist`](/pkg/os/#IsNotExist)
and
[`IsPermission`](/pkg/os/#IsPermission).

{{code "/doc/progs/go1.go" `/os\.Open/` `/}/`}}

_Updating_:
Running `go` `fix` will update code that uses the old equivalent of the current `os.FileInfo`
and `os.FileMode` API.
Code that needs system-specific file details will need to be updated by hand.
Code that uses the old POSIX error values from the `os` package
will fail to compile and will also need to be updated by hand.

### The os/signal package {#os_signal}

The `os/signal` package in Go 1 replaces the
`Incoming` function, which returned a channel
that received all incoming signals,
with the selective `Notify` function, which asks
for delivery of specific signals on an existing channel.

_Updating_:
Code must be updated by hand.
A literal translation of

	c := signal.Incoming()

is

	c := make(chan os.Signal, 1)
	signal.Notify(c) // ask for all signals

but most code should list the specific signals it wants to handle instead:

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT)

### The path/filepath package {#path_filepath}

In Go 1, the [`Walk`](/pkg/path/filepath/#Walk) function of the
`path/filepath` package
has been changed to take a function value of type
[`WalkFunc`](/pkg/path/filepath/#WalkFunc)
instead of a `Visitor` interface value.
`WalkFunc` unifies the handling of both files and directories.

	    type WalkFunc func(path string, info os.FileInfo, err error) error

The `WalkFunc` function will be called even for files or directories that could not be opened;
in such cases the error argument will describe the failure.
If a directory's contents are to be skipped,
the function should return the value [`filepath.SkipDir`](/pkg/path/filepath/#pkg-variables)

{{code "/doc/progs/go1.go" `/STARTWALK/` `/ENDWALK/`}}

_Updating_:
The change simplifies most code but has subtle consequences, so affected programs
will need to be updated by hand.
The compiler will catch code using the old interface.

### The regexp package {#regexp}

The [`regexp`](/pkg/regexp/) package has been rewritten.
It has the same interface but the specification of the regular expressions
it supports has changed from the old "egrep" form to that of
[RE2](https://code.google.com/p/re2/).

_Updating_:
Code that uses the package should have its regular expressions checked by hand.

### The runtime package {#runtime}

In Go 1, much of the API exported by package
`runtime` has been removed in favor of
functionality provided by other packages.
Code using the `runtime.Type` interface
or its specific concrete type implementations should
now use package [`reflect`](/pkg/reflect/).
Code using `runtime.Semacquire` or `runtime.Semrelease`
should use channels or the abstractions in package [`sync`](/pkg/sync/).
The `runtime.Alloc`, `runtime.Free`,
and `runtime.Lookup` functions, an unsafe API created for
debugging the memory allocator, have no replacement.

Before, `runtime.MemStats` was a global variable holding
statistics about memory allocation, and calls to `runtime.UpdateMemStats`
ensured that it was up to date.
In Go 1, `runtime.MemStats` is a struct type, and code should use
[`runtime.ReadMemStats`](/pkg/runtime/#ReadMemStats)
to obtain the current statistics.

The package adds a new function,
[`runtime.NumCPU`](/pkg/runtime/#NumCPU), that returns the number of CPUs available
for parallel execution, as reported by the operating system kernel.
Its value can inform the setting of `GOMAXPROCS`.
The `runtime.Cgocalls` and `runtime.Goroutines` functions
have been renamed to `runtime.NumCgoCall` and `runtime.NumGoroutine`.

_Updating_:
Running `go` `fix` will update code for the function renamings.
Other code will need to be updated by hand.

### The strconv package {#strconv}

In Go 1, the
[`strconv`](/pkg/strconv/)
package has been significantly reworked to make it more Go-like and less C-like,
although `Atoi` lives on (it's similar to
`int(ParseInt(x, 10, 0))`, as does
`Itoa(x)` (`FormatInt(int64(x), 10)`).
There are also new variants of some of the functions that append to byte slices rather than
return strings, to allow control over allocation.

This table summarizes the renamings; see the
[package documentation](/pkg/strconv/)
for full details.

<table class="codetable" frame="border" summary="strconv renames">
<colgroup align="left" width="50%"></colgroup>
<colgroup align="left" width="50%"></colgroup>
<tbody><tr>
<th align="left">Old call</th>
<th align="left">New call</th>
</tr>
<tr>
<td colspan="2"><hr></hr></td>
</tr>
<tr><td>Atob(x)</td> <td>ParseBool(x)</td></tr>
<tr>
<td colspan="2"><hr></hr></td>
</tr>
<tr><td>Atof32(x)</td> <td>ParseFloat(x, 32)§</td></tr>
<tr><td>Atof64(x)</td> <td>ParseFloat(x, 64)</td></tr>
<tr><td>AtofN(x, n)</td> <td>ParseFloat(x, n)</td></tr>
<tr>
<td colspan="2"><hr></hr></td>
</tr>
<tr><td>Atoi(x)</td> <td>Atoi(x)</td></tr>
<tr><td>Atoi(x)</td> <td>ParseInt(x, 10, 0)§</td></tr>
<tr><td>Atoi64(x)</td> <td>ParseInt(x, 10, 64)</td></tr>
<tr>
<td colspan="2"><hr></hr></td>
</tr>
<tr><td>Atoui(x)</td> <td>ParseUint(x, 10, 0)§</td></tr>
<tr><td>Atoui64(x)</td> <td>ParseUint(x, 10, 64)</td></tr>
<tr>
<td colspan="2"><hr></hr></td>
</tr>
<tr><td>Btoi64(x, b)</td> <td>ParseInt(x, b, 64)</td></tr>
<tr><td>Btoui64(x, b)</td> <td>ParseUint(x, b, 64)</td></tr>
<tr>
<td colspan="2"><hr></hr></td>
</tr>
<tr><td>Btoa(x)</td> <td>FormatBool(x)</td></tr>
<tr>
<td colspan="2"><hr></hr></td>
</tr>
<tr><td>Ftoa32(x, f, p)</td> <td>FormatFloat(float64(x), f, p, 32)</td></tr>
<tr><td>Ftoa64(x, f, p)</td> <td>FormatFloat(x, f, p, 64)</td></tr>
<tr><td>FtoaN(x, f, p, n)</td> <td>FormatFloat(x, f, p, n)</td></tr>
<tr>
<td colspan="2"><hr></hr></td>
</tr>
<tr><td>Itoa(x)</td> <td>Itoa(x)</td></tr>
<tr><td>Itoa(x)</td> <td>FormatInt(int64(x), 10)</td></tr>
<tr><td>Itoa64(x)</td> <td>FormatInt(x, 10)</td></tr>
<tr>
<td colspan="2"><hr></hr></td>
</tr>
<tr><td>Itob(x, b)</td> <td>FormatInt(int64(x), b)</td></tr>
<tr><td>Itob64(x, b)</td> <td>FormatInt(x, b)</td></tr>
<tr>
<td colspan="2"><hr></hr></td>
</tr>
<tr><td>Uitoa(x)</td> <td>FormatUint(uint64(x), 10)</td></tr>
<tr><td>Uitoa64(x)</td> <td>FormatUint(x, 10)</td></tr>
<tr>
<td colspan="2"><hr></hr></td>
</tr>
<tr><td>Uitob(x, b)</td> <td>FormatUint(uint64(x), b)</td></tr>
<tr><td>Uitob64(x, b)</td> <td>FormatUint(x, b)</td></tr>
</tbody></table>

_Updating_:
Running `go` `fix` will update almost all code affected by the change.
\
§ `Atoi` persists but `Atoui` and `Atof32` do not, so
they may require
a cast that must be added by hand; the `go` `fix` tool will warn about it.

### The template packages {#templates}

The `template` and `exp/template/html` packages have moved to
[`text/template`](/pkg/text/template/) and
[`html/template`](/pkg/html/template/).
More significant, the interface to these packages has been simplified.
The template language is the same, but the concept of "template set" is gone
and the functions and methods of the packages have changed accordingly,
often by elimination.

Instead of sets, a `Template` object
may contain multiple named template definitions,
in effect constructing
name spaces for template invocation.
A template can invoke any other template associated with it, but only those
templates associated with it.
The simplest way to associate templates is to parse them together, something
made easier with the new structure of the packages.

_Updating_:
The imports will be updated by fix tool.
Single-template uses will be otherwise be largely unaffected.
Code that uses multiple templates in concert will need to be updated by hand.
The [examples](/pkg/text/template/#pkg-examples) in
the documentation for `text/template` can provide guidance.

### The testing package {#testing}

The testing package has a type, `B`, passed as an argument to benchmark functions.
In Go 1, `B` has new methods, analogous to those of `T`, enabling
logging and failure reporting.

{{code "/doc/progs/go1.go" `/func.*Benchmark/` `/^}/`}}

_Updating_:
Existing code is unaffected, although benchmarks that use `println`
or `panic` should be updated to use the new methods.

### The testing/script package {#testing_script}

The testing/script package has been deleted. It was a dreg.

_Updating_:
No code is likely to be affected.

### The unsafe package {#unsafe}

In Go 1, the functions
`unsafe.Typeof`, `unsafe.Reflect`,
`unsafe.Unreflect`, `unsafe.New`, and
`unsafe.NewArray` have been removed;
they duplicated safer functionality provided by
package [`reflect`](/pkg/reflect/).

_Updating_:
Code using these functions must be rewritten to use
package [`reflect`](/pkg/reflect/).
The changes to [encoding/gob](/change/2646dc956207) and the [protocol buffer library](https://code.google.com/p/goprotobuf/source/detail?r=5340ad310031)
may be helpful as examples.

### The url package {#url}

In Go 1 several fields from the [`url.URL`](/pkg/net/url/#URL) type
were removed or replaced.

The [`String`](/pkg/net/url/#URL.String) method now
predictably rebuilds an encoded URL string using all of `URL`'s
fields as necessary. The resulting string will also no longer have
passwords escaped.

The `Raw` field has been removed. In most cases the `String`
method may be used in its place.

The old `RawUserinfo` field is replaced by the `User`
field, of type [`*net.Userinfo`](/pkg/net/url/#Userinfo).
Values of this type may be created using the new [`net.User`](/pkg/net/url/#User)
and [`net.UserPassword`](/pkg/net/url/#UserPassword)
functions. The `EscapeUserinfo` and `UnescapeUserinfo`
functions are also gone.

The `RawAuthority` field has been removed. The same information is
available in the `Host` and `User` fields.

The `RawPath` field and the `EncodedPath` method have
been removed. The path information in rooted URLs (with a slash following the
schema) is now available only in decoded form in the `Path` field.
Occasionally, the encoded data may be required to obtain information that
was lost in the decoding process. These cases must be handled by accessing
the data the URL was built from.

URLs with non-rooted paths, such as `"mailto:dev@golang.org?subject=Hi"`,
are also handled differently. The `OpaquePath` boolean field has been
removed and a new `Opaque` string field introduced to hold the encoded
path for such URLs. In Go 1, the cited URL parses as:

	    URL{
	        Scheme: "mailto",
	        Opaque: "dev@golang.org",
	        RawQuery: "subject=Hi",
	    }

A new [`RequestURI`](/pkg/net/url/#URL.RequestURI) method was
added to `URL`.

The `ParseWithReference` function has been renamed to `ParseWithFragment`.

_Updating_:
Code that uses the old fields will fail to compile and must be updated by hand.
The semantic changes make it difficult for the fix tool to update automatically.

## The go command {#cmd_go}

Go 1 introduces the [go command](/cmd/go/), a tool for fetching,
building, and installing Go packages and commands. The `go` command
does away with makefiles, instead using Go source code to find dependencies and
determine build conditions. Most existing Go programs will no longer require
makefiles to be built.

See [How to Write Go Code](/doc/code.html) for a primer on the
`go` command and the [go command documentation](/cmd/go/)
for the full details.

_Updating_:
Projects that depend on the Go project's old makefile-based build
infrastructure (`Make.pkg`, `Make.cmd`, and so on) should
switch to using the `go` command for building Go code and, if
necessary, rewrite their makefiles to perform any auxiliary build tasks.

## The cgo command {#cmd_cgo}

In Go 1, the [cgo command](/cmd/cgo)
uses a different `_cgo_export.h`
file, which is generated for packages containing `//export` lines.
The `_cgo_export.h` file now begins with the C preamble comment,
so that exported function definitions can use types defined there.
This has the effect of compiling the preamble multiple times, so a
package using `//export` must not put function definitions
or variable initializations in the C preamble.

## Packaged releases {#releases}

One of the most significant changes associated with Go 1 is the availability
of prepackaged, downloadable distributions.
They are available for many combinations of architecture and operating system
(including Windows) and the list will grow.
Installation details are described on the
[Getting Started](/doc/install) page, while
the distributions themselves are listed on the
[downloads page](/dl/).
