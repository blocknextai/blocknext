package organizations

import (
	resultPkg "github.com/blocknextai/go-packages/result"
	"github.com/blocknextai/platform-api/internal/common/application/cqrs"
	commonHTTP "github.com/blocknextai/platform-api/internal/common/presentation/http"
	"github.com/blocknextai/platform-api/internal/organizations/application/organizations/deleteorganization"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type DeleteOrganizationRequest struct {
	OrganizationID uuid.UUID `uri:"organizationId"`
}

func NewDeleteOrganizationHandler(handler cqrs.Handler[*deleteorganization.DeleteOrganizationCommand, *deleteorganization.DeleteOrganizationResponse]) fiber.Handler {
	return func(c fiber.Ctx) error {
		request := new(DeleteOrganizationRequest)
		if err := c.Bind().All(request); err != nil {
			return commonHTTP.ErrInvalidRequest
		}

		userID := commonHTTP.GetUserID(c)

		result, err := handler.Handle(c.RequestCtx(), &deleteorganization.DeleteOrganizationCommand{
			UserID:         userID,
			OrganizationID: request.OrganizationID,
		})

		if err != nil {
			return err
		}

		return c.Status(fiber.StatusOK).JSON(resultPkg.Ok(result, resultPkg.WithMessage("organization deleted")))
	}
}
