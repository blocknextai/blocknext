package config

import (
	"net"
	"net/url"
	"strings"
	"time"

	"github.com/blocknextai/go-packages/cast"
)

type DatabaseOptions struct {
	Host            string        `env:"HOST,notEmpty"`
	Port            int           `env:"PORT"`
	User            string        `env:"USER,notEmpty"`
	Password        string        `env:"PASSWORD,notEmpty"`
	Name            string        `env:"NAME,notEmpty"`
	SSLMode         string        `env:"SSL_MODE"`
	MaxOpenConns    int           `env:"MAX_OPEN_CONNS"`
	MaxIdleConns    int           `env:"MAX_IDLE_CONNS"`
	ConnMaxLifetime time.Duration `env:"CONN_MAX_LIFETIME"`
	ConnMaxIdleTime time.Duration `env:"CONN_MAX_IDLE_TIME"`
	QueryTimeout    time.Duration `env:"QUERY_TIMEOUT"`
}

func (o DatabaseOptions) DSN() string {
	u := url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(o.User, o.Password),
		Host:   net.JoinHostPort(o.Host, cast.ToString(o.Port)),
		Path:   o.Name,
	}
	if strings.TrimSpace(o.SSLMode) != "" {
		q := u.Query()
		q.Set("sslmode", o.SSLMode)
		u.RawQuery = q.Encode()
	}
	return u.String()
}
