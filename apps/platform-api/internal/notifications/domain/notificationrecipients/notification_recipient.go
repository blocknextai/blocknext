package notificationrecipients

import (
	"time"

	"github.com/blocknextai/go-packages/database"
	bnuuid "github.com/blocknextai/go-packages/uuid"
	"github.com/google/uuid"
)

type NotificationRecipient struct {
	database.BaseEntity

	NotificationID uuid.UUID
	UserID         uuid.UUID
	ReadAt         *time.Time
	SeenAt         *time.Time
}

func New(
	notificationID uuid.UUID,
	userID uuid.UUID,
) (*NotificationRecipient, error) {
	utcNow := time.Now().UTC()

	recipient := &NotificationRecipient{
		BaseEntity: database.BaseEntity{
			ID:        bnuuid.NewV7(),
			CreatedAt: utcNow,
			UpdatedAt: utcNow,
			DeletedAt: nil,
		},
		NotificationID: notificationID,
		UserID:         userID,
	}

	return recipient.validateThenReturn()
}

func (r *NotificationRecipient) MarkRead() (*NotificationRecipient, error) {
	if r.ReadAt != nil {
		return r, nil
	}

	now := time.Now().UTC()

	if r.SeenAt == nil {
		r.SeenAt = new(now)
	}

	r.ReadAt = new(now)
	r.UpdatedAt = now

	return r.validateThenReturn()
}

func (r *NotificationRecipient) Delete() (*NotificationRecipient, error) {
	now := time.Now().UTC()

	r.DeletedAt = new(now)
	r.UpdatedAt = now

	return r.validateThenReturn()
}

func (r *NotificationRecipient) validateThenReturn() (*NotificationRecipient, error) {
	if r.NotificationID == uuid.Nil {
		return nil, ErrInvalidNotificationID
	}

	if r.UserID == uuid.Nil {
		return nil, ErrInvalidUserID
	}

	return r, nil
}
