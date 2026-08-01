package config

type SecretManagerOptions struct {
	SecretKey string `env:"SECRET_KEY,notEmpty"`
}
