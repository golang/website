<!--{
  "Title": "Go Fuzzing technical details",
  "Breadcrumb": true
}-->

This document provides an overview of the technical details of the native fuzzing implementation, and is intended to be a resource for contributors.

## General architecture

The fuzzer uses a single coordinator process, which manages the corpus, and multiple worker processes, which mutate inputs and execute the fuzz target. The coordinator and worker processes communicate using a piped JSON-based RPC protocol, and shared regions of virtual memory.

For each worker the coordinator creates a goroutine which spawns the worker process and sets up the cross-process communication. Each goroutine then reads from a shared channel which is fed by the main coordinator loop, sending instructions it reads from the channel to the relevant worker processes.

The main coordinator loop picks inputs from the corpus, sending them to the shared worker channel. Whichever worker picks up that input from the channel will send a fuzzing request to the corresponding worker process. This process sits in a loop, mutating the input and executing the fuzz target until an execution either causes an increase in the coverage counters, causes a panic or crash, or passes a predetermined deadline.

If the worker process executes a mutated input which causes an increase in coverage counters or a recoverable panic, it signals this to the coordinator which is then able to reconstruct the mutated input. The coordinator will attempt to [minimize the input](#input-minimization), then either add it to the corpus for further fuzzing, in the case of finding increased coverage, or write it to the testdata directory, in the case of an input which causes an error or panic.

If a non-recoverable error occurs while fuzzing which causes the worker process to shut down (e.g. infinite loop, os.Exit, memory exhaustion, etc), minimization will not be attempted, and the failing input will be written to the testdata directory and reported.

<img alt="Sequence diagram of the interaction between coordinator and worker, as described above." src="/security/fuzz/seq-diagram.png"/>

### Cross-process communication

When spawning the child worker processes, the coordinator sets up two methods of communication: a pipe, which is used to pass JSON-based RPC messages, and a shared memory region, which is used to pass inputs and RNG state. Each worker process has its own pipe and shared memory region.

The RPC pipe is used by the coordinator to control the worker process, sending it either fuzzing or minimization instructions, and by the worker to relay results of its operations to the coordinator (i.e. whether the input expanded coverage, caused a crash, was successfully minimized, etc).

The shared memory region is used to pass specific information back and forth with the workers. The coordinator uses the region to pass the corpus entry to fuzz to the worker, and is used by the worker to store its current RNG state. The RNG state is used by the coordinator to reconstruct the mutations that were applied to the input by the worker when it has finished executing the target (this reconstruction happens both when the worker exits cleanly, and when it crashes.)

## Input selection

The coordinator currently does not implement any advanced form of input prioritization. It cycles through the entire corpus, looping after it exhausts the entries.

Similarly, the coordinator does not implement any type of corpus minimization (not to be confused with input minimization, [discussed below](#input-minimization)).

## Coverage guidance

The fuzzer uses [libFuzzer compatible](https://clang.llvm.org/docs/SanitizerCoverage.html#inline-8bit-counters) inline 8 bit coverage counters. These counters are inserted during compilation at each code edge, and are incremented on entry. Counters are not protected against overflow, so that they don't become saturated.

Similarly to AFL and libFuzzer, when tracking coverage, the counters are quantized to the nearest power of two. This allows the fuzzer to differentiate between insignificant and significant changes in execution flow. In order to track these changes, the fuzzer holds a slice of bytes which map to the inline counters, the bits of which indicate if there is at least one input in the corpus which increments the related counter at least 2^bit-position times. These bytes can become saturated, if there are inputs which cause counters to hit each quantized value, at which point the related counter fails to provide further useful coverage information.

As coverage counters are added to every edge during compilation, code not being fuzzed is also instrumented, which can cause the worker to detect coverage expansion that is unrelated to the target being executed (for instance if some new code path is triggered in a goroutine unrelated to the fuzz target). The worker attempts to reduce this in two ways: firstly it resets all counters immediately before executing the fuzz target and then snapshots the counters immediately after the target returns, and secondly by explicitly ignoring a set of packages which are likely to be "noisy"

A number of packages explicitly do not have counters inserted, since they are likely to introduce counter noise that is unrelated to the target being executed. These packages are:

* `context`
* `internal/fuzz`
* `reflect`
* `runtime`
* `sync`
* `sync/atomic`
* `syscall`
* `testing`
* `time`

## Mutation engine

When the worker receives a new input, it applies mutations to the input before executing the target with the input. After each mutation, the fuzz target is executed with the new input, and if coverage is not expanded, further mutations are applied. In order to prevent inputs from massively diverging from their initial state, after five mutations are applied to an input, it is reset to its original state before further mutations are applied. For example for the input `hello world`, the mutation strategy may look like the following:

```
0. hello world [initial state]
1. kello world [replace first byte]
2. world kello [swap two chunks]
3. world ke    [delete last three bytes]
4. owrld ke    [shuffle first three bytes]
5. owrldx ke   [insert random byte]
6. ello world  [reset to initial state, delete first byte]
...
```

The mutators attempt to bias towards producing smaller inputs, rather than larger inputs, in order to prevent rapid growth of the corpus size.

There are numerous mutators for `[]byte` and `string` types, and a smaller number of mutators for all the `int`, `uint`, and `float` types.

There are currently no execution driven mutation strategies implemented (such as input-to-comparison correspondence), nor dictionary based mutators.

## Input minimization

In order to prevent the corpus from ballooning (which bogs down the fuzzer both in terms of performance, and reducing the probability that a mutation will actually touch interesting data) we attempt to minimize each input discovered which expands coverage or causes a recoverable crash (non-recoverable crashes, such as those caused by memory exhaustion, are not minimized, as the process would be extremely slow). The employed strategy for minimization is rather simple, sequentially attempting to remove bytes from the input while maintaining the initial coverage found. In particular the minimization mechanism uses the following strategy:

1. Attempt to cut an exponentially smaller chunk of bytes off the end of the input
2. Attempt to remove each individual byte
3. Attempt to remove each possible subset of bytes
4. Attempt to replace each non-human readable byte with a human readable byte (i.e. something in the ASCII set of bytes)
