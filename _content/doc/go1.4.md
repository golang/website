---
template: false
title: Go 1.4 Release Notes
---

## Introduction to Go 1.4 {#introduction}

The latest Go release, version 1.4, arrives as scheduled six months after 1.3.

It contains only one tiny language change,
in the form of a backwards-compatible simple variant of `for`-`range` loop,
and a possibly breaking change to the compiler involving methods on pointers-to-pointers.

The release focuses primarily on implementation work, improving the garbage collector
and preparing the ground for a fully concurrent collector to be rolled out in the
next few releases.
Stacks are now contiguous, reallocated when necessary rather than linking on new
"segments";
this release therefore eliminates the notorious "hot stack split" problem.
There are some new tools available including support in the `go` command
for build-time source code generation.
The release also adds support for ARM processors on Android and Native Client (NaCl)
and for AMD64 on Plan 9.

As always, Go 1.4 keeps the [promise
of compatibility](/doc/go1compat.html),
and almost everything
will continue to compile and run without change when moved to 1.4.

## Changes to the language {#language}

### For-range loops {#forrange}

Up until Go 1.3, `for`-`range` loop had two forms

	for i, v := range x {
		...
	}

and

	for i := range x {
		...
	}

If one was not interested in the loop values, only the iteration itself, it was still
necessary to mention a variable (probably the [blank identifier](/ref/spec#Blank_identifier), as in
`for` `_` `=` `range` `x`), because
the form

	for range x {
		...
	}

was not syntactically permitted.

This situation seemed awkward, so as of Go 1.4 the variable-free form is now legal.
The pattern arises rarely but the code can be cleaner when it does.

_Updating_: The change is strictly backwards compatible to existing Go
programs, but tools that analyze Go parse trees may need to be modified to accept
this new form as the
`Key` field of [`RangeStmt`](/pkg/go/ast/#RangeStmt)
may now be `nil`.

### Method calls on \*\*T {#methodonpointertopointer}

Given these declarations,

	type T int
	func (T) M() {}
	var x **T

both `gc` and `gccgo` accepted the method call

	x.M()

which is a double dereference of the pointer-to-pointer `x`.
The Go specification allows a single dereference to be inserted automatically,
but not two, so this call is erroneous according to the language definition.
It has therefore been disallowed in Go 1.4, which is a breaking change,
although very few programs will be affected.

_Updating_: Code that depends on the old, erroneous behavior will no longer
compile but is easy to fix by adding an explicit dereference.

## Changes to the supported operating systems and architectures {#os}

### Android {#android}

Go 1.4 can build binaries for ARM processors running the Android operating system.
It can also build a `.so` library that can be loaded by an Android application
using the supporting packages in the [mobile](https://golang.org/x/mobile) subrepository.
A brief description of the plans for this experimental port are available
[here](/s/go14android).

### NaCl on ARM {#naclarm}

The previous release introduced Native Client (NaCl) support for the 32-bit x86
(`GOARCH=386`)
and 64-bit x86 using 32-bit pointers (GOARCH=amd64p32).
The 1.4 release adds NaCl support for ARM (GOARCH=arm).

### Plan9 on AMD64 {#plan9amd64}

This release adds support for the Plan 9 operating system on AMD64 processors,
provided the kernel supports the `nsec` system call and uses 4K pages.

## Changes to the compatibility guidelines {#compatibility}

The [`unsafe`](/pkg/unsafe/) package allows one
to defeat Go's type system by exploiting internal details of the implementation
or machine representation of data.
It was never explicitly specified what use of `unsafe` meant
with respect to compatibility as specified in the
[Go compatibility guidelines](go1compat.html).
The answer, of course, is that we can make no promise of compatibility
for code that does unsafe things.

We have clarified this situation in the documentation included in the release.
The [Go compatibility guidelines](go1compat.html) and the
docs for the [`unsafe`](/pkg/unsafe/) package
are now explicit that unsafe code is not guaranteed to remain compatible.

_Updating_: Nothing technical has changed; this is just a clarification
of the documentation.

## Changes to the implementations and tools {#impl}

### Changes to the runtime {#runtime}

Prior to Go 1.4, the runtime (garbage collector, concurrency support, interface management,
maps, slices, strings, ...) was mostly written in C, with some assembler support.
In 1.4, much of the code has been translated to Go so that the garbage collector can scan
the stacks of programs in the runtime and get accurate information about what variables
are active.
This change was large but should have no semantic effect on programs.

This rewrite allows the garbage collector in 1.4 to be fully precise,
meaning that it is aware of the location of all active pointers in the program.
This means the heap will be smaller as there will be no false positives keeping non-pointers alive.
Other related changes also reduce the heap size, which is smaller by 10%-30% overall
relative to the previous release.

A consequence is that stacks are no longer segmented, eliminating the "hot split" problem.
When a stack limit is reached, a new, larger stack is allocated, all active frames for
the goroutine are copied there, and any pointers into the stack are updated.
Performance can be noticeably better in some cases and is always more predictable.
Details are available in [the design document](/s/contigstacks).

The use of contiguous stacks means that stacks can start smaller without triggering performance issues,
so the default starting size for a goroutine's stack in 1.4 has been reduced from 8192 bytes to 2048 bytes.

As preparation for the concurrent garbage collector scheduled for the 1.5 release,
writes to pointer values in the heap are now done by a function call,
called a write barrier, rather than directly from the function updating the value.
In this next release, this will permit the garbage collector to mediate writes to the heap while it is running.
This change has no semantic effect on programs in 1.4, but was
included in the release to test the compiler and the resulting performance.

The implementation of interface values has been modified.
In earlier releases, the interface contained a word that was either a pointer or a one-word
scalar value, depending on the type of the concrete object stored.
This implementation was problematical for the garbage collector,
so as of 1.4 interface values always hold a pointer.
In running programs, most interface values were pointers anyway,
so the effect is minimal, but programs that store integers (for example) in
interfaces will see more allocations.

As of Go 1.3, the runtime crashes if it finds a memory word that should contain
a valid pointer but instead contains an obviously invalid pointer (for example, the value 3).
Programs that store integers in pointer values may run afoul of this check and crash.
In Go 1.4, setting the [`GODEBUG`](/pkg/runtime/) variable
`invalidptr=0` disables
the crash as a workaround, but we cannot guarantee that future releases will be
able to avoid the crash; the correct fix is to rewrite code not to alias integers and pointers.

### Assembly {#asm}

The language accepted by the assemblers `cmd/5a`, `cmd/6a`
and `cmd/8a` has had several changes,
mostly to make it easier to deliver type information to the runtime.

First, the `textflag.h` file that defines flags for `TEXT` directives
has been copied from the linker source directory to a standard location so it can be
included with the simple directive

	#include "textflag.h"

The more important changes are in how assembler source can define the necessary
type information.
For most programs it will suffice to move data
definitions (`DATA` and `GLOBL` directives)
out of assembly into Go files
and to write a Go declaration for each assembly function.
The [assembly document](/doc/asm#runtime) describes what to do.

_Updating_:
Assembly files that include `textflag.h` from its old
location will still work, but should be updated.
For the type information, most assembly routines will need no change,
but all should be examined.
Assembly source files that define data,
functions with non-empty stack frames, or functions that return pointers
need particular attention.
A description of the necessary (but simple) changes
is in the [assembly document](/doc/asm#runtime).

More information about these changes is in the [assembly document](/doc/asm).

### Status of gccgo {#gccgo}

The release schedules for the GCC and Go projects do not coincide.
GCC release 4.9 contains the Go 1.2 version of gccgo.
The next release, GCC 5, will likely have the Go 1.4 version of gccgo.

### Internal packages {#internalpackages}

Go's package system makes it easy to structure programs into components with clean boundaries,
but there are only two forms of access: local (unexported) and global (exported).
Sometimes one wishes to have components that are not exported,
for instance to avoid acquiring clients of interfaces to code that is part of a public repository
but not intended for use outside the program to which it belongs.

The Go language does not have the power to enforce this distinction, but as of Go 1.4 the
[`go`](/cmd/go/) command introduces
a mechanism to define "internal" packages that may not be imported by packages outside
the source subtree in which they reside.

To create such a package, place it in a directory named `internal` or in a subdirectory of a directory
named internal.
When the `go` command sees an import of a package with `internal` in its path,
it verifies that the package doing the import
is within the tree rooted at the parent of the `internal` directory.
For example, a package `.../a/b/c/internal/d/e/f`
can be imported only by code in the directory tree rooted at `.../a/b/c`.
It cannot be imported by code in `.../a/b/g` or in any other repository.

For Go 1.4, the internal package mechanism is enforced for the main Go repository;
from 1.5 and onward it will be enforced for any repository.

Full details of the mechanism are in
[the design document](/s/go14internal).

### Canonical import paths {#canonicalimports}

Code often lives in repositories hosted by public services such as `github.com`,
meaning that the import paths for packages begin with the name of the hosting service,
`github.com/rsc/pdf` for example.
One can use
[an existing mechanism](/cmd/go/#hdr-Remote_import_paths)
to provide a "custom" or "vanity" import path such as
`rsc.io/pdf`, but
that creates two valid import paths for the package.
That is a problem: one may inadvertently import the package through the two
distinct paths in a single program, which is wasteful;
miss an update to a package because the path being used is not recognized to be
out of date;
or break clients using the old path by moving the package to a different hosting service.

Go 1.4 introduces an annotation for package clauses in Go source that identify a canonical
import path for the package.
If an import is attempted using a path that is not canonical,
the [`go`](/cmd/go/) command
will refuse to compile the importing package.

The syntax is simple: put an identifying comment on the package line.
For our example, the package clause would read:

	package pdf // import "rsc.io/pdf"

With this in place,
the `go` command will
refuse to compile a package that imports `github.com/rsc/pdf`,
ensuring that the code can be moved without breaking users.

The check is at build time, not download time, so if `go` `get`
fails because of this check, the mis-imported package has been copied to the local machine
and should be removed manually.

To complement this new feature, a check has been added at update time to verify
that the local package's remote repository matches that of its custom import.
The `go` `get` `-u` command will fail to
update a package if its remote repository has changed since it was first
downloaded.
The new `-f` flag overrides this check.

Further information is in
[the design document](/s/go14customimport).

### Import paths for the subrepositories {#subrepo}

The Go project subrepositories (`code.google.com/p/go.tools` and so on)
are now available under custom import paths replacing `code.google.com/p/go.` with `golang.org/x/`,
as in `golang.org/x/tools`.
We will add canonical import comments to the code around June 1, 2015,
at which point Go 1.4 and later will stop accepting the old `code.google.com` paths.

_Updating_: All code that imports from subrepositories should change
to use the new `golang.org` paths.
Go 1.0 and later can resolve and import the new paths, so updating will not break
compatibility with older releases.
Code that has not updated will stop compiling with Go 1.4 around June 1, 2015.

### The go generate subcommand {#gogenerate}

The [`go`](/cmd/go/) command has a new subcommand,
[`go generate`](/cmd/go/#hdr-Generate_Go_files_by_processing_source),
to automate the running of tools to generate source code before compilation.
For example, it can be used to run the [`yacc`](/cmd/yacc)
compiler-compiler on a `.y` file to produce the Go source file implementing the grammar,
or to automate the generation of `String` methods for typed constants using the new
[stringer](https://godoc.org/golang.org/x/tools/cmd/stringer)
tool in the `golang.org/x/tools` subrepository.

For more information, see the
[design document](/s/go1.4-generate).

### Change to file name handling {#filenames}

Build constraints, also known as build tags, control compilation by including or excluding files
(see the documentation [`/go/build`](/pkg/go/build/)).
Compilation can also be controlled by the name of the file itself by "tagging" the file with
a suffix (before the `.go` or `.s` extension) with an underscore
and the name of the architecture or operating system.
For instance, the file `gopher_arm.go` will only be compiled if the target
processor is an ARM.

Before Go 1.4, a file called just `arm.go` was similarly tagged, but this behavior
can break sources when new architectures are added, causing files to suddenly become tagged.
In 1.4, therefore, a file will be tagged in this manner only if the tag (architecture or operating
system name) is preceded by an underscore.

_Updating_: Packages that depend on the old behavior will no longer compile correctly.
Files with names like `windows.go` or `amd64.go` should either
have explicit build tags added to the source or be renamed to something like
`os_windows.go` or `support_amd64.go`.

### Other changes to the go command {#gocmd}

There were a number of minor changes to the
[`cmd/go`](/cmd/go/)
command worth noting.

  - Unless [`cgo`](/cmd/cgo/) is being used to build the package,
    the `go` command now refuses to compile C source files,
    since the relevant C compilers
    ([`6c`](/cmd/6c/) etc.)
    are intended to be removed from the installation in some future release.
    (They are used today only to build part of the runtime.)
    It is difficult to use them correctly in any case, so any extant uses are likely incorrect,
    so we have disabled them.
  - The [`go` `test`](/cmd/go/#hdr-Test_packages)
    subcommand has a new flag, `-o`, to set the name of the resulting binary,
    corresponding to the same flag in other subcommands.
    The non-functional `-file` flag has been removed.
  - The [`go` `test`](/cmd/go/#hdr-Test_packages)
    subcommand will compile and link all `*_test.go` files in the package,
    even when there are no `Test` functions in them.
    It previously ignored such files.
  - The behavior of the
    [`go` `build`](/cmd/go/#hdr-Test_packages)
    subcommand's
    `-a` flag has been changed for non-development installations.
    For installations running a released distribution, the `-a` flag will no longer
    rebuild the standard library and commands, to avoid overwriting the installation's files.

### Changes to package source layout {#pkg}

In the main Go source repository, the source code for the packages was kept in
the directory `src/pkg`, which made sense but differed from
other repositories, including the Go subrepositories.
In Go 1.4, the`  pkg ` level of the source tree is now gone, so for example
the [`fmt`](/pkg/fmt/) package's source, once kept in
directory `src/pkg/fmt`, now lives one level higher in `src/fmt`.

_Updating_: Tools like `godoc` that discover source code
need to know about the new location. All tools and services maintained by the Go team
have been updated.

### SWIG {#swig}

Due to runtime changes in this release, Go 1.4 requires SWIG 3.0.3.

### Miscellany {#misc}

The standard repository's top-level `misc` directory used to contain
Go support for editors and IDEs: plugins, initialization scripts and so on.
Maintaining these was becoming time-consuming
and needed external help because many of the editors listed were not used by
members of the core team.
It also required us to make decisions about which plugin was best for a given
editor, even for editors we do not use.

The Go community at large is much better suited to managing this information.
In Go 1.4, therefore, this support has been removed from the repository.
Instead, there is a curated, informative list of what's available on
a [wiki page](/wiki/IDEsAndTextEditorPlugins).

## Performance {#performance}

Most programs will run about the same speed or slightly faster in 1.4 than in 1.3;
some will be slightly slower.
There are many changes, making it hard to be precise about what to expect.

As mentioned above, much of the runtime was translated to Go from C,
which led to some reduction in heap sizes.
It also improved performance slightly because the Go compiler is better
at optimization, due to things like inlining, than the C compiler used to build
the runtime.

The garbage collector was sped up, leading to measurable improvements for
garbage-heavy programs.
On the other hand, the new write barriers slow things down again, typically
by about the same amount but, depending on their behavior, some programs
may be somewhat slower or faster.

Library changes that affect performance are documented below.

## Changes to the standard library {#library}

### New packages {#new_packages}

There are no new packages in this release.

### Major changes to the library {#major_library_changes}

#### bufio.Scanner {#scanner}

The [`Scanner`](/pkg/bufio/#Scanner) type in the
[`bufio`](/pkg/bufio/) package
has had a bug fixed that may require changes to custom
[`split functions`](/pkg/bufio/#SplitFunc).
The bug made it impossible to generate an empty token at EOF; the fix
changes the end conditions seen by the split function.
Previously, scanning stopped at EOF if there was no more data.
As of 1.4, the split function will be called once at EOF after input is exhausted,
so the split function can generate a final empty token
as the documentation already promised.

_Updating_: Custom split functions may need to be modified to
handle empty tokens at EOF as desired.

#### syscall {#syscall}

The [`syscall`](/pkg/syscall/) package is now frozen except
for changes needed to maintain the core repository.
In particular, it will no longer be extended to support new or different system calls
that are not used by the core.
The reasons are described at length in [a
separate document](/s/go1.4-syscall).

A new subrepository, [golang.org/x/sys](https://golang.org/x/sys),
has been created to serve as the location for new developments to support system
calls on all kernels.
It has a nicer structure, with three packages that each hold the implementation of
system calls for one of
[Unix](https://godoc.org/golang.org/x/sys/unix),
[Windows](https://godoc.org/golang.org/x/sys/windows) and
[Plan 9](https://godoc.org/golang.org/x/sys/plan9).
These packages will be curated more generously, accepting all reasonable changes
that reflect kernel interfaces in those operating systems.
See the documentation and the article mentioned above for more information.

_Updating_: Existing programs are not affected as the `syscall`
package is largely unchanged from the 1.3 release.
Future development that requires system calls not in the `syscall` package
should build on `golang.org/x/sys` instead.

### Minor changes to the library {#minor_library_changes}

The following list summarizes a number of minor changes to the library, mostly additions.
See the relevant package documentation for more information about each change.

  - The [`archive/zip`](/pkg/archive/zip/) package's
    [`Writer`](/pkg/archive/zip/#Writer) now supports a
    [`Flush`](/pkg/archive/zip/#Writer.Flush) method.
  - The [`compress/flate`](/pkg/compress/flate/),
    [`compress/gzip`](/pkg/compress/gzip/),
    and [`compress/zlib`](/pkg/compress/zlib/)
    packages now support a `Reset` method
    for the decompressors, allowing them to reuse buffers and improve performance.
    The [`compress/gzip`](/pkg/compress/gzip/) package also has a
    [`Multistream`](/pkg/compress/gzip/#Reader.Multistream) method to control support
    for multistream files.
  - The [`crypto`](/pkg/crypto/) package now has a
    [`Signer`](/pkg/crypto/#Signer) interface, implemented by the
    `PrivateKey` types in
    [`crypto/ecdsa`](/pkg/crypto/ecdsa) and
    [`crypto/rsa`](/pkg/crypto/rsa).
  - The [`crypto/tls`](/pkg/crypto/tls/) package
    now supports ALPN as defined in [RFC 7301](https://tools.ietf.org/html/rfc7301).
  - The [`crypto/tls`](/pkg/crypto/tls/) package
    now supports programmatic selection of server certificates
    through the new [`CertificateForName`](/pkg/crypto/tls/#Config.CertificateForName) function
    of the [`Config`](/pkg/crypto/tls/#Config) struct.
  - Also in the crypto/tls package, the server now supports
    [TLS\_FALLBACK\_SCSV](https://tools.ietf.org/html/draft-ietf-tls-downgrade-scsv-00)
    to help clients detect fallback attacks.
    (The Go client does not support fallback at all, so it is not vulnerable to
    those attacks.)
  - The [`database/sql`](/pkg/database/sql/) package can now list all registered
    [`Drivers`](/pkg/database/sql/#Drivers).
  - The [`debug/dwarf`](/pkg/debug/dwarf/) package now supports
    [`UnspecifiedType`](/pkg/debug/dwarf/#UnspecifiedType)s.
  - In the [`encoding/asn1`](/pkg/encoding/asn1/) package,
    optional elements with a default value will now only be omitted if they have that value.
  - The [`encoding/csv`](/pkg/encoding/csv/) package no longer
    quotes empty strings but does quote the end-of-data marker `\.` (backslash dot).
    This is permitted by the definition of CSV and allows it to work better with Postgres.
  - The [`encoding/gob`](/pkg/encoding/gob/) package has been rewritten to eliminate
    the use of unsafe operations, allowing it to be used in environments that do not permit use of the
    [`unsafe`](/pkg/unsafe/) package.
    For typical uses it will be 10-30% slower, but the delta is dependent on the type of the data and
    in some cases, especially involving arrays, it can be faster.
    There is no functional change.
  - The [`encoding/xml`](/pkg/encoding/xml/) package's
    [`Decoder`](/pkg/encoding/xml/#Decoder) can now report its input offset.
  - In the [`fmt`](/pkg/fmt/) package,
    formatting of pointers to maps has changed to be consistent with that of pointers
    to structs, arrays, and so on.
    For instance, `&map[string]int{"one":` `1}` now prints by default as
    `&map[one:` `1]` rather than as a hexadecimal pointer value.
  - The [`image`](/pkg/image/) package's
    [`Image`](/pkg/image/#Image)
    implementations like
    [`RGBA`](/pkg/image/#RGBA) and
    [`Gray`](/pkg/image/#Gray) have specialized
    [`RGBAAt`](/pkg/image/#RGBA.RGBAAt) and
    [`GrayAt`](/pkg/image/#Gray.GrayAt) methods alongside the general
    [`At`](/pkg/image/#Image.At) method.
  - The [`image/png`](/pkg/image/png/) package now has an
    [`Encoder`](/pkg/image/png/#Encoder)
    type to control the compression level used for encoding.
  - The [`math`](/pkg/math/) package now has a
    [`Nextafter32`](/pkg/math/#Nextafter32) function.
  - The [`net/http`](/pkg/net/http/) package's
    [`Request`](/pkg/net/http/#Request) type
    has a new [`BasicAuth`](/pkg/net/http/#Request.BasicAuth) method
    that returns the username and password from authenticated requests using the
    HTTP Basic Authentication
    Scheme.
  - The [`net/http`](/pkg/net/http/) package's
    [`Transport`](/pkg/net/http/#Request) type
    has a new [`DialTLS`](/pkg/net/http/#Transport.DialTLS) hook
    that allows customizing the behavior of outbound TLS connections.
  - The [`net/http/httputil`](/pkg/net/http/httputil/) package's
    [`ReverseProxy`](/pkg/net/http/httputil/#ReverseProxy) type
    has a new field,
    [`ErrorLog`](/pkg/net/http/#ReverseProxy.ErrorLog), that
    provides user control of logging.
  - The [`os`](/pkg/os/) package
    now implements symbolic links on the Windows operating system
    through the [`Symlink`](/pkg/os/#Symlink) function.
    Other operating systems already have this functionality.
    There is also a new [`Unsetenv`](/pkg/os/#Unsetenv) function.
  - The [`reflect`](/pkg/reflect/) package's
    [`Type`](/pkg/reflect/#Type) interface
    has a new method, [`Comparable`](/pkg/reflect/#type.Comparable),
    that reports whether the type implements general comparisons.
  - Also in the [`reflect`](/pkg/reflect/) package, the
    [`Value`](/pkg/reflect/#Value) interface is now three instead of four words
    because of changes to the implementation of interfaces in the runtime.
    This saves memory but has no semantic effect.
  - The [`runtime`](/pkg/runtime/) package
    now implements monotonic clocks on Windows,
    as it already did for the other systems.
  - The [`runtime`](/pkg/runtime/) package's
    [`Mallocs`](/pkg/runtime/#MemStats.Mallocs) counter
    now counts very small allocations that were missed in Go 1.3.
    This may break tests using [`ReadMemStats`](/pkg/runtime/#ReadMemStats)
    or [`AllocsPerRun`](/pkg/testing/#AllocsPerRun)
    due to the more accurate answer.
  - In the [`runtime`](/pkg/runtime/) package,
    an array [`PauseEnd`](/pkg/runtime/#MemStats.PauseEnd)
    has been added to the
    [`MemStats`](/pkg/runtime/#MemStats)
    and [`GCStats`](/pkg/runtime/#GCStats) structs.
    This array is a circular buffer of times when garbage collection pauses ended.
    The corresponding pause durations are already recorded in
    [`PauseNs`](/pkg/runtime/#MemStats.PauseNs)
  - The [`runtime/race`](/pkg/runtime/race/) package
    now supports FreeBSD, which means the
    [`go`](/pkg/cmd/go/) command's `-race`
    flag now works on FreeBSD.
  - The [`sync/atomic`](/pkg/sync/atomic/) package
    has a new type, [`Value`](/pkg/sync/atomic/#Value).
    `Value` provides an efficient mechanism for atomic loads and
    stores of values of arbitrary type.
  - In the [`syscall`](/pkg/syscall/) package's
    implementation on Linux, the
    [`Setuid`](/pkg/syscall/#Setuid)
    and [`Setgid`](/pkg/syscall/#Setgid) have been disabled
    because those system calls operate on the calling thread, not the whole process, which is
    different from other platforms and not the expected result.
  - The [`testing`](/pkg/testing/) package
    has a new facility to provide more control over running a set of tests.
    If the test code contains a function
    <pre>
    func TestMain(m *<a href="/pkg/testing/#M"><code>testing.M</code></a>)
    </pre>
    that function will be called instead of running the tests directly.
    The `M` struct contains methods to access and run the tests.
  - Also in the [`testing`](/pkg/testing/) package,
    a new [`Coverage`](/pkg/testing/#Coverage)
    function reports the current test coverage fraction,
    enabling individual tests to report how much they are contributing to the
    overall coverage.
  - The [`text/scanner`](/pkg/text/scanner/) package's
    [`Scanner`](/pkg/text/scanner/#Scanner) type
    has a new function,
    [`IsIdentRune`](/pkg/text/scanner/#Scanner.IsIdentRune),
    allowing one to control the definition of an identifier when scanning.
  - The [`text/template`](/pkg/text/template/) package's boolean
    functions `eq`, `lt`, and so on have been generalized to allow comparison
    of signed and unsigned integers, simplifying their use in practice.
    (Previously one could only compare values of the same signedness.)
    All negative values compare less than all unsigned values.
  - The `time` package now uses the standard symbol for the micro prefix,
    the micro symbol (U+00B5 'µ'), to print microsecond durations.
    [`ParseDuration`](/pkg/time/#ParseDuration) still accepts `us`
    but the package no longer prints microseconds as `us`.
    \
    _Updating_: Code that depends on the output format of durations
    but does not use ParseDuration will need to be updated.
