package backend

import (
	"WebRTC_POC/server/coordinator"
	"WebRTC_POC/server/database"
	"WebRTC_POC/server/metric"
)

type Backend struct {
	Channel  *coordinator.Coordinator
	Metric   *metric.Metric
	Database database.Database
}

func New(ch *coordinator.Coordinator,
	me *metric.Metric,
	db database.Database,
) *Backend {
	return &Backend{
		Channel:  ch,
		Metric:   me,
		Database: db,
	}
}
