package config

import (
	"time"
)

type TaskRunnerMode string

const (
	TaskRunnerModeEmbedded TaskRunnerMode = "embedded"
	TaskRunnerModeQueue    TaskRunnerMode = "queue"
)

type TaskQueueProvider string

const (
	TaskQueueProviderRedis TaskQueueProvider = "redis"
)

type TaskLockProvider string

const (
	TaskLockProviderMemory TaskLockProvider = "memory"
	TaskLockProviderRedis  TaskLockProvider = "redis"
)

type TaskSemaphoreProvider string

const (
	TaskSemaphoreProviderMemory TaskSemaphoreProvider = "memory"
	TaskSemaphoreProviderRedis  TaskSemaphoreProvider = "redis"
)

type TaskRunnerWorkerOptions struct {
	PoolSize      int `env:"POOL_SIZE"`
	PoolQueueSize int `env:"POOL_QUEUE_SIZE"`
}

type TaskRunnerQueueOptions struct {
	Provider          TaskQueueProvider `env:"PROVIDER"`
	IdleTimeout       time.Duration     `env:"IDLE_TIMEOUT"`
	StaleClaimTimeout time.Duration     `env:"STALE_CLAIM_TIMEOUT"`
	MaxRetries        int               `env:"MAX_RETRIES"`

	Redis TaskQueueRedisOptions `envPrefix:"REDIS_"`
}

type TaskQueueRedisOptions struct {
	StreamName      string        `env:"STREAM_NAME"`
	ConsumerGroup   string        `env:"CONSUMER_GROUP"`
	ConsumerName    string        `env:"CONSUMER_NAME"`
	BlockTimeout    time.Duration `env:"BLOCK_TIMEOUT"`
	PrefetchCount   int           `env:"PREFETCH_COUNT"`
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

type TaskRunnerLeaderOptions struct {
	Provider      TaskLockProvider `env:"PROVIDER"`
	Key           string           `env:"KEY"`
	TTL           time.Duration    `env:"TTL"`
	PollInterval  time.Duration    `env:"POLL_INTERVAL"`
	RenewInterval time.Duration    `env:"RENEW_INTERVAL"`

	Redis TaskLockRedisOptions `envPrefix:"REDIS_"`
}

type TaskLockRedisOptions struct {
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

type TaskRunnerSemaphoreOptions struct {
	Provider TaskSemaphoreProvider `env:"PROVIDER"`
	TTL      time.Duration         `env:"TTL"`

	Redis TaskSemaphoreRedisOptions `envPrefix:"REDIS_"`
}

type TaskSemaphoreRedisOptions struct {
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

type TaskRunnerOptions struct {
	Mode TaskRunnerMode `env:"MODE"`

	ShutdownTimeout    time.Duration `env:"SHUTDOWN_TIMEOUT"`
	RecoveryInterval   time.Duration `env:"RECOVERY_INTERVAL"`
	HeartbeatInterval  time.Duration `env:"HEARTBEAT_INTERVAL"`
	MaxConcurrentTasks int64         `env:"MAX_CONCURRENT_TASKS"`

	Worker    TaskRunnerWorkerOptions    `envPrefix:"WORKER_"`
	Queue     TaskRunnerQueueOptions     `envPrefix:"QUEUE_"`
	Leader    TaskRunnerLeaderOptions    `envPrefix:"LEADER_"`
	Semaphore TaskRunnerSemaphoreOptions `envPrefix:"SEMAPHORE_"`
}
