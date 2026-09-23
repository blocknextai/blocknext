package notificationrecipients

import (
	"context"
	"time"

	"github.com/google/uuid"

	resultPkg "github.com/blocknextai/go-packages/result"
	"github.com/blocknextai/platform-api/internal/modules/notifications/internal/domain/notificationrecipients"
)

type GetAllNotificationsQuery struct {
	UserID         uuid.UUID
	OrganizationID *uuid.UUID
	Search         resultPkg.SearchRequest
	Pagination     resultPkg.PaginationRequest
}

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

func (s *Service) GetAllNotifications(ctx context.Context, request *GetAllNotificationsQuery) (*GetAllNotificationsResponse, error) {
	items, totalCount, err := s.notificationRecipientRepository.GetInboxByUserID(
		ctx,
		request.UserID,
		request.OrganizationID,
		request.Search.Query,
		request.Pagination.Offset,
		request.Pagination.Limit,
	)
	if err != nil {
		return nil, err
	}

	return &GetAllNotificationsResponse{
		Items:      MapInboxItemsToResponse(items),
		TotalCount: totalCount,
	}, nil
}

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
