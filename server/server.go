package server

import (
	"WebRTC_POC/server/backend"
	"WebRTC_POC/server/controller"
	"WebRTC_POC/server/coordinator"
	"WebRTC_POC/server/database/memdb"
	"WebRTC_POC/server/interceptor"
	"WebRTC_POC/server/interceptor/auth"
	"WebRTC_POC/server/interceptor/cors"
	logg "WebRTC_POC/server/interceptor/log"
	"WebRTC_POC/server/interceptor/parse"
	"WebRTC_POC/server/logging"
	"WebRTC_POC/server/profiling"
	"WebRTC_POC/server/profiling/metric"
	"WebRTC_POC/test/client"
	"fmt"
	"log"
	"net/http"
)

type PDN struct {
	apiServer       *http.Server
	profilingServer *profiling.Server
}

func New() *PDN {
	if err := logging.SetLogLevel("debug"); err != nil {
		log.Fatalf("Failed to set log level: %v", err)
	}
	logger := logging.New("PDN")

	cm := coordinator.New()
	me, err := metric.New()
	if err != nil {

	}
	db := memdb.New()

	be := backend.New(cm, me, db)
	con := controller.New(be)
	mw := interceptor.New(parse.New(), auth.New(), cors.New(), logg.New(logger))

	mux := http.NewServeMux()

	fs := client.New()
	mux.Handle("/test/", fs)
	mux.Handle("/channel/", interceptor.WithInterceptors(con, mw))

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
		logging.DefaultLogger().Infof(`Profiler starts to run on :%d`, 8081)
		err := s.profilingServer.Start()
		if err != nil {
			logging.DefaultLogger().Error(err)
			return
		}
	}()

	logging.DefaultLogger().Infof(`PDN starts to run on :%d`, 8080)
	return s.apiServer.ListenAndServe()
}
