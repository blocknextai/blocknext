package auth

import (
	"sync"
	"time"

	"github.com/blocknextai/go-packages/auth/jwt"
	"github.com/blocknextai/go-packages/database"
	"github.com/blocknextai/go-packages/hashing"
	"github.com/blocknextai/platform-api/internal/eventbus/publishing"
	"github.com/blocknextai/platform-api/internal/modules/account/internal/application/auth/passwordpolicy"
	"github.com/blocknextai/platform-api/internal/modules/account/internal/application/auth/tokenissuer"
	"github.com/blocknextai/platform-api/internal/modules/account/internal/application/auth/verificationtokenissuer"
	accountApplicationSessions "github.com/blocknextai/platform-api/internal/modules/account/internal/application/sessions"
	accountDomainLinkedAccounts "github.com/blocknextai/platform-api/internal/modules/account/internal/domain/linkedaccounts"
	accountDomainPasswordCredentials "github.com/blocknextai/platform-api/internal/modules/account/internal/domain/passwordcredentials"
	accountDomainUserNonces "github.com/blocknextai/platform-api/internal/modules/account/internal/domain/usernonces"
	accountDomainUserPreferences "github.com/blocknextai/platform-api/internal/modules/account/internal/domain/userpreferences"
	accountDomainUsers "github.com/blocknextai/platform-api/internal/modules/account/internal/domain/users"
	accountDomainVerificationTokens "github.com/blocknextai/platform-api/internal/modules/account/internal/domain/verificationtokens"
)

type Service struct {
	userNonceRepository          accountDomainUserNonces.UserNonceRepository
	authProviderRegistry         AuthProviderRegistry
	transactionManager           database.TransactionManager
	userRepository               accountDomainUsers.UserRepository
	linkedAccountRepository      accountDomainLinkedAccounts.LinkedAccountRepository
	userPreferenceRepository     accountDomainUserPreferences.UserPreferenceRepository
	tokenIssuer                  tokenissuer.Service
	eventBusPublisherService     publishing.PublisherService
	verificationTokenIssuer      verificationtokenissuer.Service
	verificationTokenTTL         time.Duration
	passwordCredentialRepository accountDomainPasswordCredentials.PasswordCredentialRepository
	hasher                       hashing.Hasher
	tokenTTL                     time.Duration
	verificationTokenRepository  accountDomainVerificationTokens.VerificationTokenRepository
	passwordEnabled              bool
	magicLinkEnabled             bool
	passwordPolicy               passwordpolicy.Policy
	resetTokenTTL                time.Duration
	dummyHashOnce                sync.Once
	dummyHash                    string
	sessionTokenIssuer           tokenissuer.Service
	sessionService               accountApplicationSessions.SessionService
	authJWTService               jwt.AuthJWTService
}

func NewService(
	userNonceRepository accountDomainUserNonces.UserNonceRepository,
	authProviderRegistry AuthProviderRegistry,
	transactionManager database.TransactionManager,
	userRepository accountDomainUsers.UserRepository,
	linkedAccountRepository accountDomainLinkedAccounts.LinkedAccountRepository,
	userPreferenceRepository accountDomainUserPreferences.UserPreferenceRepository,
	tokenIssuer tokenissuer.Service,
	eventBusPublisherService publishing.PublisherService,
	verificationTokenIssuer verificationtokenissuer.Service,
	verificationTokenTTL time.Duration,
	passwordCredentialRepository accountDomainPasswordCredentials.PasswordCredentialRepository,
	hasher hashing.Hasher,
	tokenTTL time.Duration,
	verificationTokenRepository accountDomainVerificationTokens.VerificationTokenRepository,
	passwordEnabled bool,
	magicLinkEnabled bool,
	passwordPolicy passwordpolicy.Policy,
	resetTokenTTL time.Duration,
	sessionTokenIssuer tokenissuer.Service,
	sessionService accountApplicationSessions.SessionService,
	authJWTService jwt.AuthJWTService,
) *Service {
	return &Service{
		userNonceRepository:          userNonceRepository,
		authProviderRegistry:         authProviderRegistry,
		transactionManager:           transactionManager,
		userRepository:               userRepository,
		linkedAccountRepository:      linkedAccountRepository,
		userPreferenceRepository:     userPreferenceRepository,
		tokenIssuer:                  tokenIssuer,
		eventBusPublisherService:     eventBusPublisherService,
		verificationTokenIssuer:      verificationTokenIssuer,
		verificationTokenTTL:         verificationTokenTTL,
		passwordCredentialRepository: passwordCredentialRepository,
		hasher:                       hasher,
		tokenTTL:                     tokenTTL,
		verificationTokenRepository:  verificationTokenRepository,
		passwordEnabled:              passwordEnabled,
		magicLinkEnabled:             magicLinkEnabled,
		passwordPolicy:               passwordPolicy,
		resetTokenTTL:                resetTokenTTL,
		sessionTokenIssuer:           sessionTokenIssuer,
		sessionService:               sessionService,
		authJWTService:               authJWTService,
	}
}
