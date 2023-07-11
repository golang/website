---
title: WASI support in Go
date: 2023-09-13
by:
- Johan Brandhorst-Satzkorn, Julien Fabre, Damian Gryski, Evan Phoenix, and Achille Roussel
summary: Go 1.21 adds a new port targeting the WASI preview 1 syscall API
---

Go 1.21 adds a new port targeting the WASI preview 1 syscall API through the
new `GOOS` value `wasip1`. This port builds on the existing WebAssembly
port introduced in Go 1.11.

## What is WebAssembly?

[WebAssembly (Wasm)](https://webassembly.org/) is a binary instruction format
originally designed for the web. It represents a standard that allows
developers to run high-performance, low-level code directly in web browsers at
near-native speeds.

Go first added support for compiling to Wasm in the 1.11 release, through the
`js/wasm` port. This allowed Go code compiled using the Go compiler to be
executed in web browsers, but it required a JavaScript execution environment.

As the use of Wasm has grown, so have use cases outside of the browser. Many
cloud providers are now offering services that allow the user to execute Wasm
executables directly, leveraging the new
[WebAssembly System Interface (WASI)](https://wasi.dev/) syscall API.

## The WebAssembly System Interface

WASI defines a syscall API for Wasm executables, allowing them to interact with
system resources such as the filesystem, the system clock, random data
utilities, and more. The latest release of the WASI spec is called
`wasi_snapshot_preview1`, from which we derive the `GOOS` name `wasip1`. New
versions of the API are being developed, and supporting them in the Go
compiler in the future will likely mean adding a new `GOOS`.

The creation of WASI has allowed a number of Wasm runtimes (hosts) to
standardize their syscall API around it. Examples of Wasm/WASI hosts include
[Wasmtime](https://wasmtime.dev), [Wazero](https://wazero.io/),
[WasmEdge](https://wasmedge.org/), [Wasmer](https://wasmer.io/), and
[NodeJS](https://nodejs.org). There are also a number of cloud providers
offering hosting of Wasm/WASI executables.

## How can we use it with Go?

Make sure that you have installed at least version 1.21 of Go. For this demo,
we’ll use [the Wasmtime host](https://docs.wasmtime.dev/cli-install.html) to
execute our binary. Let’s start with a simple `main.go`:

```go
package main

import "fmt"

func main() {
	fmt.Println("Hello world!")
}
```

We can build it for `wasip1` using the command:

```shell
$ GOOS=wasip1 GOARCH=wasm go build -o main.wasm main.go
```

This will produce a file, `main.wasm` which we can execute with `wasmtime`:

```shell
$ wasmtime main.wasm
Hello world!
```

That’s all it takes to get started with Wasm/WASI! You can expect almost all
the features of Go to just work with `wasip1`. To learn more about the details
of how WASI works with Go, please see
[the proposal](https://go.dev/issue/58141).

## Running go tests with wasip1

Building and running a binary is easy, but sometimes we want to be able to run
`go test` directly without having to build and execute the binary manually.
Similar to the `js/wasm` port, the standard library distribution included
in your Go installation comes with a file that makes this very easy. Add the
`misc/wasm` directory to your `PATH` when running Go tests and it will
run the tests using the Wasm host of your choice. This works by `go test`
[automatically executing](https://pkg.go.dev/cmd/go#hdr-Compile_and_run_Go_program)
`misc/wasm/go_wasip1_wasm_exec` when it finds this file in the `PATH`.

```shell
$ export PATH=$PATH:$(go env GOROOT)/misc/wasm
$ GOOS=wasip1 GOARCH=wasm go test ./...
```

This will run `go test` using Wasmtime. The Wasm host used can be controlled
using the environment variable `GOWASIRUNTIME`. Currently supported values
for this variable are `wazero`, `wasmedge`, `wasmtime`, and `wasmer`. This
script is subject to breaking changes between Go versions. Note that Go
`wasip1` binaries don’t execute perfectly on all hosts yet (see
[#59907](https://go.dev/issue/59907) and
[#60097](https://go.dev/issue/60097)).

This functionality also works when using `go run`:

```shell
$ GOOS=wasip1 GOARCH=wasm go run ./main.go
Hello world!
```

## Wrapping Wasm functions in Go with go:wasmimport

In addition to the new `wasip1/wasm` port, Go 1.21 introduces a new compiler
directive: `go:wasmimport`. It instructs the compiler to translate calls to
the annotated function into a call to the function specified by the host
module name and function name. This new compiler functionality is what allowed
us to define the `wasip1` syscall API in Go to support the new port, but it
isn’t limited to being used in the standard library.

For example, the wasip1 syscall API defines the
[`random_get` function](https://github.com/WebAssembly/WASI/blob/a51a66df5b1db01cf9e873f5537bc5bd552cf770/legacy/preview1/docs.md#-random_getbuf-pointeru8-buf_len-size---result-errno),
and it is exposed to the Go standard library through
[a function wrapper](https://cs.opensource.google/go/go/+/refs/tags/go1.21.0:src/runtime/os_wasip1.go;l=73-75)
defined in the runtime package. It looks like this:

```go
//go:wasmimport wasi_snapshot_preview1 random_get
//go:noescape
func random_get(buf unsafe.Pointer, bufLen size) errno
```

This function wrapper is then wrapped in
[a more ergonomic function](https://cs.opensource.google/go/go/+/refs/tags/go1.21.0:src/runtime/os_wasip1.go;l=183-187)
for use in the standard library:

```go
func getRandomData(r []byte) {
	if random_get(unsafe.Pointer(&r[0]), size(len(r))) != 0 {
		throw("random_get failed")
	}
}
```

This way, a user can call `getRandomData` with a byte slice and it will
eventually make its way to the host-defined `random_get` function. In the same
way, users can define their own wrappers for host functions.

To learn more about the intricacies of wrapping Wasm functions in Go, please
see [the `go:wasmimport` proposal](https://go.dev/issue/59149).

## Limitations

While the `wasip1` port passes all standard library tests, there are some
notable fundamental limitations of the Wasm architecture that may surprise
users.

Wasm is a single threaded architecture with no parallelism. The scheduler can
still schedule goroutines to run concurrently, and standard in/out/error is
non-blocking, so a goroutine can execute while another reads or writes, but any
host function calls (such as requesting random data using the example above)
will cause all goroutines to block until the host function call has returned.

A notable missing feature in the `wasip1` API is a full implementation of
network sockets. `wasip1` only defines functions that operate on already opened
sockets, making it impossible to support some of the most popular features of
the Go standard library, such as HTTP servers. Hosts like Wasmer and WasmEdge
implement extensions to the `wasip1` API, allowing the opening of network
sockets. While these extensions are not implemented by the Go compiler, there
exists a third party library,
[`github.com/stealthrocket/net`](https://github.com/stealthrocket/net), which
uses `go:wasmimport` to allow the use of `net.Dial` and `net.Listen` on
supported Wasm hosts. This enables the creation of `net/http` servers and other
network related functionality when using this package.

## The future of Wasm in Go

The addition of the `wasip1/wasm` port is just the beginning of the Wasm
capabilities we would like to bring to Go. Please keep an eye out on
[the issue tracker](https://github.com/golang/go/issues?q=is%3Aopen+is%3Aissue+label%3Aarch-wasm)
for proposals around exporting Go functions to Wasm (`go:wasmexport`), a 32-bit
port and future WASI API compatibility.

## Get involved

If you are experimenting with and want to contribute to Wasm and Go, please get
involved! The Go issue tracker tracks all in-progress work and the #webassembly
channel on [the Gophers Slack](https://invite.slack.golangbridge.org/) is a
great place to discuss Go and WebAssembly. We look forward to hearing from you!
