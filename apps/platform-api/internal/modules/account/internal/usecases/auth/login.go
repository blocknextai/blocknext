package auth

import (
	"context"
	"errors"
	"log/slog"

	"github.com/blocknextai/platform-api/internal/common/validation"
	accountApplicationAuth "github.com/blocknextai/platform-api/internal/modules/account/internal/application/auth"
	accountDomain "github.com/blocknextai/platform-api/internal/modules/account/internal/domain"
	accountDomainLinkedAccounts "github.com/blocknextai/platform-api/internal/modules/account/internal/domain/linkedaccounts"
	accountDomainPasswordCredentials "github.com/blocknextai/platform-api/internal/modules/account/internal/domain/passwordcredentials"
)

type LoginCommand struct {
	Email     string
	Password  string
	IPAddress string
	UserAgent string
}

func (c *LoginCommand) Validate() error {
	if err := validation.Email(c.Email); err != nil {
		return err
	}
	return validation.Password(c.Password)
}

const (
	requireVerifiedMail = false
	dummyPassword       = "constant-time-dummy"
)

func (s *Service) Login(ctx context.Context, command *LoginCommand) (*accountApplicationAuth.AccessTokenResponse, error) {
	if err := command.Validate(); err != nil {
		return nil, err
	}

	normalizedEmail := accountDomain.NormalizeEmail(command.Email)

	linkedAccount, hashToCompare, err := s.resolveCredential(ctx, normalizedEmail)
	if err != nil {
		return nil, err
	}

	passwordMatches := s.hasher.Compare(command.Password, hashToCompare)

	if linkedAccount == nil || !passwordMatches {
		return nil, accountApplicationAuth.ErrInvalidEmailOrPassword
	}

	if requireVerifiedMail && !linkedAccount.IsVerified {
		return nil, accountApplicationAuth.ErrEmailNotVerified
	}

	return s.tokenIssuer.IssueTokens(ctx, linkedAccount.UserID, accountDomain.AuthProviderPassword, command.IPAddress, command.UserAgent)
}

func (s *Service) resolveCredential(ctx context.Context, email string) (*accountDomainLinkedAccounts.LinkedAccount, string, error) {
	linkedAccount, err := s.linkedAccountRepository.GetByAuthProviderAndProviderID(ctx, accountDomain.AuthProviderEmail, email)
	if err != nil {
		if errors.Is(err, accountDomainLinkedAccounts.ErrLinkedAccountNotFound) {
			return nil, s.getDummyHash(), nil
		}
		return nil, "", err
	}

	credential, err := s.passwordCredentialRepository.GetByUserID(ctx, linkedAccount.UserID)
	if err != nil {
		if errors.Is(err, accountDomainPasswordCredentials.ErrPasswordCredentialNotFound) {
			return nil, s.getDummyHash(), nil
		}
		return nil, "", err
	}

	return linkedAccount, credential.PasswordHash, nil
}

func (s *Service) getDummyHash() string {
	s.dummyHashOnce.Do(func() {
		hash, err := s.hasher.Generate(dummyPassword)
		if err != nil {
			slog.Error("Failed to seed dummy bcrypt hash; constant-time path degraded",
				"component", "login.Handler",
				"error", err)
			return
		}
		s.dummyHash = hash
	})
	return s.dummyHash
}
