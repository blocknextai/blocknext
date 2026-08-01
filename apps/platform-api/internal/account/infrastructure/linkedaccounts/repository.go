package linkedaccounts

import (
	"context"
	"database/sql"

	"github.com/blocknextai/go-packages/database"
	accountDomain "github.com/blocknextai/platform-api/internal/account/domain"
	"github.com/blocknextai/platform-api/internal/account/domain/linkedaccounts"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

const (
	tableName = "account.linked_accounts"
	columns   = "id, user_id, auth_provider, provider_id, identifier, display_name, is_primary, is_verified, created_at, updated_at, deleted_at"
)

var (
	queryGetByIDAndUserID = database.BuildQuery(`
		SELECT `, columns, `
		FROM `, tableName, `
		WHERE
			id = $1
			AND user_id = $2
			AND deleted_at IS NULL
		LIMIT 1
	`)

	queryGetByAuthProviderAndProviderID = database.BuildQuery(`
		SELECT `, columns, `
		FROM `, tableName, `
		WHERE
			auth_provider = $1
			AND provider_id = $2
			AND deleted_at IS NULL
		LIMIT 1
	`)

	queryGetByAuthProviderAndUserID = database.BuildQuery(`
		SELECT `, columns, `
		FROM `, tableName, `
		WHERE
			auth_provider = $1
			AND user_id = $2
			AND deleted_at IS NULL
		LIMIT 1
	`)

	queryGetByProviderID = database.BuildQuery(`
		SELECT `, columns, `
		FROM `, tableName, `
		WHERE
			provider_id = $1
			AND deleted_at IS NULL
		LIMIT 1
	`)

	queryGetAllByUserID = database.BuildQuery(`
		SELECT `, columns, `
		FROM `, tableName, `
		WHERE
			user_id = $1
			AND deleted_at IS NULL
		ORDER BY
			created_at ASC
	`)

	queryGetAllByUserIDs = database.BuildQuery(`
		SELECT `, columns, `
		FROM `, tableName, `
		WHERE
			user_id = ANY($1)
			AND deleted_at IS NULL
		ORDER BY
			created_at ASC
	`)

	queryGetByIdentifier = database.BuildQuery(`
		SELECT `, columns, `
		FROM `, tableName, `
		WHERE
			identifier = $1
			AND deleted_at IS NULL
		LIMIT 1
	`)

	queryCreate = database.BuildQuery(`
		INSERT INTO `, tableName, ` (`, columns, `)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`)

	queryUpdate = database.BuildQuery(`
		UPDATE `, tableName, `
		SET
			provider_id = $2,
			identifier = $3,
			display_name = $4,
			is_verified = $5,
			updated_at = $6
		WHERE
			id = $1
			AND deleted_at IS NULL
	`)

	queryDelete = database.BuildQuery(`
		UPDATE `, tableName, `
		SET
			deleted_at = $2
		WHERE
			id = $1
			AND deleted_at IS NULL
	`)
)

type LinkedAccountRepository struct {
	database.BaseRepository
}

func NewLinkedAccountRepository(db *sql.DB) linkedaccounts.LinkedAccountRepository {
	return &LinkedAccountRepository{
		BaseRepository: database.NewBaseRepository(db),
	}
}

func (r *LinkedAccountRepository) scan(row interface{ Scan(dest ...any) error }, _ ...any) (*linkedaccounts.LinkedAccount, error) {
	var linkedAccount linkedaccounts.LinkedAccount
	var displayName sql.NullString
	var deletedAt sql.NullTime

	err := row.Scan(
		&linkedAccount.ID,
		&linkedAccount.UserID,
		&linkedAccount.AuthProvider,
		&linkedAccount.ProviderID,
		&linkedAccount.Identifier,
		&displayName,
		&linkedAccount.IsPrimary,
		&linkedAccount.IsVerified,
		&linkedAccount.CreatedAt,
		&linkedAccount.UpdatedAt,
		&deletedAt,
	)

	if err != nil {
		return nil, err
	}

	linkedAccount.DisplayName = database.ScanNullString(displayName)
	linkedAccount.DeletedAt = database.ScanNullTime(deletedAt)

	return &linkedAccount, nil
}

func (r *LinkedAccountRepository) getOne(ctx context.Context, query string, args ...any) (*linkedaccounts.LinkedAccount, error) {
	return database.GetOne(ctx, r.Executor(ctx), query, r.scan, linkedaccounts.ErrLinkedAccountNotFound, args...)
}

func (r *LinkedAccountRepository) getMany(ctx context.Context, query string, args ...any) ([]*linkedaccounts.LinkedAccount, error) {
	return database.GetMany(ctx, r.Executor(ctx), query, r.scan, args...)
}

func (r *LinkedAccountRepository) exec(ctx context.Context, query string, args ...any) error {
	return database.Exec(ctx, r.Executor(ctx), query, args...)
}

func (r *LinkedAccountRepository) execWithRowCheck(ctx context.Context, query string, args ...any) error {
	return database.ExecWithRowCheck(ctx, r.Executor(ctx), query, linkedaccounts.ErrLinkedAccountNotFound, args...)
}

func (r *LinkedAccountRepository) GetByIDAndUserID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*linkedaccounts.LinkedAccount, error) {
	return r.getOne(ctx, queryGetByIDAndUserID, id, userID)
}

func (r *LinkedAccountRepository) GetAllByUserID(ctx context.Context, userID uuid.UUID) ([]*linkedaccounts.LinkedAccount, error) {
	return r.getMany(ctx, queryGetAllByUserID, userID)
}

func (r *LinkedAccountRepository) GetByAuthProviderAndProviderID(ctx context.Context, authProvider accountDomain.AuthProvider, providerID string) (*linkedaccounts.LinkedAccount, error) {
	return r.getOne(ctx, queryGetByAuthProviderAndProviderID, authProvider, providerID)
}

func (r *LinkedAccountRepository) GetByAuthProviderAndUserID(ctx context.Context, authProvider accountDomain.AuthProvider, userID uuid.UUID) (*linkedaccounts.LinkedAccount, error) {
	return r.getOne(ctx, queryGetByAuthProviderAndUserID, authProvider, userID)
}

func (r *LinkedAccountRepository) GetByProviderID(ctx context.Context, providerID string) (*linkedaccounts.LinkedAccount, error) {
	return r.getOne(ctx, queryGetByProviderID, providerID)
}

func (r *LinkedAccountRepository) GetByIdentifier(ctx context.Context, identifier string) (*linkedaccounts.LinkedAccount, error) {
	return r.getOne(ctx, queryGetByIdentifier, identifier)
}

func (r *LinkedAccountRepository) Create(ctx context.Context, linkedAccount *linkedaccounts.LinkedAccount) error {
	err := r.exec(ctx, queryCreate,
		linkedAccount.ID,
		linkedAccount.UserID,
		linkedAccount.AuthProvider,
		linkedAccount.ProviderID,
		linkedAccount.Identifier,
		linkedAccount.DisplayName,
		linkedAccount.IsPrimary,
		linkedAccount.IsVerified,
		linkedAccount.CreatedAt,
		linkedAccount.UpdatedAt,
		linkedAccount.DeletedAt,
	)
	if err != nil {
		if database.IsUniqueViolationOn(err, "uq_account_linked_accounts_provider_lookup") {
			return linkedaccounts.ErrLinkedAccountProviderTaken
		}
		if database.IsUniqueViolationOn(err, "uq_account_linked_accounts_user_provider") {
			return linkedaccounts.ErrLinkedAccountUserProviderTaken
		}
		return err
	}

	return nil
}

func (r *LinkedAccountRepository) Update(ctx context.Context, linkedAccount *linkedaccounts.LinkedAccount) error {
	return r.execWithRowCheck(ctx, queryUpdate,
		linkedAccount.ID,
		linkedAccount.ProviderID,
		linkedAccount.Identifier,
		linkedAccount.DisplayName,
		linkedAccount.IsVerified,
		linkedAccount.UpdatedAt,
	)
}

func (r *LinkedAccountRepository) Delete(ctx context.Context, linkedAccount *linkedaccounts.LinkedAccount) error {
	return r.execWithRowCheck(ctx, queryDelete,
		linkedAccount.ID,
		linkedAccount.DeletedAt,
	)
}

func (r *LinkedAccountRepository) GetAllByUserIDs(ctx context.Context, userIDs []uuid.UUID) ([]*linkedaccounts.LinkedAccount, error) {
	if len(userIDs) == 0 {
		return []*linkedaccounts.LinkedAccount{}, nil
	}
	return r.getMany(ctx, queryGetAllByUserIDs, pq.Array(userIDs))
}
