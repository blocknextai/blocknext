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

type MagicLinkRequestCommand struct {
	Email string
}

func (c *MagicLinkRequestCommand) Validate() error {
	return validation.Email(c.Email)
}

type MagicLinkRequestResponse struct{}

func (s *Service) MagicLinkRequest(ctx context.Context, command *MagicLinkRequestCommand) (*MagicLinkRequestResponse, error) {
	if err := command.Validate(); err != nil {
		return nil, err
	}

	normalizedEmail := accountDomain.NormalizeEmail(command.Email)

	linkedAccount, err := s.linkedAccountRepository.GetByAuthProviderAndProviderID(ctx, accountDomain.AuthProviderEmail, normalizedEmail)
	if err != nil {
		if errors.Is(err, accountDomainLinkedAccounts.ErrLinkedAccountNotFound) {
			return &MagicLinkRequestResponse{}, nil
		}
		return nil, err
	}

	err = s.transactionManager.ExecuteInTransaction(ctx, func(txCtx context.Context) error {
		plainToken, err := s.verificationTokenIssuer.Issue(txCtx, linkedAccount.UserID, accountDomainVerificationTokens.PurposeMagicLink, normalizedEmail, s.tokenTTL)
		if err != nil {
			return err
		}

		return s.eventBusPublisherService.Enqueue(txCtx, accountDomainEmails.MagicLinkCreatedDomainEvent{
			Email: normalizedEmail,
			Token: plainToken,
		})
	})
	if err != nil {
		return nil, err
	}

	return &MagicLinkRequestResponse{}, nil
}
