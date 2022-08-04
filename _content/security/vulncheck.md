---
title: Vulnerability Detection For Go
layout: article
---

## Overview

Writing secure and reliable software requires knowing about vulnerabilities in
your dependencies. This page provides an overview of the Go vulnerability
detection package,
[golang.org/x/vuln/vulncheck](https://pkg.go.dev/golang.org/x/vuln/vulncheck),
which enables Go developers to scan dependencies in their Go projects for
public vulnerabilities.

This package is also available as a CLI tool,
[govulncheck](https://pkg.go.dev/golang.org/x/vuln/cmd/govulncheck).

## Package vulncheck

Package vulncheck delivers support for detection and understanding of how user
programs exercise known vulnerabilities. It can detect packages and modules
with known vulnerabilities that are transitively imported by the user program.
What makes vulncheck unique is that it also tries to find runtime call stacks,
sequences of currently active function calls made by the program, demonstrating
how user code reaches a vulnerability without actually executing the code. This
feature brings several benefits.

### Vulnerabilities at the package and function levels

vulncheck is more accurate than standard package-level vulnerability detection:
just because a vulnerable _symbol_ (function or method) is imported, this does
not mean the symbol is actually used. Consider the following illustrative
code.

{{raw `
    package main

    import "some/third/party/pkg/p"

    func main() {
        p.G()
    }

    // package some/third/party/pkg/p
    package p

    // F has a known vulnerability
    func F() { … }

    // G is safe and does not transitively invoke F
    func G() { … }
`}}

Package-level vulnerability detection would issue a warning for the above code
saying that the package `p` has some vulnerabilities and that it has been
transitively imported by the user package `main`. This warning is useful, but
it does not tell the whole story.  Since `F` never gets called, the above user
code is not affected by `p`'s vulnerabilities. Programmers might choose to not
address the import of `p` at all, or postpone the fix until the next release,
if they knew that no vulnerabilities of `p` are in fact exercised.
Package-level detection is inherently limited in providing the programmers with
such knowledge of vulnerabilities in their code.

### Understanding Vulnerabilities

To drive home this point, let us assume we also have an accompanying test code.

{{raw `
    package main

    import (
       "testing"
       "some/third/party/pkg/p"
    )

    func TestFoo(t *testing.T) {
       p.F()
       …
    }

    func TestBar(t *testing.T) {
       p.F()
       …
    }
`}}

Here, the vulnerable symbol `F` is indeed used by the code, but only in tests.
Programmers might be fine with vulnerabilities being potentially triggered in
tests or, say, sandboxed environments. Package and module level detection do
not provide programmers with the level of detail necessary to make such
decisions in an informed manner. vulncheck, on the other hand, is designed
precisely for that. vulncheck reports call stacks demonstrating how the
vulnerable symbols are reachable by user code. For the above example, vulncheck
communicates `[TestFoo, p.F]` and `[TestBar, p.F]` call stacks to programmers.
If we exclude above tests, vulncheck does not report any call stacks, which the
programmers can interpret as no vulnerabilities are in fact reachable.

As shown by the above example, vulncheck's motivation for reporting call stacks
to programmers goes beyond just improving the precision of vulnerability
detection. Its overarching goal is to help programmers understand
vulnerabilities, their impact, and their potential remedies. Although the
simplest fix often is just to update the corresponding vulnerable package to
its healthy version, if any, there are still some very important open
questions:

- Could my systems have been breached and, if so, where does the breach occur?
- Do I need to escalate the issue?
- Do I need to alarm my customers?

Call stacks reported by vulncheck can help programmers answer those questions,
because vulnerabilities can be buried deep in unfamiliar places in the code.
Package and module level detection on its own is very often not helpful when
addressing these questions.

## Vulnerability Graphs

The main output of vulncheck are subgraphs of the program call graph, package
import graph, and module require graph that lead to vulnerabilities. We refer
to such subgraphs as _vulnerability graphs_. A vulnerability call graph
contains only nodes and edges of the original call graph that show how
vulnerable symbols are reachable from the program entry points. At the call
graph level, entry points are `main`s, `init`s, as well as exported functions
and methods of user packages. Consider the following example:

{{raw `
    package p

    import "some/package/q"

    type X struct { ... }

    func (x X) Foo() { ... } // makes no further calls

    func A(x X) {
       q.D(x)
       q.E(x)
    }

    func B(x X) {
       q.E(x)
    }

    func C(x X) {
       x.Foo()
    }


    // package some/package/q
    package q

    import "vulnerable/package/vuln"

    type I interface {
       Foo()
    }

    type Y struct { ... }

    func (y Y) Foo() {
    vuln.V()
    }

    func D(i I) {
       i.Foo()
       y := Y{...}
       y.Foo()
    }

    func E(i I) {
       i.Foo()
       D(i)
    }

    // package vulnerable/package/vuln
    package vuln

    func V() {...} // known to be vulnerable, makes no further calls
`}}

vulncheck's `Source` function takes this program as input and first constructs
its call graph, shown below. We omit package information of each function for
brevity.

{{raw `
        A _   B    C
        |  \  |    |
        |   \ |    |
        D <-- E    |
        | \   \    |
        |  \   \   |
        |   \   \  |
      Y.Foo  -> X.Foo
        |
        V
`}}

The entry points are functions `A`, `B`, and `C` of the input package `p`.
These functions mainly pass `X` values to exported functions of package `q`
that in turn call `X.Foo`. `D` also calls `Y.Foo`.

The call to `Y.Foo` is problematic as it itself makes a call to the vulnerable
function `V` of `vuln`. We thus have a call to a vulnerable function in a
dependent package that is not under control of the author of the package `p`.
This can be hard to trace down for programmers by relying on just package-level
vulnerability detection. vulncheck detects this and computes the following
vulnerability call graph.

{{raw `
        A _   B
        |  \  |
        |   \ |
        D <-- E
        |
        Y.Foo
        |
        V
`}}

Functions `C` and `X.Foo` are not in the vulnerability graph as they do not
transitively lead to `V`. In general, all edges not leading to vulnerable
symbols are omitted, as well as nodes appearing exclusively along those edges.
The same principles are used to create vulnerability graphs of package imports
and module require graphs.

### Evidence of vulnerability uses

Clients of vulncheck can present the vulnerability graphs, such as the one
above, to the programmers as a way of showing how vulnerabilities are reachable
in their code. However, vulnerability graphs can get large for big projects,
which would make it hard for programmers to manually inspect vulnerabilities.
In response, vulncheck also provides `CallStacks` functionality for extracting
call stacks from vulnerability call graphs.

For each pair of a vulnerable symbol and an entry point, `CallStacks` traverses
the vulnerability graph searching for call stacks starting at the entry point
and ending with a call to the vulnerable symbol. To avoid exponential
explosion, each node is visited at most once. The extracted stacks for a
particular vulnerability are heuristically ordered by how easy is to understand
them: shorter call stacks with less dynamic call sites appear earlier in the
extracted results. For the vulnerability call graph shown earlier, there are
two stacks reported for `V`.

{{raw `
        A      B
        |      |
        D      E
        |      |
        Y.Foo   D
        |      |
        V     Y.Foo
               |
               V
`}}

Note that the call stack `[A, E, D, Y.Foo, V]` is not reported since a shorter
extracted stack `[A, D, Y.Foo, V]` starting at `A` already goes through `D`.The
clients of vulncheck can present (a subset) of representative calls stacks to
programmers as a more succinct evidence of vulnerability uses. For instance,
[govulncheck](https://pkg.go.dev/golang.org/x/vuln/cmd/govulncheck) by default
shows only the first call stack extracted by `CallStacks`.

Few notes. vulncheck can also analyze Go binaries with some limitations (see
Limitations section). Vulnerabilities are modeled using the shared
[golang.org/x/vuln/osv](https://golang.org/x/vuln/osv) format and an existing
vulnerability database is available at
[https://vuln.go.dev](https://vuln.go.dev). For more details on vulncheck data
structures and APIs, please see
[here](https://pkg.go.dev/golang.org/x/vuln/vulncheck).

## Call Graph Construction

One of the main technical challenges in vulncheck is to statically compute call
graph information of a Go program. As this is an undecidable problem, we can
only hope for an approximate solution. More precise call graph algorithms will
require more execution time. On the other hand, a really fast algorithm could
easily be very imprecise, either missing call stacks or often reporting ones
that do not appear at runtime. vulncheck strikes the balance between precision
and volume of used computational resources with the
[Variable Type Analysis](https://dl.acm.org/doi/pdf/10.1145/354222.353189)
(VTA) algorithm.

### Variable type analysis

VTA is an over-approximate call graph algorithm. VTA does not miss a call stack
realizable in practice (see Limitations section for exceptions to this), but it
might sometimes report a call stack leading to a vulnerability that cannot be
exercised in practice. Our experiments suggest this does not happen too often.

Consider again the program from the previous section. Existing algorithms, such
as [CHA](https://pkg.go.dev/golang.org/x/tools/go/callgraph/cha) or
[RTA](https://pkg.go.dev/golang.org/x/tools/go/callgraph/rta), would say that
`i.Foo()` call in `E` resolves to `X.Foo` and `Y.Foo` because types `X` and `Y`
implement interface `I` and are used in the program. If vulncheck relied on
these two algorithms, it would report vulnerable call stack `[B, E, Y.Foo, V]`
that is in fact not realizable in practice. VTA, as we hinted earlier,
correctly resolves that call to only `X.Foo` that does not lead to `V`.

VTA works on an abstract representation of a program where variables are
represented by their types. The types are then propagated around the program
based on variable usage. For the running example, parameter `x` of `B` is
abstracted via type `X` which is then propagated to parameter `i` of `E`. The
values actually stored to the variable are not taken into account, only their
types. This can lead to imprecision when types reaching an interface variable
depend on valuation of, say, involved conditional statements or complicated
aliasing. However, types of concrete variables are always the same, regardless
of the complexity of the surrounding logic. For instance, values reaching a
variable of type `X` always have precisely the type `X`. This rather unique
property of Go's type system enables VTA to produce precise call graph
information.

### Achieving scale

In the current example, VTA propagates type `X` from `B` to `E` because the
call `E(i)` is static. VTA knows what function the identifier `E` resolves to.
But what if that call was dynamic? After all, VTA is supposed to construct the
call graph so how can it then propagate types across function boundaries? One
solution is to rely on a fix-point where the results of type propagation are
also used to establish function call edges on which type propagation then needs
to be repeated, and so on. This could be very expensive, so VTA relies on an
initial approximation of the call graph to scale. Note that the initial call
graph is only used to propagate types over function calls. We choose CHA as the
initial call graph. As CHA can be rather imprecise, as shown on the earlier
example, it could cause VTA to be overly imprecise as well. To counter that,
vulncheck bootstraps VTA by VTA. After computing VTA on top of CHA, we feed the
more precise resulting call graph to VTA again, toning down excessive
imprecision initially introduced by CHA.

Package VTA can be found at
[golang.org/x/tools/go/callgraph/vta](https://pkg.go.dev/golang.org/x/tools/go/callgraph/vta).

## Limitations

As VTA can produce call stacks that are not realizable in practice, vulncheck
can claim that a vulnerable symbol is reachable while in fact it is not. We
also note that VTA might miss some call stacks that go through _unsafe_ and
_reflect_ packages.

Because binaries do not contain detailed call information, vulncheck cannot
compute vulnerability call graphs and call stack witnesses for Go binaries.

vulncheck currently does not detect vulnerable packages and symbols that have
been vendored rather than imported. Also, there is currently no support for
silencing vulnerability findings. If you are interested in any of these features,
please let us know by [filing an issue](https://golang.org/s/govulncheck-feedback).
