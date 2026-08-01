package credentials

import (
	resultPkg "github.com/blocknextai/go-packages/result"
	"github.com/blocknextai/platform-api/internal/common/application/cqrs"
	commonHTTP "github.com/blocknextai/platform-api/internal/common/presentation/http"
	"github.com/blocknextai/platform-api/internal/nodeengine/application/credentials/getcredentialbyid"
	"github.com/gofiber/fiber/v3"
)

type GetCredentialByIDRequest struct {
	ID string `uri:"id"`
}

func NewGetCredentialByIDHandler(handler cqrs.Handler[*getcredentialbyid.GetCredentialByIDQuery, *getcredentialbyid.GetCredentialByIDResponse]) fiber.Handler {
	return func(c fiber.Ctx) error {
		request := new(GetCredentialByIDRequest)
		if err := c.Bind().All(request); err != nil {
			return commonHTTP.ErrInvalidRequest
		}

		result, err := handler.Handle(c.RequestCtx(), &getcredentialbyid.GetCredentialByIDQuery{
			ID: request.ID,
		})

		if err != nil {
			return err
		}

		return c.Status(fiber.StatusOK).JSON(resultPkg.Ok(result))
	}
}
