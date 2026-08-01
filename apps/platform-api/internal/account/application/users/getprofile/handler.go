package getprofile

import (
	"context"
	"errors"

	"github.com/blocknextai/go-packages/rbac"
	accountApplicationUsers "github.com/blocknextai/platform-api/internal/account/application/users"
	accountDomainUsers "github.com/blocknextai/platform-api/internal/account/domain/users"
)

type Handler struct {
	userRepository accountDomainUsers.UserRepository
}

func New(
	userRepository accountDomainUsers.UserRepository,
) *Handler {
	return &Handler{
		userRepository: userRepository,
	}
}

func (h *Handler) Handle(ctx context.Context, request *GetProfileQuery) (*GetProfileResponse, error) {
	user, err := h.userRepository.GetByID(ctx, request.UserID)
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
