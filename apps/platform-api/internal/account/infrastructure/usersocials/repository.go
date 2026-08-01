package usersocials

import (
	"context"
	"database/sql"

	"github.com/blocknextai/go-packages/database"
	"github.com/blocknextai/platform-api/internal/account/domain/usersocials"
	"github.com/google/uuid"
)

const (
	tableName = "account.socials"
	columns   = "id, user_id, platform, url, sort_order, created_at, updated_at, deleted_at"
)

var (
	queryGetByUserID = database.BuildQuery(`
		SELECT `, columns, `
		FROM `, tableName, `
		WHERE
			user_id = $1
			AND deleted_at IS NULL
		ORDER BY
			sort_order ASC
	`)

	queryCreate = database.BuildQuery(`
		INSERT INTO `, tableName, ` (`, columns, `)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
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

type UserSocialRepository struct {
	database.BaseRepository
}

func NewUserSocialRepository(db *sql.DB) usersocials.UserSocialRepository {
	return &UserSocialRepository{
		BaseRepository: database.NewBaseRepository(db),
	}
}

func (r *UserSocialRepository) scan(row interface{ Scan(dest ...any) error }, _ ...any) (*usersocials.UserSocial, error) {
	var social usersocials.UserSocial
	var deletedAt sql.NullTime

	err := row.Scan(
		&social.ID,
		&social.UserID,
		&social.Platform,
		&social.URL,
		&social.SortOrder,
		&social.CreatedAt,
		&social.UpdatedAt,
		&deletedAt,
	)

	if err != nil {
		return nil, err
	}

	social.DeletedAt = database.ScanNullTime(deletedAt)

	return &social, nil
}

func (r *UserSocialRepository) getMany(ctx context.Context, query string, args ...any) ([]*usersocials.UserSocial, error) {
	return database.GetMany(ctx, r.Executor(ctx), query, r.scan, args...)
}

func (r *UserSocialRepository) exec(ctx context.Context, query string, args ...any) error {
	return database.Exec(ctx, r.Executor(ctx), query, args...)
}

func (r *UserSocialRepository) execWithRowCheck(ctx context.Context, query string, args ...any) error {
	return database.ExecWithRowCheck(ctx, r.Executor(ctx), query, usersocials.ErrSocialNotFound, args...)
}

func (r *UserSocialRepository) GetAllByUserID(ctx context.Context, userID uuid.UUID) ([]*usersocials.UserSocial, error) {
	return r.getMany(ctx, queryGetByUserID, userID)
}

func (r *UserSocialRepository) Create(ctx context.Context, social *usersocials.UserSocial) error {
	return r.exec(ctx, queryCreate,
		social.ID,
		social.UserID,
		social.Platform,
		social.URL,
		social.SortOrder,
		social.CreatedAt,
		social.UpdatedAt,
		social.DeletedAt,
	)
}

func (r *UserSocialRepository) Delete(ctx context.Context, social *usersocials.UserSocial) error {
	return r.execWithRowCheck(ctx, queryDelete,
		social.ID,
		social.UpdatedAt,
		social.DeletedAt,
	)
}
