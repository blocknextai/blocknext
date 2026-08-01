package login

import (
	"context"
	"errors"
	"log/slog"
	"sync"

	"github.com/blocknextai/go-packages/database"
	"github.com/blocknextai/go-packages/hashing"
	accountApplicationAuth "github.com/blocknextai/platform-api/internal/account/application/auth"
	"github.com/blocknextai/platform-api/internal/account/application/auth/tokenissuer"
	accountDomain "github.com/blocknextai/platform-api/internal/account/domain"
	accountDomainLinkedAccounts "github.com/blocknextai/platform-api/internal/account/domain/linkedaccounts"
	accountDomainPasswordCredentials "github.com/blocknextai/platform-api/internal/account/domain/passwordcredentials"
)

const (
	requireVerifiedMail = false
	dummyPassword       = "constant-time-dummy"
)

type Handler struct {
	linkedAccountRepository      accountDomainLinkedAccounts.LinkedAccountRepository
	passwordCredentialRepository accountDomainPasswordCredentials.PasswordCredentialRepository
	tokenIssuer                  tokenissuer.Service
	hasher                       hashing.Hasher
	transactionManager           database.TransactionManager

	dummyHashOnce sync.Once
	dummyHash     string
}

func New(
	linkedAccountRepository accountDomainLinkedAccounts.LinkedAccountRepository,
	passwordCredentialRepository accountDomainPasswordCredentials.PasswordCredentialRepository,
	tokenIssuer tokenissuer.Service,
	hasher hashing.Hasher,
	transactionManager database.TransactionManager,
) *Handler {
	return &Handler{
		linkedAccountRepository:      linkedAccountRepository,
		passwordCredentialRepository: passwordCredentialRepository,
		tokenIssuer:                  tokenIssuer,
		hasher:                       hasher,
		transactionManager:           transactionManager,
	}
}

func (h *Handler) Handle(ctx context.Context, command *LoginCommand) (*accountApplicationAuth.AccessTokenResponse, error) {
	normalizedEmail := accountDomain.NormalizeEmail(command.Email)

	linkedAccount, hashToCompare, err := h.resolveCredential(ctx, normalizedEmail)
	if err != nil {
		return nil, err
	}

	passwordMatches := h.hasher.Compare(command.Password, hashToCompare)

	if linkedAccount == nil || !passwordMatches {
		return nil, accountApplicationAuth.ErrInvalidEmailOrPassword
	}

	if requireVerifiedMail && !linkedAccount.IsVerified {
		return nil, accountApplicationAuth.ErrEmailNotVerified
	}

	return h.tokenIssuer.IssueTokens(ctx, linkedAccount.UserID, accountDomain.AuthProviderPassword, command.IPAddress, command.UserAgent)
}

func (h *Handler) resolveCredential(ctx context.Context, email string) (*accountDomainLinkedAccounts.LinkedAccount, string, error) {
	linkedAccount, err := h.linkedAccountRepository.GetByAuthProviderAndProviderID(ctx, accountDomain.AuthProviderEmail, email)
	if err != nil {
		if errors.Is(err, accountDomainLinkedAccounts.ErrLinkedAccountNotFound) {
			return nil, h.getDummyHash(), nil
		}
		return nil, "", err
	}

	credential, err := h.passwordCredentialRepository.GetByUserID(ctx, linkedAccount.UserID)
	if err != nil {
		if errors.Is(err, accountDomainPasswordCredentials.ErrPasswordCredentialNotFound) {
			return nil, h.getDummyHash(), nil
		}
		return nil, "", err
	}

	return linkedAccount, credential.PasswordHash, nil
}

func (h *Handler) getDummyHash() string {
	h.dummyHashOnce.Do(func() {
		hash, err := h.hasher.Generate(dummyPassword)
		if err != nil {
			slog.Error("Failed to seed dummy bcrypt hash; constant-time path degraded",
				"component", "login.Handler",
				"error", err)
			return
		}
		h.dummyHash = hash
	})
	return h.dummyHash
}
