package taskexecutions

import (
	"context"
	"database/sql"
	"time"

	"github.com/blocknextai/go-packages/database"
	"github.com/blocknextai/go-packages/json"
	commonDomain "github.com/blocknextai/platform-api/internal/common/domain"
	"github.com/blocknextai/platform-api/internal/executions/domain/taskexecutions"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

const (
	tableName = "executions.task_executions"
	columns   = "id, organization_id, triggered_by_user_id, flow_trigger_id, execution_context, context_item_id, status, execution_type, error_message, nodes, edges, started_at, completed_at, created_at, updated_at, deleted_at"
)

var (
	queryGetByIDAndOrganizationID = database.BuildQuery(`
		SELECT `, columns, `
		FROM `, tableName, `
		WHERE
			id = $1
			AND organization_id = $2
			AND deleted_at IS NULL
		LIMIT 1
	`)

	queryGetByID = database.BuildQuery(`
		SELECT `, columns, `
		FROM `, tableName, `
		WHERE
			id = $1
			AND deleted_at IS NULL
		LIMIT 1
	`)

	queryGetAllByStatuses = database.BuildQuery(`
		SELECT `, columns, `
		FROM `, tableName, `
		WHERE
			status = ANY($1)
			AND deleted_at IS NULL
		ORDER BY
			updated_at DESC
	`)

	queryGetAllByIDsAndOrganizationID = database.BuildQuery(`
		SELECT `, columns, `
		FROM `, tableName, `
		WHERE
			organization_id = $1
			AND id = ANY($2)
			AND deleted_at IS NULL
		ORDER BY
			updated_at DESC
	`)

	queryCreate = database.BuildQuery(`
		INSERT INTO `, tableName, ` (`, columns, `)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)
	`)

	queryUpdate = database.BuildQuery(`
		UPDATE `, tableName, `
		SET
			status = $2,
			error_message = $3,
			completed_at = $4,
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

	queryBulkDelete = database.BuildQuery(`
		UPDATE `, tableName, `
		SET
			updated_at = $3,
			deleted_at = $4
		WHERE
			organization_id = $1
			AND id = ANY($2)
			AND deleted_at IS NULL
	`)

	queryGetAllByOrganizationID = database.BuildQuery(`
		SELECT `, columns, `, COUNT(*) OVER() AS total
		FROM `, tableName, `
		WHERE
			organization_id = $1
			AND deleted_at IS NULL
			AND ($2 = '' OR status ILIKE '%' || $2 || '%' OR execution_type ILIKE '%' || $2 || '%')
		ORDER BY
			updated_at DESC
		LIMIT $3 OFFSET $4
	`)
)

type TaskExecutionRepository struct {
	database.BaseRepository
}

func NewTaskExecutionRepository(db *sql.DB) taskexecutions.TaskExecutionRepository {
	return &TaskExecutionRepository{
		BaseRepository: database.NewBaseRepository(db),
	}
}

func (r *TaskExecutionRepository) scan(row interface{ Scan(dest ...any) error }, extras ...any) (*taskexecutions.TaskExecution, error) {
	var te taskexecutions.TaskExecution
	var deletedAt sql.NullTime
	var startedAt sql.NullTime
	var completedAt sql.NullTime
	var triggeredByUserID *uuid.UUID
	var flowTriggerID *uuid.UUID
	var errorMessage sql.NullString
	var executionContext string
	var nodes, edges []byte

	dest := []any{
		&te.ID,
		&te.OrganizationID,
		&triggeredByUserID,
		&flowTriggerID,
		&executionContext,
		&te.ContextItemID,
		&te.Status,
		&te.ExecutionType,
		&errorMessage,
		&nodes,
		&edges,
		&startedAt,
		&completedAt,
		&te.CreatedAt,
		&te.UpdatedAt,
		&deletedAt,
	}
	dest = append(dest, extras...)

	if err := row.Scan(dest...); err != nil {
		return nil, err
	}

	te.ExecutionContext = commonDomain.ExecutionContext(executionContext)
	te.TriggeredByUserID = triggeredByUserID
	te.FlowTriggerID = flowTriggerID

	if err := database.UnmarshalJSONColumn(nodes, &te.Nodes); err != nil {
		return nil, err
	}
	if err := database.UnmarshalJSONColumn(edges, &te.Edges); err != nil {
		return nil, err
	}

	te.ErrorMessage = database.ScanNullString(errorMessage)
	te.StartedAt = database.ScanNullTime(startedAt)
	te.CompletedAt = database.ScanNullTime(completedAt)
	te.DeletedAt = database.ScanNullTime(deletedAt)

	return &te, nil
}

func (r *TaskExecutionRepository) getOne(ctx context.Context, query string, args ...any) (*taskexecutions.TaskExecution, error) {
	return database.GetOne(ctx, r.Executor(ctx), query, r.scan, taskexecutions.ErrTaskExecutionNotFound, args...)
}

func (r *TaskExecutionRepository) getMany(ctx context.Context, query string, args ...any) ([]*taskexecutions.TaskExecution, error) {
	return database.GetMany(ctx, r.Executor(ctx), query, r.scan, args...)
}

func (r *TaskExecutionRepository) getManyWithTotal(ctx context.Context, query string, args ...any) ([]*taskexecutions.TaskExecution, int64, error) {
	return database.GetManyWithTotal(ctx, r.Executor(ctx), query, r.scan, args...)
}

func (r *TaskExecutionRepository) exec(ctx context.Context, query string, args ...any) error {
	return database.Exec(ctx, r.Executor(ctx), query, args...)
}

func (r *TaskExecutionRepository) execWithRowCheck(ctx context.Context, query string, args ...any) error {
	return database.ExecWithRowCheck(ctx, r.Executor(ctx), query, taskexecutions.ErrTaskExecutionNotFound, args...)
}

func (r *TaskExecutionRepository) GetAllByOrganizationID(
	ctx context.Context,
	organizationID uuid.UUID,
	searchQuery string,
	offset int,
	limit int,
) ([]*taskexecutions.TaskExecution, int64, error) {
	return r.getManyWithTotal(ctx, queryGetAllByOrganizationID, organizationID, searchQuery, limit, offset)
}

func (r *TaskExecutionRepository) GetByIDAndOrganizationID(ctx context.Context, id uuid.UUID, organizationID uuid.UUID) (*taskexecutions.TaskExecution, error) {
	return r.getOne(ctx, queryGetByIDAndOrganizationID, id, organizationID)
}

func (r *TaskExecutionRepository) GetByID(ctx context.Context, id uuid.UUID) (*taskexecutions.TaskExecution, error) {
	return r.getOne(ctx, queryGetByID, id)
}

func (r *TaskExecutionRepository) GetAllByStatuses(ctx context.Context, statuses []string) ([]*taskexecutions.TaskExecution, error) {
	if len(statuses) == 0 {
		return []*taskexecutions.TaskExecution{}, nil
	}
	return r.getMany(ctx, queryGetAllByStatuses, pq.Array(statuses))
}

func (r *TaskExecutionRepository) GetAllByIDsAndOrganizationID(ctx context.Context, ids []uuid.UUID, organizationID uuid.UUID) ([]*taskexecutions.TaskExecution, error) {
	if len(ids) == 0 {
		return []*taskexecutions.TaskExecution{}, nil
	}
	return r.getMany(ctx, queryGetAllByIDsAndOrganizationID, organizationID, pq.Array(ids))
}

func (r *TaskExecutionRepository) Create(ctx context.Context, te *taskexecutions.TaskExecution) error {
	nodesJSON, err := json.Marshal(te.Nodes)
	if err != nil {
		return err
	}
	edgesJSON, err := json.Marshal(te.Edges)
	if err != nil {
		return err
	}

	return r.exec(ctx, queryCreate,
		te.ID,
		te.OrganizationID,
		te.TriggeredByUserID,
		te.FlowTriggerID,
		te.ExecutionContext.String(),
		te.ContextItemID,
		te.Status,
		te.ExecutionType,
		te.ErrorMessage,
		nodesJSON,
		edgesJSON,
		te.StartedAt,
		te.CompletedAt,
		te.CreatedAt,
		te.UpdatedAt,
		te.DeletedAt,
	)
}

func (r *TaskExecutionRepository) Update(ctx context.Context, te *taskexecutions.TaskExecution) error {
	return r.execWithRowCheck(ctx, queryUpdate,
		te.ID,
		te.Status,
		te.ErrorMessage,
		te.CompletedAt,
		te.UpdatedAt,
	)
}

func (r *TaskExecutionRepository) Delete(ctx context.Context, te *taskexecutions.TaskExecution) error {
	return r.execWithRowCheck(ctx, queryDelete,
		te.ID,
		te.UpdatedAt,
		te.DeletedAt,
	)
}

func (r *TaskExecutionRepository) BulkDelete(ctx context.Context, ids []uuid.UUID, organizationID uuid.UUID, now time.Time) error {
	if len(ids) == 0 {
		return nil
	}
	return r.exec(ctx, queryBulkDelete,
		organizationID,
		pq.Array(ids),
		now,
		now,
	)
}
