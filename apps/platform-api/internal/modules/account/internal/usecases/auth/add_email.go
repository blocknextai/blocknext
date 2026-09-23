package auth

import (
	"context"
	"errors"

	"github.com/google/uuid"

	"github.com/blocknextai/go-packages/apperror"
	"github.com/blocknextai/platform-api/internal/common/validation"
	accountDomain "github.com/blocknextai/platform-api/internal/modules/account/internal/domain"
	accountDomainEmails "github.com/blocknextai/platform-api/internal/modules/account/internal/domain/emails"
	accountDomainLinkedAccounts "github.com/blocknextai/platform-api/internal/modules/account/internal/domain/linkedaccounts"
	accountDomainUsers "github.com/blocknextai/platform-api/internal/modules/account/internal/domain/users"
	accountDomainVerificationTokens "github.com/blocknextai/platform-api/internal/modules/account/internal/domain/verificationtokens"
)

type AddEmailCommand struct {
	UserID uuid.UUID
	Email  string
}

func (c *AddEmailCommand) Validate() error {
	return validation.Email(c.Email)
}

var (
	ErrEmailAlreadyLinked = apperror.Conflict("email already linked to this account")
)

type AddEmailResponse struct{}

func (s *Service) AddEmail(ctx context.Context, command *AddEmailCommand) (*AddEmailResponse, error) {
	if err := command.Validate(); err != nil {
		return nil, err
	}

	normalizedEmail := accountDomain.NormalizeEmail(command.Email)

	existingForUser, err := s.linkedAccountRepository.GetByAuthProviderAndUserID(ctx, accountDomain.AuthProviderEmail, command.UserID)
	if err != nil && !errors.Is(err, accountDomainLinkedAccounts.ErrLinkedAccountNotFound) {
		return nil, err
	}
	if existingForUser != nil {
		return nil, ErrEmailAlreadyLinked
	}

	emailTaken, err := s.linkedAccountRepository.GetByAuthProviderAndProviderID(ctx, accountDomain.AuthProviderEmail, normalizedEmail)
	if err != nil && !errors.Is(err, accountDomainLinkedAccounts.ErrLinkedAccountNotFound) {
		return nil, err
	}
	if emailTaken != nil {
		return nil, accountDomainUsers.ErrEmailAlreadyInUse
	}

	err = s.transactionManager.ExecuteInTransaction(ctx, func(txCtx context.Context) error {
		linkedAccount, err := accountDomainLinkedAccounts.NewLinkedAccount(
			command.UserID,
			accountDomain.AuthProviderEmail,
			normalizedEmail,
			normalizedEmail,
			new(normalizedEmail),
			false,
			false,
		)
		if err != nil {
			return err
		}

		if err = s.linkedAccountRepository.Create(txCtx, linkedAccount); err != nil {
			return err
		}

		plainToken, err := s.verificationTokenIssuer.Issue(txCtx, command.UserID, accountDomainVerificationTokens.PurposeEmailVerify, normalizedEmail, s.verificationTokenTTL)
		if err != nil {
			return err
		}

		return s.eventBusPublisherService.Enqueue(txCtx, accountDomainEmails.EmailAddedDomainEvent{
			Email: normalizedEmail,
			Token: plainToken,
		})
	})

	if err != nil {
		return nil, err
	}

	return &AddEmailResponse{}, nil
}
