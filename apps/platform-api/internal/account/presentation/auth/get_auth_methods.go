package auth

import (
	resultPkg "github.com/blocknextai/go-packages/result"
	"github.com/blocknextai/platform-api/internal/account/application/auth/getauthmethods"
	"github.com/blocknextai/platform-api/internal/common/application/cqrs"
	"github.com/gofiber/fiber/v3"
)

func NewGetAuthMethodsHandler(handler cqrs.Handler[*getauthmethods.GetAuthMethodsQuery, *getauthmethods.GetAuthMethodsResponse]) fiber.Handler {
	return func(c fiber.Ctx) error {
		result, err := handler.Handle(c.RequestCtx(), &getauthmethods.GetAuthMethodsQuery{})

		if err != nil {
			return err
		}

		return c.Status(fiber.StatusOK).JSON(resultPkg.Ok(result))
	}
}
