package frontend

import "net/http"

func New() http.Handler {
	mux := http.NewServeMux()
	fs := http.FileServer(http.Dir("./frontend/static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./frontend/static/demo.html")
	})
	return mux
}
