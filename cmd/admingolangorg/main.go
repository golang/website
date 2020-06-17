// Copyright 2020 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// The admingolangorg command serves an administrative interface for owners of
// the golang-org Google Cloud project.
package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"strings"

	"cloud.google.com/go/datastore"
	"golang.org/x/website/internal/memcache"
	"golang.org/x/website/internal/short"
)

func main() {
	http.HandleFunc("/", short.AdminHandler(getClients()))
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	log.Printf("Listening on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func getClients() (*datastore.Client, *memcache.Client) {
	ctx := context.Background()

	datastoreClient, err := datastore.NewClient(ctx, "")
	if err != nil {
		if strings.Contains(err.Error(), "missing project") {
			log.Fatalf("Missing datastore project. Set the DATASTORE_PROJECT_ID env variable. Use `gcloud beta emulators datastore` to start a local datastore.")
		}
		log.Fatalf("datastore.NewClient: %v.", err)
	}

	redisAddr := os.Getenv("GOLANGORG_REDIS_ADDR")
	if redisAddr == "" {
		log.Fatalf("Missing redis server for golangorg in production mode. set GOLANGORG_REDIS_ADDR environment variable.")
	}
	memcacheClient := memcache.New(redisAddr)

	return datastoreClient, memcacheClient
}
