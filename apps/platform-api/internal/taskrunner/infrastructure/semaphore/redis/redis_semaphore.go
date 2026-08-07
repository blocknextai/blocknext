package redis

import (
	"context"
	"strings"
	"time"

	"github.com/blocknextai/go-packages/redisclient"
	taskRunnerDomainTaskRunner "github.com/blocknextai/platform-api/internal/taskrunner/domain/taskrunner"
	"github.com/blocknextai/platform-api/internal/taskrunner/infrastructure/semaphore"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

const (
	OrganizationSemaphorePrefix = "semaphore:org:"
)

type RedisSemaphore struct {
	client          *redis.Client
	ttl             time.Duration
	acquireScript   *redis.Script
	releaseScript   *redis.Script
	heartbeatScript *redis.Script
}

func NewRedisSemaphore(addr string, password string, db int, poolOptions redisclient.PoolOptions, ttl time.Duration) (taskRunnerDomainTaskRunner.SemaphoreManager, error) {
	client, err := redisclient.New(addr, password, db, poolOptions)
	if err != nil {
		return nil, err
	}

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

	return semaphore.AcquireWithBackoff(ctx, func(ctx context.Context) (bool, error) {
		return sm.tryAcquire(ctx, key, taskID, maxConcurrentExecutions)
	})
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

func (sm *RedisSemaphore) ReleaseSemaphore(ctx context.Context, organizationID uuid.UUID, taskID uuid.UUID, token chan struct{}) error {
	semaphore.Drain(token)

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
