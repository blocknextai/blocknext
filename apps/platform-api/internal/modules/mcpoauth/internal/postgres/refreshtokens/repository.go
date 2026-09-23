package refreshtokens

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"

	"github.com/blocknextai/go-packages/database"
	mcpOAuthDomainOAuth2 "github.com/blocknextai/platform-api/internal/modules/mcpoauth/internal/domain/oauth2"
	mcpOAuthDomainRefreshTokens "github.com/blocknextai/platform-api/internal/modules/mcpoauth/internal/domain/refreshtokens"
)

const (
	tableName = "mcpoauth.refresh_tokens"
	columns   = "id, grant_id, client_id, token_hash, scopes, resource, expires_at, used_at, revoked_at, created_at, updated_at, deleted_at"
)

var (
	queryGetByTokenHash = database.BuildQuery(`
		SELECT `, columns, `
		FROM `, tableName, `
		WHERE
			token_hash = $1
			AND deleted_at IS NULL
		LIMIT 1
	`)

	queryCreate = database.BuildQuery(`
		INSERT INTO `, tableName, ` (`, columns, `)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	`)

	queryUpdate = database.BuildQuery(`
		UPDATE `, tableName, `
		SET
			used_at = $2,
			revoked_at = $3,
			updated_at = $4
		WHERE
			id = $1
			AND used_at IS NULL
			AND deleted_at IS NULL
	`)

	queryRevokeAllByGrantID = database.BuildQuery(`
		UPDATE `, tableName, `
		SET
			revoked_at = $2,
			updated_at = $2
		WHERE
			grant_id = $1
			AND revoked_at IS NULL
			AND deleted_at IS NULL
	`)
)

type RefreshTokenRepository struct {
	database.BaseRepository
}

func NewRefreshTokenRepository(db *sql.DB) mcpOAuthDomainRefreshTokens.RefreshTokenRepository {
	return &RefreshTokenRepository{
		BaseRepository: database.NewBaseRepository(db),
	}
}

func (r *RefreshTokenRepository) scan(row interface{ Scan(dest ...any) error }, extras ...any) (*mcpOAuthDomainRefreshTokens.RefreshToken, error) {
	var rt mcpOAuthDomainRefreshTokens.RefreshToken
	var scopes []string
	var usedAt sql.NullTime
	var revokedAt sql.NullTime
	var deletedAt sql.NullTime

	dest := []any{
		&rt.ID,
		&rt.GrantID,
		&rt.ClientID,
		&rt.TokenHash,
		pq.Array(&scopes),
		&rt.Resource,
		&rt.ExpiresAt,
		&usedAt,
		&revokedAt,
		&rt.CreatedAt,
		&rt.UpdatedAt,
		&deletedAt,
	}
	dest = append(dest, extras...)

	if err := row.Scan(dest...); err != nil {
		return nil, err
	}

	rt.Scopes = mcpOAuthDomainOAuth2.ScopesFromStrings(scopes)
	rt.UsedAt = database.ScanNullTime(usedAt)
	rt.RevokedAt = database.ScanNullTime(revokedAt)
	rt.DeletedAt = database.ScanNullTime(deletedAt)

	return &rt, nil
}

func (r *RefreshTokenRepository) getOne(ctx context.Context, query string, args ...any) (*mcpOAuthDomainRefreshTokens.RefreshToken, error) {
	return database.GetOne(ctx, r.Executor(ctx), query, r.scan, mcpOAuthDomainRefreshTokens.ErrRefreshTokenNotFound, args...)
}

func (r *RefreshTokenRepository) exec(ctx context.Context, query string, args ...any) error {
	return database.Exec(ctx, r.Executor(ctx), query, args...)
}

func (r *RefreshTokenRepository) execWithRowCheck(ctx context.Context, query string, args ...any) error {
	return database.ExecWithRowCheck(ctx, r.Executor(ctx), query, mcpOAuthDomainRefreshTokens.ErrRefreshTokenUsed, args...)
}

func (r *RefreshTokenRepository) GetByTokenHash(ctx context.Context, tokenHash string) (*mcpOAuthDomainRefreshTokens.RefreshToken, error) {
	return r.getOne(ctx, queryGetByTokenHash, tokenHash)
}

func (r *RefreshTokenRepository) Create(ctx context.Context, token *mcpOAuthDomainRefreshTokens.RefreshToken) error {
	return r.exec(ctx, queryCreate,
		token.ID,
		token.GrantID,
		token.ClientID,
		token.TokenHash,
		pq.Array(token.Scopes.Strings()),
		token.Resource,
		token.ExpiresAt,
		token.UsedAt,
		token.RevokedAt,
		token.CreatedAt,
		token.UpdatedAt,
		token.DeletedAt,
	)
}

func (r *RefreshTokenRepository) Update(ctx context.Context, token *mcpOAuthDomainRefreshTokens.RefreshToken) error {
	return r.execWithRowCheck(ctx, queryUpdate, token.ID, token.UsedAt, token.RevokedAt, token.UpdatedAt)
}

func (r *RefreshTokenRepository) RevokeAllByGrantID(ctx context.Context, grantID uuid.UUID, revokedAt time.Time) error {
	return r.exec(ctx, queryRevokeAllByGrantID, grantID, revokedAt)
}
