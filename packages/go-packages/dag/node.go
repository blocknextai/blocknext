package dag

// Node represents a single workflow node in the DAG, including its type,
// instructions, parameters, settings, credentials, and position.
type Node struct {
	ID                 string         `json:"id"`
	Type               string         `json:"type"`
	NodeID             string         `json:"nodeId"`
	Instruction        string         `json:"instruction,omitempty"`
	RuntimeInstruction string         `json:"runtimeInstruction,omitempty"`
	RuntimePrompt      string         `json:"runtimePrompt,omitempty"`
	Parameters         map[string]any `json:"parameters,omitempty"`
	Settings           *NodeSettings  `json:"settings,omitempty"`
	Credentials        map[string]any `json:"credentials,omitempty"`
	Position           *NodePosition  `json:"position,omitempty"`
}
