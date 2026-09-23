package authorizationrequests

import (
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"

	resultPkg "github.com/blocknextai/go-packages/result"
	commonHTTP "github.com/blocknextai/platform-api/internal/common/http"
	authorizationrequestsUseCases "github.com/blocknextai/platform-api/internal/modules/mcpoauth/internal/usecases/authorizationrequests"
)

type ApproveAuthorizationRequestRequest struct {
	AuthorizationRequestID uuid.UUID `uri:"authorizationRequestId"`
	OrganizationID         uuid.UUID `json:"organizationId"`
	Scopes                 []string  `json:"scopes"`
}

func NewApproveAuthorizationRequestHandler(service *authorizationrequestsUseCases.Service) fiber.Handler {
	return func(c fiber.Ctx) error {
		request := new(ApproveAuthorizationRequestRequest)
		if err := c.Bind().All(request); err != nil {
			return commonHTTP.ErrInvalidRequest
		}

		result, err := service.ApproveAuthorizationRequest(c.RequestCtx(), &authorizationrequestsUseCases.ApproveAuthorizationRequestCommand{
			AuthorizationRequestID: request.AuthorizationRequestID,
			UserID:                 commonHTTP.GetUserID(c),
			OrganizationID:         request.OrganizationID,
			Scopes:                 request.Scopes,
		})
		if err != nil {
			return err
		}

		return c.Status(fiber.StatusOK).JSON(resultPkg.Ok(result, resultPkg.WithMessage("authorization request approved")))
	}
}
