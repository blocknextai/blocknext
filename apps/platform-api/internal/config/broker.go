package config

import (
	"time"
)

type BrokerType string

const (
	BrokerTypeMemory BrokerType = "memory"
	BrokerTypeRedis  BrokerType = "redis"
)

type BrokerOptions struct {
	Type  BrokerType         `env:"TYPE"`
	Redis BrokerRedisOptions `envPrefix:"REDIS_"`
}

type BrokerRedisOptions struct {
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
