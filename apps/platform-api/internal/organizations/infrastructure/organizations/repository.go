package organizations

import (
	"context"
	"database/sql"

	"github.com/blocknextai/go-packages/database"
	"github.com/blocknextai/platform-api/internal/organizations/domain/organizations"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

const (
	tableName = "organizations.organizations"
	columns   = "id, title, description, is_verified, created_at, updated_at, deleted_at"
)

var (
	queryGetAllByUserID = database.BuildQuery(`
		SELECT `, columns, `
		FROM `, tableName, `
		WHERE
			id IN (SELECT organization_id FROM organizations.users WHERE user_id = $1 AND deleted_at IS NULL)
			AND deleted_at IS NULL
		ORDER BY
			updated_at DESC
	`)

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
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`)

	queryUpdate = database.BuildQuery(`
		UPDATE `, tableName, `
		SET
			title = $2,
			description = $3,
			updated_at = $4
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

	queryGetAllByIDs = database.BuildQuery(`
		SELECT `, columns, `
		FROM `, tableName, `
		WHERE
			id = ANY($1)
			AND deleted_at IS NULL
	`)
)

type OrganizationRepository struct {
	database.BaseRepository
}

func NewOrganizationRepository(db *sql.DB) organizations.OrganizationRepository {
	return &OrganizationRepository{
		BaseRepository: database.NewBaseRepository(db),
	}
}

func (r *OrganizationRepository) scan(row interface{ Scan(dest ...any) error }, extras ...any) (*organizations.Organization, error) {
	var org organizations.Organization
	var description sql.NullString
	var deletedAt sql.NullTime

	dest := []any{
		&org.ID,
		&org.Title,
		&description,
		&org.IsVerified,
		&org.CreatedAt,
		&org.UpdatedAt,
		&deletedAt,
	}
	dest = append(dest, extras...)

	if err := row.Scan(dest...); err != nil {
		return nil, err
	}

	org.Description = database.ScanNullString(description)
	org.DeletedAt = database.ScanNullTime(deletedAt)

	return &org, nil
}

func (r *OrganizationRepository) getOne(ctx context.Context, query string, args ...any) (*organizations.Organization, error) {
	return database.GetOne(ctx, r.Executor(ctx), query, r.scan, organizations.ErrOrganizationNotFound, args...)
}

func (r *OrganizationRepository) getMany(ctx context.Context, query string, args ...any) ([]*organizations.Organization, error) {
	return database.GetMany(ctx, r.Executor(ctx), query, r.scan, args...)
}

func (r *OrganizationRepository) exec(ctx context.Context, query string, args ...any) error {
	return database.Exec(ctx, r.Executor(ctx), query, args...)
}

func (r *OrganizationRepository) execWithRowCheck(ctx context.Context, query string, args ...any) error {
	return database.ExecWithRowCheck(ctx, r.Executor(ctx), query, organizations.ErrOrganizationNotFound, args...)
}

func (r *OrganizationRepository) GetAllByUserID(ctx context.Context, userID uuid.UUID) ([]*organizations.Organization, error) {
	return r.getMany(ctx, queryGetAllByUserID, userID)
}

func (r *OrganizationRepository) GetByID(ctx context.Context, id uuid.UUID) (*organizations.Organization, error) {
	return r.getOne(ctx, queryGetByID, id)
}

func (r *OrganizationRepository) Create(ctx context.Context, organization *organizations.Organization) error {
	return r.exec(ctx, queryCreate,
		organization.ID,
		organization.Title,
		organization.Description,
		organization.IsVerified,
		organization.CreatedAt,
		organization.UpdatedAt,
		organization.DeletedAt,
	)
}

func (r *OrganizationRepository) Update(ctx context.Context, organization *organizations.Organization) error {
	return r.execWithRowCheck(ctx, queryUpdate,
		organization.ID,
		organization.Title,
		organization.Description,
		organization.UpdatedAt,
	)
}

func (r *OrganizationRepository) Delete(ctx context.Context, organization *organizations.Organization) error {
	return r.execWithRowCheck(ctx, queryDelete,
		organization.ID,
		organization.UpdatedAt,
		organization.DeletedAt,
	)
}

func (r *OrganizationRepository) GetAllByIDs(ctx context.Context, ids []uuid.UUID) ([]*organizations.Organization, error) {
	if len(ids) == 0 {
		return []*organizations.Organization{}, nil
	}
	return r.getMany(ctx, queryGetAllByIDs, pq.Array(ids))
}
