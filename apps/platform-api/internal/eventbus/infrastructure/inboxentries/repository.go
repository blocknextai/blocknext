package inboxentries

import (
	"context"
	"database/sql"
	"errors"

	"github.com/blocknextai/go-packages/database"
	inboxEntriesDomain "github.com/blocknextai/platform-api/internal/eventbus/domain/inboxentries"
)

const (
	tableName = "eventbus.inbox_entries"
)

var (
	queryMarkProcessed = database.BuildQuery(`
		INSERT INTO `, tableName, ` (handler_key, event_id, processed_at)
		VALUES ($1, $2, $3)
		ON CONFLICT (handler_key, event_id) DO NOTHING
		RETURNING event_id
	`)
)

type Repository struct {
	database.BaseRepository
}

func NewRepository(db *sql.DB) inboxEntriesDomain.Repository {
	return &Repository{
		BaseRepository: database.NewBaseRepository(db),
	}
}

func (r *Repository) MarkProcessed(ctx context.Context, entry *inboxEntriesDomain.InboxEntry) (bool, error) {
	var inserted string
	err := r.Executor(ctx).QueryRowContext(ctx, queryMarkProcessed, entry.HandlerKey, entry.EventID, entry.ProcessedAt).Scan(&inserted)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
