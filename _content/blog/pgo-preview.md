---
title: Profile-guided optimization preview
date: 2023-02-08
by:
- Michael Pratt
summary:  Introduction to profile-guided optimization, available as a preview in Go 1.20.
---

When you build a Go binary, the Go compiler performs optimizations to try to generate the best performing binary it can.
For example, constant propagation can evaluate constant expressions at compile time, avoiding runtime evaluation cost.
Escape analysis avoids heap allocations for locally-scoped objects, avoiding GC overheads.
Inlining copies the body of simple functions into callers, often enabling further optimization in the caller (such as additional constant propagation or better escape analysis).

Go improves optimizations from release to release, but this is not always an easy task.
Some optimizations are tunable, but the compiler can't just "turn it up to 11" on every function because overly aggressive optimizations can actually hurt performance or cause excessive build times.
Other optimizations require the compiler to make a judgment call about what the "common" and "uncommon" paths in a function are.
The compiler must make a best guess based on static heuristics because it can't know which cases will be common at run time.

Or can it?

With no definitive information about how the code is used in a production environment, the compiler can operate only on the source code of packages.
But we do have a tool to evaluate production behavior: [profiling](https://go.dev/doc/diagnostics#profiling).
If we provide a profile to the compiler, it can make more informed decisions: more aggressively optimizing the most frequently used functions, or more accurately selecting common cases.

Using profiles of application behavior for compiler optimization is known as _Profile-Guided Optimization (PGO)_ (also known as Feedback-Directed Optimization (FDO)).

Go 1.20 includes initial support for PGO as a preview.
See the [profile-guided optimization user guide](https://go.dev/doc/pgo) for complete documentation.
There are still some rough edges that may prevent production use, but we would love for you to try it out and [send us any feedback or issues you encounter](https://go.dev/issue/new).

## Example

Let's build a service that converts Markdown to HTML: users upload Markdown source to `/render`, which returns the HTML conversion.
We can use [`gitlab.com/golang-commonmark/markdown`](https://pkg.go.dev/gitlab.com/golang-commonmark/markdown) to implement this easily.

### Set up

```
$ go mod init example.com/markdown
$ go get gitlab.com/golang-commonmark/markdown@bf3e522c626a
```

In `main.go`:

```
package main

import (
	"bytes"
	"io"
	"log"
	"net/http"
	_ "net/http/pprof"

	"gitlab.com/golang-commonmark/markdown"
)

func render(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	src, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("error reading body: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	md := markdown.New(
		markdown.XHTMLOutput(true),
		markdown.Typographer(true),
		markdown.Linkify(true),
		markdown.Tables(true),
	)

	var buf bytes.Buffer
	if err := md.Render(&buf, src); err != nil {
		log.Printf("error converting markdown: %v", err)
		http.Error(w, "Malformed markdown", http.StatusBadRequest)
		return
	}

	if _, err := io.Copy(w, &buf); err != nil {
		log.Printf("error writing response: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func main() {
	http.HandleFunc("/render", render)
	log.Printf("Serving on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

```

Build and run the server:

```
$ go build -o markdown.nopgo.exe
$ ./markdown.nopgo.exe
2023/01/19 14:26:24 Serving on port 8080...
```

Let's try sending some Markdown from another terminal.
We can use the README from the Go project as a sample document:

```
$ curl -o README.md -L "https://raw.githubusercontent.com/golang/go/c16c2c49e2fa98ae551fc6335215fadd62d33542/README.md"
$ curl --data-binary @README.md http://localhost:8080/render
<h1>The Go Programming Language</h1>
<p>Go is an open source programming language that makes it easy to build simple,
reliable, and efficient software.</p>
...
```

### Profiling

Now that we have a working service, let's collect a profile and rebuild with PGO to see if we get better performance.

In `main.go`, we imported [net/http/pprof](https://pkg.go.dev/net/http/pprof) which automatically adds a `/debug/pprof/profile` endpoint to the server for fetching a CPU profile.

Normally you want to collect a profile from your production environment so that the compiler gets a representative view of behavior in production.
Since this example doesn't have a "production" environment, we will create a simple program to generate load while we collect a profile.
Copy the source of [this program](https://go.dev/play/p/yYH0kfsZcpL) to `load/main.go` and start the load generator (make sure the server is still running!).

```
$ go run example.com/markdown/load
```

While that is running, download a profile from the server:

```
$ curl -o cpu.pprof "http://localhost:8080/debug/pprof/profile?seconds=30"
```

Once this completes, kill the load generator and the server.

### Using the profile

We can ask the Go toolchain to build with PGO using the `-pgo` flag to `go build`.
`-pgo` takes either the path to the profile to use, or `auto`, which will use the `default.pgo` file in the main package directory.

We recommend committing `default.pgo` profiles to your repository.
Storing profiles alongside your source code ensures that users automatically have access to the profile simply by fetching the repository (either via the version control system, or via `go get`) and that builds remain reproducible.
In Go 1.20, `-pgo=off` is the default, so users still need to add `-pgo=auto`, but a future version of Go is expected to change the default to `-pgo=auto`, automatically giving anyone that builds the binary the benefit of PGO.

Let's build:

```
$ mv cpu.pprof default.pgo
$ go build -pgo=auto -o markdown.withpgo.exe
```

### Evaluation

We will use a Go benchmark version of the load generator to evaluate the effect of PGO on performance.
Copy [this benchmark](https://go.dev/play/p/6FnQmHfRjbh) to `load/bench_test.go`.

First, we will benchmark the server without PGO. Start that server:

```
$ ./markdown.nopgo.exe
```

While that is running, run several benchmark iterations:

```
$ go test example.com/markdown/load -bench=. -count=20 -source ../README.md > nopgo.txt
```

Once that completes, kill the original server and start the version with PGO:

```
$ ./markdown.withpgo.exe
```

While that is running, run several benchmark iterations:

```
$ go test example.com/markdown/load -bench=. -count=20 -source ../README.md > withpgo.txt
```

Once that completes, let's compare the results:

```
$ go install golang.org/x/perf/cmd/benchstat@latest
$ benchstat nopgo.txt withpgo.txt
goos: linux
goarch: amd64
pkg: example.com/markdown/load
cpu: Intel(R) Xeon(R) W-2135 CPU @ 3.70GHz
        │  nopgo.txt  │            withpgo.txt             │
        │   sec/op    │   sec/op     vs base               │
Load-12   393.8µ ± 1%   383.6µ ± 1%  -2.59% (p=0.000 n=20)
```

The new version is around 2.6% faster!
In Go 1.20, workloads typically get between 2% and 4% CPU usage improvements from enabling PGO.
Profiles contain a wealth of information about application behavior and Go 1.20 just begins to crack the surface by using this information for inlining.
Future releases will continue improving performance as more parts of the compiler take advantage of PGO.

## Next steps

In this example, after collecting a profile, we rebuilt our server using the exact same source code used in the original build.
In a real-world scenario, there is always ongoing development.
So we may collect a profile from production, which is running last week's code, and use it to build with today's source code.
That is perfectly fine!
PGO in Go can handle minor changes to source code without issue.

For much more information on using PGO, best practices and caveats to be aware of, please see the [profile-guided optimization user guide](https://go.dev/doc/pgo).

Please send us your feedback!
PGO is still in preview and we'd love to hear about anything that is difficult to use, doesn't work correctly, etc.
Please file issues at https://go.dev/issue/new.
