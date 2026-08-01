package gettaskexecutionbyid

import (
	"time"

	"github.com/google/uuid"
)

type GetTaskExecutionByIDResponse struct {
	ID              uuid.UUID        `json:"id,omitempty"`
	Workflow        Workflow         `json:"workflow"`
	TriggeredByUser *TriggeredByUser `json:"triggeredByUser,omitempty"`
	Status          string           `json:"status,omitempty"`
	ExecutionType   string           `json:"executionType,omitempty"`
	ErrorMessage    *string          `json:"errorMessage,omitempty"`
	StartedAt       *time.Time       `json:"startedAt,omitempty"`
	CompletedAt     *time.Time       `json:"completedAt,omitempty"`
	NodeExecutions  []NodeExecution  `json:"nodeExecutions,omitempty"`
}

type Workflow struct {
	ID    uuid.UUID `json:"id,omitempty"`
	Title string    `json:"title,omitempty"`
}

type LinkedAccount struct {
	AuthProvider string  `json:"authProvider"`
	DisplayName  *string `json:"displayName"`
	IsPrimary    bool    `json:"isPrimary"`
}

type TriggeredByUser struct {
	ID             uuid.UUID       `json:"id,omitempty"`
	Alias          string          `json:"alias,omitempty"`
	IsVerified     bool            `json:"isVerified"`
	LinkedAccounts []LinkedAccount `json:"linkedAccounts"`
}

type NodeExecution struct {
	ID           uuid.UUID        `json:"id,omitempty"`
	NodeType     string           `json:"nodeType,omitempty"`
	Status       string           `json:"status,omitempty"`
	Inputs       []map[string]any `json:"inputs,omitempty"`
	Outputs      []map[string]any `json:"outputs,omitempty"`
	ErrorMessage *string          `json:"errorMessage,omitempty"`
	StartedAt    *time.Time       `json:"startedAt,omitempty"`
	CompletedAt  *time.Time       `json:"completedAt,omitempty"`
}
