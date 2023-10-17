---
title: Experimenting with project templates
date: 2023-07-31
by:
- Cameron Balahan
summary: Announcing golang.org/x/tools/cmd/gonew, an experimental tool for starting new Go projects from predefined templates
---

When you start a new project in Go, you might begin by cloning an existing project.
That way, you can start with something that already works,
making incremental changes instead of starting from scratch.

For a long time now, we have heard from Go developers that getting started
is often the hardest part.
New developers coming from other languages expect guidance on a default project layout,
experienced developers working on teams expect consistency in their projects’ dependencies,
and developers of all kinds expect an easy way to try new products and services
without having to copy and paste from samples on the web.

To that end, today we published `gonew`, an experimental tool for instantiating
new projects in Go from predefined templates.
Anyone can write templates, which are packaged and distributed as modules,
leveraging the Go module proxy and checksum database for better security and availability.

The prototype `gonew` is intentionally minimal:
what we have released today is an extremely limited prototype meant to provide
a base from which we can gather feedback and community direction.
Try it out, [tell us what you think](https://go.dev/s/gonew-feedback),
and help us build a more useful tool for everyone.

## Getting started

Start by installing `gonew` using [`go install`](https://pkg.go.dev/cmd/go#hdr-Compile_and_install_packages_and_dependencies):

```
$ go install golang.org/x/tools/cmd/gonew@latest
```

To copy an existing template, run `gonew` in your new project’s parent
directory with two arguments:
first, the path to the template you wish to copy,
and second, the module name of the project you are creating. For example:

```
$ gonew golang.org/x/example/helloserver example.com/myserver
$ cd ./myserver
```

And then you can read and edit the files in `./myserver` to customize.

We’ve written two templates to get you started:

- [hello](https://pkg.go.dev/golang.org/x/example/hello):
  A command line tool that prints a greeting,
  with customization flags.
- [helloserver](https://pkg.go.dev/golang.org/x/example/helloserver): An HTTP server that serves greetings.

## Write your own templates

Writing your own template is as easy as [creating any other module](/doc/tutorial/create-module) in Go.
Check out the examples we linked above to get started.

There are also examples available from the [Google Cloud](https://github.com/GoogleCloudPlatform/go-templates)
and [Service Weaver](https://github.com/ServiceWeaver/template) teams.

## Next steps

Please try out `gonew` and let us know how we can make it better and more useful.
Remember, `gonew` is just an experiment for now;
we need your [feedback to get it right](https://go.dev/s/gonew-feedback).

