package auth

import (
	"context"
	"errors"

	"github.com/google/uuid"

	"github.com/blocknextai/go-packages/apperror"
	"github.com/blocknextai/platform-api/internal/common/validation"
	accountApplicationAuth "github.com/blocknextai/platform-api/internal/modules/account/internal/application/auth"
	accountDomain "github.com/blocknextai/platform-api/internal/modules/account/internal/domain"
	accountDomainEmails "github.com/blocknextai/platform-api/internal/modules/account/internal/domain/emails"
	accountDomainLinkedAccounts "github.com/blocknextai/platform-api/internal/modules/account/internal/domain/linkedaccounts"
	accountDomainPasswordCredentials "github.com/blocknextai/platform-api/internal/modules/account/internal/domain/passwordcredentials"
	accountDomainUsers "github.com/blocknextai/platform-api/internal/modules/account/internal/domain/users"
	accountDomainVerificationTokens "github.com/blocknextai/platform-api/internal/modules/account/internal/domain/verificationtokens"
)

type ChangeEmailCommand struct {
	UserID          uuid.UUID
	NewEmail        string
	CurrentPassword string
}

var ErrSameEmail = apperror.Validation("new email is the same as current email")

func (c *ChangeEmailCommand) Validate() error {
	if err := validation.Email(c.NewEmail); err != nil {
		return err
	}
	return validation.Password(c.CurrentPassword)
}

type ChangeEmailResponse struct{}

func (s *Service) ChangeEmail(ctx context.Context, command *ChangeEmailCommand) (*ChangeEmailResponse, error) {
	if err := command.Validate(); err != nil {
		return nil, err
	}

	normalizedNewEmail := accountDomain.NormalizeEmail(command.NewEmail)

	currentLinkedAccount, err := s.linkedAccountRepository.GetByAuthProviderAndUserID(ctx, accountDomain.AuthProviderEmail, command.UserID)
	if err != nil {
		if errors.Is(err, accountDomainLinkedAccounts.ErrLinkedAccountNotFound) {
			return nil, accountApplicationAuth.ErrPasswordAuthDisabled
		}
		return nil, err
	}

	if currentLinkedAccount.ProviderID == normalizedNewEmail {
		return nil, ErrSameEmail
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

	existing, err := s.linkedAccountRepository.GetByAuthProviderAndProviderID(ctx, accountDomain.AuthProviderEmail, normalizedNewEmail)
	if err != nil && !errors.Is(err, accountDomainLinkedAccounts.ErrLinkedAccountNotFound) {
		return nil, err
	}
	if existing != nil {
		return nil, accountDomainUsers.ErrEmailAlreadyInUse
	}

	err = s.transactionManager.ExecuteInTransaction(ctx, func(txCtx context.Context) error {
		plainToken, err := s.verificationTokenIssuer.Issue(txCtx, command.UserID, accountDomainVerificationTokens.PurposeEmailChange, normalizedNewEmail, s.tokenTTL)
		if err != nil {
			return err
		}

		if err := s.eventBusPublisherService.Enqueue(txCtx, accountDomainUsers.EmailChangedDomainEvent{
			UserID:   command.UserID,
			OldEmail: currentLinkedAccount.ProviderID,
			NewEmail: normalizedNewEmail,
		}); err != nil {
			return err
		}

		return s.eventBusPublisherService.Enqueue(txCtx, accountDomainEmails.EmailChangeRequestedDomainEvent{
			Email: normalizedNewEmail,
			Token: plainToken,
		})
	})
	if err != nil {
		return nil, err
	}

	return &ChangeEmailResponse{}, nil
}
