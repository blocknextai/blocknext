package forgot

import (
	"context"
	"errors"
	"time"

	"github.com/blocknextai/go-packages/database"
	"github.com/blocknextai/platform-api/internal/account/application/auth/verificationtokenissuer"
	accountDomain "github.com/blocknextai/platform-api/internal/account/domain"
	accountDomainEmails "github.com/blocknextai/platform-api/internal/account/domain/emails"
	accountDomainLinkedAccounts "github.com/blocknextai/platform-api/internal/account/domain/linkedaccounts"
	accountDomainVerificationTokens "github.com/blocknextai/platform-api/internal/account/domain/verificationtokens"
	"github.com/blocknextai/platform-api/internal/eventbus/application/publishing"
)

type Handler struct {
	linkedAccountRepository  accountDomainLinkedAccounts.LinkedAccountRepository
	tokenIssuer              verificationtokenissuer.Service
	resetTokenTTL            time.Duration
	eventBusPublisherService publishing.PublisherService
	transactionManager       database.TransactionManager
}

func New(
	linkedAccountRepository accountDomainLinkedAccounts.LinkedAccountRepository,
	tokenIssuer verificationtokenissuer.Service,
	resetTokenTTL time.Duration,
	eventBusPublisherService publishing.PublisherService,
	transactionManager database.TransactionManager,
) *Handler {
	return &Handler{
		linkedAccountRepository:  linkedAccountRepository,
		tokenIssuer:              tokenIssuer,
		resetTokenTTL:            resetTokenTTL,
		eventBusPublisherService: eventBusPublisherService,
		transactionManager:       transactionManager,
	}
}

func (h *Handler) Handle(ctx context.Context, command *ForgotCommand) (*ForgotResponse, error) {
	normalizedEmail := accountDomain.NormalizeEmail(command.Email)

	linkedAccount, err := h.linkedAccountRepository.GetByAuthProviderAndProviderID(ctx, accountDomain.AuthProviderEmail, normalizedEmail)
	if err != nil {
		if errors.Is(err, accountDomainLinkedAccounts.ErrLinkedAccountNotFound) {
			return &ForgotResponse{}, nil
		}
		return nil, err
	}

	err = h.transactionManager.ExecuteInTransaction(ctx, func(txCtx context.Context) error {
		plainToken, err := h.tokenIssuer.Issue(txCtx, linkedAccount.UserID, accountDomainVerificationTokens.PurposePasswordReset, normalizedEmail, h.resetTokenTTL)
		if err != nil {
			return err
		}

		return h.eventBusPublisherService.Enqueue(txCtx, accountDomainEmails.PasswordResetRequestedDomainEvent{
			Email: normalizedEmail,
			Token: plainToken,
		})
	})
	if err != nil {
		return nil, err
	}

	return &ForgotResponse{}, nil
}
