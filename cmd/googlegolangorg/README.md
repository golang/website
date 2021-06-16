This trivial App Engine app serves the small meta+redirector HTML pages
for https://google.golang.org/. For example:

- https://google.golang.org/appengine
- https://google.golang.org/cloud
- https://google.golang.org/api/
- https://google.golang.org/api/storage/v1
- https://google.golang.org/grpc

The page includes a meta tag to instruct the go tool to translate e.g. the
import path "google.golang.org/appengine" to "github.com/golang/appengine".
See `go help importpath` for the mechanics.

To update the public site, run:

```
gcloud app --project=golang-org deploy --promote app.yaml
```
