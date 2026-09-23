package notificationrecipients

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type NotificationRecipientRepository interface {
	CreateMany(ctx context.Context, recipients []*NotificationRecipient) error
	GetByIDAndUserID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*NotificationRecipient, error)
	GetInboxByUserID(ctx context.Context, userID uuid.UUID, organizationID *uuid.UUID, searchQuery string, offset int, limit int) ([]*InboxItem, int64, error)
	CountsByUserID(ctx context.Context, userID uuid.UUID, organizationID *uuid.UUID) (unread int64, unseen int64, err error)
	MarkAllReadByUserID(ctx context.Context, userID uuid.UUID, organizationID *uuid.UUID, readAt time.Time) (int64, error)
	MarkAllSeenByUserID(ctx context.Context, userID uuid.UUID, organizationID *uuid.UUID, seenAt time.Time) (int64, error)
	Update(ctx context.Context, recipient *NotificationRecipient) error
	Delete(ctx context.Context, recipient *NotificationRecipient) error
}
