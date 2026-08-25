package adapter

import (
	"github.com/blocknextai/go-packages/json"
	mcpsdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

func errorResult(message string) *mcpsdk.CallToolResult {
	return &mcpsdk.CallToolResult{
		IsError: true,
		Content: []mcpsdk.Content{
			&mcpsdk.TextContent{Text: message},
		},
	}
}

func successResult(outputs []map[string]any) *mcpsdk.CallToolResult {
	content := make([]mcpsdk.Content, 0, len(outputs))
	for _, output := range outputs {
		text, err := json.Marshal(output)
		if err != nil {
			return errorResult("failed to marshal output: " + err.Error())
		}
		content = append(content, &mcpsdk.TextContent{Text: string(text)})
	}

	return &mcpsdk.CallToolResult{
		StructuredContent: map[string]any{
			itemsKey: outputs,
		},
		Content: content,
	}
}
