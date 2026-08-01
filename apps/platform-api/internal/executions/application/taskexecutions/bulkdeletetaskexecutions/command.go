package bulkdeletetaskexecutions

import (
	"github.com/google/uuid"
)

type BulkDeleteTaskExecutionsCommand struct {
	IDs            []uuid.UUID
	OrganizationID uuid.UUID
}
