package database

import (
	"context"
	"database/sql"
)

// TransactionFromContext returns the transaction stored in ctx, if any, and
// reports whether one was present.
func TransactionFromContext(ctx context.Context) (*sql.Tx, bool) {
	tx, ok := ctx.Value(txContextKey).(*sql.Tx)
	return tx, ok
}
