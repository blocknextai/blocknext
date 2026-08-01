package apikeys

import (
	"context"
	"database/sql"

	"github.com/blocknextai/go-packages/database"
	apiKeysDomain "github.com/blocknextai/platform-api/internal/apikeys/domain/apikeys"
	commonDomain "github.com/blocknextai/platform-api/internal/common/domain"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

const (
	tableName = "apikeys.api_keys"
	columns   = "id, owner_type, owner_id, name, key_hash, scopes, last_used_at, created_at, updated_at, deleted_at"
)

var (
	queryGetByOwnerAndID = database.BuildQuery(`
		SELECT `, columns, `
		FROM `, tableName, `
		WHERE
			owner_type = $1
			AND owner_id = $2
			AND id = $3
			AND deleted_at IS NULL
		LIMIT 1
	`)

	queryGetByKeyHash = database.BuildQuery(`
		SELECT `, columns, `
		FROM `, tableName, `
		WHERE
			key_hash = $1
			AND deleted_at IS NULL
		LIMIT 1
	`)

	queryGetAllByOwner = database.BuildQuery(`
		SELECT `, columns, `, COUNT(*) OVER() AS total
		FROM `, tableName, `
		WHERE
			owner_type = $1
			AND owner_id = $2
			AND deleted_at IS NULL
			AND ($3 = '' OR name ILIKE '%' || $3 || '%')
		ORDER BY created_at DESC
		LIMIT $4 OFFSET $5
	`)

	queryCreate = database.BuildQuery(`
		INSERT INTO `, tableName, ` (`, columns, `)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
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

	queryUpdate = database.BuildQuery(`
		UPDATE `, tableName, `
		SET
			key_hash = $2,
			last_used_at = $3,
			updated_at = $4
		WHERE
			id = $1
			AND deleted_at IS NULL
	`)
)

type APIKeyRepository struct {
	database.BaseRepository
}

func NewAPIKeyRepository(db *sql.DB) apiKeysDomain.APIKeyRepository {
	return &APIKeyRepository{
		BaseRepository: database.NewBaseRepository(db),
	}
}

func scopesToStrings(scopes apiKeysDomain.Scopes) []string {
	values := make([]string, 0, len(scopes))
	for _, scope := range scopes {
		values = append(values, string(scope))
	}
	return values
}

func stringsToScopes(values []string) apiKeysDomain.Scopes {
	scopes := make(apiKeysDomain.Scopes, 0, len(values))
	for _, value := range values {
		scopes = append(scopes, apiKeysDomain.Scope(value))
	}
	return scopes
}

func (r *APIKeyRepository) scan(row interface{ Scan(dest ...any) error }, extras ...any) (*apiKeysDomain.APIKey, error) {
	var k apiKeysDomain.APIKey
	var scopes []string
	var lastUsedAt sql.NullTime
	var deletedAt sql.NullTime

	dest := []any{
		&k.ID,
		&k.OwnerType,
		&k.OwnerID,
		&k.Name,
		&k.KeyHash,
		pq.Array(&scopes),
		&lastUsedAt,
		&k.CreatedAt,
		&k.UpdatedAt,
		&deletedAt,
	}
	dest = append(dest, extras...)

	if err := row.Scan(dest...); err != nil {
		return nil, err
	}

	k.Scopes = stringsToScopes(scopes)
	k.LastUsedAt = database.ScanNullTime(lastUsedAt)
	k.DeletedAt = database.ScanNullTime(deletedAt)

	return &k, nil
}

func (r *APIKeyRepository) getOne(ctx context.Context, query string, args ...any) (*apiKeysDomain.APIKey, error) {
	return database.GetOne(ctx, r.Executor(ctx), query, r.scan, apiKeysDomain.ErrAPIKeyNotFound, args...)
}

func (r *APIKeyRepository) getManyWithTotal(ctx context.Context, query string, args ...any) ([]*apiKeysDomain.APIKey, int64, error) {
	return database.GetManyWithTotal(ctx, r.Executor(ctx), query, r.scan, args...)
}

func (r *APIKeyRepository) exec(ctx context.Context, query string, args ...any) error {
	return database.Exec(ctx, r.Executor(ctx), query, args...)
}

func (r *APIKeyRepository) execWithRowCheck(ctx context.Context, query string, args ...any) error {
	return database.ExecWithRowCheck(ctx, r.Executor(ctx), query, apiKeysDomain.ErrAPIKeyNotFound, args...)
}

func (r *APIKeyRepository) GetByOwnerAndID(ctx context.Context, ownerType commonDomain.OwnerType, ownerID uuid.UUID, id uuid.UUID) (*apiKeysDomain.APIKey, error) {
	return r.getOne(ctx, queryGetByOwnerAndID, ownerType, ownerID, id)
}

func (r *APIKeyRepository) GetByKeyHash(ctx context.Context, keyHash string) (*apiKeysDomain.APIKey, error) {
	return r.getOne(ctx, queryGetByKeyHash, keyHash)
}

func (r *APIKeyRepository) GetAllByOwner(ctx context.Context, ownerType commonDomain.OwnerType, ownerID uuid.UUID, searchQuery string, offset int, limit int) ([]*apiKeysDomain.APIKey, int64, error) {
	return r.getManyWithTotal(ctx, queryGetAllByOwner, ownerType, ownerID, searchQuery, limit, offset)
}

func (r *APIKeyRepository) Create(ctx context.Context, apiKey *apiKeysDomain.APIKey) error {
	scopes := scopesToStrings(apiKey.Scopes)

	err := r.exec(ctx, queryCreate,
		apiKey.ID,
		apiKey.OwnerType,
		apiKey.OwnerID,
		apiKey.Name,
		apiKey.KeyHash,
		pq.Array(scopes),
		apiKey.LastUsedAt,
		apiKey.CreatedAt,
		apiKey.UpdatedAt,
		apiKey.DeletedAt,
	)
	if err != nil && database.IsUniqueViolationOn(err, "idx_api_keys_key_hash") {
		return apiKeysDomain.ErrAPIKeyAlreadyExists
	}
	return err
}

func (r *APIKeyRepository) Delete(ctx context.Context, apiKey *apiKeysDomain.APIKey) error {
	return r.execWithRowCheck(ctx, queryDelete, apiKey.ID, apiKey.UpdatedAt, apiKey.DeletedAt)
}

func (r *APIKeyRepository) Update(ctx context.Context, apiKey *apiKeysDomain.APIKey) error {
	err := r.execWithRowCheck(ctx, queryUpdate, apiKey.ID, apiKey.KeyHash, apiKey.LastUsedAt, apiKey.UpdatedAt)
	if err != nil && database.IsUniqueViolationOn(err, "idx_api_keys_key_hash") {
		return apiKeysDomain.ErrAPIKeyAlreadyExists
	}
	return err
}
