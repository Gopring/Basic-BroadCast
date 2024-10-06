package controller

import (
	"WebRTC_POC/service"
	"encoding/json"
	"io"
	"net/http"
)

const (
	broadcast = "/channel/broadcast"
	view      = "/channel/view"
)

type Request struct {
	Key string `json:"key"`
	Sdp string `json:"sdp"`
}

type Controller struct {
	s *service.Service
}

func NewController(service *service.Service) *Controller {
	return &Controller{
		s: service,
	}
}

func (c *Controller) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case broadcast:
		c.Broadcast(w, r, c.s)
	case view:
		c.View(w, r, c.s)
	default:
		http.Error(w, "wrong path", http.StatusNotFound)
	}
}

func (c *Controller) Broadcast(w http.ResponseWriter, r *http.Request, s *service.Service) {
	d, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "failed read body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	req := Request{}
	if err = json.Unmarshal(d, &req); err != nil {
		http.Error(w, "failed parse body", http.StatusBadRequest)
	}
	if err = s.Broadcast(req.Key, req.Sdp); err != nil {
		http.Error(w, "failed broadcast", http.StatusInternalServerError)
	}
}

func (c *Controller) View(w http.ResponseWriter, r *http.Request, s *service.Service) {
	d, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "failed read body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	req := Request{}
	if err = json.Unmarshal(d, &req); err != nil {
		http.Error(w, "failed parse body", http.StatusBadRequest)
	}
	if err = s.View(req.Key, req.Sdp); err != nil {
		http.Error(w, "failed view", http.StatusInternalServerError)
	}
}
