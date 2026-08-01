package config

import (
	"net"
	"time"
)

type HTTPServerOptions struct {
	Host              string        `env:"HOST"`
	Port              string        `env:"PORT"`
	IsPrefork         bool          `env:"IS_PREFORK"`
	ReadTimeout       time.Duration `env:"READ_TIMEOUT"`
	WriteTimeout      time.Duration `env:"WRITE_TIMEOUT"`
	IdleTimeout       time.Duration `env:"IDLE_TIMEOUT"`
	ReduceMemoryUsage bool          `env:"REDUCE_MEMORY_USAGE"`
	BodyLimit         int           `env:"BODY_LIMIT"`
	ReadBufferSize    int           `env:"READ_BUFFER_SIZE"`
	WriteBufferSize   int           `env:"WRITE_BUFFER_SIZE"`
	TrustProxyEnabled bool          `env:"TRUST_PROXY_ENABLED"`
	TrustProxies      string        `env:"TRUST_PROXIES"`
	ProxyHeader       string        `env:"PROXY_HEADER"`
	AllowOrigins      string        `env:"ALLOW_ORIGINS"`
	AllowHeaders      string        `env:"ALLOW_HEADERS"`
	AllowMethods      string        `env:"ALLOW_METHODS"`
	ExposeHeaders     string        `env:"EXPOSE_HEADERS"`
	CORSMaxAge        int           `env:"CORS_MAX_AGE"`
	MaxRequests       int           `env:"MAX_REQUESTS"`
	ExpirationTime    time.Duration `env:"EXPIRATION_TIME"`
	//Metrics           MetricsOptions `envPrefix:"METRICS_"`
}

func (s HTTPServerOptions) Address() string {
	return net.JoinHostPort(s.Host, s.Port)
}

/*func (s HTTPServerOptions) MetricsAddress() string {
	return net.JoinHostPort(s.Host, s.Metrics.Port)
}*/
