---
title: Announcing App Engine’s New Go 1.11 Runtime
date: 2018-10-16
by:
- Eno Compton
- Tyler Bui-Palsulich
tags:
- appengine
summary: Google Cloud is announcing a new Go 1.11 runtime for App Engine, with fewer limits on app structure.
---


[App Engine](https://cloud.google.com/appengine/) launched
[experimental support for Go](https://blog.golang.org/go-and-google-app-engine)
in 2011. In the subsequent years, the Go community has grown significantly and
has settled on idiomatic
patterns for cloud-based applications. Today, Google Cloud is
[announcing a new Go 1.11 runtime](https://cloud.google.com/blog/products/application-development/go-1-11-is-now-available-on-app-engine)
for the App Engine standard environment that provides all the
power of App Engine—things like paying only for what you use, automatic scaling,
and managed infrastructure—while supporting idiomatic Go.

Starting with Go 1.11, Go on App Engine has no limits on application structure,
supported packages, `context.Context` values, or HTTP clients. Write your Go
application however you prefer, add an `app.yaml` file, and your app is ready
to deploy on App Engine.
[Specifying Dependencies](https://cloud.google.com/appengine/docs/standard/go111/specifying-dependencies)
describes how the new runtime
supports [vendoring](/cmd/go/#hdr-Vendor_Directories) and
[modules](/doc/go1.11#modules) (experimental) for dependency
management.

Along with [Cloud Functions support for Go](https://twitter.com/kelseyhightower/status/1035278586754813952)
(more on that in a future post), App Engine provides a compelling way to run Go
code on Google Cloud Platform (GCP) with no concern for the underlying
infrastructure.

Let’s take a look at creating a small application for App Engine. For the
example here, we assume a `GOPATH`-based workflow, although Go modules have
[experimental support](https://cloud.google.com/appengine/docs/standard/go111/specifying-dependencies)
as well.

First, you create the application in your `GOPATH`:

{{code "appengine/main.go"}}

The code contains an idiomatic setup for a small HTTP server that responds with
“Hello, 世界.” If you have previous App Engine experience, you’ll notice the
absence of any call to `appengine.Main()`, which is now entirely optional.
Furthermore, the application code is completely portable—there are no ties to
the infrastructure that your application is deployed on.

If you need to use external dependencies, you can add those dependencies to a
`vendor` directory or to a `go.mod` file, both of which the new runtime
supports.

With the application code complete, create an `app.yaml` file to specify
the runtime:

	runtime: go111

Finally, set your machine up with a Google Cloud Platform account:

  - Create an account with [GCP](https://cloud.google.com).
  - [Create a project](https://cloud.google.com/resource-manager/docs/creating-managing-projects).
  - Install the [Cloud SDK](https://cloud.google.com/sdk/) on your system.

With all the setup complete, you can deploy using one command:

	gcloud app deploy

We think Go developers will find the new Go 1.11 runtime for App Engine an
exciting addition to the available options to run Go applications. There is a
[free tier](https://cloud.google.com/free/). Check out the
[getting started guide](https://cloud.google.com/appengine/docs/standard/go111/building-app/)
or the
[migration guide](https://cloud.google.com/appengine/docs/standard/go111/go-differences)
and deploy an app to the new runtime today!
