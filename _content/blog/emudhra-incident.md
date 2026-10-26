---
title: Our Response To The eMudhra Key Disclosure Incident
date: 2026-10-26
by:
- Roland Shoemaker, on behalf of the Go Security team
tags:
- security
- release
- technical
summary: How the Go team responded to the eMudhra certificate compromise with our first-ever URGENT track security release.
---

On October 19, the Go team initiated an expedited [security release](https://groups.google.com/g/golang-announce/c/66pojueaQek) to protect Go users and their applications from a critical public key infrastructure (PKI) compromise involving certificates issued by eMudhra.

Within approximately seven and a half hours of being notified, we shipped a targeted patch, which was included in our first-ever [URGENT track](/doc/security/policy#urgent) security release.

Here is the context behind the incident, why Go stepped in at the runtime level to protect our users and how our incident response unfolded.

## eMudhra Incident

On October 17, Mozilla added eight intermediate certificates (six used for TLS, and two used for code signing) issued by a root certificate controlled by eMudhra to OneCRL, effectively blocking their usage in Firefox. Soon after, a number of security researchers noticed this and began investigating what might’ve caused this.

On October 19, the Go Security team was made aware of these investigations and their conclusions, which showed that the keys associated with these intermediates were trivially weak, allowing recovery of the private keys in microseconds.

Due to the fact that the eMudhra root certificate was present in most trusted root programs, recovering the private keys would allow any attacker to sign arbitrary TLS certificates for any domain name, without the usual requirement to show proof of control of said domain name. This would allow an attacker to man-in-the-middle (MITM) TLS connections and present a valid certificate for a domain name they did not control, effectively allowing them to impersonate any TLS server the targeted machine attempted to connect to.

Go ships a TLS certificate verifier in the standard library. When users use [`net/http`](https://pkg.go.dev/net/http) or [`crypto/tls`](https://pkg.go.dev/crypto/tls) to connect to a TLS service, the verifier is called to check if the certificates the other side of the connection presents are valid. On Windows and macOS we delegate these checks to the platform APIs, and on Linux, which lacks platform APIs, we use our own verifier implementation.

Since Linux does not have a uniform built-in mechanism to handle revocation checking, and in order to provide a timely fix for users on other platforms that might not receive platform updates, we decided the best way to protect Go users was to do an [URGENT track](/doc/security/policy#urgent) security release which similarly removed trust.

## Incident Timeline

A typical Go security release follows a multi-day process, with advanced pre-announcements in order to allow those who rely on Go to prepare. In this case, because we believed that private key exposure was imminent, our team chose to use the URGENT track for the first time in order to ship protections before active public exploitation could begin.

**11:45 PDT**

Our Security team was alerted by an independent external researcher, citing unreported actions by Mozilla and other organizations to remove trust in eight eMudhra intermediates. It became clear quickly that this was due to weak private keys they used.

**12:30 PDT**

Mozilla independently reported the issue to our Security team, confirming the issue, and asking that we not immediately disclose the issue, but encouraging us to release a fix as soon as possible.

**14:00 PDT**

We finalized a patch which blocked the usage of certificates which chained to the compromised intermediates, and thanks to the researcher who initially alerted us to the issue, had a set of certificates signed using the recovered keys to test against.

**15:00 PDT**

We began the process to release minor Go versions 1.27.3 and 1.26.10 containing our patch. Due to the urgent nature of this issue and its exploitability, and with the belief that the private keys would likely become publicly known within hours, we decided an immediate release was the most prudent option.

**19:13 PDT**

We finalized the release and announced it publicly. At this time the underlying issue was still not widely known. In an effort to give other defenders a chance to protect as many users as possible without revealing to attackers how the issue could be exploited, our patch and communications were somewhat opaque, warning users that we were “blocking trust in a small set of intermediate certificates due to specific security concerns around their operation” and urging users to immediately upgrade to this release.

In under seven and a half hours from initial notification, patched Go releases were live worldwide, giving users everything they needed to secure their deployments before the keys became public.

By the morning of October 20, numerous independent researchers had publicly disclosed the weak keys.

## Lessons Learned and Future Improvements

This was our first ever URGENT track security release, and during the release process we ran into a number of issues which caused the release to take longer than expected.

While we are happy we were able to get a release out the same day we were notified of the issue, in the modern security environment we feel that we should be able to provide releases in a more compressed timeline. The issues we ran into caused our release process to take just under four and a quarter hours from start to finish.

These issues had previously not been problematic, since the typical release process is split up over a number of days.

We’ve subsequently made a number of changes to improve our release efficiency, which should significantly reduce the time it takes us to push out a release, enhancing our ability to quickly respond to urgent security threats.

PRIVATE track security fixes are pre-announced three to seven days before their release. URGENT track fixes happen on a shorter time schedule, but we could have pre-announced the upcoming release when we began the release process, giving users at least a few hours to prepare.

While on Windows and macOS we rely on the revocation mechanisms the platform APIs provide, on Linux we don’t currently have a similar mechanism. The lack of this support was one of our major drivers to issue an expedited security release, as we did not want to leave Linux users exposed to this issue longer than was necessary.

Hence we are planning on prioritizing work to bring parity to revocation support on Linux, to help mitigate this class of issues in the future. Part of this work will involve investigating adding support for the [upki](https://github.com/rustls/upki) project, which aims to “provide reliable, privacy-preserving and efficient certificate revocation” for Linux (see [issue #82402](/issue/82402)).

## Upgrading and Verifying Your Deployments

If you run Go binaries which interact with network services using TLS (such as directly via `net/http` or `crypto/tls` or via a dependency such as gRPC), ensure you are running on the latest minor releases. Binaries rebuilt with Go 1.27.3 and 1.26.10 will automatically reject certificates issued by the compromised eMudhra intermediates.

You can verify your current toolchain version using `go version` and you can check the version of Go a binary was built with using `go version -m <path_to_your_binary>`.

The [govulncheck tool](https://pkg.go.dev/golang.org/x/vuln/cmd/govulncheck) can also be used to determine if a binary is affected by this issue:

```
$ GOTOOLCHAIN=go1.27.3+auto go install golang.org/x/vuln/cmd/govulncheck@latest
$ govulncheck -mode=binary <path_to_your_binary>
```

## Securing Go Deployments Everywhere

For us, securing Go and its users extends beyond fixing direct security vulnerabilities in the runtime and libraries we ship. In order to keep the ecosystem secure, we believe we also need to provide users with the right tools to write secure software, look towards the future of cryptographic algorithms and protocols, and proactively respond to ecosystem wide issues when we can.

We are grateful to the independent security researchers and the Mozilla security team for this rapid notification and coordination during this incident.
