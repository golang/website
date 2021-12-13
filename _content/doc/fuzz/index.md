<!--{
  "Title": "Go Fuzzing"
}-->

<!-- Potential pages:
  - What fuzzing is and is not good for
  - Common gotchas / Strategies for inefficient fuzzing executions
  - Rules of a fuzz test (one per package, function signature, etc)
  - Corpus entry file format
  - Commands
    - go clean -fuzzcache
    - tool for converting corpus entries
  - Fuzzing customization (links to docs for -fuzzminimizetime, -fuzz, etc)
  - Technical discussion around how the coordinator/worker work (this may make
    more sense as a blog post?)
-->

Go supports fuzzing in its standard toolchain beginning in Go 1.18.

## Overview

Fuzzing is a type of automated testing which continuously manipulates inputs to
a program to find bugs. Go fuzzing uses coverage guidance to intelligently walk
through the code being fuzzed to find and report failures to the user. Since it
can reach edge cases which humans often miss, fuzz testing can be particularly
valuable for finding security exploits and vulnerabilities.

Below is an example of a [fuzz test](#glos-fuzz-test), highlighting it's main
components.

<img alt="Example code showing the overall fuzz test, with a fuzz target within it. Before the fuzz target is a corpus addition with f.Add, and the parameters of the fuzz target are highlighted as the fuzzing arguments." src="/doc/fuzz/example.png" style="display: block; width: 600px; height: auto;"/>

## Writing and running fuzz tests

### Requirements

Below are rules that fuzz tests must follow.

- A fuzz test must be a function named like `FuzzXxx`, which accepts only a
  `*testing.F`, and has no return value.
- Fuzz tests must be in \*\_test.go files to run.
- A [fuzz target](#glos-fuzz-target) must be a method call to
  <code>[(\*testing.F).Fuzz](https://pkg.go.dev/testing#F.Fuzz)</code> which
  accepts a `*testing.T` as the first parameter, followed by the fuzzing
  arguments. There is no return value.
- There must be exactly one fuzz target per fuzz test.
- All [seed corpus](#glos-seed-corpus) entries must have types which are
  identical to the [fuzzing arguments](#fuzzing-arguments), in the same order.
  This is true for calls to
  <code>[(\*testing.F).Add](https://pkg.go.dev/testing#F.Add)</code> and any
  corpus files in the testdata/fuzz directory of the fuzz test.
- The fuzzing arguments can only be the following types:
  - string, []byte
  - int, int8, int16, int32/rune, int64
  - uint, uint8/byte, uint16, uint32, uint64
  - float32, float64
  - bool

### Suggestions

Below are suggestions that will help you get the most out of fuzzing.

- Fuzzing should be run on a platform that supports coverage instrumentation
  (currently AMD64 and ARM64) so that the corpus can meaningfully grow as it
  runs, and more code can be covered while fuzzing.
- Fuzz targets should be fast and deterministic so the fuzzing engine can work
  efficiently, and new failures and code coverage can be easily reproduced.
- Since the fuzz target is invoked in parallel across multiple workers and in
  nondeterministic order, the state of a fuzz target should not persist past
  the end of each call, and the behavior of a fuzz target should not depend on
  global state.

## Resources

- **Tutorial**:
  - For an introductory tutorial of fuzzing with Go, please see [the blog
    post](https://go.dev/blog/fuzz-beta).
  - More to come soon!
- **Documentation**:
  - The [`testing`](https://pkg.go.dev//testing#hdr-Fuzzing) package docs
    describes the `testing.F` type which is used when writing fuzz tests.
  - The [`cmd/go`](https://pkg.go.dev/cmd/go) package docs describe the flags
    associated with fuzzing.
- **Technical details**:
  - [Design draft](https://golang.org/s/draft-fuzzing-design)
  - [Proposal](https://golang.org/issue/44551)

## Glossary {#glossary}

<a id="glos-corpus-entry"></a>
**corpus entry:** An input in the corpus which can be used while fuzzing. This
can be a specially-formatted file, or a call to
<code>[(\*testing.F).Add](https://pkg.go.dev/testing#F.Add)</code>.

<a id="glos-coverage-guidance"></a>
**coverage guidance:** A method of fuzzing which uses expansions in code
coverage to determine which corpus entries are worth keeping for future use.

<a id="glos-fuzz-target"></a>
**fuzz target:** The function of the fuzz test which is executed for corpus
entries and generated values while fuzzing. It is provided to the fuzz test by
passing the function to
<code>[(\*testing.F).Fuzz](https://pkg.go.dev/testing#F.Fuzz)</code>.

<a id="glos-fuzz-test"></a>
**fuzz test:** A function in a test file of the form `func FuzzXxx(*testing.F)`
which can be used for fuzzing.

<a id="glos-fuzzing"></a>
**fuzzing:** A type of automated testing which continuously manipulates inputs
to a program to find issues such as bugs or
[vulnerabilities](#glos-vulnerability) to which the code may be susceptible.

<a id="glos-fuzzing-arguments"></a>
**fuzzing arguments:** The types which will be passed to the fuzz target, and
mutated by the [mutator](#glos-mutator).

<a id="glos-fuzzing-engine"></a>
**fuzzing engine:** A tool that manages fuzzing, including maintaining the
corpus, invoking the mutator, identifying new coverage, and reporting failures.

<a id="glos-generated-corpus"></a>
**generated corpus:** A corpus which is maintained by the fuzzing engine over
time while fuzzing to keep track of progress. It is stored in `$GOCACHE`/fuzz.

<a id="glos-mutator"></a>
**mutator:** A tool used while fuzzing which randomly manipulates corpus entries
before passing them to a fuzz target.

<a id="glos-package"></a>
**package:** A collection of source files in the same directory that are
compiled together. See the [Packages section](/ref/spec#Packages) in the Go
Language Specification.

<a id="glos-seed-corpus"></a>
**seed corpus:** A user-provided corpus for a fuzz test which can be used to
guide the fuzzing engine. It is composed of the corpus entries provided by f.Add
calls within the fuzz test, and the files in the testdata/fuzz/{FuzzTestName}
directory within the package.

<a id="glos-test-file"></a>
**test file:** A file of the format xxx_test.go that may contain tests, benchmarks, examples and fuzz tests.

<a id="glos-vulnerability"></a>
**vulnerability:** A security-sensitive weakness in code which can be exploited
by an attacker.

## Feedback

If you experience any problems or have an idea for a feature, please [file an
issue](https://github.com/golang/go/issues/new?&labels=fuzz).

For discussion and general feedback about the feature, you can also participate
in the [#fuzzing channel](https://gophers.slack.com/archives/CH5KV1AKE) in
Gophers Slack.
