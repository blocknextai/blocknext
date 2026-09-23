package http

import (
	"strings"

	"github.com/gofiber/fiber/v3"

	commonHTTP "github.com/blocknextai/platform-api/internal/common/http"
	serversUseCases "github.com/blocknextai/platform-api/internal/modules/mcp/internal/usecases/servers"
)

type GetProtectedResourceMetadataRequest struct {
	ServerID string `uri:"serverId"`
}

func NewGetProtectedResourceMetadataHandler(service *serversUseCases.Service) fiber.Handler {
	return func(c fiber.Ctx) error {
		request := new(GetProtectedResourceMetadataRequest)
		if err := c.Bind().All(request); err != nil {
			return commonHTTP.ErrInvalidRequest
		}

		result, err := service.GetProtectedResourceMetadata(c.RequestCtx(), &serversUseCases.GetProtectedResourceMetadataQuery{
			ServerID:     request.ServerID,
			ResourcePath: strings.TrimPrefix(c.Path(), protectedResourceMetadataPath),
		})
		if err != nil {
			return err
		}

		return c.Status(fiber.StatusOK).JSON(result)
	}
}
