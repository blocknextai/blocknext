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

type ForgotCommand struct {
	Email string
}

func (c *ForgotCommand) Validate() error {
	return validation.Email(c.Email)
}

type ForgotResponse struct{}

func (s *Service) ForgotPassword(ctx context.Context, command *ForgotCommand) (*ForgotResponse, error) {
	if err := command.Validate(); err != nil {
		return nil, err
	}

	normalizedEmail := accountDomain.NormalizeEmail(command.Email)

	linkedAccount, err := s.linkedAccountRepository.GetByAuthProviderAndProviderID(ctx, accountDomain.AuthProviderEmail, normalizedEmail)
	if err != nil {
		if errors.Is(err, accountDomainLinkedAccounts.ErrLinkedAccountNotFound) {
			return &ForgotResponse{}, nil
		}
		return nil, err
	}

	err = s.transactionManager.ExecuteInTransaction(ctx, func(txCtx context.Context) error {
		plainToken, err := s.verificationTokenIssuer.Issue(txCtx, linkedAccount.UserID, accountDomainVerificationTokens.PurposePasswordReset, normalizedEmail, s.resetTokenTTL)
		if err != nil {
			return err
		}

		return s.eventBusPublisherService.Enqueue(txCtx, accountDomainEmails.PasswordResetRequestedDomainEvent{
			Email: normalizedEmail,
			Token: plainToken,
		})
	})
	if err != nil {
		return nil, err
	}

	return &ForgotResponse{}, nil
}
