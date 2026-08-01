package password

import (
	resultPkg "github.com/blocknextai/go-packages/result"
	"github.com/blocknextai/platform-api/internal/account/application/auth/password/change"
	"github.com/blocknextai/platform-api/internal/common/application/cqrs"
	commonHTTP "github.com/blocknextai/platform-api/internal/common/presentation/http"
	"github.com/gofiber/fiber/v3"
)

type ChangePasswordRequest struct {
	CurrentPassword string `json:"currentPassword"`
	NewPassword     string `json:"newPassword"`
}

func NewChangePasswordHandler(handler cqrs.Handler[*change.ChangePasswordCommand, *change.ChangePasswordResponse]) fiber.Handler {
	return func(c fiber.Ctx) error {
		request := new(ChangePasswordRequest)
		if err := c.Bind().All(request); err != nil {
			return commonHTTP.ErrInvalidRequest
		}

		userID := commonHTTP.GetUserID(c)

		result, err := handler.Handle(c.RequestCtx(), &change.ChangePasswordCommand{
			UserID:          userID,
			CurrentPassword: request.CurrentPassword,
			NewPassword:     request.NewPassword,
		})

		if err != nil {
			return err
		}

		return c.Status(fiber.StatusOK).JSON(resultPkg.Ok(result, resultPkg.WithMessage("password changed")))
	}
}
