package notifications

import (
	resultPkg "github.com/blocknextai/go-packages/result"
	"github.com/blocknextai/platform-api/internal/common/application/cqrs"
	commonHTTP "github.com/blocknextai/platform-api/internal/common/presentation/http"
	"github.com/blocknextai/platform-api/internal/notifications/application/notificationrecipients/markallnotificationsread"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

func NewMarkAllUserNotificationsReadHandler(handler cqrs.Handler[*markallnotificationsread.MarkAllNotificationsReadCommand, *markallnotificationsread.MarkAllNotificationsReadResponse]) fiber.Handler {
	return func(c fiber.Ctx) error {
		userID := commonHTTP.GetUserID(c)

		result, err := handler.Handle(c.RequestCtx(), &markallnotificationsread.MarkAllNotificationsReadCommand{
			UserID:         userID,
			OrganizationID: nil,
		})
		if err != nil {
			return err
		}

		return c.Status(fiber.StatusOK).JSON(resultPkg.Ok(result))
	}
}

type MarkAllOrganizationNotificationsReadRequest struct {
	OrganizationID uuid.UUID `uri:"organizationId"`
}

func NewMarkAllOrganizationNotificationsReadHandler(handler cqrs.Handler[*markallnotificationsread.MarkAllNotificationsReadCommand, *markallnotificationsread.MarkAllNotificationsReadResponse]) fiber.Handler {
	return func(c fiber.Ctx) error {
		request := new(MarkAllOrganizationNotificationsReadRequest)
		if err := c.Bind().All(request); err != nil {
			return commonHTTP.ErrInvalidRequest
		}

		userID := commonHTTP.GetUserID(c)

		result, err := handler.Handle(c.RequestCtx(), &markallnotificationsread.MarkAllNotificationsReadCommand{
			UserID:         userID,
			OrganizationID: new(request.OrganizationID),
		})
		if err != nil {
			return err
		}

		return c.Status(fiber.StatusOK).JSON(resultPkg.Ok(result))
	}
}
