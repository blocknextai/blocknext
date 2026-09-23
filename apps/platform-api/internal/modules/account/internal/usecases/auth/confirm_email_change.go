package auth

import (
	"context"
	"errors"
	"strings"

	"github.com/blocknextai/go-packages/apperror"
	"github.com/blocknextai/platform-api/internal/modules/account/internal/application/auth/verificationtoken"
	accountDomain "github.com/blocknextai/platform-api/internal/modules/account/internal/domain"
	accountDomainLinkedAccounts "github.com/blocknextai/platform-api/internal/modules/account/internal/domain/linkedaccounts"
	accountDomainUsers "github.com/blocknextai/platform-api/internal/modules/account/internal/domain/users"
	accountDomainVerificationTokens "github.com/blocknextai/platform-api/internal/modules/account/internal/domain/verificationtokens"
)

type ConfirmEmailChangeCommand struct {
	Token string
}

var (
	errTokenIsRequired = apperror.Validation("token is required")
)

func (c *ConfirmEmailChangeCommand) Validate() error {
	if strings.TrimSpace(c.Token) == "" {
		return errTokenIsRequired
	}
	return nil
}

type ConfirmEmailChangeResponse struct{}

func (s *Service) ConfirmEmailChange(ctx context.Context, command *ConfirmEmailChangeCommand) (*ConfirmEmailChangeResponse, error) {
	if err := command.Validate(); err != nil {
		return nil, err
	}

	tokenHash := verificationtoken.Hash(command.Token)

	err := s.transactionManager.ExecuteInTransaction(ctx, func(txCtx context.Context) error {
		token, err := s.verificationTokenRepository.GetByTokenHash(txCtx, tokenHash)
		if err != nil {
			return err
		}

		if token.Purpose != accountDomainVerificationTokens.PurposeEmailChange {
			return accountDomainVerificationTokens.ErrInvalidPurpose
		}

		consumed, err := token.Consume()
		if err != nil {
			return err
		}

		if err = s.verificationTokenRepository.Update(txCtx, consumed); err != nil {
			return err
		}

		taken, err := s.linkedAccountRepository.GetByAuthProviderAndProviderID(txCtx, accountDomain.AuthProviderEmail, token.Email)
		if err != nil && !errors.Is(err, accountDomainLinkedAccounts.ErrLinkedAccountNotFound) {
			return err
		}
		if taken != nil && taken.UserID != token.UserID {
			return accountDomainUsers.ErrEmailAlreadyInUse
		}

		linkedAccount, err := s.linkedAccountRepository.GetByAuthProviderAndUserID(txCtx, accountDomain.AuthProviderEmail, token.UserID)
		if err != nil {
			return err
		}

		updated, err := linkedAccount.ChangeEmailAndVerify(token.Email)
		if err != nil {
			return err
		}

		return s.linkedAccountRepository.Update(txCtx, updated)
	})

	if err != nil {
		return nil, err
	}

	return &ConfirmEmailChangeResponse{}, nil
}
