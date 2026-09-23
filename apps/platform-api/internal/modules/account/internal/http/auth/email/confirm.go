package email

import (
	"github.com/gofiber/fiber/v3"

	resultPkg "github.com/blocknextai/go-packages/result"
	commonHTTP "github.com/blocknextai/platform-api/internal/common/http"
	authUseCases "github.com/blocknextai/platform-api/internal/modules/account/internal/usecases/auth"
)

type ConfirmEmailChangeRequest struct {
	Token string `json:"token"`
}

func NewConfirmEmailChangeHandler(service *authUseCases.Service) fiber.Handler {
	return func(c fiber.Ctx) error {
		request := new(ConfirmEmailChangeRequest)
		if err := c.Bind().All(request); err != nil {
			return commonHTTP.ErrInvalidRequest
		}

		result, err := service.ConfirmEmailChange(c.RequestCtx(), &authUseCases.ConfirmEmailChangeCommand{
			Token: request.Token,
		})

		if err != nil {
			return err
		}

		return c.Status(fiber.StatusOK).JSON(resultPkg.Ok(result, resultPkg.WithMessage("email change confirmed")))
	}
}
