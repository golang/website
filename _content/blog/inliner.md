---
title: "//go:fix inline and the source-level inliner"
date: 2026-03-10
by:
- Alan Donovan
tags:
- go fix
- go vet
- analysis framework
- modernizers
- source-level inliner
- static analysis
summary: "How Go 1.26's source-level inliner works, and how it can help you with self-service API migrations."
---

<style>
.beforeafter {
  justify-content: center;
  display: grid;
  gap: 1em;
  margin: 1em;
  grid-template-columns: minmax(min-content, 1fr) auto minmax(min-content, 1fr);
  font-size: 180%;
  @media screen and (max-width: 57.7rem) {
    grid-template-columns: 1fr;
  }
}
#content .beforeafter pre {
  margin: 0em; /* Handled by grid gap */
}
.beforeafter-context {
  grid-column: 1 / -1;
}
#content .beforeafter > pre:nth-of-type(1) { background: var(--color-diff-old); }
#content .beforeafter > pre:nth-of-type(2) { background: var(--color-diff-new); }
.beforeafter-arrow {
  place-self: center;
  /* Undo unnecessary grid gap. */
  margin: -0.5em;
}
.beforeafter-arrow::before {
  content: "⟶";
  @media screen and (max-width: 57.7rem) {
    content: "⇓";
  }
}
</style>

Go 1.26 contains an all-new implementation of the `go fix` subcommand,
designed to help you keep your Go code up-to-date and modern. For an
introduction, start by reading our [recent post](gofix) on the topic.
In this post, we’ll look at one particular feature, the source-level
inliner.

While `go fix` has several bespoke modernizers for specific new
language and library features,
the source-level inliner is the first fruit of our efforts to provide
“[self-service](gofix#self-service)” modernizers and analyzers.
It enables any package author to express simple API migrations and
updates in a straightforward and safe way.
We’ll first explain what the source-level inliner is and how you can use it,
then we’ll dive into some aspects of the problem and the technology behind it.

## Source-level inlining

In 2023, we built an [algorithm](https://pkg.go.dev/golang.org/x/tools/internal/refactor/inline) for source-level inlining of function calls in Go. To “inline” a call means to replace the call by a copy of the body of the called function, substituting arguments for parameters. We call it “source-level” inlining because it durably modifies the source code. By contrast, the inlining algorithm found in a typical compiler, including Go’s, applies a similar transformation, but to the compiler’s ephemeral [intermediate representation](https://en.wikipedia.org/wiki/Intermediate_representation), to generate more efficient code.

If you’ve ever invoked [gopls](/gopls/)’ "[Inline call](/gopls/features/transformation#refactorinlinecall-inline-call-to-function)" interactive refactoring, you’ve used the source-level inliner. (In VS Code, this code action can be found on the “Source Action…” menu.) The before-and-after screenshots below show the effect of inlining the call to `sum` from the function named `six`.

<center>
<img src="/gopls/assets/inline-before.png"/>

<img src="/gopls/assets/inline-after.png"/>
</center>

The inliner is a crucial building block for a number of source transformation tools. For example, gopls uses it for the “Change signature” and “Remove unused parameter” refactorings because, as we’ll see below, it takes care of many subtle correctness issues that arise when refactoring function calls.

This same inliner is also one of the analyzers in the all-new `go fix` command.
In `go fix`, it enables self-service API migration and upgrades using a new `//go:fix inline` directive comment.
Let's take a look at a few examples of how this works and what it can be used for.

### Example: renaming `ioutil.ReadFile`

In Go 1.16, the `ioutil.ReadFile` function, which reads the content of a file, was deprecated in favor of the new `os.ReadFile` function. In effect, the function was renamed, though of course Go’s [compatibility promise](/doc/go1compat) prevents us from ever removing the old name.

```go
package ioutil

import "os"

// ReadFile reads the file named by filename…
// Deprecated: As of Go 1.16, this function simply calls [os.ReadFile].
func ReadFile(filename string) ([]byte, error) {
	return os.ReadFile(filename)
}
```

Ideally, we would like to change every Go program in the world to stop using `ioutil.ReadFile` and to call `os.ReadFile` instead. The inliner can help us do that. First we annotate the old function with `//go:fix inline`. This comment tells the tool that any time it sees a call to this function, it should inline the call.

```go
package ioutil

import "os"

// ReadFile reads the file named by filename…
// Deprecated: As of Go 1.16, this function simply calls [os.ReadFile].
//go:fix inline
func ReadFile(filename string) ([]byte, error) {
	return os.ReadFile(filename)
}
```

When we run `go fix` on a file containing a call to `ioutil.ReadFile`, it applies the replacement:

```
$ go fix -diff ./...
-import "io/ioutil"
+import "os"

-	data, err := ioutil.ReadFile("hello.txt")
+	data, err := os.ReadFile("hello.txt")
```

The call has been inlined, in effect replacing a call to one function by a call to another.

Because the inliner replaces a function call by a copy of the body of
the called function, not by some arbitrary expression, in principle
the transformation should not change the program’s behavior
(barring code that inspects the call stack, of course).
This differs from other tools that allow for arbitrary rewrites,
such as `gofmt -r`, which are very powerful but need to be watched closely.

For many years now, our Google colleagues on the teams supporting
Java, Kotlin, and C++ have been using source-level inliner tools like this.
To date, these tools have eliminated millions of calls to deprecated
functions in Google’s code base.
Users simply add the directives, and wait.
During the night, robots quietly prepare, test, and submit batches of
code changes across a monorepo of billions of lines of code.
If all goes well, by the morning the old code is no longer in use and can be
safely deleted.
Go’s inliner is a relative newcomer, but it has already been used to
prepare more than 18,000 changelists to Google’s monorepo.

### Example: fixing API design flaws

With a little creativity, a variety of migrations can be expressed as inlinings.
Consider this hypothetical `oldmath` package:

```go
// Package oldmath is the bad old math package.
package oldmath

// Sub returns x - y.
func Sub(y, x int) int

// Inf returns positive infinity.
func Inf() float64

// Neg returns -x.
func Neg(x int) int
```

It has several design flaws: the `Sub` function declares its parameters in the wrong order; the `Inf` function implicitly prefers one of the two infinities; and the `Neg` function is redundant with `Sub`. Fortunately we have a `newmath` package that avoids these mistakes, and we’d like to get users to switch to it. The first step is to implement the old API in terms of the new package and to deprecate the old functions. Then we add inliner directives:

```
// Package oldmath is the bad old math package.
package oldmath

import "newmath"

// Sub returns x - y.
// Deprecated: the parameter order is confusing.
//go:fix inline
func Sub(y, x int) int {
	return newmath.Sub(x, y)
}

// Inf returns positive infinity.
// Deprecated: there are two infinite values; be explicit.
//go:fix inline
func Inf() float64 {
	return newmath.Inf(+1)
}

// Neg returns -x.
// Deprecated: this function is unnecessary.
//go:fix inline
func Neg(x int) int {
	return newmath.Sub(0, x)
}
```

Now, when users of `oldmath` run the `go fix` command on their code, it will replace all calls to the old functions by their new counterparts. By the way, gopls has included `inline` in its analyzer suite for some time, so if your editor uses gopls, the moment you add the `//go:fix inline` directives you should start seeing a diagnostic at each call site, such as “call of `oldmath.Sub` should be inlined”, along with a suggested fix that inlines that particular call.

For example, this old code:
```
import "oldmath"

var nine = oldmath.Sub(1, 10) // diagnostic: "call to oldmath.Sub should be inlined"
```
will be transformed to:
```
import "newmath"

var nine = newmath.Sub(10, 1)
```
Observe that after the fix, the arguments to `Sub` are in the logical order. This is progress! If you’re in luck, the inliner will succeed at removing every call to the functions in `oldmath`, perhaps allowing you to delete it as a dependency.

The `inline` analyzer works on types and constants too. If our `oldmath` package had originally declared a data type for rational numbers and a constant for π, we could use the following forwarding declarations to migrate them to the `newmath` package while preserving the behavior of existing code:
```
package oldmath

//go:fix inline
type Rational = newmath.Rational

//go:fix inline
const Pi = newmath.Pi
```

Each time the `inline` analyzer encounters a reference to `oldmath.Rational` or `oldmath.Pi`, it will update them to refer instead to `newmath`.

## Under the hood of the inliner

At a glance, source inlining seems straightforward: just replace the
call with the body of the callee function, introduce variables for the
function parameters, and bind the call arguments to those variables.
But handling all of the complexities and corner cases correctly
while producing acceptable results is no small technical challenge:
the inliner is about 7,000 lines of dense, compiler-like logic.
Let’s look at six aspects of the problem that make it so tricky.

### 1. Parameter elimination

One of the inliner’s most important tasks is to attempt to replace each occurrence of a parameter in the callee by its corresponding argument from the call. In the simplest case, the argument is a trivial literal such as `0` or `""`, so the replacement is straightforward and the parameter can be eliminated.

<div class="beforeafter">
<div class="beforeafter-context"><pre>
//go:fix inline
func show(prefix, item string) {
	fmt.Println(prefix, item)
}
</pre></div>
<pre>
show("", "hello")
</pre>
<div class="beforeafter-arrow"></div>
<pre>
fmt.Println("", "hello")
</pre>
</div>

For less trivial literals such as `404` or `"go.dev"`, the replacement is equally straightforward, so long as the parameter appears in the callee at most once. But if it appears multiple times, it would be bad style to sprinkle copies of these magic values throughout the code as it would obscure the relationship between them; a later change to only one of them might create an inconsistency.

In such cases the inliner must tread carefully and emit a more conservative result. Whenever one or more parameters cannot be completely substituted for any reason, the inliner inserts an explicit “parameter binding” declaration:

<div class="beforeafter">
<div class="beforeafter-context"><pre>
//go:fix inline
func printPair(before, x, y, after string) {
	fmt.Println(before, x, after)
	fmt.Println(before, y, after)
}
</pre></div>
<pre>
printPair("[", "one", "two", "]")
</pre>
<div class="beforeafter-arrow"></div>
<pre>
// a “parameter binding” declaration
var before, after = "[", "]"
fmt.Println(before, "one", after)
fmt.Println(before, "two", after)
</pre>
</div>

### 2. Side effects

In Go, as in all imperative programming languages, calling a function may have the side effect of updating variables, which in turn may affect the behavior of other functions. Consider the call to `add` below:

```go
func add(x, y int) int { return y + x }

z = add(f(), g())
```

A trivial inlining of the call would replace `x` with `f()` and `y` with `g()`, with this result:

```
z = g() + f()
```

But this result is incorrect because evaluation of `g()` now occurs before `f()`; if the two functions have side effects, those effects will now be observed in a different order and may affect the result of the expression. Of course, it is bad form to write code that relies on effect ordering among call arguments, but that doesn’t mean people don’t do it, and our tools have to get it right.

So, the inliner must attempt to prove that `f()` and `g()` do not have side effects on each other. On success, it can safely proceed with the result above. Otherwise, it must fall back to an explicit parameter binding:

```
var x = f()
z = g() + x
```

When considering side effects, it’s not only the argument expressions that matter. Also significant is the order in which parameters are evaluated relative to other code in the callee. Consider this call to `add2`:

```go
//go:fix inline
func add2(x, y int) int {
	return x + other() + y
}

add2(f(), g())
```

This time, parameters `x` and `y` are used in the same order they are declared, so the substitution `f() + other() + g()` won’t change the order of effects of `f()` and `g()`—but it will change the order of any effects of `other()` and `g()`. Furthermore, if the function body uses a parameter within a loop, substitution might change the cardinality of effects.

The inliner uses a novel [hazard analysis](https://cs.opensource.google/go/x/tools/+/refs/tags/v0.42.0:internal/refactor/inline/inline.go;l=1978;drc=e3a69ffcdbb984f50100e76ebca6ff53cf88de9c) to model the order of effects in each callee function. Nonetheless, its ability to construct the necessary safety proofs is quite limited. For example, if the calls `f()` and `g()` are simple accessors, it would be perfectly safe to call them in either order. Indeed, an optimizing compiler might use its knowledge of the internals of `f` and `g` to safely reorder the two calls. But unlike a compiler, which generates object code that reflects the source at a specific moment, the purpose of the inliner is to make permanent changes to the source, so it can’t take advantage of ephemeral details. As an extreme example, consider this `start` function:

```
func start() { /* TODO: implement */ }
```

An optimizing compiler is free to delete each call to `start()` because it has no effects today, but the inliner is not, because it may become important tomorrow.

<!-- There's a bit of a contradiction here since the hazard analysis uses implementation details du jour. -->

In short, the inliner may produce results that—to the informed eye of a project maintainer—are clearly too conservative. In such cases, the fixed code would benefit stylistically from a little manual cleanup.

### 3. “Fallible” constant expressions

You might imagine (as I once did) that it would always be safe to replace a parameter variable by a constant argument of the same type. Surprisingly, this turns out not to be the case, because some checks previously done at run time would now happen—and fail—at compile time. Consider this call to the `index` function:

```
//go:fix inline
func index(s string, i int) byte {
	return s[i]
}

index("", 0)
```

A naive inliner might replace `s` with `""` and `i` with `0`, resulting in `""[0]`, but this is not actually a legal Go expression because this particular index is out of bounds for this particular string. Because the expression `""[0]` is composed of constants, it is evaluated at compile time, and a program that contains it will not even build. By contrast, the original program would fail only if execution reaches this call to `index`, which presumably in a working program it does not.

Consequently, the inliner must keep track of all expressions and their operands that might become constant during parameter substitution, triggering additional compile-time checks. It builds a [constraint system](https://cs.opensource.google/go/x/tools/+/master:internal/refactor/inline/falcon.go;l=43;drc=1aca71e85510ecc45dddbc335b30b64298c2a31e) and attempts to solve it. Each unsatisfied constraint is resolved by adding an explicit binding for the constrained parameters.

<!--
  The fundamental reason for falcon is that we can’t type-check the result
  since in a “separate analysis” system we don’t have type information
  for all dependencies. See hidden comment within section
  [gofix#synergistic-fixes](gofix#synergistic-fixes).
-->

### 4. Shadowing

Typical argument expressions contain one or more identifiers that refer to symbols (variables, functions, and so on) in the caller’s file. The inliner must make sure that each name in the argument expression would refer to the same symbol after parameter substitution; in other words, none of the caller’s names is *shadowed* in the callee. If this fails, the inliner must again insert parameter bindings, as in this example:

<div class="beforeafter">
<div class="beforeafter-context"><pre>
//go:fix inline
func f(val string) {
	x := 123
	fmt.Println(val, x)
}
</pre></div>
<pre>
x := "hello"
f(x)
</pre>
<div class="beforeafter-arrow"></div>
<pre>
x := "hello"
{
	// another “parameter binding” declaration
	// to read the caller's x before shadowing it
	var val string = x
	x := 123
	fmt.Println(val, x)
}
</pre>
</div>

Conversely, the inliner must also check that each name in the *callee* function body would refer to the same thing when it is spliced into the call site. In other words, none of the callee’s names is shadowed or missing in the caller. For missing names, the inliner may need to insert additional imports.

### 5. Unused variables

When an argument expression has no effects and its corresponding parameter is never used, the expression may be eliminated. However, if the expression contains the last reference to a local variable at the caller, this may cause a compile error because the variable is now unused.

<div class="beforeafter">
<div class="beforeafter-context"><pre>
//go:fix inline
func f(_ int) { print("hello") }
</pre></div>
<pre>
x := 42
f(x)
</pre>
<div class="beforeafter-arrow"></div>
<pre>
x := 42 // error: unused variable: x
print("hello")
</pre>
</div>

So the inliner must account for references to local variables and avoid removing the last one. (Of course it is still possible that two different inliner fixes each remove the *second*-to-last reference to a variable, so the two fixes are valid in isolation but not together; see the discussion of [semantic conflicts](gofix#merging-fixes-and-conflicts) in the previous post. Unfortunately manual cleanup is inevitably required in this case.)

### 6. Defer

In some cases, it is simply impossible to inline away the call.
Consider a call to a function that uses a `defer` statement:
if we were to eliminate the call, the deferred function would execute
when the *caller* function returns, which is too late.
All we can safely do when the callee uses `defer` is to
put the body of the callee in a function literal and immediately call it.
This function literal, `func() { … }()`, delimits the lifetime of the
`defer` statement, as in this example:

<div class="beforeafter">
<div class="beforeafter-context"><pre>
//go:fix inline
func callee() {
	defer f()
	…
}
</pre></div>
<pre>
callee()
</pre>
<div class="beforeafter-arrow"></div>
<pre>
func() {
	defer f()
	…
}()
</pre>
</div>

If you invoke the inliner in gopls, you’ll see that it makes the change shown above and introduces the function literal. This result may be appropriate in an interactive setting, since you are likely to immediately tweak the code (or undo the fix) as you prefer, but it is rarely desirable in a batch tool, so as a matter of policy the analyzer in `go fix` refuses to inline such “literalized” calls.

### An optimizing compiler for “tidiness”

We’ve now seen half a dozen examples of how the inliner handles tricky semantic edge cases correctly.
(Many thanks to Rob Findley, Jonathan Amsterdam, Olena Synenka, and Lasse Folger for insights, discussions, reviews, features, and fixes.)
By putting all of the smarts into the inliner, users can simply apply an “Inline call” refactoring in their IDE or add a `//go:fix inline` directive to their own functions and be confident that the resulting code transformations can be applied with only the most cursory review.

Although we have made good progress toward that goal, we have not yet fully attained it, and it is likely that we never will. Consider a compiler. A sound compiler produces correct output for any input and never miscompiles your code; this is the fundamental expectation that every user should have of their compiler. An *optimizing* compiler produces code carefully chosen for speed without compromising on safety. Similarly, an inliner is a bit like an optimizing compiler whose goal is not speed but *tidiness*: inlining a call must never change the behavior of your program, and ideally it produces code that is maximally neat and tidy. Unfortunately, an optimizing compiler is [provably](https://en.wikipedia.org/wiki/Rice%27s_theorem) never done: showing that two different programs are equivalent is an undecidable problem, and there will always be improvements that an expert knows are safe but the compiler cannot prove. So too with the inliner: there will always be cases where the inliner’s output is too fussy or otherwise stylistically inferior to that of a human expert, and there will always be more “tidiness optimizations” to add.

## Try it out!

We hope this tour of the inliner gives you a sense of some of the challenges involved, and of our priorities and directions in providing sound, self-service code transformation tools. Please try out the inliner, either interactively in your IDE, or through `//go:fix inline` directives and the `go fix` command, and share with us your experiences and any ideas you have for further improvements or new tools.
