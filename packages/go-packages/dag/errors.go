package dag

import (
	"errors"
)

var (
	// ErrCycleDetected indicates that the graph contains a cycle and cannot be sorted.
	ErrCycleDetected = errors.New("cycle detected")
	// ErrEmptyDAG indicates that the graph has no nodes.
	ErrEmptyDAG = errors.New("empty dag")
	// ErrInvalidNode indicates that a node or edge has missing or invalid identifiers.
	ErrInvalidNode = errors.New("invalid node")
	// ErrOrphanedNode indicates that no nodes are reachable from any start node.
	ErrOrphanedNode = errors.New("orphaned node")
)
