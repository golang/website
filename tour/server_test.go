// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"log"
	"net/http"
	"testing"

	"golang.org/x/website/internal/webtest"
)

func TestWeb(t *testing.T) {
	if err := initTour("SocketTransport"); err != nil {
		log.Fatal(err)
	}
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/lesson/", lessonHandler)
	registerStatic()

	webtest.TestHandler(t, "testdata/*.txt", http.DefaultServeMux)
}
