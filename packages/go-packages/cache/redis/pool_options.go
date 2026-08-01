// Package redis provides a Redis-backed implementation of the cache.Service
// interface, including atomic semaphore operations via Lua scripts.
package redis

import (
	"time"
)

// PoolOptions configures the Redis client connection pool.
type PoolOptions struct {
	// PoolSize is the maximum number of socket connections.
	PoolSize int
	// MinIdleConns is the minimum number of idle connections to maintain.
	MinIdleConns int
	// MaxIdleConns is the maximum number of idle connections.
	MaxIdleConns int
	// PoolTimeout is how long a client waits for a free connection.
	PoolTimeout time.Duration
	// ConnMaxIdleTime is the maximum amount of time a connection may be idle.
	ConnMaxIdleTime time.Duration
	// ConnMaxLifetime is the maximum amount of time a connection may be reused.
	ConnMaxLifetime time.Duration
}
