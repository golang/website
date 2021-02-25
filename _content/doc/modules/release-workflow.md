<!--{
  "Title": "Module release and versioning workflow"
}-->

When you develop modules for use by other developers, you can follow a workflow
that helps ensure a reliable, consistent experience for developers using the
module. This topic describes the high-level steps in that workflow.

For an overview of module development, see [Developing and publishing
modules](developing).

**See also**

* If you're merely wanting to use external packages in your code, be sure to
  see [Managing dependencies](/doc/modules/managing-dependencies).
* With each new version, you signal the changes to your module with its
  version number. For more, see [Module version numbering](/doc/modules/version-numbers).

## Common workflow steps {#common-steps}

The following sequence illustrates release and versioning workflow steps for an
example new module. For more about each step, see the sections in this topic.

1.  **Begin a module** and organize its sources to make it easier for developers
    to use and for you to maintain.

    If you're brand new to developing modules, check out [Tutorial: Create a Go
    module](/doc/tutorial/create-module).

    In Go's decentralized module publishing system, how you organize your code
    matters. For more, see [Managing module source](/doc/modules/managing-source).

1.  Set up to **write local client code** that calls functions in the
    unpublished module.

    Before you publish a module, it's unavailable for the typical dependency
    management workflow using commands such as `go get`. A good way to test your
    module code at this stage is to try it while it is in a directory local to
    your calling code.

    See [Coding against an unpublished module](#unpublished) for more about
    local development.

1.  When the module's code is ready for other developers to try it out,
    **begin publishing v0 pre-releases** such as alphas and betas. See
    [Publishing pre-release versions](#pre-release) for more.

1.  **Release a v0** that's not guaranteed to be stable, but which users can try
    out. For more, see [Publishing the first (unstable) version](#first-unstable).

1.  After your v0 version is published, you can (and should!) continue to
    **release new versions** of it.

    These new versions might include bug fixes (patch releases), additions to
    the module's public API (minor releases), and even breaking changes. Because
    a v0 release makes no guarantees of stability or backward compatibility, you
    can make breaking changes in its versions.

    For more, see [Publishing bug fixes](#bug-fixes) and [Publishing
    non-breaking API changes](#non-breaking).

1.  When you're getting a stable version ready for release, you **publish
    pre-releases as alphas and betas**. For more, see [Publishing pre-release
    versions](#pre-release).

1.  Release a v1 as the **first stable release**.

    This is the first release that makes commitments about the module's
    stability. For more, see [Publishing the first stable
    version](#first-stable).

1.  In the v1 version, **continue to fix bugs** and, where necessary, make
    additions to the module's public API.

    For more, see [Publishing bug fixes](#bug-fixes) and [Publishing
    non-breaking API changes](#non-breaking).

1.  When it can't be avoided, publish breaking changes in a **new major version**.

    A major version update -- such as from v1.x.x to v2.x.x -- can be a very
    disruptive upgrade for your module's users. It should be a last resort. For
    more, see [Publishing breaking API changes](#breaking).

## Coding against an unpublished module {#unpublished}

When you begin developing a module or a new version of a module, you won't yet
have published it. Before you publish a module, you won't be able to use Go
commands to add the module as a dependency. Instead, at first, when writing
client code in a different module that calls functions in the unpublished
module, you'll need to reference a copy of the module on the local file system.

You can reference a module locally from the client module's go.mod file by using
the `replace` directive in the client module's go.mod file. For more
information, see in [Requiring module code in a local
directory](managing-dependencies#local_directory).

## Publishing pre-release versions {#pre-release}

You can publish pre-release versions to make a module available for others to
try it out and give you feedback. A pre-release version includes no guarantee of
stability.

Pre-release version numbers are appended with a pre-release identifier. For more
on version numbers, see [Module version numbering](/doc/modules/version-numbers).

Here are two examples:

```
v0.2.1-beta.1
v1.2.3-alpha
```

When making a pre-release available, keep in mind that developers using the
pre-release will need to explicitly specify it by version with the `go get`
command. That's because, by default, the `go` command prefers release versions
over pre-release versions when locating the module you're asking for. So
developers must get the pre-release by specifying it explicitly, as in the
following example:

```
go get example.com/theirmodule@v1.2.3-alpha
```

You publish a pre-release by tagging the module code in your repository,
specifying the pre-release identifier in the tag. For more, see [Publishing a
module](publishing).

## Publishing the first (unstable) version {#first-unstable}

As when you publish a pre-release version, you can publish release versions that
don't guarantee stability or backward compatibility, but give your users an
opportunity to try out the module and give you feedback.

Unstable releases are those whose version numbers are in the v0.x.x range. A v0
version makes no stability or backward compatibility guarantees. But it gives
you a way to get feedback and refine your API before making stability
commitments with v1 and later. For more see, [Module version
numbering](version-numbers).

As with other published versions, you can increment the minor and patch parts of
the v0 version number as you make changes toward releasing a stable v1 version.
For example, after releasing a v.0.0.0, you might release a v0.0.1 with the
first set of bug fixes.

Here's an example version number:

```
v0.1.3
```

You publish an unstable release by tagging the module code in your repository,
specifying a v0 version number in the tag. For more, see [Publishing a
module](publishing).

## Publishing the first stable version {#first-stable}

Your first stable release will have a v1.x.x version number. The first stable
release follows pre-release and v0 releases through which you got feedback,
fixed bugs, and stabilized the module for users.

With a v1 release, you're making the following commitments to developers using
your module:

* They can upgrade to the major version's subsequent minor and patch releases
  without breaking their own code.
* You won't be making further changes to the module's public API -- including
  its function and method signatures -- that break backward compatibility.
* You won't be removing any exported types, which would break backward
  compatibility.
* Future changes to your API (such as adding a new field to a struct) will be
  backward compatible and will be included in a new minor release.
* Bug fixes (such as a security fix) will be included in a patch release or as
  part of a minor release.

**Note:** While your first major version might be a v0 release, a v0 version
does not signal stability or backward compatibility guarantees. As a result,
when you increment from v0 to v1, you needn't be mindful of breaking backward
compatibility because the v0 release was not considered stable.

For more about version numbers, see [Module version numbering](/doc/modules/version-numbers).

Here's an example of a stable version number:

```
v1.0.0
```

You publish a first stable release by tagging the module code in your
repository, specifying a v1 version number in the tag. For more, see [Publishing
a module](publishing).

## Publishing bug fixes {#bug-fixes}

You can publish a release in which the changes are limited to bug fixes. This is
known as a patch release.

A _patch release_ includes only minor changes. In particular, it includes no
changes to the module's public API. Developers of consuming code can upgrade to
this version safely and without needing to change their code.

**Note:** Your patch release should try not to upgrade any of that module's own
transitive dependencies by more than a patch release. Otherwise, someone
upgrading to the patch of your module could wind up accidentally pulling in a
more invasive change to a transitive dependency that they use.

A patch release increments the patch part of the module's version number. For
more see, [Module version numbering](/doc/modules/version-numbers).

In the following example, v1.0.1 is a patch release.

Old version: `v1.0.0`

New version: `v1.0.1`

You publish a patch release by tagging the module code in your repository,
incrementing the patch version number in the tag. For more, see [Publishing a
module](publishing).

## Publishing non-breaking API changes {#non-breaking}

You can make non-breaking changes to your module's public API and publish those
changes in a _minor_ version release.

This version changes the API, but not in a way that breaks calling code. This
might include changes to a module’s own dependencies or the addition of new
functions, methods, struct fields, or types. Even with the changes it includes,
this kind of release guarantees backward compatibility and stability for
existing code that calls the module's functions.

A minor release increments the minor part of the module's version number. For
more, see [Module version numbering](/doc/modules/version-numbers).

In the following example, v1.1.0 is a minor release.

Old version: `v1.0.1`

New version: `v1.1.0`

You publish a minor release by tagging the module code in your repository,
incrementing the minor version number in the tag. For more, see [Publishing a
module](publishing).

## Publishing breaking API changes {#breaking}

You can publish a version that breaks backward compatibility by publishing a
_major_ version release.

A major version release doesn't guarantee backward compatibility, typically
because it includes changes to the module's public API that would break code
using the module's previous versions.

Given the disruptive effect a major version upgrade can have on code relying on
the module, you should avoid a major version update if you can. For more about
major version updates, see [Developing a major version update](/doc/modules/major-version).
For strategies to avoid making breaking changes, see the blog post [Keeping your
modules compatible](https://blog.golang.org/module-compatibility).

Where publishing other kinds of versions requires essentially tagging the module
code with the version number, publishing a major version update requires more
steps.

1.  Before beginning development of the new major version, in your repository
    create a place for the new version's source.

    One way to do this is to create a new branch in your repository that is
    specifically for the new major version and its subsequent minor and patch
    versions. For more, see [Managing module source](/doc/modules/managing-source).

1.  In the module's go.mod file, revise the module path to append the new major
    version number, as in the following example:

    ```
    example.com/mymodule/v2
    ```

    Given that the module path is the module's identifier, this change
    effectively creates a new module. It also changes the package path, ensuring
    that developers won't unintentionally import a version that breaks their
    code. Instead, those wanting to upgrade will explicitly replace occurrences
    of the old path with the new one.

1.  In your code, change any package paths where you're importing packages in
    the module you're updating, including packages in the module you're updating.
    You need to do this because you changed your module path.

1.  As with any new release, you should publish pre-release versions to get
    feedback and bug reports before publishing an official release.

1.  Publish the new major version by tagging the module code in your repository,
    incrementing the major version number in the tag -- such as from v1.5.2 to
    v2.0.0.

    For more, see [Publishing a module](/doc/modules/publishing).
