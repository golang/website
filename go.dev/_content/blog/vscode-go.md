---
title: The VS Code Go extension joins the Go project
date: 2020-06-09
by:
- The Go team
summary: Announcement of VS Code Go’s move to the Go project.
---


When the Go project began, “an overarching goal was that Go do more to help the
working programmer by enabling tooling, automating mundane tasks such as code
formatting, and removing obstacles to working on large code bases”
([Go FAQ](/doc/faq#What_is_the_purpose_of_the_project)).
Today, more than a decade later, we continue to be guided by that same goal,
especially as it pertains to the programmer’s most critical tool: their editor.

Throughout the past decade, Go developers have relied on a variety of editors
and dozens of independently authored tools and plugins. Much of Go’s early
success can be attributed to the fantastic development tools created by the Go
community. The
[VS Code extension for Go](https://github.com/microsoft/vscode-go), built using
many of these tools, is now used by 41 percent of Go developers
([Go developer survey](https://blog.golang.org/survey2019-results)).

As the VS Code Go extension grows in popularity and as
[the ecosystem expands](https://www.youtube.com/watch?v=EFJfdWzBHwE), it
requires
[more maintenance and support](https://twitter.com/ramyanexus/status/1154470078978486272).
Over the past few years, the Go team has collaborated with the VS Code team to
help the Go extension maintainers. The Go team also began a new initiative to
improve the tools powering all Go editor extensions, with a focus on supporting
the
[Language Server Protocol](https://microsoft.github.io/language-server-protocol/)
with [`gopls`](/s/gopls) and
[the Debug Adapter Protocol with Delve](https://github.com/go-delve/delve/issues/1515).

Through this collaborative work between the VS Code and Go teams, we realized
that the Go team is uniquely positioned to evolve the Go development experience
alongside the Go language.

As a result, we’re happy to announce the next phase in the Go team’s
partnership with the VS Code team: **The VS Code extension for Go is officially
joining the Go project**. With this come two critical changes:

1. The publisher of the plugin is shifting from "Microsoft" to "Go Team at Google".
2. The project’s repository is moving to join the rest of the Go project at [https://github.com/golang/vscode-go](https://github.com/golang/vscode-go).

We cannot overstate our gratitude to those who have helped
build and maintain this beloved extension. We know that innovative ideas and
features come from you, our users. The Go team’s primary aim as owners of the
extension is to reduce the burden of maintenance work on the Go community.
We’ll make sure the builds stay green, the issues get triaged, and the docs get
updated. Go team members will keep contributors abreast of relevant language
changes, and we’ll smooth the rough edges between the extension’s different
dependencies.

Please continue to share your thoughts with us by filing
[issues](https://github.com/golang/vscode-go/issues) and making
[contributions](https://github.com/golang/vscode-go/blob/master/docs/contributing.md)
to the project. The process for contributing will now be the same as for the
[rest of the Go project](/doc/contribute.html). Go team
members will offer general help in the #vscode channel on
[Gophers Slack](https://invite.slack.golangbridge.org/), and we’ve also created
a #vscode-dev channel to discuss issues and brainstorm ideas with contributors.

We’re excited about this new step forward, and we hope you are too.
By maintaining a major Go editor extension, as well as the Go tooling and
language, the Go team will be able to provide all Go users, regardless of their
editor, a more cohesive and refined development experience.

As always, our goal remains the same: Every user should have an excellent
experience writing Go code.

*See the accompanying post from the [Visual Studio Code team](https://aka.ms/go-blog-vscode-202006).*
