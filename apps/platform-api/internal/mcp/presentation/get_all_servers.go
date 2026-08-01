package presentation

import (
	resultPkg "github.com/blocknextai/go-packages/result"
	"github.com/blocknextai/platform-api/internal/common/application/cqrs"
	"github.com/blocknextai/platform-api/internal/mcp/application/getallservers"
	"github.com/gofiber/fiber/v3"
)

func NewGetAllServersHandler(handler cqrs.Handler[*getallservers.GetAllServersQuery, *getallservers.GetAllServersResponse]) fiber.Handler {
	return func(c fiber.Ctx) error {
		result, err := handler.Handle(c.RequestCtx(), &getallservers.GetAllServersQuery{})
		if err != nil {
			return err
		}
		return c.Status(fiber.StatusOK).JSON(resultPkg.Ok(result))
	}
}
