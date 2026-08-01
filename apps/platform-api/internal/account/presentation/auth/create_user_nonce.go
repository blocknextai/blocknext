package auth

import (
	resultPkg "github.com/blocknextai/go-packages/result"
	"github.com/blocknextai/platform-api/internal/account/application/auth/createusernonce"
	"github.com/blocknextai/platform-api/internal/account/domain"
	"github.com/blocknextai/platform-api/internal/common/application/cqrs"
	commonHTTP "github.com/blocknextai/platform-api/internal/common/presentation/http"
	"github.com/gofiber/fiber/v3"
)

type CreateUserNonceRequest struct {
	AuthProvider domain.AuthProvider `json:"authProvider"`
	ProviderID   *string             `json:"providerId"`
}

func NewCreateUserNonceHandler(handler cqrs.Handler[*createusernonce.CreateUserNonceCommand, *createusernonce.CreateUserNonceResponse]) fiber.Handler {
	return func(c fiber.Ctx) error {
		request := new(CreateUserNonceRequest)
		if err := c.Bind().All(request); err != nil {
			return commonHTTP.ErrInvalidRequest
		}

		result, err := handler.Handle(c.RequestCtx(), &createusernonce.CreateUserNonceCommand{
			AuthProvider: request.AuthProvider,
			ProviderID:   request.ProviderID,
		})

		if err != nil {
			return err
		}

		return c.Status(fiber.StatusCreated).JSON(resultPkg.Ok(result))
	}
}
