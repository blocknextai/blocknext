package password

import (
	"github.com/gofiber/fiber/v3"

	resultPkg "github.com/blocknextai/go-packages/result"
	commonHTTP "github.com/blocknextai/platform-api/internal/common/http"
	authUseCases "github.com/blocknextai/platform-api/internal/modules/account/internal/usecases/auth"
)

type ChangePasswordRequest struct {
	CurrentPassword string `json:"currentPassword"`
	NewPassword     string `json:"newPassword"`
}

func NewChangePasswordHandler(service *authUseCases.Service) fiber.Handler {
	return func(c fiber.Ctx) error {
		request := new(ChangePasswordRequest)
		if err := c.Bind().All(request); err != nil {
			return commonHTTP.ErrInvalidRequest
		}

		userID := commonHTTP.GetUserID(c)

		result, err := service.ChangePassword(c.RequestCtx(), &authUseCases.ChangePasswordCommand{
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
