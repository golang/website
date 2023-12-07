---
template: false
title: Go 1.8 Release Notes
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

## Introduction to Go 1.8 {#introduction}

The latest Go release, version 1.8, arrives six months after [Go 1.7](go1.7).
Most of its changes are in the implementation of the toolchain, runtime, and libraries.
There are [two minor changes](#language) to the language specification.
As always, the release maintains the Go 1 [promise of compatibility](/doc/go1compat.html).
We expect almost all Go programs to continue to compile and run as before.

The release [adds support for 32-bit MIPS](#ports),
[updates the compiler back end](#compiler) to generate more efficient code,
[reduces GC pauses](#gc) by eliminating stop-the-world stack rescanning,
[adds HTTP/2 Push support](#h2push),
[adds HTTP graceful shutdown](#http_shutdown),
[adds more context support](#more_context),
[enables profiling mutexes](#mutex_prof),
and [simplifies sorting slices](#sort_slice).

## Changes to the language {#language}

When explicitly converting a value from one struct type to another,
as of Go 1.8 the tags are ignored. Thus two structs that differ
only in their tags may be converted from one to the other:

	func example() {
		type T1 struct {
			X int `json:"foo"`
		}
		type T2 struct {
			X int `json:"bar"`
		}
		var v1 T1
		var v2 T2
		v1 = T1(v2) // now legal
	}

<!-- CL 17711 -->
The language specification now only requires that implementations
support up to 16-bit exponents in floating-point constants. This does not affect
either the “[`gc`](/cmd/compile/)” or
`gccgo` compilers, both of
which still support 32-bit exponents.

## Ports {#ports}

Go now supports 32-bit MIPS on Linux for both big-endian
(`linux/mips`) and little-endian machines
(`linux/mipsle`) that implement the MIPS32r1 instruction set with FPU
or kernel FPU emulation. Note that many common MIPS-based routers lack an FPU and
have firmware that doesn't enable kernel FPU emulation; Go won't run on such machines.

On DragonFly BSD, Go now requires DragonFly 4.4.4 or later. <!-- CL 29491, CL 29971 -->

On OpenBSD, Go now requires OpenBSD 5.9 or later. <!-- CL 34093 -->

The Plan 9 port's networking support is now much more complete
and matches the behavior of Unix and Windows with respect to deadlines
and cancellation. For Plan 9 kernel requirements, see the
[Plan 9 wiki page](/wiki/Plan9).

Go 1.8 now only supports OS X 10.8 or later. This is likely the last
Go release to support 10.8. Compiling Go or running
binaries on older OS X versions is untested.

Go 1.8 will be the last release to support Linux on ARMv5E and ARMv6 processors:
Go 1.9 will likely require the ARMv6K (as found in the Raspberry Pi 1) or later.
To identify whether a Linux system is ARMv6K or later, run
“`go` `tool` `dist` `-check-armv6k`”
(to facilitate testing, it is also possible to just copy the `dist` command to the
system without installing a full copy of Go 1.8)
and if the program terminates with output "ARMv6K supported." then the system
implements ARMv6K or later.
Go on non-Linux ARM systems already requires ARMv6K or later.

<!-- CL 31596, go.dev/issue/17528 -->
`zos` is now a recognized value for `GOOS`,
reserved for the z/OS operating system.

### Known Issues {#known_issues}

There are some instabilities on FreeBSD and NetBSD that are known but not understood.
These can lead to program crashes in rare cases.
See
[issue 15658](/issue/15658) and
[issue 16511](/issue/16511).
Any help in solving these issues would be appreciated.

## Tools {#tools}

### Assembler {#cmd_asm}

For 64-bit x86 systems, the following instructions have been added:
`VBROADCASTSD`,
`BROADCASTSS`,
`MOVDDUP`,
`MOVSHDUP`,
`MOVSLDUP`,
`VMOVDDUP`,
`VMOVSHDUP`, and
`VMOVSLDUP`.

For 64-bit PPC systems, the common vector scalar instructions have been
added:
`LXS`,
`LXSDX`,
`LXSI`,
`LXSIWAX`,
`LXSIWZX`,
`LXV`,
`LXVD2X`,
`LXVDSX`,
`LXVW4X`,
`MFVSR`,
`MFVSRD`,
`MFVSRWZ`,
`MTVSR`,
`MTVSRD`,
`MTVSRWA`,
`MTVSRWZ`,
`STXS`,
`STXSDX`,
`STXSI`,
`STXSIWX`,
`STXV`,
`STXVD2X`,
`STXVW4X`,
`XSCV`,
`XSCVDPSP`,
`XSCVDPSPN`,
`XSCVDPSXDS`,
`XSCVDPSXWS`,
`XSCVDPUXDS`,
`XSCVDPUXWS`,
`XSCVSPDP`,
`XSCVSPDPN`,
`XSCVSXDDP`,
`XSCVSXDSP`,
`XSCVUXDDP`,
`XSCVUXDSP`,
`XSCVX`,
`XSCVXP`,
`XVCV`,
`XVCVDPSP`,
`XVCVDPSXDS`,
`XVCVDPSXWS`,
`XVCVDPUXDS`,
`XVCVDPUXWS`,
`XVCVSPDP`,
`XVCVSPSXDS`,
`XVCVSPSXWS`,
`XVCVSPUXDS`,
`XVCVSPUXWS`,
`XVCVSXDDP`,
`XVCVSXDSP`,
`XVCVSXWDP`,
`XVCVSXWSP`,
`XVCVUXDDP`,
`XVCVUXDSP`,
`XVCVUXWDP`,
`XVCVUXWSP`,
`XVCVX`,
`XVCVXP`,
`XXLAND`,
`XXLANDC`,
`XXLANDQ`,
`XXLEQV`,
`XXLNAND`,
`XXLNOR`,
`XXLOR`,
`XXLORC`,
`XXLORQ`,
`XXLXOR`,
`XXMRG`,
`XXMRGHW`,
`XXMRGLW`,
`XXPERM`,
`XXPERMDI`,
`XXSEL`,
`XXSI`,
`XXSLDWI`,
`XXSPLT`, and
`XXSPLTW`.

### Yacc {#tool_yacc}

<!-- CL 27324, CL 27325 -->
The `yacc` tool (previously available by running
“`go` `tool` `yacc`”) has been removed.
As of Go 1.7 it was no longer used by the Go compiler.
It has moved to the “tools” repository and is now available at
[`golang.org/x/tools/cmd/goyacc`](https://godoc.org/golang.org/x/tools/cmd/goyacc).

### Fix {#tool_fix}

<!-- CL 28872 -->
The `fix` tool has a new “`context`”
fix to change imports from “`golang.org/x/net/context`”
to “[`context`](/pkg/context/)”.

### Pprof {#tool_pprof}

<!-- CL 33157 -->
The `pprof` tool can now profile TLS servers
and skip certificate validation by using the “`https+insecure`”
URL scheme.

<!-- CL 23781 -->
The callgrind output now has instruction-level granularity.

### Trace {#tool_trace}

<!-- CL 23324 -->
The `trace` tool has a new `-pprof` flag for
producing pprof-compatible blocking and latency profiles from an
execution trace.

<!-- CL 30017, CL 30702 -->
Garbage collection events are now shown more clearly in the
execution trace viewer. Garbage collection activity is shown on its
own row and GC helper goroutines are annotated with their roles.

### Vet {#tool_vet}

Vet is stricter in some ways and looser where it
previously caused false positives.

Vet now checks for copying an array of locks,
duplicate JSON and XML struct field tags,
non-space-separated struct tags,
deferred calls to HTTP `Response.Body.Close`
before checking errors, and
indexed arguments in `Printf`.
It also improves existing checks.


### Compiler Toolchain {#compiler}

Go 1.7 introduced a new compiler back end for 64-bit x86 systems.
In Go 1.8, that back end has been developed further and is now used for
all architectures.

The new back end, based on
[static single assignment form](https://en.wikipedia.org/wiki/Static_single_assignment_form) (SSA),
generates more compact, more efficient code
and provides a better platform for optimizations
such as bounds check elimination.
The new back end reduces the CPU time required by
our benchmark programs by 20-30%
on 32-bit ARM systems. For 64-bit x86 systems, which already used the SSA back end in
Go 1.7, the gains are a more modest 0-10%. Other architectures will likely
see improvements closer to the 32-bit ARM numbers.

The temporary `-ssa=0` compiler flag introduced in Go 1.7
to disable the new back end has been removed in Go 1.8.

In addition to enabling the new compiler back end for all systems,
Go 1.8 also introduces a new compiler front end. The new compiler
front end should not be noticeable to users but is the foundation for
future performance work.

The compiler and linker have been optimized and run faster in this
release than in Go 1.7, although they are still slower than we would
like and will continue to be optimized in future releases.
Compared to the previous release, Go 1.8 is
[about 15% faster](https://dave.cheney.net/2016/11/19/go-1-8-toolchain-improvements).

### Cgo {#cmd_cgo}

<!-- CL 31141 -->
The Go tool now remembers the value of the `CGO_ENABLED` environment
variable set during `make.bash` and applies it to all future compilations
by default to fix issue [#12808](/issue/12808).
When doing native compilation, it is rarely necessary to explicitly set
the `CGO_ENABLED` environment variable as `make.bash`
will detect the correct setting automatically. The main reason to explicitly
set the `CGO_ENABLED` environment variable is when your environment
supports cgo, but you explicitly do not want cgo support, in which case, set
`CGO_ENABLED=0` during `make.bash` or `all.bash`.

<!-- CL 29991 -->
The environment variable `PKG_CONFIG` may now be used to
set the program to run to handle `#cgo` `pkg-config`
directives. The default is `pkg-config`, the program
always used by earlier releases. This is intended to make it easier
to cross-compile
[cgo](/cmd/cgo/) code.

<!-- CL 32354 -->
The [cgo](/cmd/cgo/) tool now supports a `-srcdir`
option, which is used by the [go](/cmd/go/) command.

<!-- CL 31768, 31811 -->
If [cgo](/cmd/cgo/) code calls `C.malloc`, and
`malloc` returns `NULL`, the program will now
crash with an out of memory error.
`C.malloc` will never return `nil`.
Unlike most C functions, `C.malloc` may not be used in a
two-result form returning an errno value.

<!-- CL 33237 -->
If [cgo](/cmd/cgo/) is used to call a C function passing a
pointer to a C union, and if the C union can contain any pointer
values, and if [cgo pointer
checking](/cmd/cgo/#hdr-Passing_pointers) is enabled (as it is by default), the union value is now
checked for Go pointers.

### Gccgo {#gccgo}

Due to the alignment of Go's semiannual release schedule with GCC's
annual release schedule,
GCC release 6 contains the Go 1.6.1 version of gccgo.
We expect that the next release, GCC 7, will contain the Go 1.8
version of gccgo.

### Default GOPATH {#gopath}

The
[`GOPATH`
environment variable](/cmd/go/#hdr-GOPATH_environment_variable) now has a default value if it
is unset. It defaults to
`$HOME/go` on Unix and
`%USERPROFILE%/go` on Windows.

### Go get {#go_get}

<!-- CL 34818 -->
The “`go` `get`” command now always respects
HTTP proxy environment variables, regardless of whether
the <code style="white-space:nowrap">-insecure</code> flag is used. In previous releases, the
<code style="white-space:nowrap">-insecure</code> flag had the side effect of not using proxies.

### Go bug {#go_bug}

The new
“[`go` `bug`](/cmd/go/#hdr-Print_information_for_bug_reports)”
command starts a bug report on GitHub, prefilled
with information about the current system.

### Go doc {#cmd_doc}

<!-- CL 25419 -->
The
“[`go` `doc`](/cmd/go/#hdr-Show_documentation_for_package_or_symbol)”
command now groups constants and variables with their type,
following the behavior of
[`godoc`](/cmd/godoc/).

<!-- CL 25420 -->
In order to improve the readability of `doc`'s
output, each summary of the first-level items is guaranteed to
occupy a single line.

<!-- CL 31852 -->
Documentation for a specific method in an interface definition can
now be requested, as in
“`go` `doc` `net.Conn.SetDeadline`”.

### Plugins {#plugin}

Go now provides early support for plugins with a “`plugin`”
build mode for generating plugins written in Go, and a
new [`plugin`](/pkg/plugin/) package for
loading such plugins at run time. Plugin support is currently only
available on Linux. Please report any issues.

## Runtime {#runtime}

### Argument Liveness {#liveness}

<!-- Issue 15843 -->
The garbage collector no longer considers
arguments live throughout the entirety of a function. For more
information, and for how to force a variable to remain live, see
the [`runtime.KeepAlive`](/pkg/runtime/#KeepAlive)
function added in Go 1.7.

_Updating:_
Code that sets a finalizer on an allocated object may need to add
calls to `runtime.KeepAlive` in functions or methods
using that object.
Read the
[`KeepAlive`
documentation](/pkg/runtime/#KeepAlive) and its example for more details.

### Concurrent Map Misuse {#mapiter}

In Go 1.6, the runtime
[added lightweight,
best-effort detection of concurrent misuse of maps](/doc/go1.6#runtime). This release
improves that detector with support for detecting programs that
concurrently write to and iterate over a map.

As always, if one goroutine is writing to a map, no other goroutine should be
reading (which includes iterating) or writing the map concurrently.
If the runtime detects this condition, it prints a diagnosis and crashes the program.
The best way to find out more about the problem is to run the program
under the
[race detector](/blog/race-detector),
which will more reliably identify the race
and give more detail.

### MemStats Documentation {#memstats}

<!-- CL 28972 -->
The [`runtime.MemStats`](/pkg/runtime/#MemStats)
type has been more thoroughly documented.

## Performance {#performance}

As always, the changes are so general and varied that precise statements
about performance are difficult to make.
Most programs should run a bit faster,
due to speedups in the garbage collector and
optimizations in the standard library.

There have been optimizations to implementations in the
[`bytes`](/pkg/bytes/),
[`crypto/aes`](/pkg/crypto/aes/),
[`crypto/cipher`](/pkg/crypto/cipher/),
[`crypto/elliptic`](/pkg/crypto/elliptic/),
[`crypto/sha256`](/pkg/crypto/sha256/),
[`crypto/sha512`](/pkg/crypto/sha512/),
[`encoding/asn1`](/pkg/encoding/asn1/),
[`encoding/csv`](/pkg/encoding/csv/),
[`encoding/hex`](/pkg/encoding/hex/),
[`encoding/json`](/pkg/encoding/json/),
[`hash/crc32`](/pkg/hash/crc32/),
[`image/color`](/pkg/image/color/),
[`image/draw`](/pkg/image/draw/),
[`math`](/pkg/math/),
[`math/big`](/pkg/math/big/),
[`reflect`](/pkg/reflect/),
[`regexp`](/pkg/regexp/),
[`runtime`](/pkg/runtime/),
[`strconv`](/pkg/strconv/),
[`strings`](/pkg/strings/),
[`syscall`](/pkg/syscall/),
[`text/template`](/pkg/text/template/), and
[`unicode/utf8`](/pkg/unicode/utf8/)
packages.

### Garbage Collector {#gc}

Garbage collection pauses should be significantly shorter than they
were in Go 1.7, usually under 100 microseconds and often as low as
10 microseconds.
See the
[document on eliminating stop-the-world stack re-scanning](https://github.com/golang/proposal/blob/master/design/17503-eliminate-rescan.md)
for details. More work remains for Go 1.9.

### Defer {#defer}

<!-- CL 29656, CL 29656 -->

The overhead of [deferred
function calls](/ref/spec/#Defer_statements) has been reduced by about half.

### Cgo {#cgoperf}

The overhead of calls from Go into C has been reduced by about half.

## Standard library {#library}

### Examples {#examples}

Examples have been added to the documentation across many packages.

### Sort {#sort_slice}

The [sort](/pkg/sort/) package
now includes a convenience function
[`Slice`](/pkg/sort/#Slice) to sort a
slice given a _less_ function.
In many cases this means that writing a new sorter type is not
necessary.

Also new are
[`SliceStable`](/pkg/sort/#SliceStable) and
[`SliceIsSorted`](/pkg/sort/#SliceIsSorted).

### HTTP/2 Push {#h2push}

The [net/http](/pkg/net/http/) package now includes a
mechanism to
send HTTP/2 server pushes from a
[`Handler`](/pkg/net/http/#Handler).
Similar to the existing `Flusher` and `Hijacker`
interfaces, an HTTP/2
[`ResponseWriter`](/pkg/net/http/#ResponseWriter)
now implements the new
[`Pusher`](/pkg/net/http/#Pusher) interface.

### HTTP Server Graceful Shutdown {#http_shutdown}

<!-- CL 32329 -->
The HTTP Server now has support for graceful shutdown using the new
[`Server.Shutdown`](/pkg/net/http/#Server.Shutdown)
method and abrupt shutdown using the new
[`Server.Close`](/pkg/net/http/#Server.Close)
method.

### More Context Support {#more_context}

Continuing [Go 1.7's adoption](/doc/go1.7#context)
of [`context.Context`](/pkg/context/#Context)
into the standard library, Go 1.8 adds more context support
to existing packages:

  - The new [`Server.Shutdown`](/pkg/net/http/#Server.Shutdown)
    takes a context argument.
  - There have been [significant additions](#database_sql) to the
    [database/sql](/pkg/database/sql/) package with context support.
  - All nine of the new `Lookup` methods on the new
    [`net.Resolver`](/pkg/net/#Resolver) now
    take a context.

### Mutex Contention Profiling {#mutex_prof}

The runtime and tools now support profiling contended mutexes.

Most users will want to use the new `-mutexprofile`
flag with “[`go` `test`](/cmd/go/#hdr-Description_of_testing_flags)”,
and then use [pprof](/cmd/pprof/) on the resultant file.

Lower-level support is also available via the new
[`MutexProfile`](/pkg/runtime/#MutexProfile)
and
[`SetMutexProfileFraction`](/pkg/runtime/#SetMutexProfileFraction).

A known limitation for Go 1.8 is that the profile only reports contention for
[`sync.Mutex`](/pkg/sync/#Mutex),
not
[`sync.RWMutex`](/pkg/sync/#RWMutex).

### Minor changes to the library {#minor_library_changes}

As always, there are various minor changes and updates to the library,
made with the Go 1 [promise of compatibility](/doc/go1compat)
in mind. The following sections list the user visible changes and additions.
Optimizations and minor bug fixes are not listed.

#### [archive/tar](/pkg/archive/tar/)

<!-- CL 28471, CL 31440, CL 31441, CL 31444, CL 28418, CL 31439 -->
The tar implementation corrects many bugs in corner cases of the file format.
The [`Reader`](/pkg/archive/tar/#Reader)
is now able to process tar files in the PAX format with entries larger than 8GB.
The [`Writer`](/pkg/archive/tar/#Writer)
no longer produces invalid tar files in some situations involving long pathnames.

#### [compress/flate](/pkg/compress/flate/)

<!-- CL 31640, CL 31174, CL 32149 -->
There have been some minor fixes to the encoder to improve the
compression ratio in certain situations. As a result, the exact
encoded output of `DEFLATE` may be different from Go 1.7. Since
`DEFLATE` is the underlying compression of gzip, png, zlib, and zip,
those formats may have changed outputs.

<!-- CL 31174 -->
The encoder, when operating in
[`NoCompression`](/pkg/compress/flate/#NoCompression)
mode, now produces a consistent output that is not dependent on
the size of the slices passed to the
[`Write`](/pkg/compress/flate/#Writer.Write)
method.

<!-- CL 28216 -->
The decoder, upon encountering an error, now returns any
buffered data it had uncompressed along with the error.

#### [compress/gzip](/pkg/compress/gzip/)

The [`Writer`](/pkg/compress/gzip/#Writer)
now encodes a zero `MTIME` field when
the [`Header.ModTime`](/pkg/compress/gzip/#Header)
field is the zero value.
In previous releases of Go, the `Writer` would encode
a nonsensical value.
Similarly,
the [`Reader`](/pkg/compress/gzip/#Reader)
now reports a zero encoded `MTIME` field as a zero
`Header.ModTime`.

#### [context](/pkg/context/)

<!-- CL 30370 -->
The [`DeadlineExceeded`](/pkg/context#DeadlineExceeded)
error now implements
[`net.Error`](/pkg/net/#Error)
and reports true for both the `Timeout` and
`Temporary` methods.

#### [crypto/tls](/pkg/crypto/tls/)

<!-- CL 25159, CL 31318 -->
The new method
[`Conn.CloseWrite`](/pkg/crypto/tls/#Conn.CloseWrite)
allows TLS connections to be half closed.

<!-- CL 28075 -->
The new method
[`Config.Clone`](/pkg/crypto/tls/#Config.Clone)
clones a TLS configuration.

<!-- CL 30790 -->
The new [`Config.GetConfigForClient`](/pkg/crypto/tls/#Config.GetConfigForClient)
callback allows selecting a configuration for a client dynamically, based
on the client's
[`ClientHelloInfo`](/pkg/crypto/tls/#ClientHelloInfo). <!-- CL 31391, CL 32119 -->
The [`ClientHelloInfo`](/pkg/crypto/tls/#ClientHelloInfo)
struct now has new
fields `Conn`, `SignatureSchemes` (using
the new
type [`SignatureScheme`](/pkg/crypto/tls/#SignatureScheme)),
`SupportedProtos`, and `SupportedVersions`.

<!-- CL 32115 -->
The new [`Config.GetClientCertificate`](/pkg/crypto/tls/#Config.GetClientCertificate)
callback allows selecting a client certificate based on the server's
TLS `CertificateRequest` message, represented by the new
[`CertificateRequestInfo`](/pkg/crypto/tls/#CertificateRequestInfo).

<!-- CL 27434 -->
The new
[`Config.KeyLogWriter`](/pkg/crypto/tls/#Config.KeyLogWriter)
allows debugging TLS connections
in [WireShark](https://www.wireshark.org/) and
similar tools.

<!-- CL 32115 -->
The new
[`Config.VerifyPeerCertificate`](/pkg/crypto/tls/#Config.VerifyPeerCertificate)
callback allows additional validation of a peer's presented certificate.

<!-- CL 18130 -->
The `crypto/tls` package now implements basic
countermeasures against CBC padding oracles. There should be
no explicit secret-dependent timings, but it does not attempt to
normalize memory accesses to prevent cache timing leaks.

The `crypto/tls` package now supports
X25519 and <!-- CL 30824, CL 30825 -->
ChaCha20-Poly1305. <!-- CL 30957, CL 30958 -->
ChaCha20-Poly1305 is now prioritized unless <!-- CL 32871 -->
hardware support for AES-GCM is present.

<!-- CL 27315, CL 35290 -->
AES-128-CBC cipher suites with SHA-256 are also
now supported, but disabled by default.

#### [crypto/x509](/pkg/crypto/x509/)

<!-- CL 24743 -->
PSS signatures are now supported.

<!-- CL 32644 -->
[`UnknownAuthorityError`](/pkg/crypto/x509/#UnknownAuthorityError)
now has a `Cert` field, reporting the untrusted
certificate.

Certificate validation is more permissive in a few cases and
stricter in a few other cases. <!--
crypto/x509: allow a leaf certificate to be specified directly as root (CL 27393)
crypto/x509: check that the issuer name matches the issuer's subject name (CL 23571)
crypto/x509: don't accept a root that already appears in a chain. (CL 32121)
crypto/x509: fix name constraints handling (CL 30155)
crypto/x509: parse all names in an RDN (CL 30810)
crypto/x509: recognise ISO OID for RSA+SHA1 (CL 27394)
crypto/x509: require a NULL parameters for RSA public keys (CL 16166, CL 27312)
crypto/x509: return error for missing SerialNumber (CL 27238)
-->

<!-- CL 30375 -->
Root certificates will now also be looked for
at `/etc/pki/ca-trust/extracted/pem/tls-ca-bundle.pem`
on Linux, to support RHEL and CentOS.

#### [database/sql](/pkg/database/sql/)

The package now supports `context.Context`. There are new methods
ending in `Context` such as
[`DB.QueryContext`](/pkg/database/sql/#DB.QueryContext) and
[`DB.PrepareContext`](/pkg/database/sql/#DB.PrepareContext)
that take context arguments. Using the new `Context` methods ensures that
connections are closed and returned to the connection pool when the
request is done; enables canceling in-progress queries
should the driver support that; and allows the database
pool to cancel waiting for the next available connection.

The [`IsolationLevel`](/pkg/database/sql#IsolationLevel)
can now be set when starting a transaction by setting the isolation level
on [`TxOptions.Isolation`](/pkg/database/sql#TxOptions.Isolation) and passing
it to [`DB.BeginTx`](/pkg/database/sql#DB.BeginTx).
An error will be returned if an isolation level is selected that the driver
does not support. A read-only attribute may also be set on the transaction
by setting [`TxOptions.ReadOnly`](/pkg/database/sql/#TxOptions.ReadOnly)
to true.

Queries now expose the SQL column type information for drivers that support it.
Rows can return [`ColumnTypes`](/pkg/database/sql#Rows.ColumnTypes)
which can include SQL type information, column type lengths, and the Go type.

A [`Rows`](/pkg/database/sql/#Rows)
can now represent multiple result sets. After
[`Rows.Next`](/pkg/database/sql/#Rows.Next) returns false,
[`Rows.NextResultSet`](/pkg/database/sql/#Rows.NextResultSet)
may be called to advance to the next result set. The existing `Rows`
should continue to be used after it advances to the next result set.

[`NamedArg`](/pkg/database/sql/#NamedArg) may be used
as query arguments. The new function [`Named`](/pkg/database/sql/#Named)
helps create a [`NamedArg`](/pkg/database/sql/#NamedArg)
more succinctly.

If a driver supports the new
[`Pinger`](/pkg/database/sql/driver/#Pinger)
interface, the
[`DB.Ping`](/pkg/database/sql/#DB.Ping)
and
[`DB.PingContext`](/pkg/database/sql/#DB.PingContext)
methods will use that interface to check whether a
database connection is still valid.

The new `Context` query methods work for all drivers, but
`Context` cancellation is not responsive unless the driver has been
updated to use them. The other features require driver support in
[`database/sql/driver`](/pkg/database/sql/driver).
Driver authors should review the new interfaces. Users of existing
driver should review the driver documentation to see what
it supports and any system specific documentation on each feature.

#### [debug/pe](/pkg/debug/pe/)

<!-- CL 22720, CL 27212, CL 22181, CL 22332, CL 22336, Issue 15345 -->
The package has been extended and is now used by
[the Go linker](/cmd/link/) to read `gcc`-generated object files.
The new
[`File.StringTable`](/pkg/debug/pe/#File.StringTable)
and
[`Section.Relocs`](/pkg/debug/pe/#Section.Relocs)
fields provide access to the COFF string table and COFF relocations.
The new
[`File.COFFSymbols`](/pkg/debug/pe/#File.COFFSymbols)
allows low-level access to the COFF symbol table.

#### [encoding/base64](/pkg/encoding/base64/)

<!-- CL 24964 -->
The new
[`Encoding.Strict`](/pkg/encoding/base64/#Encoding.Strict)
method returns an `Encoding` that causes the decoder
to return an error when the trailing padding bits are not zero.

#### [encoding/binary](/pkg/encoding/binary/)

<!-- CL 28514 -->
[`Read`](/pkg/encoding/binary/#Read)
and
[`Write`](/pkg/encoding/binary/#Write)
now support booleans.

#### [encoding/json](/pkg/encoding/json/)

<!-- CL 18692  -->
[`UnmarshalTypeError`](/pkg/encoding/json/#UnmarshalTypeError)
now includes the struct and field name.

<!-- CL 31932 -->
A nil [`Marshaler`](/pkg/encoding/json/#Marshaler)
now marshals as a JSON `null` value.

<!-- CL 21811 -->
A [`RawMessage`](/pkg/encoding/json/#RawMessage) value now
marshals the same as its pointer type.

<!-- CL 30371 -->
[`Marshal`](/pkg/encoding/json/#Marshal)
encodes floating-point numbers using the same format as in ES6,
preferring decimal (not exponential) notation for a wider range of values.
In particular, all floating-point integers up to 2<sup>64</sup> format the
same as the equivalent `int64` representation.

<!-- CL 30944 -->
In previous versions of Go, unmarshaling a JSON `null` into an
[`Unmarshaler`](/pkg/encoding/json/#Unmarshaler)
was considered a no-op; now the `Unmarshaler`'s
`UnmarshalJSON` method is called with the JSON literal
`null` and can define the semantics of that case.

#### [encoding/pem](/pkg/encoding/pem/)

<!-- CL 27391 -->
[`Decode`](/pkg/encoding/pem/#Decode)
is now strict about the format of the ending line.

#### [encoding/xml](/pkg/encoding/xml/)

<!-- CL 30946 -->
[`Unmarshal`](/pkg/encoding/xml/#Unmarshal)
now has wildcard support for collecting all attributes using
the new `",any,attr"` struct tag.

#### [expvar](/pkg/expvar/)

<!-- CL 30917 -->
The new methods
[`Int.Value`](/pkg/expvar/#Int.Value),
[`String.Value`](/pkg/expvar/#String.Value),
[`Float.Value`](/pkg/expvar/#Float.Value), and
[`Func.Value`](/pkg/expvar/#Func.Value)
report the current value of an exported variable.

<!-- CL 24722 -->
The new
function [`Handler`](/pkg/expvar/#Handler)
returns the package's HTTP handler, to enable installing it in
non-standard locations.

#### [fmt](/pkg/fmt/)

<!-- CL 30611 -->
[`Scanf`](/pkg/fmt/#Scanf),
[`Fscanf`](/pkg/fmt/#Fscanf), and
[`Sscanf`](/pkg/fmt/#Sscanf) now
handle spaces differently and more consistently than
previous releases. See the
[scanning documentation](/pkg/fmt/#hdr-Scanning)
for details.

#### [go/doc](/pkg/go/doc/)

<!-- CL 29870 -->
The new [`IsPredeclared`](/pkg/go/doc/#IsPredeclared)
function reports whether a string is a predeclared identifier.

#### [go/types](/pkg/go/types/)

<!-- CL 30715 -->
The new function
[`Default`](/pkg/go/types/#Default)
returns the default "typed" type for an "untyped" type.

<!-- CL 31939 -->
The alignment of `complex64` now matches
the [Go compiler](/cmd/compile/).

#### [html/template](/pkg/html/template/)

<!-- CL 14336 -->
The package now validates
the `"type"` attribute on
a `<script>` tag.

#### [image/png](/pkg/image/png/)

<!-- CL 32143, CL 32140 -->
[`Decode`](/pkg/image/png/#Decode)
(and `DecodeConfig`)
now supports True Color and grayscale transparency.

<!-- CL 29872 -->
[`Encoder`](/pkg/image/png/#Encoder)
is now faster and creates smaller output
when encoding paletted images.

#### [math/big](/pkg/math/big/)

<!-- CL 30706 -->
The new method
[`Int.Sqrt`](/pkg/math/big/#Int.Sqrt)
calculates ⌊√x⌋.

The new method
[`Float.Scan`](/pkg/math/big/#Float.Scan)
is a support routine for
[`fmt.Scanner`](/pkg/fmt/#Scanner).

[`Int.ModInverse`](/pkg/math/big/#Int.ModInverse)
now supports negative numbers.

#### [math/rand](/pkg/math/rand/)

<!-- CL 27253, CL 33456 -->
The new [`Rand.Uint64`](/pkg/math/rand/#Rand.Uint64)
method returns `uint64` values. The
new [`Source64`](/pkg/math/rand/#Source64)
interface describes sources capable of generating such values
directly; otherwise the `Rand.Uint64` method
constructs a `uint64` from two calls
to [`Source`](/pkg/math/rand/#Source)'s
`Int63` method.

#### [mime](/pkg/mime/)

<!-- CL 32175 -->
[`ParseMediaType`](/pkg/mime/#ParseMediaType)
now preserves unnecessary backslash escapes as literals,
in order to support MSIE.
When MSIE sends a full file path (in “intranet mode”), it does not
escape backslashes: “`C:\dev\go\foo.txt`”, not
“`C:\\dev\\go\\foo.txt`”.
If we see an unnecessary backslash escape, we now assume it is from MSIE
and intended as a literal backslash.
No known MIME generators emit unnecessary backslash escapes
for simple token characters like numbers and letters.

#### [mime/quotedprintable](/pkg/mime/quotedprintable/)

The
[`Reader`](/pkg/mime/quotedprintable/#Reader)'s
parsing has been relaxed in two ways to accept
more input seen in the wild. <!-- CL 32174 -->
First, it accepts an equals sign (`=`) not followed
by two hex digits as a literal equal sign. <!-- CL 27530 -->
Second, it silently ignores a trailing equals sign at the end of
an encoded input.

#### [net](/pkg/net/)

<!-- CL 30164, CL 33473 -->
The [`Conn`](/pkg/net/#Conn) documentation
has been updated to clarify expectations of an interface
implementation. Updates in the `net/http` packages
depend on implementations obeying the documentation.

_Updating:_ implementations of the `Conn` interface should verify
they implement the documented semantics. The
[golang.org/x/net/nettest](https://godoc.org/golang.org/x/net/nettest)
package will exercise a `Conn` and validate it behaves properly.

<!-- CL 32099 -->
The new method
[`UnixListener.SetUnlinkOnClose`](/pkg/net/#UnixListener.SetUnlinkOnClose)
sets whether the underlying socket file should be removed from the file system when
the listener is closed.

<!-- CL 29951 -->
The new [`Buffers`](/pkg/net/#Buffers) type permits
writing to the network more efficiently from multiple discontiguous buffers
in memory. On certain machines, for certain types of connections,
this is optimized into an OS-specific batch write operation (such as `writev`).

<!-- CL 29440 -->
The new [`Resolver`](/pkg/net/#Resolver) looks up names and numbers
and supports [`context.Context`](/pkg/context/#Context).
The [`Dialer`](/pkg/net/#Dialer) now has an optional
[`Resolver` field](/pkg/net/#Dialer.Resolver).

<!-- CL 29892 -->
[`Interfaces`](/pkg/net/#Interfaces) is now supported on Solaris.

<!-- CL 29233, CL 24901 -->
The Go DNS resolver now supports `resolv.conf`'s “`rotate`”
and “`option` `ndots:0`” options. The “`ndots`” option is
now respected in the same way as `libresolve`.

#### [net/http](/pkg/net/http/)

Server changes:

  - The server now supports graceful shutdown support, [mentioned above](#http_shutdown).
  - <!-- CL 32024 -->
    The [`Server`](/pkg/net/http/#Server)
    adds configuration options
    `ReadHeaderTimeout` and `IdleTimeout`
    and documents `WriteTimeout`.
  - <!-- CL 32014 -->
    [`FileServer`](/pkg/net/http/#FileServer)
    and
    [`ServeContent`](/pkg/net/http/#ServeContent)
    now support HTTP `If-Match` conditional requests,
    in addition to the previous `If-None-Match`
    support for ETags properly formatted according to RFC 7232, section 2.3.

There are several additions to what a server's `Handler` can do:

  - <!-- CL 31173 -->
    The [`Context`](/pkg/context/#Context)
    returned
    by [`Request.Context`](/pkg/net/http/#Request.Context)
    is canceled if the underlying `net.Conn`
    closes. For instance, if the user closes their browser in the
    middle of a slow request, the `Handler` can now
    detect that the user is gone. This complements the
    existing [`CloseNotifier`](/pkg/net/http/#CloseNotifier)
    support. This functionality requires that the underlying
    [`net.Conn`](/pkg/net/#Conn) implements
    [recently clarified interface documentation](#net).
  - <!-- CL 32479 -->
    To serve trailers produced after the header has already been written,
    see the new
    [`TrailerPrefix`](/pkg/net/http/#TrailerPrefix)
    mechanism.
  - <!-- CL 33099 -->
    A `Handler` can now abort a response by panicking
    with the error
    [`ErrAbortHandler`](/pkg/net/http/#ErrAbortHandler).
  - <!-- CL 30812 -->
    A `Write` of zero bytes to a
    [`ResponseWriter`](/pkg/net/http/#ResponseWriter)
    is now defined as a
    way to test whether a `ResponseWriter` has been hijacked:
    if so, the `Write` returns
    [`ErrHijacked`](/pkg/net/http/#ErrHijacked)
    without printing an error
    to the server's error log.

Client & Transport changes:

  - <!-- CL 28930, CL 31435 -->
    The [`Client`](/pkg/net/http/#Client)
    now copies most request headers on redirect. See
    [the documentation](/pkg/net/http/#Client)
    on the `Client` type for details.
  - <!-- CL 29072 -->
    The [`Transport`](/pkg/net/http/#Transport)
    now supports international domain names. Consequently, so do
    [Get](/pkg/net/http/#Get) and other helpers.
  - <!-- CL 31733, CL 29852 -->
    The `Client` now supports 301, 307, and 308 redirects.
    For example, `Client.Post` now follows 301
    redirects, converting them to `GET` requests
    without bodies, like it did for 302 and 303 redirect responses
    previously.
    The `Client` now also follows 307 and 308
    redirects, preserving the original request method and body, if
    any. If the redirect requires resending the request body, the
    request must have the new
    [`Request.GetBody`](/pkg/net/http/#Request)
    field defined.
    [`NewRequest`](/pkg/net/http/#NewRequest)
    sets `Request.GetBody` automatically for common
    body types.
  - <!-- CL 32482 -->
    The `Transport` now rejects requests for URLs with
    ports containing non-digit characters.
  - <!-- CL 27117 -->
    The `Transport` will now retry non-idempotent
    requests if no bytes were written before a network failure
    and the request has no body.
  - <!-- CL 32481 -->
    The
    new [`Transport.ProxyConnectHeader`](/pkg/net/http/#Transport)
    allows configuration of header values to send to a proxy
    during a `CONNECT` request.
  - <!-- CL 28077 -->
    The [`DefaultTransport.Dialer`](/pkg/net/http/#DefaultTransport)
    now enables `DualStack` ("[Happy Eyeballs](https://tools.ietf.org/html/rfc6555)") support,
    allowing the use of IPv4 as a backup if it looks like IPv6 might be
    failing.
  - <!-- CL 31726 -->
    The [`Transport`](/pkg/net/http/#Transport)
    no longer reads a byte of a non-nil
    [`Request.Body`](/pkg/net/http/#Request.Body)
    when the
    [`Request.ContentLength`](/pkg/net/http/#Request.ContentLength)
    is zero to determine whether the `ContentLength`
    is actually zero or just undefined.
    To explicitly signal that a body has zero length,
    either set it to `nil`, or set it to the new value
    [`NoBody`](/pkg/net/http/#NoBody).
    The new `NoBody` value is intended for use by `Request`
    constructor functions; it is used by
    [`NewRequest`](/pkg/net/http/#NewRequest).

#### [net/http/httptrace](/pkg/net/http/httptrace/)

<!-- CL 30359 -->
There is now support for tracing a client request's TLS handshakes with
the new
[`ClientTrace.TLSHandshakeStart`](/pkg/net/http/httptrace/#ClientTrace.TLSHandshakeStart)
and
[`ClientTrace.TLSHandshakeDone`](/pkg/net/http/httptrace/#ClientTrace.TLSHandshakeDone).

#### [net/http/httputil](/pkg/net/http/httputil/)

<!-- CL 32356 -->
The [`ReverseProxy`](/pkg/net/http/httputil/#ReverseProxy)
has a new optional hook,
[`ModifyResponse`](/pkg/net/http/httputil/#ReverseProxy.ModifyResponse),
for modifying the response from the back end before proxying it to the client.

#### [net/mail](/pkg/net/mail/)

<!-- CL 32176 -->
Empty quoted strings are once again allowed in the name part of
an address. That is, Go 1.4 and earlier accepted
`""` `<gopher@example.com>`,
but Go 1.5 introduced a bug that rejected this address.
The address is recognized again.

<!-- CL 31581 -->
The
[`Header.Date`](/pkg/net/mail/#Header.Date)
method has always provided a way to parse
the `Date:` header.
A new function
[`ParseDate`](/pkg/net/mail/#ParseDate)
allows parsing dates found in other
header lines, such as the `Resent-Date:` header.

#### [net/smtp](/pkg/net/smtp/)

<!-- CL 33143 -->
If an implementation of the
[`Auth.Start`](/pkg/net/smtp/#Auth)
method returns an empty `toServer` value,
the package no longer sends
trailing whitespace in the SMTP `AUTH` command,
which some servers rejected.

#### [net/url](/pkg/net/url/)

<!-- CL 31322 -->
The new functions
[`PathEscape`](/pkg/net/url/#PathEscape)
and
[`PathUnescape`](/pkg/net/url/#PathUnescape)
are similar to the query escaping and unescaping functions but
for path elements.

<!-- CL 28933 -->
The new methods
[`URL.Hostname`](/pkg/net/url/#URL.Hostname)
and
[`URL.Port`](/pkg/net/url/#URL.Port)
return the hostname and port fields of a URL,
correctly handling the case where the port may not be present.

<!-- CL 28343 -->
The existing method
[`URL.ResolveReference`](/pkg/net/url/#URL.ResolveReference)
now properly handles paths with escaped bytes without losing
the escaping.

<!-- CL 31467 -->
The `URL` type now implements
[`encoding.BinaryMarshaler`](/pkg/encoding/#BinaryMarshaler) and
[`encoding.BinaryUnmarshaler`](/pkg/encoding/#BinaryUnmarshaler),
making it possible to process URLs in [gob data](/pkg/encoding/gob/).

<!-- CL 29610, CL 31582 -->
Following RFC 3986,
[`Parse`](/pkg/net/url/#Parse)
now rejects URLs like `this_that:other/thing` instead of
interpreting them as relative paths (`this_that` is not a valid scheme).
To force interpretation as a relative path,
such URLs should be prefixed with “`./`”.
The `URL.String` method now inserts this prefix as needed.

#### [os](/pkg/os/)

<!-- CL 16551 -->
The new function
[`Executable`](/pkg/os/#Executable) returns
the path name of the running executable.

<!-- CL 30614 -->
An attempt to call a method on
an [`os.File`](/pkg/os/#File) that has
already been closed will now return the new error
value [`os.ErrClosed`](/pkg/os/#ErrClosed).
Previously it returned a system-specific error such
as `syscall.EBADF`.

<!-- CL 31358 -->
On Unix systems, [`os.Rename`](/pkg/os/#Rename)
will now return an error when used to rename a directory to an
existing empty directory.
Previously it would fail when renaming to a non-empty directory
but succeed when renaming to an empty directory.
This makes the behavior on Unix correspond to that of other systems.

<!-- CL 32451 -->
On Windows, long absolute paths are now transparently converted to
extended-length paths (paths that start with “`\\?\`”).
This permits the package to work with files whose path names are
longer than 260 characters.

<!-- CL 29753 -->
On Windows, [`os.IsExist`](/pkg/os/#IsExist)
will now return `true` for the system
error `ERROR_DIR_NOT_EMPTY`.
This roughly corresponds to the existing handling of the Unix
error `ENOTEMPTY`.

<!-- CL 32152 -->
On Plan 9, files that are not served by `#M` will now
have [`ModeDevice`](/pkg/os/#ModeDevice) set in
the value returned
by [`FileInfo.Mode`](/pkg/os/#FileInfo).

#### [path/filepath](/pkg/path/filepath/)

A number of bugs and corner cases on Windows were fixed:
[`Abs`](/pkg/path/filepath/#Abs) now calls `Clean` as documented,
[`Glob`](/pkg/path/filepath/#Glob) now matches
“`\\?\c:\*`”,
[`EvalSymlinks`](/pkg/path/filepath/#EvalSymlinks) now
correctly handles “`C:.`”, and
[`Clean`](/pkg/path/filepath/#Clean) now properly
handles a leading “`..`” in the path.

#### [reflect](/pkg/reflect/)

<!-- CL 30088 -->
The new function
[`Swapper`](/pkg/reflect/#Swapper) was
added to support [`sort.Slice`](#sortslice).

#### [strconv](/pkg/strconv/)

<!-- CL 31210 -->
The [`Unquote`](/pkg/strconv/#Unquote)
function now strips carriage returns (`\r`) in
backquoted raw strings, following the
[Go language semantics](/ref/spec#String_literals).

#### [syscall](/pkg/syscall/)

<!-- CL 25050, CL 25022 -->
The [`Getpagesize`](/pkg/syscall/#Getpagesize)
now returns the system's size, rather than a constant value.
Previously it always returned 4KB.

<!-- CL 31446 -->
The signature
of [`Utimes`](/pkg/syscall/#Utimes) has
changed on Solaris to match all the other Unix systems'
signature. Portable code should continue to use
[`os.Chtimes`](/pkg/os/#Chtimes) instead.

<!-- CL 32319 -->
The `X__cmsg_data` field has been removed from
[`Cmsghdr`](/pkg/syscall/#Cmsghdr).

#### [text/template](/pkg/text/template/)

<!-- CL 31462 -->
[`Template.Execute`](/pkg/text/template/#Template.Execute)
can now take a
[`reflect.Value`](/pkg/reflect/#Value) as its data
argument, and
[`FuncMap`](/pkg/text/template/#FuncMap)
functions can also accept and return `reflect.Value`.

#### [time](/pkg/time/)

<!-- CL 20118 -->
The new function
[`Until`](/pkg/time/#Until) complements
the analogous `Since` function.

<!-- CL 29338 -->
[`ParseDuration`](/pkg/time/#ParseDuration)
now accepts long fractional parts.

<!-- CL 33429 -->
[`Parse`](/pkg/time/#Parse)
now rejects dates before the start of a month, such as June 0;
it already rejected dates beyond the end of the month, such as
June 31 and July 32.

<!-- CL 33029 -->
<!-- CL 34816 -->
The `tzdata` database has been updated to version
2016j for systems that don't already have a local time zone
database.


#### [testing](/pkg/testing/)

<!-- CL 29970 -->
The new method
[`T.Name`](/pkg/testing/#T.Name)
(and `B.Name`) returns the name of the current
test or benchmark.

<!-- CL 32483 -->
The new function
[`CoverMode`](/pkg/testing/#CoverMode)
reports the test coverage mode.

<!-- CL 32615 -->
Tests and benchmarks are now marked as failed if the race
detector is enabled and a data race occurs during execution.
Previously, individual test cases would appear to pass,
and only the overall execution of the test binary would fail.

<!-- CL 32455 -->
The signature of the
[`MainStart`](/pkg/testing/#MainStart)
function has changed, as allowed by the documentation. It is an
internal detail and not part of the Go 1 compatibility promise.
If you're not calling `MainStart` directly but see
errors, that likely means you set the
normally-empty `GOROOT` environment variable and it
doesn't match the version of your `go` command's binary.

#### [unicode](/pkg/unicode/)

<!-- CL 30935 -->
[`SimpleFold`](/pkg/unicode/#SimpleFold)
now returns its argument unchanged if the provided input was an invalid rune.
Previously, the implementation failed with an index bounds check panic.
