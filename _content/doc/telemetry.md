---
title: "Go Telemetry"
layout: article
breadcrumb: true
date: 2024-02-07:00:00Z
---

<style>
.DocInfo {
  background-color: var(--color-background-info);
  padding: 1.5rem 2rem 1.5rem 4rem;
  border-left: 0.875rem solid var(--color-border);
  position: relative;
}
.DocInfo:before {
  content: "ⓘ";
  position: absolute;
  top: 1rem;
  left: 1rem;
  font-size: 2rem;
}
</style>

Table of Contents:

 [Background](#background)\
 [Overview](#overview)\
 [Configuration](#config)\
 [Counters](#counters)\
 [Reporting and Uploading](#reports)\
 [Charts](#charts) \
 [Telemetry Proposals](#proposals)\
 [IDE Prompting](#ide) \
 [Frequently Asked Questions](#faq)

## Background {#background}

Go telemetry is a way for Go toolchain programs to collect data about their
performance and usage. Here "Go toolchain" means developer tools maintained
by the Go team, including the `go` command and supplemental tools such as the
Go language server [`gopls`] or Go security tool [`govulncheck`]. Go telemetry is
only intended for use in programs maintained by the Go team and their selected
dependencies like [Delve].

By default, telemetry data is kept only on the local computer, but users may
opt in to uploading an approved subset of telemetry data to [telemetry.go.dev].
Uploaded data helps the Go team improve the Go language and its tools,
by helping us understand usage and breakages.

The word "telemetry" has acquired negative connotations in the world of open
source software, in many cases deservedly so. Yet measuring the user experience
is an important element of modern software engineering, and data sources such
as GitHub issues or annual surveys are coarse and lagging indicators,
insufficient for the types of questions the Go team needs to be able to answer.
Go telemetry is designed to help programs in the toolchain collect useful data
about their reliability, performance, and usage, while maintaining the
transparency and privacy that users expect from the Go project. To learn more
about the design process and motivation for telemetry, please see the
[telemetry blog posts](https://research.swtch.com/telemetry).
To learn more about telemetry and privacy, please see the
[telemetry privacy policy](https://telemetry.go.dev/privacy).

This page explains how Go telemetry works, in some detail. For quick answers to
frequently asked questions, see the [FAQ](#faq).

<div class="DocInfo">
Using Go 1.23 or later, to <strong>opt in</strong> to uploading telemetry data
to the Go team, run:
<pre>
go telemetry on
</pre>
To completely disable telemetry, including local collection, run:
<pre>
go telemetry off
</pre>
To revert to the default mode of local-only telemetry, run:
<pre>
go telemetry local
</pre>
Prior to Go 1.23, this can also be done with the
<code>golang.org/x/telemetry/cmd/gotelemetry</code> command. See <a
href="#config">Configuration</a> for more details.
</div>

## Overview {#overview}

Go telemetry uses three core data types:

- [_Counters_](#counters) are lightweight counts of named events, instrumented
  in the toolchain program. If collection is enabled (the [mode](#config)
  is **local** or **on**), counters are written to a memory-mapped file in the
  local file system.
- [_Reports_](#reports) are aggregated summaries of counters for a given week.
  If uploading is enabled (the [mode](#config) is **on**), reports for
  [approved counters](#proposals) are uploaded to [telemetry.go.dev], where
  they are publicly accessible.
- [_Charts_](#charts) summarize uploaded reports for all users.
  Charts can be viewed at [telemetry.go.dev].

All local Go telemetry data and configuration is stored in the directory
<code>[os.UserConfigDir()](/pkg/os#UserConfigDir)/go/telemetry</code>
directory. Below, we'll refer to this directory as `<gotelemetry>`.

The diagram below illustrates this data flow.

<div class="image">
  <center>
    <img max-width="800px" src="/doc/telemetry/dataflow.png" />
  </center>
</div>

In the rest of this document, we'll explore the components of this diagram. But
first, let's learn more about the configuration that controls it.

## Configuration {#config}

The behavior of Go telemetry is controlled by a single value: the telemetry
_mode_. The possible values for `mode` are `local` (the default), `on`, or
`off`:

- When `mode` is `local`, telemetry data is collected and stored on the local
  computer, but never uploaded to remote servers.
- When `mode` is `on`, data is collected, and may be uploaded depending on
  [sampling](#uploads).
- When `mode` is `off`, data is neither collected nor uploaded.

With Go 1.23 or later, the following commands interact with the telemetry mode:

- `go telemetry`: see the current mode.
- `go telemetry on`: set the mode to `on`.
- `go telemetry off`: set the mode to `off`.
- `go telemetry local`: set the mode to `local`.

Information about telemetry configuration is also available via read-only Go
environment variables:

- `go env GOTELEMETRY` reports the telemetry mode.
- `go env GOTELEMETRYDIR` reports the directory holding telemetry configuration
  and data.

The [`gotelemetry`](/pkg/golang.org/x/telemetry/cmd/gotelemetry) command can
also be used to configure the telemetry mode, as well as to inspect local
telemetry data. Use this command to install it:

```
go install golang.org/x/telemetry/cmd/gotelemetry@latest
```

For the complete usage information of the `gotelemetry` command line tool,
see its [package documentation](/pkg/golang.org/x/telemetry/cmd/gotelemetry).

## Counters {#counters}

As mentioned above, Go telemetry is instrumented via _counters_. Counters come
in two variants: basic counters and stack counters.

### Basic counters

A _basic counter_ is an incrementable value with a name that describes the
event that it counts. For example, the `gopls/client:vscode` counter records
the number of times a `gopls` session is initiated by VS Code. Alongside this
counter we may have `gopls/client:neovim`, `gopls/client:eglot`, and so on, to
record sessions with different editors or language clients. If you used
multiple editors throughout the week, you might record the following counter
data:

    gopls/client:vscode 8
    gopls/client:neovim 5
    gopls/client:eglot  2

When counters are related in this way, we sometimes refer to the part before
the `:` the _chart name_ (`gopls/client` in this case), and the part after `:`
as the _bucket name_ (`vscode`). We'll see why this matters when we discuss
[charts](#charts).

Basic counters can also represent a _histogram_. For example, the {{raw
`<code>gopls/completion/latency:&lt;50ms</code>`}} counter records the number
of times an autocompletion takes less than 50ms.

{{raw `
<pre>
gopls/completion/latency:&lt;10ms
gopls/completion/latency:&lt;50ms
gopls/completion/latency:&lt;100ms
...
</pre>
`}}

This pattern for recording histogram data is a convention: there's nothing
special about the {{raw `<code>&lt;50ms</code>`}} bucket name. These types of
counters are commonly used to measure performance.

### Stack counters

A _stack counter_ is a counter that also records the current call stack of the
Go toolchain program when the count is incremented. For example, the
`crash/crash` stack counter records the call stack when a toolchain program
crashes:

    crash/crash
    golang.org/x/tools/gopls/internal/golang.hoverBuiltin:+22
    golang.org/x/tools/gopls/internal/golang.Hover:+94
    golang.org/x/tools/gopls/internal/server.Hover:+42
    ...

Stack counters typically measure events where program invariants are violated.
The most common example of this is a crash, but another example is the
`gopls/bug` stack counter, which counts unusual situations identified in
advance by the programmer, such as a recovered panic or an error that "can't
happen". Stack counters include only the names and line numbers of functions
within Go toolchain programs. They don't include any information about user
inputs, such as the names or contents of a user's source code.

Stack counters can help track down rare or tricky bugs that don't get reported
by other means. Since introducing the `gopls/bug` counter, we've found
[dozens of instances](https://github.com/golang/go/issues?q=label%3Agopls%2Ftelemetry-wins)
of "unreachable" code that was reached in practice, and tracking down these
exceptions has led to the discovery (and fix) of many user-visible bugs that
were either not obvious to the user or too difficult to report. Especially with
prerelease testing, stack counters can help us improve the product more
efficiently than we could without automation.

### Counter files

All counter data is written to the `<gotelemetry>/local` directory, in
files named according to the following schema:

```
[program name]@[program version]-[go version]-[GOOS]-[GOARCH]-[date].v1.count
```

- The **program name** is the basename of the program's package path, as reported
  by [debug.BuildInfo].
- The **program version** and **go version** are also reported by [debug.BuildInfo].
- The **GOOS** and **GOARCH** values are reported by
  [`runtime.GOOS`](/pkg/runtime#GOOS) and
  [`runtime.GOARCH`](/pkg/runtime#GOARCH).
- The **date** is the date the counter file was created, in `YYYY-MM-DD` format.

These files are memory mapped into each running instance of the instrumented
programs. The use of a memory-mapped file means that even if the program
immediately crashes, or several copies of instrumented tools are running
simultaneously, the counters are recorded safely.

## Reporting and uploading {#reports}

Approximately once a week, counter data gets aggregated into reports named
`<date>.json` in the `<gotelemetry>/local` directory. These reports sum all of
counts for the previous week, grouped by the same program identifiers used for
the counter file (program name, program version, go version, GOOS, and GOARCH).

Local reports can be viewed as charts with the
[`gotelemetry view`](/pkg/golang.org/x/telemetry/cmd/gotelemetry) command.
Here's an example summary of the `gopls/completion/latency` counter:

<div class="image">
  <center>
    <img max-width="800px" src="/doc/telemetry/gopls-latency.png" />
  </center>
</div>

### Uploading {#uploads}

If telemetry uploading is enabled, the weekly reporting process will also
generate reports containing the subset of counters present in the
[upload config](https://telemetry.go.dev/config). These counters must be
approved by the public review process described in the next section. After it
has been successfully uploaded, a copy of the uploaded reports are stored in
the `<gotelemetry>/upload` directory.

Once enough users opt in to uploading telemetry data, the upload process will
randomly skip uploading for a fraction of reports, to reduce collection amounts
and increase privacy while maintaining statistical significance.

## Charts {#charts}

In addition to accepting uploads, the [telemetry.go.dev] website makes uploaded
data publicly available. Each day, uploaded reports are processed into two
outputs, which are available on the [telemetry.go.dev] homepage.

- _merged_ reports merged counters from all uploads received on the given day.
- _charts_ plot uploaded data as specified in the [chart config], which was
  produced as part of the proposal process. Recall from the discussion of
  [counters](#counters) that counter names such as `foo:bar` are decomposed
  into the chart name `foo` and bucket name `bar`. Each chart aggregates
  counters with the same chart name into the corresponding buckets.

Charts are specified in the format of the [chartconfig] package. For example,
here's the chart config for the `gopls/client` chart.

    title: Editor Distribution
    counter: gopls/client:{vscode,vscodium,vscode-insiders,code-server,eglot,govim,neovim,coc.nvim,sublimetext,other}
    description: measure editor distribution for gopls users.
    type: partition
    issue: https://go.dev/issue/61038
    issue: https://go.dev/issue/62214 # add vscode-insiders
    program: golang.org/x/tools/gopls
    version: v0.13.0 # temporarily back-version to demonstrate config generation.

This configuration describes the chart to be produced, enumerates the set of
counters to be aggregated, and specifies the program versions to which the
chart applies. Additionally, the [proposal process](#proposals) requires that
an accepted proposal be associated with the chart. Here's the chart resulting
from that config:

<div class="image">
  <center>
    <img src="/doc/telemetry/gopls-clients.png" />
  </center>
</div>

## The telemetry proposal process {#proposals}

Changes to the upload configuration or set of charts on [telemetry.go.dev] must
go through the _telemetry proposal process_, which is intended to ensure
transparency around changes to the telemetry configuration.

Notably, there is actually no distinction between upload configuration and
chart configuration in this process. Upload configuration is itself expressed
in terms of the aggregations that we want to render on telemetry.go.dev, based
on the principle that we should only collect data that we want to _see_.

The proposal process is as follows:

1. The proposer creates a CL modifying [config.txt] of the [chartconfig]
   package to contain the desired new counter aggregations.
2. The proposer files a [proposal] to merge this CL.
3. Once discussion on the issue resolves, the proposal is approved or declined
   by a member of the Go team.
4. An automatic process regenerates the upload config to allow uploading of the
   counters required for the new chart. This process will also regularly add
   new versions of the relevant programs to the upload config as they are
   released.

In order to be approved, new charts can't carry sensitive user information,
and additionally must be both useful and feasible. In order to be useful,
charts must serve a specific purpose, with actionable outcomes. In order to be
feasible, it must be possible to reliably collect the requisite data, and the
resulting measurements must be statistically significant. To demonstrate
feasibility, the proposer may be asked to instrument the target program with
counters and collect them locally first.

The full set of such proposals is available at the
[proposal project](https://github.com/orgs/golang/projects/29) on GitHub.

## IDE Prompting {#ide}

For telemetry to answer the types of questions we want to ask of it, the set of
users opting in to uploading need not be large--approximately 16,000
participants would allow for statistically significant measurements at the
desired level of granularity. However, there is still a cost to assembling this
healthy sample: we need to ask a large number of Go developers if they want to
opt in.

Furthermore, even if a large number of users choose to opt in _now_ (perhaps
after reading a Go blog post), those users may be skewed toward experienced Go
developers, and over time that initial sample will grow even more skewed.
Also, as people replace their computers, they must actively choose to opt in
again. In the telemetry blog post series, this is referred to as the
["campaign cost"](https://research.swtch.com/telemetry-opt-in#campaign) of
the opt-in model.

To help keep the sample of participating users fresh, the Go language server
[`gopls`] supports a prompt that asks users to opt in to Go telemetry.
Here's what that looks like from VS Code:

<div class="image">
  <center>
    <img width="600px" src="/doc/telemetry/prompt.png" />
  </center>
</div>

If users choose "Yes", their telemetry [mode](#config) will be set to `on`,
just as if they had run
[`gotelemetry on`](/pkg/golang.org/x/telemetry/cmd/gotelemetry). In this way,
opting in is as easy as possible, and we can continually reach a large and
stratified sample of Go developers.

## Frequently Asked Question {#faq}

**Q: How do I enable or disable Go telemetry?**

A: Use the `gotelemetry` command, which can be installed with `go install
golang.org/x/telemetry/cmd/gotelemetry@latest`. Run `gotelemetry off` to
disable everything, even local collection. Run `gotelemetry on` to enable
everything, including uploading approved counters to [telemetry.go.dev]. See
the [Configuration](#config) section for more info.

**Q: Where does local data get stored?**

A: In the <code>[os.UserConfigDir()](/pkg/os#UserConfigDir)/go/telemetry</code> directory.

**Q: How often does data get uploaded, if I opt in?**

A: Approximately once a week.

**Q: What data gets uploaded, if I opt in?**

A: Only counters that are listed in the
[upload config](https://telemetry.go.dev/config) may be uploaded. 
This is generated from the [chart config], which may be more readable.

**Q: How do counters get added to the upload config?**

A: Through the [public proposal process](#proposals).

**Q: Where can I see telemetry data that has been uploaded?**

A: Uploaded data is available as charts or merged summaries at [telemetry.go.dev].

**Q: Where is the source code for Go telemetry?**

A: At [golang.org/x/telemetry](/pkg/golang.org/x/telemetry).

[`gopls`]: /pkg/golang.org/x/tools/gopls
[`govulncheck`]: /pkg/golang.org/x/vuln/cmd/govulncheck
[Delve]: /pkg/github.com/go-delve/delve#section-readme
[debug.BuildInfo]: /pkg/runtime/debug#BuildInfo
[proposal]: /issue/new?assignees=&labels=Telemetry-Proposal&projects=golang%2F29&template=12-telemetry.yml&title=x%2Ftelemetry%2Fconfig%3A+proposal+title
[telemetry.go.dev]: https://telemetry.go.dev
[chartconfig]: /pkg/golang.org/x/telemetry/internal/chartconfig
[config.txt]: https://go.googlesource.com/telemetry/+/refs/heads/master/internal/chartconfig/config.txt
