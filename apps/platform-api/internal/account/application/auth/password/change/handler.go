package change

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
	accountDomainUsers "github.com/blocknextai/platform-api/internal/account/domain/users"
	"github.com/blocknextai/platform-api/internal/eventbus/application/publishing"
)

type Handler struct {
	linkedAccountRepository      accountDomainLinkedAccounts.LinkedAccountRepository
	passwordCredentialRepository accountDomainPasswordCredentials.PasswordCredentialRepository
	passwordPolicy               passwordpolicy.Policy
	hasher                       hashing.Hasher
	eventBusPublisherService     publishing.PublisherService
	transactionManager           database.TransactionManager
}

func New(
	linkedAccountRepository accountDomainLinkedAccounts.LinkedAccountRepository,
	passwordCredentialRepository accountDomainPasswordCredentials.PasswordCredentialRepository,
	passwordPolicy passwordpolicy.Policy,
	hasher hashing.Hasher,
	eventBusPublisherService publishing.PublisherService,
	transactionManager database.TransactionManager,
) *Handler {
	return &Handler{
		linkedAccountRepository:      linkedAccountRepository,
		passwordCredentialRepository: passwordCredentialRepository,
		passwordPolicy:               passwordPolicy,
		hasher:                       hasher,
		eventBusPublisherService:     eventBusPublisherService,
		transactionManager:           transactionManager,
	}
}

func (h *Handler) Handle(ctx context.Context, command *ChangePasswordCommand) (*ChangePasswordResponse, error) {
	if _, err := h.linkedAccountRepository.GetByAuthProviderAndUserID(ctx, accountDomain.AuthProviderEmail, command.UserID); err != nil {
		if errors.Is(err, accountDomainLinkedAccounts.ErrLinkedAccountNotFound) {
			return nil, accountApplicationAuth.ErrPasswordAuthDisabled
		}
		return nil, err
	}

	credential, err := h.passwordCredentialRepository.GetByUserID(ctx, command.UserID)
	if err != nil {
		if errors.Is(err, accountDomainPasswordCredentials.ErrPasswordCredentialNotFound) {
			return nil, accountApplicationAuth.ErrPasswordNotSet
		}
		return nil, err
	}

	if !h.hasher.Compare(command.CurrentPassword, credential.PasswordHash) {
		return nil, accountApplicationAuth.ErrInvalidEmailOrPassword
	}

	if err := h.passwordPolicy.Check(ctx, command.NewPassword, nil); err != nil {
		return nil, err
	}

	newHash, err := h.hasher.Generate(command.NewPassword)
	if err != nil {
		return nil, accountApplicationAuth.ErrFailedToHashPassword
	}

	err = h.transactionManager.ExecuteInTransaction(ctx, func(txCtx context.Context) error {
		updated, err := credential.ChangePassword(newHash)
		if err != nil {
			return err
		}

		if err := h.passwordCredentialRepository.Update(txCtx, updated); err != nil {
			return err
		}

		return h.eventBusPublisherService.Enqueue(txCtx, accountDomainUsers.PasswordChangedDomainEvent{
			UserID: command.UserID,
		})
	})

	if err != nil {
		return nil, err
	}

	return &ChangePasswordResponse{}, nil
}
