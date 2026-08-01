package triggers

import (
	"github.com/blocknextai/go-packages/dag"
)

type RuntimeConfig struct {
	RuntimePrompt string     `json:"runtimePrompt,omitempty"`
	Nodes         []dag.Node `json:"nodes,omitempty"`
}
