package config

import (
	"github.com/caarlos0/env/v11"
)

type TaskWorkerConfig struct {
	Shared *SharedConfig `env:"-"`

	TaskRunner TaskRunnerOptions `envPrefix:"TASK_RUNNER_"`
}

func LoadTaskWorker() (*TaskWorkerConfig, error) {
	cfg := &TaskWorkerConfig{}
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
