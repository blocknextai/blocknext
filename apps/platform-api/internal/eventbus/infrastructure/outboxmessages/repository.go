package outboxmessages

import (
	"context"
	"database/sql"
	"time"

	"github.com/blocknextai/go-packages/database"
	outboxMessagesDomain "github.com/blocknextai/platform-api/internal/eventbus/domain/outboxmessages"
)

const (
	tableName = "eventbus.outbox_messages"
	columns   = "id, event_name, payload, occurred_at, status, attempts, next_attempt_at, locked_at, last_error, processed_at"
)

var (
	queryCreate = database.BuildQuery(`
		INSERT INTO `, tableName, ` (`, columns, `)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`)

	queryClaimPending = database.BuildQuery(`
		UPDATE `, tableName, ` AS o
		SET
			status = 'processing',
			attempts = o.attempts + 1,
			locked_at = $1
		FROM (
			SELECT id AS due_id
			FROM `, tableName, `
			WHERE
				status = 'pending'
				AND next_attempt_at <= $1
			ORDER BY occurred_at
			LIMIT $2
			FOR UPDATE SKIP LOCKED
		) AS due
		WHERE o.id = due.due_id
		RETURNING `, columns, `
	`)

	queryUpdate = database.BuildQuery(`
		UPDATE `, tableName, `
		SET
			status = $2,
			attempts = $3,
			next_attempt_at = $4,
			locked_at = $5,
			last_error = $6,
			processed_at = $7
		WHERE id = $1
	`)

	queryReclaimStuck = database.BuildQuery(`
		UPDATE `, tableName, `
		SET
			status = 'pending',
			locked_at = NULL
		WHERE
			status = 'processing'
			AND locked_at < $1
	`)
)

type Repository struct {
	database.BaseRepository
}

func NewRepository(db *sql.DB) outboxMessagesDomain.Repository {
	return &Repository{
		BaseRepository: database.NewBaseRepository(db),
	}
}

func (r *Repository) scan(row interface{ Scan(dest ...any) error }, extras ...any) (*outboxMessagesDomain.OutboxMessage, error) {
	var message outboxMessagesDomain.OutboxMessage
	var lockedAt sql.NullTime
	var lastError sql.NullString
	var processedAt sql.NullTime

	dest := []any{
		&message.ID,
		&message.EventName,
		&message.Payload,
		&message.OccurredAt,
		&message.Status,
		&message.Attempts,
		&message.NextAttemptAt,
		&lockedAt,
		&lastError,
		&processedAt,
	}
	dest = append(dest, extras...)

	if err := row.Scan(dest...); err != nil {
		return nil, err
	}

	message.LockedAt = database.ScanNullTime(lockedAt)
	message.LastError = database.ScanNullString(lastError)
	message.ProcessedAt = database.ScanNullTime(processedAt)

	return &message, nil
}

func (r *Repository) Create(ctx context.Context, message *outboxMessagesDomain.OutboxMessage) error {
	return database.Exec(ctx, r.Executor(ctx), queryCreate,
		message.ID,
		message.EventName,
		message.Payload,
		message.OccurredAt,
		message.Status,
		message.Attempts,
		message.NextAttemptAt,
		message.LockedAt,
		message.LastError,
		message.ProcessedAt,
	)
}

func (r *Repository) ClaimPending(ctx context.Context, limit int, now time.Time) ([]*outboxMessagesDomain.OutboxMessage, error) {
	return database.GetMany(ctx, r.Executor(ctx), queryClaimPending, r.scan, now, limit)
}

func (r *Repository) Update(ctx context.Context, message *outboxMessagesDomain.OutboxMessage) error {
	return database.Exec(ctx, r.Executor(ctx), queryUpdate,
		message.ID,
		message.Status,
		message.Attempts,
		message.NextAttemptAt,
		message.LockedAt,
		message.LastError,
		message.ProcessedAt,
	)
}

func (r *Repository) ReclaimStuck(ctx context.Context, lockedBefore time.Time) (int, error) {
	result, err := r.Executor(ctx).ExecContext(ctx, queryReclaimStuck, lockedBefore)
	if err != nil {
		return 0, err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}
	return int(affected), nil
}
