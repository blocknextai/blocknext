package notifications

import (
	"context"
)

type NotificationRepository interface {
	Create(ctx context.Context, notification *Notification) error
}
