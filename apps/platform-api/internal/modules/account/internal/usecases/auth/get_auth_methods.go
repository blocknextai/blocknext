package auth

import (
	"context"

	"github.com/blocknextai/platform-api/internal/modules/account/internal/domain"
)

type GetAuthMethodsQuery struct{}

type GetAuthMethodsResponse struct {
	Crypto    []domain.AuthProvider `json:"crypto"`
	OAuth     []domain.AuthProvider `json:"oauth"`
	Password  bool                  `json:"password"`
	MagicLink bool                  `json:"magicLink"`
}

func (s *Service) GetAuthMethods(ctx context.Context, request *GetAuthMethodsQuery) (*GetAuthMethodsResponse, error) {
	providers := s.authProviderRegistry.GetAllProviderKeys()
	return MapAuthMethodsToResponse(providers, s.passwordEnabled, s.magicLinkEnabled), nil
}

func MapAuthMethodsToResponse(
	providers []domain.AuthProvider,
	passwordEnabled bool,
	magicLinkEnabled bool,
) *GetAuthMethodsResponse {
	crypto := make([]domain.AuthProvider, 0)
	oauth := make([]domain.AuthProvider, 0)

	for _, provider := range providers {
		category := provider.Category()

		if category == domain.AuthProviderCategoryCrypto {
			crypto = append(crypto, provider)
			continue
		}

		if category == domain.AuthProviderCategoryOAuth {
			oauth = append(oauth, provider)
		}
	}

	return &GetAuthMethodsResponse{
		Crypto:    crypto,
		OAuth:     oauth,
		Password:  passwordEnabled,
		MagicLink: magicLinkEnabled,
	}
}
