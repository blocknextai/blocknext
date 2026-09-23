package organizations

import (
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"

	resultPkg "github.com/blocknextai/go-packages/result"
	commonHTTP "github.com/blocknextai/platform-api/internal/common/http"
	organizationsUseCases "github.com/blocknextai/platform-api/internal/modules/organizations/internal/usecases/organizations"
)

type UpdateOrganizationRequest struct {
	OrganizationID uuid.UUID `uri:"organizationId"`
	Title          string    `json:"title"`
	Description    *string   `json:"description"`
}

func NewUpdateOrganizationHandler(service *organizationsUseCases.Service) fiber.Handler {
	return func(c fiber.Ctx) error {
		request := new(UpdateOrganizationRequest)
		if err := c.Bind().All(request); err != nil {
			return commonHTTP.ErrInvalidRequest
		}

		userID := commonHTTP.GetUserID(c)

		result, err := service.UpdateOrganization(c.RequestCtx(), &organizationsUseCases.UpdateOrganizationCommand{
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
