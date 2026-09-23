package config

import (
	"time"
)

type MCPOptions struct {
	MaxExecutionTime time.Duration    `env:"MAX_EXECUTION_TIME"`
	Server           MCPServerOptions `envPrefix:"SERVER_"`
	OAuth            MCPOAuthOptions  `envPrefix:"OAUTH_"`
}

type MCPServerOptions struct {
	URLTemplate string `env:"URL_TEMPLATE,notEmpty"`
}

type MCPOAuthOptions struct {
	IssuerURL               string        `env:"ISSUER_URL,notEmpty"`
	ResourceURL             string        `env:"RESOURCE_URL,notEmpty"`
	ConsentURLTemplate      string        `env:"CONSENT_URL_TEMPLATE,notEmpty"`
	SigningKey              string        `env:"SIGNING_KEY,notEmpty"`
	AccessTokenTTL          time.Duration `env:"ACCESS_TOKEN_TTL"`
	RefreshTokenTTL         time.Duration `env:"REFRESH_TOKEN_TTL"`
	AuthorizationCodeTTL    time.Duration `env:"AUTHORIZATION_CODE_TTL"`
	AuthorizationRequestTTL time.Duration `env:"AUTHORIZATION_REQUEST_TTL"`
}
