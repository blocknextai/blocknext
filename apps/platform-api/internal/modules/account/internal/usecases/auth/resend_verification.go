package auth

import (
	"context"
	"errors"

	"github.com/blocknextai/platform-api/internal/common/validation"
	accountDomain "github.com/blocknextai/platform-api/internal/modules/account/internal/domain"
	accountDomainEmails "github.com/blocknextai/platform-api/internal/modules/account/internal/domain/emails"
	accountDomainLinkedAccounts "github.com/blocknextai/platform-api/internal/modules/account/internal/domain/linkedaccounts"
	accountDomainVerificationTokens "github.com/blocknextai/platform-api/internal/modules/account/internal/domain/verificationtokens"
)

type ResendVerificationCommand struct {
	Email string
}

func (c *ResendVerificationCommand) Validate() error {
	return validation.Email(c.Email)
}

type ResendVerificationResponse struct{}

func (s *Service) ResendVerification(ctx context.Context, command *ResendVerificationCommand) (*ResendVerificationResponse, error) {
	if err := command.Validate(); err != nil {
		return nil, err
	}

	normalizedEmail := accountDomain.NormalizeEmail(command.Email)

	linkedAccount, err := s.linkedAccountRepository.GetByAuthProviderAndProviderID(ctx, accountDomain.AuthProviderEmail, normalizedEmail)
	if err != nil {
		if errors.Is(err, accountDomainLinkedAccounts.ErrLinkedAccountNotFound) {
			return &ResendVerificationResponse{}, nil
		}
		return nil, err
	}

	if linkedAccount.IsVerified {
		return &ResendVerificationResponse{}, nil
	}

	err = s.transactionManager.ExecuteInTransaction(ctx, func(txCtx context.Context) error {
		plainToken, err := s.verificationTokenIssuer.Issue(txCtx, linkedAccount.UserID, accountDomainVerificationTokens.PurposeEmailVerify, normalizedEmail, s.verificationTokenTTL)
		if err != nil {
			return err
		}

		return s.eventBusPublisherService.Enqueue(txCtx, accountDomainEmails.EmailVerificationRequestedDomainEvent{
			Email: normalizedEmail,
			Token: plainToken,
		})
	})
	if err != nil {
		return nil, err
	}

	return &ResendVerificationResponse{}, nil
}
