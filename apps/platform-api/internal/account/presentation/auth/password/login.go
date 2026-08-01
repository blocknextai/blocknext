package password

import (
	resultPkg "github.com/blocknextai/go-packages/result"
	accountApplicationAuth "github.com/blocknextai/platform-api/internal/account/application/auth"
	"github.com/blocknextai/platform-api/internal/account/application/auth/password/login"
	"github.com/blocknextai/platform-api/internal/common/application/cqrs"
	commonHTTP "github.com/blocknextai/platform-api/internal/common/presentation/http"
	"github.com/gofiber/fiber/v3"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func NewLoginHandler(handler cqrs.Handler[*login.LoginCommand, *accountApplicationAuth.AccessTokenResponse]) fiber.Handler {
	return func(c fiber.Ctx) error {
		request := new(LoginRequest)
		if err := c.Bind().All(request); err != nil {
			return commonHTTP.ErrInvalidRequest
		}

		result, err := handler.Handle(c.RequestCtx(), &login.LoginCommand{
			Email:     request.Email,
			Password:  request.Password,
			IPAddress: c.IP(),
			UserAgent: c.Get(fiber.HeaderUserAgent),
		})

		if err != nil {
			return err
		}

		return c.Status(fiber.StatusOK).JSON(resultPkg.Ok(result, resultPkg.WithMessage("logged in")))
	}
}
