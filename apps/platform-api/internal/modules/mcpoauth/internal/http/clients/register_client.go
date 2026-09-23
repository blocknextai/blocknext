package clients

import (
	"github.com/gofiber/fiber/v3"

	commonHTTP "github.com/blocknextai/platform-api/internal/common/http"
	mcpOAuthDomainOAuth2 "github.com/blocknextai/platform-api/internal/modules/mcpoauth/internal/domain/oauth2"
	clientsUseCases "github.com/blocknextai/platform-api/internal/modules/mcpoauth/internal/usecases/clients"
)

type RegisterClientRequest struct {
	ClientName              string   `json:"client_name"`
	RedirectURIs            []string `json:"redirect_uris"`
	GrantTypes              []string `json:"grant_types"`
	ResponseTypes           []string `json:"response_types"`
	Scope                   string   `json:"scope"`
	TokenEndpointAuthMethod string   `json:"token_endpoint_auth_method"`
	LogoURI                 string   `json:"logo_uri"`
	ClientURI               string   `json:"client_uri"`
}

func NewRegisterClientHandler(service *clientsUseCases.Service) fiber.Handler {
	return func(c fiber.Ctx) error {
		request := new(RegisterClientRequest)
		if err := c.Bind().All(request); err != nil {
			return writeOAuthError(c, mcpOAuthDomainOAuth2.InvalidClientMetadataError, commonHTTP.ErrInvalidRequest)
		}

		result, err := service.RegisterClient(c.RequestCtx(), &clientsUseCases.RegisterClientCommand{
			ClientName:              request.ClientName,
			RedirectURIs:            request.RedirectURIs,
			GrantTypes:              request.GrantTypes,
			ResponseTypes:           request.ResponseTypes,
			Scope:                   request.Scope,
			TokenEndpointAuthMethod: request.TokenEndpointAuthMethod,
			LogoURI:                 request.LogoURI,
			ClientURI:               request.ClientURI,
		})
		if err != nil {
			if errorCode, ok := errorCodeOf(err); ok {
				return writeOAuthError(c, errorCode, err)
			}
			return err
		}

		c.Set(fiber.HeaderCacheControl, cacheControlNoStore)

		return c.Status(fiber.StatusCreated).JSON(result)
	}
}
