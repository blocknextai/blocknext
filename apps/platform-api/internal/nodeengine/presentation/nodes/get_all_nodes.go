package nodes

import (
	resultPkg "github.com/blocknextai/go-packages/result"
	"github.com/blocknextai/platform-api/internal/common/application/cqrs"
	"github.com/blocknextai/platform-api/internal/nodeengine/application/nodes/getallnodes"
	"github.com/gofiber/fiber/v3"
)

func NewGetAllNodesHandler(handler cqrs.Handler[*getallnodes.GetAllNodesQuery, *getallnodes.GetAllNodesResponse]) fiber.Handler {
	return func(c fiber.Ctx) error {
		result, err := handler.Handle(c.RequestCtx(), &getallnodes.GetAllNodesQuery{})
		if err != nil {
			return err
		}
		return c.Status(fiber.StatusOK).JSON(resultPkg.Ok(result))
	}
}
