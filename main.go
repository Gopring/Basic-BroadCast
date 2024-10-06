package main

import (
	"WebRTC_POC/controller"
	"WebRTC_POC/frontend"
	"WebRTC_POC/service"
	"fmt"
	"log"
	"net/http"
)

type Server struct {
	server *http.Server
}

func NewServer() *Server {
	svc := &service.Service{}
	con := controller.NewController(svc)
	fs := frontend.NewServer()
	mux := http.NewServeMux()
	mux.Handle("/channel", con)
	mux.Handle("/", fs)
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", 8080),
		Handler: mux,
	}
	return &Server{
		server: srv,
	}
}

func (s *Server) Run() {
	log.Fatal(s.server.ListenAndServe())
}

func main() {
	server := NewServer()
	server.Run()
}
