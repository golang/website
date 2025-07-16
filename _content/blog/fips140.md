---
title: The FIPS 140-3 Go Cryptographic Module
date: 2025-07-15
by:
- Filippo Valsorda (Geomys)
- Daniel McCarney (Geomys)
- Roland Shoemaker (Google)
summary: Go now has a built-in, native FIPS 140-3 compliant mode.
---

FIPS 140 is a standard for cryptography implementations and, although it doesn’t
necessarily improve security, FIPS 140 compliance is a requirement in certain
regulated environments that are increasingly adopting Go. Until now, FIPS 140
compliance has been a significant source of friction for Go users, requiring
unsupported solutions with safety, developer experience, functionality, release
velocity, and compliance issues.

Go is addressing this growing need with native FIPS 140 support built right into
the standard library and the `go` command, making Go the easiest, most secure
way to comply with FIPS 140. The FIPS 140-3 validated Go Cryptographic Module
now underlies Go’s built-in crypto libraries, starting with the Go Cryptographic
Module v1.0.0 that is included in Go 1.24, released last February.

The v1.0.0 module has been awarded [Cryptographic Algorithm Validation Program
(CAVP) certificate A6650][], was submitted to the Cryptographic Module
Validation Program (CMVP), and reached the [Modules In Process List][] in May.
Modules on the MIP list are awaiting NIST review and can already be deployed in
certain regulated environments.

[Geomys][] led the implementation effort in collaboration with the Go Security
Team, and is pursuing a broadly applicable FIPS 140-3 validation for the benefit
of the Go community. Google and other industry stakeholders have a contractual
relationship with Geomys to include specific Operating Environments in the
certificate.

Further details on the module are available in the
[documentation](/doc/security/fips140).

Some Go users currently rely on the [Go+BoringCrypto][] GOEXPERIMENT, or on one
of its forks, as part of their FIPS 140 compliance strategy. Unlike the FIPS
140-3 Go Cryptographic Module, Go+BoringCrypto was never officially supported
and had significant developer experience issues, since it was produced
exclusively for the internal needs of Google. It will be removed in a future
release once Google migrates to the native module.

## A native developer experience

The module integrates completely transparently into Go applications. In fact,
every Go program built with Go 1.24 already uses it for all FIPS 140-3 approved
algorithms! The module is just another name for the
`crypto/internal/fips140/...` packages of the standard library, which provide
the implementation of operations exposed by packages such as `crypto/ecdsa` and
`crypto/rand`.

These packages involve no cgo, meaning they cross-compile like any other Go
program, they pay no FFI performance overhead, and they don’t suffer from
[memory management security issues][], unlike Go+BoringCrypto and its forks.

When starting a Go binary, the module can be put into FIPS 140-3 mode with the
`fips140=on` [GODEBUG option][], which can be set as an environment variable or
through the `go.mod` file. If FIPS 140-3 mode is enabled, the module will use
the NIST DRBG for randomness, `crypto/tls` will automatically only negotiate
FIPS 140-3 approved TLS versions and algorithms, and it will perform the
mandatory self-tests on initialization and during key generation. That’s it;
there are no other behavior differences.

There is also an experimental stricter mode, `fips140=only`, which causes all
non-approved algorithms to return errors or panic. We understand this might be
too inflexible for most deployments and are [looking for
feedback](/issue/74630) on what a policy enforcement framework
might look like.

Finally, applications can use the [`GOFIPS140` environment
variable](/doc/security/fips140#the-gofips140-environment-variable)
to build against older, validated versions of the `crypto/internal/fips140/...`
packages. `GOFIPS140` works like `GOOS` and `GOARCH`, and if set to
`GOFIPS140=v1.0.0` the program will be built against the v1.0.0 snapshot of the
packages as they were submitted for validation to CMVP. This snapshot ships with
the rest of the Go standard library, as `lib/fips140/v1.0.0.zip`.

When using `GOFIPS140`, the `fips140` GODEBUG defaults to `on`, so putting it
all together, all that’s needed to build against the FIPS 140-3 module and run
in FIPS 140-3 mode is `GOFIPS140=v1.0.0 go build`. That’s it.

If a toolchain is built with `GOFIPS140` set, all builds it produces will
default to that value.

The `GOFIPS140` version used to build a binary can be verified with
`go version -m`.

Future versions of Go will continue shipping and working with v1.0.0 of the Go
Cryptographic Module until the next version is fully certified by Geomys, but
some new cryptography features might not be available when building against old
modules. Starting with Go 1.24.3, you can use `GOFIPS140=inprocess` to
dynamically select the latest module for which a Geomys validation has reached
the In Process stage. Geomys plans to validate new module versions at least
every year—to avoid leaving FIPS 140 builds too far behind—and every time a
vulnerability in the module can’t be mitigated in the calling standard library
code.

## Uncompromising security

Our first priority in developing the module has been matching or exceeding the
security of the existing Go standard library cryptography packages. It might be
surprising, but sometimes the easiest way to achieve and demonstrate compliance
with the FIPS 140 security requirements is not to exceed them. We declined to
accept that.

For example, `crypto/ecdsa` [always produced hedged signatures][]. Hedged
signatures generate nonces by combining the private key, the message, and random
bytes. Like [deterministic ECDSA][RFC 6979], they protect against failure of the
random number generator, which would otherwise leak the private key(!). Unlike
deterministic ECDSA, they are also resistant to [API issues][] and [fault
attacks][], and they don’t leak message equality. FIPS 186-5 introduced support
for [RFC 6979][] deterministic ECDSA, but not for hedged ECDSA.

Instead of downgrading to regular randomized or deterministic ECDSA signatures
in FIPS 140-3 mode (or worse, across modes), we [switched the hedging
algorithm][] and connected dots across half a dozen documents to [prove the new
one is a compliant composition of a DRBG and traditional ECDSA][]. While at it,
we also [added opt-in support for deterministic signatures][].

Another example is random number generation. FIPS 140-3 has strict rules on how
cryptographic randomness is generated, which essentially enforce the use of a
userspace [CSPRNG][]. Conversely, we believe the kernel is best suited to
produce secure random bytes, because it’s best positioned to collect entropy
from the system, and to detect when processes or even virtual machines are
cloned (which could lead to reuse of supposedly random bytes). Hence,
[crypto/rand][] routes every read operation to the kernel.

To square this circle, in FIPS 140-3 mode we maintain a compliant userspace NIST
DRBG based on AES-256-CTR, and then inject into it 128 bits sourced from the
kernel at every read operation. This extra entropy is considered “uncredited”
additional data for FIPS 140-3 purposes, but in practice makes it as strong as
reading directly from the kernel—even if slower.

Finally, all of the Go Cryptographic Module v1.0.0 was in scope for the [recent
security audit by Trail of Bits](/blog/tob-crypto-audit), and was
not affected by the only non-informational finding.

Combined with the memory safety guarantees provided by the Go compiler and
runtime, we believe this delivers on our goal of making Go one of the easiest,
most secure solutions for FIPS 140 compliance.

## Broad platform support

A FIPS 140-3 module is only compliant if operated on a tested or “Vendor
Affirmed” Operating Environment, essentially a combination of operating system
and hardware platform. To enable as many Go use cases as possible, the Geomys
validation is tested on [one of the most comprehensive sets of Operating
Environments][] in the industry.

Geomys’s laboratory tested various Linux flavors (Alpine Linux on Podman, Amazon
Linux, Google Prodimage, Oracle Linux, Red Hat Enterprise Linux, and SUSE Linux
Enterprise Server), macOS, Windows, and FreeBSD on a mix of x86-64 (AMD and
Intel), ARMv8/9 (Ampere Altra, Apple M, AWS Graviton, and Qualcomm Snapdragon),
ARMv7, MIPS, z/ Architecture, and POWER, for a total of 23 tested environments.

Some of these were paid for by stakeholders, others were funded by Geomys for
the benefit of the Go community.

Moreover, the Geomys validation lists a broad set of generic platforms as Vendor
Affirmed Operating Environments:
* Linux 3.10+ on x86-64 and ARMv7/8/9,
* macOS 11–15 on Apple M processors,
* FreeBSD 12–14 on x86-64,
* Windows 10 and Windows Server 2016–2022 on x86-64, and
* Windows 11 and Windows Server 2025 on x86-64 and ARMv8/9.

## Comprehensive algorithm coverage

It may be surprising, but even using a FIPS 140-3 approved algorithm implemented
by a FIPS 140-3 module on a supported Operating Environment is not necessarily
enough for compliance; the algorithm must have been specifically covered by
testing as part of validation. Hence, to make it as easy as possible to build
FIPS 140 compliant applications in Go, all FIPS 140-3 approved algorithms in the
standard library are implemented by the Go Cryptographic Module and were tested
as part of the validation, from digital signatures to the TLS key schedule.

The post-quantum ML-KEM key exchange (FIPS 203), [introduced in Go 1.24][mlkem
relnote], is also validated, meaning `crypto/tls` can establish FIPS 140-3
compliant post-quantum secure connections with X25519MLKEM768.

In some cases, we validated the same algorithms under multiple different NIST
designations, to make it possible to use them in full compliance for different
purposes. For example, [HKDF is tested and validated under *four* names][hkdf]:
SP 800-108 Feedback KDF, SP 800-56C two-step KDF, Implementation Guidance D.P
OneStepNoCounter KDF, and SP 800-133 Section 6.3 KDF.

Finally, we validated some internal algorithms such as CMAC Counter KDF, to make
it possible to expose future functionality such as [XAES-256-GCM][].

Overall, the native FIPS 140-3 module delivers a better compliance profile than
Go+BoringCrypto, while making more algorithms available to FIPS 140-3 restricted
applications.

We look forward to the new native Go Cryptographic Module making it easier and
safer for Go developers to run FIPS 140 compliant workloads.

[Geomys]: https://geomys.org
[Cryptographic Algorithm Validation Program (CAVP) certificate A6650]: https://csrc.nist.gov/projects/cryptographic-algorithm-validation-program/details?validation=39260
[Modules In Process List]: https://csrc.nist.gov/Projects/cryptographic-module-validation-program/modules-in-process/modules-in-process-list
[Go+BoringCrypto]: /doc/security/fips140#goboringcrypto
[memory management security issues]: /blog/tob-crypto-audit#cgo-memory-management
[GODEBUG option]: /doc/godebug
[always produced hedged signatures]: https://cs.opensource.google/go/go/+/refs/tags/go1.23.0:src/crypto/ecdsa/ecdsa.go;l=417
[API issues]: https://github.com/MystenLabs/ed25519-unsafe-libs
[fault attacks]: https://en.wikipedia.org/wiki/Differential_fault_analysis
[RFC 6979]: https://www.rfc-editor.org/rfc/rfc6979
[switched the hedging algorithm]: https://github.com/golang/go/commit/9776d028f4b99b9a935dae9f63f32871b77c49af
[prove the new one is a compliant composition of a DRBG and traditional ECDSA]: https://github.com/cfrg/draft-irtf-cfrg-det-sigs-with-noise/issues/6#issuecomment-2067819904
[added opt-in support for deterministic signatures]: /doc/go1.24#cryptoecdsapkgcryptoecdsa
[CSPRNG]: https://en.wikipedia.org/wiki/Cryptographically_secure_pseudorandom_number_generator
[crypto/rand]: https://pkg.go.dev/crypto/rand
[one of the most comprehensive sets of Operating Environments]: https://csrc.nist.gov/projects/cryptographic-algorithm-validation-program/details?product=19371&displayMode=Aggregated
[mlkem relnote]: /doc/go1.24#crypto-mlkem
[hkdf]: https://words.filippo.io/dispatches/fips-hkdf/
[XAES-256-GCM]: https://c2sp.org/XAES-256-GCM
