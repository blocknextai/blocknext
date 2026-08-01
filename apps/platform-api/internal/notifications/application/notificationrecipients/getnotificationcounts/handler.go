package getnotificationcounts

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

func (h *Handler) Handle(ctx context.Context, request *GetNotificationCountsQuery) (*GetNotificationCountsResponse, error) {
	unread, unseen, err := h.notificationRecipientRepository.CountsByUserID(ctx, request.UserID, request.OrganizationID)
	if err != nil {
		return nil, err
	}

	return &GetNotificationCountsResponse{
		Unread: unread,
		Unseen: unseen,
	}, nil
}
