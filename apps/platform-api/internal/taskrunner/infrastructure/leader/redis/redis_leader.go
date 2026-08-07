package redis

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/blocknextai/go-packages/apperror"
	"github.com/blocknextai/go-packages/redisclient"
	taskRunnerDomainTaskRunner "github.com/blocknextai/platform-api/internal/taskrunner/domain/taskrunner"
	"github.com/redis/go-redis/v9"
)

var (
	ErrLockNotHeld = apperror.Internal("leader lock no longer held by this instance")
)

type Options struct {
	Key           string
	InstanceID    string
	TTL           time.Duration
	PollInterval  time.Duration
	RenewInterval time.Duration
}

type RedisLeader struct {
	client                 *redis.Client
	compareAndRenewScript  *redis.Script
	compareAndDeleteScript *redis.Script
	options                Options
}

func NewRedisLeader(addr string, password string, db int, poolOptions redisclient.PoolOptions, options Options) (taskRunnerDomainTaskRunner.LeaderRunner, error) {
	client, err := redisclient.New(addr, password, db, poolOptions)
	if err != nil {
		return nil, err
	}

	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}

	compareAndRenewScript := redis.NewScript(`
		if redis.call("GET", KEYS[1]) == ARGV[1] then
			return redis.call("EXPIRE", KEYS[1], ARGV[2])
		end
		return 0
	`)

	compareAndDeleteScript := redis.NewScript(`
		if redis.call("GET", KEYS[1]) == ARGV[1] then
			return redis.call("DEL", KEYS[1])
		end
		return 0
	`)

	return &RedisLeader{
		client:                 client,
		compareAndRenewScript:  compareAndRenewScript,
		compareAndDeleteScript: compareAndDeleteScript,
		options:                options,
	}, nil
}

func (l *RedisLeader) Ping(ctx context.Context) error {
	return l.client.Ping(ctx).Err()
}

func (l *RedisLeader) Run(ctx context.Context, leaderFunc func(leaderCtx context.Context)) {
	pollTicker := time.NewTicker(l.options.PollInterval)
	defer pollTicker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		acquired, err := l.tryAcquire(ctx)
		if err != nil {
			slog.WarnContext(ctx, "leader: failed to attempt acquire",
				"component", "leader_runner", "key", l.options.Key, "error", err)
		}

		if !acquired {
			select {
			case <-ctx.Done():
				return
			case <-pollTicker.C:
			}
			continue
		}

		slog.InfoContext(ctx, "leader: acquired leadership",
			"component", "leader_runner", "key", l.options.Key, "instance", l.options.InstanceID)

		leaderCtx, cancel := context.WithCancel(ctx)
		funcDone := make(chan struct{})
		go func() {
			defer close(funcDone)
			leaderFunc(leaderCtx)
		}()

		l.holdLeadership(ctx, cancel)

		<-funcDone
		l.release(context.Background())

		if ctx.Err() != nil {
			return
		}
	}
}

func (l *RedisLeader) holdLeadership(ctx context.Context, leaderCancel context.CancelFunc) {
	renewTicker := time.NewTicker(l.options.RenewInterval)
	defer renewTicker.Stop()

	for {
		select {
		case <-ctx.Done():
			leaderCancel()
			return
		case <-renewTicker.C:
		}

		if err := l.tryRenew(ctx); err != nil {
			slog.WarnContext(ctx, "leader: lost leadership",
				"component", "leader_runner", "key", l.options.Key, "error", err)
			leaderCancel()
			return
		}
	}
}

func (l *RedisLeader) tryAcquire(ctx context.Context) (bool, error) {
	set, err := l.client.SetNX(ctx, l.options.Key, l.options.InstanceID, l.options.TTL).Result()
	if err != nil {
		return false, err
	}
	return set, nil
}

func (l *RedisLeader) tryRenew(ctx context.Context) error {
	result, err := l.compareAndRenewScript.Run(ctx, l.client,
		[]string{l.options.Key},
		l.options.InstanceID,
		int(l.options.TTL.Seconds()),
	).Int()
	if err != nil {
		return err
	}
	if result == 0 {
		return ErrLockNotHeld
	}
	return nil
}

func (l *RedisLeader) release(ctx context.Context) {
	if _, err := l.compareAndDeleteScript.Run(ctx, l.client,
		[]string{l.options.Key},
		l.options.InstanceID,
	).Result(); err != nil && !errors.Is(err, redis.Nil) {
		slog.WarnContext(ctx, "leader: failed to release lock",
			"component", "leader_runner", "key", l.options.Key, "error", err)
	}
}
