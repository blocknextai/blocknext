package deletelinkedaccount

import (
	"context"
	"errors"

	"github.com/blocknextai/go-packages/database"
	accountApplicationLinkedAccounts "github.com/blocknextai/platform-api/internal/account/application/linkedaccounts"
	accountDomain "github.com/blocknextai/platform-api/internal/account/domain"
	accountDomainLinkedAccounts "github.com/blocknextai/platform-api/internal/account/domain/linkedaccounts"
	accountDomainPasswordCredentials "github.com/blocknextai/platform-api/internal/account/domain/passwordcredentials"
	"github.com/google/uuid"
)

type Handler struct {
	linkedAccountRepository      accountDomainLinkedAccounts.LinkedAccountRepository
	passwordCredentialRepository accountDomainPasswordCredentials.PasswordCredentialRepository
	transactionManager           database.TransactionManager
}

func New(
	linkedAccountRepository accountDomainLinkedAccounts.LinkedAccountRepository,
	passwordCredentialRepository accountDomainPasswordCredentials.PasswordCredentialRepository,
	transactionManager database.TransactionManager,
) *Handler {
	return &Handler{
		linkedAccountRepository:      linkedAccountRepository,
		passwordCredentialRepository: passwordCredentialRepository,
		transactionManager:           transactionManager,
	}
}

func (h *Handler) Handle(ctx context.Context, command *DeleteLinkedAccountCommand) (*DeleteLinkedAccountResponse, error) {
	err := h.transactionManager.ExecuteInTransaction(ctx, func(txCtx context.Context) error {
		linkedAccount, err := h.linkedAccountRepository.GetByIDAndUserID(txCtx, command.LinkedAccountID, command.UserID)
		if err != nil {
			return accountDomainLinkedAccounts.ErrLinkedAccountNotFound
		}

		if linkedAccount.IsPrimary {
			return accountApplicationLinkedAccounts.ErrCannotDeletePrimaryLinkedAccount
		}

		if linkedAccount.AuthProvider == accountDomain.AuthProviderEmail {
			if err := h.deletePasswordLinkedAccount(txCtx, command.UserID); err != nil {
				return err
			}
		}

		if linkedAccount.AuthProvider == accountDomain.AuthProviderEmail || linkedAccount.AuthProvider == accountDomain.AuthProviderPassword {
			if err := h.deletePasswordCredential(txCtx, command.UserID); err != nil {
				return err
			}
		}

		deletedAccount, err := linkedAccount.Delete()
		if err != nil {
			return accountApplicationLinkedAccounts.ErrFailedToDeleteLinkedAccount.WithCause(err)
		}

		if err := h.linkedAccountRepository.Delete(txCtx, deletedAccount); err != nil {
			return accountApplicationLinkedAccounts.ErrFailedToDeleteLinkedAccount.WithCause(err)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &DeleteLinkedAccountResponse{}, nil
}

func (h *Handler) deletePasswordLinkedAccount(ctx context.Context, userID uuid.UUID) error {
	passwordLinkedAccount, err := h.linkedAccountRepository.GetByAuthProviderAndUserID(ctx, accountDomain.AuthProviderPassword, userID)
	if err != nil {
		if errors.Is(err, accountDomainLinkedAccounts.ErrLinkedAccountNotFound) {
			return nil
		}
		return err
	}

	deleted, err := passwordLinkedAccount.Delete()
	if err != nil {
		return accountApplicationLinkedAccounts.ErrFailedToDeleteLinkedAccount.WithCause(err)
	}

	if err := h.linkedAccountRepository.Delete(ctx, deleted); err != nil {
		return accountApplicationLinkedAccounts.ErrFailedToDeleteLinkedAccount.WithCause(err)
	}

	return nil
}

func (h *Handler) deletePasswordCredential(ctx context.Context, userID uuid.UUID) error {
	credential, err := h.passwordCredentialRepository.GetByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, accountDomainPasswordCredentials.ErrPasswordCredentialNotFound) {
			return nil
		}
		return err
	}

	deleted, err := credential.Delete()
	if err != nil {
		return err
	}

	return h.passwordCredentialRepository.Delete(ctx, deleted)
}
