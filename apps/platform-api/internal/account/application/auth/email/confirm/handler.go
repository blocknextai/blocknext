package confirm

import (
	"context"
	"errors"

	"github.com/blocknextai/go-packages/database"
	"github.com/blocknextai/platform-api/internal/account/application/auth/verificationtoken"
	accountDomain "github.com/blocknextai/platform-api/internal/account/domain"
	accountDomainLinkedAccounts "github.com/blocknextai/platform-api/internal/account/domain/linkedaccounts"
	accountDomainUsers "github.com/blocknextai/platform-api/internal/account/domain/users"
	accountDomainVerificationTokens "github.com/blocknextai/platform-api/internal/account/domain/verificationtokens"
)

type Handler struct {
	linkedAccountRepository     accountDomainLinkedAccounts.LinkedAccountRepository
	verificationTokenRepository accountDomainVerificationTokens.VerificationTokenRepository
	transactionManager          database.TransactionManager
}

func New(
	linkedAccountRepository accountDomainLinkedAccounts.LinkedAccountRepository,
	verificationTokenRepository accountDomainVerificationTokens.VerificationTokenRepository,
	transactionManager database.TransactionManager,
) *Handler {
	return &Handler{
		linkedAccountRepository:     linkedAccountRepository,
		verificationTokenRepository: verificationTokenRepository,
		transactionManager:          transactionManager,
	}
}

func (h *Handler) Handle(ctx context.Context, command *ConfirmEmailChangeCommand) (*ConfirmEmailChangeResponse, error) {
	tokenHash := verificationtoken.Hash(command.Token)

	err := h.transactionManager.ExecuteInTransaction(ctx, func(txCtx context.Context) error {
		token, err := h.verificationTokenRepository.GetByTokenHash(txCtx, tokenHash)
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

		if err = h.verificationTokenRepository.Update(txCtx, consumed); err != nil {
			return err
		}

		taken, err := h.linkedAccountRepository.GetByAuthProviderAndProviderID(txCtx, accountDomain.AuthProviderEmail, token.Email)
		if err != nil && !errors.Is(err, accountDomainLinkedAccounts.ErrLinkedAccountNotFound) {
			return err
		}
		if taken != nil && taken.UserID != token.UserID {
			return accountDomainUsers.ErrEmailAlreadyInUse
		}

		linkedAccount, err := h.linkedAccountRepository.GetByAuthProviderAndUserID(txCtx, accountDomain.AuthProviderEmail, token.UserID)
		if err != nil {
			return err
		}

		updated, err := linkedAccount.ChangeEmailAndVerify(token.Email)
		if err != nil {
			return err
		}

		return h.linkedAccountRepository.Update(txCtx, updated)
	})

	if err != nil {
		return nil, err
	}

	return &ConfirmEmailChangeResponse{}, nil
}
