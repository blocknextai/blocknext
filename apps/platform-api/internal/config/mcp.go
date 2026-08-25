package config

import (
	"time"
)

type MCPOptions struct {
	MaxExecutionTime time.Duration `env:"MAX_EXECUTION_TIME"`

	Server MCPServerOptions `envPrefix:"SERVER_"`
}

type MCPServerOptions struct {
	URLTemplate string `env:"URL_TEMPLATE"`
}
