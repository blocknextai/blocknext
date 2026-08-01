package change

import (
	"context"
	"errors"
	"time"

	"github.com/blocknextai/go-packages/database"
	"github.com/blocknextai/go-packages/hashing"
	accountApplicationAuth "github.com/blocknextai/platform-api/internal/account/application/auth"
	"github.com/blocknextai/platform-api/internal/account/application/auth/verificationtokenissuer"
	accountDomain "github.com/blocknextai/platform-api/internal/account/domain"
	accountDomainEmails "github.com/blocknextai/platform-api/internal/account/domain/emails"
	accountDomainLinkedAccounts "github.com/blocknextai/platform-api/internal/account/domain/linkedaccounts"
	accountDomainPasswordCredentials "github.com/blocknextai/platform-api/internal/account/domain/passwordcredentials"
	accountDomainUsers "github.com/blocknextai/platform-api/internal/account/domain/users"
	accountDomainVerificationTokens "github.com/blocknextai/platform-api/internal/account/domain/verificationtokens"
	"github.com/blocknextai/platform-api/internal/eventbus/application/publishing"
)

type Handler struct {
	linkedAccountRepository      accountDomainLinkedAccounts.LinkedAccountRepository
	passwordCredentialRepository accountDomainPasswordCredentials.PasswordCredentialRepository
	tokenIssuer                  verificationtokenissuer.Service
	hasher                       hashing.Hasher
	tokenTTL                     time.Duration
	eventBusPublisherService     publishing.PublisherService
	transactionManager           database.TransactionManager
}

func New(
	linkedAccountRepository accountDomainLinkedAccounts.LinkedAccountRepository,
	passwordCredentialRepository accountDomainPasswordCredentials.PasswordCredentialRepository,
	tokenIssuer verificationtokenissuer.Service,
	hasher hashing.Hasher,
	tokenTTL time.Duration,
	eventBusPublisherService publishing.PublisherService,
	transactionManager database.TransactionManager,
) *Handler {
	return &Handler{
		linkedAccountRepository:      linkedAccountRepository,
		passwordCredentialRepository: passwordCredentialRepository,
		tokenIssuer:                  tokenIssuer,
		hasher:                       hasher,
		tokenTTL:                     tokenTTL,
		eventBusPublisherService:     eventBusPublisherService,
		transactionManager:           transactionManager,
	}
}

func (h *Handler) Handle(ctx context.Context, command *ChangeEmailCommand) (*ChangeEmailResponse, error) {
	normalizedNewEmail := accountDomain.NormalizeEmail(command.NewEmail)

	currentLinkedAccount, err := h.linkedAccountRepository.GetByAuthProviderAndUserID(ctx, accountDomain.AuthProviderEmail, command.UserID)
	if err != nil {
		if errors.Is(err, accountDomainLinkedAccounts.ErrLinkedAccountNotFound) {
			return nil, accountApplicationAuth.ErrPasswordAuthDisabled
		}
		return nil, err
	}

	if currentLinkedAccount.ProviderID == normalizedNewEmail {
		return nil, ErrSameEmail
	}

	credential, err := h.passwordCredentialRepository.GetByUserID(ctx, command.UserID)
	if err != nil {
		if errors.Is(err, accountDomainPasswordCredentials.ErrPasswordCredentialNotFound) {
			return nil, accountApplicationAuth.ErrPasswordNotSet
		}
		return nil, err
	}

	if !h.hasher.Compare(command.CurrentPassword, credential.PasswordHash) {
		return nil, accountApplicationAuth.ErrInvalidEmailOrPassword
	}

	existing, err := h.linkedAccountRepository.GetByAuthProviderAndProviderID(ctx, accountDomain.AuthProviderEmail, normalizedNewEmail)
	if err != nil && !errors.Is(err, accountDomainLinkedAccounts.ErrLinkedAccountNotFound) {
		return nil, err
	}
	if existing != nil {
		return nil, accountDomainUsers.ErrEmailAlreadyInUse
	}

	err = h.transactionManager.ExecuteInTransaction(ctx, func(txCtx context.Context) error {
		plainToken, err := h.tokenIssuer.Issue(txCtx, command.UserID, accountDomainVerificationTokens.PurposeEmailChange, normalizedNewEmail, h.tokenTTL)
		if err != nil {
			return err
		}

		if err := h.eventBusPublisherService.Enqueue(txCtx, accountDomainUsers.EmailChangedDomainEvent{
			UserID:   command.UserID,
			OldEmail: currentLinkedAccount.ProviderID,
			NewEmail: normalizedNewEmail,
		}); err != nil {
			return err
		}

		return h.eventBusPublisherService.Enqueue(txCtx, accountDomainEmails.EmailChangeRequestedDomainEvent{
			Email: normalizedNewEmail,
			Token: plainToken,
		})
	})
	if err != nil {
		return nil, err
	}

	return &ChangeEmailResponse{}, nil
}
