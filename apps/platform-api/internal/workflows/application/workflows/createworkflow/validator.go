package createworkflow

import (
	"strings"

	workflowsDomainWorkflows "github.com/blocknextai/platform-api/internal/workflows/domain/workflows"
)

const (
	MaxTitleLength = 255
)

func (c *CreateWorkflowCommand) Validate() error {
	if strings.TrimSpace(c.Title) == "" {
		return workflowsDomainWorkflows.ErrTitleIsRequired
	}

	if len(strings.TrimSpace(c.Title)) > MaxTitleLength {
		return workflowsDomainWorkflows.ErrTitleTooLong
	}

	return nil
}
