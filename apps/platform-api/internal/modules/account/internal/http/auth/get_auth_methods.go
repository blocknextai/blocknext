package auth

import (
	"github.com/gofiber/fiber/v3"

	resultPkg "github.com/blocknextai/go-packages/result"
	authUseCases "github.com/blocknextai/platform-api/internal/modules/account/internal/usecases/auth"
)

func NewGetAuthMethodsHandler(service *authUseCases.Service) fiber.Handler {
	return func(c fiber.Ctx) error {
		result, err := service.GetAuthMethods(c.RequestCtx(), &authUseCases.GetAuthMethodsQuery{})

		if err != nil {
			return err
		}

		return c.Status(fiber.StatusOK).JSON(resultPkg.Ok(result))
	}
}
