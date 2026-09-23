package sessions

import (
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"

	resultPkg "github.com/blocknextai/go-packages/result"
	commonHTTP "github.com/blocknextai/platform-api/internal/common/http"
	sessionsUseCases "github.com/blocknextai/platform-api/internal/modules/account/internal/usecases/sessions"
)

type RevokeSessionRequest struct {
	SessionID uuid.UUID `uri:"sessionId"`
}

func NewRevokeSessionHandler(service *sessionsUseCases.Service) fiber.Handler {
	return func(c fiber.Ctx) error {
		request := new(RevokeSessionRequest)
		if err := c.Bind().All(request); err != nil {
			return commonHTTP.ErrInvalidRequest
		}

		userID := commonHTTP.GetUserID(c)

		result, err := service.RevokeSession(c.RequestCtx(), &sessionsUseCases.RevokeSessionCommand{
			UserID:    userID,
			SessionID: request.SessionID,
		})

		if err != nil {
			return err
		}

		return c.Status(fiber.StatusOK).JSON(resultPkg.Ok(result, resultPkg.WithMessage("session revoked")))
	}
}
