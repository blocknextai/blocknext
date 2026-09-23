package organizationusers

import (
	"github.com/gofiber/fiber/v3"

	resultPkg "github.com/blocknextai/go-packages/result"
	organizationusersUseCases "github.com/blocknextai/platform-api/internal/modules/organizations/internal/usecases/organizationusers"
)

func NewGetOrganizationRolesHandler(service *organizationusersUseCases.Service) fiber.Handler {
	return func(c fiber.Ctx) error {

		result, err := service.GetRoles(c.RequestCtx(), &organizationusersUseCases.GetRolesQuery{})

		if err != nil {
			return err
		}

		return c.Status(fiber.StatusOK).JSON(resultPkg.Ok(result))
	}
}
