package linkedaccounts

import (
	"context"
	"errors"

	"github.com/google/uuid"

	accountApplicationLinkedAccounts "github.com/blocknextai/platform-api/internal/modules/account/internal/application/linkedaccounts"
	accountDomainLinkedaccounts "github.com/blocknextai/platform-api/internal/modules/account/internal/domain/linkedaccounts"
	accountDomainUsers "github.com/blocknextai/platform-api/internal/modules/account/internal/domain/users"
)

type GetAllLinkedAccountsQuery struct {
	UserID uuid.UUID
}

type LinkedAccount struct {
	ID           uuid.UUID `json:"id"`
	AuthProvider string    `json:"authProvider"`
	Identifier   string    `json:"identifier"`
	DisplayName  *string   `json:"displayName"`
	IsPrimary    bool      `json:"isPrimary"`
	IsVerified   bool      `json:"isVerified"`
}

type GetAllLinkedAccountsResponse = []LinkedAccount

func (s *Service) GetAllLinkedAccounts(ctx context.Context, request *GetAllLinkedAccountsQuery) (*GetAllLinkedAccountsResponse, error) {
	user, err := s.userRepository.GetByID(ctx, request.UserID)
	hasRecordNotFound := errors.Is(err, accountDomainUsers.ErrUserNotFound)
	hasUser := user != nil && !hasRecordNotFound

	if !hasUser {
		return nil, accountDomainUsers.ErrUserNotFound
	} else if err != nil {
		return nil, accountApplicationLinkedAccounts.ErrFailedToGetUser.WithCause(err)
	}

	linkedAccounts, err := s.linkedAccountRepository.GetAllByUserID(ctx, request.UserID)
	if err != nil {
		return nil, err
	}

	return MapLinkedAccountsToResponse(linkedAccounts), nil
}

func MapLinkedAccountsToResponse(linkedAccounts []*accountDomainLinkedaccounts.LinkedAccount) *GetAllLinkedAccountsResponse {
	result := make(GetAllLinkedAccountsResponse, 0, len(linkedAccounts))
	for _, linkedAccount := range linkedAccounts {
		result = append(result, LinkedAccount{
			ID:           linkedAccount.ID,
			AuthProvider: linkedAccount.AuthProvider.String(),
			Identifier:   linkedAccount.Identifier,
			DisplayName:  linkedAccount.DisplayName,
			IsPrimary:    linkedAccount.IsPrimary,
			IsVerified:   linkedAccount.IsVerified,
		})
	}
	return &result
}
