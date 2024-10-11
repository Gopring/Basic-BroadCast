package parse

import (
	"WebRTC_POC/server/logging"
	"WebRTC_POC/types/request"
	"encoding/json"
	"io"
	"net/http"
)

type Parse struct {
}

func New() *Parse {
	return &Parse{}
}

func (p *Parse) Intercept(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
		nr := r.WithContext(ctx)
		next.ServeHTTP(w, nr)
	})
}
