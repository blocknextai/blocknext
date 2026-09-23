package tokens

import (
	"github.com/gofiber/fiber/v3"

	commonHTTP "github.com/blocknextai/platform-api/internal/common/http"
	mcpOAuthDomainOAuth2 "github.com/blocknextai/platform-api/internal/modules/mcpoauth/internal/domain/oauth2"
	tokensUseCases "github.com/blocknextai/platform-api/internal/modules/mcpoauth/internal/usecases/tokens"
)

type ExchangeTokenRequest struct {
	Authorization string `header:"Authorization"`
	GrantType     string `form:"grant_type"`
	ClientID      string `form:"client_id"`
	ClientSecret  string `form:"client_secret"`
	Code          string `form:"code"`
	RedirectURI   string `form:"redirect_uri"`
	CodeVerifier  string `form:"code_verifier"`
	RefreshToken  string `form:"refresh_token"`
	Scope         string `form:"scope"`
	Resource      string `form:"resource"`
}

func NewExchangeTokenHandler(service *tokensUseCases.Service) fiber.Handler {
	return func(c fiber.Ctx) error {
		request := new(ExchangeTokenRequest)
		if err := c.Bind().All(request); err != nil {
			return writeOAuthError(c, mcpOAuthDomainOAuth2.InvalidRequestError, commonHTTP.ErrInvalidRequest)
		}

		clientID, clientSecret := clientCredentials(request.Authorization, request.ClientID, request.ClientSecret)

		result, err := service.ExchangeToken(c.RequestCtx(), &tokensUseCases.ExchangeTokenCommand{
			GrantType:    request.GrantType,
			ClientID:     clientID,
			ClientSecret: clientSecret,
			Code:         request.Code,
			RedirectURI:  request.RedirectURI,
			CodeVerifier: request.CodeVerifier,
			RefreshToken: request.RefreshToken,
			Scope:        request.Scope,
			Resource:     request.Resource,
		})
		if err != nil {
			if errorCode, ok := errorCodeOf(err); ok {
				return writeOAuthError(c, errorCode, err)
			}
			return err
		}

		c.Set(fiber.HeaderCacheControl, cacheControlNoStore)

		return c.Status(fiber.StatusOK).JSON(result)
	}
}
