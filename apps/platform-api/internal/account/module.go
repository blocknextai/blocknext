package account

import (
	"database/sql"
	"time"

	"github.com/blocknextai/go-packages/auth/jwt"
	"github.com/blocknextai/go-packages/cache"
	"github.com/blocknextai/go-packages/database"
	pkgEmail "github.com/blocknextai/go-packages/email"
	cachemiddleware "github.com/blocknextai/go-packages/fiber/middleware/cache"
	"github.com/blocknextai/go-packages/hashing"
	accountApplicationAuth "github.com/blocknextai/platform-api/internal/account/application/auth"
	accountApplicationLinkedAccounts "github.com/blocknextai/platform-api/internal/account/application/linkedaccounts"
	accountApplicationSessions "github.com/blocknextai/platform-api/internal/account/application/sessions"
	accountApplicationUsers "github.com/blocknextai/platform-api/internal/account/application/users"
	accountInfrastructure "github.com/blocknextai/platform-api/internal/account/infrastructure"
	accountInfrastructureLinkedAccounts "github.com/blocknextai/platform-api/internal/account/infrastructure/linkedaccounts"
	accountInfrastructurePasswordCredentials "github.com/blocknextai/platform-api/internal/account/infrastructure/passwordcredentials"
	accountInfrastructureSessions "github.com/blocknextai/platform-api/internal/account/infrastructure/sessions"
	accountInfrastructureUserNonces "github.com/blocknextai/platform-api/internal/account/infrastructure/usernonces"
	accountInfrastructureUserPreferences "github.com/blocknextai/platform-api/internal/account/infrastructure/userpreferences"
	accountInfrastructureUsers "github.com/blocknextai/platform-api/internal/account/infrastructure/users"
	accountInfrastructureUserSocials "github.com/blocknextai/platform-api/internal/account/infrastructure/usersocials"
	accountInfrastructureVerificationTokens "github.com/blocknextai/platform-api/internal/account/infrastructure/verificationtokens"
	accountPresentation "github.com/blocknextai/platform-api/internal/account/presentation"
	commonApplicationAuth "github.com/blocknextai/platform-api/internal/common/application/auth"
	"github.com/blocknextai/platform-api/internal/common/presentation/auth"
	"github.com/blocknextai/platform-api/internal/config"
	"github.com/blocknextai/platform-api/internal/eventbus"
	"github.com/blocknextai/platform-api/internal/eventbus/application/idempotency"
	"github.com/blocknextai/platform-api/internal/eventbus/application/publishing"
	web3DomainCrypto "github.com/blocknextai/platform-api/internal/web3/domain/crypto"
	"github.com/gofiber/fiber/v3"
)

type Dependencies struct {
	DB                       *sql.DB
	TransactionManager       database.TransactionManager
	EventBus                 *eventbus.Bus
	EventBusPublisherService publishing.PublisherService
	EventBusInboxService     *idempotency.InboxService
	CacheService             cache.Service

	AccessTokenExpirationTime time.Duration
	AuthOptions               config.AuthOptions
	PlatformUIBaseURL         string

	JWTService        jwt.AuthJWTService
	EmailSender       pkgEmail.EmailSender
	PasswordHasher    hashing.Hasher
	SignatureVerifier web3DomainCrypto.SignatureVerifier
}

type Module struct {
	SessionService        accountApplicationSessions.SessionService
	UserService           accountApplicationUsers.UserService
	LinkedAccountService  accountApplicationLinkedAccounts.LinkedAccountService
	UserPermissionChecker commonApplicationAuth.UserPermissionChecker

	handlers *accountInfrastructure.Handlers
}

func NewModule(deps Dependencies) (*Module, error) {
	userRepository := accountInfrastructureUsers.NewUserRepository(deps.DB)
	userNonceRepository := accountInfrastructureUserNonces.NewUserNonceRepository(deps.DB)
	linkedAccountRepository := accountInfrastructureLinkedAccounts.NewLinkedAccountRepository(deps.DB)
	userSocialRepository := accountInfrastructureUserSocials.NewUserSocialRepository(deps.DB)
	userPreferenceRepository := accountInfrastructureUserPreferences.NewUserPreferenceRepository(deps.DB)
	sessionRepository := accountInfrastructureSessions.NewSessionRepository(deps.DB)
	passwordCredentialRepository := accountInfrastructurePasswordCredentials.NewPasswordCredentialRepository(deps.DB)
	verificationTokenRepository := accountInfrastructureVerificationTokens.NewVerificationTokenRepository(deps.DB)

	sessionService := accountApplicationSessions.NewSessionService(sessionRepository, deps.CacheService, deps.AccessTokenExpirationTime)
	userService := accountApplicationUsers.NewUserService(userRepository)
	linkedAccountService := accountApplicationLinkedAccounts.NewLinkedAccountService(linkedAccountRepository)
	userPermissionChecker := accountApplicationAuth.NewUserPermissionChecker(userRepository, deps.CacheService)

	authProviderRegistry, err := accountInfrastructure.NewAuthProviderRegistry(deps.AuthOptions, deps.SignatureVerifier, userNonceRepository)
	if err != nil {
		return nil, err
	}

	handlers := accountInfrastructure.RegisterInfrastructure(accountInfrastructure.RegisterInfrastructureDeps{
		TransactionManager:       deps.TransactionManager,
		EventBus:                 deps.EventBus,
		EventBusPublisherService: deps.EventBusPublisherService,
		EventBusInboxService:     deps.EventBusInboxService,

		AuthOptions:       deps.AuthOptions,
		PlatformUIBaseURL: deps.PlatformUIBaseURL,

		UserRepository:               userRepository,
		UserNonceRepository:          userNonceRepository,
		LinkedAccountRepository:      linkedAccountRepository,
		UserSocialRepository:         userSocialRepository,
		UserPreferenceRepository:     userPreferenceRepository,
		SessionRepository:            sessionRepository,
		PasswordCredentialRepository: passwordCredentialRepository,
		VerificationTokenRepository:  verificationTokenRepository,
		SessionService:               sessionService,
		AuthProviderRegistry:         authProviderRegistry,
		AuthJWTService:               deps.JWTService,
		PasswordHasher:               deps.PasswordHasher,
		EmailSender:                  deps.EmailSender,
	})

	return &Module{
		SessionService:        sessionService,
		UserService:           userService,
		LinkedAccountService:  linkedAccountService,
		UserPermissionChecker: userPermissionChecker,
		handlers:              handlers,
	}, nil
}

func (m *Module) Register(router fiber.Router, authMiddleware *auth.AuthMiddleware, cacheMiddleware *cachemiddleware.Middleware) {
	accountPresentation.RegisterPresentation(router, authMiddleware, cacheMiddleware, m.handlers)
}
