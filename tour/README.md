# Go Tour

[![Go Reference](https://pkg.go.dev/badge/golang.org/x/website/tour.svg)](https://pkg.go.dev/golang.org/x/website/tour)

A Tour of Go is an introduction to the Go programming language. Visit
https://go.dev/tour/ to start the tour.

## Download/Install

To install the tour from source, first
[install Go](https://go.dev/doc/install) and then run:

	go install golang.org/x/website/tour@latest

This will place a `tour` binary in your
[GOPATH](https://go.dev/cmd/go/#hdr-GOPATH_and_Modules)'s `bin` directory.
The tour program can be run offline.

## Send Patches

This repository uses Gerrit for code changes. To learn how to submit changes to
this repository, see https://go.dev/doc/contribute.

To run the tour server locally:

	go run .

Your browser should now open. If not, please visit [http://localhost:3999/](http://localhost:3999).

## Report Issues

The issue tracker for the tour's code is located at https://github.com/golang/go/issues.
Prefix your issue with "x/website/tour:" in the subject line, so it is easy to find.

Issues with the tour's content itself should be reported in the issue tracker
at https://github.com/golang/tour/issues.

## Deploying

Each time a CL is reviewed and submitted, the tour is automatically deployed to App Engine
as part of the main go.dev web site.
See [../cmd/golangorg/README.md](../cmd/golangorg/README.md#deploying-to-go_dev-and-golang_org) for details.

## License

Unless otherwise noted, the go-tour source files are distributed
under the BSD-style license found in the LICENSE file.
