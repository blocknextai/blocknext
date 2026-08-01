// Package cache provides a Fiber storage adapter backed by a cache.Service,
// implementing the fiber.Storage interface with an optional key prefix.
package cache

import (
	"context"
	"strings"
	"time"

	"github.com/blocknextai/go-packages/cache"
)

// Storage adapts a cache.Service to the Fiber storage interface using a key
// prefix.
type Storage struct {
	service cache.Service
	prefix  string
}

// New creates a Storage backed by service, prefixing every key with prefix.
func New(service cache.Service, prefix string) *Storage {
	return &Storage{
		service: service,
		prefix:  prefix,
	}
}

// Get returns the value stored for key, or nil if the key is empty or absent.
func (s *Storage) Get(key string) ([]byte, error) {
	return s.GetWithContext(context.Background(), key)
}

// GetWithContext returns the value stored for key using the given context, or
// nil if the key is empty or absent.
func (s *Storage) GetWithContext(ctx context.Context, key string) ([]byte, error) {
	if strings.TrimSpace(key) == "" {
		return nil, nil
	}

	value, err := s.service.Get(ctx, s.fullKey(key))
	if err != nil {
		return nil, err
	}
	if value == "" {
		return nil, nil
	}
	return []byte(value), nil
}

// Set stores val for key with the given expiration. Empty keys or values are
// ignored.
func (s *Storage) Set(key string, val []byte, exp time.Duration) error {
	return s.SetWithContext(context.Background(), key, val, exp)
}

// SetWithContext stores val for key with the given expiration using the
// provided context. Empty keys or values are ignored.
func (s *Storage) SetWithContext(ctx context.Context, key string, val []byte, exp time.Duration) error {
	if strings.TrimSpace(key) == "" || len(val) == 0 {
		return nil
	}
	return s.service.Set(ctx, s.fullKey(key), string(val), exp)
}

// Delete removes the value stored for key. An empty key is ignored.
func (s *Storage) Delete(key string) error {
	return s.DeleteWithContext(context.Background(), key)
}

// DeleteWithContext removes the value stored for key using the given context.
// An empty key is ignored.
func (s *Storage) DeleteWithContext(ctx context.Context, key string) error {
	if strings.TrimSpace(key) == "" {
		return nil
	}
	return s.service.Delete(ctx, s.fullKey(key))
}

// Reset is a no-op; the shared backend is not flushed by this adapter.
func (s *Storage) Reset() error {
	return s.ResetWithContext(context.Background())
}

// ResetWithContext is a no-op; the shared backend is not flushed by this
// adapter.
func (s *Storage) ResetWithContext(_ context.Context) error {
	return nil
}

// Close is a no-op; the underlying service lifecycle is managed elsewhere.
func (s *Storage) Close() error {
	return nil
}

func (s *Storage) fullKey(key string) string {
	if s.prefix == "" {
		return key
	}
	var builder strings.Builder
	builder.Grow(len(s.prefix) + len(key))
	builder.WriteString(s.prefix)
	builder.WriteString(key)
	return builder.String()
}
