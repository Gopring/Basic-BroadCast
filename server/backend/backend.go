package backend

import (
	"WebRTC_POC/server/coordinator"
	"WebRTC_POC/server/database"
	"WebRTC_POC/server/profiling/metric"
)

type Backend struct {
	Coordinator *coordinator.Coordinator
	Metric      *metric.Metrics
	Database    database.Database
}

func New(ch *coordinator.Coordinator,
	me *metric.Metrics,
	db database.Database,
) *Backend {
	return &Backend{
		Coordinator: ch,
		Metric:      me,
		Database:    db,
	}
}
