package database

import (
	"slices"
	"strings"

	"github.com/lib/pq"
)

// BuildQuery concatenates the given parts into a single query string.
//
// SECURITY: BuildQuery performs no escaping, quoting, or placeholder
// substitution — it is a raw string join. Every part MUST be a trusted, static
// SQL fragment. Never pass untrusted or user-supplied input as a part: doing so
// results in SQL injection. Values must always be bound through the args of the
// execution helpers (Exec, GetOne, GetMany, GetManyWithTotal), never spliced
// into the query via BuildQuery. Dynamic identifiers (table/column names, sort
// columns) cannot be parameterized and must be passed through AllowedIdentifier
// (or QuoteIdentifier for an already-validated name) before reaching BuildQuery.
func BuildQuery(parts ...string) string {
	var b strings.Builder
	for _, p := range parts {
		b.WriteString(p)
	}
	return b.String()
}

// QuoteIdentifier returns name quoted for safe use as a PostgreSQL identifier
// (table or column name) inside a query fragment. Prefer AllowedIdentifier when
// name originates from user input, since quoting alone does not constrain which
// identifiers a caller may reference.
func QuoteIdentifier(name string) string {
	return pq.QuoteIdentifier(name)
}

// AllowedIdentifier validates name against the allowed set and returns it quoted
// for safe use as a PostgreSQL identifier. It returns ErrDisallowedIdentifier if
// name is not one of the allowed values. Use this for dynamic identifiers
// derived from user input (for example a sort column), which cannot be bound as
// query parameters.
func AllowedIdentifier(name string, allowed ...string) (string, error) {
	if slices.Contains(allowed, name) {
		return pq.QuoteIdentifier(name), nil
	}
	return "", ErrDisallowedIdentifier
}

// SortDirection validates a sort direction and returns its canonical uppercase
// form ("ASC" or "DESC"), ignoring surrounding whitespace and case. It returns
// ErrInvalidSortDirection for any other value, so a user-supplied direction can
// be safely embedded in an ORDER BY clause.
func SortDirection(direction string) (string, error) {
	switch strings.ToUpper(strings.TrimSpace(direction)) {
	case "ASC":
		return "ASC", nil
	case "DESC":
		return "DESC", nil
	default:
		return "", ErrInvalidSortDirection
	}
}
