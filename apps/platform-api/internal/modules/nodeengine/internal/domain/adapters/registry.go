package adapters

import (
	"sync"
)

var (
	once     sync.Once
	instance *AdapterRegistry
)

type AdapterRegistry struct {
	adaptersMap map[string]TriggerAdapter
	adapters    []TriggerAdapter
}

func GetRegistry() *AdapterRegistry {
	once.Do(func() {
		instance = &AdapterRegistry{
			adaptersMap: make(map[string]TriggerAdapter),
			adapters:    make([]TriggerAdapter, 0),
		}
	})
	return instance
}

func (r *AdapterRegistry) RegisterAdapter(adapter TriggerAdapter) {
	r.adaptersMap[adapter.GetID()] = adapter
	r.adapters = append(r.adapters, adapter)
}

func (r *AdapterRegistry) GetAdapter(adapterType string) (TriggerAdapter, bool) {
	adapter, ok := r.adaptersMap[adapterType]
	return adapter, ok
}

func (r *AdapterRegistry) GetAllAdapters() []TriggerAdapter {
	return r.adapters
}

func RegisterAdapter(adapter TriggerAdapter) {
	if adapter.GetDisabled() {
		return
	}

	GetRegistry().RegisterAdapter(adapter)
}

func GetAdapter(adapterType string) (TriggerAdapter, bool) {
	return GetRegistry().GetAdapter(adapterType)
}

func GetAllAdapters() []TriggerAdapter {
	return GetRegistry().GetAllAdapters()
}

func Adapt(adapterType string, raw map[string]any) (*TriggerContext, error) {
	adapter, ok := GetAdapter(adapterType)
	if !ok {
		return nil, ErrAdapterNotFound
	}
	return adapter.Adapt(raw)
}
