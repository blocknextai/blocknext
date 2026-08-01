package users

import (
	"context"

	"github.com/google/uuid"
)

type UserRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*User, error)
	Create(ctx context.Context, user *User) error
	GetAllByIDs(ctx context.Context, ids []uuid.UUID) ([]*User, error)
}
