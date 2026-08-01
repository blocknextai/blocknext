package organizationusers

import (
	"context"
	"database/sql"
	"time"

	"github.com/blocknextai/go-packages/database"
	"github.com/blocknextai/go-packages/rbac"
	"github.com/blocknextai/platform-api/internal/organizations/domain/organizationusers"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

const (
	tableName = "organizations.users"
	columns   = "id, organization_id, user_id, role, alias, created_at, updated_at, deleted_at"
)

var (
	queryGetAllByIDs = database.BuildQuery(`
		SELECT `, columns, `
		FROM `, tableName, `
		WHERE
			id = ANY($1)
			AND deleted_at IS NULL
	`)

	queryGetAllByOrganizationID = database.BuildQuery(`
		SELECT `, columns, `, COUNT(*) OVER() AS total
		FROM `, tableName, `
		WHERE
			organization_id = $1
			AND deleted_at IS NULL
			AND ($2 = '' OR COALESCE(alias, '') ILIKE '%' || $2 || '%')
		ORDER BY
			CASE role
				WHEN 'organization:owner' THEN 100
				WHEN 'organization:admin' THEN 75
				WHEN 'organization:editor' THEN 50
				WHEN 'organization:viewer' THEN 25
				ELSE 0
			END DESC,
			created_at DESC
		LIMIT $3 OFFSET $4
	`)

	queryGetAllByOrganizationIDAndUserIDs = database.BuildQuery(`
		SELECT `, columns, `
		FROM `, tableName, `
		WHERE
			organization_id = $1
			AND user_id = ANY($2)
			AND deleted_at IS NULL
	`)

	queryGetAllByUserID = database.BuildQuery(`
		SELECT `, columns, `
		FROM `, tableName, `
		WHERE
			user_id = $1
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

	queryGetByIDAndOrganizationID = database.BuildQuery(`
		SELECT `, columns, `
		FROM `, tableName, `
		WHERE
			id = $1
			AND organization_id = $2
			AND deleted_at IS NULL
		LIMIT 1
	`)

	queryGetByOrganizationIDAndUserID = database.BuildQuery(`
		SELECT `, columns, `
		FROM `, tableName, `
		WHERE
			organization_id = $1
			AND user_id = $2
			AND deleted_at IS NULL
		LIMIT 1
	`)

	queryCountByOrganizationID = database.BuildQuery(`
		SELECT COUNT(*)
		FROM `, tableName, `
		WHERE
			organization_id = $1
			AND deleted_at IS NULL
	`)

	queryHasOwner = database.BuildQuery(`
		SELECT COUNT(*)
		FROM `, tableName, `
		WHERE
			organization_id = $1
			AND role = $2
			AND deleted_at IS NULL
	`)

	queryCreate = database.BuildQuery(`
		INSERT INTO `, tableName, ` (`, columns, `)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`)

	queryUpdate = database.BuildQuery(`
		UPDATE `, tableName, `
		SET
			role = $2,
			alias = $3,
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

	queryDeleteAllByOrganizationID = database.BuildQuery(`
		UPDATE `, tableName, `
		SET
			updated_at = $2,
			deleted_at = $2
		WHERE
			organization_id = $1
			AND deleted_at IS NULL
	`)
)

type OrganizationUserRepository struct {
	database.BaseRepository
}

func NewOrganizationUserRepository(db *sql.DB) organizationusers.OrganizationUserRepository {
	return &OrganizationUserRepository{
		BaseRepository: database.NewBaseRepository(db),
	}
}

func (r *OrganizationUserRepository) scan(row interface{ Scan(dest ...any) error }, extras ...any) (*organizationusers.OrganizationUser, error) {
	var ou organizationusers.OrganizationUser
	var alias sql.NullString
	var deletedAt sql.NullTime

	dest := []any{
		&ou.ID,
		&ou.OrganizationID,
		&ou.UserID,
		&ou.Role,
		&alias,
		&ou.CreatedAt,
		&ou.UpdatedAt,
		&deletedAt,
	}
	dest = append(dest, extras...)

	if err := row.Scan(dest...); err != nil {
		return nil, err
	}

	if alias.Valid {
		ou.Alias = alias.String
	}

	ou.DeletedAt = database.ScanNullTime(deletedAt)

	return &ou, nil
}

func (r *OrganizationUserRepository) getOne(ctx context.Context, query string, args ...any) (*organizationusers.OrganizationUser, error) {
	return database.GetOne(ctx, r.Executor(ctx), query, r.scan, organizationusers.ErrOrganizationUserNotFound, args...)
}

func (r *OrganizationUserRepository) getMany(ctx context.Context, query string, args ...any) ([]*organizationusers.OrganizationUser, error) {
	return database.GetMany(ctx, r.Executor(ctx), query, r.scan, args...)
}

func (r *OrganizationUserRepository) getManyWithTotal(ctx context.Context, query string, args ...any) ([]*organizationusers.OrganizationUser, int64, error) {
	return database.GetManyWithTotal(ctx, r.Executor(ctx), query, r.scan, args...)
}

func (r *OrganizationUserRepository) exec(ctx context.Context, query string, args ...any) error {
	return database.Exec(ctx, r.Executor(ctx), query, args...)
}

func (r *OrganizationUserRepository) execWithRowCheck(ctx context.Context, query string, args ...any) error {
	return database.ExecWithRowCheck(ctx, r.Executor(ctx), query, organizationusers.ErrOrganizationUserNotFound, args...)
}

func (r *OrganizationUserRepository) GetAllByIDs(ctx context.Context, ids []uuid.UUID) ([]*organizationusers.OrganizationUser, error) {
	if len(ids) == 0 {
		return []*organizationusers.OrganizationUser{}, nil
	}
	return r.getMany(ctx, queryGetAllByIDs, pq.Array(ids))
}

func (r *OrganizationUserRepository) GetAllByOrganizationID(ctx context.Context, organizationID uuid.UUID, searchQuery string, offset int, limit int) ([]*organizationusers.OrganizationUser, int64, error) {
	return r.getManyWithTotal(ctx, queryGetAllByOrganizationID, organizationID, searchQuery, limit, offset)
}

func (r *OrganizationUserRepository) GetAllByOrganizationIDAndUserIDs(ctx context.Context, organizationID uuid.UUID, userIDs []uuid.UUID) ([]*organizationusers.OrganizationUser, error) {
	if len(userIDs) == 0 {
		return []*organizationusers.OrganizationUser{}, nil
	}
	return r.getMany(ctx, queryGetAllByOrganizationIDAndUserIDs, organizationID, pq.Array(userIDs))
}

func (r *OrganizationUserRepository) GetAllByUserID(ctx context.Context, userID uuid.UUID) ([]*organizationusers.OrganizationUser, error) {
	return r.getMany(ctx, queryGetAllByUserID, userID)
}

func (r *OrganizationUserRepository) GetByID(ctx context.Context, id uuid.UUID) (*organizationusers.OrganizationUser, error) {
	return r.getOne(ctx, queryGetByID, id)
}

func (r *OrganizationUserRepository) GetByIDAndOrganizationID(ctx context.Context, id uuid.UUID, organizationID uuid.UUID) (*organizationusers.OrganizationUser, error) {
	return r.getOne(ctx, queryGetByIDAndOrganizationID, id, organizationID)
}

func (r *OrganizationUserRepository) GetByOrganizationIDAndUserID(ctx context.Context, organizationID uuid.UUID, userID uuid.UUID) (*organizationusers.OrganizationUser, error) {
	return r.getOne(ctx, queryGetByOrganizationIDAndUserID, organizationID, userID)
}

func (r *OrganizationUserRepository) CountByOrganizationID(ctx context.Context, organizationID uuid.UUID) (int64, error) {
	var count int64
	err := r.Executor(ctx).QueryRowContext(ctx, queryCountByOrganizationID, organizationID).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *OrganizationUserRepository) HasOwner(ctx context.Context, organizationID uuid.UUID) (bool, error) {
	var count int64
	err := r.Executor(ctx).QueryRowContext(ctx, queryHasOwner, organizationID, rbac.OrganizationOwnerRole.Name).Scan(&count)
	if err != nil {
		return false, err
	}
	return count >= 1, nil
}

func (r *OrganizationUserRepository) Create(ctx context.Context, organizationUser *organizationusers.OrganizationUser) error {
	err := r.exec(ctx, queryCreate,
		organizationUser.ID,
		organizationUser.OrganizationID,
		organizationUser.UserID,
		organizationUser.Role,
		organizationUser.Alias,
		organizationUser.CreatedAt,
		organizationUser.UpdatedAt,
		organizationUser.DeletedAt,
	)
	if err != nil && database.IsUniqueViolationOn(err, "uq_organizations_users_org_user") {
		return organizationusers.ErrOrganizationUserAlreadyExists
	}
	return err
}

func (r *OrganizationUserRepository) Update(ctx context.Context, organizationUser *organizationusers.OrganizationUser) error {
	return r.execWithRowCheck(ctx, queryUpdate,
		organizationUser.ID,
		organizationUser.Role,
		organizationUser.Alias,
		organizationUser.UpdatedAt,
	)
}

func (r *OrganizationUserRepository) Delete(ctx context.Context, organizationUser *organizationusers.OrganizationUser) error {
	return r.execWithRowCheck(ctx, queryDelete,
		organizationUser.ID,
		organizationUser.UpdatedAt,
		organizationUser.DeletedAt,
	)
}

func (r *OrganizationUserRepository) DeleteAllByOrganizationID(ctx context.Context, organizationID uuid.UUID, deletedAt time.Time) error {
	return r.exec(ctx, queryDeleteAllByOrganizationID, organizationID, deletedAt)
}
