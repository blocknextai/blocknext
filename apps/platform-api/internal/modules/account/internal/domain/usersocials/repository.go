package usersocials

import (
	"context"

	"github.com/google/uuid"
)

type UserSocialRepository interface {
	GetAllByUserID(ctx context.Context, userID uuid.UUID) ([]*UserSocial, error)
	Create(ctx context.Context, social *UserSocial) error
	Delete(ctx context.Context, social *UserSocial) error
}
