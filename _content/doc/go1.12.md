---
template: false
title: Go 1.12 Release Notes
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

## Introduction to Go 1.12 {#introduction}

The latest Go release, version 1.12, arrives six months after [Go 1.11](go1.11).
Most of its changes are in the implementation of the toolchain, runtime, and libraries.
As always, the release maintains the Go 1 [promise of compatibility](/doc/go1compat).
We expect almost all Go programs to continue to compile and run as before.

## Changes to the language {#language}

There are no changes to the language specification.

## Ports {#ports}

<!-- CL 138675 -->
The race detector is now supported on `linux/arm64`.

Go 1.12 is the last release that is supported on FreeBSD 10.x, which has
already reached end-of-life. Go 1.13 will require FreeBSD 11.2+ or FreeBSD
12.0+.
FreeBSD 12.0+ requires a kernel with the COMPAT\_FREEBSD11 option set (this is the default).

<!-- CL 146898 -->
cgo is now supported on `linux/ppc64`.

<!-- CL 146023 -->
`hurd` is now a recognized value for `GOOS`, reserved
for the GNU/Hurd system for use with `gccgo`.

### Windows {#windows}

Go's new `windows/arm` port supports running Go on Windows 10
IoT Core on 32-bit ARM chips such as the Raspberry Pi 3.

### AIX {#aix}

Go now supports AIX 7.2 and later on POWER8 architectures (`aix/ppc64`). External linking, cgo, pprof and the race detector aren't yet supported.

### Darwin {#darwin}

Go 1.12 is the last release that will run on macOS 10.10 Yosemite.
Go 1.13 will require macOS 10.11 El Capitan or later.

<!-- CL 141639 -->
`libSystem` is now used when making syscalls on Darwin,
ensuring forward-compatibility with future versions of macOS and iOS. <!-- CL 153338 -->
The switch to `libSystem` triggered additional App Store
checks for private API usage. Since it is considered private,
`syscall.Getdirentries` now always fails with
`ENOSYS` on iOS.
Additionally, [`syscall.Setrlimit`](/pkg/syscall/#Setrlimit)
reports `invalid` `argument` in places where it historically
succeeded. These consequences are not specific to Go and users should expect
behavioral parity with `libSystem`'s implementation going forward.

## Tools {#tools}

### `go tool vet` no longer supported {#vet}

The `go vet` command has been rewritten to serve as the
base for a range of different source code analysis tools. See
the [golang.org/x/tools/go/analysis](https://godoc.org/golang.org/x/tools/go/analysis)
package for details. A side-effect is that `go tool vet`
is no longer supported. External tools that use `go tool
  vet` must be changed to use `go
  vet`. Using `go vet` instead of `go tool
  vet` should work with all supported versions of Go.

As part of this change, the experimental `-shadow` option
is no longer available with `go vet`. Checking for
variable shadowing may now be done using

	go get -u golang.org/x/tools/go/analysis/passes/shadow/cmd/shadow
	go vet -vettool=$(which shadow)


### Tour {#tour}

<!-- CL 152657 -->
The Go tour is no longer included in the main binary distribution. To
run the tour locally, instead of running `go` `tool` `tour`,
manually install it:

	go get -u golang.org/x/tour
	tour


### Build cache requirement {#gocache}

The [build cache](/cmd/go/#hdr-Build_and_test_caching) is now
required as a step toward eliminating
`$GOPATH/pkg`. Setting the environment variable
`GOCACHE=off` will cause `go` commands that write to the
cache to fail.

### Binary-only packages {#binary-only}

Go 1.12 is the last release that will support binary-only packages.

### Cgo {#cgo}

Go 1.12 will translate the C type `EGLDisplay` to the Go type `uintptr`.
This change is similar to how Go 1.10 and newer treats Darwin's CoreFoundation
and Java's JNI types. See the
[cgo documentation](/cmd/cgo/#hdr-Special_cases)
for more information.

<!-- CL 152657 -->
Mangled C names are no longer accepted in packages that use Cgo. Use the Cgo
names instead. For example, use the documented cgo name `C.char`
rather than the mangled name `_Ctype_char` that cgo generates.

### Modules {#modules}

<!-- CL 148517 -->
When `GO111MODULE` is set to `on`, the `go`
command now supports module-aware operations outside of a module directory,
provided that those operations do not need to resolve import paths relative to
the current directory or explicitly edit the `go.mod` file.
Commands such as `go`&nbsp;`get`,
`go`&nbsp;`list`, and
`go`&nbsp;`mod`&nbsp;`download` behave as if in a
module with initially-empty requirements.
In this mode, `go`&nbsp;`env`&nbsp;`GOMOD` reports
the system's null device (`/dev/null` or `NUL`).

<!-- CL 146382 -->
`go` commands that download and extract modules are now safe to
invoke concurrently.
The module cache (`GOPATH/pkg/mod`) must reside in a filesystem that
supports file locking.

<!-- CL 147282, 147281 -->
The `go` directive in a `go.mod` file now indicates the
version of the language used by the files within that module.
It will be set to the current release
(`go`&nbsp;`1.12`) if no existing version is
present.
If the `go` directive for a module specifies a
version _newer_ than the toolchain in use, the `go` command
will attempt to build the packages regardless, and will note the mismatch only if
that build fails.

<!-- CL 147282, 147281 -->
This changed use of the `go` directive means that if you
use Go 1.12 to build a module, thus recording `go 1.12`
in the `go.mod` file, you will get an error when
attempting to build the same module with Go 1.11 through Go 1.11.3.
Go 1.11.4 or later will work fine, as will releases older than Go 1.11.
If you must use Go 1.11 through 1.11.3, you can avoid the problem by
setting the language version to 1.11, using the Go 1.12 go tool,
via `go mod edit -go=1.11`.

<!-- CL 152739 -->
When an import cannot be resolved using the active modules,
the `go` command will now try to use the modules mentioned in the
main module's `replace` directives before consulting the module
cache and the usual network sources.
If a matching replacement is found but the `replace` directive does
not specify a version, the `go` command uses a pseudo-version
derived from the zero `time.Time` (such
as `v0.0.0-00010101000000-000000000000`).

### Compiler toolchain {#compiler}

<!-- CL 134155, 134156 -->
The compiler's live variable analysis has improved. This may mean that
finalizers will be executed sooner in this release than in previous
releases. If that is a problem, consider the appropriate addition of a
[`runtime.KeepAlive`](/pkg/runtime/#KeepAlive) call.

<!-- CL 147361 -->
More functions are now eligible for inlining by default, including
functions that do nothing but call another function.
This extra inlining makes it additionally important to use
[`runtime.CallersFrames`](/pkg/runtime/#CallersFrames)
instead of iterating over the result of
[`runtime.Callers`](/pkg/runtime/#Callers) directly.

```
// Old code which no longer works correctly (it will miss inlined call frames).
var pcs [10]uintptr
n := runtime.Callers(1, pcs[:])
for _, pc := range pcs[:n] {
	f := runtime.FuncForPC(pc)
	if f != nil {
		fmt.Println(f.Name())
	}
}
```

```
// New code which will work correctly.
var pcs [10]uintptr
n := runtime.Callers(1, pcs[:])
frames := runtime.CallersFrames(pcs[:n])
for {
	frame, more := frames.Next()
	fmt.Println(frame.Function)
	if !more {
		break
	}
}
```

<!-- CL 153477 -->
Wrappers generated by the compiler to implement method expressions
are no longer reported
by [`runtime.CallersFrames`](/pkg/runtime/#CallersFrames)
and [`runtime.Stack`](/pkg/runtime/#Stack). They
are also not printed in panic stack traces.
This change aligns the `gc` toolchain to match
the `gccgo` toolchain, which already elided such wrappers
from stack traces.
Clients of these APIs might need to adjust for the missing
frames. For code that must interoperate between 1.11 and 1.12
releases, you can replace the method expression `x.M`
with the function literal ` func (...) { x.M(...) }  `.

<!-- CL 144340 -->
The compiler now accepts a `-lang` flag to set the Go language
version to use. For example, `-lang=go1.8` causes the compiler to
emit an error if the program uses type aliases, which were added in Go 1.9.
Language changes made before Go 1.12 are not consistently enforced.

<!-- CL 147160 -->
The compiler toolchain now uses different conventions to call Go
functions and assembly functions. This should be invisible to users,
except for calls that simultaneously cross between Go and
assembly _and_ cross a package boundary. If linking results
in an error like "relocation target not defined for ABIInternal (but
is defined for ABI0)", please refer to the
[compatibility section](https://github.com/golang/proposal/blob/master/design/27539-internal-abi.md#compatibility)
of the ABI design document.

<!-- CL 145179 -->
There have been many improvements to the DWARF debug information
produced by the compiler, including improvements to argument
printing and variable location information.

<!-- CL 61511 -->
Go programs now also maintain stack frame pointers on `linux/arm64`
for the benefit of profiling tools like `perf`. The frame pointer
maintenance has a small run-time overhead that varies but averages around 3%.
To build a toolchain that does not use frame pointers, set
`GOEXPERIMENT=noframepointer` when running `make.bash`.

<!-- CL 142717 -->
The obsolete "safe" compiler mode (enabled by the `-u` gcflag) has been removed.

### `godoc` and `go` `doc` {#godoc}

In Go 1.12, `godoc` no longer has a command-line interface and
is only a web server. Users should use `go` `doc`
for command-line help output instead. Go 1.12 is the last release that will
include the `godoc` webserver; in Go 1.13 it will be available
via `go` `get`.

<!-- CL 141977 -->
`go` `doc` now supports the `-all` flag,
which will cause it to print all exported APIs and their documentation,
as the `godoc` command line used to do.

<!-- CL 140959 -->
`go` `doc` also now includes the `-src` flag,
which will show the target's source code.

### Trace {#trace}

<!-- CL 60790 -->
The trace tool now supports plotting mutator utilization curves,
including cross-references to the execution trace. These are useful
for analyzing the impact of the garbage collector on application
latency and throughput.

### Assembler {#assembler}

<!-- CL 147218 -->
On `arm64`, the platform register was renamed from
`R18` to `R18_PLATFORM` to prevent accidental
use, as the OS could choose to reserve this register.

## Runtime {#runtime}

<!-- CL 138959 -->
Go 1.12 significantly improves the performance of sweeping when a
large fraction of the heap remains live. This reduces allocation
latency immediately following a garbage collection.

<!-- CL 139719 -->
The Go runtime now releases memory back to the operating system more
aggressively, particularly in response to large allocations that
can't reuse existing heap space.

<!-- CL 146342, CL 146340, CL 146345, CL 146339, CL 146343, CL 146337, CL 146341, CL 146338 -->
The Go runtime's timer and deadline code is faster and scales better
with higher numbers of CPUs. In particular, this improves the
performance of manipulating network connection deadlines.

<!-- CL 135395 -->
On Linux, the runtime now uses `MADV_FREE` to release unused
memory. This is more efficient but may result in higher reported
RSS. The kernel will reclaim the unused data when it is needed.
To revert to the Go 1.11 behavior (`MADV_DONTNEED`), set the
environment variable `GODEBUG=madvdontneed=1`.

<!-- CL 149578 -->
Adding cpu._extension_=off to the
[GODEBUG](/doc/diagnostics.html#godebug) environment
variable now disables the use of optional CPU instruction
set extensions in the standard library and runtime. This is not
yet supported on Windows.

<!-- CL 158337 -->
Go 1.12 improves the accuracy of memory profiles by fixing
overcounting of large heap allocations.

<!-- CL 159717 -->
Tracebacks, `runtime.Caller`,
and `runtime.Callers` no longer include
compiler-generated initialization functions. Doing a traceback
during the initialization of a global variable will now show a
function named `PKG.init.ializers`.

## Standard library {#library}

### TLS 1.3 {#tls_1_3}

Go 1.12 adds opt-in support for TLS 1.3 in the `crypto/tls` package as
specified by [RFC 8446](https://www.rfc-editor.org/info/rfc8446). It can
be enabled by adding the value `tls13=1` to the `GODEBUG`
environment variable. It will be enabled by default in Go 1.13.

To negotiate TLS 1.3, make sure you do not set an explicit `MaxVersion` in
[`Config`](/pkg/crypto/tls/#Config) and run your program with
the environment variable `GODEBUG=tls13=1` set.

All TLS 1.2 features except `TLSUnique` in
[`ConnectionState`](/pkg/crypto/tls/#ConnectionState)
and renegotiation are available in TLS 1.3 and provide equivalent or
better security and performance. Note that even though TLS 1.3 is backwards
compatible with previous versions, certain legacy systems might not work
correctly when attempting to negotiate it. RSA certificate keys too small
to be secure (including 512-bit keys) will not work with TLS 1.3.

TLS 1.3 cipher suites are not configurable. All supported cipher suites are
safe, and if `PreferServerCipherSuites` is set in
[`Config`](/pkg/crypto/tls/#Config) the preference order
is based on the available hardware.

Early data (also called "0-RTT mode") is not currently supported as a
client or server. Additionally, a Go 1.12 server does not support skipping
unexpected early data if a client sends it. Since TLS 1.3 0-RTT mode
involves clients keeping state regarding which servers support 0-RTT,
a Go 1.12 server cannot be part of a load-balancing pool where some other
servers do support 0-RTT. If switching a domain from a server that supported
0-RTT to a Go 1.12 server, 0-RTT would have to be disabled for at least the
lifetime of the issued session tickets before the switch to ensure
uninterrupted operation.

In TLS 1.3 the client is the last one to speak in the handshake, so if it causes
an error to occur on the server, it will be returned on the client by the first
[`Read`](/pkg/crypto/tls/#Conn.Read), not by
[`Handshake`](/pkg/crypto/tls/#Conn.Handshake). For
example, that will be the case if the server rejects the client certificate.
Similarly, session tickets are now post-handshake messages, so are only
received by the client upon its first
[`Read`](/pkg/crypto/tls/#Conn.Read).

### Minor changes to the library {#minor_library_changes}

As always, there are various minor changes and updates to the library,
made with the Go 1 [promise of compatibility](/doc/go1compat)
in mind.

<!-- TODO: CL 115677: https://golang.org/cl/115677: cmd/vet: check embedded field tags too -->

#### [bufio](/pkg/bufio/)

<!-- CL 149297 -->
`Reader`'s [`UnreadRune`](/pkg/bufio/#Reader.UnreadRune) and
[`UnreadByte`](/pkg/bufio/#Reader.UnreadByte) methods will now return an error
if they are called after [`Peek`](/pkg/bufio/#Reader.Peek).

<!-- bufio -->

#### [bytes](/pkg/bytes/)

<!-- CL 137855 -->
The new function [`ReplaceAll`](/pkg/bytes/#ReplaceAll) returns a copy of
a byte slice with all non-overlapping instances of a value replaced by another.

<!-- CL 145098 -->
A pointer to a zero-value [`Reader`](/pkg/bytes/#Reader) is now
functionally equivalent to [`NewReader`](/pkg/bytes/#NewReader)`(nil)`.
Prior to Go 1.12, the former could not be used as a substitute for the latter in all cases.

<!-- bytes -->

#### [crypto/rand](/pkg/crypto/rand/)

<!-- CL 139419 -->
A warning will now be printed to standard error the first time
`Reader.Read` is blocked for more than 60 seconds waiting
to read entropy from the kernel.

<!-- CL 120055 -->
On FreeBSD, `Reader` now uses the `getrandom`
system call if available, `/dev/urandom` otherwise.

<!-- crypto/rand -->

#### [crypto/rc4](/pkg/crypto/rc4/)

<!-- CL 130397 -->
This release removes the assembly implementations, leaving only
the pure Go version. The Go compiler generates code that is
either slightly better or slightly worse, depending on the exact
CPU. RC4 is insecure and should only be used for compatibility
with legacy systems.

<!-- crypto/rc4 -->

#### [crypto/tls](/pkg/crypto/tls/)

<!-- CL 143177 -->
If a client sends an initial message that does not look like TLS, the server
will no longer reply with an alert, and it will expose the underlying
`net.Conn` in the new field `Conn` of
[`RecordHeaderError`](/pkg/crypto/tls/#RecordHeaderError).

<!-- crypto/tls -->

#### [database/sql](/pkg/database/sql/)

<!-- CL 145738 -->
A query cursor can now be obtained by passing a
[`*Rows`](/pkg/database/sql/#Rows)
value to the [`Row.Scan`](/pkg/database/sql/#Row.Scan) method.

<!-- database/sql -->

#### [expvar](/pkg/expvar/)

<!-- CL 139537 -->
The new [`Delete`](/pkg/expvar/#Map.Delete) method allows
for deletion of key/value pairs from a [`Map`](/pkg/expvar/#Map).

<!-- expvar -->

#### [fmt](/pkg/fmt/)

<!-- CL 142737 -->
Maps are now printed in key-sorted order to ease testing. The ordering rules are:

  - When applicable, nil compares low
  - ints, floats, and strings order by <
  - NaN compares less than non-NaN floats
  - bool compares false before true
  - Complex compares real, then imaginary
  - Pointers compare by machine address
  - Channel values compare by machine address
  - Structs compare each field in turn
  - Arrays compare each element in turn
  - Interface values compare first by `reflect.Type` describing the concrete type
    and then by concrete value as described in the previous rules.


<!-- CL 129777 -->
When printing maps, non-reflexive key values like `NaN` were previously
displayed as `<nil>`. As of this release, the correct values are printed.

<!-- fmt -->

#### [go/doc](/pkg/go/doc/)

<!-- CL 140958 -->
To address some outstanding issues in [`cmd/doc`](/cmd/doc/),
this package has a new [`Mode`](/pkg/go/doc/#Mode) bit,
`PreserveAST`, which controls whether AST data is cleared.

<!-- go/doc -->

#### [go/token](/pkg/go/token/)

<!-- CL 134075 -->
The [`File`](/pkg/go/token#File) type has a new
[`LineStart`](/pkg/go/token#File.LineStart) field,
which returns the position of the start of a given line. This is especially useful
in programs that occasionally handle non-Go files, such as assembly, but wish to use
the `token.Pos` mechanism to identify file positions.

<!-- go/token -->

#### [image](/pkg/image/)

<!-- CL 118755 -->
The [`RegisterFormat`](/pkg/image/#RegisterFormat) function is now safe for concurrent use.

<!-- image -->

#### [image/png](/pkg/image/png/)

<!-- CL 134235 -->
Paletted images with fewer than 16 colors now encode to smaller outputs.

<!-- image/png -->

#### [io](/pkg/io/)

<!-- CL 139457 -->
The new [`StringWriter`](/pkg/io#StringWriter) interface wraps the
[`WriteString`](/pkg/io/#WriteString) function.

<!-- io -->

#### [math](/pkg/math/)

<!-- CL 153059 -->
The functions
[`Sin`](/pkg/math/#Sin),
[`Cos`](/pkg/math/#Cos),
[`Tan`](/pkg/math/#Tan),
and [`Sincos`](/pkg/math/#Sincos) now
apply Payne-Hanek range reduction to huge arguments. This
produces more accurate answers, but they will not be bit-for-bit
identical with the results in earlier releases.

<!-- math -->

#### [math/bits](/pkg/math/bits/)

<!-- CL 123157 -->
New extended precision operations [`Add`](/pkg/math/bits/#Add), [`Sub`](/pkg/math/bits/#Sub), [`Mul`](/pkg/math/bits/#Mul), and [`Div`](/pkg/math/bits/#Div) are available in `uint`, `uint32`, and `uint64` versions.

<!-- math/bits -->

#### [net](/pkg/net/)

<!-- CL 146659 -->
The
[`Dialer.DualStack`](/pkg/net/#Dialer.DualStack) setting is now ignored and deprecated;
RFC 6555 Fast Fallback ("Happy Eyeballs") is now enabled by default. To disable, set
[`Dialer.FallbackDelay`](/pkg/net/#Dialer.FallbackDelay) to a negative value.

<!-- CL 107196 -->
Similarly, TCP keep-alives are now enabled by default if
[`Dialer.KeepAlive`](/pkg/net/#Dialer.KeepAlive) is zero.
To disable, set it to a negative value.

<!-- CL 113997 -->
On Linux, the [`splice` system call](https://man7.org/linux/man-pages/man2/splice.2.html) is now used when copying from a
[`UnixConn`](/pkg/net/#UnixConn) to a
[`TCPConn`](/pkg/net/#TCPConn).

<!-- net -->

#### [net/http](/pkg/net/http/)

<!-- CL 143177 -->
The HTTP server now rejects misdirected HTTP requests to HTTPS servers with a plaintext "400 Bad Request" response.

<!-- CL 130115 -->
The new [`Client.CloseIdleConnections`](/pkg/net/http/#Client.CloseIdleConnections)
method calls the `Client`'s underlying `Transport`'s `CloseIdleConnections`
if it has one.

<!-- CL 145398 -->
The [`Transport`](/pkg/net/http/#Transport) no longer rejects HTTP responses which declare
HTTP Trailers but don't use chunked encoding. Instead, the declared trailers are now just ignored.

<!-- CL 152080 -->
<!-- CL 151857 -->
The [`Transport`](/pkg/net/http/#Transport) no longer handles `MAX_CONCURRENT_STREAMS` values
advertised from HTTP/2 servers as strictly as it did during Go 1.10 and Go 1.11. The default behavior is now back
to how it was in Go 1.9: each connection to a server can have up to `MAX_CONCURRENT_STREAMS` requests
active and then new TCP connections are created as needed. In Go 1.10 and Go 1.11 the `http2` package
would block and wait for requests to finish instead of creating new connections.
To get the stricter behavior back, import the
[`golang.org/x/net/http2`](https://godoc.org/golang.org/x/net/http2) package
directly and set
[`Transport.StrictMaxConcurrentStreams`](https://godoc.org/golang.org/x/net/http2#Transport.StrictMaxConcurrentStreams) to
`true`.

<!-- net/http -->

#### [net/url](/pkg/net/url/)

<!-- CL 159157, CL 160178 -->
[`Parse`](/pkg/net/url/#Parse),
[`ParseRequestURI`](/pkg/net/url/#ParseRequestURI),
and
[`URL.Parse`](/pkg/net/url/#URL.Parse)
now return an
error for URLs containing ASCII control characters, which includes NULL,
tab, and newlines.

<!-- net/url -->

#### [net/http/httputil](/pkg/net/http/httputil/)

<!-- CL 146437 -->
The [`ReverseProxy`](/pkg/net/http/httputil/#ReverseProxy) now automatically
proxies WebSocket requests.

<!-- net/http/httputil -->

#### [os](/pkg/os/)

<!-- CL 125443 -->
The new [`ProcessState.ExitCode`](/pkg/os/#ProcessState.ExitCode) method
returns the process's exit code.

<!-- CL 135075 -->
`ModeCharDevice` has been added to the `ModeType` bitmask, allowing for
`ModeDevice | ModeCharDevice` to be recovered when masking a
[`FileMode`](/pkg/os/#FileMode) with `ModeType`.

<!-- CL 139418 -->
The new function [`UserHomeDir`](/pkg/os/#UserHomeDir) returns the
current user's home directory.

<!-- CL 146020 -->
[`RemoveAll`](/pkg/os/#RemoveAll) now supports paths longer than 4096 characters
on most Unix systems.

<!-- CL 130676 -->
[`File.Sync`](/pkg/os/#File.Sync) now uses `F_FULLFSYNC` on macOS
to correctly flush the file contents to permanent storage.
This may cause the method to run more slowly than in previous releases.

<!--CL 155517 -->
[`File`](/pkg/os/#File) now supports
a [`SyscallConn`](/pkg/os/#File.SyscallConn)
method returning
a [`syscall.RawConn`](/pkg/syscall/#RawConn)
interface value. This may be used to invoke system-specific
operations on the underlying file descriptor.

<!-- os -->

#### [path/filepath](/pkg/path/filepath/)

<!-- CL 145220 -->
The [`IsAbs`](/pkg/path/filepath/#IsAbs) function now returns true when passed
a reserved filename on Windows such as `NUL`.
[List of reserved names.](https://docs.microsoft.com/en-us/windows/desktop/fileio/naming-a-file#naming-conventions)

<!-- path/filepath -->

#### [reflect](/pkg/reflect/)

<!-- CL 33572 -->
A new [`MapIter`](/pkg/reflect#MapIter) type is
an iterator for ranging over a map. This type is exposed through the
[`Value`](/pkg/reflect#Value) type's new
[`MapRange`](/pkg/reflect#Value.MapRange) method.
This follows the same iteration semantics as a range statement, with `Next`
to advance the iterator, and `Key`/`Value` to access each entry.

<!-- reflect -->

#### [regexp](/pkg/regexp/)

<!-- CL 139784 -->
[`Copy`](/pkg/regexp/#Regexp.Copy) is no longer necessary
to avoid lock contention, so it has been given a partial deprecation comment.
[`Copy`](/pkg/regexp/#Regexp.Copy)
may still be appropriate if the reason for its use is to make two copies with
different [`Longest`](/pkg/regexp/#Regexp.Longest) settings.

<!-- regexp -->

#### [runtime/debug](/pkg/runtime/debug/)

<!-- CL 144220 -->
A new [`BuildInfo`](/pkg/runtime/debug/#BuildInfo) type
exposes the build information read from the running binary, available only in
binaries built with module support. This includes the main package path, main
module information, and the module dependencies. This type is given through the
[`ReadBuildInfo`](/pkg/runtime/debug/#ReadBuildInfo) function
on [`BuildInfo`](/pkg/runtime/debug/#BuildInfo).

<!-- runtime/debug -->

#### [strings](/pkg/strings/)

<!-- CL 137855 -->
The new function [`ReplaceAll`](/pkg/strings/#ReplaceAll) returns a copy of
a string with all non-overlapping instances of a value replaced by another.

<!-- CL 145098 -->
A pointer to a zero-value [`Reader`](/pkg/strings/#Reader) is now
functionally equivalent to [`NewReader`](/pkg/strings/#NewReader)`(nil)`.
Prior to Go 1.12, the former could not be used as a substitute for the latter in all cases.

<!-- CL 122835 -->
The new [`Builder.Cap`](/pkg/strings/#Builder.Cap) method returns the capacity of the builder's underlying byte slice.

<!-- CL 131495 -->
The character mapping functions [`Map`](/pkg/strings/#Map),
[`Title`](/pkg/strings/#Title),
[`ToLower`](/pkg/strings/#ToLower),
[`ToLowerSpecial`](/pkg/strings/#ToLowerSpecial),
[`ToTitle`](/pkg/strings/#ToTitle),
[`ToTitleSpecial`](/pkg/strings/#ToTitleSpecial),
[`ToUpper`](/pkg/strings/#ToUpper), and
[`ToUpperSpecial`](/pkg/strings/#ToUpperSpecial)
now always guarantee to return valid UTF-8. In earlier releases, if the input was invalid UTF-8 but no character replacements
needed to be applied, these routines incorrectly returned the invalid UTF-8 unmodified.

<!-- strings -->

#### [syscall](/pkg/syscall/)

<!-- CL 138595 -->
64-bit inodes are now supported on FreeBSD 12. Some types have been adjusted accordingly.

<!-- CL 125456 -->
The Unix socket
([`AF_UNIX`](https://blogs.msdn.microsoft.com/commandline/2017/12/19/af_unix-comes-to-windows/))
address family is now supported for compatible versions of Windows.

<!-- CL 147117 -->
The new function [`Syscall18`](/pkg/syscall/?GOOS=windows&GOARCH=amd64#Syscall18)
has been introduced for Windows, allowing for calls with up to 18 arguments.

<!-- syscall -->

#### [syscall/js](/pkg/syscall/js/)

<!-- CL 153559 -->

The `Callback` type and `NewCallback` function have been renamed;
they are now called
[`Func`](/pkg/syscall/js/?GOOS=js&GOARCH=wasm#Func) and
[`FuncOf`](/pkg/syscall/js/?GOOS=js&GOARCH=wasm#FuncOf), respectively.
This is a breaking change, but WebAssembly support is still experimental
and not yet subject to the
[Go 1 compatibility promise](/doc/go1compat). Any code using the
old names will need to be updated.

<!-- CL 141644 -->
If a type implements the new
[`Wrapper`](/pkg/syscall/js/?GOOS=js&GOARCH=wasm#Wrapper)
interface,
[`ValueOf`](/pkg/syscall/js/?GOOS=js&GOARCH=wasm#ValueOf)
will use it to return the JavaScript value for that type.

<!-- CL 143137 -->
The meaning of the zero
[`Value`](/pkg/syscall/js/?GOOS=js&GOARCH=wasm#Value)
has changed. It now represents the JavaScript `undefined` value
instead of the number zero.
This is a breaking change, but WebAssembly support is still experimental
and not yet subject to the
[Go 1 compatibility promise](/doc/go1compat). Any code relying on
the zero [`Value`](/pkg/syscall/js/?GOOS=js&GOARCH=wasm#Value)
to mean the number zero will need to be updated.

<!-- CL 144384 -->
The new
[`Value.Truthy`](/pkg/syscall/js/?GOOS=js&GOARCH=wasm#Value.Truthy)
method reports the
[JavaScript "truthiness"](https://developer.mozilla.org/en-US/docs/Glossary/Truthy)
of a given value.

<!-- syscall/js -->

#### [testing](/pkg/testing/)

<!-- CL 139258 -->
The [`-benchtime`](/cmd/go/#hdr-Testing_flags) flag now supports setting an explicit iteration count instead of a time when the value ends with an "`x`". For example, `-benchtime=100x` runs the benchmark 100 times.

<!-- testing -->

#### [text/template](/pkg/text/template/)

<!-- CL 142217 -->
When executing a template, long context values are no longer truncated in errors.

`executing "tmpl" at <.very.deep.context.v...>: map has no entry for key "notpresent"`

is now

`executing "tmpl" at <.very.deep.context.value.notpresent>: map has no entry for key "notpresent"`

<!-- CL 143097 -->
If a user-defined function called by a template panics, the
panic is now caught and returned as an error by
the `Execute` or `ExecuteTemplate` method.

<!-- text/template -->

#### [time](/pkg/time/)

<!-- CL 151299 -->
The time zone database in `$GOROOT/lib/time/zoneinfo.zip`
has been updated to version 2018i. Note that this ZIP file is
only used if a time zone database is not provided by the operating
system.

<!-- time -->

#### [unsafe](/pkg/unsafe/)

<!-- CL 146058 -->
It is invalid to convert a nil `unsafe.Pointer` to `uintptr` and back with arithmetic.
(This was already invalid, but will now cause the compiler to misbehave.)

<!-- unsafe -->
