package redis

import (
	"context"
	"log/slog"
	"strings"
	"time"

	"github.com/blocknextai/go-packages/apperror"
	"github.com/blocknextai/go-packages/json"
	taskRunnerDomainNode "github.com/blocknextai/platform-api/internal/taskrunner/domain/node"
	taskRunnerDomainTask "github.com/blocknextai/platform-api/internal/taskrunner/domain/task"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

var (
	ErrFailedToMarshalTaskEvent = apperror.Internal("failed to marshal task event")
	ErrFailedToPublishTaskEvent = apperror.Internal("failed to publish task event")
	ErrFailedToMarshalNodeEvent = apperror.Internal("failed to marshal node event")
	ErrFailedToPublishNodeEvent = apperror.Internal("failed to publish node event")
)

type PoolOptions struct {
	PoolSize        int
	MinIdleConns    int
	MaxIdleConns    int
	PoolTimeout     time.Duration
	ConnMaxIdleTime time.Duration
	ConnMaxLifetime time.Duration
}
type redisBroadcaster struct {
	client *redis.Client
}

func New(addr string, password string, db int, poolOptions PoolOptions) (*redisBroadcaster, error) {
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

	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return &redisBroadcaster{
		client: client,
	}, nil
}

func (b *redisBroadcaster) Ping(ctx context.Context) error {
	return b.client.Ping(ctx).Err()
}

func (b *redisBroadcaster) PublishTaskEvent(ctx context.Context, event *taskRunnerDomainTask.TaskEvent) error {
	eventJSON, err := json.Marshal(event)
	if err != nil {
		return ErrFailedToMarshalTaskEvent
	}

	var builder strings.Builder
	builder.WriteString("organization:")
	builder.WriteString(event.OrganizationID.String())
	organizationKey := builder.String()

	err = b.client.Publish(ctx, organizationKey, string(eventJSON)).Err()
	if err != nil {
		return ErrFailedToPublishTaskEvent
	}

	return nil
}

func (b *redisBroadcaster) PublishNodeEvent(ctx context.Context, event *taskRunnerDomainNode.NodeEvent) error {
	eventJSON, err := json.Marshal(event)
	if err != nil {
		return ErrFailedToMarshalNodeEvent
	}

	var builder strings.Builder
	builder.WriteString("organization:")
	builder.WriteString(event.OrganizationID.String())
	organizationKey := builder.String()

	err = b.client.Publish(ctx, organizationKey, string(eventJSON)).Err()
	if err != nil {
		return ErrFailedToPublishNodeEvent
	}

	return nil
}

func (b *redisBroadcaster) Close() error {
	return b.client.Close()
}

func (b *redisBroadcaster) Subscribe(ctx context.Context, organizationID uuid.UUID) (<-chan string, error) {
	var builder strings.Builder
	builder.WriteString("organization:")
	builder.WriteString(organizationID.String())
	organizationKey := builder.String()

	pubsub := b.client.Subscribe(ctx, organizationKey)
	ch := make(chan string, 64)

	go func() {
		defer func() {
			err := pubsub.Close()
			if err != nil {
				slog.ErrorContext(ctx, "Failed to close pubsub",
					"component", "redis_message_broker",
					"error", err)
			}
			close(ch)
		}()

		for msg := range pubsub.Channel() {
			select {
			case ch <- msg.Payload:
			case <-ctx.Done():
				return
			}
		}
	}()

	return ch, nil
}
