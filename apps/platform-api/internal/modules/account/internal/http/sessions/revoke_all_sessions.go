package sessions

import (
	"github.com/gofiber/fiber/v3"

	resultPkg "github.com/blocknextai/go-packages/result"
	commonHTTP "github.com/blocknextai/platform-api/internal/common/http"
	sessionsUseCases "github.com/blocknextai/platform-api/internal/modules/account/internal/usecases/sessions"
)

func NewRevokeAllSessionsHandler(service *sessionsUseCases.Service) fiber.Handler {
	return func(c fiber.Ctx) error {
		userID := commonHTTP.GetUserID(c)

		result, err := service.RevokeAllSessions(c.RequestCtx(), &sessionsUseCases.RevokeAllSessionsCommand{
			UserID: userID,
		})

		if err != nil {
			return err
		}

		return c.Status(fiber.StatusOK).JSON(resultPkg.Ok(result, resultPkg.WithMessage("all sessions revoked")))
	}
}
