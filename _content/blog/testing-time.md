---
title: Testing Time (and other asynchronicities)
date: 2025-08-26
by:
- Damien Neil
tags:
- concurrency
- testing
- synctest
summary: A discussion of testing asynchronous code
  and an exploration of the `testing/synctest` package.
  Based on the GopherCon Europe 2025 talk with the same title.
---

In Go 1.24, we introduced the [`testing/synctest`](/pkg/testing/synctest)
package as an experimental package.
This package can significantly simplify writing tests for concurrent,
asynchronous code.
In Go 1.25, the `testing/synctest` package has graduated from experiment
to general availability.

What follows is the blog version of my talk on
the [`testing/synctest`](/pkg/testing/synctest) package
at GopherCon Europe 2025 in Berlin.

## What is an asynchronous function?

A synchronous function is pretty simple.
You call it, it does something, and it returns.

An asynchronous function is different.
You call it, it returns, and then it does something.

As a concrete, if somewhat artificial, example,
the following `Cleanup` function is synchronous.
You call it, it deletes a cache directory, and it returns.

```
func (c *Cache) Cleanup() {
    os.RemoveAll(c.cacheDir)
}
```

`CleanupInBackground` is an asynchronous function.
You call it, it returns, and the cache directory is deleted...sooner or later.

```
func (c *Cache) CleanupInBackground() {
    go os.RemoveAll(c.cacheDir)
}
```

Sometimes an asynchronous function does something in the future.
For example, the `context` package's `WithDeadline` function
returns a context which will be canceled in the future.

```
package context

// WithDeadline returns a derived context
// with a deadline no later than d.
func WithDeadline(parent Context, d time.Time) (Context, CancelFunc)
```

When I talk about testing concurrent code,
I mean testing these sorts of asynchronous operations,
both ones which use real time and ones which do not.

## Tests

A test verifies that a system behaves as we expect.
There's a lot of terminology describing types
of test--unit tests, integration tests, and so on--but
for our purposes here every kind of test reduces to three steps:

1. Set up some initial conditions.
2. Tell the system under test to do something.
3. Verify the result.

Testing a synchronous function is straightforward:

- You call the function;
- the function does something and returns;
- you verify the result.

Testing an asynchronous function, however, is tricky:

- You call the function;
- it returns;
- you wait for it to finish doing whatever it does;
- you verify the result.

If you don't wait for the correct amount of time,
you may find yourself verifying the result of an operation that hasn't happened yet
or has only happened partially.
This never ends well.

Testing an asynchronous function is especially tricky
when you want to assert that something has *not* happened.
You can verify that the thing has not happened yet,
but how do you know with certainty that it isn't going to happen later?

## An example

To make things a little more concrete,
let's work with a real-world example.
Consider the `context` package's `WithDeadline` function again.

```
package context

// WithDeadline returns a derived context
// with a deadline no later than d.
func WithDeadline(parent Context, d time.Time) (Context, CancelFunc)
```

There are two obvious tests to write for `WithDeadline`.

1. The context is *not* canceled *before* the deadline.
2. The context *is* canceled *after* the deadline.

Let's write a test.

To keep the amount of code marginally less overwhelming,
we'll just test the second case:
After the deadline expires, the context is canceled.

```
func TestWithDeadlineAfterDeadline(t *testing.T) {
    deadline := time.Now().Add(1 * time.Second)
    ctx, _ := context.WithDeadline(t.Context(), deadline)

    time.Sleep(time.Until(deadline))

    if err := ctx.Err(); err != context.DeadlineExceeded {
        t.Fatalf("context not canceled after deadline")
    }
}
```

This test is simple:

1. Use `context.WithDeadline` to create a context with a deadline one second in the future.
2. Wait until the deadline.
3. Verify that the context is canceled.

Unfortunately, this test obviously has a problem.
It sleeps until the exact moment the deadline expires.
Odds are good that the context has not been canceled yet by the time we examine it.
At best, this test will be very flaky.

Let's fix it.

```
time.Sleep(time.Until(deadline) + 100*time.Millisecond)
```

We can sleep until 100ms after the deadline.
A hundred milliseconds is an eternity in computer terms.
This should be fine.

Unfortunately, we still have two problems.

First, this test takes 1.1 seconds to execute.
That's slow.
This is a simple test.
It should execute in milliseconds at the most.

Second, this test is flaky.
A hundred milliseconds is an eternity in computer terms,
but on an overloaded continuous integration (CI) system
it isn't unusual to see pauses much longer than that.
This test will probably pass consistently on a developer's workstation,
but I would expect occasional failures in a CI system.

## Slow or flaky: Pick two

Tests that use real time are always slow or flaky.
Usually they're both.
If the test waits longer than necessary, it is slow.
If it doesn't wait long enough, it is flaky.
You can make the test more slow and less flaky,
or less slow and more flaky,
but you can't make it fast and reliable.

We have a lot of tests in the `net/http` package which use this approach.
They're all slow and/or flaky, which is what started me down the road
which brings us here today.

## Write synchronous functions?

The simplest way to test an asynchronous function is not to do it.
Synchronous functions are easy to test.
If you can transform an asynchronous function into a synchronous one,
it will be easier to test.

For example, if we consider our cache cleanup functions from earlier,
the synchronous `Cleanup` is obviously better than
the asynchronous `CleanupInBackground`.
The synchronous function is easier to test,
and the caller can easily start a new goroutine to run it in the background if needed.
As a general rule,
the higher up the call stack you can push your concurrency,
the better.

```
// CleanupInBackground is hard to test.
cache.CleanupInBackground()

// Cleanup is easy to test,
// and easy to run in the background when needed.
go cache.Cleanup()
```


Unfortunately, this sort of transformation isn't always possible.
For example, `context.WithDeadline` is an inherently asynchronous API.

## Instrument code for testability?

A better approach is to make our code more testable.

Here's an example of what this might look like for our `WithDeadline` test:

```
func TestWithDeadlineAfterDeadline(t *testing.T) {
    clock := fakeClock()
    timeout := 1 * time.Second
    deadline := clock.Now().Add(timeout)

    ctx, _ := context.WithDeadlineClock(
        t.Context(), deadline, clock)

    clock.Advance(timeout)
    context.WaitUntilIdle(ctx)
    if err := ctx.Err(); err != context.DeadlineExceeded {
        t.Fatalf("context not canceled after deadline")
    }
}
```

Instead of using real time, we use a fake time implementation.
Using fake time avoids unnecessarily slow tests,
because we never wait around doing nothing.
It also helps avoid test flakiness,
since the current time only changes when the test adjusts it.

There are various fake time packages out there,
or you can write your own.

To use fake time, we need to modify our API to accept a fake clock.
I've added a `context.WithDeadlineClock` function here,
that takes an additional clock parameter:

```
ctx, _ := context.WithDeadlineClock(
    t.Context(), deadline, clock)
```

When we advance our fake clock, we have a problem.
Advancing time is an asynchrounous operation.
Sleeping goroutines may wake up,
timers may send on their channels,
and timer functions may run.
We need to wait for that work to finish before we can test
the expected behavior of the system.

I've added a `context.WaitUntilIdle` function here,
which waits for any background work related to a context to complete:

```
clock.Advance(timeout)
context.WaitUntilIdle(ctx)
```

This is a simple example, but it demonstrates
the two fundamental principles of writing testable concurrent code:

1. Use fake time (if you use time).
2. Have some way to wait for quiescence,
   which is a fancy way of saying
   "all background activity has stopped and the system is stable".

The interesting question, of course, is how we do this.
I've glossed over the details in this example because
there are some big downsides to this approach.

It's hard.
Using a fake clock isn't difficult,
but identifying when background concurrent work is finished
and it is safe to examine the state of the system is.

Your code becomes less idiomatic.
You can't use standard time package functions.
You need to be very careful to keep track of everything happening
in the background.

You need to instrument not just your code,
but any other packages you use.
If you call any third-party concurrent code,
you're probably out of luck.

Worst of all, it can be just about impossible
to retrofit this approach into an existing codebase.

I attempted to apply this approach to Go's HTTP implementation,
and while I had some success at doing so in places,
the HTTP/2 server simply defeated me.
In particular, adding instrumentation to detect quiescence
without extensive rewriting proved infeasible,
or at least beyond my skills.

## Horrible runtime hacks?

What do we do if we can't make our code testable?

What if instead of instrumenting our code,
we had a way to observe the behavior of the uninstrumented system?

A Go program consists of a set of goroutines.
Those goroutines have states.
We just need to wait until all the goroutines have stopped running.

Unfortunately, the Go runtime doesn't provide any way to tell what
those goroutines are doing. Or does it?

The `runtime` package contains a function that gives us a stack trace
for every running goroutine, as well as their states.
This is text intended for human consumption,
but we could parse that output.
Could we use this to detect quiescence?

Now, of course this is a terrible idea.
There is no guarantee that the format of these stack traces will be stable over time.
You should not do this.

I did it.
And it worked.
In fact, it worked surprisingly well.

With a simple implementation of a fake clock,
a small amount of instrumentation to keep track of what goroutines were part of the test,
and some horrifying abuse of `runtime.Stack`,
I finally had a way to write fast, reliable tests for the `http` package.

The underlying implementation of these tests was horrible,
but it demonstrated that there was a useful concept here.

## A better way

Go may have built-in concurrency,
but testing programs that use that concurrency is hard.

We're faced with an unfortunate choice:
We can write simple, idiomatic code, but it will be impossible to test quickly and reliably;
or we can write testable code, but it will be complicated and unidiomatic.

So we asked ourselves what we can do to make this better.

As we saw earlier, the two fundamental features required to write testable concurrent code are
fake time and a way to wait for quiescence.

We need a better way to to wait for quiescence.
We should be able to ask the runtime when background goroutines have finished their work.
We also want to be able to limit the scope of this query to a single test,
so that unrelated tests do not interfere with each other.

We also need better support for testing programs using fake time.

It isn't hard to make a fake time implementation,
but code which uses an implementation like this is not idiomatic.

Idiomatic code will use a `time.Timer`,
but it is not possible to create a fake `Timer`.
We asked ourselves whether we should provide a way for tests to
create a fake `Timer`, where the test controls when the timer fires.

A testing implementation of time needs to define an entirely new version of the `time` package,
and pass that to every function that operates on time.
We considered whether we should define a common time interface,
in the same way that `net.Conn` is a common interface describing a network connection.

What we realized, however, is that unlike network connections,
there is only one possible implementation of fake time.
A fake network may want to introduce latency or errors.
Time, in contrast, does only one thing: It moves forward.
Tests need to control the rate at which time progresses,
but a timer scheduled to fire ten seconds in the future
should always fire ten (possibly fake) seconds in the future.

In addition, we don't want to upset the entire Go ecosystem.
Most programs today use functions in the time package.
We want to keep those programs not only working,
but idiomatic.

This led to the conclusion that what we need is a way for a test to
tell the time package to use a fake clock,
in much the same way that the Go playground uses a fake clock.
Unlike the playground,
we need to limit the scope of that change to a single test.
(It may not be obvious that the Go playground uses a fake clock,
because we turn any fake delays into real delays on the front end,
but it does.)

## The `synctest` experiment

And so in Go 1.24 we introduced [`testing/synctest`](/pkg/testing/synctest),
a new, experimental package to simplify testing concurrent programs.
Over the months following the release of Go 1.24
we gathered feedback from early adopters.
(Thank you to everyone who tried it out!)
We made a number of changes to address problems and shortcomings.
And now, in Go 1.25, we've released the `testing/synctest` package
as part of the standard library.

It lets you run a function in what we're calling a "bubble".
Within the bubble, the time package uses a fake clock,
and the `synctest` package provides a function to wait for the bubble to quiesce.

## The `synctest` package

The `synctest` package contains just two functions.

```
package synctest

// Test executes f in a new bubble.
// Goroutines in the bubble use a fake clock.
func Test(t *testing.T, f func(*testing.T))

// Wait waits for background activity in the bubble to complete.
func Wait()
```

[`Test`](/pkg/testing/synctest#Test) executes a function in a new bubble.

[`Wait`](/pkg/testing/synctest#Wait) blocks until every goroutine in the bubble is blocked
waiting for some other goroutine in the bubble.
We call that state being "durably blocked".

## Testing with synctest

Let's look at an example of synctest in action.

```
func TestWithDeadlineAfterDeadline(t *testing.T) {
    synctest.Test(t, func(t *testing.T) {
        deadline := time.Now().Add(1 * time.Second)
        ctx, _ := context.WithDeadline(t.Context(), deadline)

        time.Sleep(time.Until(deadline))
        synctest.Wait()
        if err := ctx.Err(); err != context.DeadlineExceeded {
            t.Fatalf("context not canceled after deadline")
        }
    })
}
```

This might look a little familiar.
This is the naïve test for `context.WithDeadline` that we looked at earlier.
The only changes are that we've wrapped the test in
a `synctest.Test` call to execute it in a bubble
and we have added a `synctest.Wait` call.

This test is fast and reliable.
It runs almost instantaneously.
It precisely tests the expected behavior of the system under test.
It also requires no modification of the `context` package.

Using the `synctest` package,
we can write simple, idiomatic code
and test it reliably.

This is a very simple example, of course,
but this is a real test of real production code.
If `synctest` had existed when the `context` package was written,
we would have had a much easier time writing tests for it.

## Time

Time in the bubble behaves much the same as the fake time in the Go playground.
Time starts at midnight, January 1, 2000 UTC.
If you need to run a test at some specific point in time for some reason,
you can just sleep until then.

```
func TestAtSpecificTime(t *testing.T) {
   synctest.Test(t, func(t *testing.T) {
       // 2000-01-01 00:00:00 +0000 UTC
       t.Log(time.Now().In(time.UTC))

       // This does not take 25 years.
       time.Sleep(time.Until(
           time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)))

       // 2025-01-01 00:00:00 +0000 UTC
       t.Log(time.Now().In(time.UTC))
   })
}
```

Time only passes when every goroutine in the bubble has blocked.
You can think of the bubble as simulating an infinitely fast computer:
Any amount of computation takes no time.

The following test will always print that zero seconds
of fake time have elapsed since the start of the test,
no matter how much real time has passed.

```
func TestExpensiveWork(t *testing.T) {
   synctest.Test(t, func(t *testing.T) {
       start := time.Now()
       for range 1e7 {
           // do expensive work
       }
       t.Log(time.Since(start)) // 0s
   })
}
```

In the next test, the `time.Sleep` call will return immediately,
rather than waiting for ten real seconds.
The test will always print that exactly ten fake seconds
have passed since the start of the test.

```
func TestSleep(t *testing.T) {
   synctest.Test(t, func(t *testing.T) {
       start := time.Now()
       time.Sleep(10 * time.Second)
       t.Log(time.Since(start)) // 10s
   })
}
```

## Waiting for quiescence

The [`synctest.Wait`](/pkg/testing/synctest#Wait) function
lets us wait for background activity to complete.

```
func TestWait(t *testing.T) {
   synctest.Test(t, func(t *testing.T) {
       done := false
       go func() {
           done = true
       }()

       // Wait for the above goroutine to finish.
       synctest.Wait()

       t.Log(done) // true
   })
}
```

If we didn't have the `Wait` call in the above test,
we would have a race condition:
One goroutine modifies the `done` variable
while another reads from it without synchronization.
The `Wait` call provides that synchronization.

You may be familiar with the `-race` test flag,
which enables the data race detector.
The race detector is aware of the synchronization provided by `Wait`,
and does not complain about this test.
If we forgot the `Wait` call, the race detector would correctly complain.

The `synctest.Wait` function provides synchronization,
but the passage of time does not.

In the next example, one goroutine writes to the `done` variable
while another sleeps for one nanosecond before reading from it.
It should be obvious that when run with a real clock outside a synctest bubble,
this code contains a race condition.
Inside a synctest bubble,
while the fake clock ensures that the goroutine completes before `time.Sleep` returns,
the race detector will still report the data race,
just like it would if this code were run outside a synctest bubble.

```
func TestTimeDataRace(t *testing.T) {
   synctest.Test(t, func(t *testing.T) {
       done := false
       go func() {
           done = true // write
       }()

       time.Sleep(1 * time.Nanosecond)

       t.Log(done)     // read (unsynchronized)
   })
}
```


Adding a `Wait` call provides explicit synchronization and fixes the data race:

```
time.Sleep(1 * time.Nanosecond)
synctest.Wait() // synchronize
t.Log(done)     // read
```

## Example: `io.Copy`

Taking advantage of the synchronization provided by `synctest.Wait` allows us
to write simpler tests with less explicit synchronization.

For example, consider this test of [`io.Copy`](/pkg/io#Copy).

```
func TestIOCopy(t *testing.T) {
   synctest.Test(t, func(t *testing.T) {
       srcReader, srcWriter := io.Pipe()
       defer srcWriter.Close()

       var dst bytes.Buffer
       go io.Copy(&dst, srcReader)

       data := "1234"
       srcWriter.Write([]byte("1234"))
       synctest.Wait()

       if got, want := dst.String(), data; got != want {
           t.Errorf("Copy wrote %q, want %q", got, want)
       }
   })
}
```

The `io.Copy` function copies data from an `io.Reader` to an `io.Writer`.
You might not immediately think of `io.Copy` as a concurrent function,
since it blocks until the copy has completed.
However, providing data to `io.Copy`'s reader is an asynchronous operation:

- `Copy` calls the reader's `Read` method;
- `Read` returns some data;
- and the data is written to the writer at a later time.

In this test, we are verifying that `io.Copy` writes new data to the writer
without waiting to fill its buffer.

Looking at the test step by step,
we first create an `io.Pipe` to serve as the source `io.Copy` reads from:

```
srcReader, srcWriter := io.Pipe()
defer srcWriter.Close()
```

We call `io.Copy` in a new goroutine,
copying from the read end of the pipe into a `bytes.Buffer`:

```
var dst bytes.Buffer
go io.Copy(&dst, srcReader)
```

We write to the other end of the pipe,
and wait for `io.Copy` to handle the data:

```
data := "1234"
srcWriter.Write([]byte("1234"))
synctest.Wait()
```

Finally, we verify that the destination buffer contains the desired data:

```
if got, want := dst.String(), data; got != want {
    t.Errorf("Copy wrote %q, want %q", got, want)
}
```

We don't need to add a mutex or other synchronization around the destination buffer,
because `synctest.Wait` ensures that it is never accessed concurrently.

This test demonstrates a few important points.

Even synchronous functions like `io.Copy`,
which do not perform additional background work after they return,
may exhibit asynchronous behaviors.

Using `synctest.Wait`, we can test those behaviors.

Note also that this test does not work with time.
Many asynchronous systems involve time, but not all.

## Bubble exit

The `synctest.Test` function waits for all goroutines in the bubble to exit
before returning.
Time stops advancing after the root goroutine (the goroutine started by `Test`) returns.

In the next example, `Test` waits for the background goroutine to run and exit
before it returns:

```
func TestWaitForGoroutine(t *testing.T) {
    synctest.Test(t, func(t *testing.T) {
        go func() {
            // This runs before synctest.Test returns.
        }()
    })
}
```

In this example, we schedule a `time.AfterFunc` for a time in the future.
The bubble's root goroutine returns before that time is reached,
so the `AfterFunc` never runs:

```
func TestDoNotWaitForTimer(t *testing.T) {
    synctest.Test(t, func(t *testing.T) {
        time.AfterFunc(1 * time.Nanosecond, func() {
            // This never runs.
        })
    })
}
```

In the next example, we start a goroutine that sleeps.
The root goroutine returns and time stops advancing.
The bubble is now deadlocked,
because `Test` is waiting for all goroutines in the bubble to finish
but the sleeping goroutine is waiting for time to advance.

```
func TestDeadlock(t *testing.T) {
    synctest.Test(t, func(t *testing.T) {
        go func() {
            // This sleep never returns and the test deadlocks.
            time.Sleep(1 * time.Nanosecond)
        }()
    })
}
```

## Deadlocks

The `synctest` package panics when a bubble is deadlocked
due to every goroutine in the bubble being durably blocked on
another goroutine in the bubble.

```
--- FAIL: Test (0.00s)
--- FAIL: TestDeadlock (0.00s)
panic: deadlock: main bubble goroutine has exited but blocked goroutines remain [recovered, repanicked]

goroutine 7 [running]:
(stacks elided for clarity)

goroutine 10 [sleep (durable), synctest bubble 1]:
time.Sleep(0x1)
	/Users/dneil/src/go/src/runtime/time.go:361 +0x130
_.TestDeadlock.func1.1()
	/tmp/s/main_test.go:13 +0x20
created by _.TestDeadlock.func1 in goroutine 9
	/tmp/s/main_test.go:11 +0x24
FAIL	_	0.173s
FAIL
```

The runtime will print stack traces for every goroutine in the deadlocked bubble.

When printing the status of a bubbled goroutine,
the runtime indicates when the goroutine is durably blocked.
You can see that the sleeping goroutine in this test is durably blocked.

## Durable blocking

"Durably blocking" is a core concept in synctest.

A goroutine is durably blocked when it is not only blocked,
but when it can only be unblocked by another goroutine in the same bubble.

When every goroutine in a bubble is durably blocked:

1. `synctest.Wait` returns.
2. If there is no `synctest.Wait` call in progress,
   fake time advances instantly to the next point that will wake a goroutine.
3. If there is no goroutine that can be woken by advancing time,
   the bubble is deadlocked and the test fails.

It is important for us to make a distinction between a goroutine which is merely blocked
and one which is *durably* blocked.
We don't want to declare a deadlock when a goroutine is temporarily blocked on
some event arising outside its bubble.

Let's look at some ways in which a goroutine can block non-durably.

### Not durably blocking: I/O (files, pipes, network connections, etc.)

The most important limitation is that I/O is not durably blocking,
including network I/O.
A goroutine reading from a network connection may be blocked,
but it will be unblocked by data arriving on that connection.

This is obviously true for a connection to some network service,
but it is also true for a loopback connection,
even when the reader and writer are both in the same bubble.

When we write data to a network socket,
even a loopback socket,
the data is passed to the kernel for delivery.
There is a period of time between the write system call returning
and the kernel notifying the other side of the connection that data is available.
The Go runtime cannot distinguish between a goroutine blocked waiting for
data that is already in the kernel's buffers
and one blocked waiting for data that will not arrive.

This means that tests of networked programs using synctest
usually cannot use real network connections.
Instead, they should use an in-memory fake.

I'm not going to go over the process of creating a fake network here,
but the `synctest` package documentation contains
[a complete worked example](/pkg/testing/synctest#hdr-Example__HTTP_100_Continue)
of testing an HTTP client and server communicating over a fake network.

### Not durably blocking: syscalls, cgo calls, anything that isn't Go

Syscalls and cgo calls are not durably blocking.
We can only reason about the state of goroutines executing Go code.

### Not durably blocking: Mutexes

Perhaps surprisingly, mutexes are not durably blocking.
This is a decision born of practicality:
Mutexes are often used to guard global state,
so a bubbled goroutine will often need to acquire a mutex held outside its bubble.
Mutexes are highly performance-sensitive,
so adding additional instrumentation to them
risks slowing down non-test programs.

We can test programs that use mutexes with synctest,
but the fake clock will not advance while a goroutine is blocked on mutex acquisition.
This hasn't posed a problem in any case we've encountered,
but it is something to be aware of.

### Durably blocking: `time.Sleep`

So what is durably blocking?

`time.Sleep` is obviously durable,
since time can only advance when every goroutine in the bubble is durably blocked.

### Durably blocking: send or receive on channels created in the same bubble

Channel operations on channels created within the same bubble are durable.

We make a distinction between bubbled channels (created in a bubble)
and unbubbled channels (created outside any bubble).
This means that a function using a global channel for synchronization,
for example to control access to a globally cached resource,
can be safely called from within a bubble.

Trying to operate on a bubbled channel from outside its bubble is an error.

### Durably blocking: `sync.WaitGroup` belonging to the same bubble

We also associate `sync.WaitGroup`s with bubbles.

`WaitGroup` doesn't have a constructor,
so we make the association with the bubble implicitly on the first call to `Go` or `Add`.

As with channels,
waiting on a `WaitGroup` belonging to the same bubble is durably blocking,
and waiting on one from outside the bubble is not.
Calling `Go` or `Add` on a `WaitGroup` belonging to a different bubble is an error.

### Durably blocking: `sync.Cond.Wait`

Waiting on a `sync.Cond` is always durably blocking.
Waking up a goroutine waiting on a `Cond` in a different bubble is an error.

### Durably blocking: `select{}`

Finally, an empty select is durably blocking.
(A select with cases is durably blocking if all the operations in it are so.)

That's the complete list of durably blocking operations.
It isn't very long,
but it's enough to handle almost all real-world programs.

The rule is that a goroutine is durably blocked when it is blocked,
and we can guarantee that it can only be unblocked
by another goroutine in its bubble.

In cases where it is possible to attempt to wake a bubbled goroutine from outside its bubble,
we panic.
For example, it is an error to operate on a bubbled channel from outside its bubble.

## Changes from 1.24 to 1.25

We released an experimental version of the `synctest` package in Go 1.24.
To ensure that early adopters were aware of the experimental status of the package,
you needed to set a GOEXPERIMENT flag to make the package visible.

The feedback we received from those early adopters was invaluable,
both in demonstrating that the package is useful
and in uncovering areas where the API needed work.

These are some of the changes made between the experimental version
and the version released in Go 1.25.

### Replaced Run with Test

The original version of the API created a bubble with a `Run` function:

```
// Run executes f in a new bubble.
func Run(f func())
```

It became clear that we needed a way to create a `*testing.T`
that is scoped to a bubble.
For example, `t.Cleanup` should run cleanup functions in the same bubble
they are registered in, not after the bubble exits.
We renamed `Run` to `Test` and made it create a `T` scoped to the lifetime
of the new bubble.

### Time stops when a bubble's root goroutine returns

We originally continued to advance time within a bubble for so long as
the bubble contained any goroutines waiting for future events.
This turned out to be very confusing when a long-lived goroutine never returned,
such as a goroutine reading forever from a `time.Ticker`.
We now stop advancing time when a bubble's root goroutine returns.
If the bubble is blocked waiting for time to advance,
this results in a deadlock and a panic which can be analyzed.

### Removed cases where "durable" wasn't

We cleaned up the definition of "durably blocking".
The original implementation had cases where a durably blocked goroutine could
be unblocked from outside the bubble.
For example, channels recorded whether they were created in a bubble,
but not which in which bubble they were created,
so one bubble could unblock a channel in a different bubble.
The current implementation contains no cases we know of
where a durably blocked goroutine can be unblocked from outside its bubble.

### Better stack traces

We made improvements to the information printed in stack traces.
When a bubble deadlocks, we by default now only print stacks for the goroutines in that bubble.
Stack traces also clearly indicate which goroutines in a bubble are durably blocked.

### Randomized events happening at the same time

We made improvements to the randomization of events happening at the same time.
Originally, timers scheduled to fire at the same instant
would always do so in the order they were created.
This ordering is now randomized.

## Future work

We're pretty happy with the synctest package at the moment.

Aside from the inevitable bug fixes,
we don't currently expect any major changes to it in the future.
Of course, with wider adoption it is always possible that we'll discover something
that needs doing.

One possible area of work is to improve the detection of durably blocked goroutines.
It would be nice if we could make mutex operations durably blocking,
with a restriction that a mutex acquired in a bubble must be released
from within the same bubble.

Testing networked code with synctest requires a fake network.
The `net.Pipe` function can create a fake `net.Conn`,
but there is currently no standard library function that creates
a fake `net.Listener` or `net.PacketConn`.
In addition, the `net.Conn` returned by `net.Pipe` is synchronous--every write blocks
until a read consumes the data--which is not representative of real network behavior.
Perhaps we should add a good fake implementations of common network interfaces
to the standard library.

## Conclusion

That's the `synctest` package.

I can't say that it makes testing concurrent code simple,
because concurrency is never simple.
What it does is let you write the simplest possible concurrent code,
using idiomatic Go,
and the standard time package,
and then write fast, reliable tests for it.

I hope you find it useful.
