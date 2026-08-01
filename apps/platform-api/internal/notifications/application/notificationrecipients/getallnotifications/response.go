package getallnotifications

import (
	"time"

	"github.com/google/uuid"
)

type NotificationResponse struct {
	ID        uuid.UUID  `json:"id"`
	Level     string     `json:"level"`
	Title     string     `json:"title"`
	Body      *string    `json:"body,omitempty"`
	ActionURL *string    `json:"actionUrl,omitempty"`
	ReadAt    *time.Time `json:"readAt,omitempty"`
	SeenAt    *time.Time `json:"seenAt,omitempty"`
	CreatedAt time.Time  `json:"createdAt"`
}

type GetAllNotificationsResponse struct {
	Items      []*NotificationResponse
	TotalCount int64
}
