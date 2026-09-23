package users

import (
	"context"
	"errors"

	"github.com/google/uuid"

	"github.com/blocknextai/go-packages/rbac"
	accountApplicationUsers "github.com/blocknextai/platform-api/internal/modules/account/internal/application/users"
	accountDomainUsers "github.com/blocknextai/platform-api/internal/modules/account/internal/domain/users"
)

type GetProfileQuery struct {
	UserID uuid.UUID
}

type GetProfileResponse struct {
	Role        string   `json:"role"`
	Permissions []string `json:"permissions"`
	IsVerified  bool     `json:"isVerified"`
}

func (s *Service) GetProfile(ctx context.Context, request *GetProfileQuery) (*GetProfileResponse, error) {
	user, err := s.userRepository.GetByID(ctx, request.UserID)
	hasRecordNotFound := errors.Is(err, accountDomainUsers.ErrUserNotFound)
	hasUser := user != nil && !hasRecordNotFound

	if !hasUser {
		return nil, accountDomainUsers.ErrUserNotFound
	} else if err != nil {
		return nil, accountApplicationUsers.ErrFailedToGetUser.WithCause(err)
	}

	permissions := rbac.UserPermissions(user.Role)

	return MapUserToResponse(user, permissions), nil
}

func MapUserToResponse(user *accountDomainUsers.User, permissions []string) *GetProfileResponse {
	return &GetProfileResponse{
		Role:        user.Role,
		Permissions: permissions,
		IsVerified:  user.IsVerified,
	}
}
