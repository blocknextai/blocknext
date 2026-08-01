package messages

import (
	"context"
	"database/sql"

	"github.com/blocknextai/go-packages/database"
	"github.com/blocknextai/go-packages/json"
	"github.com/blocknextai/platform-api/internal/workflows/domain/generation/messages"
	"github.com/google/uuid"
)

const (
	tableName = "workflows.generation_messages"
	columns   = "id, session_id, role, content, metadata, created_at, updated_at, deleted_at"
)

var (
	queryGetBySessionID = database.BuildQuery(`
		SELECT `, columns, `, COUNT(*) OVER() AS total
		FROM `, tableName, `
		WHERE
			session_id = $1
			AND deleted_at IS NULL
		ORDER BY
			created_at ASC
		LIMIT $2 OFFSET $3
	`)

	queryCreate = database.BuildQuery(`
		INSERT INTO `, tableName, ` (`, columns, `)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`)

	queryDeleteBySessionID = database.BuildQuery(`
		DELETE FROM `, tableName, `
		WHERE
			session_id = $1
	`)
)

type MessageRepository struct {
	database.BaseRepository
}

func NewMessageRepository(db *sql.DB) messages.MessageRepository {
	return &MessageRepository{
		BaseRepository: database.NewBaseRepository(db),
	}
}

func (r *MessageRepository) scan(row interface{ Scan(dest ...any) error }, extras ...any) (*messages.GenerationMessage, error) {
	var msg messages.GenerationMessage
	var metadata []byte
	var deletedAt sql.NullTime

	dest := []any{
		&msg.ID,
		&msg.SessionID,
		&msg.Role,
		&msg.Content,
		&metadata,
		&msg.CreatedAt,
		&msg.UpdatedAt,
		&deletedAt,
	}
	dest = append(dest, extras...)

	if err := row.Scan(dest...); err != nil {
		return nil, err
	}

	if err := database.UnmarshalJSONColumn(metadata, &msg.Metadata); err != nil {
		return nil, err
	}

	msg.DeletedAt = database.ScanNullTime(deletedAt)

	return &msg, nil
}

func (r *MessageRepository) getManyWithTotal(ctx context.Context, query string, args ...any) ([]*messages.GenerationMessage, int64, error) {
	return database.GetManyWithTotal(ctx, r.Executor(ctx), query, r.scan, args...)
}

func (r *MessageRepository) exec(ctx context.Context, query string, args ...any) error {
	return database.Exec(ctx, r.Executor(ctx), query, args...)
}

func (r *MessageRepository) GetAllBySessionID(ctx context.Context, sessionID uuid.UUID, offset int, limit int) ([]*messages.GenerationMessage, int64, error) {
	return r.getManyWithTotal(ctx, queryGetBySessionID, sessionID, limit, offset)
}

func (r *MessageRepository) Create(ctx context.Context, message *messages.GenerationMessage) error {
	var metadataJSON []byte
	if message.Metadata != nil {
		var marshalErr error
		metadataJSON, marshalErr = json.Marshal(message.Metadata)
		if marshalErr != nil {
			return marshalErr
		}
	}

	return r.exec(ctx, queryCreate,
		message.ID,
		message.SessionID,
		message.Role,
		message.Content,
		metadataJSON,
		message.CreatedAt,
		message.UpdatedAt,
		message.DeletedAt,
	)
}

func (r *MessageRepository) DeleteBySessionID(ctx context.Context, sessionID uuid.UUID) error {
	return r.exec(ctx, queryDeleteBySessionID, sessionID)
}
