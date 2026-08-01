package email

import (
	resultPkg "github.com/blocknextai/go-packages/result"
	"github.com/blocknextai/platform-api/internal/account/application/auth/email/change"
	"github.com/blocknextai/platform-api/internal/common/application/cqrs"
	commonHTTP "github.com/blocknextai/platform-api/internal/common/presentation/http"
	"github.com/gofiber/fiber/v3"
)

type ChangeEmailRequest struct {
	NewEmail        string `json:"newEmail"`
	CurrentPassword string `json:"currentPassword"`
}

func NewChangeEmailHandler(handler cqrs.Handler[*change.ChangeEmailCommand, *change.ChangeEmailResponse]) fiber.Handler {
	return func(c fiber.Ctx) error {
		request := new(ChangeEmailRequest)
		if err := c.Bind().All(request); err != nil {
			return commonHTTP.ErrInvalidRequest
		}

		userID := commonHTTP.GetUserID(c)

		result, err := handler.Handle(c.RequestCtx(), &change.ChangeEmailCommand{
			UserID:          userID,
			NewEmail:        request.NewEmail,
			CurrentPassword: request.CurrentPassword,
		})

		if err != nil {
			return err
		}

		return c.Status(fiber.StatusAccepted).JSON(resultPkg.Ok(result, resultPkg.WithMessage("email change requested")))
	}
}
