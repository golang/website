---
title: Gopls on by default in the VS Code Go extension
date: 2021-02-01
by:
- Go tools team
tags:
- tools
- gopls
summary: Gopls, which provides IDE features for Go to many editors, is now used by default in VS Code Go.
---


We're happy to announce that the VS Code Go extension now enables the [gopls
language server](https://github.com/golang/tools/blob/master/gopls/README.md)
by default, to deliver more robust IDE features and better support for Go
modules.

{{image "gopls/features.gif" 635}}
_(`gopls` provides IDE features, such as intelligent autocompletion, signature help, refactoring, and workspace symbol search.)_

When [Go modules](using-go-modules) were
released two years ago, they completely changed the landscape of Go developer
tooling. Tools like `goimports` and `godef` previously depended on the fact
that code was stored in your `$GOPATH`. When the Go team began rewriting these
tools to work with modules, we immediately realized that we needed a more
systematic approach to bridge the gap.

As a result, we began working on a single Go
[language server](https://microsoft.github.io/language-server-protocol/),
`gopls`, which provides IDE features, such as autocompletion, formatting, and
diagnostics to any compatible editor frontend. This persistent and unified
server is a [fundamental
shift](https://www.youtube.com/watch?v=EFJfdWzBHwE&t=1s) from the earlier
collections of command-line tools.

In addition to working on `gopls`, we sought other ways of creating a stable
ecosystem of editor tooling. Last year, the Go team took responsibility for the
[Go extension for VS Code](https://blog.golang.org/vscode-go). As part of this
work, we smoothed the extension’s integration with the language server—automating
`gopls` updates, rearranging and clarifying `gopls` settings, improving the
troubleshooting workflow, and soliciting feedback through a survey. We’ve also
continued to foster a community of active users and contributors who have
helped us improve the stability, performance, and user experience of the Go
extension.

## Announcement

January 28 marked a major milestone in both the `gopls` and VS Code Go
journeys, as `gopls` is now enabled by default in the Go extension for VS Code.

In advance of this switch we spent a long time iterating on the design, feature
set, and user experience of `gopls`, focusing on improving performance and
stability. For more than a year, `gopls` has been the default in most plugins for
Vim, Emacs, and other editors. We’ve had 24 `gopls` releases, and we’re
incredibly grateful to our users for consistently providing feedback and
reporting issues on each and every one.

We’ve also dedicated time to smoothing the new user experience. We hope that VS
Code Go with `gopls` will be intuitive with clear error messages, but if you have
a question or need to adjust some configuration, you’ll be able to find answers
in our [updated documentation](https://github.com/golang/vscode-go/blob/master/README.md).
We have also recorded [a screencast](https://www.youtube.com/watch?v=1MXIGYrMk80)
to help you get started, as well as
[animations](https://github.com/golang/vscode-go/blob/master/docs/features.md)
to show off some hard-to-find features.

Gopls is the best way of working with Go code, especially with Go modules.
With the upcoming arrival of Go 1.16, in which modules are enabled by default,
VS Code Go users will have the best possible experience out-of-the-box.

Still, this switch does not mean that `gopls` is complete. We will continue
working on bug fixes, new features, and general stability. Our next area of
focus will be improving the user experience when [working with multiple
modules](https://github.com/golang/tools/blob/master/gopls/doc/workspace.md).
Feedback from our larger user base will help inform our next steps.

## So, what should you do?

If you use VS Code, you don’t need to do anything.
When you get the next VS Code Go update, `gopls` will be enabled automatically.

If you use another editor, you are likely using `gopls` already. If not, see
[the `gopls` user guide](https://github.com/golang/tools/blob/master/gopls/README.md)
to learn how to enable `gopls` in your preferred editor. The Language Server
Protocol ensures that `gopls` will continue to offer the same features to every
editor.

If `gopls` is not working for you, please see our [detailed troubleshooting
guide](https://github.com/golang/vscode-go/blob/master/docs/troubleshooting.md)
and file an issue. If you need to, you can always [disable `gopls` in VS
Code](https://github.com/golang/vscode-go/blob/master/docs/settings.md#gouselanguageserver).

## Thank you

To our existing users, thank you for bearing with us as we rewrote our caching
layer for the third time. To our new users, we look forward to hearing your
experience reports and feedback.

Finally, no discussion of Go tooling is complete without mentioning the
valuable contributions of the Go tools community. Thank you for the lengthy
discussions, detailed bug reports, integration tests, and most importantly,
thank you for the fantastic contributions. The most exciting `gopls` features
come from our passionate open-source contributors, and we are appreciative of
your hard work and dedication.

## Learn more

Watch [the screencast](https://www.youtube.com/watch?v=1MXIGYrMk80) for a
walk-through of how to get started with `gopls` and VS Code Go, and see the
[VS Code Go README](https://github.com/golang/vscode-go/blob/master/README.md)
for additional information.

If you’d like to read about `gopls` in more detail, see the
[`gopls` README](https://github.com/golang/tools/blob/master/gopls/README.md).
