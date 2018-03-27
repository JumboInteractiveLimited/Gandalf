package goserver

import (
	"fmt"
	"net/http"
)

func setupServer() http.Handler {
	http.HandleFunc("/robots.txt", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		fmt.Fprintf(w, "User-agent: *\nDisallow: /")
	})
	return nil
}
