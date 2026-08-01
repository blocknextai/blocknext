package adapters

// TriggerContext holds the extracted data from a webhook trigger payload.
// This is a flow-level context injected into Task and available to all nodes
// during execution via task.TriggerContext in node_executor.go.
//
// IMPORTANT: If you add/remove/rename fields here, update TriggerVariables below.
type TriggerContext struct {
	Source  string         `json:"source"`
	Sender  string         `json:"sender"`
	Prompt  string         `json:"prompt"`
	Payload map[string]any `json:"payload"`
}

// TriggerVariables are the available variable patterns for use in node instructions.
var TriggerVariables = []string{
	"$trigger.source",
	"$trigger.sender",
	"$trigger.prompt",
	"$trigger.payload",
}
