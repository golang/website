---
template: false
title: Go 1.6 Release Notes
---

<!--
Edit .,s;^PKG:([a-z][A-Za-z0-9_/]+);<a href="/pkg/\1/"><code>\1</code></a>;g
Edit .,s;^([a-z][A-Za-z0-9_/]+)\.([A-Z][A-Za-z0-9_]+\.)?([A-Z][A-Za-z0-9_]+)([ .',]|$);<a href="/pkg/\1/#\2\3"><code>\3</code></a>\4;g
-->

<style>
  main ul li { margin: 0.5em 0; }
</style>

## Introduction to Go 1.6 {#introduction}

The latest Go release, version 1.6, arrives six months after 1.5.
Most of its changes are in the implementation of the language, runtime, and libraries.
There are no changes to the language specification.
As always, the release maintains the Go 1 [promise of compatibility](/doc/go1compat.html).
We expect almost all Go programs to continue to compile and run as before.

The release adds new ports to [Linux on 64-bit MIPS and Android on 32-bit x86](#ports);
defined and enforced [rules for sharing Go pointers with C](#cgo);
transparent, automatic [support for HTTP/2](#http2);
and a new mechanism for [template reuse](#template).

## Changes to the language {#language}

There are no language changes in this release.

## Ports {#ports}

Go 1.6 adds experimental ports to
Linux on 64-bit MIPS (`linux/mips64` and `linux/mips64le`).
These ports support `cgo` but only with internal linking.

Go 1.6 also adds an experimental port to Android on 32-bit x86 (`android/386`).

On FreeBSD, Go 1.6 defaults to using `clang`, not `gcc`, as the external C compiler.

On Linux on little-endian 64-bit PowerPC (`linux/ppc64le`),
Go 1.6 now supports `cgo` with external linking and
is roughly feature complete.

On NaCl, Go 1.5 required SDK version pepper-41.
Go 1.6 adds support for later SDK versions.

On 32-bit x86 systems using the `-dynlink` or `-shared` compilation modes,
the register CX is now overwritten by certain memory references and should
be avoided in hand-written assembly.
See the [assembly documentation](/doc/asm#x86) for details.

## Tools {#tools}

### Cgo {#cgo}

There is one major change to [`cgo`](/cmd/cgo/), along with one minor change.

The major change is the definition of rules for sharing Go pointers with C code,
to ensure that such C code can coexist with Go's garbage collector.
Briefly, Go and C may share memory allocated by Go
when a pointer to that memory is passed to C as part of a `cgo` call,
provided that the memory itself contains no pointers to Go-allocated memory,
and provided that C does not retain the pointer after the call returns.
These rules are checked by the runtime during program execution:
if the runtime detects a violation, it prints a diagnosis and crashes the program.
The checks can be disabled by setting the environment variable
`GODEBUG=cgocheck=0`, but note that the vast majority of
code identified by the checks is subtly incompatible with garbage collection
in one way or another.
Disabling the checks will typically only lead to more mysterious failure modes.
Fixing the code in question should be strongly preferred
over turning off the checks.
See the [`cgo` documentation](/cmd/cgo/#hdr-Passing_pointers) for more details.

The minor change is
the addition of explicit `C.complexfloat` and `C.complexdouble` types,
separate from Go's `complex64` and `complex128`.
Matching the other numeric types, C's complex types and Go's complex type are
no longer interchangeable.

### Compiler Toolchain {#compiler}

The compiler toolchain is mostly unchanged.
Internally, the most significant change is that the parser is now hand-written
instead of generated from [yacc](/cmd/yacc/).

The compiler, linker, and `go` command have a new flag `-msan`,
analogous to `-race` and only available on linux/amd64,
that enables interoperation with the [Clang MemorySanitizer](https://clang.llvm.org/docs/MemorySanitizer.html).
Such interoperation is useful mainly for testing a program containing suspect C or C++ code.

The linker has a new option `-libgcc` to set the expected location
of the C compiler support library when linking [`cgo`](/cmd/cgo/) code.
The option is only consulted when using `-linkmode=internal`,
and it may be set to `none` to disable the use of a support library.

The implementation of [build modes started in Go 1.5](/doc/go1.5#link) has been expanded to more systems.
This release adds support for the `c-shared` mode on `android/386`, `android/amd64`,
`android/arm64`, `linux/386`, and `linux/arm64`;
for the `shared` mode on `linux/386`, `linux/arm`, `linux/amd64`, and `linux/ppc64le`;
and for the new `pie` mode (generating position-independent executables) on
`android/386`, `android/amd64`, `android/arm`, `android/arm64`, `linux/386`,
`linux/amd64`, `linux/arm`, `linux/arm64`, and `linux/ppc64le`.
See the [design document](/s/execmodes) for details.

As a reminder, the linker's `-X` flag changed in Go 1.5.
In Go 1.4 and earlier, it took two arguments, as in

	-X importpath.name value

Go 1.5 added an alternative syntax using a single argument
that is itself a `name=value` pair:

	-X importpath.name=value

In Go 1.5 the old syntax was still accepted, after printing a warning
suggesting use of the new syntax instead.
Go 1.6 continues to accept the old syntax and print the warning.
Go 1.7 will remove support for the old syntax.

### Gccgo {#gccgo}

The release schedules for the GCC and Go projects do not coincide.
GCC release 5 contains the Go 1.4 version of gccgo.
The next release, GCC 6, will have the Go 1.6.1 version of gccgo.

### Go command {#go_command}

The [`go`](/cmd/go) command's basic operation
is unchanged, but there are a number of changes worth noting.

Go 1.5 introduced experimental support for vendoring,
enabled by setting the `GO15VENDOREXPERIMENT` environment variable to `1`.
Go 1.6 keeps the vendoring support, no longer considered experimental,
and enables it by default.
It can be disabled explicitly by setting
the `GO15VENDOREXPERIMENT` environment variable to `0`.
Go 1.7 will remove support for the environment variable.

The most likely problem caused by enabling vendoring by default happens
in source trees containing an existing directory named `vendor` that
does not expect to be interpreted according to new vendoring semantics.
In this case, the simplest fix is to rename the directory to anything other
than `vendor` and update any affected import paths.

For details about vendoring,
see the documentation for the [`go` command](/cmd/go/#hdr-Vendor_Directories)
and the [design document](/s/go15vendor).

There is a new build flag, `-msan`,
that compiles Go with support for the LLVM memory sanitizer.
This is intended mainly for use when linking against C or C++ code
that is being checked with the memory sanitizer.

### Go doc command {#doc_command}

Go 1.5 introduced the
[`go doc`](/cmd/go/#hdr-Show_documentation_for_package_or_symbol) command,
which allows references to packages using only the package name, as in
`go` `doc` `http`.
In the event of ambiguity, the Go 1.5 behavior was to use the package
with the lexicographically earliest import path.
In Go 1.6, ambiguity is resolved by preferring import paths with
fewer elements, breaking ties using lexicographic comparison.
An important effect of this change is that original copies of packages
are now preferred over vendored copies.
Successful searches also tend to run faster.

### Go vet command {#vet_command}

The [`go vet`](/cmd/vet) command now diagnoses
passing function or method values as arguments to `Printf`,
such as when passing `f` where `f()` was intended.

## Performance {#performance}

As always, the changes are so general and varied that precise statements
about performance are difficult to make.
Some programs may run faster, some slower.
On average the programs in the Go 1 benchmark suite run a few percent faster in Go 1.6
than they did in Go 1.5.
The garbage collector's pauses are even lower than in Go 1.5,
especially for programs using
a large amount of memory.

There have been significant optimizations bringing more than 10% improvements
to implementations of the
[`compress/bzip2`](/pkg/compress/bzip2/),
[`compress/gzip`](/pkg/compress/gzip/),
[`crypto/aes`](/pkg/crypto/aes/),
[`crypto/elliptic`](/pkg/crypto/elliptic/),
[`crypto/ecdsa`](/pkg/crypto/ecdsa/), and
[`sort`](/pkg/sort/) packages.

## Standard library {#library}

### HTTP/2 {#http2}

Go 1.6 adds transparent support in the
[`net/http`](/pkg/net/http/) package
for the new [HTTP/2 protocol](https://http2.github.io/).
Go clients and servers will automatically use HTTP/2 as appropriate when using HTTPS.
There is no exported API specific to details of the HTTP/2 protocol handling,
just as there is no exported API specific to HTTP/1.1.

Programs that must disable HTTP/2 can do so by setting
[`Transport.TLSNextProto`](/pkg/net/http/#Transport) (for clients)
or
[`Server.TLSNextProto`](/pkg/net/http/#Server) (for servers)
to a non-nil, empty map.

Programs that must adjust HTTP/2 protocol-specific details can import and use
[`golang.org/x/net/http2`](https://golang.org/x/net/http2),
in particular its
[ConfigureServer](https://godoc.org/golang.org/x/net/http2/#ConfigureServer)
and
[ConfigureTransport](https://godoc.org/golang.org/x/net/http2/#ConfigureTransport)
functions.

### Runtime {#runtime}

The runtime has added lightweight, best-effort detection of concurrent misuse of maps.
As always, if one goroutine is writing to a map, no other goroutine should be
reading or writing the map concurrently.
If the runtime detects this condition, it prints a diagnosis and crashes the program.
The best way to find out more about the problem is to run the program
under the
[race detector](/blog/race-detector),
which will more reliably identify the race
and give more detail.

For program-ending panics, the runtime now by default
prints only the stack of the running goroutine,
not all existing goroutines.
Usually only the current goroutine is relevant to a panic,
so omitting the others significantly reduces irrelevant output
in a crash message.
To see the stacks from all goroutines in crash messages, set the environment variable
`GOTRACEBACK` to `all`
or call
[`debug.SetTraceback`](/pkg/runtime/debug/#SetTraceback)
before the crash, and rerun the program.
See the [runtime documentation](/pkg/runtime/#hdr-Environment_Variables) for details.

_Updating_:
Uncaught panics intended to dump the state of the entire program,
such as when a timeout is detected or when explicitly handling a received signal,
should now call `debug.SetTraceback("all")` before panicking.
Searching for uses of
[`signal.Notify`](/pkg/os/signal/#Notify) may help identify such code.

On Windows, Go programs in Go 1.5 and earlier forced
the global Windows timer resolution to 1ms at startup
by calling `timeBeginPeriod(1)`.
Go no longer needs this for good scheduler performance,
and changing the global timer resolution caused problems on some systems,
so the call has been removed.

When using `-buildmode=c-archive` or
`-buildmode=c-shared` to build an archive or a shared
library, the handling of signals has changed.
In Go 1.5 the archive or shared library would install a signal handler
for most signals.
In Go 1.6 it will only install a signal handler for the
synchronous signals needed to handle run-time panics in Go code:
SIGBUS, SIGFPE, SIGSEGV.
See the [os/signal](/pkg/os/signal) package for more
details.

### Reflect {#reflect}

The
[`reflect`](/pkg/reflect/) package has
[resolved a long-standing incompatibility](/issue/12367)
between the gc and gccgo toolchains
regarding embedded unexported struct types containing exported fields.
Code that walks data structures using reflection, especially to implement
serialization in the spirit
of the
[`encoding/json`](/pkg/encoding/json/) and
[`encoding/xml`](/pkg/encoding/xml/) packages,
may need to be updated.

The problem arises when using reflection to walk through
an embedded unexported struct-typed field
into an exported field of that struct.
In this case, `reflect` had incorrectly reported
the embedded field as exported, by returning an empty `Field.PkgPath`.
Now it correctly reports the field as unexported
but ignores that fact when evaluating access to exported fields
contained within the struct.

_Updating_:
Typically, code that previously walked over structs and used

	f.PkgPath != ""

to exclude inaccessible fields
should now use

	f.PkgPath != "" && !f.Anonymous

For example, see the changes to the implementations of
[`encoding/json`](https://go-review.googlesource.com/#/c/14011/2/src/encoding/json/encode.go) and
[`encoding/xml`](https://go-review.googlesource.com/#/c/14012/2/src/encoding/xml/typeinfo.go).

### Sorting {#sort}

In the
[`sort`](/pkg/sort/)
package,
the implementation of
[`Sort`](/pkg/sort/#Sort)
has been rewritten to make about 10% fewer calls to the
[`Interface`](/pkg/sort/#Interface)'s
`Less` and `Swap`
methods, with a corresponding overall time savings.
The new algorithm does choose a different ordering than before
for values that compare equal (those pairs for which `Less(i,` `j)` and `Less(j,` `i)` are false).

_Updating_:
The definition of `Sort` makes no guarantee about the final order of equal values,
but the new behavior may still break programs that expect a specific order.
Such programs should either refine their `Less` implementations
to report the desired order
or should switch to
[`Stable`](/pkg/sort/#Stable),
which preserves the original input order
of equal values.

### Templates {#template}

In the
[text/template](/pkg/text/template/) package,
there are two significant new features to make writing templates easier.

First, it is now possible to [trim spaces around template actions](/pkg/text/template/#hdr-Text_and_spaces),
which can make template definitions more readable.
A minus sign at the beginning of an action says to trim space before the action,
and a minus sign at the end of an action says to trim space after the action.
For example, the template

	{{23 -}}
	   <
	{{- 45}}

formats as `23<45`.

Second, the new [`{{block}}` action](/pkg/text/template/#hdr-Actions),
combined with allowing redefinition of named templates,
provides a simple way to define pieces of a template that
can be replaced in different instantiations.
There is [an example](/pkg/text/template/#example_Template_block)
in the `text/template` package that demonstrates this new feature.

### Minor changes to the library {#minor_library_changes}

  - The [`archive/tar`](/pkg/archive/tar/) package's
    implementation corrects many bugs in rare corner cases of the file format.
    One visible change is that the
    [`Reader`](/pkg/archive/tar/#Reader) type's
    [`Read`](/pkg/archive/tar/#Reader.Read) method
    now presents the content of special file types as being empty,
    returning `io.EOF` immediately.
  - In the [`archive/zip`](/pkg/archive/zip/) package, the
    [`Reader`](/pkg/archive/zip/#Reader) type now has a
    [`RegisterDecompressor`](/pkg/archive/zip/#Reader.RegisterDecompressor) method,
    and the
    [`Writer`](/pkg/archive/zip/#Writer) type now has a
    [`RegisterCompressor`](/pkg/archive/zip/#Writer.RegisterCompressor) method,
    enabling control over compression options for individual zip files.
    These take precedence over the pre-existing global
    [`RegisterDecompressor`](/pkg/archive/zip/#RegisterDecompressor) and
    [`RegisterCompressor`](/pkg/archive/zip/#RegisterCompressor) functions.
  - The [`bufio`](/pkg/bufio/) package's
    [`Scanner`](/pkg/bufio/#Scanner) type now has a
    [`Buffer`](/pkg/bufio/#Scanner.Buffer) method,
    to specify an initial buffer and maximum buffer size to use during scanning.
    This makes it possible, when needed, to scan tokens larger than
    `MaxScanTokenSize`.
    Also for the `Scanner`, the package now defines the
    [`ErrFinalToken`](/pkg/bufio/#ErrFinalToken) error value, for use by
    [split functions](/pkg/bufio/#SplitFunc) to abort processing or to return a final empty token.
  - The [`compress/flate`](/pkg/compress/flate/) package
    has deprecated its
    [`ReadError`](/pkg/compress/flate/#ReadError) and
    [`WriteError`](/pkg/compress/flate/#WriteError) error implementations.
    In Go 1.5 they were only rarely returned when an error was encountered;
    now they are never returned, although they remain defined for compatibility.
  - The [`compress/flate`](/pkg/compress/flate/),
    [`compress/gzip`](/pkg/compress/gzip/), and
    [`compress/zlib`](/pkg/compress/zlib/) packages
    now report
    [`io.ErrUnexpectedEOF`](/pkg/io/#ErrUnexpectedEOF) for truncated input streams, instead of
    [`io.EOF`](/pkg/io/#EOF).
  - The [`crypto/cipher`](/pkg/crypto/cipher/) package now
    overwrites the destination buffer in the event of a GCM decryption failure.
    This is to allow the AESNI code to avoid using a temporary buffer.
  - The [`crypto/tls`](/pkg/crypto/tls/) package
    has a variety of minor changes.
    It now allows
    [`Listen`](/pkg/crypto/tls/#Listen)
    to succeed when the
    [`Config`](/pkg/crypto/tls/#Config)
    has a nil `Certificates`, as long as the `GetCertificate` callback is set,
    it adds support for RSA with AES-GCM cipher suites,
    and
    it adds a
    [`RecordHeaderError`](/pkg/crypto/tls/#RecordHeaderError)
    to allow clients (in particular, the [`net/http`](/pkg/net/http/) package)
    to report a better error when attempting a TLS connection to a non-TLS server.
  - The [`crypto/x509`](/pkg/crypto/x509/) package
    now permits certificates to contain negative serial numbers
    (technically an error, but unfortunately common in practice),
    and it defines a new
    [`InsecureAlgorithmError`](/pkg/crypto/x509/#InsecureAlgorithmError)
    to give a better error message when rejecting a certificate
    signed with an insecure algorithm like MD5.
  - The [`debug/dwarf`](/pkg/debug/dwarf) and
    [`debug/elf`](/pkg/debug/elf/) packages
    together add support for compressed DWARF sections.
    User code needs no updating: the sections are decompressed automatically when read.
  - The [`debug/elf`](/pkg/debug/elf/) package
    adds support for general compressed ELF sections.
    User code needs no updating: the sections are decompressed automatically when read.
    However, compressed
    [`Sections`](/pkg/debug/elf/#Section) do not support random access:
    they have a nil `ReaderAt` field.
  - The [`encoding/asn1`](/pkg/encoding/asn1/) package
    now exports
    [tag and class constants](/pkg/encoding/asn1/#pkg-constants)
    useful for advanced parsing of ASN.1 structures.
  - Also in the [`encoding/asn1`](/pkg/encoding/asn1/) package,
    [`Unmarshal`](/pkg/encoding/asn1/#Unmarshal) now rejects various non-standard integer and length encodings.
  - The [`encoding/base64`](/pkg/encoding/base64) package's
    [`Decoder`](/pkg/encoding/base64/#Decoder) has been fixed
    to process the final bytes of its input. Previously it processed as many four-byte tokens as
    possible but ignored the remainder, up to three bytes.
    The `Decoder` therefore now handles inputs in unpadded encodings (like
    [RawURLEncoding](/pkg/encoding/base64/#RawURLEncoding)) correctly,
    but it also rejects inputs in padded encodings that are truncated or end with invalid bytes,
    such as trailing spaces.
  - The [`encoding/json`](/pkg/encoding/json/) package
    now checks the syntax of a
    [`Number`](/pkg/encoding/json/#Number)
    before marshaling it, requiring that it conforms to the JSON specification for numeric values.
    As in previous releases, the zero `Number` (an empty string) is marshaled as a literal 0 (zero).
  - The [`encoding/xml`](/pkg/encoding/xml/) package's
    [`Marshal`](/pkg/encoding/xml/#Marshal)
    function now supports a `cdata` attribute, such as `chardata`
    but encoding its argument in one or more `<![CDATA[ ... ]]>` tags.
  - Also in the [`encoding/xml`](/pkg/encoding/xml/) package,
    [`Decoder`](/pkg/encoding/xml/#Decoder)'s
    [`Token`](/pkg/encoding/xml/#Decoder.Token) method
    now reports an error when encountering EOF before seeing all open tags closed,
    consistent with its general requirement that tags in the input be properly matched.
    To avoid that requirement, use
    [`RawToken`](/pkg/encoding/xml/#Decoder.RawToken).
  - The [`fmt`](/pkg/fmt/) package now allows
    any integer type as an argument to
    [`Printf`](/pkg/fmt/#Printf)'s `*` width and precision specification.
    In previous releases, the argument to `*` was required to have type `int`.
  - Also in the [`fmt`](/pkg/fmt/) package,
    [`Scanf`](/pkg/fmt/#Scanf) can now scan hexadecimal strings using %X, as an alias for %x.
    Both formats accept any mix of upper- and lower-case hexadecimal.
  - The [`image`](/pkg/image/)
    and
    [`image/color`](/pkg/image/color/) packages
    add
    [`NYCbCrA`](/pkg/image/#NYCbCrA)
    and
    [`NYCbCrA`](/pkg/image/color/#NYCbCrA)
    types, to support Y'CbCr images with non-premultiplied alpha.
  - The [`io`](/pkg/io/) package's
    [`MultiWriter`](/pkg/io/#MultiWriter)
    implementation now implements a `WriteString` method,
    for use by
    [`WriteString`](/pkg/io/#WriteString).
  - In the [`math/big`](/pkg/math/big/) package,
    [`Int`](/pkg/math/big/#Int) adds
    [`Append`](/pkg/math/big/#Int.Append)
    and
    [`Text`](/pkg/math/big/#Int.Text)
    methods to give more control over printing.
  - Also in the [`math/big`](/pkg/math/big/) package,
    [`Float`](/pkg/math/big/#Float) now implements
    [`encoding.TextMarshaler`](/pkg/encoding/#TextMarshaler) and
    [`encoding.TextUnmarshaler`](/pkg/encoding/#TextUnmarshaler),
    allowing it to be serialized in a natural form by the
    [`encoding/json`](/pkg/encoding/json/) and
    [`encoding/xml`](/pkg/encoding/xml/) packages.
  - Also in the [`math/big`](/pkg/math/big/) package,
    [`Float`](/pkg/math/big/#Float)'s
    [`Append`](/pkg/math/big/#Float.Append) method now supports the special precision argument -1.
    As in
    [`strconv.ParseFloat`](/pkg/strconv/#ParseFloat),
    precision -1 means to use the smallest number of digits necessary such that
    [`Parse`](/pkg/math/big/#Float.Parse)
    reading the result into a `Float` of the same precision
    will yield the original value.
  - The [`math/rand`](/pkg/math/rand/) package
    adds a
    [`Read`](/pkg/math/rand/#Read)
    function, and likewise
    [`Rand`](/pkg/math/rand/#Rand) adds a
    [`Read`](/pkg/math/rand/#Rand.Read) method.
    These make it easier to generate pseudorandom test data.
    Note that, like the rest of the package,
    these should not be used in cryptographic settings;
    for such purposes, use the [`crypto/rand`](/pkg/crypto/rand/) package instead.
  - The [`net`](/pkg/net/) package's
    [`ParseMAC`](/pkg/net/#ParseMAC) function now accepts 20-byte IP-over-InfiniBand (IPoIB) link-layer addresses.
  - Also in the [`net`](/pkg/net/) package,
    there have been a few changes to DNS lookups.
    First, the
    [`DNSError`](/pkg/net/#DNSError) error implementation now implements
    [`Error`](/pkg/net/#Error),
    and in particular its new
    [`IsTemporary`](/pkg/net/#DNSError.IsTemporary)
    method returns true for DNS server errors.
    Second, DNS lookup functions such as
    [`LookupAddr`](/pkg/net/#LookupAddr)
    now return rooted domain names (with a trailing dot)
    on Plan 9 and Windows, to match the behavior of Go on Unix systems.
  - The [`net/http`](/pkg/net/http/) package has
    a number of minor additions beyond the HTTP/2 support already discussed.
    First, the
    [`FileServer`](/pkg/net/http/#FileServer) now sorts its generated directory listings by file name.
    Second, the
    [`ServeFile`](/pkg/net/http/#ServeFile) function now refuses to serve a result
    if the request's URL path contains “..” (dot-dot) as a path element.
    Programs should typically use `FileServer` and
    [`Dir`](/pkg/net/http/#Dir)
    instead of calling `ServeFile` directly.
    Programs that need to serve file content in response to requests for URLs containing dot-dot can
    still call [`ServeContent`](/pkg/net/http/#ServeContent).
    Third, the
    [`Client`](/pkg/net/http/#Client) now allows user code to set the
    `Expect:` `100-continue` header (see
    [`Transport.ExpectContinueTimeout`](/pkg/net/http/#Transport)).
    Fourth, there are
    [five new error codes](/pkg/net/http/#pkg-constants):
    `StatusPreconditionRequired` (428),
    `StatusTooManyRequests` (429),
    `StatusRequestHeaderFieldsTooLarge` (431), and
    `StatusNetworkAuthenticationRequired` (511) from RFC 6585,
    as well as the recently-approved
    `StatusUnavailableForLegalReasons` (451).
    Fifth, the implementation and documentation of
    [`CloseNotifier`](/pkg/net/http/#CloseNotifier)
    has been substantially changed.
    The [`Hijacker`](/pkg/net/http/#Hijacker)
    interface now works correctly on connections that have previously
    been used with `CloseNotifier`.
    The documentation now describes when `CloseNotifier`
    is expected to work.
  - Also in the [`net/http`](/pkg/net/http/) package,
    there are a few changes related to the handling of a
    [`Request`](/pkg/net/http/#Request) data structure with its `Method` field set to the empty string.
    An empty `Method` field has always been documented as an alias for `"GET"`
    and it remains so.
    However, Go 1.6 fixes a few routines that did not treat an empty
    `Method` the same as an explicit `"GET"`.
    Most notably, in previous releases
    [`Client`](/pkg/net/http/#Client) followed redirects only with
    `Method` set explicitly to `"GET"`;
    in Go 1.6 `Client` also follows redirects for the empty `Method`.
    Finally,
    [`NewRequest`](/pkg/net/http/#NewRequest) accepts a `method` argument that has not been
    documented as allowed to be empty.
    In past releases, passing an empty `method` argument resulted
    in a `Request` with an empty `Method` field.
    In Go 1.6, the resulting `Request` always has an initialized
    `Method` field: if its argument is an empty string, `NewRequest`
    sets the `Method` field in the returned `Request` to `"GET"`.
  - The [`net/http/httptest`](/pkg/net/http/httptest/) package's
    [`ResponseRecorder`](/pkg/net/http/httptest/#ResponseRecorder) now initializes a default Content-Type header
    using the same content-sniffing algorithm as in
    [`http.Server`](/pkg/net/http/#Server).
  - The [`net/url`](/pkg/net/url/) package's
    [`Parse`](/pkg/net/url/#Parse) is now stricter and more spec-compliant regarding the parsing
    of host names.
    For example, spaces in the host name are no longer accepted.
  - Also in the [`net/url`](/pkg/net/url/) package,
    the [`Error`](/pkg/net/url/#Error) type now implements
    [`net.Error`](/pkg/net/#Error).
  - The [`os`](/pkg/os/) package's
    [`IsExist`](/pkg/os/#IsExist),
    [`IsNotExist`](/pkg/os/#IsNotExist),
    and
    [`IsPermission`](/pkg/os/#IsPermission)
    now return correct results when inquiring about an
    [`SyscallError`](/pkg/os/#SyscallError).
  - On Unix-like systems, when a write
    to [`os.Stdout`
    or `os.Stderr`](/pkg/os/#pkg-variables) (more precisely, an `os.File`
    opened for file descriptor 1 or 2) fails due to a broken pipe error,
    the program will raise a `SIGPIPE` signal.
    By default this will cause the program to exit; this may be changed by
    calling the
    [`os/signal`](/pkg/os/signal)
    [`Notify`](/pkg/os/signal/#Notify) function
    for `syscall.SIGPIPE`.
    A write to a broken pipe on a file descriptor other 1 or 2 will simply
    return `syscall.EPIPE` (possibly wrapped in
    [`os.PathError`](/pkg/os#PathError)
    and/or [`os.SyscallError`](/pkg/os#SyscallError))
    to the caller.
    The old behavior of raising an uncatchable `SIGPIPE` signal
    after 10 consecutive writes to a broken pipe no longer occurs.
  - In the [`os/exec`](/pkg/os/exec/) package,
    [`Cmd`](/pkg/os/exec/#Cmd)'s
    [`Output`](/pkg/os/exec/#Cmd.Output) method continues to return an
    [`ExitError`](/pkg/os/exec/#ExitError) when a command exits with an unsuccessful status.
    If standard error would otherwise have been discarded,
    the returned `ExitError` now holds a prefix and suffix
    (currently 32 kB) of the failed command's standard error output,
    for debugging or for inclusion in error messages.
    The `ExitError`'s
    [`String`](/pkg/os/exec/#ExitError.String)
    method does not show the captured standard error;
    programs must retrieve it from the data structure
    separately.
  - On Windows, the [`path/filepath`](/pkg/path/filepath/) package's
    [`Join`](/pkg/path/filepath/#Join) function now correctly handles the case when the base is a relative drive path.
    For example, ``Join(`c:`,`` `` `a`) `` now
    returns `` `c:a` `` instead of `` `c:\a` `` as in past releases.
    This may affect code that expects the incorrect result.
  - In the [`regexp`](/pkg/regexp/) package,
    the
    [`Regexp`](/pkg/regexp/#Regexp) type has always been safe for use by
    concurrent goroutines.
    It uses a [`sync.Mutex`](/pkg/sync/#Mutex) to protect
    a cache of scratch spaces used during regular expression searches.
    Some high-concurrency servers using the same `Regexp` from many goroutines
    have seen degraded performance due to contention on that mutex.
    To help such servers, `Regexp` now has a
    [`Copy`](/pkg/regexp/#Regexp.Copy) method,
    which makes a copy of a `Regexp` that shares most of the structure
    of the original but has its own scratch space cache.
    Two goroutines can use different copies of a `Regexp`
    without mutex contention.
    A copy does have additional space overhead, so `Copy`
    should only be used when contention has been observed.
  - The [`strconv`](/pkg/strconv/) package adds
    [`IsGraphic`](/pkg/strconv/#IsGraphic),
    similar to [`IsPrint`](/pkg/strconv/#IsPrint).
    It also adds
    [`QuoteToGraphic`](/pkg/strconv/#QuoteToGraphic),
    [`QuoteRuneToGraphic`](/pkg/strconv/#QuoteRuneToGraphic),
    [`AppendQuoteToGraphic`](/pkg/strconv/#AppendQuoteToGraphic),
    and
    [`AppendQuoteRuneToGraphic`](/pkg/strconv/#AppendQuoteRuneToGraphic),
    analogous to
    [`QuoteToASCII`](/pkg/strconv/#QuoteToASCII),
    [`QuoteRuneToASCII`](/pkg/strconv/#QuoteRuneToASCII),
    and so on.
    The `ASCII` family escapes all space characters except ASCII space (U+0020).
    In contrast, the `Graphic` family does not escape any Unicode space characters (category Zs).
  - In the [`testing`](/pkg/testing/) package,
    when a test calls
    [t.Parallel](/pkg/testing/#T.Parallel),
    that test is paused until all non-parallel tests complete, and then
    that test continues execution with all other parallel tests.
    Go 1.6 changes the time reported for such a test:
    previously the time counted only the parallel execution,
    but now it also counts the time from the start of testing
    until the call to `t.Parallel`.
  - The [`text/template`](/pkg/text/template/) package
    contains two minor changes, in addition to the [major changes](#template)
    described above.
    First, it adds a new
    [`ExecError`](/pkg/text/template/#ExecError) type
    returned for any error during
    [`Execute`](/pkg/text/template/#Template.Execute)
    that does not originate in a `Write` to the underlying writer.
    Callers can distinguish template usage errors from I/O errors by checking for
    `ExecError`.
    Second, the
    [`Funcs`](/pkg/text/template/#Template.Funcs) method
    now checks that the names used as keys in the
    [`FuncMap`](/pkg/text/template/#FuncMap)
    are identifiers that can appear in a template function invocation.
    If not, `Funcs` panics.
  - The [`time`](/pkg/time/) package's
    [`Parse`](/pkg/time/#Parse) function has always rejected any day of month larger than 31,
    such as January 32.
    In Go 1.6, `Parse` now also rejects February 29 in non-leap years,
    February 30, February 31, April 31, June 31, September 31, and November 31.
