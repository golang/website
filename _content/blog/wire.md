---
title: Compile-time Dependency Injection With Go Cloud's Wire
date: 2018-10-09
by:
- Robert van Gent
summary: How to use Wire, a dependency injection tool for Go.
---

## Overview

The Go team recently [announced](https://blog.golang.org/go-cloud) the
open source project [Go Cloud](https://github.com/google/go-cloud),
with portable Cloud APIs and tools for [open cloud](https://cloud.google.com/open-cloud/) development.
This post goes into more detail about Wire,
a dependency injection tool used in Go Cloud.

## What problem does Wire solve?

[Dependency injection](https://en.wikipedia.org/wiki/Dependency_injection)
is a standard technique for producing flexible and loosely coupled code,
by explicitly providing components with all of the dependencies they need to work.
In Go, this often takes the form of passing dependencies to constructors:

	// NewUserStore returns a UserStore that uses cfg and db as dependencies.
	func NewUserStore(cfg *Config, db *mysql.DB) (*UserStore, error) {...}

This technique works great at small scale,
but larger applications can have a complex graph of dependencies,
resulting in a big block of initialization code that's order-dependent but
otherwise not very interesting.
It's often hard to break up this code cleanly,
especially because some dependencies are used multiple times.
Replacing one implementation of a service with another can be painful because
it involves modifying the dependency graph by adding a whole new set of
dependencies (and their dependencies...),
and removing unused old ones.
In practice, making changes to initialization code in applications with
large dependency graphs is tedious and slow.

Dependency injection tools like Wire aim to simplify the management of initialization code.
You describe your services and their dependencies,
either as code or as configuration, then Wire processes the resulting graph
to figure out ordering and how to pass each service what it needs.
Make changes to an application's dependencies by changing a function signature
or adding or removing an initializer,
and then let Wire do the tedious work of generating initialization code
for the entire dependency graph.

## Why is this part of Go Cloud?

Go Cloud's goal is to make it easier to write portable Cloud applications
by providing idiomatic Go APIs for useful Cloud services.
For example, [blob.Bucket](https://godoc.org/github.com/google/go-cloud/blob)
provides a storage API with implementations for Amazon's S3 and Google Cloud Storage (GCS);
applications written using `blob.Bucket` can swap implementations without
changing their application logic.
However, the initialization code is inherently provider-specific,
and each provider has a different set of dependencies.

For example, [constructing a GCS `blob.Bucket`](https://godoc.org/github.com/google/go-cloud/blob/gcsblob#OpenBucket)
requires a `gcp.HTTPClient`,
which eventually requires `google.Credentials`,
while [constructing one for S3](https://godoc.org/github.com/google/go-cloud/blob/s3blob)
requires an `aws.Config`,
which eventually requires AWS credentials.
Thus, updating an application to use a different `blob.Bucket` implementation
involves exactly the kind of tedious update to the dependency graph that we described above.
The driving use case for Wire is to make it easy to swap implementations
of Go Cloud portable APIs,
but it's also a general-purpose tool for dependency injection.

## Hasn't this been done already?

There are a number of dependency injection frameworks out there.
For Go, [Uber's dig](https://github.com/uber-go/dig) and [Facebook's inject](https://github.com/facebookgo/inject)
both use reflection to do runtime dependency injection.
Wire was primarily inspired by Java's [Dagger 2](https://google.github.io/dagger/),
and uses code generation rather than reflection or [service locators](https://en.wikipedia.org/wiki/Service_locator_pattern).

We think this approach has several advantages:

  - Runtime dependency injection can be hard to follow and debug when the
    dependency graph gets complex.
    Using code generation means that the initialization code that's executed
    at runtime is regular,
    idiomatic Go code that's easy to understand and debug.
    Nothing is obfuscated by an intervening framework doing "magic".
    In particular, problems like forgetting a dependency become compile-time errors,
    not run-time errors.
  - Unlike [service locators](https://en.wikipedia.org/wiki/Service_locator_pattern),
    there's no need to make up arbitrary names or keys to register services.
    Wire uses Go types to connect components with their dependencies.
  - It's easier to avoid dependency bloat. Wire's generated code will only
    import the dependencies you need,
    so your binary won't have unused imports.
    Runtime dependency injectors can't identify unused dependencies until runtime.
  - Wire's dependency graph is knowable statically, which provides opportunities for tooling and visualization.

## How does it work?

Wire has two basic concepts: providers and injectors.

_Providers_ are ordinary Go functions that "provide" values given their dependencies,
which are described simply as parameters to the function.
Here's some sample code that defines three providers:

	// NewUserStore is the same function we saw above; it is a provider for UserStore,
	// with dependencies on *Config and *mysql.DB.
	func NewUserStore(cfg *Config, db *mysql.DB) (*UserStore, error) {...}

	// NewDefaultConfig is a provider for *Config, with no dependencies.
	func NewDefaultConfig() *Config {...}

	// NewDB is a provider for *mysql.DB based on some connection info.
	func NewDB(info *ConnectionInfo) (*mysql.DB, error) {...}

Providers that are commonly used together can be grouped into `ProviderSets`.
For example, it's common to use a default `*Config` when creating a `*UserStore`,
so we can group `NewUserStore` and `NewDefaultConfig` in a `ProviderSet`:

	var UserStoreSet = wire.ProviderSet(NewUserStore, NewDefaultConfig)

_Injectors_ are generated functions that call providers in dependency order.
You write the injector's signature, including any needed inputs as arguments,
and insert a call to `wire.Build` with the list of providers or provider
sets that are needed to construct the end result:

	func initUserStore() (*UserStore, error) {
		// We're going to get an error, because NewDB requires a *ConnectionInfo
		// and we didn't provide one.
		wire.Build(UserStoreSet, NewDB)
		return nil, nil  // These return values are ignored.
	}

Now we run go generate to execute wire:

	$ go generate
	wire.go:2:10: inject initUserStore: no provider found for ConnectionInfo (required by provider of *mysql.DB)
	wire: generate failed

Oops! We didn't include a `ConnectionInfo` or tell Wire how to build one.
Wire helpfully tells us the line number and types involved.
We can either add a provider for it to `wire.Build`,
or add it as an argument:

	func initUserStore(info ConnectionInfo) (*UserStore, error) {
		wire.Build(UserStoreSet, NewDB)
		return nil, nil  // These return values are ignored.
	}

Now `go generate` will create a new file with the generated code:

	// File: wire_gen.go
	// Code generated by Wire. DO NOT EDIT.
	//go:generate wire
	//+build !wireinject

	func initUserStore(info ConnectionInfo) (*UserStore, error) {
		defaultConfig := NewDefaultConfig()
		db, err := NewDB(info)
		if err != nil {
			return nil, err
		}
		userStore, err := NewUserStore(defaultConfig, db)
		if err != nil {
			return nil, err
		}
		return userStore, nil
	}

Any non-injector declarations are copied into the generated file.
There is no dependency on Wire at runtime:
all of the written code is just normal Go code.

As you can see, the output is very close to what a developer would write themselves.
This was a trivial example with just three components,
so writing the initializer by hand wouldn't be too painful,
but Wire saves a lot of manual toil for components and applications with
more complex dependency graphs.

## How can I get involved and learn more?

The [Wire README](https://github.com/google/wire/blob/master/README.md)
goes into more detail about how to use Wire and its more advanced features.
There's also a [tutorial](https://github.com/google/wire/tree/master/_tutorial)
that walks through using Wire in a simple application.

We appreciate any input you have about your experience with Wire!
[Wire's](https://github.com/google/wire) development is conducted on GitHub,
so you can [file an issue](https://github.com/google/wire/issues/new/choose)
to tell us what could be better.
For updates and discussion about the project,
join [the Go Cloud mailing list](https://groups.google.com/forum/#!forum/go-cloud).

Thank you for taking the time to learn about Go Cloud's Wire.
We’re excited to work with you to make Go the language of choice for developers
building portable cloud applications.
