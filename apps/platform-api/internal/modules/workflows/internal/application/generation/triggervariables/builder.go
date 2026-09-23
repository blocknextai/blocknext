package triggervariables

import (
	"log/slog"
	"sync"

	"github.com/blocknextai/go-packages/json"
	nodeengineContract "github.com/blocknextai/platform-api/internal/modules/nodeengine/contract"
)

type TriggerVariablesContextBuilder struct {
	adapterService nodeengineContract.AdapterService
	once           sync.Once
	cached         string
}

func NewTriggerVariablesContextBuilder(
	adapterService nodeengineContract.AdapterService,
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
