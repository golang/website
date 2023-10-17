---
title: The App Engine SDK and workspaces (GOPATH)
date: 2013-01-09
by:
- Andrew Gerrand
tags:
- appengine
- tools
- gopath
summary: App Engine SDK 1.7.4 adds support for GOPATH-style workspaces.
---

## Introduction

When we released Go 1 we introduced the [go tool](/cmd/go/) and,
with it, the concept of workspaces.
Workspaces (specified by the GOPATH environment variable) are a convention
for organizing code that simplifies fetching,
building, and installing Go packages.
If you're not familiar with workspaces, please read [this article](/doc/code.html)
or watch [this screencast](http://www.youtube.com/watch?v=XCsL89YtqCs) before reading on.

Until recently, the tools in the App Engine SDK were not aware of workspaces.
Without workspaces the "[go get](/cmd/go/#hdr-Download_and_install_packages_and_dependencies)"
command cannot function,
and so app authors had to install and update their app dependencies manually. It was a pain.

This has all changed with version 1.7.4 of the App Engine SDK.
The [dev\_appserver](https://developers.google.com/appengine/docs/go/tools/devserver)
and [appcfg](https://developers.google.com/appengine/docs/go/tools/uploadinganapp)
tools are now workspace-aware.
When running locally or uploading an app,
these tools now search for dependencies in the workspaces specified by the
GOPATH environment variable.
This means you can now use "go get" while building App Engine apps,
and switch between normal Go programs and App Engine apps without changing
your environment or habits.

For example, let's say you want to build an app that uses OAuth 2.0 to authenticate
with a remote service.
A popular OAuth 2.0 library for Go is the [oauth2](https://godoc.org/golang.org/x/oauth2) package,
which you can install to your workspace with this command:

	go get golang.org/x/oauth2

When writing your App Engine app, import the oauth package just as you would in a regular Go program:

	import "golang.org/x/oauth2"

Now, whether running your app with the dev\_appserver or deploying it with appcfg,
the tools will find the oauth package in your workspace. It just works.

## Hybrid stand-alone/App Engine apps

The Go App Engine SDK builds on Go's standard [net/http](/pkg/net/http/)
package to serve web requests and,
as a result, many Go web servers can be run on App Engine with only a few changes.
For example, [godoc](/cmd/godoc/) is included in the
Go distribution as a stand-alone program,
but it can also run as an App Engine app (godoc serves [golang.org](/) from App Engine).

But wouldn't it be nice if you could write a program that is both a stand-alone
web server and an App Engine app? By using [build constraints](/pkg/go/build/#hdr-Build_Constraints), you can.

Build constraints are line comments that determine whether a file should
be included in a package.
They are most often used in code that handles a variety of operating systems
or processor architectures.
For instance, the [path/filepath](/pkg/path/filepath/)
package includes the file [symlink.go](/src/pkg/path/filepath/symlink.go),
which specifies a build constraint to ensure that it is not built on Windows
systems (which do not have symbolic links):

	// +build !windows

The App Engine SDK introduces a new build constraint term: "appengine". Files that specify

	// +build appengine

will be built by the App Engine SDK and ignored by the go tool. Conversely, files that specify

	// +build !appengine

are ignored by the App Engine SDK, while the go tool will happily build them.

The [goprotobuf](http://code.google.com/p/goprotobuf/) library uses this
mechanism to provide two implementations of a key part of its encode/decode machinery:
[pointer\_unsafe.go](http://code.google.com/p/goprotobuf/source/browse/proto/pointer_unsafe.go)
is the faster version that cannot be used on App Engine because it uses
the [unsafe package](/pkg/unsafe/),
while [pointer\_reflect.go](http://code.google.com/p/goprotobuf/source/browse/proto/pointer_reflect.go)
is a slower version that avoids unsafe by using the [reflect package](/pkg/reflect/) instead.

Let's take a simple Go web server and turn it into a hybrid app. This is main.go:

	package main

	import (
	    "fmt"
	    "net/http"
	)

	func main() {
	    http.HandleFunc("/", handler)
	    http.ListenAndServe("localhost:8080", nil)
	}

	func handler(w http.ResponseWriter, r *http.Request) {
	    fmt.Fprint(w, "Hello!")
	}

Build this with the go tool and you'll get a stand-alone web server executable.

The App Engine infrastructure provides its own main function that runs its
equivalent to ListenAndServe.
To convert main.go to an App Engine app, drop the call to ListenAndServe
and register the handler in an init function (which runs before main). This is app.go:

	package main

	import (
	    "fmt"
	    "net/http"
	)

	func init() {
	    http.HandleFunc("/", handler)
	}

	func handler(w http.ResponseWriter, r *http.Request) {
	    fmt.Fprint(w, "Hello!")
	}

To make this a hybrid app, we need to split it into an App Engine-specific part,
an stand-alone binary-specific part, and the parts common to both versions.
In this case, there is no App Engine-specific part,
so we split it into just two files:

app.go specifies and registers the handler function.
It is identical to the code listing above,
and requires no build constraints as it should be included in all versions of the program.

main.go runs the web server. It includes the "!appengine" build constraint,
as it must only be included when building the stand-alone binary.

	// +build !appengine

	package main

	import "net/http"

	func main() {
	    http.ListenAndServe("localhost:8080", nil)
	}

To see a more complex hybrid app, take a look at the [present tool](https://godoc.org/golang.org/x/tools/present).

## Conclusions

We hope these changes will make it easier to work on apps with external dependencies,
and to maintain code bases that contain both stand-alone programs and App Engine apps.
