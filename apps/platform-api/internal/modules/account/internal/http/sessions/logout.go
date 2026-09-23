package sessions

import (
	"github.com/gofiber/fiber/v3"

	resultPkg "github.com/blocknextai/go-packages/result"
	commonHTTP "github.com/blocknextai/platform-api/internal/common/http"
	sessionsUseCases "github.com/blocknextai/platform-api/internal/modules/account/internal/usecases/sessions"
)

func NewLogoutHandler(service *sessionsUseCases.Service) fiber.Handler {
	return func(c fiber.Ctx) error {
		sessionID := commonHTTP.GetSessionID(c)

		result, err := service.Logout(c.RequestCtx(), &sessionsUseCases.LogoutCommand{
			SessionID: sessionID,
		})

		if err != nil {
			return err
		}

		return c.Status(fiber.StatusOK).JSON(resultPkg.Ok(result, resultPkg.WithMessage("logged out")))
	}
}
