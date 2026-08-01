package redis

import (
	"context"
	"strings"
	"time"

	taskRunnerDomainTaskRunner "github.com/blocknextai/platform-api/internal/taskrunner/domain/taskrunner"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

const (
	OrganizationSemaphorePrefix = "semaphore:org:"
	BaseDelay                   = 50 * time.Millisecond
	MaxDelay                    = 2 * time.Second
)

type PoolOptions struct {
	PoolSize        int
	MinIdleConns    int
	MaxIdleConns    int
	PoolTimeout     time.Duration
	ConnMaxIdleTime time.Duration
	ConnMaxLifetime time.Duration
}

type RedisSemaphore struct {
	client          *redis.Client
	ttl             time.Duration
	acquireScript   *redis.Script
	releaseScript   *redis.Script
	heartbeatScript *redis.Script
}

func NewRedisSemaphore(addr string, password string, db int, poolOptions PoolOptions, ttl time.Duration) (taskRunnerDomainTaskRunner.SemaphoreManager, error) {
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

	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}

	acquireScript := redis.NewScript(`
		local t = redis.call('TYPE', KEYS[1])
		if t.ok ~= 'none' and t.ok ~= 'zset' then
			redis.call('DEL', KEYS[1])
		end
		local time = redis.call('TIME')
		local now = tonumber(time[1])
		local ttl = tonumber(ARGV[3])
		local expiresAt = now + ttl
		redis.call('ZREMRANGEBYSCORE', KEYS[1], '-inf', '(' .. now)
		if redis.call('ZSCORE', KEYS[1], ARGV[1]) then
			redis.call('ZADD', KEYS[1], expiresAt, ARGV[1])
			redis.call('EXPIRE', KEYS[1], ttl)
			return 1
		end
		if redis.call('ZCARD', KEYS[1]) >= tonumber(ARGV[2]) then
			return 0
		end
		redis.call('ZADD', KEYS[1], expiresAt, ARGV[1])
		redis.call('EXPIRE', KEYS[1], ttl)
		return 1
	`)

	releaseScript := redis.NewScript(`
		local t = redis.call('TYPE', KEYS[1])
		if t.ok ~= 'none' and t.ok ~= 'zset' then
			redis.call('DEL', KEYS[1])
			return 1
		end
		redis.call('ZREM', KEYS[1], ARGV[1])
		if redis.call('ZCARD', KEYS[1]) == 0 then
			redis.call('DEL', KEYS[1])
		end
		return 1
	`)

	heartbeatScript := redis.NewScript(`
		local t = redis.call('TYPE', KEYS[1])
		if t.ok ~= 'none' and t.ok ~= 'zset' then
			redis.call('DEL', KEYS[1])
			return 0
		end
		if redis.call('ZSCORE', KEYS[1], ARGV[1]) == false then
			return 0
		end
		local time = redis.call('TIME')
		local now = tonumber(time[1])
		local ttl = tonumber(ARGV[2])
		redis.call('ZREMRANGEBYSCORE', KEYS[1], '-inf', '(' .. now)
		redis.call('ZADD', KEYS[1], now + ttl, ARGV[1])
		redis.call('EXPIRE', KEYS[1], ttl)
		return 1
	`)

	return &RedisSemaphore{
		client:          client,
		ttl:             ttl,
		acquireScript:   acquireScript,
		releaseScript:   releaseScript,
		heartbeatScript: heartbeatScript,
	}, nil
}

func (sm *RedisSemaphore) Ping(ctx context.Context) error {
	return sm.client.Ping(ctx).Err()
}

func (sm *RedisSemaphore) AcquireSemaphore(ctx context.Context, organizationID uuid.UUID, taskID uuid.UUID, maxConcurrentExecutions int64) (chan struct{}, error) {
	key := sm.buildSemaphoreKey(organizationID)
	currentDelay := BaseDelay

	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		acquired, err := sm.tryAcquire(ctx, key, taskID, maxConcurrentExecutions)
		if err != nil {
			return nil, err
		}

		if acquired {
			semaphore := make(chan struct{}, 1)
			semaphore <- struct{}{}
			return semaphore, nil
		}

		select {
		case <-time.After(currentDelay):
			if currentDelay < MaxDelay {
				currentDelay = min(currentDelay*2, MaxDelay)
			}
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}
}

func (sm *RedisSemaphore) tryAcquire(ctx context.Context, key string, taskID uuid.UUID, max int64) (bool, error) {
	result, err := sm.acquireScript.Run(ctx, sm.client,
		[]string{key},
		taskID.String(),
		max,
		int(sm.ttl.Seconds()),
	).Int()
	if err != nil {
		return false, err
	}
	return result == 1, nil
}

func (sm *RedisSemaphore) ReleaseSemaphore(ctx context.Context, organizationID uuid.UUID, taskID uuid.UUID, semaphore chan struct{}) error {
	select {
	case <-semaphore:
	default:
	}

	key := sm.buildSemaphoreKey(organizationID)
	_, err := sm.releaseScript.Run(ctx, sm.client, []string{key}, taskID.String()).Result()
	return err
}

func (sm *RedisSemaphore) HeartbeatSemaphore(ctx context.Context, organizationID uuid.UUID, taskID uuid.UUID) error {
	key := sm.buildSemaphoreKey(organizationID)
	_, err := sm.heartbeatScript.Run(ctx, sm.client,
		[]string{key},
		taskID.String(),
		int(sm.ttl.Seconds()),
	).Result()
	return err
}

func (sm *RedisSemaphore) buildSemaphoreKey(organizationID uuid.UUID) string {
	var builder strings.Builder
	id := organizationID.String()
	builder.Grow(len(OrganizationSemaphorePrefix) + len(id))
	builder.WriteString(OrganizationSemaphorePrefix)
	builder.WriteString(id)
	return builder.String()
}
