package notificationrecipients

import (
	"github.com/blocknextai/platform-api/internal/modules/notifications/internal/domain/notifications"
)

type InboxItem struct {
	Recipient    *NotificationRecipient
	Notification *notifications.Notification
}
