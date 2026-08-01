package config

import (
	"github.com/caarlos0/env/v11"
)

type MCPAPIConfig struct {
	Shared *SharedConfig `env:"-"`

	HTTPServer HTTPServerOptions `envPrefix:"MCP_API_"`
	MCP        MCPOptions        `envPrefix:"MCP_"`
}

func LoadMCPAPI() (*MCPAPIConfig, error) {
	cfg := &MCPAPIConfig{}
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
