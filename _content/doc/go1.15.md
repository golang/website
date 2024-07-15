---
template: false
title: Go 1.15 Release Notes
---

<!--
NOTE: In this document and others in this directory, the convention is to
set fixed-width phrases with non-fixed-width spaces, as in
`hello` `world`.
Do not send CLs removing the interior tags from such phrases.
-->

<style>
  main ul li { margin: 0.5em 0; }
</style>

## Introduction to Go 1.15 {#introduction}

The latest Go release, version 1.15, arrives six months after [Go 1.14](go1.14).
Most of its changes are in the implementation of the toolchain, runtime, and libraries.
As always, the release maintains the Go 1 [promise of compatibility](/doc/go1compat.html).
We expect almost all Go programs to continue to compile and run as before.

Go 1.15 includes [substantial improvements to the linker](#linker),
improves [allocation for small objects at high core counts](#runtime), and
deprecates [X.509 CommonName](#commonname).
`GOPROXY` now supports skipping proxies that return errors and
a new [embedded tzdata package](#time/tzdata) has been added.

## Changes to the language {#language}

There are no changes to the language.

## Ports {#ports}

### Darwin {#darwin}

As [announced](go1.14#darwin) in the Go 1.14 release
notes, Go 1.15 requires macOS 10.12 Sierra or later; support for
previous versions has been discontinued.

<!-- golang.org/issue/37610, golang.org/issue/37611, CL 227582, and CL 227198  -->
As [announced](/doc/go1.14#darwin) in the Go 1.14 release
notes, Go 1.15 drops support for 32-bit binaries on macOS, iOS,
iPadOS, watchOS, and tvOS (the `darwin/386`
and `darwin/arm` ports). Go continues to support the
64-bit `darwin/amd64` and `darwin/arm64` ports.

### Windows {#windows}

<!-- CL 214397 and CL 230217 -->
Go now generates Windows ASLR executables when `-buildmode=pie`
cmd/link flag is provided. Go command uses `-buildmode=pie`
by default on Windows.

<!-- CL 227003 -->
The `-race` and `-msan` flags now always
enable `-d=checkptr`, which checks uses
of `unsafe.Pointer`. This was previously the case on all
OSes except Windows.

<!-- CL 211139 -->
Go-built DLLs no longer cause the process to exit when it receives a
signal (such as Ctrl-C at a terminal).

### Android {#android}

<!-- CL 235017, golang.org/issue/38838 -->
When linking binaries for Android, Go 1.15 explicitly selects
the `lld` linker available in recent versions of the NDK.
The `lld` linker avoids crashes on some devices, and is
planned to become the default NDK linker in a future NDK version.

### OpenBSD {#openbsd}

<!-- CL 234381 -->
Go 1.15 adds support for OpenBSD 6.7 on `GOARCH=arm`
and `GOARCH=arm64`. Previous versions of Go already
supported OpenBSD 6.7 on `GOARCH=386`
and `GOARCH=amd64`.

### RISC-V {#riscv}

<!-- CL 226400, CL 226206, and others -->
There has been progress in improving the stability and performance
of the 64-bit RISC-V port on Linux (`GOOS=linux`,
`GOARCH=riscv64`). It also now supports asynchronous
preemption.

### 386 {#386}

<!-- golang.org/issue/40255 -->
Go 1.15 is the last release to support x87-only floating-point
hardware (`GO386=387`). Future releases will require at
least SSE2 support on 386, raising Go's
minimum `GOARCH=386` requirement to the Intel Pentium 4
(released in 2000) or AMD Opteron/Athlon 64 (released in 2003).

## Tools {#tools}

### Go command {#go-command}

<!-- golang.org/issue/37367 -->
The `GOPROXY` environment variable now supports skipping proxies
that return errors. Proxy URLs may now be separated with either commas
(`,`) or pipe characters (`|`). If a proxy URL is
followed by a comma, the `go` command will only try the next proxy
in the list after a 404 or 410 HTTP response. If a proxy URL is followed by a
pipe character, the `go` command will try the next proxy in the
list after any error. Note that the default value of `GOPROXY`
remains `https://proxy.golang.org,direct`, which does not fall
back to `direct` in case of errors.

#### `go` `test` {#go-test}

<!-- https://golang.org/issue/36134 -->
Changing the `-timeout` flag now invalidates cached test results. A
cached result for a test run with a long timeout will no longer count as
passing when `go` `test` is re-invoked with a short one.

#### Flag parsing {#go-flag-parsing}

<!-- https://golang.org/cl/211358 -->
Various flag parsing issues in `go` `test` and
`go` `vet` have been fixed. Notably, flags specified
in `GOFLAGS` are handled more consistently, and
the `-outputdir` flag now interprets relative paths relative to the
working directory of the `go` command (rather than the working
directory of each individual test).

#### Module cache {#module-cache}

<!-- https://golang.org/cl/219538 -->
The location of the module cache may now be set with
the `GOMODCACHE` environment variable. The default value of
`GOMODCACHE` is `GOPATH[0]/pkg/mod`, the location of the
module cache before this change.

<!-- https://golang.org/cl/221157 -->
A workaround is now available for Windows "Access is denied" errors in
`go` commands that access the module cache, caused by external
programs concurrently scanning the file system (see
[issue #36568](/issue/36568)). The workaround is
not enabled by default because it is not safe to use when Go versions lower
than 1.14.2 and 1.13.10 are running concurrently with the same module cache.
It can be enabled by explicitly setting the environment variable
`GODEBUG=modcacheunzipinplace=1`.

### Vet {#vet}

#### New warning for string(x) {#vet-string-int}

<!-- CL 212919, 232660 -->
The vet tool now warns about conversions of the
form `string(x)` where `x` has an integer type
other than `rune` or `byte`.
Experience with Go has shown that many conversions of this form
erroneously assume that `string(x)` evaluates to the
string representation of the integer `x`.
It actually evaluates to a string containing the UTF-8 encoding of
the value of `x`.
For example, `string(9786)` does not evaluate to the
string `"9786"`; it evaluates to the
string `"\xe2\x98\xba"`, or `"☺"`.

Code that is using `string(x)` correctly can be rewritten
to `string(rune(x))`.
Or, in some cases, calling `utf8.EncodeRune(buf, x)` with
a suitable byte slice `buf` may be the right solution.
Other code should most likely use `strconv.Itoa`
or `fmt.Sprint`.

This new vet check is enabled by default when
using `go` `test`.

We are considering prohibiting the conversion in a future release of Go.
That is, the language would change to only
permit `string(x)` for integer `x` when the
type of `x` is `rune` or `byte`.
Such a language change would not be backward compatible.
We are using this vet check as a first trial step toward changing
the language.

#### New warning for impossible interface conversions {#vet-impossible-interface}

<!-- CL 218779, 232660 -->
The vet tool now warns about type assertions from one interface type
to another interface type when the type assertion will always fail.
This will happen if both interface types implement a method with the
same name but with a different type signature.

There is no reason to write a type assertion that always fails, so
any code that triggers this vet check should be rewritten.

This new vet check is enabled by default when
using `go` `test`.

We are considering prohibiting impossible interface type assertions
in a future release of Go.
Such a language change would not be backward compatible.
We are using this vet check as a first trial step toward changing
the language.

## Runtime {#runtime}

<!-- CL 221779 -->
If `panic` is invoked with a value whose type is derived from any
of: `bool`, `complex64`, `complex128`, `float32`, `float64`,
`int`, `int8`, `int16`, `int32`, `int64`, `string`,
`uint`, `uint8`, `uint16`, `uint32`, `uint64`, `uintptr`,
then the value will be printed, instead of just its address.
Previously, this was only true for values of exactly these types.

<!-- CL 228900 -->
On a Unix system, if the `kill` command
or `kill` system call is used to send
a `SIGSEGV`, `SIGBUS`,
or `SIGFPE` signal to a Go program, and if the signal
is not being handled via
[`os/signal.Notify`](/pkg/os/signal/#Notify),
the Go program will now reliably crash with a stack trace.
In earlier releases the behavior was unpredictable.

<!-- CL 221182, CL 229998 -->
Allocation of small objects now performs much better at high core
counts, and has lower worst-case latency.

<!-- CL 216401 -->
Converting a small integer value into an interface value no longer
causes allocation.

<!-- CL 216818 -->
Non-blocking receives on closed channels now perform as well as
non-blocking receives on open channels.

## Compiler {#compiler}

<!-- CL 229578 -->
Package `unsafe`'s [safety
rules](/pkg/unsafe/#Pointer) allow converting an `unsafe.Pointer`
into `uintptr` when calling certain
functions. Previously, in some cases, the compiler allowed multiple
chained conversions (for example, `syscall.Syscall(…,`
`uintptr(uintptr(ptr)),` `…)`). The compiler
now requires exactly one conversion. Code that used multiple
conversions should be updated to satisfy the safety rules.

<!-- CL 230544, CL 231397 -->
Go 1.15 reduces typical binary sizes by around 5% compared to Go
1.14 by eliminating certain types of GC metadata and more
aggressively eliminating unused type metadata.

<!-- CL 219357, CL 231600 -->
The toolchain now mitigates
[Intel
CPU erratum SKX102](https://www.intel.com/content/www/us/en/support/articles/000055650/processors.html) on `GOARCH=amd64` by aligning
functions to 32 byte boundaries and padding jump instructions. While
this padding increases binary sizes, this is more than made up for
by the binary size improvements mentioned above.

<!-- CL 222661 -->
Go 1.15 adds a `-spectre` flag to both the
compiler and the assembler, to allow enabling Spectre mitigations.
These should almost never be needed and are provided mainly as a
“defense in depth” mechanism.
See the [Spectre wiki page](/wiki/Spectre) for details.

<!-- CL 228578 -->
The compiler now rejects `//go:` compiler directives that
have no meaning for the declaration they are applied to with a
"misplaced compiler directive" error. Such misapplied directives
were broken before, but were silently ignored by the compiler.

<!-- CL 206658, CL 205066 -->
The compiler's `-json` optimization logging now reports
large (>= 128 byte) copies and includes explanations of escape
analysis decisions.

## Linker {#linker}

This release includes substantial improvements to the Go linker,
which reduce linker resource usage (both time and memory) and
improve code robustness/maintainability.

For a representative set of large Go programs, linking is 20% faster
and requires 30% less memory on average, for `ELF`-based
OSes (Linux, FreeBSD, NetBSD, OpenBSD, Dragonfly, and Solaris)
running on `amd64` architectures, with more modest
improvements for other architecture/OS combinations.

The key contributors to better linker performance are a newly
redesigned object file format, and a revamping of internal
phases to increase concurrency (for example, applying relocations to
symbols in parallel). Object files in Go 1.15 are slightly larger
than their 1.14 equivalents.

These changes are part of a multi-release project
to [modernize the Go
linker](/s/better-linker), meaning that there will be additional linker
improvements expected in future releases.

<!-- CL 207877 -->
The linker now defaults to internal linking mode
for `-buildmode=pie` on
`linux/amd64` and `linux/arm64`, so these
configurations no longer require a C linker. External linking
mode (which was the default in Go 1.14 for
`-buildmode=pie`) can still be requested with
`-ldflags=-linkmode=external` flag.

## Objdump {#objdump}

<!-- CL 225459 -->
The [objdump](/cmd/objdump/) tool now supports
disassembling in GNU assembler syntax with the `-gnu`
flag.

## Standard library {#library}

### New embedded tzdata package {#time_tzdata}

<!-- CL 224588 -->
Go 1.15 includes a new package,
[`time/tzdata`](/pkg/time/tzdata/),
that permits embedding the timezone database into a program.
Importing this package (as `import _ "time/tzdata"`)
permits the program to find timezone information even if the
timezone database is not available on the local system.
You can also embed the timezone database by building
with `-tags timetzdata`.
Either approach increases the size of the program by about 800 KB.

### Cgo {#cgo}

<!-- CL 235817 -->
Go 1.15 will translate the C type `EGLConfig` to the
Go type `uintptr`. This change is similar to how Go
1.12 and newer treats `EGLDisplay`, Darwin's CoreFoundation and
Java's JNI types. See the [cgo
documentation](/cmd/cgo/#hdr-Special_cases) for more information.

<!-- CL 250940 -->
In Go 1.15.3 and later, cgo will not permit Go code to allocate an
undefined struct type (a C struct defined as just `struct
  S;` or similar) on the stack or heap.
Go code will only be permitted to use pointers to those types.
Allocating an instance of such a struct and passing a pointer, or a
full struct value, to C code was always unsafe and unlikely to work
correctly; it is now forbidden.
The fix is to either rewrite the Go code to use only pointers, or to
ensure that the Go code sees the full definition of the struct by
including the appropriate C header file.

### X.509 CommonName deprecation {#commonname}

<!-- CL 231379 -->
The deprecated, legacy behavior of treating the `CommonName`
field on X.509 certificates as a host name when no Subject Alternative Names
are present is now disabled by default. It can be temporarily re-enabled by
adding the value `x509ignoreCN=0` to the `GODEBUG`
environment variable.

Note that if the `CommonName` is an invalid host name, it's always
ignored, regardless of `GODEBUG` settings. Invalid names include
those with any characters other than letters, digits, hyphens and underscores,
and those with empty labels or trailing dots.

### Minor changes to the library {#minor_library_changes}

As always, there are various minor changes and updates to the library,
made with the Go 1 [promise of compatibility](/doc/go1compat)
in mind.

#### [bufio](/pkg/bufio/)

<!-- CL 225357, CL 225557 -->
When a [`Scanner`](/pkg/bufio/#Scanner) is
used with an invalid
[`io.Reader`](/pkg/io/#Reader) that
incorrectly returns a negative number from `Read`,
the `Scanner` will no longer panic, but will instead
return the new error
[`ErrBadReadCount`](/pkg/bufio/#ErrBadReadCount).

<!-- bufio -->

#### [context](/pkg/context/)

<!-- CL 223777 -->
Creating a derived `Context` using a nil parent is now explicitly
disallowed. Any attempt to do so with the
[`WithValue`](/pkg/context/#WithValue),
[`WithDeadline`](/pkg/context/#WithDeadline), or
[`WithCancel`](/pkg/context/#WithCancel) functions
will cause a panic.

<!-- context -->

#### [crypto](/pkg/crypto/)

<!-- CL 231417, CL 225460 -->
The `PrivateKey` and `PublicKey` types in the
[`crypto/rsa`](/pkg/crypto/rsa/),
[`crypto/ecdsa`](/pkg/crypto/ecdsa/), and
[`crypto/ed25519`](/pkg/crypto/ed25519/) packages
now have an `Equal` method to compare keys for equivalence
or to make type-safe interfaces for public keys. The method signature
is compatible with
[`go-cmp`'s
definition of equality](https://pkg.go.dev/github.com/google/go-cmp/cmp#Equal).

<!-- CL 224937 -->
[`Hash`](/pkg/crypto/#Hash) now implements
[`fmt.Stringer`](/pkg/fmt/#Stringer).

<!-- crypto -->

#### [crypto/ecdsa](/pkg/crypto/ecdsa/)

<!-- CL 217940 -->
The new [`SignASN1`](/pkg/crypto/ecdsa/#SignASN1)
and [`VerifyASN1`](/pkg/crypto/ecdsa/#VerifyASN1)
functions allow generating and verifying ECDSA signatures in the standard
ASN.1 DER encoding.

<!-- crypto/ecdsa -->

#### [crypto/elliptic](/pkg/crypto/elliptic/)

<!-- CL 202819 -->
The new [`MarshalCompressed`](/pkg/crypto/elliptic/#MarshalCompressed)
and [`UnmarshalCompressed`](/pkg/crypto/elliptic/#UnmarshalCompressed)
functions allow encoding and decoding NIST elliptic curve points in compressed format.

<!-- crypto/elliptic -->

#### [crypto/rsa](/pkg/crypto/rsa/)

<!-- CL 226203 -->
[`VerifyPKCS1v15`](/pkg/crypto/rsa/#VerifyPKCS1v15)
now rejects invalid short signatures with missing leading zeroes, according to RFC 8017.

<!-- crypto/rsa -->

#### [crypto/tls](/pkg/crypto/tls/)

<!-- CL 214977 -->
The new
[`Dialer`](/pkg/crypto/tls/#Dialer)
type and its
[`DialContext`](/pkg/crypto/tls/#Dialer.DialContext)
method permit using a context to both connect and handshake with a TLS server.

<!-- CL 229122 -->
The new
[`VerifyConnection`](/pkg/crypto/tls/#Config.VerifyConnection)
callback on the [`Config`](/pkg/crypto/tls/#Config) type
allows custom verification logic for every connection. It has access to the
[`ConnectionState`](/pkg/crypto/tls/#ConnectionState)
which includes peer certificates, SCTs, and stapled OCSP responses.

<!-- CL 230679 -->
Auto-generated session ticket keys are now automatically rotated every 24 hours,
with a lifetime of 7 days, to limit their impact on forward secrecy.

<!-- CL 231317 -->
Session ticket lifetimes in TLS 1.2 and earlier, where the session keys
are reused for resumed connections, are now limited to 7 days, also to
limit their impact on forward secrecy.

<!-- CL 231038 -->
The client-side downgrade protection checks specified in RFC 8446 are now
enforced. This has the potential to cause connection errors for clients
encountering middleboxes that behave like unauthorized downgrade attacks.

<!-- CL 208226 -->
[`SignatureScheme`](/pkg/crypto/tls/#SignatureScheme),
[`CurveID`](/pkg/crypto/tls/#CurveID), and
[`ClientAuthType`](/pkg/crypto/tls/#ClientAuthType)
now implement [`fmt.Stringer`](/pkg/fmt/#Stringer).

<!-- CL 236737 -->
The [`ConnectionState`](/pkg/crypto/tls/#ConnectionState)
fields `OCSPResponse` and `SignedCertificateTimestamps`
are now repopulated on client-side resumed connections.

<!-- CL 227840 -->
[`tls.Conn`](/pkg/crypto/tls/#Conn)
now returns an opaque error on permanently broken connections, wrapping
the temporary
[`net.Error`](/pkg/net/http/#Error). To access the
original `net.Error`, use
[`errors.As`](/pkg/errors/#As) (or
[`errors.Unwrap`](/pkg/errors/#Unwrap)) instead of a
type assertion.

<!-- crypto/tls -->

#### [crypto/x509](/pkg/crypto/x509/)

<!-- CL 231378, CL 231380, CL 231381 -->
If either the name on the certificate or the name being verified (with
[`VerifyOptions.DNSName`](/pkg/crypto/x509/#VerifyOptions.DNSName)
or [`VerifyHostname`](/pkg/crypto/x509/#Certificate.VerifyHostname))
are invalid, they will now be compared case-insensitively without further
processing (without honoring wildcards or stripping trailing dots).
Invalid names include those with any characters other than letters,
digits, hyphens and underscores, those with empty labels, and names on
certificates with trailing dots.

<!-- CL 217298 -->
The new [`CreateRevocationList`](/pkg/crypto/x509/#CreateRevocationList)
function and [`RevocationList`](/pkg/crypto/x509/#RevocationList) type
allow creating RFC 5280-compliant X.509 v2 Certificate Revocation Lists.

<!-- CL 227098 -->
[`CreateCertificate`](/pkg/crypto/x509/#CreateCertificate)
now automatically generates the `SubjectKeyId` if the template
is a CA and doesn't explicitly specify one.

<!-- CL 228777 -->
[`CreateCertificate`](/pkg/crypto/x509/#CreateCertificate)
now returns an error if the template specifies `MaxPathLen` but is not a CA.

<!-- CL 205237 -->
On Unix systems other than macOS, the `SSL_CERT_DIR`
environment variable can now be a colon-separated list.

<!-- CL 227037 -->
On macOS, binaries are now always linked against
`Security.framework` to extract the system trust roots,
regardless of whether cgo is available. The resulting behavior should be
more consistent with the OS verifier.

<!-- crypto/x509 -->

#### [crypto/x509/pkix](/pkg/crypto/x509/pkix/)

<!-- CL 229864, CL 240543 -->
[`Name.String`](/pkg/crypto/x509/pkix/#Name.String)
now prints non-standard attributes from
[`Names`](/pkg/crypto/x509/pkix/#Name.Names) if
[`ExtraNames`](/pkg/crypto/x509/pkix/#Name.ExtraNames) is nil.

<!-- crypto/x509/pkix -->

#### [database/sql](/pkg/database/sql/)

<!-- CL 145758 -->
The new [`DB.SetConnMaxIdleTime`](/pkg/database/sql/#DB.SetConnMaxIdleTime)
method allows removing a connection from the connection pool after
it has been idle for a period of time, without regard to the total
lifespan of the connection. The [`DBStats.MaxIdleTimeClosed`](/pkg/database/sql/#DBStats.MaxIdleTimeClosed)
field shows the total number of connections closed due to
`DB.SetConnMaxIdleTime`.

<!-- CL 214317 -->
The new [`Row.Err`](/pkg/database/sql/#Row.Err) getter
allows checking for query errors without calling
`Row.Scan`.

<!-- database/sql -->

#### [database/sql/driver](/pkg/database/sql/driver/)

<!-- CL 174122 -->
The new [`Validator`](/pkg/database/sql/driver/#Validator)
interface may be implemented by `Conn` to allow drivers
to signal if a connection is valid or if it should be discarded.

<!-- database/sql/driver -->

#### [debug/pe](/pkg/debug/pe/)

<!-- CL 222637 -->
The package now defines the
`IMAGE_FILE`, `IMAGE_SUBSYSTEM`,
and `IMAGE_DLLCHARACTERISTICS` constants used by the
PE file format.

<!-- debug/pe -->

#### [encoding/asn1](/pkg/encoding/asn1/)

<!-- CL 226984 -->
[`Marshal`](/pkg/encoding/asn1/#Marshal) now sorts the components
of SET OF according to X.690 DER.

<!-- CL 227320 -->
[`Unmarshal`](/pkg/encoding/asn1/#Unmarshal) now rejects tags and
Object Identifiers which are not minimally encoded according to X.690 DER.

<!-- encoding/asn1 -->

#### [encoding/json](/pkg/encoding/json/)

<!-- CL 199837 -->
The package now has an internal limit to the maximum depth of
nesting when decoding. This reduces the possibility that a
deeply nested input could use large quantities of stack memory,
or even cause a "goroutine stack exceeds limit" panic.

<!-- encoding/json -->

#### [flag](/pkg/flag/)

<!-- CL 221427 -->
When the `flag` package sees `-h` or `-help`,
and those flags are not defined, it now prints a usage message.
If the [`FlagSet`](/pkg/flag/#FlagSet) was created with
[`ExitOnError`](/pkg/flag/#ExitOnError),
[`FlagSet.Parse`](/pkg/flag/#FlagSet.Parse) would then
exit with a status of 2. In this release, the exit status for `-h`
or `-help` has been changed to 0. In particular, this applies to
the default handling of command line flags.

#### [fmt](/pkg/fmt/)

<!-- CL 215001 -->
The printing verbs `%#g` and `%#G` now preserve
trailing zeros for floating-point values.

<!-- fmt -->

#### [go/format](/pkg/go/format/)

<!-- golang.org/issue/37476, CL 231461, CL 240683 -->
The [`Source`](/pkg/go/format/#Source) and
[`Node`](/pkg/go/format/#Node) functions
now canonicalize number literal prefixes and exponents as part
of formatting Go source code. This matches the behavior of the
[`gofmt`](/pkg/cmd/gofmt/) command as it
was implemented [since Go 1.13](/doc/go1.13#gofmt).

<!-- go/format -->

#### [html/template](/pkg/html/template/)

<!-- CL 226097 -->
The package now uses Unicode escapes (`\uNNNN`) in all
JavaScript and JSON contexts. This fixes escaping errors in
`application/ld+json` and `application/json`
contexts.

<!-- html/template -->

#### [io/ioutil](/pkg/io/ioutil/)

<!-- CL 212597 -->
[`TempDir`](/pkg/io/ioutil/#TempDir) and
[`TempFile`](/pkg/io/ioutil/#TempFile)
now reject patterns that contain path separators.
That is, calls such as `ioutil.TempFile("/tmp",` `"../base*")` will no longer succeed.
This prevents unintended directory traversal.

<!-- io/ioutil -->

#### [math/big](/pkg/math/big/)

<!-- CL 230397 -->
The new [`Int.FillBytes`](/pkg/math/big/#Int.FillBytes)
method allows serializing to fixed-size pre-allocated byte slices.

<!-- math/big -->

#### [math/cmplx](/pkg/math/cmplx/)

<!-- CL 220689 -->
The functions in this package were updated to conform to the C99 standard
(Annex G IEC 60559-compatible complex arithmetic) with respect to handling
of special arguments such as infinity, NaN and signed zero.

<!-- math/cmplx-->

#### [net](/pkg/net/)

<!-- CL 228645 -->
If an I/O operation exceeds a deadline set by
the [`Conn.SetDeadline`](/pkg/net/#Conn),
`Conn.SetReadDeadline`,
or `Conn.SetWriteDeadline` methods, it will now
return an error that is or wraps
[`os.ErrDeadlineExceeded`](/pkg/os/#ErrDeadlineExceeded).
This may be used to reliably detect whether an error is due to
an exceeded deadline.
Earlier releases recommended calling the `Timeout`
method on the error, but I/O operations can return errors for
which `Timeout` returns `true` although a
deadline has not been exceeded.

<!-- CL 228641 -->
The new [`Resolver.LookupIP`](/pkg/net/#Resolver.LookupIP)
method supports IP lookups that are both network-specific and accept a context.

#### [net/http](/pkg/net/http/)

<!-- CL 231418, CL 231419 -->
Parsing is now stricter as a hardening measure against request smuggling attacks:
non-ASCII white space is no longer trimmed like SP and HTAB, and support for the
"`identity`" `Transfer-Encoding` was dropped.

<!-- net/http -->

#### [net/http/httputil](/pkg/net/http/httputil/)

<!-- CL 230937 -->
[`ReverseProxy`](/pkg/net/http/httputil/#ReverseProxy)
now supports not modifying the `X-Forwarded-For`
header when the incoming `Request.Header` map entry
for that field is `nil`.

<!-- CL 224897 -->
When a Switching Protocol (like WebSocket) request handled by
[`ReverseProxy`](/pkg/net/http/httputil/#ReverseProxy)
is canceled, the backend connection is now correctly closed.

#### [net/http/pprof](/pkg/net/http/pprof/)

<!-- CL 147598, CL 229537 -->
All profile endpoints now support a "`seconds`" parameter. When present,
the endpoint profiles for the specified number of seconds and reports the difference.
The meaning of the "`seconds`" parameter in the `cpu` profile and
the trace endpoints is unchanged.

#### [net/url](/pkg/net/url/)

<!-- CL 227645 -->
The new [`URL`](/pkg/net/url/#URL) field
`RawFragment` and method [`EscapedFragment`](/pkg/net/url/#URL.EscapedFragment)
provide detail about and control over the exact encoding of a particular fragment.
These are analogous to
`RawPath` and [`EscapedPath`](/pkg/net/url/#URL.EscapedPath).

<!-- CL 207082 -->
The new [`URL`](/pkg/net/url/#URL)
method [`Redacted`](/pkg/net/url/#URL.Redacted)
returns the URL in string form with any password replaced with `xxxxx`.

#### [os](/pkg/os/)

<!-- CL -->
If an I/O operation exceeds a deadline set by
the [`File.SetDeadline`](/pkg/os/#File.SetDeadline),
[`File.SetReadDeadline`](/pkg/os/#File.SetReadDeadline),
or [`File.SetWriteDeadline`](/pkg/os/#File.SetWriteDeadline)
methods, it will now return an error that is or wraps
[`os.ErrDeadlineExceeded`](/pkg/os/#ErrDeadlineExceeded).
This may be used to reliably detect whether an error is due to
an exceeded deadline.
Earlier releases recommended calling the `Timeout`
method on the error, but I/O operations can return errors for
which `Timeout` returns `true` although a
deadline has not been exceeded.

<!-- CL 232862 -->
Packages `os` and `net` now automatically
retry system calls that fail with `EINTR`. Previously
this led to spurious failures, which became more common in Go
1.14 with the addition of asynchronous preemption. Now this is
handled transparently.

<!-- CL 229101 -->
The [`os.File`](/pkg/os/#File) type now
supports a [`ReadFrom`](/pkg/os/#File.ReadFrom)
method. This permits the use of the `copy_file_range`
system call on some systems when using
[`io.Copy`](/pkg/io/#Copy) to copy data
from one `os.File` to another. A consequence is that
[`io.CopyBuffer`](/pkg/io/#CopyBuffer)
will not always use the provided buffer when copying to a
`os.File`. If a program wants to force the use of
the provided buffer, it can be done by writing
`io.CopyBuffer(struct{ io.Writer }{dst}, src, buf)`.

#### [plugin](/pkg/plugin/)

<!-- CL 182959 -->
DWARF generation is now supported (and enabled by default) for `-buildmode=plugin` on macOS.

<!-- CL 191617 -->
Building with `-buildmode=plugin` is now supported on `freebsd/amd64`.

#### [reflect](/pkg/reflect/)

<!-- CL 228902 -->
Package `reflect` now disallows accessing methods of all
non-exported fields, whereas previously it allowed accessing
those of non-exported, embedded fields. Code that relies on the
previous behavior should be updated to instead access the
corresponding promoted method of the enclosing variable.

#### [regexp](/pkg/regexp/)

<!-- CL 187919 -->
The new [`Regexp.SubexpIndex`](/pkg/regexp/#Regexp.SubexpIndex)
method returns the index of the first subexpression with the given name
within the regular expression.

<!-- regexp -->

#### [runtime](/pkg/runtime/)

<!-- CL 216557 -->
Several functions, including
[`ReadMemStats`](/pkg/runtime/#ReadMemStats)
and
[`GoroutineProfile`](/pkg/runtime/#GoroutineProfile),
no longer block if a garbage collection is in progress.

#### [runtime/pprof](/pkg/runtime/pprof/)

<!-- CL 189318 -->
The goroutine profile now includes the profile labels associated with each
goroutine at the time of profiling. This feature is not yet implemented for
the profile reported with `debug=2`.

#### [strconv](/pkg/strconv/)

<!-- CL 216617 -->
[`FormatComplex`](/pkg/strconv/#FormatComplex) and [`ParseComplex`](/pkg/strconv/#ParseComplex) are added for working with complex numbers.

[`FormatComplex`](/pkg/strconv/#FormatComplex) converts a complex number into a string of the form (a+bi), where a and b are the real and imaginary parts.

[`ParseComplex`](/pkg/strconv/#ParseComplex) converts a string into a complex number of a specified precision. `ParseComplex` accepts complex numbers in the format `N+Ni`.

<!-- strconv -->

#### [sync](/pkg/sync/)

<!-- CL 205899, golang.org/issue/33762 -->
The new method
[`Map.LoadAndDelete`](/pkg/sync/#Map.LoadAndDelete)
atomically deletes a key and returns the previous value if present.

<!-- CL 205899 -->
The method
[`Map.Delete`](/pkg/sync/#Map.Delete)
is more efficient.

<!-- sync -->

#### [syscall](/pkg/syscall/)

<!-- CL 231638 -->
On Unix systems, functions that use
[`SysProcAttr`](/pkg/syscall/#SysProcAttr)
will now reject attempts to set both the `Setctty`
and `Foreground` fields, as they both use
the `Ctty` field but do so in incompatible ways.
We expect that few existing programs set both fields.

Setting the `Setctty` field now requires that the
`Ctty` field be set to a file descriptor number in the
child process, as determined by the `ProcAttr.Files` field.
Using a child descriptor always worked, but there were certain
cases where using a parent file descriptor also happened to work.
Some programs that set `Setctty` will need to change
the value of `Ctty` to use a child descriptor number.

<!-- CL 220578 -->
It is [now possible](/pkg/syscall/#Proc.Call) to call
system calls that return floating point values
on `windows/amd64`.

#### [testing](/pkg/testing/)

<!-- golang.org/issue/28135 -->
The `testing.T` type now has a
[`Deadline`](/pkg/testing/#T.Deadline) method
that reports the time at which the test binary will have exceeded its
timeout.

<!-- golang.org/issue/34129 -->
A `TestMain` function is no longer required to call
`os.Exit`. If a `TestMain` function returns,
the test binary will call `os.Exit` with the value returned
by `m.Run`.

<!-- CL 226877, golang.org/issue/35998 -->
The new methods
[`T.TempDir`](/pkg/testing/#T.TempDir) and
[`B.TempDir`](/pkg/testing/#B.TempDir)
return temporary directories that are automatically cleaned up
at the end of the test.

<!-- CL 229085 -->
`go` `test` `-v` now groups output by
test name, rather than printing the test name on each line.

<!-- testing -->

#### [text/template](/pkg/text/template/)

<!-- CL 226097 -->
[`JSEscape`](/pkg/text/template/#JSEscape) now
consistently uses Unicode escapes (`\u00XX`), which are
compatible with JSON.

<!-- text/template -->

#### [time](/pkg/time/)

<!-- CL 220424, CL 217362, golang.org/issue/33184 -->
The new method
[`Ticker.Reset`](/pkg/time/#Ticker.Reset)
supports changing the duration of a ticker.

<!-- CL 227878 -->
When returning an error, [`ParseDuration`](/pkg/time/#ParseDuration) now quotes the original value.

<!-- time -->
