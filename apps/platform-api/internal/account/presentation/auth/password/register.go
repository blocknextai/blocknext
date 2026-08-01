package password

import (
	resultPkg "github.com/blocknextai/go-packages/result"
	accountApplicationAuth "github.com/blocknextai/platform-api/internal/account/application/auth"
	"github.com/blocknextai/platform-api/internal/account/application/auth/password/register"
	"github.com/blocknextai/platform-api/internal/common/application/cqrs"
	commonHTTP "github.com/blocknextai/platform-api/internal/common/presentation/http"
	"github.com/gofiber/fiber/v3"
)

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func NewRegisterHandler(handler cqrs.Handler[*register.RegisterCommand, *accountApplicationAuth.AccessTokenResponse]) fiber.Handler {
	return func(c fiber.Ctx) error {
		request := new(RegisterRequest)
		if err := c.Bind().All(request); err != nil {
			return commonHTTP.ErrInvalidRequest
		}

		result, err := handler.Handle(c.RequestCtx(), &register.RegisterCommand{
			Email:     request.Email,
			Password:  request.Password,
			IPAddress: c.IP(),
			UserAgent: c.Get(fiber.HeaderUserAgent),
		})

		if err != nil {
			return err
		}

		return c.Status(fiber.StatusCreated).JSON(resultPkg.Ok(result, resultPkg.WithMessage("registered")))
	}
}
