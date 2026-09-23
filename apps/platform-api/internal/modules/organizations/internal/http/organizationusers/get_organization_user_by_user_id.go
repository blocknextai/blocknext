package organizationusers

import (
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"

	resultPkg "github.com/blocknextai/go-packages/result"
	commonHTTP "github.com/blocknextai/platform-api/internal/common/http"
	organizationusersUseCases "github.com/blocknextai/platform-api/internal/modules/organizations/internal/usecases/organizationusers"
)

type GetOrganizationUserByUserIDRequest struct {
	OrganizationID uuid.UUID `uri:"organizationId"`
	UserID         uuid.UUID `uri:"userId"`
}

func NewGetOrganizationUserByUserIDHandler(service *organizationusersUseCases.Service) fiber.Handler {
	return func(c fiber.Ctx) error {
		request := new(GetOrganizationUserByUserIDRequest)
		if err := c.Bind().All(request); err != nil {
			return commonHTTP.ErrInvalidRequest
		}

		result, err := service.GetOrganizationUserByUserID(c.RequestCtx(), &organizationusersUseCases.GetOrganizationUserByUserIDQuery{
			OrganizationID: request.OrganizationID,
			UserID:         request.UserID,
		})

		if err != nil {
			return err
		}

		return c.Status(fiber.StatusOK).JSON(resultPkg.Ok(result))
	}
}
