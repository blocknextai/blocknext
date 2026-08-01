package gettaskexecutionbyid

import (
	"github.com/google/uuid"
)

type GetTaskExecutionByIDQuery struct {
	ID             uuid.UUID
	OrganizationID uuid.UUID
}
