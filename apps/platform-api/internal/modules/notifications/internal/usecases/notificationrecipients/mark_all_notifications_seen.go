package notificationrecipients

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type MarkAllNotificationsSeenCommand struct {
	UserID         uuid.UUID
	OrganizationID *uuid.UUID
}

type MarkAllNotificationsSeenResponse struct {
	UpdatedCount int64 `json:"updatedCount"`
}

func (s *Service) MarkAllNotificationsSeen(ctx context.Context, command *MarkAllNotificationsSeenCommand) (*MarkAllNotificationsSeenResponse, error) {
	var updatedCount int64

	err := s.transactionManager.ExecuteInTransaction(ctx, func(txCtx context.Context) error {
		count, err := s.notificationRecipientRepository.MarkAllSeenByUserID(txCtx, command.UserID, command.OrganizationID, time.Now().UTC())
		if err != nil {
			return err
		}
		updatedCount = count
		return nil
	})

	if err != nil {
		return nil, err
	}

	return &MarkAllNotificationsSeenResponse{
		UpdatedCount: updatedCount,
	}, nil
}
