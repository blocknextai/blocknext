package workflows

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"

	"github.com/blocknextai/go-packages/dag"
)

func ContentHash(nodes []dag.Node, edges []dag.Edge) string {
	payload := map[string]any{
		"nodes": nodes,
		"edges": edges,
	}

	encoded, err := json.Marshal(payload)
	if err != nil {
		return ""
	}

	sum := sha256.Sum256(encoded)
	return hex.EncodeToString(sum[:])
}
