package notifications

import (
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"

	resultPkg "github.com/blocknextai/go-packages/result"
	commonHTTP "github.com/blocknextai/platform-api/internal/common/http"
	notificationrecipientsUseCases "github.com/blocknextai/platform-api/internal/modules/notifications/internal/usecases/notificationrecipients"
)

type DeleteNotificationRequest struct {
	RecipientID uuid.UUID `uri:"recipientId"`
}

func NewDeleteNotificationHandler(service *notificationrecipientsUseCases.Service) fiber.Handler {
	return func(c fiber.Ctx) error {
		request := new(DeleteNotificationRequest)
		if err := c.Bind().All(request); err != nil {
			return commonHTTP.ErrInvalidRequest
		}

		userID := commonHTTP.GetUserID(c)

		result, err := service.DeleteNotification(c.RequestCtx(), &notificationrecipientsUseCases.DeleteNotificationCommand{
			UserID:      userID,
			RecipientID: request.RecipientID,
		})
		if err != nil {
			return err
		}

		return c.Status(fiber.StatusOK).JSON(resultPkg.Ok(result, resultPkg.WithMessage("notification deleted")))
	}
}
