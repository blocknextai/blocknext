package auth

import (
	"github.com/gofiber/fiber/v3"

	resultPkg "github.com/blocknextai/go-packages/result"
	commonHTTP "github.com/blocknextai/platform-api/internal/common/http"
	authUseCases "github.com/blocknextai/platform-api/internal/modules/account/internal/usecases/auth"
)

type RefreshTokenRequest struct {
	RefreshToken string `json:"refreshToken"`
}

func NewRefreshTokenHandler(service *authUseCases.Service) fiber.Handler {
	return func(c fiber.Ctx) error {
		request := new(RefreshTokenRequest)
		if err := c.Bind().All(request); err != nil {
			return commonHTTP.ErrInvalidRequest
		}

		result, err := service.RefreshToken(c.RequestCtx(), &authUseCases.RefreshTokenCommand{
			RefreshToken: request.RefreshToken,
			IPAddress:    c.IP(),
			UserAgent:    c.Get(fiber.HeaderUserAgent),
		})

		if err != nil {
			return err
		}

		return c.Status(fiber.StatusOK).JSON(resultPkg.Ok(result))
	}
}
