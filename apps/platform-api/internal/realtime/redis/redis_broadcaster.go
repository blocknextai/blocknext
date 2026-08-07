package redis

import (
	"context"
	"log/slog"
	"strings"

	"github.com/blocknextai/go-packages/apperror"
	"github.com/blocknextai/go-packages/redisclient"
	"github.com/blocknextai/platform-api/internal/realtime/events"
	taskRunnerDomainNode "github.com/blocknextai/platform-api/internal/taskrunner/domain/node"
	taskRunnerDomainTask "github.com/blocknextai/platform-api/internal/taskrunner/domain/task"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

const organizationKeyPrefix = "organization:"

var (
	ErrFailedToPublishTaskEvent = apperror.Internal("failed to publish task event")
	ErrFailedToPublishNodeEvent = apperror.Internal("failed to publish node event")
)

type redisBroadcaster struct {
	client *redis.Client
}

func New(addr string, password string, db int, poolOptions redisclient.PoolOptions) (*redisBroadcaster, error) {
	client, err := redisclient.New(addr, password, db, poolOptions)
	if err != nil {
		return nil, err
	}

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
	payload, err := events.MarshalTask(event)
	if err != nil {
		return err
	}

	err = b.client.Publish(ctx, buildOrganizationKey(event.OrganizationID), payload).Err()
	if err != nil {
		return ErrFailedToPublishTaskEvent
	}

	return nil
}

func (b *redisBroadcaster) PublishNodeEvent(ctx context.Context, event *taskRunnerDomainNode.NodeEvent) error {
	payload, err := events.MarshalNode(event)
	if err != nil {
		return err
	}

	err = b.client.Publish(ctx, buildOrganizationKey(event.OrganizationID), payload).Err()
	if err != nil {
		return ErrFailedToPublishNodeEvent
	}

	return nil
}

func (b *redisBroadcaster) Close() error {
	return b.client.Close()
}

func (b *redisBroadcaster) Subscribe(ctx context.Context, organizationID uuid.UUID) (<-chan string, error) {
	pubsub := b.client.Subscribe(ctx, buildOrganizationKey(organizationID))
	ch := make(chan string, events.SubscriberBuffer)

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

func buildOrganizationKey(organizationID uuid.UUID) string {
	id := organizationID.String()

	var builder strings.Builder
	builder.Grow(len(organizationKeyPrefix) + len(id))
	builder.WriteString(organizationKeyPrefix)
	builder.WriteString(id)

	return builder.String()
}
