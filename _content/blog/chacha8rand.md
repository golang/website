---
title: "Secure Randomness in Go 1.22"
date: 2024-05-02
by:
- Russ Cox
- Filippo Valsorda
summary: ChaCha8Rand is a new cryptographically secure pseudorandom number generator used in Go 1.22.
---

Computers aren't random.
On the contrary, hardware designers work very hard to make sure computers run every program the same way every time.
So when a program does need random numbers, that requires extra effort.
Traditionally, computer scientists and programming languages
have distinguished between two different kinds of random numbers:
statistical and cryptographic randomness.
In Go, those are provided by [`math/rand`](/pkg/math/rand/)
and [`crypto/rand`](/pkg/crypto/rand), respectively.
This post is about how Go 1.22 brings the two closer together,
by using a cryptographic random number source in `math/rand`
(as well as `math/rand/v2`, as mentioned in our [previous post](/blog/randv2)).
The result is better randomness and far less damage when
developers accidentally use `math/rand` instead of `crypto/rand`.

Before we can explain what Go 1.22 did, let's take a closer look
at statistical randomness compared to cryptographic randomness.

## Statistical Randomness

Random numbers that pass basic statistical tests
are usually appropriate for use cases like simulations, sampling,
numerical analysis, non-cryptographic randomized algorithms,
[random testing](/doc/security/fuzz/),
[shuffling inputs](https://en.wikipedia.org/wiki/Fisher%E2%80%93Yates_shuffle),
and
[random exponential backoff](https://en.wikipedia.org/wiki/Exponential_backoff#Collision_avoidance).
Very basic, easy to compute mathematical formulas turn out to work
well enough for these use cases.
Because the methods are so simple, however, an observer who
knows what algorithm is being used can typically predict the rest
of the sequence after seeing enough values.

Essentially all programming environments provide a mechanism for generating
statistical random numbers
that traces back through C to
Research Unix Third Edition (V3), which added a pair of functions: `srand` and `rand`.
The manual page included
a note that read:

> _WARNING   The author of this routine has been writing
random-number generators for many years and has never been
known to write one that worked._

This note was partly a joke but also an acknowledgement that such
generators are [inherently not random](https://www.tuhs.org/pipermail/tuhs/2024-March/029587.html).

The source code of the generator makes clear how trivial it is.
Translated from PDP-11 assembly to modern C, it was:

	uint16 ranx;

	void
	srand(uint16 seed)
	{
	    ranx = seed;
	}

	int16
	rand(void)
	{
	    ranx = 13077*ranx + 6925;
	    return ranx & ~0x8000;
	}

Calling `srand` seeds the generator with a single integer seed,
and `rand` returns the next number from the generator.
The AND in the return statement clears the sign bit to make sure the result is positive.

This function is an instance of the general class of
[linear congruential generators (LCGs)](https://en.wikipedia.org/wiki/Linear_congruential_generator),
which Knuth analyzes in _The Art of Computer Programming_, Volume 2, section 3.2.1.
The main benefit of LCGs is that constants can be chosen such that they
emit every possible output value once before repeating,
as the Unix implementation did for 15-bit outputs.
A serious problem with LCGs, however, is that the high bits of the state do not affect the low bits at all,
so every truncation of the sequence to _k_ bits necessarily repeats with a smaller period.
The low bit must toggle: 0, 1, 0, 1, 0, 1.
The low two bits must count up or down: 0, 1, 2, 3, 0, 1, 2, 3, or else 0, 3, 2, 1, 0, 3, 2, 1.
There are four possible three-bit sequences; the original Unix implementation repeats 0, 5, 6, 3, 4, 1, 2, 7.
(These problems can be avoided by reducing the value modulo a prime,
but that would have been quite expensive at the time.
See S. K. Park and K. W. Miller's 1988 CACM paper
“[Random number generators: good ones are hard to find](https://dl.acm.org/doi/10.1145/63039.63042)”
for a short analysis
and the first chapter of Knuth Volume 2 for a longer one.)

Even with these known problems,
the `srand` and `rand` functions were included in the first C standard,
and equivalent functionality was included in essentially every language since then.
LCGs were once the dominant implementation strategy,
although they've fallen off in popularity due to some important drawbacks.
One significant remaining use is [`java.util.Random`](https://github.com/openjdk/jdk8u-dev/blob/master/jdk/src/share/classes/java/util/Random.java),
which powers [`java.lang.Math.random`](https://github.com/openjdk/jdk8u-dev/blob/master/jdk/src/share/classes/java/util/Random.java).

Another thing you can see from the implementation above
is that the internal state is completely exposed by the result of `rand`.
An observer who knows the algorithm and sees a single result
can easily compute all future results.
If you are running a server that calculates some random values
that become public and some random values that must stay secret,
using this kind of generator would be disastrous:
the secrets wouldn't be secret.

More modern random generators aren't as terrible as the original Unix one,
but they're still not completely unpredictable.
To make that point, next we will look at the original `math/rand` generator from Go 1
and the PCG generator we added in `math/rand/v2`.

## The Go 1 Generator

The generator used in Go 1's `math/rand` is an instance of what is called a
[linear-feedback shift register](https://en.wikipedia.org/wiki/Linear-feedback_shift_register).
The algorithm is based on an idea by George Marsaglia,
tweaked by Don Mitchell and Jim Reeds,
and further customized by Ken Thompson for Plan 9 and then Go.
It has no official name, so this post calls it the Go 1 generator.

The Go 1 generator's internal state is a slice `vec` of 607 uint64s.
In that slice, there are two distinguished elements: `vec[606]`, the last element, is called the “tap”,
and `vec[334]` is called the “feed”.
To generate the next random number,
the generator adds the tap and the feed
to produce a value `x`,
stores `x` back into the feed,
shifts the entire slice one position to the right
(the tap moves to `vec[0]` and `vec[i]` moves to `vec[i+1]`),
and returns `x`.
The generator is called “linear feedback” because the tap is _added_ to the feed;
the entire state is a “shift register” because each step shifts the slice entries.

Of course, actually moving every slice entry forward would be prohibitively expensive,
so instead the implementation leaves the slice data in place
and moves the tap and feed positions backward
on each step. The code looks like:

{{raw `
	func (r *rngSource) Uint64() uint64 {
		r.tap--
		if r.tap < 0 {
			r.tap += len(r.vec)
		}

		r.feed--
		if r.feed < 0 {
			r.feed += len(r.vec)
		}

		x := r.vec[r.feed] + r.vec[r.tap]
		r.vec[r.feed] = x
		return uint64(x)
	}
`}}

Generating the next number is quite cheap: two subtractions, two conditional adds, two loads, one add, one store.

Unfortunately, because the generator directly returns one slice element from its internal state vector,
reading 607 values from the generator completely exposes all its state.
With those values, you can predict all the future values, by filling in your own `vec`
and then running the algorithm.
You can also recover all the previous values, by running the algorithm backward
(subtracting the tap from the feed and shifting the slice to the left).

As a complete demonstration, here is an [insecure program](/play/p/v0QdGjUAtzC)
generating pseudorandom authentication
tokens along with code that predicts the next token given a sequence of earlier tokens.
As you can see, the Go 1 generator provides no security at all (nor was it meant to).
The quality of the generated numbers also depends on the initial setting of `vec`.

## The PCG Generator

For `math/rand/v2`, we wanted to provide a more modern statistical random generator
and settled on Melissa O'Neill's PCG algorithm, published in 2014 in her paper
“[PCG: A Family of Simple Fast Space-Efficient Statistically Good Algorithms for Random Number Generation](https://www.pcg-random.org/pdf/hmc-cs-2014-0905.pdf)”.
The exhaustive analysis in the paper can make it hard to notice at first glance how utterly trivial the generators are:
PCG is a post-processed 128-bit LCG.

If the state `p.x` were a `uint128` (hypothetically), the code to compute the next value would be:

	const (
		pcgM = 0x2360ed051fc65da44385df649fccf645
		pcgA = 0x5851f42d4c957f2d14057b7ef767814f
	)

	type PCG struct {
		x uint128
	}

	func (p *PCG) Uint64() uint64 {
		p.x = p.x * pcgM + pcgA
		return scramble(p.x)
	}

The entire state is a single 128-bit number,
and the update is a 128-bit multiply and add.
In the return statement, the `scramble` function reduces the 128-bit state
down to a 64-bit state.
The original PCG used (again using a hypothetical `uint128` type):

	func scramble(x uint128) uint64 {
		return bits.RotateLeft(uint64(x>>64) ^ uint64(x), -int(x>>122))
	}

This code XORs the two halves of the 128-bit state together
and then rotates the result according to the top six bits of the state.
This version is called PCG-XSL-RR, for “xor shift low, right rotate”.

Based on a [suggestion from O'Neill during proposal discussion](/issue/21835#issuecomment-739065688),
Go's PCG uses a new scramble function based on multiplication,
which mixes the bits more aggressively:

	func scramble(x uint128) uint64 {
		hi, lo := uint64(x>>64), uint64(x)
		hi ^= hi >> 32
		hi *= 0xda942042e4dd58b5
		hi ^= hi >> 48
		hi *= lo | 1
	}

O'Neill calls PCG with this scrambler PCG-DXSM, for “double xorshift multiply.”
Numpy uses this form of PCG as well.

Although PCG uses more computation to generate each value,
it uses significantly less state: two uint64s instead of 607.
It is also much less sensitive to the initial values of that state,
and [it passes many statistical tests that other generators do not](https://www.pcg-random.org/statistical-tests.html).
In many ways it is an ideal statistical generator.

Even so, PCG is not unpredictable.
While the scrambling of bits to prepare the result does not
expose the state directly like in the LCG and Go 1 generators,
[PCG-XSL-RR can still be reversed](https://pdfs.semanticscholar.org/4c5e/4a263d92787850edd011d38521966751a179.pdf),
and it would not be surprising if PCG-DXSM could too.
For secrets, we need something different.

## Cryptographic Randomness

_Cryptographic random numbers_ need to be utterly unpredictable
in practice, even to an observer who knows how they are generated
and has observed any number of previously generated values.
The safety of cryptographic protocols, secret keys, modern commerce,
online privacy, and more all critically depend on access to cryptographic
randomness.

Providing cryptographic randomness is ultimately the job of the
operating system, which can gather true randomness from physical devices—timings
of the mouse, keyboard, disks, and network, and more recently
[electrical noise measured directly by the CPU itself](https://web.archive.org/web/20141230024150/http://www.cryptography.com/public/pdf/Intel_TRNG_Report_20120312.pdf).
Once the operating system has gathered a meaningful
amount of randomness—say, at least 256 bits—it can use cryptographic
hashing or encryption algorithms to stretch that seed into
an arbitrarily long sequence of random numbers.
(In practice the operating system is also constantly gathering and
adding new randomness to the sequence too.)

The exact operating system interfaces have evolved over time.
A decade ago, most systems provided a device file named
`/dev/random` or something similar.
Today, in recognition of how fundamental randomness has become,
operating systems provide a direct system call instead.
(This also allows programs to read randomness even
when cut off from the file system.)
In Go, the [`crypto/rand`](/pkg/crypto/rand/) package abstracts away those details,
providing the same interface on every operating system: [`rand.Read`](/pkg/crypto/rand/#Read).

It would not be practical for `math/rand` to ask the operating system for
randomness each time it needs a `uint64`.
But we can use cryptographic techniques to define an in-process
random generator that improves on LCGs, the Go 1 generator, and even PCG.

## The ChaCha8Rand Generator

Our new generator, which we unimaginatively named ChaCha8Rand for specification purposes
and implemented as `math/rand/v2`'s [`rand.ChaCha8`](/pkg/math/rand/v2/#ChaCha8),
is a lightly modified version of Daniel J. Bernstein's [ChaCha stream cipher](https://cr.yp.to/chacha.html).
ChaCha is widely used in a 20-round form called ChaCha20, including in TLS and SSH.
Jean-Philippe Aumasson's paper “[Too Much Crypto](https://eprint.iacr.org/2019/1492.pdf)”
argues persuasively that the 8-round form ChaCha8 is secure too (and it's roughly 2.5X faster).
We used ChaCha8 as the core of ChaCha8Rand.

Most stream ciphers, including ChaCha8, work by defining a function that is given
a key and a block number and produces a fixed-size block of apparently random data.
The cryptographic standard these aim for (and usually meet) is for this output to be indistinguishable
from actual random data in the absence of some kind of exponentially costly brute force search.
A message is encrypted or decrypted by XOR'ing successive blocks of input data
with successive randomly generated blocks.
To use ChaCha8 as a `rand.Source`,
we use the generated blocks directly instead of XOR'ing them with input data
(this is equivalent to encrypting or decrypting all zeros).

We changed a few details to make ChaCha8Rand more suitable for generating random numbers. Briefly:

 - ChaCha8Rand takes a 32-byte seed, used as the ChaCha8 key.
 - ChaCha8 generates 64-byte blocks, with calculations treating a block as 16 `uint32`s.
   A common implementation is to compute four blocks at a time using [SIMD instructions](https://en.wikipedia.org/wiki/Single_instruction,_multiple_data)
   on 16 vector registers of four `uint32`s each.
   This produces four interleaved blocks that must be unshuffled for XOR'ing with the input data.
   ChaCha8Rand defines that the interleaved blocks are the random data stream,
   removing the cost of the unshuffle.
   (For security purposes, this can be viewed as standard ChaCha8 followed by a reshuffle.)
 - ChaCha8 finishes a block by adding certain values to each `uint32` in the block.
   Half the values are key material and the other half are known constants.
   ChaCha8Rand defines that the known constants are not re-added,
   removing half of the final adds.
   (For security purposes, this can be viewed as standard ChaCha8 followed by subtracting the known constants.)
 - Every 16th generated block, ChaCha8Rand takes the final 32 bytes of the block for itself,
   making them the key for the next 16 blocks.
   This provides a kind of [forward secrecy](https://en.wikipedia.org/wiki/Forward_secrecy):
   if a system is compromised by an attack that
   recovers the entire memory state of the generator, only values generated
   since the last rekeying can be recovered. The past is inaccessible.
   ChaCha8Rand as defined so far must generate 4 blocks at a time,
   but we chose to do this key rotation every 16 blocks to leave open the
   possibility of faster implementations using 256-bit or 512-bit vectors,
   which could generate 8 or 16 blocks at a time.

We wrote and published a [C2SP specification for ChaCha8Rand](https://c2sp.org/chacha8rand),
along with test cases.
This will enable other implementations to share repeatability with the Go implementation
for a given seed.

The Go runtime now maintains a per-core ChaCha8Rand state (300 bytes),
seeded with operating system-supplied cryptographic randomness,
so that random numbers can be generated quickly without any lock contention.
Dedicating 300 bytes per core may sound expensive,
but on a 16-core system, it is about the same as storing a single shared Go 1 generator state (4,872 bytes).
The speed is worth the memory.
This per-core ChaCha8Rand generator is now used in three different places in the Go standard library:

 1. The `math/rand/v2` package functions, such as
   [`rand.Float64`](/pkg/math/rand/v2/#Float64) and
   [`rand.N`](/pkg/math/rand/v2/#N), always use ChaCha8Rand.

 2. The `math/rand` package functions, such as
   [`rand.Float64`](/pkg/math/rand/#Float64) and
   [`rand.Intn`](/pkg/math/rand/#Intn),
   use ChaCha8Rand when
   [`rand.Seed`](/pkg/math/rand/#Seed) has not been called.
   Applying ChaCha8Rand in `math/rand` improves the security of programs
   even before they update to `math/rand/v2`,
   provided they are not calling `rand.Seed`.
   (If `rand.Seed` is called, the implementation is required to fall back to the Go 1 generator for compatibility.)

 3. The runtime chooses the hash seed for each new map
    using ChaCha8Rand instead of a less secure [wyrand-based generator](https://github.com/wangyi-fudan/wyhash)
    it previously used.
    Random seeds are needed because if
    an attacker knows the specific hash function used by a map implementation,
    they can prepare input that drives the map into quadratic behavior
    (see Crosby and Wallach's “[Denial of Service via Algorithmic Complexity Attacks](https://www.usenix.org/conference/12th-usenix-security-symposium/denial-service-algorithmic-complexity-attacks)”).
    Using a per-map seed, instead of one global seed for all maps,
    also avoids [other degenerate behaviors](https://accidentallyquadratic.tumblr.com/post/153545455987/rust-hash-iteration-reinsertion).
    It is not strictly clear that maps need a cryptographically random seed,
    but it's also not clear that they don't. It seemed prudent and was trivial to switch.

Code that needs its own ChaCha8Rand instances can create its own [`rand.ChaCha8`](/pkg/math/rand/v2/#ChaCha8) directly.

## Fixing Security Mistakes

Go aims to help developers write code that is secure by default.
When we observe a common mistake with security consequences,
we look for ways to reduce the risk of that mistake
or eliminate it entirely.
In this case, `math/rand`'s global generator was far too predictable,
leading to serious problems in a variety of contexts.

For example, when Go 1.20 deprecated [`math/rand`’s `Read`](/pkg/math/rand/#Read),
we heard from developers who discovered (thanks to tooling pointing out
use of deprecated functionality) they had been
using it in places where [`crypto/rand`’s `Read`](/pkg/crypto/rand/#Read)
was definitely needed, like generating key material.
Using Go 1.20, that mistake
is a serious security problem that merits a detailed investigation
to understand the damage.
Where were the keys used?
How were the keys exposed?
Were other random outputs exposed that might allow an attacker to derive the keys?
And so on.
Using Go 1.22, that mistake is just a mistake.
It's still better to use `crypto/rand`,
because the operating system kernel can do a better job keeping the random values
secret from various kinds of prying eyes,
the kernel is continually adding new entropy to its generator,
and the kernel has had more scrutiny.
But accidentally using `math/rand` is no longer a security catastrophe.

There are also a variety of use cases that don't seem like “crypto”
but nonetheless need unpredictable randomness.
These cases are made more robust by using ChaCha8Rand instead of the Go 1 generator.

For example, consider generating a
[random UUID](https://en.wikipedia.org/wiki/Universally_unique_identifier#Version_4_(random)).
Since UUIDs are not secret, using `math/rand` might seem fine.
But if `math/rand` has been seeded with the current time,
then running it at the same instant on different computers
will produce the same value, making them not “universally unique”.
This is especially likely on systems where the current time
is only available with millisecond precision.
Even with auto-seeding using OS-provided entropy,
as introduced in Go 1.20,
the Go 1 generator's seed is only a 63-bit integer,
so a program that generates a UUID at startup
can only generate 2⁶³ possible UUIDs and is
likely to see collisions after 2³¹ or so UUIDs.
Using Go 1.22, the new ChaCha8Rand generator
is seeded from 256 bits of entropy and can generate
2²⁵⁶ possible first UUIDs.
It does not need to worry about collisions.

As another example, consider load balancing in a front-end server
that randomly assigns incoming requests to back-end servers.
If an attacker can observe the assignments and knows the
predictable algorithm generating them,
then the attacker could send a stream
of mostly cheap requests but arrange for all the expensive requests
to land on a single back-end server.
This is an unlikely but plausible problem using the Go 1 generator.
Using Go 1.22, it's not a problem at all.

In all these examples, Go 1.22 has eliminated or greatly reduced
security problems.

## Performance

The security benefits of ChaCha8Rand do have a small cost,
but ChaCha8Rand is still in the same ballpark as both the Go 1 generator and PCG.
The following graphs compare the performance of the three generators,
across a variety of hardware, running two operations:
the primitive operation “Uint64,” which returns the next `uint64` in the random stream,
and the higher-level operation “N(1000),” which returns a random value in the range [0, 1000).

<div style="background-color: white;">
<img src="chacha8rand/amd.svg">
<img src="chacha8rand/intel.svg">
<img src="chacha8rand/amd32.svg">
<img src="chacha8rand/intel32.svg">
<img src="chacha8rand/m1.svg">
<img src="chacha8rand/m3.svg">
<img src="chacha8rand/taut2a.svg">
</div>

The “running 32-bit code” graphs show modern 64-bit x86 chips
executing code built with `GOARCH=386`, meaning they are
running in 32-bit mode.
In that case, the fact that PCG requires 128-bit multiplications
makes it slower than ChaCha8Rand, which only uses 32-bit SIMD arithmetic.
Actual 32-bit systems matter less every year,
but it is still interesting that ChaCha8Rand is faster than PCG
on those systems.

On some systems, “Go 1: Uint64” is faster than “PCG: Uint64”,
but “Go 1: N(1000)” is slower than “PCG: N(1000)”.
This happens because “Go 1: N(1000)” is using `math/rand`'s algorithm for
reducing a random `int64` down to a value in the range [0, 1000),
and that algorithm does two 64-bit integer divide operations.
In contrast, “PCG: N(1000)” and “ChaCha8: N(1000)” use the [faster `math/rand/v2` algorithm](/blog/randv2#problem.rand),
which almost always avoids the divisions.
Removing the 64-bit divisions dominates the algorithm change
for 32-bit execution and on the Ampere.

Overall, ChaCha8Rand is slower than the Go 1 generator,
but it is never more than twice as slow, and on typical servers the
difference is never more than 3ns.
Very few programs will be bottlenecked by this difference,
and many programs will enjoy the improved security.

## Conclusion

Go 1.22 makes your programs more secure without any code changes.
We did this by identifying the common mistake of accidentally using `math/rand`
instead of `crypto/rand` and then strengthening `math/rand`.
This is one small step in Go's ongoing journey to keep programs
safe by default.

These kinds of mistakes are not unique to Go.
For example, the npm `keypair` package tries to generate an RSA key pair
using Web Crypto APIs, but if they're not available, it falls back to JavaScript's `Math.random`.
This is hardly an isolated case,
and the security of our systems cannot depend on developers not making mistakes.
Instead, we hope that eventually all programming languages
will move to cryptographically strong pseudorandom generators
even for “mathematical” randomness,
eliminating this kind of mistake, or at least greatly reducing its blast radius.
Go 1.22's [ChaCha8Rand](https://c2sp.org/chacha8rand) implementation
proves that this approach is competitive with other generators.

