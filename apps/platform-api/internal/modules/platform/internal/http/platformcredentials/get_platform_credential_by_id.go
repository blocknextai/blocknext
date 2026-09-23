package platformcredentials

import (
	"github.com/gofiber/fiber/v3"

	resultPkg "github.com/blocknextai/go-packages/result"
	platformcredentialsUseCases "github.com/blocknextai/platform-api/internal/modules/platform/internal/usecases/platformcredentials"
)

func NewGetPlatformCredentialByIDHandler(service *platformcredentialsUseCases.Service) fiber.Handler {
	return func(c fiber.Ctx) error {
		id := c.Params("id")

		result, err := service.GetPlatformCredentialByID(c.RequestCtx(), &platformcredentialsUseCases.GetPlatformCredentialByIDQuery{
			ID: id,
		})
		if err != nil {
			return err
		}

		return c.Status(fiber.StatusOK).JSON(resultPkg.Ok(result))
	}
}
