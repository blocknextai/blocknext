package getalltaskexecutions

import (
	"time"

	"github.com/google/uuid"
)

type TaskExecutionResponse struct {
	ID            uuid.UUID  `json:"id,omitempty"`
	Workflow      Workflow   `json:"workflow"`
	Status        string     `json:"status,omitempty"`
	ExecutionType string     `json:"executionType,omitempty"`
	ErrorMessage  *string    `json:"errorMessage,omitempty"`
	StartedAt     *time.Time `json:"startedAt,omitempty"`
	CompletedAt   *time.Time `json:"completedAt,omitempty"`
}

type Workflow struct {
	ID    uuid.UUID `json:"id,omitempty"`
	Title string    `json:"title,omitempty"`
}

type GetAllTaskExecutionsResponse struct {
	Items      []*TaskExecutionResponse `json:"items"`
	TotalCount int64                    `json:"totalCount"`
}
