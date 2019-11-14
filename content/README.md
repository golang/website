# content

The `static` directory contains static content. This content is used alongside
the content in the main Go repository by the [`golangorg`](golang.org/x/website/cmd/golangorg)
binary when serving the [golang.org](https://golang.org) website.
The details of the directory to path mapping are documented at the top of
[`cmd/golangorg/main.go`](https://go.googlesource.com/website/+/refs/heads/master/cmd/golangorg/main.go).

TODO(dmitshur): The process below can be simplified.
See [golang.org/issue/29206#issuecomment-536099768](https://golang.org/issue/29206#issuecomment-536099768).

## Development mode

In production, CSS/JS/template assets need to be compiled into the `golangorg`
binary. It can be tedious to recompile assets every time, but you can pass a
flag to load CSS/JS/templates from disk every time a page loads:

```
golangorg -templates=$GOPATH/src/golang.org/x/website/content/static -http=:6060
```

## Recompiling static assets

Files such as `static/style.css`, `static/doc/copyright.html` and so on are not
present in the final binary. They are embedded into `static/static.go` by running
`go generate`. To compile a change and test it in your browser:

1) Make changes to an existing file such as `static/style.css`.

2) If a new file is being added to the `static` directory, add it to the `files`
slice in `static/gen`.

3) Run `go generate golang.org/x/website/content/static` so `static/static.go` is
up to date.

4) Run `go run golang.org/x/website/cmd/golangorg -http=:6060` and view your changes
in the browser at http://localhost:6060. You may need to disable your browser's cache
to avoid reloading a stale file.

A test exists to catch a possible mistake of forgetting to regenerate static assets:

```
website $ go test ./...
--- FAIL: TestStaticIsUpToDate (0.06s)
    gen_test.go:27: static.go is stale.  Run:
          $ go generate golang.org/x/website/content/static
          $ git diff
        to see the differences.
FAIL
FAIL	golang.org/x/website/content/static	0.650s
```
