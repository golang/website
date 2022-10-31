---
title: "Automatic cipher suite ordering in crypto/tls"
date: 2021-09-15
by:
- Filippo Valsorda
summary: Go 1.17 is making TLS configuration easier and safer by automating TLS cipher suite preference ordering.
---

The Go standard library provides `crypto/tls`,
a robust implementation of Transport Layer Security (TLS),
the most important security protocol on the Internet,
and the fundamental component of HTTPS.
In Go 1.17 we made its configuration easier, more secure,
and more efficient by automating the priority order of cipher suites.


## How cipher suites work

Cipher suites date back to TLS’s predecessor Secure Socket Layer (SSL),
which [called them “cipher kinds”](https://datatracker.ietf.org/doc/html/draft-hickman-netscape-ssl-00#appendix-C.4).
They are the intimidating-looking identifiers like
`TLS_RSA_WITH_AES_256_CBC_SHA` and
`TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305_SHA256`
that spell out the algorithms used to exchange keys,
authenticate certificates, and encrypt records in a TLS connection.

Cipher suites are _negotiated_ during the TLS handshake:
the client sends the list of cipher suites it supports in its first message,
the Client Hello, and the server picks one from that list,
communicating its choice to the client.
The client sends the list of supported cipher suites in its own preference order,
and the server is free to pick from it however it wants.
Most commonly, the server will pick the first mutually supported cipher
suite either in client preference order or in server preference order,
based on its configuration.

Cipher suites are really only one of many negotiated parameters—supported
curves/groups and signature algorithms are additionally negotiated through
their own extensions—but are the most complex and famous ones,
and the only ones that developers and administrators were trained over the
years to have opinions on.

In TLS 1.0–1.2, all these parameters interact in a complex web of inter-dependencies:
for example supported certificates depend on supported signature algorithms,
supported curves, and supported cipher suites.
In TLS 1.3 this was all drastically simplified:
cipher suites only specify symmetric encryption algorithms,
while supported curves/groups govern the key exchange and supported signature
algorithms apply to the certificate.


## A complex choice abdicated to developers

Most HTTPS and TLS servers delegate the choice of cipher suites and preference
order to the server operator or the applications developer.
This is a complex choice that requires up-to-date and specialized knowledge for many reasons.

Some older cipher suites have insecure components,
some require extremely careful and complex implementations to be secure,
and some are only secure if the client applies certain mitigations or even
has certain hardware.
Beyond the security of the individual components,
different cipher suites can provide drastically different security properties
for the whole connection,
as cipher suites without ECDHE or DHE don’t provide forward secrecy—the
property that connections can’t be retroactively or passively decrypted
with the certificate’s key.
Finally, the choice of supported cipher suites impacts compatibility and performance,
and making changes without an up-to-date understanding of the ecosystem
can lead to breaking connections from legacy clients,
increasing the resources consumed by the server,
or draining the batteries of mobile clients.

This choice is so arcane and delicate that there are dedicated tools to guide operators,
such as the excellent [Mozilla SSL Configuration Generator](https://ssl-config.mozilla.org/).

How did we get here and why is it like this?

To start, individual cryptographic components used to break much more often.
In 2011, when the BEAST attack broke CBC cipher suites in such a way that
only clients could mitigate the attack,
servers moved to preferring RC4, which was unaffected.
In 2013, when it became clear that RC4 was broken,
servers went back to CBC.
When Lucky Thirteen made it clear it was extremely hard to implement CBC
cipher suites due to their backwards MAC-then-encrypt design...
Well, there wasn’t anything else on the table so implementations had to
[carefully jump through hoops](https://www.imperialviolet.org/2013/02/04/luckythirteen.html)
to implement CBC and kept [failing at that daunting task for years](https://blog.cloudflare.com/yet-another-padding-oracle-in-openssl-cbc-ciphersuites/).
Configurable cipher suites and [cryptographic agility](https://www.imperialviolet.org/2016/05/16/agility.html)
used to provide some reassurance that when a component broke it could be
replaced on the fly.

Modern cryptography is significantly different.
Protocols can still break from time to time,
but it’s rarely an individual abstracted component that fails.
_None of the AEAD-based cipher suites introduced starting with TLS 1.2 in
2008 have been broken._ These days cryptographic agility is a liability:
it introduces complexity that can lead to weaknesses or downgrades,
and it is only necessary for performance and compliance reasons.

Patching used to be different, too. Today we acknowledge that promptly applying
software patches for disclosed vulnerabilities is the cornerstone of secure
software deployments,
but ten years ago it was not standard practice.
Changing configuration was seen as a much more rapid option to respond to
vulnerable cipher suites,
so the operator, through configuration, was put fully in charge of them.
We now have the opposite problem: there are fully patched and updated servers
that still behave weirdly,
suboptimally, or insecurely, because their configurations haven't been touched in years.

Finally, it was understood that servers tended to update more slowly than clients,
and therefore were less reliable judges of the best choice of cipher suite.
However, it’s servers who have the last word on cipher suite selection,
so the default became to make servers defer to the client preference order,
instead of having strong opinions.
This is still partially true: browsers managed to make automatic updates
happen and are much more up-to-date than the average server.
On the other hand, a number of legacy devices are now out of support and
are stuck on old TLS client configurations,
which often makes an up-to-date server better equipped to choose than some of its clients.

Regardless of how we got here, it’s a failure of cryptography engineering
to require application developers and server operators to become experts
in the nuances of cipher suite selection,
and to stay up-to-date on the latest developments to keep their configs up-to-date.
If they are deploying our security patches,
that should be enough.

The Mozilla SSL Configuration Generator is great, and it should not exist.

Is this getting any better?

There is good news and bad news for how things are trending in the past few years.
The bad news is that ordering is getting even more nuanced,
because there are sets of cipher suites that have equivalent security properties.
The best choice within such a set depends on the available hardware and
is hard to express in a config file.
In other systems, what started as a simple list of cipher suites now depends
on [more complex syntax](https://boringssl.googlesource.com/boringssl/+/c3b373bf4f4b2e2fba2578d1d5b5fe04e410f7cb/include/openssl/ssl.h#1457)
or additional flags like [SSL\_OP\_PRIORITIZE\_CHACHA](https://www.openssl.org/docs/man1.1.1/man3/SSL_CTX_clear_options.html#:~:text=session-,ssl_op_prioritize_chacha,-When).

The good news is that TLS 1.3 drastically simplified cipher suites,
and it uses a disjoint set from TLS 1.0–1.2.
All TLS 1.3 cipher suites are secure, so application developers and server
operators shouldn’t have to worry about them at all.
Indeed, some TLS libraries like BoringSSL and Go’s `crypto/tls` don’t
allow configuring them at all.


## Go’s crypto/tls and cipher suites

Go does allow configuring cipher suites in TLS 1.0–1.2.
Applications have always been able to set the enabled cipher suites and
preference order with [`Config.CipherSuites`](https://pkg.go.dev/crypto/tls#Config.CipherSuites).
Servers prioritize the client’s preference order by default,
unless [`Config.PreferServerCipherSuites`](https://pkg.go.dev/crypto/tls#Config.PreferServerCipherSuites) is set.

When we implemented TLS 1.3 in Go 1.12, [we didn’t make TLS 1.3 cipher suites configurable](/issue/29349),
because they are a disjoint set from the TLS 1.0–1.2 ones and most importantly
they are all secure,
so there is no need to delegate a choice to the application.
`Config.PreferServerCipherSuites` still controls which side’s preference order is used,
and the local side’s preferences depend on the available hardware.

In Go 1.14 we [exposed supported cipher suites](https://pkg.go.dev/crypto/tls#CipherSuites),
but explicitly chose to return them in a neutral order (sorted by their ID),
so that we wouldn’t end up tied to representing our priority logic in
terms of a static sort order.

In Go 1.16, we started actively [preferring ChaCha20Poly1305 cipher suites over AES-GCM on the server](/cl/262857)
when we detect that either the client or the server lacks hardware support for AES-GCM.
This is because AES-GCM is hard to implement efficiently and securely without
dedicated hardware support (such as the AES-NI and CLMUL instruction sets).

**Go 1.17, recently released, takes over cipher suite preference ordering for all Go users.**
While `Config.CipherSuites` still controls which TLS 1.0–1.2 cipher suites are enabled,
it is not used for ordering, and `Config.PreferServerCipherSuites` is now ignored.
Instead, `crypto/tls` [makes all ordering decisions](/cl/314609),
based on the available cipher suites, the local hardware,
and the inferred remote hardware capabilities.

The [current TLS 1.0–1.2 ordering logic](https://cs.opensource.google/go/go/+/9d0819b27ca248f9949e7cf6bf7cb9fe7cf574e8:src/crypto/tls/cipher_suites.go;l=206-270)
follows the following rules:



1. ECDHE is preferred over the static RSA key exchange.

    The most important property of a cipher suite is enabling forward secrecy.
    We don’t implement “classic” finite-field Diffie-Hellman,
    because it’s complex, slower, weaker, and [subtly broken](https://datatracker.ietf.org/doc/draft-bartle-tls-deprecate-ffdh/) in TLS 1.0–1.2,
    so that means prioritizing the Elliptic Curve Diffie-Hellman key exchange
    over the legacy static RSA key exchange.
    (The latter simply encrypts the connection’s secret using the certificate’s
    public key, making it possible to decrypt if the certificate is compromised
    in the future.)

2. AEAD modes are preferred over CBC for encryption.

    Even if we do implement partial countermeasures for Lucky13
    ([my first contribution to the Go standard library, back in 2015!](/cl/18130)),
    the CBC suites are [a nightmare to get right](https://blog.cloudflare.com/yet-another-padding-oracle-in-openssl-cbc-ciphersuites/),
    so all other more important things being equal,
    we pick AES-GCM and ChaCha20Poly1305 instead.

3. 3DES, CBC-SHA256, and RC4 are only used if nothing else is available, in that preference order.

    3DES has 64-bit blocks, which makes it fundamentally vulnerable to
    [birthday attacks](https://sweet32.info) given enough traffic.
    3DES is listed under [`InsecureCipherSuites`](https://pkg.go.dev/crypto/tls#InsecureCipherSuites),
    but it’s enabled by default for compatibility.
    (An additional benefit of controlling preference orders is that
    we can afford to keep less secure cipher suites enabled by default
    without worrying about applications or clients selecting them
    except as a last resort.
    This is safe because there are no downgrade attacks that rely on
    the availability of a weaker cipher suite to attack peers
    that support better alternatives.)


    The CBC cipher suites are vulnerable to Lucky13-style side channel attacks
    and we only partially implement the [complex](https://www.imperialviolet.org/2013/02/04/luckythirteen.html)
    countermeasures discussed above for the SHA-1 hash, not for SHA-256.
    CBC-SHA1 suites have compatibility value, justifying the extra complexity,
    while the CBC-SHA256 ones don’t, so they are disabled by default.


    RC4 has [practically exploitable biases](https://www.rc4nomore.com)
    that can lead to plaintext recovery without side channels.
    It doesn’t get any worse than this, so RC4 is disabled by default.

4. ChaCha20Poly1305 is preferred over AES-GCM for encryption,
    unless both sides have hardware support for the latter.

    As we discussed above, AES-GCM is hard to implement efficiently and
    securely without hardware support.
    If we detect that there isn’t local hardware support or (on the server)
    that the client has not prioritized AES-GCM,
    we pick ChaCha20Poly1305 instead.

5. AES-128 is preferred over AES-256 for encryption.

    AES-256 has a larger key than AES-128, which is usually good,
    but it also performs more rounds of the core encryption function,
    making it slower.
    (The extra rounds in AES-256 are independent of the key size change;
    they are an attempt to provide a wider margin against cryptanalysis.)
    The larger key is only useful in multi-user and post-quantum settings,
    which are not relevant to TLS, which generates sufficiently random IVs
    and has no post-quantum key exchange support.
    Since the larger key doesn’t have any benefit,
    we prefer AES-128 for its speed.


[TLS 1.3's ordering logic](https://cs.opensource.google/go/go/+/9d0819b27ca248f9949e7cf6bf7cb9fe7cf574e8:src/crypto/tls/cipher_suites.go;l=342-355)
needs only the last two rules,
because TLS 1.3 eliminated the problematic algorithms the first three rules
are guarding against.


## FAQs

_What if a cipher suite turns out to be broken?_ Just like any other vulnerability,
it will be fixed in a security release for all supported Go versions.
All applications need to be prepared to apply security fixes to operate securely.
Historically, broken cipher suites are increasingly rare.

_Why leave enabled TLS 1.0–1.2 cipher suites configurable?_ There is a
meaningful tradeoff between _baseline_ security and legacy compatibility
to make in choosing which cipher suites to enable,
and that’s a choice we can’t make ourselves without either cutting out
an unacceptable slice of the ecosystem,
or reducing the security guarantees of modern users.

_Why not make TLS 1.3 cipher suites configurable?_ Conversely,
there is no tradeoff to make with TLS 1.3,
as all its cipher suites provide strong security.
This lets us leave them all enabled and pick the fastest based on the specifics
of the connection without requiring the developer’s involvement.


## Key takeaways

Starting in Go 1.17, `crypto/tls` is taking over the order in which available
cipher suites are selected.
With a regularly updated Go version, this is safer than letting potentially
outdated clients pick the order,
lets us optimize performance, and it lifts significant complexity from Go developers.

This is consistent with our general philosophy of making cryptographic decisions whenever we can,
instead of delegating them to developers,
and with our [cryptography principles](/design/cryptography-principles).
Hopefully other TLS libraries will adopt similar changes,
making delicate cipher suite configuration a thing of the past.
