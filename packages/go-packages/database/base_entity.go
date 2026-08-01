package database

import (
	"time"

	"github.com/google/uuid"
)

// BaseEntity holds the common fields shared by persisted entities, including
// identity and audit timestamps.
type BaseEntity struct {
	ID        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}
