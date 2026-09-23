package nodes

import (
	"github.com/gofiber/fiber/v3"

	resultPkg "github.com/blocknextai/go-packages/result"
	nodesUseCases "github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/usecases/nodes"
)

func NewGetAllNodesHandler(service *nodesUseCases.Service) fiber.Handler {
	return func(c fiber.Ctx) error {
		result, err := service.GetAllNodes(c.RequestCtx(), &nodesUseCases.GetAllNodesQuery{})
		if err != nil {
			return err
		}
		return c.Status(fiber.StatusOK).JSON(resultPkg.Ok(result))
	}
}
