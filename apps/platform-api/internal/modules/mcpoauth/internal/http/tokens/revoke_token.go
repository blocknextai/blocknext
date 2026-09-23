package tokens

import (
	"github.com/gofiber/fiber/v3"

	commonHTTP "github.com/blocknextai/platform-api/internal/common/http"
	mcpOAuthDomainOAuth2 "github.com/blocknextai/platform-api/internal/modules/mcpoauth/internal/domain/oauth2"
	tokensUseCases "github.com/blocknextai/platform-api/internal/modules/mcpoauth/internal/usecases/tokens"
)

type RevokeTokenRequest struct {
	Authorization string `header:"Authorization"`
	Token         string `form:"token"`
	TokenTypeHint string `form:"token_type_hint"`
	ClientID      string `form:"client_id"`
	ClientSecret  string `form:"client_secret"`
}

func NewRevokeTokenHandler(service *tokensUseCases.Service) fiber.Handler {
	return func(c fiber.Ctx) error {
		request := new(RevokeTokenRequest)
		if err := c.Bind().All(request); err != nil {
			return writeOAuthError(c, mcpOAuthDomainOAuth2.InvalidRequestError, commonHTTP.ErrInvalidRequest)
		}

		clientID, clientSecret := clientCredentials(request.Authorization, request.ClientID, request.ClientSecret)

		_, err := service.RevokeToken(c.RequestCtx(), &tokensUseCases.RevokeTokenCommand{
			Token:         request.Token,
			TokenTypeHint: request.TokenTypeHint,
			ClientID:      clientID,
			ClientSecret:  clientSecret,
		})
		if err != nil {
			if errorCode, ok := errorCodeOf(err); ok {
				return writeOAuthError(c, errorCode, err)
			}
			return err
		}

		c.Set(fiber.HeaderCacheControl, cacheControlNoStore)

		return c.SendStatus(fiber.StatusOK)
	}
}
