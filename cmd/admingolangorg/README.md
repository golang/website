# admingolangorg

This app serves as the [admin interface](https://admin-dot-golang-org.appspot.com) for the go.dev/s link shortener. It can also remove unwanted playground snippets.

## Deployment:

To update the public site, run:

```
gcloud app --project=golang-org deploy --promote app.yaml
```
