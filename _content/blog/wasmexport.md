---
title: Extensible Wasm Applications with Go
date: 2025-02-13
by:
- Cherry Mui
summary: Go 1.24 enhances WebAssembly capabilities with function export and reactor mode
---

Go 1.24 enhances its WebAssembly (Wasm) capabilities with the
addition of the `go:wasmexport` directive and the ability to build a reactor
for WebAssembly System Interface (WASI).
These features enable Go developers to export Go functions to Wasm,
facilitating better integration with Wasm hosts and expanding the possibilities
for Go-based Wasm applications.

## WebAssembly and the WebAssembly System Interface

[WebAssembly (Wasm)](https://webassembly.org/) is a binary instruction format
that was initially created for web browsers, providing the execution of
high-performance, low-level code at speeds approaching native performance.
Since then, Wasm's utility has expanded, and it is now used in various
environments beyond the browser.
Notably, cloud providers offer services that directly execute Wasm
executables, taking advantage of the
[WebAssembly System Interface (WASI)](https://wasi.dev/) system call API.
WASI allows these executables to interact with system resources.

Go first added support for compiling to Wasm in the 1.11 release, through the
`js/wasm` port.
Go 1.21 added a new port targeting the WASI preview 1 syscall API through the
new `GOOS=wasip1` port.

## Exporting Go Functions to Wasm with `go:wasmexport`

Go 1.24 introduces a new compiler directive, `go:wasmexport`, which allows
developers to export Go functions to be called from outside of the
Wasm module, typically from a host application that runs the Wasm runtime.
This directive instructs the compiler to make the annotated function available
as a Wasm [export](https://webassembly.github.io/spec/core/valid/modules.html?highlight=export#exports)
in the resulting Wasm binary.

To use the `go:wasmexport` directive, simply add it to a function definition:

```
//go:wasmexport add
func add(a, b int32) int32 { return a + b }
```

With this, the Wasm module will have an exported function named `add` that
can be called from the host.

This is analogous to the [cgo `export` directive](/cmd/cgo#hdr-C_references_to_Go),
which makes the function available to be called from C,
though `go:wasmexport` uses a different, simpler mechanism.

## Building a WASI Reactor

A WASI reactor is a WebAssembly module that operates continuously, and
can be called upon multiple times to react on events or requests.
Unlike a "command" module, which terminates after its main function finishes,
a reactor instance remains live after initialization, and its exports remain
accessible.

With Go 1.24, one can build a WASI reactor with the `-buildmode=c-shared` build
flag.

```
$ GOOS=wasip1 GOARCH=wasm go build -buildmode=c-shared -o reactor.wasm
```

The build flag signals to the linker not to generate the `_start` function
(the entry point for a command module), and instead generate an
`_initialize` function, which performs runtime and package initialization,
along with any exported functions and their dependencies.
The `_initialize` function must be called before any other exported functions.
The `main` function will not be automatically invoked.

To use a WASI reactor, the host application first initializes it by calling
`_initialize`, then simply invoke the exported functions.
Here is an example using [Wazero](https://wazero.io/), a Go-based Wasm runtime
implementation:

```
// Create a Wasm runtime, set up WASI.
r := wazero.NewRuntime(ctx)
defer r.Close(ctx)
wasi_snapshot_preview1.MustInstantiate(ctx, r)

// Configure the module to initialize the reactor.
config := wazero.NewModuleConfig().WithStartFunctions("_initialize")

// Instantiate the module.
wasmModule, _ := r.InstantiateWithConfig(ctx, wasmFile, config)

// Call the exported function.
fn := wasmModule.ExportedFunction("add")
var a, b int32 = 1, 2
res, _ := fn.Call(ctx, api.EncodeI32(a), api.EncodeI32(b))
c := api.DecodeI32(res[0])
fmt.Printf("add(%d, %d) = %d\n", a, b, c)

// The instance is still alive. We can call the function again.
res, _ = fn.Call(ctx, api.EncodeI32(b), api.EncodeI32(c))
fmt.Printf("add(%d, %d) = %d\n", b, c, api.DecodeI32(res[0]))
```

The `go:wasmexport` directive and the reactor build mode allow applications to
be extended by calling into Go-based Wasm code.
This is particularly valuable for applications that have adopted Wasm as a
plugin or extension mechanism with well-defined interfaces.
By exporting Go functions, applications can leverage the Go Wasm modules to
provide functionality without needing to recompile the entire application.
Furthermore, building as a reactor ensures that the exported functions can be
called multiple times without requiring reinitialization, making it suitable
for long-running applications or services.

## Supporting rich types between the host and the client

Go 1.24 also relaxes the constraints on types that can be used as input and
result parameters with `go:wasmimport` functions.
For example, one can pass a bool, a string, a pointer to an `int32`, or a
pointer to a struct which embeds `structs.HostLayout` and contains supported
field types
(see the [documentation](/cmd/compile#hdr-WebAssembly_Directives) for detail).
This allows Go Wasm applications to be written in a more natural and ergonomic
way, and removes some unnecessary type conversions.

## Limitations

While Go 1.24 has made significant enhancements to its Wasm capabilities,
there are still some notable limitations.

Wasm is a single-threaded architecture with no parallelism.
A `go:wasmexport` function can spawn new goroutines.
But if a function creates a background goroutine, it will not continue
executing when the `go:wasmexport` function returns, until calling back into
the Go-based Wasm module.

While some type restrictions have been relaxed in Go 1.24, there are still
limitations on the types that can be used with `go:wasmimport` and
`go:wasmexport` functions.
Due to the unfortunate mismatch between the 64-bit architecture of the client
and the 32-bit architecture of the host, it is not possible to pass pointers in
memory.
For example, a `go:wasmimport` function cannot take a pointer to a struct that
contains a pointer-typed field.

## Conclusion

The addition of the ability to build a WASI reactor and export Go functions to
Wasm in Go 1.24 represent a significant step forward for Go's WebAssembly
capabilities.
These features empower developers to create more versatile and powerful Go-based
Wasm applications, opening up new possibilities for Go in the Wasm ecosystem.
