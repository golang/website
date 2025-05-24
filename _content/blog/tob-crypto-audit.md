---
title: "Go Cryptography Security Audit"
date: 2025-05-19
by:
- Roland Shoemaker and Filippo Valsorda
summary: Go's cryptography libraries underwent an audit by Trail of Bits.
---

Go ships with a full suite of cryptography packages in the standard library to help developers build secure applications. Google recently contracted the independent security firm [Trail of Bits](https://www.trailofbits.com/) to complete an audit of the core set of packages that are also validated as part of the [new native FIPS 140-3 module](/doc/go1.24#fips140). The audit produced a single low-severity finding, in the legacy and unsupported [Go+BoringCrypto integration](/doc/security/fips140#goboringcrypto), and a handful of informational findings. The full text of the audit report can be found [here](https://github.com/trailofbits/publications/blob/d47e8fafa7e3323e5620d228f2f3f3bf58ed5978/reviews/2025-03-google-gocryptographiclibraries-securityreview.pdf).

The scope of the audit included our implementations of key exchange (ECDH and post-quantum ML-KEM), digital signature (ECDSA, RSA, and Ed25519), encryption (AES-GCM, AES-CBC, and AES-CTR), hashing (SHA-1, SHA-2, and SHA-3), key derivation (HKDF and PBKDF2), and authentication (HMAC), as well as the cryptographic random number generator. Low-level big integer and elliptic curve implementations, with their delicate assembly cores, were included. Higher level protocols like TLS and X.509 were not in scope. Three Trail of Bits engineers worked on the audit for a month.

We are proud of the security track record of the Go cryptography packages, and of the outcome of this audit, which is just one of many ways we gain assurance of the packages’ correctness. First, we aggressively limit their complexity, guided by the [Cryptography Principles](/design/cryptography-principles) which for example prioritize security over performance. Further, we [thoroughly test them](https://www.youtube.com/watch?v=lkEH3V3PkS0) with an array of different techniques. We make a point of leveraging safe APIs even for internal packages, and naturally we can rely on the Go language properties to avoid memory management issues. Finally, we focus on readability to make maintenance easier and code review and audits more effective.

## One low-severity finding

The only potentially exploitable issue, TOB-GOCL-3, has *low severity*, meaning it had minor impact and was difficult to trigger. This issue has been fixed in Go 1.24.

Crucially, TOB-GOCL-3 ([discussed further below](#cgo-memory-management)) concerns memory management in the [legacy Go+BoringCrypto GOEXPERIMENT](/doc/security/fips140#goboringcrypto), which is not enabled by default and unsupported for use outside of Google.

## Five informational findings

The remaining findings are *informational*, meaning they do not pose an immediate risk but are relevant to security best practices. We addressed these in the current Go 1.25 development tree.

Findings TOB-GOCL-1, TOB-GOCL-2, and TOB-GOCL-6 concern possible timing side-channels in various cryptographic operations. Of these three findings, only TOB-GOCL-2 affects operations that were expected to be constant time due to operating on secret values, but it only affects Power ISA targets (GOARCH ppc64 and ppc64le). TOB-GOCL-4 highlights misuse risk in an internal API, should it be repurposed beyond its current use case. TOB-GOCL-5 points out a missing check for a limit that is impractical to reach.

## Timing Side-Channels

Findings TOB-GOCL-1, TOB-GOCL-2, and TOB-GOCL-6 concern minor timing side-channels. TOB-GOCL-1 and TOB-GOCL-6 are related to functions which we do not use for sensitive values, but could be used for such values in the future, and TOB-GOCL-2 is related to the assembly implementation of P-256 ECDSA on Power ISA.

### `crypto/ecdh,crypto/ecdsa`: conversion from bytes to field elements is not constant time (TOB-GOCL-1)

The internal implementation of NIST elliptic curves provided a method to convert field elements between an internal and external representation which operated in variable time.

All usages of this method operated on public inputs which are not considered secret (public ECDH values, and ECDSA public keys), so we determined that this was not a security issue. That said, we decided to [make the method constant time anyway](/cl/650579), in order to prevent accidentally using this method in the future with secret values, and so that we don't have to think about whether it is an issue or not.

### `crypto/ecdsa`: P-256 conditional negation is not constant time in Power ISA assembly (TOB-GOCL-2, CVE-2025-22866)

Beyond the [first class Go platforms](/wiki/PortingPolicy#first-class-ports), Go also supports a number of additional platforms, including some less common architectures. During the review of our assembly implementations of various underlying cryptographic primitives, the Trail of Bits team found one issue that affected the ECDSA implementation on the ppc64 and ppc64le architectures.

Due to the usage of a conditional branching instruction in the implementation of the conditional negation of P-256 points, the function operated in variable-time, rather than constant-time, as expected. The fix for this was relatively simple, [replacing the conditional branching instruction](/cl/643735) with a pattern we already use elsewhere to conditionally select the correct result in constant time. We assigned this issue CVE-2025-22866.

To prioritize the code that reaches most of our users, and due to the specialized knowledge required to target specific ISAs, we generally rely on community contributions to maintain assembly for non-first class platforms. We thank our partners at IBM for helping provide review for our fix.

### `crypto/ed25519`: Scalar.SetCanonicalBytes is not constant time (TOB-GOCL-6)

The internal edwards25519 package provided a method to convert between an internal and external representation of scalars which operated in variable time.

This method was only used on signature inputs to ed25519.Verify, which are not considered secret, so we determined that this was not a security issue. That said, similarly to the TOB-GOCL-1 finding, we decided to [make the method constant time anyway](/cl/648035), in order to prevent accidentally using this method in the future with secret values, and because we are aware that people fork this code outside of the standard library, and may be using it with secret values.

## Cgo Memory Management

Finding TOB-GOCL-3 concerns a memory management issue in the Go+BoringCrypto integration.

### `crypto/ecdh`: custom finalizer may free memory at the start of a C function call using this memory (TOB-GOCL-3)

During the review, there were a number of questions about our cgo-based Go+BoringCrypto integration, which provides a FIPS 140-2 compliant cryptography mode for internal usage at Google. The Go+BoringCrypto code is not supported by the Go team for external use, but has been critical for Google’s internal usage of Go.

The Trail of Bits team found one vulnerability and one [non-security relevant bug](/cl/644120), both of which were results of the manual memory management required to interact with a C library. Since the Go team does not support usage of this code outside of Google, we have chosen not to issue a CVE or Go vulnerability database entry for this issue, but we [fixed it in Go 1.24](/cl/644119).

This kind of pitfall is one of the many reasons that we decided to move away from the Go+BoringCrypto integration. We have been working on a [native FIPS 140-3 mode](/doc/security/fips140) that uses the regular pure Go cryptography packages, allowing us to avoid the complex cgo semantics in favor of the traditional Go memory model.

## Implementation Completeness

Findings TOB-GOCL-4 and TOB-GOCL-5 concern limited implementations of two specifications, [NIST SP 800-90A](https://csrc.nist.gov/pubs/sp/800/90/a/r1/final) and [RFC 8018](https://datatracker.ietf.org/doc/html/rfc8018).

### `crypto/internal/fips140/drbg`: CTR\_DRBG API presents multiple misuse risks (TOB-GOCL-4)

As part of the [native FIPS 140-3 mode](/doc/security/fips140) that we are introducing, we needed an implementation of the NIST CTR\_DRBG (an AES-CTR based deterministic random bit generator) to provide compliant randomness.

Since we only need a small subset of the functionality of the NIST SP 800-90A Rev. 1 CTR\_DRBG for our purposes, we did not implement the full specification, in particular omitting the derivation function and personalization strings. These features can be critical to safely use the DRBG in generic contexts.

As our implementation is tightly scoped to the specific use case we need, and since the implementation is not publicly exported, we determined that this was acceptable and worth the decreased complexity of the implementation. We do not expect this implementation to ever be used for other purposes internally, and have [added a warning to the documentation](/cl/647815) that details these limitations.

### `crypto/pbkdf2`: PBKDF2 does not enforce output length limitations (TOB-GOCL-5)

In Go 1.24, we began the process of moving packages from [golang.org/x/crypto](https://golang.org/x/crypto) into the standard library, ending a confusing pattern where first-party, high-quality, and stable Go cryptography packages were kept outside of the standard library for no particular reason.

As part of this process we moved the [golang.org/x/crypto/pbkdf2](https://golang.org/x/crypto/pbkdf2) package into the standard library, as crypto/pbkdf2. While reviewing this package, the Trail of Bits team noticed that we did not enforce the limit on the size of derived keys defined in [RFC 8018](https://datatracker.ietf.org/doc/html/rfc8018).

The limit is `(2^32 - 1) * <hash length>`, after which the key would loop. When using SHA-256, exceeding the limit would take a key of more than 137GB. We do not expect anyone has ever used PBKDF2 to generate a key this large, especially because PBKDF2 runs the iterations at every block, but for the sake of correctness, we [now enforce the limit as defined by the standard](/cl/644122).

# What’s Next

The results of this audit validate the effort the Go team has put into developing high-quality, easy to use cryptography libraries and should provide confidence to our users who rely on them to build safe and secure software.

We’re not resting on our laurels, though: the Go contributors are continuing to develop and improve the cryptography libraries we provide users.

Go 1.24 now includes a FIPS 140-3 mode written in pure Go, which is currently undergoing CMVP testing. This will provide a supported FIPS 140-3 compliant mode for all users of Go, replacing the currently unsupported Go+BoringCrypto integration.

We are also working to implement modern post-quantum cryptography, introducing a ML-KEM-768 and ML-KEM-1024 implementation in Go 1.24 in the [crypto/mlkem package](/pkg/crypto/mlkem), and adding support to the crypto/tls package for the hybrid X25519MLKEM768 key exchange.

Finally, we are planning on introducing new easier to use high-level cryptography APIs, designed to reduce the barrier for picking and using high-quality algorithms for basic use cases. We plan to begin with offering a simple password hashing API that removes the need for users to decide which of the myriad of possible algorithms they should be relying on, with mechanisms to automatically migrate to newer algorithms as the state-of-the-art changes.
