package marknotificationread

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

func (h *Handler) Handle(ctx context.Context, command *MarkNotificationReadCommand) (*MarkNotificationReadResponse, error) {
	err := h.transactionManager.ExecuteInTransaction(ctx, func(txCtx context.Context) error {
		recipient, err := h.notificationRecipientRepository.GetByIDAndUserID(txCtx, command.RecipientID, command.UserID)
		if err != nil {
			return err
		}

		updated, err := recipient.MarkRead()
		if err != nil {
			return err
		}

		return h.notificationRecipientRepository.Update(txCtx, updated)
	})

	if err != nil {
		return nil, err
	}

	return &MarkNotificationReadResponse{
		ID: command.RecipientID,
	}, nil
}
