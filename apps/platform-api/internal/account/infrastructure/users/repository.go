package users

import (
	"context"
	"database/sql"

	"github.com/blocknextai/go-packages/database"
	"github.com/blocknextai/platform-api/internal/account/domain/users"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

const (
	tableName = "account.users"
	columns   = "id, role, is_verified, is_banned, created_at, updated_at, deleted_at"
)

var (
	queryGetByID = database.BuildQuery(`
		SELECT `, columns, `
		FROM `, tableName, `
		WHERE
			id = $1
			AND deleted_at IS NULL
		LIMIT 1
	`)

	queryCreate = database.BuildQuery(`
		INSERT INTO `, tableName, ` (`, columns, `)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`)

	queryGetAllByIDs = database.BuildQuery(`
		SELECT `, columns, `
		FROM `, tableName, `
		WHERE
			id = ANY($1)
			AND deleted_at IS NULL
	`)
)

type UserRepository struct {
	database.BaseRepository
}

func NewUserRepository(db *sql.DB) users.UserRepository {
	return &UserRepository{
		BaseRepository: database.NewBaseRepository(db),
	}
}

func (r *UserRepository) scan(row interface{ Scan(dest ...any) error }, _ ...any) (*users.User, error) {
	var user users.User
	var deletedAt sql.NullTime

	err := row.Scan(
		&user.ID,
		&user.Role,
		&user.IsVerified,
		&user.IsBanned,
		&user.CreatedAt,
		&user.UpdatedAt,
		&deletedAt,
	)

	if err != nil {
		return nil, err
	}

	user.DeletedAt = database.ScanNullTime(deletedAt)

	return &user, nil
}

func (r *UserRepository) getOne(ctx context.Context, query string, args ...any) (*users.User, error) {
	return database.GetOne(ctx, r.Executor(ctx), query, r.scan, users.ErrUserNotFound, args...)
}

func (r *UserRepository) getMany(ctx context.Context, query string, args ...any) ([]*users.User, error) {
	return database.GetMany(ctx, r.Executor(ctx), query, r.scan, args...)
}

func (r *UserRepository) exec(ctx context.Context, query string, args ...any) error {
	return database.Exec(ctx, r.Executor(ctx), query, args...)
}

func (r *UserRepository) GetByID(ctx context.Context, id uuid.UUID) (*users.User, error) {
	return r.getOne(ctx, queryGetByID, id)
}

func (r *UserRepository) Create(ctx context.Context, user *users.User) error {
	return r.exec(ctx, queryCreate,
		user.ID,
		user.Role,
		user.IsVerified,
		user.IsBanned,
		user.CreatedAt,
		user.UpdatedAt,
		user.DeletedAt,
	)
}

func (r *UserRepository) GetAllByIDs(ctx context.Context, ids []uuid.UUID) ([]*users.User, error) {
	if len(ids) == 0 {
		return []*users.User{}, nil
	}
	return r.getMany(ctx, queryGetAllByIDs, pq.Array(ids))
}
