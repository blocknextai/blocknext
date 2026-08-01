package redis

import (
	"context"
	"errors"
	"log/slog"
	"strings"
	"time"

	"github.com/blocknextai/go-packages/apperror"
	"github.com/blocknextai/go-packages/json"
	taskRunnerDomainTask "github.com/blocknextai/platform-api/internal/taskrunner/domain/task"
	"github.com/blocknextai/platform-api/internal/taskrunner/domain/taskqueue"
	"github.com/redis/go-redis/v9"
)

var (
	ErrFailedToMarshalEnvelope   = apperror.Internal("failed to marshal task envelope")
	ErrFailedToUnmarshalEnvelope = apperror.Internal("failed to unmarshal task envelope")
	ErrFailedToPushTask          = apperror.Internal("failed to push task to queue")
	ErrFailedToAckTask           = apperror.Internal("failed to ack task")
)

const (
	envelopeField           = "envelope"
	orphanCleanupInterval   = time.Hour
	orphanConsumerIdleLimit = 24 * time.Hour
)

type Options struct {
	StreamName    string
	ConsumerGroup string
	ConsumerName  string
	BlockTimeout  time.Duration
	PrefetchCount int
}

type PoolOptions struct {
	PoolSize        int
	MinIdleConns    int
	MaxIdleConns    int
	PoolTimeout     time.Duration
	ConnMaxIdleTime time.Duration
	ConnMaxLifetime time.Duration
}

type RedisTaskQueue struct {
	client  *redis.Client
	options Options
}

func NewRedisTaskQueue(addr string, password string, db int, poolOptions PoolOptions, options Options) (taskqueue.TaskQueue, error) {
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

	if err := ensureConsumerGroup(ctx, client, options.StreamName, options.ConsumerGroup); err != nil {
		return nil, err
	}

	return &RedisTaskQueue{
		client:  client,
		options: options,
	}, nil
}

func ensureConsumerGroup(ctx context.Context, client *redis.Client, stream string, group string) error {
	err := client.XGroupCreateMkStream(ctx, stream, group, "$").Err()
	if err == nil {
		return nil
	}
	if strings.Contains(err.Error(), "BUSYGROUP") {
		return nil
	}
	return err
}

func (q *RedisTaskQueue) Ping(ctx context.Context) error {
	return q.client.Ping(ctx).Err()
}

func (q *RedisTaskQueue) Push(ctx context.Context, envelope taskRunnerDomainTask.TaskEnvelope) error {
	body, err := json.Marshal(envelope)
	if err != nil {
		return ErrFailedToMarshalEnvelope
	}

	if err := q.client.XAdd(ctx, &redis.XAddArgs{
		Stream: q.options.StreamName,
		Values: map[string]any{envelopeField: body},
	}).Err(); err != nil {
		slog.ErrorContext(ctx, "failed to push task to redis stream",
			"component", "redis_task_queue",
			"stream", q.options.StreamName,
			"error", err)
		return ErrFailedToPushTask
	}

	return nil
}

func (q *RedisTaskQueue) Subscribe(ctx context.Context) (<-chan taskqueue.TaskMessage, error) {
	ch := make(chan taskqueue.TaskMessage, q.options.PrefetchCount)

	go q.consumeLoop(ctx, ch)
	go q.cleanupOrphanConsumersLoop(ctx)

	return ch, nil
}

func (q *RedisTaskQueue) cleanupOrphanConsumersLoop(ctx context.Context) {
	ticker := time.NewTicker(orphanCleanupInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			q.cleanupOrphanConsumers(ctx)
		}
	}
}

func (q *RedisTaskQueue) cleanupOrphanConsumers(ctx context.Context) {
	consumers, err := q.client.XInfoConsumers(ctx, q.options.StreamName, q.options.ConsumerGroup).Result()
	if err != nil {
		slog.WarnContext(ctx, "failed to list consumers for orphan cleanup",
			"component", "redis_task_queue",
			"error", err)
		return
	}

	deleted := 0
	for _, consumer := range consumers {
		if consumer.Name == q.options.ConsumerName {
			continue
		}
		if consumer.Pending > 0 {
			continue
		}
		if time.Duration(consumer.Idle)*time.Millisecond < orphanConsumerIdleLimit {
			continue
		}

		if err := q.client.XGroupDelConsumer(ctx, q.options.StreamName, q.options.ConsumerGroup, consumer.Name).Err(); err != nil {
			slog.WarnContext(ctx, "failed to delete orphan consumer",
				"component", "redis_task_queue",
				"consumer", consumer.Name,
				"error", err)
			continue
		}
		deleted++
	}

	if deleted > 0 {
		slog.InfoContext(ctx, "cleaned up orphan consumers",
			"component", "redis_task_queue",
			"count", deleted)
	}
}

func (q *RedisTaskQueue) consumeLoop(ctx context.Context, ch chan<- taskqueue.TaskMessage) {
	defer close(ch)

	q.drainOwnPEL(ctx, ch)

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		streams, err := q.client.XReadGroup(ctx, &redis.XReadGroupArgs{
			Group:    q.options.ConsumerGroup,
			Consumer: q.options.ConsumerName,
			Streams:  []string{q.options.StreamName, ">"},
			Count:    int64(q.options.PrefetchCount),
			Block:    q.options.BlockTimeout,
		}).Result()

		if err != nil {
			if errors.Is(err, redis.Nil) || errors.Is(err, context.Canceled) {
				continue
			}
			slog.ErrorContext(ctx, "failed to read from redis stream",
				"component", "redis_task_queue",
				"stream", q.options.StreamName,
				"group", q.options.ConsumerGroup,
				"error", err)
			select {
			case <-ctx.Done():
				return
			case <-time.After(time.Second):
			}
			continue
		}

		for _, stream := range streams {
			for _, msg := range stream.Messages {
				envelope, decodeErr := q.decodeEnvelope(msg.Values)
				if decodeErr != nil {
					slog.ErrorContext(ctx, "failed to decode task envelope, dropping message",
						"component", "redis_task_queue",
						"message_id", msg.ID,
						"error", decodeErr)
					if ackErr := q.Ack(ctx, msg.ID); ackErr != nil {
						slog.ErrorContext(ctx, "failed to ack malformed message",
							"component", "redis_task_queue",
							"message_id", msg.ID,
							"error", ackErr)
					}
					continue
				}

				select {
				case ch <- taskqueue.TaskMessage{ID: msg.ID, Envelope: envelope}:
				case <-ctx.Done():
					return
				}
			}
		}
	}
}

func (q *RedisTaskQueue) drainOwnPEL(ctx context.Context, ch chan<- taskqueue.TaskMessage) {
	lastID := "0"
	totalDrained := 0

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		streams, err := q.client.XReadGroup(ctx, &redis.XReadGroupArgs{
			Group:    q.options.ConsumerGroup,
			Consumer: q.options.ConsumerName,
			Streams:  []string{q.options.StreamName, lastID},
			Count:    int64(q.options.PrefetchCount),
		}).Result()

		if err != nil {
			if errors.Is(err, redis.Nil) {
				break
			}
			slog.ErrorContext(ctx, "failed to drain own PEL on startup",
				"component", "redis_task_queue",
				"consumer", q.options.ConsumerName,
				"error", err)
			return
		}

		if len(streams) == 0 || len(streams[0].Messages) == 0 {
			break
		}

		for _, stream := range streams {
			for _, msg := range stream.Messages {
				envelope, decodeErr := q.decodeEnvelope(msg.Values)
				if decodeErr != nil {
					slog.ErrorContext(ctx, "failed to decode pending message, acking",
						"component", "redis_task_queue",
						"message_id", msg.ID,
						"error", decodeErr)
					if ackErr := q.Ack(ctx, msg.ID); ackErr != nil {
						slog.ErrorContext(ctx, "failed to ack malformed pending message",
							"component", "redis_task_queue",
							"message_id", msg.ID,
							"error", ackErr)
					}
					lastID = msg.ID
					continue
				}

				select {
				case ch <- taskqueue.TaskMessage{ID: msg.ID, Envelope: envelope}:
					totalDrained++
				case <-ctx.Done():
					return
				}
				lastID = msg.ID
			}
		}
	}

	if totalDrained > 0 {
		slog.InfoContext(ctx, "drained own pending entries on startup",
			"component", "redis_task_queue",
			"consumer", q.options.ConsumerName,
			"count", totalDrained)
	}
}

func (q *RedisTaskQueue) decodeEnvelope(values map[string]any) (taskRunnerDomainTask.TaskEnvelope, error) {
	raw, ok := values[envelopeField].(string)
	if !ok {
		return taskRunnerDomainTask.TaskEnvelope{}, ErrFailedToUnmarshalEnvelope
	}
	var envelope taskRunnerDomainTask.TaskEnvelope
	if err := json.Unmarshal([]byte(raw), &envelope); err != nil {
		return taskRunnerDomainTask.TaskEnvelope{}, ErrFailedToUnmarshalEnvelope
	}
	return envelope, nil
}

func (q *RedisTaskQueue) Ack(ctx context.Context, msgID string) error {
	if err := q.client.XAck(ctx, q.options.StreamName, q.options.ConsumerGroup, msgID).Err(); err != nil {
		slog.ErrorContext(ctx, "failed to ack task message",
			"component", "redis_task_queue",
			"message_id", msgID,
			"error", err)
		return ErrFailedToAckTask
	}
	return nil
}

func (q *RedisTaskQueue) ReclaimIdle(ctx context.Context, idleTimeout time.Duration) (int, error) {
	reclaimed := 0
	start := "0-0"

	for {
		select {
		case <-ctx.Done():
			return reclaimed, ctx.Err()
		default:
		}

		msgs, nextStart, err := q.client.XAutoClaim(ctx, &redis.XAutoClaimArgs{
			Stream:   q.options.StreamName,
			Group:    q.options.ConsumerGroup,
			Consumer: q.options.ConsumerName,
			MinIdle:  idleTimeout,
			Start:    start,
			Count:    int64(q.options.PrefetchCount),
		}).Result()
		if err != nil {
			return reclaimed, err
		}

		reclaimed += len(msgs)

		if nextStart == "0-0" || nextStart == "" {
			break
		}
		start = nextStart
	}

	return reclaimed, nil
}

func (q *RedisTaskQueue) Close() error {
	return q.client.Close()
}
