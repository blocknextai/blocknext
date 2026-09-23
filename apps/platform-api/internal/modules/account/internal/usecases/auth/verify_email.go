package auth

import (
	"context"
	"strings"

	"github.com/google/uuid"

	"github.com/blocknextai/go-packages/apperror"
	accountApplicationAuth "github.com/blocknextai/platform-api/internal/modules/account/internal/application/auth"
	"github.com/blocknextai/platform-api/internal/modules/account/internal/application/auth/verificationtoken"
	accountDomain "github.com/blocknextai/platform-api/internal/modules/account/internal/domain"
	accountDomainVerificationTokens "github.com/blocknextai/platform-api/internal/modules/account/internal/domain/verificationtokens"
)

type VerifyCommand struct {
	Token     string
	IPAddress string
	UserAgent string
}

var (
	verifyErrTokenIsRequired = apperror.Validation("token is required")
)

func (c *VerifyCommand) Validate() error {
	if strings.TrimSpace(c.Token) == "" {
		return verifyErrTokenIsRequired
	}
	return nil
}

func (s *Service) VerifyEmail(ctx context.Context, command *VerifyCommand) (*accountApplicationAuth.AccessTokenResponse, error) {
	if err := command.Validate(); err != nil {
		return nil, err
	}

	tokenHash := verificationtoken.Hash(command.Token)

	var userID uuid.UUID

	err := s.transactionManager.ExecuteInTransaction(ctx, func(txCtx context.Context) error {
		token, err := s.verificationTokenRepository.GetByTokenHash(txCtx, tokenHash)
		if err != nil {
			return err
		}

		if token.Purpose != accountDomainVerificationTokens.PurposeEmailVerify {
			return accountDomainVerificationTokens.ErrInvalidPurpose
		}

		consumed, err := token.Consume()
		if err != nil {
			return err
		}

		if err = s.verificationTokenRepository.Update(txCtx, consumed); err != nil {
			return err
		}

		linkedAccount, err := s.linkedAccountRepository.GetByAuthProviderAndProviderID(txCtx, accountDomain.AuthProviderEmail, token.Email)
		if err != nil {
			return err
		}

		if !linkedAccount.IsVerified {
			verified, err := linkedAccount.MarkVerified()
			if err != nil {
				return err
			}

			if err = s.linkedAccountRepository.Update(txCtx, verified); err != nil {
				return err
			}
		}

		userID = token.UserID
		return nil
	})

	if err != nil {
		return nil, err
	}

	return s.tokenIssuer.IssueTokens(ctx, userID, accountDomain.AuthProviderPassword, command.IPAddress, command.UserAgent)
}
