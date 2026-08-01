package config

type FileGatewayOptions struct {
	BaseURL string                 `env:"BASE_URL"`
	Auth    FileGatewayAuthOptions `envPrefix:"AUTH_"`
}

type FileGatewayAuthOptions struct {
	ServiceKey string `env:"SERVICE_KEY"`
}
