# golangorg

## Local Development

For local development, simply build and run. It serves on localhost:6060.
You can specify the domain name as the first path element, such as
http://localhost:6060/go.dev/blog.

	go run .

## Deploying to golang.org

Each time a CL is reviewed and submitted, the site is automatically deployed to App Engine.
If the CL is submitted with a Website-Publish +1 vote,
the new deployment automatically becomes https://golang.org/.
Otherwise, the new deployment can be found in the
[App Engine versions list](https://console.cloud.google.com/appengine/versions?project=golang-org&serviceId=default) and verified and manually promoted.

If the automatic deployment is not working, or to check on the status of a pending deployment,
see the “website-redeploy-golang-org” trigger in the
[Cloud Build console](https://console.cloud.google.com/cloud-build/builds?project=golang-org).
