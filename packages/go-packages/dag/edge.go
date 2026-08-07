package dag

// Edge represents a directed connection from a source node to a target node.
// SourceHandle names which of the source node's outputs the edge leaves from;
// it is empty for a node with a single, unnamed output.
type Edge struct {
	ID           string `json:"id,omitempty"`
	Source       string `json:"source"`
	SourceHandle string `json:"sourceHandle,omitempty"`
	Target       string `json:"target"`
}
