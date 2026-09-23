package auth

import (
	"context"
	"errors"

	"github.com/google/uuid"

	"github.com/blocknextai/go-packages/apperror"
	"github.com/blocknextai/platform-api/internal/common/validation"
	accountApplicationAuth "github.com/blocknextai/platform-api/internal/modules/account/internal/application/auth"
	accountDomain "github.com/blocknextai/platform-api/internal/modules/account/internal/domain"
	accountDomainLinkedAccounts "github.com/blocknextai/platform-api/internal/modules/account/internal/domain/linkedaccounts"
	accountDomainPasswordCredentials "github.com/blocknextai/platform-api/internal/modules/account/internal/domain/passwordcredentials"
	accountDomainUsers "github.com/blocknextai/platform-api/internal/modules/account/internal/domain/users"
)

type ChangePasswordCommand struct {
	UserID          uuid.UUID
	CurrentPassword string
	NewPassword     string
}

var ErrSameAsCurrent = apperror.Validation("new password must differ from current password")

func (c *ChangePasswordCommand) Validate() error {
	if err := validation.Password(c.CurrentPassword); err != nil {
		return err
	}
	if err := validation.NewPassword(c.NewPassword); err != nil {
		return err
	}
	if c.CurrentPassword == c.NewPassword {
		return ErrSameAsCurrent
	}
	return nil
}

type ChangePasswordResponse struct{}

func (s *Service) ChangePassword(ctx context.Context, command *ChangePasswordCommand) (*ChangePasswordResponse, error) {
	if err := command.Validate(); err != nil {
		return nil, err
	}

	if _, err := s.linkedAccountRepository.GetByAuthProviderAndUserID(ctx, accountDomain.AuthProviderEmail, command.UserID); err != nil {
		if errors.Is(err, accountDomainLinkedAccounts.ErrLinkedAccountNotFound) {
			return nil, accountApplicationAuth.ErrPasswordAuthDisabled
		}
		return nil, err
	}

	credential, err := s.passwordCredentialRepository.GetByUserID(ctx, command.UserID)
	if err != nil {
		if errors.Is(err, accountDomainPasswordCredentials.ErrPasswordCredentialNotFound) {
			return nil, accountApplicationAuth.ErrPasswordNotSet
		}
		return nil, err
	}

	if !s.hasher.Compare(command.CurrentPassword, credential.PasswordHash) {
		return nil, accountApplicationAuth.ErrInvalidEmailOrPassword
	}

	if err := s.passwordPolicy.Check(ctx, command.NewPassword, nil); err != nil {
		return nil, err
	}

	newHash, err := s.hasher.Generate(command.NewPassword)
	if err != nil {
		return nil, accountApplicationAuth.ErrFailedToHashPassword
	}

	err = s.transactionManager.ExecuteInTransaction(ctx, func(txCtx context.Context) error {
		updated, err := credential.ChangePassword(newHash)
		if err != nil {
			return err
		}

		if err := s.passwordCredentialRepository.Update(txCtx, updated); err != nil {
			return err
		}

		return s.eventBusPublisherService.Enqueue(txCtx, accountDomainUsers.PasswordChangedDomainEvent{
			UserID: command.UserID,
		})
	})

	if err != nil {
		return nil, err
	}

	return &ChangePasswordResponse{}, nil
}
