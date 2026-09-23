package features

import (
	"github.com/gofiber/fiber/v3"

	resultPkg "github.com/blocknextai/go-packages/result"
	featuresUseCases "github.com/blocknextai/platform-api/internal/modules/platform/internal/usecases/features"
)

func NewGetFeaturesHandler(service *featuresUseCases.Service) fiber.Handler {
	return func(c fiber.Ctx) error {
		result, err := service.GetFeatures(c.RequestCtx(), &featuresUseCases.GetFeaturesQuery{})
		if err != nil {
			return err
		}

		return c.Status(fiber.StatusOK).JSON(resultPkg.Ok(result))
	}
}
