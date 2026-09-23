package email

import (
	"github.com/gofiber/fiber/v3"

	resultPkg "github.com/blocknextai/go-packages/result"
	commonHTTP "github.com/blocknextai/platform-api/internal/common/http"
	authUseCases "github.com/blocknextai/platform-api/internal/modules/account/internal/usecases/auth"
)

type ChangeEmailRequest struct {
	NewEmail        string `json:"newEmail"`
	CurrentPassword string `json:"currentPassword"`
}

func NewChangeEmailHandler(service *authUseCases.Service) fiber.Handler {
	return func(c fiber.Ctx) error {
		request := new(ChangeEmailRequest)
		if err := c.Bind().All(request); err != nil {
			return commonHTTP.ErrInvalidRequest
		}

		userID := commonHTTP.GetUserID(c)

		result, err := service.ChangeEmail(c.RequestCtx(), &authUseCases.ChangeEmailCommand{
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
