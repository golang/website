---
title: 'About'
date: 2019-06-25T17:51:23-04:00
draft: false
---

[Go.dev](https://go.dev) is a companion website to [golang.org](https://golang.org). [golang.org](https://golang.org) is the home of the open source project and distribtion. [Go.dev](https://go.dev) is the hub for Go users. Go.dev is a portal to the entire Go ecosystem providing centalized and curated resources. 

{{% gopher gopher=pink align=right %}}
Go.dev provides:

1. Centralized information for Go packages and modules published on index.golang.org. 
2. Essential learning resources
3. Critical use cases & case studies

Go.dev is currently in [MVP](https://en.wikipedia.org/wiki/Minimum_viable_product) status. We're proud of what we've built and are excited to share it with the community. We hope you find value and joy in using go.dev. As a result, go.dev only has a small portion of features we intend to build. We are actively seeking feedback. If you have any ideas, suggestions or issues, please let us know.

## Sharing feedback / Reporting an issue

On the footer of every page there are two links, "Share Feedback" and "Report an issue". These links will enable you to capture a screenshot of the page you are on and annotate that screenshot, then send this directly to the go.dev team. 

Alternatively, you can send your bugs, ideas, feature requests and questions to [go-discovery-feedback@google.com](mailto:go-discovery-feedback@google.com). 

## Adding a package
To add a package or module, simply fetch it from proxy.golang.org. Documentation is generated based on Go source code downloaded from proxy.golang.org/\<module\>@\<version\>.zip. New module versions are fetched from index.golang.org and added to the go.dev site every few minutes.

The [guidelines](https://blog.golang.org/godoc-documenting-go-code) for writing documentation for the [godoc](https://golang.org/cmd/godoc/) tool apply to go.dev. 

It's important to write a good summary of the package in the first sentence of the package comment. The go.dev site indexes the first sentence and displays it in search.

## Removing a package
If you would like a package to be removed, please send an email to [go-discovery-feedback@google.com](mailto:go-discovery-feedback@google.com), with the import path or module path that you want to remove. 

## License policy
Information for a given package or module may be limited if we are not able to detect a suitable license. See our [license policy](https://pkg.go.dev/license-policy) for more information.
