package add

import (
	"context"
	"errors"
	"time"

	"github.com/blocknextai/go-packages/database"
	"github.com/blocknextai/platform-api/internal/account/application/auth/verificationtokenissuer"
	accountDomain "github.com/blocknextai/platform-api/internal/account/domain"
	accountDomainEmails "github.com/blocknextai/platform-api/internal/account/domain/emails"
	accountDomainLinkedAccounts "github.com/blocknextai/platform-api/internal/account/domain/linkedaccounts"
	accountDomainUsers "github.com/blocknextai/platform-api/internal/account/domain/users"
	accountDomainVerificationTokens "github.com/blocknextai/platform-api/internal/account/domain/verificationtokens"
	"github.com/blocknextai/platform-api/internal/eventbus/application/publishing"
)

type Handler struct {
	linkedAccountRepository  accountDomainLinkedAccounts.LinkedAccountRepository
	tokenIssuer              verificationtokenissuer.Service
	verificationTokenTTL     time.Duration
	eventBusPublisherService publishing.PublisherService
	transactionManager       database.TransactionManager
}

func New(
	linkedAccountRepository accountDomainLinkedAccounts.LinkedAccountRepository,
	tokenIssuer verificationtokenissuer.Service,
	verificationTokenTTL time.Duration,
	eventBusPublisherService publishing.PublisherService,
	transactionManager database.TransactionManager,
) *Handler {
	return &Handler{
		linkedAccountRepository:  linkedAccountRepository,
		tokenIssuer:              tokenIssuer,
		verificationTokenTTL:     verificationTokenTTL,
		eventBusPublisherService: eventBusPublisherService,
		transactionManager:       transactionManager,
	}
}

func (h *Handler) Handle(ctx context.Context, command *AddEmailCommand) (*AddEmailResponse, error) {
	normalizedEmail := accountDomain.NormalizeEmail(command.Email)

	existingForUser, err := h.linkedAccountRepository.GetByAuthProviderAndUserID(ctx, accountDomain.AuthProviderEmail, command.UserID)
	if err != nil && !errors.Is(err, accountDomainLinkedAccounts.ErrLinkedAccountNotFound) {
		return nil, err
	}
	if existingForUser != nil {
		return nil, ErrEmailAlreadyLinked
	}

	emailTaken, err := h.linkedAccountRepository.GetByAuthProviderAndProviderID(ctx, accountDomain.AuthProviderEmail, normalizedEmail)
	if err != nil && !errors.Is(err, accountDomainLinkedAccounts.ErrLinkedAccountNotFound) {
		return nil, err
	}
	if emailTaken != nil {
		return nil, accountDomainUsers.ErrEmailAlreadyInUse
	}

	err = h.transactionManager.ExecuteInTransaction(ctx, func(txCtx context.Context) error {
		linkedAccount, err := accountDomainLinkedAccounts.NewLinkedAccount(
			command.UserID,
			accountDomain.AuthProviderEmail,
			normalizedEmail,
			normalizedEmail,
			new(normalizedEmail),
			false,
			false,
		)
		if err != nil {
			return err
		}

		if err = h.linkedAccountRepository.Create(txCtx, linkedAccount); err != nil {
			return err
		}

		plainToken, err := h.tokenIssuer.Issue(txCtx, command.UserID, accountDomainVerificationTokens.PurposeEmailVerify, normalizedEmail, h.verificationTokenTTL)
		if err != nil {
			return err
		}

		return h.eventBusPublisherService.Enqueue(txCtx, accountDomainEmails.EmailAddedDomainEvent{
			Email: normalizedEmail,
			Token: plainToken,
		})
	})

	if err != nil {
		return nil, err
	}

	return &AddEmailResponse{}, nil
}
