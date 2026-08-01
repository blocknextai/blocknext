package authproviders

import (
	"github.com/blocknextai/go-packages/apperror"
	"github.com/blocknextai/platform-api/internal/account/application/auth/createusertoken"
	accountDomain "github.com/blocknextai/platform-api/internal/account/domain"
)

var (
	ErrProviderNotFound = apperror.NotFound("auth provider not found")
)

type AuthProviderRegistry struct {
	providers map[accountDomain.AuthProvider]createusertoken.AuthProvider
}

func NewAuthProviderRegistry() *AuthProviderRegistry {
	return &AuthProviderRegistry{
		providers: make(map[accountDomain.AuthProvider]createusertoken.AuthProvider),
	}
}

func (r *AuthProviderRegistry) Register(authProvider accountDomain.AuthProvider, provider createusertoken.AuthProvider) {
	r.providers[authProvider] = provider
}

func (r *AuthProviderRegistry) GetProvider(authProvider accountDomain.AuthProvider) (createusertoken.AuthProvider, error) {
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
