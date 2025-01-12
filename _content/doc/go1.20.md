---
path: /doc/go1.20
template: false
title: Go 1.20 Release Notes
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

## Introduction to Go 1.20 {#introduction}

The latest Go release, version 1.20, arrives six months after [Go 1.19](/doc/go1.19).
Most of its changes are in the implementation of the toolchain, runtime, and libraries.
As always, the release maintains the Go 1 [promise of compatibility](/doc/go1compat).
We expect almost all Go programs to continue to compile and run as before.

## Changes to the language {#language}

Go 1.20 includes four changes to the language.

<!-- https://go.dev/issue/46505 -->
Go 1.17 added [conversions from slice to an array pointer](/ref/spec#Conversions_from_slice_to_array_or_array_pointer).
Go 1.20 extends this to allow conversions from a slice to an array:
given a slice `x`, `[4]byte(x)` can now be written
instead of `*(*[4]byte)(x)`.

<!-- https://go.dev/issue/53003 -->
The [`unsafe` package](/ref/spec/#Package_unsafe) defines
three new functions `SliceData`, `String`, and `StringData`.
Along with Go 1.17's `Slice`, these functions now provide the complete ability to
construct and deconstruct slice and string values, without depending on their exact representation.

<!-- https://go.dev/issue/8606 -->
The specification now defines that struct values are compared one field at a time,
considering fields in the order they appear in the struct type definition,
and stopping at the first mismatch.
The specification could previously have been read as if
all fields needed to be compared beyond the first mismatch.
Similarly, the specification now defines that array values are compared
one element at a time, in increasing index order.
In both cases, the difference affects whether certain comparisons must panic.
Existing programs are unchanged: the new spec wording describes
what the implementations have always done.

<!-- https://go.dev/issue/56548 -->
[Comparable types](/ref/spec#Comparison_operators) (such as ordinary interfaces)
may now satisfy `comparable` constraints, even if the type arguments
are not strictly comparable (comparison may panic at runtime).
This makes it possible to instantiate a type parameter constrained by `comparable`
(e.g., a type parameter for a user-defined generic map key) with a non-strictly comparable type argument
such as an interface type, or a composite type containing an interface type.

## Ports {#ports}

### Windows {#windows}

<!-- https://go.dev/issue/57003, https://go.dev/issue/57004 -->
Go 1.20 is the last release that will run on any release of Windows 7, 8, Server 2008 and Server 2012.
Go 1.21 will require at least Windows 10 or Server 2016.

### Darwin and iOS {#darwin}

<!-- https://go.dev/issue/23011 -->
Go 1.20 is the last release that will run on macOS 10.13 High Sierra or 10.14 Mojave.
Go 1.21 will require macOS 10.15 Catalina or later.

### FreeBSD/RISC-V {#freebsd-riscv}

<!-- https://go.dev/issue/53466 -->
Go 1.20 adds experimental support for FreeBSD on RISC-V (`GOOS=freebsd`, `GOARCH=riscv64`).

## Tools {#tools}

### Go command {#go-command}

<!-- CL 432535, https://go.dev/issue/47257 -->
The directory `$GOROOT/pkg` no longer stores
pre-compiled package archives for the standard library:
`go` `install` no longer writes them,
the `go` build no longer checks for them,
and the Go distribution no longer ships them.
Instead, packages in the standard library are built as needed
and cached in the build cache, just like packages outside `GOROOT`.
This change reduces the size of the Go distribution and also
avoids C toolchain skew for packages that use cgo.

<!-- CL 448357: cmd/go: print test2json start events -->
The implementation of `go` `test` `-json`
has been improved to make it more robust.
Programs that run `go` `test` `-json`
do not need any updates.
Programs that invoke `go` `tool` `test2json`
directly should now run the test binary with `-v=test2json`
(for example, `go` `test` `-v=test2json`
or `./pkg.test` `-test.v=test2json`)
instead of plain `-v`.

<!-- CL 448357: cmd/go: print test2json start events -->
A related change to `go` `test` `-json`
is the addition of an event with `Action` set to `start`
at the beginning of each test program's execution.
When running multiple tests using the `go` command,
these start events are guaranteed to be emitted in the same order as
the packages named on the command line.

<!-- https://go.dev/issue/45454, CL 421434 -->
The `go` command now defines
architecture feature build tags, such as `amd64.v2`,
to allow selecting a package implementation file based on the presence
or absence of a particular architecture feature.
See [`go` `help` `buildconstraint`](/cmd/go#hdr-Build_constraints) for details.

<!-- https://go.dev/issue/50332 -->
The `go` subcommands now accept
`-C` `<dir>` to change directory to \<dir>
before performing the command, which may be useful for scripts that need to
execute commands in multiple different modules.

<!-- https://go.dev/issue/41696, CL 416094 -->
The `go` `build` and `go` `test`
commands no longer accept the `-i` flag,
which has been [deprecated since Go 1.16](/issue/41696).

<!-- https://go.dev/issue/38687, CL 421440 -->
The `go` `generate` command now accepts
`-skip` `<pattern>` to skip `//go:generate` directives
matching `<pattern>`.

<!-- https://go.dev/issue/41583 -->
The `go` `test` command now accepts
`-skip` `<pattern>` to skip tests, subtests, or examples
matching `<pattern>`.

<!-- https://go.dev/issue/37015 -->
When the main module is located within `GOPATH/src`,
`go` `install` no longer installs libraries for
non-`main` packages to `GOPATH/pkg`,
and `go` `list` no longer reports a `Target`
field for such packages. (In module mode, compiled packages are stored in the
[build cache](https://pkg.go.dev/cmd/go#hdr-Build_and_test_caching)
only, but [a bug](/issue/37015) had caused
the `GOPATH` install targets to unexpectedly remain in effect.)

<!-- https://go.dev/issue/55022 -->
The `go` `build`, `go` `install`,
and other build-related commands now support a `-pgo` flag that enables
profile-guided optimization, which is described in more detail in the
[Compiler](#compiler) section below.
The `-pgo` flag specifies the file path of the profile.
Specifying `-pgo=auto` causes the `go` command to search
for a file named `default.pgo` in the main package's directory and
use it if present.
This mode currently requires a single main package to be specified on the
command line, but we plan to lift this restriction in a future release.
Specifying `-pgo=off` turns off profile-guided optimization.

<!-- https://go.dev/issue/51430 -->
The `go` `build`, `go` `install`,
and other build-related commands now support a `-cover`
flag that builds the specified target with code coverage instrumentation.
This is described in more detail in the
[Cover](#cover) section below.

#### `go` `version` {#go-version}

<!-- https://go.dev/issue/48187 -->
The `go` `version` `-m` command
now supports reading more types of Go binaries, most notably, Windows DLLs
built with `go` `build` `-buildmode=c-shared`
and Linux binaries without execute permission.

### Cgo {#cgo}

<!-- CL 450739 -->
The `go` command now disables `cgo` by default
on systems without a C toolchain.
More specifically, when the `CGO_ENABLED` environment variable is unset,
the `CC` environment variable is unset,
and the default C compiler (typically `clang` or `gcc`)
is not found in the path,
`CGO_ENABLED` defaults to `0`.
As always, you can override the default by setting `CGO_ENABLED` explicitly.

The most important effect of the default change is that when Go is installed
on a system without a C compiler, it will now use pure Go builds for packages
in the standard library that use cgo, instead of using pre-distributed package archives
(which have been removed, as [noted above](#go-command))
or attempting to use cgo and failing.
This makes Go work better in some minimal container environments
as well as on macOS, where pre-distributed package archives have
not been used for cgo-based packages since Go 1.16.

The packages in the standard library that use cgo are [`net`](/pkg/net/),
[`os/user`](/pkg/os/user/), and
[`plugin`](/pkg/plugin/).
On macOS, the `net` and `os/user` packages have been rewritten not to use cgo:
the same code is now used for cgo and non-cgo builds as well as cross-compiled builds.
On Windows, the `net` and `os/user` packages have never used cgo.
On other systems, builds with cgo disabled will use a pure Go version of these packages.

A consequence is that, on macOS, if Go code that uses
the `net` package is built
with `-buildmode=c-archive`, linking the resulting
archive into a C program requires passing `-lresolv` when
linking the C code.

On macOS, the race detector has been rewritten not to use cgo:
race-detector-enabled programs can be built and run without Xcode.
On Linux and other Unix systems, and on Windows, a host C toolchain
is required to use the race detector.

### Cover {#cover}

<!-- CL 436236, CL 401236, CL 438503 -->
Go 1.20 supports collecting code coverage profiles for programs
(applications and integration tests), as opposed to just unit tests.

To collect coverage data for a program, build it with `go`
`build`'s `-cover` flag, then run the resulting
binary with the environment variable `GOCOVERDIR` set
to an output directory for coverage profiles.
See the
['coverage for integration tests' landing page](/doc/build-cover) for more on how to get started.
For details on the design and implementation, see the
[proposal](/issue/51430).

### Vet {#vet}

#### Improved detection of loop variable capture by nested functions {#vet-loopclosure}

<!-- CL 447256, https://go.dev/issue/55972: extend the loopclosure analysis to parallel subtests -->
The `vet` tool now reports references to loop variables following
a call to [`T.Parallel()`](/pkg/testing/#T.Parallel)
within subtest function bodies. Such references may observe the value of the
variable from a different iteration (typically causing test cases to be
skipped) or an invalid state due to unsynchronized concurrent access.

<!-- CL 452615 -->
The tool also detects reference mistakes in more places. Previously it would
only consider the last statement of the loop body, but now it recursively
inspects the last statements within if, switch, and select statements.

#### New diagnostic for incorrect time formats {#vet-timeformat}

<!-- CL 354010, https://go.dev/issue/48801: check for time formats with 2006-02-01 -->
The vet tool now reports use of the time format 2006-02-01 (yyyy-dd-mm)
with [`Time.Format`](/pkg/time/#Time.Format) and
[`time.Parse`](/pkg/time/#Parse).
This format does not appear in common date standards, but is frequently
used by mistake when attempting to use the ISO 8601 date format
(yyyy-mm-dd).

## Runtime {#runtime}

<!-- CL 422634 -->
Some of the garbage collector's internal data structures were reorganized to
be both more space and CPU efficient.
This change reduces memory overheads and improves overall CPU performance by
up to 2%.

<!-- CL 417558, https://go.dev/issue/53892 -->
The garbage collector behaves less erratically with respect to goroutine
assists in some circumstances.

<!-- https://go.dev/issue/51430 -->
Go 1.20 adds a new `runtime/coverage` package
containing APIs for writing coverage profile data at
runtime from long-running and/or server programs that
do not terminate via `os.Exit()`.

## Compiler {#compiler}

<!-- https://go.dev/issue/55022 -->
Go 1.20 adds preview support for profile-guided optimization (PGO).
PGO enables the toolchain to perform application- and workload-specific
optimizations based on run-time profile information.
Currently, the compiler supports pprof CPU profiles, which can be collected
through usual means, such as the `runtime/pprof` or
`net/http/pprof` packages.
To enable PGO, pass the path of a pprof profile file via the
`-pgo` flag to `go` `build`,
as mentioned [above](#go-command).
Go 1.20 uses PGO to more aggressively inline functions at hot call sites.
Benchmarks for a representative set of Go programs show enabling
profile-guided inlining optimization improves performance about 3–4%.
See the [PGO user guide](/doc/pgo) for detailed documentation.
We plan to add more profile-guided optimizations in future releases.
Note that profile-guided optimization is a preview, so please use it
with appropriate caution.

The Go 1.20 compiler upgraded its front-end to use a new way of handling the
compiler's internal data, which fixes several generic-types issues and enables
type declarations within generic functions and methods.

<!-- https://go.dev/issue/56103, CL 445598 -->
The compiler now [rejects anonymous interface cycles](/issue/56103)
with a compiler error by default.
These arise from tricky uses of [embedded interfaces](/ref/spec#Embedded_interfaces)
and have always had subtle correctness issues,
yet we have no evidence that they're actually used in practice.
Assuming no reports from users adversely affected by this change,
we plan to update the language specification for Go 1.22 to formally disallow them
so tools authors can stop supporting them too.

<!-- https://go.dev/issue/49569 -->
Go 1.18 and 1.19 saw regressions in build speed, largely due to the addition
of support for generics and follow-on work. Go 1.20 improves build speeds by
up to 10%, bringing it back in line with Go 1.17.
Relative to Go 1.19, generated code performance is also generally slightly improved.

## Linker {#linker}

<!-- https://go.dev/issue/54197, CL 420774 -->
On Linux, the linker now selects the dynamic interpreter for `glibc`
or `musl` at link time.

<!-- https://go.dev/issue/35006 -->
On Windows, the Go linker now supports modern LLVM-based C toolchains.

<!-- https://go.dev/issue/37762, CL 317917 -->
Go 1.20 uses `go:` and `type:` prefixes for compiler-generated
symbols rather than `go.` and `type.`.
This avoids confusion for user packages whose name starts with `go.`.
The [`debug/gosym`](/pkg/debug/gosym) package understands
this new naming convention for binaries built with Go 1.20 and newer.

## Bootstrap {#bootstrap}

<!-- https://go.dev/issue/44505 -->
When building a Go release from source and `GOROOT_BOOTSTRAP` is not set,
previous versions of Go looked for a Go 1.4 or later bootstrap toolchain in the directory
`$HOME/go1.4` (`%HOMEDRIVE%%HOMEPATH%\go1.4` on Windows).
Go 1.18 and Go 1.19 looked first for `$HOME/go1.17` or `$HOME/sdk/go1.17`
before falling back to `$HOME/go1.4`,
in anticipation of requiring Go 1.17 for use when bootstrapping Go 1.20.
Go 1.20 does require a Go 1.17 release for bootstrapping, but we realized that we should
adopt the latest point release of the bootstrap toolchain, so it requires Go 1.17.13.
Go 1.20 looks for `$HOME/go1.17.13` or `$HOME/sdk/go1.17.13`
before falling back to `$HOME/go1.4`
(to support systems that hard-coded the path $HOME/go1.4 but have installed
a newer Go toolchain there).
In the future, we plan to move the bootstrap toolchain forward approximately once a year,
and in particular we expect that Go 1.22 will require the final point release of Go 1.20 for bootstrap.

## Standard library {#library}

### New crypto/ecdh package {#crypto_ecdh}

<!-- https://go.dev/issue/52221, CL 398914, CL 450335, https://go.dev/issue/56052 -->
Go 1.20 adds a new [`crypto/ecdh`](/pkg/crypto/ecdh/) package
to provide explicit support for Elliptic Curve Diffie-Hellman key exchanges
over NIST curves and Curve25519.

Programs should use `crypto/ecdh` instead of the lower-level functionality in
[`crypto/elliptic`](/pkg/crypto/elliptic/) for ECDH, and
third-party modules for more advanced use cases.

### Wrapping multiple errors {#errors}

<!-- CL 432898 -->
Go 1.20 expands support for error wrapping to permit an error to
wrap multiple other errors.

An error `e` can wrap more than one error by providing
an `Unwrap` method that returns a `[]error`.

The [`errors.Is`](/pkg/errors/#Is) and
[`errors.As`](/pkg/errors/#As) functions
have been updated to inspect multiply wrapped errors.

The [`fmt.Errorf`](/pkg/fmt/#Errorf) function
now supports multiple occurrences of the `%w` format verb,
which will cause it to return an error that wraps all of those error operands.

The new function [`errors.Join`](/pkg/errors/#Join)
returns an error wrapping a list of errors.

### HTTP ResponseController {#http_responsecontroller}

<!-- CL 436890, https://go.dev/issue/54136 -->
The new
[`"net/http".ResponseController`](/pkg/net/http/#ResponseController)
type provides access to extended per-request functionality not handled by the
[`"net/http".ResponseWriter`](/pkg/net/http/#ResponseWriter) interface.

Previously, we have added new per-request functionality by defining optional
interfaces which a `ResponseWriter` can implement, such as
[`Flusher`](/pkg/net/http/#Flusher). These interfaces
are not discoverable and clumsy to use.

The `ResponseController` type provides a clearer, more discoverable way
to add per-handler controls. Two such controls also added in Go 1.20 are
`SetReadDeadline` and `SetWriteDeadline`, which allow setting
per-request read and write deadlines. For example:

	func RequestHandler(w ResponseWriter, r *Request) {
	  rc := http.NewResponseController(w)
	  rc.SetWriteDeadline(time.Time{}) // disable Server.WriteTimeout when sending a large response
	  io.Copy(w, bigData)
	}

### New ReverseProxy Rewrite hook {#reverseproxy_rewrite}

<!-- https://go.dev/issue/53002, CL 407214 -->
The [`httputil.ReverseProxy`](/pkg/net/http/httputil/#ReverseProxy)
forwarding proxy includes a new
[`Rewrite`](/pkg/net/http/httputil/#ReverseProxy.Rewrite)
hook function, superseding the
previous `Director` hook.

The `Rewrite` hook accepts a
[`ProxyRequest`](/pkg/net/http/httputil/#ProxyRequest) parameter,
which includes both the inbound request received by the proxy and the outbound
request that it will send.
Unlike `Director` hooks, which only operate on the outbound request,
this permits `Rewrite` hooks to avoid certain scenarios where
a malicious inbound request may cause headers added by the hook
to be removed before forwarding.
See [issue #50580](/issue/50580).

The [`ProxyRequest.SetURL`](/pkg/net/http/httputil/#ProxyRequest.SetURL)
method routes the outbound request to a provided destination
and supersedes the `NewSingleHostReverseProxy` function.
Unlike `NewSingleHostReverseProxy`, `SetURL`
also sets the `Host` header of the outbound request.

<!-- https://go.dev/issue/50465, CL 407414 -->
The
[`ProxyRequest.SetXForwarded`](/pkg/net/http/httputil/#ProxyRequest.SetXForwarded)
method sets the `X-Forwarded-For`, `X-Forwarded-Host`,
and `X-Forwarded-Proto` headers of the outbound request.
When using a `Rewrite`, these headers are not added by default.

An example of a `Rewrite` hook using these features is:

	proxyHandler := &httputil.ReverseProxy{
	  Rewrite: func(r *httputil.ProxyRequest) {
	    r.SetURL(outboundURL) // Forward request to outboundURL.
	    r.SetXForwarded()     // Set X-Forwarded-* headers.
	    r.Out.Header.Set("X-Additional-Header", "header set by the proxy")
	  },
	}

<!-- CL 407375 -->
[`ReverseProxy`](/pkg/net/http/httputil/#ReverseProxy) no longer adds a `User-Agent` header
to forwarded requests when the incoming request does not have one.

### Minor changes to the library {#minor_library_changes}

As always, there are various minor changes and updates to the library,
made with the Go 1 [promise of compatibility](/doc/go1compat)
in mind.
There are also various performance improvements, not enumerated here.

#### [archive/tar](/pkg/archive/tar/)

<!-- https://go.dev/issue/55356, CL 449937 -->
When the `GODEBUG=tarinsecurepath=0` environment variable is set,
[`Reader.Next`](/pkg/archive/tar/#Reader.Next) method
will now return the error [`ErrInsecurePath`](/pkg/archive/tar/#ErrInsecurePath)
for an entry with a file name that is an absolute path,
refers to a location outside the current directory, contains invalid
characters, or (on Windows) is a reserved name such as `NUL`.
A future version of Go may disable insecure paths by default.

<!-- archive/tar -->

#### [archive/zip](/pkg/archive/zip/)

<!-- https://go.dev/issue/55356 -->
When the `GODEBUG=zipinsecurepath=0` environment variable is set,
[`NewReader`](/pkg/archive/zip/#NewReader) will now return the error
[`ErrInsecurePath`](/pkg/archive/zip/#ErrInsecurePath)
when opening an archive which contains any file name that is an absolute path,
refers to a location outside the current directory, contains invalid
characters, or (on Windows) is a reserved names such as `NUL`.
A future version of Go may disable insecure paths by default.

<!-- CL 449955 -->
Reading from a directory file that contains file data will now return an error.
The zip specification does not permit directory files to contain file data,
so this change only affects reading from invalid archives.

<!-- archive/zip -->

#### [bytes](/pkg/bytes/)

<!-- CL 407176 -->
The new
[`CutPrefix`](/pkg/bytes/#CutPrefix) and
[`CutSuffix`](/pkg/bytes/#CutSuffix) functions
are like [`TrimPrefix`](/pkg/bytes/#TrimPrefix)
and [`TrimSuffix`](/pkg/bytes/#TrimSuffix)
but also report whether the string was trimmed.

<!-- CL 359675, https://go.dev/issue/45038 -->
The new [`Clone`](/pkg/bytes/#Clone) function
allocates a copy of a byte slice.

<!-- bytes -->

#### [context](/pkg/context/)

<!-- https://go.dev/issue/51365, CL 375977 -->
The new [`WithCancelCause`](/pkg/context/#WithCancelCause) function
provides a way to cancel a context with a given error.
That error can be retrieved by calling the new [`Cause`](/pkg/context/#Cause) function.

<!-- context -->

#### [crypto/ecdsa](/pkg/crypto/ecdsa/)

<!-- CL 353849 -->
When using supported curves, all operations are now implemented in constant time.
This led to an increase in CPU time between 5% and 30%, mostly affecting P-384 and P-521.

<!-- https://go.dev/issue/56088, CL 450816 -->
The new [`PrivateKey.ECDH`](/pkg/crypto/ecdsa/#PrivateKey.ECDH) method
converts an `ecdsa.PrivateKey` to an `ecdh.PrivateKey`.

<!-- crypto/ecdsa -->

#### [crypto/ed25519](/pkg/crypto/ed25519/)

<!-- CL 373076, CL 404274, https://go.dev/issue/31804 -->
The [`PrivateKey.Sign`](/pkg/crypto/ed25519/#PrivateKey.Sign) method
and the
[`VerifyWithOptions`](/pkg/crypto/ed25519/#VerifyWithOptions) function
now support signing pre-hashed messages with Ed25519ph,
indicated by an
[`Options.HashFunc`](/pkg/crypto/ed25519/#Options.HashFunc)
that returns
[`crypto.SHA512`](/pkg/crypto/#SHA512).
They also now support Ed25519ctx and Ed25519ph with context,
indicated by setting the new
[`Options.Context`](/pkg/crypto/ed25519/#Options.Context)
field.

<!-- crypto/ed25519 -->

#### [crypto/rsa](/pkg/crypto/rsa/)

<!-- CL 418874, https://go.dev/issue/19974 -->
The new field [`OAEPOptions.MGFHash`](/pkg/crypto/rsa/#OAEPOptions.MGFHash)
allows configuring the MGF1 hash separately for OAEP decryption.

<!-- https://go.dev/issue/20654 -->
crypto/rsa now uses a new, safer, constant-time backend. This causes a CPU
runtime increase for decryption operations between approximately 15%
(RSA-2048 on amd64) and 45% (RSA-4096 on arm64), and more on 32-bit architectures.
Encryption operations are approximately 20x slower than before (but still 5-10x faster than decryption).
Performance is expected to improve in future releases.
Programs must not modify or manually generate the fields of
[`PrecomputedValues`](/pkg/crypto/rsa/#PrecomputedValues).

<!-- crypto/rsa -->

#### [crypto/subtle](/pkg/crypto/subtle/)

<!-- https://go.dev/issue/53021, CL 421435 -->
The new function [`XORBytes`](/pkg/crypto/subtle/#XORBytes)
XORs two byte slices together.

<!-- crypto/subtle -->

#### [crypto/tls](/pkg/crypto/tls/)

<!-- CL 426455, CL 427155, CL 426454, https://go.dev/issue/46035 -->
Parsed certificates are now shared across all clients actively using that certificate.
The memory savings can be significant in programs that make many concurrent connections to a
server or collection of servers sharing any part of their certificate chains.

<!-- https://go.dev/issue/48152, CL 449336 -->
For a handshake failure due to a certificate verification failure,
the TLS client and server now return an error of the new type
[`CertificateVerificationError`](/pkg/crypto/tls/#CertificateVerificationError),
which includes the presented certificates.

<!-- crypto/tls -->

#### [crypto/x509](/pkg/crypto/x509/)

<!-- CL 450816, CL 450815 -->
[`ParsePKCS8PrivateKey`](/pkg/crypto/x509/#ParsePKCS8PrivateKey)
and
[`MarshalPKCS8PrivateKey`](/pkg/crypto/x509/#MarshalPKCS8PrivateKey)
now support keys of type [`*crypto/ecdh.PrivateKey`](/pkg/crypto/ecdh.PrivateKey).
[`ParsePKIXPublicKey`](/pkg/crypto/x509/#ParsePKIXPublicKey)
and
[`MarshalPKIXPublicKey`](/pkg/crypto/x509/#MarshalPKIXPublicKey)
now support keys of type [`*crypto/ecdh.PublicKey`](/pkg/crypto/ecdh.PublicKey).
Parsing NIST curve keys still returns values of type
`*ecdsa.PublicKey` and `*ecdsa.PrivateKey`.
Use their new `ECDH` methods to convert to the `crypto/ecdh` types.

<!-- CL 449235 -->
The new [`SetFallbackRoots`](/pkg/crypto/x509/#SetFallbackRoots)
function allows a program to define a set of fallback root certificates in case an
operating system verifier or standard platform root bundle is unavailable at runtime.
It will most commonly be used with a new package, [golang.org/x/crypto/x509roots/fallback](/pkg/golang.org/x/crypto/x509roots/fallback),
which will provide an up to date root bundle.

<!-- crypto/x509 -->

#### [debug/elf](/pkg/debug/elf/)

<!-- CL 429601 -->
Attempts to read from a `SHT_NOBITS` section using
[`Section.Data`](/pkg/debug/elf/#Section.Data)
or the reader returned by [`Section.Open`](/pkg/debug/elf/#Section.Open)
now return an error.

<!-- CL 420982 -->
Additional [`R_LARCH_*`](/pkg/debug/elf/#R_LARCH) constants are defined for use with LoongArch systems.

<!-- CL 420982, CL 435415, CL 425555 -->
Additional [`R_PPC64_*`](/pkg/debug/elf/#R_PPC64) constants are defined for use with PPC64 ELFv2 relocations.

<!-- CL 411915 -->
The constant value for [`R_PPC64_SECTOFF_LO_DS`](/pkg/debug/elf/#R_PPC64_SECTOFF_LO_DS) is corrected, from 61 to 62.

<!-- debug/elf -->

#### [debug/gosym](/pkg/debug/gosym/)

<!-- https://go.dev/issue/37762, CL 317917 -->
Due to a change of [Go's symbol naming conventions](#linker), tools that
process Go binaries should use Go 1.20's `debug/gosym` package to
transparently handle both old and new binaries.

<!-- debug/gosym -->

#### [debug/pe](/pkg/debug/pe/)

<!-- CL 421357 -->
Additional [`IMAGE_FILE_MACHINE_RISCV*`](/pkg/debug/pe/#IMAGE_FILE_MACHINE_RISCV128) constants are defined for use with RISC-V systems.

<!-- debug/pe -->

#### [encoding/binary](/pkg/encoding/binary/)

<!-- CL 420274 -->
The [`ReadVarint`](/pkg/encoding/binary/#ReadVarint) and
[`ReadUvarint`](/pkg/encoding/binary/#ReadUvarint)
functions will now return `io.ErrUnexpectedEOF` after reading a partial value,
rather than `io.EOF`.

<!-- encoding/binary -->

#### [encoding/xml](/pkg/encoding/xml/)

<!-- https://go.dev/issue/53346, CL 424777 -->
The new [`Encoder.Close`](/pkg/encoding/xml/#Encoder.Close) method
can be used to check for unclosed elements when finished encoding.

<!-- CL 103875, CL 105636 -->
The decoder now rejects element and attribute names with more than one colon,
such as `<a:b:c>`,
as well as namespaces that resolve to an empty string, such as `xmlns:a=""`.

<!-- CL 107255 -->
The decoder now rejects elements that use different namespace prefixes in the opening and closing tag,
even if those prefixes both denote the same namespace.

<!-- encoding/xml -->

#### [errors](/pkg/errors/)

<!-- https://go.dev/issue/53435 -->
The new [`Join`](/pkg/errors/#Join) function returns an error wrapping a list of errors.

<!-- errors -->

#### [fmt](/pkg/fmt/)

<!-- https://go.dev/issue/53435 -->
The [`Errorf`](/pkg/fmt/#Errorf) function supports multiple occurrences of
the `%w` format verb, returning an error that unwraps to the list of all arguments to `%w`.

<!-- https://go.dev/issue/51668, CL 400875 -->
The new [`FormatString`](/pkg/fmt/#FormatString) function recovers the
formatting directive corresponding to a [`State`](/pkg/fmt/#State),
which can be useful in [`Formatter`](/pkg/fmt/#Formatter).
implementations.

<!-- fmt -->

#### [go/ast](/pkg/go/ast/)

<!-- CL 426091, https://go.dev/issue/50429 -->
The new [`RangeStmt.Range`](/pkg/go/ast/#RangeStmt.Range) field
records the position of the `range` keyword in a range statement.

<!-- CL 427955, https://go.dev/issue/53202 -->
The new [`File.FileStart`](/pkg/go/ast/#File.FileStart)
and [`File.FileEnd`](/pkg/go/ast/#File.FileEnd) fields
record the position of the start and end of the entire source file.

<!-- go/ast -->

#### [go/token](/pkg/go/token/)

<!-- CL 410114, https://go.dev/issue/53200 -->
The new [`FileSet.RemoveFile`](/pkg/go/token/#FileSet.RemoveFile) method
removes a file from a `FileSet`.
Long-running programs can use this to release memory associated
with files they no longer need.

<!-- go/token -->

#### [go/types](/pkg/go/types/)

<!-- CL 454575 -->
The new [`Satisfies`](/pkg/go/types/#Satisfies) function reports
whether a type satisfies a constraint.
This change aligns with the [new language semantics](#language)
that distinguish satisfying a constraint from implementing an interface.

<!-- go/types -->

#### [html/template](/pkg/html/template/)

<!-- https://go.dev/issue/59153 -->
<!-- CL 481993 -->
Go 1.20.3 and later
[disallow actions in ECMAScript 6 template literals.](/pkg/html/template#hdr-Security_Model)
This behavior can be reverted by the `GODEBUG=jstmpllitinterp=1` setting.

<!-- html/template -->

#### [io](/pkg/io/)

<!-- https://go.dev/issue/45899, CL 406776 -->
The new [`OffsetWriter`](/pkg/io/#OffsetWriter) wraps an underlying
[`WriterAt`](/pkg/io/#WriterAt)
and provides `Seek`, `Write`, and `WriteAt` methods
that adjust their effective file offset position by a fixed amount.

<!-- io -->

#### [io/fs](/pkg/io/fs/)

<!-- CL 363814, https://go.dev/issue/47209 -->
The new error [`SkipAll`](/pkg/io/fs/#SkipAll)
terminates a [`WalkDir`](/pkg/io/fs/#WalkDir)
immediately but successfully.

<!-- io -->

#### [math/big](/pkg/math/big/)

<!-- https://go.dev/issue/52182 -->
The [math/big](/pkg/math/big/) package's wide scope and
input-dependent timing make it ill-suited for implementing cryptography.
The cryptography packages in the standard library no longer call non-trivial
[Int](/pkg/math/big#Int) methods on attacker-controlled inputs.
In the future, the determination of whether a bug in math/big is
considered a security vulnerability will depend on its wider impact on the
standard library.

<!-- math/big -->

#### [math/rand](/pkg/math/rand/)

<!-- https://go.dev/issue/54880, CL 436955, https://go.dev/issue/56319 -->
The [math/rand](/pkg/math/rand/) package now automatically seeds
the global random number generator
(used by top-level functions like `Float64` and `Int`) with a random value,
and the top-level [`Seed`](/pkg/math/rand/#Seed) function has been deprecated.
Programs that need a reproducible sequence of random numbers
should prefer to allocate their own random source, using `rand.New(rand.NewSource(seed))`.

Programs that need the earlier consistent global seeding behavior can set
`GODEBUG=randautoseed=0` in their environment.

<!-- https://go.dev/issue/20661 -->
The top-level [`Read`](/pkg/math/rand/#Read) function has been deprecated.
In almost all cases, [`crypto/rand.Read`](/pkg/crypto/rand/#Read) is more appropriate.

<!-- math/rand -->

#### [mime](/pkg/mime/)

<!-- https://go.dev/issue/48866 -->
The [`ParseMediaType`](/pkg/mime/#ParseMediaType) function now allows duplicate parameter names,
so long as the values of the names are the same.

<!-- mime -->

#### [mime/multipart](/pkg/mime/multipart/)

<!-- CL 431675 -->
Methods of the [`Reader`](/pkg/mime/multipart/#Reader) type now wrap errors
returned by the underlying `io.Reader`.

<!-- https://go.dev/issue/59153 -->
<!-- CL 481985 -->
In Go 1.19.8 and later, this package sets limits the size
of the MIME data it processes to protect against malicious inputs.
`Reader.NextPart` and `Reader.NextRawPart` limit the
number of headers in a part to 10000 and `Reader.ReadForm` limits
the total number of headers in all `FileHeaders` to 10000.
These limits may be adjusted with the `GODEBUG=multipartmaxheaders`
setting.
`Reader.ReadForm` further limits the number of parts in a form to 1000.
This limit may be adjusted with the `GODEBUG=multipartmaxparts`
setting.

<!-- mime/multipart -->

#### [net](/pkg/net/)

<!-- https://go.dev/issue/50101, CL 446179 -->
The [`LookupCNAME`](/pkg/net/#LookupCNAME)
function now consistently returns the contents
of a `CNAME` record when one exists. Previously on Unix systems and
when using the pure Go resolver, `LookupCNAME` would return an error
if a `CNAME` record referred to a name that with no `A`,
`AAAA`, or `CNAME` record. This change modifies
`LookupCNAME` to match the previous behavior on Windows,
allowing `LookupCNAME` to succeed whenever a
`CNAME` exists.

<!-- https://go.dev/issue/53482, CL 413454 -->
[`Interface.Flags`](/pkg/net/#Interface.Flags) now includes the new flag `FlagRunning`,
indicating an operationally active interface. An interface which is administratively
configured but not active (for example, because the network cable is not connected)
will have `FlagUp` set but not `FlagRunning`.

<!-- https://go.dev/issue/55301, CL 444955 -->
The new [`Dialer.ControlContext`](/pkg/net/#Dialer.ControlContext) field contains a callback function
similar to the existing [`Dialer.Control`](/pkg/net/#Dialer.Control) hook, that additionally
accepts the dial context as a parameter.
`Control` is ignored when `ControlContext` is not nil.

<!-- CL 428955 -->
The Go DNS resolver recognizes the `trust-ad` resolver option.
When `options trust-ad` is set in `resolv.conf`,
the Go resolver will set the AD bit in DNS queries. The resolver does not
make use of the AD bit in responses.

<!-- CL 448075 -->
DNS resolution will detect changes to `/etc/nsswitch.conf`
and reload the file when it changes. Checks are made at most once every
five seconds, matching the previous handling of `/etc/hosts`
and `/etc/resolv.conf`.

<!-- net -->

#### [net/http](/pkg/net/http/)

<!-- https://go.dev/issue/51914 -->
The [`ResponseWriter.WriteHeader`](/pkg/net/http/#ResponseWriter.WriteHeader) function now supports sending
`1xx` status codes.

<!-- https://go.dev/issue/41773, CL 356410 -->
The new [`Server.DisableGeneralOptionsHandler`](/pkg/net/http/#Server.DisableGeneralOptionsHandler) configuration setting
allows disabling the default `OPTIONS *` handler.

<!-- https://go.dev/issue/54299, CL 447216 -->
The new [`Transport.OnProxyConnectResponse`](/pkg/net/http/#Transport.OnProxyConnectResponse) hook is called
when a `Transport` receives an HTTP response from a proxy
for a `CONNECT` request.

<!-- https://go.dev/issue/53960, CL 418614  -->
The HTTP server now accepts HEAD requests containing a body,
rather than rejecting them as invalid.

<!-- https://go.dev/issue/53896 -->
HTTP/2 stream errors returned by `net/http` functions may be converted
to a [`golang.org/x/net/http2.StreamError`](/pkg/golang.org/x/net/http2/#StreamError) using
[`errors.As`](/pkg/errors/#As).

<!-- https://go.dev/cl/397734 -->
Leading and trailing spaces are trimmed from cookie names,
rather than being rejected as invalid.
For example, a cookie setting of "name =value"
is now accepted as setting the cookie "name".

<!-- https://go.dev/issue/52989 -->
A [`Cookie`](/pkg/net/http#Cookie) with an empty Expires field is now considered valid.
[`Cookie.Valid`](/pkg/net/http#Cookie.Valid) only checks Expires when it is set.

<!-- net/http -->

#### [net/netip](/pkg/net/netip/)

<!-- https://go.dev/issue/51766, https://go.dev/issue/51777, CL 412475 -->
The new [`IPv6LinkLocalAllRouters`](/pkg/net/netip/#IPv6LinkLocalAllRouters)
and [`IPv6Loopback`](/pkg/net/netip/#IPv6Loopback) functions
are the `net/netip` equivalents of
[`net.IPv6loopback`](/pkg/net/#IPv6loopback) and
[`net.IPv6linklocalallrouters`](/pkg/net/#IPv6linklocalallrouters).

<!-- net/netip -->

#### [os](/pkg/os/)

<!-- CL 448897 -->
On Windows, the name `NUL` is no longer treated as a special case in
[`Mkdir`](/pkg/os/#Mkdir) and
[`Stat`](/pkg/os/#Stat).

<!-- https://go.dev/issue/52747, CL 405275 -->
On Windows, [`File.Stat`](/pkg/os/#File.Stat)
now uses the file handle to retrieve attributes when the file is a directory.
Previously it would use the path passed to
[`Open`](/pkg/os/#Open), which may no longer be the file
represented by the file handle if the file has been moved or replaced.
This change modifies `Open` to open directories without the
`FILE_SHARE_DELETE` access, which match the behavior of regular files.

<!-- https://go.dev/issue/36019, CL 405275 -->
On Windows, [`File.Seek`](/pkg/os/#File.Seek) now supports
seeking to the beginning of a directory.

<!-- os -->

#### [os/exec](/pkg/os/exec/)

<!-- https://go.dev/issue/50436, CL 401835 -->
The new [`Cmd`](/pkg/os/exec/#Cmd) fields
[`Cancel`](/pkg/os/exec/#Cmd.Cancel) and
[`WaitDelay`](/pkg/os/exec/#Cmd.WaitDelay)
specify the behavior of the `Cmd` when its associated
`Context` is canceled or its process exits with I/O pipes still
held open by a child process.

<!-- os/exec -->

#### [path/filepath](/pkg/path/filepath/)

<!-- CL 363814, https://go.dev/issue/47209 -->
The new error [`SkipAll`](/pkg/path/filepath/#SkipAll)
terminates a [`Walk`](/pkg/path/filepath/#Walk)
immediately but successfully.

<!-- https://go.dev/issue/56219, CL 449239 -->
The new [`IsLocal`](/pkg/path/filepath/#IsLocal) function reports whether a path is
lexically local to a directory.
For example, if `IsLocal(p)` is `true`,
then `Open(p)` will refer to a file that is lexically
within the subtree rooted at the current directory.

<!-- io -->

#### [reflect](/pkg/reflect/)

<!-- https://go.dev/issue/46746, CL 423794 -->
The new [`Value.Comparable`](/pkg/reflect/#Value.Comparable) and
[`Value.Equal`](/pkg/reflect/#Value.Equal) methods
can be used to compare two `Value`s for equality.
`Comparable` reports whether `Equal` is a valid operation for a given `Value` receiver.

<!-- https://go.dev/issue/48000, CL 389635 -->
The new [`Value.Grow`](/pkg/reflect/#Value.Grow) method
extends a slice to guarantee space for another `n` elements.

<!-- https://go.dev/issue/52376, CL 411476 -->
The new [`Value.SetZero`](/pkg/reflect/#Value.SetZero) method
sets a value to be the zero value for its type.

<!-- CL 425184 -->
Go 1.18 introduced [`Value.SetIterKey`](/pkg/reflect/#Value.SetIterKey)
and [`Value.SetIterValue`](/pkg/reflect/#Value.SetIterValue) methods.
These are optimizations: `v.SetIterKey(it)` is meant to be equivalent to `v.Set(it.Key())`.
The implementations incorrectly omitted a check for use of unexported fields that was present in the unoptimized forms.
Go 1.20 corrects these methods to include the unexported field check.

<!-- reflect -->

#### [regexp](/pkg/regexp/)

<!-- CL 444817 -->
Go 1.19.2 and Go 1.18.7 included a security fix to the regular expression parser,
making it reject very large expressions that would consume too much memory.
Because Go patch releases do not introduce new API,
the parser returned [`syntax.ErrInternalError`](/pkg/regexp/syntax/#ErrInternalError) in this case.
Go 1.20 adds a more specific error, [`syntax.ErrLarge`](/pkg/regexp/syntax/#ErrLarge),
which the parser now returns instead.

<!-- regexp -->

#### [runtime/cgo](/pkg/runtime/cgo/)

<!-- https://go.dev/issue/46731, CL 421879 -->
Go 1.20 adds new [`Incomplete`](/pkg/runtime/cgo/#Incomplete) marker type.
Code generated by cgo will use `cgo.Incomplete` to mark an incomplete C type.

<!-- runtime/cgo -->

#### [runtime/metrics](/pkg/runtime/metrics/)

<!-- https://go.dev/issue/47216, https://go.dev/issue/49881 -->
Go 1.20 adds new [supported metrics](/pkg/runtime/metrics/#hdr-Supported_metrics),
including the current `GOMAXPROCS` setting (`/sched/gomaxprocs:threads`),
the number of cgo calls executed (`/cgo/go-to-c-calls:calls`),
total mutex block time (`/sync/mutex/wait/total:seconds`), and various measures of time
spent in garbage collection.

<!-- CL 427615 -->
Time-based histogram metrics are now less precise, but take up much less memory.

<!-- runtime/metrics -->

#### [runtime/pprof](/pkg/runtime/pprof/)

<!-- CL 443056 -->
Mutex profile samples are now pre-scaled, fixing an issue where old mutex profile
samples would be scaled incorrectly if the sampling rate changed during execution.

<!-- CL 416975 -->
Profiles collected on Windows now include memory mapping information that fixes
symbolization issues for position-independent binaries.

<!-- runtime/pprof -->

#### [runtime/trace](/pkg/runtime/trace/)

<!-- CL 447135, https://go.dev/issue/55022 -->
The garbage collector's background sweeper now yields less frequently,
resulting in many fewer extraneous events in execution traces.

<!-- runtime/trace -->

#### [strings](/pkg/strings/)

<!-- CL 407176, https://go.dev/issue/42537 -->
The new
[`CutPrefix`](/pkg/strings/#CutPrefix) and
[`CutSuffix`](/pkg/strings/#CutSuffix) functions
are like [`TrimPrefix`](/pkg/strings/#TrimPrefix)
and [`TrimSuffix`](/pkg/strings/#TrimSuffix)
but also report whether the string was trimmed.

<!-- strings -->

#### [sync](/pkg/sync/)

<!-- CL 399094, https://go.dev/issue/51972 -->
The new [`Map`](/pkg/sync/#Map) methods [`Swap`](/pkg/sync/#Map.Swap),
[`CompareAndSwap`](/pkg/sync/#Map.CompareAndSwap), and
[`CompareAndDelete`](/pkg/sync/#Map.CompareAndDelete)
allow existing map entries to be updated atomically.

<!-- sync -->

#### [syscall](/pkg/syscall/)

<!-- CL 411596 -->
On FreeBSD, compatibility shims needed for FreeBSD 11 and earlier have been removed.

<!-- CL 407574 -->
On Linux, additional [`CLONE_*`](/pkg/syscall/#CLONE_CLEAR_SIGHAND) constants
are defined for use with the [`SysProcAttr.Cloneflags`](/pkg/syscall/#SysProcAttr.Cloneflags) field.

<!-- CL 417695 -->
On Linux, the new [`SysProcAttr.CgroupFD`](/pkg/syscall/#SysProcAttr.CgroupFD)
and [`SysProcAttr.UseCgroupFD`](/pkg/syscall/#SysProcAttr.UseCgroupFD) fields
provide a way to place a child process into a specific cgroup.

<!-- syscall -->

#### [testing](/pkg/testing/)

<!-- https://go.dev/issue/43620, CL 420254 -->
The new method [`B.Elapsed`](/pkg/testing/#B.Elapsed)
reports the current elapsed time of the benchmark, which may be useful for
calculating rates to report with `ReportMetric`.

<!-- https://go.dev/issue/48515, CL 352349 -->
Calling [`T.Run`](/pkg/testing/#T.Run)
from a function passed
to [`T.Cleanup`](/pkg/testing/#T.Cleanup)
was never well-defined, and will now panic.

<!-- testing -->

#### [time](/pkg/time/)

<!-- https://go.dev/issue/52746, CL 412495 -->
The new time layout constants [`DateTime`](/pkg/time/#DateTime),
[`DateOnly`](/pkg/time/#DateOnly), and
[`TimeOnly`](/pkg/time/#TimeOnly)
provide names for three of the most common layout strings used in a survey of public Go source code.

<!-- CL 382734, https://go.dev/issue/50770 -->
The new [`Time.Compare`](/pkg/time/#Time.Compare) method
compares two times.

<!-- CL 425037 -->
[`Parse`](/pkg/time/#Parse)
now ignores sub-nanosecond precision in its input,
instead of reporting those digits as an error.

<!-- CL 444277 -->
The [`Time.MarshalJSON`](/pkg/time/#Time.MarshalJSON) method
is now more strict about adherence to RFC 3339.

<!-- time -->

#### [unicode/utf16](/pkg/unicode/utf16/)

<!-- https://go.dev/issue/51896, CL 409054 -->
The new [`AppendRune`](/pkg/unicode/utf16/#AppendRune)
function appends the UTF-16 encoding of a given rune to a uint16 slice,
analogous to [`utf8.AppendRune`](/pkg/unicode/utf8/#AppendRune).

<!-- unicode/utf16 -->

<!-- Silence false positives from x/build/cmd/relnote: -->
<!-- https://go.dev/issue/45964 was documented in Go 1.18 release notes but closed recently -->
<!-- https://go.dev/issue/52114 is an accepted proposal to add golang.org/x/net/http2.Transport.DialTLSContext; it's not a part of the Go release -->
<!-- CL 431335: cmd/api: make check pickier about api/*.txt -->
<!-- CL 447896 api: add newline to 55301.txt; modified api/next/55301.txt -->
<!-- CL 449215 api/next/54299: add missing newline; modified api/next/54299.txt -->
<!-- CL 433057 cmd: update vendored golang.org/x/tools for multiple error wrapping -->
<!-- CL 423362 crypto/internal/boring: update to newer boringcrypto, add arm64 -->
<!-- https://go.dev/issue/53481 x/cryptobyte ReadUint64, AddUint64 -->
<!-- https://go.dev/issue/51994 x/crypto/ssh -->
<!-- https://go.dev/issue/55358 x/exp/slices -->
<!-- https://go.dev/issue/54714 x/sys/unix -->
<!-- https://go.dev/issue/50035 https://go.dev/issue/54237 x/time/rate -->
<!-- CL 345488 strconv optimization -->
<!-- CL 428757 reflect deprecation, rolled back -->
<!-- https://go.dev/issue/49390 compile -l -N is fully supported -->
<!-- https://go.dev/issue/54619 x/tools -->
<!-- CL 448898 reverted -->
<!-- https://go.dev/issue/54850 x/net/http2 Transport.MaxReadFrameSize -->
<!-- https://go.dev/issue/56054 x/net/http2 SETTINGS_HEADER_TABLE_SIZE -->
<!-- CL 450375 reverted -->
<!-- CL 453259 tracking deprecations in api -->
<!-- CL 453260 tracking darwin port in api -->
<!-- CL 453615 fix deprecation comment in archive/tar -->
<!-- CL 453616 fix deprecation comment in archive/zip -->
<!-- CL 453617 fix deprecation comment in encoding/csv -->
<!-- https://go.dev/issue/54661 x/tools/go/analysis -->
<!-- CL 423359, https://go.dev/issue/51317 arena -->
