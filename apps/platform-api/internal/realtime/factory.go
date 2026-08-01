package realtime

import (
	"github.com/blocknextai/go-packages/apperror"
	"github.com/blocknextai/platform-api/internal/config"
	realtimeRedis "github.com/blocknextai/platform-api/internal/realtime/redis"
)

var (
	ErrInvalidBrokerType = apperror.Internal("invalid broker type")
)

func New(options config.BrokerOptions) (Broadcaster, error) {
	if options.Type == config.BrokerTypeRedis {
		return realtimeRedis.New(
			options.Redis.Address,
			options.Redis.Password,
			options.Redis.DB,
			realtimeRedis.PoolOptions{
				PoolSize:        options.Redis.PoolSize,
				MinIdleConns:    options.Redis.MinIdleConns,
				MaxIdleConns:    options.Redis.MaxIdleConns,
				PoolTimeout:     options.Redis.PoolTimeout,
				ConnMaxIdleTime: options.Redis.ConnMaxIdleTime,
				ConnMaxLifetime: options.Redis.ConnMaxLifetime,
			},
		)
	}

	return nil, ErrInvalidBrokerType
}
