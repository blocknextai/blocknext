package idempotency

import (
	"context"

	"github.com/google/uuid"
)

type eventIDContextKey struct{}

func WithEventID(ctx context.Context, eventID uuid.UUID) context.Context {
	return context.WithValue(ctx, eventIDContextKey{}, eventID)
}

func EventIDFromContext(ctx context.Context) (uuid.UUID, bool) {
	eventID, ok := ctx.Value(eventIDContextKey{}).(uuid.UUID)
	return eventID, ok
}
