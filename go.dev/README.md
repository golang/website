# go.dev

## Contributing

```
git clone sso://partner-code/go.dev && (cd go.dev && f=`git rev-parse --git-dir`/hooks/commit-msg ; mkdir -p $(dirname $f) ; curl -Lo $f https://gerrit-review.googlesource.com/tools/hooks/commit-msg ; chmod +x $f)
```

- `data/learn` contains links for the Learn pages, as all content is currently external.
- `content/solutions` contains Use Cases and Case Studies.
  - Please include relevant resources using the same `name` attribute for images.
- `themes/default` contains the site layout.

### Style Guides
- [CSS](https://golang.org/wiki/CSSStyleGuide)
- [JavaScript](https://google.github.io/styleguide/jsguide.html)

## Deploying

All commits targeting `master` will trigger a CI test defined in `cloudbuild.ci.yaml`.
All commits pushed to `master` will be automatically deployed to https://dev.go.dev.

## Code repo
https://partner-code.git.corp.google.com/go.dev

## Commands

- Running the server:  `hugo server -D`
- Pushing to staging:  `git push -f origin HEAD:staging`

## Where things live
- Javascript:
- Carousels: /static/js/carousels.js
- Tab navigation, filtering, listeners: /static/js/site.js
- Solutions page template: /layouts/solutions/single.html
- Home page template: /layouts/index.html
- Site wide styles: /assets/css/styles.css
- Site configuration: /config.toml
- Promotional components (modal, snackbar, etc) are in this branch: `messaging-components`