package config

import (
	"github.com/caarlos0/env/v11"
)

type ConfigurationOptions struct {
	AppEnv   AppEnv          `env:"APP_ENV"`
	API      APIOptions      `envPrefix:"FILE_GATEWAY_API_"`
	Cache    CacheOptions    `envPrefix:"CACHE_"`
	Auth     AuthOptions     `envPrefix:"FILE_GATEWAY_AUTH_"`
	Storage  StorageOptions  `envPrefix:"FILE_GATEWAY_STORAGE_"`
	Download DownloadOptions `envPrefix:"FILE_GATEWAY_DOWNLOAD_"`
}

func Load() (*ConfigurationOptions, error) {
	cfg := &ConfigurationOptions{}
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
