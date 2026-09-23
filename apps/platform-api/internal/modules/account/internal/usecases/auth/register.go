package auth

import (
	"context"
	"errors"

	"github.com/google/uuid"

	"github.com/blocknextai/platform-api/internal/common/validation"
	accountApplicationAuth "github.com/blocknextai/platform-api/internal/modules/account/internal/application/auth"
	accountDomain "github.com/blocknextai/platform-api/internal/modules/account/internal/domain"
	accountDomainEmails "github.com/blocknextai/platform-api/internal/modules/account/internal/domain/emails"
	accountDomainLinkedAccounts "github.com/blocknextai/platform-api/internal/modules/account/internal/domain/linkedaccounts"
	accountDomainPasswordCredentials "github.com/blocknextai/platform-api/internal/modules/account/internal/domain/passwordcredentials"
	accountDomainUserPreferences "github.com/blocknextai/platform-api/internal/modules/account/internal/domain/userpreferences"
	accountDomainUsers "github.com/blocknextai/platform-api/internal/modules/account/internal/domain/users"
	accountDomainVerificationTokens "github.com/blocknextai/platform-api/internal/modules/account/internal/domain/verificationtokens"
)

type RegisterCommand struct {
	Email     string
	Password  string
	IPAddress string
	UserAgent string
}

func (c *RegisterCommand) Validate() error {
	if err := validation.Email(c.Email); err != nil {
		return err
	}
	return validation.NewPassword(c.Password)
}

func (s *Service) Register(ctx context.Context, command *RegisterCommand) (*accountApplicationAuth.AccessTokenResponse, error) {
	if err := command.Validate(); err != nil {
		return nil, err
	}

	normalizedEmail := accountDomain.NormalizeEmail(command.Email)

	if err := s.passwordPolicy.Check(ctx, command.Password, []string{normalizedEmail}); err != nil {
		return nil, err
	}

	existing, err := s.linkedAccountRepository.GetByAuthProviderAndProviderID(ctx, accountDomain.AuthProviderEmail, normalizedEmail)
	if err != nil && !errors.Is(err, accountDomainLinkedAccounts.ErrLinkedAccountNotFound) {
		return nil, err
	}

	if existing != nil {
		if err := s.eventBusPublisherService.Enqueue(ctx, accountDomainEmails.RegistrationExistingEmailNotifiedDomainEvent{
			Email: normalizedEmail,
		}); err != nil {
			return nil, err
		}
		return &accountApplicationAuth.AccessTokenResponse{}, nil
	}

	passwordHash, err := s.hasher.Generate(command.Password)
	if err != nil {
		return nil, accountApplicationAuth.ErrFailedToHashPassword
	}

	var userID uuid.UUID

	err = s.transactionManager.ExecuteInTransaction(ctx, func(txCtx context.Context) error {
		user, err := accountDomainUsers.NewUser()
		if err != nil {
			return err
		}

		userID = user.ID

		if err = s.userRepository.Create(txCtx, user); err != nil {
			return err
		}

		preferences, err := accountDomainUserPreferences.NewDefault(user.ID)
		if err != nil {
			return err
		}

		if err = s.userPreferenceRepository.Upsert(txCtx, preferences); err != nil {
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

		if err = s.linkedAccountRepository.Create(txCtx, emailLinkedAccount); err != nil {
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

		if err = s.linkedAccountRepository.Create(txCtx, passwordLinkedAccount); err != nil {
			return err
		}

		credential, err := accountDomainPasswordCredentials.NewPasswordCredential(user.ID, passwordHash)
		if err != nil {
			return err
		}

		if err = s.passwordCredentialRepository.Create(txCtx, credential); err != nil {
			return err
		}

		plainToken, err := s.verificationTokenIssuer.Issue(txCtx, user.ID, accountDomainVerificationTokens.PurposeEmailVerify, normalizedEmail, s.verificationTokenTTL)
		if err != nil {
			return err
		}

		if err := s.eventBusPublisherService.Enqueue(txCtx, accountDomainUsers.UserCreatedDomainEvent{
			UserID:      user.ID,
			Identifier:  normalizedEmail,
			DisplayName: normalizedEmail,
		}); err != nil {
			return err
		}

		return s.eventBusPublisherService.Enqueue(txCtx, accountDomainEmails.EmailVerificationRequestedDomainEvent{
			Email: normalizedEmail,
			Token: plainToken,
		})
	})

	if err != nil {
		return nil, err
	}

	return s.sessionTokenIssuer.IssueTokens(ctx, userID, accountDomain.AuthProviderPassword, command.IPAddress, command.UserAgent)
}
