package cors

import "net/http"

type Cors struct {
}

func New() *Cors {
	return &Cors{}
}

func (c *Cors) Intercept(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		h.ServeHTTP(w, r)
	})
}
