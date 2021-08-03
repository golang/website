---
title: Fuzzing is Beta Ready
date: 2021-06-03
by:
- Katie Hockman
- Jay Conrod
tags:
- fuzz
- testing
summary: Native Go fuzzing is now ready for beta testing in the dev.fuzz development branch.
---


We are excited to announce that native fuzzing is ready for beta testing in its
development branch, [dev.fuzz](https://github.com/golang/go/tree/dev.fuzz)!

Fuzzing is a type of automated testing which continuously manipulates inputs to
a program to find issues such as panics or bugs. These semi-random data
mutations can discover new code coverage that existing unit tests may miss, and
uncover edge case bugs which would otherwise go unnoticed. Since fuzzing can
reach these edge cases, fuzz testing is particularly valuable for finding
security exploits and vulnerabilities.

See
[golang.org/s/draft-fuzzing-design](https://golang.org/s/draft-fuzzing-design)
for more details about this feature.


## Getting started

To get started, you may run the following

	$ go get golang.org/dl/gotip
	$ gotip download dev.fuzz

This builds the Go toolchain from the dev.fuzz development branch, and won’t be
needed once the code is merged to the master branch in the future. After running
this, `gotip` can act as a drop-in replacement for the `go` command. You can now
run commands like

	$ gotip test -fuzz=FuzzFoo

There will be ongoing development and bug fixes in the dev.fuzz branch, so you
should regularly run `gotip download dev.fuzz` to use the latest code.

For compatibility with released versions of Go, use the gofuzzbeta build tag
when committing source files containing fuzz targets to your repository. This
tag is enabled by default at build-time in the dev.fuzz branch. See the [go
command documentation about build
tags](https://golang.org/cmd/go/#hdr-Build_constraints) if you have questions
about how to use them.

	// +build gofuzzbeta

## Writing a fuzz target

A fuzz target must be in a \*\_test.go file as a function in the form `FuzzXxx`.
This function must be passed a` *testing.F` argument, much like a `*testing.T`
argument is passed to a `TestXxx` function.

Below is an example of a fuzz target that’s testing the behavior of the [net/url
package](https://pkg.go.dev/net/url#ParseQuery).

	// +build gofuzzbeta

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

You can read more about the fuzzing APIs with go doc

	gotip doc testing
	gotip doc testing.F
	gotip doc testing.F.Add
	gotip doc testing.F.Fuzz

## Expectations

This is a beta release in a development branch, so you should expect some bugs
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
cache, so it may occupy a large amount of storage (ie. several GBs). You can
clear the fuzz cache by running `gotip clean -fuzzcache`.

## What’s next?

This feature will not be available in the upcoming Go release (1.17), but there
are plans to land this in a future Go release. We hope that this working
prototype will allow Go developers to start writing fuzz targets and provide
helpful feedback about the design in preparation for a merge to master.

If you experience any problems or have an idea for a feature request, please
[file an
issue](https://github.com/golang/go/issues/new/?&labels=fuzz&title=%5Bdev%2Efuzz%5D&milestone=backlog).

For discussion and general feedback about the feature, you can also participate
in the [#fuzzing channel](https://gophers.slack.com/archives/CH5KV1AKE) in
Gophers Slack.

Happy fuzzing!
