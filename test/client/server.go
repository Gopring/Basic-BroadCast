package client

import (
	"net/http"
)

func New() http.Handler {
	mux := http.NewServeMux()
	fs := http.FileServer(http.Dir("./test/client/static"))
	mux.Handle("/test/", http.StripPrefix("/test", fs))
	mux.HandleFunc("/test/html", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./test/client/static/demo.html")
	})
	return mux
}
