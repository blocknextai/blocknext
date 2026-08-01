package adapters

import (
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/adapters"
)

type AdapterService interface {
	GetAllAdapters() []adapters.TriggerAdapter
	GetAdapterByType(adapterType string) (adapters.TriggerAdapter, bool)
	GetTriggerVariables() []string
}

type adapterService struct{}

func NewAdapterService() AdapterService {
	return &adapterService{}
}

func (s *adapterService) GetAllAdapters() []adapters.TriggerAdapter {
	return adapters.GetAllAdapters()
}

func (s *adapterService) GetAdapterByType(adapterType string) (adapters.TriggerAdapter, bool) {
	return adapters.GetAdapter(adapterType)
}

func (s *adapterService) GetTriggerVariables() []string {
	return adapters.TriggerVariables
}
