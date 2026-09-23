package notificationrecipients

import (
	"context"

	"github.com/google/uuid"
)

type GetNotificationCountsQuery struct {
	UserID         uuid.UUID
	OrganizationID *uuid.UUID
}

type GetNotificationCountsResponse struct {
	Unread int64 `json:"unread"`
	Unseen int64 `json:"unseen"`
}

func (s *Service) GetNotificationCounts(ctx context.Context, request *GetNotificationCountsQuery) (*GetNotificationCountsResponse, error) {
	unread, unseen, err := s.notificationRecipientRepository.CountsByUserID(ctx, request.UserID, request.OrganizationID)
	if err != nil {
		return nil, err
	}

	return &GetNotificationCountsResponse{
		Unread: unread,
		Unseen: unseen,
	}, nil
}
