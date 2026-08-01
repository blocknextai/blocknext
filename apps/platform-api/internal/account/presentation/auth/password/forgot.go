package password

import (
	resultPkg "github.com/blocknextai/go-packages/result"
	"github.com/blocknextai/platform-api/internal/account/application/auth/password/forgot"
	"github.com/blocknextai/platform-api/internal/common/application/cqrs"
	commonHTTP "github.com/blocknextai/platform-api/internal/common/presentation/http"
	"github.com/gofiber/fiber/v3"
)

type ForgotRequest struct {
	Email string `json:"email"`
}

func NewForgotHandler(handler cqrs.Handler[*forgot.ForgotCommand, *forgot.ForgotResponse]) fiber.Handler {
	return func(c fiber.Ctx) error {
		request := new(ForgotRequest)
		if err := c.Bind().All(request); err != nil {
			return commonHTTP.ErrInvalidRequest
		}

		result, err := handler.Handle(c.RequestCtx(), &forgot.ForgotCommand{
			Email: request.Email,
		})

		if err != nil {
			return err
		}

		return c.Status(fiber.StatusAccepted).JSON(resultPkg.Ok(result, resultPkg.WithMessage("password reset email sent")))
	}
}
