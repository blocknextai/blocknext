package platformcredentials

import (
	resultPkg "github.com/blocknextai/go-packages/result"
	"github.com/blocknextai/platform-api/internal/common/application/cqrs"
	"github.com/blocknextai/platform-api/internal/platform/application/platformcredentials/getallplatformcredentials"
	"github.com/gofiber/fiber/v3"
)

func NewGetAllPlatformCredentialsHandler(handler cqrs.Handler[*getallplatformcredentials.GetAllPlatformCredentialsQuery, *getallplatformcredentials.GetAllPlatformCredentialsResponse]) fiber.Handler {
	return func(c fiber.Ctx) error {
		result, err := handler.Handle(c.RequestCtx(), &getallplatformcredentials.GetAllPlatformCredentialsQuery{})
		if err != nil {
			return err
		}

		return c.Status(fiber.StatusOK).JSON(resultPkg.Ok(result))
	}
}
