package password

import (
	"github.com/gofiber/fiber/v3"

	resultPkg "github.com/blocknextai/go-packages/result"
	commonHTTP "github.com/blocknextai/platform-api/internal/common/http"
	authUseCases "github.com/blocknextai/platform-api/internal/modules/account/internal/usecases/auth"
)

type SetPasswordRequest struct {
	Password string `json:"password"`
}

func NewSetPasswordHandler(service *authUseCases.Service) fiber.Handler {
	return func(c fiber.Ctx) error {
		request := new(SetPasswordRequest)
		if err := c.Bind().All(request); err != nil {
			return commonHTTP.ErrInvalidRequest
		}

		userID := commonHTTP.GetUserID(c)

		result, err := service.SetPassword(c.RequestCtx(), &authUseCases.SetPasswordCommand{
			UserID:   userID,
			Password: request.Password,
		})

		if err != nil {
			return err
		}

		return c.Status(fiber.StatusAccepted).JSON(resultPkg.Ok(result, resultPkg.WithMessage("password set")))
	}
}
