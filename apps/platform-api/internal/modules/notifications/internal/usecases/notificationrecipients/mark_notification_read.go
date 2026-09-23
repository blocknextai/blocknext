package notificationrecipients

import (
	"context"

	"github.com/google/uuid"
)

type MarkNotificationReadCommand struct {
	UserID      uuid.UUID
	RecipientID uuid.UUID
}

type MarkNotificationReadResponse struct {
	ID uuid.UUID `json:"id"`
}

func (s *Service) MarkNotificationRead(ctx context.Context, command *MarkNotificationReadCommand) (*MarkNotificationReadResponse, error) {
	err := s.transactionManager.ExecuteInTransaction(ctx, func(txCtx context.Context) error {
		recipient, err := s.notificationRecipientRepository.GetByIDAndUserID(txCtx, command.RecipientID, command.UserID)
		if err != nil {
			return err
		}

		updated, err := recipient.MarkRead()
		if err != nil {
			return err
		}

		return s.notificationRecipientRepository.Update(txCtx, updated)
	})

	if err != nil {
		return nil, err
	}

	return &MarkNotificationReadResponse{
		ID: command.RecipientID,
	}, nil
}
