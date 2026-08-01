// Package metrics provides a Fiber middleware that records Prometheus metrics
// for HTTP requests and exposes them over an HTTP endpoint.
package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	promcollectors "github.com/prometheus/client_golang/prometheus/collectors"
)

const (
	labelMethod = "method"
	labelRoute  = "route"
	labelStatus = "status"
)

type collectors struct {
	requestsTotal    *prometheus.CounterVec
	requestDuration  *prometheus.HistogramVec
	responseSize     *prometheus.HistogramVec
	requestsInFlight prometheus.Gauge
}

func newCollectors(o *options) *collectors {
	labels := []string{labelMethod, labelRoute, labelStatus}

	return &collectors{
		requestsTotal: prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: o.namespace,
			Subsystem: o.subsystem,
			Name:      "requests_total",
			Help:      "Total number of HTTP requests processed, partitioned by method, route and status code.",
		}, labels),
		requestDuration: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Namespace: o.namespace,
			Subsystem: o.subsystem,
			Name:      "request_duration_seconds",
			Help:      "HTTP request latency in seconds, partitioned by method, route and status code.",
			Buckets:   o.durationBuckets,
		}, labels),
		responseSize: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Namespace: o.namespace,
			Subsystem: o.subsystem,
			Name:      "response_size_bytes",
			Help:      "HTTP response size in bytes, partitioned by method, route and status code.",
			Buckets:   o.sizeBuckets,
		}, labels),
		requestsInFlight: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: o.namespace,
			Subsystem: o.subsystem,
			Name:      "requests_in_flight",
			Help:      "Current number of HTTP requests being served.",
		}),
	}
}

func (c *collectors) register(registry prometheus.Registerer) error {
	for _, collector := range []prometheus.Collector{
		c.requestsTotal,
		c.requestDuration,
		c.responseSize,
		c.requestsInFlight,
	} {
		if err := registry.Register(collector); err != nil {
			return err
		}
	}

	return nil
}

func (c *collectors) registerGoCollectors(registry prometheus.Registerer) error {
	for _, collector := range []prometheus.Collector{
		promcollectors.NewGoCollector(),
		promcollectors.NewProcessCollector(promcollectors.ProcessCollectorOpts{}),
	} {
		if err := registry.Register(collector); err != nil {
			return err
		}
	}

	return nil
}
