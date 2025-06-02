---
title: Go 1.25 Release Notes
template: false
---

<style>
  main ul li { margin: 0.5em 0; }
</style>

## DRAFT RELEASE NOTES — Introduction to Go 1.N {#introduction}

**Go 1.25 is not yet released. These are work-in-progress release notes.
Go 1.25 is expected to be released in August 2025.**

## Tools {#tools}

### Go command {#go-command}

The `go build` `-asan` option now defaults to doing leak detection at
program exit.
This will report an error if memory allocated by C is not freed and is
not referenced by any other memory allocated by either C or Go.
These new error reports may be disabled by setting
`ASAN_OPTIONS=detect_leaks=0` in the environment when running the
program.

<!-- go.dev/issue/71294 -->

The new `work` package pattern matches all packages in the work (formerly called main)
modules: either the single work module in module mode or the set of workspace modules
in workspace mode.

<!-- go.dev/issue/65847 -->

When the go command updates the `go` line in a `go.mod` or `go.work` file,
it [no longer](/ref/mod#go-mod-file-toolchain) adds a toolchain line
specifying the command's current version.

### Vet {#vet}

The `go vet` command includes new analyzers:

<!-- go.dev/issue/18022 -->

- [waitgroup](https://pkg.go.dev/golang.org/x/tools/go/analysis/passes/waitgroup),
  which reports misplaced calls to [`sync.WaitGroup.Add`](/pkg/sync#WaitGroup.Add); and

<!-- go.dev/issue/28308 -->

- [hostport](https://pkg.go.dev/golang.org/x/tools/go/analysis/passes/hostport),
  which reports uses of `fmt.Sprintf("%s:%d", host, port)` to
  construct addresses for [`net.Dial`](/pkg/net#Dial), as these will not work with
  IPv6; instead it suggests using [`net.JoinHostPort`](/pkg/net#JoinHostPort).

## Runtime {#runtime}

### Container-aware `GOMAXPROCS`

<!-- go.dev/issue/73193 -->

The default behavior of the `GOMAXPROCS` has changed. In prior versions of Go,
`GOMAXPROCS` defaults to the number of logical CPUs available at startup
([`runtime.NumCPU`](/pkg/runtime#NumCPU)). Go 1.25 introduces two changes:

1. On Linux, the runtime considers the CPU bandwidth limit of the cgroup
   containing the process, if any. If the CPU bandwidth limit is lower than the
   number of logical CPUs available, `GOMAXPROCS` will default to the lower
   limit. In container runtime systems like Kubernetes, cgroup CPU bandwidth
   limits generally correspond to the "CPU limit" option. The Go runtime does
   not consider the "CPU requests" option.

2. On all OSes, the runtime periodically updates `GOMAXPROCS` if the number
   of logical CPUs available or the cgroup CPU bandwidth limit change.

Both of these behaviors are automatically disabled if `GOMAXPROCS` is set
manually via the `GOMAXPROCS` environment variable or a call to
[`runtime.GOMAXPROCS`](/pkg/runtime#GOMAXPROCS). They can also be disabled explicitly with the [GODEBUG
settings](/doc/godebug) `containermaxprocs=0` and `updatemaxprocs=0`,
respectively.

In order to support reading updated cgroup limits, the runtime will keep cached
file descriptors for the cgroup files for the duration of the process lifetime.

### New experimental garbage collector

<!-- go.dev/issue/73581 -->

A new garbage collector is now available as an experiment. This garbage
collector's design improves the performance of marking and scanning small objects
through better locality and CPU scalability. Benchmark result vary, but we expect
somewhere between a 10—40% reduction in garbage collection overhead in real-world
programs that heavily use the garbage collector.

The new garbage collector may be enabled by setting `GOEXPERIMENT=greenteagc`
at build time. We expect the design to continue to evolve and improve. To that
end, we encourage Go developers to try it out and report back their experiences.
See the [GitHub issue](/issue/73581) for more details on the design and
instructions for sharing feedback.

### Change to unhandled panic output

<!-- go.dev/issue/71517 -->

The message printed when a program exits due to an unhandled panic
that was recovered and repanicked no longer repeats the text of
the panic value.

Previously, a program which panicked with `panic("PANIC")`,
recovered the panic, and then repanicked with the original
value would print:

    panic: PANIC [recovered]
      panic: PANIC

This program will now print:

    panic: PANIC [recovered, repanicked]

### VMA names on Linux

<!-- go.dev/issue/71546 -->

On Linux systems with kernel support for anonymous VMA names
(`CONFIG_ANON_VMA_NAME`), the Go runtime will annotate anonymous memory
mappings with context about their purpose. e.g., `[anon: Go: heap]` for heap
memory. This can be disabled with the [GODEBUG setting](/doc/godebug)
`decoratemappings=0`.

## Compiler {#compiler}

<!-- https://go.dev/issue/26379 -->

The compiler and linker in Go 1.25 now generate debug information
using [DWARF version 5](https://dwarfstd.org/dwarf5std.html); the
newer DWARF version reduces the space required for debugging
information in Go binaries.
DWARF 5 generation is gated by the "dwarf5" GOEXPERIMENT; this
functionality can be disabled (for now) using GOEXPERIMENT=nodwarf5.

<!-- https://go.dev/issue/72860, CL 657715 -->

The compiler [has been fixed](/cl/657715)
to ensure that nil pointer checks are performed promptly. Programs like the following,
which used to execute successfully, will now panic with a nil-pointer exception:

```
package main

import "os"

func main() {
	f, err := os.Open("nonExistentFile")
	name := f.Name()
	if err != nil {
		return
	}
	println(name)
}
```

This program is incorrect in that it uses the result of `os.Open` before checking
the error. The main result of `os.Open` can be a nil pointer if the error result is non-nil.
But because of [a compiler bug](/issue/72860), this program ran successfully under
Go versions 1.21 through 1.24 (in violation of the Go spec). It will no longer run
successfully in Go 1.25. If this change is affecting your code, the solution is to put
the non-nil error check earlier in your code, preferably immediately after
the error-generating statement.

<!-- CLs 653856, 657937, 663795, 664299 -->

The compiler can now allocate the backing store for slices on the
stack in more situations, which improves performance. This change has
the potential to amplify the effects of incorrect
[unsafe.Pointer](/pkg/unsafe#Pointer) usage, see for example [issue
73199](/issue/73199). In order to track down these problems, the
[bisect tool](https://pkg.go.dev/golang.org/x/tools/cmd/bisect) can be
used to find the allocation causing trouble using the
`-compile=variablemake` flag. All such new stack allocations can also
be turned off using `-gcflags=all=-d=variablemakehash=n`.

## Standard library {#library}

### New testing/synctest package

<!-- go.dev/issue/67434, go.dev/issue/73567 -->
The new [testing/synctest](/pkg/testing/synctest) package
provides support for testing concurrent code.

The [`synctest.Test`](/pkg/synctest#Test) function runs a test function in an isolated
"bubble". Within the bubble, [time](/pkg/time) package functions
operate on a fake clock.

The [`synctest.Wait`](/pkg/synctest#Wait) function waits for all goroutines in the
current bubble to block.

### Minor changes to the library {#minor_library_changes}

#### [`archive/tar`](/pkg/archive/tar/)

The [`*Writer.AddFS`](/pkg/archive/tar#Writer.AddFS) implementation now supports symbolic links
for filesystems that implement [`io/fs.ReadLinkFS`](/pkg/io/fs#ReadLinkFS).

#### [`crypto`](/pkg/crypto/)

[`MessageSigner`](/pkg/crypto#MessageSigner) is a new signing interface that can be implemented by signers that wish to hash the message to be signed themselves. A new function is also introduced, [`SignMessage`](/pkg/crypto#SignMessage) which attempts to update a [`Signer`](/pkg/crypto#Signer) interface to [`MessageSigner`](/pkg/crypto#MessageSigner), using the [`MessageSigner.SignMessage`](/pkg/crypto#MessageSigner.SignMessage) method if successful, and [`Signer.Sign`](/pkg/crypto#Signer.Sign) if not. This can be used when code wishes to support both [`Signer`](/pkg/crypto#Signer) and [`MessageSigner`](/pkg/crypto#MessageSigner).

#### [`crypto/ecdsa`](/pkg/crypto/ecdsa/)

The new [`ParseRawPrivateKey`](/pkg/crypto/ecdsa#ParseRawPrivateKey), [`ParseUncompressedPublicKey`](/pkg/crypto/ecdsa#ParseUncompressedPublicKey), [`PrivateKey.Bytes`](/pkg/crypto/ecdsa#PrivateKey.Bytes),
and [`PublicKey.Bytes`](/pkg/crypto/ecdsa#PublicKey.Bytes) functions and methods implement low-level encodings,
replacing the need to use crypto/elliptic or math/big functions and methods.

#### [`crypto/elliptic`](/pkg/crypto/elliptic/)

The hidden and undocumented `Inverse` and `CombinedMult` methods on some [`Curve`](/pkg/crypto/elliptic#Curve)
implementations have been removed.

#### [`crypto/sha3`](/pkg/crypto/sha3/)

The new [`SHA3.Clone`](/pkg/crypto/sha3#SHA3.Clone) method implements [hash.Cloner](/pkg/hash#Cloner).

#### [`crypto/tls`](/pkg/crypto/tls/)

The new [`ConnectionState.CurveID`](/pkg/crypto/tls#ConnectionState.CurveID) field exposes the key exchange mechanism used
to establish the connection.

The new [`Config.GetEncryptedClientHelloKeys`](/pkg/crypto/tls#Config.GetEncryptedClientHelloKeys) callback can be used to set the
[EncryptedClientHelloKey]s for a server to use when a client sends an Encrypted
Client Hello extension.

SHA-1 signature algorithms are now disallowed in TLS 1.2 handshakes, per
[RFC 9155](https://www.rfc-editor.org/rfc/rfc9155.html).
They can be re-enabled with the `tlssha1=1` GODEBUG option.

When [FIPS 140-3 mode](/doc/security/fips140) is enabled, Extended Master Secret
is now required in TLS 1.2, and Ed25519 and X25519MLKEM768 are now allowed.

TLS servers now prefer the highest supported protocol version, even if it isn't the client's most preferred protocol version.

#### [`crypto/x509`](/pkg/crypto/x509/)

[`CreateCertificate`](/pkg/crypto/x509#CreateCertificate), [`CreateCertificateRequest`](/pkg/crypto/x509#CreateCertificateRequest), and [`CreateRevocationList`](/pkg/crypto/x509#CreateRevocationList) can now accept a [`crypto.MessageSigner`](/pkg/crypto#MessageSigner) signing interface as well as [`crypto.Signer`](/pkg/crypto#Signer). This allows these functions to use signers which implement "one-shot" signing interfaces, where hashing is done as part of the signing operation, instead of by the caller.

[`CreateCertificate`](/pkg/crypto/x509#CreateCertificate) now uses truncated SHA-256 to populate the `SubjectKeyId` if
it is missing. The GODEBUG setting `x509sha256skid=0` reverts to SHA-1.

#### [`debug/elf`](/pkg/debug/elf/)

The [`debug/elf`](/pkg/debug/elf) package adds two new constants:
- [`PT_RISCV_ATTRIBUTES`](/pkg/debug/elf#PT_RISCV_ATTRIBUTES)
- [`SHT_RISCV_ATTRIBUTES`](/pkg/debug/elf#SHT_RISCV_ATTRIBUTES)
  for RISC-V ELF parsing.

#### [`go/ast`](/pkg/go/ast/)

The [`ast.FilterPackage`](/pkg/ast#FilterPackage), [`ast.PackageExports`](/pkg/ast#PackageExports), and
[`ast.MergePackageFiles`](/pkg/ast#MergePackageFiles) functions, and the [`MergeMode`](/pkg/go/ast#MergeMode) type and its
constants, are all deprecated, as they are for use only with the
long-deprecated [`ast.Object`](/pkg/ast#Object) and [`ast.Package`](/pkg/ast#Package) machinery.

The new [`PreorderStack`](/pkg/go/ast#PreorderStack) function, like [`Inspect`](/pkg/go/ast#Inspect), traverses a syntax
tree and provides control over descent into subtrees, but as a
convenience it also provides the stack of enclosing nodes at each
point.

#### [`go/parser`](/pkg/go/parser/)

The [`ParseDir`](/pkg/go/parser#ParseDir) function is deprecated.

#### [`go/token`](/pkg/go/token/)

The new [`FileSet.AddExistingFiles`](/pkg/go/token#FileSet.AddExistingFiles) method enables existing Files to be
added to a FileSet, or a FileSet to be constructed for an arbitrary
set of Files, alleviating the problems associated with a single global
FileSet in long-lived applications.

#### [`go/types`](/pkg/go/types/)

[`Var`](/pkg/go/types#Var) now has a [`Var.Kind`](/pkg/go/types#Var.Kind) method that classifies the variable as one
of: package-level, receiver, parameter, result, or local variable, or
a struct field.

The new [`LookupSelection`](/pkg/go/types#LookupSelection) function looks up the field or method of a
given name and receiver type, like the existing [`LookupFieldOrMethod`](/pkg/go/types#LookupFieldOrMethod)
function, but returns the result in the form of a [`Selection`](/pkg/go/types#Selection).

#### [`hash`](/pkg/hash/)

The new [XOF](/pkg/hash#XOF) interface can be implemented by "extendable output
functions", which are hash functions with arbitrary or unlimited output length
such as [SHAKE](https://pkg.go.dev/crypto/sha3#SHAKE).

Hashes implementing the new [`Cloner`](/pkg/hash#Cloner) interface can return a copy of their state.
All standard library [`Hash`](/pkg/hash#Hash) implementations now implement [`Cloner`](/pkg/hash#Cloner).

#### [`hash/maphash`](/pkg/hash/maphash/)

The new [`Hash.Clone`](/pkg/hash/maphash#Hash.Clone) method implements [hash.Cloner](/pkg/hash#Cloner).

#### [`io/fs`](/pkg/io/fs/)

A new [`ReadLinkFS`](/pkg/io/fs#ReadLinkFS) interface provides the ability to read symbolic links in a filesystem.

#### [`log/slog`](/pkg/log/slog/)

[`GroupAttrs`](/pkg/log/slog#GroupAttrs) creates a group [`Attr`](/pkg/log/slog#Attr) from a slice of [`Attr`](/pkg/log/slog#Attr) values.

[`Record`](/pkg/log/slog#Record) now has a Source() method, returning its source location or nil if unavailable.

#### [`mime/multipart`](/pkg/mime/multipart/)

The new helper function [`FieldContentDisposition`](/pkg/mime/multipart#FieldContentDisposition) builds multipart
Content-Disposition header fields.

#### [`net`](/pkg/net/)

On Windows, the [`TCPConn.File`](/pkg/net#TCPConn.File), [`UDPConn.File`](/pkg/net#UDPConn.File), [`UnixConn.File`](/pkg/net#UnixConn.File),
[`IPConn.File`](/pkg/net#IPConn.File), [`TCPListener.File`](/pkg/net#TCPListener.File), and [`UnixListener.File`](/pkg/net#UnixListener.File)
methods are now supported.

[`LookupMX`](/pkg/net#LookupMX) and [`*Resolver.LookupMX`](/pkg/net#Resolver.LookupMX) now return DNS names that look
like valid IP address, as well as valid domain names.
Previously if a name server returned an IP address as a DNS name,
LookupMX would discard it, as required by the RFCs.
However, name servers in practice do sometimes return IP addresses.

On Windows, the [`ListenMulticastUDP`](/pkg/net#ListenMulticastUDP) now supports IPv6 addresses.

On Windows, the [`FileConn`](/pkg/net#FileConn), [`FilePacketConn`](/pkg/net#FilePacketConn), [`FileListener`](/pkg/net#FileListener)
functions are now supported.

#### [`net/http`](/pkg/net/http/)

The new [`CrossOriginProtection`](/pkg/net/http#CrossOriginProtection) implements protections against [Cross-Site
Request Forgery (CSRF)](https://developer.mozilla.org/en-US/docs/Web/Security/Attacks/CSRF) by rejecting non-safe cross-origin browser requests.
It uses [modern browser Fetch metadata](https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Sec-Fetch-Site), doesn't require tokens
or cookies, and supports origin-based and pattern-based bypasses.

#### [`os`](/pkg/os/)

On Windows, [`NewFile`](/pkg/os#NewFile) now supports handles opened for asynchronous I/O (that is,
[`syscall.FILE_FLAG_OVERLAPPED`](/pkg/syscall#FILE_FLAG_OVERLAPPED) is specified in the [`syscall.CreateFile`](/pkg/syscall#CreateFile) call).
These handles are associated with the Go runtime's I/O completion port,
which provides the following benefits for the resulting [`File`](/pkg/os#File):

- I/O methods ([`File.Read`](/pkg/os#File.Read), [`File.Write`](/pkg/os#File.Write), [`File.ReadAt`](/pkg/os#File.ReadAt), and [`File.WriteAt`](/pkg/os#File.WriteAt)) do not block an OS thread.
- Deadline methods ([`File.SetDeadline`](/pkg/os#File.SetDeadline), [`File.SetReadDeadline`](/pkg/os#File.SetReadDeadline), and [`File.SetWriteDeadline`](/pkg/os#File.SetWriteDeadline)) are supported.

This enhancement is especially beneficial for applications that communicate via named pipes on Windows.

Note that a handle can only be associated with one completion port at a time.
If the handle provided to [`NewFile`](/pkg/os#NewFile) is already associated with a completion port,
the returned [`File`](/pkg/os#File) is downgraded to synchronous I/O mode.
In this case, I/O methods will block an OS thread, and the deadline methods have no effect.

The filesystems returned by [`DirFS`](/pkg/os#DirFS) and [`*Root.FS`](/pkg/os#Root.FS) implement the new [`io/fs.ReadLinkFS`](/pkg/io/fs#ReadLinkFS) interface.
[`CopyFS`](/pkg/os#CopyFS) supports symlinks when copying filesystems that implement [`io/fs.ReadLinkFS`](/pkg/io/fs#ReadLinkFS).

The [`os.Root`](/pkg/os#Root) type supports the following additional methods:

  * [`os.Root.Chmod`](/pkg/os#Root.Chmod)
  * [`os.Root.Chown`](/pkg/os#Root.Chown)
  * [`os.Root.Chtimes`](/pkg/os#Root.Chtimes)
  * [`os.Root.Lchown`](/pkg/os#Root.Lchown)
  * [`os.Root.Link`](/pkg/os#Root.Link)
  * [`os.Root.MkdirAll`](/pkg/os#Root.MkdirAll)
  * [`os.Root.ReadFile`](/pkg/os#Root.ReadFile)
  * [`os.Root.Readlink`](/pkg/os#Root.Readlink)
  * [`os.Root.RemoveAll`](/pkg/os#Root.RemoveAll)
  * [`os.Root.Rename`](/pkg/os#Root.Rename)
  * [`os.Root.Symlink`](/pkg/os#Root.Symlink)
  * [`os.Root.WriteFile`](/pkg/os#Root.WriteFile)

<!-- go.dev/issue/73126 is documented as part of 67002 -->

#### [`reflect`](/pkg/reflect/)

The new [`TypeAssert`](/pkg/reflect#TypeAssert) function permits converting a [`Value`](/pkg/reflect#Value) directly to a Go value
of the given type. This is like using a type assertion on the result of [`Value.Interface`](/pkg/reflect#Value.Interface),
but avoids unnecessary memory allocations.

#### [`regexp/syntax`](/pkg/regexp/syntax/)

The `\p{name}` and `\P{name}` character class syntaxes now accept the names
Any, ASCII, Assigned, Cn, and LC, as well as Unicode category aliases like `\p{Letter}` for `\pL`.
Following [Unicode TR18](https://unicode.org/reports/tr18/), they also now use
case-insensitive name lookups, ignoring spaces, underscores, and hyphens.

#### [`runtime`](/pkg/runtime/)

Cleanup functions scheduled by [`AddCleanup`](/pkg/runtime#AddCleanup) are now executed
concurrently and in parallel, making cleanups more viable for heavy
use like the [`unique`](/pkg/unique) package. Note that individual cleanups should
still shunt their work to a new goroutine if they must execute or
block for a long time to avoid blocking the cleanup queue.

When `GODEBUG=checkfinalizers=1` is set, the runtime will run
diagnostics on each garbage collection cycle to find common issues
with how the program might use finalizers and cleanups, such as those
described [in the GC
guide](/doc/gc-guide#Finalizers_cleanups_and_weak_pointers). In this
mode, the runtime will also regularly report the finalizer and
cleanup queue lengths to stderr to help identify issues with
long-running finalizers and/or cleanups.

The new [`SetDefaultGOMAXPROCS`](/pkg/runtime#SetDefaultGOMAXPROCS) function sets `GOMAXPROCS` to the runtime
default value, as if the `GOMAXPROCS` environment variable is not set. This is
useful for enabling the [new `GOMAXPROCS` default](#runtime) if it has been
disabled by the `GOMAXPROCS` environment variable or a prior call to
[`GOMAXPROCS`](/pkg/runtime#GOMAXPROCS).

#### [`runtime/pprof`](/pkg/runtime/pprof/)

The mutex profile for contention on runtime-internal locks now correctly points
to the end of the critical section that caused the delay. This matches the
profile's behavior for contention on `sync.Mutex` values. The
`runtimecontentionstacks` setting for `GODEBUG`, which allowed opting in to the
unusual behavior of Go 1.22 through 1.24 for this part of the profile, is now
gone.

#### [`runtime/trace`](/pkg/runtime/trace/)

<!-- go.dev/issue/63185 -->
The new [`FlightRecorder`](/pkg/runtime/trace#FlightRecorder) provides a
lightweight way to capture a trace of last few seconds of execution at a
specific moment in time. When a significant event occurs, a program may call
[`FlightRecorder.WriteTo`](/pkg/runtime/trac#FlightRecorder.WriteTo) to
snapshot available trace data. The length of time and amount of data captured
by a [`FlightRecorder`](/pkg/runtime/trace#FlightRecorder) may be configured
within the [`FlightRecorderConfig`](/pkg/runtime/trace#FlightRecorderConfig).

#### [`sync`](/pkg/sync/)

The new method on [`WaitGroup`](/pkg/sync#WaitGroup), [`WaitGroup.Go`](/pkg/sync#WaitGroup.Go),
makes the common pattern of creating and counting goroutines more convenient.

#### [`testing`](/pkg/testing/)

The new methods [`T.Attr`](/pkg/testing#T.Attr), [`B.Attr`](/pkg/testing#B.Attr), and [`F.Attr`](/pkg/testing#F.Attr) emit an
attribute to the test log. An attribute is an arbitrary
key and value associated with a test.

For example, in a test named `TestAttr`,
`t.Attr("key", "value")` emits:

```
=== ATTR  TestAttr key value
```

<!-- go.dev/issue/59928 -->

The new [`Output`](/pkg/testing#Output) method of [`testing.T`](/pkg/testing#T), [`testing.B`](/pkg/testing#B) and [`testing.F`](/pkg/testing#F) provides a Writer
that writes to the same test output stream as [`TB.Log`](/pkg/testing#TB.Log), but omits the file and line number.

#### [`testing/fstest`](/pkg/testing/fstest/)

[`MapFS`](/pkg/testing/fstest#MapFS) implements the new [`io/fs.ReadLinkFS`](/pkg/io/fs#ReadLinkFS) interface.
[`TestFS`](/pkg/testing/fstest#TestFS) will verify the functionality of the [`io/fs.ReadLinkFS`](/pkg/io/fs#ReadLinkFS) interface if implemented.
[`TestFS`](/pkg/testing/fstest#TestFS) will no longer follow symlinks to avoid unbounded recursion.

#### [`testing/synctest`](/pkg/testing/synctest/)

<!-- testing/synctest -->

#### [`unicode`](/pkg/unicode/)

The new [`CategoryAliases`](/pkg/unicode#CategoryAliases) map provides access to category alias names, such as “Letter” for “L”.
The new categories [`Cn`](/pkg/unicode#Cn) and [`LC`](/pkg/unicode#LC) define unassigned codepoints and cased letters, respectively.
These have always been defined by Unicode but were inadvertently omitted in earlier versions of Go.
The [`C`](/pkg/unicode#C) category now includes [`Cn`](/pkg/unicode#Cn), meaning it has added all unassigned code points.

#### [`unique`](/pkg/unique/)

The [`unique`](/pkg/unique) package now reclaims interned values more eagerly,
more efficiently, and in parallel. As a consequence, applications using
[`Make`](/pkg/unique#Make) are now less likely to experience memory blow-up when lots of
truly unique values are interned.

Values passed to [`Make`](/pkg/unique#Make) containing [Handle]s previously required multiple
garbage collection cycles to collect, proportional to the depth of the chain
of [`Handle`](/pkg/unique#Handle) values. Now, they are collected promptly in a single cycle, once
unused.

## Ports {#ports}

### Darwin

<!-- go.dev/issue/69839 -->
As [announced](/doc/go1.24#darwin) in the Go 1.24 release notes, Go 1.25 requires macOS 12 Monterey or later; support for previous versions has been discontinued.

### Windows

<!-- go.dev/issue/71671 -->
Go 1.25 is the last release that contains the [broken](/doc/go1.24#windows) 32-bit windows/arm port (`GOOS=windows` `GOARCH=arm`). It will be removed in Go 1.26.

### RISC-V

<!-- CL 420114 -->
The linux/riscv64 port now supports the `plugin` build mode.

<!--
Output from relnote todo that was generated and reviewed on 2025-05-23, plus summary info from bug/CL: -->

### TODO

**Please turn these into proper release notes**

<!-- TODO: CL 660996 has a RELNOTE comment without a suggested text (from RELNOTE comment in https://go.dev/cl/660996) -->
cmd/link/internal/ld: introduce -funcalign=N option\
This patch adds linker option -funcalign=N that allows to set alignment
for function entries.\
For \#72130.

<!-- TODO: accepted proposal https://go.dev/issue/32816 (from https://go.dev/cl/645155, https://go.dev/cl/645455, https://go.dev/cl/645955, https://go.dev/cl/646255, https://go.dev/cl/646455, https://go.dev/cl/646495, https://go.dev/cl/646655, https://go.dev/cl/646875, https://go.dev/cl/647298, https://go.dev/cl/647299, https://go.dev/cl/647736, https://go.dev/cl/648581, https://go.dev/cl/648715, https://go.dev/cl/648976, https://go.dev/cl/648995, https://go.dev/cl/649055, https://go.dev/cl/649056, https://go.dev/cl/649057, https://go.dev/cl/649456, https://go.dev/cl/649476, https://go.dev/cl/650755, https://go.dev/cl/651615, https://go.dev/cl/651617, https://go.dev/cl/651655, https://go.dev/cl/653436) -->
cmd/fix: automate migrations for simple deprecations

<!-- TODO: accepted proposal https://go.dev/issue/34055 (from https://go.dev/cl/625577) -->
cmd/go: allow serving module under the subdirectory of git repository\
cmd/go: add subdirectory support to go-import meta tag\
This CL adds ability to specify a subdirectory in the go-import meta tag.
A go-import meta tag now will support:
\<meta name="go-import" content="root-path vcs repo-url subdir">\
Fixes: \#34055

<!-- TODO: accepted proposal https://go.dev/issue/42965 (from https://go.dev/cl/643355, https://go.dev/cl/670656, https://go.dev/cl/670975, https://go.dev/cl/674076) -->
cmd/go: add global ignore mechanism for Go tooling ecosystem

<!-- TODO: accepted proposal https://go.dev/issue/51430 (from https://go.dev/cl/644997, https://go.dev/cl/646355) -->
cmd/cover: extend coverage testing to include applications

<!-- TODO: accepted proposal https://go.dev/issue/60905 (from https://go.dev/cl/645795) -->
all: add GOARM64=v8.1 and so on\
runtime: check LSE support on ARM64 at runtime init\
Check presence of LSE support on ARM64 chip if we targeted it at compile
time.\
Related to \#69124\
Updates \#60905\
Fixes \#71411

<!-- TODO: accepted proposal https://go.dev/issue/61476 (from https://go.dev/cl/633417) -->
all: add GORISCV64 environment variable\
cmd/go: add rva23u64 as a valid value for GORISCV64\
The RVA23 profile was ratified on the 21st of October 2024.
https://riscv.org/announcements/2024/10/risc-v-announces-ratification-of-the-rva23-profile-standard/
Now that it's ratified we can add rva23u64 as a valid value for the
GORISCV64 environment variable. This will allow the compiler and
assembler to generate instructions made mandatory by the new profile
without a runtime check.  Examples of such instructions include those
introduced by the Vector and Zicond extensions.
Setting GORISCV64=rva23u64 defines the riscv64.rva20u64,
riscv64.rva22u64 and riscv64.rva23u64 build tags, sets the internal
variable buildcfg.GORISCV64 to 23 and defines the macros
GORISCV64_rva23u64, hasV, hasZba, hasZbb, hasZbs, hasZfa, and
hasZicond for use in assembly language code.\
Updates \#61476

<!-- TODO: accepted proposal https://go.dev/issue/61716 (from https://go.dev/cl/644475) -->
math/rand/v2: revised API for math/rand\
rand: deprecate in favor of math/rand/v2\
For golang/go#61716\
Fixes golang/go#71373

<!-- TODO: accepted proposal https://go.dev/issue/64876 (from https://go.dev/cl/649435) -->
cmd/go: enable GOCACHEPROG by default\
cmd/go/internal/cacheprog: drop Request.ObjectID\
ObjectID was a misnaming of OutputID from cacheprog's initial
implementation. It was maintained for compatibility with existing
cacheprog users in 1.24 but can be removed in 1.25.

<!-- TODO: accepted proposal https://go.dev/issue/68106 (from https://go.dev/cl/628175, https://go.dev/cl/674158, https://go.dev/cl/674436, https://go.dev/cl/674437, https://go.dev/cl/674555, https://go.dev/cl/674556, https://go.dev/cl/674575, https://go.dev/cl/675075, https://go.dev/cl/675076, https://go.dev/cl/675155, https://go.dev/cl/675235) -->
cmd/go: doc -http should start a pkgsite instance and open a browser

<!-- TODO: accepted proposal https://go.dev/issue/69712 (from https://go.dev/cl/619955) -->
cmd/go: -json flag for go version -m\
cmd/go: support -json flag in go version\
It supports features described in the issue:
- add -json flag for 'go version -m' to print json encoding of
  runtime/debug.BuildSetting to standard output.
- report an error when specifying -json flag without -m.
- print build settings on seperated line for each binary\
  Fixes \#69712

<!-- TODO: accepted proposal https://go.dev/issue/70123 (from https://go.dev/cl/657116) -->
crypto: mechanism to enable FIPS mode

<!-- TODO: accepted proposal https://go.dev/issue/70128 (from https://go.dev/cl/645716, https://go.dev/cl/647455, https://go.dev/cl/651215, https://go.dev/cl/651256, https://go.dev/cl/652136, https://go.dev/cl/652215, https://go.dev/cl/653095, https://go.dev/cl/653139, https://go.dev/cl/653156, https://go.dev/cl/654395) -->
spec: remove notion of core types

<!-- TODO: accepted proposal https://go.dev/issue/70200 (from https://go.dev/cl/674916) -->
cmd/go: add fips140 module selection mechanism\
lib/fips140: set inprocess.txt to v1.0.0

<!-- TODO: accepted proposal https://go.dev/issue/70464 (from https://go.dev/cl/630137) -->
testing: panic in AllocsPerRun during parallel test\
testing: panic in AllocsPerRun if parallel tests are running\
If other tests are running, AllocsPerRun's result will be inherently flaky.
Saw this with CL 630136 and \#70327.\
Proposed in \#70464.\
Fixes \#70464.

<!-- TODO: accepted proposal https://go.dev/issue/71845 (from https://go.dev/cl/665796, https://go.dev/cl/666935) -->
encoding/json/v2: add new JSON API behind a GOEXPERIMENT=jsonv2 guard

<!-- TODO: accepted proposal https://go.dev/issue/71867 (from https://go.dev/cl/666476, https://go.dev/cl/666755, https://go.dev/cl/673119, https://go.dev/cl/673696) -->
cmd/go, cmd/distpack: build and run tools that are not necessary for builds as needed and don't include in binary distribution

<!-- Items that don't need to be mentioned in Go 1.25 release notes but are picked up by relnote todo

TODO: accepted proposal https://go.dev/issue/30999 (from https://go.dev/cl/671795)
net: reject leading zeros in IP address parsers
net: don't test with leading 0 in ipv4 addresses
Updates \#30999
Fixes \#73378

TODO: accepted proposal https://go.dev/issue/36532 (from https://go.dev/cl/647555)
testing: reconsider adding Context method to testing.T
database/sql: use t.Context in tests
Replace "context.WithCancel(context.Background())" with "t.Context()".
Updates \#36532

TODO: accepted proposal https://go.dev/issue/48429 (from https://go.dev/cl/648577)
cmd/go: track tool dependencies in go.mod
cmd/go: document -modfile and other flags for 'go tool'
Mention -modfile, -C, -overlay, and -modcacherw in the 'go tool'
documentation. We let a reference to 'go help build' give a pointer to
more detailed information.
The -modfile flag in particular is newly useful with the Go 1.24 support
for user-defined tools with 'go tool'.
Updates \#48429
Updates \#33926
Updates \#71663
Fixes \#71502

TODO: accepted proposal https://go.dev/issue/51572 (from https://go.dev/cl/651996)
cmd/go: add 'unix' build tag but not \*\_unix.go file support
os, syscall: use unix build tag where appropriate
These newly added files may use the unix build tag instead of explitly
listing all unix-like GOOS values.
For \#51572

TODO: accepted proposal https://go.dev/issue/53757 (from https://go.dev/cl/644575)
x/sync/errgroup: propagate panics and Goexits through Wait
errgroup: propagate panic and Goexit through Wait
Recovered panic values are wrapped and saved in Group.
Goexits are detected by a sentinel value set after the given function
returns normally. Wait propagates the first instance of a panic or
Goexit.
According to the runtime.Goexit after the code will not be executed,
with a bool, if f not call runtime.Goexit, is true,
determine whether to propagate runtime.Goexit.
Fixes golang/go#53757

TODO: accepted proposal https://go.dev/issue/54743 (from https://go.dev/cl/532415)
ssh: add server side support for Diffie Hellman Group Exchange

TODO: accepted proposal https://go.dev/issue/57792 (from https://go.dev/cl/649716, https://go.dev/cl/651737)
x/crypto/x509roots: new module

TODO: accepted proposal https://go.dev/issue/58523 (from https://go.dev/cl/538235)
ssh: expose negotiated algorithms
Fixes golang/go#58523
Fixes golang/go#46638

TODO: accepted proposal https://go.dev/issue/61537 (from https://go.dev/cl/531935)
ssh: export supported algorithms
Fixes golang/go#61537

TODO: accepted proposal https://go.dev/issue/61901 (from https://go.dev/cl/647875)
bytes, strings: add iterator forms of existing functions

TODO: accepted proposal https://go.dev/issue/61940 (from https://go.dev/cl/650235)
all: fix links to Go wiki
The Go wiki on GitHub has moved to go.dev in golang/go#61940.

TODO: accepted proposal https://go.dev/issue/64207 (from https://go.dev/cl/647015, https://go.dev/cl/652235)
all: end support for macOS 10.15 in Go 1.23

TODO: accepted proposal https://go.dev/issue/67839 (from https://go.dev/cl/646535)
x/sys/unix: access to ELF auxiliary vector
runtime: adjust comments for auxv getAuxv
github.com/cilium/ebpf no longer accesses getAuxv using linkname but now
uses the golang.org/x/sys/unix.Auxv wrapper introduced in
go.dev/cl/644295.
Also adjust the list of users to include x/sys/unix.
Updates \#67839
Updates \#67401

TODO: accepted proposal https://go.dev/issue/68780 (from https://go.dev/cl/659835)
x/term: support pluggable history
term: support pluggable history
Expose a new History interface that allows replacement of the default
ring buffer to customize what gets added or not; as well as to allow
saving/restoring history on either the default ringbuffer or a custom
replacement.
Fixes golang/go#68780

TODO: accepted proposal https://go.dev/issue/69095 (from https://go.dev/cl/649320, https://go.dev/cl/649321, https://go.dev/cl/649337, https://go.dev/cl/649376, https://go.dev/cl/649377, https://go.dev/cl/649378, https://go.dev/cl/649379, https://go.dev/cl/649380, https://go.dev/cl/649397, https://go.dev/cl/649398, https://go.dev/cl/649419, https://go.dev/cl/649497, https://go.dev/cl/649498, https://go.dev/cl/649618, https://go.dev/cl/649675, https://go.dev/cl/649676, https://go.dev/cl/649677, https://go.dev/cl/649695, https://go.dev/cl/649696, https://go.dev/cl/649697, https://go.dev/cl/649698, https://go.dev/cl/649715, https://go.dev/cl/649717, https://go.dev/cl/649718, https://go.dev/cl/649755, https://go.dev/cl/649775, https://go.dev/cl/649795, https://go.dev/cl/649815, https://go.dev/cl/649835, https://go.dev/cl/651336, https://go.dev/cl/651736, https://go.dev/cl/651737, https://go.dev/cl/658018)
all, x/build/cmd/relui: automate go directive maintenance in golang.org/x repositories

TODO: accepted proposal https://go.dev/issue/70859 (from https://go.dev/cl/666056, https://go.dev/cl/670835, https://go.dev/cl/672015, https://go.dev/cl/672016, https://go.dev/cl/672017)
x/tools/go/ast/inspector: add Cursor, to enable partial and multi-level traversals

-->
[cross-site request forgery (csrf)]: https://developer.mozilla.org/en-US/docs/Web/Security/Attacks/CSRF
[sec-fetch-site]: https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Sec-Fetch-Site
