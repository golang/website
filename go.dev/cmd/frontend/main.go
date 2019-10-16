package main

import (
	"log"
	"net"
	"net/http"
	"os"
)

func main() {
	fs := http.FileServer(http.Dir("public/"))
	http.Handle("/", fs)
	http.HandleFunc("/explore", exploreHandler)

	p := listenPort()
	l, err := net.Listen("tcp", ":" + p)
	if err != nil {
		log.Fatalf("net.Listen(%q, %q) = _, %v", "tcp", p, err)
	}
	defer l.Close()
	log.Printf("Listening on http://%v/\n", l.Addr().String())
	log.Print(http.Serve(l, nil))
}

func listenPort() string {
	if p := os.Getenv("PORT"); p != "" {
		return p
	}
	return "0"
}

func exploreHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Host {
	case "dev.go.dev":
		http.Redirect(w, r, "https://dev-pkg.go.dev/", http.StatusFound)
	case "staging.go.dev":
		http.Redirect(w, r, "https://staging-pkg.go.dev/", http.StatusFound)
	default:
		http.Redirect(w, r, "https://pkg.go.dev/", http.StatusFound)
	}
}
