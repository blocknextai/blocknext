package nodeschema

import (
	"log/slog"
	"sync"

	"github.com/blocknextai/go-packages/json"
	nodeEngineApplicationNodes "github.com/blocknextai/platform-api/internal/nodeengine/application/nodes"
)

type NodeSchemaContextBuilder struct {
	nodeService nodeEngineApplicationNodes.NodeService
	once        sync.Once
	cached      string
}

func NewNodeSchemaContextBuilder(
	nodeService nodeEngineApplicationNodes.NodeService,
) *NodeSchemaContextBuilder {
	return &NodeSchemaContextBuilder{
		nodeService: nodeService,
	}
}

func (b *NodeSchemaContextBuilder) Build() string {
	b.once.Do(func() {
		nodes := b.nodeService.GetAllNodes()

		data, err := json.Marshal(nodes)
		if err != nil {
			slog.Error("Failed to marshal node schema context",
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
