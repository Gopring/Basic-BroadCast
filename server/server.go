package server

import (
	"WebRTC_POC/controller"
	"WebRTC_POC/frontend"
	"WebRTC_POC/server/backend"
	"WebRTC_POC/server/coordinator"
	"WebRTC_POC/server/database/memdb"
	"WebRTC_POC/server/metric"
	"fmt"
	"net/http"
)

type Server struct {
	server *http.Server
}

func New() *Server {
	cm := coordinator.New()
	me := metric.New()
	db := memdb.New()
	be := backend.New(cm, me, db)
	con := controller.New(be)

	mux := http.NewServeMux()

	mux.Handle("/channel", con)

	fs := frontend.New()
	mux.Handle("/", fs)
	return &Server{
		server: &http.Server{
			Addr:    fmt.Sprintf(":%d", 8080),
			Handler: mux,
		},
	}
}

func (s *Server) Start() error {
	return s.server.ListenAndServe()
}
