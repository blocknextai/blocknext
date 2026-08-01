package passwordcredentials

import (
	"context"
	"database/sql"

	"github.com/blocknextai/go-packages/database"
	"github.com/blocknextai/platform-api/internal/account/domain/passwordcredentials"
	"github.com/google/uuid"
)

const (
	tableName = "account.password_credentials"
	columns   = "id, user_id, password_hash, created_at, updated_at, deleted_at"
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

	queryCreate = database.BuildQuery(`
		INSERT INTO `, tableName, ` (`, columns, `)
		VALUES ($1, $2, $3, $4, $5, $6)
	`)

	queryUpdate = database.BuildQuery(`
		UPDATE `, tableName, `
		SET
			password_hash = $2,
			updated_at = $3
		WHERE
			user_id = $1
			AND deleted_at IS NULL
	`)

	queryDelete = database.BuildQuery(`
		UPDATE `, tableName, `
		SET
			updated_at = $2,
			deleted_at = $3
		WHERE
			user_id = $1
			AND deleted_at IS NULL
	`)
)

type PasswordCredentialRepository struct {
	database.BaseRepository
}

func NewPasswordCredentialRepository(db *sql.DB) passwordcredentials.PasswordCredentialRepository {
	return &PasswordCredentialRepository{
		BaseRepository: database.NewBaseRepository(db),
	}
}

func (r *PasswordCredentialRepository) scan(row interface{ Scan(dest ...any) error }, _ ...any) (*passwordcredentials.PasswordCredential, error) {
	var credential passwordcredentials.PasswordCredential
	var deletedAt sql.NullTime

	err := row.Scan(
		&credential.ID,
		&credential.UserID,
		&credential.PasswordHash,
		&credential.CreatedAt,
		&credential.UpdatedAt,
		&deletedAt,
	)

	if err != nil {
		return nil, err
	}

	credential.DeletedAt = database.ScanNullTime(deletedAt)

	return &credential, nil
}

func (r *PasswordCredentialRepository) getOne(ctx context.Context, query string, args ...any) (*passwordcredentials.PasswordCredential, error) {
	return database.GetOne(ctx, r.Executor(ctx), query, r.scan, passwordcredentials.ErrPasswordCredentialNotFound, args...)
}

func (r *PasswordCredentialRepository) exec(ctx context.Context, query string, args ...any) error {
	return database.Exec(ctx, r.Executor(ctx), query, args...)
}

func (r *PasswordCredentialRepository) execWithRowCheck(ctx context.Context, query string, args ...any) error {
	return database.ExecWithRowCheck(ctx, r.Executor(ctx), query, passwordcredentials.ErrPasswordCredentialNotFound, args...)
}

func (r *PasswordCredentialRepository) GetByUserID(ctx context.Context, userID uuid.UUID) (*passwordcredentials.PasswordCredential, error) {
	return r.getOne(ctx, queryGetByUserID, userID)
}

func (r *PasswordCredentialRepository) Create(ctx context.Context, credential *passwordcredentials.PasswordCredential) error {
	return r.exec(ctx, queryCreate,
		credential.ID,
		credential.UserID,
		credential.PasswordHash,
		credential.CreatedAt,
		credential.UpdatedAt,
		credential.DeletedAt,
	)
}

func (r *PasswordCredentialRepository) Update(ctx context.Context, credential *passwordcredentials.PasswordCredential) error {
	return r.execWithRowCheck(ctx, queryUpdate,
		credential.UserID,
		credential.PasswordHash,
		credential.UpdatedAt,
	)
}

func (r *PasswordCredentialRepository) Delete(ctx context.Context, credential *passwordcredentials.PasswordCredential) error {
	return r.execWithRowCheck(ctx, queryDelete,
		credential.UserID,
		credential.UpdatedAt,
		credential.DeletedAt,
	)
}
