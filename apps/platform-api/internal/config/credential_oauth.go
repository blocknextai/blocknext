package config

import (
	"time"
)

type CredentialOAuthOptions struct {
	OAuth2RedirectURL string        `env:"OAUTH2_REDIRECT_URL"`
	StateTTL          time.Duration `env:"STATE_TTL"`
}
