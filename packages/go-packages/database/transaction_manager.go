package database

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"time"
)

// TransactionManager runs a function within a single database transaction.
type TransactionManager interface {
	ExecuteInTransaction(ctx context.Context, fn func(ctx context.Context) error) error
}

type transactionManager struct {
	db           *sql.DB
	queryTimeout time.Duration
}

// NewTransactionManager returns a TransactionManager backed by db that applies
// the given per-transaction query timeout when it is positive.
func NewTransactionManager(db *sql.DB, queryTimeout time.Duration) TransactionManager {
	return &transactionManager{db: db, queryTimeout: queryTimeout}
}

// ExecuteInTransaction runs fn inside a transaction. When ctx already carries a
// transaction the call joins it instead of opening a second one: fn runs against
// the caller's transaction and the outermost call remains responsible for commit
// and rollback. This keeps a nested call's work atomic with its caller's, and
// avoids taking transaction-scoped locks (advisory locks, SELECT FOR UPDATE) on
// a second connection that would deadlock against the transaction already
// holding them.
//
// An error returned from a joined fn therefore rolls back the whole outer
// transaction, not just the nested portion.
func (tm *transactionManager) ExecuteInTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	if _, ok := TransactionFromContext(ctx); ok {
		return fn(ctx)
	}

	if tm.queryTimeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, tm.queryTimeout)
		defer cancel()
	}

	tx, err := tm.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	txCtx := context.WithValue(ctx, txContextKey, tx)

	err = fn(txCtx)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil && !errors.Is(rbErr, sql.ErrTxDone) {
			slog.ErrorContext(ctx, "Transaction rollback failed",
				"component", "TransactionManager",
				"rollbackError", rbErr,
				"originalError", err,
			)
		}
		return err
	}

	return tx.Commit()
}
