---
title: "Go on App Engine: tools, tests, and concurrency"
date: 2013-12-13
by:
- Andrew Gerrand
- Johan Euphrosine
tags:
- appengine
summary: Announcing improvements to Go on App Engine.
---

## Background

When we [launched Go for App Engine](https://blog.golang.org/go-and-google-app-engine)
in May 2011 the SDK was just a modified version of the Python SDK.
At the time, there was no canonical way to build or organize Go programs, so it
made sense to take the Python approach. Since then Go 1.0 was released,
including the [go tool](/cmd/go/) and a
[convention](/doc/code.html) for organizing Go programs.

In January 2013 we announced
[better integration](https://blog.golang.org/the-app-engine-sdk-and-workspaces-gopath)
between the Go App Engine SDK and the go tool, promoting the use of
conventional import paths in App Engine apps and making it possible to use "go
get" to fetch app dependencies.

With the recent release of App Engine 1.8.8 we are pleased to announce more
improvements to the developer experience for Go on App Engine.

## The goapp tool

The Go App Engine SDK now includes the "goapp" tool, an App Engine-specific
version of the "go" tool. The new name permits users to keep both the regular
"go" tool and the "goapp" tool in their system PATH.

In addition to the existing "go" tool [commands](/cmd/go/),
the "goapp" tool provides new commands for working with App Engine apps.
The "[goapp serve](https://developers.google.com/appengine/docs/go/tools/devserver)"
command starts the local development server and the
"[goapp deploy](https://developers.google.com/appengine/docs/go/tools/uploadinganapp)"
command uploads an app to App Engine.

The main advantages offered by the "goapp serve" and "goapp deploy" commands
are a simplified user interface and consistency with existing commands like
"go get" and "go fmt".
For example, to run a local instance of the app in the current directory, run:

	$ goapp serve

To upload it to App Engine:

	$ goapp deploy

You can also specify the Go import path to serve or deploy:

	$ goapp serve github.com/user/myapp

You can even specify a YAML file to serve or deploy a specific
[module](https://developers.google.com/appengine/docs/go/modules/):

	$ goapp deploy mymodule.yaml

These commands can replace most uses of `dev_appserver.py` and `appcfg.py`,
although the Python tools are still available for their less common uses.

## Local unit testing

The Go App Engine SDK now supports local unit testing, using Go's native
[testing package](https://developers.google.com/appengine/docs/go/tools/localunittesting)
and the "[go test](/cmd/go/#hdr-Test_packages)" command
(provided as "goapp test" by the SDK).

Furthermore, you can now write tests that use App Engine services.
The [aetest package](https://developers.google.com/appengine/docs/go/tools/localunittesting#Go_Introducing_the_aetest_package)
provides an appengine.Context value that delegates requests to a temporary
instance of the development server.

For more information about using "goapp test" and the aetest package, see the
[Local Unit Testing for Go documentation](https://developers.google.com/appengine/docs/go/tools/localunittesting).
Note that the aetest package is still in its early days;
we hope to add more features over time.

## Better concurrency support

It is now possible to configure the number of concurrent requests served by
each of your app's dynamic instances by setting the
[`max_concurrent_requests`](https://developers.google.com/appengine/docs/go/modules/#max_concurrent_requests) option
(available to [Automatic Scaling modules](https://developers.google.com/appengine/docs/go/modules/#automatic_scaling) only).

Here's an example `app.yaml` file:

	application: maxigopher
	version: 1
	runtime: go
	api_version: go1
	automatic_scaling:
	  max_concurrent_requests: 100

This configures each instance of the app to serve up to 100 requests
concurrently (up from the default of 10). You can configure Go instances to
serve up to a maximum of 500 concurrent requests.

This setting allows your instances to handle more simultaneous requests by
taking advantage of Go's efficient handling of concurrency, which should yield
better instance utilization and ultimately fewer billable instance hours.

## Conclusion

With these changes Go on App Engine is more convenient and efficient than ever,
and we hope you enjoy the improvements. Please join the
[google-appengine-go group](http://groups.google.com/group/google-appengine-go/)
to raise questions or discuss these changes with the engineering team and the
rest of the community.
