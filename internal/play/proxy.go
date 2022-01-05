// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package play

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"

	"golang.org/x/website/internal/web"
)

const playgroundURL = "https://play.golang.org"

type Request struct {
	Body string
}

type Response struct {
	Errors string
	Events []Event
}

type Event struct {
	Message string
	Kind    string        // "stdout" or "stderr"
	Delay   time.Duration // time to wait before printing Message
}

const expires = 7 * 24 * time.Hour // 1 week
var cacheControlHeader = fmt.Sprintf("public, max-age=%d", int(expires.Seconds()))

// RegisterHandlers registers handlers for the playground endpoints.
func RegisterHandlers(mux *http.ServeMux, godevSite, chinaSite *web.Site) {
	mux.Handle("/play/", playHandler(godevSite))
	mux.Handle("golang.google.cn/play/", playHandler(chinaSite))
	for _, host := range []string{"golang.org", "go.dev/_", "golang.google.cn/_"} {
		mux.HandleFunc(host+"/compile", compile)
		if host != "golang.google.cn" {
			mux.HandleFunc(host+"/share", share)
		}
		mux.HandleFunc(host+"/fmt", fmtHandler)
	}
}

func compile(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "I only answer to POST requests.", http.StatusMethodNotAllowed)
		return
	}

	ctx := r.Context()

	body := r.FormValue("body")
	res := &Response{}
	req := &Request{Body: body}
	if err := makeCompileRequest(ctx, backend(r), req, res); err != nil {
		log.Printf("ERROR compile error %s: %v", backend(r), err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	var out interface{}
	switch r.FormValue("version") {
	case "2":
		out = res
	default: // "1"
		out = struct {
			CompileErrors string `json:"compile_errors"`
			Output        string `json:"output"`
		}{res.Errors, flatten(res.Events)}
	}
	b, err := json.Marshal(out)
	if err != nil {
		log.Printf("ERROR encoding response: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	expiresTime := time.Now().Add(expires).UTC()
	w.Header().Set("Expires", expiresTime.Format(time.RFC1123))
	w.Header().Set("Cache-Control", cacheControlHeader)
	w.Write(b)
}

// makeCompileRequest sends the given Request to the playground compile
// endpoint and stores the response in the given Response.
func makeCompileRequest(ctx context.Context, backend string, req *Request, res *Response) error {
	reqJ, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("marshaling request: %v", err)
	}
	hReq, _ := http.NewRequest("POST", "https://"+backend+"/compile", bytes.NewReader(reqJ))
	hReq.Header.Set("Content-Type", "application/json")
	hReq = hReq.WithContext(ctx)

	client := &http.Client{
		Timeout: 20 * time.Second,
	}
	r, err := client.Do(hReq)
	if err != nil {
		return fmt.Errorf("making request: %v", err)
	}
	defer r.Body.Close()

	if r.StatusCode != http.StatusOK {
		b, _ := ioutil.ReadAll(r.Body)
		return fmt.Errorf("bad status: %v body:\n%s", r.Status, b)
	}

	if err := json.NewDecoder(r.Body).Decode(res); err != nil {
		return fmt.Errorf("unmarshaling response: %v", err)
	}
	return nil
}

// flatten takes a sequence of Events and returns their contents, concatenated.
func flatten(seq []Event) string {
	var buf bytes.Buffer
	for _, e := range seq {
		buf.WriteString(e.Message)
	}
	return buf.String()
}

var validID = regexp.MustCompile(`^[A-Za-z0-9_\-]+$`)

func share(w http.ResponseWriter, r *http.Request) {
	if id := r.FormValue("id"); r.Method == "GET" && validID.MatchString(id) {
		simpleProxy(w, r, playgroundURL+"/p/"+id+".go")
		return
	}

	simpleProxy(w, r, playgroundURL+"/share")
}

func fmtHandler(w http.ResponseWriter, r *http.Request) {
	simpleProxy(w, r, "https://"+backend(r)+"/fmt")
}

func simpleProxy(w http.ResponseWriter, r *http.Request, url string) {
	if r.Method == "GET" {
		r.Body = nil
	} else if len(r.Form) > 0 {
		r.Body = io.NopCloser(strings.NewReader(r.Form.Encode()))
	}
	req, _ := http.NewRequest(r.Method, url, r.Body)
	req.Header.Set("Content-Type", r.Header.Get("Content-Type"))
	req = req.WithContext(r.Context())
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("ERROR share error: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	copyHeader := func(k string) {
		if v := resp.Header.Get(k); v != "" {
			w.Header().Set(k, v)
		}
	}
	copyHeader("Content-Type")
	copyHeader("Content-Length")
	defer resp.Body.Close()
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

func backend(r *http.Request) string {
	b := r.URL.Query().Get("backend")
	if !isDomainElem(b) {
		return "play.golang.org"
	}
	return b + "play.golang.org"
}

func isDomainElem(s string) bool {
	for i := 0; i < len(s); i++ {
		if !('a' <= s[i] && s[i] <= 'z' || '0' <= s[i] && s[i] <= '9') {
			return false
		}
	}
	return s != ""
}
