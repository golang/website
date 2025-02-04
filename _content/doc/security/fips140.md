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
algorithms, and the use of CMVP-validated cryptographic modules tested in the
target operating environments.

The mechanisms described in this page facilitate compliance for Go applications.

Applications that have no need for FIPS 140-3 compliance can safely ignore them,
and should not enable FIPS 140-3 mode.

## The Go Cryptographic Module

The Go Cryptographic Module is a collection of standard library Go packages
under `crypto/internal/fips140/...` that implement FIPS 140-3 approved
algorithms.

Public API packages such as `crypto/ecdsa` and `crypto/rand` transparently use
the Go Cryptographic Module to implement FIPS 140-3 algorithms.

Go Cryptographic Module version v1.0.0 is currently under test with a
CMVP-accredited laboratory.

## FIPS 140-3 mode

The run-time `fips140` [GODEBUG](/doc/godebug) option controls whether the Go
Cryptographic Module operates in FIPS 140-3 mode. It defaults to `off`. It can't
be changed after the program has started.

When operating in FIPS 140-3 mode (the `fips140` GODEBUG setting is `on`):

 - The Go Cryptographic Module automatically performs an integrity self-check at
   `init` time, comparing the checksum of the module's object file computed at
   build time with the symbols loaded in memory.

 - All algorithms perform known-answer self-tests according to the relevant FIPS
   140-3 Implementation Guidance, either at `init` time, or on first use.

 - Pairwise consistency tests are performed on generated cryptographic keys.
   Note that this can cause a slowdown of up to 2x for certain key types, which
   is especially relevant for ephemeral keys.

 - [`crypto/rand.Reader`](/pkg/crypto/rand/#Reader) is implemented in terms of a
   NIST SP 800-90A DRBG. To guarantee the same level of security as
   `GODEBUG=fips140=off`, random bytes are sourced from the platform's CSPRNG at
   every `Read` and mixed into the output as uncredited additional data.

 - The [`crypto/tls`](/pkg/crypto/tls/) package will ignore and not negotiate
   any protocol version, cipher suite, signature algorithm, or key exchange
   mechanism that is not compliant with NIST SP 800-52r2.

 - [`crypto/rsa.SignPSS`](/pkg/crypto/rsa/#SignPSS) with
   [`PSSSaltLengthAuto`](/pkg/crypto/rsa/#PSSSaltLengthAuto) will cap the length
   of the salt at the length of the hash.

When `GODEBUG=fips140=only` is used, in addition to the above, cryptographic
algorithms that are not FIPS 140-3 compliant will return an error or panic. Note
that this mode is a best effort and can't guarantee compliance with all FIPS
140-3 requirements.

`GODEBUG=fips140=on` and `only` are not supported on OpenBSD, Wasm, AIX, and
32-bit Windows platforms.

## The `crypto/fips140` package

The [`crypto/fips140.Enabled`](/pkg/crypto/fips140/#Enabled) function reports
whether FIPS 140-3 mode is active.

## The `GOFIPS140` environment variable

The `GOFIPS140` environment variable can be used with `go build`, `go install`,
and `go test` to select the version of the Go Cryptographic Module to be linked
into the executable program.

 - `off` is the default, and uses the `crypto/internal/fips140/...` packages in
   the standard library tree in use.

 - `latest` is like `off`, but enables FIPS 140-3 mode by default.

 - `v1.0.0` uses Go Cryptographic Module version v1.0.0, frozen in early 2025
   and first shipped with Go 1.24. It enables FIPS 140-3 mode by default.

## Go+BoringCrypto

The previous, unsupported mechanism to use the BoringCrypto module for certain
FIPS 140-3 approved algorithms is currently still available, but it is meant to
be removed and replaced with the mechanism described in this page in a future
release.

Go+BoringCrypto is incompatible with the native FIPS 140-3 mode.
