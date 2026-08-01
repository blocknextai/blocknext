package organizationusers

import (
	resultPkg "github.com/blocknextai/go-packages/result"
	"github.com/blocknextai/platform-api/internal/common/application/cqrs"
	commonHTTP "github.com/blocknextai/platform-api/internal/common/presentation/http"
	"github.com/blocknextai/platform-api/internal/organizations/application/organizationusers/getorganizationme"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type GetOrganizationMeRequest struct {
	OrganizationID uuid.UUID `uri:"organizationId"`
}

func NewGetOrganizationMeHandler(handler cqrs.Handler[*getorganizationme.GetOrganizationMeQuery, *getorganizationme.GetOrganizationMeResponse]) fiber.Handler {
	return func(c fiber.Ctx) error {
		request := new(GetOrganizationMeRequest)
		if err := c.Bind().All(request); err != nil {
			return commonHTTP.ErrInvalidRequest
		}

		userID := commonHTTP.GetUserID(c)

		result, err := handler.Handle(c.RequestCtx(), &getorganizationme.GetOrganizationMeQuery{
			OrganizationID: request.OrganizationID,
			UserID:         userID,
		})

		if err != nil {
			return err
		}

		return c.Status(fiber.StatusOK).JSON(resultPkg.Ok(result))
	}
}
