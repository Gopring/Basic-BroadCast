package server

import (
	"WebRTC_POC/controller"
	"WebRTC_POC/frontend"
	"WebRTC_POC/server/backend"
	"WebRTC_POC/server/coordinator"
	"WebRTC_POC/server/database/memdb"
	"WebRTC_POC/server/interceptor"
	"WebRTC_POC/server/interceptor/auth"
	"WebRTC_POC/server/interceptor/cors"
	"WebRTC_POC/server/logging"
	"WebRTC_POC/server/profiling"
	"WebRTC_POC/server/profiling/metric"
	"fmt"
	"net/http"
)

type PDN struct {
	apiServer       *http.Server
	profilingServer *profiling.Server
}

func New() *PDN {
	err := logging.SetLogLevel("debug")
	if err != nil {

	}
	cm := coordinator.New()
	me, err := metric.New()
	if err != nil {

	}
	db := memdb.New()

	be := backend.New(cm, me, db)
	con := controller.New(be)
	mw := interceptor.New(auth.New(), cors.New())

	mux := http.NewServeMux()

	mux.Handle("/channel", interceptor.WithInterceptors(con, mw))

	fs := frontend.New()
	mux.Handle("/", fs)

	ps := profiling.New(me)

	return &PDN{
		apiServer: &http.Server{
			Addr:    fmt.Sprintf(":%d", 8080),
			Handler: mux,
		},
		profilingServer: ps,
	}
}

func (s *PDN) Start() error {

	go func() {
		s.profilingServer.Start()
	}()
	fmt.Printf("PDN starts to run on :%d\n", 8080)
	return s.apiServer.ListenAndServe()
}
