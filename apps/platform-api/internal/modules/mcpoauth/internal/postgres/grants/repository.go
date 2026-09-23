package grants

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/lib/pq"

	"github.com/blocknextai/go-packages/database"
	mcpOAuthDomainGrants "github.com/blocknextai/platform-api/internal/modules/mcpoauth/internal/domain/grants"
	mcpOAuthDomainOAuth2 "github.com/blocknextai/platform-api/internal/modules/mcpoauth/internal/domain/oauth2"
)

const (
	tableName = "mcpoauth.grants"
	columns   = "id, client_id, user_id, organization_id, scopes, resource, revoked_at, created_at, updated_at, deleted_at"
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

	queryGetActiveByClientIDAndUserIDAndOrganizationIDAndResource = database.BuildQuery(`
		SELECT `, columns, `
		FROM `, tableName, `
		WHERE
			client_id = $1
			AND user_id = $2
			AND organization_id = $3
			AND resource = $4
			AND revoked_at IS NULL
			AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT 1
	`)

	queryGetAllActiveByUserID = database.BuildQuery(`
		SELECT `, columns, `, COUNT(*) OVER() AS total
		FROM `, tableName, `
		WHERE
			user_id = $1
			AND revoked_at IS NULL
			AND deleted_at IS NULL
			AND ($2 = '' OR client_id ILIKE '%' || $2 || '%')
		ORDER BY created_at DESC
		LIMIT $3 OFFSET $4
	`)

	queryCreate = database.BuildQuery(`
		INSERT INTO `, tableName, ` (`, columns, `)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`)

	queryUpdate = database.BuildQuery(`
		UPDATE `, tableName, `
		SET
			scopes = $2,
			revoked_at = $3,
			updated_at = $4
		WHERE
			id = $1
			AND deleted_at IS NULL
	`)
)

type GrantRepository struct {
	database.BaseRepository
}

func NewGrantRepository(db *sql.DB) mcpOAuthDomainGrants.GrantRepository {
	return &GrantRepository{
		BaseRepository: database.NewBaseRepository(db),
	}
}

func (r *GrantRepository) scan(row interface{ Scan(dest ...any) error }, extras ...any) (*mcpOAuthDomainGrants.Grant, error) {
	var g mcpOAuthDomainGrants.Grant
	var scopes []string
	var revokedAt sql.NullTime
	var deletedAt sql.NullTime

	dest := []any{
		&g.ID,
		&g.ClientID,
		&g.UserID,
		&g.OrganizationID,
		pq.Array(&scopes),
		&g.Resource,
		&revokedAt,
		&g.CreatedAt,
		&g.UpdatedAt,
		&deletedAt,
	}
	dest = append(dest, extras...)

	if err := row.Scan(dest...); err != nil {
		return nil, err
	}

	g.Scopes = mcpOAuthDomainOAuth2.ScopesFromStrings(scopes)
	g.RevokedAt = database.ScanNullTime(revokedAt)
	g.DeletedAt = database.ScanNullTime(deletedAt)

	return &g, nil
}

func (r *GrantRepository) getOne(ctx context.Context, query string, args ...any) (*mcpOAuthDomainGrants.Grant, error) {
	return database.GetOne(ctx, r.Executor(ctx), query, r.scan, mcpOAuthDomainGrants.ErrGrantNotFound, args...)
}

func (r *GrantRepository) getManyWithTotal(ctx context.Context, query string, args ...any) ([]*mcpOAuthDomainGrants.Grant, int64, error) {
	return database.GetManyWithTotal(ctx, r.Executor(ctx), query, r.scan, args...)
}

func (r *GrantRepository) exec(ctx context.Context, query string, args ...any) error {
	return database.Exec(ctx, r.Executor(ctx), query, args...)
}

func (r *GrantRepository) execWithRowCheck(ctx context.Context, query string, args ...any) error {
	return database.ExecWithRowCheck(ctx, r.Executor(ctx), query, mcpOAuthDomainGrants.ErrGrantNotFound, args...)
}

func (r *GrantRepository) GetByID(ctx context.Context, id uuid.UUID) (*mcpOAuthDomainGrants.Grant, error) {
	return r.getOne(ctx, queryGetByID, id)
}

func (r *GrantRepository) GetActiveByClientIDAndUserIDAndOrganizationIDAndResource(
	ctx context.Context,
	clientID string,
	userID uuid.UUID,
	organizationID uuid.UUID,
	resource string,
) (*mcpOAuthDomainGrants.Grant, error) {
	return r.getOne(ctx, queryGetActiveByClientIDAndUserIDAndOrganizationIDAndResource, clientID, userID, organizationID, resource)
}

func (r *GrantRepository) GetAllActiveByUserID(ctx context.Context, userID uuid.UUID, searchQuery string, offset int, limit int) ([]*mcpOAuthDomainGrants.Grant, int64, error) {
	return r.getManyWithTotal(ctx, queryGetAllActiveByUserID, userID, searchQuery, limit, offset)
}

func (r *GrantRepository) Create(ctx context.Context, grant *mcpOAuthDomainGrants.Grant) error {
	return r.exec(ctx, queryCreate,
		grant.ID,
		grant.ClientID,
		grant.UserID,
		grant.OrganizationID,
		pq.Array(grant.Scopes.Strings()),
		grant.Resource,
		grant.RevokedAt,
		grant.CreatedAt,
		grant.UpdatedAt,
		grant.DeletedAt,
	)
}

func (r *GrantRepository) Update(ctx context.Context, grant *mcpOAuthDomainGrants.Grant) error {
	return r.execWithRowCheck(ctx, queryUpdate, grant.ID, pq.Array(grant.Scopes.Strings()), grant.RevokedAt, grant.UpdatedAt)
}
