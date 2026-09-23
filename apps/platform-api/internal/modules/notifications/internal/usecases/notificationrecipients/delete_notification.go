package notificationrecipients

import (
	"context"

	"github.com/google/uuid"
)

type DeleteNotificationCommand struct {
	UserID      uuid.UUID
	RecipientID uuid.UUID
}

type DeleteNotificationResponse struct {
	ID uuid.UUID `json:"id"`
}

func (s *Service) DeleteNotification(ctx context.Context, command *DeleteNotificationCommand) (*DeleteNotificationResponse, error) {
	err := s.transactionManager.ExecuteInTransaction(ctx, func(txCtx context.Context) error {
		recipient, err := s.notificationRecipientRepository.GetByIDAndUserID(txCtx, command.RecipientID, command.UserID)
		if err != nil {
			return err
		}

		deleted, err := recipient.Delete()
		if err != nil {
			return err
		}

		return s.notificationRecipientRepository.Delete(txCtx, deleted)
	})

	if err != nil {
		return nil, err
	}

	return &DeleteNotificationResponse{
		ID: command.RecipientID,
	}, nil
}
