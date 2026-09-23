package toolinvocations

import (
	apikeysContract "github.com/blocknextai/platform-api/internal/modules/apikeys/contract"
	"github.com/blocknextai/platform-api/internal/modules/executions/internal/domain/toolinvocations"
)

type Service struct {
	toolInvocationRepository toolinvocations.ToolInvocationRepository
	apiKeyService            apikeysContract.APIKeyService
}

func NewService(
	toolInvocationRepository toolinvocations.ToolInvocationRepository,
	apiKeyService apikeysContract.APIKeyService,
) *Service {
	return &Service{
		toolInvocationRepository: toolInvocationRepository,
		apiKeyService:            apiKeyService,
	}
}
