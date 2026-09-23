package usecases

import (
	"github.com/blocknextai/go-packages/auth/jwt"
	"github.com/blocknextai/go-packages/database"
	pkgEmail "github.com/blocknextai/go-packages/email"
	"github.com/blocknextai/go-packages/hashing"
	"github.com/blocknextai/platform-api/internal/config"
	"github.com/blocknextai/platform-api/internal/eventbus"
	"github.com/blocknextai/platform-api/internal/eventbus/idempotency"
	"github.com/blocknextai/platform-api/internal/eventbus/publishing"
	"github.com/blocknextai/platform-api/internal/modules/account/internal/application/auth/mailer"
	"github.com/blocknextai/platform-api/internal/modules/account/internal/application/auth/passwordpolicy"
	"github.com/blocknextai/platform-api/internal/modules/account/internal/application/auth/tokenissuer"
	"github.com/blocknextai/platform-api/internal/modules/account/internal/application/auth/verificationtokenissuer"
	"github.com/blocknextai/platform-api/internal/modules/account/internal/application/events/emailadded"
	"github.com/blocknextai/platform-api/internal/modules/account/internal/application/events/emailchangerequested"
	"github.com/blocknextai/platform-api/internal/modules/account/internal/application/events/emailverificationrequested"
	"github.com/blocknextai/platform-api/internal/modules/account/internal/application/events/magiclinkcreated"
	"github.com/blocknextai/platform-api/internal/modules/account/internal/application/events/passwordresetrequested"
	"github.com/blocknextai/platform-api/internal/modules/account/internal/application/events/registrationexistingemailnotified"
	"github.com/blocknextai/platform-api/internal/modules/account/internal/application/events/userwelcome"
	accountApplicationSessions "github.com/blocknextai/platform-api/internal/modules/account/internal/application/sessions"
	"github.com/blocknextai/platform-api/internal/modules/account/internal/domain/linkedaccounts"
	"github.com/blocknextai/platform-api/internal/modules/account/internal/domain/passwordcredentials"
	"github.com/blocknextai/platform-api/internal/modules/account/internal/domain/sessions"
	"github.com/blocknextai/platform-api/internal/modules/account/internal/domain/usernonces"
	"github.com/blocknextai/platform-api/internal/modules/account/internal/domain/userpreferences"
	"github.com/blocknextai/platform-api/internal/modules/account/internal/domain/users"
	"github.com/blocknextai/platform-api/internal/modules/account/internal/domain/usersocials"
	"github.com/blocknextai/platform-api/internal/modules/account/internal/domain/verificationtokens"
	authUseCases "github.com/blocknextai/platform-api/internal/modules/account/internal/usecases/auth"
	linkedaccountsUseCases "github.com/blocknextai/platform-api/internal/modules/account/internal/usecases/linkedaccounts"
	sessionsUseCases "github.com/blocknextai/platform-api/internal/modules/account/internal/usecases/sessions"
	userpreferencesUseCases "github.com/blocknextai/platform-api/internal/modules/account/internal/usecases/userpreferences"
	usersUseCases "github.com/blocknextai/platform-api/internal/modules/account/internal/usecases/users"
	usersocialsUseCases "github.com/blocknextai/platform-api/internal/modules/account/internal/usecases/usersocials"
)

type Services struct {
	Auth            *authUseCases.Service
	LinkedAccounts  *linkedaccountsUseCases.Service
	Sessions        *sessionsUseCases.Service
	Userpreferences *userpreferencesUseCases.Service
	Users           *usersUseCases.Service
	UserSocials     *usersocialsUseCases.Service
}

type ServiceDependencies struct {
	TransactionManager       database.TransactionManager
	EventBus                 *eventbus.Bus
	EventBusPublisherService publishing.PublisherService
	EventBusInboxService     *idempotency.InboxService

	AuthOptions       config.AuthOptions
	PlatformUIBaseURL string

	UserRepository               users.UserRepository
	UserNonceRepository          usernonces.UserNonceRepository
	LinkedAccountRepository      linkedaccounts.LinkedAccountRepository
	UserSocialRepository         usersocials.UserSocialRepository
	UserPreferenceRepository     userpreferences.UserPreferenceRepository
	SessionRepository            sessions.SessionRepository
	PasswordCredentialRepository passwordcredentials.PasswordCredentialRepository
	VerificationTokenRepository  verificationtokens.VerificationTokenRepository
	SessionService               accountApplicationSessions.SessionService
	AuthProviderRegistry         authUseCases.AuthProviderRegistry
	AuthJWTService               jwt.AuthJWTService
	PasswordHasher               hashing.Hasher
	EmailSender                  pkgEmail.EmailSender
}

func NewServices(deps ServiceDependencies) *Services {
	tokenIssuer := tokenissuer.NewService(deps.AuthJWTService, deps.SessionService, deps.UserRepository)
	verificationIssuer := verificationtokenissuer.NewService(deps.VerificationTokenRepository)
	policyChecker := passwordpolicy.NewChecker()
	mailerInstance := mailer.NewMailer(deps.EmailSender, deps.PlatformUIBaseURL)

	emailverificationrequested.New(mailerInstance, deps.EventBus)
	emailadded.New(mailerInstance, deps.EventBus)
	passwordresetrequested.New(mailerInstance, deps.EventBus)
	emailchangerequested.New(mailerInstance, deps.EventBus)
	magiclinkcreated.New(mailerInstance, deps.EventBus)
	registrationexistingemailnotified.New(mailerInstance, deps.EventBus)
	userwelcome.New(mailerInstance, deps.EventBus)

	return &Services{
		Auth:            authUseCases.NewService(deps.UserNonceRepository, deps.AuthProviderRegistry, deps.TransactionManager, deps.UserRepository, deps.LinkedAccountRepository, deps.UserPreferenceRepository, tokenIssuer, deps.EventBusPublisherService, verificationIssuer, deps.AuthOptions.Password.VerificationTokenTTL, deps.PasswordCredentialRepository, deps.PasswordHasher, deps.AuthOptions.Password.VerificationTokenTTL, deps.VerificationTokenRepository, deps.AuthOptions.Password.Enabled, deps.AuthOptions.MagicLink.Enabled, policyChecker, deps.AuthOptions.Password.PasswordResetTokenTTL, tokenIssuer, deps.SessionService, deps.AuthJWTService),
		LinkedAccounts:  linkedaccountsUseCases.NewService(deps.UserRepository, deps.UserNonceRepository, deps.LinkedAccountRepository, deps.AuthProviderRegistry, deps.TransactionManager, deps.PasswordCredentialRepository),
		Sessions:        sessionsUseCases.NewService(deps.SessionRepository, deps.SessionService),
		Userpreferences: userpreferencesUseCases.NewService(deps.UserPreferenceRepository),
		Users:           usersUseCases.NewService(deps.UserRepository),
		UserSocials:     usersocialsUseCases.NewService(deps.UserSocialRepository, deps.TransactionManager),
	}
}
