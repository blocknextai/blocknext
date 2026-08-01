// Package database provides helpers for working with PostgreSQL through the
// standard database/sql package, including connection pool setup, transaction
// management, query building, and generic row scanning utilities.
package database

import (
	"database/sql"
	"log/slog"
	"time"

	_ "github.com/lib/pq"
)

const (
	defaultMaxOpenConns    = 25
	defaultMaxIdleConns    = 10
	defaultConnMaxLifetime = 30 * time.Minute
	defaultConnMaxIdleTime = 5 * time.Minute
)

// NewDB opens a PostgreSQL connection using the given DSN, verifies it with a
// ping, and configures the connection pool. Any non-positive pool setting falls
// back to its default value.
func NewDB(
	dsn string,
	maxOpenConns int,
	maxIdleConns int,
	connMaxLifetime time.Duration,
	connMaxIdleTime time.Duration,
) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		slog.Error("Failed to open database connection",
			"component", "database",
			"error", err,
		)
		return nil, err
	}

	if err := db.Ping(); err != nil {
		slog.Error("Failed to ping database",
			"component", "database",
			"error", err,
		)
		return nil, err
	}

	if maxOpenConns <= 0 {
		maxOpenConns = defaultMaxOpenConns
	}
	if maxIdleConns <= 0 {
		maxIdleConns = defaultMaxIdleConns
	}
	if connMaxLifetime <= 0 {
		connMaxLifetime = defaultConnMaxLifetime
	}
	if connMaxIdleTime <= 0 {
		connMaxIdleTime = defaultConnMaxIdleTime
	}

	db.SetMaxOpenConns(maxOpenConns)
	db.SetMaxIdleConns(maxIdleConns)
	db.SetConnMaxLifetime(connMaxLifetime)
	db.SetConnMaxIdleTime(connMaxIdleTime)

	slog.Info("Database connection pool configured",
		"component", "database",
		"max_open_conns", maxOpenConns,
		"max_idle_conns", maxIdleConns,
		"max_lifetime", connMaxLifetime,
		"max_idle_time", connMaxIdleTime,
	)

	slog.Info("Successfully connected to database", "component", "database")

	return db, nil
}
