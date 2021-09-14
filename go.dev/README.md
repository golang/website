# go.dev

### Style Guides

- [CSS](https://golang.org/wiki/CSSStyleGuide)
- [JavaScript](https://google.github.io/styleguide/jsguide.html)

## Installation/Usage

To serve the go.dev pages, run

	go run ./cmd/golangorg

and load http://localhost:6060/go.dev/

## Deploying

Each time a CL is reviewed and submitted, the blog is automatically deployed to App Engine.
If the CL is submitted with a Website-Publish +1 vote,
the new deployment automatically becomes https://go.dev/.
Otherwise, the new deployment can be found in the
[App Engine versions list](https://console.cloud.google.com/appengine/versions?project=go-discovery&serviceId=go-dev) and verified and manually promoted.

If the automatic deployment is not working, or to check on the status of a pending deployment,
see the “website-redeploy-go-dev” trigger in the
[Cloud Build console](https://console.cloud.google.com/cloud-build/builds?project=go-discovery).
