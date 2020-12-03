# Go website

[![Go Reference](https://pkg.go.dev/badge/golang.org/x/website.svg)](https://pkg.go.dev/golang.org/x/website)

This repository holds the Go website server code and content.

## Checkout and Run

To download and run the golang.org web server locally:

 - `git clone https://go.googlesource.com/website`
 - `cd website`
 - `go run ./cmd/golangorg`
 - Open http://localhost:6060/ in your browser.

See [cmd/golangorg/README.md](cmd/golangorg/README.md) for more details.

## Changing Content

To make basic changes to the golang.org website content:

 - Make the changes you want in the `content/static` directory.
 - Stop any running `go run ./cmd/golangorg`.
 - `go generate ./content/static`
 - `go run ./cmd/golangorg`
 - Open http://localhost:6060/ in your browser.

See [content/README.md](content/README.md) for more sophisticated instructions.

## JS/CSS Formatting

This repository uses [prettier](https://prettier.io/) to format JS and CSS files.

The version of `prettier` used is 1.18.2.

It is encouraged that all JS and CSS code be run through this before submitting
a change. However, it is not a strict requirement enforced by CI.

## Report Issues / Send Patches

This repository uses Gerrit for code changes. To learn how to submit changes to
this repository, see https://golang.org/doc/contribute.html.

The main issue tracker for the website repository is located at
https://github.com/golang/go/issues. Prefix your issue with "x/website:" in the
subject line, so it is easy to find.
