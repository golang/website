---
title: Get familiar with workspaces
date: 2022-04-05
by:
- Beth Brown, for the Go team
tags:
- go
- workspaces
- go1.18
summary: Learn about Go workspaces and some of the workflows they enable.
---

Go 1.18 adds workspace mode to Go, which lets you work on multiple modules
simultaneously.

You can get Go 1.18 by visiting the [download](https://go.dev/dl/) page. The
[release notes](https://go.dev/doc/go1.18) have more details about all the
changes.

## Workspaces

[Workspaces](https://go.dev/ref/mod#workspaces) in Go 1.18 let you work on
multiple modules simultaneously without having to edit `go.mod` files for each
module. Each module within a workspace is treated as a main module when
resolving dependencies.

Previously, to add a feature to one module and use it in another module, you
needed to either publish the changes to the first module, or [edit the
go.mod](https://go.dev/doc/tutorial/call-module-code) file of the dependent
module with a `replace` directive for your local, unpublished module changes. In
order to publish without errors, you had to remove the `replace` directive from
the dependent module's `go.mod` file after you published the local changes to
the first module.

With Go workspaces, you control all your dependencies using a `go.work` file in
the root of your workspace directory. The `go.work` file has `use` and
`replace` directives that override the individual `go.mod` files, so there is
no need to edit each `go.mod` file individually.

You create a workspace by running `go work init` with a list of module
directories as space-separated arguments. The workspace doesn't need to contain
the modules you're working with. The` init` command creates a `go.work` file
that lists modules in the workspace.  If you run `go work init` without
arguments, the command creates an empty workspace.

To add modules to the workspace, run `go work use [moddir]` or manually edit
the `go.work` file. Run `go work use -r .` to recursively add directories in the
argument directory with a `go.mod` file to your workspace. If a directory
doesn't have a `go.mod` file, or no longer exists, the `use` directive for that
directory is removed from your `go.work` file.

The syntax of a `go.work` file is similar to a `go.mod` file and contains the
following directives:

-  `go`: the go toolchain version e.g. `go 1.18`
-  `use`: adds a module on disk to the set of main modules in a workspace.
    Its argument is a relative path to the directory containing the module's
    `go.mod` file. A `use` directive doesn't add modules in subdirectories of
    the specified directory.
-  `replace`: Similar to a `replace` directive in a `go.mod` file, a
    `replace` directive in a `go.work` file replaces the contents of a
    _specific version_ of a module, or _all versions_ of a module, with
    contents found elsewhere.

## Workflows

Workspaces are flexible and support a variety of workflows. The following
sections are a brief overview of the ones we think will be the most common.

### Add a feature to an upstream module and use it in your own module

1. Create a directory for your workspace.
1. Clone the upstream module you want to edit.
1. Add your feature to the local version of the upstream module.
1. Run `go work init [path-to-upstream-mod-dir]` in the workspace folder.
1. Make changes to your own module in order to implement the feature added
    to the upstream module.
1. Run `go work use [path-to-your-module]` in the workspace folder.

   The `go work use` command adds the path to your module to your `go.work`
   file:

    ```
    go 1.18

    use (
           ./path-to-upstream-mod-dir
           ./path-to-your-module
    )
    ```

1. Run and test your module using the new feature added to the upstream module.
1. Publish the upstream module with the new feature.
1. Publish your module using the new feature.

### Work with multiple interdependent modules in the same repository

While working on multiple modules in the same repository, the `go.work` file
defines the workspace instead of using `replace` directives in each module's
`go.mod` file.

1. Create a directory for your workspace.
1. Clone the repository with the modules you want to edit. The modules don't
    have to be in your workspace folder as you specify the relative path to
    each with the `use` directive.
1. Run `go work init [path-to-module-one] [path-to-module-two]` in your
    workspace directory.

   Example: You are working on `example.com/x/tools/groundhog` which depends
   on other packages in the `example.com/x/tools` module.

   You clone the repository and then run `go work init tools tools/groundhog` in
  your workspace folder.

   The contents of your `go.work` file resemble the following:

   ```
   go 1.18

   use (
           ./tools
           ./tools/groundhog
   )
   ```

   Any local changes made in the `tools` module will be used by
    `tools/groundhog` in your workspace.

### Switching between dependency configurations

To test your modules with different dependency configurations you can either
create multiple workspaces with separate `go.work` files, or keep one workspace
and comment out the `use` directives you don't want in a single `go.work`
file.

To create multiple workspaces:


1. Create separate directories for different dependency needs.
1. Run `go work init` in each of your workspace directories.
1. Add the dependencies you want within each directory via `go work use
    [path-to-dependency]`.
1. Run `go run [path-to-your-module]` in each workspace directory to use the
    dependencies specified by its `go.work` file.

To test out different dependencies within the same workspace, open the `go.work`
file and add or comment out the desired dependencies.

### Still using GOPATH?

Maybe using workspaces will change your mind. `GOPATH` users can resolve their
dependencies using a `go.work` file located at the base of their `GOPATH`
directory. Workspaces don't aim to completely recreate all `GOPATH` workflows,
but they can create a setup that shares some of the convenience of `GOPATH`
while still providing the benefits of modules.

To create a workspace for GOPATH:

1. Run `go work init` in the root of your `GOPATH` directory.
1. To use a local module or specific version as a dependency in your
    workspace, run `go work use [path-to-module]`.
1. To replace existing dependencies in your modules' `go.mod` files use
   `go work replace [path-to-module]`.
1. To add all the modules in your GOPATH or any directory, run `go work use
    -r` to recursively add directories with a `go.mod` file to your workspace.
    If a directory doesn't have a `go.mod` file, or no longer exists, the `use`
    directive for that directory is removed from your `go.work` file.

> Note: If you have projects without `go.mod` files that you want to add to
the workspace, change into their project directory and run `go mod init`,
then add the new module to your workspace with `go work use [path-to-module].`

## Workspace commands

Along with `go work init` and `go work use`, Go 1.18 introduces the following
commands for workspaces:

-  `go work sync`: pushes the dependencies in the `go.work` file back into
    the `go.mod` files of each workspace module.
-  `go work edit`: provides a command-line interface for editing `go.work`,
    for use primarily by tools or scripts.

Module-aware build commands and some `go mod` subcommands examine the `GOWORK`
environment variable to determine if they are in a workspace context.

Workspace mode is enabled if the `GOWORK` variable names a path to a file that
ends in `.work`. To determine which `go.work` file is being used, run
`go env GOWORK`. The output is empty if the `go` command is not in workspace
mode.

When workspace mode is enabled, the `go.work` file is parsed to determine the
three parameters for workspace mode: A Go version, a list of directories, and a
list of replacements.

Some commands to try in workspace mode (provided you already know what they
do!):

```
go work init
go work sync
go work use
go list
go build
go test
go run
go vet
```

## Editor experience improvements

We're particularly excited about the upgrades to Go's language server
[gopls](https://pkg.go.dev/golang.org/x/tools/gopls) and the
[VSCode Go extension](https://marketplace.visualstudio.com/items?itemName=golang.go)
that make working with multiple modules in an LSP-compatible editor a smooth
and rewarding experience.

Find references, code completion, and go to definitions work across modules
within the workspace. Version [0.8.1](https://github.com/golang/tools/releases/tag/gopls%2Fv0.8.1)
of `gopls` introduces diagnostics, completion, formatting, and hover for
`go.work` files. You can take advantage of these gopls features with any
[LSP](https://microsoft.github.io/language-server-protocol/)-compatible editor.

#### Editor specific notes

-  The latest [vscode-go
    release](https://github.com/golang/vscode-go/releases/tag/v0.32.0) allows
    quick access to your workspace's `go.work` file via the Go status bar's
    Quick Pick menu.

![Access the go.work file via the Go status bar's Quick Pick menu](https://user-images.githubusercontent.com/4999471/157268414-fba63843-5a14-44ba-be82-d42765568856.gif)

-  [GoLand](https://www.jetbrains.com/go/) supports workspaces and has
    plans to add syntax highlighting and code completion for `go.work` files.

For more information on using `gopls` with different editors see the `gopls`[
documentation](https://pkg.go.dev/golang.org/x/tools/gopls#readme-editors).

## What's next?

-  Download and install [Go 1.18](https://go.dev/dl/).
-  Try using [workspaces](https://go.dev/ref/mod#workspaces) with the [Go
    workspaces Tutorial](https://go.dev/doc/tutorial/workspaces).
-  If you encounter any problems with workspaces, or want to suggest
    something, file an [issue](https://github.com/golang/go/issues/new/choose).
-  Read the
    [workspace maintenance documentation](https://pkg.go.dev/cmd/go#hdr-Workspace_maintenance).
-  Explore module commands for [working outside of a single
    module](https://go.dev/ref/mod#commands-outside) including `go work init`,
    `go work sync` and more.
