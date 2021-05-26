# go.dev

### Style Guides

- [CSS](https://golang.org/wiki/CSSStyleGuide)
- [JavaScript](https://google.github.io/styleguide/jsguide.html)

## Installation/Usage

To serve the go.dev pages, run

	go run ./cmd/frontend

## Deploying

All commits pushed to `master` will be automatically deployed to https://go.dev.

For now moment, the deployment is not automatic. Instead, after submitting,
visit the [Cloud Build triggers list](https://console.cloud.google.com/cloud-build/triggers?project=go-discovery),
find the one named “Redeploy-go-dev-on-website-commit”, which should say “Disabled” in the status column,
and then click “RUN”.

## Commands

- Running the server: `go run ./cmd/frontend`
