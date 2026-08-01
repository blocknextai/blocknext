package sessions

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/blocknextai/go-packages/database"
	"github.com/blocknextai/platform-api/internal/account/domain/sessions"
	"github.com/google/uuid"
)

const (
	tableName = "account.sessions"
	columns   = "id, user_id, auth_provider, ip_address, user_agent, is_revoked, refresh_token_hash, refresh_token_expires_at, token_family, token_generation, created_at, updated_at, deleted_at"
)

var (
	queryCreate = database.BuildQuery(`
		INSERT INTO `, tableName, ` (`, columns, `)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
	`)

	queryGetActiveByUserID = database.BuildQuery(`
		SELECT `, columns, `, COUNT(*) OVER() AS total
		FROM `, tableName, `
		WHERE
			user_id = $1
			AND is_revoked = FALSE
			AND deleted_at IS NULL
			AND ($2 = '' OR ip_address ILIKE '%' || $2 || '%' OR user_agent ILIKE '%' || $2 || '%')
		ORDER BY
			created_at DESC
		LIMIT $3 OFFSET $4
	`)

	queryGetByID = database.BuildQuery(`
		SELECT `, columns, `
		FROM `, tableName, `
		WHERE
			id = $1
			AND deleted_at IS NULL
	`)

	queryGetByIDAndUserID = database.BuildQuery(`
		SELECT `, columns, `
		FROM `, tableName, `
		WHERE
			id = $1
			AND user_id = $2
			AND deleted_at IS NULL
	`)

	queryRevokeByID = database.BuildQuery(`
		UPDATE `, tableName, `
		SET
			is_revoked = TRUE,
			updated_at = $2
		WHERE
			id = $1
			AND deleted_at IS NULL
	`)

	queryRevokeAllByUserID = database.BuildQuery(`
		UPDATE `, tableName, `
		SET
			is_revoked = TRUE,
			updated_at = $2
		WHERE
			user_id = $1
			AND is_revoked = FALSE
			AND deleted_at IS NULL
		RETURNING id
	`)

	queryRevokeAllByTokenFamily = database.BuildQuery(`
		UPDATE `, tableName, `
		SET
			is_revoked = TRUE,
			updated_at = $2
		WHERE
			token_family = $1
			AND is_revoked = FALSE
			AND deleted_at IS NULL
		RETURNING id
	`)

	queryUpdateRefreshToken = database.BuildQuery(`
		UPDATE `, tableName, `
		SET
			refresh_token_hash = $2,
			refresh_token_expires_at = $3,
			token_generation = $4,
			updated_at = $5
		WHERE
			id = $1
			AND token_generation = $6
			AND deleted_at IS NULL
	`)
)

type SessionRepository struct {
	database.BaseRepository
}

func NewSessionRepository(db *sql.DB) sessions.SessionRepository {
	return &SessionRepository{
		BaseRepository: database.NewBaseRepository(db),
	}
}

func (r *SessionRepository) scan(row interface{ Scan(dest ...any) error }, extras ...any) (*sessions.Session, error) {
	var s sessions.Session
	var deletedAt sql.NullTime
	var refreshTokenHash sql.NullString
	var refreshTokenExpiresAt sql.NullTime
	var tokenFamily uuid.NullUUID
	var tokenGeneration sql.NullInt32

	dest := []any{
		&s.ID,
		&s.UserID,
		&s.AuthProvider,
		&s.IPAddress,
		&s.UserAgent,
		&s.IsRevoked,
		&refreshTokenHash,
		&refreshTokenExpiresAt,
		&tokenFamily,
		&tokenGeneration,
		&s.CreatedAt,
		&s.UpdatedAt,
		&deletedAt,
	}
	dest = append(dest, extras...)

	if err := row.Scan(dest...); err != nil {
		return nil, err
	}

	if tokenGeneration.Valid {
		s.TokenGeneration = int(tokenGeneration.Int32)
	}

	s.RefreshTokenHash = database.ScanNullString(refreshTokenHash)
	s.RefreshTokenExpiresAt = database.ScanNullTime(refreshTokenExpiresAt)
	s.TokenFamily = database.ScanNullUUID(tokenFamily)
	s.DeletedAt = database.ScanNullTime(deletedAt)

	return &s, nil
}

func (r *SessionRepository) getOne(ctx context.Context, query string, args ...any) (*sessions.Session, error) {
	return database.GetOne(ctx, r.Executor(ctx), query, r.scan, sessions.ErrSessionNotFound, args...)
}

func (r *SessionRepository) getManyWithTotal(ctx context.Context, query string, args ...any) ([]*sessions.Session, int64, error) {
	return database.GetManyWithTotal(ctx, r.Executor(ctx), query, r.scan, args...)
}

func (r *SessionRepository) exec(ctx context.Context, query string, args ...any) error {
	return database.Exec(ctx, r.Executor(ctx), query, args...)
}

func (r *SessionRepository) execWithRowCheck(ctx context.Context, query string, args ...any) error {
	return database.ExecWithRowCheck(ctx, r.Executor(ctx), query, sessions.ErrSessionNotFound, args...)
}

func (r *SessionRepository) execReturnIDs(ctx context.Context, query string, args ...any) (ids []uuid.UUID, err error) {
	rows, err := r.Executor(ctx).QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer func() {
		err = errors.Join(err, rows.Close())
	}()

	for rows.Next() {
		var id uuid.UUID
		if scanErr := rows.Scan(&id); scanErr != nil {
			return nil, scanErr
		}
		ids = append(ids, id)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return ids, nil
}

func (r *SessionRepository) Create(ctx context.Context, session *sessions.Session) error {
	return r.exec(ctx, queryCreate,
		session.ID,
		session.UserID,
		session.AuthProvider,
		session.IPAddress,
		session.UserAgent,
		session.IsRevoked,
		session.RefreshTokenHash,
		session.RefreshTokenExpiresAt,
		session.TokenFamily,
		session.TokenGeneration,
		session.CreatedAt,
		session.UpdatedAt,
		session.DeletedAt,
	)
}

func (r *SessionRepository) GetActiveByUserID(ctx context.Context, userID uuid.UUID, searchQuery string, offset int, limit int) ([]*sessions.Session, int64, error) {
	return r.getManyWithTotal(ctx, queryGetActiveByUserID, userID, searchQuery, limit, offset)
}

func (r *SessionRepository) GetByID(ctx context.Context, id uuid.UUID) (*sessions.Session, error) {
	return r.getOne(ctx, queryGetByID, id)
}

func (r *SessionRepository) GetByIDAndUserID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*sessions.Session, error) {
	return r.getOne(ctx, queryGetByIDAndUserID, id, userID)
}

func (r *SessionRepository) RevokeByID(ctx context.Context, id uuid.UUID, updatedAt time.Time) error {
	return r.execWithRowCheck(ctx, queryRevokeByID, id, updatedAt)
}

func (r *SessionRepository) RevokeAllByUserID(ctx context.Context, userID uuid.UUID, updatedAt time.Time) ([]uuid.UUID, error) {
	return r.execReturnIDs(ctx, queryRevokeAllByUserID, userID, updatedAt)
}

func (r *SessionRepository) RevokeAllByTokenFamily(ctx context.Context, tokenFamily uuid.UUID, updatedAt time.Time) ([]uuid.UUID, error) {
	return r.execReturnIDs(ctx, queryRevokeAllByTokenFamily, tokenFamily, updatedAt)
}

func (r *SessionRepository) UpdateRefreshToken(ctx context.Context, session *sessions.Session, expectedGeneration int) error {
	return database.ExecWithRowCheck(ctx, r.Executor(ctx), queryUpdateRefreshToken, sessions.ErrRefreshTokenReused,
		session.ID,
		session.RefreshTokenHash,
		session.RefreshTokenExpiresAt,
		session.TokenGeneration,
		session.UpdatedAt,
		expectedGeneration,
	)
}
