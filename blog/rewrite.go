// Copyright 2013 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import "net/http"

// Register HTTP handlers that redirect old blog paths to their new locations.
func init() {
	for p := range urlMap {
		dest := "/" + urlMap[p]
		http.HandleFunc(p, func(w http.ResponseWriter, r *http.Request) {
			http.Redirect(w, r, dest, http.StatusMovedPermanently)
		})
	}
}

var urlMap = map[string]string{
	"/2010/03/go-whats-new-in-march-2010.html":               "go-whats-new-in-march-2010",
	"/2010/04/json-rpc-tale-of-interfaces.html":              "json-rpc-tale-of-interfaces",
	"/2010/04/third-party-libraries-goprotobuf-and.html":     "third-party-libraries-goprotobuf-and",
	"/2010/05/go-at-io-frequently-asked-questions.html":      "go-at-io-frequently-asked-questions",
	"/2010/05/new-talk-and-tutorials.html":                   "new-talk-and-tutorials",
	"/2010/05/upcoming-google-io-go-events.html":             "upcoming-google-io-go-events",
	"/2010/06/go-programming-session-video-from.html":        "go-programming-session-video-from",
	"/2010/07/gos-declaration-syntax.html":                   "gos-declaration-syntax",
	"/2010/07/share-memory-by-communicating.html":            "share-memory-by-communicating",
	"/2010/08/defer-panic-and-recover.html":                  "defer-panic-and-recover",
	"/2010/09/go-concurrency-patterns-timing-out-and.html":   "go-concurrency-patterns-timing-out-and",
	"/2010/09/go-wins-2010-bossie-award.html":                "go-wins-2010-bossie-award",
	"/2010/09/introducing-go-playground.html":                "introducing-go-playground",
	"/2010/10/real-go-projects-smarttwitter-and-webgo.html":  "real-go-projects-smarttwitter-and-webgo",
	"/2010/11/debugging-go-code-status-report.html":          "debugging-go-code-status-report",
	"/2010/11/go-one-year-ago-today.html":                    "go-one-year-ago-today",
	"/2011/01/go-slices-usage-and-internals.html":            "go-slices-usage-and-internals",
	"/2011/01/json-and-go.html":                              "json-and-go",
	"/2011/03/c-go-cgo.html":                                 "c-go-cgo",
	"/2011/03/go-becomes-more-stable.html":                   "go-becomes-more-stable",
	"/2011/03/gobs-of-data.html":                             "gobs-of-data",
	"/2011/03/godoc-documenting-go-code.html":                "godoc-documenting-go-code",
	"/2011/04/go-at-heroku.html":                             "go-at-heroku",
	"/2011/04/introducing-gofix.html":                        "introducing-gofix",
	"/2011/05/gif-decoder-exercise-in-go-interfaces.html":    "gif-decoder-exercise-in-go-interfaces",
	"/2011/05/go-and-google-app-engine.html":                 "go-and-google-app-engine",
	"/2011/05/go-at-google-io-2011-videos.html":              "go-at-google-io-2011-videos",
	"/2011/06/first-class-functions-in-go-and-new-go.html":   "first-class-functions-in-go-and-new-go",
	"/2011/06/profiling-go-programs.html":                    "profiling-go-programs",
	"/2011/06/spotlight-on-external-go-libraries.html":       "spotlight-on-external-go-libraries",
	"/2011/07/error-handling-and-go.html":                    "error-handling-and-go",
	"/2011/07/go-for-app-engine-is-now-generally.html":       "go-for-app-engine-is-now-generally",
	"/2011/09/go-image-package.html":                         "go-image-package",
	"/2011/09/go-imagedraw-package.html":                     "go-imagedraw-package",
	"/2011/09/laws-of-reflection.html":                       "laws-of-reflection",
	"/2011/09/two-go-talks-lexical-scanning-in-go-and.html":  "two-go-talks-lexical-scanning-in-go-and",
	"/2011/10/debugging-go-programs-with-gnu-debugger.html":  "debugging-go-programs-with-gnu-debugger",
	"/2011/10/go-app-engine-sdk-155-released.html":           "go-app-engine-sdk-155-released",
	"/2011/10/learn-go-from-your-browser.html":               "learn-go-from-your-browser",
	"/2011/10/preview-of-go-version-1.html":                  "preview-of-go-version-1",
	"/2011/11/go-programming-language-turns-two.html":        "go-programming-language-turns-two",
	"/2011/11/writing-scalable-app-engine.html":              "writing-scalable-app-engine",
	"/2011/12/building-stathat-with-go.html":                 "building-stathat-with-go",
	"/2011/12/from-zero-to-go-launching-on-google.html":      "from-zero-to-go-launching-on-google",
	"/2011/12/getting-to-know-go-community.html":             "getting-to-know-go-community",
	"/2012/03/go-version-1-is-released.html":                 "go-version-1-is-released",
	"/2012/07/gccgo-in-gcc-471.html":                         "gccgo-in-gcc-471",
	"/2012/07/go-videos-from-google-io-2012.html":            "go-videos-from-google-io-2012",
	"/2012/08/go-updates-in-app-engine-171.html":             "go-updates-in-app-engine-171",
	"/2012/08/organizing-go-code.html":                       "organizing-go-code",
	"/2012/11/go-turns-three.html":                           "go-turns-three",
	"/2013/01/concurrency-is-not-parallelism.html":           "concurrency-is-not-parallelism",
	"/2013/01/go-fmt-your-code.html":                         "go-fmt-your-code",
	"/2013/01/the-app-engine-sdk-and-workspaces-gopath.html": "the-app-engine-sdk-and-workspaces-gopath",
	"/2013/01/two-recent-go-talks.html":                      "two-recent-go-talks",
	"/2013/02/getthee-to-go-meetup.html":                     "getthee-to-go-meetup",
	"/2013/02/go-maps-in-action.html":                        "go-maps-in-action",
	"/2013/03/two-recent-go-articles.html":                   "two-recent-go-articles",
	"/2013/03/the-path-to-go-1.html":                         "the-path-to-go-1",
	"/2013/05/go-11-is-released.html":                        "go-11-is-released.article",
	"/2013/05/advanced-go-concurrency-patterns.html":         "advanced-go-concurrency-patterns.article",
}
