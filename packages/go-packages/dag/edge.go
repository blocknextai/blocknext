package dag

// Edge represents a directed connection from a source node to a target node,
// optionally guarded by a condition.
type Edge struct {
	ID        string `json:"id,omitempty"`
	Source    string `json:"source"`
	Target    string `json:"target"`
	Condition string `json:"condition,omitempty"`
}
