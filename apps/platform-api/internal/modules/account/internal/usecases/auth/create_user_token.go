package auth

import (
	"context"
	"errors"

	"github.com/google/uuid"

	"github.com/blocknextai/go-packages/apperror"
	accountApplicationAuth "github.com/blocknextai/platform-api/internal/modules/account/internal/application/auth"
	"github.com/blocknextai/platform-api/internal/modules/account/internal/domain"
	accountDomainLinkedAccounts "github.com/blocknextai/platform-api/internal/modules/account/internal/domain/linkedaccounts"
	accountDomainUserPreferences "github.com/blocknextai/platform-api/internal/modules/account/internal/domain/userpreferences"
	accountDomainUsers "github.com/blocknextai/platform-api/internal/modules/account/internal/domain/users"
)

type CreateUserTokenCommand struct {
	AuthProvider domain.AuthProvider
	Payload      map[string]any
	IPAddress    string
	UserAgent    string
}

var (
	ErrInvalidPayload = apperror.Validation("invalid payload")
)

func (c *CreateUserTokenCommand) Validate() error {
	if c.Payload == nil {
		return ErrInvalidPayload
	}

	return nil
}

func (s *Service) CreateUserToken(ctx context.Context, request *CreateUserTokenCommand) (*accountApplicationAuth.AccessTokenResponse, error) {
	if err := request.Validate(); err != nil {
		return nil, err
	}

	var response *accountApplicationAuth.AccessTokenResponse
	var userID uuid.UUID

	err := s.transactionManager.ExecuteInTransaction(ctx, func(txCtx context.Context) error {
		authProvider, err := s.authProviderRegistry.GetProvider(request.AuthProvider)
		if err != nil {
			return err
		}

		authResponse, err := authProvider.Validate(txCtx, Request{
			AuthProvider: request.AuthProvider,
			Payload:      request.Payload,
		})
		if err != nil {
			return err
		}

		linkedAccount, err := s.linkedAccountRepository.GetByAuthProviderAndProviderID(txCtx, request.AuthProvider, authResponse.ProviderID)
		if err != nil && !errors.Is(err, accountDomainLinkedAccounts.ErrLinkedAccountNotFound) {
			return err
		}

		if linkedAccount == nil {
			user, err := accountDomainUsers.NewUser()
			if err != nil {
				return err
			}

			err = s.userRepository.Create(txCtx, user)
			if err != nil {
				return err
			}

			preferences, err := accountDomainUserPreferences.NewDefault(user.ID)
			if err != nil {
				return err
			}

			if err = s.userPreferenceRepository.Upsert(txCtx, preferences); err != nil {
				return err
			}

			linkedAccount, err = accountDomainLinkedAccounts.NewLinkedAccount(
				user.ID,
				request.AuthProvider,
				authResponse.ProviderID,
				authResponse.Identifier,
				&authResponse.DisplayName,
				true,
				true,
			)
			if err != nil {
				return err
			}

			err = s.linkedAccountRepository.Create(txCtx, linkedAccount)
			if err != nil {
				return err
			}

			userID = user.ID

			if err = s.eventBusPublisherService.Enqueue(txCtx, accountDomainUsers.UserCreatedDomainEvent{
				UserID:      user.ID,
				Identifier:  linkedAccount.Identifier,
				DisplayName: authResponse.DisplayName,
			}); err != nil {
				return err
			}
		} else {
			user, err := s.userRepository.GetByID(txCtx, linkedAccount.UserID)
			if err != nil {
				return err
			}

			userID = user.ID
		}

		deletedNonce, err := authResponse.Nonce.Delete()
		if err != nil {
			return err
		}

		if err = s.userNonceRepository.Delete(txCtx, deletedNonce); err != nil {
			return err
		}

		issued, err := s.tokenIssuer.IssueTokens(txCtx, userID, request.AuthProvider, request.IPAddress, request.UserAgent)
		if err != nil {
			return err
		}
		response = issued
		return nil
	})

	if err != nil {
		return nil, err
	}

	return response, nil
}
