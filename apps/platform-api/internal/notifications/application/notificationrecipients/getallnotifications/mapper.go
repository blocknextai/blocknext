package getallnotifications

import (
	"github.com/blocknextai/platform-api/internal/notifications/domain/notificationrecipients"
)

func MapInboxItemsToResponse(items []*notificationrecipients.InboxItem) []*NotificationResponse {
	responses := make([]*NotificationResponse, 0, len(items))
	for _, item := range items {
		responses = append(responses, &NotificationResponse{
			ID:        item.Recipient.ID,
			Level:     item.Notification.Level.String(),
			Title:     item.Notification.Title,
			Body:      item.Notification.Body,
			ActionURL: item.Notification.ActionURL,
			ReadAt:    item.Recipient.ReadAt,
			SeenAt:    item.Recipient.SeenAt,
			CreatedAt: item.Recipient.CreatedAt,
		})
	}
	return responses
}
