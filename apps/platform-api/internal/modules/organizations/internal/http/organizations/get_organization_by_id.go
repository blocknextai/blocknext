package organizations

import (
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"

	resultPkg "github.com/blocknextai/go-packages/result"
	commonHTTP "github.com/blocknextai/platform-api/internal/common/http"
	organizationsUseCases "github.com/blocknextai/platform-api/internal/modules/organizations/internal/usecases/organizations"
)

type GetOrganizationByIDRequest struct {
	OrganizationID uuid.UUID `uri:"organizationId"`
}

func NewGetOrganizationByIDHandler(service *organizationsUseCases.Service) fiber.Handler {
	return func(c fiber.Ctx) error {
		request := new(GetOrganizationByIDRequest)
		if err := c.Bind().All(request); err != nil {
			return commonHTTP.ErrInvalidRequest
		}

		userID := commonHTTP.GetUserID(c)

		result, err := service.GetOrganizationByID(c.RequestCtx(), &organizationsUseCases.GetOrganizationByIDQuery{
			UserID:         userID,
			OrganizationID: request.OrganizationID,
		})

		if err != nil {
			return err
		}

		return c.Status(fiber.StatusOK).JSON(resultPkg.Ok(result))
	}
}
