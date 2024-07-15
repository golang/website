---
template: false
title: Go 1.16 Release Notes
---

<!--
NOTE: In this document and others in this directory, the convention is to
set fixed-width phrases with non-fixed-width spaces, as in
`hello` `world`.
Do not send CLs removing the interior tags from such phrases.
-->

<style>
  main ul li { margin: 0.5em 0; }
</style>

## Introduction to Go 1.16 {#introduction}

The latest Go release, version 1.16, arrives six months after [Go 1.15](/doc/go1.15).
Most of its changes are in the implementation of the toolchain, runtime, and libraries.
As always, the release maintains the Go 1 [promise of compatibility](/doc/go1compat.html).
We expect almost all Go programs to continue to compile and run as before.

## Changes to the language {#language}

There are no changes to the language.

## Ports {#ports}

### Darwin and iOS {#darwin}

<!-- golang.org/issue/38485, golang.org/issue/41385, CL 266373, more CLs -->
Go 1.16 adds support of 64-bit ARM architecture on macOS (also known as
Apple Silicon) with `GOOS=darwin`, `GOARCH=arm64`.
Like the `darwin/amd64` port, the `darwin/arm64`
port supports cgo, internal and external linking, `c-archive`,
`c-shared`, and `pie` build modes, and the race
detector.

<!-- CL 254740 -->
The iOS port, which was previously `darwin/arm64`, has
been renamed to `ios/arm64`. `GOOS=ios`
implies the
`darwin` build tag, just as `GOOS=android`
implies the `linux` build tag. This change should be
transparent to anyone using gomobile to build iOS apps.

The introduction of `GOOS=ios` means that file names
like `x_ios.go` will now only be built for
`GOOS=ios`; see
[`go`
`help` `buildconstraint`](/cmd/go/#hdr-Build_constraints) for details.
Existing packages that use file names of this form will have to
rename the files.

<!-- golang.org/issue/42100, CL 263798 -->
Go 1.16 adds an `ios/amd64` port, which targets the iOS
simulator running on AMD64-based macOS. Previously this was
unofficially supported through `darwin/amd64` with
the `ios` build tag set. See also
[`misc/ios/README`](/misc/ios/README) for
details about how to build programs for iOS and iOS simulator.

<!-- golang.org/issue/23011 -->
Go 1.16 is the last release that will run on macOS 10.12 Sierra.
Go 1.17 will require macOS 10.13 High Sierra or later.

### NetBSD {#netbsd}

<!-- golang.org/issue/30824 -->
Go now supports the 64-bit ARM architecture on NetBSD (the
`netbsd/arm64` port).

### OpenBSD {#openbsd}

<!-- golang.org/issue/40995 -->
Go now supports the MIPS64 architecture on OpenBSD
(the `openbsd/mips64` port). This port does not yet
support cgo.

<!-- golang.org/issue/36435, many CLs -->
On the 64-bit x86 and 64-bit ARM architectures on OpenBSD (the
`openbsd/amd64` and `openbsd/arm64` ports), system
calls are now made through `libc`, instead of directly using
the `SYSCALL`/`SVC` instruction. This ensures
forward-compatibility with future versions of OpenBSD. In particular,
OpenBSD 6.9 onwards will require system calls to be made through
`libc` for non-static Go binaries.

### 386 {#386}

<!-- golang.org/issue/40255, golang.org/issue/41848, CL 258957, and CL 260017 -->
As [announced](go1.15#386) in the Go 1.15 release notes,
Go 1.16 drops support for x87 mode compilation (`GO386=387`).
Support for non-SSE2 processors is now available using soft float
mode (`GO386=softfloat`).
Users running on non-SSE2 processors should replace `GO386=387`
with `GO386=softfloat`.

### RISC-V {#riscv}

<!-- golang.org/issue/36641, CL 267317 -->
The `linux/riscv64` port now supports cgo and
`-buildmode=pie`. This release also includes performance
optimizations and code generation improvements for RISC-V.

## Tools {#tools}

### Go command {#go-command}

#### Modules {#modules}

<!-- golang.org/issue/41330 -->
Module-aware mode is enabled by default, regardless of whether a
`go.mod` file is present in the current working directory or a
parent directory. More precisely, the `GO111MODULE` environment
variable now defaults to `on`. To switch to the previous behavior,
set `GO111MODULE` to `auto`.

<!-- golang.org/issue/40728 -->
Build commands like `go` `build` and `go`
`test` no longer modify `go.mod` and `go.sum`
by default. Instead, they report an error if a module requirement or checksum
needs to be added or updated (as if the `-mod=readonly` flag were
used). Module requirements and sums may be adjusted with `go`
`mod` `tidy` or `go` `get`.

<!-- golang.org/issue/40276 -->
`go` `install` now accepts arguments with
version suffixes (for example, `go` `install`
`example.com/cmd@v1.0.0`). This causes `go`
`install` to build and install packages in module-aware mode,
ignoring the `go.mod` file in the current directory or any parent
directory, if there is one. This is useful for installing executables without
affecting the dependencies of the main module.

<!-- golang.org/issue/40276 -->
`go` `install`, with or without a version suffix (as
described above), is now the recommended way to build and install packages in
module mode. `go` `get` should be used with the
`-d` flag to adjust the current module's dependencies without
building packages, and use of `go` `get` to build and
install packages is deprecated. In a future release, the `-d` flag
will always be enabled.

<!-- golang.org/issue/24031 -->
`retract` directives may now be used in a `go.mod` file
to indicate that certain published versions of the module should not be used
by other modules. A module author may retract a version after a severe problem
is discovered or if the version was published unintentionally.

<!-- golang.org/issue/26603 -->
The `go` `mod` `vendor`
and `go` `mod` `tidy` subcommands now accept
the `-e` flag, which instructs them to proceed despite errors in
resolving missing packages.

<!-- golang.org/issue/36465 -->
The `go` command now ignores requirements on module versions
excluded by `exclude` directives in the main module. Previously,
the `go` command used the next version higher than an excluded
version, but that version could change over time, resulting in
non-reproducible builds.

<!-- golang.org/issue/43052, golang.org/issue/43985 -->
In module mode, the `go` command now disallows import paths that
include non-ASCII characters or path elements with a leading dot character
(`.`). Module paths with these characters were already disallowed
(see [Module paths and versions](/ref/mod#go-mod-file-ident)),
so this change affects only paths within module subdirectories.

#### Embedding Files {#embed}

The `go` command now supports including
static files and file trees as part of the final executable,
using the new `//go:embed` directive.
See the documentation for the new
[`embed`](/pkg/embed/)
package for details.

#### `go` `test` {#go-test}

<!-- golang.org/issue/29062 -->
When using `go` `test`, a test that
calls `os.Exit(0)` during execution of a test function
will now be considered to fail.
This will help catch cases in which a test calls code that calls
`os.Exit(0)` and thereby stops running all future tests.
If a `TestMain` function calls `os.Exit(0)`
that is still considered to be a passing test.

<!-- golang.org/issue/39484 -->
`go` `test` reports an error when the `-c`
or `-i` flags are used together with unknown flags. Normally,
unknown flags are passed to tests, but when `-c` or `-i`
are used, tests are not run.

#### `go` `get` {#go-get}

<!-- golang.org/issue/37519 -->
The `go` `get` `-insecure` flag is
deprecated and will be removed in a future version. This flag permits
fetching from repositories and resolving custom domains using insecure
schemes such as HTTP, and also bypasses module sum validation using the
checksum database. To permit the use of insecure schemes, use the
`GOINSECURE` environment variable instead. To bypass module
sum validation, use `GOPRIVATE` or `GONOSUMDB`.
See `go` `help` `environment` for details.

<!-- golang.org/cl/263267 -->
`go` `get` `example.com/mod@patch` now
requires that some version of `example.com/mod` already be
required by the main module.
(However, `go` `get` `-u=patch` continues
to patch even newly-added dependencies.)

#### `GOVCS` environment variable {#govcs}

<!-- golang.org/issue/266420 -->
`GOVCS` is a new environment variable that limits which version
control tools the `go` command may use to download source code.
This mitigates security issues with tools that are typically used in trusted,
authenticated environments. By default, `git` and `hg`
may be used to download code from any repository. `svn`,
`bzr`, and `fossil` may only be used to download code
from repositories with module paths or package paths matching patterns in
the `GOPRIVATE` environment variable. See
[`go`
`help` `vcs`](/cmd/go/#hdr-Controlling_version_control_with_GOVCS) for details.

#### The `all` pattern {#all-pattern}

<!-- golang.org/cl/240623 -->
When the main module's `go.mod` file
declares `go` `1.16` or higher, the `all`
package pattern now matches only those packages that are transitively imported
by a package or test found in the main module. (Packages imported by _tests
of_ packages imported by the main module are no longer included.) This is
the same set of packages retained
by `go` `mod` `vendor` since Go 1.11.

#### The `-toolexec` build flag {#toolexec}

<!-- golang.org/cl/263357 -->
When the `-toolexec` build flag is specified to use a program when
invoking toolchain programs like compile or asm, the environment variable
`TOOLEXEC_IMPORTPATH` is now set to the import path of the package
being built.

#### The `-i` build flag {#i-flag}

<!-- golang.org/issue/41696 -->
The `-i` flag accepted by `go` `build`,
`go` `install`, and `go` `test` is
now deprecated. The `-i` flag instructs the `go` command
to install packages imported by packages named on the command line. Since
the build cache was introduced in Go 1.10, the `-i` flag no longer
has a significant effect on build times, and it causes errors when the install
directory is not writable.

#### The `list` command {#list-buildid}

<!-- golang.org/cl/263542 -->
When the `-export` flag is specified, the `BuildID`
field is now set to the build ID of the compiled package. This is equivalent
to running `go` `tool` `buildid` on
`go` `list` `-exported` `-f` `{{.Export}}`,
but without the extra step.

#### The `-overlay` flag {#overlay-flag}

<!-- golang.org/issue/39958 -->
The `-overlay` flag specifies a JSON configuration file containing
a set of file path replacements. The `-overlay` flag may be used
with all build commands and `go` `mod` subcommands.
It is primarily intended to be used by editor tooling such as gopls to
understand the effects of unsaved changes to source files. The config file
maps actual file paths to replacement file paths and the `go`
command and its builds will run as if the actual file paths exist with the
contents given by the replacement file paths, or don't exist if the replacement
file paths are empty.

### Cgo {#cgo}

<!-- CL 252378 -->
The [cgo](/cmd/cgo) tool will no longer try to translate
C struct bitfields into Go struct fields, even if their size can be
represented in Go. The order in which C bitfields appear in memory
is implementation dependent, so in some cases the cgo tool produced
results that were silently incorrect.

### Vet {#vet}

#### New warning for invalid testing.T use in goroutines {#vet-testing-T}

<!-- CL 235677 -->
The vet tool now warns about invalid calls to the `testing.T`
method `Fatal` from within a goroutine created during the test.
This also warns on calls to `Fatalf`, `FailNow`, and
`Skip{,f,Now}` methods on `testing.T` tests or
`testing.B` benchmarks.

Calls to these methods stop the execution of the created goroutine and not
the `Test*` or `Benchmark*` function. So these are
[required](/pkg/testing/#T.FailNow) to be called by the goroutine
running the test or benchmark function. For example:

	func TestFoo(t *testing.T) {
	    go func() {
	        if condition() {
	            t.Fatal("oops") // This exits the inner func instead of TestFoo.
	        }
	        ...
	    }()
	}

Code calling `t.Fatal` (or a similar method) from a created
goroutine should be rewritten to signal the test failure using
`t.Error` and exit the goroutine early using an alternative
method, such as using a `return` statement. The previous example
could be rewritten as:

	func TestFoo(t *testing.T) {
	    go func() {
	        if condition() {
	            t.Error("oops")
	            return
	        }
	        ...
	    }()
	}

#### New warning for frame pointer {#vet-frame-pointer}

<!-- CL 248686, CL 276372 -->
The vet tool now warns about amd64 assembly that clobbers the BP
register (the frame pointer) without saving and restoring it,
contrary to the calling convention. Code that doesn't preserve the
BP register must be modified to either not use BP at all or preserve
BP by saving and restoring it. An easy way to preserve BP is to set
the frame size to a nonzero value, which causes the generated
prologue and epilogue to preserve the BP register for you.
See [CL 248260](/cl/248260) for example
fixes.

#### New warning for asn1.Unmarshal {#vet-asn1-unmarshal}

<!-- CL 243397 -->
The vet tool now warns about incorrectly passing a non-pointer or nil argument to
[`asn1.Unmarshal`](/pkg/encoding/asn1/#Unmarshal).
This is like the existing checks for
[`encoding/json.Unmarshal`](/pkg/encoding/json/#Unmarshal)
and [`encoding/xml.Unmarshal`](/pkg/encoding/xml/#Unmarshal).

## Runtime {#runtime}

The new [`runtime/metrics`](/pkg/runtime/metrics/) package
introduces a stable interface for reading
implementation-defined metrics from the Go runtime.
It supersedes existing functions like
[`runtime.ReadMemStats`](/pkg/runtime/#ReadMemStats)
and
[`debug.GCStats`](/pkg/runtime/debug/#GCStats)
and is significantly more general and efficient.
See the package documentation for more details.

<!-- CL 254659 -->
Setting the `GODEBUG` environment variable
to `inittrace=1` now causes the runtime to emit a single
line to standard error for each package `init`,
summarizing its execution time and memory allocation. This trace can
be used to find bottlenecks or regressions in Go startup
performance.
The [`GODEBUG`
documentation](/pkg/runtime/#hdr-Environment_Variables) describes the format.

<!-- CL 267100 -->
On Linux, the runtime now defaults to releasing memory to the
operating system promptly (using `MADV_DONTNEED`), rather
than lazily when the operating system is under memory pressure
(using `MADV_FREE`). This means process-level memory
statistics like RSS will more accurately reflect the amount of
physical memory being used by Go processes. Systems that are
currently using `GODEBUG=madvdontneed=1` to improve
memory monitoring behavior no longer need to set this environment
variable.

<!-- CL 220419, CL 271987 -->
Go 1.16 fixes a discrepancy between the race detector and
the [Go memory model](/ref/mem). The race detector now
more precisely follows the channel synchronization rules of the
memory model. As a result, the detector may now report races it
previously missed.

## Compiler {#compiler}

<!-- CL 256459, CL 264837, CL 266203, CL 256460 -->
The compiler can now inline functions with
non-labeled `for` loops, method values, and type
switches. The inliner can also detect more indirect calls where
inlining is possible.

## Linker {#linker}

<!-- CL 248197 -->
This release includes additional improvements to the Go linker,
reducing linker resource usage (both time and memory) and improving
code robustness/maintainability. These changes form the second half
of a two-release project to
[modernize the Go
linker](/s/better-linker).

The linker changes in 1.16 extend the 1.15 improvements to all
supported architecture/OS combinations (the 1.15 performance improvements
were primarily focused on `ELF`-based OSes and
`amd64` architectures). For a representative set of
large Go programs, linking is 20-25% faster than 1.15 and requires
5-15% less memory on average for `linux/amd64`, with larger
improvements for other architectures and OSes. Most binaries are
also smaller as a result of more aggressive symbol pruning.

<!-- CL 255259 -->
On Windows, `go build -buildmode=c-shared` now generates Windows
ASLR DLLs by default. ASLR can be disabled with `--ldflags=-aslr=false`.

## Standard library {#library}

### Embedded Files {#library-embed}

The new [`embed`](/pkg/embed/) package
provides access to files embedded in the program during compilation
using the new [`//go:embed` directive](#embed).

### File Systems {#fs}

The new [`io/fs`](/pkg/io/fs/) package
defines the [`fs.FS`](/pkg/io/fs/#FS) interface,
an abstraction for read-only trees of files.
The standard library packages have been adapted to make use
of the interface as appropriate.

On the producer side of the interface,
the new [`embed.FS`](/pkg/embed/#FS) type
implements `fs.FS`, as does
[`zip.Reader`](/pkg/archive/zip/#Reader).
The new [`os.DirFS`](/pkg/os/#DirFS) function
provides an implementation of `fs.FS` backed by a tree
of operating system files.

On the consumer side,
the new [`http.FS`](/pkg/net/http/#FS)
function converts an `fs.FS` to an
[`http.FileSystem`](/pkg/net/http/#FileSystem).
Also, the [`html/template`](/pkg/html/template/)
and [`text/template`](/pkg/text/template/)
packages’ [`ParseFS`](/pkg/html/template/#ParseFS)
functions and methods read templates from an `fs.FS`.

For testing code that implements `fs.FS`,
the new [`testing/fstest`](/pkg/testing/fstest/)
package provides a [`TestFS`](/pkg/testing/fstest/#TestFS)
function that checks for and reports common mistakes.
It also provides a simple in-memory file system implementation,
[`MapFS`](/pkg/testing/fstest/#MapFS),
which can be useful for testing code that accepts `fs.FS`
implementations.

### Deprecation of io/ioutil {#ioutil}

The [`io/ioutil`](/pkg/io/ioutil/) package has
turned out to be a poorly defined and hard to understand collection
of things. All functionality provided by the package has been moved
to other packages. The `io/ioutil` package remains and
will continue to work as before, but we encourage new code to use
the new definitions in the [`io`](/pkg/io/) and
[`os`](/pkg/os/) packages.
Here is a list of the new locations of the names exported
by `io/ioutil`:

  - [`Discard`](/pkg/io/ioutil/#Discard)
    => [`io.Discard`](/pkg/io/#Discard)
  - [`NopCloser`](/pkg/io/ioutil/#NopCloser)
    => [`io.NopCloser`](/pkg/io/#NopCloser)
  - [`ReadAll`](/pkg/io/ioutil/#ReadAll)
    => [`io.ReadAll`](/pkg/io/#ReadAll)
  - [`ReadDir`](/pkg/io/ioutil/#ReadDir)
    => [`os.ReadDir`](/pkg/os/#ReadDir)
    (note: returns a slice of
    [`os.DirEntry`](/pkg/os/#DirEntry)
    rather than a slice of
    [`fs.FileInfo`](/pkg/io/fs/#FileInfo))
  - [`ReadFile`](/pkg/io/ioutil/#ReadFile)
    => [`os.ReadFile`](/pkg/os/#ReadFile)
  - [`TempDir`](/pkg/io/ioutil/#TempDir)
    => [`os.MkdirTemp`](/pkg/os/#MkdirTemp)
  - [`TempFile`](/pkg/io/ioutil/#TempFile)
    => [`os.CreateTemp`](/pkg/os/#CreateTemp)
  - [`WriteFile`](/pkg/io/ioutil/#WriteFile)
    => [`os.WriteFile`](/pkg/os/#WriteFile)


### Minor changes to the library {#minor_library_changes}

As always, there are various minor changes and updates to the library,
made with the Go 1 [promise of compatibility](/doc/go1compat)
in mind.

#### [archive/zip](/pkg/archive/zip/)

<!-- CL 243937 -->
The new [`Reader.Open`](/pkg/archive/zip/#Reader.Open)
method implements the [`fs.FS`](/pkg/io/fs/#FS)
interface.

#### [crypto/dsa](/pkg/crypto/dsa/)

<!-- CL 257939 -->
The [`crypto/dsa`](/pkg/crypto/dsa/) package is now deprecated.
See [issue #40337](/issue/40337).

<!-- crypto/dsa -->

#### [crypto/hmac](/pkg/crypto/hmac/)

<!-- CL 261960 -->
[`New`](/pkg/crypto/hmac/#New) will now panic if
separate calls to the hash generation function fail to return new values.
Previously, the behavior was undefined and invalid outputs were sometimes
generated.

<!-- crypto/hmac -->

#### [crypto/tls](/pkg/crypto/tls/)

<!-- CL 256897 -->
I/O operations on closing or closed TLS connections can now be detected
using the new [`net.ErrClosed`](/pkg/net/#ErrClosed)
error. A typical use would be `errors.Is(err, net.ErrClosed)`.

<!-- CL 266037 -->
A default write deadline is now set in
[`Conn.Close`](/pkg/crypto/tls/#Conn.Close)
before sending the "close notify" alert, in order to prevent blocking
indefinitely.

<!-- CL 239748 -->
Clients now return a handshake error if the server selects
[
an ALPN protocol](/pkg/crypto/tls/#ConnectionState.NegotiatedProtocol) that was not in
[
the list advertised by the client](/pkg/crypto/tls/#Config.NextProtos).

<!-- CL 262857 -->
Servers will now prefer other available AEAD cipher suites (such as ChaCha20Poly1305)
over AES-GCM cipher suites if either the client or server doesn't have AES hardware
support, unless both [
`Config.PreferServerCipherSuites`](/pkg/crypto/tls/#Config.PreferServerCipherSuites)
and [`Config.CipherSuites`](/pkg/crypto/tls/#Config.CipherSuites)
are set. The client is assumed not to have AES hardware support if it does
not signal a preference for AES-GCM cipher suites.

<!-- CL 246637 -->
[`Config.Clone`](/pkg/crypto/tls/#Config.Clone) now
returns nil if the receiver is nil, rather than panicking.

<!-- crypto/tls -->

#### [crypto/x509](/pkg/crypto/x509/)

The `GODEBUG=x509ignoreCN=0` flag will be removed in Go 1.17.
It enables the legacy behavior of treating the `CommonName`
field on X.509 certificates as a host name when no Subject Alternative
Names are present.

<!-- CL 235078 -->
[`ParseCertificate`](/pkg/crypto/x509/#ParseCertificate) and
[`CreateCertificate`](/pkg/crypto/x509/#CreateCertificate)
now enforce string encoding restrictions for the `DNSNames`,
`EmailAddresses`, and `URIs` fields. These fields
can only contain strings with characters within the ASCII range.

<!-- CL 259697 -->
[`CreateCertificate`](/pkg/crypto/x509/#CreateCertificate)
now verifies the generated certificate's signature using the signer's
public key. If the signature is invalid, an error is returned, instead of
a malformed certificate.

<!-- CL 257939 -->
DSA signature verification is no longer supported. Note that DSA signature
generation was never supported.
See [issue #40337](/issue/40337).

<!-- CL 257257 -->
On Windows, [`Certificate.Verify`](/pkg/crypto/x509/#Certificate.Verify)
will now return all certificate chains that are built by the platform
certificate verifier, instead of just the highest ranked chain.

<!-- CL 262343 -->
The new [`SystemRootsError.Unwrap`](/pkg/crypto/x509/#SystemRootsError.Unwrap)
method allows accessing the [`Err`](/pkg/crypto/x509/#SystemRootsError.Err)
field through the [`errors`](/pkg/errors) package functions.

<!-- CL 230025 -->
On Unix systems, the `crypto/x509` package is now more
efficient in how it stores its copy of the system cert pool.
Programs that use only a small number of roots will use around a
half megabyte less memory.

<!-- crypto/x509 -->

#### [debug/elf](/pkg/debug/elf/)

<!-- CL 255138 -->
More [`DT`](/pkg/debug/elf/#DT_NULL)
and [`PT`](/pkg/debug/elf/#PT_NULL)
constants have been added.

<!-- debug/elf -->

#### [encoding/asn1](/pkg/encoding/asn1)

<!-- CL 255881 -->
[`Unmarshal`](/pkg/encoding/asn1/#Unmarshal) and
[`UnmarshalWithParams`](/pkg/encoding/asn1/#UnmarshalWithParams)
now return an error instead of panicking when the argument is not
a pointer or is nil. This change matches the behavior of other
encoding packages such as [`encoding/json`](/pkg/encoding/json).

#### [encoding/json](/pkg/encoding/json/)

<!-- CL 234818 -->
The `json` struct field tags understood by
[`Marshal`](/pkg/encoding/json/#Marshal),
[`Unmarshal`](/pkg/encoding/json/#Unmarshal),
and related functionality now permit semicolon characters within
a JSON object name for a Go struct field.

<!-- encoding/json -->

#### [encoding/xml](/pkg/encoding/xml/)

<!-- CL 264024 -->
The encoder has always taken care to avoid using namespace prefixes
beginning with `xml`, which are reserved by the XML
specification.
Now, following the specification more closely, that check is
case-insensitive, so that prefixes beginning
with `XML`, `XmL`, and so on are also
avoided.

<!-- encoding/xml -->

#### [flag](/pkg/flag/)

<!-- CL 240014 -->
The new [`Func`](/pkg/flag/#Func) function
allows registering a flag implemented by calling a function,
as a lighter-weight alternative to implementing the
[`Value`](/pkg/flag/#Value) interface.

<!-- flag -->

#### [go/build](/pkg/go/build/)

<!-- CL 243941, CL 283636 -->
The [`Package`](/pkg/go/build/#Package)
struct has new fields that report information
about `//go:embed` directives in the package:
[`EmbedPatterns`](/pkg/go/build/#Package.EmbedPatterns),
[`EmbedPatternPos`](/pkg/go/build/#Package.EmbedPatternPos),
[`TestEmbedPatterns`](/pkg/go/build/#Package.TestEmbedPatterns),
[`TestEmbedPatternPos`](/pkg/go/build/#Package.TestEmbedPatternPos),
[`XTestEmbedPatterns`](/pkg/go/build/#Package.XTestEmbedPatterns),
[`XTestEmbedPatternPos`](/pkg/go/build/#Package.XTestEmbedPatternPos).

<!-- CL 240551 -->
The [`Package`](/pkg/go/build/#Package) field
[`IgnoredGoFiles`](/pkg/go/build/#Package.IgnoredGoFiles)
will no longer include files that start with "\_" or ".",
as those files are always ignored.
`IgnoredGoFiles` is for files ignored because of
build constraints.

<!-- CL 240551 -->
The new [`Package`](/pkg/go/build/#Package)
field [`IgnoredOtherFiles`](/pkg/go/build/#Package.IgnoredOtherFiles)
has a list of non-Go files ignored because of build constraints.

<!-- go/build -->

#### [go/build/constraint](/pkg/go/build/constraint/)

<!-- CL 240604 -->
The new
[`go/build/constraint`](/pkg/go/build/constraint/)
package parses build constraint lines, both the original
`// +build` syntax and the `//go:build`
syntax that will be introduced in Go 1.17.
This package exists so that tools built with Go 1.16 will be able
to process Go 1.17 source code.
See [https://golang.org/design/draft-gobuild](/design/draft-gobuild)
for details about the build constraint syntaxes and the planned
transition to the `//go:build` syntax.
Note that `//go:build` lines are **not** supported
in Go 1.16 and should not be introduced into Go programs yet.

<!-- go/build/constraint -->

#### [html/template](/pkg/html/template/)

<!-- CL 243938 -->
The new [`template.ParseFS`](/pkg/html/template/#ParseFS)
function and [`template.Template.ParseFS`](/pkg/html/template/#Template.ParseFS)
method are like [`template.ParseGlob`](/pkg/html/template/#ParseGlob)
and [`template.Template.ParseGlob`](/pkg/html/template/#Template.ParseGlob),
but read the templates from an [`fs.FS`](/pkg/io/fs/#FS).

<!-- html/template -->

#### [io](/pkg/io/)

<!-- CL 261577 -->
The package now defines a
[`ReadSeekCloser`](/pkg/io/#ReadSeekCloser) interface.

<!-- CL 263141 -->
The package now defines
[`Discard`](/pkg/io/#Discard),
[`NopCloser`](/pkg/io/#NopCloser), and
[`ReadAll`](/pkg/io/#ReadAll),
to be used instead of the same names in the
[`io/ioutil`](/pkg/io/ioutil/) package.

<!-- io -->

#### [log](/pkg/log/)

<!-- CL 264460 -->
The new [`Default`](/pkg/log/#Default) function
provides access to the default [`Logger`](/pkg/log/#Logger).

<!-- log -->

#### [log/syslog](/pkg/log/syslog/)

<!-- CL 264297 -->
The [`Writer`](/pkg/log/syslog/#Writer)
now uses the local message format
(omitting the host name and using a shorter time stamp)
when logging to custom Unix domain sockets,
matching the format already used for the default log socket.

<!-- log/syslog -->

#### [mime/multipart](/pkg/mime/multipart/)

<!-- CL 247477 -->
The [`Reader`](/pkg/mime/multipart/#Reader)'s
[`ReadForm`](/pkg/mime/multipart/#Reader.ReadForm)
method no longer rejects form data
when passed the maximum int64 value as a limit.

<!-- mime/multipart -->

#### [net](/pkg/net/)

<!-- CL 250357 -->
The case of I/O on a closed network connection, or I/O on a network
connection that is closed before any of the I/O completes, can now
be detected using the new [`ErrClosed`](/pkg/net/#ErrClosed)
error. A typical use would be `errors.Is(err, net.ErrClosed)`.
In earlier releases the only way to reliably detect this case was to
match the string returned by the `Error` method
with `"use of closed network connection"`.

<!-- CL 255898 -->
In previous Go releases the default TCP listener backlog size on Linux systems,
set by `/proc/sys/net/core/somaxconn`, was limited to a maximum of `65535`.
On Linux kernel version 4.1 and above, the maximum is now `4294967295`.

<!-- CL 238629 -->
On Linux, host name lookups no longer use DNS before checking
`/etc/hosts` when `/etc/nsswitch.conf`
is missing; this is common on musl-based systems and makes
Go programs match the behavior of C programs on those systems.

<!-- net -->

#### [net/http](/pkg/net/http/)

<!-- CL 233637 -->
In the [`net/http`](/pkg/net/http/) package, the
behavior of [`StripPrefix`](/pkg/net/http/#StripPrefix)
has been changed to strip the prefix from the request URL's
`RawPath` field in addition to its `Path` field.
In past releases, only the `Path` field was trimmed, and so if the
request URL contained any escaped characters the URL would be modified to
have mismatched `Path` and `RawPath` fields.
In Go 1.16, `StripPrefix` trims both fields.
If there are escaped characters in the prefix part of the request URL the
handler serves a 404 instead of its previous behavior of invoking the
underlying handler with a mismatched `Path`/`RawPath` pair.

<!-- CL 252497 -->
The [`net/http`](/pkg/net/http/) package now rejects HTTP range requests
of the form `"Range": "bytes=--N"` where `"-N"` is a negative suffix length, for
example `"Range": "bytes=--2"`. It now replies with a `416 "Range Not Satisfiable"` response.

<!-- CL 256498, golang.org/issue/36990 -->
Cookies set with [`SameSiteDefaultMode`](/pkg/net/http/#SameSiteDefaultMode)
now behave according to the current spec (no attribute is set) instead of
generating a SameSite key without a value.

<!-- CL 250039 -->
The [`Client`](/pkg/net/http/#Client) now sends
an explicit `Content-Length:` `0`
header in `PATCH` requests with empty bodies,
matching the existing behavior of `POST` and `PUT`.

<!-- CL 249440 -->
The [`ProxyFromEnvironment`](/pkg/net/http/#ProxyFromEnvironment)
function no longer returns the setting of the `HTTP_PROXY`
environment variable for `https://` URLs when
`HTTPS_PROXY` is unset.

<!-- 259917 -->
The [`Transport`](/pkg/net/http/#Transport)
type has a new field
[`GetProxyConnectHeader`](/pkg/net/http/#Transport.GetProxyConnectHeader)
which may be set to a function that returns headers to send to a
proxy during a `CONNECT` request.
In effect `GetProxyConnectHeader` is a dynamic
version of the existing field
[`ProxyConnectHeader`](/pkg/net/http/#Transport.ProxyConnectHeader);
if `GetProxyConnectHeader` is not `nil`,
then `ProxyConnectHeader` is ignored.

<!-- CL 243939 -->
The new [`http.FS`](/pkg/net/http/#FS)
function converts an [`fs.FS`](/pkg/io/fs/#FS)
to an [`http.FileSystem`](/pkg/net/http/#FileSystem).

<!-- net/http -->

#### [net/http/httputil](/pkg/net/http/httputil/)

<!-- CL 260637 -->
[`ReverseProxy`](/pkg/net/http/httputil/#ReverseProxy)
now flushes buffered data more aggressively when proxying
streamed responses with unknown body lengths.

<!-- net/http/httputil -->

#### [net/smtp](/pkg/net/smtp/)

<!-- CL 247257 -->
The [`Client`](/pkg/net/smtp/#Client)'s
[`Mail`](/pkg/net/smtp/#Client.Mail)
method now sends the `SMTPUTF8` directive to
servers that support it, signaling that addresses are encoded in UTF-8.

<!-- net/smtp -->

#### [os](/pkg/os/)

<!-- CL 242998 -->
[`Process.Signal`](/pkg/os/#Process.Signal) now
returns [`ErrProcessDone`](/pkg/os/#ErrProcessDone)
instead of the unexported `errFinished` when the process has
already finished.

<!-- CL 261540 -->
The package defines a new type
[`DirEntry`](/pkg/os/#DirEntry)
as an alias for [`fs.DirEntry`](/pkg/io/fs/#DirEntry).
The new [`ReadDir`](/pkg/os/#ReadDir)
function and the new
[`File.ReadDir`](/pkg/os/#File.ReadDir)
method can be used to read the contents of a directory into a
slice of [`DirEntry`](/pkg/os/#DirEntry).
The [`File.Readdir`](/pkg/os/#File.Readdir)
method (note the lower case `d` in `dir`)
still exists, returning a slice of
[`FileInfo`](/pkg/os/#FileInfo), but for
most programs it will be more efficient to switch to
[`File.ReadDir`](/pkg/os/#File.ReadDir).

<!-- CL 263141 -->
The package now defines
[`CreateTemp`](/pkg/os/#CreateTemp),
[`MkdirTemp`](/pkg/os/#MkdirTemp),
[`ReadFile`](/pkg/os/#ReadFile), and
[`WriteFile`](/pkg/os/#WriteFile),
to be used instead of functions defined in the
[`io/ioutil`](/pkg/io/ioutil/) package.

<!-- CL 243906 -->
The types [`FileInfo`](/pkg/os/#FileInfo),
[`FileMode`](/pkg/os/#FileMode), and
[`PathError`](/pkg/os/#PathError)
are now aliases for types of the same name in the
[`io/fs`](/pkg/io/fs/) package.
Function signatures in the [`os`](/pkg/os/)
package have been updated to refer to the names in the
[`io/fs`](/pkg/io/fs/) package.
This should not affect any existing code.

<!-- CL 243911 -->
The new [`DirFS`](/pkg/os/#DirFS) function
provides an implementation of
[`fs.FS`](/pkg/io/fs/#FS) backed by a tree
of operating system files.

<!-- os -->

#### [os/signal](/pkg/os/signal/)

<!-- CL 219640 -->
The new
[`NotifyContext`](/pkg/os/signal/#NotifyContext)
function allows creating contexts that are canceled upon arrival of
specific signals.

<!-- os/signal -->

#### [path](/pkg/path/)

<!-- CL 264397, golang.org/issues/28614 -->
The [`Match`](/pkg/path/#Match) function now
returns an error if the unmatched part of the pattern has a
syntax error. Previously, the function returned early on a failed
match, and thus did not report any later syntax error in the
pattern.

<!-- path -->

#### [path/filepath](/pkg/path/filepath/)

<!-- CL 267887 -->
The new function
[`WalkDir`](/pkg/path/filepath/#WalkDir)
is similar to
[`Walk`](/pkg/path/filepath/#Walk),
but is typically more efficient.
The function passed to `WalkDir` receives a
[`fs.DirEntry`](/pkg/io/fs/#DirEntry)
instead of a
[`fs.FileInfo`](/pkg/io/fs/#FileInfo).
(To clarify for those who recall the `Walk` function
as taking an [`os.FileInfo`](/pkg/os/#FileInfo),
`os.FileInfo` is now an alias for `fs.FileInfo`.)

<!-- CL 264397, golang.org/issues/28614 -->
The [`Match`](/pkg/path/filepath#Match) and
[`Glob`](/pkg/path/filepath#Glob) functions now
return an error if the unmatched part of the pattern has a
syntax error. Previously, the functions returned early on a failed
match, and thus did not report any later syntax error in the
pattern.

<!-- path/filepath -->

#### [reflect](/pkg/reflect/)

<!-- CL 192331 -->
The Zero function has been optimized to avoid allocations. Code
which incorrectly compares the returned Value to another Value
using == or DeepEqual may get different results than those
obtained in previous Go versions. The documentation
for [`reflect.Value`](/pkg/reflect#Value)
describes how to compare two `Value`s correctly.

<!-- reflect -->

#### [runtime/debug](/pkg/runtime/debug/)

<!-- CL 249677 -->
The [`runtime.Error`](/pkg/runtime#Error) values
used when `SetPanicOnFault` is enabled may now have an
`Addr` method. If that method exists, it returns the memory
address that triggered the fault.

<!-- runtime/debug -->

#### [strconv](/pkg/strconv/)

<!-- CL 260858 -->
[`ParseFloat`](/pkg/strconv/#ParseFloat) now uses
the [Eisel-Lemire
algorithm](https://nigeltao.github.io/blog/2020/eisel-lemire.html), improving performance by up to a factor of 2. This can
also speed up decoding textual formats like [`encoding/json`](/pkg/encoding/json/).

<!-- strconv -->

#### [syscall](/pkg/syscall/)

<!-- CL 263271 -->
[`NewCallback`](/pkg/syscall/?GOOS=windows#NewCallback)
and
[`NewCallbackCDecl`](/pkg/syscall/?GOOS=windows#NewCallbackCDecl)
now correctly support callback functions with multiple
sub-`uintptr`-sized arguments in a row. This may
require changing uses of these functions to eliminate manual
padding between small arguments.

<!-- CL 261917 -->
[`SysProcAttr`](/pkg/syscall/?GOOS=windows#SysProcAttr) on Windows has a new `NoInheritHandles` field that disables inheriting handles when creating a new process.

<!-- CL 269761, golang.org/issue/42584 -->
[`DLLError`](/pkg/syscall/?GOOS=windows#DLLError) on Windows now has an `Unwrap` method for unwrapping its underlying error.

<!-- CL 210639 -->
On Linux,
[`Setgid`](/pkg/syscall/#Setgid),
[`Setuid`](/pkg/syscall/#Setuid),
and related calls are now implemented.
Previously, they returned an `syscall.EOPNOTSUPP` error.

<!-- CL 210639 -->
On Linux, the new functions
[`AllThreadsSyscall`](/pkg/syscall/#AllThreadsSyscall)
and [`AllThreadsSyscall6`](/pkg/syscall/#AllThreadsSyscall6)
may be used to make a system call on all Go threads in the process.
These functions may only be used by programs that do not use cgo;
if a program uses cgo, they will always return
[`syscall.ENOTSUP`](/pkg/syscall/#ENOTSUP).

<!-- syscall -->

#### [testing/iotest](/pkg/testing/iotest/)

<!-- CL 199501 -->
The new
[`ErrReader`](/pkg/testing/iotest/#ErrReader)
function returns an
[`io.Reader`](/pkg/io/#Reader) that always
returns an error.

<!-- CL 243909 -->
The new
[`TestReader`](/pkg/testing/iotest/#TestReader)
function tests that an [`io.Reader`](/pkg/io/#Reader)
behaves correctly.

<!-- testing/iotest -->

#### [text/template](/pkg/text/template/)

<!-- CL 254257, golang.org/issue/29770 -->
Newlines characters are now allowed inside action delimiters,
permitting actions to span multiple lines.

<!-- CL 243938 -->
The new [`template.ParseFS`](/pkg/text/template/#ParseFS)
function and [`template.Template.ParseFS`](/pkg/text/template/#Template.ParseFS)
method are like [`template.ParseGlob`](/pkg/text/template/#ParseGlob)
and [`template.Template.ParseGlob`](/pkg/text/template/#Template.ParseGlob),
but read the templates from an [`fs.FS`](/pkg/io/fs/#FS).

<!-- text/template -->

#### [text/template/parse](/pkg/text/template/parse/)

<!-- CL 229398, golang.org/issue/34652 -->
A new [`CommentNode`](/pkg/text/template/parse/#CommentNode)
was added to the parse tree. The [`Mode`](/pkg/text/template/parse/#Mode)
field in the `parse.Tree` enables access to it.

<!-- text/template/parse -->

#### [time/tzdata](/pkg/time/tzdata/)

<!-- CL 261877 -->
The slim timezone data format is now used for the timezone database in
`$GOROOT/lib/time/zoneinfo.zip` and the embedded copy in this
package. This reduces the size of the timezone database by about 350 KB.

<!-- time/tzdata -->

#### [unicode](/pkg/unicode/)

<!-- CL 248765 -->
The [`unicode`](/pkg/unicode/) package and associated
support throughout the system has been upgraded from Unicode 12.0.0 to
[Unicode 13.0.0](https://www.unicode.org/versions/Unicode13.0.0/),
which adds 5,930 new characters, including four new scripts, and 55 new emoji.
Unicode 13.0.0 also designates plane 3 (U+30000-U+3FFFF) as the tertiary
ideographic plane.

<!-- unicode -->
