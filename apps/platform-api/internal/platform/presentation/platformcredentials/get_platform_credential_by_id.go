package platformcredentials

import (
	resultPkg "github.com/blocknextai/go-packages/result"
	"github.com/blocknextai/platform-api/internal/common/application/cqrs"
	"github.com/blocknextai/platform-api/internal/platform/application/platformcredentials/getplatformcredentialbyid"
	"github.com/gofiber/fiber/v3"
)

func NewGetPlatformCredentialByIDHandler(handler cqrs.Handler[*getplatformcredentialbyid.GetPlatformCredentialByIDQuery, *getplatformcredentialbyid.GetPlatformCredentialByIDResponse]) fiber.Handler {
	return func(c fiber.Ctx) error {
		id := c.Params("id")

		result, err := handler.Handle(c.RequestCtx(), &getplatformcredentialbyid.GetPlatformCredentialByIDQuery{
			ID: id,
		})
		if err != nil {
			return err
		}

		return c.Status(fiber.StatusOK).JSON(resultPkg.Ok(result))
	}
}
