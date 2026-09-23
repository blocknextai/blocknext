package platformcredentials

import (
	"github.com/gofiber/fiber/v3"

	resultPkg "github.com/blocknextai/go-packages/result"
	platformcredentialsUseCases "github.com/blocknextai/platform-api/internal/modules/platform/internal/usecases/platformcredentials"
)

func NewGetAllPlatformCredentialsHandler(service *platformcredentialsUseCases.Service) fiber.Handler {
	return func(c fiber.Ctx) error {
		result, err := service.GetAllPlatformCredentials(c.RequestCtx(), &platformcredentialsUseCases.GetAllPlatformCredentialsQuery{})
		if err != nil {
			return err
		}

		return c.Status(fiber.StatusOK).JSON(resultPkg.Ok(result))
	}
}
