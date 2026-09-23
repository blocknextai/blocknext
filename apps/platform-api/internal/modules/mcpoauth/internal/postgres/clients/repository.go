package clients

import (
	"context"
	"database/sql"

	"github.com/lib/pq"

	"github.com/blocknextai/go-packages/database"
	mcpOAuthDomainClients "github.com/blocknextai/platform-api/internal/modules/mcpoauth/internal/domain/clients"
	mcpOAuthDomainOAuth2 "github.com/blocknextai/platform-api/internal/modules/mcpoauth/internal/domain/oauth2"
)

const (
	tableName = "mcpoauth.clients"
	columns   = "id, client_id, client_secret_hash, name, redirect_uris, grant_types, response_types, scopes, token_endpoint_auth_method, logo_uri, client_uri, created_at, updated_at, deleted_at"
)

var (
	queryGetByClientID = database.BuildQuery(`
		SELECT `, columns, `
		FROM `, tableName, `
		WHERE
			client_id = $1
			AND deleted_at IS NULL
		LIMIT 1
	`)

	queryCreate = database.BuildQuery(`
		INSERT INTO `, tableName, ` (`, columns, `)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
	`)
)

type ClientRepository struct {
	database.BaseRepository
}

func NewClientRepository(db *sql.DB) mcpOAuthDomainClients.ClientRepository {
	return &ClientRepository{
		BaseRepository: database.NewBaseRepository(db),
	}
}

func grantTypesToStrings(grantTypes []mcpOAuthDomainOAuth2.GrantType) []string {
	values := make([]string, 0, len(grantTypes))
	for _, grantType := range grantTypes {
		values = append(values, grantType.String())
	}
	return values
}

func stringsToGrantTypes(values []string) []mcpOAuthDomainOAuth2.GrantType {
	grantTypes := make([]mcpOAuthDomainOAuth2.GrantType, 0, len(values))
	for _, value := range values {
		grantTypes = append(grantTypes, mcpOAuthDomainOAuth2.GrantType(value))
	}
	return grantTypes
}

func responseTypesToStrings(responseTypes []mcpOAuthDomainOAuth2.ResponseType) []string {
	values := make([]string, 0, len(responseTypes))
	for _, responseType := range responseTypes {
		values = append(values, responseType.String())
	}
	return values
}

func stringsToResponseTypes(values []string) []mcpOAuthDomainOAuth2.ResponseType {
	responseTypes := make([]mcpOAuthDomainOAuth2.ResponseType, 0, len(values))
	for _, value := range values {
		responseTypes = append(responseTypes, mcpOAuthDomainOAuth2.ResponseType(value))
	}
	return responseTypes
}

func (r *ClientRepository) scan(row interface{ Scan(dest ...any) error }, extras ...any) (*mcpOAuthDomainClients.Client, error) {
	var c mcpOAuthDomainClients.Client
	var clientSecretHash sql.NullString
	var grantTypes []string
	var responseTypes []string
	var scopes []string
	var logoURI sql.NullString
	var clientURI sql.NullString
	var deletedAt sql.NullTime

	dest := []any{
		&c.ID,
		&c.ClientID,
		&clientSecretHash,
		&c.Name,
		pq.Array(&c.RedirectURIs),
		pq.Array(&grantTypes),
		pq.Array(&responseTypes),
		pq.Array(&scopes),
		&c.TokenEndpointAuthMethod,
		&logoURI,
		&clientURI,
		&c.CreatedAt,
		&c.UpdatedAt,
		&deletedAt,
	}
	dest = append(dest, extras...)

	if err := row.Scan(dest...); err != nil {
		return nil, err
	}

	c.ClientSecretHash = database.ScanNullString(clientSecretHash)
	c.GrantTypes = stringsToGrantTypes(grantTypes)
	c.ResponseTypes = stringsToResponseTypes(responseTypes)
	c.Scopes = mcpOAuthDomainOAuth2.ScopesFromStrings(scopes)
	c.LogoURI = database.ScanNullString(logoURI)
	c.ClientURI = database.ScanNullString(clientURI)
	c.DeletedAt = database.ScanNullTime(deletedAt)

	return &c, nil
}

func (r *ClientRepository) getOne(ctx context.Context, query string, args ...any) (*mcpOAuthDomainClients.Client, error) {
	return database.GetOne(ctx, r.Executor(ctx), query, r.scan, mcpOAuthDomainClients.ErrClientNotFound, args...)
}

func (r *ClientRepository) exec(ctx context.Context, query string, args ...any) error {
	return database.Exec(ctx, r.Executor(ctx), query, args...)
}

func (r *ClientRepository) GetByClientID(ctx context.Context, clientID string) (*mcpOAuthDomainClients.Client, error) {
	return r.getOne(ctx, queryGetByClientID, clientID)
}

func (r *ClientRepository) Create(ctx context.Context, client *mcpOAuthDomainClients.Client) error {
	err := r.exec(ctx, queryCreate,
		client.ID,
		client.ClientID,
		client.ClientSecretHash,
		client.Name,
		pq.Array(client.RedirectURIs),
		pq.Array(grantTypesToStrings(client.GrantTypes)),
		pq.Array(responseTypesToStrings(client.ResponseTypes)),
		pq.Array(client.Scopes.Strings()),
		client.TokenEndpointAuthMethod,
		client.LogoURI,
		client.ClientURI,
		client.CreatedAt,
		client.UpdatedAt,
		client.DeletedAt,
	)
	if err != nil && database.IsUniqueViolationOn(err, "idx_mcpoauth_clients_client_id") {
		return mcpOAuthDomainClients.ErrClientAlreadyExists
	}
	return err
}
