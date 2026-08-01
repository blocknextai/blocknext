package verificationtokens

import (
	"context"
	"database/sql"
	"time"

	"github.com/blocknextai/go-packages/database"
	"github.com/blocknextai/platform-api/internal/account/domain/verificationtokens"
	"github.com/google/uuid"
)

const (
	tableName = "account.verification_tokens"
	columns   = "id, user_id, purpose, token_hash, email, expires_at, created_at, updated_at, deleted_at"
)

var (
	queryGetByTokenHash = database.BuildQuery(`
		SELECT `, columns, `
		FROM `, tableName, `
		WHERE
			token_hash = $1
			AND deleted_at IS NULL
		LIMIT 1
	`)

	queryCreate = database.BuildQuery(`
		INSERT INTO `, tableName, ` (`, columns, `)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`)

	queryUpdate = database.BuildQuery(`
		UPDATE `, tableName, `
		SET
			deleted_at = $2,
			updated_at = $3
		WHERE
			id = $1
			AND deleted_at IS NULL
	`)

	queryInvalidateAllForUserAndPurpose = database.BuildQuery(`
		UPDATE `, tableName, `
		SET
			deleted_at = $3,
			updated_at = $3
		WHERE
			user_id = $1
			AND purpose = $2
			AND deleted_at IS NULL
	`)
)

type VerificationTokenRepository struct {
	database.BaseRepository
}

func NewVerificationTokenRepository(db *sql.DB) verificationtokens.VerificationTokenRepository {
	return &VerificationTokenRepository{
		BaseRepository: database.NewBaseRepository(db),
	}
}

func (r *VerificationTokenRepository) scan(row interface{ Scan(dest ...any) error }, _ ...any) (*verificationtokens.VerificationToken, error) {
	var token verificationtokens.VerificationToken
	var deletedAt sql.NullTime

	err := row.Scan(
		&token.ID,
		&token.UserID,
		&token.Purpose,
		&token.TokenHash,
		&token.Email,
		&token.ExpiresAt,
		&token.CreatedAt,
		&token.UpdatedAt,
		&deletedAt,
	)
	if err != nil {
		return nil, err
	}

	token.DeletedAt = database.ScanNullTime(deletedAt)

	return &token, nil
}

func (r *VerificationTokenRepository) getOne(ctx context.Context, query string, args ...any) (*verificationtokens.VerificationToken, error) {
	return database.GetOne(ctx, r.Executor(ctx), query, r.scan, verificationtokens.ErrVerificationTokenNotFound, args...)
}

func (r *VerificationTokenRepository) exec(ctx context.Context, query string, args ...any) error {
	return database.Exec(ctx, r.Executor(ctx), query, args...)
}

func (r *VerificationTokenRepository) execWithRowCheck(ctx context.Context, query string, args ...any) error {
	return database.ExecWithRowCheck(ctx, r.Executor(ctx), query, verificationtokens.ErrVerificationTokenNotFound, args...)
}

func (r *VerificationTokenRepository) GetByTokenHash(ctx context.Context, tokenHash string) (*verificationtokens.VerificationToken, error) {
	return r.getOne(ctx, queryGetByTokenHash, tokenHash)
}

func (r *VerificationTokenRepository) Create(ctx context.Context, token *verificationtokens.VerificationToken) error {
	err := r.exec(ctx, queryCreate,
		token.ID,
		token.UserID,
		token.Purpose,
		token.TokenHash,
		token.Email,
		token.ExpiresAt,
		token.CreatedAt,
		token.UpdatedAt,
		token.DeletedAt,
	)
	if err != nil && database.IsUniqueViolationOn(err, "uq_account_verification_tokens_hash") {
		return verificationtokens.ErrVerificationTokenAlreadyExists
	}
	return err
}

func (r *VerificationTokenRepository) Update(ctx context.Context, token *verificationtokens.VerificationToken) error {
	return r.execWithRowCheck(ctx, queryUpdate,
		token.ID,
		token.DeletedAt,
		token.UpdatedAt,
	)
}

func (r *VerificationTokenRepository) InvalidateAllForUserAndPurpose(ctx context.Context, userID uuid.UUID, purpose verificationtokens.Purpose) error {
	return r.exec(ctx, queryInvalidateAllForUserAndPurpose, userID, purpose, time.Now().UTC())
}
