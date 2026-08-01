package infrastructure

import (
	"github.com/blocknextai/file-gateway-api/internal/config"
	"github.com/blocknextai/go-packages/apperror"
	"github.com/blocknextai/go-packages/cache"
	"github.com/blocknextai/go-packages/cache/redis"
)

var (
	ErrInvalidCacheType = apperror.Internal("invalid cache type")
)

func NewCacheService(options config.CacheOptions) (cache.Service, error) {
	if options.Type == config.CacheTypeRedis {
		return redis.New(
			options.Redis.Address,
			options.Redis.Password,
			options.Redis.DB,
			redis.PoolOptions{
				PoolSize:        options.Redis.PoolSize,
				MinIdleConns:    options.Redis.MinIdleConns,
				MaxIdleConns:    options.Redis.MaxIdleConns,
				PoolTimeout:     options.Redis.PoolTimeout,
				ConnMaxIdleTime: options.Redis.ConnMaxIdleTime,
				ConnMaxLifetime: options.Redis.ConnMaxLifetime,
			},
		)
	}

	return nil, ErrInvalidCacheType
}
