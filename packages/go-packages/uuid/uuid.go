// Package uuid provides helpers for generating version 4 and version 7 UUIDs.
package uuid

import (
	"log/slog"

	googleuuid "github.com/google/uuid"
)

// NewV4 returns a randomly generated version 4 UUID.
func NewV4() googleuuid.UUID {
	return googleuuid.New()
}

// NewV7 returns a time-ordered version 7 UUID, falling back to a version 4 UUID
// if generation fails.
func NewV7() googleuuid.UUID {
	id, err := googleuuid.NewV7()
	if err != nil {
		slog.Warn("failed to generate UUID v7, falling back to v4",
			"package", "uuid",
			"function", "NewV7",
			"error", err,
		)
		return NewV4()
	}
	return id
}
