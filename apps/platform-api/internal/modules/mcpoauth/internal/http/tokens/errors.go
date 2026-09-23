package tokens

import (
	"errors"

	"github.com/gofiber/fiber/v3"

	mcpOAuthDomainAuthorizationCodes "github.com/blocknextai/platform-api/internal/modules/mcpoauth/internal/domain/authorizationcodes"
	mcpOAuthDomainClients "github.com/blocknextai/platform-api/internal/modules/mcpoauth/internal/domain/clients"
	mcpOAuthDomainGrants "github.com/blocknextai/platform-api/internal/modules/mcpoauth/internal/domain/grants"
	mcpOAuthDomainOAuth2 "github.com/blocknextai/platform-api/internal/modules/mcpoauth/internal/domain/oauth2"
	mcpOAuthDomainRefreshTokens "github.com/blocknextai/platform-api/internal/modules/mcpoauth/internal/domain/refreshtokens"
)

const (
	cacheControlNoStore          = "no-store"
	basicAuthenticationChallenge = `Basic realm="mcp-oauth"`
)

func errorCodeOf(err error) (mcpOAuthDomainOAuth2.ErrorCode, bool) {
	switch {
	case errors.Is(err, mcpOAuthDomainClients.ErrClientNotFound),
		errors.Is(err, mcpOAuthDomainClients.ErrInvalidClientID),
		errors.Is(err, mcpOAuthDomainClients.ErrInvalidClientCredentials),
		errors.Is(err, mcpOAuthDomainClients.ErrInvalidClientMetadataDocument),
		errors.Is(err, mcpOAuthDomainClients.ErrClientMetadataDocumentHostNotAllowed):
		return mcpOAuthDomainOAuth2.InvalidClientError, true
	case errors.Is(err, mcpOAuthDomainClients.ErrUnauthorizedGrantType):
		return mcpOAuthDomainOAuth2.UnauthorizedClientError, true
	case errors.Is(err, mcpOAuthDomainOAuth2.ErrUnsupportedGrantType):
		return mcpOAuthDomainOAuth2.UnsupportedGrantTypeError, true
	case errors.Is(err, mcpOAuthDomainAuthorizationCodes.ErrMissingCode),
		errors.Is(err, mcpOAuthDomainAuthorizationCodes.ErrMissingCodeVerifier),
		errors.Is(err, mcpOAuthDomainRefreshTokens.ErrMissingRefreshToken):
		return mcpOAuthDomainOAuth2.InvalidRequestError, true
	case errors.Is(err, mcpOAuthDomainGrants.ErrScopeExceedsGrant):
		return mcpOAuthDomainOAuth2.InvalidScopeError, true
	case errors.Is(err, mcpOAuthDomainGrants.ErrResourceMismatch):
		return mcpOAuthDomainOAuth2.InvalidTargetError, true
	case errors.Is(err, mcpOAuthDomainAuthorizationCodes.ErrAuthorizationCodeNotFound),
		errors.Is(err, mcpOAuthDomainAuthorizationCodes.ErrAuthorizationCodeUsed),
		errors.Is(err, mcpOAuthDomainAuthorizationCodes.ErrAuthorizationCodeExpired),
		errors.Is(err, mcpOAuthDomainAuthorizationCodes.ErrAuthorizationCodeClientMismatch),
		errors.Is(err, mcpOAuthDomainAuthorizationCodes.ErrAuthorizationCodeRedirectMismatch),
		errors.Is(err, mcpOAuthDomainAuthorizationCodes.ErrInvalidCodeVerifier),
		errors.Is(err, mcpOAuthDomainRefreshTokens.ErrRefreshTokenNotFound),
		errors.Is(err, mcpOAuthDomainRefreshTokens.ErrRefreshTokenUsed),
		errors.Is(err, mcpOAuthDomainRefreshTokens.ErrRefreshTokenRevoked),
		errors.Is(err, mcpOAuthDomainRefreshTokens.ErrRefreshTokenExpired),
		errors.Is(err, mcpOAuthDomainRefreshTokens.ErrRefreshTokenClientMismatch),
		errors.Is(err, mcpOAuthDomainRefreshTokens.ErrRefreshTokenResourceMismatch),
		errors.Is(err, mcpOAuthDomainGrants.ErrGrantNotFound),
		errors.Is(err, mcpOAuthDomainGrants.ErrGrantRevoked):
		return mcpOAuthDomainOAuth2.InvalidGrantError, true
	default:
		return "", false
	}
}

func writeOAuthError(c fiber.Ctx, errorCode mcpOAuthDomainOAuth2.ErrorCode, err error) error {
	c.Set(fiber.HeaderCacheControl, cacheControlNoStore)

	status := fiber.StatusBadRequest
	if errorCode == mcpOAuthDomainOAuth2.InvalidClientError {
		status = fiber.StatusUnauthorized
		c.Set(fiber.HeaderWWWAuthenticate, basicAuthenticationChallenge)
	}

	return c.Status(status).JSON(mcpOAuthDomainOAuth2.ErrorResponse{
		Error:            errorCode,
		ErrorDescription: err.Error(),
	})
}
