package createusertoken

import (
	"context"
	"errors"

	"github.com/blocknextai/go-packages/database"
	accountApplicationAuth "github.com/blocknextai/platform-api/internal/account/application/auth"
	"github.com/blocknextai/platform-api/internal/account/application/auth/tokenissuer"
	accountDomainLinkedAccounts "github.com/blocknextai/platform-api/internal/account/domain/linkedaccounts"
	accountDomainUserNonces "github.com/blocknextai/platform-api/internal/account/domain/usernonces"
	accountDomainUserPreferences "github.com/blocknextai/platform-api/internal/account/domain/userpreferences"
	accountDomainUsers "github.com/blocknextai/platform-api/internal/account/domain/users"
	"github.com/blocknextai/platform-api/internal/eventbus/application/publishing"
	"github.com/google/uuid"
)

type Handler struct {
	userRepository           accountDomainUsers.UserRepository
	userNonceRepository      accountDomainUserNonces.UserNonceRepository
	linkedAccountRepository  accountDomainLinkedAccounts.LinkedAccountRepository
	userPreferenceRepository accountDomainUserPreferences.UserPreferenceRepository
	tokenIssuer              tokenissuer.Service
	authProviderRegistry     AuthProviderRegistry
	eventBusPublisherService publishing.PublisherService
	transactionManager       database.TransactionManager
}

func New(
	userRepository accountDomainUsers.UserRepository,
	userNonceRepository accountDomainUserNonces.UserNonceRepository,
	linkedAccountRepository accountDomainLinkedAccounts.LinkedAccountRepository,
	userPreferenceRepository accountDomainUserPreferences.UserPreferenceRepository,
	tokenIssuer tokenissuer.Service,
	authProviderRegistry AuthProviderRegistry,
	eventBusPublisherService publishing.PublisherService,
	transactionManager database.TransactionManager,
) *Handler {
	return &Handler{
		userRepository:           userRepository,
		userNonceRepository:      userNonceRepository,
		linkedAccountRepository:  linkedAccountRepository,
		userPreferenceRepository: userPreferenceRepository,
		tokenIssuer:              tokenIssuer,
		authProviderRegistry:     authProviderRegistry,
		eventBusPublisherService: eventBusPublisherService,
		transactionManager:       transactionManager,
	}
}

func (h *Handler) Handle(ctx context.Context, request *CreateUserTokenCommand) (*accountApplicationAuth.AccessTokenResponse, error) {
	var response *accountApplicationAuth.AccessTokenResponse
	var userID uuid.UUID

	err := h.transactionManager.ExecuteInTransaction(ctx, func(txCtx context.Context) error {
		authProvider, err := h.authProviderRegistry.GetProvider(request.AuthProvider)
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

		linkedAccount, err := h.linkedAccountRepository.GetByAuthProviderAndProviderID(txCtx, request.AuthProvider, authResponse.ProviderID)
		if err != nil && !errors.Is(err, accountDomainLinkedAccounts.ErrLinkedAccountNotFound) {
			return err
		}

		if linkedAccount == nil {
			user, err := accountDomainUsers.NewUser()
			if err != nil {
				return err
			}

			err = h.userRepository.Create(txCtx, user)
			if err != nil {
				return err
			}

			preferences, err := accountDomainUserPreferences.NewDefault(user.ID)
			if err != nil {
				return err
			}

			if err = h.userPreferenceRepository.Upsert(txCtx, preferences); err != nil {
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

			err = h.linkedAccountRepository.Create(txCtx, linkedAccount)
			if err != nil {
				return err
			}

			userID = user.ID

			if err = h.eventBusPublisherService.Enqueue(txCtx, accountDomainUsers.UserCreatedDomainEvent{
				UserID:      user.ID,
				Identifier:  linkedAccount.Identifier,
				DisplayName: authResponse.DisplayName,
			}); err != nil {
				return err
			}
		} else {
			user, err := h.userRepository.GetByID(txCtx, linkedAccount.UserID)
			if err != nil {
				return err
			}

			userID = user.ID
		}

		deletedNonce, err := authResponse.Nonce.Delete()
		if err != nil {
			return err
		}

		if err = h.userNonceRepository.Delete(txCtx, deletedNonce); err != nil {
			return err
		}

		issued, err := h.tokenIssuer.IssueTokens(txCtx, userID, request.AuthProvider, request.IPAddress, request.UserAgent)
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
