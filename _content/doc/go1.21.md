---
path: /doc/go1.21
template: false
title: Go 1.21 Release Notes
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

## Introduction to Go 1.21 {#introduction}

The latest Go release, version 1.21, arrives six months after [Go 1.20](/doc/go1.20).
Most of its changes are in the implementation of the toolchain, runtime, and libraries.
As always, the release maintains the Go 1 [promise of compatibility](/doc/go1compat);
in fact, Go 1.21 [improves upon that promise](#tools).
We expect almost all Go programs to continue to compile and run as before.

<!-- https://go.dev/issue/57631 -->
Go 1.21 introduces a small change to the numbering of releases.
In the past, we used Go 1._N_ to refer to both the overall Go language version and release family
as well as the first release in that family.
Starting in Go 1.21, the first release is now Go 1._N_.0.
Today we are releasing both the Go 1.21 language and its initial implementation, the Go 1.21.0 release.
These notes refer to “Go 1.21”; tools like `go` `version` will report “`go1.21.0`”
(until you upgrade to Go 1.21.1).
See “[Go versions](/doc/toolchain#version)” in the “Go Toolchains” documentation for details
about the new version numbering.

## Changes to the language {#language}

Go 1.21 adds three new built-ins to the language.

  - <!-- https://go.dev/issue/59488 -->
    The new functions `min` and `max` compute the
    smallest (or largest, for `max`) value of a fixed number
    of given arguments.
    See the language spec for
    [details](/ref/spec#Min_and_max).
  - <!-- https://go.dev/issue/56351 -->
    The new function `clear` deletes all elements from a
    map or zeroes all elements of a slice.
    See the language spec for
    [details](/ref/spec#Clear).


<!-- https://go.dev/issue/57411 -->
Package initialization order is now specified more precisely. The
new algorithm is:

- Sort all packages by import path.
- Repeat until the list of packages is empty:
  - Find the first package in the list for which all imports are
    already initialized.
  - Initialize that package and remove it from the list.

This may change the behavior of some programs that rely on a
specific initialization ordering that was not expressed by explicit
imports. The behavior of such programs was not well defined by the
spec in past releases. The new rule provides an unambiguous definition.


Multiple improvements that increase the power and precision of type inference have been made.

  - <!-- https://go.dev/issue/59338 -->
    A (possibly partially instantiated generic) function may now be called with arguments that are
    themselves (possibly partially instantiated) generic functions.
    The compiler will attempt to infer the missing type arguments of the callee (as before) and,
    for each argument that is a generic function that is not fully instantiated,
    its missing type arguments (new).
    Typical use cases are calls to generic functions operating on containers
    (such as [slices.IndexFunc](/pkg/slices#IndexFunc)) where a function argument
    may also be generic, and where the type argument of the called function and its arguments
    are inferred from the container type.
    More generally, a generic function may now be used without explicit instantiation when
    it is assigned to a variable or returned as a result value if the type arguments can
    be inferred from the assignment.
  - <!-- https://go.dev/issue/60353, https://go.dev/issue/57192, https://go.dev/issue/52397, https://go.dev/issue/41176 -->
    Type inference now also considers methods when a value is assigned to an interface:
    type arguments for type parameters used in method signatures may be inferred from
    the corresponding parameter types of matching methods.
  - <!-- https://go.dev/issue/51593 https://go.dev/issue/39661 -->
    Similarly, since a type argument must implement all the methods of its corresponding constraint,
    the methods of the type argument and constraint are matched which may lead to the inference of
    additional type arguments.
  - <!-- https://go.dev/issue/58671 -->
    If multiple untyped constant arguments of different kinds (such as an untyped int and
    an untyped floating-point constant) are passed to parameters with the same (not otherwise
    specified) type parameter type, instead of an error, now type inference determines the
    type using the same approach as an operator with untyped constant operands.
    This change brings the types inferred from untyped constant arguments in line with the
    types of constant expressions.
  - <!-- https://go.dev/issue/59750 -->
    Type inference is now precise when matching corresponding types in assignments:
    component types (such as the elements of slices, or the parameter types in function signatures)
    must be identical (given suitable type arguments) to match, otherwise inference fails.
    This change produces more accurate error messages:
    where in the past type inference may have succeeded incorrectly and lead to an invalid assignment,
    the compiler now reports an inference error if two types can't possibly match.

<!-- https://go.dev/issue/58650 -->
More generally, the description of
[type inference](/ref/spec#Type_inference)
in the language spec has been clarified.
Together, all these changes make type inference more powerful and inference failures less surprising.

<!-- https://go.dev/issue/57969 -->

Go 1.21 includes a preview of a language change we are considering for a future version of Go:
making for loop variables per-iteration instead of per-loop, to avoid accidental sharing bugs.
For details about how to try that language change, see [the LoopvarExperiment wiki page](/wiki/LoopvarExperiment).

<!-- https://go.dev/issue/25448 -->

Go 1.21 now defines that if a goroutine is panicking and recover was called directly by a deferred
function, the return value of recover is guaranteed not to be nil. To ensure this, calling panic
with a nil interface value (or an untyped nil) causes a run-time panic of type
[`*runtime.PanicNilError`](/pkg/runtime/#PanicNilError).

To support programs written for older versions of Go, nil panics can be re-enabled by setting
`GODEBUG=panicnil=1`.
This setting is enabled automatically when compiling a program whose main package
is in a module that declares `go` `1.20` or earlier.

## Tools {#tools}

Go 1.21 adds improved support for backwards compatibility and forwards compatibility
in the Go toolchain.

<!-- https://go.dev/issue/56986 -->
To improve backwards compatibility, Go 1.21 formalizes
Go's use of the GODEBUG environment variable to control
the default behavior for changes that are non-breaking according to the
[compatibility policy](/doc/go1compat)
but nonetheless may cause existing programs to break.
(For example, programs that depend on buggy behavior may break
when a bug is fixed, but bug fixes are not considered breaking changes.)
When Go must make this kind of behavior change,
it now chooses between the old and new behavior based on the
`go` line in the workspace's `go.work` file
or else the main module's `go.mod` file.
Upgrading to a new Go toolchain but leaving the `go` line
set to its original (older) Go version preserves the behavior of the older
toolchain.
With this compatibility support, the latest Go toolchain should always
be the best, most secure, implementation of an older version of Go.
See “[Go, Backwards Compatibility, and GODEBUG](/doc/godebug)” for details.

<!-- https://go.dev/issue/57001 -->
To improve forwards compatibility, Go 1.21 now reads the `go` line
in a `go.work` or `go.mod` file as a strict
minimum requirement: `go` `1.21.0` means
that the workspace or module cannot be used with Go 1.20 or with Go 1.21rc1.
This allows projects that depend on fixes made in later versions of Go
to ensure that they are not used with earlier versions.
It also gives better error reporting for projects that make use of new Go features:
when the problem is that a newer Go version is needed,
that problem is reported clearly, instead of attempting to build the code
and printing errors about unresolved imports or syntax errors.

To make these new stricter version requirements easier to manage,
the `go` command can now invoke not just the toolchain
bundled in its own release but also other Go toolchain versions found in the PATH
or downloaded on demand.
If a `go.mod` or `go.work` `go` line
declares a minimum requirement on a newer version of Go, the `go`
command will find and run that version automatically.
The new `toolchain` directive sets a suggested minimum toolchain to use,
which may be newer than the strict `go` minimum.
See “[Go Toolchains](/doc/toolchain)” for details.

### Go command {#go-command}

<!-- https://go.dev/issue/58099, CL 474236 -->
The `-pgo` build flag now defaults to `-pgo=auto`,
and the restriction of specifying a single main package on the command
line is now removed. If a file named `default.pgo` is present
in the main package's directory, the `go` command will use
it to enable profile-guided optimization for building the corresponding
program.

The `-C` `dir` flag must now be the first
flag on the command-line when used.

<!-- https://go.dev/issue/37708, CL 463837 -->
The new `go` `test` option
`-fullpath` prints full path names in test log messages,
rather than just base names.

<!-- https://go.dev/issue/15513, CL 466397 -->
The `go` `test` `-c` flag now
supports writing test binaries for multiple packages, each to
`pkg.test` where `pkg` is the package name.
It is an error if more than one test package being compiled has a given package name.

<!-- https://go.dev/issue/15513, CL 466397 -->
The `go` `test` `-o` flag now
accepts a directory argument, in which case test binaries are written to that
directory instead of the current directory.

<!-- https://go.dev/issue/31544, CL 477839 -->
When using an external (C) linker with cgo enabled, the `runtime/cgo` package is
now supplied to the Go linker as an additional dependency to ensure that the Go
runtime is compatible with any additional libraries added by the C linker.

### Cgo {#cgo}

<!-- CL 490819 -->
In files that `import "C"`, the Go toolchain now
correctly reports errors for attempts to declare Go methods on C types.

## Runtime {#runtime-changes}

<!-- https://go.dev/issue/7181 -->
When printing very deep stacks, the runtime now prints the first 50
(innermost) frames followed by the bottom 50 (outermost) frames,
rather than just printing the first 100 frames. This makes it easier
to see how deeply recursive stacks started, and is especially
valuable for debugging stack overflows.

<!-- https://go.dev/issue/59960 -->
On Linux platforms that support transparent huge pages, the Go runtime
now manages which parts of the heap may be backed by huge pages more
explicitly. This leads to better utilization of memory: small heaps
should see less memory used (up to 50% in pathological cases) while
large heaps should see fewer broken huge pages for dense parts of the
heap, improving CPU usage and latency by up to 1%. A consequence of this
change is that the runtime no longer tries to work around a particular
problematic Linux configuration setting, which may result in higher
memory overheads. The recommended fix is to adjust the OS's huge page
settings according to the [GC guide](/doc/gc-guide#Linux_transparent_huge_pages).
However, other workarounds are available as well. See the [section on
`max_ptes_none`](/doc/gc-guide#Linux_THP_max_ptes_none_workaround).

<!-- https://go.dev/issue/57069, https://go.dev/issue/56966 -->
As a result of runtime-internal garbage collection tuning,
applications may see up to a 40% reduction in application tail latency
and a small decrease in memory use. Some applications may also observe
a small loss in throughput.
The memory use decrease should be proportional to the loss in
throughput, such that the previous release's throughput/memory
tradeoff may be recovered (with little change to latency) by
increasing `GOGC` and/or `GOMEMLIMIT` slightly.

<!-- https://go.dev/issue/51676 -->
Calls from C to Go on threads created in C require some setup to prepare for
Go execution. On Unix platforms, this setup is now preserved across multiple
calls from the same thread. This significantly reduces the overhead of
subsequent C to Go calls from ~1-3 microseconds per call to ~100-200
nanoseconds per call.

## Compiler {#compiler}

Profile-guide optimization (PGO), added as a preview in Go 1.20, is now ready
for general use. PGO enables additional optimizations on code identified as
hot by profiles of production workloads. As mentioned in the
[Go command section](#go-command), PGO is enabled by default for
binaries that contain a `default.pgo` profile in the main
package directory. Performance improvements vary depending on application
behavior, with most programs from a representative set of Go programs seeing
between 2 and 7% improvement from enabling PGO. See the
[PGO user guide](/doc/pgo) for detailed documentation.

<!-- https://go.dev/issue/59959 -->

PGO builds can now devirtualize some interface method calls, adding a
concrete call to the most common callee. This enables further optimization,
such as inlining the callee.

<!-- CL 497455 -->

Go 1.21 improves build speed by up to 6%, largely thanks to building the
compiler itself with PGO.

## Assembler {#assembler}

<!-- https://go.dev/issue/58378 -->

On amd64, frameless nosplit assembly functions are no longer automatically marked as `NOFRAME`.
Instead, the `NOFRAME` attribute must be explicitly specified if desired,
which is already the behavior on other architectures supporting frame pointers.
With this, the runtime now maintains the frame pointers for stack transitions.

<!-- CL 476295 -->

The verifier that checks for incorrect uses of `R15` when dynamic linking on amd64 has been improved.

## Linker {#linker}

<!-- https://go.dev/issue/57302, CL 461749, CL 457455 -->
On windows/amd64, the linker (with help from the compiler) now emits
SEH unwinding data by default, which improves the integration
of Go applications with Windows debuggers and other tools.

<!-- CL 463395, CL 461315 -->

In Go 1.21 the linker (with help from the compiler) is now capable of
deleting dead (unreferenced) global map variables, if the number of
entries in the variable initializer is sufficiently large, and if the
initializer expressions are side-effect free.

## Standard library {#library}

### New log/slog package {#slog}

<!--
 https://go.dev/issue/59060, https://go.dev/issue/59141, https://go.dev/issue/59204, https://go.dev/issue/59280,
        https://go.dev/issue/59282, https://go.dev/issue/59339, https://go.dev/issue/59345, https://go.dev/issue/61200,
        CL 477295, CL 484096, CL 486376, CL 486415, CL 487855, CL 508195
 -->
The new [log/slog](/pkg/log/slog) package provides structured logging with levels.
Structured logging emits key-value pairs
to enable fast, accurate processing of large amounts of log data.
The package supports integration with popular log analysis tools and services.

### New testing/slogtest package {#slogtest}

<!-- CL 487895 -->
The new [testing/slogtest](/pkg/testing/slogtest) package can help
to validate [slog.Handler](/pkg/log/slog#Handler) implementations.

### New slices package {#slices}

<!-- https://go.dev/issue/45955, https://go.dev/issue/54768 -->
<!-- https://go.dev/issue/57348, https://go.dev/issue/57433 -->
<!-- https://go.dev/issue/58565, https://go.dev/issue/60091 -->
<!-- https://go.dev/issue/60546 -->
<!-- CL 467417, CL 468855, CL 483175, CL 496078, CL 498175, CL 502955 -->
The new [slices](/pkg/slices) package provides many common
operations on slices, using generic functions that work with slices
of any element type.

### New maps package {#maps}

<!-- https://go.dev/issue/57436, CL 464343 -->
The new [maps](/pkg/maps/) package provides several
common operations on maps, using generic functions that work with
maps of any key or element type.

### New cmp package {#cmp}

<!-- https://go.dev/issue/59488, CL 496356 -->
The new [cmp](/pkg/cmp/) package defines the type
constraint [`Ordered`](/pkg/cmp/#Ordered) and
two new generic functions
[`Less`](/pkg/cmp/#Less)
and [`Compare`](/pkg/cmp/#Compare) that are
useful with [ordered
types](/ref/spec/#Comparison_operators).

### Minor changes to the library {#minor_library_changes}

As always, there are various minor changes and updates to the library,
made with the Go 1 [promise of compatibility](/doc/go1compat)
in mind.
There are also various performance improvements, not enumerated here.

#### [archive/tar](/pkg/archive/tar/)

<!-- https://go.dev/issue/54451, CL 491175 -->
The implementation of the
[`io/fs.FileInfo`](/pkg/io/fs/#FileInfo)
interface returned by
[`Header.FileInfo`](/pkg/archive/tar/#Header.FileInfo)
now implements a `String` method that calls
[`io/fs.FormatFileInfo`](/pkg/io/fs/#FormatFileInfo).

<!-- archive/tar -->

#### [archive/zip](/pkg/archive/zip/)

<!-- https://go.dev/issue/54451, CL 491175 -->
The implementation of the
[`io/fs.FileInfo`](/pkg/io/fs/#FileInfo)
interface returned by
[`FileHeader.FileInfo`](/pkg/archive/zip/#FileHeader.FileInfo)
now implements a `String` method that calls
[`io/fs.FormatFileInfo`](/pkg/io/fs/#FormatFileInfo).

<!-- https://go.dev/issue/54451, CL 491175 -->
The implementation of the
[`io/fs.DirEntry`](/pkg/io/fs/#DirEntry)
interface returned by the
[`io/fs.ReadDirFile.ReadDir`](/pkg/io/fs/#ReadDirFile.ReadDir)
method of the
[`io/fs.File`](/pkg/io/fs/#File)
returned by
[`Reader.Open`](/pkg/archive/zip/#Reader.Open)
now implements a `String` method that calls
[`io/fs.FormatDirEntry`](/pkg/io/fs/#FormatDirEntry).

<!-- archive/zip -->

#### [bytes](/pkg/bytes/)

<!-- https://go.dev/issue/53685, CL 474635 -->
The [`Buffer`](/pkg/bytes/#Buffer) type
has two new methods:
[`Available`](/pkg/bytes/#Buffer.Available)
and [`AvailableBuffer`](/pkg/bytes/#Buffer.AvailableBuffer).
These may be used along with the
[`Write`](/pkg/bytes/#Buffer.Write)
method to append directly to the `Buffer`.

<!-- bytes -->

#### [context](/pkg/context/)

<!-- https://go.dev/issue/40221, CL 479918 -->
The new [`WithoutCancel`](/pkg/context/#WithoutCancel)
function returns a copy of a context that is not canceled when the original
context is canceled.

<!-- https://go.dev/issue/56661, CL 449318 -->
The new [`WithDeadlineCause`](/pkg/context/#WithDeadlineCause)
and [`WithTimeoutCause`](/pkg/context/#WithTimeoutCause)
functions provide a way to set a context cancellation cause when a deadline or
timer expires. The cause may be retrieved with the
[`Cause`](/pkg/context/#Cause) function.

<!-- https://go.dev/issue/57928, CL 482695 -->
The new [`AfterFunc`](/pkg/context/#AfterFunc)
function registers a function to run after a context has been cancelled.

<!-- CL 455455 -->
An optimization means that the results of calling
[`Background`](/pkg/context/#Background)
and [`TODO`](/pkg/context/#TODO) and
converting them to a shared type can be considered equal.
In previous releases they were always different. Comparing
[`Context`](/pkg/context/#Context) values
for equality has never been well-defined, so this is not
considered to be an incompatible change.

#### [crypto/ecdsa](/pkg/crypto/ecdsa/)

<!-- CL 492955 -->
[`PublicKey.Equal`](/pkg/crypto/ecdsa/#PublicKey.Equal) and
[`PrivateKey.Equal`](/pkg/crypto/ecdsa/#PrivateKey.Equal)
now execute in constant time.

<!-- crypto/ecdsa -->

#### [crypto/elliptic](/pkg/crypto/elliptic/)

<!-- CL 459977 -->
All of the [`Curve`](/pkg/crypto/elliptic/#Curve) methods have been deprecated, along with [`GenerateKey`](/pkg/crypto/elliptic/#GenerateKey), [`Marshal`](/pkg/crypto/elliptic/#Marshal), and [`Unmarshal`](/pkg/crypto/elliptic/#Unmarshal). For ECDH operations, the new [`crypto/ecdh`](/pkg/crypto/ecdh/) package should be used instead. For lower-level operations, use third-party modules such as [filippo.io/nistec](https://pkg.go.dev/filippo.io/nistec).

<!-- crypto/elliptic -->

#### [crypto/rand](/pkg/crypto/rand/)

<!-- CL 463123 -->
The [`crypto/rand`](/pkg/crypto/rand/) package now uses the `getrandom` system call on NetBSD 10.0 and later.

<!-- crypto/rand -->

#### [crypto/rsa](/pkg/crypto/rsa/)

<!-- CL 471259, CL 492935 -->
The performance of private RSA operations (decryption and signing) is now better than Go 1.19 for `GOARCH=amd64` and `GOARCH=arm64`. It had regressed in Go 1.20.

Due to the addition of private fields to [`PrecomputedValues`](/pkg/crypto/rsa/#PrecomputedValues), [`PrivateKey.Precompute`](/pkg/crypto/rsa/#PrivateKey.Precompute) must be called for optimal performance even if deserializing (for example from JSON) a previously-precomputed private key.

<!-- CL 492955 -->
[`PublicKey.Equal`](/pkg/crypto/rsa/#PublicKey.Equal) and
[`PrivateKey.Equal`](/pkg/crypto/rsa/#PrivateKey.Equal)
now execute in constant time.

<!-- https://go.dev/issue/56921, CL 459976 -->
The [`GenerateMultiPrimeKey`](/pkg/crypto/rsa/#GenerateMultiPrimeKey) function and the [`PrecomputedValues.CRTValues`](/pkg/crypto/rsa/#PrecomputedValues.CRTValues) field have been deprecated. [`PrecomputedValues.CRTValues`](/pkg/crypto/rsa/#PrecomputedValues.CRTValues) will still be populated when [`PrivateKey.Precompute`](/pkg/crypto/rsa/#PrivateKey.Precompute) is called, but the values will not be used during decryption operations.

<!-- crypto/rsa -->

<!-- CL 483815 reverted -->

#### [crypto/sha256](/pkg/crypto/sha256/)

<!-- https://go.dev/issue/50543, CL 408795 -->
SHA-224 and SHA-256 operations now use native instructions when available when `GOARCH=amd64`, providing a performance improvement on the order of 3-4x.

<!-- crypto/sha256 -->

<!-- CL 481478 reverted -->
<!-- CL 483816 reverted -->

#### [crypto/tls](/pkg/crypto/tls/)

<!-- CL 497895 -->
Servers now skip verifying client certificates (including not running
[`Config.VerifyPeerCertificate`](/pkg/crypto/tls/#Config.VerifyPeerCertificate))
for resumed connections, besides checking the expiration time. This makes
session tickets larger when client certificates are in use. Clients were
already skipping verification on resumption, but now check the expiration
time even if [`Config.InsecureSkipVerify`](/pkg/crypto/tls/#Config.InsecureSkipVerify)
is set.

<!-- https://go.dev/issue/60105, CL 496818, CL 496820, CL 496822, CL 496821, CL 501675 -->
Applications can now control the content of session tickets.

  - The new [`SessionState`](/pkg/crypto/tls/#SessionState) type
    describes a resumable session.
  - The [`SessionState.Bytes`](/pkg/crypto/tls/#SessionState.Bytes)
    method and [`ParseSessionState`](/pkg/crypto/tls/#ParseSessionState)
    function serialize and deserialize a `SessionState`.
  - The [`Config.WrapSession`](/pkg/crypto/tls/#Config.WrapSession) and
    [`Config.UnwrapSession`](/pkg/crypto/tls/#Config.UnwrapSession)
    hooks convert a `SessionState` to and from a ticket on the server side.
  - The [`Config.EncryptTicket`](/pkg/crypto/tls/#Config.EncryptTicket)
    and [`Config.DecryptTicket`](/pkg/crypto/tls/#Config.DecryptTicket)
    methods provide a default implementation of `WrapSession` and
    `UnwrapSession`.
  - The [`ClientSessionState.ResumptionState`](/pkg/crypto/tls/#ClientSessionState.ResumptionState) method and
    [`NewResumptionState`](/pkg/crypto/tls/#NewResumptionState) function
    may be used by a `ClientSessionCache` implementation to store and
    resume sessions on the client side.


<!-- CL 496817 -->
To reduce the potential for session tickets to be used as a tracking
mechanism across connections, the server now issues new tickets on every
resumption (if they are supported and not disabled) and tickets don't bear
an identifier for the key that encrypted them anymore. If passing a large
number of keys to [`Conn.SetSessionTicketKeys`](/pkg/crypto/tls/#Conn.SetSessionTicketKeys),
this might lead to a noticeable performance cost.

<!-- CL 497376 -->
Both clients and servers now implement the Extended Master Secret extension (RFC 7627).
The deprecation of [`ConnectionState.TLSUnique`](/pkg/crypto/tls/#ConnectionState.TLSUnique)
has been reverted, and is now set for resumed connections that support Extended Master Secret.

<!-- https://go.dev/issue/44886, https://go.dev/issue/60107, CL 493655, CL 496995, CL 514997 -->
The new [`QUICConn`](/pkg/crypto/tls/#QUICConn) type
provides support for QUIC implementations, including 0-RTT support. Note
that this is not itself a QUIC implementation, and 0-RTT is still not
supported in TLS.

<!-- https://go.dev/issue/46308, CL 497377 -->
The new [`VersionName`](/pkg/crypto/tls/#VersionName) function
returns the name for a TLS version number.

<!-- https://go.dev/issue/52113, CL 410496 -->
The TLS alert codes sent from the server for client authentication failures have
been improved. Previously, these failures always resulted in a "bad certificate" alert.
Now, certain failures will result in more appropriate alert codes,
as defined by RFC 5246 and RFC 8446:

  - For TLS 1.3 connections, if the server is configured to require client authentication using
    [RequireAnyClientCert](/pkg/crypto/tls/#RequireAnyClientCert) or
    [RequireAndVerifyClientCert](/pkg/crypto/tls/#RequireAndVerifyClientCert),
    and the client does not provide any certificate, the server will now return the "certificate required" alert.
  - If the client provides a certificate that is not signed by the set of trusted certificate authorities
    configured on the server, the server will return the "unknown certificate authority" alert.
  - If the client provides a certificate that is either expired or not yet valid,
    the server will return the "expired certificate" alert.
  - In all other scenarios related to client authentication failures, the server still returns "bad certificate".


<!-- crypto/tls -->

#### [crypto/x509](/pkg/crypto/x509/)

<!-- https://go.dev/issue/53573, CL 468875 -->
[`RevocationList.RevokedCertificates`](/pkg/crypto/x509/#RevocationList.RevokedCertificates) has been deprecated and replaced with the new [`RevokedCertificateEntries`](/pkg/crypto/x509/#RevocationList.RevokedCertificateEntries) field, which is a slice of [`RevocationListEntry`](/pkg/crypto/x509/#RevocationListEntry). [`RevocationListEntry`](/pkg/crypto/x509/#RevocationListEntry) contains all of the fields in [`pkix.RevokedCertificate`](/pkg/crypto/x509/pkix#RevokedCertificate), as well as the revocation reason code.

<!-- CL 478216 -->
Name constraints are now correctly enforced on non-leaf certificates, and
not on the certificates where they are expressed.

<!-- crypto/x509 -->

#### [debug/elf](/pkg/debug/elf/)

<!-- https://go.dev/issue/56892, CL 452617 -->
The new
[`File.DynValue`](/pkg/debug/elf/#File.DynValue)
method may be used to retrieve the numeric values listed with a
given dynamic tag.

<!-- https://go.dev/issue/56887, CL 452496 -->
The constant flags permitted in a `DT_FLAGS_1`
dynamic tag are now defined with type
[`DynFlag1`](/pkg/debug/elf/#DynFlag1). These
tags have names starting with `DF_1`.

<!-- CL 473256 -->
The package now defines the constant
[`COMPRESS_ZSTD`](/pkg/debug/elf/#COMPRESS_ZSTD).

<!-- https://go.dev/issue/60348, CL 496918 -->
The package now defines the constant
[`R_PPC64_REL24_P9NOTOC`](/pkg/debug/elf/#R_PPC64_REL24_P9NOTOC).

<!-- debug/elf -->

#### [debug/pe](/pkg/debug/pe/)

<!-- CL 488475 -->
Attempts to read from a section containing uninitialized data
using
[`Section.Data`](/pkg/debug/pe/#Section.Data)
or the reader returned by [`Section.Open`](/pkg/debug/pe/#Section.Open)
now return an error.

<!-- debug/pe -->

#### [embed](/pkg/embed/)

<!-- https://go.dev/issue/57803, CL 483235 -->
The [`io/fs.File`](/pkg/io/fs/#File)
returned by
[`FS.Open`](/pkg/embed/#FS.Open) now
has a `ReadAt` method that
implements [`io.ReaderAt`](/pkg/io/#ReaderAt).

<!-- https://go.dev/issue/54451, CL 491175 -->
Calling <code>[FS.Open](/pkg/embed/FS.Open).[Stat](/pkg/io/fs/#File.Stat)</code>
will return a type that now implements a `String`
method that calls
[`io/fs.FormatFileInfo`](/pkg/io/fs/#FormatFileInfo).

<!-- embed -->

#### [encoding/binary](/pkg/encoding/binary/)

<!-- https://go.dev/issue/57237, CL 463218, CL 463985 -->
The new
[`NativeEndian`](/pkg/encoding/binary/#NativeEndian)
variable may be used to convert between byte slices and integers
using the current machine's native endianness.

<!-- encoding/binary -->

#### [errors](/pkg/errors/)

<!-- https://go.dev/issue/41198, CL 473935 -->
The new
[`ErrUnsupported`](/pkg/errors/#ErrUnsupported)
error provides a standardized way to indicate that a requested
operation may not be performed because it is unsupported.
For example, a call to
[`os.Link`](/pkg/os/#Link) when using a
file system that does not support hard links.

<!-- errors -->

#### [flag](/pkg/flag/)

<!-- https://go.dev/issue/53747, CL 476015 -->
The new [`BoolFunc`](/pkg/flag/#BoolFunc)
function and
[`FlagSet.BoolFunc`](/pkg/flag/#FlagSet.BoolFunc)
method define a flag that does not require an argument and calls
a function when the flag is used. This is similar to
[`Func`](/pkg/flag/#Func) but for a
boolean flag.

<!-- CL 480215 -->
A flag definition
(via [`Bool`](/pkg/flag/#Bool),
[`BoolVar`](/pkg/flag/#BoolVar),
[`Int`](/pkg/flag/#Int),
[`IntVar`](/pkg/flag/#IntVar), etc.)
will panic if [`Set`](/pkg/flag/#Set) has
already been called on a flag with the same name. This change is
intended to detect cases where [changes in
initialization order](#language) cause flag operations to occur in a
different order than expected. In many cases the fix to this
problem is to introduce a explicit package dependence to
correctly order the definition before any
[`Set`](/pkg/flag/#Set) operations.

<!-- flag -->

#### [go/ast](/pkg/go/ast/)

<!-- https://go.dev/issue/28089, CL 487935 -->
The new [`IsGenerated`](/pkg/go/ast/#IsGenerated) predicate
reports whether a file syntax tree contains the
[special comment](/s/generatedcode)
that conventionally indicates that the file was generated by a tool.

<!-- https://go.dev/issue/59033, CL 476276 -->
The new
[`File.GoVersion`](/pkg/go/ast/#File.GoVersion)
field records the minimum Go version required by
any `//go:build` or `// +build`
directives.

<!-- go/ast -->

#### [go/build](/pkg/go/build/)

<!-- https://go.dev/issue/56986, CL 453603 -->
The package now parses build directives (comments that start
with `//go:`) in file headers (before
the `package` declaration). These directives are
available in the new
[`Package`](/pkg/go/build#Package) fields
[`Directives`](/pkg/go/build#Package.Directives),
[`TestDirectives`](/pkg/go/build#Package.TestDirectives),
and
[`XTestDirectives`](/pkg/go/build#Package.XTestDirectives).

<!-- go/build -->

#### [go/build/constraint](/pkg/go/build/constraint/)

<!-- https://go.dev/issue/59033, CL 476275 -->
The new
[`GoVersion`](/pkg/go/build/constraint/#GoVersion)
function returns the minimum Go version implied by a build
expression.

<!-- go/build/constraint -->

#### [go/token](/pkg/go/token/)

<!-- https://go.dev/issue/57708, CL 464515 -->
The new [`File.Lines`](/pkg/go/token/#File.Lines) method
returns the file's line-number table in the same form as accepted by
`File.SetLines`.

<!-- go/token -->

#### [go/types](/pkg/go/types/)

<!-- https://go.dev/issue/61175, CL 507975 -->
The new [`Package.GoVersion`](/pkg/go/types/#Package.GoVersion)
method returns the Go language version used to check the package.

<!-- go/types -->

#### [hash/maphash](/pkg/hash/maphash/)

<!-- https://go.dev/issue/47342, CL 468795 -->
The `hash/maphash` package now has a pure Go implementation, selectable with the `purego` build tag.

<!-- hash/maphash -->

#### [html/template](/pkg/html/template/)

<!-- https://go.dev/issue/59584, CL 496395 -->
The new error
[`ErrJSTemplate`](/pkg/html/template/#ErrJSTemplate)
is returned when an action appears in a JavaScript template
literal. Previously an unexported error was returned.

<!-- html/template -->

#### [io/fs](/pkg/io/fs/)

<!-- https://go.dev/issue/54451, CL 489555 -->
The new
[`FormatFileInfo`](/pkg/io/fs/#FormatFileInfo)
function returns a formatted version of a
[`FileInfo`](/pkg/io/fs/#FileInfo).
The new
[`FormatDirEntry`](/pkg/io/fs/#FormatDirEntry)
function returns a formatted version of a
[`DirEntry`](/pkg/io/fs/#FileInfo).
The implementation of
[`DirEntry`](/pkg/io/fs/#DirEntry)
returned by
[`ReadDir`](/pkg/io/fs/#ReadDir) now
implements a `String` method that calls
[`FormatDirEntry`](/pkg/io/fs/#FormatDirEntry),
and the same is true for
the [`DirEntry`](/pkg/io/fs/#DirEntry)
value passed to
[`WalkDirFunc`](/pkg/io/fs/#WalkDirFunc).

<!-- io/fs -->

<!-- https://go.dev/issue/56491 rolled back by https://go.dev/issue/60519 -->
<!-- CL 459435 reverted by CL 467255 -->
<!-- CL 467515 reverted by CL 499416 -->

#### [math/big](/pkg/math/big/)

<!-- https://go.dev/issue/56984, CL 453115, CL 500116 -->
The new [`Int.Float64`](/pkg/math/big/#Int.Float64)
method returns the nearest floating-point value to a
multi-precision integer, along with an indication of any
rounding that occurred.

<!-- math/big -->

#### [net](/pkg/net/)

<!-- https://go.dev/issue/59166, https://go.dev/issue/56539 -->
<!-- CL 471136, CL 471137, CL 471140 -->
On Linux, the [net](/pkg/net/) package can now use
Multipath TCP when the kernel supports it. It is not used by
default. To use Multipath TCP when available on a client, call
the
[`Dialer.SetMultipathTCP`](/pkg/net/#Dialer.SetMultipathTCP)
method before calling the
[`Dialer.Dial`](/pkg/net/#Dialer.Dial) or
[`Dialer.DialContext`](/pkg/net/#Dialer.DialContext)
methods. To use Multipath TCP when available on a server, call
the
[`ListenConfig.SetMultipathTCP`](/pkg/net/#ListenConfig.SetMultipathTCP)
method before calling the
[`ListenConfig.Listen`](/pkg/net/#ListenConfig.Listen)
method. Specify the network as `"tcp"` or
`"tcp4"` or `"tcp6"` as usual. If
Multipath TCP is not supported by the kernel or the remote host,
the connection will silently fall back to TCP. To test whether a
particular connection is using Multipath TCP, use the
[`TCPConn.MultipathTCP`](/pkg/net/#TCPConn.MultipathTCP)
method.

In a future Go release we may enable Multipath TCP by default on
systems that support it.

<!-- net -->

#### [net/http](/pkg/net/http/)

<!-- CL 472636 -->
The new [`ResponseController.EnableFullDuplex`](/pkg/net/http#ResponseController.EnableFullDuplex)
method allows server handlers to concurrently read from an HTTP/1
request body while writing the response. Normally, the HTTP/1 server
automatically consumes any remaining request body before starting to
write the response, to avoid deadlocking clients which attempt to
write a complete request before reading the response. The
`EnableFullDuplex` method disables this behavior.

<!-- https://go.dev/issue/44855, CL 382117 -->
The new [`ErrSchemeMismatch`](/pkg/net/http/#ErrSchemeMismatch) error is returned by [`Client`](/pkg/net/http/#Client) and [`Transport`](/pkg/net/http/#Transport) when the server responds to an HTTPS request with an HTTP response.

<!-- CL 494122 -->
The [net/http](/pkg/net/http/) package now supports
[`errors.ErrUnsupported`](/pkg/errors/#ErrUnsupported),
in that the expression
`errors.Is(http.ErrNotSupported, errors.ErrUnsupported)`
will return true.

<!-- net/http -->

#### [os](/pkg/os/)

<!-- https://go.dev/issue/32558, CL 219638 -->
Programs may now pass an empty `time.Time` value to
the [`Chtimes`](/pkg/os/#Chtimes) function
to leave either the access time or the modification time unchanged.

<!-- CL 480135 -->
On Windows the
[`File.Chdir`](/pkg/os#File.Chdir) method
now changes the current directory to the file, rather than
always returning an error.

<!-- CL 495079 -->
On Unix systems, if a non-blocking descriptor is passed
to [`NewFile`](/pkg/os/#NewFile), calling
the [`File.Fd`](/pkg/os/#File.Fd) method
will now return a non-blocking descriptor. Previously the
descriptor was converted to blocking mode.

<!-- CL 477215 -->
On Windows calling
[`Truncate`](/pkg/os/#Truncate) on a
non-existent file used to create an empty file. It now returns
an error indicating that the file does not exist.

<!-- https://go.dev/issue/56899, CL 463219 -->
On Windows calling
[`TempDir`](/pkg/os/#TempDir) now uses
GetTempPath2W when available, instead of GetTempPathW. The
new behavior is a security hardening measure that prevents
temporary files created by processes running as SYSTEM to
be accessed by non-SYSTEM processes.

<!-- CL 493036 -->
On Windows the os package now supports working with files whose
names, stored as UTF-16, can't be represented as valid UTF-8.

<!-- CL 463177 -->
On Windows [`Lstat`](/pkg/os/#Lstat) now resolves
symbolic links for paths ending with a path separator, consistent with its
behavior on POSIX platforms.

<!-- https://go.dev/issue/54451, CL 491175 -->
The implementation of the
[`io/fs.DirEntry`](/pkg/io/fs/#DirEntry)
interface returned by the
[`ReadDir`](/pkg/os/#ReadDir) function and
the [`File.ReadDir`](/pkg/os/#File.ReadDir)
method now implements a `String` method that calls
[`io/fs.FormatDirEntry`](/pkg/io/fs/#FormatDirEntry).

<!-- https://go.dev/issue/53761, CL 416775, CL 498015-->
The implementation of the
[`io/fs.FS`](/pkg/io/fs/#FS) interface returned by
the [`DirFS`](/pkg/os/#DirFS) function now implements
the [`io/fs.ReadFileFS`](/pkg/io/fs/#ReadFileFS) and
the [`io/fs.ReadDirFS`](/pkg/io/fs/#ReadDirFS)
interfaces.

<!-- os -->

#### [path/filepath](/pkg/path/filepath/)

The implementation of the
[`io/fs.DirEntry`](/pkg/io/fs/#DirEntry)
interface passed to the function argument of
[`WalkDir`](/pkg/path/filepath/#WalkDir)
now implements a `String` method that calls
[`io/fs.FormatDirEntry`](/pkg/io/fs/#FormatDirEntry).

<!-- path/filepath -->

<!-- CL 459455 reverted -->

#### [reflect](/pkg/reflect/)

<!-- CL 408826, CL 413474 -->
In Go 1.21, [`ValueOf`](/pkg/reflect/#ValueOf)
no longer forces its argument to be allocated on the heap, allowing
a `Value`'s content to be allocated on the stack. Most
operations on a `Value` also allow the underlying value
to be stack allocated.

<!-- https://go.dev/issue/55002 -->
The new [`Value`](/pkg/reflect/#Value)
method [`Value.Clear`](/pkg/reflect/#Value.Clear)
clears the contents of a map or zeros the contents of a slice.
This corresponds to the new `clear` built-in
[added to the language](#language).

<!-- https://go.dev/issue/56906, CL 452762 -->
The [`SliceHeader`](/pkg/reflect/#SliceHeader)
and [`StringHeader`](/pkg/reflect/#StringHeader)
types are now deprecated. In new code
prefer [`unsafe.Slice`](/pkg/unsafe/#Slice),
[`unsafe.SliceData`](/pkg/unsafe/#SliceData),
[`unsafe.String`](/pkg/unsafe/#String),
or [`unsafe.StringData`](/pkg/unsafe/#StringData).

<!-- reflect -->

#### [regexp](/pkg/regexp/)

<!-- https://go.dev/issue/46159, CL 479401 -->
[`Regexp`](/pkg/regexp#Regexp) now defines
[`MarshalText`](/pkg/regexp#Regexp.MarshalText)
and [`UnmarshalText`](/pkg/regexp#Regexp.UnmarshalText)
methods. These implement
[`encoding.TextMarshaler`](/pkg/encoding#TextMarshaler)
and
[`encoding.TextUnmarshaler`](/pkg/encoding#TextUnmarshaler)
and will be used by packages such as
[encoding/json](/pkg/encoding/json).

<!-- regexp -->

#### [runtime](/pkg/runtime/)

<!-- https://go.dev/issue/38651, CL 435337 -->
Textual stack traces produced by Go programs, such as those
produced when crashing, calling `runtime.Stack`, or
collecting a goroutine profile with `debug=2`, now
include the IDs of the goroutines that created each goroutine in
the stack trace.

<!-- https://go.dev/issue/57441, CL 474915 -->
Crashing Go applications can now opt-in to Windows Error Reporting (WER) by setting the environment variable
`GOTRACEBACK=wer` or calling [`debug.SetTraceback("wer")`](/pkg/runtime/debug/#SetTraceback)
before the crash. Other than enabling WER, the runtime will behave as with `GOTRACEBACK=crash`.
On non-Windows systems, `GOTRACEBACK=wer` is ignored.

<!-- CL 447778 -->
`GODEBUG=cgocheck=2`, a thorough checker of cgo pointer passing rules,
is no longer available as a [debug option](/pkg/runtime#hdr-Environment_Variables).
Instead, it is available as an experiment using `GOEXPERIMENT=cgocheck2`.
In particular this means that this mode has to be selected at build time instead of startup time.

`GODEBUG=cgocheck=1` is still available (and is still the default).

<!-- https://go.dev/issue/46787, CL 367296 -->
A new type `Pinner` has been added to the runtime
package. `Pinner`s may be used to "pin" Go memory
such that it may be used more freely by non-Go code. For instance,
passing Go values that reference pinned Go memory to C code is
now allowed. Previously, passing any such nested reference was
disallowed by the
[cgo pointer passing rules.](https://pkg.go.dev/cmd/cgo#hdr-Passing_pointers)
See [the docs](/pkg/runtime#Pinner) for more details.

<!-- CL 472195 no release note needed -->

<!-- runtime -->

#### [runtime/metrics](/pkg/runtime/metrics/)

<!-- https://go.dev/issue/56857, CL 497315 -->
A few previously-internal GC metrics, such as live heap size, are
now available.
`GOGC` and `GOMEMLIMIT` are also now
available as metrics.

<!-- runtime/metrics -->

#### [runtime/trace](/pkg/runtime/trace/)

<!-- https://go.dev/issue/16638 -->
Collecting traces on amd64 and arm64 now incurs a substantially
smaller CPU cost: up to a 10x improvement over the previous release.

<!-- CL 494495 -->
Traces now contain explicit stop-the-world events for every reason
the Go runtime might stop-the-world, not just garbage collection.

<!-- runtime/trace -->

#### [sync](/pkg/sync/)

<!-- https://go.dev/issue/56102, CL 451356 -->
The new [`OnceFunc`](/pkg/sync/#OnceFunc),
[`OnceValue`](/pkg/sync/#OnceValue), and
[`OnceValues`](/pkg/sync/#OnceValues)
functions capture a common use of [Once](/pkg/sync/#Once) to
lazily initialize a value on first use.

#### [syscall](/pkg/syscall/)

<!-- CL 480135 -->
On Windows the
[`Fchdir`](/pkg/syscall#Fchdir) function
now changes the current directory to its argument, rather than
always returning an error.

<!-- https://go.dev/issue/46259, CL 458335 -->
On FreeBSD
[`SysProcAttr`](/pkg/syscall#SysProcAttr)
has a new field `Jail` that may be used to put the
newly created process in a jailed environment.

<!-- CL 493036 -->
On Windows the syscall package now supports working with files whose
names, stored as UTF-16, can't be represented as valid UTF-8.
The [`UTF16ToString`](/pkg/syscall#UTF16ToString)
and [`UTF16FromString`](/pkg/syscall#UTF16FromString)
functions now convert between UTF-16 data and
[WTF-8](https://simonsapin.github.io/wtf-8/) strings.
This is backward compatible as WTF-8 is a superset of the UTF-8
format that was used in earlier releases.

<!-- CL 476578, CL 476875, CL 476916 -->
Several error values match the new
[`errors.ErrUnsupported`](/pkg/errors/#ErrUnsupported),
such that `errors.Is(err, errors.ErrUnsupported)`
returns true.

  - `ENOSYS`
  - `ENOTSUP`
  - `EOPNOTSUPP`
  - `EPLAN9` (Plan 9 only)
  - `ERROR_CALL_NOT_IMPLEMENTED` (Windows only)
  - `ERROR_NOT_SUPPORTED` (Windows only)
  - `EWINDOWS` (Windows only)


<!-- syscall -->

#### [testing](/pkg/testing/)

<!-- https://go.dev/issue/37708, CL 463837 -->
The new `-test.fullpath` option will print full path
names in test log messages, rather than just base names.

<!-- https://go.dev/issue/52600, CL 475496 -->
The new [`Testing`](/pkg/testing/#Testing) function reports whether the program is a test created by `go` `test`.

<!-- testing -->

#### [testing/fstest](/pkg/testing/fstest/)

<!-- https://go.dev/issue/54451, CL 491175 -->
Calling <code>[Open](/pkg/testing/fstest/MapFS.Open).[Stat](/pkg/io/fs/#File.Stat)</code>
will return a type that now implements a `String`
method that calls
[`io/fs.FormatFileInfo`](/pkg/io/fs/#FormatFileInfo).

<!-- testing/fstest -->

#### [unicode](/pkg/unicode/)

<!-- CL 456837 -->
The [`unicode`](/pkg/unicode/) package and
associated support throughout the system has been upgraded to
[Unicode 15.0.0](https://www.unicode.org/versions/Unicode15.0.0/).

<!-- unicode -->

## Ports {#ports}

### Darwin {#darwin}

<!-- https://go.dev/issue/57125 -->
As [announced](go1.20#darwin) in the Go 1.20 release notes,
Go 1.21 requires macOS 10.15 Catalina or later;
support for previous versions has been discontinued.

### Windows {#windows}

<!-- https://go.dev/issue/57003, https://go.dev/issue/57004 -->
As [announced](go1.20#windows) in the Go 1.20 release notes,
Go 1.21 requires at least Windows 10 or Windows Server 2016;
support for previous versions has been discontinued.

### ARM

<!-- CL 470695 -->

When building the Go distribution with `GOARCH=arm` when not running
on an ARM system (that is, when building a cross-compiler to ARM), the
default value for the `GOARM` environment variable is now always set
to `7`.
Previously the default depended on characteristics of the build
system.

When not building a cross-compiler, the default value is determined
by examining the build system.
That was true before and remains true in Go 1.21.
What has changed is the behavior when building a cross-compiler.

### WebAssembly {#wasm}

<!-- https://go.dev/issue/38248, https://go.dev/issue/59149, CL 489255 -->
The new `go:wasmimport` directive can now be used in Go programs
to import functions from the WebAssembly host.

<!-- https://go.dev/issue/56100 -->

The Go scheduler now interacts much more efficiently with the
JavaScript event loop, especially in applications that block
frequently on asynchronous events.

### WebAssembly System Interface {#wasip1}

<!-- https://go.dev/issue/58141 -->
Go 1.21 adds an experimental port to the [
WebAssembly System Interface (WASI)](https://wasi.dev/), Preview 1
(`GOOS=wasip1`, `GOARCH=wasm`).

As a result of the addition of the new `GOOS` value
"`wasip1`", Go files named `*_wasip1.go`
will now be [ignored
by Go tools](/pkg/go/build/#hdr-Build_Constraints) except when that `GOOS` value is being
used.
If you have existing filenames matching that pattern, you will
need to rename them.

### ppc64/ppc64le {#PPC64}

<!-- go.dev/issue/44549 -->
On Linux, `GOPPC64=power10` now generates PC-relative instructions, prefixed
instructions, and other new Power10 instructions. On AIX, `GOPPC64=power10`
generates Power10 instructions, but does not generate PC-relative instructions.

When building position-independent binaries for `GOPPC64=power10`
`GOOS=linux` `GOARCH=ppc64le`, users can expect reduced binary
sizes in most cases, in some cases 3.5%. Position-independent binaries are built for
ppc64le with the following `-buildmode` values:
`c-archive`, `c-shared`, `shared`, `pie`, `plugin`.

### loong64 {#loong64}

<!-- go.dev/issue/53301, CL 455075, CL 425474, CL 425476, CL 425478, CL 489576 -->
The `linux/loong64` port now supports `-buildmode=c-archive`,
`-buildmode=c-shared` and `-buildmode=pie`.

<!--
 proposals for x repos that don't need to be mentioned here but
     are picked up by the relnote tool.
 -->
<!-- https://go.dev/issue/54232 -->
<!-- https://go.dev/issue/57051 -->
<!-- https://go.dev/issue/57792 -->
<!-- https://go.dev/issue/57906 -->
<!-- https://go.dev/issue/58668 -->
<!-- https://go.dev/issue/59016 -->
<!-- https://go.dev/issue/59676 -->
<!-- https://go.dev/issue/60409 -->
<!-- https://go.dev/issue/61176 -->
<!-- changes to cmd/api that don't need release notes. -->
<!-- CL 469115, CL 469135, CL 499981 -->
<!-- proposals that don't need release notes. -->
<!-- https://go.dev/issue/10275 -->
<!-- https://go.dev/issue/59719 -->
