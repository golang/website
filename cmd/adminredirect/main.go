// Copyright 2023 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// The adminredirect app redirects traffic to the new admin host.
package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	redirect := os.Getenv("REDIRECT")
	if redirect == "" {
		log.Fatalf("redirect not set")
	}

	h := func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, redirect, http.StatusSeeOther)
	}
	log.Printf("Listening on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, http.HandlerFunc(h)))
}
