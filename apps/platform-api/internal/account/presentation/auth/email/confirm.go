package email

import (
	resultPkg "github.com/blocknextai/go-packages/result"
	"github.com/blocknextai/platform-api/internal/account/application/auth/email/confirm"
	"github.com/blocknextai/platform-api/internal/common/application/cqrs"
	commonHTTP "github.com/blocknextai/platform-api/internal/common/presentation/http"
	"github.com/gofiber/fiber/v3"
)

type ConfirmEmailChangeRequest struct {
	Token string `json:"token"`
}

func NewConfirmEmailChangeHandler(handler cqrs.Handler[*confirm.ConfirmEmailChangeCommand, *confirm.ConfirmEmailChangeResponse]) fiber.Handler {
	return func(c fiber.Ctx) error {
		request := new(ConfirmEmailChangeRequest)
		if err := c.Bind().All(request); err != nil {
			return commonHTTP.ErrInvalidRequest
		}

		result, err := handler.Handle(c.RequestCtx(), &confirm.ConfirmEmailChangeCommand{
			Token: request.Token,
		})

		if err != nil {
			return err
		}

		return c.Status(fiber.StatusOK).JSON(resultPkg.Ok(result, resultPkg.WithMessage("email change confirmed")))
	}
}
