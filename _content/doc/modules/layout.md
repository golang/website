<!--{
  "Title": "Organizing a Go module"
}-->

A common question developers new to Go have is "How do I organize my Go
project?", in terms of the layout of files and folders. The goal of this
document is to provide some guidelines that will help answer this question. To
make the most of this document, make sure you're familiar with the basics of Go
modules by reading [the tutorial](/doc/tutorial/create-module) and
[managing module source](/doc/modules/managing-source).

Go projects can include packages, command-line programs or a combination of the
two. This guide is organized by project type.

### Basic package

A basic Go package has all its code in the project's root directory. The project
consists of a single module, which consists of a single package. The package
name matches the last path component of the module name. For a very simple
package requiring a single Go file, the project structure is:

```
project-root-directory/
  go.mod
  modname.go
  modname_test.go
```

_[throughout this document, file/package names are entirely arbitrary]_

Assuming this directory is uploaded to a GitHub repository at
`github.com/someuser/modname`, the `module` line in the `go.mod` file should say
`module github.com/someuser/modname`.

The code in `modname.go` declares the package with:

```
package modname

// ... package code here
```

Users can then rely on this package by `import`-ing it in their Go code with:

```
import "github.com/someuser/modname"
```

A Go package can be split into multiple files, all residing within the same
directory, e.g.:

```
project-root-directory/
  go.mod
  modname.go
  modname_test.go
  auth.go
  auth_test.go
  hash.go
  hash_test.go
```

All the files in the directory declare `package modname`.

### Basic command

A basic executable program (or command-line tool) is structured according to its
complexity and code size. The simplest program can consist of a single Go file
where `func main` is defined. Larger programs can have their code split across
multiple files, all declaring `package main`:

```
project-root-directory/
  go.mod
  auth.go
  auth_test.go
  client.go
  main.go
```

Here the `main.go` file contains `func main`, but this is just a convention. The
"main" file can also be called `modname.go` (for an appropriate value of
`modname`) or anything else.

Assuming this directory is uploaded to a GitHub repository at
`github.com/someuser/modname`, the `module` line in the `go.mod` file should
say:

```
module github.com/someuser/modname
```

And a user should be able to install it on their machine with:

```
$ go install github.com/someuser/modname@latest
```

### Package or command with supporting packages

Larger packages or commands may benefit from splitting off some functionality
into supporting packages. Initially, it's recommended placing such packages into
a directory named `internal`;
[this prevents](https://pkg.go.dev/cmd/go#hdr-Internal_Directories) other
modules from depending on packages we don't necessarily want to expose and
support for external uses. Since other projects cannot import code from our
`internal` directory, we're free to refactor its API and generally move things
around without breaking external users. The project structure for a package is
thus:

```
project-root-directory/
  internal/
    auth/
      auth.go
      auth_test.go
    hash/
      hash.go
      hash_test.go
  go.mod
  modname.go
  modname_test.go
```

The `modname.go` file declares `package modname`, `auth.go` declares `package
auth` and so on. `modname.go` can import the `auth` package as follows:

```
import "github.com/someuser/modname/internal/auth"
```

The layout for a command with supporting packages in an `internal` directory is
very similar, except that the file(s) in the root directory declare `package
main`.

### Multiple packages

A module can consist of multiple importable packages; each package has its own
directory, and can be structured hierarchically. Here's a sample project
structure:

```
project-root-directory/
  go.mod
  modname.go
  modname_test.go
  auth/
    auth.go
    auth_test.go
    token/
      token.go
      token_test.go
  hash/
    hash.go
  internal/
    trace/
      trace.go
```

As a reminder, we assume that the `module` line in `go.mod` says:

```
module github.com/someuser/modname
```

The `modname` package resides in the root directory, declares `package modname`
and can be imported by users with:

```
import "github.com/someuser/modname"
```

Sub-packages can be imported by users as follows:

```
import "github.com/someuser/modname/auth"
import "github.com/someuser/modname/auth/token"
import "github.com/someuser/modname/hash"
```

Package `trace` that resides in `internal/trace` cannot be imported outside this
module. It's recommended to keep packages in `internal` as much as possible.

### Multiple commands

Multiple programs in the same repository will typically have separate directories:

```
project-root-directory/
  go.mod
  internal/
    ... shared internal packages
  prog1/
    main.go
  prog2/
    main.go
```

In each directory, the program's Go files declare `package main`. A top-level
`internal` directory can contain shared packages used by all commands in the
repository.

Users can install these programs as follows:

```
$ go install github.com/someuser/modname/prog1@latest
$ go install github.com/someuser/modname/prog2@latest
```

A common convention is placing all commands in a repository into a `cmd`
directory; while this isn't strictly necessary in a repository that consists
only of commands, it's very useful in a mixed repository that has both commands
and importable packages, as we will discuss next.

### Packages and commands in the same repository

Sometimes a repository will provide both importable packages and installable
commands with related functionality. Here's a sample project structure for such
a repository:

```
project-root-directory/
  go.mod
  modname.go
  modname_test.go
  auth/
    auth.go
    auth_test.go
  internal/
    ... internal packages
  cmd/
    prog1/
      main.go
    prog2/
      main.go
```

Assuming this module is called `github.com/someuser/modname`, users can now both
import packages from it:

```
import "github.com/someuser/modname"
import "github.com/someuser/modname/auth"
```

And install programs from it:

```
$ go install github.com/someuser/modname/cmd/prog1@latest
$ go install github.com/someuser/modname/cmd/prog2@latest
```

### Server project

Go is a common language choice for implementing *servers*. There is a very large
variance in the structure of such projects, given the many aspects of server
development: protocols (REST? gRPC?), deployments, front-end files,
containerization, scripts and so on. We will focus our guidance here on the
parts of the project written in Go.

Server projects typically won't have packages for export, since a server is
usually a self-contained binary (or a group of binaries). Therefore, it's
recommended to keep the Go packages implementing the server's logic in the
`internal` directory. Moreover, since the project is likely to have many other
directories with non-Go files, it's a good idea to keep all Go commands together
in a `cmd` directory:

```
project-root-directory/
  go.mod
  internal/
    auth/
      ...
    metrics/
      ...
    model/
      ...
  cmd/
    api-server/
      main.go
    metrics-analyzer/
      main.go
    ...
  ... the project's other directories with non-Go code
```

In case the server repository grows packages that become useful for sharing with
other projects, it's best to split these off to separate modules.
