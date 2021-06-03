# Go Tour

[![Go Reference](https://pkg.go.dev/badge/golang.org/x/tour.svg)](https://pkg.go.dev/golang.org/x/tour)

A Tour of Go is an introduction to the Go programming language. Visit
https://tour.golang.org to start the tour.

## Download/Install

To install the tour from source, first
[install Go](https://golang.org/doc/install) and then run:

	$ go get golang.org/x/tour

This will place a `tour` binary in your
[workspace](https://golang.org/doc/code.html#Workspaces)'s `bin` directory.
The tour program can be run offline.

## Contributing

Contributions should follow the same procedure as for the Go project:
https://golang.org/doc/contribute.html

To run the tour server locally:

```sh
go run .
```

Your browser should now open. If not, please visit [http://localhost:3999/](http://localhost:3999).


## Report Issues / Send Patches

This repository uses Gerrit for code changes. To learn how to submit changes to
this repository, see https://golang.org/doc/contribute.html.

The issue tracker for the tour's code is located at https://github.com/golang/go/issues.
Prefix your issue with "x/tour:" in the subject line, so it is easy to find.

Issues with the tour's content itself should be reported in the issue tracker
at https://github.com/golang/tour/issues.

## Deploying

Each time a CL is reviewed and submitted, the tour is automatically deployed to App Engine.
If the CL is submitted with a Website-Publish +1 vote,
the new deployment automatically becomes https://tour.golang.org/.
Otherwise, the new deployment can be found in the
[App Engine versions list](https://console.cloud.google.com/appengine/versions?project=golang-org&serviceId=tour) and verified and manually promoted.

If the automatic deployment is not working, or to check on the status of a pending deployment,
see the “website-redeploy-tour” trigger in the
[Cloud Build console](https://console.cloud.google.com/cloud-build/builds?project=golang-org).

## License

Unless otherwise noted, the go-tour source files are distributed
under the BSD-style license found in the LICENSE file.
