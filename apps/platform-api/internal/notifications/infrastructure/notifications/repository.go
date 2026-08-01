package notifications

import (
	"context"
	"database/sql"

	"github.com/blocknextai/go-packages/database"
	"github.com/blocknextai/go-packages/json"
	notificationsDomainNotifications "github.com/blocknextai/platform-api/internal/notifications/domain/notifications"
)

const (
	tableName = "notifications.notifications"
	columns   = "id, type, level, audience_type, audience_id, title, body, data, action_url, created_at, updated_at, deleted_at"
)

var (
	queryCreate = database.BuildQuery(`
		INSERT INTO `, tableName, ` (`, columns, `)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	`)
)

type NotificationRepository struct {
	database.BaseRepository
}

func NewNotificationRepository(db *sql.DB) notificationsDomainNotifications.NotificationRepository {
	return &NotificationRepository{
		BaseRepository: database.NewBaseRepository(db),
	}
}

func (r *NotificationRepository) Create(ctx context.Context, notification *notificationsDomainNotifications.Notification) error {
	var data []byte
	if notification.Data != nil {
		marshaled, err := json.Marshal(notification.Data)
		if err != nil {
			return err
		}
		data = marshaled
	}

	return database.Exec(ctx, r.Executor(ctx), queryCreate,
		notification.ID,
		notification.Type,
		notification.Level,
		notification.AudienceType,
		notification.AudienceID,
		notification.Title,
		notification.Body,
		data,
		notification.ActionURL,
		notification.CreatedAt,
		notification.UpdatedAt,
		notification.DeletedAt,
	)
}
