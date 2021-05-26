// +build OMIT

package main

import (
	"log"
	"net/http"
)

func main() {
	// START OMIT
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if pusher, ok := w.(http.Pusher); ok {
			// Push is supported.
			if err := pusher.Push("/app.js", nil); err != nil {
				log.Printf("Failed to push: %v", err)
			}
		}
		// ...
	})
	// END OMIT

	// START1 OMIT
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if pusher, ok := w.(http.Pusher); ok {
			// Push is supported.
			options := &http.PushOptions{
				Header: http.Header{
					"Accept-Encoding": r.Header["Accept-Encoding"],
				},
			}
			if err := pusher.Push("/app.js", options); err != nil {
				log.Printf("Failed to push: %v", err)
			}
		}
		// ...
	})
	// END1 OMIT
}
