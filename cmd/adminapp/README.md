# adminapp

This app serves as the [admin interface](https://admin.go.dev/) for the go.dev/s link shortener. It can also remove unwanted playground snippets.

## Deployment:

To update the public site, run:

```
gcloud app --project=symbolic-datum-552 deploy --promote app.yaml
```

We used to use admin.golang.org. Now we run a redirector there. See ../adminredirect.
