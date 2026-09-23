package linkedaccounts

import (
	"context"
	"errors"

	"github.com/google/uuid"

	accountApplicationLinkedAccounts "github.com/blocknextai/platform-api/internal/modules/account/internal/application/linkedaccounts"
	accountDomain "github.com/blocknextai/platform-api/internal/modules/account/internal/domain"
	accountDomainLinkedAccounts "github.com/blocknextai/platform-api/internal/modules/account/internal/domain/linkedaccounts"
	accountDomainPasswordCredentials "github.com/blocknextai/platform-api/internal/modules/account/internal/domain/passwordcredentials"
)

type DeleteLinkedAccountCommand struct {
	UserID          uuid.UUID
	LinkedAccountID uuid.UUID
}

type DeleteLinkedAccountResponse struct{}

func (s *Service) DeleteLinkedAccount(ctx context.Context, command *DeleteLinkedAccountCommand) (*DeleteLinkedAccountResponse, error) {
	err := s.transactionManager.ExecuteInTransaction(ctx, func(txCtx context.Context) error {
		linkedAccount, err := s.linkedAccountRepository.GetByIDAndUserID(txCtx, command.LinkedAccountID, command.UserID)
		if err != nil {
			return accountDomainLinkedAccounts.ErrLinkedAccountNotFound
		}

		if linkedAccount.IsPrimary {
			return accountApplicationLinkedAccounts.ErrCannotDeletePrimaryLinkedAccount
		}

		if linkedAccount.AuthProvider == accountDomain.AuthProviderEmail {
			if err := s.deletePasswordLinkedAccount(txCtx, command.UserID); err != nil {
				return err
			}
		}

		if linkedAccount.AuthProvider == accountDomain.AuthProviderEmail || linkedAccount.AuthProvider == accountDomain.AuthProviderPassword {
			if err := s.deletePasswordCredential(txCtx, command.UserID); err != nil {
				return err
			}
		}

		deletedAccount, err := linkedAccount.Delete()
		if err != nil {
			return accountApplicationLinkedAccounts.ErrFailedToDeleteLinkedAccount.WithCause(err)
		}

		if err := s.linkedAccountRepository.Delete(txCtx, deletedAccount); err != nil {
			return accountApplicationLinkedAccounts.ErrFailedToDeleteLinkedAccount.WithCause(err)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &DeleteLinkedAccountResponse{}, nil
}

func (s *Service) deletePasswordLinkedAccount(ctx context.Context, userID uuid.UUID) error {
	passwordLinkedAccount, err := s.linkedAccountRepository.GetByAuthProviderAndUserID(ctx, accountDomain.AuthProviderPassword, userID)
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

	if err := s.linkedAccountRepository.Delete(ctx, deleted); err != nil {
		return accountApplicationLinkedAccounts.ErrFailedToDeleteLinkedAccount.WithCause(err)
	}

	return nil
}

func (s *Service) deletePasswordCredential(ctx context.Context, userID uuid.UUID) error {
	credential, err := s.passwordCredentialRepository.GetByUserID(ctx, userID)
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

	return s.passwordCredentialRepository.Delete(ctx, deleted)
}
