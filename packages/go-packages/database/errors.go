package database

import (
	"errors"

	"github.com/lib/pq"
)

const (
	pgUniqueViolation = "23505"
)

var (
	// ErrDisallowedIdentifier indicates that a dynamic SQL identifier was not
	// present in the caller-supplied allowlist.
	ErrDisallowedIdentifier = errors.New("disallowed sql identifier")
	// ErrInvalidSortDirection indicates that a sort direction was neither ASC
	// nor DESC.
	ErrInvalidSortDirection = errors.New("invalid sort direction")
)

// IsUniqueViolationOn reports whether err is a PostgreSQL unique-violation error
// raised by the named constraint.
func IsUniqueViolationOn(err error, constraint string) bool {
	var pqErr *pq.Error
	return errors.As(err, &pqErr) && string(pqErr.Code) == pgUniqueViolation && pqErr.Constraint == constraint
}
