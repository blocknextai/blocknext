package database

import (
	"context"
	"database/sql"
	"errors"
)

// ScanFunc scans a single row into a value of type T, optionally using extras
// to receive additional scan destinations such as a total-count column.
type ScanFunc[T any] func(row interface{ Scan(dest ...any) error }, extras ...any) (T, error)

// GetOne runs the query, scans the single resulting row with scan, and returns
// notFoundErr if no row is found.
func GetOne[T any](ctx context.Context, exec Executor, query string, scan ScanFunc[T], notFoundErr error, args ...any) (T, error) {
	var zero T
	row := exec.QueryRowContext(ctx, query, args...)
	result, err := scan(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return zero, notFoundErr
		}
		return zero, err
	}
	return result, nil
}

// GetMany runs the query and scans every resulting row with scan, returning the
// collected slice.
func GetMany[T any](ctx context.Context, exec Executor, query string, scan ScanFunc[T], args ...any) (result []T, err error) {
	rows, err := exec.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer func() {
		err = errors.Join(err, rows.Close())
	}()

	result = make([]T, 0)
	for rows.Next() {
		item, scanErr := scan(rows)
		if scanErr != nil {
			return nil, scanErr
		}
		result = append(result, item)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return result, nil
}

// GetManyWithTotal runs the query and scans every resulting row with scan,
// passing a pointer to total so the scan can also read a total-count column,
// and returns the collected slice along with the total.
func GetManyWithTotal[T any](ctx context.Context, exec Executor, query string, scan ScanFunc[T], args ...any) (result []T, total int64, err error) {
	rows, err := exec.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer func() {
		err = errors.Join(err, rows.Close())
	}()

	result = make([]T, 0)
	for rows.Next() {
		item, scanErr := scan(rows, &total)
		if scanErr != nil {
			return nil, 0, scanErr
		}
		result = append(result, item)
	}
	if err = rows.Err(); err != nil {
		return nil, 0, err
	}
	return result, total, nil
}
