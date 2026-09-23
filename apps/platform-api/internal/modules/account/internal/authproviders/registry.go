package authproviders

import (
	"github.com/blocknextai/go-packages/apperror"
	accountDomain "github.com/blocknextai/platform-api/internal/modules/account/internal/domain"
	authUseCases "github.com/blocknextai/platform-api/internal/modules/account/internal/usecases/auth"
)

var (
	ErrProviderNotFound = apperror.NotFound("auth provider not found")
)

type AuthProviderRegistry struct {
	providers map[accountDomain.AuthProvider]authUseCases.AuthProvider
}

func NewAuthProviderRegistry() *AuthProviderRegistry {
	return &AuthProviderRegistry{
		providers: make(map[accountDomain.AuthProvider]authUseCases.AuthProvider),
	}
}

func (r *AuthProviderRegistry) Register(authProvider accountDomain.AuthProvider, provider authUseCases.AuthProvider) {
	r.providers[authProvider] = provider
}

func (r *AuthProviderRegistry) GetProvider(authProvider accountDomain.AuthProvider) (authUseCases.AuthProvider, error) {
	provider, exists := r.providers[authProvider]
	if !exists {
		return nil, ErrProviderNotFound
	}
	return provider, nil
}

func (r *AuthProviderRegistry) GetAllProviderKeys() []accountDomain.AuthProvider {
	keys := make([]accountDomain.AuthProvider, 0, len(r.providers))
	for key := range r.providers {
		keys = append(keys, key)
	}
	return keys
}
