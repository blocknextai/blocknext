package auth

import (
	"context"
	"errors"
	"strings"

	"github.com/google/uuid"

	"github.com/blocknextai/go-packages/apperror"
	"github.com/blocknextai/platform-api/internal/common/validation"
	accountApplicationAuth "github.com/blocknextai/platform-api/internal/modules/account/internal/application/auth"
	"github.com/blocknextai/platform-api/internal/modules/account/internal/application/auth/verificationtoken"
	accountDomain "github.com/blocknextai/platform-api/internal/modules/account/internal/domain"
	accountDomainLinkedAccounts "github.com/blocknextai/platform-api/internal/modules/account/internal/domain/linkedaccounts"
	accountDomainPasswordCredentials "github.com/blocknextai/platform-api/internal/modules/account/internal/domain/passwordcredentials"
	accountDomainVerificationTokens "github.com/blocknextai/platform-api/internal/modules/account/internal/domain/verificationtokens"
)

type ResetCommand struct {
	Token       string
	NewPassword string
}

var (
	resetErrTokenIsRequired = apperror.Validation("token is required")
)

func (c *ResetCommand) Validate() error {
	if strings.TrimSpace(c.Token) == "" {
		return resetErrTokenIsRequired
	}
	return validation.NewPassword(c.NewPassword)
}

type ResetResponse struct{}

func (s *Service) ResetPassword(ctx context.Context, command *ResetCommand) (*ResetResponse, error) {
	if err := command.Validate(); err != nil {
		return nil, err
	}

	if err := s.passwordPolicy.Check(ctx, command.NewPassword, nil); err != nil {
		return nil, err
	}

	tokenHash := verificationtoken.Hash(command.Token)

	newHash, err := s.hasher.Generate(command.NewPassword)
	if err != nil {
		return nil, accountApplicationAuth.ErrFailedToHashPassword
	}

	err = s.transactionManager.ExecuteInTransaction(ctx, func(txCtx context.Context) error {
		token, err := s.verificationTokenRepository.GetByTokenHash(txCtx, tokenHash)
		if err != nil {
			return err
		}

		if token.Purpose != accountDomainVerificationTokens.PurposePasswordReset {
			return accountDomainVerificationTokens.ErrInvalidPurpose
		}

		consumed, err := token.Consume()
		if err != nil {
			return err
		}

		if err = s.verificationTokenRepository.Update(txCtx, consumed); err != nil {
			return err
		}

		credential, err := s.passwordCredentialRepository.GetByUserID(txCtx, token.UserID)
		if err != nil {
			if errors.Is(err, accountDomainPasswordCredentials.ErrPasswordCredentialNotFound) {
				return s.setupPassword(txCtx, token.UserID, newHash)
			}
			return err
		}

		updated, err := credential.ChangePassword(newHash)
		if err != nil {
			return err
		}

		return s.passwordCredentialRepository.Update(txCtx, updated)
	})

	if err != nil {
		return nil, err
	}

	return &ResetResponse{}, nil
}

func (s *Service) setupPassword(ctx context.Context, userID uuid.UUID, passwordHash string) error {
	existingLA, err := s.linkedAccountRepository.GetByAuthProviderAndUserID(ctx, accountDomain.AuthProviderPassword, userID)
	if err != nil && !errors.Is(err, accountDomainLinkedAccounts.ErrLinkedAccountNotFound) {
		return err
	}

	if existingLA == nil {
		passwordLA, err := accountDomainLinkedAccounts.NewLinkedAccount(
			userID,
			accountDomain.AuthProviderPassword,
			userID.String(),
			accountDomain.AuthProviderPassword.String(),
			new(accountDomain.AuthProviderPassword.String()),
			false,
			true,
		)
		if err != nil {
			return err
		}

		if err := s.linkedAccountRepository.Create(ctx, passwordLA); err != nil {
			return err
		}
	}

	credential, err := accountDomainPasswordCredentials.NewPasswordCredential(userID, passwordHash)
	if err != nil {
		return err
	}

	return s.passwordCredentialRepository.Create(ctx, credential)
}
