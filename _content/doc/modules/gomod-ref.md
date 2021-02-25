<!--{
  "Title": "go.mod file reference"
}-->

Each Go module is defined by a go.mod file that describes the module's
properties, including its dependencies on other modules and on versions of Go.

These properties include:

* The current module's **module path**. This should be a location from which
the module can be down loaded by Go tools, such as the module code's
repository location. This serves as a unique identifier, when combined
with the module's version number. It is also the prefix of the package path for
all packages in the module. For more about how Go locates the module, see the
<a href="/ref/mod#vcs-find">Go Modules Reference</a>.
* The minimum **version of Go** required by the current module.
* A list of minimum versions of other **modules required** by the current module.
* Instructions, optionally, to **replace** a required module with another
  module version or a local directory, or to **exclude** a specific version of
  a required module.

Go generates a go.mod file when you run the [`go mod init`
command](/ref/mod#go-mod-init). The following example creates a go.mod file,
setting the module's module path to example.com/mymodule:

```
$ go mod init example.com/mymodule
```

Use `go` commands to manage dependencies. The commands ensure that the
requirements described in your go.mod file remain consistent and the content of
your go.mod file is valid. These commands include the [`go get`](/ref/mod#go-get)
and [`go mod tidy`](/ref/mod#go-mod-tidy) and [`go mod edit`](/ref/mod#go-mod-edit)
commands.

For reference on `go` commands, see [Command go](/cmd/go/).
You can get help from the command line by typing `go help` _command-name_, as
with `go help mod tidy`.

**See also**

* Go tools make changes to your go.mod file as you use them to manage
  dependencies. For more, see [Managing dependencies](/doc/modules/managing-dependencies).
* For more details and constraints related to go.mod files, see the [Go modules
  reference](/ref/mod#go-mod-file).

## Example {#example}

A go.mod file includes directives shown in the following example. These are
described in this topic.

```
module example.com/mymodule

go 1.14

require (
    example.com/othermodule v1.2.3
    example.com/thismodule v1.2.3
    example.com/thatmodule v1.2.3
)

replace example.com/thatmodule => ../thatmodule
exclude example.com/thismodule v1.3.0
```

## module {#module}

Declares the module's module path, which is the module's unique identifier
(when combined with the module version number). This becomes the import prefix
for all packages the module contains.

### Syntax {#module-syntax}

<pre>module <var>module-path</var></pre>

<dl>
    <dt>module-path</dt>
    <dd>The module's module path, usually the repository location from which
      the module can be downloaded by Go tools. For module versions v2 and
      later, this value must end with the major version number, such as
      <code>/v2</code>.</dd>
</dl>

### Examples {#module-examples}

The following examples substitute `example.com` for a repository domain from
which the module could be downloaded.

* Module declaration for a v0 or v1 module:
  ```
  module example.com/mymodule
  ```
* Module path for a v2 module:
  ```
  module example.com/mymodule/v2
  ```

### Notes {#module-notes}

The module path should be a path from which Go tools can download the module
source. In practice, this is typically the module source's repository domain
and path to the module code within the repository. The <code>go</code> command
relies on this form when downloading module versions to resolve dependencies
on the module user's behalf.

Even if you're not at first intending to make your module available for use
from other code, using its repository path is a best practice that will help
you avoid having to rename the module if you publish it later.

If at first you don't know the module's eventual repository location, consider
temporarily using a safe substitute, such as the name of a domain you own or
`example.com`, along with a path following from the module's name or source
directory.

For example, if you're developing in a `stringtools` directory, your temporary
module path might be `example.com/stringtools`, as in the following example:

```
go mod init example.com/stringtools
```

## go {#go}

Specifies the minimum Go version required by the module.

### Syntax {#go-syntax}

<pre>go <var>minimum-go-version</var></pre>

<dl>
    <dt>minimum-go-version</dt>
    <dd>The minimum version of Go required to compile packages in this module.</dd>
</dl>

### Examples {#go-examples}

* Module must run on Go version 1.14 or later:
  ```
  go 1.14
  ```

### Notes {#go-notes}

The version number also controls some module-related features in the `go`
command. For example, if the vendor directory is present, that directory will be
used automatically if the version number is 1.14 or higher.

For more about version numbers, see [Module version numbering](/doc/modules/version-numbers).

## require {#require}

Declares a module as dependency required by the current module, specifying the
minimum version of the module required.

### Syntax {#require-syntax}

<pre>require <var>module-path</var> <var>module-version</var></pre>

<dl>
    <dt>module-path</dt>
    <dd>The module's module path, usually a concatenation of the module source's
      repository domain and the module name. For module versions v2 and later,
      this value must end with the major version number, such as <code>/v2</code>.</dd>
    <dt>module-version</dt>
    <dd>The module's version. This can be either a release version number, such
      as v1.2.3, or a Go-generated pseudo-version number, such as
      v0.0.0-20200921210052-fa0125251cc4.</dd>
</dl>

### Examples {#require-examples}

* Requiring a released version v1.2.3:
    ```
    require example.com/othermodule v1.2.3
    ```
* Requiring a version not yet tagged in its repository by using a pseudo-version
  number generated by Go tools:
    ```
    require example.com/othermodule v0.0.0-20200921210052-fa0125251cc4
    ```

### Notes {#require-notes}

When you run a `go` command such as `go get`, Go inserts `require` directives
for each module containing imported packages. When a module isn't yet tagged in
its repository, Go assigns a pseudo-version number it generates when you run the
command.

You can have Go require a module from a location other than its repository by
using the [`replace` directive](#replace).

For more about version numbers, see [Module version numbering](/doc/modules/version-numbers).

For more about managing dependencies, see the following:

* [Adding a dependency](/doc/modules/managing-dependencies#adding_dependency)
* [Getting a specific dependency version](/doc/modules/managing-dependencies#getting_version)
* [Discovering available updates](/doc/modules/managing-dependencies#discovering_updates)
* [Upgrading or downgrading a dependency](/doc/modules/managing-dependencies#upgrading)
* [Synchronizing your code's dependencies](/doc/modules/managing-dependencies#synchronizing)

## replace {#replace}

Replaces the content of a module at a specific version (or all versions) with
another module version or with a local directory. Go tools will use the
replacement path when resolving the dependency.

### Syntax {#replace-syntax}

<pre>replace <var>module-path</var> <var>[module-version]</var> => <var>replacement-path</var> <var>[replacement-version]</var></pre>

<dl>
    <dt>module-path</dt>
    <dd>The module path of the module to replace.</dd>
    <dt>module-version</dt>
    <dd>Optional. A specific version to replace. If this version number is
      omitted, all versions of the module are replaced with the content on the
      right side of the arrow.</dd>
    <dt>replacement-path</dt>
    <dd>The path at which Go should look for the required module. This can be a
      module path or a path to a directory on the file system local to the
      replacement module. If this is a module path, you must specify a
      <em>replacement-version</em> value. If this is a local path, you may not use a
      <em>replacement-version</em> value.</dd>
    <dt>replacement-version</dt>
    <dd>The version of the replacement module. The replacement version may only
      be specified if <em>replacement-path</em> is a module path (not a local directory).</dd>
</dl>

### Examples {#replace-examples}

* Replacing with a fork of the module repository

  In the following example, any version of example.com/othermodule is replaced
  with the specified fork of its code.

  ```
  require example.com/othermodule v1.2.3

  replace example.com/othermodule => example.com/myfork/othermodule v1.2.3-fixed
  ```

  When you replace one module path with another, do not change import statements
  for packages in the module you're replacing.

  For more on using a forked copy of module code, see [Requiring external module
  code from your own repository fork](/doc/modules/managing-dependencies#external_fork).

* Replacing with a different version number

  The following example specifies that version v1.2.3 should be used instead of
  any other version of the module.

  ```
  require example.com/othermodule v1.2.2

  replace example.com/othermodule => example.com/othermodule v1.2.3
  ```

  The following example replaces module version v1.2.5 with version v1.2.3 of
  the same module.

  ```
  replace example.com/othermodule v1.2.5 => example.com/othermodule v1.2.3
  ```

* Replacing with local code

  The following example specifies that a local directory should be used as a
  replacement for all versions of the module.

  ```
  require example.com/othermodule v1.2.3

  replace example.com/othermodule => ../othermodule
  ```

  The following example specifies that a local directory should be used as a
  replacement for v1.2.5 only.

  ```
  require example.com/othermodule v1.2.5

  replace example.com/othermodule v1.2.5 => ../othermodule
  ```

  For more on using a local copy of module code, see [Requiring module code in a
  local directory](/doc/modules/managing-dependencies#local_directory).

### Notes {#replace-notes}

Use the `replace` directive to temporarily substitute a module path value with
another value when you want Go to use the other path to find the module's
source. This has the effect of redirecting Go's search for the module to the
replacement's location. You needn't change package import paths to use the
replacement path.

Use the `exclude` and `replace` directives to control build-time dependency
resolution when building the current module. These directives are ignored in
modules that depend on the current module.

The `replace` directive can be useful in situations such as the following:

* You're developing a new module whose code is not yet in the repository. You
  want to test with clients using a local version.
* You've identified an issue with a dependency, have cloned the dependency's
  repository, and you're testing a fix with the local repository.

For more on replacing a required module, including using Go tools to make the
change, see:

* [Requiring external module code from your own repository
fork](/doc/modules/managing-dependencies#external_fork)
* [Requiring module code in a local
directory](/doc/modules/managing-dependencies#local_directory)

For more about version numbers, see [Module version
numbering](/doc/modules/version-numbers).

## exclude {#exclude}

Specifies a module or module version to exclude from the current module's
dependency graph.

### Syntax {#exclude-syntax}

<pre>exclude <var>module-path</var> <var>module-version</var></pre>

<dl>
    <dt>module-path</dt>
    <dd>The module path of the module to exclude.</dd>
    <dt>module-version</dt>
    <dd>The specific version to exclude.</dd>
</dl>

### Example {#exclude-example}

* Exclude example.com/theirmodule version v1.3.0

  ```
  exclude example.com/theirmodule v1.3.0
  ```

### Notes {#exclude-notes}

Use the `exclude` directive to exclude a specific version of a module that is
indirectly required but can't be loaded for some reason. For example, you might
use it to exclude a version of a module that has an invalid checksum.

Use the `exclude` and `replace` directives to control build-time dependency
resolution when building the current module (the main module you're building).
These directives are ignored in modules that depend on the current module.

You can use the [`go mod edit`](/ref/mod#go-mod-edit) command
to exclude a module, as in the following example.

```
go mod edit -exclude=example.com/theirmodule@v1.3.0
```

For more about version numbers, see [Module version numbering](/doc/modules/version-numbers).
