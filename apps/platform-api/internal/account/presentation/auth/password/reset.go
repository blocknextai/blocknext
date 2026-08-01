package password

import (
	resultPkg "github.com/blocknextai/go-packages/result"
	"github.com/blocknextai/platform-api/internal/account/application/auth/password/reset"
	"github.com/blocknextai/platform-api/internal/common/application/cqrs"
	commonHTTP "github.com/blocknextai/platform-api/internal/common/presentation/http"
	"github.com/gofiber/fiber/v3"
)

type ResetRequest struct {
	Token       string `json:"token"`
	NewPassword string `json:"newPassword"`
}

func NewResetHandler(handler cqrs.Handler[*reset.ResetCommand, *reset.ResetResponse]) fiber.Handler {
	return func(c fiber.Ctx) error {
		request := new(ResetRequest)
		if err := c.Bind().All(request); err != nil {
			return commonHTTP.ErrInvalidRequest
		}

		result, err := handler.Handle(c.RequestCtx(), &reset.ResetCommand{
			Token:       request.Token,
			NewPassword: request.NewPassword,
		})

		if err != nil {
			return err
		}

		return c.Status(fiber.StatusOK).JSON(resultPkg.Ok(result, resultPkg.WithMessage("password reset")))
	}
}
