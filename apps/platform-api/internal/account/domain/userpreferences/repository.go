package userpreferences

import (
	"context"

	"github.com/google/uuid"
)

type UserPreferenceRepository interface {
	GetByUserID(ctx context.Context, userID uuid.UUID) (*UserPreference, error)
	Upsert(ctx context.Context, preference *UserPreference) error
}
