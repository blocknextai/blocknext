package adapters

import (
	nodeEngineApplicationAdapters "github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/application/adapters"
)

type Service struct {
	adapterService     nodeEngineApplicationAdapters.AdapterService
	webhookURLTemplate string
}

func NewService(
	adapterService nodeEngineApplicationAdapters.AdapterService,
	webhookURLTemplate string,
) *Service {
	return &Service{
		adapterService:     adapterService,
		webhookURLTemplate: webhookURLTemplate,
	}
}
