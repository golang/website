// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build go1.8

// The server command demonstrates a server with HTTP/2 Server Push support.
package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

var httpAddr = flag.String("http", ":8080", "Listen address")

func main() {
	flag.Parse()
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		pusher, ok := w.(http.Pusher)
		if ok {
			// Push is supported. Try pushing rather than
			// waiting for the browser request these static assets.
			if err := pusher.Push("/static/app.js", nil); err != nil {
				log.Printf("Failed to push: %v", err)
			}
			if err := pusher.Push("/static/style.css", nil); err != nil {
				log.Printf("Failed to push: %v", err)
			}
		}
		fmt.Fprintf(w, indexHTML)
	})
	log.Fatal(http.ListenAndServeTLS(*httpAddr, "cert.pem", "key.pem", nil))
}

const indexHTML = `<html>
<head>
	<title>Hello World</title>
	<script src="/static/app.js"></script>
	<link rel="stylesheet" href="/static/style.css"">
</head>
<body>
Hello, gopher!
</body>
</html>
`
