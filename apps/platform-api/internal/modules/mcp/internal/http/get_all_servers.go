package http

import (
	"github.com/gofiber/fiber/v3"

	resultPkg "github.com/blocknextai/go-packages/result"
	serversUseCases "github.com/blocknextai/platform-api/internal/modules/mcp/internal/usecases/servers"
)

func NewGetAllServersHandler(service *serversUseCases.Service) fiber.Handler {
	return func(c fiber.Ctx) error {
		result, err := service.GetAllServers(c.RequestCtx(), &serversUseCases.GetAllServersQuery{})
		if err != nil {
			return err
		}
		return c.Status(fiber.StatusOK).JSON(resultPkg.Ok(result))
	}
}
