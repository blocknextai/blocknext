package email

import (
	"github.com/gofiber/fiber/v3"

	resultPkg "github.com/blocknextai/go-packages/result"
	commonHTTP "github.com/blocknextai/platform-api/internal/common/http"
	authUseCases "github.com/blocknextai/platform-api/internal/modules/account/internal/usecases/auth"
)

type VerifyRequest struct {
	Token string `json:"token"`
}

func NewVerifyHandler(service *authUseCases.Service) fiber.Handler {
	return func(c fiber.Ctx) error {
		request := new(VerifyRequest)
		if err := c.Bind().All(request); err != nil {
			return commonHTTP.ErrInvalidRequest
		}

		result, err := service.VerifyEmail(c.RequestCtx(), &authUseCases.VerifyCommand{
			Token:     request.Token,
			IPAddress: c.IP(),
			UserAgent: c.Get(fiber.HeaderUserAgent),
		})

		if err != nil {
			return err
		}

		return c.Status(fiber.StatusOK).JSON(resultPkg.Ok(result, resultPkg.WithMessage("email verified")))
	}
}
