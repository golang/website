---
title: Go 1.27 Release Notes
template: false
---

<style>
  main ul li { margin: 0.5em 0; }
</style>

## DRAFT RELEASE NOTES — Introduction to Go 1.27 {#introduction}

**Go 1.27 is not yet released. These are work-in-progress release notes.
Go 1.27 is expected to be released in August 2026.**

## Changes to the language {#language}

<!-- go.dev/issue/9859 -->

A key in a [struct literal](/ref/spec#Composite_literals) may now be any
valid [field selector](/issue/9859) for the struct type, not just a
(top-level) field name of the struct.

<!-- go.dev/issue/77273 -->

Go 1.27 now supports [generic methods](/issue/77273):
a [method declaration](/ref/spec#Method_declarations) may declare its own
[type parameters](/ref/spec#Type_parameter_declarations).
This widely anticipated change allows adding generic functions within
the namespace of a particular data type where before one had to declare
such functions with a scope of the entire package.
Note that methods of [interfaces](/ref/spec#Interface_types) may not declare
type parameters nor can interface methods be implemented by generic methods.

<!-- go.dev/issue/77245 -->

Function type inference has been [generalized](/issue/77245) to apply in all
contexts where a generic function is [assigned](/ref/spec#Assignability) to a
variable of (or converted to) a matching function type.

## Tools {#tools}

<!-- go.dev/issue/77177 -->

Response file (`@file`) parsing is now supported for the `compile`, `link`, `asm`, `cgo`, `cover`, and `pack` tools.
The response file contains whitespace-separated arguments with support for single-quoted and double-quoted strings, escape sequences, and backslash-newline line continuation.
The format is compatible with GCC's response file implementation to ensure interoperability with existing build systems.

### Go command {#go-command}

`go test` now invokes the `stdversion` vet check by default.
This reports the use of standard library symbols that are too new
for the Go version in force in the referring file,
as determined by `go` directive in `go.mod` and build tags on the file.

<!-- go.dev/issue/78090 -->

The `go` command no longer has support for the `bzr` version control system.
It will no longer be able to directly fetch modules hosted on `bzr` servers.

The `go fix` command contains several new modernizers (`atomictypes`, `embedlit`, `slicesbackward`, and `unsafefuncs`).
The existing `fmtappendf` analyzer was removed due to stylistic concerns. <!-- #77581 -->
The existing `waitgroup` analyzer was renamed to `waitgroupgo` to avoid ambiguity.

<!-- go.dev/issue/63696 -->

The `go doc` command now supports `package@version` syntax, such as
`go doc example.com/pkg@v1.2.3`.

<!-- go.dev/issue/26715 -->

The `go doc` command now accepts the `-ex` command-line option to
list executable examples of the given package or symbol.
When an example name is passed on the command line (such as
`go doc bytes.ExampleBuffer`), `go doc` now prints the example source
code along with comments.

<!-- go.dev/issue/79422 -->

Starting with the Go 1.27 tool chain, the `go` command now recognizes a `GODEBUG` setting
for which support was removed (such as `asynctimerchan`, see below) if it appears in `go.mod`
files (`go debug` entries) and `.go` source files (`//go:debug` comments).
It accepts these settings if they are set to the final default value established before
the setting was removed.
If they are set to an old value, the `go` command will fail.
This change is in the spirit of the [Go 1 compatibility guarantee](/doc/go1compat)
and allows existing programs that set supported `GODEBUG` settings to continue to
build and run without changes even when the respective setting support has been removed.

<!-- go.dev/issue/56471 -->

For modules specifying `go 1.27` or later in their `go.mod` file, `go mod tidy`
now automatically merges duplicate require blocks. This ensures the file
maintains a clean, standard structure containing at most two require blocks:
one for direct dependencies and one for indirect dependencies.

Existing comment blocks attached to dependencies are preserved during this
consolidation. If a comment block is associated with a mixed set of directives
(containing both direct and indirect dependencies), the comment block is merged
and attached to the new direct dependency block.

Previously, if a `go.mod` file accumulated multiple disjoint require blocks
(often due to manual edits, unresolved Git merge conflicts, or legacy upgrades)
`go mod tidy` would leave the extra blocks intact or inadvertently create new
ones. The tool now strictly enforces the two-block layout, consolidating
disparate requirements into their respective blocks and cleaning up the
structure of the module file automatically.

### Trace

<!-- go.dev/issue/78921 -->

`go tool trace`'s `-http` argument now restricts the listen address to localhost when passed only a port (e.g., `-http=:6060`).
This change makes `go tool trace` consistent with the behavior of `go tool pprof`'s `-http` flag.
To listen on all addresses, explicitly include the specified address (e.g., `-http=0.0.0.0:6060`).

## Runtime {#runtime}

<!-- CL 742580 -->

Tracebacks for modules with `go` directives configuring Go 1.27 or later will now
include [runtime/pprof](https://pkg.go.dev/runtime/pprof) goroutine labels in
the header line. This behavior can be disabled with `GODEBUG=tracebacklabels=0`
(added in [Go 1.26](/doc/godebug#go-126)). This opt-out is expected to be
kept indefinitely in case goroutine labels acquire sensitive information that
shouldn't be made available in tracebacks.

<!-- CL 781580 -->

The `asynctimerchan` `GODEBUG` setting (added in [Go 1.23](/doc/godebug#go-123))
has been removed permanently. Channels created by package [time](https://pkg.go.dev/time)
are now always unbuffered (synchronous), irrespective of `GODEBUG` settings.

### Faster memory allocation

<!-- go.dev.issue/79286 -->

The compiler will now generate calls to size-specialized memory allocation
routines, reducing the cost of some small (<80 byte) memory allocations by
up to 30%.
Improvements vary depending on the workload, but the overall improvement is
expected to be ~1% in real allocation-heavy programs.
This causes the binary size to increase by about 60 KB (independent of the
workload).
Please [file an issue](/issue/new) if you notice any regressions.
You may set `GOEXPERIMENT=nosizespecializedmalloc` at build time to disable
it.
This opt-out setting is expected to be removed in Go 1.28.

## Compiler {#compiler}

The compiler now resolves a relative filename in a `//line` or `/*line*/`
directive against the directory of the file containing the directive,
matching [`go/scanner`](/pkg/go/scanner). Absolute filenames are unaffected.
See [#70478](/issue/70478).

## Linker {#linker}

<!-- CL 751260, go.dev/issue/58722 -->

When targeting macOS, the linker now accepts `-macos` and `-macsdk`
command-line options, which specify the OS and SDK versions in the
`LC_BUILD_VERSION` load command.
By default, it selects the oldest supported macOS version (currently
[13.0.0](#darwin)) and a recent SDK version (currently 26.2.0).

## Standard library {#library}

### New encoding/json/v2 and encoding/json/jsontext packages

<!-- go.dev/issue/71497 -->

Two new packages are now available:

  - The [`encoding/json/v2`](/pkg/encoding/json/v2) package is a major
    revision of [`encoding/json`](/pkg/encoding/json). It provides
    [`Marshal`](/pkg/encoding/json/v2#Marshal),
    [`MarshalWrite`](/pkg/encoding/json/v2#MarshalWrite),
    [`MarshalEncode`](/pkg/encoding/json/v2#MarshalEncode),
    [`Unmarshal`](/pkg/encoding/json/v2#Unmarshal),
    [`UnmarshalRead`](/pkg/encoding/json/v2#UnmarshalRead), and
    [`UnmarshalDecode`](/pkg/encoding/json/v2#UnmarshalDecode),
    all of which accept variadic [`Options`](/pkg/encoding/json/v2#Options)
    arguments to configure marshaling and unmarshaling behavior.

  - The [`encoding/json/jsontext`](/pkg/encoding/json/jsontext) package
    provides lower-level syntactic processing of JSON.
    The [`Encoder`](/pkg/encoding/json/jsontext#Encoder) and
    [`Decoder`](/pkg/encoding/json/jsontext#Decoder) types operate on
    JSON as a sequence of
    [`Token`](/pkg/encoding/json/jsontext#Token) and
    [`Value`](/pkg/encoding/json/jsontext#Value),
    maintaining a state machine to ensure the produced or consumed
    sequence is valid JSON text.

The v2 package chooses stricter, more interoperable defaults than v1:
it rejects invalid UTF-8 in JSON strings and rejects duplicate names within
a JSON object. See the v1 [`encoding/json`](/pkg/encoding/json#hdr-Migrating_to_v2) package
documentation for the complete set of behavioral differences and
the options available to adjust them.

The [`encoding/json`](/pkg/encoding/json) package is now backed by the
v2 implementation. Marshaling and unmarshaling behavior is preserved, but
the exact text of error messages may differ. The package also gains a number of
new [`Options`](/pkg/encoding/json#Options) that can configure v2 to operate
with v1 semantics to avoid requiring a full migration to the new API.
The v1 API will continue to be supported and users are not required to migrate.

Marshal performance is broadly at parity with the previous implementation,
while unmarshal performance is significantly faster.

Users who encounter compatibility problems with the new implementation
may disable it by setting `GOEXPERIMENT=nojsonv2` at build time,
restoring the original v1 implementation.
This opt-out is expected to be removed in a future release.

See the [proposal issue](/issue/71497) for background and additional detail.
If you need to disable the new implementation, [please file an issue](/issue/new).

### New uuid package

<!-- https://go.dev/issue/62026 --->

The new [`uuid`](/pkg/uuid) package generates and parses UUIDs.

### New crypto/mldsa package

<!-- https://go.dev/issue/77626, https://go.dev/issue/78888 --->

The new [`crypto/mldsa`](/pkg/crypto/mldsa) package implements the post-quantum ML-DSA signature
scheme specified in FIPS 204.

[`crypto/x509`](/pkg/crypto/x509) now supports ML-DSA private keys, public keys, and signatures.

[`crypto/tls`](/pkg/crypto/tls) now supports ML-DSA signatures in TLS 1.3, with the new
[MLDSA44], [MLDSA65], and [MLDSA87] [SignatureScheme] values.

### Minor changes to the library {#minor_library_changes}

#### [`bytes`](/pkg/bytes/)

<!-- 6-stdlib/99-minor/bytes/71151.md -->

The new [`CutLast`](/pkg/bytes#CutLast) function slices a []byte
around the last occurrence of a separator.
It can replace and simplify some common uses of LastIndex.

#### [`crypto`](/pkg/crypto/)

<!-- 6-stdlib/99-minor/crypto/77626.md -->

The new [`MLDSAMu`](/pkg/crypto#MLDSAMu) [`Hash`](/pkg/crypto#Hash) value is meant to be used as a signaling mechanism for
External μ ML-DSA signing.

#### [`crypto/ecdsa`](/pkg/crypto/ecdsa/)

<!-- 6-stdlib/99-minor/crypto/ecdsa/hashlen.md -->

[`PrivateKey.Sign`](/pkg/crypto/ecdsa#PrivateKey.Sign) now checks that the length of the hash is correct, if opts is
not nil.

#### [`crypto/mldsa`](/pkg/crypto/mldsa/)

<!-- 6-stdlib/99-minor/crypto/mldsa/77626.md -->
<!-- crypto/mldsa is documented in doc/next/6-stdlib/70-mldsa.md. -->

#### [`crypto/tls`](/pkg/crypto/tls/)

<!-- 6-stdlib/99-minor/crypto/tls/77363.md -->

The new [`QUICConfig.ClientHelloInfoConn`](/pkg/crypto/tls#QUICConfig.ClientHelloInfoConn) field specifies the [`net.Conn`](/pkg/net#Conn) to use
for the [`ClientHelloInfo.Conn`](/pkg/crypto/tls#ClientHelloInfo.Conn) field during QUIC server handshakes.

<!-- 6-stdlib/99-minor/crypto/tls/78543.md -->

The [`MLKEM1024`](/pkg/crypto/tls#MLKEM1024) key exchange is now supported. It can be enabled by adding it to
[`Config.CurvePreferences`](/pkg/crypto/tls#Config.CurvePreferences).

<!-- 6-stdlib/99-minor/crypto/tls/78888.md -->
<!-- crypto/tls ML-DSA support is documented in doc/next/6-stdlib/70-mldsa.md. -->

<!-- 6-stdlib/99-minor/crypto/tls/79367.md -->

[`Config.Rand`](/pkg/crypto/tls#Config.Rand) is now deprecated.
For deterministic testing, use [`testing/cryptotest.SetGlobalRandom`](/pkg/testing/cryptotest#SetGlobalRandom).

<!-- 6-stdlib/99-minor/crypto/tls/tlsmlkem.md -->

Post-quantum hybrid key exchanges can now be explicitly enabled in
[`Config.CurvePreferences`](/pkg/crypto/tls#Config.CurvePreferences) even if the `tlsmlkem=0` or `tlssecpmlkem=0` `GODEBUG`
options are used. Those options were always meant to only apply to the default
set used when [`Config.CurvePreferences`](/pkg/crypto/tls#Config.CurvePreferences) is nil.

#### [`crypto/x509`](/pkg/crypto/x509/)

<!-- 6-stdlib/99-minor/crypto/x509/75260.md -->

When parsing into [`pkix.Name`](/pkg/pkix#Name) fields, a wider range of
[`pkix.AttributeTypeAndValue.Value`](/pkg/pkix#AttributeTypeAndValue.Value) types is now supported, and unknown types are
parsed into [`asn1.RawValue`](/pkg/asn1#RawValue).

<!-- 6-stdlib/99-minor/crypto/x509/76133.md -->

The new [`Certificate.RawSignatureAlgorithm`](/pkg/crypto/x509#Certificate.RawSignatureAlgorithm), [`CertificateRequest.RawSignatureAlgorithm`](/pkg/crypto/x509#CertificateRequest.RawSignatureAlgorithm),
and [`RevocationList.RawSignatureAlgorithm`](/pkg/crypto/x509#RevocationList.RawSignatureAlgorithm) fields expose the DER-encoded
AlgorithmIdentifier of the signature algorithm, including when the
SignatureAlgorithm field is [`UnknownSignatureAlgorithm`](/pkg/crypto/x509#UnknownSignatureAlgorithm).

<!-- 6-stdlib/99-minor/crypto/x509/77865.md -->

[`SystemCertPool`](/pkg/crypto/x509#SystemCertPool) now respects SSL_CERT_FILE and SSL_CERT_DIR on Windows and
Darwin. When these environment variables are set, roots are loaded from disk and
instead of using the platform certificate verification APIs, the native Go
verifier is used. This behavior can be disabled with
`GODEBUG=x509sslcertoverrideplatform=0`.

<!-- 6-stdlib/99-minor/crypto/x509/78888.md -->
<!-- crypto/x509 ML-DSA support is documented in doc/next/6-stdlib/70-mldsa.md. -->

#### [`crypto/x509/pkix`](/pkg/crypto/x509/pkix/)

<!-- 6-stdlib/99-minor/crypto/x509/pkix/33093.md -->

[`RDNSequence.String`](/pkg/crypto/x509/pkix#RDNSequence.String) (and therefore [`Name.String`](/pkg/crypto/x509/pkix#Name.String)) now renders string-typed
attribute values as strings even when the attribute's OID is unrecognized.
Previously such values were always hex-encoded in their DER form.
See [#33093](/issue/33093).

#### [`database/sql`](/pkg/database/sql/)

<!-- 6-stdlib/99-minor/database/sql/67546.md -->

The new [`ConvertAssign`](/pkg/database/sql#ConvertAssign) function gives database drivers access
to the type conversions performed by [`Rows.Scan`](/pkg/database/sql#Rows.Scan).

#### [`database/sql/driver`](/pkg/database/sql/driver/)

<!-- 6-stdlib/99-minor/database/sql/driver/67546.md -->

Drivers may implement the new [`RowsColumnScanner`](/pkg/database/sql/driver#RowsColumnScanner) interface
to scan directly into user-provided destinations.

#### [`go/constant`](/pkg/go/constant/)

<!-- 6-stdlib/99-minor/go/constant/79042.md -->

The new [`StringLen`](/pkg/go/constant#StringLen) function returns the length of a string [`Value`](/pkg/go/constant#Value). For an [`Unknown`](/pkg/go/constant#Unknown) value, the length is 0.

#### [`go/scanner`](/pkg/go/scanner/)

<!-- 6-stdlib/99-minor/go/scanner/74958.md -->

The scanner now allows retrieving the end position of a token via the new [`Scanner.End`](/pkg/go/scanner#Scanner.End) method.

#### [`go/token`](/pkg/go/token/)

<!-- 6-stdlib/99-minor/go/token/76285.md -->

[`File`](/pkg/go/token#File) now has a String method.

#### [`go/types`](/pkg/go/types/)

<!-- 6-stdlib/99-minor/go/types/69420.md -->

The [`Hasher`](/pkg/go/types#Hasher) type is an implementation of `maphash.Hasher` for [Type]s
that respects the [`Identical`](/pkg/go/types#Identical) equivalence relation, allowing `Types`
to be used in hash tables and similar data structures (see `container/hash`).
[`HasherIgnoreTags`](/pkg/go/types#HasherIgnoreTags) is the analogous hasher for [`IdenticalIgnoreTags`](/pkg/go/types#IdenticalIgnoreTags).

<!-- 6-stdlib/99-minor/go/types/76472.md -->
<!-- CL 736441 -->

The `gotypesalias` `GODEBUG` setting (added in [Go 1.22](/doc/godebug#go-122))
has been removed permanently and the package [go/types](https://pkg.go.dev/go/types)
now always produces an [Alias](https://pkg.go.dev/go/types#Alias) type node for
[alias declarations](/ref/spec#Alias_declarations) irrespective of `GODEBUG` settings.

<!-- 6-stdlib/99-minor/go/types/79287.md -->
<!-- nothing to see here but some String methods -->

#### [`hash/maphash`](/pkg/hash/maphash/)

<!-- 6-stdlib/99-minor/hash/maphash/70471.md -->

The [`Hasher`](/pkg/hash/maphash#Hasher) interface type defines the contract between values of a
particular type and future hash-based data structures such as hash
tables and Bloom filters; see [#70471](/issue/70471).

#### [`math/big`](/pkg/math/big/)

<!-- 6-stdlib/99-minor/math/big/76821.md -->
<!-- go.dev/issue/76821 -->

[`Int`](/pkg/math/big#Int) now has method [`Int.Divide`](/pkg/math/big#Int.Divide) to compute quotient and remainder of two [`Int`](/pkg/math/big#Int) values.
It supports rounding modes [`Trunc`](/pkg/math/big#Trunc), [`Floor`](/pkg/math/big#Floor), [`Round`](/pkg/math/big#Round) and [`Ceil`](/pkg/math/big#Ceil).

#### [`math/rand/v2`](/pkg/math/rand/v2/)

<!-- 6-stdlib/99-minor/math/rand/v2/77853.md -->

add the generic method [`*Rand`](/pkg/math/rand/v2#Rand).N, matching the behavior of the top-level N function.

#### [`net`](/pkg/net/)

<!-- 6-stdlib/99-minor/net/78137.md -->

[`UnixConn`](/pkg/net#UnixConn) read methods now return [`io.EOF`](/pkg/io#EOF) directly instead of wrapping it in [`net.OpError`](/pkg/net#OpError) when the underlying read returns EOF.

#### [`net/http`](/pkg/net/http/)

<!-- 6-stdlib/99-minor/net/http/21753.md -->

[`Transport`](/pkg/net/http#Transport) and [`Server`](/pkg/net/http#Server) support TLS ALPN protocol negotiation on
user-provided [`net.Conn`](/pkg/net#Conn) connections which implement a
`ConnectionState() tls.ConnectionState` method.

<!-- 6-stdlib/99-minor/net/http/75500.md -->

HTTP/2 server now accepts client priority signals, as defined in RFC 9218,
allowing it to prioritize serving HTTP/2 streams with higher priority. If the
old behavior is preferred, where streams are served in a round-robin manner
regardless of priority, [`Server.DisableClientPriority`](/pkg/net/http#Server.DisableClientPriority) can be set to `true`.

<!-- 6-stdlib/99-minor/net/http/77370.md -->

HTTP/1 [`Response.Body`](/pkg/net/http#Response.Body) now automatically drains any unread content upon being
closed, up to a conservative limit, to allow better connection reuse. For most
programs, this change should be a no-op, or result in a performance improvement.
In rare cases, programs that do not benefit from connection reuse might
experience performance degradation if they had been improperly allowing an
excessive amount of idle connections to linger; usually by setting
[`Transport.MaxIdleConns`](/pkg/net/http#Transport.MaxIdleConns) to `0` or using different [Client]s for different
requests, thereby bypassing [`Transport.MaxIdleConns`](/pkg/net/http#Transport.MaxIdleConns) limit. In these cases,
setting [`Transport.DisableKeepAlives`](/pkg/net/http#Transport.DisableKeepAlives) to `true` will disable connection reuse.
However, such performance degradation usually indicates improper configuration
or usage of [`Transport`](/pkg/net/http#Transport) or [`Client`](/pkg/net/http#Client) in the first place, and a deeper look would
likely be beneficial.

#### [`net/http/httptest`](/pkg/net/http/httptest/)

<!-- 6-stdlib/99-minor/net/http/httptest/76608.md -->

[`NewTestServer`](/pkg/net/http/httptest#NewTestServer) creates a [`Server`](/pkg/net/http/httptest#Server) configured to use an in-memory
fake network suitable for use with the [`testing/synctest`](/pkg/testing/synctest) package.

#### [`net/url`](/pkg/net/url/)

<!-- 6-stdlib/99-minor/net/url/73450.md -->

The new [`URL.Clone`](/pkg/net/url#URL.Clone) method creates a deep copy of a URL.
The new [`Values.Clone`](/pkg/net/url#Values.Clone) method creates a deep copy of Values.

#### [`runtime/secret`](/pkg/runtime/secret/)

Goroutines that are created while in [`secret mode`](/pkg/runtime/secret#Do)
will now themselves execute in secret mode.

#### [`strings`](/pkg/strings/)

<!-- 6-stdlib/99-minor/strings/71151.md -->

The new [`CutLast`](/pkg/strings#CutLast) function slices a string
around the last occurrence of a separator.
It can replace and simplify some common uses of LastIndex.

#### [`testing/synctest`](/pkg/testing/synctest/)

<!-- 6-stdlib/99-minor/testing/synctest/77169.md -->

The new [`Sleep`](/pkg/testing/synctest#Sleep) helper function combines [`time.Sleep`](/pkg/time#Sleep) and [`testing/synctest.Wait`](/pkg/testing/synctest#Wait).

#### [`unicode`](/pkg/unicode/)

<!-- 6-stdlib/99-minor/unicode/77266.md -->

The unicode package and associated support throughout the system has been upgraded from Unicode 15 to Unicode 17.
See the [Unicode 16.0.0](https://www.unicode.org/versions/Unicode16.0.0/) and
[Unicode 17.0.0](https://www.unicode.org/versions/Unicode17.0.0/)
release notes for information about the changes.

#### [`uuid`](/pkg/uuid/)

<!-- 6-stdlib/99-minor/uuid/62026.md -->
<!-- uuid is documented in its own section. -->

## Ports {#ports}

### Darwin {#darwin}

<!-- go.dev/issue/75836 -->

As [announced](go1.26#darwin) in the Go 1.26 release notes,
Go 1.27 requires macOS 13 Ventura or later;
support for previous versions has been discontinued.

### Linux {#linux}

<!-- go.dev/issue/76244 -->

On ppc64, the ABI has been migrated to ELFv2. This change
has no effect for those building and running pure Go
binaries.

On ppc64, for maximum backward compatibility, existing users
should explicitly disable cgo, external linking, and use the
default build mode. That is, set `CGO_ENABLED=0` in the
environment, and pass the option `-ldflags="-linkmode=internal"`
to go build and similar. Linux ELFv2 support was added in
3.13, RHEL7 backported this support to its 3.10 kernel.

On ppc64, external linking, cgo, and PIE binaries are now
supported. Using these features requires an ELFv2 compatible
runtime (libc, kernel, and all linked and loaded libraries).

## TODO

Please convert these into documentation in the right places.
Some of them may not need any documentation or may be false
positives from automation.

### TODO: CL 774621 has a RELNOTE comment without a suggested text (from RELNOTE comment in [/cl/774621](/cl/774621))

- `internal/goexperiment,runtime: drop goroutineleakprofile experiment`

### TODO: accepted proposal [/issue/62728](/issue/62728) (from [/cl/601535](/cl/601535), [/cl/628615](/cl/628615), [/cl/751940](/cl/751940))

- `testing: annotate output text type`
- `testing: annotate output text type`
- `cmd/internal/test2json: generate and validate test artifacts`
- `testing: escapes framing markers`

### TODO: accepted proposal [/issue/63741](/issue/63741) (from [/cl/723102](/cl/723102))

- `doc/godebug: allow carve out for GODEBUGs introduced in security releases`
- `doc: document GODEBUG carve out for security releases`

### TODO: accepted proposal [/issue/69985](/issue/69985) (from [/cl/777220](/cl/777220))

- `crypto/tls: add X25519MLKEM768 and use by default; remove x25519Kyber768Draft00`
- `crypto/tls: let Config.CurvePreferences override GODEBUG options`

### TODO: accepted proposal [/issue/71206](/issue/71206) (from [/cl/777220](/cl/777220))

- `crypto/tls: add support for NIST curve based ML-KEM hybrids`
- `crypto/tls: let Config.CurvePreferences override GODEBUG options`

### TODO: accepted proposal [/issue/74609](/issue/74609) (from [/cl/774620](/cl/774620), [/cl/774621](/cl/774621), [/cl/775621](/cl/775621))

- `runtime/pprof,runtime: new goroutine leak profile`
- `internal/buildcfg: enable goroutineleakprofile GOEXPERIMENT by default`
- `internal/goexperiment,runtime: drop goroutineleakprofile experiment`
- `internal/goexperiment: actually delete goroutineleakprofile experiment`

### TODO: accepted proposal [/issue/75154](/issue/75154) (from [/cl/747160](/cl/747160))

- `crypto/sha3: make the zero value of SHA3 useable`
- `crypto/sha3: ensure unwrapped *sha3.Digest are usable`

### TODO: accepted proposal [/issue/75316](/issue/75316) (from [/cl/777380](/cl/777380), [/cl/777381](/cl/777381), [/cl/777382](/cl/777382), [/cl/777383](/cl/777383), [/cl/777384](/cl/777384))

- `crypto: remove in Go 1.27 GODEBUGs introduced in Go 1.23 and earlier`
- `crypto/tls: remove the tlsunsafeekm GODEBUG setting`
- `crypto/tls: remove tlsrsakex GODEBUG setting`
- `crypto/tls: remove tls3des GODEBUG setting`
- `crypto/tls: remove the tls10server GODEBUG setting`
- `crypto/tls: remove the x509keypairleaf GODEBUG setting`
