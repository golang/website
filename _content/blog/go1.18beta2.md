---
title: Announcing Go 1.18 Beta 2
date: 2022-01-31
by:
- Jeremy Faller and Steve Francia, for the Go team
summary: Go 1.18 Beta 2 is our second preview of Go 1.18. Please try it and let us know if you find problems.
---

We are encouraged by all the excitement around Go’s upcoming 1.18 release,
which adds support for
[generics](https://go.dev/blog/why-generics),
[fuzzing](https://go.dev/blog/fuzz-beta), and the new
[Go workspace mode](https://go.dev/design/45713-workspace).

We released Go 1.18 beta 1 two months ago,
and it is now the most downloaded Go beta ever,
with twice as many downloads as any previous release.
Beta 1 has also proved very reliable;
in fact, we are already running it in production here at Google.

Your feedback on Beta 1 helped us identify obscure bugs
in the new support for generics and ensure a more stable final release.
We've resolved these issues in today's release of Go 1.18 Beta 2,
and we encourage everyone to try it out.
The easiest way to install it alongside your existing Go toolchain is to run:

	go install golang.org/dl/go1.18beta2@latest
	go1.18beta2 download

After that, you can run `go1.18beta2` as a drop-in replacement for `go`.
For more download options, visit https://go.dev/dl/#go1.18beta2.

Because we are taking the time to issue a second beta,
we now expect that the Go 1.18 release candidate will be issued in February,
with the final Go 1.18 release in March.

The Go language server `gopls` and the VS Code Go extension
now support generics.
To install `gopls` with generics, see
[this documentation](https://github.com/golang/tools/blob/master/gopls/doc/advanced.md#working-with-generic-code),
and to configure the VS Code Go extension, follow [this instruction](https://github.com/golang/vscode-go/blob/master/docs/advanced.md#using-go118).


As always, especially for beta releases,
if you notice any problems, please [file an issue](https://go.dev/issue/new).

