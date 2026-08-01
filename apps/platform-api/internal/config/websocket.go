package config

type WebSocketOptions struct {
	MaxConnectionsPerRoom int `env:"MAX_CONNECTIONS_PER_ROOM"`
}
