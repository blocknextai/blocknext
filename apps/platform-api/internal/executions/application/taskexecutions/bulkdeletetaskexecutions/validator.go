package bulkdeletetaskexecutions

import (
	"github.com/blocknextai/go-packages/apperror"
)

var (
	ErrIDsIsRequired = apperror.Validation("ids is required")
	ErrIDsIsEmpty    = apperror.Validation("ids is empty")
	ErrTooManyIDs    = apperror.Validation("too many ids")
)

const MaxBulkDeleteLimit = 50

func (c *BulkDeleteTaskExecutionsCommand) Validate() error {
	if c.IDs == nil {
		return ErrIDsIsRequired
	}

	if len(c.IDs) == 0 {
		return ErrIDsIsEmpty
	}

	if len(c.IDs) > MaxBulkDeleteLimit {
		return ErrTooManyIDs
	}

	return nil
}
