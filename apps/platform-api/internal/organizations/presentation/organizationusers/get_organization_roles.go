package organizationusers

import (
	resultPkg "github.com/blocknextai/go-packages/result"
	"github.com/blocknextai/platform-api/internal/common/application/cqrs"
	"github.com/blocknextai/platform-api/internal/organizations/application/organizationusers/getroles"
	"github.com/gofiber/fiber/v3"
)

func NewGetOrganizationRolesHandler(handler cqrs.Handler[*getroles.GetRolesQuery, *getroles.GetRolesResponse]) fiber.Handler {
	return func(c fiber.Ctx) error {

		result, err := handler.Handle(c.RequestCtx(), &getroles.GetRolesQuery{})

		if err != nil {
			return err
		}

		return c.Status(fiber.StatusOK).JSON(resultPkg.Ok(result))
	}
}
