package config

import (
	"time"
)

type DownloadOptions struct {
	MaxSize int64         `env:"MAX_SIZE"`
	Timeout time.Duration `env:"TIMEOUT"`
}
