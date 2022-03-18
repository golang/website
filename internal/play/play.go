// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package play

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"sync/atomic"
	"time"

	"golang.org/x/website/internal/web"
)

type playVersion struct {
	Name    string
	Backend string
}

var playVersions atomic.Value

func playHandler(site *web.Site) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/play/" && !strings.HasPrefix(r.URL.Path, "/play/p/") || r.URL.Path == "/play/p/" {
			http.Redirect(w, r, "/play/", http.StatusFound)
			return
		}
		if r.Host == "golang.google.cn" && strings.HasPrefix(r.URL.Path, "/play/p/") {
			site.ServeError(w, r, errors.New("Sorry, but shared playground snippets are not visible in China."))
			return
		}
		if strings.HasSuffix(r.URL.Path, ".go") {
			simpleProxy(w, r, "https://"+backend(r)+strings.TrimPrefix(r.URL.Path, "/play"))
			return
		}
		site.ServePage(w, r, web.Page{
			"URL":          r.URL.Path,
			"layout":       "play",
			"title":        "Go Playground",
			"playVersions": playVersions.Load().([]playVersion),
		})
	})
}

func init() {
	// Default set until we get the full list from play.golang.org.
	playVersions.Store([]playVersion{
		{Name: "Go release", Backend: ""},
		{Name: "Go previous release", Backend: "goprev"},
		{Name: "Go dev branch", Backend: "gotip"},
	})

	go watchVersions()
}

func readVersions() (list []playVersion, err error) {
	defer func() {
		if e := recover(); e != nil {
			list = nil
			err = fmt.Errorf("readVersions: %v", e)
		}
	}()

	client := &http.Client{
		Timeout: 2 * time.Minute,
	}

	list = append([]playVersion(nil), playVersions.Load().([]playVersion)...)
	for i, v := range list {
		resp, err := client.Get("https://" + v.Backend + "play.golang.org/version")
		if err != nil {
			log.Printf("readVersions: %v", err)
			continue
		}
		js, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			log.Printf("readVersions: %v", err)
			continue
		}
		if resp.StatusCode != 200 {
			log.Printf("readVersions: %s\n%s", resp.Status, js)
			continue
		}
		var data struct{ Name string }
		if err := json.Unmarshal(js, &data); err != nil {
			log.Printf("readVersions: %v", err)
			continue
		}
		if data.Name != "" {
			list[i].Name = data.Name
		}
	}
	return list, nil
}

func watchVersions() {
	for ; ; time.Sleep(1 * time.Minute) {
		list, err := readVersions()
		if err != nil {
			log.Print(err)
			continue
		}
		playVersions.Store(list)
	}
}
