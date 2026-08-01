package nodeexecutions

import (
	"context"
	"database/sql"

	"github.com/blocknextai/go-packages/database"
	"github.com/blocknextai/go-packages/json"
	"github.com/blocknextai/platform-api/internal/executions/domain/nodeexecutions"
	"github.com/google/uuid"
)

const (
	tableName = "executions.node_executions"
	columns   = "id, task_id, node_type, node_id, status, inputs, outputs, function_calling_outputs, error_message, started_at, completed_at, created_at, updated_at, deleted_at"
)

var (
	queryGetByIDAndTaskID = database.BuildQuery(`
		SELECT `, columns, `
		FROM `, tableName, `
		WHERE
			id = $1
			AND task_id = $2
			AND deleted_at IS NULL
		LIMIT 1
	`)

	queryGetAllByTaskID = database.BuildQuery(`
		SELECT `, columns, `
		FROM `, tableName, `
		WHERE
			task_id = $1
			AND deleted_at IS NULL
		ORDER BY
			started_at ASC
	`)

	queryCreate = database.BuildQuery(`
		INSERT INTO `, tableName, ` (`, columns, `)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
	`)

	queryUpdate = database.BuildQuery(`
		UPDATE `, tableName, `
		SET
			status = $2,
			inputs = $3,
			outputs = $4,
			function_calling_outputs = $5,
			error_message = $6,
			started_at = $7,
			completed_at = $8,
			updated_at = $9
		WHERE
			id = $1
			AND deleted_at IS NULL
	`)
)

type NodeExecutionRepository struct {
	database.BaseRepository
}

func NewNodeExecutionRepository(db *sql.DB) nodeexecutions.NodeExecutionRepository {
	return &NodeExecutionRepository{
		BaseRepository: database.NewBaseRepository(db),
	}
}

func (r *NodeExecutionRepository) scan(row interface{ Scan(dest ...any) error }, _ ...any) (*nodeexecutions.NodeExecution, error) {
	var ne nodeexecutions.NodeExecution
	var deletedAt sql.NullTime
	var startedAt sql.NullTime
	var completedAt sql.NullTime
	var errorMessage sql.NullString
	var inputs, outputs, functionCallingOutputs []byte

	err := row.Scan(
		&ne.ID,
		&ne.TaskID,
		&ne.NodeType,
		&ne.NodeID,
		&ne.Status,
		&inputs,
		&outputs,
		&functionCallingOutputs,
		&errorMessage,
		&startedAt,
		&completedAt,
		&ne.CreatedAt,
		&ne.UpdatedAt,
		&deletedAt,
	)
	if err != nil {
		return nil, err
	}

	if err := database.UnmarshalJSONColumn(inputs, &ne.Inputs); err != nil {
		return nil, err
	}
	if err := database.UnmarshalJSONColumn(outputs, &ne.Outputs); err != nil {
		return nil, err
	}
	if err := database.UnmarshalJSONColumn(functionCallingOutputs, &ne.FunctionCallingOutputs); err != nil {
		return nil, err
	}

	ne.ErrorMessage = database.ScanNullString(errorMessage)
	ne.StartedAt = database.ScanNullTime(startedAt)
	ne.CompletedAt = database.ScanNullTime(completedAt)
	ne.DeletedAt = database.ScanNullTime(deletedAt)

	return &ne, nil
}

func (r *NodeExecutionRepository) getOne(ctx context.Context, query string, args ...any) (*nodeexecutions.NodeExecution, error) {
	return database.GetOne(ctx, r.Executor(ctx), query, r.scan, nodeexecutions.ErrNodeExecutionNotFound, args...)
}

func (r *NodeExecutionRepository) getMany(ctx context.Context, query string, args ...any) ([]*nodeexecutions.NodeExecution, error) {
	return database.GetMany(ctx, r.Executor(ctx), query, r.scan, args...)
}

func (r *NodeExecutionRepository) exec(ctx context.Context, query string, args ...any) error {
	return database.Exec(ctx, r.Executor(ctx), query, args...)
}

func (r *NodeExecutionRepository) execWithRowCheck(ctx context.Context, query string, args ...any) error {
	return database.ExecWithRowCheck(ctx, r.Executor(ctx), query, nodeexecutions.ErrNodeExecutionNotFound, args...)
}

func (r *NodeExecutionRepository) GetByIDAndTaskID(ctx context.Context, nodeExecutionID uuid.UUID, taskID uuid.UUID) (*nodeexecutions.NodeExecution, error) {
	return r.getOne(ctx, queryGetByIDAndTaskID, nodeExecutionID, taskID)
}

func (r *NodeExecutionRepository) GetAllByTaskID(ctx context.Context, taskID uuid.UUID) ([]*nodeexecutions.NodeExecution, error) {
	return r.getMany(ctx, queryGetAllByTaskID, taskID)
}

func (r *NodeExecutionRepository) Create(ctx context.Context, ne *nodeexecutions.NodeExecution) error {
	inputsJSON, err := json.Marshal(ne.Inputs)
	if err != nil {
		return err
	}
	outputsJSON, err := json.Marshal(ne.Outputs)
	if err != nil {
		return err
	}
	fcOutputsJSON, err := json.Marshal(ne.FunctionCallingOutputs)
	if err != nil {
		return err
	}

	return r.exec(ctx, queryCreate,
		ne.ID,
		ne.TaskID,
		ne.NodeType,
		ne.NodeID,
		ne.Status,
		inputsJSON,
		outputsJSON,
		fcOutputsJSON,
		ne.ErrorMessage,
		ne.StartedAt,
		ne.CompletedAt,
		ne.CreatedAt,
		ne.UpdatedAt,
		ne.DeletedAt,
	)
}

func (r *NodeExecutionRepository) Update(ctx context.Context, ne *nodeexecutions.NodeExecution) error {
	inputsJSON, err := json.Marshal(ne.Inputs)
	if err != nil {
		return err
	}
	outputsJSON, err := json.Marshal(ne.Outputs)
	if err != nil {
		return err
	}
	fcOutputsJSON, err := json.Marshal(ne.FunctionCallingOutputs)
	if err != nil {
		return err
	}

	return r.execWithRowCheck(ctx, queryUpdate,
		ne.ID,
		ne.Status,
		inputsJSON,
		outputsJSON,
		fcOutputsJSON,
		ne.ErrorMessage,
		ne.StartedAt,
		ne.CompletedAt,
		ne.UpdatedAt,
	)
}
