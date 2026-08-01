package triggers

import (
	"context"
	"database/sql"

	"github.com/blocknextai/go-packages/database"
	"github.com/blocknextai/go-packages/json"
	"github.com/blocknextai/platform-api/internal/triggers/domain/triggers"
	"github.com/google/uuid"
)

const (
	tableName = "triggers.triggers"
	columns   = "id, organization_id, triggered_by_user_id, execution_context, context_item_id, type, cron_pattern, timezone, webhook_token_hash, webhook_secret, runtime_config, is_active, created_at, updated_at, deleted_at"
)

var (
	queryGetAllActive = database.BuildQuery(`
		SELECT `, columns, `
		FROM `, tableName, `
		WHERE
			is_active = true
			AND deleted_at IS NULL
		ORDER BY
			updated_at DESC
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

	queryGetAllByOrganizationID = database.BuildQuery(`
		SELECT `, columns, `, COUNT(*) OVER() AS total
		FROM `, tableName, `
		WHERE
			organization_id = $1
			AND deleted_at IS NULL
			AND ($2 = '' OR type ILIKE '%' || $2 || '%' OR COALESCE(cron_pattern, '') ILIKE '%' || $2 || '%')
		ORDER BY
			updated_at DESC
		LIMIT $3 OFFSET $4
	`)

	queryGetByWebhookTokenHash = database.BuildQuery(`
		SELECT `, columns, `
		FROM `, tableName, `
		WHERE
			webhook_token_hash = $1
			AND is_active = true
			AND deleted_at IS NULL
		LIMIT 1
	`)

	queryCreate = database.BuildQuery(`
		INSERT INTO `, tableName, ` (`, columns, `)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
	`)

	queryUpdate = database.BuildQuery(`
		UPDATE `, tableName, `
		SET
			is_active = $2,
			cron_pattern = $3,
			timezone = $4,
			webhook_token_hash = $5,
			webhook_secret = $6,
			runtime_config = $7,
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

type TriggerRepository struct {
	database.BaseRepository
}

func NewTriggerRepository(db *sql.DB) triggers.TriggerRepository {
	return &TriggerRepository{
		BaseRepository: database.NewBaseRepository(db),
	}
}

func (r *TriggerRepository) scan(row interface{ Scan(dest ...any) error }, extras ...any) (*triggers.Trigger, error) {
	var t triggers.Trigger
	var cronPattern sql.NullString
	var timezone sql.NullString
	var webhookTokenHash sql.NullString
	var webhookSecret sql.NullString
	var runtimeConfig []byte
	var deletedAt sql.NullTime

	dest := []any{
		&t.ID,
		&t.OrganizationID,
		&t.TriggeredByUserID,
		&t.ExecutionContext,
		&t.ContextItemID,
		&t.Type,
		&cronPattern,
		&timezone,
		&webhookTokenHash,
		&webhookSecret,
		&runtimeConfig,
		&t.IsActive,
		&t.CreatedAt,
		&t.UpdatedAt,
		&deletedAt,
	}
	dest = append(dest, extras...)

	if err := row.Scan(dest...); err != nil {
		return nil, err
	}

	if runtimeConfig != nil {
		var rc triggers.RuntimeConfig
		if err := database.UnmarshalJSONColumn(runtimeConfig, &rc); err == nil {
			t.RuntimeConfig = new(rc)
		}
	}

	t.CronPattern = database.ScanNullString(cronPattern)
	t.Timezone = database.ScanNullString(timezone)
	t.WebhookTokenHash = database.ScanNullString(webhookTokenHash)
	t.WebhookSecret = database.ScanNullString(webhookSecret)
	t.DeletedAt = database.ScanNullTime(deletedAt)

	return &t, nil
}

func (r *TriggerRepository) getOne(ctx context.Context, query string, args ...any) (*triggers.Trigger, error) {
	return database.GetOne(ctx, r.Executor(ctx), query, r.scan, triggers.ErrTriggerNotFound, args...)
}

func (r *TriggerRepository) getMany(ctx context.Context, query string, args ...any) ([]*triggers.Trigger, error) {
	return database.GetMany(ctx, r.Executor(ctx), query, r.scan, args...)
}

func (r *TriggerRepository) getManyWithTotal(ctx context.Context, query string, args ...any) ([]*triggers.Trigger, int64, error) {
	return database.GetManyWithTotal(ctx, r.Executor(ctx), query, r.scan, args...)
}

func (r *TriggerRepository) exec(ctx context.Context, query string, args ...any) error {
	return database.Exec(ctx, r.Executor(ctx), query, args...)
}

func (r *TriggerRepository) execWithRowCheck(ctx context.Context, query string, args ...any) error {
	return database.ExecWithRowCheck(ctx, r.Executor(ctx), query, triggers.ErrTriggerNotFound, args...)
}

func (r *TriggerRepository) GetAllActive(ctx context.Context) ([]*triggers.Trigger, error) {
	return r.getMany(ctx, queryGetAllActive)
}

func (r *TriggerRepository) GetByIDAndOrganizationID(ctx context.Context, id uuid.UUID, organizationID uuid.UUID) (*triggers.Trigger, error) {
	return r.getOne(ctx, queryGetByIDAndOrganizationID, id, organizationID)
}

func (r *TriggerRepository) GetByWebhookTokenHash(ctx context.Context, tokenHash string) (*triggers.Trigger, error) {
	return r.getOne(ctx, queryGetByWebhookTokenHash, tokenHash)
}

func (r *TriggerRepository) GetAllByOrganizationID(
	ctx context.Context,
	organizationID uuid.UUID,
	searchQuery string,
	offset int,
	limit int,
) ([]*triggers.Trigger, int64, error) {
	return r.getManyWithTotal(ctx, queryGetAllByOrganizationID, organizationID, searchQuery, limit, offset)
}

func (r *TriggerRepository) Create(ctx context.Context, trigger *triggers.Trigger) error {
	var runtimeConfigJSON []byte
	if trigger.RuntimeConfig != nil {
		var err error
		runtimeConfigJSON, err = json.Marshal(trigger.RuntimeConfig)
		if err != nil {
			return err
		}
	}

	err := r.exec(ctx, queryCreate,
		trigger.ID,
		trigger.OrganizationID,
		trigger.TriggeredByUserID,
		trigger.ExecutionContext,
		trigger.ContextItemID,
		trigger.Type,
		trigger.CronPattern,
		trigger.Timezone,
		trigger.WebhookTokenHash,
		trigger.WebhookSecret,
		runtimeConfigJSON,
		trigger.IsActive,
		trigger.CreatedAt,
		trigger.UpdatedAt,
		trigger.DeletedAt,
	)
	if err != nil && database.IsUniqueViolationOn(err, "uq_triggers_webhook_token_hash") {
		return triggers.ErrWebhookTokenTaken
	}
	return err
}

func (r *TriggerRepository) Update(ctx context.Context, trigger *triggers.Trigger) error {
	var runtimeConfigJSON []byte
	if trigger.RuntimeConfig != nil {
		var runtimeConfigErr error
		runtimeConfigJSON, runtimeConfigErr = json.Marshal(trigger.RuntimeConfig)
		if runtimeConfigErr != nil {
			return runtimeConfigErr
		}
	}

	err := r.execWithRowCheck(ctx, queryUpdate,
		trigger.ID,
		trigger.IsActive,
		trigger.CronPattern,
		trigger.Timezone,
		trigger.WebhookTokenHash,
		trigger.WebhookSecret,
		runtimeConfigJSON,
		trigger.UpdatedAt,
	)
	if err != nil && database.IsUniqueViolationOn(err, "uq_triggers_webhook_token_hash") {
		return triggers.ErrWebhookTokenTaken
	}
	return err
}

func (r *TriggerRepository) Delete(ctx context.Context, trigger *triggers.Trigger) error {
	return r.execWithRowCheck(ctx, queryDelete,
		trigger.ID,
		trigger.UpdatedAt,
		trigger.DeletedAt,
	)
}
