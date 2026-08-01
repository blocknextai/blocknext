package sessions

import (
	resultPkg "github.com/blocknextai/go-packages/result"
	"github.com/blocknextai/platform-api/internal/account/application/sessions/logout"
	"github.com/blocknextai/platform-api/internal/common/application/cqrs"
	commonHTTP "github.com/blocknextai/platform-api/internal/common/presentation/http"
	"github.com/gofiber/fiber/v3"
)

func NewLogoutHandler(handler cqrs.Handler[*logout.LogoutCommand, *logout.LogoutResponse]) fiber.Handler {
	return func(c fiber.Ctx) error {
		sessionID := commonHTTP.GetSessionID(c)

		result, err := handler.Handle(c.RequestCtx(), &logout.LogoutCommand{
			SessionID: sessionID,
		})

		if err != nil {
			return err
		}

		return c.Status(fiber.StatusOK).JSON(resultPkg.Ok(result, resultPkg.WithMessage("logged out")))
	}
}
