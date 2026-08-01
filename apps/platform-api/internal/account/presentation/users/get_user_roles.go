package users

import (
	resultPkg "github.com/blocknextai/go-packages/result"
	"github.com/blocknextai/platform-api/internal/account/application/users/getroles"
	"github.com/blocknextai/platform-api/internal/common/application/cqrs"
	"github.com/gofiber/fiber/v3"
)

func NewGetRolesHandler(handler cqrs.Handler[*getroles.GetRolesQuery, *getroles.GetRolesResponse]) fiber.Handler {
	return func(c fiber.Ctx) error {

		result, err := handler.Handle(c.RequestCtx(), &getroles.GetRolesQuery{})

		if err != nil {
			return err
		}

		return c.Status(fiber.StatusOK).JSON(resultPkg.Ok(result))
	}
}
