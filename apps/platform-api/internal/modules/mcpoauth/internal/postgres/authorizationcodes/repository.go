package authorizationcodes

import (
	"context"
	"database/sql"

	"github.com/lib/pq"

	"github.com/blocknextai/go-packages/database"
	mcpOAuthDomainAuthorizationCodes "github.com/blocknextai/platform-api/internal/modules/mcpoauth/internal/domain/authorizationcodes"
	mcpOAuthDomainOAuth2 "github.com/blocknextai/platform-api/internal/modules/mcpoauth/internal/domain/oauth2"
)

const (
	tableName = "mcpoauth.authorization_codes"
	columns   = "id, grant_id, client_id, code_hash, redirect_uri, code_challenge, code_challenge_method, resource, scopes, expires_at, used_at, created_at, updated_at, deleted_at"
)

var (
	queryGetByCodeHash = database.BuildQuery(`
		SELECT `, columns, `
		FROM `, tableName, `
		WHERE
			code_hash = $1
			AND deleted_at IS NULL
		LIMIT 1
	`)

	queryCreate = database.BuildQuery(`
		INSERT INTO `, tableName, ` (`, columns, `)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
	`)

	queryUpdate = database.BuildQuery(`
		UPDATE `, tableName, `
		SET
			used_at = $2,
			updated_at = $3
		WHERE
			id = $1
			AND used_at IS NULL
			AND deleted_at IS NULL
	`)
)

type AuthorizationCodeRepository struct {
	database.BaseRepository
}

func NewAuthorizationCodeRepository(db *sql.DB) mcpOAuthDomainAuthorizationCodes.AuthorizationCodeRepository {
	return &AuthorizationCodeRepository{
		BaseRepository: database.NewBaseRepository(db),
	}
}

func (r *AuthorizationCodeRepository) scan(row interface{ Scan(dest ...any) error }, extras ...any) (*mcpOAuthDomainAuthorizationCodes.AuthorizationCode, error) {
	var ac mcpOAuthDomainAuthorizationCodes.AuthorizationCode
	var scopes []string
	var usedAt sql.NullTime
	var deletedAt sql.NullTime

	dest := []any{
		&ac.ID,
		&ac.GrantID,
		&ac.ClientID,
		&ac.CodeHash,
		&ac.RedirectURI,
		&ac.CodeChallenge,
		&ac.CodeChallengeMethod,
		&ac.Resource,
		pq.Array(&scopes),
		&ac.ExpiresAt,
		&usedAt,
		&ac.CreatedAt,
		&ac.UpdatedAt,
		&deletedAt,
	}
	dest = append(dest, extras...)

	if err := row.Scan(dest...); err != nil {
		return nil, err
	}

	ac.Scopes = mcpOAuthDomainOAuth2.ScopesFromStrings(scopes)
	ac.UsedAt = database.ScanNullTime(usedAt)
	ac.DeletedAt = database.ScanNullTime(deletedAt)

	return &ac, nil
}

func (r *AuthorizationCodeRepository) getOne(ctx context.Context, query string, args ...any) (*mcpOAuthDomainAuthorizationCodes.AuthorizationCode, error) {
	return database.GetOne(ctx, r.Executor(ctx), query, r.scan, mcpOAuthDomainAuthorizationCodes.ErrAuthorizationCodeNotFound, args...)
}

func (r *AuthorizationCodeRepository) exec(ctx context.Context, query string, args ...any) error {
	return database.Exec(ctx, r.Executor(ctx), query, args...)
}

func (r *AuthorizationCodeRepository) execWithRowCheck(ctx context.Context, query string, args ...any) error {
	return database.ExecWithRowCheck(ctx, r.Executor(ctx), query, mcpOAuthDomainAuthorizationCodes.ErrAuthorizationCodeUsed, args...)
}

func (r *AuthorizationCodeRepository) GetByCodeHash(ctx context.Context, codeHash string) (*mcpOAuthDomainAuthorizationCodes.AuthorizationCode, error) {
	return r.getOne(ctx, queryGetByCodeHash, codeHash)
}

func (r *AuthorizationCodeRepository) Create(ctx context.Context, code *mcpOAuthDomainAuthorizationCodes.AuthorizationCode) error {
	return r.exec(ctx, queryCreate,
		code.ID,
		code.GrantID,
		code.ClientID,
		code.CodeHash,
		code.RedirectURI,
		code.CodeChallenge,
		code.CodeChallengeMethod,
		code.Resource,
		pq.Array(code.Scopes.Strings()),
		code.ExpiresAt,
		code.UsedAt,
		code.CreatedAt,
		code.UpdatedAt,
		code.DeletedAt,
	)
}

func (r *AuthorizationCodeRepository) Update(ctx context.Context, code *mcpOAuthDomainAuthorizationCodes.AuthorizationCode) error {
	return r.execWithRowCheck(ctx, queryUpdate, code.ID, code.UsedAt, code.UpdatedAt)
}
