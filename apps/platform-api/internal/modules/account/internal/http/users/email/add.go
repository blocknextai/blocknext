package email

import (
	"github.com/gofiber/fiber/v3"

	resultPkg "github.com/blocknextai/go-packages/result"
	commonHTTP "github.com/blocknextai/platform-api/internal/common/http"
	authUseCases "github.com/blocknextai/platform-api/internal/modules/account/internal/usecases/auth"
)

type AddEmailRequest struct {
	Email string `json:"email"`
}

func NewAddEmailHandler(service *authUseCases.Service) fiber.Handler {
	return func(c fiber.Ctx) error {
		request := new(AddEmailRequest)
		if err := c.Bind().All(request); err != nil {
			return commonHTTP.ErrInvalidRequest
		}

		userID := commonHTTP.GetUserID(c)

		result, err := service.AddEmail(c.RequestCtx(), &authUseCases.AddEmailCommand{
			UserID: userID,
			Email:  request.Email,
		})

		if err != nil {
			return err
		}

		return c.Status(fiber.StatusAccepted).JSON(resultPkg.Ok(result, resultPkg.WithMessage("email added")))
	}
}
