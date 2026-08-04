package features

import (
	resultPkg "github.com/blocknextai/go-packages/result"
	"github.com/blocknextai/platform-api/internal/common/application/cqrs"
	"github.com/blocknextai/platform-api/internal/platform/application/features/getfeatures"
	"github.com/gofiber/fiber/v3"
)

func NewGetFeaturesHandler(handler cqrs.Handler[*getfeatures.GetFeaturesQuery, *getfeatures.GetFeaturesResponse]) fiber.Handler {
	return func(c fiber.Ctx) error {
		result, err := handler.Handle(c.RequestCtx(), &getfeatures.GetFeaturesQuery{})
		if err != nil {
			return err
		}

		return c.Status(fiber.StatusOK).JSON(resultPkg.Ok(result))
	}
}
