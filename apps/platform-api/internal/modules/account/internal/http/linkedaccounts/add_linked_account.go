package linkedaccounts

import (
	"github.com/gofiber/fiber/v3"

	resultPkg "github.com/blocknextai/go-packages/result"
	commonHTTP "github.com/blocknextai/platform-api/internal/common/http"
	accountDomain "github.com/blocknextai/platform-api/internal/modules/account/internal/domain"
	linkedaccountsUseCases "github.com/blocknextai/platform-api/internal/modules/account/internal/usecases/linkedaccounts"
)

type AddLinkedAccountRequest struct {
	AuthProvider accountDomain.AuthProvider `json:"authProvider"`
	Payload      map[string]any             `json:"payload"`
}

func NewAddLinkedAccountHandler(service *linkedaccountsUseCases.Service) fiber.Handler {
	return func(c fiber.Ctx) error {
		request := new(AddLinkedAccountRequest)
		if err := c.Bind().All(request); err != nil {
			return commonHTTP.ErrInvalidRequest
		}

		userID := commonHTTP.GetUserID(c)

		result, err := service.AddLinkedAccount(c.RequestCtx(), &linkedaccountsUseCases.AddLinkedAccountCommand{
			UserID:       userID,
			AuthProvider: request.AuthProvider,
			Payload:      request.Payload,
		})
		if err != nil {
			return err
		}

		return c.Status(fiber.StatusCreated).JSON(resultPkg.Ok(result, resultPkg.WithMessage("linked account added")))
	}
}
