// Copyright 2013 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build appengine

// This file implements an App Engine blog server.

package main

import (
	"net/http"

	"golang.org/x/tools/blog"
)

func init() {
	config.ContentPath = "content/"
	config.TemplatePath = "template/"
	s, err := blog.NewServer(config)
	if err != nil {
		panic(err)
	}
	http.Handle("/", s)
}
