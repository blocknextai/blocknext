package sessions

import (
	resultPkg "github.com/blocknextai/go-packages/result"
	"github.com/blocknextai/platform-api/internal/account/application/sessions/revokesession"
	"github.com/blocknextai/platform-api/internal/common/application/cqrs"
	commonHTTP "github.com/blocknextai/platform-api/internal/common/presentation/http"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type RevokeSessionRequest struct {
	SessionID uuid.UUID `uri:"sessionId"`
}

func NewRevokeSessionHandler(handler cqrs.Handler[*revokesession.RevokeSessionCommand, *revokesession.RevokeSessionResponse]) fiber.Handler {
	return func(c fiber.Ctx) error {
		request := new(RevokeSessionRequest)
		if err := c.Bind().All(request); err != nil {
			return commonHTTP.ErrInvalidRequest
		}

		userID := commonHTTP.GetUserID(c)

		result, err := handler.Handle(c.RequestCtx(), &revokesession.RevokeSessionCommand{
			UserID:    userID,
			SessionID: request.SessionID,
		})

		if err != nil {
			return err
		}

		return c.Status(fiber.StatusOK).JSON(resultPkg.Ok(result, resultPkg.WithMessage("session revoked")))
	}
}
