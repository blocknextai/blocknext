package clients

import (
	"errors"

	"github.com/gofiber/fiber/v3"

	mcpOAuthDomainClients "github.com/blocknextai/platform-api/internal/modules/mcpoauth/internal/domain/clients"
	mcpOAuthDomainOAuth2 "github.com/blocknextai/platform-api/internal/modules/mcpoauth/internal/domain/oauth2"
)

const (
	cacheControlNoStore = "no-store"
)

func errorCodeOf(err error) (mcpOAuthDomainOAuth2.ErrorCode, bool) {
	switch {
	case errors.Is(err, mcpOAuthDomainClients.ErrInvalidRedirectURI):
		return mcpOAuthDomainOAuth2.InvalidRedirectURIError, true
	case errors.Is(err, mcpOAuthDomainClients.ErrInvalidClientName),
		errors.Is(err, mcpOAuthDomainClients.ErrInvalidGrantTypes),
		errors.Is(err, mcpOAuthDomainClients.ErrInvalidResponseTypes),
		errors.Is(err, mcpOAuthDomainClients.ErrInvalidScopes),
		errors.Is(err, mcpOAuthDomainClients.ErrInvalidTokenEndpointAuthMethod):
		return mcpOAuthDomainOAuth2.InvalidClientMetadataError, true
	default:
		return "", false
	}
}

func writeOAuthError(c fiber.Ctx, errorCode mcpOAuthDomainOAuth2.ErrorCode, err error) error {
	c.Set(fiber.HeaderCacheControl, cacheControlNoStore)

	return c.Status(fiber.StatusBadRequest).JSON(mcpOAuthDomainOAuth2.ErrorResponse{
		Error:            errorCode,
		ErrorDescription: err.Error(),
	})
}
