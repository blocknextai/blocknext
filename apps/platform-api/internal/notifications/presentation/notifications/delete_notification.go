package notifications

import (
	resultPkg "github.com/blocknextai/go-packages/result"
	"github.com/blocknextai/platform-api/internal/common/application/cqrs"
	commonHTTP "github.com/blocknextai/platform-api/internal/common/presentation/http"
	"github.com/blocknextai/platform-api/internal/notifications/application/notificationrecipients/deletenotification"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type DeleteNotificationRequest struct {
	RecipientID uuid.UUID `uri:"recipientId"`
}

func NewDeleteNotificationHandler(handler cqrs.Handler[*deletenotification.DeleteNotificationCommand, *deletenotification.DeleteNotificationResponse]) fiber.Handler {
	return func(c fiber.Ctx) error {
		request := new(DeleteNotificationRequest)
		if err := c.Bind().All(request); err != nil {
			return commonHTTP.ErrInvalidRequest
		}

		userID := commonHTTP.GetUserID(c)

		result, err := handler.Handle(c.RequestCtx(), &deletenotification.DeleteNotificationCommand{
			UserID:      userID,
			RecipientID: request.RecipientID,
		})
		if err != nil {
			return err
		}

		return c.Status(fiber.StatusOK).JSON(resultPkg.Ok(result, resultPkg.WithMessage("notification deleted")))
	}
}
