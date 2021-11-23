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
	resp, err := client.Get("https://play.golang.org/versions")
	if err != nil {
		return nil, fmt.Errorf("readVersions: %v", err)
	}
	defer resp.Body.Close()

	js, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("readVersions: %v", err)
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("readVersions: %s\n%s", resp.Status, js)
	}

	if err := json.Unmarshal(js, &list); err != nil || len(list) < 2 {
		return nil, fmt.Errorf("readVersions: bad response: %#q", js)
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
