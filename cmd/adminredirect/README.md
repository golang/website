This app redirects traffic at https://admin.golang.org/ to the new admin app.

## Deployment

To update the redirector, run:

```
gcloud app --project=golang-org deploy --promote app.yaml
```
