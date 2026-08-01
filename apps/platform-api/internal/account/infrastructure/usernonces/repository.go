package usernonces

import (
	"context"
	"database/sql"

	"github.com/blocknextai/go-packages/database"
	accountDomain "github.com/blocknextai/platform-api/internal/account/domain"
	"github.com/blocknextai/platform-api/internal/account/domain/usernonces"
)

const (
	tableName = "account.nonces"
	columns   = "id, auth_provider, provider_id, nonce, code_verifier, code_challenge, code_challenge_method, expires_at, created_at, updated_at, deleted_at"
)

var (
	queryGetByNonceAndAuthProvider = database.BuildQuery(`
		SELECT `, columns, `
		FROM `, tableName, `
		WHERE
			nonce = $1
			AND auth_provider = $2
			AND expires_at > NOW()
			AND deleted_at IS NULL
		ORDER BY
			created_at DESC
		LIMIT 1
	`)

	queryCreate = database.BuildQuery(`
		INSERT INTO `, tableName, ` (`, columns, `)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
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

type UserNonceRepository struct {
	database.BaseRepository
}

func NewUserNonceRepository(db *sql.DB) usernonces.UserNonceRepository {
	return &UserNonceRepository{
		BaseRepository: database.NewBaseRepository(db),
	}
}

func (r *UserNonceRepository) scan(row interface{ Scan(dest ...any) error }, _ ...any) (*usernonces.UserNonce, error) {
	var userNonce usernonces.UserNonce
	var providerID sql.NullString
	var deletedAt sql.NullTime

	err := row.Scan(
		&userNonce.ID,
		&userNonce.AuthProvider,
		&providerID,
		&userNonce.Nonce,
		&userNonce.CodeVerifier,
		&userNonce.CodeChallenge,
		&userNonce.CodeChallengeMethod,
		&userNonce.ExpiresAt,
		&userNonce.CreatedAt,
		&userNonce.UpdatedAt,
		&deletedAt,
	)

	if err != nil {
		return nil, err
	}

	userNonce.ProviderID = database.ScanNullString(providerID)
	userNonce.DeletedAt = database.ScanNullTime(deletedAt)

	return &userNonce, nil
}

func (r *UserNonceRepository) getOne(ctx context.Context, query string, args ...any) (*usernonces.UserNonce, error) {
	return database.GetOne(ctx, r.Executor(ctx), query, r.scan, usernonces.ErrUserNonceNotFound, args...)
}

func (r *UserNonceRepository) exec(ctx context.Context, query string, args ...any) error {
	return database.Exec(ctx, r.Executor(ctx), query, args...)
}

func (r *UserNonceRepository) execWithRowCheck(ctx context.Context, query string, args ...any) error {
	return database.ExecWithRowCheck(ctx, r.Executor(ctx), query, usernonces.ErrUserNonceNotFound, args...)
}

func (r *UserNonceRepository) GetByNonceAndAuthProvider(ctx context.Context, nonce string, authProvider accountDomain.AuthProvider) (*usernonces.UserNonce, error) {
	return r.getOne(ctx, queryGetByNonceAndAuthProvider, nonce, authProvider)
}

func (r *UserNonceRepository) Create(ctx context.Context, userNonce *usernonces.UserNonce) error {
	err := r.exec(ctx, queryCreate,
		userNonce.ID,
		userNonce.AuthProvider,
		userNonce.ProviderID,
		userNonce.Nonce,
		userNonce.CodeVerifier,
		userNonce.CodeChallenge,
		userNonce.CodeChallengeMethod,
		userNonce.ExpiresAt,
		userNonce.CreatedAt,
		userNonce.UpdatedAt,
		userNonce.DeletedAt,
	)
	if err != nil && database.IsUniqueViolationOn(err, "uq_account_nonces_nonce") {
		return usernonces.ErrNonceAlreadyExists
	}
	return err
}

func (r *UserNonceRepository) Delete(ctx context.Context, userNonce *usernonces.UserNonce) error {
	return r.execWithRowCheck(ctx, queryDelete,
		userNonce.ID,
		userNonce.UpdatedAt,
		userNonce.DeletedAt,
	)
}
