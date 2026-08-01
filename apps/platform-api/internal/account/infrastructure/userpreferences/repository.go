package userpreferences

import (
	"context"
	"database/sql"

	"github.com/blocknextai/go-packages/database"
	"github.com/blocknextai/platform-api/internal/account/domain/userpreferences"
	"github.com/google/uuid"
)

const (
	tableName = "account.user_preferences"
	columns   = "id, user_id, theme_mode, theme_color, language, created_at, updated_at, deleted_at"
)

var (
	queryGetByUserID = database.BuildQuery(`
		SELECT `, columns, `
		FROM `, tableName, `
		WHERE
			user_id = $1
			AND deleted_at IS NULL
		LIMIT 1
	`)

	queryUpsert = database.BuildQuery(`
		INSERT INTO `, tableName, ` (`, columns, `)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		ON CONFLICT (user_id) WHERE deleted_at IS NULL
		DO UPDATE SET
			theme_mode = $3,
			theme_color = $4,
			language = $5,
			updated_at = $7
	`)
)

type UserPreferenceRepository struct {
	database.BaseRepository
}

func NewUserPreferenceRepository(db *sql.DB) userpreferences.UserPreferenceRepository {
	return &UserPreferenceRepository{
		BaseRepository: database.NewBaseRepository(db),
	}
}

func (r *UserPreferenceRepository) scan(row interface{ Scan(dest ...any) error }, _ ...any) (*userpreferences.UserPreference, error) {
	var pref userpreferences.UserPreference
	var deletedAt sql.NullTime

	err := row.Scan(
		&pref.ID,
		&pref.UserID,
		&pref.ThemeMode,
		&pref.ThemeColor,
		&pref.Language,
		&pref.CreatedAt,
		&pref.UpdatedAt,
		&deletedAt,
	)
	if err != nil {
		return nil, err
	}

	pref.DeletedAt = database.ScanNullTime(deletedAt)

	return &pref, nil
}

func (r *UserPreferenceRepository) getOne(ctx context.Context, query string, args ...any) (*userpreferences.UserPreference, error) {
	return database.GetOne(ctx, r.Executor(ctx), query, r.scan, userpreferences.ErrPreferenceNotFound, args...)
}

func (r *UserPreferenceRepository) exec(ctx context.Context, query string, args ...any) error {
	return database.Exec(ctx, r.Executor(ctx), query, args...)
}

func (r *UserPreferenceRepository) GetByUserID(ctx context.Context, userID uuid.UUID) (*userpreferences.UserPreference, error) {
	return r.getOne(ctx, queryGetByUserID, userID)
}

func (r *UserPreferenceRepository) Upsert(ctx context.Context, pref *userpreferences.UserPreference) error {
	return r.exec(ctx, queryUpsert,
		pref.ID,
		pref.UserID,
		pref.ThemeMode,
		pref.ThemeColor,
		pref.Language,
		pref.CreatedAt,
		pref.UpdatedAt,
		pref.DeletedAt,
	)
}
