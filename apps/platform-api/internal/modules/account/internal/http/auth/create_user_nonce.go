package auth

import (
	"github.com/gofiber/fiber/v3"

	resultPkg "github.com/blocknextai/go-packages/result"
	commonHTTP "github.com/blocknextai/platform-api/internal/common/http"
	"github.com/blocknextai/platform-api/internal/modules/account/internal/domain"
	authUseCases "github.com/blocknextai/platform-api/internal/modules/account/internal/usecases/auth"
)

type CreateUserNonceRequest struct {
	AuthProvider domain.AuthProvider `json:"authProvider"`
	ProviderID   *string             `json:"providerId"`
}

func NewCreateUserNonceHandler(service *authUseCases.Service) fiber.Handler {
	return func(c fiber.Ctx) error {
		request := new(CreateUserNonceRequest)
		if err := c.Bind().All(request); err != nil {
			return commonHTTP.ErrInvalidRequest
		}

		result, err := service.CreateUserNonce(c.RequestCtx(), &authUseCases.CreateUserNonceCommand{
			AuthProvider: request.AuthProvider,
			ProviderID:   request.ProviderID,
		})

		if err != nil {
			return err
		}

		return c.Status(fiber.StatusCreated).JSON(resultPkg.Ok(result))
	}
}
