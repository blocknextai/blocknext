package bulkdeletetaskexecutions

import (
	"github.com/google/uuid"
)

type BulkDeleteTaskExecutionsResponse struct {
	DeletedCount int         `json:"deletedCount"`
	DeletedIDs   []uuid.UUID `json:"deletedIds"`
	FailedIDs    []uuid.UUID `json:"failedIds"`
}
