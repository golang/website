---
title: Evolving the Go Standard Library with math/rand/v2
date: 2024-05-01
by:
- Russ Cox
summary: Go 1.22 adds math/rand/v2 and charts a course for the evolution of the Go standard library.
---

Since Go 1 was [released in March 2012](/blog/go1),
changes to the standard library have been
constrained by Go's [compatibility promise](/doc/go1compat).
Overall, compatibility has been a boon for Go users,
providing a stable base for production systems,
documentation, tutorials, books, and more.
Over time, however, we've realized mistakes in the original APIs
that cannot be fixed compatibly; in other cases,
best practices and convention have changed.
We need a plan for making important, breaking changes too.

This blog post is about Go 1.22's new [`math/rand/v2`](/pkg/math/rand/v2/) package,
the first “v2” in the standard library.
It brings needed improvements to the [`math/rand`](/pkg/math/rand/) API,
but more importantly it sets an example for how we can
revise other standard library packages as the need arises.

(In Go, `math/rand` and `math/rand/v2` are two different packages
with different import paths.
Go 1 and every release after it have included `math/rand`; Go 1.22 added `math/rand/v2`.
A Go program can import either package, or both.)

This post discusses the specific rationale for the changes in `math/rand/v2`
and then [reflects on the general principles](#principles) that will guide
new versions of other packages.

## Pseudorandom Number Generators {#pseudo}

Before we look at `math/rand`, which is an API for a pseudorandom number generator,
let's take a moment to understand what that means.

A pseudorandom number generator is a deterministic program
that generates a long sequence of
seemingly random numbers from a small seed input,
although the numbers are not in fact random at all.
In the case of `math/rand`, the seed is a single int64,
and the algorithm produces a sequence of int64s
using a variant of a
[linear-feedback shift register (LFSR)](https://en.wikipedia.org/wiki/Linear-feedback_shift_register).
The algorithm is based on an idea by George Marsaglia,
tweaked by Don Mitchell and Jim Reeds,
and further customized by Ken Thompson for Plan 9 and then Go.
It has no official name, so this post calls it the Go 1 generator.

The goal is for these generators to be fast,
repeatable, and random enough to support simulations,
shuffling, and other non-cryptographic use cases.
Repeatability is particularly important for uses like
numerical simulations or randomized testing.
For example, a randomized tester might pick a seed
(perhaps based on the current time), generate
a large random test input, and repeat.
When the tester finds a failure, it only needs to print the seed
to allow repeating the test with that specific large input.

Repeatability also matters over time: given a particular
seed, a new version of Go needs to generate the same
sequence of values that an older version did.
We didn't realize this when we released Go 1;
instead, we discovered it the hard way,
when we tried to make a change in Go 1.2
and got reports that we had broken certain tests
and other use cases.
At that point, we decided Go 1 compatibility included
the specific random outputs for a given seed
and [added a test](/change/5aca0514941ce7dd0f3cea8d8ffe627dbcd542ca).

It is not a goal for these kinds of generators to produce
random numbers suitable for deriving cryptographic keys
or other important secrets.
Because the seed is only 63 bits,
any output drawn from the generator, no matter how long,
will also only contain 63 bits of entropy.
For example, using `math/rand` to generate a 128-bit or 256-bit AES key
would be a serious mistake,
since the key would be easier to brute force.
For that kind of use, you need a cryptographically strong
random number generator, as provided by [`crypto/rand`](/pkg/crypto/rand/).

That's enough background that we can move on to what needed
fixing in the `math/rand` package.

## Problems with `math/rand` {#problems}

Over time, we noticed more and more problems with `math/rand`.
The most serious were the following.

### Generator Algorithm {#problem.generator}

The generator itself needed replacement.

The initial implementation of Go, while production ready, was in many ways a “pencil sketch”
of the entire system, working well enough to serve as a base for future development:
the compiler and runtime were written in C; the garbage collector was a conservative, single-threaded,
stop-the-world collector; and the libraries used basic implementations throughout.
From Go 1 through around Go 1.5, we went back and drew the “fully inked”
version of each of these: we converted the compiler and runtime to Go; we wrote a new, precise, parallel,
concurrent garbage collection with microsecond pause times; and we replaced
standard library implementations with more sophisticated, optimized algorithms
as needed.

Unfortunately, the repeatability requirement in `math/rand`
meant that we couldn't replace the generator there without
breaking compatibility.
We were stuck with the Go 1 generator,
which is reasonably fast (about 1.8ns per number on my M3 Mac)
but maintains an internal state of almost 5 kilobytes.
In contrast, Melissa O'Neill's [PCG family of generators](https://www.pcg-random.org/)
generates better random numbers in about 2.1ns per number
with only 16 bytes of internal state.
We also wanted to explore using
Daniel J. Bernstein's [ChaCha stream cipher](https://cr.yp.to/chacha.html)
as a generator.
A [follow-up post](/blog/chacha8rand) discusses that generator specifically.

### Source Interface {#problem.source}

The [`rand.Source` interface](/pkg/math/rand/#Source) was wrong.
That interface defines the
concept of a low-level random number generator that generates
non-negative `int64` values:

{{raw `
	% go doc -src math/rand.Source
	package rand // import "math/rand"

	// A Source represents a source of uniformly-distributed
	// pseudo-random int64 values in the range [0, 1<<63).
	//
	// A Source is not safe for concurrent use by multiple goroutines.
	type Source interface {
		Int63() int64
		Seed(seed int64)
	}

	func NewSource(seed int64) Source
	%
`}}

(In the doc comment, “[0, N)” denotes a
[half-open interval](https://en.wikipedia.org/wiki/Interval_(mathematics)#Definitions_and_terminology),
meaning the range includes 0 but ends just before 2⁶³.)

The [`rand.Rand` type](/pkg/math/rand/#Rand) wraps a `Source`
to implement a richer set of operations, such as
generating [an integer between 0 and N](/pkg/math/rand/#Rand.Intn),
generating [floating-point numbers](/pkg/math/rand/#Rand.Float64), and so on.

We defined the `Source` interface to return a shortened 63-bit value
instead of a uint64 because that's what the Go 1 generator and
other widely-used generators produce,
and it matches the convention set by the C standard library.
But this was a mistake: more modern generators produce full-width uint64s,
which is a more convenient interface.

Another problem is the `Seed` method hard-coding an `int64` seed:
some generators are seeded by larger values,
and the interface provides no way to handle that.

### Seeding Responsibility {#problem.seed}

A bigger problem with `Seed` is that responsibility for seeding the global generator was unclear.
Most users don't use `Source` and `Rand` directly.
Instead, the `math/rand` package provides a global generator
accessed by top-level functions like [`Intn`](/pkg/math/rand/#Intn).
Following the C standard library, the global generator defaults to
behaving as if `Seed(1)` is called at startup.
This is good for repeatability but bad for programs that want
their random outputs to be different from one run to the next.
The package documentation suggests using `rand.Seed(time.Now().UnixNano())` in that case,
to make the generator's output time-dependent,
but what code should do this?

Probably the main package should be in charge of how `math/rand` is seeded:
it would be unfortunate for imported libraries to configure global state themselves,
since their choices might conflict with other libraries or the main package.
But what happens if a library needs some random data and wants to use `math/rand`?
What if the main package doesn't even know `math/rand` is being used?
We found that in practice many libraries add init functions
that seed the global generator with the current time, “just to be sure”.

Library packages seeding the global generator themselves causes a new problem.
Suppose package main imports two packages that both use `math/rand`:
package A assumes the global generator will be seeded by package main,
but package B seeds it in an `init` func.
And suppose that package main doesn't seed the generator itself.
Now package A's correct operation depends on the coincidence that package B is also
imported in the program.
If package main stops importing package B, package A will stop getting random values.
We observed this happening in practice in large codebases.

In retrospect, it was clearly a mistake to follow the C standard library here:
seeding the global generator automatically would remove the confusion
about who seeds it, and users would stop being surprised by repeatable
output when they didn't want that.

### Scalability {#problem.scale}

The global generator also did not scale well.
Because top-level functions like [`rand.Intn`](/pkg/math/rand/#Intn)
can be called simultaneously from multiple goroutines,
the implementation needed a lock protecting the shared generator state.
In parallel usage, acquiring and releasing this lock was more expensive
than the actual generation.
It would make sense instead to have a per-thread generator state,
but doing so would break repeatability
in programs without concurrent use of `math/rand`.

### The `Rand` implementation was missing important optimizations {#problem.rand}

The [`rand.Rand` type](/pkg/math/rand/#Rand) wraps a `Source`
to implement a richer set of operations.
For example, here is the Go 1 implementation of `Int63n`, which returns
a random integer in the range [0, `n`).

{{raw `
	func (r *Rand) Int63n(n int64) int64 {
		if n <= 0 {
			panic("invalid argument to Int63n")
		}
		max := int64((1<<63 - 1)  - (1<<63)%uint64(n))
		v := r.Int63()
		for v > max {
			v = r.Int63()
		}
		return v % n
	}
`}}

The actual conversion is easy: `v % n`.
However, no algorithm can convert 2⁶³ equally likely values
into `n` equally likely values unless 2⁶³ is a multiple of `n`:
otherwise some outputs will necessarily happen more often
than others. (As a simpler example, try converting 4 equally likely values into 3.)
The code computes `max` such that
`max+1` is the largest multiple of `n` less than or equal to 2⁶³,
and then the loop rejects random values greater than or equal to `max+1`.
Rejecting these too-large values ensures that all `n` outputs are equally likely.
For small `n`, needing to reject any value at all is rare;
rejection becomes more common and more important for larger values.
Even without the rejection loop, the two (slow) modulus operations
can make the conversion more expensive than generating the random value `v`
in the first place.

In 2018, [Daniel Lemire found an algorithm](https://arxiv.org/abs/1805.10941)
that avoids the divisions nearly all the time
(see also his [2019 blog post](https://lemire.me/blog/2019/06/06/nearly-divisionless-random-integer-generation-on-various-systems/)).
In `math/rand`, adopting Lemire's algorithm would make `Intn(1000)` 20-30% faster,
but we can't: the faster algorithm generates different values than the standard conversion,
breaking repeatability.

Other methods are also slower than they could be, constrained by repeatability.
For example, the `Float64` method could easily be sped up by about 10%
if we could change the generated value stream.
(This was the change we tried to make in Go 1.2 and rolled back, mentioned earlier.)

### The `Read` Mistake {#problem.read}

As mentioned earlier, `math/rand` is not intended for
and not suitable for generating cryptographic secrets.
The `crypto/rand` package does that, and its fundamental
primitive is its [`Read` function](/pkg/crypto/rand/#Read)
and [`Reader`](/pkg/crypto/rand/#Reader) variable.

In 2015, we accepted a proposal to make
`rand.Rand` implement `io.Reader` as well,
along with [adding a top-level `Read` function](/pkg/math/rand/#Read).
This seemed reasonable at the time,
but in retrospect we did not pay enough attention to the
software engineering aspects of this change.
Now, if you want to read random data, you have
two choices: `math/rand.Read` and `crypto/rand.Read`.
If the data is going to be used for key material,
it is very important to use `crypto/rand`,
but now it is possible to use `math/rand` instead,
potentially with disastrous consequences.

Tools like `goimports` and `gopls` have a special case
to make sure they prefer to use `rand.Read` from
`crypto/rand` instead of `math/rand`, but that's not a complete fix.
It would be better to remove `Read` entirely.

## Fixing `math/rand` directly {#fix.v1}

Making a new, incompatible major version of a package is never our first choice:
that new version only benefits programs that switch to it,
leaving all existing usage of the old major version behind.
In contrast, fixing a problem in the existing package has much more impact,
since it fixes all the existing usage.
We should never create a `v2` without doing as much as possible to fix `v1`.
In the case of `math/rand`, we were able to partly address
a few of the problems described above:

- Go 1.8 introduced an optional [`Source64` interface](/pkg/math/rand/#Uint64) with a `Uint64` method.
  If a `Source` also implements `Source64`, then `Rand` uses that method
  when appropriate.
  This “extension interface” pattern provides a compatible (if slightly awkward)
  way to revise an interface after the fact.

- Go 1.20 automatically seeded the top-level generator and
  deprecated [`rand.Seed`](/pkg/math/rand/#Seed).
  Although this may seem like an incompatible change
  given our focus on repeatability of the output stream,
  [we reasoned](/issue/56319) that any imported package that called [`rand.Int`](/pkg/math/rand/#Int)
  at init time or inside any computation would also
  visibly change the output stream, and surely adding or removing
  such a call cannot be considered a breaking change.
  And if that's true, then auto-seeding is no worse,
  and it would eliminate this source of fragility for future programs.
  We also added a [GODEBUG setting](/doc/godebug) to opt
  back into the old behavior.
  Then we marked the top-level `rand.Seed` as [deprecated](/wiki/Deprecated).
  (Programs that need seeded repeatability can still use
  `rand.New(rand.NewSource(seed))` to obtain a local generator
  instead of using the global one.)

- Having eliminated repeatability of the global output stream,
  Go 1.20 was also able to make the global generator scale better
  in programs that don't call `rand.Seed`,
  replacing the Go 1 generator with a very cheap per-thread
  [wyrand generator](https://github.com/wangyi-fudan/wyhash)
  already used inside the Go runtime. This removed the global mutex
  and made the top-level functions scale much better.
  Programs that do call `rand.Seed` fall back to the
  mutex-protected Go 1 generator.

- We were able to adopt Lemire's optimization in the Go runtime,
  and we also used it inside [`rand.Shuffle`](/pkg/math/rand/#Shuffle),
  which was implemented after Lemire's paper was published.

- Although we couldn't remove [`rand.Read`](/pkg/math/rand/#Read) entirely,
  Go 1.20 marked it [deprecated](/wiki/Deprecated) in favor of
  `crypto/rand`.
  We have since heard from people who discovered that they were accidentally
  using `math/rand.Read` in cryptographic contexts when their editors
  flagged the use of the deprecated function.

These fixes are imperfect and incomplete but also real improvements
that helped all users of the existing `math/rand` package.
For more complete fixes, we needed to turn our attention to `math/rand/v2`.

## Fixing the rest in `math/rand/v2` {#fix.v2}

Defining `math/rand/v2` took
significant planning,
then a [GitHub Discussion](/issue/60751)
and then a [proposal discussion](/issue/61716).
It is the same
as `math/rand` with the following breaking changes
addressing the problems outlined above:

- We removed the Go 1 generator entirely, replacing it with two new generators,
  [PCG](/pkg/math/rand/v2/#PCG) and [ChaCha8](/pkg/math/rand/v2/#ChaCha8).
  The new types are named for their algorithms (avoiding the generic name `NewSource`)
  so that if another important algorithm needs to be added, it will fit well into the
  naming scheme.

  Adopting a suggestion from the proposal discussion, the new types implement the
  [`encoding.BinaryMarshaler`](/pkg/encoding/#BinaryMarshaler)
  and
  [`encoding.BinaryUnmarshaler`](/pkg/encoding/#BinaryUnmarshaler)
  interfaces.

- We changed the `Source` interface, replacing the `Int63` method with a `Uint64` method
  and deleting the `Seed` method. Implementations that support seeding can provide
  their own concrete methods, like [`PCG.Seed`](/pkg/math/rand/v2/#PCG.Seed) and
  [`ChaCha8.Seed`](/pkg/math/rand/v2/#ChaCha8.Seed).
  Note that the two take different seed types, and neither is a single `int64`.

- We removed the top-level `Seed` function: the global functions like `Int` can only be used
  in auto-seeded form now.

- Removing the top-level `Seed` also let us hard-code the use of scalable,
  per-thread generators by the top-level methods,
  avoiding a GODEBUG check at each use.

- We implemented Lemire's optimization for `Intn` and related functions.
  The concrete `rand.Rand` API is now locked in to that value stream,
  so we will not be able to take advantage of any optimizations yet to be discovered,
  but at least we are up to date once again.
  We also implemented the `Float32` and `Float64` optimizations we wanted to use back in Go 1.2.

- During the proposal discussion, a contributor pointed out detectable bias in the
  implementations of `ExpFloat64` and `NormFloat64`.
  We fixed that bias and locked in the new value streams.

- `Perm` and `Shuffle` used different shuffling algorithms and produced different value streams,
  because `Shuffle` happened second and used a faster algorithm.
  Deleting `Perm` entirely would have made migration harder for users.
  Instead we implemented `Perm` in terms of `Shuffle`, which still lets us
  delete an implementation.

- We renamed `Int31`, `Int63`, `Intn`, `Int31n`, and `Int63n` to
  `Int32`, `Int64`, `IntN`, `Int32N`, and `Int64N`.
  The 31 and 63 in the names were unnecessarily pedantic
  and confusing, and the capitalized N is more idiomatic for a second
  “word” in the name in Go.

- We added `Uint`, `Uint32`, `Uint64`, `UintN`, `Uint32N`, and `Uint64N`
  top-level functions and methods.
  We needed to add `Uint64` to provide direct access to the core `Source`
  functionality, and it seemed inconsistent not to add the others.

- Adopting another suggestion from the proposal discussion,
  we added a new top-level, generic function `N` that is like
  `Int64N` or `Uint64N` but works for any integer type.
  In the old API, to create a random duration of up to 5 seconds,
  it was necessary to write:

      d := time.Duration(rand.Int63n(int64(5*time.Second)))

  Using `N`, the equivalent code is:

      d := rand.N(5 * time.Second)

  `N` is only a top-level function; there is no `N` method on `rand.Rand`
  because there are no generic methods in Go.
  (Generic methods are not likely in the future, either;
  they conflict badly with interfaces, and a complete implementation
  would require either run-time code generation or slow execution.)

- To ameliorate misuse of `math/rand` in cryptographic contexts,
  we made `ChaCha8` the default generator used in global functions,
  and we also changed the Go runtime to use it (replacing wyrand).
  Programs are still strongly encouraged to use `crypto/rand`
  to generate cryptographic secrets,
  but accidentally using `math/rand/v2` is not as catastrophic
  as using `math/rand` would be.
  Even in `math/rand`, the global functions now use the `ChaCha8` generator when not explicitly seeded.

## Principles for evolving the Go standard library {#principles}

As mentioned at the start this post, one of the goals for this work
was to establish principles and a pattern for how we approach all
v2 packages in the standard library.
There will not be a glut of v2 packages
in the next few Go releases.
Instead, we will handle one package
at a time, making sure we set a quality bar that will last for another decade.
Many packages will not need a v2 at all.
But for those that do, our approach boils down to three principles.

First, a new, incompatible version of a package will use
`that/package/v2` as its import path,
following
[semantic import versioning](https://research.swtch.com/vgo-import)
just like a v2 module outside the standard library would.
This allows uses of the original package and the v2 package
to coexist in a single program,
which is critical for a
[gradual conversion](/talks/2016/refactor.article) to the new API.

Second, all changes must be rooted in
respect for existing usage and users:
we must not introduce needless churn,
whether in the form of unnecessary changes to an existing package or
an entirely new package that must be learned instead.
In practice, that means we take the existing package
as the starting point
and only make changes that are well motivated
and provide a value that justifies the cost to users of updating.

Third, the v2 package must not leave v1 users behind.
Ideally, the v2 package should be able to do everything the v1 package
could do,
and when v2 is released, the v1 package should be rewritten
to be a thin wrapper around v2.
This would ensure that existing uses of v1 continue to benefit
from bug fixes and performance optimizations in v2.
Of course, given that v2 is introducing breaking changes,
this is not always possible, but it is always something to consider carefully.
For `math/rand/v2`, we arranged for the auto-seeded v1 functions to
call the v2 generator, but we were unable to share other code
due to the repeatability violations.
Ultimately `math/rand` is not a lot of code and does not require
regular maintenance, so the duplication is manageable.
In other contexts, more work to avoid duplication could be worthwhile.
For example, in the
[encoding/json/v2 design (still in progress)](/issue/63397),
although the default semantics and the API are changed,
the package provides configuration knobs that
make it possible to implement the v1 API.
When we eventually ship `encoding/json/v2`,
`encoding/json` (v1) will become a thin wrapper around it,
ensuring that users who don't migrate from v1 still
benefit from optimizations and security fixes in v2.

A [follow-up blog post](/blog/chacha8rand) presents the `ChaCha8` generator in more detail.
