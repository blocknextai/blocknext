package config

type MCPOptions struct {
	Server MCPServerOptions `envPrefix:"SERVER_"`
}

type MCPServerOptions struct {
	URLTemplate string `env:"URL_TEMPLATE"`
}
