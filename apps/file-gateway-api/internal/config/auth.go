package config

import (
	"time"
)

type AuthOptions struct {
	ServiceKey string         `env:"SERVICE_KEY,notEmpty"`
	JWT        AuthJWTOptions `envPrefix:"JWT_"`
}

type AuthJWTOptions struct {
	Secret              string        `env:"SECRET,notEmpty"`
	Issuer              string        `env:"ISSUER,notEmpty"`
	Audience            string        `env:"AUDIENCE,notEmpty"`
	TokenExpirationTime time.Duration `env:"TOKEN_EXPIRATION_TIME"`
	TokenLeeway         time.Duration `env:"TOKEN_LEEWAY"`
}
