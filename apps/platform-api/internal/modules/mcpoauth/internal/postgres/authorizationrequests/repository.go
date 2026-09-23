package authorizationrequests

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/lib/pq"

	"github.com/blocknextai/go-packages/database"
	mcpOAuthDomainAuthorizationRequests "github.com/blocknextai/platform-api/internal/modules/mcpoauth/internal/domain/authorizationrequests"
	mcpOAuthDomainOAuth2 "github.com/blocknextai/platform-api/internal/modules/mcpoauth/internal/domain/oauth2"
)

const (
	tableName = "mcpoauth.authorization_requests"
	columns   = "id, client_id, redirect_uri, scopes, state, code_challenge, code_challenge_method, resource, status, user_id, organization_id, expires_at, resolved_at, created_at, updated_at, deleted_at"
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
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)
	`)

	queryUpdate = database.BuildQuery(`
		UPDATE `, tableName, `
		SET
			scopes = $2,
			status = $3,
			user_id = $4,
			organization_id = $5,
			resolved_at = $6,
			updated_at = $7
		WHERE
			id = $1
			AND deleted_at IS NULL
	`)
)

type AuthorizationRequestRepository struct {
	database.BaseRepository
}

func NewAuthorizationRequestRepository(db *sql.DB) mcpOAuthDomainAuthorizationRequests.AuthorizationRequestRepository {
	return &AuthorizationRequestRepository{
		BaseRepository: database.NewBaseRepository(db),
	}
}

func (r *AuthorizationRequestRepository) scan(row interface{ Scan(dest ...any) error }, extras ...any) (*mcpOAuthDomainAuthorizationRequests.AuthorizationRequest, error) {
	var ar mcpOAuthDomainAuthorizationRequests.AuthorizationRequest
	var scopes []string
	var state sql.NullString
	var userID uuid.NullUUID
	var organizationID uuid.NullUUID
	var resolvedAt sql.NullTime
	var deletedAt sql.NullTime

	dest := []any{
		&ar.ID,
		&ar.ClientID,
		&ar.RedirectURI,
		pq.Array(&scopes),
		&state,
		&ar.CodeChallenge,
		&ar.CodeChallengeMethod,
		&ar.Resource,
		&ar.Status,
		&userID,
		&organizationID,
		&ar.ExpiresAt,
		&resolvedAt,
		&ar.CreatedAt,
		&ar.UpdatedAt,
		&deletedAt,
	}
	dest = append(dest, extras...)

	if err := row.Scan(dest...); err != nil {
		return nil, err
	}

	ar.Scopes = mcpOAuthDomainOAuth2.ScopesFromStrings(scopes)
	ar.State = database.ScanNullString(state)
	if userID.Valid {
		ar.UserID = new(userID.UUID)
	}
	if organizationID.Valid {
		ar.OrganizationID = new(organizationID.UUID)
	}
	ar.ResolvedAt = database.ScanNullTime(resolvedAt)
	ar.DeletedAt = database.ScanNullTime(deletedAt)

	return &ar, nil
}

func (r *AuthorizationRequestRepository) getOne(ctx context.Context, query string, args ...any) (*mcpOAuthDomainAuthorizationRequests.AuthorizationRequest, error) {
	return database.GetOne(ctx, r.Executor(ctx), query, r.scan, mcpOAuthDomainAuthorizationRequests.ErrAuthorizationRequestNotFound, args...)
}

func (r *AuthorizationRequestRepository) exec(ctx context.Context, query string, args ...any) error {
	return database.Exec(ctx, r.Executor(ctx), query, args...)
}

func (r *AuthorizationRequestRepository) execWithRowCheck(ctx context.Context, query string, args ...any) error {
	return database.ExecWithRowCheck(ctx, r.Executor(ctx), query, mcpOAuthDomainAuthorizationRequests.ErrAuthorizationRequestNotFound, args...)
}

func (r *AuthorizationRequestRepository) GetByID(ctx context.Context, id uuid.UUID) (*mcpOAuthDomainAuthorizationRequests.AuthorizationRequest, error) {
	return r.getOne(ctx, queryGetByID, id)
}

func (r *AuthorizationRequestRepository) Create(ctx context.Context, request *mcpOAuthDomainAuthorizationRequests.AuthorizationRequest) error {
	return r.exec(ctx, queryCreate,
		request.ID,
		request.ClientID,
		request.RedirectURI,
		pq.Array(request.Scopes.Strings()),
		request.State,
		request.CodeChallenge,
		request.CodeChallengeMethod,
		request.Resource,
		request.Status,
		request.UserID,
		request.OrganizationID,
		request.ExpiresAt,
		request.ResolvedAt,
		request.CreatedAt,
		request.UpdatedAt,
		request.DeletedAt,
	)
}

func (r *AuthorizationRequestRepository) Update(ctx context.Context, request *mcpOAuthDomainAuthorizationRequests.AuthorizationRequest) error {
	return r.execWithRowCheck(ctx, queryUpdate,
		request.ID,
		pq.Array(request.Scopes.Strings()),
		request.Status,
		request.UserID,
		request.OrganizationID,
		request.ResolvedAt,
		request.UpdatedAt,
	)
}
