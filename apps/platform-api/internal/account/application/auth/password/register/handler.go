package register

import (
	"context"
	"errors"
	"time"

	"github.com/blocknextai/go-packages/database"
	"github.com/blocknextai/go-packages/hashing"
	accountApplicationAuth "github.com/blocknextai/platform-api/internal/account/application/auth"
	"github.com/blocknextai/platform-api/internal/account/application/auth/passwordpolicy"
	"github.com/blocknextai/platform-api/internal/account/application/auth/tokenissuer"
	"github.com/blocknextai/platform-api/internal/account/application/auth/verificationtokenissuer"
	accountDomain "github.com/blocknextai/platform-api/internal/account/domain"
	accountDomainEmails "github.com/blocknextai/platform-api/internal/account/domain/emails"
	accountDomainLinkedAccounts "github.com/blocknextai/platform-api/internal/account/domain/linkedaccounts"
	accountDomainPasswordCredentials "github.com/blocknextai/platform-api/internal/account/domain/passwordcredentials"
	accountDomainUserPreferences "github.com/blocknextai/platform-api/internal/account/domain/userpreferences"
	accountDomainUsers "github.com/blocknextai/platform-api/internal/account/domain/users"
	accountDomainVerificationTokens "github.com/blocknextai/platform-api/internal/account/domain/verificationtokens"
	"github.com/blocknextai/platform-api/internal/eventbus/application/publishing"
	"github.com/google/uuid"
)

type Handler struct {
	userRepository               accountDomainUsers.UserRepository
	linkedAccountRepository      accountDomainLinkedAccounts.LinkedAccountRepository
	passwordCredentialRepository accountDomainPasswordCredentials.PasswordCredentialRepository
	userPreferenceRepository     accountDomainUserPreferences.UserPreferenceRepository
	tokenIssuer                  verificationtokenissuer.Service
	sessionTokenIssuer           tokenissuer.Service
	passwordPolicy               passwordpolicy.Policy
	hasher                       hashing.Hasher
	verificationTokenTTL         time.Duration
	eventBusPublisherService     publishing.PublisherService
	transactionManager           database.TransactionManager
}

func New(
	userRepository accountDomainUsers.UserRepository,
	linkedAccountRepository accountDomainLinkedAccounts.LinkedAccountRepository,
	passwordCredentialRepository accountDomainPasswordCredentials.PasswordCredentialRepository,
	userPreferenceRepository accountDomainUserPreferences.UserPreferenceRepository,
	tokenIssuer verificationtokenissuer.Service,
	sessionTokenIssuer tokenissuer.Service,
	passwordPolicy passwordpolicy.Policy,
	hasher hashing.Hasher,
	verificationTokenTTL time.Duration,
	eventBusPublisherService publishing.PublisherService,
	transactionManager database.TransactionManager,
) *Handler {
	return &Handler{
		userRepository:               userRepository,
		linkedAccountRepository:      linkedAccountRepository,
		passwordCredentialRepository: passwordCredentialRepository,
		userPreferenceRepository:     userPreferenceRepository,
		tokenIssuer:                  tokenIssuer,
		sessionTokenIssuer:           sessionTokenIssuer,
		passwordPolicy:               passwordPolicy,
		hasher:                       hasher,
		verificationTokenTTL:         verificationTokenTTL,
		eventBusPublisherService:     eventBusPublisherService,
		transactionManager:           transactionManager,
	}
}

func (h *Handler) Handle(ctx context.Context, command *RegisterCommand) (*accountApplicationAuth.AccessTokenResponse, error) {
	normalizedEmail := accountDomain.NormalizeEmail(command.Email)

	if err := h.passwordPolicy.Check(ctx, command.Password, []string{normalizedEmail}); err != nil {
		return nil, err
	}

	existing, err := h.linkedAccountRepository.GetByAuthProviderAndProviderID(ctx, accountDomain.AuthProviderEmail, normalizedEmail)
	if err != nil && !errors.Is(err, accountDomainLinkedAccounts.ErrLinkedAccountNotFound) {
		return nil, err
	}

	if existing != nil {
		if err := h.eventBusPublisherService.Enqueue(ctx, accountDomainEmails.RegistrationExistingEmailNotifiedDomainEvent{
			Email: normalizedEmail,
		}); err != nil {
			return nil, err
		}
		return &accountApplicationAuth.AccessTokenResponse{}, nil
	}

	passwordHash, err := h.hasher.Generate(command.Password)
	if err != nil {
		return nil, accountApplicationAuth.ErrFailedToHashPassword
	}

	var userID uuid.UUID

	err = h.transactionManager.ExecuteInTransaction(ctx, func(txCtx context.Context) error {
		user, err := accountDomainUsers.NewUser()
		if err != nil {
			return err
		}

		userID = user.ID

		if err = h.userRepository.Create(txCtx, user); err != nil {
			return err
		}

		preferences, err := accountDomainUserPreferences.NewDefault(user.ID)
		if err != nil {
			return err
		}

		if err = h.userPreferenceRepository.Upsert(txCtx, preferences); err != nil {
			return err
		}

		emailLinkedAccount, err := accountDomainLinkedAccounts.NewLinkedAccount(
			user.ID,
			accountDomain.AuthProviderEmail,
			normalizedEmail,
			normalizedEmail,
			new(normalizedEmail),
			true,
			false,
		)
		if err != nil {
			return err
		}

		if err = h.linkedAccountRepository.Create(txCtx, emailLinkedAccount); err != nil {
			return err
		}

		passwordLinkedAccount, err := accountDomainLinkedAccounts.NewLinkedAccount(
			user.ID,
			accountDomain.AuthProviderPassword,
			user.ID.String(),
			accountDomain.AuthProviderPassword.String(),
			new(accountDomain.AuthProviderPassword.String()),
			false,
			true,
		)
		if err != nil {
			return err
		}

		if err = h.linkedAccountRepository.Create(txCtx, passwordLinkedAccount); err != nil {
			return err
		}

		credential, err := accountDomainPasswordCredentials.NewPasswordCredential(user.ID, passwordHash)
		if err != nil {
			return err
		}

		if err = h.passwordCredentialRepository.Create(txCtx, credential); err != nil {
			return err
		}

		plainToken, err := h.tokenIssuer.Issue(txCtx, user.ID, accountDomainVerificationTokens.PurposeEmailVerify, normalizedEmail, h.verificationTokenTTL)
		if err != nil {
			return err
		}

		if err := h.eventBusPublisherService.Enqueue(txCtx, accountDomainUsers.UserCreatedDomainEvent{
			UserID:      user.ID,
			Identifier:  normalizedEmail,
			DisplayName: normalizedEmail,
		}); err != nil {
			return err
		}

		return h.eventBusPublisherService.Enqueue(txCtx, accountDomainEmails.EmailVerificationRequestedDomainEvent{
			Email: normalizedEmail,
			Token: plainToken,
		})
	})

	if err != nil {
		return nil, err
	}

	return h.sessionTokenIssuer.IssueTokens(ctx, userID, accountDomain.AuthProviderPassword, command.IPAddress, command.UserAgent)
}
