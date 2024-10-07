package controller

import (
	"WebRTC_POC/server/backend"
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
	backend *backend.Backend
}

func New(be *backend.Backend) *Controller {
	return &Controller{
		backend: be,
	}
}

func (c *Controller) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case broadcast:
		Broadcast(w, r, c.backend)
	case view:
		View(w, r, c.backend)
	default:
		http.Error(w, "wrong path", http.StatusNotFound)
	}
}

func Broadcast(w http.ResponseWriter, r *http.Request, be *backend.Backend) {
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
	if err = be.Coordinator.Broadcast(req.Key, req.Sdp); err != nil {
		http.Error(w, "failed broadcast", http.StatusInternalServerError)
	}
}

func View(w http.ResponseWriter, r *http.Request, be *backend.Backend) {
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
	if err = be.Coordinator.View(req.Key, req.Sdp); err != nil {
		http.Error(w, "failed view", http.StatusInternalServerError)
	}
}
