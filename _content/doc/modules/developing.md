<!--{
  "Title": "Developing and publishing modules"
}-->

You can collect related packages into modules, then publish the modules for
other developers to use. This topic gives an overview of developing and
publishing modules.

To support developing, publishing, and using modules, you use:

*   A **workflow** through which you develop and publish modules, revising them
	with new versions over time. See [Workflow for developing and publishing
	modules](#workflow).
*	**Design practices** that help a module's users understand it and upgrade
	to new versions in a stable way. See [Design and development](#design).
*   A **decentralized system for publishing** modules and retrieving their code.
	You make your module available for other developers to use from your own
	repository and publish with a version number. See [Decentralized
	publishing](#decentralized).
*   A **package search engine** and documentation browser (pkg.go.dev) at which
	developers can find your module. See [Package discovery](#discovery).
*   A module **version numbering convention** to communicate expectations of
	stability and backward compatibility to developers using your module. See
	[Versioning](#versioning).
*   **Go tools** that make it easier for other developers to manage
	dependencies, including getting your module's source, upgrading, and so on.
	See [Managing dependencies](/doc/modules/managing-dependencies).

**See also**

*   If you're interested simply in using packages developed by others, this
	isn't the topic for you. Instead, see [Managing
	dependencies](managing-dependencies).
*   For a tutorial that includes a few module development basics, see
	[Tutorial: Create a Go module](/doc/tutorial/create-module).

## Workflow for developing and publishing modules {#workflow}

When you want to publish your modules for others, you adopt a few conventions to
make using those modules easier.

The following high-level steps are described in more detail in [Module release
and versioning workflow](release-workflow).

1. Design and code the packages that the module will include.
1. Commit code to your repository using conventions that ensure it's available
	to others via Go tools.
1. Publish the module to make it discoverable by developers.
1. Over time, revise the module with versions that use a version numbering
	convention that signals each version's stability and backward compatibility.

## Design and development {#design}

Your module will be easier for developers to find and use if the functions and
packages in it form a coherent whole. When you're designing a module's public
API, try to keep its functionality focused and discrete.

Also, designing and developing your module with backward compatibility in mind
helps its users upgrade while minimizing churn to their own code. You can use
certain techniques in code to avoid releasing a version that breaks backward
compatibility. For more about those techniques, see [Keeping your modules
compatible](https://blog.golang.org/module-compatibility) on the Go blog.

Before you publish a module, you can reference it on the local file system using
the replace directive. This makes it easier to write client code that calls
functions in the module while the module is still in development. For more
information, see "Coding against an unpublished module" in [Module release and
versioning workflow](release-workflow#unpublished).

## Decentralized publishing {#decentralized}

In Go, you publish your module by tagging its code in your repository to make it
available for other developers to use. You don't need to push your module to a
centralized service because Go tools can download your module directly from your
repository (located using the module's path, which is a URL with the scheme
omitted) or from a proxy server.

After importing your package in their code, developers use Go tools (including
the `go get` command) to download your module's code to compile with. To support
this model, you follow conventions and best practices that make it possible for
Go tools (on behalf of another developer) to retrieve your module's source from
your repository. For example, Go tools use the module's module path you specify,
along with the module version number you use to tag the module for release, to
locate and download the module for its users.

For more about source and publishing conventions and best practices, see
[Managing module source](/doc/modules/managing-source).

For step-by-step instructions on publishing a module, see [Publishing a
module](publishing).

## Package discovery {#discovery}

After you've published your module and someone has fetched it with Go tools, it
will become visible on the Go package discovery site at
[pkg.go.dev](https://pkg.go.dev/). There, developers can search the site to find
it and read its documentation.

To begin using the module, a developer imports packages from the module, then
runs the `go get` command to download its source code to compile with.

For more about how developers find and use modules, see [Managing
dependencies](managing-dependencies).

## Versioning {#versioning}

As you revise and improve your module over time, you assign version numbers
(based on the semantic versioning model) designed to signal each version's
stability and backward compatibility. This helps developers using your module
determine when the module is stable and whether an upgrade may include
significant changes in behavior. You indicate a module's version number by
tagging the module's source in the repository with the number.

For more on developing major version updates, see [Developing a major version
update](major-version).

For more about how you use the semantic versioning model for Go modules, see
[Module version numbering](/doc/modules/version-numbers).
