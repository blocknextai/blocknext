package notificationrecipients

import (
	"context"
	"database/sql"
	"strconv"
	"strings"
	"time"

	"github.com/blocknextai/go-packages/database"
	notificationsDomainRecipients "github.com/blocknextai/platform-api/internal/notifications/domain/notificationrecipients"
	notificationsDomainNotifications "github.com/blocknextai/platform-api/internal/notifications/domain/notifications"
	"github.com/google/uuid"
)

const (
	tableName = "notifications.recipients"
	columns   = "r.id, r.notification_id, r.user_id, r.read_at, r.seen_at, r.created_at, r.updated_at, r.deleted_at"
	// IMPORTANT: we must remove this but required for the JOIN
	notificationColumns = "n.id, n.type, n.level, n.audience_type, n.audience_id, n.title, n.body, n.data, n.action_url, n.created_at, n.updated_at, n.deleted_at"
)

var (
	queryGetByIDAndUserID = database.BuildQuery(`
		SELECT `, columns, `
		FROM `, tableName, ` r
		WHERE
			r.id = $1
			AND r.user_id = $2
			AND r.deleted_at IS NULL
		LIMIT 1
	`)

	queryGetInboxByUserID = database.BuildQuery(`
		SELECT `, columns, `,
			`, notificationColumns, `,
			COUNT(*) OVER() AS total
		FROM `, tableName, ` r
		JOIN notifications.notifications n
			ON n.id = r.notification_id AND n.deleted_at IS NULL
		WHERE
			r.user_id = $1
			AND r.deleted_at IS NULL
			AND (
				($2::uuid IS NULL AND n.audience_type = 'user')
				OR (n.audience_type = 'organization' AND n.audience_id = $2)
			)
			AND ($3 = '' OR n.title ILIKE '%' || $3 || '%')
		ORDER BY r.created_at DESC
		LIMIT $4 OFFSET $5
	`)

	queryCountsByUserID = database.BuildQuery(`
		SELECT
			COUNT(*) FILTER (WHERE r.read_at IS NULL) AS unread,
			COUNT(*) FILTER (WHERE r.seen_at IS NULL) AS unseen
		FROM `, tableName, ` r
		JOIN notifications.notifications n ON n.id = r.notification_id AND n.deleted_at IS NULL
		WHERE
			r.user_id = $1
			AND r.deleted_at IS NULL
			AND (
				($2::uuid IS NULL AND n.audience_type = 'user')
				OR (n.audience_type = 'organization' AND n.audience_id = $2)
			)
	`)

	queryMarkAllReadByUserID = database.BuildQuery(`
		UPDATE `, tableName, `
		SET
			read_at = $2,
			seen_at = COALESCE(seen_at, $2),
			updated_at = $2
		WHERE
			user_id = $1
			AND deleted_at IS NULL
			AND read_at IS NULL
			AND notification_id IN (
				SELECT id
				FROM notifications.notifications
				WHERE
					deleted_at IS NULL
					AND (
						($3::uuid IS NULL AND audience_type = 'user')
						OR (audience_type = 'organization' AND audience_id = $3)
					)
			)
	`)

	queryMarkAllSeenByUserID = database.BuildQuery(`
		UPDATE `, tableName, `
		SET
			seen_at = $2,
			updated_at = $2
		WHERE
			user_id = $1
			AND deleted_at IS NULL
			AND seen_at IS NULL
			AND notification_id IN (
				SELECT id
				FROM notifications.notifications
				WHERE
					deleted_at IS NULL
					AND (
						($3::uuid IS NULL AND audience_type = 'user')
						OR (audience_type = 'organization' AND audience_id = $3)
					)
			)
	`)

	queryUpdate = database.BuildQuery(`
		UPDATE `, tableName, `
		SET
			read_at = $2,
			seen_at = $3,
			updated_at = $4
		WHERE
			id = $1
			AND deleted_at IS NULL
	`)

	queryDelete = database.BuildQuery(`
		UPDATE `, tableName, `
		SET
			updated_at = $2,
			deleted_at = $3
		WHERE
			id = $1
			AND deleted_at IS NULL
	`)
)

type NotificationRecipientRepository struct {
	database.BaseRepository
}

func NewNotificationRecipientRepository(db *sql.DB) notificationsDomainRecipients.NotificationRecipientRepository {
	return &NotificationRecipientRepository{
		BaseRepository: database.NewBaseRepository(db),
	}
}

func (r *NotificationRecipientRepository) CreateMany(
	ctx context.Context,
	recipients []*notificationsDomainRecipients.NotificationRecipient,
) error {
	if len(recipients) == 0 {
		return nil
	}

	var builder strings.Builder
	builder.WriteString("INSERT INTO ")
	builder.WriteString(tableName)
	builder.WriteString(" (id, notification_id, user_id, created_at, updated_at) VALUES ")

	args := make([]any, 0, len(recipients)*5)
	for i, recipient := range recipients {
		if i > 0 {
			builder.WriteString(", ")
		}

		base := i * 5
		builder.WriteString("($")
		builder.WriteString(strconv.Itoa(base + 1))
		builder.WriteString(", $")
		builder.WriteString(strconv.Itoa(base + 2))
		builder.WriteString(", $")
		builder.WriteString(strconv.Itoa(base + 3))
		builder.WriteString(", $")
		builder.WriteString(strconv.Itoa(base + 4))
		builder.WriteString(", $")
		builder.WriteString(strconv.Itoa(base + 5))
		builder.WriteString(")")

		args = append(args, recipient.ID, recipient.NotificationID, recipient.UserID, recipient.CreatedAt, recipient.UpdatedAt)
	}

	builder.WriteString(" ON CONFLICT DO NOTHING")

	return database.Exec(ctx, r.Executor(ctx), builder.String(), args...)
}

func (r *NotificationRecipientRepository) scan(row interface{ Scan(dest ...any) error }, extras ...any) (*notificationsDomainRecipients.NotificationRecipient, error) {
	var recipient notificationsDomainRecipients.NotificationRecipient
	var readAt, seenAt, deletedAt sql.NullTime

	dest := []any{
		&recipient.ID,
		&recipient.NotificationID,
		&recipient.UserID,
		&readAt,
		&seenAt,
		&recipient.CreatedAt,
		&recipient.UpdatedAt,
		&deletedAt,
	}
	dest = append(dest, extras...)

	if err := row.Scan(dest...); err != nil {
		return nil, err
	}

	recipient.ReadAt = database.ScanNullTime(readAt)
	recipient.SeenAt = database.ScanNullTime(seenAt)
	recipient.DeletedAt = database.ScanNullTime(deletedAt)

	return &recipient, nil
}

func (r *NotificationRecipientRepository) scanInboxItem(row interface{ Scan(dest ...any) error }, extras ...any) (*notificationsDomainRecipients.InboxItem, error) {
	var recipient notificationsDomainRecipients.NotificationRecipient
	var readAt, seenAt, recipientDeletedAt sql.NullTime

	var notification notificationsDomainNotifications.Notification
	var body, actionURL sql.NullString
	var data []byte
	var notificationDeletedAt sql.NullTime

	dest := []any{
		&recipient.ID,
		&recipient.NotificationID,
		&recipient.UserID,
		&readAt,
		&seenAt,
		&recipient.CreatedAt,
		&recipient.UpdatedAt,
		&recipientDeletedAt,

		&notification.ID,
		&notification.Type,
		&notification.Level,
		&notification.AudienceType,
		&notification.AudienceID,
		&notification.Title,
		&body,
		&data,
		&actionURL,
		&notification.CreatedAt,
		&notification.UpdatedAt,
		&notificationDeletedAt,
	}
	dest = append(dest, extras...)

	if err := row.Scan(dest...); err != nil {
		return nil, err
	}

	if err := database.UnmarshalJSONColumn(data, &notification.Data); err != nil {
		return nil, err
	}

	recipient.ReadAt = database.ScanNullTime(readAt)
	recipient.SeenAt = database.ScanNullTime(seenAt)
	recipient.DeletedAt = database.ScanNullTime(recipientDeletedAt)

	notification.Body = database.ScanNullString(body)
	notification.ActionURL = database.ScanNullString(actionURL)
	notification.DeletedAt = database.ScanNullTime(notificationDeletedAt)

	return &notificationsDomainRecipients.InboxItem{
		Recipient:    &recipient,
		Notification: &notification,
	}, nil
}

func (r *NotificationRecipientRepository) GetByIDAndUserID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*notificationsDomainRecipients.NotificationRecipient, error) {
	return database.GetOne(ctx, r.Executor(ctx), queryGetByIDAndUserID, r.scan, notificationsDomainRecipients.ErrNotificationRecipientNotFound, id, userID)
}

func (r *NotificationRecipientRepository) GetInboxByUserID(ctx context.Context, userID uuid.UUID, organizationID *uuid.UUID, searchQuery string, offset int, limit int) ([]*notificationsDomainRecipients.InboxItem, int64, error) {
	return database.GetManyWithTotal(ctx, r.Executor(ctx), queryGetInboxByUserID, r.scanInboxItem, userID, organizationID, searchQuery, limit, offset)
}

func (r *NotificationRecipientRepository) CountsByUserID(ctx context.Context, userID uuid.UUID, organizationID *uuid.UUID) (int64, int64, error) {
	var unread, unseen int64
	err := r.Executor(ctx).QueryRowContext(ctx, queryCountsByUserID, userID, organizationID).Scan(&unread, &unseen)
	if err != nil {
		return 0, 0, err
	}
	return unread, unseen, nil
}

func (r *NotificationRecipientRepository) MarkAllReadByUserID(ctx context.Context, userID uuid.UUID, organizationID *uuid.UUID, readAt time.Time) (int64, error) {
	result, err := r.Executor(ctx).ExecContext(ctx, queryMarkAllReadByUserID, userID, readAt, organizationID)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func (r *NotificationRecipientRepository) MarkAllSeenByUserID(ctx context.Context, userID uuid.UUID, organizationID *uuid.UUID, seenAt time.Time) (int64, error) {
	result, err := r.Executor(ctx).ExecContext(ctx, queryMarkAllSeenByUserID, userID, seenAt, organizationID)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func (r *NotificationRecipientRepository) Update(ctx context.Context, recipient *notificationsDomainRecipients.NotificationRecipient) error {
	return database.ExecWithRowCheck(ctx, r.Executor(ctx), queryUpdate, notificationsDomainRecipients.ErrNotificationRecipientNotFound,
		recipient.ID,
		recipient.ReadAt,
		recipient.SeenAt,
		recipient.UpdatedAt,
	)
}

func (r *NotificationRecipientRepository) Delete(ctx context.Context, recipient *notificationsDomainRecipients.NotificationRecipient) error {
	return database.ExecWithRowCheck(ctx, r.Executor(ctx), queryDelete, notificationsDomainRecipients.ErrNotificationRecipientNotFound,
		recipient.ID,
		recipient.UpdatedAt,
		recipient.DeletedAt,
	)
}
