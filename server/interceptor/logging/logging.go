package logging

import (
	"WebRTC_POC/server/logging"
	"net/http"
	"strconv"
	"sync/atomic"
)

type Logging struct {
	requestID int32
}

func New() *Logging {
	return &Logging{
		requestID: 0,
	}
}

func (l *Logging) generateID() string {
	next := atomic.AddInt32(&l.requestID, 1)
	return strconv.Itoa(int(next))
}

func (l *Logging) Intercept(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := logging.With(r.Context(), logging.New(l.generateID()))
		nr := r.WithContext(ctx)
		logging.From(nr.Context()).Debugf(`request %d in`, l.requestID)
		h.ServeHTTP(w, nr)
		logging.From(nr.Context()).Debugf(`request %d out`, l.requestID)
	})
}
