package config

import (
	"time"
)

type AuthOptions struct {
	Google    AuthGoogleOptions    `envPrefix:"GOOGLE_"`
	Github    AuthGithubOptions    `envPrefix:"GITHUB_"`
	X         AuthXOptions         `envPrefix:"X_"`
	Facebook  AuthFacebookOptions  `envPrefix:"FACEBOOK_"`
	Password  AuthPasswordOptions  `envPrefix:"PASSWORD_"`
	MagicLink AuthMagicLinkOptions `envPrefix:"MAGIC_LINK_"`
}

type AuthGoogleOptions struct {
	Enabled      bool   `env:"ENABLED"`
	ClientID     string `env:"CLIENT_ID"`
	ClientSecret string `env:"CLIENT_SECRET"`
	RedirectURI  string `env:"REDIRECT_URI"`
}

type AuthGithubOptions struct {
	Enabled      bool   `env:"ENABLED"`
	ClientID     string `env:"CLIENT_ID"`
	ClientSecret string `env:"CLIENT_SECRET"`
	RedirectURI  string `env:"REDIRECT_URI"`
}

type AuthXOptions struct {
	Enabled      bool   `env:"ENABLED"`
	ClientID     string `env:"CLIENT_ID"`
	ClientSecret string `env:"CLIENT_SECRET"`
	RedirectURI  string `env:"REDIRECT_URI"`
}

type AuthFacebookOptions struct {
	Enabled      bool   `env:"ENABLED"`
	ClientID     string `env:"CLIENT_ID"`
	ClientSecret string `env:"CLIENT_SECRET"`
	RedirectURI  string `env:"REDIRECT_URI"`
}

type AuthPasswordOptions struct {
	Enabled               bool          `env:"ENABLED"`
	BcryptCost            int           `env:"BCRYPT_COST"`
	VerificationTokenTTL  time.Duration `env:"VERIFICATION_TOKEN_TTL"`
	PasswordResetTokenTTL time.Duration `env:"PASSWORD_RESET_TOKEN_TTL"`
}

type AuthMagicLinkOptions struct {
	Enabled  bool          `env:"ENABLED"`
	TokenTTL time.Duration `env:"TOKEN_TTL"`
}
