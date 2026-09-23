package auth

import (
	"github.com/gofiber/fiber/v3"

	resultPkg "github.com/blocknextai/go-packages/result"
	commonHTTP "github.com/blocknextai/platform-api/internal/common/http"
	accountDomain "github.com/blocknextai/platform-api/internal/modules/account/internal/domain"
	authUseCases "github.com/blocknextai/platform-api/internal/modules/account/internal/usecases/auth"
)

type CreateUserTokenRequest struct {
	AuthProvider accountDomain.AuthProvider `json:"authProvider"`
	Payload      map[string]any             `json:"payload"`
}

func NewCreateUserTokenHandler(service *authUseCases.Service) fiber.Handler {
	return func(c fiber.Ctx) error {
		request := new(CreateUserTokenRequest)
		if err := c.Bind().All(request); err != nil {
			return commonHTTP.ErrInvalidRequest
		}

		result, err := service.CreateUserToken(c.RequestCtx(), &authUseCases.CreateUserTokenCommand{
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
