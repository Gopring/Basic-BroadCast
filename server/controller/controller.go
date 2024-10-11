package controller

import (
	"WebRTC_POC/server/backend"
	"WebRTC_POC/server/logging"
	"WebRTC_POC/server/rtc/connection"
	"WebRTC_POC/types/request"
	"fmt"
	"net/http"
)

const (
	broadcast = "/channel/broadcast"
	view      = "/channel/view"
)

type Controller struct {
	be *backend.Backend
}

func New(b *backend.Backend) *Controller {
	return &Controller{
		be: b,
	}
}

func (c *Controller) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case broadcast:
		c.Broadcast(w, r)
	case view:
		c.View(w, r)
	default:
		http.Error(w, "wrong path", http.StatusNotFound)
	}
}

func (c *Controller) Broadcast(w http.ResponseWriter, r *http.Request) {
	req := request.From(r.Context())
	if req == nil {
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		return
	}

	ch, err := c.be.Coordinator.CreateChannel(req.ID)
	if err != nil {
		http.Error(w, "channel already exists", http.StatusInternalServerError)
		return
	}

	conn, err := connection.NewInboundConnection(r.Context(), ch.Config)
	if err != nil {
		c.be.Coordinator.RemoveChannel(req.ID)
		http.Error(w, "failed to make connection", http.StatusInternalServerError)
		return
	}

	ch.SetUpstream(r.Context(), conn, req.ID)

	err = conn.StartICE(r.Context(), req.SDP)
	if err != nil {
		c.be.Coordinator.RemoveChannel(req.ID)
		http.Error(w, "failed to ICE", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	_, err = fmt.Fprint(w, conn.ServerSDP())
	if err != nil {
		logging.From(r.Context()).Error(err)
		return
	}
}

func (c *Controller) View(w http.ResponseWriter, r *http.Request) {
	req := request.From(r.Context())
	if req == nil {
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		return
	}

	ch, err := c.be.Coordinator.GetChannel(req.ID)
	if err != nil {
		http.Error(w, "channel not exists", http.StatusInternalServerError)
		return
	}

	conn, err := connection.NewOutboundConnection(r.Context(), ch.Config)
	if err != nil {
		http.Error(w, "failed to make connection", http.StatusInternalServerError)
		return
	}

	err = ch.SetDownstream(r.Context(), conn)
	if err != nil {
		http.Error(w, "failed to set down stream", http.StatusInternalServerError)
		return
	}

	err = conn.StartICE(r.Context(), req.SDP)
	if err != nil {
		http.Error(w, "failed to ICE", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	_, err = fmt.Fprint(w, conn.ServerSDP())
	if err != nil {
		logging.From(r.Context()).Error(err)
		return
	}
}
