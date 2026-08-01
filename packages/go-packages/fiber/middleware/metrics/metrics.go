package metrics

import (
	"errors"
	"log/slog"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/adaptor"
	"github.com/gofiber/fiber/v3/middleware/healthcheck"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// ErrRegisterCollectors is returned when the Prometheus collectors cannot be
// registered with the configured registry.
var (
	ErrRegisterCollectors = errors.New("metrics: failed to register collectors")
)

// Middleware records Prometheus HTTP metrics and exposes them for scraping.
type Middleware struct {
	options    *options
	collectors *collectors
}

// New creates a Middleware configured with the given options and registers its
// collectors. It returns an error joined with ErrRegisterCollectors if
// registration fails.
func New(opts ...Option) (*Middleware, error) {
	o := defaultOptions()
	for _, opt := range opts {
		opt(o)
	}

	o.skipPaths[o.metricsPath] = struct{}{}
	o.skipPaths[healthcheck.LivenessEndpoint] = struct{}{}
	o.skipPaths[healthcheck.ReadinessEndpoint] = struct{}{}

	c := newCollectors(o)
	if err := c.register(o.registry); err != nil {
		slog.Error("Failed to register metrics collectors",
			"component", "metrics_middleware",
			"error", err,
		)
		return nil, errors.Join(ErrRegisterCollectors, err)
	}

	if o.goCollectors {
		if err := c.registerGoCollectors(o.registry); err != nil {
			slog.Error("Failed to register Go runtime collectors",
				"component", "metrics_middleware",
				"error", err,
			)
			return nil, errors.Join(ErrRegisterCollectors, err)
		}
	}

	return &Middleware{
		options:    o,
		collectors: c,
	}, nil
}

// Collect returns a Fiber handler that records request count, duration and
// response size metrics for each non-skipped request.
func (m *Middleware) Collect() fiber.Handler {
	col := m.collectors

	return func(c fiber.Ctx) error {
		if m.shouldSkip(c.Path()) {
			return c.Next()
		}

		col.requestsInFlight.Inc()
		defer col.requestsInFlight.Dec()

		start := time.Now()
		err := c.Next()

		resp := c.Response()
		method := c.Method()
		route := c.Route().Path
		if route == "" {
			route = c.Path()
		}
		status := strconv.Itoa(resp.StatusCode())

		size := resp.Header.ContentLength()
		if size < 0 {
			size = len(resp.Body())
		}

		col.requestsTotal.WithLabelValues(method, route, status).Inc()
		col.requestDuration.WithLabelValues(method, route, status).Observe(time.Since(start).Seconds())
		col.responseSize.WithLabelValues(method, route, status).Observe(float64(size))

		return err
	}
}

// Expose returns a Fiber handler that serves the collected metrics in the
// Prometheus text exposition format.
func (m *Middleware) Expose() fiber.Handler {
	return adaptor.HTTPHandler(promhttp.HandlerFor(m.options.registry, promhttp.HandlerOpts{
		Registry: m.options.registry,
	}))
}

func (m *Middleware) shouldSkip(path string) bool {
	_, skip := m.options.skipPaths[path]
	return skip
}
