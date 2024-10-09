package metric

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const (
	channelID = "channel_id"
	apiKey    = "api_key"
)

type Metrics struct {
	registry *prometheus.Registry

	upstreamTotal   *prometheus.GaugeVec
	downstreamTotal *prometheus.GaugeVec

	viewRequestTotal      *prometheus.CounterVec
	broadcastRequestTotal *prometheus.CounterVec
}

func New() (*Metrics, error) {
	reg := prometheus.NewRegistry()
	if err := reg.Register(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{})); err != nil {
		return nil, fmt.Errorf("register process collector: %w", err)
	}
	if err := reg.Register(collectors.NewGoCollector()); err != nil {
		return nil, fmt.Errorf("register go collector: %w", err)
	}
	return &Metrics{
		registry: reg,
		upstreamTotal: promauto.With(reg).NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "PDN",
			Subsystem: "channel",
			Name:      "active_upstream_total",
			Help:      "Total number of upstreams",
		}, []string{apiKey, channelID}),
		downstreamTotal: promauto.With(reg).NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "PDN",
			Subsystem: "channel",
			Name:      "active_downstream_total",
			Help:      "Total number of upstreams",
		}, []string{apiKey, channelID}),
		broadcastRequestTotal: promauto.With(reg).NewCounterVec(prometheus.CounterOpts{
			Namespace: "PDN",
			Subsystem: "request",
			Name:      "broadcast_request_total",
			Help:      "Total number of broadcast request",
		}, []string{apiKey, channelID}),
		viewRequestTotal: promauto.With(reg).NewCounterVec(prometheus.CounterOpts{
			Namespace: "PDN",
			Subsystem: "request",
			Name:      "view_request_total",
			Help:      "Total number of broadcast request",
		}, []string{apiKey, channelID}),
	}, nil
}

func (m *Metrics) Registry() *prometheus.Registry {
	return m.registry
}

func (m *Metrics) AddUpstream(key string, id string) {
	m.upstreamTotal.With(prometheus.Labels{
		apiKey:    key,
		channelID: id,
	}).Inc()
}

func (m *Metrics) RemoveUpstream(key string, id string) {
	m.upstreamTotal.With(prometheus.Labels{
		apiKey:    key,
		channelID: id,
	}).Dec()
}

func (m *Metrics) AddDownstream(key string, id string) {
	m.downstreamTotal.With(prometheus.Labels{
		apiKey:    key,
		channelID: id,
	}).Inc()
}

func (m *Metrics) RemoveDownstream(key string, id string) {
	m.downstreamTotal.With(prometheus.Labels{
		apiKey:    key,
		channelID: id,
	}).Dec()
}

func (m *Metrics) AddBroadcastRequest(key string, id string) {
	m.broadcastRequestTotal.With(prometheus.Labels{
		apiKey:    key,
		channelID: id,
	}).Inc()
}

func (m *Metrics) AddViewRequest(key string, id string) {
	m.viewRequestTotal.With(prometheus.Labels{
		apiKey:    key,
		channelID: id,
	}).Inc()
}
