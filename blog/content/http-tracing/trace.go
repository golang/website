// +build OMIT

package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptrace"
)

func trace() {
	// START OMIT
	req, _ := http.NewRequest("GET", "http://example.com", nil)
	trace := &httptrace.ClientTrace{
		DNSDone: func(dnsInfo httptrace.DNSDoneInfo) {
			fmt.Printf("DNS Info: %+v\n", dnsInfo)
		},
		GotConn: func(connInfo httptrace.GotConnInfo) {
			fmt.Printf("Got Conn: %+v\n", connInfo)
		},
	}
	req = req.WithContext(httptrace.WithClientTrace(req.Context(), trace))
	if _, err := http.DefaultTransport.RoundTrip(req); err != nil {
		log.Fatal(err)
	}
	// END OMIT
}
