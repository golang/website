---
title: Go 1.28 Release Notes
template: false
---

<style>
  main ul li { margin: 0.5em 0; }
</style>

## DRAFT RELEASE NOTES — Introduction to Go 1.28 {#introduction}

**Go 1.28 is not yet released. These are work-in-progress release notes.
Go 1.28 is expected to be released in February 2027.**

## Changes to the language {#language}

<!-- go.dev/issue/12854 -->

Go 1.28 allows [omitting the type](/ref/spec#Omitting_composite_literal_types)
of a [composite literal](/ref/spec#Composite_literals) whenever its type can be
[inferred](/issue/12854) from [assignment context](/ref/spec#Assignability).
Previously, a composite literal could omit its type only when appearing within
an enclosing array, slice, or map literal.

```
type Point struct{ X, Y float64 }

var origin Point  = {}      // Point{}
var ptr    *Point = {1, 2}  // &Point{1, 2}
```

<!-- go.dev/issue/3939 -->

Starting with Go 1.28, the [conversion](/ref/spec#Conversions_to_and_from_a_string_type)
of an integer value `x` to a string type is valid only if `x`'s underlying type is `byte`
or `rune`, or if `x` is an untyped rune constant.
The ability to convert any integer value to a string type has been removed.
This new restriction lifts a long-existing `vet` check to the language level.
See [issue #3939](/issue/3939) for a discussion of this change.

```
string('a')     // "a"
string(97)      // invalid: 97 is not a rune constant
```

## Tools {#tools}

### Go command {#go-command}

<!-- go.dev/issue/81250 -->
The `GOFIPS140=certified` alias now selects v1.26.0 of the
[FIPS 140-3 Go Cryptographic Module](/doc/security/fips140), which obtained
[CMVP Certificate #5581](https://csrc.nist.gov/projects/cryptographic-module-validation-program/certificate/5581).

Support GOAMD64=local, like GOAMD64=vx,
where x is the value for the local micro-architecture.

### Vet {#vet}

The new [`scannererr`](https://pkg.go.dev/golang.org/x/tools/go/analysis/passes/scannererr)
analyzer checks for failure to handle scanner errors after a loop
around [`bufio.Scanner.Scan`](/pkg/bufio#Scanner.Scan), which may cause scanning or I/O errors to
go unreported. <!-- /issue/17747/ -->

The [`sqlrowserr`](https://pkg.go.dev/golang.org/x/tools/go/analysis/passes/sqlrowserr)
analyzer performs a similar check for loops around [`sql.Rows.Next`](/pkg/sql#Rows.Next),
so that iteration errors are correctly distinguished from a smaller result.

## Runtime {#runtime}

### Memory profile {#memory-profile}
<!-- go.dev/issue/79179 -->

A new profile type that breaks down the process's resident set size (RSS) is
now available.
The new profile type, named `memory`, is supported in the [`runtime/pprof`](/pkg/runtime/pprof) package
and as a [`net/http/pprof`](/pkg/net/http/pprof) endpoint, `/debug/pprof/memory`.

Go developers are often confused when the operating system reports much more
memory usage than the heap profile shows.
That's because the existing `inuse_space` heap profile only shows a sample of
objects that were live at the most recently completed garbage collection.
With the default garbage collection configuration, this is often less than
half of RSS.

The new memory profile is a superset of the `inuse_space` heap profile that
closes this gap by breaking down RSS into Go and non-Go memory.
Go memory is further split into the heap, stacks, and runtime metadata.
The heap includes the same live object samples as the heap profile, as well
as dead objects the most recently completed garbage collection found
unreachable.
On Linux, non-Go memory is further split into file-backed and shared memory,
and anonymous memory.

The runtime work that supports the memory profile can be turned off using
`GODEBUG=memoryprofile=0`. If you encounter a need to set this `GODEBUG`, please [file an issue](/issue/new).

Special thanks to Felix Geisendörfer and Nick Ripley at Datadog for
contributing this work.

## Compiler {#compiler}

<!-- go.dev/issue/83032 -->

[Generic type aliases](/ref/spec#Alias_declarations) were introduced in
[Go 1.24](/doc/go1.24#language), but due to an oversight, the compiler and
[`go/types`](/pkg/go/types) checked for language version Go 1.23 instead of Go 1.24.
Starting with Go 1.28, declaring a generic type alias requires language
version Go 1.24 or later. Packages targeting Go 1.23 (for instance, via
`go 1.23` in `go.mod`) that declare generic type aliases will now be
rejected unless updated to at least `go 1.24`.
See [issue #83032](/issue/83032).

## Linker {#linker}

<!-- go.dev/issue/62577 -->
<!-- go.dev/issue/82729 -->
On macOS and iOS, the linker now supports emitting DWARF debug information
into a separate file in `<out>.dSYM/Contents/Resources/DWARF/<out>`, instead
of combining it into the output binary.
By default, split DWARF is enabled for `-buildmode=c-shared` and
`-buildmode=plugin`, and disabled for executables.
(On macOS, `-buildmode=c-shared` previously defaulted to `-w` to omit DWARF;
it now emits split DWARF by default.)
When split DWARF is enabled, `go build` and `go test -c` output the `.dSYM`
directory alongside the built binary, whereas `go install` does not install
the `.dSYM` directory.
`go clean` removes it along with the binary.

The `.dSYM` directory follows the platform convention used by the C toolchain
and the `dsymutil` tool.
LLDB understands this convention and will load the debug information automatically
from the split DWARF file.
(Note: LLDB does not understand compressed DWARF, so the binary needs to be
built with the `-ldflags=-compressdwarf=0` build option.)
<!-- TODO: mention Delve once Delve supports it -->
The [`debug/macho.File.DWARF`](/pkg/debug/macho#File.DWARF) method is also
updated to load the debug information from the split DWARF file if available.

A linker command line option `-splitdwarf` can be used to enable or disable
split DWARF.
This option may be removed in a future Go release.

## Bootstrap {#boostrap}

<!-- go.dev/issue/69315 -->
As mentioned in the [Go 1.26 release notes](/doc/go1.26#bootstrap), Go 1.28 now requires
Go 1.26.0 or later for bootstrap.
We expect that Go 1.30 will require a minor release of Go 1.28.0 or later for bootstrap.

## Standard library {#library}

### SIMD {#simd}

The new [`simd`](/pkg/simd) and [`simd/archsimd`](/pkg/simd/archsimd) packages provide access to SIMD
(single instruction, multiple data) operations. The [`simd`](/pkg/simd) package
provides portable, vector-size-agnostic types and operations, and the
[`simd/archsimd`](/pkg/simd/archsimd) package provides architecture-specific vector types and
operations on amd64, arm64, and WebAssembly.

These packages, and the compiler support for them, can be disabled by
setting `GOEXPERIMENT=nosimdstable` at build time. This setting is
expected to be removed in a future release.

Additional, experimental API, including arm64 SVE support, is available
by setting `GOEXPERIMENT=simd` at build time. That API is not subject
to the Go 1 compatibility promise.

TODO: Expand.
<!-- go.dev/issue/83009 -->

### Minor changes to the library {#minor_library_changes}

#### [`crypto/ecdsa`](/pkg/crypto/ecdsa/)

<!-- 6-stdlib/99-minor/crypto/ecdsa/hashlen.md -->
[`SignASN1`](/pkg/crypto/ecdsa#SignASN1), [`Sign`](/pkg/crypto/ecdsa#Sign), and [`PrivateKey.Sign`](/pkg/crypto/ecdsa#PrivateKey.Sign) now return an error, and [`VerifyASN1`](/pkg/crypto/ecdsa#VerifyASN1)
and [`Verify`](/pkg/crypto/ecdsa#Verify) now return false, if the hash is shorter than 16 bytes or longer
than 64 bytes. Such inputs are not the output of a supported hash function, and
suggest a message is being passed instead of its hash, which can be severely
insecure.

#### [`crypto/tls`](/pkg/crypto/tls/)

<!-- 6-stdlib/99-minor/crypto/tls/81199.md -->
The `signature_algorithms_cert` extension is no longer sent, because some
middleboxes reject it when it includes ML-DSA. As a result, support for SHA-1
signatures in certificates is no longer advertised by default. The extension
can be re-enabled, without ML-DSA, with the `tlssigalgscert=1` GODEBUG option.

<!-- 6-stdlib/99-minor/crypto/tls/82402.md -->
Clients now check server certificates against a local
[UPKI revocation cache](https://c2sp.org/upki-revocation) and reject
certificates the cache reports as revoked. For programs targeting Go 1.28
or later, checking is on by default whenever a cache is found, so a
program's behavior can change when its host
starts running a cache updater such as `upki fetch`, which some operating
system distributions ship. The cache is looked for in the user cache
directory and then the system cache directory: `~/.cache/upki/revocation`
and `/var/cache/upki/revocation` on Linux and other Unix systems,
`~/Library/Caches/dev.rustls.upki/revocation` and
`/Library/Caches/dev.rustls.upki/revocation` on macOS, and
`%LocalAppData%\rustls\upki\cache\revocation` and
`%ProgramData%\rustls\upki\cache\revocation` on Windows. Certificates the
cache has no information about, and certificates without embedded Signed
Certificate Timestamps, are accepted. Revocation lookups require a cache
and a certificate chain verified by the TLS client. This includes the ECH
public-name certificate when ECH is rejected and
[`Config.EncryptedClientHelloRejectionVerify`](/pkg/crypto/tls#Config.EncryptedClientHelloRejectionVerify) is nil, even when
[`Config.InsecureSkipVerify`](/pkg/crypto/tls#Config.InsecureSkipVerify) is set.

The new [`Config.UPKIRevocation`](/pkg/crypto/tls#Config.UPKIRevocation) field, of type [`UPKIRevocationConfig`](/pkg/crypto/tls#UPKIRevocationConfig),
selects a cache directory explicitly, adjusts the policy for inconclusive
checks, or opts out with its `Disabled` field. Automatic discovery skips
missing or inaccessible cache paths. Applications can set
[`UPKIRevocationConfig.MissingData`](/pkg/crypto/tls#UPKIRevocationConfig.MissingData) to [`UPKIRevocationPolicyReject`](/pkg/crypto/tls#UPKIRevocationPolicyReject) to require
an available index for lookups, or to [`UPKIRevocationPolicyAllow`](/pkg/crypto/tls#UPKIRevocationPolicyAllow) to tolerate
an unavailable index in an explicitly configured cache. A malformed embedded
SCT extension causes verification to fail when a cache is explicitly
configured or found during discovery. A revoked certificate is
reported as a [`CertificateVerificationError`](/pkg/crypto/tls#CertificateVerificationError) wrapping [`ErrUPKIRevoked`](/pkg/crypto/tls#ErrUPKIRevoked).
The [GODEBUG setting](/doc/godebug) `tlsupki=0` disables checking with
the platform-default cache. Programs targeting Go 1.27 or earlier use this
setting by default and can opt in with `tlsupki=1`.

#### [`crypto/x509`](/pkg/crypto/x509/)

<!-- 6-stdlib/99-minor/crypto/x509/79860.md -->
[`ParseCertificate`](/pkg/crypto/x509#ParseCertificate), [`ParseCertificateRequest`](/pkg/crypto/x509#ParseCertificateRequest), [`ParseRevocationList`](/pkg/crypto/x509#ParseRevocationList) and
[`ParsePKIXPublicKey`](/pkg/crypto/x509#ParsePKIXPublicKey) now reject trailing data after the parameters of an
AlgorithmIdentifier, after the RSAPublicKey structure of an RSA public key
and after its public exponent, and after the named curve of an ECDSA public key.

<!-- 6-stdlib/99-minor/crypto/x509/81798.md -->
[`ParseCertificate`](/pkg/crypto/x509#ParseCertificate) and [`ParseCertificateRequest`](/pkg/crypto/x509#ParseCertificateRequest) now reject trailing data
after the value of the key usage, basic constraints, subject alternative
name, extended key usage and certificate policies extensions.

<!-- 6-stdlib/99-minor/crypto/x509/params.md -->
Signature algorithm parameters other than NULL or absent for RSA PKCS #1 v1.5,
and other than absent for ECDSA, are no longer accepted in certificates,
certificate requests, and revocation lists.

#### [`debug/macho`](/pkg/debug/macho/)

<!-- 6-stdlib/99-minor/debug/macho/62577.md -->
If a Mach-O file does not contain DWARF sections, [`File.DWARF`](/pkg/debug/macho#File.DWARF) now loads
the debug information from the associated split DWARF file in
`<file>.dSYM/Contents/Resources/DWARF/<file>` if available.

#### [`debug/pe`](/pkg/debug/pe/)

<!-- 6-stdlib/99-minor/debug/pe/81594.md -->
The [`ImportDirectory`](/pkg/debug/pe#ImportDirectory) type is now marked as deprecated.
It was exported by accident, is an internal detail of the package,
and should not be used directly.

#### [`embed`](/pkg/embed/)

<!-- 6-stdlib/99-minor/embed/80822.md -->
<!-- go.dev/issue/80822 -->

The `ReadAt` methods of files returned by [`embed.FS`](/pkg/embed#FS) now return
[`io.EOF`](/pkg/io#EOF) when called with an offset at or beyond the end of the
file, as required by [`io.ReaderAt`](/pkg/io#ReaderAt).

#### [`encoding/base32`](/pkg/encoding/base32/)

<!-- 6-stdlib/99-minor/encoding/base32/20235.md -->
[`Encoding.EncodedLen`](/pkg/encoding/base32#Encoding.EncodedLen) now panics if the encoded length overflows int,
rather than returning a negative number.

#### [`encoding/base64`](/pkg/encoding/base64/)

<!-- 6-stdlib/99-minor/encoding/base64/20235.md -->
[`Encoding.EncodedLen`](/pkg/encoding/base64#Encoding.EncodedLen) now panics if the encoded length overflows int,
rather than returning a negative number.

#### [`flag`](/pkg/flag/)

<!-- 6-stdlib/99-minor/flag/65675.md -->
The new [`FlagSet.All`](/pkg/flag#FlagSet.All) method returns an iterator over all the flags in the FlagSet;
the function [`All`](/pkg/flag#All) does the same for the global [`CommandLine`](/pkg/flag#CommandLine) flag set.
The new [`Flag.IsSet`](/pkg/flag#Flag.IsSet) field indicates whether the flag has been set.

#### [`go/analysis`](/pkg/go/analysis/)

<!-- 6-stdlib/99-minor/go/analysis/82465.md -->
<!-- go.dev/issue/82465 -->

The API of the Go analysis framework, first published in 2019 at
[`golang.org/x/tools/go/analysis`](/pkg/golang.org/x/tools/go/analysis) as part of the Go 1.12 reorganization
of `cmd/vet`, is now part of the standard library.
The new [`go/analysis`](/pkg/go/analysis) package declares the types, such as [`Analyzer`](/pkg/go/analysis#Analyzer),
[`Pass`](/pkg/go/analysis#Pass), and [`Diagnostic`](/pkg/go/analysis#Diagnostic), that define the interface between analysis
drivers (such as `cmd/vet`, `gopls`, and `staticcheck`), and specific
analyzers (sometimes called checks or linters).
Aside from the package name, the actual API is identical.

<!-- I will expand this into a whole section of the
     release notes when the #82464 work lands. --adonovan
-->

#### [`go/build`](/pkg/go/build/)

<!-- 6-stdlib/99-minor/go/build/75129.md -->
The [`build.ArchChar`](/pkg/build#ArchChar) function is deprecated, as it no longer varies by architecture.

#### [`go/doc`](/pkg/go/doc/)

<!-- 6-stdlib/99-minor/go/doc/80118.md -->
The new [`Deprecated`](/pkg/go/doc#Deprecated) function reports the deprecation notice in a doc comment, if any.

#### [`go/parser`](/pkg/go/parser/)

<!-- 6-stdlib/99-minor/go/parser/79802.md -->
The [`ResolveFile`](/pkg/go/parser#ResolveFile) function performs legacy syntactic identifier
resolution on an [`ast.File`](/pkg/ast#File) syntax tree. This resolution is
deprecated; new applications should instead use [`go/types`](/pkg/go/types).
This function is provided to ease the transition away from the old
resolution in applications that mostly parse files in
[`SkipObjectResolution`](/pkg/go/parser#SkipObjectResolution) mode but later encounter some cases where the
old resolution is needed.

#### [`go/types`](/pkg/go/types/)

<!-- 6-stdlib/99-minor/go/types/60994.md -->
The new [`CommonUnderlying`](/pkg/go/types#CommonUnderlying) function returns the common underlying type
(if any) of the types in a type set. If its argument is a type
parameter, the result is the single underlying type (if it exists) of
all types in its constraint's type set; otherwise, the result is the
type's underlying type. (This is equivalent to the "core type" concept
formerly used in the Go specification.)

<!-- 6-stdlib/99-minor/go/types/72059.md -->
The new [`StrictlyComparable`](/pkg/go/types#StrictlyComparable) function reports whether values of a given type are [strictly comparable](/ref/spec#Comparison_operators).

<!-- 6-stdlib/99-minor/go/types/79603.md -->
Added [`NewTypeList`](/pkg/go/types#NewTypeList) and [`NewTypeParamList`](/pkg/go/types#NewTypeParamList) constructors to
create [`TypeList`](/pkg/go/types#TypeList) and [`TypeParamList`](/pkg/go/types#TypeParamList) instances
from slices of their respective element types.

<!-- 6-stdlib/99-minor/go/types/80449.md -->
The new [`Scope.Objects`](/pkg/go/types#Scope.Objects) method returns the sequence of Objects in the Scope.

#### [`log/slog`](/pkg/log/slog/)

<!-- 6-stdlib/99-minor/log/slog/81838.md -->
<!-- go.dev/issue/81838 -->
In Go 1.27.1, direct calls to [`Handler.WithGroup`](/pkg/log/slog#Handler.WithGroup) that passed `""`, including
those inside [`MultiHandler`](/pkg/log/slog#MultiHandler), would erroneously insert a group with name `""`
instead of doing nothing. This bug is now fixed.

#### [`maps`](/pkg/maps/)

<!-- 6-stdlib/99-minor/maps/81698.md -->
The new [`Grow`](/pkg/maps#Grow) function increases the capacity of an existing map
so that a known number of additional entries can be added without
repeated incremental growth, for example before a [`Copy`](/pkg/maps#Copy) into a map
that cannot be replaced.

#### [`math`](/pkg/math/)

<!-- 6-stdlib/99-minor/math/82653.md -->
The new [`FMA32`](/pkg/math#FMA32) function computes x\*y + z for float32 values with only one
rounding. It uses a hardware fused multiply-add instruction where available.

#### [`net`](/pkg/net/)

<!-- 6-stdlib/99-minor/net/29678.md -->
The [`IPMask`](/pkg/net#IPMask), [`IPNet`](/pkg/net#IPNet), and [`HardwareAddr`](/pkg/net#HardwareAddr) types each now implement
a [`encoding.TextUnmarshaler`](/pkg/encoding#TextUnmarshaler) method. This method will recognize
the current marshaling of those types, as well as a planned future
marshaling format. The future marshaling format will be readable
in JSON format, unlike the current format which is a base64 encoding.
The future marshaling format will become the default in a later Go release.
This is being done as a two-step process so that JSON generated
by this release of Go can still be read by earlier releases of Go.

<!-- 6-stdlib/99-minor/net/63118.md -->
<!-- https://go.dev/issue/63118 -->
The unused [`DNSConfigError`](/pkg/net#DNSConfigError) type is deprecated.

#### [`net/http`](/pkg/net/http/)

<!-- 6-stdlib/99-minor/net/http/53208.md -->
<!-- go.dev/issue/53208 -->
The new [`HTTP2Config.EnableConnectProtocol`](/pkg/net/http#HTTP2Config.EnableConnectProtocol) field enables support for the
Extended CONNECT protocol defined in [RFC 8441](https://www.rfc-editor.org/rfc/rfc8441).
When enabled, HTTP/2 servers advertise support for Extended CONNECT requests.

<!-- 6-stdlib/99-minor/net/http/77441.md -->
The new method [`Server.HandleRequest`](/pkg/net/http#Server.HandleRequest) handles a request as if it was received by the server.

<!-- 6-stdlib/99-minor/net/http/79040.md -->
The long-deprecated [`Transport.CancelRequest`](/pkg/net/http#Transport.CancelRequest) now does nothing.
Use [`NewRequestWithContext`](/pkg/net/http#NewRequestWithContext) (available since Go 1.13)
to create a cancellable request.

<!-- 6-stdlib/99-minor/net/http/79656.md -->
<!-- go.dev/issue/79656 -->

[`ProxyFromEnvironment`](/pkg/net/http#ProxyFromEnvironment) now prioritizes lowercase proxy environment variables
(`http_proxy`, `https_proxy`, and `no_proxy`) over the uppercase versions if
both are set to different values. This change also affects default proxy lookups
in [`DefaultTransport`](/pkg/net/http#DefaultTransport) and [`DefaultClient`](/pkg/net/http#DefaultClient) by extension.

<!-- 6-stdlib/99-minor/net/http/80058.md -->
<!-- go.dev/issue/80058 -->
The new [`MethodQuery`](/pkg/net/http#MethodQuery) constant represents the HTTP QUERY method defined in
[RFC 10008](https://www.rfc-editor.org/rfc/rfc10008).

<!-- 6-stdlib/99-minor/net/http/80735.md -->
<!-- go.dev/issue/80735 -->
The HTTP/2 server now lets a connection's serve goroutine exit while
the connection has nothing to do, restarting it on demand. This
reduces memory usage for servers with many idle connections or many
handlers idling mid-response (such as long polls).

<!-- 6-stdlib/99-minor/net/http/82272.md -->
<!-- go.dev/issue/82272 -->

[`Transport`](/pkg/net/http#Transport) uses an entirely rewritten connection pool,
which manages HTTP/1 and HTTP/2 connections in a single pool.

Transport connection limits now consistently apply to all connections:
[`Transport.MaxIdleConns`](/pkg/net/http#Transport.MaxIdleConns), [`Transport.MaxIdleConnsPerHost`](/pkg/net/http#Transport.MaxIdleConnsPerHost), and [`Transport.MaxConnsPerHost`](/pkg/net/http#Transport.MaxConnsPerHost)
all apply to the combined total of HTTP/1 and HTTP/2 connections.
Previously, these limits were applied inconsistently or not at all to HTTP/2 connections.

Concurrent requests to the same host now coalesce onto a single dial
when HTTP/2 is enabled. If HTTP/2 is negotiated, the requests will share
the new connection.

#### [`net/url`](/pkg/net/url/)

<!-- 6-stdlib/99-minor/net/url/79946.md -->
<!-- go.dev/issue/79946 -->
The new [`MustParse`](/pkg/net/url#MustParse) function returns the parsed [`URL`](/pkg/net/url#URL)
or panics on error.
It can simplify some common uses of [`Parse`](/pkg/net/url#Parse).
For instance, initializing variables with valid URL string constants.

#### [`regexp`](/pkg/regexp/)

<!-- 6-stdlib/99-minor/regexp/61902.md -->
<!-- go.dev/issue/61902 -->
New [`Regexp.All`](/pkg/regexp#Regexp.All), [`Regexp.AllIndex`](/pkg/regexp#Regexp.AllIndex), [`Regexp.AllString`](/pkg/regexp#Regexp.AllString), [`Regexp.AllStringIndex`](/pkg/regexp#Regexp.AllStringIndex),
[`Regexp.AllStringSubmatch`](/pkg/regexp#Regexp.AllStringSubmatch), [`Regexp.AllStringSubmatchIndex`](/pkg/regexp#Regexp.AllStringSubmatchIndex), [`Regexp.AllSubmatch`](/pkg/regexp#Regexp.AllSubmatch), and
[`Regexp.AllSubmatchIndex`](/pkg/regexp#Regexp.AllSubmatchIndex) iterator methods on [`Regexp`](/pkg/regexp#Regexp) allow incremental matching of the input.

The new [`Regexp.SplitSeq`](/pkg/regexp#Regexp.SplitSeq) method iterates over the results of splitting the input with the
regular expression.

#### [`simd`](/pkg/simd/)

<!-- 6-stdlib/99-minor/simd/83009.md -->
The new [`simd`](/pkg/simd) package provides portable, vector-size-agnostic SIMD
types and operations. On architectures supported by [`simd/archsimd`](/pkg/simd/archsimd),
the simd package uses multi-compilation and runtime feature detection to
transparently utilize the largest vectors and latest CPU features
supported by the hardware a program is running on, regardless of where
it was compiled. On unsupported architectures, or if hardware vector
acceleration isn't available, it falls back to emulation.
<!-- go.dev/issue/83009 -->
<!-- TODO: expand into a release note section for simd and simd/archsimd. -->

#### [`simd/archsimd`](/pkg/simd/archsimd/)

<!-- 6-stdlib/99-minor/simd/archsimd/83009.md -->
The new [`simd/archsimd`](/pkg/simd/archsimd) package provides access to architecture-specific
SIMD (single instruction, multiple data) vector types and operations on
amd64, arm64, and wasm.
<!-- go.dev/issue/83009 -->

#### [`strconv`](/pkg/strconv/)

<!-- 6-stdlib/99-minor/strconv/82205.md -->
When [`Unquote`](/pkg/strconv#Unquote) is applied to a string containing a single-quoted rune
literal, it now correctly returns the rune's UTF-8 encoding, even when
the literal uses a hexadecimal or octal escape, as in
``Unquote(`'\xff'`)``. (Previously it would return a string whose
sole byte had the same numeric value as the rune, which is incorrect
for non-ASCII.)

#### [`sync`](/pkg/sync/)

<!-- 6-stdlib/99-minor/sync/80645.md -->
<!-- go.dev/issue/80645 -->
The new [`Counter`](/pkg/sync#Counter) type is a sharded counter of int64, uint64, or float64
values that is optimized for frequent concurrent calls to [`Counter.Add`](/pkg/sync#Counter.Add) and
infrequent calls to [`Counter.Value`](/pkg/sync#Counter.Value). Updates are spread across per-processor
shards, so Add remains fast even when many goroutines update the same counter,
while Value sums the shards and takes time proportional to GOMAXPROCS.
The zero value of a Counter is ready to use.

#### [`syscall`](/pkg/syscall/)

<!-- 6-stdlib/99-minor/syscall/68595.md -->
On Linux, [`SysProcAttr`](/pkg/syscall#SysProcAttr) now supports enforcement of Landlock policies
and the `PR_SET_NO_NEW_PRIVS` flag for newly spawned subprocesses.

<!-- 6-stdlib/99-minor/syscall/80415.md -->
<!-- go.dev/issue/80415 -->
On Windows, [`SysProcAttr`](/pkg/syscall#SysProcAttr) has several new fields to support additional process
creation features: `Jobs` (`[]Handle`) specifies a list of job object handles
to associate with the new process, `Console` (`Handle`) specifies a console
handle to connect to the process, `Desktop` (`string`) specifies the name of
the desktop on which the process will start, and `AppContainer`
(`*AppContainer`) specifies a pointer to a new [`AppContainer`](/pkg/syscall#AppContainer) struct to create
the process within an AppContainer using the provided SID and capabilities.

#### [`testing/fstest`](/pkg/testing/fstest/)

<!-- 6-stdlib/99-minor/testing/fstest/80822.md -->
<!-- go.dev/issue/80822 -->

The `ReadAt` methods of files returned by [`testing/fstest.MapFile.Open`](/pkg/testing/fstest#MapFile.Open)
now return [`io.EOF`](/pkg/io#EOF) when called with an offset at or beyond the end
of the file, as required by [`io.ReaderAt`](/pkg/io#ReaderAt).

#### [`testing/synctest`](/pkg/testing/synctest/)

<!-- 6-stdlib/99-minor/testing/synctest/77320.md -->
The new [`Subtest`](/pkg/testing/synctest#Subtest) convenience function combines [testing.(*T).Run] and [`testing/synctest.Test`](/pkg/testing/synctest#Test).

#### [`time`](/pkg/time/)

<!-- 6-stdlib/99-minor/time/26032.md -->
[`Parse`](/pkg/time#Parse) and [`ParseInLocation`](/pkg/time#ParseInLocation) now accept a numeric time zone abbreviation
with minutes, such as "+0530", where the layout has a zone abbreviation
(MST). A numeric abbreviation that the location does not define now gets
the offset it denotes rather than a zero offset. Numeric abbreviations
with three or more than four digits are no longer accepted.

<!-- 6-stdlib/99-minor/time/67066.md -->
<!-- go.dev/issue/67066 -->

On Windows, the monotonic clock used by [`Now`](/pkg/time#Now) now has a resolution of
100 nanoseconds. Previously, its resolution was on the order of milliseconds.
This improves the precision of elapsed-time measurements, such as those
made with [`Since`](/pkg/time#Since).
