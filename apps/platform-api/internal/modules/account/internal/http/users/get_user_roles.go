package users

import (
	"github.com/gofiber/fiber/v3"

	resultPkg "github.com/blocknextai/go-packages/result"
	usersUseCases "github.com/blocknextai/platform-api/internal/modules/account/internal/usecases/users"
)

func NewGetRolesHandler(service *usersUseCases.Service) fiber.Handler {
	return func(c fiber.Ctx) error {

		result, err := service.GetRoles(c.RequestCtx(), &usersUseCases.GetRolesQuery{})

		if err != nil {
			return err
		}

		return c.Status(fiber.StatusOK).JSON(resultPkg.Ok(result))
	}
}
