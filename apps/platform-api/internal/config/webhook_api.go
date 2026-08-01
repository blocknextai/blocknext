package config

import (
	"github.com/caarlos0/env/v11"
)

type WebhookAPIConfig struct {
	Shared *SharedConfig `env:"-"`

	HTTPServer HTTPServerOptions `envPrefix:"WEBHOOK_API_"`
	TaskRunner TaskRunnerOptions `envPrefix:"TASK_RUNNER_"`
}

func LoadWebhookAPI() (*WebhookAPIConfig, error) {
	cfg := &WebhookAPIConfig{}
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
