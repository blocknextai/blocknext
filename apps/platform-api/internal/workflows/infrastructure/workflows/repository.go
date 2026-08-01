package workflows

import (
	"context"
	"database/sql"

	"github.com/blocknextai/go-packages/database"
	"github.com/blocknextai/go-packages/json"
	"github.com/blocknextai/platform-api/internal/workflows/domain/workflows"
	"github.com/google/uuid"
)

const (
	tableName = "workflows.workflows"
	columns   = "id, organization_id, owner_id, title, description, is_pinned, is_active, nodes, edges, created_at, updated_at, deleted_at"
)

var (
	queryGetAll = database.BuildQuery(`
		SELECT `, columns, `, COUNT(*) OVER() AS total
		FROM `, tableName, `
		WHERE
			organization_id = $1
			AND deleted_at IS NULL
			AND ($2 = '' OR title ILIKE '%' || $2 || '%' OR COALESCE(description, '') ILIKE '%' || $2 || '%')
		ORDER BY
			is_pinned DESC,
			updated_at DESC
		LIMIT $3 OFFSET $4
	`)

	queryGetByOrganizationIDAndID = database.BuildQuery(`
		SELECT `, columns, `
		FROM `, tableName, `
		WHERE
			organization_id = $1
			AND id = $2
			AND deleted_at IS NULL
		LIMIT 1
	`)

	queryGetCountByOrganizationID = database.BuildQuery(`
		SELECT COUNT(*)
		FROM `, tableName, `
		WHERE
			organization_id = $1
			AND deleted_at IS NULL
	`)

	queryCreate = database.BuildQuery(`
		INSERT INTO `, tableName, ` (`, columns, `)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	`)

	queryUpdate = database.BuildQuery(`
		UPDATE `, tableName, `
		SET
			title = $2,
			description = $3,
			is_pinned = $4,
			is_active = $5,
			nodes = $6,
			edges = $7,
			updated_at = $8
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

type WorkflowRepository struct {
	database.BaseRepository
}

func NewWorkflowRepository(db *sql.DB) workflows.WorkflowRepository {
	return &WorkflowRepository{
		BaseRepository: database.NewBaseRepository(db),
	}
}

func (r *WorkflowRepository) scan(row interface{ Scan(dest ...any) error }, extras ...any) (*workflows.Workflow, error) {
	var w workflows.Workflow
	var deletedAt sql.NullTime
	var description sql.NullString
	var nodes, edges []byte

	dest := []any{
		&w.ID,
		&w.OrganizationID,
		&w.OwnerID,
		&w.Title,
		&description,
		&w.IsPinned,
		&w.IsActive,
		&nodes,
		&edges,
		&w.CreatedAt,
		&w.UpdatedAt,
		&deletedAt,
	}
	dest = append(dest, extras...)

	if err := row.Scan(dest...); err != nil {
		return nil, err
	}

	if err := database.UnmarshalJSONColumn(nodes, &w.Nodes); err != nil {
		return nil, err
	}
	if err := database.UnmarshalJSONColumn(edges, &w.Edges); err != nil {
		return nil, err
	}

	w.DeletedAt = database.ScanNullTime(deletedAt)
	w.Description = database.ScanNullString(description)

	return &w, nil
}

func (r *WorkflowRepository) getOne(ctx context.Context, query string, args ...any) (*workflows.Workflow, error) {
	return database.GetOne(ctx, r.Executor(ctx), query, r.scan, workflows.ErrWorkflowNotFound, args...)
}

func (r *WorkflowRepository) getManyWithTotal(ctx context.Context, query string, args ...any) ([]*workflows.Workflow, int64, error) {
	return database.GetManyWithTotal(ctx, r.Executor(ctx), query, r.scan, args...)
}

func (r *WorkflowRepository) exec(ctx context.Context, query string, args ...any) error {
	return database.Exec(ctx, r.Executor(ctx), query, args...)
}

func (r *WorkflowRepository) execWithRowCheck(ctx context.Context, query string, args ...any) error {
	return database.ExecWithRowCheck(ctx, r.Executor(ctx), query, workflows.ErrWorkflowNotFound, args...)
}

func (r *WorkflowRepository) GetAllByOrganizationID(ctx context.Context, organizationID uuid.UUID, searchQuery string, offset int, limit int) ([]*workflows.Workflow, int64, error) {
	return r.getManyWithTotal(ctx, queryGetAll, organizationID, searchQuery, limit, offset)
}

func (r *WorkflowRepository) GetByOrganizationIDAndID(ctx context.Context, organizationID uuid.UUID, workflowID uuid.UUID) (*workflows.Workflow, error) {
	return r.getOne(ctx, queryGetByOrganizationIDAndID, organizationID, workflowID)
}

func (r *WorkflowRepository) GetCountByOrganizationID(ctx context.Context, organizationID uuid.UUID) (int64, error) {
	var count int64
	err := r.Executor(ctx).QueryRowContext(ctx, queryGetCountByOrganizationID, organizationID).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *WorkflowRepository) Create(ctx context.Context, workflow *workflows.Workflow) error {
	nodesJSON, err := json.Marshal(workflow.Nodes)
	if err != nil {
		return err
	}
	edgesJSON, err := json.Marshal(workflow.Edges)
	if err != nil {
		return err
	}

	return r.exec(ctx, queryCreate,
		workflow.ID,
		workflow.OrganizationID,
		workflow.OwnerID,
		workflow.Title,
		workflow.Description,
		workflow.IsPinned,
		workflow.IsActive,
		nodesJSON,
		edgesJSON,
		workflow.CreatedAt,
		workflow.UpdatedAt,
		workflow.DeletedAt,
	)
}

func (r *WorkflowRepository) Update(ctx context.Context, workflow *workflows.Workflow) error {
	nodesJSON, err := json.Marshal(workflow.Nodes)
	if err != nil {
		return err
	}
	edgesJSON, err := json.Marshal(workflow.Edges)
	if err != nil {
		return err
	}

	return r.execWithRowCheck(ctx, queryUpdate,
		workflow.ID,
		workflow.Title,
		workflow.Description,
		workflow.IsPinned,
		workflow.IsActive,
		nodesJSON,
		edgesJSON,
		workflow.UpdatedAt,
	)
}

func (r *WorkflowRepository) Delete(ctx context.Context, workflow *workflows.Workflow) error {
	return r.execWithRowCheck(ctx, queryDelete,
		workflow.ID,
		workflow.UpdatedAt,
		workflow.DeletedAt,
	)
}
