package credentials

import (
	"context"
	"database/sql"

	"github.com/blocknextai/go-packages/database"
	commonDomain "github.com/blocknextai/platform-api/internal/common/domain"
	"github.com/blocknextai/platform-api/internal/credentials/domain/credentials"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

const (
	tableName = "credentials.credentials"
	columns   = "id, owner_type, owner_id, source_type, key, title, data, created_at, updated_at, deleted_at"
)

var (
	ownerFilter = `
		owner_type = $1
		AND owner_id = $2
		AND deleted_at IS NULL
	`

	queryGetAllByOwner = database.BuildQuery(`
		SELECT `, columns, `, COUNT(*) OVER() AS total
		FROM `, tableName, `
		WHERE`, ownerFilter, `
			AND ($3 = '' OR title ILIKE '%' || $3 || '%')
		ORDER BY
			updated_at DESC
		LIMIT $4 OFFSET $5
	`)

	queryGetAllByOwnerAndKeys = database.BuildQuery(`
		SELECT `, columns, `
		FROM `, tableName, `
		WHERE
			owner_type = $1
			AND owner_id = $2
			AND key = ANY($3)
			AND deleted_at IS NULL
		ORDER BY
			updated_at DESC
	`)

	queryGetByIDAndOwner = database.BuildQuery(`
		SELECT `, columns, `
		FROM `, tableName, `
		WHERE
			id = $1
			AND owner_type = $2
			AND owner_id = $3
			AND deleted_at IS NULL
		LIMIT 1
	`)

	queryCreate = database.BuildQuery(`
		INSERT INTO `, tableName, ` (`, columns, `)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`)

	queryUpdate = database.BuildQuery(`
		UPDATE `, tableName, `
		SET
			key = $2,
			title = $3,
			data = $4,
			updated_at = $5
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

type CredentialRepository struct {
	database.BaseRepository
}

func NewCredentialRepository(db *sql.DB) credentials.CredentialRepository {
	return &CredentialRepository{
		BaseRepository: database.NewBaseRepository(db),
	}
}

func (r *CredentialRepository) scan(row interface{ Scan(dest ...any) error }, extras ...any) (*credentials.Credential, error) {
	var cred credentials.Credential
	var deletedAt sql.NullTime

	dest := []any{
		&cred.ID,
		&cred.OwnerType,
		&cred.OwnerID,
		&cred.SourceType,
		&cred.Key,
		&cred.Title,
		&cred.Data,
		&cred.CreatedAt,
		&cred.UpdatedAt,
		&deletedAt,
	}
	dest = append(dest, extras...)

	if err := row.Scan(dest...); err != nil {
		return nil, err
	}

	cred.DeletedAt = database.ScanNullTime(deletedAt)

	return &cred, nil
}

func (r *CredentialRepository) getOne(ctx context.Context, query string, args ...any) (*credentials.Credential, error) {
	return database.GetOne(ctx, r.Executor(ctx), query, r.scan, credentials.ErrCredentialNotFound, args...)
}

func (r *CredentialRepository) getMany(ctx context.Context, query string, args ...any) ([]*credentials.Credential, error) {
	return database.GetMany(ctx, r.Executor(ctx), query, r.scan, args...)
}

func (r *CredentialRepository) getManyWithTotal(ctx context.Context, query string, args ...any) ([]*credentials.Credential, int64, error) {
	return database.GetManyWithTotal(ctx, r.Executor(ctx), query, r.scan, args...)
}

func (r *CredentialRepository) exec(ctx context.Context, query string, args ...any) error {
	return database.Exec(ctx, r.Executor(ctx), query, args...)
}

func (r *CredentialRepository) execWithRowCheck(ctx context.Context, query string, args ...any) error {
	return database.ExecWithRowCheck(ctx, r.Executor(ctx), query, credentials.ErrCredentialNotFound, args...)
}

func (r *CredentialRepository) GetAllByOwner(ctx context.Context, ownerType commonDomain.OwnerType, ownerID uuid.UUID, searchQuery string, offset int, limit int) ([]*credentials.Credential, int64, error) {
	return r.getManyWithTotal(ctx, queryGetAllByOwner, ownerType, ownerID, searchQuery, limit, offset)
}

func (r *CredentialRepository) GetAllByOwnerAndKeys(ctx context.Context, ownerType commonDomain.OwnerType, ownerID uuid.UUID, keys []string) ([]*credentials.Credential, error) {
	if len(keys) == 0 {
		return []*credentials.Credential{}, nil
	}
	return r.getMany(ctx, queryGetAllByOwnerAndKeys, ownerType, ownerID, pq.Array(keys))
}

func (r *CredentialRepository) GetByIDAndOwner(ctx context.Context, id uuid.UUID, ownerType commonDomain.OwnerType, ownerID uuid.UUID) (*credentials.Credential, error) {
	return r.getOne(ctx, queryGetByIDAndOwner, id, ownerType, ownerID)
}

func (r *CredentialRepository) Create(ctx context.Context, credential *credentials.Credential) error {
	return r.exec(ctx, queryCreate,
		credential.ID,
		credential.OwnerType,
		credential.OwnerID,
		credential.SourceType,
		credential.Key,
		credential.Title,
		credential.Data,
		credential.CreatedAt,
		credential.UpdatedAt,
		credential.DeletedAt,
	)
}

func (r *CredentialRepository) Update(ctx context.Context, credential *credentials.Credential) error {
	return r.execWithRowCheck(ctx, queryUpdate,
		credential.ID,
		credential.Key,
		credential.Title,
		credential.Data,
		credential.UpdatedAt,
	)
}

func (r *CredentialRepository) Delete(ctx context.Context, credential *credentials.Credential) error {
	return r.execWithRowCheck(ctx, queryDelete,
		credential.ID,
		credential.UpdatedAt,
		credential.DeletedAt,
	)
}
