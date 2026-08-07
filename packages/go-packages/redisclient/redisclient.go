// Package redisclient builds the go-redis client that every Redis-backed
// provider in this project shares, so pool configuration and the connectivity
// check are defined in one place.
package redisclient

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

// pingTimeout bounds the connectivity check so an unreachable server fails the
// caller's startup instead of hanging it.
const pingTimeout = 5 * time.Second

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

// New connects to the Redis server at addr using the given password, database
// index and pool options. It pings the server to verify connectivity and
// returns an error if the connection cannot be established.
func New(addr string, password string, db int, poolOptions PoolOptions) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:            addr,
		Password:        password,
		DB:              db,
		PoolSize:        poolOptions.PoolSize,
		MinIdleConns:    poolOptions.MinIdleConns,
		MaxIdleConns:    poolOptions.MaxIdleConns,
		PoolTimeout:     poolOptions.PoolTimeout,
		ConnMaxIdleTime: poolOptions.ConnMaxIdleTime,
		ConnMaxLifetime: poolOptions.ConnMaxLifetime,
	})

	ctx, cancel := context.WithTimeout(context.Background(), pingTimeout)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return client, nil
}
