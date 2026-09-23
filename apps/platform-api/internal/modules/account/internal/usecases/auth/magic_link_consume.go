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

type MagicLinkConsumeCommand struct {
	Token     string
	IPAddress string
	UserAgent string
}

var (
	magicLinkConsumeErrTokenIsRequired = apperror.Validation("token is required")
)

func (c *MagicLinkConsumeCommand) Validate() error {
	if strings.TrimSpace(c.Token) == "" {
		return magicLinkConsumeErrTokenIsRequired
	}
	return nil
}

func (s *Service) MagicLinkConsume(ctx context.Context, command *MagicLinkConsumeCommand) (*accountApplicationAuth.AccessTokenResponse, error) {
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

		if token.Purpose != accountDomainVerificationTokens.PurposeMagicLink {
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

		userID = linkedAccount.UserID
		return nil
	})

	if err != nil {
		return nil, err
	}

	return s.tokenIssuer.IssueTokens(ctx, userID, accountDomain.AuthProviderEmail, command.IPAddress, command.UserAgent)
}
