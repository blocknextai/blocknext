package dag

// Node represents a single workflow node in the DAG, including its type,
// instructions, parameters, settings, credentials, and canvas placement.
// Title is the name the author gave the node on the canvas; empty means the
// catalog name. HandleLayout is canvas-only: it names the sides the node's
// input and output handles sit on, in "<in>-<out>" form ("l-r", "t-b", …).
type Node struct {
	ID                 string         `json:"id"`
	Type               string         `json:"type"`
	NodeID             string         `json:"nodeId"`
	Title              string         `json:"title,omitempty"`
	Instruction        string         `json:"instruction,omitempty"`
	RuntimeInstruction string         `json:"runtimeInstruction,omitempty"`
	RuntimePrompt      string         `json:"runtimePrompt,omitempty"`
	Parameters         map[string]any `json:"parameters,omitempty"`
	Settings           *NodeSettings  `json:"settings,omitempty"`
	Credentials        map[string]any `json:"credentials,omitempty"`
	Position           *NodePosition  `json:"position,omitempty"`
	HandleLayout       string         `json:"handleLayout,omitempty"`
}
