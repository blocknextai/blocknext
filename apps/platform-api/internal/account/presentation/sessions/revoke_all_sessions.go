package sessions

import (
	resultPkg "github.com/blocknextai/go-packages/result"
	"github.com/blocknextai/platform-api/internal/account/application/sessions/revokeallsessions"
	"github.com/blocknextai/platform-api/internal/common/application/cqrs"
	commonHTTP "github.com/blocknextai/platform-api/internal/common/presentation/http"
	"github.com/gofiber/fiber/v3"
)

func NewRevokeAllSessionsHandler(handler cqrs.Handler[*revokeallsessions.RevokeAllSessionsCommand, *revokeallsessions.RevokeAllSessionsResponse]) fiber.Handler {
	return func(c fiber.Ctx) error {
		userID := commonHTTP.GetUserID(c)

		result, err := handler.Handle(c.RequestCtx(), &revokeallsessions.RevokeAllSessionsCommand{
			UserID: userID,
		})

		if err != nil {
			return err
		}

		return c.Status(fiber.StatusOK).JSON(resultPkg.Ok(result, resultPkg.WithMessage("all sessions revoked")))
	}
}
