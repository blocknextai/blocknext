package outboxmessages

import (
	"context"
	"time"
)

type Repository interface {
	Create(ctx context.Context, message *OutboxMessage) error
	ClaimPending(ctx context.Context, limit int, now time.Time) ([]*OutboxMessage, error)
	Update(ctx context.Context, message *OutboxMessage) error
	ReclaimStuck(ctx context.Context, lockedBefore time.Time) (int, error)
}
