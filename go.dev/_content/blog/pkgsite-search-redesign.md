---
title: "A new search experience on pkg.go.dev"
date: 2021-11-09
by:
- Julie Qiu
summary: Package search on pkg.go.dev has been updated, and you can now search for symbols!
---

We are excited to launch a new search experience on [pkg.go.dev](https://pkg.go.dev/).

These changes were motivated by
[feedback we've received](/issue/47321) about the search
page, and we hope you enjoy them. This blog post provides an overview of
what you can expect to see on the site.

## Grouping related package search results

Search results for packages in the same module are now grouped together. The
most relevant package for the search request is highlighted. This change was
made to reduce noise when several packages in the same module may be relevant
to a search. For example, searching for "markdown" shows a row listing "Other
packages in module" for several of the results.

{{image "pkgsite-search-redesign/markdown.png" 850}}

Results for different major versions of the same module are also now grouped
together. The highest major version containing a tagged release is highlighted.
For example, searching for "github" shows the v39 module, with older versions
listed as "Other major versions".

{{image "pkgsite-search-redesign/github.png" 850}}

Lastly, we have reorganized information related to imports, versions, and
licenses. We also added links to these tabs directly from the search results
page.

## Introducing symbol search

Over the past year, we have introduced more information about symbols on
pkg.go.dev and worked on improving the way that information is presented. We
launched the ability to view the API history of any package. We also label
symbols that are deprecated in the documentation index and hide them by
default in the package documentation.

With this search update, pkg.go.dev now also supports searching for symbols in
Go packages. When a user types a symbol into the search bar, they will be
brought to a new search tab for symbol search results. There are a few
different ways in which pkg.go.dev identifies that users are searching for a
symbol. We've added examples to the pkg.go.dev homepage, and detailed
instructions to the [search help page](https://pkg.go.dev/search-help).

{{image "pkgsite-search-redesign/httpclient.png" 850}}

## Feedback

We’re excited to share this new search experience with you, and we would love
to hear your feedback!

As always, please use the “Report an Issue” button at the bottom of every page
on the site to share your input.

If you’re interested in contributing to this project, pkg.go.dev is open
source! Check out the
[contribution guidelines](https://go.googlesource.com/pkgsite/+/refs/heads/master/CONTRIBUTING.md)
to find out more.
