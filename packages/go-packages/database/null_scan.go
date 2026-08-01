package database

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

// ScanNullString returns a pointer to the string value, or nil if it is not valid.
func ScanNullString(ns sql.NullString) *string {
	if !ns.Valid {
		return nil
	}
	return &ns.String
}

// ScanNullTime returns a pointer to the time value, or nil if it is not valid.
func ScanNullTime(ns sql.NullTime) *time.Time {
	if !ns.Valid {
		return nil
	}
	return &ns.Time
}

// ScanNullInt64 returns a pointer to the int64 value, or nil if it is not valid.
func ScanNullInt64(ns sql.NullInt64) *int64 {
	if !ns.Valid {
		return nil
	}
	return &ns.Int64
}

// ScanNullUUID returns a pointer to the UUID value, or nil if it is not valid.
func ScanNullUUID(nu uuid.NullUUID) *uuid.UUID {
	if !nu.Valid {
		return nil
	}
	return &nu.UUID
}
