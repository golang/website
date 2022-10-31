// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package short implements a simple URL shortener, serving shortened urls
// from /s/key. An administrative handler is provided for other services to use.
package short

// TODO(adg): collect statistics on URL visits

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"cloud.google.com/go/datastore"
	"golang.org/x/website/internal/backport/html/template"
	"golang.org/x/website/internal/memcache"
)

const (
	prefix  = "/s"
	kind    = "Link"
	baseURL = "https://go.dev" + prefix
)

// Link represents a short link.
type Link struct {
	Key, Target string
}

var validKey = regexp.MustCompile(`^[a-zA-Z0-9-_.]+$`)

type server struct {
	datastore *datastore.Client
	memcache  *memcache.CodecClient
}

func newServer(dc *datastore.Client, mc *memcache.Client) *server {
	return &server{
		datastore: dc,
		memcache:  mc.WithCodec(memcache.JSON),
	}
}

func RegisterHandlers(mux *http.ServeMux, host string, dc *datastore.Client, mc *memcache.Client) {
	s := newServer(dc, mc)
	mux.HandleFunc(host+prefix+"/", s.linkHandler)
}

// linkHandler services requests to short URLs.
//
//	https://go.dev/s/key[/remaining/path]
//
// It consults memcache and datastore for the Link for key.
// It then sends a redirects or an error message.
// If the remaining path part is not empty, the redirects
// will be the relative path from the resolved Link.
func (h server) linkHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	key, remainingPath, err := extractKey(r)
	if err != nil { // invalid key or url
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	var link Link
	if err := h.memcache.Get(ctx, cacheKey(key), &link); err != nil {
		k := datastore.NameKey(kind, key, nil)
		err = h.datastore.Get(ctx, k, &link)
		switch err {
		case datastore.ErrNoSuchEntity:
			http.Error(w, "not found", http.StatusNotFound)
			return
		default: // != nil
			log.Printf("ERROR %q: %v", key, err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		case nil:
			item := &memcache.Item{
				Key:    cacheKey(key),
				Object: &link,
			}
			if err := h.memcache.Set(ctx, item); err != nil {
				log.Printf("WARNING %q: %v", key, err)
			}
		}
	}

	target := link.Target
	if remainingPath != "" {
		target += remainingPath
	}
	http.Redirect(w, r, target, http.StatusFound)
}

func extractKey(r *http.Request) (key, remainingPath string, err error) {
	path := r.URL.Path
	if !strings.HasPrefix(path, prefix+"/") {
		return "", "", errors.New("invalid path")
	}

	key, remainingPath = path[len(prefix)+1:], ""
	if slash := strings.Index(key, "/"); slash > 0 {
		key, remainingPath = key[:slash], key[slash:]
	}

	if !validKey.MatchString(key) {
		return "", "", errors.New("invalid key")
	}
	return key, remainingPath, nil
}

// AdminHandler serves an administrative interface for managing shortener entries.
// Be careful. It is the caller’s responsibility to ensure that the handler is
// only exposed to authorized users.
func AdminHandler(dc *datastore.Client, mc *memcache.Client) http.HandlerFunc {
	s := newServer(dc, mc)
	return s.adminHandler
}

var (
	adminTemplate = template.Must(template.New("admin").Parse(templateHTML))

	//go:embed admin.html
	templateHTML string
)

// adminHandler serves an administrative interface.
// Be careful. Ensure that this handler is only be exposed to authorized users.
func (h server) adminHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var newLink *Link
	var doErr error
	if r.Method == "POST" {
		key := r.FormValue("key")
		switch r.FormValue("do") {
		case "Add":
			newLink = &Link{key, r.FormValue("target")}
			doErr = h.putLink(ctx, newLink)
		case "Delete":
			k := datastore.NameKey(kind, key, nil)
			doErr = h.datastore.Delete(ctx, k)
		default:
			http.Error(w, "unknown action", http.StatusBadRequest)
		}
		err := h.memcache.Delete(ctx, cacheKey(key))
		if err != nil && err != memcache.ErrCacheMiss {
			log.Printf("WARNING %q: %v", key, err)
		}
	}

	var links []*Link
	q := datastore.NewQuery(kind).Order("Key")
	if _, err := h.datastore.GetAll(ctx, q, &links); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("ERROR %v", err)
		return
	}

	// Put the new link in the list if it's not there already.
	// (Eventual consistency means that it might not show up
	// immediately, which might be confusing for the user.)
	if newLink != nil && doErr == nil {
		found := false
		for i := range links {
			if links[i].Key == newLink.Key {
				found = true
				break
			}
		}
		if !found {
			links = append([]*Link{newLink}, links...)
		}
		newLink = nil
	}

	var data = struct {
		BaseURL string
		Prefix  string
		Links   []*Link
		New     *Link
		Error   error
	}{baseURL, prefix, links, newLink, doErr}
	if err := adminTemplate.Execute(w, &data); err != nil {
		log.Printf("ERROR adminTemplate: %v", err)
	}
}

// putLink validates the provided link and puts it into the datastore.
func (h server) putLink(ctx context.Context, link *Link) error {
	if !validKey.MatchString(link.Key) {
		return fmt.Errorf("invalid key %q; must match %s", link.Key, validKey.String())
	}
	if _, err := url.Parse(link.Target); err != nil {
		return fmt.Errorf("bad target %q: %v", link.Target, err)
	}
	k := datastore.NameKey(kind, link.Key, nil)
	_, err := h.datastore.Put(ctx, k, link)
	return err
}

// cacheKey returns a short URL key as a memcache key.
func cacheKey(key string) string {
	return "link-" + key
}
