---
title: "Go Concurrency Patterns: Context"
date: 2014-07-29
by:
- Sameer Ajmani
tags:
- concurrency
- cancellation
- context
summary: An introduction to the Go context package.
---

## Introduction

In Go servers, each incoming request is handled in its own goroutine.
Request handlers often start additional goroutines to access backends such as
databases and RPC services.
The set of goroutines working on a request typically needs access to
request-specific values such as the identity of the end user, authorization
tokens, and the request's deadline.
When a request is canceled or times out, all the goroutines working on that
request should exit quickly so the system can reclaim any resources they are
using.

At Google, we developed a `context` package that makes it easy to pass
request-scoped values, cancellation signals, and deadlines across API boundaries
to all the goroutines involved in handling a request.
The package is publicly available as
[context](/pkg/context).
This article describes how to use the package and provides a complete working
example.

## Context

The core of the `context` package is the `Context` type:

{{code "context/interface.go" `/A Context/` `/^}/`}}

(This description is condensed; the
[godoc](/pkg/context) is authoritative.)

The `Done` method returns a channel that acts as a cancellation signal to
functions running on behalf of the `Context`: when the channel is closed, the
functions should abandon their work and return.
The `Err` method returns an error indicating why the `Context` was canceled.
The [Pipelines and Cancellation](/blog/pipelines) article discusses the `Done`
channel idiom in more detail.

A `Context` does _not_ have a `Cancel` method for the same reason the `Done`
channel is receive-only: the function receiving a cancellation signal is usually
not the one that sends the signal.
In particular, when a parent operation starts goroutines for sub-operations,
those sub-operations should not be able to cancel the parent.
Instead, the `WithCancel` function (described below) provides a way to cancel a
new `Context` value.

A `Context` is safe for simultaneous use by multiple goroutines.
Code can pass a single `Context` to any number of goroutines and cancel that
`Context` to signal all of them.

The `Deadline` method allows functions to determine whether they should start
work at all; if too little time is left, it may not be worthwhile.
Code may also use a deadline to set timeouts for I/O operations.

`Value` allows a `Context` to carry request-scoped data.
That data must be safe for simultaneous use by multiple goroutines.

### Derived contexts

The `context` package provides functions to _derive_ new `Context` values from
existing ones.
These values form a tree: when a `Context` is canceled, all `Contexts` derived
from it are also canceled.

`Background` is the root of any `Context` tree; it is never canceled:

{{code "context/interface.go" `/Background returns/` `/func Background/`}}

`WithCancel` and `WithTimeout` return derived `Context` values that can be
canceled sooner than the parent `Context`.
The `Context` associated with an incoming request is typically canceled when the
request handler returns.
`WithCancel` is also useful for canceling redundant requests when using multiple
replicas.
`WithTimeout` is useful for setting a deadline on requests to backend servers:

{{code "context/interface.go" `/WithCancel/` `/func WithTimeout/`}}

`WithValue` provides a way to associate request-scoped values with a `Context`:

{{code "context/interface.go" `/WithValue/` `/func WithValue/`}}

The best way to see how to use the `context` package is through a worked
example.

## Example: Google Web Search

Our example is an HTTP server that handles URLs like
`/search?q=golang&timeout=1s` by forwarding the query "golang" to the
[Google Web Search API](https://developers.google.com/web-search/docs/) and
rendering the results.
The `timeout` parameter tells the server to cancel the request after that
duration elapses.

The code is split across three packages:

  - [server](context/server/server.go) provides the `main` function and the handler for `/search`.
  - [userip](context/userip/userip.go) provides functions for extracting a user IP address from a request and associating it with a `Context`.
  - [google](context/google/google.go) provides the `Search` function for sending a query to Google.

### The server program

The [server](context/server/server.go) program handles requests like
`/search?q=golang` by serving the first few Google search results for `golang`.
It registers `handleSearch` to handle the `/search` endpoint.
The handler creates an initial `Context` called `ctx` and arranges for it to be
canceled when the handler returns.
If the request includes the `timeout` URL parameter, the `Context` is canceled
automatically when the timeout elapses:

{{code "context/server/server.go" `/func handleSearch/` `/defer cancel/`}}

The handler extracts the query from the request and extracts the client's IP
address by calling on the `userip` package.
The client's IP address is needed for backend requests, so `handleSearch`
attaches it to `ctx`:

{{code "context/server/server.go" `/Check the search query/` `/userip.NewContext/`}}

The handler calls `google.Search` with `ctx` and the `query`:

{{code "context/server/server.go" `/Run the Google search/` `/elapsed/`}}

If the search succeeds, the handler renders the results:

{{code "context/server/server.go" `/resultsTemplate/` `/(?m)}$/`}}

### Package userip

The [userip](context/userip/userip.go) package provides functions for
extracting a user IP address from a request and associating it with a `Context`.
A `Context` provides a key-value mapping, where the keys and values are both of
type `interface{}`.
Key types must support equality, and values must be safe for simultaneous use by
multiple goroutines.
Packages like `userip` hide the details of this mapping and provide
strongly-typed access to a specific `Context` value.

To avoid key collisions, `userip` defines an unexported type `key` and uses
a value of this type as the context key:

{{code "context/userip/userip.go" `/The key type/` `/const userIPKey/`}}

`FromRequest` extracts a `userIP` value from an `http.Request`:

{{code "context/userip/userip.go" `/func FromRequest/` `/}/`}}

`NewContext` returns a new `Context` that carries a provided `userIP` value:

{{code "context/userip/userip.go" `/func NewContext/` `/}/`}}

`FromContext` extracts a `userIP` from a `Context`:

{{code "context/userip/userip.go" `/func FromContext/` `/}/`}}

### Package google

The [google.Search](context/google/google.go) function makes an HTTP request
to the [Google Web Search API](https://developers.google.com/web-search/docs/)
and parses the JSON-encoded result.
It accepts a `Context` parameter `ctx` and returns immediately if `ctx.Done` is
closed while the request is in flight.

The Google Web Search API request includes the search query and the user IP as
query parameters:

{{code "context/google/google.go" `/func Search/` `/q.Encode/`}}

`Search` uses a helper function, `httpDo`, to issue the HTTP request and cancel
it if `ctx.Done` is closed while the request or response is being processed.
`Search` passes a closure to `httpDo` handle the HTTP response:

{{code "context/google/google.go" `/var results/` `/return results/`}}

The `httpDo` function runs the HTTP request and processes its response in a new
goroutine.
It cancels the request if `ctx.Done` is closed before the goroutine exits:

{{code "context/google/google.go" `/func httpDo/` `/^}/`}}

## Adapting code for Contexts

Many server frameworks provide packages and types for carrying request-scoped
values.
We can define new implementations of the `Context` interface to bridge between
code using existing frameworks and code that expects a `Context` parameter.

For example, Gorilla's
[github.com/gorilla/context](http://www.gorillatoolkit.org/pkg/context)
package allows handlers to associate data with incoming requests by providing a
mapping from HTTP requests to key-value pairs.
In [gorilla.go](context/gorilla/gorilla.go), we provide a `Context`
implementation whose `Value` method returns the values associated with a
specific HTTP request in the Gorilla package.

Other packages have provided cancellation support similar to `Context`.
For example, [Tomb](https://godoc.org/gopkg.in/tomb.v2) provides a `Kill`
method that signals cancellation by closing a `Dying` channel.
`Tomb` also provides methods to wait for those goroutines to exit, similar to
`sync.WaitGroup`.
In [tomb.go](context/tomb/tomb.go), we provide a `Context` implementation that
is canceled when either its parent `Context` is canceled or a provided `Tomb` is
killed.

## Conclusion

At Google, we require that Go programmers pass a `Context` parameter as the
first argument to every function on the call path between incoming and outgoing
requests.
This allows Go code developed by many different teams to interoperate well.
It provides simple control over timeouts and cancellation and ensures that
critical values like security credentials transit Go programs properly.

Server frameworks that want to build on `Context` should provide implementations
of `Context` to bridge between their packages and those that expect a `Context`
parameter.
Their client libraries would then accept a `Context` from the calling code.
By establishing a common interface for request-scoped data and cancellation,
`Context` makes it easier for package developers to share code for creating
scalable services.
