package config

import (
	"time"
)

type JWTOptions struct {
	Issuer                     string        `env:"ISSUER"`
	Audience                   string        `env:"AUDIENCE"`
	SecretKey                  string        `env:"SECRET,notEmpty"`
	AccessTokenExpirationTime  time.Duration `env:"ACCESS_TOKEN_EXPIRATION_TIME"`
	RefreshTokenExpirationTime time.Duration `env:"REFRESH_TOKEN_EXPIRATION_TIME"`
	Leeway                     time.Duration `env:"LEEWAY"`
}
