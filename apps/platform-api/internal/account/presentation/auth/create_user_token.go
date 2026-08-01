package auth

import (
	resultPkg "github.com/blocknextai/go-packages/result"
	accountApplicationAuth "github.com/blocknextai/platform-api/internal/account/application/auth"
	"github.com/blocknextai/platform-api/internal/account/application/auth/createusertoken"
	accountDomain "github.com/blocknextai/platform-api/internal/account/domain"
	"github.com/blocknextai/platform-api/internal/common/application/cqrs"
	commonHTTP "github.com/blocknextai/platform-api/internal/common/presentation/http"
	"github.com/gofiber/fiber/v3"
)

type CreateUserTokenRequest struct {
	AuthProvider accountDomain.AuthProvider `json:"authProvider"`
	Payload      map[string]any             `json:"payload"`
}

func NewCreateUserTokenHandler(handler cqrs.Handler[*createusertoken.CreateUserTokenCommand, *accountApplicationAuth.AccessTokenResponse]) fiber.Handler {
	return func(c fiber.Ctx) error {
		request := new(CreateUserTokenRequest)
		if err := c.Bind().All(request); err != nil {
			return commonHTTP.ErrInvalidRequest
		}

		result, err := handler.Handle(c.RequestCtx(), &createusertoken.CreateUserTokenCommand{
			AuthProvider: request.AuthProvider,
			Payload:      request.Payload,
			IPAddress:    c.IP(),
			UserAgent:    c.Get(fiber.HeaderUserAgent),
		})

		if err != nil {
			return err
		}

		return c.Status(fiber.StatusCreated).JSON(resultPkg.Ok(result))
	}
}
