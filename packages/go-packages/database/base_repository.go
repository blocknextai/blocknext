package database

import (
	"context"
	"database/sql"
)

// BaseRepository provides shared database access for repositories, resolving
// the correct Executor depending on whether a transaction is active.
type BaseRepository struct {
	db *sql.DB
}

// NewBaseRepository returns a BaseRepository backed by the given database handle.
func NewBaseRepository(db *sql.DB) BaseRepository {
	return BaseRepository{db: db}
}

// Executor returns the Executor to use for the given context, preferring the
// transaction stored in the context if one is present, otherwise the underlying
// database handle.
func (r *BaseRepository) Executor(ctx context.Context) Executor {
	if tx, ok := TransactionFromContext(ctx); ok {
		return tx
	}
	return r.db
}

// DB returns the underlying database handle.
func (r *BaseRepository) DB() *sql.DB {
	return r.db
}
