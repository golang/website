---
path: /doc/go1.22
title: Go 1.22 Release Notes
---

<!--
NOTE: In this document and others in this directory, the convention is to
set fixed-width phrases with non-fixed-width spaces, as in
<code>hello</code> <code>world</code>.
Do not send CLs removing the interior tags from such phrases.
-->

<style>
  main ul li { margin: 0.5em 0; }
</style>

## Introduction to Go 1.22 {#introduction}

The latest Go release, version 1.22, arrives six months after [Go 1.21](/doc/go1.21).
Most of its changes are in the implementation of the toolchain, runtime, and libraries.
As always, the release maintains the Go 1 [promise of compatibility](/doc/go1compat).
We expect almost all Go programs to continue to compile and run as before.

## Changes to the language {#language}

<!-- loop variable scope -->
<!-- range over int -->
Go 1.22 makes two changes to "for" loops.

  - Previously, the variables declared by a "for" loop were created once and updated by each iteration. In Go 1.22, each iteration of the loop creates new variables, to avoid accidental sharing bugs.
    The [transition support tooling](/wiki/LoopvarExperiment#my-test-fails-with-the-change-how-can-i-debug-it)
    described in the proposal continues to work in the same way it did in Go 1.21.
  - "For" loops may now range over integers.
    For [example](/play/p/ky02zZxgk_r?v=gotip):

    	package main

    	import "fmt"

    	func main() {
    	  for i := range 10 {
    	    fmt.Println(10 - i)
    	  }
    	  fmt.Println("go1.22 has lift-off!")
    	}

    See the spec for [details](/ref/spec#For_range).

<!-- range over func GOEXPERIMENT; https://go.dev/issue/61405, https://go.dev/issue/61897, CLs 510541,539277,540263,543319 -->


Go 1.22 includes a preview of a language change we are considering
for a future version of Go: [range-over-function iterators](/wiki/RangefuncExperiment).
Building with `GOEXPERIMENT=rangefunc` enables this feature.

## Tools {#tools}

### Go command {#go-command}

<!-- https://go.dev/issue/60056 -->

Commands in [workspaces](/ref/mod#workspaces) can now
use a `vendor` directory containing the dependencies of the
workspace. The directory is created by
[`go` `work` `vendor`](/pkg/cmd/go#hdr-Make_vendored_copy_of_dependencies),
and used by build commands when the `-mod` flag is set to
`vendor`, which is the default when a workspace `vendor`
directory is present.

Note that the `vendor` directory's contents for a workspace are different
from those of a single module: if the directory at the root of a workspace also
contains one of the modules in the workspace, its `vendor` directory
can contain the dependencies of either the workspace or of the module,
but not both.

<!-- CL 518775, https://go.dev/issue/60915 -->

`go` `get` is no longer supported outside of a module in the
legacy `GOPATH` mode (that is, with `GO111MODULE=off`).
Other build commands, such as `go` `build` and
`go` `test`, will continue to work indefinitely
for legacy `GOPATH` programs.

<!-- CL 518776 -->

`go` `mod` `init` no longer attempts to import
module requirements from configuration files for other vendoring tools
(such as `Gopkg.lock`).

<!-- CL 495447 -->

`go` `test` `-cover` now prints coverage summaries
for covered packages that do not have their own test files. Prior to Go 1.22 a
`go` `test` `-cover` run for such a package would
report

`?     mymod/mypack    [no test files]`

and now with Go 1.22, functions in the package are treated as uncovered:

`mymod/mypack    coverage: 0.0% of statements`

Note that if a package contains no executable code at all, we can't report
a meaningful coverage percentage; for such packages the `go` tool
will continue to report that there are no test files.

<!-- CL 522239, https://go.dev/issue/46330 -->

`go` build commands that invoke the linker now error out if an
external (C) linker will be used but cgo is not enabled. (The Go runtime
requires cgo support to ensure that it is compatible with any additional
libraries added by the C linker.)

### Trace {#trace}

<!-- https://go.dev/issue/63960 -->

The `trace` tool's web UI has been gently refreshed as part of the
work to support the new tracer, resolving several issues and improving the
readability of various sub-pages.
The web UI now supports exploring traces in a thread-oriented view.
The trace viewer also now displays the full duration of all system calls.
\
These improvements only apply for viewing traces produced by programs built with
Go 1.22 or newer.
A future release will bring some of these improvements to traces produced by older
version of Go.

### Vet {#vet}

#### References to loop variables {#vet-loopclosure}

<!-- CL 539016, https://go.dev/issue/63888: cmd/vet: do not report variable capture for loop variables with the new lifetime rules -->
The behavior of the `vet` tool has changed to match
the new semantics (see above) of loop variables in Go 1.22.
When analyzing a file that requires Go 1.22 or newer
(due to its go.mod file or a per-file build constraint),
`vet` no longer reports references to
loop variables from within a function literal that
might outlive the iteration of the loop.
In Go 1.22, loop variables are created anew for each iteration,
so such references are no longer at risk of using a variable
after it has been updated by the loop.

#### New warnings for missing values after append {#vet-appends}

<!-- CL 498416, https://go.dev/issue/60448: add a new analyzer for check missing values after append -->
The `vet` tool now reports calls to
[`append`](/pkg/builtin/#append) that pass
no values to be appended to the slice, such as `slice = append(slice)`.
Such a statement has no effect, and experience has shown that is nearly always a mistake.

#### New warnings for deferring `time.Since` {#vet-defers}

<!-- CL 527095, https://go.dev/issue/60048: time.Since should not be used in defer statement -->
The vet tool now reports a non-deferred call to
[`time.Since(t)`](/pkg/time/#Since) within a `defer` statement.
This is equivalent to calling `time.Now().Sub(t)` before the `defer` statement,
not when the deferred function is called. In nearly all cases, the correct code
requires deferring the `time.Since` call. For example:

	t := time.Now()
	defer log.Println(time.Since(t)) // non-deferred call to time.Since
	tmp := time.Since(t); defer log.Println(tmp) // equivalent to the previous defer

	defer func() {
	  log.Println(time.Since(t)) // a correctly deferred call to time.Since
	}()

#### New warnings for mismatched key-value pairs in `log/slog` calls {#vet-slog}

<!-- CL 496156, https://go.dev/issue/59407: log/slog: add vet checks for variadic ...any inputs -->
The vet tool now reports invalid arguments in calls to functions and methods
in the structured logging package, [`log/slog`](/pkg/log/slog),
that accept alternating key/value pairs.
It reports calls where an argument in a key position is neither a
`string` nor a `slog.Attr`, and where a final key is missing its value.

## Runtime {#runtime}

<!-- CL 543255 -->
The runtime now keeps type-based garbage collection metadata nearer to each
heap object, improving the CPU performance (latency or throughput) of Go programs
by 1–3%.
This change also reduces the memory overhead of the majority Go programs by
approximately 1% by deduplicating redundant metadata.
Some programs may see a smaller improvement because this change adjusts the size
class boundaries of the memory allocator, so some objects may be moved up a size
class.

A consequence of this change is that some objects' addresses that were previously
always aligned to a 16 byte (or higher) boundary will now only be aligned to an 8
byte boundary.
Some programs that use assembly instructions that require memory addresses to be
more than 8-byte aligned and rely on the memory allocator's previous alignment behavior
may break, but we expect such programs to be rare.
Such programs may be built with `GOEXPERIMENT=noallocheaders` to revert
to the old metadata layout and restore the previous alignment behavior, but package
owners should update their assembly code to avoid the alignment assumption, as this
workaround will be removed in a future release.

<!-- CL 525475 -->
On the `windows/amd64 port`, programs linking or loading Go libraries built with
`-buildmode=c-archive` or `-buildmode=c-shared` can now use
the `SetUnhandledExceptionFilter` Win32 function to catch exceptions not handled
by the Go runtime. Note that this was already supported on the `windows/386` port.

## Compiler {#compiler}

<!-- https://go.dev/issue/61577 -->
[Profile-guided Optimization (PGO)](/doc/pgo) builds
can now devirtualize a higher proportion of calls than previously possible.
Most programs from a representative set of Go programs now see between 2 and
14% improvement at runtime from enabling PGO.

<!-- https://go.dev/cl/528321 -->
The compiler now interleaves devirtualization and inlining, so interface
method calls are better optimized.

<!-- https://go.dev/issue/61502 -->
Go 1.22 also includes a preview of an enhanced implementation of the compiler's inlining phase that uses heuristics to boost inlinability at call sites deemed "important" (for example, in loops) and discourage inlining at call sites deemed "unimportant" (for example, on panic paths).
Building with `GOEXPERIMENT=newinliner` enables the new call-site
heuristics; see [issue #61502](/issue/61502) for
more info and to provide feedback.

## Linker {#linker}

<!-- CL 493136 -->
The linker's `-s` and `-w` flags are now behave more
consistently across all platforms.
The `-w` flag suppresses DWARF debug information generation.
The `-s` flag suppresses symbol table generation.
The `-s` flag also implies the `-w` flag, which can be
negated with `-w=0`.
That is, `-s` `-w=0` will generate a binary with DWARF
debug information generation but without the symbol table.

<!-- CL 511475 -->
On ELF platforms, the `-B` linker flag now accepts a special form:
with `-B` `gobuildid`, the linker will generate a GNU
build ID (the ELF `NT_GNU_BUILD_ID` note) derived from the Go
build ID.

<!-- CL 534555 -->
On Windows, when building with `-linkmode=internal`, the linker now
preserves SEH information from C object files by copying the `.pdata`
and `.xdata` sections into the final binary.
This helps with debugging and profiling binaries using native tools, such as WinDbg.
Note that until now, C functions' SEH exception handlers were not being honored,
so this change may cause some programs to behave differently.
`-linkmode=external` is not affected by this change, as external linkers
already preserve SEH information.

## Bootstrap {#bootstrap}

As mentioned in the [Go 1.20 release notes](/doc/go1.20#bootstrap), Go 1.22 now requires
the final point release of Go 1.20 or later for bootstrap.
We expect that Go 1.24 will require the final point release of Go 1.22 or later for bootstrap.

## Standard library {#library}

### New math/rand/v2 package {#math_rand_v2}

<!-- CL 502495 -->
<!-- CL 502497 -->
<!-- CL 502498 -->
<!-- CL 502499 -->
<!-- CL 502500 -->
<!-- CL 502505 -->
<!-- CL 502506 -->
<!-- CL 516857 -->
<!-- CL 516859 -->

Go 1.22 includes the first “v2” package in the standard library,
[`math/rand/v2`](/pkg/math/rand/v2/).
The changes compared to [`math/rand`](/pkg/math/rand/) are
detailed in [proposal #61716](/issue/61716). The most important changes are:

  - The `Read` method, deprecated in `math/rand`,
    was not carried forward for `math/rand/v2`.
    (It remains available in `math/rand`.)
    The vast majority of calls to `Read` should use
    [`crypto/rand`’s `Read`](/pkg/crypto/rand/#Read) instead.
    Otherwise a custom `Read` can be constructed using the `Uint64` method.
  - The global generator accessed by top-level functions is unconditionally randomly seeded.
    Because the API guarantees no fixed sequence of results,
    optimizations like per-thread random generator states are now possible.
  - The [`Source`](/pkg/math/rand/v2/#Source)
    interface now has a single `Uint64` method;
    there is no `Source64` interface.
  - Many methods now use faster algorithms that were not possible to adopt in `math/rand`
    because they changed the output streams.
  - The
    `Intn`,
    `Int31`,
    `Int31n`,
    `Int63`,
    and
    `Int64n`
    top-level functions and methods from `math/rand`
    are spelled more idiomatically in `math/rand/v2`:
    `IntN`,
    `Int32`,
    `Int32N`,
    `Int64`,
    and
    `Int64N`.
    There are also new top-level functions and methods
    `Uint32`,
    `Uint32N`,
    `Uint64`,
    `Uint64N`,
    and
    `UintN`.
  - The
    new generic function [`N`](/pkg/math/rand/v2/#N)
    is like
    [`Int64N`](/pkg/math/rand/v2/#Int64N) or
    [`Uint64N`](/pkg/math/rand/v2/#Uint64N)
    but works for any integer type.
    For example a random duration from 0 up to 5 minutes is
    `rand.N(5*time.Minute)`.
  - The Mitchell & Reeds LFSR generator provided by
    [`math/rand`’s `Source`](/pkg/math/rand/#Source)
    has been replaced by two more modern pseudo-random generator sources:
    [`ChaCha8`](/pkg/math/rand/v2/#ChaCha8) and
    [`PCG`](/pkg/math/rand/v2/#PCG).
    ChaCha8 is a new, cryptographically strong random number generator
    roughly similar to PCG in efficiency.
    ChaCha8 is the algorithm used for the top-level functions in `math/rand/v2`.
    As of Go 1.22, `math/rand`'s top-level functions (when not explicitly seeded)
    and the Go runtime also use ChaCha8 for randomness.

We plan to include an API migration tool in a future release, likely Go 1.23.

### New go/version package {#go-version}

<!-- https://go.dev/issue/62039, https://go.dev/cl/538895 -->
The new [`go/version`](/pkg/go/version/) package implements functions
for validating and comparing Go version strings.

### Enhanced routing patterns {#enhanced_routing_patterns}

<!-- https://go.dev/issue/61410 -->
HTTP routing in the standard library is now more expressive.
The patterns used by [`net/http.ServeMux`](/pkg/net/http#ServeMux) have been enhanced to accept methods and wildcards.

Registering a handler with a method, like `"POST /items/create"`, restricts
invocations of the handler to requests with the given method. A pattern with a method takes precedence over a matching pattern without one.
As a special case, registering a handler with `  "GET" ` also registers it with `"HEAD"`.

Wildcards in patterns, like `/items/{id}`, match segments of the URL path.
The actual segment value may be accessed by calling the [`Request.PathValue`](/pkg/net/http#Request.PathValue) method.
A wildcard ending in "...", like `/files/{path...}`, must occur at the end of a pattern and matches all the remaining segments.

A pattern that ends in "/" matches all paths that have it as a prefix, as always.
To match the exact pattern including the trailing slash, end it with `{$}`,
as in `/exact/match/{$}`.

If two patterns overlap in the requests that they match, then the more specific pattern takes precedence.
If neither is more specific, the patterns conflict.
This rule generalizes the original precedence rules and maintains the property that the order in which
patterns are registered does not matter.

This change breaks backwards compatibility in small ways, some obvious—patterns with "{" and "}" behave differently—
and some less so—treatment of escaped paths has been improved.
The change is controlled by a [`GODEBUG`](/doc/godebug) field named `httpmuxgo121`.
Set `httpmuxgo121=1` to restore the old behavior.

### Minor changes to the library {#minor_library_changes}

As always, there are various minor changes and updates to the library,
made with the Go 1 [promise of compatibility](/doc/go1compat)
in mind.
There are also various performance improvements, not enumerated here.

[archive/tar](/pkg/archive/tar/)

:   <!-- https://go.dev/issue/58000, CL 513316 -->
    The new method [`Writer.AddFS`](/pkg/archive/tar#Writer.AddFS) adds all of the files from an [`fs.FS`](/pkg/io/fs#FS) to the archive.

<!-- archive/tar -->

[archive/zip](/pkg/archive/zip/)

:   <!-- https://go.dev/issue/54898, CL 513438 -->
    The new method [`Writer.AddFS`](/pkg/archive/zip#Writer.AddFS) adds all of the files from an [`fs.FS`](/pkg/io/fs#FS) to the archive.

<!-- archive/zip -->

[bufio](/pkg/bufio/)

:   <!-- https://go.dev/issue/56381, CL 498117 -->
    When a [`SplitFunc`](/pkg/bufio#SplitFunc) returns [`ErrFinalToken`](/pkg/bufio#ErrFinalToken) with a `nil` token, [`Scanner`](/pkg/bufio#Scanner) will now stop immediately.
    Previously, it would report a final empty token before stopping, which was usually not desired.
    Callers that do want to report a final empty token can do so by returning `[]byte{}` rather than `nil`.

<!-- bufio -->

[cmp](/pkg/cmp/)

:   <!-- https://go.dev/issue/60204 -->
    <!-- CL 504883 -->
    The new function `Or` returns the first in a sequence of values that is not the zero value.

<!-- cmp -->

[crypto/tls](/pkg/crypto/tls/)

:   <!-- https://go.dev/issue/43922, CL 544155 -->
    [`ConnectionState.ExportKeyingMaterial`](/pkg/crypto/tls#ConnectionState.ExportKeyingMaterial) will now
    return an error unless TLS 1.3 is in use, or the `extended_master_secret` extension is supported by both the server and
    client. `crypto/tls` has supported this extension since Go 1.20. This can be disabled with the
    `tlsunsafeekm=1` GODEBUG setting.

    <!-- https://go.dev/issue/62459, CL 541516 -->
    By default, the minimum version offered by `crypto/tls` servers is now TLS 1.2 if not specified with
    [`config.MinimumVersion`](/pkg/crypto/tls#Config.MinimumVersion), matching the behavior of `crypto/tls`
    clients. This change can be reverted with the `tls10server=1` GODEBUG setting.

    <!-- https://go.dev/issue/63413, CL 541517 -->
    By default, cipher suites without ECDHE support are no longer offered by either clients or servers during pre-TLS 1.3
    handshakes. This change can be reverted with the `tlsrsakex=1` GODEBUG setting.

<!-- crypto/tls -->

[crypto/x509](/pkg/crypto/x509/)

:   <!-- https://go.dev/issue/57178 -->
    The new [`CertPool.AddCertWithConstraint`](/pkg/crypto/x509#CertPool.AddCertWithConstraint)
    method can be used to add customized constraints to root certificates to be applied during chain building.

    <!-- https://go.dev/issue/58922, CL 519315-->
    On Android, root certificates will now be loaded from `/data/misc/keychain/certs-added` as well as `/system/etc/security/cacerts`.

    <!-- https://go.dev/issue/60665, CL 520535 -->
    A new type, [`OID`](/pkg/crypto/x509#OID), supports ASN.1 Object Identifiers with individual
    components larger than 31 bits. A new field which uses this type, [`Policies`](/pkg/crypto/x509#Certificate.Policies),
    is added to the `Certificate` struct, and is now populated during parsing. Any OIDs which cannot be represented
    using a [`asn1.ObjectIdentifier`](/pkg/encoding/asn1#ObjectIdentifier) will appear in `Policies`,
    but not in the old `PolicyIdentifiers` field.
    When calling [`CreateCertificate`](/pkg/crypto/x509#CreateCertificate), the `Policies` field is ignored, and
    policies are taken from the `PolicyIdentifiers` field. Using the `x509usepolicies=1` GODEBUG setting inverts this,
    populating certificate policies from the `Policies` field, and ignoring the `PolicyIdentifiers` field. We may change the
    default value of `x509usepolicies` in Go 1.23, making `Policies` the default field for marshaling.

<!-- crypto/x509 -->

[database/sql](/pkg/database/sql/)

:   <!-- https://go.dev/issue/60370, CL 501700 -->
    The new [`Null[T]`](/pkg/database/sql/#Null) type
    provide a way to scan nullable columns for any column types.

<!-- database/sql -->

[debug/elf](/pkg/debug/elf/)

:   <!-- https://go.dev/issue/61974, CL 469395 -->
    Constant `R_MIPS_PC32` is defined for use with MIPS64 systems.

    <!-- https://go.dev/issue/63725, CL 537615 -->
    Additional `R_LARCH_*` constants are defined for use with LoongArch systems.

<!-- debug/elf -->

[encoding](/pkg/encoding/)

:   <!-- https://go.dev/issue/53693, https://go.dev/cl/504884 -->
    The new methods `AppendEncode` and `AppendDecode` added to
    each of the `Encoding` types in the packages
    [`encoding/base32`](/pkg/encoding/base32),
    [`encoding/base64`](/pkg/encoding/base64), and
    [`encoding/hex`](/pkg/encoding/hex)
    simplify encoding and decoding from and to byte slices by taking care of byte slice buffer management.

    <!-- https://go.dev/cl/505236 -->
    The methods
    [`base32.Encoding.WithPadding`](/pkg/encoding/base32#Encoding.WithPadding) and
    [`base64.Encoding.WithPadding`](/pkg/encoding/base64#Encoding.WithPadding)
    now panic if the `padding` argument is a negative value other than
    `NoPadding`.

<!-- encoding -->

[encoding/json](/pkg/encoding/json/)

:   <!-- https://go.dev/cl/521675 -->
    Marshaling and encoding functionality now escapes
    `'\b'` and `'\f'` characters as
    `\b` and `\f` instead of
    `\u0008` and `\u000c`.

<!-- encoding/json -->

[go/ast](/pkg/go/ast/)

:   <!-- https://go.dev/issue/52463, https://go/dev/cl/504915 -->
    The following declarations related to
    [syntactic identifier resolution](https://pkg.go.dev/go/ast#Object)
    are now [deprecated](/issue/52463):
    `Ident.Obj`,
    `Object`,
    `Scope`,
    `File.Scope`,
    `File.Unresolved`,
    `Importer`,
    `Package`,
    `NewPackage`.
    In general, identifiers cannot be accurately resolved without type information.
    Consider, for example, the identifier `K`
    in `T{K: ""}`: it could be the name of a local variable
    if T is a map type, or the name of a field if T is a struct type.
    New programs should use the [go/types](/pkg/go/types)
    package to resolve identifiers; see
    [`Object`](https://pkg.go.dev/go/types#Object),
    [`Info.Uses`](https://pkg.go.dev/go/types#Info.Uses), and
    [`Info.Defs`](https://pkg.go.dev/go/types#Info.Defs) for details.

    <!-- https://go.dev/issue/60061 -->
    The new [`ast.Unparen`](https://pkg.go.dev/go/ast#Unparen)
    function removes any enclosing
    [parentheses](https://pkg.go.dev/go/ast#ParenExpr) from
    an [expression](https://pkg.go.dev/go/ast#Expr).

<!-- go/ast -->

[go/types](/pkg/go/types/)

:   <!-- https://go.dev/issue/63223, CL 521956, CL 541737 -->
    The new [`Alias`](/pkg/go/types#Alias) type represents type aliases.
    Previously, type aliases were not represented explicitly, so a reference to a type alias was equivalent
    to spelling out the aliased type, and the name of the alias was lost.
    The new representation retains the intermediate `Alias`.
    This enables improved error reporting (the name of a type alias can be reported), and allows for better handling
    of cyclic type declarations involving type aliases.
    In a future release, `Alias` types will also carry [type parameter information](/issue/46477).
    The new function [`Unalias`](/pkg/go/types#Unalias) returns the actual type denoted by an
    `Alias` type (or any other [`Type`](/pkg/go/types#Type) for that matter).

    Because `Alias` types may break existing type switches that do not know to check for them,
    this functionality is controlled by a [`GODEBUG`](/doc/godebug) field named `gotypesalias`.
    With `gotypesalias=0`, everything behaves as before, and `Alias` types are never created.
    With `gotypesalias=1`, `Alias` types are created and clients must expect them.
    The default is `gotypesalias=0`.
    In a future release, the default will be changed to `gotypesalias=1`.
    _Clients of [`go/types`](/pkg/go/types) are urged to adjust their code as soon as possible
    to work with `gotypesalias=1` to eliminate problems early._

    <!-- https://go.dev/issue/62605, CL 540056 -->
    The [`Info`](/pkg/go/types#Info) struct now exports the
    [`FileVersions`](/pkg/go/types#Info.FileVersions) map
    which provides per-file Go version information.

    <!-- https://go.dev/issue/62037, CL 541575 -->
    The new helper method [`PkgNameOf`](/pkg/go/types#Info.PkgNameOf) returns the local package name
    for the given import declaration.

    <!-- https://go.dev/issue/61035, multiple CLs, see issue for details -->
    The implementation of [`SizesFor`](/pkg/go/types#SizesFor) has been adjusted to compute
    the same type sizes as the compiler when the compiler argument for `SizesFor` is `"gc"`.
    The default [`Sizes`](/pkg/go/types#Sizes) implementation used by the type checker is now
    `types.SizesFor("gc", "amd64")`.

    <!-- https://go.dev/issue/64295, CL 544035 -->
    The start position ([`Pos`](/pkg/go/types#Scope.Pos))
    of the lexical environment block ([`Scope`](/pkg/go/types#Scope))
    that represents a function body has changed:
    it used to start at the opening curly brace of the function body,
    but now starts at the function's `func` token.

[html/template](/pkg/html/template/)

:   <!-- https://go.dev/issue/61619, CL 507995 -->
    JavaScript template literals may now contain Go template actions, and parsing a template containing one will
    no longer return `ErrJSTemplate`. Similarly the GODEBUG setting `jstmpllitinterp` no
    longer has any effect.

<!-- html/template -->

[io](/pkg/io/)

:   <!-- https://go.dev/issue/61870, CL 526855 -->
    The new [`SectionReader.Outer`](/pkg/io#SectionReader.Outer) method returns the [`ReaderAt`](/pkg/io#ReaderAt), offset, and size passed to [`NewSectionReader`](/pkg/io#NewSectionReader).

<!-- io -->

[log/slog](/pkg/log/slog/)

:   <!-- https://go.dev/issue/62418 -->
    The new [`SetLogLoggerLevel`](/pkg/log/slog#SetLogLoggerLevel) function
    controls the level for the bridge between the `slog` and `log` packages. It sets the minimum level
    for calls to the top-level `slog` logging functions, and it sets the level for calls to `log.Logger`
    that go through `slog`.

[math/big](/pkg/math/big/)

:   <!-- https://go.dev/issue/50489, CL 539299 -->
    The new method [`Rat.FloatPrec`](/pkg/math/big#Rat.FloatPrec) computes the number of fractional decimal digits
    required to represent a rational number accurately as a floating-point number, and whether accurate decimal representation
    is possible in the first place.

<!-- math/big -->

[net](/pkg/net/)

:   <!-- https://go.dev/issue/58808 -->
    When [`io.Copy`](/pkg/io#Copy) copies
    from a `TCPConn` to a `UnixConn`,
    it will now use Linux's `splice(2)` system call if possible,
    using the new method [`TCPConn.WriteTo`](/pkg/net#TCPConn.WriteTo).

    <!-- CL 467335 -->
    The Go DNS Resolver, used when building with "-tags=netgo",
    now searches for a matching name in the Windows hosts file,
    located at `%SystemRoot%\System32\drivers\etc\hosts`,
    before making a DNS query.

<!-- net -->

[net/http](/pkg/net/http/)

:   <!-- https://go.dev/issue/51971 -->
    The new functions
    [`ServeFileFS`](/pkg/net/http#ServeFileFS),
    [`FileServerFS`](/pkg/net/http#FileServerFS), and
    [`NewFileTransportFS`](/pkg/net/http#NewFileTransportFS)
    are versions of the existing
    `ServeFile`, `FileServer`, and `NewFileTransport`,
    operating on an `fs.FS`.

    <!-- https://go.dev/issue/61679 -->
    The HTTP server and client now reject requests and responses containing
    an invalid empty `Content-Length` header.
    The previous behavior may be restored by setting
    [`GODEBUG`](/doc/godebug) field `httplaxcontentlength=1`.

    <!-- https://go.dev/issue/61410, CL 528355 -->
    The new method
    [`Request.PathValue`](/pkg/net/http#Request.PathValue)
    returns path wildcard values from a request
    and the new method
    [`Request.SetPathValue`](/pkg/net/http#Request.SetPathValue)
    sets path wildcard values on a request.

<!-- net/http -->

[net/http/cgi](/pkg/net/http/cgi/)

:   <!-- CL 539615 -->
    When executing a CGI process, the `PATH_INFO` variable is now
    always set to the empty string or a value starting with a `/` character,
    as required by RFC 3875. It was previously possible for some combinations of
    [`Handler.Root`](/pkg/net/http/cgi#Handler.Root)
    and request URL to violate this requirement.

<!-- net/http/cgi -->

[net/netip](/pkg/net/netip/)

:   <!-- https://go.dev/issue/61642 -->
    The new [`AddrPort.Compare`](/pkg/net/netip#AddrPort.Compare)
    method compares two `AddrPort`s.

<!-- net/netip -->

[os](/pkg/os/)

:   <!-- CL 516555 -->
    On Windows, the [`Stat`](/pkg/os#Stat) function now follows all reparse points
    that link to another named entity in the system.
    It was previously only following `IO_REPARSE_TAG_SYMLINK` and
    `IO_REPARSE_TAG_MOUNT_POINT` reparse points.

    <!-- CL 541015 -->
    On Windows, passing [`O_SYNC`](/pkg/os#O_SYNC) to [`OpenFile`](/pkg/os#OpenFile) now causes write operations to go directly to disk, equivalent to `O_SYNC` on Unix platforms.

    <!-- CL 452995 -->
    On Windows, the [`ReadDir`](/pkg/os#ReadDir),
    [`File.ReadDir`](/pkg/os#File.ReadDir),
    [`File.Readdir`](/pkg/os#File.Readdir),
    and [`File.Readdirnames`](/pkg/os#File.Readdirnames) functions
    now read directory entries in batches to reduce the number of system calls,
    improving performance up to 30%.

    <!-- https://go.dev/issue/58808 -->
    When [`io.Copy`](/pkg/io#Copy) copies
    from a `File` to a `net.UnixConn`,
    it will now use Linux's `sendfile(2)` system call if possible,
    using the new method [`File.WriteTo`](/pkg/os#File.WriteTo).

<!-- os -->

[os/exec](/pkg/os/exec/)

:   <!-- CL 528037 -->
    On Windows, [`LookPath`](/pkg/os/exec#LookPath) now
    ignores empty entries in `%PATH%`, and returns
    `ErrNotFound` (instead of `ErrNotExist`) if
    no executable file extension is found to resolve an otherwise-unambiguous
    name.

    <!-- CL 528038, CL 527820 -->
    On Windows, [`Command`](/pkg/os/exec#Command) and
    [`Cmd.Start`](/pkg/os/exec#Cmd.Start) no
    longer call `LookPath` if the path to the executable is already
    absolute and has an executable file extension. In addition,
    `Cmd.Start` no longer writes the resolved extension back to
    the [`Path`](/pkg/os/exec#Cmd.Path) field,
    so it is now safe to call the `String` method concurrently
    with a call to `Start`.

<!-- os/exec -->

[reflect](/pkg/reflect/)

:   <!-- https://go.dev/issue/61827, CL 517777 -->
    The [`Value.IsZero`](/pkg/reflect/#Value.IsZero)
    method will now return true for a floating-point or complex
    negative zero, and will return true for a struct value if a
    blank field (a field named `_`) somehow has a
    non-zero value.
    These changes make `IsZero` consistent with comparing
    a value to zero using the language `==` operator.

    <!-- https://go.dev/issue/59599, CL 511035 -->
    The [`PtrTo`](/pkg/reflect/#PtrTo) function is deprecated,
    in favor of [`PointerTo`](/pkg/reflect/#PointerTo).

    <!-- https://go.dev/issue/60088, CL 513478 -->
    The new function [`TypeFor`](/pkg/reflect/#TypeFor)
    returns the [`Type`](/pkg/reflect/#Type) that represents
    the type argument T.
    Previously, to get the `reflect.Type` value for a type, one had to use
    `reflect.TypeOf((*T)(nil)).Elem()`.
    This may now be written as `reflect.TypeFor[T]()`.

<!-- reflect -->

[runtime/metrics](/pkg/runtime/metrics/)

:   <!-- https://go.dev/issue/63340 -->
    Four new histogram metrics
    `/sched/pauses/stopping/gc:seconds`,
    `/sched/pauses/stopping/other:seconds`,
    `/sched/pauses/total/gc:seconds`, and
    `/sched/pauses/total/other:seconds` provide additional details
    about stop-the-world pauses.
    The "stopping" metrics report the time taken from deciding to stop the
    world until all goroutines are stopped.
    The "total" metrics report the time taken from deciding to stop the world
    until it is started again.

    <!-- https://go.dev/issue/63340 -->
    The `/gc/pauses:seconds` metric is deprecated, as it is
    equivalent to the new `/sched/pauses/total/gc:seconds` metric.

    <!-- https://go.dev/issue/57071 -->
    `/sync/mutex/wait/total:seconds` now includes contention on
    runtime-internal locks in addition to
    [`sync.Mutex`](/pkg/sync#Mutex) and
    [`sync.RWMutex`](/pkg/sync#RWMutex).

<!-- runtime/metrics -->

[runtime/pprof](/pkg/runtime/pprof/)

:   <!-- https://go.dev/issue/61015 -->
    Mutex profiles now scale contention by the number of goroutines blocked on the mutex.
    This provides a more accurate representation of the degree to which a mutex is a bottleneck in
    a Go program.
    For instance, if 100 goroutines are blocked on a mutex for 10 milliseconds, a mutex profile will
    now record 1 second of delay instead of 10 milliseconds of delay.

    <!-- https://go.dev/issue/57071 -->
    Mutex profiles also now include contention on runtime-internal locks in addition to
    [`sync.Mutex`](/pkg/sync#Mutex) and
    [`sync.RWMutex`](/pkg/sync#RWMutex).
    Contention on runtime-internal locks is always reported at `runtime._LostContendedRuntimeLock`.
    A future release will add complete stack traces in these cases.

    <!-- https://go.dev/issue/50891 -->
    CPU profiles on Darwin platforms now contain the process's memory map, enabling the disassembly
    view in the pprof tool.

<!-- runtime/pprof -->

[runtime/trace](/pkg/runtime/trace/)

:   <!-- https://go.dev/issue/60773 -->
    The execution tracer has been completely overhauled in this release, resolving several long-standing
    issues and paving the way for new use-cases for execution traces.

    Execution traces now use the operating system's clock on most platforms (Windows excluded) so
    it is possible to correlate them with traces produced by lower-level components.
    Execution traces no longer depend on the reliability of the platform's clock to produce a correct trace.
    Execution traces are now partitioned regularly on-the-fly and as a result may be processed in a
    streamable way.
    Execution traces now contain complete durations for all system calls.
    Execution traces now contain information about the operating system threads that goroutines executed on.
    The latency impact of starting and stopping execution traces has been dramatically reduced.
    Execution traces may now begin or end during the garbage collection mark phase.

    To allow Go developers to take advantage of these improvements, an experimental
    trace reading package is available at [golang.org/x/exp/trace](/pkg/golang.org/x/exp/trace).
    Note that this package only works on traces produced by programs built with Go 1.22 at the moment.
    Please try out the package and provide feedback on
    [the corresponding proposal issue](/issue/62627).

    If you experience any issues with the new execution tracer implementation, you may switch back to the
    old implementation by building your Go program with `GOEXPERIMENT=noexectracer2`.
    If you do, please file an issue, otherwise this option will be removed in a future release.

<!-- runtime/trace -->

[slices](/pkg/slices/)

:   <!-- https://go.dev/issue/56353 -->
    <!-- CL 504882 -->
    The new function `Concat` concatenates multiple slices.

    <!-- https://go.dev/issue/63393 -->
    <!-- CL 543335 -->
    Functions that shrink the size of a slice (`Delete`, `DeleteFunc`, `Compact`, `CompactFunc`, and `Replace`) now zero the elements between the new length and the old length.

    <!-- https://go.dev/issue/63913 -->
    <!-- CL 540155 -->
    `Insert` now always panics if the argument `i` is out of range. Previously it did not panic in this situation if there were no elements to be inserted.

<!-- slices -->

[syscall](/pkg/syscall/)

:   <!-- https://go.dev/issue/60797 -->
    The `syscall` package has been [frozen](/s/go1.4-syscall) since Go 1.4 and was marked as deprecated in Go 1.11, causing many editors to warn about any use of the package.
    However, some non-deprecated functionality requires use of the `syscall` package, such as the [`os/exec.Cmd.SysProcAttr`](/pkg/os/exec#Cmd) field.
    To avoid unnecessary complaints on such code, the `syscall` package is no longer marked as deprecated.
    The package remains frozen to most new functionality, and new code remains encouraged to use [`golang.org/x/sys/unix`](/pkg/golang.org/x/sys/unix) or [`golang.org/x/sys/windows`](/pkg/golang.org/x/sys/windows) where possible.

    <!-- https://go.dev/issue/51246, CL 520266 -->
    On Linux, the new [`SysProcAttr.PidFD`](/pkg/syscall#SysProcAttr) field allows obtaining a PID FD when starting a child process via [`StartProcess`](/pkg/syscall#StartProcess) or [`os/exec`](/pkg/os/exec).

    <!-- CL 541015 -->
    On Windows, passing [`O_SYNC`](/pkg/syscall#O_SYNC) to [`Open`](/pkg/syscall#Open) now causes write operations to go directly to disk, equivalent to `O_SYNC` on Unix platforms.

<!-- syscall -->

[testing/slogtest](/pkg/testing/slogtest/)

:   <!-- https://go.dev/issue/61758 -->
    The new [`Run`](/pkg/testing/slogtest#Run) function uses sub-tests to run test cases,
    providing finer-grained control.

<!-- testing/slogtest -->

## Ports {#ports}

### Darwin {#darwin}

<!-- CL 461697 -->
On macOS on 64-bit x86 architecture (the `darwin/amd64` port),
the Go toolchain now generates position-independent executables (PIE) by default.
Non-PIE binaries can be generated by specifying the `-buildmode=exe`
build flag.
On 64-bit ARM-based macOS (the `darwin/arm64` port),
the Go toolchain already generates PIE by default.

<!-- go.dev/issue/64207 -->
Go 1.22 is the last release that will run on macOS 10.15 Catalina. Go 1.23 will require macOS 11 Big Sur or later.

### ARM {#arm}

<!-- CL 514907 -->
The `GOARM` environment variable now allows you to select whether to use software or hardware floating point.
Previously, valid `GOARM` values were `5`, `6`, or `7`. Now those same values can
be optionally followed by `,softfloat` or `,hardfloat` to select the floating-point implementation.

This new option defaults to `softfloat` for version `5` and `hardfloat` for versions
`6` and `7`.

### Loong64 {#loong64}

<!-- CL 481315 -->
The `loong64` port now supports passing function arguments and results using registers.

<!-- CL 481315,537615,480878 -->
The `linux/loong64` port now supports the address sanitizer, memory sanitizer, new-style linker relocations, and the `plugin` build mode.

### OpenBSD {#openbsd}

<!-- CL 517935 -->
Go 1.22 adds an experimental port to OpenBSD on big-endian 64-bit PowerPC
(`openbsd/ppc64`).
