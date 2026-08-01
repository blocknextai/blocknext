package relay

import (
	"context"
	"log/slog"
	"time"

	"github.com/blocknextai/platform-api/internal/eventbus/application/idempotency"
	outboxMessagesDomain "github.com/blocknextai/platform-api/internal/eventbus/domain/outboxmessages"
)

const (
	maxBackoffShift = 8
)

// TODO: maybe move to own folder
type EventDispatcher interface {
	Dispatch(ctx context.Context, eventName string, payload []byte) error
}

type Options struct {
	PollInterval    time.Duration
	BatchSize       int
	MaxAttempts     int
	BaseBackoff     time.Duration
	StuckTimeout    time.Duration
	ReclaimInterval time.Duration
}

type RelayService struct {
	repository outboxMessagesDomain.Repository
	dispatcher EventDispatcher
	opts       Options
}

func NewRelayService(repository outboxMessagesDomain.Repository, dispatcher EventDispatcher, opts Options) *RelayService {
	return &RelayService{
		repository: repository,
		dispatcher: dispatcher,
		opts:       opts,
	}
}

func (s *RelayService) Start(ctx context.Context) {
	go s.reclaimLoop(ctx)

	slog.InfoContext(ctx, "event relay started",
		"component", "eventbus_relay",
		"poll_interval", s.opts.PollInterval,
		"batch_size", s.opts.BatchSize)

	for {
		drained, err := s.drain(ctx)
		if err != nil {
			if ctx.Err() != nil {
				return
			}
			slog.ErrorContext(ctx, "event relay poll failed", "component", "eventbus_relay", "error", err)
		}

		if drained >= s.opts.BatchSize && err == nil {
			continue
		}

		select {
		case <-ctx.Done():
			return
		case <-time.After(s.opts.PollInterval):
		}
	}
}

func (s *RelayService) drain(ctx context.Context) (int, error) {
	messages, err := s.repository.ClaimPending(ctx, s.opts.BatchSize, time.Now().UTC())
	if err != nil {
		return 0, err
	}

	for _, message := range messages {
		if ctx.Err() != nil {
			return len(messages), ctx.Err()
		}
		s.process(ctx, message)
	}

	return len(messages), nil
}

func (s *RelayService) process(ctx context.Context, message *outboxMessagesDomain.OutboxMessage) {
	dispatchCtx := idempotency.WithEventID(ctx, message.ID)

	if dispatchErr := s.dispatcher.Dispatch(dispatchCtx, message.EventName, message.Payload); dispatchErr != nil {
		s.handleFailure(ctx, message, dispatchErr)
		return
	}

	if _, err := message.MarkProcessed(time.Now().UTC()); err != nil {
		slog.ErrorContext(ctx, "invalid transition marking processed", "component", "eventbus_relay", "event_id", message.ID, "error", err)
		return
	}
	s.persist(ctx, message)
}

func (s *RelayService) handleFailure(ctx context.Context, message *outboxMessagesDomain.OutboxMessage, dispatchErr error) {
	if s.opts.MaxAttempts > 0 && message.Attempts >= s.opts.MaxAttempts {
		slog.ErrorContext(ctx, "event dead-lettered after exhausting attempts",
			"component", "eventbus_relay",
			"event_id", message.ID,
			"event_name", message.EventName,
			"attempts", message.Attempts,
			"error", dispatchErr)
		if _, err := message.MarkDead(dispatchErr.Error()); err != nil {
			slog.ErrorContext(ctx, "invalid transition marking dead", "component", "eventbus_relay", "event_id", message.ID, "error", err)
			return
		}
		s.persist(ctx, message)
		return
	}

	nextAttemptAt := time.Now().UTC().Add(s.backoff(message.Attempts - 1))
	slog.WarnContext(ctx, "event dispatch failed, rescheduling for retry",
		"component", "eventbus_relay",
		"event_id", message.ID,
		"event_name", message.EventName,
		"attempts", message.Attempts,
		"next_attempt_at", nextAttemptAt,
		"error", dispatchErr)
	if _, err := message.Reschedule(nextAttemptAt, dispatchErr.Error()); err != nil {
		slog.ErrorContext(ctx, "invalid transition rescheduling message", "component", "eventbus_relay", "event_id", message.ID, "error", err)
		return
	}
	s.persist(ctx, message)
}

func (s *RelayService) persist(ctx context.Context, message *outboxMessagesDomain.OutboxMessage) {
	if err := s.repository.Update(ctx, message); err != nil {
		slog.ErrorContext(ctx, "failed to persist outbox message state",
			"component", "eventbus_relay",
			"event_id", message.ID,
			"status", message.Status,
			"error", err)
	}
}

func (s *RelayService) backoff(attempts int) time.Duration {
	shift := min(max(attempts, 0), maxBackoffShift)
	return s.opts.BaseBackoff << uint(shift)
}

func (s *RelayService) reclaimLoop(ctx context.Context) {
	ticker := time.NewTicker(s.opts.ReclaimInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			lockedBefore := time.Now().UTC().Add(-s.opts.StuckTimeout)
			reclaimed, err := s.repository.ReclaimStuck(ctx, lockedBefore)
			if err != nil {
				slog.ErrorContext(ctx, "failed to reclaim stuck events", "component", "eventbus_relay", "error", err)
				continue
			}
			if reclaimed > 0 {
				slog.InfoContext(ctx, "reclaimed stuck events", "component", "eventbus_relay", "count", reclaimed)
			}
		}
	}
}
