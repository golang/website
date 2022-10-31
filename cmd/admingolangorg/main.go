// Copyright 2020 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// The admingolangorg command serves an administrative interface for owners of
// the golang-org Google Cloud project.
package main

import (
	"context"
	_ "embed"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"cloud.google.com/go/datastore"
	"golang.org/x/website/internal/memcache"
	"golang.org/x/website/internal/short"
	"google.golang.org/api/idtoken"
)

//go:embed index.html
var index string

func main() {
	audience := os.Getenv("IAP_AUDIENCE")
	dsClient, mcClient := getClients()
	mux := http.NewServeMux()
	mux.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(index))
		return
	}))
	mux.Handle("/shortlink", short.AdminHandler(dsClient, mcClient))
	mux.Handle("/snippet", &snippetHandler{dsClient})
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	log.Printf("Listening on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, iapAuth(audience, mux)))
}

type snippetHandler struct {
	ds *datastore.Client
}

//go:embed snippet.html
var snippetForm string

func (h *snippetHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	snippetLink := r.FormValue("snippet")
	if snippetLink == "" {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(snippetForm))
		return
	}

	prefixes := []string{
		"https://play.golang.org/p/",
		"http://play.golang.org/p/",
		"https://go.dev/play/p/",
		"http://go.dev/play/p/",
	}
	var snippetID string
	for _, p := range prefixes {
		if strings.HasPrefix(snippetLink, p) {
			snippetID = strings.TrimPrefix(snippetLink, p)
			break
		}
	}
	if !strings.Contains(snippetLink, "/") {
		snippetID = snippetLink
	}
	if snippetID == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "must specify snippet URL or ID\n")
		return
	}

	k := datastore.NameKey("Snippet", snippetID, nil)
	if h.ds.Get(r.Context(), k, new(struct{})) == datastore.ErrNoSuchEntity {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Snippet with ID %q does not exist\n", snippetID)
		return
	}
	if err := h.ds.Delete(r.Context(), k); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Unable to delete Snippet with ID %q: %v\n", snippetID, err)
		return
	}
	w.Write([]byte("snippet deleted\n"))
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

func iapAuth(audience string, h http.Handler) http.Handler {
	// https://cloud.google.com/iap/docs/signed-headers-howto#verifying_the_jwt_payload
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jwt := r.Header.Get("x-goog-iap-jwt-assertion")
		if jwt == "" {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprintf(w, "must run under IAP\n")
			return
		}

		payload, err := idtoken.Validate(r.Context(), jwt, audience)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			log.Printf("JWT validation error: %v", err)
			return
		}
		if payload.Issuer != "https://cloud.google.com/iap" {
			w.WriteHeader(http.StatusUnauthorized)
			log.Printf("Incorrect issuer: %q", payload.Issuer)
			return
		}
		if payload.Expires+30 < time.Now().Unix() || payload.IssuedAt-30 > time.Now().Unix() {
			w.WriteHeader(http.StatusUnauthorized)
			log.Printf("Bad JWT times: expires %v, issued %v", time.Unix(payload.Expires, 0), time.Unix(payload.IssuedAt, 0))
			return
		}
		h.ServeHTTP(w, r)
	})
}
