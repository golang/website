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


Fuzzing is a type of automated testing which continuously manipulates inputs to
a program to find bugs. Since fuzzing can reach edge cases which humans often
miss, fuzz testing can be particularly valuable for finding security exploits
and vulnerabilities.

Go supports fuzzing in its standard toolchain beginning in Go 1.18. It uses
coverage guidance to intelligently walk through the code being fuzzed to find
and report failures to the user.

The [`testing`](https://pkg.go.dev//testing#hdr-Fuzzing) package docs describes
the `testing.F` type which is used when writing fuzz tests, and the
[`cmd/go`](https://pkg.go.dev/cmd/go) package docs describe the flags associated
with fuzzing.

For an introductory tutorial for fuzzing with Go, please see [Blog: Fuzzing is
Beta Ready](https://go.dev/blog/fuzz-beta).

See the [design draft](https://golang.org/s/draft-fuzzing-design) and
[proposal](https://golang.org/issue/44551) for additional details.

## Glossary {#glossary}

<a id="glos-corpus-entry"></a>
**corpus entry:** An input in the corpus which can be used while fuzzing. This
can be a specially-formatted file, or a call to
<code>[(*testing.F).Add](https://pkg.go.dev/testing#F.Add)</code>.

<a id="glos-coverage-guidance"></a>
**coverage guidance:** A method of fuzzing which uses expansions in code
coverage to determine which corpus entries are worth keeping for future use.

<a id="glos-fuzz-target"></a>
**fuzz target:** The function of the fuzz test which is executed for corpus
entries and generated values while fuzzing. It is provided to the fuzz test by
passing the function to
<code>[(*testing.F).Fuzz](https://pkg.go.dev/testing#F.Fuzz)</code>.

<a id="glos-fuzz-test"></a>
**fuzz test:** A function in a test file of the form `func FuzzXxx(*testing.F)`
which can be used for fuzzing.

<a id="glos-fuzzing"></a>
**fuzzing:** A type of automated testing which continuously manipulates inputs
to a program to find issues such as bugs or
[vulnerabilities](#glos-vulnerability) to which the code may be susceptible.

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
calls within the fuzz test, and the files in the testdata/fuzz/{FuzzTargetName}
directory within the package.

<a id="glos-test-file"></a>
**test file:** A file of the format xxx_test.go that may contain tests, benchmarks, examples and fuzz tests.

<a id="glos-vulnerability"></a>
**vulnerability:** A security-sensitive weakness in code which can be exploited
by an attacker.