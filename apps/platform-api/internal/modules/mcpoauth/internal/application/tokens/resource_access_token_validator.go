package tokens

import (
	"context"

	commonAuth "github.com/blocknextai/platform-api/internal/common/auth"
	mcpOAuthApplicationMetadata "github.com/blocknextai/platform-api/internal/modules/mcpoauth/internal/application/metadata"
)

type resourceAccessTokenValidator struct {
	accessTokenValidator AccessTokenValidator
	metadataService      mcpOAuthApplicationMetadata.MetadataService
}

func NewResourceAccessTokenValidator(
	accessTokenValidator AccessTokenValidator,
	metadataService mcpOAuthApplicationMetadata.MetadataService,
) commonAuth.AccessTokenValidator {
	return &resourceAccessTokenValidator{
		accessTokenValidator: accessTokenValidator,
		metadataService:      metadataService,
	}
}

func (v *resourceAccessTokenValidator) Validate(
	ctx context.Context,
	rawToken string,
	resourcePath string,
) (*commonAuth.AuthenticatedAccessToken, error) {
	accessToken, err := v.accessTokenValidator.Validate(ctx, rawToken, []string{
		v.metadataService.GetResourceURL(resourcePath),
		v.metadataService.GetResourceURL(""),
	})
	if err != nil {
		return nil, err
	}

	return &commonAuth.AuthenticatedAccessToken{
		UserID:         accessToken.Subject,
		OrganizationID: accessToken.OrganizationID,
		Scopes:         accessToken.Scopes.Strings(),
	}, nil
}
