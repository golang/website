---
title: Fixing For Loops in Go 1.22
date: 2023-09-19
by:
- David Chase
- Russ Cox
summary: Go 1.21 shipped a preview of a change in Go 1.22 to make for loops less error-prone.
---

Go 1.21 includes a preview of a change to `for` loop scoping
that we plan to ship in Go 1.22,
removing one of the most common Go mistakes.

## The Problem

If you’ve written any amount of Go code, you’ve probably made the mistake
of keeping a reference to a loop variable past the end of its iteration,
at which point it takes on a new value that you didn’t want.
For example, consider this program:

{{raw `
	func main() {
		done := make(chan bool)

		values := []string{"a", "b", "c"}
		for _, v := range values {
			go func() {
				fmt.Println(v)
				done <- true
			}()
		}

		// wait for all goroutines to complete before exiting
		for _ = range values {
			<-done
		}
	}
`}}

The three created goroutines are all printing the same variable `v`,
so they usually print “c”, “c”, “c”, instead of printing “a”, “b”, and “c” in some order.

The [Go FAQ entry “What happens with closures running as goroutines?”](https://go.dev/doc/faq#closures_and_goroutines),
gives this example and remarks
“Some confusion may arise when using closures with concurrency.”

Although concurrency is often involved, it need not be.
This example has the same problem but no goroutines:

{{raw `
	func main() {
		var prints []func()
		for i := 1; i <= 3; i++ {
			prints = append(prints, func() { fmt.Println(i) })
		}
		for _, print := range prints {
			print()
		}
	}
`}}

This kind of mistake has caused production problems at many companies,
including a
[publicly documented issue at Lets Encrypt](https://bugzilla.mozilla.org/show_bug.cgi?id=1619047).
In that instance, the accidental capture of the loop variable was spread across
multiple functions and much more difficult to notice:

	// authz2ModelMapToPB converts a mapping of domain name to authz2Models into a
	// protobuf authorizations map
	func authz2ModelMapToPB(m map[string]authz2Model) (*sapb.Authorizations, error) {
		resp := &sapb.Authorizations{}
		for k, v := range m {
			// Make a copy of k because it will be reassigned with each loop.
			kCopy := k
			authzPB, err := modelToAuthzPB(&v)
			if err != nil {
				return nil, err
			}
			resp.Authz = append(resp.Authz, &sapb.Authorizations_MapElement{
				Domain: &kCopy,
				Authz: authzPB,
			})
		}
		return resp, nil
	}

The author of this code clearly understood the general problem, because they made a copy of `k`,
but it turns out `modelToAuthzPB` used pointers to fields in `v` when constructing its result,
so the loop also needed to make a copy of `v`.

Tools have been written to identify these mistakes, but it is hard to analyze
whether references to a variable outlive its iteration or not.
These tools must choose between false negatives and false positives.
The `loopclosure` analyzer used by `go vet` and `gopls` opts for false negatives,
only reporting when it is sure there is a problem but missing others.
Other checkers opt for false positives, accusing correct code of being incorrect.
We ran an analysis of commits adding `x := x` lines in open-source Go code,
expecting to find bug fixes.
Instead we found many unnecessary lines being added,
suggesting instead that popular checkers have significant false positive rates,
but developers add the lines anyway to keep the checkers happy.

One pair of examples we found was particularly illuminating:


This diff was in one program:

	     for _, informer := range c.informerMap {
	+        informer := informer
	         go informer.Run(stopCh)
	     }

And this diff was in another program:

	     for _, a := range alarms {
	+        a := a
	         go a.Monitor(b)
	     }

One of these two diffs is a bug fix; the other is an unnecessary change.
You can’t tell which is which unless you know more about the types
and functions involved.

## The Fix

For Go 1.22, we plan to change `for` loops to make these variables have
per-iteration scope instead of per-loop scope.
This change will fix the examples above, so that they are no longer buggy Go programs;
it will end the production problems caused by such mistakes;
and it will remove the need for imprecise tools that prompt users
to make unnecessary changes to their code.

To ensure backwards compatibility with existing code, the new semantics
will only apply in packages contained in modules that declare `go 1.22` or
later in their `go.mod` files.
This per-module decision provides developer control of a gradual update
to the new semantics throughout a codebase.
It is also possible to use `//go:build` lines to control the decision on a
per-file basis.

Old code will continue to mean exactly what it means today:
the fix only applies to new or updated code.
This will give developers control over when the semantics change
in a particular package.
As a consequence of our [forward compatibility work](toolchain),
Go 1.21 will not attempt to compile code that declares `go 1.22` or later.
We included a special case with the same effect in
the point releases Go 1.20.8 and Go 1.19.13,
so when Go 1.22 is released,
code written depending on the new semantics will never be compiled with
the old semantics, unless people are using very old, [unsupported Go versions](/doc/devel/release#policy).


## Previewing The Fix

Go 1.21 includes a preview of the scoping change.
If you compile your code with `GOEXPERIMENT=loopvar` set in your environment,
then the new semantics are applied to all loops
(ignoring the `go.mod` `go` lines).
For example, to check whether your tests still pass with the new loop semantics
applied to your package and all your dependencies:

	GOEXPERIMENT=loopvar go test

We patched our internal Go toolchain at Google to force this mode during all builds
at the start of May 2023, and in the past four months
we have had zero reports of any problems in production code.

You can also try test programs to better understand the semantics
on the Go playground by including a `// GOEXPERIMENT=loopvar` comment
at the top of the program, like in [this program](https://go.dev/play/p/YchKkkA1ETH).
(This comment only applies in the Go playground.)

## Fixing Buggy Tests

Although we’ve had no production problems,
to prepare for that switch, we did have to correct many buggy tests that were not
testing what they thought they were, like this:

	func TestAllEvenBuggy(t *testing.T) {
		testCases := []int{1, 2, 4, 6}
		for _, v := range testCases {
			t.Run("sub", func(t *testing.T) {
				t.Parallel()
				if v&1 != 0 {
					t.Fatal("odd v", v)
				}
			})
		}
	}

In Go 1.21, this test passes because `t.Parallel` blocks each subtest
until the entire loop has finished and then runs all the subtests
in parallel. When the loop has finished, `v` is always 6,
so the subtests all check that 6 is even,
so the test passes.
Of course, this test really should fail, because 1 is not even.
Fixing for loops exposes this kind of buggy test.

To help prepare for this kind of discovery, we improved the precision
of the `loopclosure` analyzer in Go 1.21 so that it can identify and
report this problem.
You can see the report [in this program](https://go.dev/play/p/WkJkgXRXg0m) on the Go playground.
If `go vet` is reporting this kind of problem in your own tests,
fixing them will prepare you better for Go 1.22.

If you run into other problems,
[the FAQ](https://github.com/golang/go/wiki/LoopvarExperiment#my-test-fails-with-the-change-how-can-i-debug-it)
has links to examples and details about using a tool we’ve written to identify
which specific loop is causing a test failure when the new semantics are applied.

## More Information

For more information about the change, see the
[design document](https://go.googlesource.com/proposal/+/master/design/60078-loopvar.md)
and the
[FAQ](https://go.dev/wiki/LoopvarExperiment).

