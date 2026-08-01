// Package cache defines the key-value store abstraction used across the
// project, including counters, expiry and atomic semaphore operations.
package cache

import (
	"context"
	"time"
)

// Service is a key-value cache backend supporting counters, expiry and atomic
// semaphore operations.
type Service interface {
	// Ping verifies connectivity to the backend.
	Ping(ctx context.Context) error
	// Get returns the value for key, or an empty string if it is absent.
	Get(ctx context.Context, key string) (string, error)
	// Set stores value for key with the given expiration.
	Set(ctx context.Context, key string, value string, expiration time.Duration) error
	// Incr atomically increments the integer value of key and returns the result.
	Incr(ctx context.Context, key string) (int64, error)
	// Decr atomically decrements the integer value of key and returns the result.
	Decr(ctx context.Context, key string) (int64, error)
	// Delete removes key.
	Delete(ctx context.Context, key string) error
	// GetAndDelete returns the value for key and removes it atomically.
	GetAndDelete(ctx context.Context, key string) (string, error)
	// Exists reports whether key is present.
	Exists(ctx context.Context, key string) (bool, error)
	// Keys returns the keys matching pattern.
	Keys(ctx context.Context, pattern string) ([]string, error)
	// Expire sets the expiration for key.
	Expire(ctx context.Context, key string, expiration time.Duration) error
	// AcquireSemaphoreAtomic atomically acquires a slot on the semaphore at key,
	// reporting whether it was acquired given the maximum count and TTL.
	AcquireSemaphoreAtomic(ctx context.Context, key string, maxCount int64, ttl time.Duration) (bool, error)
	// ReleaseSemaphoreAtomic atomically releases a slot on the semaphore at key
	// and returns the remaining count.
	ReleaseSemaphoreAtomic(ctx context.Context, key string) (int64, error)
	// Close releases the resources held by the backend.
	Close() error
}
