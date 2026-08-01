package getauthmethods

import (
	"github.com/blocknextai/platform-api/internal/account/domain"
)

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
