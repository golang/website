---
title: "Go Protobuf: The new Opaque API"
date: 2024-12-16
by:
- Michael Stapelberg
tags:
- protobuf
summary: We are adding a new generated code API to Go Protobuf.
---

[[Protocol Buffers (Protobuf)](https://en.wikipedia.org/wiki/Protocol_Buffers)
is Google's language-neutral data interchange format. See
[protobuf.dev](https://protobuf.dev/).]

Back in March 2020, we released the `google.golang.org/protobuf` module, [a
major overhaul of the Go Protobuf API](/blog/protobuf-apiv2). This
package introduced first-class [support for
reflection](https://pkg.go.dev/google.golang.org/protobuf/reflect/protoreflect),
a [`dynamicpb`](https://pkg.go.dev/google.golang.org/protobuf/types/dynamicpb)
implementation and the
[`protocmp`](https://pkg.go.dev/google.golang.org/protobuf/testing/protocmp)
package for easier testing.

That release introduced a new protobuf module with a new API. Today, we are
releasing an additional API for generated code, meaning the Go code in the
`.pb.go` files created by the protocol compiler (`protoc`). This blog post
explains our motivation for creating a new API and shows you how to use it in
your projects.

To be clear: We are not removing anything. We will continue to support the
existing API for generated code, just like we still support the older protobuf
module (by wrapping the `google.golang.org/protobuf` implementation). Go is
[committed to backwards compatibility](/blog/compat) and this
applies to Go Protobuf, too!

## Background: the (existing) Open Struct API {#background}

We now call the existing API the Open Struct API, because generated struct types
are open to direct access. In the next section, we will see how it differs from
the new Opaque API.

To work with protocol buffers, you first create a `.proto` definition file like
this one:

    edition = "2023";  // successor to proto2 and proto3

    package log;

    message LogEntry {
      string backend_server = 1;
      uint32 request_size = 2;
      string ip_address = 3;
    }

Then, you [run the protocol compiler
(`protoc`)](https://protobuf.dev/getting-started/gotutorial/) to generate code
like the following (in a `.pb.go` file):

    package logpb

    type LogEntry struct {
      BackendServer *string
      RequestSize   *uint32
      IPAddress     *string
      // …internal fields elided…
    }

    func (l *LogEntry) GetBackendServer() string { … }
    func (l *LogEntry) GetRequestSize() uint32   { … }
    func (l *LogEntry) GetIPAddress() string     { … }

Now you can import the generated `logpb` package from your Go code and call
functions like
[`proto.Marshal`](https://pkg.go.dev/google.golang.org/protobuf/proto#Marshal)
to encode `logpb.LogEntry` messages into protobuf wire format.

You can find more details in the [Generated Code API
documentation](https://protobuf.dev/reference/go/go-generated/).

### (Existing) Open Struct API: Field Presence {#presence}

An important aspect of this generated code is how *field presence* (whether a
field is set or not) is modeled. For instance, the above example models presence
using pointers, so you could set the `BackendServer` field to:

1. `proto.String("zrh01.prod")`: the field is set and contains "zrh01.prod"
1. `proto.String("")`: the field is set (non-`nil` pointer) but contains an
   empty value
1. `nil` pointer: the field is not set

If you are used to generated code not having pointers, you are probably using
`.proto` files that start with `syntax = "proto3"`. The field presence behavior
changed over the years:

* `syntax = "proto2"` uses *explicit presence* by default
* `syntax = "proto3"` used *implicit presence* by default (where cases 2 and 3
  cannot be distinguished and are both represented by an empty string), but was
  later extended to allow [opting into explicit presence with the `optional`
  keyword](https://protobuf.dev/programming-guides/proto3/#field-labels)
* `edition = "2023"`, the [successor to both proto2 and
  proto3](https://protobuf.dev/editions/overview/), uses [*explicit
  presence*](https://protobuf.dev/programming-guides/field_presence/) by default

## The new Opaque API {#opaqueapi}

We created the new *Opaque API* to uncouple the [Generated Code
API](https://protobuf.dev/reference/go/go-generated/) from the underlying
in-memory representation. The (existing) Open Struct API has no such separation:
it allows programs direct access to the protobuf message memory. For example,
one could use the `flag` package to parse command-line flag values into protobuf
message fields:

    var req logpb.LogEntry
    flag.StringVar(&req.BackendServer, "backend", os.Getenv("HOST"), "…")
    flag.Parse() // fills the BackendServer field from -backend flag

The problem with such a tight coupling is that we can never change how we lay
out protobuf messages in memory. Lifting this restriction enables many
implementation improvements, which we'll see below.

What changes with the new Opaque API? Here is how the generated code from the
above example would change:

    package logpb

    type LogEntry struct {
      xxx_hidden_BackendServer *string // no longer exported
      xxx_hidden_RequestSize   uint32  // no longer exported
      xxx_hidden_IPAddress     *string // no longer exported
      // …internal fields elided…
    }

    func (l *LogEntry) GetBackendServer() string { … }
    func (l *LogEntry) HasBackendServer() bool   { … }
    func (l *LogEntry) SetBackendServer(string)  { … }
    func (l *LogEntry) ClearBackendServer()      { … }
    // …

With the Opaque API, the struct fields are hidden and can no longer be
directly accessed. Instead, the new accessor methods allow for getting, setting,
or clearing a field.

### Opaque structs use less memory {#lessmemory}

One change we made to the memory layout is to model field presence for
elementary fields more efficiently:

* The (existing) Open Struct API uses pointers, which adds a 64-bit word to the
  space cost of the field.
* The Opaque API uses [bit
  fields](https://en.wikipedia.org/wiki/Bit_field), which require one bit per
  field (ignoring padding overhead).

Using fewer variables and pointers also lowers load on the allocator and on the
garbage collector.

The performance improvement depends heavily on the shapes of your protocol
messages: The change only affects elementary fields like integers, bools, enums,
and floats, but not strings, repeated fields, or submessages (because it is
[less
profitable](https://protobuf.dev/reference/go/opaque-faq/#memorylayout)
for those types).

Our benchmark results show that messages with few elementary fields exhibit
performance that is as good as before, whereas messages with more elementary
fields are decoded with significantly fewer allocations:

                 │ Open Struct API │             Opaque API             │
                 │    allocs/op    │  allocs/op   vs base               │
    Prod#1          360.3k ± 0%       360.3k ± 0%  +0.00% (p=0.002 n=6)
    Search#1       1413.7k ± 0%       762.3k ± 0%  -46.08% (p=0.002 n=6)
    Search#2        314.8k ± 0%       132.4k ± 0%  -57.95% (p=0.002 n=6)

Reducing allocations also makes decoding protobuf messages more efficient:

                 │ Open Struct API │             Opaque API            │
                 │   user-sec/op   │ user-sec/op  vs base              │
    Prod#1         55.55m ± 6%        55.28m ± 4%  ~ (p=0.180 n=6)
    Search#1       324.3m ± 22%       292.0m ± 6%  -9.97% (p=0.015 n=6)
    Search#2       67.53m ± 10%       45.04m ± 8%  -33.29% (p=0.002 n=6)

(All measurements done on an AMD Castle Peak Zen 2. Results on ARM and Intel
CPUs are similar.)

Note: proto3 with implicit presence similarly does not use pointers, so you will
not see a performance improvement if you are coming from proto3. If you were
using implicit presence for performance reasons, forgoing the convenience of
being able to distinguish empty fields from unset ones, then the Opaque API now
makes it possible to use explicit presence without a performance penalty.

### Motivation: Lazy Decoding {#lazydecoding}

Lazy decoding is a performance optimization where the contents of a submessage
are decoded when first accessed instead of during
[`proto.Unmarshal`](https://pkg.go.dev/google.golang.org/protobuf/proto#Unmarshal). Lazy
decoding can improve performance by avoiding unnecessarily decoding fields which
are never accessed.

Lazy decoding can't be supported safely by the (existing) Open Struct API. While
the Open Struct API provides getters, leaving the (un-decoded) struct fields
exposed would be extremely error-prone. To ensure that the decoding logic runs
immediately before the field is first accessed, we must make the field private
and mediate all accesses to it through getter and setter functions.

This approach made it possible to implement lazy decoding with the Opaque
API. Of course, not every workload will benefit from this optimization, but for
those that do benefit, the results can be spectacular: We have seen logs
analysis pipelines that discard messages based on a top-level message condition
(e.g. whether `backend_server` is one of the machines running a new Linux kernel
version) and can skip decoding deeply nested subtrees of messages.

As an example, here are the results of the micro-benchmark we included,
demonstrating how lazy decoding saves over 50% of the work and over 87% of
allocations!

                      │   nolazy    │                lazy                │
                      │   sec/op    │   sec/op     vs base               │
    Unmarshal/lazy-24   6.742µ ± 0%   2.816µ ± 0%  -58.23% (p=0.002 n=6)

                      │    nolazy    │                lazy                 │
                      │     B/op     │     B/op      vs base               │
    Unmarshal/lazy-24   3.666Ki ± 0%   1.814Ki ± 0%  -50.51% (p=0.002 n=6)

                      │   nolazy    │               lazy                │
                      │  allocs/op  │ allocs/op   vs base               │
    Unmarshal/lazy-24   64.000 ± 0%   8.000 ± 0%  -87.50% (p=0.002 n=6)


### Motivation: reduce pointer comparison mistakes {#pointercomparison}

Modeling field presence with pointers invites pointer-related bugs.

Consider an enum, declared within the `LogEntry` message:

    message LogEntry {
      enum DeviceType {
        DESKTOP = 0;
        MOBILE = 1;
        VR = 2;
      };
      DeviceType device_type = 1;
    }

A simple mistake is to compare the `device_type` enum field like so:

    if cv.DeviceType == logpb.LogEntry_DESKTOP.Enum() { // incorrect!

Did you spot the bug? The condition compares the memory address instead of the
value. Because the `Enum()` accessor allocates a new variable on each call, the
condition can never be true. The check should have read:

    if cv.GetDeviceType() == logpb.LogEntry_DESKTOP {

The new Opaque API prevents this mistake: Because fields are hidden, all access
must go through the getter.

### Motivation: reduce accidental sharing mistakes {#accidentalsharing}

Let's consider a slightly more involved pointer-related bug. Assume you are
trying to stabilize an RPC service that fails under high load. The following
part of the request middleware looks correct, but still the entire service goes
down whenever just one customer sends a high volume of requests:

	logEntry.IPAddress = req.IPAddress
	logEntry.BackendServer = proto.String(hostname)
	// The redactIP() function redacts IPAddress to 127.0.0.1,
	// unexpectedly not just in logEntry *but also* in req!
	go auditlog(redactIP(logEntry))
	if quotaExceeded(req) {
		// BUG: All requests end up here, regardless of their source.
		return fmt.Errorf("server overloaded")
	}

Did you spot the bug? The first line accidentally copied the pointer (thereby
sharing the pointed-to variable between the `logEntry` and `req` messages)
instead of its value. It should have read:

	logEntry.IPAddress = proto.String(req.GetIPAddress())

The new Opaque API prevents this problem as the setter takes a value
(`string`) instead of a pointer:

	logEntry.SetIPAddress(req.GetIPAddress())


### Motivation: Fix Sharp Edges: reflection {#reflection}

To write code that works not only with a specific message type
(e.g. `logpb.LogEntry`), but with any message type, one needs some kind of
reflection. The previous example used a function to redact IP addresses. To work
with any type of message, it could have been defined as `func
redactIP(proto.Message) proto.Message { … }`.

Many years ago, your only option to implement a function like `redactIP` was to
reach for [Go's `reflect` package](/blog/laws-of-reflection),
which resulted in very tight coupling: you had only the generator output and had
to reverse-engineer what the input protobuf message definition might have looked
like. The [`google.golang.org/protobuf` module
release](/blog/protobuf-apiv2) (from March 2020) introduced
[Protobuf
reflection](https://pkg.go.dev/google.golang.org/protobuf/reflect/protoreflect),
which should always be preferred: Go's `reflect` package traverses the data
structure's representation, which should be an implementation detail. Protobuf
reflection traverses the logical tree of protocol messages without regard to its
representation.

Unfortunately, merely *providing* protobuf reflection is not sufficient and
still leaves some sharp edges exposed: In some cases, users might accidentally
use Go reflection instead of protobuf reflection.

For example, encoding a protobuf message with the `encoding/json` package (which
uses Go reflection) was technically possible, but the result is not [canonical
Protobuf JSON
encoding](https://protobuf.dev/programming-guides/proto3/#json). Use the
[`protojson`](https://pkg.go.dev/google.golang.org/protobuf/encoding/protojson)
package instead.

The new Opaque API prevents this problem because the message struct fields are
hidden: accidental usage of Go reflection will see an empty message. This is
clear enough to steer developers towards protobuf reflection.

### Motivation: Making the ideal memory layout possible {#idealmemory}

The benchmark results from the [More Efficient Memory
Representation](#lessmemory) section have already shown that protobuf
performance heavily depends on the specific usage: How are the messages defined?
Which fields are set?

To keep Go Protobuf as fast as possible for *everyone*, we cannot implement
optimizations that help only one program, but hurt the performance of other
programs.

The Go compiler used to be in a similar situation, up until [Go 1.20 introduced
Profile-Guided Optimization (PGO)](/blog/go1.20). By recording the
production behavior (through [profiling](/blog/pprof)) and feeding
that profile back to the compiler, we allow the compiler to make better
trade-offs *for a specific program or workload*.

We think using profiles to optimize for specific workloads is a promising
approach for further Go Protobuf optimizations. The Opaque API makes those
possible: Program code uses accessors and does not need to be updated when the
memory representation changes, so we could, for example, move rarely set fields
into an overflow struct.

## Migration {#migration}

You can migrate on your own schedule, or even not at all—the (existing) Open
Struct API will not be removed. But, if you’re not on the new Opaque API, you
won’t benefit from its improved performance, or future optimizations that target
it.

We recommend you select the Opaque API for new development. Protobuf Edition
2024 (see [Protobuf Editions Overview](https://protobuf.dev/editions/overview/)
if you are not yet familiar) will make the Opaque API the default.

### The Hybrid API {#hybridapi}

Aside from the Open Struct API and Opaque API, there is also the Hybrid API,
which keeps existing code working by keeping struct fields exported, but also
enabling migration to the Opaque API by adding the new accessor methods.

With the Hybrid API, the protobuf compiler will generate code on two API levels:
the `.pb.go` is on the Hybrid API, whereas the `_protoopaque.pb.go` version is
on the Opaque API and can be selected by building with the `protoopaque` build
tag.

### Rewriting Code to the Opaque API {#rewriting}

See the [migration
guide](https://protobuf.dev/reference/go/opaque-migration/)
for detailed instructions. The high-level steps are:

1. Enable the Hybrid API.
1. Update existing code using the `open2opaque` migration tool.
1. Switch to the Opaque API.

### Advice for published generated code: Use Hybrid API {#publishing}

Small usages of protobuf can live entirely within the same repository, but
usually, `.proto` files are shared between different projects that are owned by
different teams. An obvious example is when different companies are involved: To
call Google APIs (with protobuf), use the [Google Cloud Client Libraries for
Go](https://github.com/googleapis/google-cloud-go) from your project. Switching
the Cloud Client Libraries to the Opaque API is not an option, as that would be
a breaking API change, but switching to the Hybrid API is safe.

Our advice for such packages that publish generated code (`.pb.go` files) is to
switch to the Hybrid API please! Publish both the `.pb.go` and the
`_protoopaque.pb.go` files, please. The `protoopaque` version allows your
consumers to migrate on their own schedule.

### Enabling Lazy Decoding {#enablelazy}

Lazy decoding is available (but not enabled) once you migrate to the Opaque API!
🎉

To enable: in your `.proto` file, annotate your message-typed fields with the
`[lazy = true]` annotation.

To opt out of lazy decoding (despite `.proto` annotations), the [`protolazy`
package
documentation](https://pkg.go.dev/google.golang.org/protobuf/runtime/protolazy)
describes the available opt-outs, which affect either an individual Unmarshal
operation or the entire program.

## Next Steps {#nextsteps}

By using the open2opaque tool in an automated fashion over the last few years,
we have converted the vast majority of Google’s `.proto` files and Go code to
the Opaque API. We continuously improved the Opaque API implementation as we
moved more and more production workloads to it.

Therefore, we expect you should not encounter problems when trying the Opaque
API. In case you do encounter any issues after all, please [let us know on the
Go Protobuf issue tracker](https://github.com/golang/protobuf/issues/).

Reference documentation for Go Protobuf can be found on [protobuf.dev → Go
Reference](https://protobuf.dev/reference/go/).
