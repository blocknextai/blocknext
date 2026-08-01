package config

import (
	"time"
)

type EventBusOptions struct {
	PollInterval    time.Duration `env:"POLL_INTERVAL"`
	BatchSize       int           `env:"BATCH_SIZE"`
	MaxAttempts     int           `env:"MAX_ATTEMPTS"`
	BaseBackoff     time.Duration `env:"BASE_BACKOFF"`
	StuckTimeout    time.Duration `env:"STUCK_TIMEOUT"`
	ReclaimInterval time.Duration `env:"RECLAIM_INTERVAL"`
}
