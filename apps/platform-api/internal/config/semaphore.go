package config

import (
	"time"
)

type SemaphoreProvider string

const (
	SemaphoreProviderMemory SemaphoreProvider = "memory"
	SemaphoreProviderRedis  SemaphoreProvider = "redis"
)

type SemaphoreOptions struct {
	Provider                SemaphoreProvider `env:"PROVIDER"`
	TTL                     time.Duration     `env:"TTL"`
	HeartbeatInterval       time.Duration     `env:"HEARTBEAT_INTERVAL"`
	MaxConcurrentExecutions int64             `env:"MAX_CONCURRENT_EXECUTIONS"`

	Redis SemaphoreRedisOptions `envPrefix:"REDIS_"`
}

type SemaphoreRedisOptions struct {
	Address         string        `env:"ADDRESS"`
	Password        string        `env:"PASSWORD"`
	DB              int           `env:"DB"`
	PoolSize        int           `env:"POOL_SIZE"`
	MinIdleConns    int           `env:"MIN_IDLE_CONNS"`
	MaxIdleConns    int           `env:"MAX_IDLE_CONNS"`
	PoolTimeout     time.Duration `env:"POOL_TIMEOUT"`
	ConnMaxIdleTime time.Duration `env:"CONN_MAX_IDLE_TIME"`
	ConnMaxLifetime time.Duration `env:"CONN_MAX_LIFETIME"`
}
