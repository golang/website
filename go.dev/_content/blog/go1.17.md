---
title: Go 1.17 is released
date: 2021-08-16
by:
- Matt Pearring
- Alex Rakoczy
summary: Go 1.17 adds performance improvements, module optimizations, arm64 on Windows, and more.
---


Today the Go team is thrilled to release Go 1.17, which you can get by visiting the
[download page](/dl/).

This release brings additional improvements to the compiler, namely a
[new way of passing function arguments and results](/doc/go1.17#compiler).
This change has shown about a 5% performance improvement in Go programs and reduction in binary
sizes of around 2% for amd64 platforms. Support for more platforms will come in future releases.

Go 1.17 also adds support for the
[64-bit ARM architecture on Windows](/doc/go1.17#ports), letting gophers run
Go natively on more devices.

We’ve also introduced [pruned module graphs](/doc/go1.17#go-command) in this
release. Modules that specify `go 1.17` or higher in their `go.mod` file will have their module graphs
include only the immediate dependencies of other Go 1.17 modules, not their full transitive
dependencies. This should help avoid the need to download or read `go.mod` files for otherwise
irrelevant dependencies—saving time in everyday development.

Go 1.17 comes with three small [changes to the language](/doc/go1.17#language).
The first two are new functions in the `unsafe` package to make it simpler for programs to conform
to the `unsafe.Pointer` rules: `unsafe.Add` allows for
[safer pointer arithmetic](/pkg/unsafe#Add), while `unsafe.Slice` allows for
[safer conversions of pointers to slices](/pkg/unsafe#Slice). The third change is
an extension to the language type conversion rules to allow conversions from
[slices to array pointers](/ref/spec#Conversions_from_slice_to_array_pointer),
provided the slice is at least as large as the array at runtime.

Finally there are quite a few other improvements and bug fixes, including verification improvements
to [crypto/x509](/doc/go1.17#crypto/x509), and alterations to
[URL query parsing](/doc/go1.17#semicolons). For a complete list of changes and
more information about the improvements above, see the
[full release notes](/doc/go1.17).

Thanks to everyone who contributed to this release by writing code, filing bugs, sharing feedback,
and testing the beta and release candidates. Your efforts helped to ensure that Go 1.17 is as stable
 as possible. As always, if you notice any problems, please
[file an issue](/issue/new).

We hope you enjoy the new release!

