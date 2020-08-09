// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dl

import (
	"context"
	"crypto/hmac"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"html"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"cloud.google.com/go/datastore"
	"golang.org/x/website/internal/env"
	"golang.org/x/website/internal/memcache"
)

type server struct {
	datastore *datastore.Client
	memcache  *memcache.CodecClient
}

func RegisterHandlers(mux *http.ServeMux, dc *datastore.Client, mc *memcache.Client) {
	s := server{dc, mc.WithCodec(memcache.Gob)}
	mux.HandleFunc("/dl", s.getHandler)
	mux.HandleFunc("/dl/", s.getHandler) // also serves listHandler
	mux.HandleFunc("/dl/upload", s.uploadHandler)

	// NOTE(cbro): this only needs to be run once per project,
	// and should be behind an admin login.
	// TODO(cbro): move into a locally-run program? or remove?
	// mux.HandleFunc("/dl/init", initHandler)
}

// rootKey is the ancestor of all File entities.
var rootKey = datastore.NameKey("FileRoot", "root", nil)

func (h server) listHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" && r.Method != "OPTIONS" {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	ctx := r.Context()
	d := listTemplateData{
		GoogleCN: googleCN(r),
	}

	if err := h.memcache.Get(ctx, cacheKey, &d); err != nil {
		if err != memcache.ErrCacheMiss {
			log.Printf("ERROR cache get error: %v", err)
			// NOTE(cbro): continue to hit datastore if the memcache is down.
		}

		var fs []File
		q := datastore.NewQuery("File").Ancestor(rootKey)
		if _, err := h.datastore.GetAll(ctx, q, &fs); err != nil {
			log.Printf("ERROR error listing: %v", err)
			http.Error(w, "Could not get download page. Try again in a few minutes.", 500)
			return
		}
		d.Stable, d.Unstable, d.Archive = filesToReleases(fs)
		if len(d.Stable) > 0 {
			d.Featured = filesToFeatured(d.Stable[0].Files)
		}

		item := &memcache.Item{Key: cacheKey, Object: &d, Expiration: cacheDuration}
		if err := h.memcache.Set(ctx, item); err != nil {
			log.Printf("ERROR cache set error: %v", err)
		}
	}

	if r.URL.Query().Get("mode") == "json" {
		serveJSON(w, r, d)
		return
	}

	if err := listTemplate.ExecuteTemplate(w, "root", d); err != nil {
		log.Printf("ERROR executing template: %v", err)
	}
}

// serveJSON serves a JSON representation of d. It assumes that requests are
// limited to GET and OPTIONS, the latter used for CORS requests, which this
// endpoint supports.
func serveJSON(w http.ResponseWriter, r *http.Request, d listTemplateData) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	if r.Method == "OPTIONS" {
		// Likely a CORS preflight request.
		w.WriteHeader(http.StatusNoContent)
		return
	}
	var releases []Release
	includes := strings.Split(r.URL.Query().Get("include"), ",")
	includesMap := map[string]bool{}
	for _, include := range includes {
		includesMap[include] = true
	}
	if includesMap["all"] {
		releases = append(releases, d.Stable...)
		releases = append(releases, d.Archive...)
		releases = append(releases, d.Unstable...)
	} else {
		if includesMap["stable"] {
			releases = append(releases, d.Stable...)
		}
		if includesMap["archive"] {
			releases = append(releases, d.Archive...)
		}
		if includesMap["unstable"] {
			releases = append(releases, d.Unstable...)
		}
	}
	if len(releases) == 0 {
		releases = d.Stable
	}
	w.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	enc.SetIndent("", " ")
	if err := enc.Encode(releases); err != nil {
		log.Printf("ERROR rendering JSON for releases: %v", err)
	}
}

// googleCN reports whether request r is considered
// to be served from golang.google.cn.
// TODO: This is duplicated within internal/proxy. Move to a common location.
func googleCN(r *http.Request) bool {
	if r.FormValue("googlecn") != "" {
		return true
	}
	if strings.HasSuffix(r.Host, ".cn") {
		return true
	}
	if !env.CheckCountry() {
		return false
	}
	switch r.Header.Get("X-Appengine-Country") {
	case "", "ZZ", "CN":
		return true
	}
	return false
}

func (h server) uploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	ctx := r.Context()

	// Authenticate using a user token (same as gomote).
	user := r.FormValue("user")
	if !validUser(user) {
		http.Error(w, "bad user", http.StatusForbidden)
		return
	}
	if r.FormValue("key") != h.userKey(ctx, user) {
		http.Error(w, "bad key", http.StatusForbidden)
		return
	}

	var f File
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&f); err != nil {
		log.Printf("ERROR decoding upload JSON: %v", err)
		http.Error(w, "Something broke", http.StatusInternalServerError)
		return
	}
	if f.Filename == "" {
		http.Error(w, "Must provide Filename", http.StatusBadRequest)
		return
	}
	if f.Uploaded.IsZero() {
		f.Uploaded = time.Now()
	}
	k := datastore.NameKey("File", f.Filename, rootKey)
	if _, err := h.datastore.Put(ctx, k, &f); err != nil {
		log.Printf("ERROR File entity: %v", err)
		http.Error(w, "could not put File entity", http.StatusInternalServerError)
		return
	}
	if err := h.memcache.Delete(ctx, cacheKey); err != nil {
		log.Printf("ERROR delete error: %v", err)
	}
	io.WriteString(w, "OK")
}

func (h server) getHandler(w http.ResponseWriter, r *http.Request) {
	isGoGet := (r.Method == "GET" || r.Method == "HEAD") && r.FormValue("go-get") == "1"
	// For go get, we need to serve the same meta tags at /dl for cmd/go to
	// validate against the import path.
	if r.URL.Path == "/dl" && isGoGet {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprintf(w, `<!DOCTYPE html><html><head>
<meta name="go-import" content="golang.org/dl git https://go.googlesource.com/dl">
</head></html>`)
		return
	}
	if r.URL.Path == "/dl" {
		http.Redirect(w, r, "/dl/", http.StatusFound)
		return
	}

	name := strings.TrimPrefix(r.URL.Path, "/dl/")
	var redirectURL string
	switch {
	case name == "":
		h.listHandler(w, r)
		return
	case fileRe.MatchString(name):
		// This is a /dl/{file} request to download a file. It's implemented by
		// redirecting to another host, which serves the bytes more efficiently.
		//
		// The redirect target is an internal implementation detail and may change
		// if there is a good reason to do so. Last time was in CL 76971 (in 2017).
		const downloadBaseURL = "https://dl.google.com/go/"
		http.Redirect(w, r, downloadBaseURL+name, http.StatusFound)
		return
	case name == "gotip":
		redirectURL = "https://godoc.org/golang.org/dl/gotip"
	case goGetRe.MatchString(name):
		redirectURL = "https://golang.org/dl/#" + name
	default:
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if !isGoGet {
		w.Header().Set("Location", redirectURL)
	}
	fmt.Fprintf(w, `<!DOCTYPE html>
<html>
<head>
<meta name="go-import" content="golang.org/dl git https://go.googlesource.com/dl">
<meta http-equiv="refresh" content="0; url=%s">
</head>
<body>
Nothing to see here; <a href="%s">move along</a>.
</body>
</html>
`, html.EscapeString(redirectURL), html.EscapeString(redirectURL))
}

func (h server) initHandler(w http.ResponseWriter, r *http.Request) {
	var fileRoot struct {
		Root string
	}
	ctx := r.Context()
	k := rootKey
	_, err := h.datastore.RunInTransaction(ctx, func(tx *datastore.Transaction) error {
		err := tx.Get(k, &fileRoot)
		if err != nil && err != datastore.ErrNoSuchEntity {
			return err
		}
		_, err = tx.Put(k, &fileRoot)
		return err
	}, nil)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	io.WriteString(w, "OK")
}

func (h server) userKey(c context.Context, user string) string {
	hash := hmac.New(md5.New, []byte(h.secret(c)))
	hash.Write([]byte("user-" + user))
	return fmt.Sprintf("%x", hash.Sum(nil))
}

// Code below copied from x/build/app/key

var theKey struct {
	sync.RWMutex
	builderKey
}

type builderKey struct {
	Secret string
}

func (k *builderKey) Key() *datastore.Key {
	return datastore.NameKey("BuilderKey", "root", nil)
}

func (h server) secret(ctx context.Context) string {
	// check with rlock
	theKey.RLock()
	k := theKey.Secret
	theKey.RUnlock()
	if k != "" {
		return k
	}

	// prepare to fill; check with lock and keep lock
	theKey.Lock()
	defer theKey.Unlock()
	if theKey.Secret != "" {
		return theKey.Secret
	}

	// fill
	if err := h.datastore.Get(ctx, theKey.Key(), &theKey.builderKey); err != nil {
		if err == datastore.ErrNoSuchEntity {
			// If the key is not stored in datastore, write it.
			// This only happens at the beginning of a new deployment.
			// The code is left here for SDK use and in case a fresh
			// deployment is ever needed.  "gophers rule" is not the
			// real key.
			if env.RequireDLSecretKey() {
				panic("lost key from datastore")
			}
			theKey.Secret = "gophers rule"
			h.datastore.Put(ctx, theKey.Key(), &theKey.builderKey)
			return theKey.Secret
		}
		panic("cannot load builder key: " + err.Error())
	}

	return theKey.Secret
}
