package config

import (
	"time"
)

type MCPOptions struct {
	Server MCPServerOptions `envPrefix:"SERVER_"`
	OAuth  MCPOAuthOptions  `envPrefix:"OAUTH_"`
}

type MCPServerOptions struct {
	URLTemplate string `env:"URL_TEMPLATE"`
}

type MCPOAuthOptions struct {
	IssuerURL               string        `env:"ISSUER_URL"`
	ResourceURL             string        `env:"RESOURCE_URL"`
	ConsentURLTemplate      string        `env:"CONSENT_URL_TEMPLATE"`
	AccessTokenTTL          time.Duration `env:"ACCESS_TOKEN_TTL"`
	RefreshTokenTTL         time.Duration `env:"REFRESH_TOKEN_TTL"`
	AuthorizationCodeTTL    time.Duration `env:"AUTHORIZATION_CODE_TTL"`
	AuthorizationRequestTTL time.Duration `env:"AUTHORIZATION_REQUEST_TTL"`
}
