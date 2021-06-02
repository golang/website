# go.dev

### Style Guides

- [CSS](https://golang.org/wiki/CSSStyleGuide)
- [JavaScript](https://google.github.io/styleguide/jsguide.html)

## Installation/Usage

To serve the go.dev pages, run

	go run ./cmd/frontend

## Deploying

Each time a CL is reviewed and submitted, the web site is automatically redeployed to
https://go.dev/.

If the automatic redeploy is not working, or to check on the status of a redeploy,
see the “website-redeploy-go-dev” trigger in the
[Cloud Build console](https://console.cloud.google.com/cloud-build/builds?project=golang-org).

## Commands

- Running the server: `go run ./cmd/frontend`
