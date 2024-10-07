package auth

import "net/http"

type Auth struct {
}

func New() *Auth {
	return &Auth{}
}

func (a *Auth) Intercept(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		h.ServeHTTP(w, r)
	})
}
