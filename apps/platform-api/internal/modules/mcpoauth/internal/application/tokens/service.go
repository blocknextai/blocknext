package tokens

import (
	"context"

	mcpOAuthDomainOAuth2 "github.com/blocknextai/platform-api/internal/modules/mcpoauth/internal/domain/oauth2"
)

type AccessTokenSigner interface {
	Sign(accessToken *mcpOAuthDomainOAuth2.AccessToken) (string, error)
}

type AccessTokenValidator interface {
	Validate(ctx context.Context, rawToken string, allowedResources []string) (*mcpOAuthDomainOAuth2.AccessToken, error)
}
