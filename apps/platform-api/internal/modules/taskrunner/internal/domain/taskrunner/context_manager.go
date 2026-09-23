package taskrunner

import (
	"context"

	"github.com/google/uuid"
)

type ContextManager interface {
	CreateContext(taskID uuid.UUID, parentCtx context.Context) context.Context
	CancelContext(taskID uuid.UUID)
	RemoveContext(taskID uuid.UUID)
	CancelAll()
}
