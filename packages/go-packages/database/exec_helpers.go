package database

import (
	"context"
)

// Exec runs the query on exec and discards the result, returning any error.
func Exec(ctx context.Context, exec Executor, query string, args ...any) error {
	_, err := exec.ExecContext(ctx, query, args...)
	return err
}

// ExecWithRowCheck runs the query on exec and returns notFoundErr if no rows
// were affected.
func ExecWithRowCheck(ctx context.Context, exec Executor, query string, notFoundErr error, args ...any) error {
	result, err := exec.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return notFoundErr
	}

	return nil
}
