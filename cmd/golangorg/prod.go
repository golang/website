// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build go1.16 && prod
// +build go1.16,prod

package main

import (
	"context"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"golang.org/x/website/internal/dl"
	"golang.org/x/website/internal/proxy"
	"golang.org/x/website/internal/redirect"
	"golang.org/x/website/internal/short"

	"cloud.google.com/go/datastore"
	"golang.org/x/website/internal/memcache"
)

func earlySetup() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	log.Println("initializing golang.org server ...")

	port := "8080"
	if p := os.Getenv("PORT"); p != "" {
		port = p
	}
	*httpAddr = ":" + port
}

func lateSetup(mux *http.ServeMux) {
	pres.GoogleAnalytics = os.Getenv("GOLANGORG_ANALYTICS")

	datastoreClient, memcacheClient := getClients()

	dl.RegisterHandlers(mux, datastoreClient, memcacheClient)
	short.RegisterHandlers(mux, datastoreClient, memcacheClient)

	// Register /compile and /share handlers against the default serve mux
	// so that other app modules can make plain HTTP requests to those
	// hosts. (For reasons, HTTPS communication between modules is broken.)
	proxy.RegisterHandlers(http.DefaultServeMux)

	http.HandleFunc("/_ah/health", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "ok")
	})

	http.HandleFunc("/robots.txt", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "User-agent: *\nDisallow: /search\n")
	})

	if err := redirect.LoadChangeMap("hg-git-mapping.bin"); err != nil {
		log.Fatalf("LoadChangeMap: %v", err)
	}

	log.Println("godoc initialization complete")
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
