package metadata

import (
	"github.com/gofiber/fiber/v3"

	metadataUseCases "github.com/blocknextai/platform-api/internal/modules/mcpoauth/internal/usecases/metadata"
)

func NewGetAuthorizationServerMetadataHandler(service *metadataUseCases.Service) fiber.Handler {
	return func(c fiber.Ctx) error {
		result, err := service.GetAuthorizationServerMetadata(c.RequestCtx(), &metadataUseCases.GetAuthorizationServerMetadataQuery{})
		if err != nil {
			return err
		}

		return c.Status(fiber.StatusOK).JSON(result)
	}
}
