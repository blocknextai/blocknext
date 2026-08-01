package deletenotification

import (
	"context"

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

func (h *Handler) Handle(ctx context.Context, command *DeleteNotificationCommand) (*DeleteNotificationResponse, error) {
	err := h.transactionManager.ExecuteInTransaction(ctx, func(txCtx context.Context) error {
		recipient, err := h.notificationRecipientRepository.GetByIDAndUserID(txCtx, command.RecipientID, command.UserID)
		if err != nil {
			return err
		}

		deleted, err := recipient.Delete()
		if err != nil {
			return err
		}

		return h.notificationRecipientRepository.Delete(txCtx, deleted)
	})

	if err != nil {
		return nil, err
	}

	return &DeleteNotificationResponse{
		ID: command.RecipientID,
	}, nil
}
