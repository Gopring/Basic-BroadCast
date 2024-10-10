package backend

import (
	"WebRTC_POC/server/channels"
	"WebRTC_POC/server/database"
	"WebRTC_POC/server/profiling/metric"
)

type Backend struct {
	Channels *channels.Channels
	Metric   *metric.Metrics
	Database database.Database
}

func New(ch *channels.Channels,
	me *metric.Metrics,
	db database.Database,
) *Backend {
	return &Backend{
		Channels: ch,
		Metric:   me,
		Database: db,
	}
}
