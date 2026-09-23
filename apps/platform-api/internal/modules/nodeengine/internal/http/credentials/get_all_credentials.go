package credentials

import (
	"github.com/gofiber/fiber/v3"

	resultPkg "github.com/blocknextai/go-packages/result"
	credentialsUseCases "github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/usecases/credentials"
)

func NewGetAllCredentialsHandler(service *credentialsUseCases.Service) fiber.Handler {
	return func(c fiber.Ctx) error {
		result, err := service.GetAllCredentials(c.RequestCtx(), &credentialsUseCases.GetAllCredentialsQuery{})
		if err != nil {
			return err
		}
		return c.Status(fiber.StatusOK).JSON(resultPkg.Ok(result))
	}
}
