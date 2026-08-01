package getalllinkedaccounts

import (
	"context"
	"errors"

	accountApplicationLinkedAccounts "github.com/blocknextai/platform-api/internal/account/application/linkedaccounts"
	accountDomainLinkedAccounts "github.com/blocknextai/platform-api/internal/account/domain/linkedaccounts"
	accountDomainUsers "github.com/blocknextai/platform-api/internal/account/domain/users"
)

type Handler struct {
	userRepository          accountDomainUsers.UserRepository
	linkedAccountRepository accountDomainLinkedAccounts.LinkedAccountRepository
}

func New(
	userRepository accountDomainUsers.UserRepository,
	linkedAccountRepository accountDomainLinkedAccounts.LinkedAccountRepository,
) *Handler {
	return &Handler{
		userRepository:          userRepository,
		linkedAccountRepository: linkedAccountRepository,
	}
}

func (h *Handler) Handle(ctx context.Context, request *GetAllLinkedAccountsQuery) (*GetAllLinkedAccountsResponse, error) {
	user, err := h.userRepository.GetByID(ctx, request.UserID)
	hasRecordNotFound := errors.Is(err, accountDomainUsers.ErrUserNotFound)
	hasUser := user != nil && !hasRecordNotFound

	if !hasUser {
		return nil, accountDomainUsers.ErrUserNotFound
	} else if err != nil {
		return nil, accountApplicationLinkedAccounts.ErrFailedToGetUser.WithCause(err)
	}

	linkedAccounts, err := h.linkedAccountRepository.GetAllByUserID(ctx, request.UserID)
	if err != nil {
		return nil, err
	}

	return MapLinkedAccountsToResponse(linkedAccounts), nil
}
