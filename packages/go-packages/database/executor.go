package database

import (
	"context"
	"database/sql"
)

// Executor abstracts the query-execution methods shared by *sql.DB and *sql.Tx,
// allowing repository helpers to run against either a connection or a transaction.
type Executor interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
	PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)
}
