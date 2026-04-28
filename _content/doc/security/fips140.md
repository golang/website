---
title: FIPS 140-3 Compliance
layout: article
---

Starting with Go 1.24, Go binaries can natively operate in a mode that
facilitates FIPS 140-3 compliance. Moreover, the toolchain can build against
frozen versions of the cryptography packages that constitute the Go
Cryptographic Module.

## FIPS 140-3

NIST FIPS 140-3 is a U.S. Government compliance regime for cryptography
applications that amongst other things requires the use of a set of approved
algorithms, and the use of
[CMVP](https://csrc.nist.gov/projects/cryptographic-module-validation-program)-validated
cryptographic modules tested in the target operating environments.

The mechanisms described in this page facilitate compliance for Go applications.

Applications that have no need for FIPS 140-3 compliance can safely ignore them,
and should not enable FIPS 140-3 mode.

**NOTE:** Simply using a FIPS 140-3 compliant and validated cryptographic module
may not—on its own—satisfy all relevant regulatory requirements. The Go team cannot provide any
guarantees or support around how usage of the provided FIPS 140-3
mode may, or may not, satisfy specific regulatory requirements for individual
users. Care should be taken in determining if usage of this module satisfies
your specific requirements.

## The Go Cryptographic Module

The Go Cryptographic Module is a collection of standard library Go packages
under `crypto/internal/fips140/...` that implement FIPS 140-3 approved
algorithms.

Public API packages such as `crypto/ecdsa` and `crypto/rand` transparently use
the Go Cryptographic Module to implement FIPS 140-3 algorithms.

## FIPS 140-3 mode

When operating in FIPS 140-3 mode:

 - The Go Cryptographic Module automatically performs an integrity self-check at
   `init` time, comparing the checksum of the module's object file computed at
   build time with the symbols loaded in memory.

 - All algorithms perform known-answer self-tests according to the relevant FIPS
   140-3 Implementation Guidance, either at `init` time, or on first use.

 - Pairwise consistency tests are performed on generated cryptographic keys.
   Note that this can cause a slowdown of up to 2x for certain key types, which
   is especially relevant for ephemeral keys.

 - [`crypto/rand.Reader`](/pkg/crypto/rand/#Reader) is implemented in terms of a
   NIST SP 800-90A DRBG. To guarantee the same level of security as programs not
   running in FIPS 140-3 mode, random bytes are also sourced from the platform's
   CSPRNG at every `Read` and mixed into the output as uncredited additional
   data.

 - The [`crypto/tls`](/pkg/crypto/tls/) package will ignore and not negotiate
   any protocol version, cipher suite, signature algorithm, or key exchange
   mechanism that is not FIPS 140-3 approved. (This is equivalent to the legacy
   opt-in `crypto/tls/fipsonly` Go+BoringCrypto mechanism.)

 - [`crypto/rsa.SignPSS`](/pkg/crypto/rsa/#SignPSS) with
   [`PSSSaltLengthAuto`](/pkg/crypto/rsa/#PSSSaltLengthAuto) will cap the length
   of the salt at the length of the hash.

FIPS 140-3 mode is not supported on OpenBSD, Wasm, AIX, and 32-bit Windows.

## The `crypto/fips140` package

The [`crypto/fips140.Enabled`](/pkg/crypto/fips140/#Enabled) function reports
whether FIPS 140-3 mode is active.

The [`crypto/fips140.Version`](/pkg/crypto/fips140/#Version) function returns
the version of the Go Cryptographic Module in use.

## The `GOFIPS140` environment variable

The `GOFIPS140` environment variable can be used with `go build`, `go install`,
and `go test` to select the version of the Go Cryptographic Module to be linked
into the executable program, and to enable FIPS 140-3 mode by default.

 - `off` is the default, and uses the `crypto/internal/fips140/...` packages in
   the standard library tree in use.

 - `latest` is like `off`, but enables FIPS 140-3 mode by default.

 - `v1.0.0` uses Go Cryptographic Module version v1.0.0, frozen in early 2025
   and first shipped with Go 1.24. It enables FIPS 140-3 mode by default.

 - `v1.26.0` uses Go Cryptographic Module version v1.26.0, frozen in early 2026
   and first shipped with Go 1.26. It enables FIPS 140-3 mode by default.

## The `fips140` GODEBUG option

The run-time `fips140` [GODEBUG](/doc/godebug) option controls whether the Go
Cryptographic Module operates in FIPS 140-3 mode. It can't be changed after the
program has started.

It defaults to `off` unless `GOFIPS140` is set at build time.

If set to `on`, FIPS 140-3 mode is enabled. This is possible even if `GOFIPS140`
was not set at build time.

If set to `only`, cryptographic algorithms that are not FIPS 140-3 compliant
will return an error or panic. Note that this is a best effort mode meant for
testing, assessment, and debugging. *It is not intended to be used in
production*, it is not required by the Security Policy, it introduces crashes
and potentially unhandled errors by design, and it may have false positives or
false negatives.

Most programs should not set this option directly, and should instead use
`GOFIPS140` at build time.

## Module Validations

Google currently has a contractual relationship with [Geomys](https://geomys.org/)
to facilitate at least yearly CMVP validations of the Go Cryptographic Module.
At the time of validation we will freeze the Go Cryptographic Module and create
a new module version for submission.

These validations are tested on a comprehensive set of Operating
Environments, supporting many popular operating system and hardware platform
combinations.

Off-cycle validations may be performed if security issues are discovered in
the module.

###  Validated Module Versions

List of module versions which have completed [CMVP validation](https://csrc.nist.gov/projects/cryptographic-module-validation-program/validated-modules/search?SearchMode=Basic&ModuleName=Go+Cryptographic+Module&CertificateStatus=Active&ValidationYear=0):

* v1.0.0 ([CMVP Certificate #5247](https://csrc.nist.gov/projects/cryptographic-module-validation-program/certificate/5247)), available in Go 1.24+

### In Process Module Versions

List of module versions which are currently in the [CMVP Modules In Process List](https://csrc.nist.gov/Projects/cryptographic-module-validation-program/modules-in-process/modules-in-process-list):

* v1.26.0 ([CAVP Certificate A8028](https://csrc.nist.gov/projects/cryptographic-algorithm-validation-program/details?validation=40638)), Pending Review, available in Go 1.26+

### Implementation Under Test Module Versions

List of module versions which are currently in the [CMVP Implementation Under Test List](https://csrc.nist.gov/Projects/cryptographic-module-validation-program/modules-in-process/iut-list):

_None at this time._

## Go+BoringCrypto

The previous, unsupported mechanism to use the BoringCrypto module for certain
FIPS 140-3 approved algorithms is currently still available, but it is meant to
be removed and replaced with the mechanism described in this page in a future
release.

Go+BoringCrypto is incompatible with the native FIPS 140-3 mode.
