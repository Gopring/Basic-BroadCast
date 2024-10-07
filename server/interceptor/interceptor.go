package interceptor

import "net/http"

type Interceptor interface {
	Intercept(h http.Handler) http.Handler
}

type Interceptors struct {
	interceptors []Interceptor
}

func New(ics ...Interceptor) Interceptors {
	m := Interceptors{}
	for _, i := range ics {
		m.interceptors = append(m.interceptors, i)
	}
	return m
}

func WithInterceptors(h http.Handler, m Interceptors) http.Handler {
	for _, i := range m.interceptors {
		h = i.Intercept(h)
	}
	return h
}
