# go.dev

## Contributing

- `data/learn` contains links for the Learn pages, as all content is currently external.
- `content/solutions` contains Use Cases and Case Studies.
  - Please include relevant resources using the same `name` attribute for images.
- `themes/default` contains the site layout.

## Deploying

All commits targeting `master` will trigger a CI test defined in `cloudbuild.ci.yaml`.

All commits pushed to `master` will be automatically deployed to https://dev.go.dev.
