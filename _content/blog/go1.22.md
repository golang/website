---
title: Go 1.22 is released!
date: 2024-02-06
by:
- Eli Bendersky, on behalf of the Go team
summary: Go 1.22 enhances for loops, brings new standard library functionality and improves performance.
---

Today the Go team is thrilled to release Go 1.22,
which you can get by visiting the [download page](/dl/).

Go 1.22 comes with several important new features and improvements. Here are
some of the notable changes; for the full list, refer to the [release
notes](/doc/go1.22).

## Language changes

The long-standing "for" loop gotcha with accidental sharing of loop variables
between iterations is now resolved. Starting with Go 1.22, the following code
will print "a", "b", and "c" in some order:

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

For more information about this change and the tooling that helps keep code from
breaking accidentally, see the earlier [loop variable blog
post](/blog/loopvar-preview).

The second language change is support for ranging over integers:

{{raw `
	package main

	import "fmt"

	func main() {
		for i := range 10 {
			fmt.Println(10 - i)
		}
		fmt.Println("go1.22 has lift-off!")
	}
`}}

The values of `i` in this countdown program go from 0 to 9, inclusive. For more
details, please refer to [the spec](/ref/spec#For_range).

## Improved performance

Memory optimization in the Go runtime improves CPU performance by 1-3%, while
also reducing the memory overhead of most Go programs by around 1%.

In Go 1.21, [we shipped](/blog/pgo) profile-guided optimization (PGO) for the Go
compiler and this functionality continues to improve. One of the optimizations
added in 1.22 is improved devirtualization, allowing static dispatch of more
interface method calls. Most programs will see improvements between 2-14% with
PGO enabled.

## Standard library additions

- A new [math/rand/v2](/pkg/math/rand/v2) package
  provides a cleaner, more consistent API and uses higher-quality,
  faster pseudo-random generation algorithms. See
  [the proposal](/issue/61716) for additional details.
- The patterns used by [net/http.ServeMux](/pkg/net/http#ServeMux)
  now accept methods and wildcards.

  For example, the router accepts a pattern like `GET /task/{id}/`, which
  matches only `GET` requests and captures the value of the `{id}` segment
  in a map that can be accessed through [Request](/pkg/net/http#Request) values.
- A new `Null[T]` type in [database/sql](/pkg/database/sql) provides
  a way to scan nullable columns.
- A `Concat` function was added in package [slices](/pkg/slices), to
  concatenate multiple slices of any type.

---

Thanks to everyone who contributed to this release by writing code and
documentation, filing bugs, sharing feedback, and testing the release
candidates. Your efforts helped to ensure that Go 1.22 is as stable as possible.
As always, if you notice any problems, please [file an issue](/issue/new).

Enjoy Go 1.22!
