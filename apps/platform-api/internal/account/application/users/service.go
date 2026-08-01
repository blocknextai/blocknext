package users

import (
	"context"

	"github.com/blocknextai/platform-api/internal/account/domain/users"
	"github.com/google/uuid"
)

type UserService interface {
	GetByID(ctx context.Context, id uuid.UUID) (*users.User, error)
	GetAllByIDs(ctx context.Context, ids []uuid.UUID) ([]*users.User, error)
}

type userService struct {
	userRepository users.UserRepository
}

func NewUserService(userRepository users.UserRepository) UserService {
	return &userService{userRepository: userRepository}
}

func (s *userService) GetByID(ctx context.Context, id uuid.UUID) (*users.User, error) {
	return s.userRepository.GetByID(ctx, id)
}

func (s *userService) GetAllByIDs(ctx context.Context, ids []uuid.UUID) ([]*users.User, error) {
	return s.userRepository.GetAllByIDs(ctx, ids)
}
