---
title: Pkg.go.dev has a new look!
date: 2020-11-10T12:00:00Z
by:
- Julie Qiu
summary: Announcing a new user experience on pkg.go.dev.
---


Since launching [pkg.go.dev](https://pkg.go.dev), we’ve received a lot of great
feedback on design and usability.
In particular, it was clear that the way information was organized confused
users when navigating the site.

Today we’re excited to share a redesigned pkg.go.dev,
which we hope will be clearer and more helpful.
This blog post presents the highlights. For details,
see [Go issue 41585](/issue/41585).

## Consistent landing page for all paths

The main change is that the pkg.go.dev/\<path> page has been reorganized
around the idea of a path.
A path represents a directory in a particular version of a module.
Now, regardless of what’s in that directory,
every path page will have the same layout,
with the goal of making the experience consistently useful and predictable.

<figure class="image">
  <img src="pkgsite-redesign/path.png" width="800" alt="Landing page for cloud.google.com/go/storage" style="border: 1px solid black;">
  <figcaption>
    Fig 1. Landing page for
    <a href="https://pkg.go.dev/cloud.google.com/go/storage">https://pkg.go.dev/cloud.google.com/go/storage</a>.
  </figcaption>
</figure>

The path page will display the README at that path if there is one.
Previously, the overview tab only showed the README if present at the module root.
This is one of many changes we’re making to place the most important information up front.

## Documentation navigation

The documentation section now displays an index along with a sidenav.
This gives the ability to see the full package API,
while having context as they are navigating the documentation section.
There is also a new Jump To input box in the left sidenav,
to search for identifiers.


<figure class="image">
  <img src="pkgsite-redesign/nav.png" width="800" alt="Jump To feature navigating net/http" style="border: 1px solid black;">
  <figcaption>
    Fig 2. Jump To feature on
    <a href="https://pkg.go.dev/net/http">https://pkg.go.dev/net/http</a>.
  </figcaption>
</figure>

See [Go issue 41587](/issue/41587) for details on changes in the documentation section.

## Metadata on main page

The top bar on each page now shows additional metadata,
such as each package’s “imports” and “imported by” counts.
Banners also show information about the latest minor and major versions of a module.
See [Go issue 41588](/issue/41588) for details.

<figure class="image">
  <img src="pkgsite-redesign/meta.png" width="800" alt="Header metadata for github.com/russross/blackfriday" style="border: 1px solid black;">
  <figcaption>
    Fig 3. Header metadata for
    <a href="https://pkg.go.dev/github.com/russross/blackfriday">https://pkg.go.dev/github.com/russross/blackfriday</a>.
  </figcaption>
</figure>

## Video Walkthrough

Last week at [Google Open Source Live](https://opensourcelive.withgoogle.com/events/go),
we presented a walkthrough of the new site experience in our talk,
[Level Up: Go Package Discovery and Editor Tooling](https://www.youtube.com/watch?v=n7ayE29b7QA&feature=emb_logo).

{{video "https://www.youtube.com/embed/n7ayE29b7QA" 650 400}}

## Feedback

We’re excited to share this updated design with you.
As always, please let us know what you think via the “Share Feedback”
and “Report an Issue” links at the bottom of every page of the site.

And if you’re interested in contributing to this project, pkg.go.dev is open source! Check out the
[contribution guidelines](https://go.googlesource.com/pkgsite/+/refs/heads/master/CONTRIBUTING.md)
to find out more.

