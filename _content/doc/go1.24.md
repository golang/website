---
title: Go 1.24 Release Notes
template: false
---

<style>
  main ul li { margin: 0.5em 0; }
</style>

## Introduction to Go 1.24 {#introduction}

The latest Go release, version 1.24,
arrives in [February 2025](/doc/devel/release#go1.24.0),
six months after [Go 1.23](/doc/go1.23).
Most of its changes are in the implementation of the toolchain, runtime, and libraries.
As always, the release maintains the Go 1 [promise of compatibility](/doc/go1compat).
We expect almost all Go programs to continue to compile and run as before.

## Changes to the language {#language}

<!-- go.dev/issue/46477 -->
Go 1.24 now fully supports [generic type aliases](/issue/46477): a type alias
may be parameterized like a defined type.
See the [language spec](/ref/spec#Alias_declarations) for details.
For now, the feature can be disabled by setting `GOEXPERIMENT=noaliastypeparams`;
but the `aliastypeparams` setting will be removed for Go 1.25.

## Tools {#tools}

### Go command {#go-command}

<!-- go.dev/issue/48429 -->

Go modules can now track executable dependencies using `tool` directives in
go.mod. This removes the need for the previous workaround of adding tools as
blank imports to a file conventionally named "tools.go". The `go tool`
command can now run these tools in addition to tools shipped with the Go
distribution. For more information see [the
documentation](/doc/modules/managing-dependencies#tools).

The new `-tool` flag for `go get` causes a tool directive to be added to the
current module for named packages in addition to adding require directives.

The new [`tool` meta-pattern](/cmd/go#hdr-Package_lists_and_patterns) refers to
all tools in the current module. This can be used to upgrade them all with `go get tool` or to install them into your GOBIN directory with `go install tool`.

<!-- go.dev/issue/69290 -->

Executables created by `go run` and the new behavior of `go tool` are now
cached in the Go build cache. This makes repeated executions faster at the
expense of making the cache larger. See [#69290](/issue/69290).

<!-- go.dev/issue/62067 -->

The `go build` and `go install` commands now accept a `-json` flag that reports
build output and failures as structured JSON output on standard output.
For details of the reporting format, see `go help buildjson`.

Furthermore, `go test -json` now reports build output and failures in JSON,
interleaved with test result JSON.
These are distinguished by new `Action` types, but if they cause problems in
a test integration system, you can revert to the text build output with
[GODEBUG setting](/doc/godebug) `gotestjsonbuildtext=1`.

<!-- go.dev/issue/26232 -->

The new `GOAUTH` environment variable provides a flexible way to authenticate
private module fetches. See `go help goauth` for more information.

<!-- go.dev/issue/50603 -->

The `go build` command now sets the
[main module's version](/pkg/runtime/debug#BuildInfo.Main) in the compiled
binary based on the version control system tag and/or commit.
A `+dirty` suffix will be appended if there are uncommitted changes.
Use the `-buildvcs=false` flag to omit version control information from the binary.

<!-- go.dev/issue/63939 -->

The new [GODEBUG setting](/doc/godebug) [`toolchaintrace=1`](/doc/toolchain#select)
can be used to trace the `go` command's toolchain selection process.

### Cgo {#cgo}

<!-- go.dev/issue/56378, CL 579955 -->
Cgo supports new annotations for C functions to improve run time
performance.
`#cgo noescape cFunctionName` tells the compiler that memory passed to
the C function `cFunctionname` does not escape.
`#cgo nocallback cFunctionName` tells the compiler that the C function
`cFunctionName` does not call back to any Go functions.
For more information, see [the cgo documentation](/pkg/cmd/cgo#hdr-Optimizing_calls_of_C_code).

<!-- go.dev/issue/67699 -->
Cgo currently refuses to compile calls to a C function which has multiple
incompatible declarations. For instance, if `f` is declared as both `void f(int)`
and `void f(double)`, cgo will report an error instead of possibly generating an
incorrect call sequence for `f(0)`. New in this release is a better detector for
this error condition when the incompatible declarations appear in different
files. See [#67699](/issue/67699).

### Objdump

<!-- go.dev/issue/15255, go.dev/issue/36738 -->
The [objdump](/cmd/objdump) tool now supports dissassembly on 64-bit
LoongArch (`GOARCH=loong64`), RISC-V (`GOARCH=riscv64`), and S390X (`GOARCH=s390x`).

### Vet

<!-- go.dev/issue/44251 -->
The new `tests` analyzer reports common mistakes in declarations of
tests, fuzzers, benchmarks, and examples in test packages, such as
malformed names, incorrect signatures, or examples that document
non-existent identifiers. Some of these mistakes may cause tests not
to run.
This analyzer is among the subset of analyzers that are run by `go test`.

<!-- go.dev/issue/60529 -->
The existing `printf` analyzer now reports a diagnostic for calls of
the form `fmt.Printf(s)`, where `s` is a non-constant format string,
with no other arguments. Such calls are nearly always a mistake
as the value of `s` may contain the `%` symbol; use `fmt.Print` instead.
See [#60529](/issue/60529). This check tends to produce findings in existing
code, and so is only applied when the language version (as specified by the
go.mod `go` directive or `//go:build` comments) is at least Go 1.24, to avoid
causing continuous integration failures when updating to the 1.24 Go toolchain.

<!-- go.dev/issue/64127 -->
The existing `buildtag` analyzer now reports a diagnostic when
there is an invalid Go [major version build constraint](/pkg/cmd/go#hdr-Build_constraints)
within a `//go:build` directive. For example, `//go:build go1.23.1` refers to
a point release; use `//go:build go1.23` instead.
See [#64127](/issue/64127).

<!-- go.dev/issue/66387 -->
The existing `copylock` analyzer now reports a diagnostic when a
variable declared in a 3-clause "for" loop such as
`for i := iter(); done(i); i = next(i) { ... }` contains a `sync.Locker`,
such as a `sync.Mutex`. [Go 1.22](/doc/go1.22#language) changed the behavior
of these loops to create a new variable for each iteration, copying the
value from the previous iteration; this copy operation is not safe for locks.
See [#66387](/issue/66387).

### GOCACHEPROG

<!-- go.dev/issue/64876 -->
The `cmd/go` internal binary and test caching mechanism can now be implemented
by child processes implementing a JSON protocol between the `cmd/go` tool
and the child process named by the `GOCACHEPROG` environment variable.
This was previously behind a GOEXPERIMENT.
For protocol details, see [the documentation](/cmd/go/internal/cacheprog).

## Runtime {#runtime}

<!-- go.dev/issue/54766 -->
<!-- go.dev/cl/614795 -->
<!-- go.dev/issue/68578 -->

Several performance improvements to the runtime have decreased CPU overheads by
2–3% on average across a suite of representative benchmarks.
Results may vary by application.
These improvements include a new builtin `map` implementation based on
[Swiss Tables](https://abseil.io/about/design/swisstables), more efficient
memory allocation of small objects, and a new runtime-internal mutex
implementation.

The new builtin `map` implementation and new runtime-internal mutex may be
disabled by setting `GOEXPERIMENT=noswissmap` and `GOEXPERIMENT=nospinbitmutex`
at build time respectively.

## Compiler {#compiler}

<!-- go.dev/issue/60725, go.dev/issue/57926 -->
The compiler already disallowed defining new methods with receiver types that were
cgo-generated, but it was possible to circumvent that restriction via an alias type.
Go 1.24 now always reports an error if a receiver denotes a cgo-generated type,
whether directly or indirectly (through an alias type).

## Linker {#linker}

<!-- go.dev/issue/68678, go.dev/issue/68652, CL 618598, CL 618601 -->
The linker now generates a GNU build ID (the ELF `NT_GNU_BUILD_ID` note) on ELF platforms
and a UUID (the Mach-O `LC_UUID` load command) on macOS by default.
The build ID or UUID is derived from the Go build ID.
It can be disabled by the `-B none` linker flag, or overridden by the `-B 0xNNNN` linker
flag with a user-specified hexadecimal value.

## Bootstrap {#bootstrap}

<!-- go.dev/issue/64751 -->
As mentioned in the [Go 1.22 release notes](/doc/go1.22#bootstrap), Go 1.24 now requires
Go 1.22.6 or later for bootstrap.
We expect that Go 1.26 will require a point release of Go 1.24 or later for bootstrap.

## Standard library {#library}

### Directory-limited filesystem access

<!-- go.dev/issue/67002 -->
The new [`os.Root`](/pkg/os#Root) type provides the ability to perform filesystem
operations within a specific directory.

The [`os.OpenRoot`](/pkg/os#OpenRoot) function opens a directory and returns an [`os.Root`](/pkg/os#Root).
Methods on [`os.Root`](/pkg/os#Root) operate within the directory and do not permit
paths that refer to locations outside the directory, including
ones that follow symbolic links out of the directory.
The methods on `os.Root` mirror most of the file system operations available in the
`os` package, including for example [`os.Root.Open`](/pkg/os#Root.Open),
[`os.Root.Create`](/pkg/os#Root.Create),
[`os.Root.Mkdir`](/pkg/os#Root.Mkdir),
and [`os.Root.Stat`](/pkg/os#Root.Stat),

### New benchmark function

Benchmarks may now use the faster and less error-prone [`testing.B.Loop`](/pkg/testing#B.Loop) method to perform benchmark iterations like `for b.Loop() { ... }` in place of the typical loop structures involving `b.N` like `for range b.N`. This offers two significant advantages:
 - The benchmark function will execute exactly once per -count, so expensive setup and cleanup steps execute only once.
 - Function call parameters and results are kept alive, preventing the compiler from fully optimizing away the loop body.

### Improved finalizers

<!-- go.dev/issue/67535 -->
The new [`runtime.AddCleanup`](/pkg/runtime#AddCleanup) function is a
finalization mechanism that is more flexible, more efficient, and less
error-prone than [`runtime.SetFinalizer`](/pkg/runtime#SetFinalizer).
`AddCleanup` attaches a cleanup function to an object that will run once
the object is no longer reachable.
However, unlike `SetFinalizer`,
multiple cleanups may be attached to a single object,
cleanups may be attached to interior pointers,
cleanups do not generally cause leaks when objects form a cycle, and
cleanups do not delay the freeing of an object or objects it points to.
New code should prefer `AddCleanup` over `SetFinalizer`.

### New weak package {#weak}

The new [`weak`](/pkg/weak/) package provides weak pointers.

Weak pointers are a low-level primitive provided to enable the
creation of memory-efficient structures, such as weak maps for
associating values, canonicalization maps for anything not
covered by package [`unique`](/pkg/unique/), and various kinds
of caches.
For supporting these use-cases, this release also provides
[`runtime.AddCleanup`](/pkg/runtime/#AddCleanup) and
[`maphash.Comparable`](/pkg/maphash/#Comparable).

### New crypto/mlkem package {#crypto-mlkem}

<!-- go.dev/issue/70122 -->

The new [`crypto/mlkem`](/pkg/crypto/mlkem/) package implements
ML-KEM-768 and ML-KEM-1024.

ML-KEM is a post-quantum key exchange mechanism formerly known as Kyber and
specified in [FIPS 203](https://doi.org/10.6028/NIST.FIPS.203).

### New crypto/hkdf, crypto/pbkdf2, and crypto/sha3 packages {#crypto-packages}

<!-- go.dev/issue/61477, go.dev/issue/69488, go.dev/issue/69982, go.dev/issue/65269, CL 629176 -->

The new [`crypto/hkdf`](/pkg/crypto/hkdf/) package implements
the HMAC-based Extract-and-Expand key derivation function HKDF,
as defined in [RFC 5869](https://www.rfc-editor.org/rfc/rfc5869.html).

The new [`crypto/pbkdf2`](/pkg/crypto/pbkdf2/) package implements
the password-based key derivation function PBKDF2,
as defined in [RFC 8018](https://www.rfc-editor.org/rfc/rfc8018.html).

The new [`crypto/sha3`](/pkg/crypto/sha3/) package implements
the SHA-3 hash function and SHAKE and cSHAKE extendable-output functions,
as defined in [FIPS 202](http://doi.org/10.6028/NIST.FIPS.202).

All three packages are based on pre-existing `golang.org/x/crypto/...` packages.

### FIPS 140-3 compliance {#fips140}

This release includes [a new set of mechanisms to facilitate FIPS 140-3
compliance](/doc/security/fips140).

The Go Cryptographic Module is a set of internal standard library packages that
are transparently used to implement FIPS 140-3 approved algorithms. Applications
require no changes to use the Go Cryptographic Module for approved algorithms.

The new `GOFIPS140` environment variable can be used to select the Go
Cryptographic Module version to use in a build. The new `fips140` [GODEBUG
setting](/doc/godebug) can be used to enable FIPS 140-3 mode at runtime.

Go 1.24 includes Go Cryptographic Module version v1.0.0, which is currently
under test with a CMVP-accredited laboratory.

### New experimental testing/synctest package {#testing-synctest}

The new experimental [`testing/synctest`](/pkg/testing/synctest/) package
provides support for testing concurrent code.
- The [`synctest.Run`](/pkg/testing/synctest/#Run) function starts a
  group of goroutines in an isolated "bubble".
  Within the bubble, [`time`](/pkg/time) package functions operate on a
  fake clock.
- The [`synctest.Wait`](/pkg/testing/synctest#Wait) function waits for
  all goroutines in the current bubble to block.

See the package documentation for more details.

The `synctest` package is experimental and must be enabled by
setting `GOEXPERIMENT=synctest` at build time.
The package API is subject to change in future releases.
See [issue #67434](/issue/67434) for more information and
to provide feeback.

### Minor changes to the library {#minor_library_changes}

#### [`archive`](/pkg/archive/)

The `(*Writer).AddFS` implementations in both `archive/zip` and `archive/tar`
now write a directory header for an empty directory.

#### [`bytes`](/pkg/bytes/)

The [`bytes`](/pkg/bytes) package adds several functions that work with iterators:
- [`Lines`](/pkg/bytes#Lines) returns an iterator over the
  newline-terminated lines in a byte slice.
- [`SplitSeq`](/pkg/bytes#SplitSeq) returns an iterator over
  all subslices of a byte slice split around a separator.
- [`SplitAfterSeq`](/pkg/bytes#SplitAfterSeq) returns an iterator
  over subslices of a byte slice split after each instance of a
  separator.
- [`FieldsSeq`](/pkg/bytes#FieldsSeq) returns an iterator over
  subslices of a byte slice split around runs of whitespace characters,
  as defined by [`unicode.IsSpace`](/pkg/unicode#IsSpace).
- [`FieldsFuncSeq`](/pkg/bytes#FieldsFuncSeq) returns an iterator
  over subslices of a byte slice split around runs of Unicode code points
  satisfying a predicate.

#### [`crypto/aes`](/pkg/crypto/aes/)

The value returned by [`NewCipher`](/pkg/crypto/aes#NewCipher) no longer
implements the `NewCTR`, `NewGCM`, `NewCBCEncrypter`, and `NewCBCDecrypter`
methods. These methods were undocumented and not available on all architectures.
Instead, the [`Block`](/pkg/crypto/cipher#Block) value should be passed
directly to the relevant [`crypto/cipher`](/pkg/crypto/cipher/) functions.
For now, `crypto/cipher` still checks for those methods on `Block` values,
even if they are not used by the standard library anymore.

#### [`crypto/cipher`](/pkg/crypto/cipher/)

The new [`NewGCMWithRandomNonce`](/pkg/crypto/cipher#NewGCMWithRandomNonce)
function returns an [`AEAD`](/pkg/crypto/cipher#AEAD) that implements AES-GCM by
generating a random nonce during Seal and prepending it to the ciphertext.

The [`Stream`](/pkg/crypto/cipher#Stream) implementation returned by
[`NewCTR`](/pkg/crypto/cipher#NewCTR) when used with
[`crypto/aes`](/pkg/crypto/aes/) is now several times faster on amd64 and arm64.

[`NewOFB`](/pkg/crypto/cipher#NewOFB),
[`NewCFBEncrypter`](/pkg/crypto/cipher#NewCFBEncrypter), and
[`NewCFBDecrypter`](/pkg/crypto/cipher#NewCFBDecrypter) are now deprecated.
OFB and CFB mode are not authenticated, which generally enables active attacks to
manipulate and recover the plaintext. It is recommended that applications use
[`AEAD`](/pkg/crypto/cipher#AEAD) modes instead. If an unauthenticated
[`Stream`](/pkg/crypto/cipher#Stream) mode is required, use
[`NewCTR`](/pkg/crypto/cipher#NewCTR) instead.

#### [`crypto/ecdsa`](/pkg/crypto/ecdsa/)

<!-- go.dev/issue/64802 -->
[`PrivateKey.Sign`](/pkg/crypto/ecdsa#PrivateKey.Sign) now produces a
deterministic signature according to
[RFC 6979](https://www.rfc-editor.org/rfc/rfc6979.html) if the random source is nil.

#### [`crypto/md5`](/pkg/crypto/md5/)

The value returned by [`md5.New`](/pkg/md5#New) now also implements the
[`encoding.BinaryAppender`](/pkg/encoding#BinaryAppender) interface.

#### [`crypto/rand`](/pkg/crypto/rand/)

<!-- go.dev/issue/66821 -->
The [`Read`](/pkg/crypto/rand#Read) function is now guaranteed not to fail.
It will always return `nil` as the `error` result.
If `Read` were to encounter an error while reading from
[`Reader`](/pkg/crypto/rand#Reader), the program will irrecoverably crash.
Note that the platform APIs used by the default `Reader` are documented to
always succeed, so this change should only affect programs that override the
`Reader` variable. One exception are Linux kernels before version 3.17, where
the default `Reader` still opens `/dev/urandom` and may fail.

<!-- go.dev/issue/69577 -->
On Linux 6.11 and later, `Reader` now uses the `getrandom` system call via vDSO.
This is several times faster, especially for small reads.

<!-- CL 608395 -->
On OpenBSD, `Reader` now uses `arc4random_buf(3)`.

<!-- go.dev/issue/67057 -->
The new [`Text`](/pkg/crypto/rand#Text) function can be used to generate
cryptographically secure random text strings.

#### [`crypto/rsa`](/pkg/crypto/rsa/)

[`GenerateKey`](/pkg/crypto/rsa#GenerateKey) now returns an error if a key of
less than 1024 bits is requested.
All Sign, Verify, Encrypt, and Decrypt methods now return an error if used with
a key smaller than 1024 bits. Such keys are insecure and should not be used.
[GODEBUG setting](/doc/godebug) `rsa1024min=0` restores the old behavior, but we
recommend doing so only if necessary and only in tests, for example by adding a
`//go:debug rsa1024min=0` line to a test file.
A new `GenerateKey` [example](/pkg/crypto/rsa#example-GenerateKey-TestKey)
provides an easy-to-use standard 2048-bit test key.

It is now safe and more efficient to call
[`PrivateKey.Precompute`](/pkg/crypto/rsa#PrivateKey.Precompute) before
[`PrivateKey.Validate`](/pkg/crypto/rsa#PrivateKey.Validate).
`Precompute` is now faster in the presence of partially filled out
[`PrecomputedValues`](/pkg/crypto/rsa#PrecomputedValues), such as when
unmarshaling a key from JSON.

The package now rejects more invalid keys, even when `Validate` is not called,
and [`GenerateKey`](/pkg/crypto/rsa#GenerateKey) may return new errors for
broken random sources.
The [`Primes`](/pkg/crypto/rsa#PrivateKey.Primes) and
[`Precomputed`](/pkg/crypto/rsa#PrivateKey.Precomputed) fields of
[`PrivateKey`](/pkg/crypto/rsa#PrivateKey) are now used and validated even when
some values are missing.
See also the changes to `crypto/x509` parsing and marshaling of RSA keys
[described below](#cryptox509pkgcryptox509).

<!-- go.dev/issue/43923 -->
[`SignPKCS1v15`](/pkg/crypto/rsa#SignPKCS1v15) and
[`VerifyPKCS1v15`](/pkg/crypto/rsa#VerifyPKCS1v15) now support
SHA-512/224, SHA-512/256, and SHA-3.

<!-- CL 639936 -->
[`GenerateKey`](/pkg/crypto/rsa#GenerateKey) now uses a slightly different
method to generate the private exponent (Carmichael's totient instead of Euler's
totient). Rare applications that externally regenerate keys from only the prime
factors may produce different but compatible results.

<!-- CL 626957 -->
Public and private key operations are now up to two times faster on wasm.

#### [`crypto/sha1`](/pkg/crypto/sha1/)

The value returned by [`sha1.New`](/pkg/sha1#New) now also implements
the [`encoding.BinaryAppender`](/pkg/encoding#BinaryAppender) interface.

#### [`crypto/sha256`](/pkg/crypto/sha256/)

The values returned by [`sha256.New`](/pkg/sha256#New) and
[`sha256.New224`](/pkg/sha256#New224) now also implement the
[`encoding.BinaryAppender`](/pkg/encoding#BinaryAppender) interface.

#### [`crypto/sha512`](/pkg/crypto/sha512/)

The values returned by [`sha512.New`](/pkg/sha512#New),
[`sha512.New384`](/pkg/sha512#New384),
[`sha512.New512_224`](/pkg/sha512#New512_224) and
[`sha512.New512_256`](/pkg/sha512#New512_256) now also implement the
[`encoding.BinaryAppender`](/pkg/encoding#BinaryAppender) interface.

#### [`crypto/subtle`](/pkg/crypto/subtle/)

The new [`WithDataIndependentTiming`](/pkg/crypto/subtle#WithDataIndependentTiming)
function allows the user to run a function with architecture specific features
enabled which guarantee specific instructions are data value timing invariant.
This can be used to make sure that code designed to run in constant time is not
optimized by CPU-level features such that it operates in variable time.
Currently, `WithDataIndependentTiming` uses the PSTATE.DIT bit on arm64, and is
a no-op on all other architectures. [GODEBUG setting](/doc/godebug)
`dataindependenttiming=1` enables the DIT mode for the entire Go program.

<!-- CL 622276 -->
The [`XORBytes`](/pkg/crypto/subtle#XORBytes) output must overlap exactly or not
at all with the inputs. Previously, the behavior was otherwise undefined, while
now `XORBytes` will panic.

#### [`crypto/tls`](/pkg/crypto/tls/)

The TLS server now supports Encrypted Client Hello (ECH). This feature can be
enabled by populating the [`Config.EncryptedClientHelloKeys`](/pkg/crypto/tls#Config.EncryptedClientHelloKeys) field.

The new post-quantum [`X25519MLKEM768`](/pkg/crypto/tls#X25519MLKEM768) key
exchange mechanism is now supported and is enabled by default when
[`Config.CurvePreferences`](/pkg/crypto/tls#Config.CurvePreferences) is nil.
[GODEBUG setting](/doc/godebug) `tlsmlkem=0` reverts the default.
This can be useful when dealing with buggy TLS servers that do not handle large records correctly,
causing a timeout during the handshake (see [TLS post-quantum TL;DR fail](https://tldr.fail/)).

Support for the experimental `X25519Kyber768Draft00` key exchange has been removed.

<!-- go.dev/issue/69393, CL 630775 -->
Key exchange ordering is now handled entirely by the `crypto/tls` package. The
order of [`Config.CurvePreferences`](/pkg/crypto/tls#Config.CurvePreferences) is
now ignored, and the contents are only used to determine which key exchanges to
enable when the field is populated.

<!-- go.dev/issue/32936 -->
The new [`ClientHelloInfo.Extensions`](/pkg/crypto/tls#ClientHelloInfo.Extensions)
field lists the IDs of the extensions received in the Client Hello message.
This can be useful for fingerprinting TLS clients.

#### [`crypto/x509`](/pkg/crypto/x509/)

<!-- go.dev/issue/41682 -->
The `x509sha1` [GODEBUG setting](/doc/godebug) has been removed.
[`Certificate.Verify`](/pkg/crypto/x509#Certificate.Verify) no longer
supports SHA-1 based signatures.

[`OID`](/pkg/crypto/x509#OID) now implements the
[`encoding.BinaryAppender`](/pkg/encoding#BinaryAppender) and
[`encoding.TextAppender`](/pkg/encoding#TextAppender) interfaces.

The default certificate policies field has changed from
[`Certificate.PolicyIdentifiers`](/pkg/crypto/x509#Certificate.PolicyIdentifiers)
to [`Certificate.Policies`](/pkg/crypto/x509#Certificate.Policies). When parsing
certificates, both fields will be populated, but when creating certificates
policies will now be taken from the `Certificate.Policies` field instead of
the `Certificate.PolicyIdentifiers` field. This change can be reverted with
[GODEBUG setting](/doc/godebug) `x509usepolicies=0`.

<!-- go.dev/issue/67675 -->
[`CreateCertificate`](/pkg/crypto/x509#CreateCertificate) will now generate a
serial number using a RFC 5280 compliant method when passed a template with a
nil [`Certificate.SerialNumber`](/pkg/crypto/x509#Certificate.SerialNumber)
field, instead of failing.

[`Certificate.Verify`](/pkg/crypto/x509#Certificate.Verify) now supports policy
validation, as defined in RFC 5280 and RFC 9618. The new
[`VerifyOptions.CertificatePolicies`](/pkg/crypto/x509#VerifyOptions.CertificatePolicies)
field can be set to an acceptable set of policy [`OIDs`](/pkg/crypto/x509#OID).
Only certificate chains with valid policy graphs will be returned from
[`Certificate.Verify`](/pkg/crypto/x509#Certificate.Verify).

[`MarshalPKCS8PrivateKey`](/pkg/crypto/x509#MarshalPKCS8PrivateKey) now returns
an error instead of marshaling an invalid RSA key.
([`MarshalPKCS1PrivateKey`](/pkg/crypto/x509#MarshalPKCS1PrivateKey) doesn't
have an error return, and its behavior when provided invalid keys continues to
be undefined.)

[`ParsePKCS1PrivateKey`](/pkg/crypto/x509#ParsePKCS1PrivateKey) and
[`ParsePKCS8PrivateKey`](/pkg/crypto/x509#ParsePKCS8PrivateKey) now use and
validate the encoded CRT values, so might reject invalid RSA keys that were
previously accepted. Use [GODEBUG setting](/doc/godebug) `x509rsacrt=0` to
revert to recomputing the CRT values.

#### [`debug/elf`](/pkg/debug/elf/)

<!-- go.dev/issue/63952 -->

The [`debug/elf`](/pkg/debug/elf) package adds support for handling symbol
versions in dynamic ELF (Executable and Linkable Format) files.
The new [`File.DynamicVersions`](/pkg/debug/elf#File.DynamicVersions) method
returns a list of dynamic versions defined in the ELF file.
The new [`File.DynamicVersionNeeds`](/pkg/debug/elf#File.DynamicVersionNeeds)
method returns a list of dynamic versions required by this ELF file that are
defined in other ELF objects.
Finally, the new [`Symbol.HasVersion`](/pkg/debug/elf#Symbol) and
[`Symbol.VersionIndex`](/pkg/debug/elf#Symbol) fields indicate the version of a
symbol.

#### [`encoding`](/pkg/encoding/)

Two new interfaces, [`TextAppender`](/pkg/encoding#TextAppender) and [`BinaryAppender`](/pkg/encoding#BinaryAppender), have been
introduced to append the textual or binary representation of an object
to a byte slice. These interfaces provide the same functionality as
[`TextMarshaler`](/pkg/encoding#TextMarshaler) and [`BinaryMarshaler`](/pkg/encoding#BinaryMarshaler), but instead of allocating a new slice
each time, they append the data directly to an existing slice.
These interfaces are now implemented by standard library types that
already implemented `TextMarshaler` and/or `BinaryMarshaler`.

#### [`encoding/json`](/pkg/encoding/json/)

<!-- go.dev/issue/45669 -->
When marshaling, a struct field with the new `omitzero` option in the struct field
tag will be omitted if its value is zero. If the field type has an `IsZero() bool`
method, that will be used to determine whether the value is zero. Otherwise, the
value is zero if it is [the zero value for its type](/ref/spec#The_zero_value).
The `omitzero` field tag is clearer and less error-prone than `omitempty` when
the intent is to omit zero values.
In particular, unlike `omitempty`, `omitzero` omits zero-valued
[`time.Time`](/pkg/time#Time) values, which is a common source of friction.

If both `omitempty` and `omitzero` are specified, the field will be omitted if the
value is either empty or zero (or both).

[`UnmarshalTypeError.Field`](/pkg/encoding/json#UnmarshalTypeError.Field) now includes embedded structs to provide more detailed error messages.

#### [`go/types`](/pkg/go/types/)

All `go/types` data structures that expose sequences using a pair of
methods such as `Len() int` and `At(int) T` now also have methods that
return iterators, allowing you to simplify code such as this:

```
params := fn.Type.(*types.Signature).Params()
for i := 0; i < params.Len(); i++ {
   use(params.At(i))
}
```

to this:

```
for param := range fn.Signature().Params().Variables() {
   use(param)
}
```

The methods are:
[`Interface.EmbeddedTypes`](/pkg/go/types#Interface.EmbeddedTypes),
[`Interface.ExplicitMethods`](/pkg/go/types#Interface.ExplicitMethods),
[`Interface.Methods`](/pkg/go/types#Interface.Methods),
[`MethodSet.Methods`](/pkg/go/types#MethodSet.Methods),
[`Named.Methods`](/pkg/go/types#Named.Methods),
[`Scope.Children`](/pkg/go/types#Scope.Children),
[`Struct.Fields`](/pkg/go/types#Struct.Fields),
[`Tuple.Variables`](/pkg/go/types#Tuple.Variables),
[`TypeList.Types`](/pkg/go/types#TypeList.Types),
[`TypeParamList.TypeParams`](/pkg/go/types#TypeParamList.TypeParams),
[`Union.Terms`](/pkg/go/types#Union.Terms).

#### [`hash/adler32`](/pkg/hash/adler32/)

The value returned by [`New`](/pkg/hash/adler32#New) now also implements the [`encoding.BinaryAppender`](/pkg/encoding#BinaryAppender) interface.

#### [`hash/crc32`](/pkg/hash/crc32/)

The values returned by [`New`](/pkg/hash/crc32#New) and [`NewIEEE`](/pkg/hash/crc32#NewIEEE) now also implement the [`encoding.BinaryAppender`](/pkg/encoding#BinaryAppender) interface.

#### [`hash/crc64`](/pkg/hash/crc64/)

The value returned by [`New`](/pkg/hash/crc64#New) now also implements the [`encoding.BinaryAppender`](/pkg/encoding#BinaryAppender) interface.

#### [`hash/fnv`](/pkg/hash/fnv/)

The values returned by [`New32`](/pkg/hash/fnv#New32), [`New32a`](/pkg/hash/fnv#New32a), [`New64`](/pkg/hash/fnv#New64), [`New64a`](/pkg/hash/fnv#New64a), [`New128`](/pkg/hash/fnv#New128) and [`New128a`](/pkg/hash/fnv#New128a) now also implement the [`encoding.BinaryAppender`](/pkg/encoding#BinaryAppender) interface.

#### [`hash/maphash`](/pkg/hash/maphash/)

The new [`Comparable`](/pkg/hash/maphash#Comparable) and
[`WriteComparable`](/pkg/hash/maphash#WriteComparable) functions can compute the
hash of any comparable value.
These make it possible to hash anything that can be used as a Go map key.

#### [`log/slog`](/pkg/log/slog/)

The new [`DiscardHandler`](/pkg/log/slog#DiscardHandler) is a handler that is never enabled and always discards its output.

[`Level`](/pkg/log/slog#Level) and [`LevelVar`](/pkg/log/slog#LevelVar) now implement the [`encoding.TextAppender`](/pkg/encoding#TextAppender) interface.

#### [`math/big`](/pkg/math/big/)

[`Float`](/pkg/math/big#Float), [`Int`](/pkg/math/big#Int) and [`Rat`](/pkg/math/big#Rat) now implement the [`encoding.TextAppender`](/pkg/encoding#TextAppender) interface.

#### [`math/rand`](/pkg/math/rand/)

Calls to the deprecated top-level [`Seed`](/pkg/math/rand#Seed) function no longer have any effect. To
restore the old behavior use [GODEBUG setting](/doc/godebug) `randseednop=0`. For more background see
[proposal #67273](/issue/67273).

#### [`math/rand/v2`](/pkg/math/rand/v2/)

[`ChaCha8`](/pkg/math/rand/v2#ChaCha8) and [`PCG`](/pkg/math/rand/v2#PCG) now implement the [`encoding.BinaryAppender`](/pkg/encoding#BinaryAppender) interface.

#### [`net`](/pkg/net/)

[`ListenConfig`](/pkg/net#ListenConfig) now uses MPTCP by default on systems where it is supported
(currently on Linux only).

[`IP`](/pkg/net#IP) now implements the [`encoding.TextAppender`](/pkg/encoding#TextAppender) interface.

#### [`net/http`](/pkg/net/http/)

[`Transport`](/pkg/net/http#Transport)'s limit on 1xx informational responses received
in response to a request has changed.
It previously aborted a request and returned an error after
receiving more than 5 1xx responses.
It now returns an error if the total size of all 1xx responses
exceeds the [`Transport.MaxResponseHeaderBytes`](/pkg/net/http#Transport.MaxResponseHeaderBytes) configuration setting.

In addition, when a request has a
[`net/http/httptrace.ClientTrace.Got1xxResponse`](/pkg/net/http/httptrace#ClientTrace.Got1xxResponse)
trace hook, there is now no limit on the total number of 1xx responses.
The `Got1xxResponse` hook may return an error to abort a request.

[`Transport`](/pkg/net/http#Transport) and [`Server`](/pkg/net/http#Server) now have an HTTP2 field which permits
configuring HTTP/2 protocol settings.

The new [`Server.Protocols`](/pkg/net/http#Server.Protocols) and [`Transport.Protocols`](/pkg/net/http#Transport.Protocols) fields provide
a simple way to configure what HTTP protocols a server or client use.

The server and client may be configured to support unencrypted HTTP/2
connections.

When [`Server.Protocols`](/pkg/net/http#Server.Protocols) contains UnencryptedHTTP2, the server will accept
HTTP/2 connections on unencrypted ports. The server can accept both
HTTP/1 and unencrypted HTTP/2 on the same port.

When [`Transport.Protocols`](/pkg/net/http#Transport.Protocols) contains UnencryptedHTTP2 and does not contain
HTTP1, the transport will use unencrypted HTTP/2 for http:// URLs.
If the transport is configured to use both HTTP/1 and unencrypted HTTP/2,
it will use HTTP/1.

Unencrypted HTTP/2 support uses "HTTP/2 with Prior Knowledge"
(RFC 9113, section 3.3). The deprecated "Upgrade: h2c" header
is not supported.

#### [`net/netip`](/pkg/net/netip/)

[`Addr`](/pkg/net/netip#Addr), [`AddrPort`](/pkg/net/netip#AddrPort) and [`Prefix`](/pkg/net/netip#Prefix) now implement the [`encoding.BinaryAppender`](/pkg/encoding#BinaryAppender) and
[`encoding.TextAppender`](/pkg/encoding#TextAppender) interfaces.

#### [`net/url`](/pkg/net/url/)

[`URL`](/pkg/net/url#URL) now also implements the [`encoding.BinaryAppender`](/pkg/encoding#BinaryAppender) interface.

#### [`os/user`](/pkg/os/user/)

On Windows, [`Current`](/pkg/os/user#Current) can now be used in Windows Nano Server.
The implementation has been updated to avoid using functions
from the `NetApi32` library, which is not available in Nano Server.

On Windows, [`Current`](/pkg/os/user#Current), [`Lookup`](/pkg/os/user#Lookup) and [`LookupId`](/pkg/os/user#LookupId) now support the
following built-in service user accounts:
- `NT AUTHORITY\SYSTEM`
- `NT AUTHORITY\LOCAL SERVICE`
- `NT AUTHORITY\NETWORK SERVICE`

On Windows, [`Current`](/pkg/os/user#Current) has been made considerably faster when
the current user is joined to a slow domain, which is the
usual case for many corporate users. The new implementation
performance is now in the order of milliseconds, compared to
the previous implementation which could take several seconds,
or even minutes, to complete.

On Windows, [`Current`](/pkg/os/user#Current) now returns the process owner user when
the current thread is impersonating another user. Previously,
it returned an error.

#### [`regexp`](/pkg/regexp/)

[`Regexp`](/pkg/regexp#Regexp) now implements the [`encoding.TextAppender`](/pkg/encoding#TextAppender) interface.

#### [`runtime`](/pkg/runtime/)

The [`GOROOT`](/pkg/runtime#GOROOT) function is now deprecated.
In new code prefer to use the system path to locate the “go” binary,
and use `go env GOROOT` to find its GOROOT.

#### [`strings`](/pkg/strings/)

The [`strings`](/pkg/strings) package adds several functions that work with iterators:
- [`Lines`](/pkg/strings#Lines) returns an iterator over
  the newline-terminated lines in a string.
- [`SplitSeq`](/pkg/strings#SplitSeq) returns an iterator over
  all substrings of a string split around a separator.
- [`SplitAfterSeq`](/pkg/strings#SplitAfterSeq) returns an iterator
  over substrings of a string split after each instance of a
  separator.
- [`FieldsSeq`](/pkg/strings#FieldsSeq) returns an iterator over
  substrings of a string split around runs of whitespace characters,
  as defined by [`unicode.IsSpace`](/pkg/unicode#IsSpace).
- [`FieldsFuncSeq`](/pkg/strings#FieldsFuncSeq) returns an iterator
  over substrings of a string split around runs of Unicode code points
  satisfying a predicate.

#### [`sync`](/pkg/sync/)

The implementation of [`sync.Map`](/pkg/sync#Map) has been changed, improving performance,
particularly for map modifications.
For instance, modifications of disjoint sets of keys are much less likely to contend on
larger maps, and there is no longer any ramp-up time required to achieve low-contention
loads from the map.

If you encounter any problems, set `GOEXPERIMENT=nosynchashtriemap` at build
time to switch back to the old implementation and please [file an
issue](/issue/new).

#### [`testing`](/pkg/testing/)

The new [`T.Context`](/pkg/testing#T.Context) and [`B.Context`](/pkg/testing#B.Context) methods return a context that's canceled
after the test completes and before test cleanup functions run.

<!-- testing.B.Loop mentioned in 6-stdlib/6-testing-bloop.md. -->

The new [`T.Chdir`](/pkg/testing#T.Chdir) and [`B.Chdir`](/pkg/testing#B.Chdir) methods can be used to change the working
directory for the duration of a test or benchmark.

#### [`text/template`](/pkg/text/template/)

Templates now support range-over-func and range-over-int.

#### [`time`](/pkg/time/)

[`Time`](/pkg/time#Time) now implements the [`encoding.BinaryAppender`](/pkg/encoding#BinaryAppender) and [`encoding.TextAppender`](/pkg/encoding#TextAppender) interfaces.

## Ports {#ports}

### Linux {#linux}

<!-- go.dev/issue/67001 -->
As [announced](go1.23#linux) in the Go 1.23 release notes, Go 1.24 requires Linux
kernel version 3.2 or later.

### Darwin {#darwin}

<!-- go.dev/issue/69839 -->
Go 1.24 is the last release that will run on macOS 11 Big Sur.
Go 1.25 will require macOS 12 Monterey or later.

### WebAssembly {#wasm}

<!-- go.dev/issue/65199, CL 603055 -->
The `go:wasmexport` compiler directive is added for Go programs to export functions
to the WebAssembly host.

On WebAssembly System Interface Preview 1 (`GOOS=wasip1 GOARCH=wasm`), Go 1.24 supports
building a Go program as a
[reactor/library](https://github.com/WebAssembly/WASI/blob/63a46f61052a21bfab75a76558485cf097c0dbba/legacy/application-abi.md#current-unstable-abi),
by specifying the `-buildmode=c-shared` build flag.

<!-- go.dev/issue/66984, CL 626615 -->
More types are now permitted as argument or result types for `go:wasmimport` functions.
Specifically, `bool`, `string`, `uintptr`, and pointers to certain types are allowed
(see the [documentation](/pkg/cmd/compile#hdr-WebAssembly_Directives) for detail),
along with 32-bit and 64-bit integer and float types, and `unsafe.Pointer`, which
are already allowed.
These types are also permitted as argument or result types for `go:wasmexport` functions.

<!-- go.dev/issue/68024 -->
The support files for WebAssembly have been moved to `lib/wasm` from `misc/wasm`.

<!-- CL 621635, CL 621636 -->
The initial memory size is significantly reduced, especially for small WebAssembly
applications.

### Windows {#windows}

<!-- go.dev/issue/70705 -->
The 32-bit windows/arm port (`GOOS=windows GOARCH=arm`) has been marked broken.
See [issue #70705](/issue/70705) for details.

<!-- Items that don't need to be mentioned in Go 1.24 release notes but are picked up by relnote todo.

accepted proposal https://go.dev/issue/25309 (from https://go.dev/cl/594018, https://go.dev/cl/595120, https://go.dev/cl/595564, https://go.dev/cl/601778) - new x/crypto package; doesn't seem to need to be mentioned
accepted proposal https://go.dev/issue/43744 (from https://go.dev/cl/357530) - no change in the standard library
accepted proposal https://go.dev/issue/60905 (from https://go.dev/cl/610195) - CL 610195 reverted
accepted proposal https://go.dev/issue/61395 (from https://go.dev/cl/594738, https://go.dev/cl/594976) - CL 594738 made sync/atomic AND/OR operations intrinsic on amd64, but the API was already added in Go 1.23; CL 594976 is a fix; probably doesn't require a Go 1.24 release note (performance change only)
accepted proposal https://go.dev/issue/51269 (from https://go.dev/cl/627035) - may be worth mentioning in Go 1.24 release notes, or may be fine to leave out; commented at https://go.dev/issue/51269#issuecomment-2501802763; Ian confirmed it's fine to leave out
accepted proposal https://go.dev/issue/66540 (from https://go.dev/cl/603958) - a Go language spec clarification; might not need to be mentioned in Go 1.24 release notes; left a comment at https://go.dev/issue/66540#issuecomment-2502051684; Robert confirmed it indeed doesn't
accepted proposal https://go.dev/issue/34208 (from https://go.dev/cl/586241) - CL 586241 implements a fix for a Go 1.23 feature, doesn't seem to be need anything in Go 1.24 release notes
accepted proposal https://go.dev/issue/43993 (from https://go.dev/cl/626116) - CL 626116 prepares the tree towards the vet change but the vet change itself isn't implemented in Go 1.24, so nothing to say in Go 1.24 release notes
accepted proposal https://go.dev/issue/44505 (from https://go.dev/cl/609955) - CL 609955 is an internal cleanup in x/tools, no need for Go 1.24 release note
accepted proposal https://go.dev/issue/61476 (from https://go.dev/cl/608255) - CL 608255 builds on GORISCV64 added in Go 1.23; nothing to mention in Go 1.24 release notes
accepted proposal https://go.dev/issue/66315 (from https://go.dev/cl/577996) - adding Pass.Module field to x/tools/go/analysis doesn't seem like something that needs to be mentioned in Go 1.24 release notes
accepted proposal https://go.dev/issue/57786 (from https://go.dev/cl/472717) - CL 472717 is in x/net/http2 and mentions a Go 1.21 proposal; it doesn't seem to need anything in Go 1.24 release notes
accepted proposal https://go.dev/issue/54265 (from https://go.dev/cl/609915, https://go.dev/cl/610675) - CLs that refer to a Go 1.22 proposal, nothing more is needed in Go 1.24 release notes
accepted proposal https://go.dev/issue/53021 (from https://go.dev/cl/622276) - CL 622276 improves docs; proposal 53021 was in Go 1.20 so nothing more is needed in Go 1.24 release notes
accepted proposal https://go.dev/issue/51430 (from https://go.dev/cl/613375) - CL 613375 is an internal documentation comment; proposal 51430 happened in Go 1.20/1.21 so nothing more is needed in Go 1.24 release notes
accepted proposal https://go.dev/issue/38445 (from https://go.dev/cl/626495) - CL 626495 works on proposal 38445, which is about x/tools/go/package, doesn't need anything in Go 1.24 release notes
accepted proposal https://go.dev/issue/56986 (from https://go.dev/cl/618115) - CL 618115 adds documentation; it doesn't need to be mentioned in Go 1.24 release notes
accepted proposal https://go.dev/issue/60061 (from https://go.dev/cl/612038) - CL 612038 is a CL that deprecates something in x/tools/go/ast and mentions a Go 1.22 proposal; doesn't need anything in Go 1.24 release notes
accepted proposal https://go.dev/issue/61324 (from https://go.dev/cl/411907) - CL 411907 is an x/tools CL that implements a proposal for a new package there; doesn't need anything in Go 1.24 release notes
accepted proposal https://go.dev/issue/61777 (from https://go.dev/cl/601496) - CL 601496 added a WriteByteTimeout field to x/net/http2.Server; doesn't need a Go 1.24 release note
accepted proposal https://go.dev/issue/61940 (from https://go.dev/cl/600997) - CL 600997 deleted obsolete code in x/build and mentioned an accepted proposal; doesn't need a Go 1.24 release note
accepted proposal https://go.dev/issue/62113 (from https://go.dev/cl/594195) - CL 594195 made iterator-related additions in x/net/html; doesn't need a Go 1.24 release note
accepted proposal https://go.dev/issue/62484 (from https://go.dev/cl/600775) - CL 600775 documents CopyFS symlink behavior and mentions the Go 1.23 proposal; doesn't need a Go 1.24 release note
accepted proposal https://go.dev/issue/64207 (from https://go.dev/cl/605875) - an x/website CL that follows up on a Go 1.23 proposal; doesn't need a Go 1.24 release note
accepted proposal https://go.dev/issue/65236 (from https://go.dev/cl/596135) - CL 596135 adds tests for the Go 1.23 proposal 65236; doesn't need a Go 1.24 release note
accepted proposal https://go.dev/issue/67795 (from https://go.dev/cl/616218) - iteratior support for x/tools/go/ast/inspector; doesn't need a Go 1.24 release note
accepted proposal https://go.dev/issue/67812 (from https://go.dev/cl/601497) - configurable server pings for x/net/http2.Server; doesn't need a Go 1.24 release note
accepted proposal https://go.dev/issue/68232 (from https://go.dev/cl/595676) - x/sys/unix additions; doesn't need a Go 1.24 release note
accepted proposal https://go.dev/issue/68898 (from https://go.dev/cl/607495, https://go.dev/cl/620036, https://go.dev/cl/620135, https://go.dev/cl/623638) - a proposal for x/tools/go/gcexportdata to document 2 releases + tip support policy; since the change is in x/tools it doesn't need a Go 1.24 release note
accepted proposal https://go.dev/issue/69095 (from https://go.dev/cl/593683, https://go.dev/cl/608955, https://go.dev/cl/610716) - a proposal that affects maintenance and support of golang.org/x repositories; doesn't need to be mentioned in Go 1.24 release notes
accepted proposal https://go.dev/issue/68384 (from https://go.dev/cl/611875) - expanding the scope of Go Telemetry to include Delve isn't directly tied to Go 1.24 and doesn't seem to need to be mentioned in Go 1.24 release notes
accepted proposal https://go.dev/issue/69291 (from https://go.dev/cl/610939) - CL 610939 refactors code in x/tools and mentions the still-open proposal #69291 to add Reachable to x/tools/go/ssa/ssautil; doesn't need a Go 1.24 release note
accepted proposal https://go.dev/issue/69360 (from https://go.dev/cl/614158, https://go.dev/cl/614159, https://go.dev/cl/614635, https://go.dev/cl/614675) - proposal 69360 is to tag and delete gorename from x/tools; doesn't need a Go 1.24 release note
accepted proposal https://go.dev/issue/61417 (from https://go.dev/cl/605955) - a new field in x/oauth2; nothing to mention in Go 1.24 release notes
accepted proposal https://go.dev/issue/29266 (from https://go.dev/cl/632897) - a documentation-only proposal for go.dev/doc/contribute; doesn't need a Go 1.24 release note
-->
