package taskclaims

import (
	"context"
	"database/sql"
	"time"

	"github.com/blocknextai/go-packages/database"
	"github.com/blocknextai/platform-api/internal/executions/domain/taskclaims"
	"github.com/google/uuid"
)

const (
	tableName = "executions.task_claims"
	columns   = "id, task_execution_id, claimed_by, claimed_at, retry_count, created_at, updated_at, deleted_at"
)

var (
	queryCreate = database.BuildQuery(`
		INSERT INTO `, tableName, ` (`, columns, `)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`)

	queryGetByTaskExecutionID = database.BuildQuery(`
		SELECT `, columns, `
		FROM `, tableName, `
		WHERE
			task_execution_id = $1
			AND deleted_at IS NULL
		LIMIT 1
	`)

	queryGetByTaskExecutionIDForUpdate = database.BuildQuery(`
		SELECT `, columns, `
		FROM `, tableName, `
		WHERE
			task_execution_id = $1
			AND deleted_at IS NULL
		LIMIT 1
		FOR UPDATE
	`)

	queryGetAllStale = database.BuildQuery(`
		SELECT `, columns, `
		FROM `, tableName, `
		WHERE
			claimed_at IS NOT NULL
			AND claimed_at < (NOW() AT TIME ZONE 'UTC') - make_interval(secs => $1)
			AND deleted_at IS NULL
	`)

	queryUpdate = database.BuildQuery(`
		UPDATE `, tableName, `
		SET
			claimed_by = $2,
			claimed_at = $3,
			retry_count = $4,
			updated_at = $5,
			deleted_at = $6
		WHERE
			id = $1
	`)
)

type TaskClaimRepository struct {
	database.BaseRepository
}

func NewTaskClaimRepository(db *sql.DB) taskclaims.TaskClaimRepository {
	return &TaskClaimRepository{
		BaseRepository: database.NewBaseRepository(db),
	}
}

func (r *TaskClaimRepository) scan(row interface{ Scan(dest ...any) error }, _ ...any) (*taskclaims.TaskClaim, error) {
	var taskClaim taskclaims.TaskClaim
	if err := row.Scan(
		&taskClaim.ID,
		&taskClaim.TaskExecutionID,
		&taskClaim.ClaimedBy,
		&taskClaim.ClaimedAt,
		&taskClaim.RetryCount,
		&taskClaim.CreatedAt,
		&taskClaim.UpdatedAt,
		&taskClaim.DeletedAt,
	); err != nil {
		return nil, err
	}
	return &taskClaim, nil
}

func (r *TaskClaimRepository) getOne(ctx context.Context, query string, args ...any) (*taskclaims.TaskClaim, error) {
	return database.GetOne(ctx, r.Executor(ctx), query, r.scan, taskclaims.ErrTaskClaimNotFound, args...)
}

func (r *TaskClaimRepository) getMany(ctx context.Context, query string, args ...any) ([]*taskclaims.TaskClaim, error) {
	return database.GetMany(ctx, r.Executor(ctx), query, r.scan, args...)
}

func (r *TaskClaimRepository) exec(ctx context.Context, query string, args ...any) error {
	return database.Exec(ctx, r.Executor(ctx), query, args...)
}

func (r *TaskClaimRepository) execWithRowCheck(ctx context.Context, query string, args ...any) error {
	return database.ExecWithRowCheck(ctx, r.Executor(ctx), query, taskclaims.ErrTaskClaimNotFound, args...)
}

func (r *TaskClaimRepository) Create(ctx context.Context, taskClaim *taskclaims.TaskClaim) error {
	return r.exec(ctx, queryCreate,
		taskClaim.ID,
		taskClaim.TaskExecutionID,
		taskClaim.ClaimedBy,
		taskClaim.ClaimedAt,
		taskClaim.RetryCount,
		taskClaim.CreatedAt,
		taskClaim.UpdatedAt,
		taskClaim.DeletedAt,
	)
}

func (r *TaskClaimRepository) GetByTaskExecutionID(ctx context.Context, taskExecutionID uuid.UUID) (*taskclaims.TaskClaim, error) {
	return r.getOne(ctx, queryGetByTaskExecutionID, taskExecutionID)
}

func (r *TaskClaimRepository) GetByTaskExecutionIDForUpdate(ctx context.Context, taskExecutionID uuid.UUID) (*taskclaims.TaskClaim, error) {
	return r.getOne(ctx, queryGetByTaskExecutionIDForUpdate, taskExecutionID)
}

func (r *TaskClaimRepository) GetAllStale(ctx context.Context, staleAfter time.Duration) ([]*taskclaims.TaskClaim, error) {
	return r.getMany(ctx, queryGetAllStale, staleAfter.Seconds())
}

func (r *TaskClaimRepository) Update(ctx context.Context, taskClaim *taskclaims.TaskClaim) error {
	return r.execWithRowCheck(ctx, queryUpdate,
		taskClaim.ID,
		taskClaim.ClaimedBy,
		taskClaim.ClaimedAt,
		taskClaim.RetryCount,
		taskClaim.UpdatedAt,
		taskClaim.DeletedAt,
	)
}
