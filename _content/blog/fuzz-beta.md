---
title: Fuzzing is Beta Ready
date: 2021-06-03
by:
- Katie Hockman
- Jay Conrod
tags:
- fuzz
- testing
summary: Native Go fuzzing is now ready for beta testing on tip.
---


We are excited to announce that native fuzzing is ready for beta testing on tip!

Fuzzing is a type of automated testing which continuously manipulates inputs to
a program to find issues such as panics or bugs. These semi-random data
mutations can discover new code coverage that existing unit tests may miss, and
uncover edge case bugs which would otherwise go unnoticed. Since fuzzing can
reach these edge cases, fuzz testing is particularly valuable for finding
security exploits and vulnerabilities.

See
[golang.org/s/draft-fuzzing-design](/s/draft-fuzzing-design)
for more details about this feature.


## Getting started

To get started, you may run the following

	$ go install golang.org/dl/gotip@latest
	$ gotip download

This builds the Go toolchain from the master branch. After running this, `gotip`
can act as a drop-in replacement for the `go` command. You can now run commands
like

	$ gotip test -fuzz=Fuzz

## Writing a fuzz test

A fuzz test must be in a \*\_test.go file as a function in the form `FuzzXxx`.
This function must be passed a` *testing.F` argument, much like a `*testing.T`
argument is passed to a `TestXxx` function.

Below is an example of a fuzz test that’s testing the behavior of the [net/url
package](https://pkg.go.dev/net/url#ParseQuery).

	//go:build go1.18
	// +build go1.18

	package fuzz

	import (
		"net/url"
		"reflect"
		"testing"
	)

	func FuzzParseQuery(f *testing.F) {
		f.Add("x=1&y=2")
		f.Fuzz(func(t *testing.T, queryStr string) {
			query, err := url.ParseQuery(queryStr)
			if err != nil {
				t.Skip()
			}
			queryStr2 := query.Encode()
			query2, err := url.ParseQuery(queryStr2)
			if err != nil {
				t.Fatalf("ParseQuery failed to decode a valid encoded query %s: %v", queryStr2, err)
			}
			if !reflect.DeepEqual(query, query2) {
				t.Errorf("ParseQuery gave different query after being encoded\nbefore: %v\nafter: %v", query, query2)
			}
		})
	}

You can read more about fuzzing at pkg.go.dev, including [an overview
of fuzzing with Go](https://pkg.go.dev/testing@master#hdr-Fuzzing) and the
[godoc for the new `testing.F` type](https://pkg.go.dev/testing@master#F).

## Expectations

This is a new feature that's still in beta, so you should expect some bugs
and an incomplete feature set. Check the [issue tracker for issues labelled
“fuzz”](https://github.com/golang/go/issues?q=is%3Aopen+is%3Aissue+label%3Afuzz)
to stay up-to-date on existing bugs and missing features.

Please be aware that fuzzing can consume a lot of memory and may impact your
machine’s performance while it runs. `go test -fuzz` defaults to running fuzzing
in `$GOMAXPROCS` processes in parallel. You may lower the number of processes
used while fuzzing by explicitly setting the `-parallel` flag with `go test`.
Read the documentation for the `go test` command by running `gotip help
testflag` if you want more information.

Also be aware that the fuzzing engine writes values that expand test coverage to
a fuzz cache directory within `$GOCACHE/fuzz` while it runs. There is currently
no limit to the number of files or total bytes that may be written to the fuzz
cache, so it may occupy a large amount of storage (i.e. several GBs). You can
clear the fuzz cache by running `gotip clean -fuzzcache`.

## What’s next?

This feature will become available starting in Go 1.18.

If you experience any problems or have an idea for a feature, please [file an
issue](https://github.com/golang/go/issues/new/?&labels=fuzz).

For discussion and general feedback about the feature, you can also participate
in the [#fuzzing channel](https://gophers.slack.com/archives/CH5KV1AKE) in
Gophers Slack.

Happy fuzzing!
