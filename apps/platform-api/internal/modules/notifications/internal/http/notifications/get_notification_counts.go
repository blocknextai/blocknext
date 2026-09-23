package notifications

import (
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"

	resultPkg "github.com/blocknextai/go-packages/result"
	commonHTTP "github.com/blocknextai/platform-api/internal/common/http"
	notificationrecipientsUseCases "github.com/blocknextai/platform-api/internal/modules/notifications/internal/usecases/notificationrecipients"
)

func NewGetUserNotificationCountsHandler(service *notificationrecipientsUseCases.Service) fiber.Handler {
	return func(c fiber.Ctx) error {
		userID := commonHTTP.GetUserID(c)

		result, err := service.GetNotificationCounts(c.RequestCtx(), &notificationrecipientsUseCases.GetNotificationCountsQuery{
			UserID:         userID,
			OrganizationID: nil,
		})
		if err != nil {
			return err
		}

		return c.Status(fiber.StatusOK).JSON(resultPkg.Ok(result))
	}
}

type GetOrganizationNotificationCountsRequest struct {
	OrganizationID uuid.UUID `uri:"organizationId"`
}

func NewGetOrganizationNotificationCountsHandler(service *notificationrecipientsUseCases.Service) fiber.Handler {
	return func(c fiber.Ctx) error {
		request := new(GetOrganizationNotificationCountsRequest)
		if err := c.Bind().All(request); err != nil {
			return commonHTTP.ErrInvalidRequest
		}

		userID := commonHTTP.GetUserID(c)

		result, err := service.GetNotificationCounts(c.RequestCtx(), &notificationrecipientsUseCases.GetNotificationCountsQuery{
			UserID:         userID,
			OrganizationID: new(request.OrganizationID),
		})
		if err != nil {
			return err
		}

		return c.Status(fiber.StatusOK).JSON(resultPkg.Ok(result))
	}
}
