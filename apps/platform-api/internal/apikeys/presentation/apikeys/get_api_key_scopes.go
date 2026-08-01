package apikeys

import (
	resultPkg "github.com/blocknextai/go-packages/result"
	"github.com/blocknextai/platform-api/internal/apikeys/application/apikeys/getscopes"
	"github.com/blocknextai/platform-api/internal/common/application/cqrs"
	"github.com/gofiber/fiber/v3"
)

func NewGetAPIKeyScopesHandler(handler cqrs.Handler[*getscopes.GetScopesQuery, *getscopes.GetScopesResponse]) fiber.Handler {
	return func(c fiber.Ctx) error {
		result, err := handler.Handle(c.RequestCtx(), &getscopes.GetScopesQuery{})
		if err != nil {
			return err
		}

		return c.Status(fiber.StatusOK).JSON(resultPkg.Ok(result))
	}
}
