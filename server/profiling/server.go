package profiling

import (
	"WebRTC_POC/server/profiling/metric"
	"fmt"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

type Server struct {
	httpServer *http.Server
}

func New(me *metric.Metrics) *Server {
	mux := http.NewServeMux()
	if me != nil {
		mux.Handle("/metrics", promhttp.HandlerFor(me.Registry(), promhttp.HandlerOpts{}))
	}

	return &Server{
		httpServer: &http.Server{
			Addr:    fmt.Sprintf(":%d", 8081),
			Handler: mux,
		},
	}
}

func (s *Server) Start() error {
	return s.httpServer.ListenAndServe()
}
