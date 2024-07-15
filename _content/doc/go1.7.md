---
template: false
title: Go 1.7 Release Notes
---

<!--
for acme:
Edit .,s;^PKG:([a-z][A-Za-z0-9_/]+);<a href="/pkg/\1/"><code>\1</code></a>;g
Edit .,s;^([a-z][A-Za-z0-9_/]+)\.([A-Z][A-Za-z0-9_]+\.)?([A-Z][A-Za-z0-9_]+)([ .',)]|$);<a href="/pkg/\1/#\2\3"><code>\3</code></a>\4;g
Edit .,s;^FULL:([a-z][A-Za-z0-9_/]+)\.([A-Z][A-Za-z0-9_]+\.)?([A-Z][A-Za-z0-9_]+)([ .',)]|$);<a href="/pkg/\1/#\2\3"><code>\1.\2\3</code></a>\4;g
Edit .,s;^DPKG:([a-z][A-Za-z0-9_/]+);<dl id="\1"><a href="/pkg/\1/">\1</a></dl>;g
rsc last updated through 6729576
-->

<!--
NOTE: In this document and others in this directory, the convention is to
set fixed-width phrases with non-fixed-width spaces, as in
`hello` `world`.
Do not send CLs removing the interior tags from such phrases.
-->

<style>
  main ul li { margin: 0.5em 0; }
</style>

## Introduction to Go 1.7 {#introduction}

The latest Go release, version 1.7, arrives six months after 1.6.
Most of its changes are in the implementation of the toolchain, runtime, and libraries.
There is one minor change to the language specification.
As always, the release maintains the Go 1 [promise of compatibility](/doc/go1compat.html).
We expect almost all Go programs to continue to compile and run as before.

The release [adds a port to IBM LinuxOne](#ports);
[updates the x86-64 compiler back end](#compiler) to generate more efficient code;
includes the [context package](#context), promoted from the
[x/net subrepository](https://golang.org/x/net/context)
and now used in the standard library;
and [adds support in the testing package](#testing) for
creating hierarchies of tests and benchmarks.
The release also [finalizes the vendoring support](#cmd_go)
started in Go 1.5, making it a standard feature.

## Changes to the language {#language}

There is one tiny language change in this release.
The section on [terminating statements](/ref/spec#Terminating_statements)
clarifies that to determine whether a statement list ends in a terminating statement,
the “final non-empty statement” is considered the end,
matching the existing behavior of the gc and gccgo compiler toolchains.
In earlier releases the definition referred only to the “final statement,”
leaving the effect of trailing empty statements at the least unclear.
The [`go/types`](/pkg/go/types/)
package has been updated to match the gc and gccgo compiler toolchains
in this respect.
This change has no effect on the correctness of existing programs.

## Ports {#ports}

Go 1.7 adds support for macOS 10.12 Sierra.
Binaries built with versions of Go before 1.7 will not work
correctly on Sierra.

Go 1.7 adds an experimental port to [Linux on z Systems](https://en.wikipedia.org/wiki/Linux_on_z_Systems) (`linux/s390x`)
and the beginning of a port to Plan 9 on ARM (`plan9/arm`).

The experimental ports to Linux on 64-bit MIPS (`linux/mips64` and `linux/mips64le`)
added in Go 1.6 now have full support for cgo and external linking.

The experimental port to Linux on little-endian 64-bit PowerPC (`linux/ppc64le`)
now requires the POWER8 architecture or later.
Big-endian 64-bit PowerPC (`linux/ppc64`) only requires the
POWER5 architecture.

The OpenBSD port now requires OpenBSD 5.6 or later, for access to the [_getentropy_(2)](https://man.openbsd.org/getentropy.2) system call.

### Known Issues {#known_issues}

There are some instabilities on FreeBSD that are known but not understood.
These can lead to program crashes in rare cases.
See [issue 16136](/issue/16136),
[issue 15658](/issue/15658),
and [issue 16396](/issue/16396).
Any help in solving these FreeBSD-specific issues would be appreciated.

## Tools {#tools}

### Assembler {#cmd_asm}

For 64-bit ARM systems, the vector register names have been
corrected to `V0` through `V31`;
previous releases incorrectly referred to them as `V32` through `V63`.

For 64-bit x86 systems, the following instructions have been added:
`PCMPESTRI`,
`RORXL`,
`RORXQ`,
`VINSERTI128`,
`VPADDD`,
`VPADDQ`,
`VPALIGNR`,
`VPBLENDD`,
`VPERM2F128`,
`VPERM2I128`,
`VPOR`,
`VPSHUFB`,
`VPSHUFD`,
`VPSLLD`,
`VPSLLDQ`,
`VPSLLQ`,
`VPSRLD`,
`VPSRLDQ`,
and
`VPSRLQ`.

### Compiler Toolchain {#compiler}

This release includes a new code generation back end for 64-bit x86 systems,
following a [proposal from 2015](/s/go17ssa)
that has been under development since then.
The new back end, based on
[SSA](https://en.wikipedia.org/wiki/Static_single_assignment_form),
generates more compact, more efficient code
and provides a better platform for optimizations
such as bounds check elimination.
The new back end reduces the CPU time required by
our benchmark programs by 5-35%.

For this release, the new back end can be disabled by passing
`-ssa=0` to the compiler.
If you find that your program compiles or runs successfully
only with the new back end disabled, please
[file a bug report](/issue/new).

The format of exported metadata written by the compiler in package archives has changed:
the old textual format has been replaced by a more compact binary format.
This results in somewhat smaller package archives and fixes a few
long-standing corner case bugs.

For this release, the new export format can be disabled by passing
`-newexport=0` to the compiler.
If you find that your program compiles or runs successfully
only with the new export format disabled, please
[file a bug report](/issue/new).

The linker's `-X` option no longer supports the unusual two-argument form
`-X` `name` `value`,
as [announced](/doc/go1.6#compiler) in the Go 1.6 release
and in warnings printed by the linker.
Use `-X` `name=value` instead.

The compiler and linker have been optimized and run significantly faster in this release than in Go 1.6,
although they are still slower than we would like and will continue to be optimized in future releases.

Due to changes across the compiler toolchain and standard library,
binaries built with this release should typically be smaller than binaries
built with Go 1.6,
sometimes by as much as 20-30%.

On x86-64 systems, Go programs now maintain stack frame pointers
as expected by profiling tools like Linux's perf and Intel's VTune,
making it easier to analyze and optimize Go programs using these tools.
The frame pointer maintenance has a small run-time overhead that varies
but averages around 2%. We hope to reduce this cost in future releases.
To build a toolchain that does not use frame pointers, set
`GOEXPERIMENT=noframepointer` when running
`make.bash`, `make.bat`, or `make.rc`.

### Cgo {#cmd_cgo}

Packages using [cgo](/cmd/cgo/) may now include
Fortran source files (in addition to C, C++, Objective C, and SWIG),
although the Go bindings must still use C language APIs.

Go bindings may now use a new helper function `C.CBytes`.
In contrast to `C.CString`, which takes a Go `string`
and returns a `*C.byte` (a C `char*`),
`C.CBytes` takes a Go `[]byte`
and returns an `unsafe.Pointer` (a C `void*`).

Packages and binaries built using `cgo` have in past releases
produced different output on each build,
due to the embedding of temporary directory names.
When using this release with
new enough versions of GCC or Clang
(those that support the `-fdebug-prefix-map` option),
those builds should finally be deterministic.

### Gccgo {#gccgo}

Due to the alignment of Go's semiannual release schedule with GCC's annual release schedule,
GCC release 6 contains the Go 1.6.1 version of gccgo.
The next release, GCC 7, will likely have the Go 1.8 version of gccgo.

### Go command {#cmd_go}

The [`go`](/cmd/go/) command's basic operation
is unchanged, but there are a number of changes worth noting.

This release removes support for the `GO15VENDOREXPERIMENT` environment variable,
as [announced](/doc/go1.6#go_command) in the Go 1.6 release.
[Vendoring support](/s/go15vendor)
is now a standard feature of the `go` command and toolchain.

The `Package` data structure made available to
“`go` `list`” now includes a
`StaleReason` field explaining why a particular package
is or is not considered stale (in need of rebuilding).
This field is available to the `-f` or `-json`
options and is useful for understanding why a target is being rebuilt.

The “`go` `get`” command now supports
import paths referring to `git.openstack.org`.

This release adds experimental, minimal support for building programs using
[binary-only packages](/pkg/go/build#hdr-Binary_Only_Packages),
packages distributed in binary form
without the corresponding source code.
This feature is needed in some commercial settings
but is not intended to be fully integrated into the rest of the toolchain.
For example, tools that assume access to complete source code
will not work with such packages, and there are no plans to support
such packages in the “`go` `get`” command.

### Go doc {#cmd_doc}

The “`go` `doc`” command
now groups constructors with the type they construct,
following [`godoc`](/cmd/godoc/).

### Go vet {#cmd_vet}

The “`go` `vet`” command
has more accurate analysis in its `-copylock` and `-printf` checks,
and a new `-tests` check that checks the name and signature of likely test functions.
To avoid confusion with the new `-tests` check, the old, unadvertised
`-test` option has been removed; it was equivalent to `-all` `-shadow`.

The `vet` command also has a new check,
`-lostcancel`, which detects failure to call the
cancellation function returned by the `WithCancel`,
`WithTimeout`, and `WithDeadline` functions in
Go 1.7's new `context` package (see [below](#context)).
Failure to call the function prevents the new `Context`
from being reclaimed until its parent is canceled.
(The background context is never canceled.)

### Go tool dist {#cmd_dist}

The new subcommand “`go` `tool` `dist` `list`”
prints all supported operating system/architecture pairs.

### Go tool trace {#cmd_trace}

The “`go` `tool` `trace`” command,
[introduced in Go 1.5](/doc/go1.5#trace_command),
has been refined in various ways.

First, collecting traces is significantly more efficient than in past releases.
In this release, the typical execution-time overhead of collecting a trace is about 25%;
in past releases it was at least 400%.
Second, trace files now include file and line number information,
making them more self-contained and making the
original executable optional when running the trace tool.
Third, the trace tool now breaks up large traces to avoid limits
in the browser-based viewer.

Although the trace file format has changed in this release,
the Go 1.7 tools can still read traces from earlier releases.

## Performance {#performance}

As always, the changes are so general and varied that precise statements
about performance are difficult to make.
Most programs should run a bit faster,
due to speedups in the garbage collector and
optimizations in the core library.
On x86-64 systems, many programs will run significantly faster,
due to improvements in generated code brought by the
new compiler back end.
As noted above, in our own benchmarks,
the code generation changes alone typically reduce program CPU time by 5-35%.

<!-- git log -''-grep '-[0-9][0-9]\.[0-9][0-9]%' go1.6.. -->
There have been significant optimizations bringing more than 10% improvements
to implementations in the
[`crypto/sha1`](/pkg/crypto/sha1/),
[`crypto/sha256`](/pkg/crypto/sha256/),
[`encoding/binary`](/pkg/encoding/binary/),
[`fmt`](/pkg/fmt/),
[`hash/adler32`](/pkg/hash/adler32/),
[`hash/crc32`](/pkg/hash/crc32/),
[`hash/crc64`](/pkg/hash/crc64/),
[`image/color`](/pkg/image/color/),
[`math/big`](/pkg/math/big/),
[`strconv`](/pkg/strconv/),
[`strings`](/pkg/strings/),
[`unicode`](/pkg/unicode/),
and
[`unicode/utf16`](/pkg/unicode/utf16/)
packages.

Garbage collection pauses should be significantly shorter than they
were in Go 1.6 for programs with large numbers of idle goroutines,
substantial stack size fluctuation, or large package-level variables.

## Standard library {#library}

### Context {#context}

Go 1.7 moves the `golang.org/x/net/context` package
into the standard library as [`context`](/pkg/context/).
This allows the use of contexts for cancellation, timeouts, and passing
request-scoped data in other standard library packages,
including
[net](#net),
[net/http](#net_http),
and
[os/exec](#os_exec),
as noted below.

For more information about contexts, see the
[package documentation](/pkg/context/)
and the Go blog post
“[Go Concurrent Patterns: Context](/blog/context).”

### HTTP Tracing {#httptrace}

Go 1.7 introduces [`net/http/httptrace`](/pkg/net/http/httptrace/),
a package that provides mechanisms for tracing events within HTTP requests.

### Testing {#testing}

The `testing` package now supports the definition
of tests with subtests and benchmarks with sub-benchmarks.
This support makes it easy to write table-driven benchmarks
and to create hierarchical tests.
It also provides a way to share common setup and tear-down code.
See the [package documentation](/pkg/testing/#hdr-Subtests_and_Sub_benchmarks) for details.

### Runtime {#runtime}

All panics started by the runtime now use panic values
that implement both the
builtin [`error`](/ref/spec#Errors),
and
[`runtime.Error`](/pkg/runtime/#Error),
as
[required by the language specification](/ref/spec#Run_time_panics).

During panics, if a signal's name is known, it will be printed in the stack trace.
Otherwise, the signal's number will be used, as it was before Go1.7.

The new function
[`KeepAlive`](/pkg/runtime/#KeepAlive)
provides an explicit mechanism for declaring
that an allocated object must be considered reachable
at a particular point in a program,
typically to delay the execution of an associated finalizer.

The new function
[`CallersFrames`](/pkg/runtime/#CallersFrames)
translates a PC slice obtained from
[`Callers`](/pkg/runtime/#Callers)
into a sequence of frames corresponding to the call stack.
This new API should be preferred instead of direct use of
[`FuncForPC`](/pkg/runtime/#FuncForPC),
because the frame sequence can more accurately describe
call stacks with inlined function calls.

The new function
[`SetCgoTraceback`](/pkg/runtime/#SetCgoTraceback)
facilitates tighter integration between Go and C code executing
in the same process called using cgo.

On 32-bit systems, the runtime can now use memory allocated
by the operating system anywhere in the address space,
eliminating the
“memory allocated by OS not in usable range” failure
common in some environments.

The runtime can now return unused memory to the operating system on
all architectures.
In Go 1.6 and earlier, the runtime could not
release memory on ARM64, 64-bit PowerPC, or MIPS.

On Windows, Go programs in Go 1.5 and earlier forced
the global Windows timer resolution to 1ms at startup
by calling `timeBeginPeriod(1)`.
Changing the global timer resolution caused problems on some systems,
and testing suggested that the call was not needed for good scheduler performance,
so Go 1.6 removed the call.
Go 1.7 brings the call back: under some workloads the call
is still needed for good scheduler performance.

### Minor changes to the library {#minor_library_changes}

As always, there are various minor changes and updates to the library,
made with the Go 1 [promise of compatibility](/doc/go1compat)
in mind.

#### [bufio](/pkg/bufio/)

In previous releases of Go, if
[`Reader`](/pkg/bufio/#Reader)'s
[`Peek`](/pkg/bufio/#Reader.Peek) method
were asked for more bytes than fit in the underlying buffer,
it would return an empty slice and the error `ErrBufferFull`.
Now it returns the entire underlying buffer, still accompanied by the error `ErrBufferFull`.

#### [bytes](/pkg/bytes/)

The new functions
[`ContainsAny`](/pkg/bytes/#ContainsAny) and
[`ContainsRune`](/pkg/bytes/#ContainsRune)
have been added for symmetry with
the [`strings`](/pkg/strings/) package.

In previous releases of Go, if
[`Reader`](/pkg/bytes/#Reader)'s
[`Read`](/pkg/bytes/#Reader.Read) method
were asked for zero bytes with no data remaining, it would
return a count of 0 and no error.
Now it returns a count of 0 and the error
[`io.EOF`](/pkg/io/#EOF).

The
[`Reader`](/pkg/bytes/#Reader) type has a new method
[`Reset`](/pkg/bytes/#Reader.Reset) to allow reuse of a `Reader`.

#### [compress/flate](/pkg/compress/flate/)

There are many performance optimizations throughout the package.
Decompression speed is improved by about 10%,
while compression for `DefaultCompression` is twice as fast.

In addition to those general improvements,
the
`BestSpeed`
compressor has been replaced entirely and uses an
algorithm similar to [Snappy](https://github.com/google/snappy),
resulting in about a 2.5X speed increase,
although the output can be 5-10% larger than with the previous algorithm.

There is also a new compression level
`HuffmanOnly`
that applies Huffman but not Lempel-Ziv encoding.
[Forgoing Lempel-Ziv encoding](https://blog.klauspost.com/constant-time-gzipzip-compression/) means that
`HuffmanOnly` runs about 3X faster than the new `BestSpeed`
but at the cost of producing compressed outputs that are 20-40% larger than those
generated by the new `BestSpeed`.

It is important to note that both
`BestSpeed` and `HuffmanOnly` produce a compressed output that is
[RFC 1951](https://tools.ietf.org/html/rfc1951) compliant.
In other words, any valid DEFLATE decompressor will continue to be able to decompress these outputs.

Lastly, there is a minor change to the decompressor's implementation of
[`io.Reader`](/pkg/io/#Reader). In previous versions,
the decompressor deferred reporting
[`io.EOF`](/pkg/io/#EOF) until exactly no more bytes could be read.
Now, it reports
[`io.EOF`](/pkg/io/#EOF) more eagerly when reading the last set of bytes.

#### [crypto/tls](/pkg/crypto/tls/)

The TLS implementation sends the first few data packets on each connection
using small record sizes, gradually increasing to the TLS maximum record size.
This heuristic reduces the amount of data that must be received before
the first packet can be decrypted, improving communication latency over
low-bandwidth networks.
Setting
[`Config`](/pkg/crypto/tls/#Config)'s
`DynamicRecordSizingDisabled` field to true
forces the behavior of Go 1.6 and earlier, where packets are
as large as possible from the start of the connection.

The TLS client now has optional, limited support for server-initiated renegotiation,
enabled by setting the
[`Config`](/pkg/crypto/tls/#Config)'s
`Renegotiation` field.
This is needed for connecting to many Microsoft Azure servers.

The errors returned by the package now consistently begin with a
`tls:` prefix.
In past releases, some errors used a `crypto/tls:` prefix,
some used a `tls:` prefix, and some had no prefix at all.

When generating self-signed certificates, the package no longer sets the
“Authority Key Identifier” field by default.

#### [crypto/x509](/pkg/crypto/x509/)

The new function
[`SystemCertPool`](/pkg/crypto/x509/#SystemCertPool)
provides access to the entire system certificate pool if available.
There is also a new associated error type
[`SystemRootsError`](/pkg/crypto/x509/#SystemRootsError).

#### [debug/dwarf](/pkg/debug/dwarf/)

The
[`Reader`](/pkg/debug/dwarf/#Reader) type's new
[`SeekPC`](/pkg/debug/dwarf/#Reader.SeekPC) method and the
[`Data`](/pkg/debug/dwarf/#Data) type's new
[`Ranges`](/pkg/debug/dwarf/#Ranges) method
help to find the compilation unit to pass to a
[`LineReader`](/pkg/debug/dwarf/#LineReader)
and to identify the specific function for a given program counter.

#### [debug/elf](/pkg/debug/elf/)

The new
[`R_390`](/pkg/debug/elf/#R_390) relocation type
and its many predefined constants
support the S390 port.

#### [encoding/asn1](/pkg/encoding/asn1/)

The ASN.1 decoder now rejects non-minimal integer encodings.
This may cause the package to reject some invalid but formerly accepted ASN.1 data.

#### [encoding/json](/pkg/encoding/json/)

The
[`Encoder`](/pkg/encoding/json/#Encoder)'s new
[`SetIndent`](/pkg/encoding/json/#Encoder.SetIndent) method
sets the indentation parameters for JSON encoding,
like in the top-level
[`Indent`](/pkg/encoding/json/#Indent) function.

The
[`Encoder`](/pkg/encoding/json/#Encoder)'s new
[`SetEscapeHTML`](/pkg/encoding/json/#Encoder.SetEscapeHTML) method
controls whether the
`&`, `<`, and `>`
characters in quoted strings should be escaped as
`\u0026`, `\u003c`, and `\u003e`,
respectively.
As in previous releases, the encoder defaults to applying this escaping,
to avoid certain problems that can arise when embedding JSON in HTML.

In earlier versions of Go, this package only supported encoding and decoding
maps using keys with string types.
Go 1.7 adds support for maps using keys with integer types:
the encoding uses a quoted decimal representation as the JSON key.
Go 1.7 also adds support for encoding maps using non-string keys that implement
the `MarshalText`
(see
[`encoding.TextMarshaler`](/pkg/encoding/#TextMarshaler))
method,
as well as support for decoding maps using non-string keys that implement
the `UnmarshalText`
(see
[`encoding.TextUnmarshaler`](/pkg/encoding/#TextUnmarshaler))
method.
These methods are ignored for keys with string types in order to preserve
the encoding and decoding used in earlier versions of Go.

When encoding a slice of typed bytes,
[`Marshal`](/pkg/encoding/json/#Marshal)
now generates an array of elements encoded using
that byte type's
`MarshalJSON`
or
`MarshalText`
method if present,
only falling back to the default base64-encoded string data if neither method is available.
Earlier versions of Go accept both the original base64-encoded string encoding
and the array encoding (assuming the byte type also implements
`UnmarshalJSON`
or
`UnmarshalText`
as appropriate),
so this change should be semantically backwards compatible with earlier versions of Go,
even though it does change the chosen encoding.

#### [go/build](/pkg/go/build/)

To implement the go command's new support for binary-only packages
and for Fortran code in cgo-based packages,
the
[`Package`](/pkg/go/build/#Package) type
adds new fields `BinaryOnly`, `CgoFFLAGS`, and `FFiles`.

#### [go/doc](/pkg/go/doc/)

To support the corresponding change in `go` `test` described above,
[`Example`](/pkg/go/doc/#Example) struct adds an Unordered field
indicating whether the example may generate its output lines in any order.

#### [io](/pkg/io/)

The package adds new constants
`SeekStart`, `SeekCurrent`, and `SeekEnd`,
for use with
[`Seeker`](/pkg/io/#Seeker)
implementations.
These constants are preferred over `os.SEEK_SET`, `os.SEEK_CUR`, and `os.SEEK_END`,
but the latter will be preserved for compatibility.

#### [math/big](/pkg/math/big/)

The
[`Float`](/pkg/math/big/#Float) type adds
[`GobEncode`](/pkg/math/big/#Float.GobEncode) and
[`GobDecode`](/pkg/math/big/#Float.GobDecode) methods,
so that values of type `Float` can now be encoded and decoded using the
[`encoding/gob`](/pkg/encoding/gob/)
package.

#### [math/rand](/pkg/math/rand/)

The
[`Read`](/pkg/math/rand/#Read) function and
[`Rand`](/pkg/math/rand/#Rand)'s
[`Read`](/pkg/math/rand/#Rand.Read) method
now produce a pseudo-random stream of bytes that is consistent and not
dependent on the size of the input buffer.

The documentation clarifies that
Rand's [`Seed`](/pkg/math/rand/#Rand.Seed)
and [`Read`](/pkg/math/rand/#Rand.Read) methods
are not safe to call concurrently, though the global
functions [`Seed`](/pkg/math/rand/#Seed)
and [`Read`](/pkg/math/rand/#Read) are (and have
always been) safe.

#### [mime/multipart](/pkg/mime/multipart/)

The
[`Writer`](/pkg/mime/multipart/#Writer)
implementation now emits each multipart section's header sorted by key.
Previously, iteration over a map caused the section header to use a
non-deterministic order.

#### [net](/pkg/net/)

As part of the introduction of [context](#context), the
[`Dialer`](/pkg/net/#Dialer) type has a new method
[`DialContext`](/pkg/net/#Dialer.DialContext), like
[`Dial`](/pkg/net/#Dialer.Dial) but adding the
[`context.Context`](/pkg/context/#Context)
for the dial operation.
The context is intended to obsolete the `Dialer`'s
`Cancel` and `Deadline` fields,
but the implementation continues to respect them,
for backwards compatibility.

The
[`IP`](/pkg/net/#IP) type's
[`String`](/pkg/net/#IP.String) method has changed its result for invalid `IP` addresses.
In past releases, if an `IP` byte slice had length other than 0, 4, or 16, `String`
returned `"?"`.
Go 1.7 adds the hexadecimal encoding of the bytes, as in `"?12ab"`.

The pure Go [name resolution](/pkg/net/#hdr-Name_Resolution)
implementation now respects `nsswitch.conf`'s
stated preference for the priority of DNS lookups compared to
local file (that is, `/etc/hosts`) lookups.

#### [net/http](/pkg/net/http/)

[`ResponseWriter`](/pkg/net/http/#ResponseWriter)'s
documentation now makes clear that beginning to write the response
may prevent future reads on the request body.
For maximal compatibility, implementations are encouraged to
read the request body completely before writing any part of the response.

As part of the introduction of [context](#context), the
[`Request`](/pkg/net/http/#Request) has a new methods
[`Context`](/pkg/net/http/#Request.Context), to retrieve the associated context, and
[`WithContext`](/pkg/net/http/#Request.WithContext), to construct a copy of `Request`
with a modified context.

In the
[`Server`](/pkg/net/http/#Server) implementation,
[`Serve`](/pkg/net/http/#Server.Serve) records in the request context
both the underlying `*Server` using the key `ServerContextKey`
and the local address on which the request was received (a
[`Addr`](/pkg/net/#Addr)) using the key `LocalAddrContextKey`.
For example, the address on which a request received is
`req.Context().Value(http.LocalAddrContextKey).(net.Addr)`.

The server's [`Serve`](/pkg/net/http/#Server.Serve) method
now only enables HTTP/2 support if the `Server.TLSConfig` field is `nil`
or includes `"h2"` in its `TLSConfig.NextProtos`.

The server implementation now
pads response codes less than 100 to three digits
as required by the protocol,
so that `w.WriteHeader(5)` uses the HTTP response
status `005`, not just `5`.

The server implementation now correctly sends only one "Transfer-Encoding" header when "chunked"
is set explicitly, following [RFC 7230](https://tools.ietf.org/html/rfc7230#section-3.3.1).

The server implementation is now stricter about rejecting requests with invalid HTTP versions.
Invalid requests claiming to be HTTP/0.x are now rejected (HTTP/0.9 was never fully supported),
and plaintext HTTP/2 requests other than the "PRI \* HTTP/2.0" upgrade request are now rejected as well.
The server continues to handle encrypted HTTP/2 requests.

In the server, a 200 status code is sent back by the timeout handler on an empty
response body, instead of sending back 0 as the status code.

In the client, the
[`Transport`](/pkg/net/http/#Transport) implementation passes the request context
to any dial operation connecting to the remote server.
If a custom dialer is needed, the new `Transport` field
`DialContext` is preferred over the existing `Dial` field,
to allow the transport to supply a context.

The
[`Transport`](/pkg/net/http/#Transport) also adds fields
`IdleConnTimeout`,
`MaxIdleConns`,
and
`MaxResponseHeaderBytes`
to help control client resources consumed
by idle or chatty servers.

A
[`Client`](/pkg/net/http/#Client)'s configured `CheckRedirect` function can now
return `ErrUseLastResponse` to indicate that the
most recent redirect response should be returned as the
result of the HTTP request.
That response is now available to the `CheckRedirect` function
as `req.Response`.

Since Go 1, the default behavior of the HTTP client is
to request server-side compression
using the `Accept-Encoding` request header
and then to decompress the response body transparently,
and this behavior is adjustable using the
[`Transport`](/pkg/net/http/#Transport)'s `DisableCompression` field.
In Go 1.7, to aid the implementation of HTTP proxies, the
[`Response`](/pkg/net/http/#Response)'s new
`Uncompressed` field reports whether
this transparent decompression took place.

[`DetectContentType`](/pkg/net/http/#DetectContentType)
adds support for a few new audio and video content types.

#### [net/http/cgi](/pkg/net/http/cgi/)

The
[`Handler`](/pkg/net/http/cgi/#Handler)
adds a new field
`Stderr`
that allows redirection of the child process's
standard error away from the host process's
standard error.

#### [net/http/httptest](/pkg/net/http/httptest/)

The new function
[`NewRequest`](/pkg/net/http/httptest/#NewRequest)
prepares a new
[`http.Request`](/pkg/net/http/#Request)
suitable for passing to an
[`http.Handler`](/pkg/net/http/#Handler) during a test.

The
[`ResponseRecorder`](/pkg/net/http/httptest/#ResponseRecorder)'s new
[`Result`](/pkg/net/http/httptest/#ResponseRecorder.Result) method
returns the recorded
[`http.Response`](/pkg/net/http/#Response).
Tests that need to check the response's headers or trailers
should call `Result` and inspect the response fields
instead of accessing
`ResponseRecorder`'s `HeaderMap` directly.

#### [net/http/httputil](/pkg/net/http/httputil/)

The
[`ReverseProxy`](/pkg/net/http/httputil/#ReverseProxy) implementation now responds with “502 Bad Gateway”
when it cannot reach a back end; in earlier releases it responded with “500 Internal Server Error.”

Both
[`ClientConn`](/pkg/net/http/httputil/#ClientConn) and
[`ServerConn`](/pkg/net/http/httputil/#ServerConn) have been documented as deprecated.
They are low-level, old, and unused by Go's current HTTP stack
and will no longer be updated.
Programs should use
[`http.Client`](/pkg/net/http/#Client),
[`http.Transport`](/pkg/net/http/#Transport),
and
[`http.Server`](/pkg/net/http/#Server)
instead.

#### [net/http/pprof](/pkg/net/http/pprof/)

The runtime trace HTTP handler, installed to handle the path `/debug/pprof/trace`,
now accepts a fractional number in its `seconds` query parameter,
allowing collection of traces for intervals smaller than one second.
This is especially useful on busy servers.

#### [net/mail](/pkg/net/mail/)

The address parser now allows unescaped UTF-8 text in addresses
following [RFC 6532](https://tools.ietf.org/html/rfc6532),
but it does not apply any normalization to the result.
For compatibility with older mail parsers,
the address encoder, namely
[`Address`](/pkg/net/mail/#Address)'s
[`String`](/pkg/net/mail/#Address.String) method,
continues to escape all UTF-8 text following [RFC 5322](https://tools.ietf.org/html/rfc5322).

The [`ParseAddress`](/pkg/net/mail/#ParseAddress)
function and
the [`AddressParser.Parse`](/pkg/net/mail/#AddressParser.Parse)
method are stricter.
They used to ignore any characters following an e-mail address, but
will now return an error for anything other than whitespace.

#### [net/url](/pkg/net/url/)

The
[`URL`](/pkg/net/url/#URL)'s
new `ForceQuery` field
records whether the URL must have a query string,
in order to distinguish URLs without query strings (like `/search`)
from URLs with empty query strings (like `/search?`).

#### [os](/pkg/os/)

[`IsExist`](/pkg/os/#IsExist) now returns true for `syscall.ENOTEMPTY`,
on systems where that error exists.

On Windows,
[`Remove`](/pkg/os/#Remove) now removes read-only files when possible,
making the implementation behave as on
non-Windows systems.

#### [os/exec](/pkg/os/exec/)

As part of the introduction of [context](#context),
the new constructor
[`CommandContext`](/pkg/os/exec/#CommandContext)
is like
[`Command`](/pkg/os/exec/#Command) but includes a context that can be used to cancel the command execution.

#### [os/user](/pkg/os/user/)

The
[`Current`](/pkg/os/user/#Current)
function is now implemented even when cgo is not available.

The new
[`Group`](/pkg/os/user/#Group) type,
along with the lookup functions
[`LookupGroup`](/pkg/os/user/#LookupGroup) and
[`LookupGroupId`](/pkg/os/user/#LookupGroupId)
and the new field `GroupIds` in the `User` struct,
provides access to system-specific user group information.

#### [reflect](/pkg/reflect/)

Although
[`Value`](/pkg/reflect/#Value)'s
[`Field`](/pkg/reflect/#Value.Field) method has always been documented to panic
if the given field number `i` is out of range, it has instead
silently returned a zero
[`Value`](/pkg/reflect/#Value).
Go 1.7 changes the method to behave as documented.

The new
[`StructOf`](/pkg/reflect/#StructOf)
function constructs a struct type at run time.
It completes the set of type constructors, joining
[`ArrayOf`](/pkg/reflect/#ArrayOf),
[`ChanOf`](/pkg/reflect/#ChanOf),
[`FuncOf`](/pkg/reflect/#FuncOf),
[`MapOf`](/pkg/reflect/#MapOf),
[`PtrTo`](/pkg/reflect/#PtrTo),
and
[`SliceOf`](/pkg/reflect/#SliceOf).

[`StructTag`](/pkg/reflect/#StructTag)'s
new method
[`Lookup`](/pkg/reflect/#StructTag.Lookup)
is like
[`Get`](/pkg/reflect/#StructTag.Get)
but distinguishes the tag not containing the given key
from the tag associating an empty string with the given key.

The
[`Method`](/pkg/reflect/#Type.Method) and
[`NumMethod`](/pkg/reflect/#Type.NumMethod)
methods of
[`Type`](/pkg/reflect/#Type) and
[`Value`](/pkg/reflect/#Value)
no longer return or count unexported methods.

#### [strings](/pkg/strings/)

In previous releases of Go, if
[`Reader`](/pkg/strings/#Reader)'s
[`Read`](/pkg/strings/#Reader.Read) method
were asked for zero bytes with no data remaining, it would
return a count of 0 and no error.
Now it returns a count of 0 and the error
[`io.EOF`](/pkg/io/#EOF).

The
[`Reader`](/pkg/strings/#Reader) type has a new method
[`Reset`](/pkg/strings/#Reader.Reset) to allow reuse of a `Reader`.

#### [time](/pkg/time/)

[`Duration`](/pkg/time/#Duration)'s
time.Duration.String method now reports the zero duration as `"0s"`, not `"0"`.
[`ParseDuration`](/pkg/time/#ParseDuration) continues to accept both forms.

The method call `time.Local.String()` now returns `"Local"` on all systems;
in earlier releases, it returned an empty string on Windows.

The time zone database in
`$GOROOT/lib/time` has been updated
to IANA release 2016d.
This fallback database is only used when the system time zone database
cannot be found, for example on Windows.
The Windows time zone abbreviation list has also been updated.

#### [syscall](/pkg/syscall/)

On Linux, the
[`SysProcAttr`](/pkg/syscall/#SysProcAttr) struct
(as used in
[`os/exec.Cmd`](/pkg/os/exec/#Cmd)'s `SysProcAttr` field)
has a new `Unshareflags` field.
If the field is nonzero, the child process created by
[`ForkExec`](/pkg/syscall/#ForkExec)
(as used in `exec.Cmd`'s `Run` method)
will call the
[_unshare_(2)](https://man7.org/linux/man-pages/man2/unshare.2.html)
system call before executing the new program.

#### [unicode](/pkg/unicode/)

The [`unicode`](/pkg/unicode/) package and associated
support throughout the system has been upgraded from version 8.0 to
[Unicode 9.0](https://www.unicode.org/versions/Unicode9.0.0/).
