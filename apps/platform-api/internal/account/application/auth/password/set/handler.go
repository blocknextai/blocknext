package set

import (
	"context"
	"errors"

	"github.com/blocknextai/go-packages/database"
	"github.com/blocknextai/go-packages/hashing"
	accountApplicationAuth "github.com/blocknextai/platform-api/internal/account/application/auth"
	"github.com/blocknextai/platform-api/internal/account/application/auth/passwordpolicy"
	accountDomain "github.com/blocknextai/platform-api/internal/account/domain"
	accountDomainLinkedAccounts "github.com/blocknextai/platform-api/internal/account/domain/linkedaccounts"
	accountDomainPasswordCredentials "github.com/blocknextai/platform-api/internal/account/domain/passwordcredentials"
)

type Handler struct {
	linkedAccountRepository      accountDomainLinkedAccounts.LinkedAccountRepository
	passwordCredentialRepository accountDomainPasswordCredentials.PasswordCredentialRepository
	passwordPolicy               passwordpolicy.Policy
	hasher                       hashing.Hasher
	transactionManager           database.TransactionManager
}

func New(
	linkedAccountRepository accountDomainLinkedAccounts.LinkedAccountRepository,
	passwordCredentialRepository accountDomainPasswordCredentials.PasswordCredentialRepository,
	passwordPolicy passwordpolicy.Policy,
	hasher hashing.Hasher,
	transactionManager database.TransactionManager,
) *Handler {
	return &Handler{
		linkedAccountRepository:      linkedAccountRepository,
		passwordCredentialRepository: passwordCredentialRepository,
		passwordPolicy:               passwordPolicy,
		hasher:                       hasher,
		transactionManager:           transactionManager,
	}
}

func (h *Handler) Handle(ctx context.Context, command *SetPasswordCommand) (*SetPasswordResponse, error) {
	emailLinkedAccount, err := h.linkedAccountRepository.GetByAuthProviderAndUserID(ctx, accountDomain.AuthProviderEmail, command.UserID)
	if err != nil {
		if errors.Is(err, accountDomainLinkedAccounts.ErrLinkedAccountNotFound) {
			return nil, ErrEmailRequired
		}
		return nil, err
	}

	existingPassword, err := h.linkedAccountRepository.GetByAuthProviderAndUserID(ctx, accountDomain.AuthProviderPassword, command.UserID)
	if err != nil && !errors.Is(err, accountDomainLinkedAccounts.ErrLinkedAccountNotFound) {
		return nil, err
	}
	if existingPassword != nil {
		return nil, ErrPasswordAlreadyExists
	}

	if err := h.passwordPolicy.Check(ctx, command.Password, []string{emailLinkedAccount.Identifier}); err != nil {
		return nil, err
	}

	passwordHash, err := h.hasher.Generate(command.Password)
	if err != nil {
		return nil, accountApplicationAuth.ErrFailedToHashPassword
	}

	err = h.transactionManager.ExecuteInTransaction(ctx, func(txCtx context.Context) error {
		passwordLinkedAccount, err := accountDomainLinkedAccounts.NewLinkedAccount(
			command.UserID,
			accountDomain.AuthProviderPassword,
			command.UserID.String(),
			accountDomain.AuthProviderPassword.String(),
			new(accountDomain.AuthProviderPassword.String()),
			false,
			true,
		)
		if err != nil {
			return err
		}

		if err = h.linkedAccountRepository.Create(txCtx, passwordLinkedAccount); err != nil {
			return err
		}

		credential, err := accountDomainPasswordCredentials.NewPasswordCredential(command.UserID, passwordHash)
		if err != nil {
			return err
		}

		return h.passwordCredentialRepository.Create(txCtx, credential)
	})

	if err != nil {
		return nil, err
	}

	return &SetPasswordResponse{}, nil
}
