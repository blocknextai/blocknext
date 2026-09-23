package linkedaccounts

import (
	"github.com/blocknextai/go-packages/database"
	accountDomainLinkedAccounts "github.com/blocknextai/platform-api/internal/modules/account/internal/domain/linkedaccounts"
	accountDomainPasswordCredentials "github.com/blocknextai/platform-api/internal/modules/account/internal/domain/passwordcredentials"
	accountDomainUserNonces "github.com/blocknextai/platform-api/internal/modules/account/internal/domain/usernonces"
	accountDomainUsers "github.com/blocknextai/platform-api/internal/modules/account/internal/domain/users"
	authUseCases "github.com/blocknextai/platform-api/internal/modules/account/internal/usecases/auth"
)

type Service struct {
	userRepository               accountDomainUsers.UserRepository
	userNonceRepository          accountDomainUserNonces.UserNonceRepository
	linkedAccountRepository      accountDomainLinkedAccounts.LinkedAccountRepository
	authProviderRegistry         authUseCases.AuthProviderRegistry
	transactionManager           database.TransactionManager
	passwordCredentialRepository accountDomainPasswordCredentials.PasswordCredentialRepository
}

func NewService(
	userRepository accountDomainUsers.UserRepository,
	userNonceRepository accountDomainUserNonces.UserNonceRepository,
	linkedAccountRepository accountDomainLinkedAccounts.LinkedAccountRepository,
	authProviderRegistry authUseCases.AuthProviderRegistry,
	transactionManager database.TransactionManager,
	passwordCredentialRepository accountDomainPasswordCredentials.PasswordCredentialRepository,
) *Service {
	return &Service{
		userRepository:               userRepository,
		userNonceRepository:          userNonceRepository,
		linkedAccountRepository:      linkedAccountRepository,
		authProviderRegistry:         authProviderRegistry,
		transactionManager:           transactionManager,
		passwordCredentialRepository: passwordCredentialRepository,
	}
}
