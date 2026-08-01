package verify

import (
	"context"

	"github.com/blocknextai/go-packages/database"
	accountApplicationAuth "github.com/blocknextai/platform-api/internal/account/application/auth"
	"github.com/blocknextai/platform-api/internal/account/application/auth/tokenissuer"
	"github.com/blocknextai/platform-api/internal/account/application/auth/verificationtoken"
	accountDomain "github.com/blocknextai/platform-api/internal/account/domain"
	accountDomainLinkedAccounts "github.com/blocknextai/platform-api/internal/account/domain/linkedaccounts"
	accountDomainVerificationTokens "github.com/blocknextai/platform-api/internal/account/domain/verificationtokens"
	"github.com/google/uuid"
)

type Handler struct {
	linkedAccountRepository     accountDomainLinkedAccounts.LinkedAccountRepository
	verificationTokenRepository accountDomainVerificationTokens.VerificationTokenRepository
	tokenIssuer                 tokenissuer.Service
	transactionManager          database.TransactionManager
}

func New(
	linkedAccountRepository accountDomainLinkedAccounts.LinkedAccountRepository,
	verificationTokenRepository accountDomainVerificationTokens.VerificationTokenRepository,
	tokenIssuer tokenissuer.Service,
	transactionManager database.TransactionManager,
) *Handler {
	return &Handler{
		linkedAccountRepository:     linkedAccountRepository,
		verificationTokenRepository: verificationTokenRepository,
		tokenIssuer:                 tokenIssuer,
		transactionManager:          transactionManager,
	}
}

func (h *Handler) Handle(ctx context.Context, command *VerifyCommand) (*accountApplicationAuth.AccessTokenResponse, error) {
	tokenHash := verificationtoken.Hash(command.Token)

	var userID uuid.UUID

	err := h.transactionManager.ExecuteInTransaction(ctx, func(txCtx context.Context) error {
		token, err := h.verificationTokenRepository.GetByTokenHash(txCtx, tokenHash)
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

		if err = h.verificationTokenRepository.Update(txCtx, consumed); err != nil {
			return err
		}

		linkedAccount, err := h.linkedAccountRepository.GetByAuthProviderAndProviderID(txCtx, accountDomain.AuthProviderEmail, token.Email)
		if err != nil {
			return err
		}

		if !linkedAccount.IsVerified {
			verified, err := linkedAccount.MarkVerified()
			if err != nil {
				return err
			}

			if err = h.linkedAccountRepository.Update(txCtx, verified); err != nil {
				return err
			}
		}

		userID = token.UserID
		return nil
	})

	if err != nil {
		return nil, err
	}

	return h.tokenIssuer.IssueTokens(ctx, userID, accountDomain.AuthProviderPassword, command.IPAddress, command.UserAgent)
}
