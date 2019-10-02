# Go Blog

This repository holds the Go Blog server code and content.

## Download/Install

The easiest way to install is to run `go get -u golang.org/x/blog`. You can also
manually git clone the repository to \$GOPATH/src/golang.org/x/blog.

## Running Locally

To run the blog server locally:

```
go build && ./blog -reload
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
https://github.com/golang/go/issues. Prefix your issue with "x/blog:" in the
subject line, so it is easy to find.

## Deploying

1.	To deploy blog.golang.org, run:

	```
	GO111MODULE=on gcloud --project=golang-org app deploy --no-promote app.yaml
	```

	This will create a new version, which can be viewed within the
	[golang-org GCP project](https://console.cloud.google.com/appengine/versions?project=golang-org&serviceId=blog).

2.	Check that the deployed version looks OK (click the version link in GCP).

3.	If all is well, click "Migrate Traffic" to move 100% of the blog.golang.org
	traffic to the new version.

4.	You're done.
