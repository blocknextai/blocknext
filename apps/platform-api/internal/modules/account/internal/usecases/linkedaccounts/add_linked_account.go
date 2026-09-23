package linkedaccounts

import (
	"context"
	"errors"

	"github.com/google/uuid"

	"github.com/blocknextai/go-packages/apperror"
	accountApplicationLinkedAccounts "github.com/blocknextai/platform-api/internal/modules/account/internal/application/linkedaccounts"
	"github.com/blocknextai/platform-api/internal/modules/account/internal/domain"
	accountDomainLinkedAccounts "github.com/blocknextai/platform-api/internal/modules/account/internal/domain/linkedaccounts"
	accountDomainUsers "github.com/blocknextai/platform-api/internal/modules/account/internal/domain/users"
	authUseCases "github.com/blocknextai/platform-api/internal/modules/account/internal/usecases/auth"
)

type AddLinkedAccountCommand struct {
	UserID       uuid.UUID
	AuthProvider domain.AuthProvider
	Payload      map[string]any
}

var (
	ErrPayloadRequired = apperror.Validation("payload is required")
)

func (c *AddLinkedAccountCommand) Validate() error {
	if c.Payload == nil {
		return ErrPayloadRequired
	}

	return nil
}

type AddLinkedAccountResponse struct {
	ID           uuid.UUID `json:"id"`
	AuthProvider string    `json:"authProvider"`
	Identifier   string    `json:"identifier"`
	DisplayName  *string   `json:"displayName"`
	IsPrimary    bool      `json:"isPrimary"`
}

func (s *Service) AddLinkedAccount(ctx context.Context, command *AddLinkedAccountCommand) (*AddLinkedAccountResponse, error) {
	if err := command.Validate(); err != nil {
		return nil, err
	}

	var response *AddLinkedAccountResponse

	err := s.transactionManager.ExecuteInTransaction(ctx, func(txCtx context.Context) error {
		user, err := s.userRepository.GetByID(txCtx, command.UserID)
		if err != nil {
			return accountDomainUsers.ErrUserNotFound
		}

		authProvider, err := s.authProviderRegistry.GetProvider(command.AuthProvider)
		if err != nil {
			return err
		}

		authResponse, err := authProvider.Validate(txCtx, authUseCases.Request{
			AuthProvider: command.AuthProvider,
			Payload:      command.Payload,
		})
		if err != nil {
			return err
		}

		existingLinkedAccount, err := s.linkedAccountRepository.GetByProviderID(txCtx, authResponse.ProviderID)
		if err != nil && !errors.Is(err, accountDomainLinkedAccounts.ErrLinkedAccountNotFound) {
			return err
		}

		if existingLinkedAccount != nil {
			return accountApplicationLinkedAccounts.ErrLinkedAccountAlreadyExists
		}

		linkedAccount, err := accountDomainLinkedAccounts.NewLinkedAccount(
			user.ID,
			command.AuthProvider,
			authResponse.ProviderID,
			authResponse.Identifier,
			&authResponse.DisplayName,
			false,
			true,
		)
		if err != nil {
			return err
		}

		deletedNonce, err := authResponse.Nonce.Delete()
		if err != nil {
			return err
		}

		if err = s.userNonceRepository.Delete(txCtx, deletedNonce); err != nil {
			return err
		}

		err = s.linkedAccountRepository.Create(txCtx, linkedAccount)
		if err != nil {
			return accountApplicationLinkedAccounts.ErrFailedToAddLinkedAccount.WithCause(err)
		}

		response = &AddLinkedAccountResponse{
			ID:           linkedAccount.ID,
			AuthProvider: linkedAccount.AuthProvider.String(),
			Identifier:   linkedAccount.Identifier,
			DisplayName:  linkedAccount.DisplayName,
			IsPrimary:    linkedAccount.IsPrimary,
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return response, nil
}
