# Go Blog

[![Go Reference](https://pkg.go.dev/badge/golang.org/x/website/blog.svg)](https://pkg.go.dev/golang.org/x/website/blog)

This directory holds the Go Blog server code and content.

## Download/Install

The easiest way to install is to run `go get -u golang.org/x/website/blog`. You can also
manually git clone the repository to \$GOPATH/src/golang.org/x/website.

## Running Locally

To run the blog server locally:

```
go run . -reload
```

and then visit [http://localhost:8080/](http://localhost:8080) in your browser.

## Contributing

Articles are written in the [x/tools/present][present] format.
Articles on the blog should have broad interest to the Go community, and
are mainly written by Go contributors. We encourage you to share your
experiences using Go on your own website, and [to share them with the Go
community][community]. [Hugo][hugo] is a static site server written in Go that
makes it easy to write and share your stories.

[present]: https://godoc.org/golang.org/x/tools/present
[community]: https://golang.org/help/
[hugo]: https://gohugo.io/

## Report Issues / Send Patches

This repository uses Gerrit for code changes. To learn how to submit changes to
this repository, see https://golang.org/doc/contribute.html.

The main issue tracker for the blog is located at
https://github.com/golang/go/issues. Prefix your issue with "x/website/blog:" in the
subject line, so it is easy to find.

## Deploying

Each time a CL is reviewed and submitted, the blog is automatically redeployed to
https://blog.golang.org/.

If the automatic redeploy is not working, or to check on the status of a redeploy,
see the “website-redeploy-blog” trigger in the
[Cloud Build console](https://console.cloud.google.com/cloud-build/builds?project=golang-org).

