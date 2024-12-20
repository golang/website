---
title: Testing concurrent code with testing/synctest
date: 2025-02-19
by:
- Damien Neil
tags:
- concurrency
- testing
summary: Go 1.24 contains an experimental package to aid in testing concurrent code.
---

One of Go's signature features is built-in support for concurrency.
Goroutines and channels are simple and effective primitives for
writing concurrent programs.

However, testing concurrent programs can be difficult and error prone.

In Go 1.24, we are introducing a new, experimental
[`testing/synctest`](/pkg/testing/synctest) package
to support testing concurrent code. This post will explain the motivation behind
this experiment, demonstrate how to use the synctest package, and discuss its potential future.

In Go 1.24, the `testing/synctest` package is experimental and
not subject to the Go compatibility promise.
It is not visible by default.
To use it, compile your code with `GOEXPERIMENT=synctest` set in your environment.

## Testing concurrent programs is difficult

To begin with, let us consider a simple example.

The [`context.AfterFunc`](/pkg/context#AfterFunc) function
arranges for a function to be called in its own goroutine after a context is canceled.
Here is a possible test for `AfterFunc`:

{{raw `
    func TestAfterFunc(t *testing.T) {
        ctx, cancel := context.WithCancel(context.Background())

        calledCh := make(chan struct{}) // closed when AfterFunc is called
        context.AfterFunc(ctx, func() {
            close(calledCh)
        })

        // TODO: Assert that the AfterFunc has not been called.

        cancel()

        // TODO: Assert that the AfterFunc has been called.
    }
`}}

We want to check two conditions in this test:
The function is not called before the context is canceled,
and the function *is* called after the context is canceled.

Checking a negative in a concurrent system is difficult.
We can easily test that the function has not been called *yet*,
but how do we check that it *will not* be called?

A common approach is to wait for some amount of time before
concluding that an event will not happen.
Let's try introducing a helper function to our test which does this.

{{raw `
    // funcCalled reports whether the function was called.
    funcCalled := func() bool {
        select {
        case <-calledCh:
            return true
        case <-time.After(10 * time.Millisecond):
            return false
        }
    }

    if funcCalled() {
        t.Fatalf("AfterFunc function called before context is canceled")
    }

    cancel()

    if !funcCalled() {
        t.Fatalf("AfterFunc function not called after context is canceled")
    }
`}}

This test is slow:
10 milliseconds isn't a lot of time, but it adds up over many tests.

This test is also flaky:
10 milliseconds is a long time on a fast computer,
but it isn't unusual to see pauses lasting several seconds
on shared and overloaded
[CI](https://en.wikipedia.org/wiki/Continuous_integration)
systems.

We can make the test less flaky at the expense of making it slower,
and we can make it less slow at the expense of making it flakier,
but we can't make it both fast and reliable.

## Introducing the testing/synctest package

The `testing/synctest` package solves this problem.
It allows us to rewrite this test to be simple, fast, and reliable,
without any changes to the code being tested.

The package contains only two functions: `Run` and `Wait`.

`Run` calls a function in a new goroutine.
This goroutine and any goroutines started by it
exist in an isolated environment which we call a *bubble*.
`Wait` waits for every goroutine in the current goroutine's bubble
to block on another goroutine in the bubble.

Let's rewrite our test above using the `testing/synctest` package.

{{raw `
    func TestAfterFunc(t *testing.T) {
        synctest.Run(func() {
            ctx, cancel := context.WithCancel(context.Background())

            funcCalled := false
            context.AfterFunc(ctx, func() {
                funcCalled = true
            })

            synctest.Wait()
            if funcCalled {
                t.Fatalf("AfterFunc function called before context is canceled")
            }

            cancel()

            synctest.Wait()
            if !funcCalled {
                t.Fatalf("AfterFunc function not called after context is canceled")
            }
        })
    }
`}}

This is almost identical to our original test,
but we have wrapped the test in a `synctest.Run` call
and we call `synctest.Wait` before asserting that the function has been called or not.

The `Wait` function waits for every goroutine in the caller's bubble to block.
When it returns, we know that the context package has either called the function,
or will not call it until we take some further action.

This test is now both fast and reliable.

The test is simpler, too:
we have replaced the `calledCh` channel with a boolean.
Previously we needed to use a channel to avoid a data race between
the test goroutine and the `AfterFunc` goroutine,
but the `Wait` function now provides that synchronization.

The race detector understands `Wait` calls,
and this test passes when run with `-race`.
If we remove the second `Wait` call,
the race detector will correctly report a data race in the test.

## Testing time

Concurrent code often deals with time.

Testing code that works with time can be difficult.
Using real time in tests causes slow and flaky tests,
as we have seen above.
Using fake time requires avoiding `time` package functions,
and designing the code under test to work with
an optional fake clock.

The `testing/synctest` package makes it simpler to test code that uses time.

Goroutines in the bubble started by `Run` use a fake clock.
Within the bubble, functions in the `time` package operate on the
fake clock. Time advances in the bubble when all goroutines are
blocked.

To demonstrate, let's write a test for the
[`context.WithTimeout`](/pkg/context#WithTimeout) function.
`WithTimeout` creates a child of a context,
which expires after a given timeout.

{{raw `
    func TestWithTimeout(t *testing.T) {
        synctest.Run(func() {
            const timeout = 5 * time.Second
            ctx, cancel := context.WithTimeout(context.Background(), timeout)
            defer cancel()

            // Wait just less than the timeout.
            time.Sleep(timeout - time.Nanosecond)
            synctest.Wait()
            if err := ctx.Err(); err != nil {
                t.Fatalf("before timeout, ctx.Err() = %v; want nil", err)
            }

            // Wait the rest of the way until the timeout.
            time.Sleep(time.Nanosecond)
            synctest.Wait()
            if err := ctx.Err(); err != context.DeadlineExceeded {
                t.Fatalf("after timeout, ctx.Err() = %v; want DeadlineExceeded", err)
            }
        })
    }
`}}

We write this test just as if we were working with real time.
The only difference is that we wrap the test function in `synctest.Run`,
and call `synctest.Wait` after each `time.Sleep` call to wait for the context
package's timers to finish running.

## Blocking and the bubble

A key concept in `testing/synctest` is the bubble becoming *durably blocked*.
This happens when every goroutine in the bubble is blocked,
and can only be unblocked by another goroutine in the bubble.

When a bubble is durably blocked:

  - If there is an outstanding `Wait` call, it returns.
  - Otherwise, time advances to the next time that could unblock a goroutine, if any.
  - Otherwise, the bubble is deadlocked and `Run` panics.

A bubble is not durably blocked if any goroutine is blocked
but might be woken by some event from outside the bubble.

The complete list of operations which durably block a goroutine is:

  - a send or receive on a nil channel
  - a send or receive blocked on a channel created within the same bubble
  - a select statement where every case is durably blocking
  - `time.Sleep`
  - `sync.Cond.Wait`
  - `sync.WaitGroup.Wait`

### Mutexes

Operations on a `sync.Mutex` are not durably blocking.

It is common for functions to acquire a global mutex.
For example, a number of functions in the reflect package
use a global cache guarded by a mutex.
If a goroutine in a synctest bubble blocks while acquiring
a mutex held by a goroutine outside the bubble,
it is not durably blocked—it is blocked, but will be unblocked
by a goroutine from outside its bubble.

Since mutexes are usually not held for long periods of time,
we simply exclude them from `testing/synctest`'s consideration.

### Channels

Channels created within a bubble behave differently from ones created outside.

Channel operations are durably blocking only if the channel is bubbled
(created in the bubble).
Operating on a bubbled channel from outside the bubble panics.

These rules ensure that a goroutine is durably blocked only when
communicating with goroutines within its bubble.

### I/O

External I/O operations, such as reading from a network connection,
are not durably blocking.

Network reads may be unblocked by writes from outside the bubble,
possibly even from other processes.
Even if the only writer to a network connection is also in the same bubble,
the runtime cannot distinguish between a connection waiting for more data to arrive
and one where the kernel has received data and is in the process of delivering it.

Testing a network server or client with synctest will generally
require supplying a fake network implementation.
For example, the [`net.Pipe`](/pkg/net#Pipe) function
creates a pair of `net.Conn`s that use an in-memory network connection
and can be used in synctest tests.

## Bubble lifetime

The `Run` function starts a goroutine in a new bubble.
It returns when every goroutine in the bubble has exited.
It panics if the bubble is durably blocked
and cannot be unblocked by advancing time.

The requirement that every goroutine in the bubble exit before Run returns
means that tests must be careful to clean up any background goroutines
before completing.

## Testing networked code

Let's look at another example, this time using the `testing/synctest`
package to test a networked program.
For this example, we'll test the `net/http` package's handling of
the 100 Continue response.

An HTTP client sending a request can include an "Expect: 100-continue"
header to tell the server that the client has additional data to send.
The server may then respond with a 100 Continue informational response
to request the rest of the request,
or with some other status to tell the client that the content is not needed.
For example, a client uploading a large file might use this feature to
confirm that the server is willing to accept the file before sending it.

Our test will confirm that when sending an "Expect: 100-continue" header
the HTTP client does not send a request's content before the server
requests it, and that it does send the content after receiving a
100 Continue response.

Often tests of a communicating client and server can use a
loopback network connection. When working with `testing/synctest`,
however, we will usually want to use a fake network connection
to allow us to detect when all goroutines are blocked on the network.
We'll start this test by creating an `http.Transport` (an HTTP client) that uses
an in-memory network connection created by [`net.Pipe`](/pkg/net#Pipe).

{{raw `
    func Test(t *testing.T) {
        synctest.Run(func() {
            srvConn, cliConn := net.Pipe()
            defer srvConn.Close()
            defer cliConn.Close()
            tr := &http.Transport{
                DialContext: func(ctx context.Context, network, address string) (net.Conn, error) {
                    return cliConn, nil
                },
                // Setting a non-zero timeout enables "Expect: 100-continue" handling.
                // Since the following test does not sleep,
                // we will never encounter this timeout,
                // even if the test takes a long time to run on a slow machine.
                ExpectContinueTimeout: 5 * time.Second,
            }
`}}

We send a request on this transport with the "Expect: 100-continue" header set.
The request is sent in a new goroutine, since it won't complete until the end of the test.

{{raw `
            body := "request body"
            go func() {
                req, _ := http.NewRequest("PUT", "http://test.tld/", strings.NewReader(body))
                req.Header.Set("Expect", "100-continue")
                resp, err := tr.RoundTrip(req)
                if err != nil {
                    t.Errorf("RoundTrip: unexpected error %v", err)
                } else {
                    resp.Body.Close()
                }
            }()
`}}

We read the request headers sent by the client.

{{raw `
            req, err := http.ReadRequest(bufio.NewReader(srvConn))
            if err != nil {
                t.Fatalf("ReadRequest: %v", err)
            }
`}}

Now we come to the heart of the test.
We want to assert that the client will not send the request body yet.

We start a new goroutine copying the body sent to the server into a `strings.Builder`,
wait for all goroutines in the bubble to block, and verify that we haven't read anything
from the body yet.

If we forget the `synctest.Wait` call, the race detector will correctly complain
about a data race, but with the `Wait` this is safe.

{{raw `
            var gotBody strings.Builder
            go io.Copy(&gotBody, req.Body)
            synctest.Wait()
            if got := gotBody.String(); got != "" {
                t.Fatalf("before sending 100 Continue, unexpectedly read body: %q", got)
            }
`}}

We write a "100 Continue" response to the client and verify that it now sends the
request body.

{{raw `
            srvConn.Write([]byte("HTTP/1.1 100 Continue\r\n\r\n"))
            synctest.Wait()
            if got := gotBody.String(); got != body {
                t.Fatalf("after sending 100 Continue, read body %q, want %q", got, body)
            }
`}}

And finally, we finish up by sending the "200 OK" response to conclude the request.

We have started several goroutines during this test.
The `synctest.Run` call will wait for all of them to exit before returning.

{{raw `
            srvConn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
        })
    }
`}}

This test can be easily extended to test other behaviors,
such as verifying that the request body is not sent if the server does not ask for it,
or that it is sent if the server does not respond within a timeout.

## Status of the experiment

We are introducing [`testing/synctest`](/pkg/testing/synctest)
in Go 1.24 as an *experimental* package.
Depending on feedback and experience
we may release it with or without amendments,
continue the experiment,
or remove it in a future version of Go.

The package is not visible by default.
To use it, compile your code with `GOEXPERIMENT=synctest` set in your environment.

We want to hear your feedback!
If you try out `testing/synctest`,
please report your experiences, positive or negative,
on [go.dev/issue/67434](/issue/67434).
