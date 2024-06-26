---
title: Generating code
date: 2014-12-22
by:
- Rob Pike
tags:
- programming
- technical
summary: How to use go generate.
---


A property of universal computation—Turing completeness—is that a computer program can write a computer program.
This is a powerful idea that is not appreciated as often as it might be, even though it happens frequently.
It's a big part of the definition of a compiler, for instance.
It's also how the `go` `test` command works: it scans the packages to be tested,
writes out a Go program containing a test harness customized for the package,
and then compiles and runs it.
Modern computers are so fast this expensive-sounding sequence can complete in a fraction of a second.

There are lots of other examples of programs that write programs.
[Yacc](https://godoc.org/golang.org/x/tools/cmd/goyacc), for instance, reads in a description of a grammar and writes out a program to parse that grammar.
The protocol buffer "compiler" reads an interface description and emits structure definitions,
methods, and other support code.
Configuration tools of all sorts work like this too, examining metadata or the environment
and emitting scaffolding customized to the local state.

Programs that write programs are therefore important elements in software engineering,
but programs like Yacc that produce source code need to be integrated into the build
process so their output can be compiled.
When an external build tool like Make is being used, this is usually easy to do.
But in Go, whose go tool gets all necessary build information from the Go source, there is a problem.
There is simply no mechanism to run Yacc from the go tool alone.

Until now, that is.

The [latest Go release](/blog/go1.4), 1.4,
includes a new command that makes it easier to run such tools.
It's called `go` `generate`, and it works by scanning for special comments in Go source code
that identify general commands to run.
It's important to understand that `go` `generate` is not part of `go` `build`.
It contains no dependency analysis and must be run explicitly before running `go` `build`.
It is intended to be used by the author of the Go package, not its clients.

The `go` `generate` command is easy to use.
As a warmup, here's how to use it to generate a Yacc grammar.

First, install Go's Yacc tool:

	go get golang.org/x/tools/cmd/goyacc

Say you have a Yacc input file called `gopher.y` that defines a grammar for your new language.
To produce the Go source file implementing the grammar,
you would normally invoke the command like this:

	goyacc -o gopher.go -p parser gopher.y

The `-o` option names the output file while `-p` specifies the package name.

To have `go` `generate` drive the process, in any one of the regular (non-generated) `.go` files
in the same directory, add this comment anywhere in the file:

	//go:generate goyacc -o gopher.go -p parser gopher.y

This text is just the command above prefixed by a special comment recognized by `go` `generate`.
The comment must start at the beginning of the line and have no spaces between the `//` and the `go:generate`.
After that marker, the rest of the line specifies a command for `go` `generate` to run.

Now run it. Change to the source directory and run `go` `generate`, then `go` `build` and so on:

	$ cd $GOPATH/myrepo/gopher
	$ go generate
	$ go build
	$ go test

That's it.
Assuming there are no errors, the `go` `generate` command will invoke `yacc` to create `gopher.go`,
at which point the directory holds the full set of Go source files, so we can build, test, and work normally.
Every time `gopher.y` is modified, just rerun `go` `generate` to regenerate the parser.

For more details about how `go` `generate` works, including options, environment variables,
and so on, see the [design document](/s/go1.4-generate).

Go generate does nothing that couldn't be done with Make or some other build mechanism,
but it comes with the `go` tool—no extra installation required—and fits nicely into the Go ecosystem.
Just keep in mind that it is for package authors, not clients,
if only for the reason that the program it invokes might not be available on the target machine.
Also, if the containing package is intended for import by `go` `get`,
once the file is generated (and tested!) it must be checked into the
source code repository to be available to clients.

Now that we have it, let's use it for something new.
As a very different example of how `go` `generate` can help, there is a new program available in the
`golang.org/x/tools` repository called `stringer`.
It automatically writes string methods for sets of integer constants.
It's not part of the released distribution, but it's easy to install:

	$ go get golang.org/x/tools/cmd/stringer

Here's an example from the documentation for
[`stringer`](https://godoc.org/golang.org/x/tools/cmd/stringer).
Imagine we have some code that contains a set of integer constants defining different types of pills:

	package painkiller

	type Pill int

	const (
		Placebo Pill = iota
		Aspirin
		Ibuprofen
		Paracetamol
		Acetaminophen = Paracetamol
	)

For debugging, we'd like these constants to pretty-print themselves, which means we want a method with signature,

	func (p Pill) String() string

It's easy to write one by hand, perhaps like this:

	func (p Pill) String() string {
		switch p {
		case Placebo:
			return "Placebo"
		case Aspirin:
			return "Aspirin"
		case Ibuprofen:
			return "Ibuprofen"
		case Paracetamol: // == Acetaminophen
			return "Paracetamol"
		}
		return fmt.Sprintf("Pill(%d)", p)
	}

There are other ways to write this function, of course.
We could use a slice of strings indexed by Pill, or a map, or some other technique.
Whatever we do, we need to maintain it if we change the set of pills, and we need to make sure it's correct.
(The two names for paracetamol make this trickier than it might otherwise be.)
Plus the very question of which approach to take depends on the types and values:
signed or unsigned, dense or sparse, zero-based or not, and so on.

The `stringer` program takes care of all these details.
Although it can be run in isolation, it is intended to be driven by `go` `generate`.
To use it, add a generate comment to the source, perhaps near the type definition:

	//go:generate stringer -type=Pill

This rule specifies that `go` `generate` should run the `stringer` tool to generate a `String` method for type `Pill`.
The output is automatically written to `pill_string.go` (a default we could override with the
`-output` flag).

Let's run it:

{{raw `
	$ go generate
	$ cat pill_string.go
	// Code generated by stringer -type Pill pill.go; DO NOT EDIT.

	package painkiller

	import "fmt"

	const _Pill_name = "PlaceboAspirinIbuprofenParacetamol"

	var _Pill_index = [...]uint8{0, 7, 14, 23, 34}

	func (i Pill) String() string {
		if i < 0 || i+1 >= Pill(len(_Pill_index)) {
			return fmt.Sprintf("Pill(%d)", i)
		}
		return _Pill_name[_Pill_index[i]:_Pill_index[i+1]]
	}
	$
`}}

Every time we change the definition of `Pill` or the constants, all we need to do is run

	$ go generate

to update the `String` method.
And of course if we've got multiple types set up this way in the same package,
that single command will update all their `String` methods with a single command.

There's no question the generated method is ugly.
That's OK, though, because humans don't need to work on it; machine-generated code is often ugly.
It's working hard to be efficient.
All the names are smashed together into a single string,
which saves memory (only one string header for all the names, even if there are zillions of them).
Then an array, `_Pill_index`, maps from value to name by a simple, efficient technique.
Note too that `_Pill_index` is an array (not a slice; one more header eliminated) of `uint8`,
the smallest integer sufficient to span the space of values.
If there were more values, or there were negatives ones,
the generated type of `_Pill_index` might change to `uint16` or `int8`: whatever works best.

The approach used by the methods printed by `stringer` varies according to the properties of the constant set.
For instance, if the constants are sparse, it might use a map.
Here's a trivial example based on a constant set representing powers of two:

	const _Power_name = "p0p1p2p3p4p5..."

	var _Power_map = map[Power]string{
		1:    _Power_name[0:2],
		2:    _Power_name[2:4],
		4:    _Power_name[4:6],
		8:    _Power_name[6:8],
		16:   _Power_name[8:10],
		32:   _Power_name[10:12],
		...,
	}

	func (i Power) String() string {
		if str, ok := _Power_map[i]; ok {
			return str
		}
		return fmt.Sprintf("Power(%d)", i)
	}

In short, generating the method automatically allows us to do a better job than we would expect a human to do.

There are lots of other uses of `go` `generate` already installed in the Go tree.
Examples include generating Unicode tables in the `unicode` package,
creating efficient methods for encoding and decoding arrays in `encoding/gob`,
producing time zone data in the `time` package, and so on.

Please use `go` `generate` creatively.
It's there to encourage experimentation.

And even if you don't, use the new `stringer` tool to write your `String` methods for your integer constants.
Let the machine do the work.
