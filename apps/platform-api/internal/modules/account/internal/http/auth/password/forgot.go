package password

import (
	"github.com/gofiber/fiber/v3"

	resultPkg "github.com/blocknextai/go-packages/result"
	commonHTTP "github.com/blocknextai/platform-api/internal/common/http"
	authUseCases "github.com/blocknextai/platform-api/internal/modules/account/internal/usecases/auth"
)

type ForgotRequest struct {
	Email string `json:"email"`
}

func NewForgotHandler(service *authUseCases.Service) fiber.Handler {
	return func(c fiber.Ctx) error {
		request := new(ForgotRequest)
		if err := c.Bind().All(request); err != nil {
			return commonHTTP.ErrInvalidRequest
		}

		result, err := service.ForgotPassword(c.RequestCtx(), &authUseCases.ForgotCommand{
			Email: request.Email,
		})

		if err != nil {
			return err
		}

		return c.Status(fiber.StatusAccepted).JSON(resultPkg.Ok(result, resultPkg.WithMessage("password reset email sent")))
	}
}
