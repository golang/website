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

The Google Cloud project triggers a fresh deploy of the blog on each submit
but that deployment is published to a temporary URL.

To publish the blog to blog.golang.org, you need access to the
Cloud Console for the golang-org project.
Then:

1. Visit the
   [builds list](https://console.cloud.google.com/cloud-build/builds?project=golang-org&query=trigger_id%3D%22c99674d3-32c1-4aec-ade4-ae2d5a844369%22)
   and click on the build hash for the most recent build
   with trigger name “Redeploy-blog-on-blog-commit”.

   Scrolling to the bottom of the build log, you will find a URL in a log message like

       Deployed service [blog] to [https://TEMPORARYURL.appspot.com]

2. Copy that URL and load it in your browser. Check that it looks OK.

3. Assuming it does, visit the
   [AppEngine versions list](https://console.cloud.google.com/appengine/versions?project=golang-org&serviceId=blog).
   Click “Migrate Traffic” on the new entry to move 100% of the blog.golang.org
   traffic to the new version.

4. You're done.
