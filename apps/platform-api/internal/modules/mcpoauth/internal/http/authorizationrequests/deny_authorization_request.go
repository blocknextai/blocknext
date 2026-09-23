package authorizationrequests

import (
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"

	resultPkg "github.com/blocknextai/go-packages/result"
	commonHTTP "github.com/blocknextai/platform-api/internal/common/http"
	authorizationrequestsUseCases "github.com/blocknextai/platform-api/internal/modules/mcpoauth/internal/usecases/authorizationrequests"
)

type DenyAuthorizationRequestRequest struct {
	AuthorizationRequestID uuid.UUID `uri:"authorizationRequestId"`
}

func NewDenyAuthorizationRequestHandler(service *authorizationrequestsUseCases.Service) fiber.Handler {
	return func(c fiber.Ctx) error {
		request := new(DenyAuthorizationRequestRequest)
		if err := c.Bind().All(request); err != nil {
			return commonHTTP.ErrInvalidRequest
		}

		result, err := service.DenyAuthorizationRequest(c.RequestCtx(), &authorizationrequestsUseCases.DenyAuthorizationRequestCommand{
			AuthorizationRequestID: request.AuthorizationRequestID,
			UserID:                 commonHTTP.GetUserID(c),
		})
		if err != nil {
			return err
		}

		return c.Status(fiber.StatusOK).JSON(resultPkg.Ok(result, resultPkg.WithMessage("authorization request denied")))
	}
}
