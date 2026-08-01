package redis

import (
	"context"
	"errors"
	"time"

	"github.com/blocknextai/go-packages/cache"
	"github.com/redis/go-redis/v9"
)

type provider struct {
	client                 *redis.Client
	acquireSemaphoreScript *redis.Script
	releaseSemaphoreScript *redis.Script
}

// New connects to the Redis server at addr using the given password, database
// index and pool options, and returns a cache.Service. It pings the server to
// verify connectivity and returns an error if the connection cannot be
// established.
func New(addr string, password string, db int, poolOptions PoolOptions) (cache.Service, error) {
	client := redis.NewClient(&redis.Options{
		Addr:            addr,
		Password:        password,
		DB:              db,
		PoolSize:        poolOptions.PoolSize,
		MinIdleConns:    poolOptions.MinIdleConns,
		MaxIdleConns:    poolOptions.MaxIdleConns,
		PoolTimeout:     poolOptions.PoolTimeout,
		ConnMaxIdleTime: poolOptions.ConnMaxIdleTime,
		ConnMaxLifetime: poolOptions.ConnMaxLifetime,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	acquireScript := redis.NewScript(`
		local key = KEYS[1]
		local maxCount = tonumber(ARGV[1])
		local ttl = tonumber(ARGV[2])

		local current = redis.call('GET', key)
		if not current then
			current = 0
		else
			current = tonumber(current)
		end

		if current < maxCount then
			local newCount = redis.call('INCR', key)
			if newCount == 1 then
				redis.call('EXPIRE', key, ttl)
			end
			return 1
		else
			return 0
		end
	`)

	releaseScript := redis.NewScript(`
		local key = KEYS[1]

		local current = redis.call('GET', key)
		if not current then
			return 0
		end

		local newCount = redis.call('DECR', key)
		if newCount <= 0 then
			redis.call('DEL', key)
			return 0
		end

		return newCount
	`)

	return &provider{
		client:                 client,
		acquireSemaphoreScript: acquireScript,
		releaseSemaphoreScript: releaseScript,
	}, nil
}

func (r *provider) Ping(ctx context.Context) error {
	return r.client.Ping(ctx).Err()
}

func (r *provider) Get(ctx context.Context, key string) (string, error) {
	result, err := r.client.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return "", nil
	}
	return result, err
}

func (r *provider) Set(ctx context.Context, key string, value string, expiration time.Duration) error {
	return r.client.Set(ctx, key, value, expiration).Err()
}

func (r *provider) Incr(ctx context.Context, key string) (int64, error) {
	return r.client.Incr(ctx, key).Result()
}

func (r *provider) Decr(ctx context.Context, key string) (int64, error) {
	return r.client.Decr(ctx, key).Result()
}

func (r *provider) Delete(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}

func (r *provider) GetAndDelete(ctx context.Context, key string) (string, error) {
	result, err := r.client.GetDel(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return "", nil
	}
	return result, err
}

func (r *provider) Exists(ctx context.Context, key string) (bool, error) {
	count, err := r.client.Exists(ctx, key).Result()
	return count > 0, err
}

func (r *provider) Keys(ctx context.Context, pattern string) ([]string, error) {
	return r.client.Keys(ctx, pattern).Result()
}

func (r *provider) Expire(ctx context.Context, key string, expiration time.Duration) error {
	return r.client.Expire(ctx, key, expiration).Err()
}

func (r *provider) AcquireSemaphoreAtomic(ctx context.Context, key string, maxCount int64, ttl time.Duration) (bool, error) {
	result, err := r.acquireSemaphoreScript.Run(ctx, r.client, []string{key}, maxCount, ttl.Seconds()).Result()
	if err != nil {
		return false, err
	}

	acquired, ok := result.(int64)
	if !ok {
		return false, nil
	}

	return acquired == 1, nil
}

func (r *provider) ReleaseSemaphoreAtomic(ctx context.Context, key string) (int64, error) {
	result, err := r.releaseSemaphoreScript.Run(ctx, r.client, []string{key}).Result()
	if err != nil {
		return 0, err
	}

	count, ok := result.(int64)
	if !ok {
		return 0, nil
	}

	return count, nil
}

func (r *provider) Close() error {
	return r.client.Close()
}
