package handler

import (
	"WebRTC_POC/server/controller"
	"WebRTC_POC/server/logging"
	"WebRTC_POC/types/request"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const (
	broadcast = "/channel/broadcast"
	view      = "/channel/view"
)

type Handler struct {
	controller *controller.Controller
}

func New(con *controller.Controller) *Handler {
	return &Handler{
		controller: con,
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			logging.From(r.Context()).Error(err)
		}
	}(r.Body)

	d, err := io.ReadAll(r.Body)
	if err != nil {
		logging.From(r.Context()).Error(err)
		http.Error(w, "failed read body", http.StatusBadRequest)
		return
	}

	req := &request.Request{}
	if err = json.Unmarshal(d, &req); err != nil {
		logging.From(r.Context()).Error(err)
		http.Error(w, "failed parse body", http.StatusBadRequest)
		return
	}
	ctx := request.With(r.Context(), req)

	var sdp string

	switch r.URL.Path {
	case broadcast:
		sdp, err = h.controller.Broadcast(ctx)
	case view:
		sdp, err = h.controller.View(ctx)
	default:
		http.Error(w, "wrong path", http.StatusNotFound)
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	_, err = fmt.Fprint(w, sdp)
	if err != nil {
		logging.From(r.Context()).Error(err)
		return
	}

}
