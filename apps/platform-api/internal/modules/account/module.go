package account

import (
	"database/sql"
	"time"

	commonAuth "github.com/blocknextai/platform-api/internal/common/auth"

	"github.com/gofiber/fiber/v3"

	"github.com/blocknextai/go-packages/auth/jwt"
	"github.com/blocknextai/go-packages/cache"
	"github.com/blocknextai/go-packages/database"
	pkgEmail "github.com/blocknextai/go-packages/email"
	cachemiddleware "github.com/blocknextai/go-packages/fiber/middleware/cache"
	"github.com/blocknextai/go-packages/hashing"
	"github.com/blocknextai/platform-api/internal/config"
	"github.com/blocknextai/platform-api/internal/eventbus"
	"github.com/blocknextai/platform-api/internal/eventbus/idempotency"
	"github.com/blocknextai/platform-api/internal/eventbus/publishing"
	accountApplicationAuth "github.com/blocknextai/platform-api/internal/modules/account/internal/application/auth"
	accountApplicationLinkedAccounts "github.com/blocknextai/platform-api/internal/modules/account/internal/application/linkedaccounts"
	accountApplicationSessions "github.com/blocknextai/platform-api/internal/modules/account/internal/application/sessions"
	accountApplicationUsers "github.com/blocknextai/platform-api/internal/modules/account/internal/application/users"
	accountHTTP "github.com/blocknextai/platform-api/internal/modules/account/internal/http"
	accountPostgresLinkedAccounts "github.com/blocknextai/platform-api/internal/modules/account/internal/postgres/linkedaccounts"
	accountPostgresPasswordCredentials "github.com/blocknextai/platform-api/internal/modules/account/internal/postgres/passwordcredentials"
	accountPostgresSessions "github.com/blocknextai/platform-api/internal/modules/account/internal/postgres/sessions"
	accountPostgresUserNonces "github.com/blocknextai/platform-api/internal/modules/account/internal/postgres/usernonces"
	accountPostgresUserPreferences "github.com/blocknextai/platform-api/internal/modules/account/internal/postgres/userpreferences"
	accountPostgresUsers "github.com/blocknextai/platform-api/internal/modules/account/internal/postgres/users"
	accountPostgresUserSocials "github.com/blocknextai/platform-api/internal/modules/account/internal/postgres/usersocials"
	accountPostgresVerificationTokens "github.com/blocknextai/platform-api/internal/modules/account/internal/postgres/verificationtokens"
	accountUseCases "github.com/blocknextai/platform-api/internal/modules/account/internal/usecases"
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

	JWTService     jwt.AuthJWTService
	EmailSender    pkgEmail.EmailSender
	PasswordHasher hashing.Hasher
}

type Module struct {
	SessionService        accountApplicationSessions.SessionService
	UserService           accountApplicationUsers.UserService
	LinkedAccountService  accountApplicationLinkedAccounts.LinkedAccountService
	UserPermissionChecker commonAuth.UserPermissionChecker

	useCases *accountUseCases.Services
}

func NewModule(deps Dependencies) (*Module, error) {
	userRepository := accountPostgresUsers.NewUserRepository(deps.DB)
	userNonceRepository := accountPostgresUserNonces.NewUserNonceRepository(deps.DB)
	linkedAccountRepository := accountPostgresLinkedAccounts.NewLinkedAccountRepository(deps.DB)
	userSocialRepository := accountPostgresUserSocials.NewUserSocialRepository(deps.DB)
	userPreferenceRepository := accountPostgresUserPreferences.NewUserPreferenceRepository(deps.DB)
	sessionRepository := accountPostgresSessions.NewSessionRepository(deps.DB)
	passwordCredentialRepository := accountPostgresPasswordCredentials.NewPasswordCredentialRepository(deps.DB)
	verificationTokenRepository := accountPostgresVerificationTokens.NewVerificationTokenRepository(deps.DB)

	sessionService := accountApplicationSessions.NewSessionService(sessionRepository, deps.CacheService, deps.AccessTokenExpirationTime)
	userService := accountApplicationUsers.NewUserService(userRepository)
	linkedAccountService := accountApplicationLinkedAccounts.NewLinkedAccountService(linkedAccountRepository)
	userPermissionChecker := accountApplicationAuth.NewUserPermissionChecker(userRepository, deps.CacheService)

	authProviderRegistry, err := NewAuthProviderRegistry(deps.AuthOptions, userNonceRepository)
	if err != nil {
		return nil, err
	}

	useCases := accountUseCases.NewServices(accountUseCases.ServiceDependencies{
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
		useCases:              useCases,
	}, nil
}

func (m *Module) Register(router fiber.Router, authMiddleware *commonAuth.AuthMiddleware, cacheMiddleware *cachemiddleware.Middleware) {
	accountHTTP.RegisterRoutes(router, authMiddleware, cacheMiddleware, m.useCases)
}
