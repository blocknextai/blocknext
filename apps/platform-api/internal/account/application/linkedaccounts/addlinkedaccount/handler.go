package addlinkedaccount

import (
	"context"
	"errors"

	"github.com/blocknextai/go-packages/database"
	"github.com/blocknextai/platform-api/internal/account/application/auth/createusertoken"
	accountApplicationLinkedAccounts "github.com/blocknextai/platform-api/internal/account/application/linkedaccounts"
	accountDomainLinkedAccounts "github.com/blocknextai/platform-api/internal/account/domain/linkedaccounts"
	accountDomainUserNonces "github.com/blocknextai/platform-api/internal/account/domain/usernonces"
	accountDomainUsers "github.com/blocknextai/platform-api/internal/account/domain/users"
)

type Handler struct {
	userRepository          accountDomainUsers.UserRepository
	userNonceRepository     accountDomainUserNonces.UserNonceRepository
	linkedAccountRepository accountDomainLinkedAccounts.LinkedAccountRepository
	authProviderRegistry    createusertoken.AuthProviderRegistry
	transactionManager      database.TransactionManager
}

func New(
	userRepository accountDomainUsers.UserRepository,
	userNonceRepository accountDomainUserNonces.UserNonceRepository,
	linkedAccountRepository accountDomainLinkedAccounts.LinkedAccountRepository,
	authProviderRegistry createusertoken.AuthProviderRegistry,
	transactionManager database.TransactionManager,
) *Handler {
	return &Handler{
		userRepository:          userRepository,
		userNonceRepository:     userNonceRepository,
		linkedAccountRepository: linkedAccountRepository,
		authProviderRegistry:    authProviderRegistry,
		transactionManager:      transactionManager,
	}
}

func (h *Handler) Handle(ctx context.Context, command *AddLinkedAccountCommand) (*AddLinkedAccountResponse, error) {
	var response *AddLinkedAccountResponse

	err := h.transactionManager.ExecuteInTransaction(ctx, func(txCtx context.Context) error {
		user, err := h.userRepository.GetByID(txCtx, command.UserID)
		if err != nil {
			return accountDomainUsers.ErrUserNotFound
		}

		authProvider, err := h.authProviderRegistry.GetProvider(command.AuthProvider)
		if err != nil {
			return err
		}

		authResponse, err := authProvider.Validate(txCtx, createusertoken.Request{
			AuthProvider: command.AuthProvider,
			Payload:      command.Payload,
		})
		if err != nil {
			return err
		}

		existingLinkedAccount, err := h.linkedAccountRepository.GetByProviderID(txCtx, authResponse.ProviderID)
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

		if err = h.userNonceRepository.Delete(txCtx, deletedNonce); err != nil {
			return err
		}

		err = h.linkedAccountRepository.Create(txCtx, linkedAccount)
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
