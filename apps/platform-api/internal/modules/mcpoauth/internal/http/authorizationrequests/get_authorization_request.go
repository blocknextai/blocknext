package authorizationrequests

import (
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"

	resultPkg "github.com/blocknextai/go-packages/result"
	commonHTTP "github.com/blocknextai/platform-api/internal/common/http"
	authorizationrequestsUseCases "github.com/blocknextai/platform-api/internal/modules/mcpoauth/internal/usecases/authorizationrequests"
)

type GetAuthorizationRequestRequest struct {
	AuthorizationRequestID uuid.UUID `uri:"authorizationRequestId"`
}

func NewGetAuthorizationRequestHandler(service *authorizationrequestsUseCases.Service) fiber.Handler {
	return func(c fiber.Ctx) error {
		request := new(GetAuthorizationRequestRequest)
		if err := c.Bind().All(request); err != nil {
			return commonHTTP.ErrInvalidRequest
		}

		result, err := service.GetAuthorizationRequest(c.RequestCtx(), &authorizationrequestsUseCases.GetAuthorizationRequestQuery{
			AuthorizationRequestID: request.AuthorizationRequestID,
		})
		if err != nil {
			return err
		}

		return c.Status(fiber.StatusOK).JSON(resultPkg.Ok(result))
	}
}
