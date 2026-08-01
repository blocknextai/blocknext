package metrics

import "github.com/prometheus/client_golang/prometheus"

const (
	defaultSubsystem   = "http"
	defaultMetricsPath = "/metrics"
)

type options struct {
	registry        *prometheus.Registry
	namespace       string
	subsystem       string
	metricsPath     string
	skipPaths       map[string]struct{}
	durationBuckets []float64
	sizeBuckets     []float64
	goCollectors    bool
}

// Option configures the metrics middleware.
type Option func(*options)

func defaultOptions() *options {
	return &options{
		registry:        prometheus.NewRegistry(),
		subsystem:       defaultSubsystem,
		metricsPath:     defaultMetricsPath,
		skipPaths:       make(map[string]struct{}),
		durationBuckets: prometheus.DefBuckets,
		sizeBuckets:     prometheus.ExponentialBuckets(100, 10, 8),
	}
}

// WithRegistry sets the Prometheus registry used to register and expose
// metrics. A nil registry is ignored.
func WithRegistry(registry *prometheus.Registry) Option {
	return func(o *options) {
		if registry != nil {
			o.registry = registry
		}
	}
}

// WithNamespace sets the namespace prefix applied to all metric names.
func WithNamespace(namespace string) Option {
	return func(o *options) {
		o.namespace = namespace
	}
}

// WithSubsystem sets the subsystem prefix applied to all metric names. An empty
// subsystem is ignored.
func WithSubsystem(subsystem string) Option {
	return func(o *options) {
		if subsystem != "" {
			o.subsystem = subsystem
		}
	}
}

// WithMetricsPath sets the path used to expose metrics, which is also excluded
// from collection. An empty path is ignored.
func WithMetricsPath(path string) Option {
	return func(o *options) {
		if path != "" {
			o.metricsPath = path
		}
	}
}

// WithSkipPaths adds paths that are excluded from metrics collection. Empty
// paths are ignored.
func WithSkipPaths(paths ...string) Option {
	return func(o *options) {
		for _, path := range paths {
			if path != "" {
				o.skipPaths[path] = struct{}{}
			}
		}
	}
}

// WithDurationBuckets sets the histogram buckets for request duration in
// seconds. An empty slice is ignored.
func WithDurationBuckets(buckets []float64) Option {
	return func(o *options) {
		if len(buckets) > 0 {
			o.durationBuckets = buckets
		}
	}
}

// WithSizeBuckets sets the histogram buckets for response size in bytes. An
// empty slice is ignored.
func WithSizeBuckets(buckets []float64) Option {
	return func(o *options) {
		if len(buckets) > 0 {
			o.sizeBuckets = buckets
		}
	}
}

// WithGoCollectors enables the Go runtime and process Prometheus collectors.
func WithGoCollectors() Option {
	return func(o *options) {
		o.goCollectors = true
	}
}
