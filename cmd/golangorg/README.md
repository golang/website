# golangorg

## Local Development

For local development, simply build and run. It serves on localhost:6060.
You can specify the domain name as the first path element, such as
http://localhost:6060/go.dev/blog.

	go run .

## Testing

The go.dev and golang.org web sites have a suite of regression tests that can be run with:

	go test golang.org/x/website/...

Test cases that check for expected URLs, content, response codes and so on are
encoded in \*.txt files in the `testdata` directory. If there is a problem that
no existing test caught, it can be a good idea to add a new test case to avoid
repeat regressions.

These tests can be run locally, via TryBots, and they are also run when
new versions are being deployed. The `testdata/live.txt` file is special
and used only when testing a live server, because its test cases depend
on production resources.

## Screentest

The go.dev web site has a suite of visual checks that can be run with:

	go run ./cmd/screentest

These checks can be run locally and will generate visual diffs of web pages
from the set of testcases in `cmd/screentest/godev.txt`, comparing screenshots
of the live server and a locally running instance of cmd/golangorg.

## Deploying to go.dev and golang.org

Each time a CL is reviewed and submitted, the site is automatically deployed to App Engine.
If it passes its serving-readiness checks, it will be automatically promoted to handle traffic.
Whether it passes or not, the new deployment can be found in the
[App Engine versions list](https://console.cloud.google.com/appengine/versions?project=golang-org&serviceId=default).

If the automatic deployment is not working, or to check on the status of a pending deployment,
see the “website-redeploy-golang-org” trigger in the
[Cloud Build console](https://console.cloud.google.com/cloud-build/builds;region=global?project=golang-org&query=trigger_id%3D%222399003e-0cc5-4877-86de-8bc8f13fd984%22).
