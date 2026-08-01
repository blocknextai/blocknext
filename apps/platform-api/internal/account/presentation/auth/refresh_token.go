package auth

import (
	resultPkg "github.com/blocknextai/go-packages/result"
	accountApplicationAuth "github.com/blocknextai/platform-api/internal/account/application/auth"
	"github.com/blocknextai/platform-api/internal/account/application/auth/refreshtoken"
	"github.com/blocknextai/platform-api/internal/common/application/cqrs"
	commonHTTP "github.com/blocknextai/platform-api/internal/common/presentation/http"
	"github.com/gofiber/fiber/v3"
)

type RefreshTokenRequest struct {
	RefreshToken string `json:"refreshToken"`
}

func NewRefreshTokenHandler(handler cqrs.Handler[*refreshtoken.RefreshTokenCommand, *accountApplicationAuth.AccessTokenResponse]) fiber.Handler {
	return func(c fiber.Ctx) error {
		request := new(RefreshTokenRequest)
		if err := c.Bind().All(request); err != nil {
			return commonHTTP.ErrInvalidRequest
		}

		result, err := handler.Handle(c.RequestCtx(), &refreshtoken.RefreshTokenCommand{
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
