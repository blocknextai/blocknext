// Package logging provides a Fiber middleware that logs HTTP requests with
// their method, path, status, duration and request metadata.
package logging

import (
	"log/slog"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/healthcheck"
)

// New returns a Fiber handler that logs each HTTP request at info level,
// skipping the health check liveness and readiness endpoints.
func New() fiber.Handler {
	return func(c fiber.Ctx) error {
		path := c.Path()
		if path == healthcheck.LivenessEndpoint || path == healthcheck.ReadinessEndpoint {
			return c.Next()
		}

		startUTC := time.Now().UTC()
		err := c.Next()

		slog.InfoContext(c.RequestCtx(), "HTTP request",
			"component", "http",
			"method", c.Method(),
			"path", path,
			"status", c.Response().StatusCode(),
			"duration_ms", time.Since(startUTC).Milliseconds(),
			"ip", c.IP(),
			"request_id", c.Get(fiber.HeaderXRequestID),
		)

		return err
	}
}
