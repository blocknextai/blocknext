package apikeys

import (
	"github.com/gofiber/fiber/v3"

	resultPkg "github.com/blocknextai/go-packages/result"
	apikeysUseCases "github.com/blocknextai/platform-api/internal/modules/apikeys/internal/usecases/apikeys"
)

func NewGetAPIKeyScopesHandler(service *apikeysUseCases.Service) fiber.Handler {
	return func(c fiber.Ctx) error {
		result, err := service.GetScopes(c.RequestCtx(), &apikeysUseCases.GetScopesQuery{})
		if err != nil {
			return err
		}

		return c.Status(fiber.StatusOK).JSON(resultPkg.Ok(result))
	}
}
