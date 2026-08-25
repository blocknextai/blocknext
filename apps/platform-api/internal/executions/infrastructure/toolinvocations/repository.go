package toolinvocations

import (
	"context"
	"database/sql"

	"github.com/blocknextai/go-packages/database"
	"github.com/blocknextai/go-packages/json"
	"github.com/blocknextai/platform-api/internal/executions/domain/toolinvocations"
	"github.com/google/uuid"
)

const (
	tableName = "executions.tool_invocations"
	columns   = "id, organization_id, api_key_id, source, tool_id, status, parameters, credentials, outputs, error_message, started_at, completed_at, created_at, updated_at, deleted_at"
)

var (
	queryGetAllByOrganizationID = database.BuildQuery(`
		SELECT `, columns, `, COUNT(*) OVER() AS total
		FROM `, tableName, `
		WHERE
			organization_id = $1
			AND deleted_at IS NULL
			AND ($2 = '' OR tool_id ILIKE '%' || $2 || '%'  OR status ILIKE '%' || $2 || '%')
		ORDER BY
			started_at DESC
		LIMIT $3 OFFSET $4
	`)

	queryGetByIDAndOrganizationID = database.BuildQuery(`
		SELECT `, columns, `
		FROM `, tableName, `
		WHERE
			id = $1
			AND organization_id = $2
			AND deleted_at IS NULL
		LIMIT 1
	`)

	queryCreate = database.BuildQuery(`
		INSERT INTO `, tableName, ` (`, columns, `)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
	`)
)

type ToolInvocationRepository struct {
	database.BaseRepository
}

func NewToolInvocationRepository(db *sql.DB) toolinvocations.ToolInvocationRepository {
	return &ToolInvocationRepository{
		BaseRepository: database.NewBaseRepository(db),
	}
}

func (r *ToolInvocationRepository) scan(row interface{ Scan(dest ...any) error }, extra ...any) (*toolinvocations.ToolInvocation, error) {
	var ti toolinvocations.ToolInvocation
	var deletedAt sql.NullTime
	var apiKeyID uuid.NullUUID
	var errorMessage sql.NullString
	var parameters, credentials, outputs []byte

	dest := []any{
		&ti.ID,
		&ti.OrganizationID,
		&apiKeyID,
		&ti.Source,
		&ti.ToolID,
		&ti.Status,
		&parameters,
		&credentials,
		&outputs,
		&errorMessage,
		&ti.StartedAt,
		&ti.CompletedAt,
		&ti.CreatedAt,
		&ti.UpdatedAt,
		&deletedAt,
	}
	dest = append(dest, extra...)

	if err := row.Scan(dest...); err != nil {
		return nil, err
	}

	if err := database.UnmarshalJSONColumn(parameters, &ti.Parameters); err != nil {
		return nil, err
	}
	if err := database.UnmarshalJSONColumn(credentials, &ti.Credentials); err != nil {
		return nil, err
	}
	if err := database.UnmarshalJSONColumn(outputs, &ti.Outputs); err != nil {
		return nil, err
	}

	ti.APIKeyID = database.ScanNullUUID(apiKeyID)
	ti.ErrorMessage = database.ScanNullString(errorMessage)
	ti.DeletedAt = database.ScanNullTime(deletedAt)

	return &ti, nil
}

func (r *ToolInvocationRepository) getOne(ctx context.Context, query string, args ...any) (*toolinvocations.ToolInvocation, error) {
	return database.GetOne(ctx, r.Executor(ctx), query, r.scan, toolinvocations.ErrToolInvocationNotFound, args...)
}

func (r *ToolInvocationRepository) getManyWithTotal(ctx context.Context, query string, args ...any) ([]*toolinvocations.ToolInvocation, int64, error) {
	return database.GetManyWithTotal(ctx, r.Executor(ctx), query, r.scan, args...)
}

func (r *ToolInvocationRepository) exec(ctx context.Context, query string, args ...any) error {
	return database.Exec(ctx, r.Executor(ctx), query, args...)
}

func (r *ToolInvocationRepository) GetAllByOrganizationID(
	ctx context.Context,
	organizationID uuid.UUID,
	searchQuery string,
	offset int,
	limit int,
) ([]*toolinvocations.ToolInvocation, int64, error) {
	return r.getManyWithTotal(ctx, queryGetAllByOrganizationID, organizationID, searchQuery, limit, offset)
}

func (r *ToolInvocationRepository) GetByIDAndOrganizationID(ctx context.Context, id uuid.UUID, organizationID uuid.UUID) (*toolinvocations.ToolInvocation, error) {
	return r.getOne(ctx, queryGetByIDAndOrganizationID, id, organizationID)
}

func (r *ToolInvocationRepository) Create(ctx context.Context, ti *toolinvocations.ToolInvocation) error {
	parametersJSON, err := json.Marshal(ti.Parameters)
	if err != nil {
		return err
	}
	credentialsJSON, err := json.Marshal(ti.Credentials)
	if err != nil {
		return err
	}
	outputsJSON, err := json.Marshal(ti.Outputs)
	if err != nil {
		return err
	}

	return r.exec(ctx, queryCreate,
		ti.ID,
		ti.OrganizationID,
		ti.APIKeyID,
		ti.Source,
		ti.ToolID,
		ti.Status,
		parametersJSON,
		credentialsJSON,
		outputsJSON,
		ti.ErrorMessage,
		ti.StartedAt,
		ti.CompletedAt,
		ti.CreatedAt,
		ti.UpdatedAt,
		ti.DeletedAt,
	)
}
