---
template: false
title: Go 1.5 Release Notes
---

## Introduction to Go 1.5 {#introduction}

The latest Go release, version 1.5,
is a significant release, including major architectural changes to the implementation.
Despite that, we expect almost all Go programs to continue to compile and run as before,
because the release still maintains the Go 1 [promise
of compatibility](/doc/go1compat.html).

The biggest developments in the implementation are:

  - The compiler and runtime are now written entirely in Go (with a little assembler).
    C is no longer involved in the implementation, and so the C compiler that was
    once necessary for building the distribution is gone.
  - The garbage collector is now [concurrent](/s/go14gc) and provides dramatically lower
    pause times by running, when possible, in parallel with other goroutines.
  - By default, Go programs run with `GOMAXPROCS` set to the
    number of cores available; in prior releases it defaulted to 1.
  - Support for [internal packages](/s/go14internal)
    is now provided for all repositories, not just the Go core.
  - The `go` command now provides [experimental
    support](/s/go15vendor) for "vendoring" external dependencies.
  - A new `go tool trace` command supports fine-grained
    tracing of program execution.
  - A new `go doc` command (distinct from `godoc`)
    is customized for command-line use.

These and a number of other changes to the implementation and tools
are discussed below.

The release also contains one small language change involving map literals.

Finally, the timing of the [release](/s/releasesched)
strays from the usual six-month interval,
both to provide more time to prepare this major release and to shift the schedule thereafter to
time the release dates more conveniently.

## Changes to the language {#language}

### Map literals {#map_literals}

Due to an oversight, the rule that allowed the element type to be elided from slice literals was not
applied to map keys.
This has been [corrected](/cl/2591) in Go 1.5.
An example will make this clear.
As of Go 1.5, this map literal,

	m := map[Point]string{
	    Point{29.935523, 52.891566}:   "Persepolis",
	    Point{-25.352594, 131.034361}: "Uluru",
	    Point{37.422455, -122.084306}: "Googleplex",
	}

may be written as follows, without the `Point` type listed explicitly:

	m := map[Point]string{
	    {29.935523, 52.891566}:   "Persepolis",
	    {-25.352594, 131.034361}: "Uluru",
	    {37.422455, -122.084306}: "Googleplex",
	}

## The Implementation {#implementation}

### No more C {#c}

The compiler and runtime are now implemented in Go and assembler, without C.
The only C source left in the tree is related to testing or to `cgo`.
There was a C compiler in the tree in 1.4 and earlier.
It was used to build the runtime; a custom compiler was necessary in part to
guarantee the C code would work with the stack management of goroutines.
Since the runtime is in Go now, there is no need for this C compiler and it is gone.
Details of the process to eliminate C are discussed [elsewhere](/s/go13compiler).

The conversion from C was done with the help of custom tools created for the job.
Most important, the compiler was actually moved by automatic translation of
the C code into Go.
It is in effect the same program in a different language.
It is not a new implementation
of the compiler so we expect the process will not have introduced new compiler
bugs.
An overview of this process is available in the slides for
[this presentation](/talks/2015/gogo.slide).

### Compiler and tools {#compiler_and_tools}

Independent of but encouraged by the move to Go, the names of the tools have changed.
The old names `6g`, `8g` and so on are gone; instead there
is just one binary, accessible as `go` `tool` `compile`,
that compiles Go source into binaries suitable for the architecture and operating system
specified by `$GOARCH` and `$GOOS`.
Similarly, there is now one linker (`go` `tool` `link`)
and one assembler (`go` `tool` `asm`).
The linker was translated automatically from the old C implementation,
but the assembler is a new native Go implementation discussed
in more detail below.

Similar to the drop of the names `6g`, `8g`, and so on,
the output of the compiler and assembler are now given a plain `.o` suffix
rather than `.8`, `.6`, etc.

### Garbage collector {#gc}

The garbage collector has been re-engineered for 1.5 as part of the development
outlined in the [design document](/s/go14gc).
Expected latencies are much lower than with the collector
in prior releases, through a combination of advanced algorithms,
better [scheduling](/s/go15gcpacing) of the collector,
and running more of the collection in parallel with the user program.
The "stop the world" phase of the collector
will almost always be under 10 milliseconds and usually much less.

For systems that benefit from low latency, such as user-responsive web sites,
the drop in expected latency with the new collector may be important.

Details of the new collector were presented in a
[talk](/talks/2015/go-gc.pdf) at GopherCon 2015.

### Runtime {#runtime}

In Go 1.5, the order in which goroutines are scheduled has been changed.
The properties of the scheduler were never defined by the language,
but programs that depend on the scheduling order may be broken
by this change.
We have seen a few (erroneous) programs affected by this change.
If you have programs that implicitly depend on the scheduling
order, you will need to update them.

Another potentially breaking change is that the runtime now
sets the default number of threads to run simultaneously,
defined by `GOMAXPROCS`, to the number
of cores available on the CPU.
In prior releases the default was 1.
Programs that do not expect to run with multiple cores may
break inadvertently.
They can be updated by removing the restriction or by setting
`GOMAXPROCS` explicitly.
For a more detailed discussion of this change, see
the [design document](/s/go15gomaxprocs).

### Build {#build}

Now that the Go compiler and runtime are implemented in Go, a Go compiler
must be available to compile the distribution from source.
Thus, to build the Go core, a working Go distribution must already be in place.
(Go programmers who do not work on the core are unaffected by this change.)
Any Go 1.4 or later distribution (including `gccgo`) will serve.
For details, see the [design document](/s/go15bootstrap).

## Ports {#ports}

Due mostly to the industry's move away from the 32-bit x86 architecture,
the set of binary downloads provided is reduced in 1.5.
A distribution for the OS X operating system is provided only for the
`amd64` architecture, not `386`.
Similarly, the ports for Snow Leopard (Apple OS X 10.6) still work but are no
longer released as a download or maintained since Apple no longer maintains that version
of the operating system.
Also, the `dragonfly/386` port is no longer supported at all
because DragonflyBSD itself no longer supports the 32-bit 386 architecture.

There are however several new ports available to be built from source.
These include `darwin/arm` and `darwin/arm64`.
The new port `linux/arm64` is mostly in place, but `cgo`
is only supported using external linking.

Also available as experiments are `ppc64`
and `ppc64le` (64-bit PowerPC, big- and little-endian).
Both these ports support `cgo` but
only with internal linking.

On FreeBSD, Go 1.5 requires FreeBSD 8-STABLE+ because of its new use of the `SYSCALL` instruction.

On NaCl, Go 1.5 requires SDK version pepper-41. Later pepper versions are not
compatible due to the removal of the sRPC subsystem from the NaCl runtime.

On Darwin, the use of the system X.509 certificate interface can be disabled
with the `ios` build tag.

The Solaris port now has full support for cgo and the packages
[`net`](/pkg/net/) and
[`crypto/x509`](/pkg/crypto/x509/),
as well as a number of other fixes and improvements.

## Tools {#tools}

### Translating {#translate}

As part of the process to eliminate C from the tree, the compiler and
linker were translated from C to Go.
It was a genuine (machine assisted) translation, so the new programs are essentially
the old programs translated rather than new ones with new bugs.
We are confident the translation process has introduced few if any new bugs,
and in fact uncovered a number of previously unknown bugs, now fixed.

The assembler is a new program, however; it is described below.

### Renaming {#rename}

The suites of programs that were the compilers (`6g`, `8g`, etc.),
the assemblers (`6a`, `8a`, etc.),
and the linkers (`6l`, `8l`, etc.)
have each been consolidated into a single tool that is configured
by the environment variables `GOOS` and `GOARCH`.
The old names are gone; the new tools are available through the `go` `tool`
mechanism as `go tool compile`,
`go tool asm`,
`and go tool link`.
Also, the file suffixes `.6`, `.8`, etc. for the
intermediate object files are also gone; now they are just plain `.o` files.

For example, to build and link a program on amd64 for Darwin
using the tools directly, rather than through `go build`,
one would run:

	$ export GOOS=darwin GOARCH=amd64
	$ go tool compile program.go
	$ go tool link program.o

### Moving {#moving}

Because the [`go/types`](/pkg/go/types/) package
has now moved into the main repository (see below),
the [`vet`](/cmd/vet) and
[`cover`](/cmd/cover)
tools have also been moved.
They are no longer maintained in the external `golang.org/x/tools` repository,
although (deprecated) source still resides there for compatibility with old releases.

### Compiler {#compiler}

As described above, the compiler in Go 1.5 is a single Go program,
translated from the old C source, that replaces `6g`, `8g`,
and so on.
Its target is configured by the environment variables `GOOS` and `GOARCH`.

The 1.5 compiler is mostly equivalent to the old,
but some internal details have changed.
One significant change is that evaluation of constants now uses
the [`math/big`](/pkg/math/big/) package
rather than a custom (and less well tested) implementation of high precision
arithmetic.
We do not expect this to affect the results.

For the amd64 architecture only, the compiler has a new option, `-dynlink`,
that assists dynamic linking by supporting references to Go symbols
defined in external shared libraries.

### Assembler {#assembler}

Like the compiler and linker, the assembler in Go 1.5 is a single program
that replaces the suite of assemblers (`6a`,
`8a`, etc.) and the environment variables
`GOARCH` and `GOOS`
configure the architecture and operating system.
Unlike the other programs, the assembler is a wholly new program
written in Go.

The new assembler is very nearly compatible with the previous
ones, but there are a few changes that may affect some
assembler source files.
See the updated [assembler guide](/doc/asm)
for more specific information about these changes. In summary:

First, the expression evaluation used for constants is a little
different.
It now uses unsigned 64-bit arithmetic and the precedence
of operators (`+`, `-`, `<<`, etc.)
comes from Go, not C.
We expect these changes to affect very few programs but
manual verification may be required.

Perhaps more important is that on machines where
`SP` or `PC` is only an alias
for a numbered register,
such as `R13` for the stack pointer and
`R15` for the hardware program counter
on ARM,
a reference to such a register that does not include a symbol
is now illegal.
For example, `SP` and `4(SP)` are
illegal but `sym+4(SP)` is fine.
On such machines, to refer to the hardware register use its
true `R` name.

One minor change is that some of the old assemblers
permitted the notation

	constant=value

to define a named constant.
Since this is always possible to do with the traditional
C-like `#define` notation, which is still
supported (the assembler includes an implementation
of a simplified C preprocessor), the feature was removed.

### Linker {#link}

The linker in Go 1.5 is now one Go program,
that replaces `6l`, `8l`, etc.
Its operating system and instruction set are specified
by the environment variables `GOOS` and `GOARCH`.

There are several other changes.
The most significant is the addition of a `-buildmode` option that
expands the style of linking; it now supports
situations such as building shared libraries and allowing other languages
to call into Go libraries.
Some of these were outlined in a [design document](/s/execmodes).
For a list of the available build modes and their use, run

	$ go help buildmode

Another minor change is that the linker no longer records build time stamps in
the header of Windows executables.
Also, although this may be fixed, Windows cgo executables are missing some
DWARF information.

Finally, the `-X` flag, which takes two arguments,
as in

	-X importpath.name value

now also accepts a more common Go flag style with a single argument
that is itself a `name=value` pair:

	-X importpath.name=value

Although the old syntax still works, it is recommended that uses of this
flag in scripts and the like be updated to the new form.

### Go command {#go_command}

The [`go`](/cmd/go) command's basic operation
is unchanged, but there are a number of changes worth noting.

The previous release introduced the idea of a directory internal to a package
being unimportable through the `go` command.
In 1.4, it was tested with the introduction of some internal elements
in the core repository.
As suggested in the [design document](/s/go14internal),
that change is now being made available to all repositories.
The rules are explained in the design document, but in summary any
package in or under a directory named `internal` may
be imported by packages rooted in the same subtree.
Existing packages with directory elements named `internal` may be
inadvertently broken by this change, which was why it was advertised
in the last release.

Another change in how packages are handled is the experimental
addition of support for "vendoring".
For details, see the documentation for the [`go` command](/cmd/go/#hdr-Vendor_Directories)
and the [design document](/s/go15vendor).

There have also been several minor changes.
Read the [documentation](/cmd/go) for full details.

  - SWIG support has been updated such that
    `.swig` and `.swigcxx`
    now require SWIG 3.0.6 or later.
  - The `install` subcommand now removes the
    binary created by the `build` subcommand
    in the source directory, if present,
    to avoid problems having two binaries present in the tree.
  - The `std` (standard library) wildcard package name
    now excludes commands.
    A new `cmd` wildcard covers the commands.
  - A new `-asmflags` build option
    sets flags to pass to the assembler.
    However,
    the `-ccflags` build option has been dropped;
    it was specific to the old, now deleted C compiler .
  - A new `-buildmode` build option
    sets the build mode, described above.
  - A new `-pkgdir` build option
    sets the location of installed package archives,
    to help isolate custom builds.
  - A new `-toolexec` build option
    allows substitution of a different command to invoke
    the compiler and so on.
    This acts as a custom replacement for `go tool`.
  - The `test` subcommand now has a `-count`
    flag to specify how many times to run each test and benchmark.
    The [`testing`](/pkg/testing/) package
    does the work here, through the `-test.count` flag.
  - The `generate` subcommand has a couple of new features.
    The `-run` option specifies a regular expression to select which directives
    to execute; this was proposed but never implemented in 1.4.
    The executing pattern now has access to two new environment variables:
    `$GOLINE` returns the source line number of the directive
    and `$DOLLAR` expands to a dollar sign.
  - The `get` subcommand now has a `-insecure`
    flag that must be enabled if fetching from an insecure repository, one that
    does not encrypt the connection.

### Go vet command {#vet_command}

The [`go tool vet`](/cmd/vet) command now does
more thorough validation of struct tags.

### Trace command {#trace_command}

A new tool is available for dynamic execution tracing of Go programs.
The usage is analogous to how the test coverage tool works.
Generation of traces is integrated into `go test`,
and then a separate execution of the tracing tool itself analyzes the results:

	$ go test -trace=trace.out path/to/package
	$ go tool trace [flags] pkg.test trace.out

The flags enable the output to be displayed in a browser window.
For details, run `go tool trace -help`.
There is also a description of the tracing facility in this
[talk](/talks/2015/dynamic-tools.slide)
from GopherCon 2015.

### Go doc command {#doc_command}

A few releases back, the `go doc`
command was deleted as being unnecessary.
One could always run "`godoc .`" instead.
The 1.5 release introduces a new [`go doc`](/cmd/doc)
command with a more convenient command-line interface than
`godoc`'s.
It is designed for command-line usage specifically, and provides a more
compact and focused presentation of the documentation for a package
or its elements, according to the invocation.
It also provides case-insensitive matching and
support for showing the documentation for unexported symbols.
For details run "`go help doc`".

### Cgo {#cgo}

When parsing `#cgo` lines,
the invocation `${SRCDIR}` is now
expanded into the path to the source directory.
This allows options to be passed to the
compiler and linker that involve file paths relative to the
source code directory. Without the expansion the paths would be
invalid when the current working directory changes.

Solaris now has full cgo support.

On Windows, cgo now uses external linking by default.

When a C struct ends with a zero-sized field, but the struct itself is
not zero-sized, Go code can no longer refer to the zero-sized field.
Any such references will have to be rewritten.

## Performance {#performance}

As always, the changes are so general and varied that precise statements
about performance are difficult to make.
The changes are even broader ranging than usual in this release, which
includes a new garbage collector and a conversion of the runtime to Go.
Some programs may run faster, some slower.
On average the programs in the Go 1 benchmark suite run a few percent faster in Go 1.5
than they did in Go 1.4,
while as mentioned above the garbage collector's pauses are
dramatically shorter, and almost always under 10 milliseconds.

Builds in Go 1.5 will be slower by a factor of about two.
The automatic translation of the compiler and linker from C to Go resulted in
unidiomatic Go code that performs poorly compared to well-written Go.
Analysis tools and refactoring helped to improve the code, but much remains to be done.
Further profiling and optimization will continue in Go 1.6 and future releases.
For more details, see these [slides](/talks/2015/gogo.slide)
and associated [video](https://www.youtube.com/watch?v=cF1zJYkBW4A).

## Standard library {#library}

### Flag {#flag}

The flag package's
[`PrintDefaults`](/pkg/flag/#PrintDefaults)
function, and method on [`FlagSet`](/pkg/flag/#FlagSet),
have been modified to create nicer usage messages.
The format has been changed to be more human-friendly and in the usage
messages a word quoted with \`backquotes\` is taken to be the name of the
flag's operand to display in the usage message.
For instance, a flag created with the invocation,

	cpuFlag = flag.Int("cpu", 1, "run `N` processes in parallel")

will show the help message,

	-cpu N
	    	run N processes in parallel (default 1)

Also, the default is now listed only when it is not the zero value for the type.

### Floats in math/big {#math_big}

The [`math/big`](/pkg/math/big/) package
has a new, fundamental data type,
[`Float`](/pkg/math/big/#Float),
which implements arbitrary-precision floating-point numbers.
A `Float` value is represented by a boolean sign,
a variable-length mantissa, and a 32-bit fixed-size signed exponent.
The precision of a `Float` (the mantissa size in bits)
can be specified explicitly or is otherwise determined by the first
operation that creates the value.
Once created, the size of a `Float`'s mantissa may be modified with the
[`SetPrec`](/pkg/math/big/#Float.SetPrec) method.
`Floats` support the concept of infinities, such as are created by
overflow, but values that would lead to the equivalent of IEEE 754 NaNs
trigger a panic.
`Float` operations support all IEEE-754 rounding modes.
When the precision is set to 24 (53) bits,
operations that stay within the range of normalized `float32`
(`float64`)
values produce the same results as the corresponding IEEE-754
arithmetic on those values.

### Go types {#go_types}

The [`go/types`](/pkg/go/types/) package
up to now has been maintained in the `golang.org/x`
repository; as of Go 1.5 it has been relocated to the main repository.
The code at the old location is now deprecated.
There is also a modest API change in the package, discussed below.

Associated with this move, the
[`go/constant`](/pkg/go/constant/)
package also moved to the main repository;
it was `golang.org/x/tools/exact` before.
The [`go/importer`](/pkg/go/importer/) package
also moved to the main repository,
as well as some tools described above.

### Net {#net}

The DNS resolver in the net package has almost always used `cgo` to access
the system interface.
A change in Go 1.5 means that on most Unix systems DNS resolution
will no longer require `cgo`, which simplifies execution
on those platforms.
Now, if the system's networking configuration permits, the native Go resolver
will suffice.
The important effect of this change is that each DNS resolution occupies a goroutine
rather than a thread,
so a program with multiple outstanding DNS requests will consume fewer operating
system resources.

The decision of how to run the resolver applies at run time, not build time.
The `netgo` build tag that has been used to enforce the use
of the Go resolver is no longer necessary, although it still works.
A new `netcgo` build tag forces the use of the `cgo` resolver at
build time.
To force `cgo` resolution at run time set
`GODEBUG=netdns=cgo` in the environment.
More debug options are documented [here](/cl/11584).

This change applies to Unix systems only.
Windows, Mac OS X, and Plan 9 systems behave as before.

### Reflect {#reflect}

The [`reflect`](/pkg/reflect/) package
has two new functions: [`ArrayOf`](/pkg/reflect/#ArrayOf)
and [`FuncOf`](/pkg/reflect/#FuncOf).
These functions, analogous to the extant
[`SliceOf`](/pkg/reflect/#SliceOf) function,
create new types at runtime to describe arrays and functions.

### Hardening {#hardening}

Several dozen bugs were found in the standard library
through randomized testing with the
[`go-fuzz`](https://github.com/dvyukov/go-fuzz) tool.
Bugs were fixed in the
[`archive/tar`](/pkg/archive/tar/),
[`archive/zip`](/pkg/archive/zip/),
[`compress/flate`](/pkg/compress/flate/),
[`encoding/gob`](/pkg/encoding/gob/),
[`fmt`](/pkg/fmt/),
[`html/template`](/pkg/html/template/),
[`image/gif`](/pkg/image/gif/),
[`image/jpeg`](/pkg/image/jpeg/),
[`image/png`](/pkg/image/png/), and
[`text/template`](/pkg/text/template/),
packages.
The fixes harden the implementation against incorrect and malicious inputs.

### Minor changes to the library {#minor_library_changes}

  - The [`archive/zip`](/pkg/archive/zip/) package's
    [`Writer`](/pkg/archive/zip/#Writer) type now has a
    [`SetOffset`](/pkg/archive/zip/#Writer.SetOffset)
    method to specify the location within the output stream at which to write the archive.
  - The [`Reader`](/pkg/bufio/#Reader) in the
    [`bufio`](/pkg/bufio/) package now has a
    [`Discard`](/pkg/bufio/#Reader.Discard)
    method to discard data from the input.
  - In the [`bytes`](/pkg/bytes/) package,
    the [`Buffer`](/pkg/bytes/#Buffer) type
    now has a [`Cap`](/pkg/bytes/#Buffer.Cap) method
    that reports the number of bytes allocated within the buffer.
    Similarly, in both the [`bytes`](/pkg/bytes/)
    and [`strings`](/pkg/strings/) packages,
    the [`Reader`](/pkg/bytes/#Reader)
    type now has a [`Size`](/pkg/bytes/#Reader.Size)
    method that reports the original length of the underlying slice or string.
  - Both the [`bytes`](/pkg/bytes/) and
    [`strings`](/pkg/strings/) packages
    also now have a [`LastIndexByte`](/pkg/bytes/#LastIndexByte)
    function that locates the rightmost byte with that value in the argument.
  - The [`crypto`](/pkg/crypto/) package
    has a new interface, [`Decrypter`](/pkg/crypto/#Decrypter),
    that abstracts the behavior of a private key used in asymmetric decryption.
  - In the [`crypto/cipher`](/pkg/crypto/cipher/) package,
    the documentation for the [`Stream`](/pkg/crypto/cipher/#Stream)
    interface has been clarified regarding the behavior when the source and destination are
    different lengths.
    If the destination is shorter than the source, the method will panic.
    This is not a change in the implementation, only the documentation.
  - Also in the [`crypto/cipher`](/pkg/crypto/cipher/) package,
    there is now support for nonce lengths other than 96 bytes in AES's Galois/Counter mode (GCM),
    which some protocols require.
  - In the [`crypto/elliptic`](/pkg/crypto/elliptic/) package,
    there is now a `Name` field in the
    [`CurveParams`](/pkg/crypto/elliptic/#CurveParams) struct,
    and the curves implemented in the package have been given names.
    These names provide a safer way to select a curve, as opposed to
    selecting its bit size, for cryptographic systems that are curve-dependent.
  - Also in the [`crypto/elliptic`](/pkg/crypto/elliptic/) package,
    the [`Unmarshal`](/pkg/crypto/elliptic/#Unmarshal) function
    now verifies that the point is actually on the curve.
    (If it is not, the function returns nils).
    This change guards against certain attacks.
  - The [`crypto/sha512`](/pkg/crypto/sha512/)
    package now has support for the two truncated versions of
    the SHA-512 hash algorithm, SHA-512/224 and SHA-512/256.
  - The [`crypto/tls`](/pkg/crypto/tls/) package
    minimum protocol version now defaults to TLS 1.0.
    The old default, SSLv3, is still available through [`Config`](/pkg/crypto/tls/#Config) if needed.
  - The [`crypto/tls`](/pkg/crypto/tls/) package
    now supports Signed Certificate Timestamps (SCTs) as specified in RFC 6962.
    The server serves them if they are listed in the
    [`Certificate`](/pkg/crypto/tls/#Certificate) struct,
    and the client requests them and exposes them, if present,
    in its [`ConnectionState`](/pkg/crypto/tls/#ConnectionState) struct.
  - The stapled OCSP response to a [`crypto/tls`](/pkg/crypto/tls/) client connection,
    previously only available via the
    [`OCSPResponse`](/pkg/crypto/tls/#Conn.OCSPResponse) method,
    is now exposed in the [`ConnectionState`](/pkg/crypto/tls/#ConnectionState) struct.
  - The [`crypto/tls`](/pkg/crypto/tls/) server implementation
    will now always call the
    `GetCertificate` function in
    the [`Config`](/pkg/crypto/tls/#Config) struct
    to select a certificate for the connection when none is supplied.
  - Finally, the session ticket keys in the
    [`crypto/tls`](/pkg/crypto/tls/) package
    can now be changed while the server is running.
    This is done through the new
    [`SetSessionTicketKeys`](/pkg/crypto/tls/#Config.SetSessionTicketKeys)
    method of the
    [`Config`](/pkg/crypto/tls/#Config) type.
  - In the [`crypto/x509`](/pkg/crypto/x509/) package,
    wildcards are now accepted only in the leftmost label as defined in
    [the specification](https://tools.ietf.org/html/rfc6125#section-6.4.3).
  - Also in the [`crypto/x509`](/pkg/crypto/x509/) package,
    the handling of unknown critical extensions has been changed.
    They used to cause parse errors but now they are parsed and caused errors only
    in [`Verify`](/pkg/crypto/x509/#Certificate.Verify).
    The new field `UnhandledCriticalExtensions` of
    [`Certificate`](/pkg/crypto/x509/#Certificate) records these extensions.
  - The [`DB`](/pkg/database/sql/#DB) type of the
    [`database/sql`](/pkg/database/sql/) package
    now has a [`Stats`](/pkg/database/sql/#DB.Stats) method
    to retrieve database statistics.
  - The [`debug/dwarf`](/pkg/debug/dwarf/)
    package has extensive additions to better support DWARF version 4.
    See for example the definition of the new type
    [`Class`](/pkg/debug/dwarf/#Class).
  - The [`debug/dwarf`](/pkg/debug/dwarf/) package
    also now supports decoding of DWARF line tables.
  - The [`debug/elf`](/pkg/debug/elf/)
    package now has support for the 64-bit PowerPC architecture.
  - The [`encoding/base64`](/pkg/encoding/base64/) package
    now supports unpadded encodings through two new encoding variables,
    [`RawStdEncoding`](/pkg/encoding/base64/#RawStdEncoding) and
    [`RawURLEncoding`](/pkg/encoding/base64/#RawURLEncoding).
  - The [`encoding/json`](/pkg/encoding/json/) package
    now returns an [`UnmarshalTypeError`](/pkg/encoding/json/#UnmarshalTypeError)
    if a JSON value is not appropriate for the target variable or component
    to which it is being unmarshaled.
  - The `encoding/json`'s
    [`Decoder`](/pkg/encoding/json/#Decoder)
    type has a new method that provides a streaming interface for decoding
    a JSON document:
    [`Token`](/pkg/encoding/json/#Decoder.Token).
    It also interoperates with the existing functionality of `Decode`,
    which will continue a decode operation already started with `Decoder.Token`.
  - The [`flag`](/pkg/flag/) package
    has a new function, [`UnquoteUsage`](/pkg/flag/#UnquoteUsage),
    to assist in the creation of usage messages using the new convention
    described above.
  - In the [`fmt`](/pkg/fmt/) package,
    a value of type [`Value`](/pkg/reflect/#Value) now
    prints what it holds, rather than use the `reflect.Value`'s `Stringer`
    method, which produces things like `<int Value>`.
  - The [`EmptyStmt`](/pkg/ast/#EmptyStmt) type
    in the [`go/ast`](/pkg/go/ast/) package now
    has a boolean `Implicit` field that records whether the
    semicolon was implicitly added or was present in the source.
  - For forward compatibility the [`go/build`](/pkg/go/build/) package
    reserves `GOARCH` values for a number of architectures that Go might support one day.
    This is not a promise that it will.
    Also, the [`Package`](/pkg/go/build/#Package) struct
    now has a `PkgTargetRoot` field that stores the
    architecture-dependent root directory in which to install, if known.
  - The (newly migrated) [`go/types`](/pkg/go/types/)
    package allows one to control the prefix attached to package-level names using
    the new [`Qualifier`](/pkg/go/types/#Qualifier)
    function type as an argument to several functions. This is an API change for
    the package, but since it is new to the core, it is not breaking the Go 1 compatibility
    rules since code that uses the package must explicitly ask for it at its new location.
    To update, run
    [`go fix`](/cmd/go/#hdr-Run_go_tool_fix_on_packages) on your package.
  - In the [`image`](/pkg/image/) package,
    the [`Rectangle`](/pkg/image/#Rectangle) type
    now implements the [`Image`](/pkg/image/#Image) interface,
    so a `Rectangle` can serve as a mask when drawing.
  - Also in the [`image`](/pkg/image/) package,
    to assist in the handling of some JPEG images,
    there is now support for 4:1:1 and 4:1:0 YCbCr subsampling and basic
    CMYK support, represented by the new `image.CMYK` struct.
  - The [`image/color`](/pkg/image/color/) package
    adds basic CMYK support, through the new
    [`CMYK`](/pkg/image/color/#CMYK) struct,
    the [`CMYKModel`](/pkg/image/color/#CMYKModel) color model, and the
    [`CMYKToRGB`](/pkg/image/color/#CMYKToRGB) function, as
    needed by some JPEG images.
  - Also in the [`image/color`](/pkg/image/color/) package,
    the conversion of a [`YCbCr`](/pkg/image/color/#YCbCr)
    value to `RGBA` has become more precise.
    Previously, the low 8 bits were just an echo of the high 8 bits;
    now they contain more accurate information.
    Because of the echo property of the old code, the operation
    `uint8(r)` to extract an 8-bit red value worked, but is incorrect.
    In Go 1.5, that operation may yield a different value.
    The correct code is, and always was, to select the high 8 bits:
    `uint8(r>>8)`.
    Incidentally, the `image/draw` package
    provides better support for such conversions; see
    [this blog post](/blog/go-imagedraw-package)
    for more information.
  - Finally, as of Go 1.5 the closest match check in
    [`Index`](/pkg/image/color/#Palette.Index)
    now honors the alpha channel.
  - The [`image/gif`](/pkg/image/gif/) package
    includes a couple of generalizations.
    A multiple-frame GIF file can now have an overall bounds different
    from all the contained single frames' bounds.
    Also, the [`GIF`](/pkg/image/gif/#GIF) struct
    now has a `Disposal` field
    that specifies the disposal method for each frame.
  - The [`io`](/pkg/io/) package
    adds a [`CopyBuffer`](/pkg/io/#CopyBuffer) function
    that is like [`Copy`](/pkg/io/#Copy) but
    uses a caller-provided buffer, permitting control of allocation and buffer size.
  - The [`log`](/pkg/log/) package
    has a new [`LUTC`](/pkg/log/#LUTC) flag
    that causes time stamps to be printed in the UTC time zone.
    It also adds a [`SetOutput`](/pkg/log/#Logger.SetOutput) method
    for user-created loggers.
  - In Go 1.4, [`Max`](/pkg/math/#Max) was not detecting all possible NaN bit patterns.
    This is fixed in Go 1.5, so programs that use `math.Max` on data including NaNs may behave differently,
    but now correctly according to the IEEE754 definition of NaNs.
  - The [`math/big`](/pkg/math/big/) package
    adds a new [`Jacobi`](/pkg/math/big/#Jacobi)
    function for integers and a new
    [`ModSqrt`](/pkg/math/big/#Int.ModSqrt)
    method for the [`Int`](/pkg/math/big/#Int) type.
  - The mime package
    adds a new [`WordDecoder`](/pkg/mime/#WordDecoder) type
    to decode MIME headers containing RFC 204-encoded words.
    It also provides [`BEncoding`](/pkg/mime/#BEncoding) and
    [`QEncoding`](/pkg/mime/#QEncoding)
    as implementations of the encoding schemes of RFC 2045 and RFC 2047.
  - The [`mime`](/pkg/mime/) package also adds an
    [`ExtensionsByType`](/pkg/mime/#ExtensionsByType)
    function that returns the MIME extensions know to be associated with a given MIME type.
  - There is a new [`mime/quotedprintable`](/pkg/mime/quotedprintable/)
    package that implements the quoted-printable encoding defined by RFC 2045.
  - The [`net`](/pkg/net/) package will now
    [`Dial`](/pkg/net/#Dial) hostnames by trying each
    IP address in order until one succeeds.
    The <code>[Dialer](/pkg/net/#Dialer).DualStack</code>
    mode now implements Happy Eyeballs
    ([RFC 6555](https://tools.ietf.org/html/rfc6555)) by giving the
    first address family a 300ms head start; this value can be overridden by
    the new `Dialer.FallbackDelay`.
  - A number of inconsistencies in the types returned by errors in the
    [`net`](/pkg/net/) package have been
    tidied up.
    Most now return an
    [`OpError`](/pkg/net/#OpError) value
    with more information than before.
    Also, the [`OpError`](/pkg/net/#OpError)
    type now includes a `Source` field that holds the local
    network address.
  - The [`net/http`](/pkg/net/http/) package now
    has support for setting trailers from a server [`Handler`](/pkg/net/http/#Handler).
    For details, see the documentation for
    [`ResponseWriter`](/pkg/net/http/#ResponseWriter).
  - There is a new method to cancel a [`net/http`](/pkg/net/http/)
    `Request` by setting the new
    [`Request.Cancel`](/pkg/net/http/#Request)
    field.
    It is supported by `http.Transport`.
    The `Cancel` field's type is compatible with the
    [`context.Context.Done`](https://godoc.org/golang.org/x/net/context)
    return value.
  - Also in the [`net/http`](/pkg/net/http/) package,
    there is code to ignore the zero [`Time`](/pkg/time/#Time) value
    in the [`ServeContent`](/pkg/net/#ServeContent) function.
    As of Go 1.5, it now also ignores a time value equal to the Unix epoch.
  - The [`net/http/fcgi`](/pkg/net/http/fcgi/) package
    exports two new errors,
    [`ErrConnClosed`](/pkg/net/http/fcgi/#ErrConnClosed) and
    [`ErrRequestAborted`](/pkg/net/http/fcgi/#ErrRequestAborted),
    to report the corresponding error conditions.
  - The [`net/http/cgi`](/pkg/net/http/cgi/) package
    had a bug that mishandled the values of the environment variables
    `REMOTE_ADDR` and `REMOTE_HOST`.
    This has been fixed.
    Also, starting with Go 1.5 the package sets the `REMOTE_PORT`
    variable.
  - The [`net/mail`](/pkg/net/mail/) package
    adds an [`AddressParser`](/pkg/net/mail/#AddressParser)
    type that can parse mail addresses.
  - The [`net/smtp`](/pkg/net/smtp/) package
    now has a [`TLSConnectionState`](/pkg/net/smtp/#Client.TLSConnectionState)
    accessor to the [`Client`](/pkg/net/smtp/#Client)
    type that returns the client's TLS state.
  - The [`os`](/pkg/os/) package
    has a new [`LookupEnv`](/pkg/os/#LookupEnv) function
    that is similar to [`Getenv`](/pkg/os/#Getenv)
    but can distinguish between an empty environment variable and a missing one.
  - The [`os/signal`](/pkg/os/signal/) package
    adds new [`Ignore`](/pkg/os/signal/#Ignore) and
    [`Reset`](/pkg/os/signal/#Reset) functions.
  - The [`runtime`](/pkg/runtime/),
    [`runtime/trace`](/pkg/runtime/trace/),
    and [`net/http/pprof`](/pkg/net/http/pprof/) packages
    each have new functions to support the tracing facilities described above:
    [`ReadTrace`](/pkg/runtime/#ReadTrace),
    [`StartTrace`](/pkg/runtime/#StartTrace),
    [`StopTrace`](/pkg/runtime/#StopTrace),
    [`Start`](/pkg/runtime/trace/#Start),
    [`Stop`](/pkg/runtime/trace/#Stop), and
    [`Trace`](/pkg/net/http/pprof/#Trace).
    See the respective documentation for details.
  - The [`runtime/pprof`](/pkg/runtime/pprof/) package
    by default now includes overall memory statistics in all memory profiles.
  - The [`strings`](/pkg/strings/) package
    has a new [`Compare`](/pkg/strings/#Compare) function.
    This is present to provide symmetry with the [`bytes`](/pkg/bytes/) package
    but is otherwise unnecessary as strings support comparison natively.
  - The [`WaitGroup`](/pkg/sync/#WaitGroup) implementation in
    package [`sync`](/pkg/sync/)
    now diagnoses code that races a call to [`Add`](/pkg/sync/#WaitGroup.Add)
    against a return from [`Wait`](/pkg/sync/#WaitGroup.Wait).
    If it detects this condition, the implementation panics.
  - In the [`syscall`](/pkg/syscall/) package,
    the Linux `SysProcAttr` struct now has a
    `GidMappingsEnableSetgroups` field, made necessary
    by security changes in Linux 3.19.
    On all Unix systems, the struct also has new `Foreground` and `Pgid` fields
    to provide more control when exec'ing.
    On Darwin, there is now a `Syscall9` function
    to support calls with too many arguments.
  - The [`testing/quick`](/pkg/testing/quick/) will now
    generate `nil` values for pointer types,
    making it possible to use with recursive data structures.
    Also, the package now supports generation of array types.
  - In the [`text/template`](/pkg/text/template/) and
    [`html/template`](/pkg/html/template/) packages,
    integer constants too large to be represented as a Go integer now trigger a
    parse error. Before, they were silently converted to floating point, losing
    precision.
  - Also in the [`text/template`](/pkg/text/template/) and
    [`html/template`](/pkg/html/template/) packages,
    a new [`Option`](/pkg/text/template/#Template.Option) method
    allows customization of the behavior of the template during execution.
    The sole implemented option allows control over how a missing key is
    handled when indexing a map.
    The default, which can now be overridden, is as before: to continue with an invalid value.
  - The [`time`](/pkg/time/) package's
    `Time` type has a new method
    [`AppendFormat`](/pkg/time/#Time.AppendFormat),
    which can be used to avoid allocation when printing a time value.
  - The [`unicode`](/pkg/unicode/) package and associated
    support throughout the system has been upgraded from version 7.0 to
    [Unicode 8.0](https://www.unicode.org/versions/Unicode8.0.0/).
