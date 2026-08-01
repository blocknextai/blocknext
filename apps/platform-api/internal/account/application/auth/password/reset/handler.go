package reset

import (
	"context"
	"errors"

	"github.com/blocknextai/go-packages/database"
	"github.com/blocknextai/go-packages/hashing"
	accountApplicationAuth "github.com/blocknextai/platform-api/internal/account/application/auth"
	"github.com/blocknextai/platform-api/internal/account/application/auth/passwordpolicy"
	"github.com/blocknextai/platform-api/internal/account/application/auth/verificationtoken"
	accountDomain "github.com/blocknextai/platform-api/internal/account/domain"
	accountDomainLinkedAccounts "github.com/blocknextai/platform-api/internal/account/domain/linkedaccounts"
	accountDomainPasswordCredentials "github.com/blocknextai/platform-api/internal/account/domain/passwordcredentials"
	accountDomainVerificationTokens "github.com/blocknextai/platform-api/internal/account/domain/verificationtokens"
	"github.com/google/uuid"
)

type Handler struct {
	linkedAccountRepository      accountDomainLinkedAccounts.LinkedAccountRepository
	passwordCredentialRepository accountDomainPasswordCredentials.PasswordCredentialRepository
	verificationTokenRepository  accountDomainVerificationTokens.VerificationTokenRepository
	passwordPolicy               passwordpolicy.Policy
	hasher                       hashing.Hasher
	transactionManager           database.TransactionManager
}

func New(
	linkedAccountRepository accountDomainLinkedAccounts.LinkedAccountRepository,
	passwordCredentialRepository accountDomainPasswordCredentials.PasswordCredentialRepository,
	verificationTokenRepository accountDomainVerificationTokens.VerificationTokenRepository,
	passwordPolicy passwordpolicy.Policy,
	hasher hashing.Hasher,
	transactionManager database.TransactionManager,
) *Handler {
	return &Handler{
		linkedAccountRepository:      linkedAccountRepository,
		passwordCredentialRepository: passwordCredentialRepository,
		verificationTokenRepository:  verificationTokenRepository,
		passwordPolicy:               passwordPolicy,
		hasher:                       hasher,
		transactionManager:           transactionManager,
	}
}

func (h *Handler) Handle(ctx context.Context, command *ResetCommand) (*ResetResponse, error) {
	if err := h.passwordPolicy.Check(ctx, command.NewPassword, nil); err != nil {
		return nil, err
	}

	tokenHash := verificationtoken.Hash(command.Token)

	newHash, err := h.hasher.Generate(command.NewPassword)
	if err != nil {
		return nil, accountApplicationAuth.ErrFailedToHashPassword
	}

	err = h.transactionManager.ExecuteInTransaction(ctx, func(txCtx context.Context) error {
		token, err := h.verificationTokenRepository.GetByTokenHash(txCtx, tokenHash)
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

		if err = h.verificationTokenRepository.Update(txCtx, consumed); err != nil {
			return err
		}

		credential, err := h.passwordCredentialRepository.GetByUserID(txCtx, token.UserID)
		if err != nil {
			if errors.Is(err, accountDomainPasswordCredentials.ErrPasswordCredentialNotFound) {
				return h.setupPassword(txCtx, token.UserID, newHash)
			}
			return err
		}

		updated, err := credential.ChangePassword(newHash)
		if err != nil {
			return err
		}

		return h.passwordCredentialRepository.Update(txCtx, updated)
	})

	if err != nil {
		return nil, err
	}

	return &ResetResponse{}, nil
}

func (h *Handler) setupPassword(ctx context.Context, userID uuid.UUID, passwordHash string) error {
	existingLA, err := h.linkedAccountRepository.GetByAuthProviderAndUserID(ctx, accountDomain.AuthProviderPassword, userID)
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

		if err := h.linkedAccountRepository.Create(ctx, passwordLA); err != nil {
			return err
		}
	}

	credential, err := accountDomainPasswordCredentials.NewPasswordCredential(userID, passwordHash)
	if err != nil {
		return err
	}

	return h.passwordCredentialRepository.Create(ctx, credential)
}
