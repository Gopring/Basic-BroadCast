package log

import (
	"WebRTC_POC/server/logging"
	"net/http"
	"strconv"
	"sync/atomic"
)

type Logging struct {
	requestID int32
	logger    logging.Logger
}

func New(l logging.Logger) *Logging {
	return &Logging{
		requestID: 0,
		logger:    l,
	}
}

func (l *Logging) generateID() string {
	next := atomic.AddInt32(&l.requestID, 1)
	return strconv.Itoa(int(next))
}

func (l *Logging) Intercept(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := logging.With(r.Context(), l.logger.With("request_id", l.generateID()))
		nr := r.WithContext(ctx)
		logging.From(nr.Context()).Named("HTTP").Debugf("Started %s %s", nr.Method, nr.URL.Path)
		h.ServeHTTP(w, nr)
		logging.From(nr.Context()).Named("HTTP").Debugf("Completed %s %s", nr.Method, nr.URL.Path)
	})
}
