package notificationrecipients

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type MarkAllNotificationsReadCommand struct {
	UserID         uuid.UUID
	OrganizationID *uuid.UUID
}

type MarkAllNotificationsReadResponse struct {
	UpdatedCount int64 `json:"updatedCount"`
}

func (s *Service) MarkAllNotificationsRead(ctx context.Context, command *MarkAllNotificationsReadCommand) (*MarkAllNotificationsReadResponse, error) {
	var updatedCount int64

	err := s.transactionManager.ExecuteInTransaction(ctx, func(txCtx context.Context) error {
		count, err := s.notificationRecipientRepository.MarkAllReadByUserID(txCtx, command.UserID, command.OrganizationID, time.Now().UTC())
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
