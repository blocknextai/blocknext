package authorizationrequests

import (
	"github.com/gofiber/fiber/v3"

	commonHTTP "github.com/blocknextai/platform-api/internal/common/http"
	authorizationrequestsUseCases "github.com/blocknextai/platform-api/internal/modules/mcpoauth/internal/usecases/authorizationrequests"
)

const (
	cacheControlNoStore = "no-store"
)

type CreateAuthorizationRequestRequest struct {
	ClientID            string `query:"client_id"`
	RedirectURI         string `query:"redirect_uri"`
	ResponseType        string `query:"response_type"`
	Scope               string `query:"scope"`
	State               string `query:"state"`
	CodeChallenge       string `query:"code_challenge"`
	CodeChallengeMethod string `query:"code_challenge_method"`
	Resource            string `query:"resource"`
}

func NewCreateAuthorizationRequestHandler(service *authorizationrequestsUseCases.Service) fiber.Handler {
	return func(c fiber.Ctx) error {
		request := new(CreateAuthorizationRequestRequest)
		if err := c.Bind().All(request); err != nil {
			return commonHTTP.ErrInvalidRequest
		}

		result, err := service.CreateAuthorizationRequest(c.RequestCtx(), &authorizationrequestsUseCases.CreateAuthorizationRequestCommand{
			ClientID:            request.ClientID,
			RedirectURI:         request.RedirectURI,
			ResponseType:        request.ResponseType,
			Scope:               request.Scope,
			State:               request.State,
			CodeChallenge:       request.CodeChallenge,
			CodeChallengeMethod: request.CodeChallengeMethod,
			Resource:            request.Resource,
		})
		if err != nil {
			return err
		}

		c.Set(fiber.HeaderCacheControl, cacheControlNoStore)

		return c.Redirect().Status(fiber.StatusFound).To(result.RedirectURI)
	}
}
