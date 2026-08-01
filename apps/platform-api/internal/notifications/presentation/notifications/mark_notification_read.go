package notifications

import (
	resultPkg "github.com/blocknextai/go-packages/result"
	"github.com/blocknextai/platform-api/internal/common/application/cqrs"
	commonHTTP "github.com/blocknextai/platform-api/internal/common/presentation/http"
	"github.com/blocknextai/platform-api/internal/notifications/application/notificationrecipients/marknotificationread"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type MarkNotificationReadRequest struct {
	RecipientID uuid.UUID `uri:"recipientId"`
}

func NewMarkNotificationReadHandler(handler cqrs.Handler[*marknotificationread.MarkNotificationReadCommand, *marknotificationread.MarkNotificationReadResponse]) fiber.Handler {
	return func(c fiber.Ctx) error {
		request := new(MarkNotificationReadRequest)
		if err := c.Bind().All(request); err != nil {
			return commonHTTP.ErrInvalidRequest
		}

		userID := commonHTTP.GetUserID(c)

		result, err := handler.Handle(c.RequestCtx(), &marknotificationread.MarkNotificationReadCommand{
			UserID:      userID,
			RecipientID: request.RecipientID,
		})
		if err != nil {
			return err
		}

		return c.Status(fiber.StatusOK).JSON(resultPkg.Ok(result))
	}
}
