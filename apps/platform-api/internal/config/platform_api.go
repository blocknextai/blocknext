package config

import (
	"github.com/caarlos0/env/v11"
)

type PlatformAPIConfig struct {
	Shared *SharedConfig `env:"-"`

	HTTPServer HTTPServerOptions `envPrefix:"PLATFORM_API_"`
	WebSocket  WebSocketOptions  `envPrefix:"WEBSOCKET_"`
	TaskRunner TaskRunnerOptions `envPrefix:"TASK_RUNNER_"`
}

func LoadPlatformAPI() (*PlatformAPIConfig, error) {
	cfg := &PlatformAPIConfig{}
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}
	shared, err := LoadShared()
	if err != nil {
		return nil, err
	}
	cfg.Shared = shared
	return cfg, nil
}
