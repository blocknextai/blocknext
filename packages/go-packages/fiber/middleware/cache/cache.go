// Package cache provides a Fiber middleware that caches successful GET
// responses in a cache.Service backend keyed by request URL.
package cache

import (
	"crypto/sha256"
	"math/rand"
	"strings"
	"time"

	"github.com/blocknextai/go-packages/cache"
	"github.com/blocknextai/go-packages/hex"
	"github.com/gofiber/fiber/v3"
)

const (
	cacheHeader = "X-Cache"
	cacheHit    = "HIT"
	cacheMiss   = "MISS"
)

// Middleware caches HTTP responses using a cache.Service backend and a key
// prefix.
type Middleware struct {
	service cache.Service
	prefix  string
}

// New creates a Middleware that stores cached responses in service using the
// given key prefix.
func New(service cache.Service, prefix string) *Middleware {
	return &Middleware{
		service: service,
		prefix:  prefix,
	}
}

// Cache returns a Fiber handler that serves cached responses for GET requests
// and stores successful (2xx) responses for the given TTL. A non-positive TTL
// disables caching for the request.
func (m *Middleware) Cache(ttl time.Duration) fiber.Handler {
	return func(c fiber.Ctx) error {
		if c.Method() != fiber.MethodGet || ttl <= 0 {
			return c.Next()
		}

		key := m.buildCacheKey(c)
		ctx := c.RequestCtx()

		if cachedData, err := m.service.Get(ctx, key); err == nil && cachedData != "" {
			c.Set(cacheHeader, cacheHit)
			c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
			return c.SendString(cachedData)
		}

		if err := c.Next(); err != nil {
			return err
		}

		statusCode := c.Response().StatusCode()
		if statusCode >= 200 && statusCode < 300 {
			body := c.Response().Body()
			if len(body) > 0 {
				jitteredTTL := m.addJitter(ttl)
				_ = m.service.Set(ctx, key, string(body), jitteredTTL)
			}
		}

		c.Set(cacheHeader, cacheMiss)
		return nil
	}
}

func (m *Middleware) buildCacheKey(c fiber.Ctx) string {
	h := sha256.New()
	h.Write([]byte(c.OriginalURL()))
	hash := hex.Encode(h.Sum(nil))

	var builder strings.Builder
	builder.Grow(len(m.prefix) + len(c.Method()) + 1 + len(hash))
	builder.WriteString(m.prefix)
	builder.WriteString(c.Method())
	builder.WriteByte(':')
	builder.WriteString(hash)
	return builder.String()
}

func (m *Middleware) addJitter(ttl time.Duration) time.Duration {
	jitterPercent := 10.0
	jitterRange := float64(ttl) * (jitterPercent / 100.0)
	jitter := time.Duration(rand.Float64() * jitterRange)
	return ttl + jitter
}
