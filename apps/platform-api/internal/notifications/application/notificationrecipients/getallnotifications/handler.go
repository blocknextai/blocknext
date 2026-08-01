package getallnotifications

import (
	"context"

	"github.com/blocknextai/platform-api/internal/notifications/domain/notificationrecipients"
)

type Handler struct {
	notificationRecipientRepository notificationrecipients.NotificationRecipientRepository
}

func New(
	notificationRecipientRepository notificationrecipients.NotificationRecipientRepository,
) *Handler {
	return &Handler{
		notificationRecipientRepository: notificationRecipientRepository,
	}
}

func (h *Handler) Handle(ctx context.Context, request *GetAllNotificationsQuery) (*GetAllNotificationsResponse, error) {
	items, totalCount, err := h.notificationRecipientRepository.GetInboxByUserID(
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
