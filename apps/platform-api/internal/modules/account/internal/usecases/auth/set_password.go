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
)

type SetPasswordCommand struct {
	UserID   uuid.UUID
	Password string
}

func (c *SetPasswordCommand) Validate() error {
	return validation.NewPassword(c.Password)
}

var (
	ErrPasswordAlreadyExists = apperror.Conflict("password already set; use change-password instead")
	ErrEmailRequired         = apperror.Conflict("email must be linked before setting a password")
)

type SetPasswordResponse struct{}

func (s *Service) SetPassword(ctx context.Context, command *SetPasswordCommand) (*SetPasswordResponse, error) {
	if err := command.Validate(); err != nil {
		return nil, err
	}

	emailLinkedAccount, err := s.linkedAccountRepository.GetByAuthProviderAndUserID(ctx, accountDomain.AuthProviderEmail, command.UserID)
	if err != nil {
		if errors.Is(err, accountDomainLinkedAccounts.ErrLinkedAccountNotFound) {
			return nil, ErrEmailRequired
		}
		return nil, err
	}

	existingPassword, err := s.linkedAccountRepository.GetByAuthProviderAndUserID(ctx, accountDomain.AuthProviderPassword, command.UserID)
	if err != nil && !errors.Is(err, accountDomainLinkedAccounts.ErrLinkedAccountNotFound) {
		return nil, err
	}
	if existingPassword != nil {
		return nil, ErrPasswordAlreadyExists
	}

	if err := s.passwordPolicy.Check(ctx, command.Password, []string{emailLinkedAccount.Identifier}); err != nil {
		return nil, err
	}

	passwordHash, err := s.hasher.Generate(command.Password)
	if err != nil {
		return nil, accountApplicationAuth.ErrFailedToHashPassword
	}

	err = s.transactionManager.ExecuteInTransaction(ctx, func(txCtx context.Context) error {
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

		if err = s.linkedAccountRepository.Create(txCtx, passwordLinkedAccount); err != nil {
			return err
		}

		credential, err := accountDomainPasswordCredentials.NewPasswordCredential(command.UserID, passwordHash)
		if err != nil {
			return err
		}

		return s.passwordCredentialRepository.Create(txCtx, credential)
	})

	if err != nil {
		return nil, err
	}

	return &SetPasswordResponse{}, nil
}
