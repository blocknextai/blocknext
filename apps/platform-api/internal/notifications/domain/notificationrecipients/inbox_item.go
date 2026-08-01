package notificationrecipients

import (
	"github.com/blocknextai/platform-api/internal/notifications/domain/notifications"
)

type InboxItem struct {
	Recipient    *NotificationRecipient
	Notification *notifications.Notification
}
