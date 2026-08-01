package password

import (
	resultPkg "github.com/blocknextai/go-packages/result"
	"github.com/blocknextai/platform-api/internal/account/application/auth/password/set"
	"github.com/blocknextai/platform-api/internal/common/application/cqrs"
	commonHTTP "github.com/blocknextai/platform-api/internal/common/presentation/http"
	"github.com/gofiber/fiber/v3"
)

type SetPasswordRequest struct {
	Password string `json:"password"`
}

func NewSetPasswordHandler(handler cqrs.Handler[*set.SetPasswordCommand, *set.SetPasswordResponse]) fiber.Handler {
	return func(c fiber.Ctx) error {
		request := new(SetPasswordRequest)
		if err := c.Bind().All(request); err != nil {
			return commonHTTP.ErrInvalidRequest
		}

		userID := commonHTTP.GetUserID(c)

		result, err := handler.Handle(c.RequestCtx(), &set.SetPasswordCommand{
			UserID:   userID,
			Password: request.Password,
		})

		if err != nil {
			return err
		}

		return c.Status(fiber.StatusAccepted).JSON(resultPkg.Ok(result, resultPkg.WithMessage("password set")))
	}
}
