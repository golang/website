# Go website

[![Go Reference](https://pkg.go.dev/badge/golang.org/x/website.svg)](https://pkg.go.dev/golang.org/x/website)

This repo holds content and serving programs for the go.dev and golang.org web sites.

Content is in \_content/ (go.dev) and tour/ (go.dev/tour).
Server code is in cmd/ and internal/.

To run the combined go.dev+golang.org server to preview local content changes, use:

    go run ./cmd/golangorg

The supporting programs cmd/admingolangorg and cmd/googlegolangorg
are the servers for admin.golang.org and google.golang.org.
(They do not use the \_content/ directories.)

Each command directory has its own README.md explaining deployment.

## JS/TS/CSS Formatting

This repository uses [eslint](https://eslint.org/) to format JS and TS files,
and [stylelint](https://stylelint.io/) to format CSS files.

See also:

- [CSS](https://go.dev/wiki/CSSStyleGuide)
- [JavaScript](https://google.github.io/styleguide/jsguide.html)
- [TypeScript](https://google.github.io/styleguide/tsguide.html)

It is encouraged that all JS, TS, and CSS code be run through formatters before
submitting a change. However, it is not a strict requirement enforced by CI.

### Installing npm Dependencies:

1. Install [docker](https://docs.docker.com/get-docker/)
2. Create a .gitignore file at repo root
3. Add .gitignore and node_modules to .gitignore
4. Run `./npm install`

### Run ESlint

    ./npx eslint [options] [file] [dir]

### Run Stylelint

    ./npx stylelint [input] [options]

## TypeScript Support

TypeScript files served from _content are transformed into JavaScript.
Reference .ts files in html templates as module code.

  `<script type="module" src="/ts/filename.ts">`

Write unit tests for TypeScript code using the [jest](https://jestjs.io/)
testing framework.

### Run Jest

    ./npx jest [TestPathPattern]

## Deploying

Each time a CL is reviewed and submitted, the code is deployed to App Engine.
See [cmd/golangorg/README.md](cmd/golangorg/README.md#deploying-to-go_dev-and-golang_org) for details.

## Report Issues / Send Patches

This repository uses Gerrit for code changes. To learn how to submit changes to
this repository, see https://go.dev/doc/contribute.

The main issue tracker for the website repository is located at
https://github.com/golang/go/issues. Prefix your issue with "x/website:" in the
subject line, so it is easy to find.
