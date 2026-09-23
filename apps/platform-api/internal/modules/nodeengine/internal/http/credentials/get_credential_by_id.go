package credentials

import (
	"github.com/gofiber/fiber/v3"

	resultPkg "github.com/blocknextai/go-packages/result"
	commonHTTP "github.com/blocknextai/platform-api/internal/common/http"
	credentialsUseCases "github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/usecases/credentials"
)

type GetCredentialByIDRequest struct {
	ID string `uri:"id"`
}

func NewGetCredentialByIDHandler(service *credentialsUseCases.Service) fiber.Handler {
	return func(c fiber.Ctx) error {
		request := new(GetCredentialByIDRequest)
		if err := c.Bind().All(request); err != nil {
			return commonHTTP.ErrInvalidRequest
		}

		result, err := service.GetCredentialByID(c.RequestCtx(), &credentialsUseCases.GetCredentialByIDQuery{
			ID: request.ID,
		})

		if err != nil {
			return err
		}

		return c.Status(fiber.StatusOK).JSON(resultPkg.Ok(result))
	}
}
