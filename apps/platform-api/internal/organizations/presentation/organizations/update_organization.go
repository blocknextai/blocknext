package organizations

import (
	resultPkg "github.com/blocknextai/go-packages/result"
	"github.com/blocknextai/platform-api/internal/common/application/cqrs"
	commonHTTP "github.com/blocknextai/platform-api/internal/common/presentation/http"
	"github.com/blocknextai/platform-api/internal/organizations/application/organizations/updateorganization"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type UpdateOrganizationRequest struct {
	OrganizationID uuid.UUID `uri:"organizationId"`
	Title          string    `json:"title"`
	Description    *string   `json:"description"`
}

func NewUpdateOrganizationHandler(handler cqrs.Handler[*updateorganization.UpdateOrganizationCommand, *updateorganization.UpdateOrganizationResponse]) fiber.Handler {
	return func(c fiber.Ctx) error {
		request := new(UpdateOrganizationRequest)
		if err := c.Bind().All(request); err != nil {
			return commonHTTP.ErrInvalidRequest
		}

		userID := commonHTTP.GetUserID(c)

		result, err := handler.Handle(c.RequestCtx(), &updateorganization.UpdateOrganizationCommand{
			UserID:         userID,
			OrganizationID: request.OrganizationID,
			Title:          request.Title,
			Description:    request.Description,
		})

		if err != nil {
			return err
		}

		return c.Status(fiber.StatusOK).JSON(resultPkg.Ok(result, resultPkg.WithMessage("organization updated")))
	}
}
