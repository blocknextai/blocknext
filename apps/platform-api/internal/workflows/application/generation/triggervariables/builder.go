package triggervariables

import (
	"log/slog"
	"sync"

	"github.com/blocknextai/go-packages/json"
	nodeEngineApplicationAdapters "github.com/blocknextai/platform-api/internal/nodeengine/application/adapters"
)

type TriggerVariablesContextBuilder struct {
	adapterService nodeEngineApplicationAdapters.AdapterService
	once           sync.Once
	cached         string
}

func NewTriggerVariablesContextBuilder(
	adapterService nodeEngineApplicationAdapters.AdapterService,
) *TriggerVariablesContextBuilder {
	return &TriggerVariablesContextBuilder{
		adapterService: adapterService,
	}
}

func (b *TriggerVariablesContextBuilder) Build() string {
	b.once.Do(func() {
		variables := b.adapterService.GetTriggerVariables()

		data, err := json.Marshal(variables)
		if err != nil {
			slog.Error("Failed to marshal trigger variables context",
				"component", "Generation",
				"error", err,
			)
			b.cached = "[]"
			return
		}

		b.cached = string(data)
	})
	return b.cached
}
