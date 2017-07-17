# Go Blog

This repository holds the Go Blog server code and content.

## Download/Install

The easiest way to install is to run `go get -u golang.org/x/blog`. You can also
manually git clone the repository to $GOPATH/src/golang.org/x/blog.

## Running Locally

To run the blog server locally:

	go build -o blog.exe ./blog
	./blog.exe

and then visit [http://localhost:8080/](http://localhost:8080) in your browser.

Note that blog.exe caches all content and you must kill it (Ctrl-C on Unix
systems) and restart it after editing any files to view changes.

## Contributing

Articles on the blog should have broad interest to the Go community, and
are mainly written by Go contributors. We encourage you to share your
experiences using Go on your own website, and [to share them with the Go
community][community]. [Hugo][hugo] is a static site server written in Go that
makes it easy to write and share your stories.

[community]: https://golang.org/help/
[hugo]: https://gohugo.io/

## Report Issues / Send Patches

This repository uses Gerrit for code changes. To learn how to submit changes to
this repository, see https://golang.org/doc/contribute.html.

The main issue tracker for the blog is located at
https://github.com/golang/go/issues. Prefix your issue with "blog:" in the
subject line, so it is easy to find.
