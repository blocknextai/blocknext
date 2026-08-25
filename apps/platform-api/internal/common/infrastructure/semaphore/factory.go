package semaphore

import (
	"github.com/blocknextai/go-packages/apperror"
	"github.com/blocknextai/go-packages/redisclient"
	commonDomainSemaphore "github.com/blocknextai/platform-api/internal/common/domain/semaphore"
	"github.com/blocknextai/platform-api/internal/common/infrastructure/semaphore/memory"
	"github.com/blocknextai/platform-api/internal/common/infrastructure/semaphore/redis"
	"github.com/blocknextai/platform-api/internal/config"
)

var (
	ErrInvalidSemaphoreProvider = apperror.Internal("invalid semaphore provider")
)

func New(semaphoreOptions config.SemaphoreOptions) (commonDomainSemaphore.SemaphoreManager, error) {
	if semaphoreOptions.Provider == config.SemaphoreProviderMemory {
		return memory.New(semaphoreOptions.TTL), nil
	}

	if semaphoreOptions.Provider == config.SemaphoreProviderRedis {
		return redis.NewRedisSemaphore(
			semaphoreOptions.Redis.Address,
			semaphoreOptions.Redis.Password,
			semaphoreOptions.Redis.DB,
			redisclient.PoolOptions{
				PoolSize:        semaphoreOptions.Redis.PoolSize,
				MinIdleConns:    semaphoreOptions.Redis.MinIdleConns,
				MaxIdleConns:    semaphoreOptions.Redis.MaxIdleConns,
				PoolTimeout:     semaphoreOptions.Redis.PoolTimeout,
				ConnMaxIdleTime: semaphoreOptions.Redis.ConnMaxIdleTime,
				ConnMaxLifetime: semaphoreOptions.Redis.ConnMaxLifetime,
			},
			semaphoreOptions.TTL,
		)
	}

	return nil, ErrInvalidSemaphoreProvider
}
