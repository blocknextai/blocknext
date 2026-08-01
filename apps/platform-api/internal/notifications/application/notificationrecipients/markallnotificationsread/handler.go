package markallnotificationsread

import (
	"context"
	"time"

	"github.com/blocknextai/go-packages/database"
	"github.com/blocknextai/platform-api/internal/notifications/domain/notificationrecipients"
)

type Handler struct {
	notificationRecipientRepository notificationrecipients.NotificationRecipientRepository
	transactionManager              database.TransactionManager
}

func New(
	notificationRecipientRepository notificationrecipients.NotificationRecipientRepository,
	transactionManager database.TransactionManager,
) *Handler {
	return &Handler{
		notificationRecipientRepository: notificationRecipientRepository,
		transactionManager:              transactionManager,
	}
}

func (h *Handler) Handle(ctx context.Context, command *MarkAllNotificationsReadCommand) (*MarkAllNotificationsReadResponse, error) {
	var updatedCount int64

	err := h.transactionManager.ExecuteInTransaction(ctx, func(txCtx context.Context) error {
		count, err := h.notificationRecipientRepository.MarkAllReadByUserID(txCtx, command.UserID, command.OrganizationID, time.Now().UTC())
		if err != nil {
			return err
		}
		updatedCount = count
		return nil
	})

	if err != nil {
		return nil, err
	}

	return &MarkAllNotificationsReadResponse{
		UpdatedCount: updatedCount,
	}, nil
}
