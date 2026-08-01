package sessions

import (
	"context"
	"database/sql"

	"github.com/blocknextai/go-packages/database"
	"github.com/blocknextai/platform-api/internal/workflows/domain/generation/sessions"
	"github.com/google/uuid"
)

const (
	tableName = "workflows.generation_sessions"
	columns   = "id, organization_id, user_id, title, created_at, updated_at, deleted_at"
)

var (
	queryGetByIDAndOrganizationID = database.BuildQuery(`
		SELECT `, columns, `
		FROM `, tableName, `
		WHERE
			id = $1
			AND organization_id = $2
			AND deleted_at IS NULL
		LIMIT 1
	`)

	queryGetByOrganizationID = database.BuildQuery(`
		SELECT `, columns, `, COUNT(*) OVER() AS total
		FROM `, tableName, `
		WHERE
			organization_id = $1
			AND deleted_at IS NULL
			AND ($2 = '' OR title ILIKE '%' || $2 || '%')
		ORDER BY
			updated_at DESC
		LIMIT $3 OFFSET $4
	`)

	queryCreate = database.BuildQuery(`
		INSERT INTO `, tableName, ` (`, columns, `)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`)

	queryUpdate = database.BuildQuery(`
		UPDATE `, tableName, `
		SET
			title = $2,
			updated_at = $3
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

type SessionRepository struct {
	database.BaseRepository
}

func NewSessionRepository(db *sql.DB) sessions.SessionRepository {
	return &SessionRepository{
		BaseRepository: database.NewBaseRepository(db),
	}
}

func (r *SessionRepository) scan(row interface{ Scan(dest ...any) error }, extras ...any) (*sessions.GenerationSession, error) {
	var s sessions.GenerationSession
	var deletedAt sql.NullTime
	var organizationID uuid.NullUUID
	var userID uuid.NullUUID

	dest := []any{
		&s.ID,
		&organizationID,
		&userID,
		&s.Title,
		&s.CreatedAt,
		&s.UpdatedAt,
		&deletedAt,
	}
	dest = append(dest, extras...)

	if err := row.Scan(dest...); err != nil {
		return nil, err
	}

	s.DeletedAt = database.ScanNullTime(deletedAt)
	s.OrganizationID = database.ScanNullUUID(organizationID)
	s.UserID = database.ScanNullUUID(userID)

	return &s, nil
}

func (r *SessionRepository) getOne(ctx context.Context, query string, args ...any) (*sessions.GenerationSession, error) {
	return database.GetOne(ctx, r.Executor(ctx), query, r.scan, sessions.ErrSessionNotFound, args...)
}

func (r *SessionRepository) getManyWithTotal(ctx context.Context, query string, args ...any) ([]*sessions.GenerationSession, int64, error) {
	return database.GetManyWithTotal(ctx, r.Executor(ctx), query, r.scan, args...)
}

func (r *SessionRepository) exec(ctx context.Context, query string, args ...any) error {
	return database.Exec(ctx, r.Executor(ctx), query, args...)
}

func (r *SessionRepository) execWithRowCheck(ctx context.Context, query string, args ...any) error {
	return database.ExecWithRowCheck(ctx, r.Executor(ctx), query, sessions.ErrSessionNotFound, args...)
}

func (r *SessionRepository) GetByIDAndOrganizationID(ctx context.Context, id uuid.UUID, organizationID uuid.UUID) (*sessions.GenerationSession, error) {
	return r.getOne(ctx, queryGetByIDAndOrganizationID, id, organizationID)
}

func (r *SessionRepository) GetAllByOrganizationID(ctx context.Context, organizationID uuid.UUID, searchQuery string, offset int, limit int) ([]*sessions.GenerationSession, int64, error) {
	return r.getManyWithTotal(ctx, queryGetByOrganizationID, organizationID, searchQuery, limit, offset)
}

func (r *SessionRepository) Create(ctx context.Context, session *sessions.GenerationSession) error {
	return r.exec(ctx, queryCreate,
		session.ID,
		session.OrganizationID,
		session.UserID,
		session.Title,
		session.CreatedAt,
		session.UpdatedAt,
		session.DeletedAt,
	)
}

func (r *SessionRepository) Update(ctx context.Context, session *sessions.GenerationSession) error {
	return r.execWithRowCheck(ctx, queryUpdate,
		session.ID,
		session.Title,
		session.UpdatedAt,
	)
}

func (r *SessionRepository) Delete(ctx context.Context, session *sessions.GenerationSession) error {
	return r.execWithRowCheck(ctx, queryDelete,
		session.ID,
		session.UpdatedAt,
		session.DeletedAt,
	)
}
